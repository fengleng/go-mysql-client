package backend

import (
	"context"
	"errors"
	"github.com/fengleng/go-common/core/resource_polll"
	"github.com/fengleng/go-mysql-client/mysql"
	"github.com/fengleng/log"
	"strings"
	"sync"
	"time"
)

const (
	Up = iota
	Down
	ManualDown
	Unknown

	PingPeroid int64 = 4
)

var (
	// ErrConnectionPoolClosed means pool closed error
	ErrConnectionPoolClosed = errors.New("connection pool is closed")
)

type DB struct {
	sync.RWMutex
	addr     string
	user     string
	password string
	db       string
	state    int32

	connLock sync.Mutex
	*resource_polll.ResourcePool

	lastPing int64

	opt *Option
}

func Open(addr string, user string, password string, dbName string, opts ...DbOption) (*DB, error) {
	db := new(DB)
	db.addr = addr
	db.user = user
	db.password = password
	db.db = dbName
	db.opt = DefaultOption
	for _, fn := range opts {
		fn(db.opt)
	}
	db.ResourcePool = resource_polll.NewResourcePool(db.newConn, db.opt.capacity, db.opt.maxCapacity, db.opt.idleTimeout)
	err := db.Ping()
	if err != nil {
		log.Error("%v", err)
		return nil, err
	}
	return db, nil
}

func (db *DB) GetConn() (*Conn, error) {
	getCtx, cancel := context.WithTimeout(context.Background(), db.opt.connTimeout)
	defer cancel()
	r, err := db.ResourcePool.Get(getCtx)
	if err != nil {
		log.Error("%v", err)
		return nil, err
	}
	return r.(*Conn), nil
}

func (db *DB) PutConn(c *Conn) {
	if db.ResourcePool == nil {
		panic(ErrConnectionPoolClosed)
	}

	if c == nil {
		db.ResourcePool.Put(nil)
	} else if err := db.tryReuse(c); err != nil {
		c.Close()
		db.ResourcePool.Put(nil)
	} else {
		db.ResourcePool.Put(c)
	}
}

func (db *DB) Addr() string {
	return db.addr
}

func (db *DB) State() string {
	var state string
	switch db.state {
	case Up:
		state = "up"
	case Down, ManualDown:
		state = "down"
	case Unknown:
		state = "unknow"
	}
	return state
}

func (db *DB) Close() {
	if db.ResourcePool == nil {
		return
	}
	db.connLock.Lock()
	defer db.connLock.Unlock()
	db.ResourcePool.Close()
	db.ResourcePool = nil
	return
}

func (db *DB) Ping() error {
	var err error
	conn, err := db.GetConn()
	if err != nil {
		log.Error("%v", err)
		return err
	}
	err = conn.Ping()
	if err != nil {
		log.Error("%v", err)
		return err
	}
	db.PutConn(conn)
	return nil
}

//type Factory func() (Resource, error)
func (db *DB) newConn() (resource_polll.Resource, error) {
	co := new(Conn)

	if err := co.Connect(db.addr, db.user, db.password, db.db); err != nil {
		return nil, err
	}

	co.pushTimestamp = time.Now().Unix()

	return co, nil
}

func (db *DB) Execute(command string, args ...interface{}) (*mysql.Result, error) {
	var (
		err error
		r   *mysql.Result
	)
	conn, err := db.GetConn()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	defer db.PutConn(conn)
	r, err = conn.Execute(command, args...)
	if err != nil && strings.Contains(err.Error(), "broken pipe") { //连接关闭
		// retry 3 times, close dc's conn、reset dc's stats and reconnect
		conn.Close()
		conn, err = db.tryConn()
		if err != nil {
			log.Error("%v", err)
			return nil, err
		}
		r, err = conn.Execute(command, args...)
	}
	return r, err
}

func (db *DB) tryConn() (*Conn, error) {
	var (
		err  error
		conn *Conn
	)
	for i := 0; i < 3; i++ {
		var re resource_polll.Resource
		re, err = db.newConn()
		if err == nil {
			conn = re.(*Conn)
			break
		}
	}
	return conn, err
}

func (db *DB) tryReuse(co *Conn) error {
	var err error
	//reuse Connection
	if co.IsInTransaction() {
		//we can not reuse a connection in transaction status
		err = co.Rollback()
		if err != nil {
			return err
		}
	}

	if !co.IsAutoCommit() {
		//we can not  reuse a connection not in autocomit
		_, err = co.exec("set autocommit = 1")
		if err != nil {
			return err
		}
	}

	//connection may be set names early
	//we must use default utf8
	if co.GetCharset() != mysql.DEFAULT_CHARSET {
		err = co.SetCharset(mysql.DEFAULT_CHARSET, mysql.DEFAULT_COLLATION_ID)
		if err != nil {
			return err
		}
	}

	return nil
}

package backend

import (
	"log"
	"sync"
	"testing"
)

func TestDB_Ping(t *testing.T) {
	db, err := Open("127.0.0.1:3306", "root", "root", "biz")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for i := 0; i < 10000; i++ {
		err := db.Ping()
		if err != nil {
			t.Fatal(err)
		}
	}
}
func TestT2(t *testing.T) {
	var m map[string]string
	//m["fdfa"] = "fsafsafa"
	s, ok := m["fdfa"]
	t.Log(s, ok)
}

func TestDb(t *testing.T) {
	db, err := Open("127.0.0.1:3306", "root", "root", "biz")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.Ping()

	//conn, err := db.GetConn()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//res, err := conn.Execute("select *\nfrom quan_channel_live_code\n;")
	//if err != nil {
	//	t.Fatal(err)
	//}
	res, err := db.Execute("update quan_channel_live_code set updated_at = 1623406676 where id = 2727;")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestDB_Execute(t *testing.T) {
	db, err := Open("127.0.0.1:3306", "root", "root", "biz")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.Ping()
	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go exec(db, wg)
	}
	wg.Wait()

}

func exec(db *DB, wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		res, err := db.Execute("select *\nfrom quan_channel_live_code\n;")
		if err != nil {
			log.Println(err)
			return
		}
		//t.Log(db.InUse())
		log.Println(res.GetInt(0, 0))
	}
	wg.Done()
}

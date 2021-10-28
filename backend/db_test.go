package backend

import (
	"log"
	"testing"
)

func TestDB_Ping(t *testing.T) {
	db, err := Open("127.0.0.1:3306", "root", "root", "biz")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.Ping()
}
func TestT2(t *testing.T) {
	var m map[string]string
	//m["fdfa"] = "fsafsafa"
	s, ok := m["fdfa"]
	t.Log(s, ok)
}

func TestDb(t *testing.T) {
	db, err := Open("127.0.0.1:3306", "root", "root", "biz_2000772")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.Ping()

	conn, err := db.GetConn()
	if err != nil {
		t.Fatal(err)
	}
	res, err := conn.Execute("select *\nfrom biz_2000772.quan_ext_contact_follow\nwhere corp_id = '2000772'\nlimit 1\n;")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res.GetInt(0, 0))
}

func TestDB_Execute(t *testing.T) {
	db, err := Open("127.0.0.1:3306", "root", "root", "biz_2000772")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.Ping()

	for i := 0; i < 100; i++ {
		go exec(db)
	}
	//for  {
	//	//db.pushConnCount
	//	//db.tryReuse()
	//	t.Log(db.InUse())
	//	//time.Sleep(1*time.Second)
	//}

}

func exec(db *DB) {
	for i := 0; i < 100; i++ {
		res, err := db.Execute("select *\nfrom biz_2000772.quan_ext_contact_follow\nwhere corp_id = '2000772'\nlimit 1\n;")
		if err != nil {
			log.Println(err)
			return
		}
		//t.Log(db.InUse())
		log.Println(res.GetInt(0, 0))
	}
}

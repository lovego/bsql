package bsql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"testing"
	"time"
)

var testDb *DB

type Info struct {
	Account string
	Name    string
	Areas   map[string]bool
}

type User struct {
	Id        int64
	Phone     string
	Status    int8
	CreatedAt time.Time
	Info
}

type Staffs struct {
	Id        int64
	CompanyId int64
	StaffId   int64
	StaffName string
}

func init() {
	db, err := sql.Open("postgres", "postgres://develop:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	testDb = &DB{db, time.Minute}

	if _, err := db.Exec(`create table  if not exists users(
		id bigint, phone varchar(50), account varchar(100), name varchar(50), status smallint,
		created_at date, areas jsonb
	)`); err != nil {
		log.Panic(err)
	}

	if _, err := db.Exec(`create table if not exists staffs(
		id bigint, company_id bigint, staff_id bigint, staff_name varchar(50)
	)`); err != nil {
		log.Panic(err)
	}
	if _, err := db.Exec(`truncate users`); err != nil {
		log.Panic(err)
	}
	if _, err := db.Exec(`insert into users(id,phone,account,name,status,created_at,areas) values
	 (1,'18380461681','jack111','jack1',0,now(),'{"成都":true}'),
	 (2,'18380461682','jack222','jack2',0,now(),'{"成都":true}'),
	 (3,'18380461683','jack333','jack3',0,now(),'{"成都":true}')`); err != nil {
		log.Panic(err)
	}
}

func TestQuery(t *testing.T) {
	var user User
	if err := testDb.Query(&user, `select * from users where id = $1`, 1); err != nil {
		t.Fatal(err)
	}
	if user.Phone != `18380461681` {
		t.Logf("unexpected phone: %v", user.Phone)
	}

	var users []User
	var area = map[string]bool{"成都": true}
	var query = fmt.Sprintf(`select * from users where areas @> %v`, V(area))
	t.Log(query)
	if err := testDb.Query(&users, query); err != nil {
		t.Fatal(err)
	}
	for _, u := range users {
		if u.Phone != `1838046168`+strconv.FormatInt(u.Id, 10) {
			t.Logf("unexpected phone: %v", u.Phone)
		}
	}
}

func TestExec(t *testing.T) {
	if _, err := testDb.Exec(`update users set phone = '18380461689' where id = $1`, 1); err != nil {
		t.Fatal(err)
	}
	var user User
	if err := testDb.Query(&user, `select phone from users where id = $1`, 1); err != nil {
		t.Fatal(err)
	}
	if user.Phone != `18380461689` {
		t.Logf("unexpected phone: %v", user.Phone)
	}
	var ids []int64
	if err := testDb.Query(&ids, `update users set name = '杰克' returning id`); err != nil {
		t.Fatal(err)
	}
	if len(ids) != 3 {
		t.Logf("unexpected ids: %v", ids)
	}
}

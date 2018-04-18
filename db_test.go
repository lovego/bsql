package bsql

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

var testDb *DB

func init() {
	db, err := sql.Open("postgres", "postgres://develop:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	testDb = &DB{db, time.Second}

	if _, err := testDb.Exec(`
	drop table if exists students;
	create table if not exists students (
		id         bigint,
		name       varchar(50),
		friend_ids bigint[],
		cities     varchar[],
		scores     jsonb,
		money      decimal,
		status     smallint,
		created_at timestamptz,
		updated_at timestamptz default '0001-01-01Z'
	)`); err != nil {
		log.Panic(err)
	}
}

type Student struct {
	Id        int64
	Name      string
	FriendIds pq.Int64Array
	Cities    []string
	Scores    map[string]int
	Money     decimal.Decimal
	Status    int8
	timeFields
}

type timeFields struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func testStudents() []Student {
	rows := []Student{{
		Id: 1, Name: "李雷", FriendIds: []int{2}, Cities: []string{"成都", "北京"},
		Scores: `{"语文": 95, "英语": 97}`,
	}, {
		Id: 2, Name: "韩梅梅", FriendIds: []int{1, 3}, Cities: []string{"成都", "深圳"},
		Scores: `{"语文": 97, "英语": 97}`,
	}, {
		Id: 3, Name: "Tom", FriendIds: []int{2}, Cities: []string{"成都", "NewYork"},
		Scores: `{"语文": 80, "英语": 91}`,
	}}
	for _, row := range rows {
		row.Money = decimal.New(1234, -2)
		row.status = 1
		row.CreatedAt = time.Now()
	}
	return rows
}

func TestDB(t *testing.T) {
	var fields = FieldsFromStruct(Student{}, []string{"Id", "UpdatedAt"})
	var columns = Fields2Columns(fields)

	if _, err := testDb.Exec(`insert into users (
  ) values
	returning *
	 `); err != nil {
		log.Panic(err)
	}

	var user Student
	if err := testDb.Query(&user, `select * from users where id = $1`, 1); err != nil {
		t.Fatal(err)
	}
	if user.Phone != `18380461681` {
		t.Logf("unexpected phone: %v", user.Phone)
	}

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
	var user Student
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

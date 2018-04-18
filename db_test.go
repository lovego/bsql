package bsql

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
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
		Id: 1, Name: "李雷", FriendIds: []int64{2}, Cities: []string{"成都", "北京"},
		Scores: map[string]int{"语文": 95, "英语": 97},
	}, {
		Id: 2, Name: "韩梅梅", FriendIds: []int64{1, 3}, Cities: []string{"成都", "深圳"},
		Scores: map[string]int{"语文": 97, "英语": 97},
	}, {
		Id: 3, Name: "Tom", FriendIds: []int64{2}, Cities: []string{"成都", "NewYork"},
		Scores: map[string]int{"语文": 80, "英语": 91},
	}}
	for i := range rows {
		rows[i].Money = decimal.New(1234, -2)
		rows[i].Status = 1
		rows[i].CreatedAt = time.Now().Round(time.Millisecond)
	}
	return rows
}

func TestDB(t *testing.T) {
	var fields = FieldsFromStruct(Student{}, []string{"Id", "UpdatedAt"})
	var columns = strings.Join(Fields2Columns(fields), ",")

	var expect = testStudents()
	var got []Student
	sql := fmt.Sprintf(
		`insert into students (%s) values %s returning *`, columns, StructValues(expect, fields),
	)
	if err := testDb.Query(got, sql); err != nil {
		t.Log(sql)
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("exptced: %v", got)
	} else {
		t.Log(got)
	}

	var student Student
	if err := testDb.Query(&student, `select * from students where id=$1`, 1); err != nil {
		t.Fatal(err)
	}

	var area = map[string]bool{"成都": true}
	var query = fmt.Sprintf(`select * from students where areas @> %v`, V(area))
	t.Log(query)
	var students []Student
	if err := testDb.Query(&students, query); err != nil {
		t.Fatal(err)
	}

	if _, err := testDb.Exec(`update students set phone = '18380461689' where id = $1`, 1); err != nil {
		t.Fatal(err)
	}
	var user Student
	if err := testDb.Query(&user, `select phone from students where id = $1`, 1); err != nil {
		t.Fatal(err)
	}

	var ids []int64
	if err := testDb.Query(&ids, `update students set name = '杰克' returning id`); err != nil {
		t.Fatal(err)
	}
	if len(ids) != 3 {
		t.Logf("unexpected ids: %v", ids)
	}
}

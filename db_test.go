package bsql

import (
	"database/sql"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

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
		rows[i].UpdatedAt = rows[i].CreatedAt
	}
	return rows
}

func createTable(t *testing.T, db DbOrTx) {
	if _, err := db.Exec(`
	drop table if exists students;
	create table if not exists students (
		id         bigserial,
		name       varchar(50),
		friend_ids bigint[],
		cities     jsonb,
		scores     jsonb,
		money      decimal,
		status     smallint,
		created_at timestamptz,
		updated_at timestamptz default '0001-01-01Z'
	)`); err != nil {
		t.Fatal(err)
	}
}

func TestDB(t *testing.T) {
	db := getTestDB()
	defer db.Close()
	runTests(t, db)
}

func TestScanArray(t *testing.T) {
	db := getTestDB()
	defer db.Close()
	var ints struct {
		Slice pq.Int64Array
	}
	if err := db.Query(&ints, `select '{1,2,3}'::int[] as slice`); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(ints.Slice, pq.Int64Array([]int64{1, 2, 3})) {
		t.Errorf("unexpected: %v", ints.Slice)
	}

	var strs struct {
		Slice []string
	}
	if err := db.Query(&strs, `select '{"abc","de''f"}'::text[] as slice`); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(strs.Slice, []string{"abc", "de'f"}) {
		t.Errorf("unexpected: %v", strs.Slice)
	}
}

func TestScanNil(t *testing.T) {
	db := getTestDB()
	defer db.Close()
	var data int
	if err := db.Query(&data, `select null`); err != nil {
		t.Error(err)
	}
	if data != 0 {
		t.Errorf("unexpected: %v", data)
	}
}

func getTestDB() *DB {
	db, err := sql.Open("postgres", "postgres://develop:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	return &DB{db, time.Second}
}

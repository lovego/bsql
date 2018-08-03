package bsql

import (
	"context"
	"database/sql"
	"fmt"
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
	FriendIds pq.Int64Array `sql:"int[]"`
	Cities    []string
	Scores    map[string]int
	Money     decimal.Decimal
	Status    int8 `sql:"default 0"`
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
		cities     json,
		scores     json,
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
	defer db.db.Close()
	runTests(t, db)
}

func TestScanArray(t *testing.T) {
	db := getTestDB()
	defer db.db.Close()
	var ints pq.Int64Array
	if err := db.Query(&ints, `select '{1,2,3}'::int[] as slice`); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(ints, pq.Int64Array([]int64{1, 2, 3})) {
		t.Errorf("unexpected: %v", ints)
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
	defer db.db.Close()

	var i = 3
	if err := db.Query(&i, `select null`); err != nil {
		t.Error(err)
	}
	if i != 0 {
		t.Errorf("unexpected: %v", i)
	}

	var ints pq.Int64Array
	if err := db.Query(&ints, `select null`); err != nil {
		t.Error(err)
	}
	if len(ints) != 0 {
		t.Errorf("unexpected: %v", ints)
	}

	ints = pq.Int64Array{1}
	if err := db.Query(&ints, `select null`); err != nil {
		t.Error(err)
	}
	if len(ints) != 0 {
		t.Errorf("unexpected: %v", ints)
	}
}

func TestScanValueOutOfRange(t *testing.T) {
	db := getTestDB()
	defer db.db.Close()

	var i int8
	if err := db.Query(&i, `select 128`); err == nil {
		t.Error("expect error")
	} else {
		t.Log(err)
	}
}

func TestScanFloat(t *testing.T) {
	db := getTestDB()
	defer db.db.Close()

	var f float32
	if err := db.Query(&f, `select 1.23`); err != nil {
		t.Fatal(err)
	}
	if f != 1.23 {
		t.Errorf("unexpected: %v", f)
	}
}

func getTestDB() *DB {
	db, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	return &DB{db, time.Second}
}

func ExampleDB_Query() {
	var people struct {
		Name string
		Age  int
	}
	Db, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	var db = &DB{
		db:      Db,
		timeout: time.Second,
	}
	if err := db.Query(&people, `select 'jack' as name, 24 as age`); err != nil {
		log.Panic(err)
	}
	fmt.Println(people)
	// Output:
	// {jack 24}
}

func ExampleDB_QueryT() {
	var people struct {
		Name string
		Age  int
	}
	Db, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	var db = &DB{
		db:      Db,
		timeout: time.Second,
	}
	if err := db.QueryT(2*time.Second, &people, `select 'jack' as name, 24 as age`); err != nil {
		log.Panic(err)
	}
	fmt.Println(people)
	// Output:
	// {jack 24}
}

func ExampleDB_QueryCtx() {
	var people struct {
		Name string
		Age  int
	}
	Db, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	var db = &DB{
		db:      Db,
		timeout: time.Second,
	}
	if err := db.QueryCtx(
		context.Background(), `query people`, &people, `select 'jack' as name, 24 as age`,
	); err != nil {
		log.Panic(err)
	}
	fmt.Println(people)
	// Output:
	// {jack 24}
}

func ExampleDB_Exec() {
	Db, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	var db = &DB{
		db:      Db,
		timeout: time.Second,
	}
	if _, err := db.Exec(`delete from people where id = $1`, 1); err != nil {
		log.Panic(err)
	}
}

func ExampleDB_ExecT() {
	Db, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	var db = &DB{
		db:      Db,
		timeout: time.Second,
	}
	if _, err := db.ExecT(5*time.Second, `delete from people where id = $1`, 1); err != nil {
		log.Panic(err)
	}
}

func ExampleDB_ExecCtx() {
	Db, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	var db = &DB{
		db:      Db,
		timeout: time.Second,
	}
	if _, err := db.ExecCtx(
		context.Background(), `delete people`, `delete from people where id = $1`, 1,
	); err != nil {
		log.Panic(err)
	}
}

func ExampleDB_RunInTransaction() {
	Db, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	var db = &DB{
		db:      Db,
		timeout: time.Second,
	}
	if err := db.RunInTransaction(func(*Tx) error {
		return nil
	}); err != nil {
		log.Panic(err)
	}
}

func ExampleDB_RunInTransactionT() {
	Db, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	var db = &DB{
		db:      Db,
		timeout: time.Second,
	}
	if err := db.RunInTransactionT(5*time.Second, func(*Tx) error {
		return nil
	}); err != nil {
		log.Panic(err)
	}
}

func ExampleDB_RunInTransactionCtx() {
	Db, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	var db = &DB{
		db:      Db,
		timeout: time.Second,
	}
	if err := db.RunInTransactionCtx(
		context.Background(), "test RunInTransactionCtx", func(*Tx, context.Context) error {
			return nil
		},
	); err != nil {
		log.Panic(err)
	}
}

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

func ExampleNew() {
	var rawDb *sql.DB
	db := New(rawDb, time.Second)
	fmt.Println(db.timeout)
	// Output:
	// 1s
}

func ExampleDB_Query() {
	var people struct {
		Name string
		Age  int
	}
	rawDb, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	defer rawDb.Close()
	if err != nil {
		log.Panic(err)
	}
	db := New(rawDb, time.Second)
	if err := db.Query(&people, `select 'jack' as name, 24 as age`); err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v", people)
	// Output:
	// {Name:jack Age:24}
}

func ExampleDB_QueryT() {
	var people struct {
		Name string
		Age  int
	}
	rawDb, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	defer rawDb.Close()
	if err != nil {
		log.Panic(err)
	}
	db := New(rawDb, time.Second)
	if err := db.QueryT(2*time.Second, &people, `select 'jack' as name, 24 as age`); err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v", people)
	// Output:
	// {Name:jack Age:24}
}

func ExampleDB_QueryCtx() {
	var people struct {
		Name string
		Age  int
	}
	rawDb, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	defer rawDb.Close()
	if err != nil {
		log.Panic(err)
	}
	db := New(rawDb, time.Second)
	if err := db.QueryCtx(
		context.Background(), `query people`, &people, `select 'jack' as name, 24 as age`,
	); err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v", people)
	// Output:
	// {Name:jack Age:24}
}

func ExampleDB_Exec() {
	rawDb, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	defer rawDb.Close()
	if err != nil {
		log.Panic(err)
	}
	db := New(rawDb, time.Second)
	result, err := db.Exec(`
		drop table if exists students;
		create table if not exists students (
			id         bigserial,
			name       varchar(50),
			friend_ids bigint[],
			scores     json,
			status     smallint,
			created_at timestamptz,
			updated_at timestamptz default '0001-01-01Z'
		);
		insert into students (
			name, friend_ids, scores, status, created_at
		) values (
			'jack', array[55,66], '{"数学":100, "语文":90}', 0, now()
		);
	`)
	if err != nil {
		log.Panic(err)
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(affectedRows)
	// Output:
	// 1
}

func ExampleDB_ExecT() {
	rawDb, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	defer rawDb.Close()
	if err != nil {
		log.Panic(err)
	}
	db := New(rawDb, time.Second)
	result, err := db.ExecT(time.Second, `
		drop table if exists students;
		create table if not exists students (
			id         bigserial,
			name       varchar(50),
			friend_ids bigint[],
			scores     json,
			status     smallint,
			created_at timestamptz,
			updated_at timestamptz default '0001-01-01Z'
		);
		insert into students (
			name, friend_ids, scores, status, created_at
		) values (
			'jack', array[55,66], '{"数学":100, "语文":90}', 0, now()
		),(
			'rose', array[99,88], '{"数学":90, "语文":100}', 0, now()
		);
	`)
	if err != nil {
		log.Panic(err)
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(affectedRows)
	// Output:
	// 2
}

func ExampleDB_ExecCtx() {
	rawDb, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	defer rawDb.Close()
	if err != nil {
		log.Panic(err)
	}
	db := New(rawDb, time.Second)
	result, err := db.ExecCtx(
		context.Background(), `delete people`, `
		drop table if exists students;
		create table if not exists students (
			id         bigserial,
			name       varchar(50),
			friend_ids bigint[],
			scores     json,
			status     smallint,
			created_at timestamptz,
			updated_at timestamptz default '0001-01-01Z'
		);
		insert into students (
			name, friend_ids, scores, status, created_at
		) values (
			'jack', array[55,66], '{"数学":100, "语文":90}', 0, now()
		),(
			'rose', array[99,88], '{"数学":90, "语文":100}', 0, now()
		);
	`)
	if err != nil {
		log.Panic(err)
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(affectedRows)
	// Output:
	// 2
}

func ExampleDB_RunInTransaction() {
	rawDb, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	defer rawDb.Close()
	if err != nil {
		log.Panic(err)
	}
	db := New(rawDb, time.Second)
	var id int
	if err := db.RunInTransaction(func(tx *Tx) error {
		if err := tx.Query(&id, `select 10 as id`); err != nil {
			return err
		}
		return nil
		return nil
	}); err != nil {
		log.Panic(err)
	}
	fmt.Println(id)
	// Output:
	// 10
}

func ExampleDB_RunInTransactionT() {
	rawDb, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	defer rawDb.Close()
	if err != nil {
		log.Panic(err)
	}
	db := New(rawDb, time.Second)
	var id int
	if err := db.RunInTransactionT(5*time.Second, func(tx *Tx) error {
		if err := tx.Query(&id, `select 10 as id`); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}
	fmt.Println(id)
	// Output:
	// 10
}

func ExampleDB_RunInTransactionCtx() {
	rawDb, err := sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	defer rawDb.Close()
	if err != nil {
		log.Panic(err)
	}
	db := New(rawDb, time.Second)
	var id int
	if err := db.RunInTransactionCtx(
		context.Background(), "test RunInTransactionCtx", func(tx *Tx, ctx context.Context) error {
			if err := tx.Query(&id, `select 10 as id`); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		log.Panic(err)
	}
	fmt.Println(id)
	// Output:
	// 10
}

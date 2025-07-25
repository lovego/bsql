package bsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

var rawDB *sql.DB

func init() {
	db, err := sql.Open("postgres", "postgres://develop:@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	rawDB = db
}

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

func getTestStudents() []Student {
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
	runTests(t, db)
}

func getTestDB() *DB {
	return &DB{DB: rawDB, Timeout: time.Second, PutSqlInError: false}
}

func ExampleNew() {
	var rawDB *sql.DB
	db := New(rawDB, time.Second)
	fmt.Println(db.Timeout)
	// Output:
	// 1s
}

func ExampleDB_Query() {
	var people struct {
		Name string
		Age  int
	}
	db := New(rawDB, time.Second)
	if err := db.Query(&people, `select 'jack' as name, 24 as age`); err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v", people)
	// Output:
	// {Name:jack Age:24}
}

func ExampleDB_Query_no_reuse() {
	var people = []struct {
		Name string
		Age  int
	}{
		{
			Name: "hhh",
		},
	}
	db := New(rawDB, time.Second)
	if err := db.Query(&people, `select 'jack' as name, 24 as age`); err != nil {
		log.Panic(err)
	}
	for i := range people {
		fmt.Printf("%+v\n", people[i])
	}
	// Output:
	// {Name:hhh Age:0}
	// {Name:jack Age:24}
}

func ExampleDB_QueryR_reuse() {
	var people = []struct {
		Name string
		Age  int
	}{
		{
			Name: "hhh",
		},
	}
	db := New(rawDB, time.Second)
	if err := db.QueryR(&people, `select 'jack' as name, 24 as age`); err != nil {
		log.Panic(err)
	}
	for i := range people {
		fmt.Printf("%+v\n", people[i])
	}
	// Output:
	// {Name:jack Age:24}
}

func ExampleDB_QueryT() {
	var people struct {
		Name string
		Age  int
	}
	db := New(rawDB, time.Second)
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
	db := New(rawDB, time.Second)
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
	db := New(rawDB, time.Second)
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
	db := New(rawDB, time.Second)
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
	db := New(rawDB, time.Second)
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
	db := New(rawDB, time.Second)
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
	db := New(rawDB, time.Second)
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
	db := New(rawDB, time.Second)
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

package bsql

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/lovego/deep"
)

func TestRunInTransactionAndTx(t *testing.T) {
	db := getTestDB()
	if err := db.RunInTransaction(func(tx *Tx) error {
		runTests(t, tx)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func runTests(t *testing.T, db DbOrTx) {
	createTable(t, db)
	var expects = getTestStudents()
	// Rows
	var gots []Student
	var fields = FieldsFromStruct(Student{}, []string{"UpdatedAt"})
	var columns = Fields2ColumnsStr(fields)
	sql := fmt.Sprintf(
		`insert into students (%s) values %s returning *`, columns,
		StructValues(expects, fields),
	)
	if err := db.Query(&gots, sql); err != nil {
		t.Log(sql)
		t.Fatal(err)
	}
	if diff := deep.Equal(expects, gots); len(diff) > 0 {
		t.Error(strings.Join(diff, "\n"))
	}

	// Row
	var got Student
	if err := db.Query(&got, `select * from students where id=$1`, 1); err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(expects[0], got); len(diff) > 0 {
		t.Error(strings.Join(diff, "\n"))
	}

	// strings
	var gotNames []string
	sql = fmt.Sprintf(`select name from students where id in %s order by id`, Values([]int{1, 2}))
	if err := db.Query(&gotNames, sql); err != nil {
		t.Log(sql)
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotNames, []string{"李雷", "韩梅梅"}) {
		t.Errorf("unexpected: %v", gotNames)
	}

	// ints
	var gotIds []int64
	if err := db.Query(&gotIds,
		`update students set updated_at = now() returning id`,
	); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotIds, []int64{1, 2, 3}) {
		t.Errorf("unexpected: %v", gotIds)
	}
}

func ExampleTx_Query() {
	var people struct {
		Name string
		Age  int
	}
	rawTx, err := rawDB.Begin()
	if err != nil {
		log.Panic(err)
	}
	var tx = &Tx{
		tx:      rawTx,
		timeout: time.Second,
	}
	if err := tx.Query(&people, `select 'jack' as name, 24 as age`); err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v", people)
	// Output:
	// {Name:jack Age:24}
}

func ExampleTx_QueryT() {
	var people struct {
		Name string
		Age  int
	}
	rawTx, err := rawDB.Begin()
	if err != nil {
		log.Panic(err)
	}
	var tx = &Tx{
		tx:      rawTx,
		timeout: time.Second,
	}
	if err := tx.QueryT(2*time.Second, &people, `select 'jack' as name, 24 as age`); err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v", people)
	// Output:
	// {Name:jack Age:24}
}

func ExampleTx_QueryCtx() {
	var people struct {
		Name string
		Age  int
	}
	rawTx, err := rawDB.Begin()
	if err != nil {
		log.Panic(err)
	}
	var tx = &Tx{
		tx:      rawTx,
		timeout: time.Second,
	}
	if err := tx.QueryCtx(
		context.Background(), `query people`, &people, `select 'jack' as name, 24 as age`,
	); err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v", people)
	// Output:
	// {Name:jack Age:24}
}

func ExampleTx_Exec() {
	db := New(rawDB, time.Second)
	var affectedRows int64
	if err := db.RunInTransaction(func(tx *Tx) error {
		result, err := tx.Exec(`
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
			return err
		}
		affectedRows, err = result.RowsAffected()
		return err
	}); err != nil {
		log.Panic(err)
	}
	fmt.Println(affectedRows)
	// Output:
	// 1
}

func ExampleTx_ExecT() {
	db := New(rawDB, time.Second)
	var affectedRows int64
	if err := db.RunInTransaction(func(tx *Tx) error {
		result, err := tx.ExecT(time.Second, `
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
			return err
		}
		affectedRows, err = result.RowsAffected()
		return err
	}); err != nil {
		log.Panic(err)
	}
	fmt.Println(affectedRows)
	// Output:
	// 2
}

func ExampleTx_ExecCtx() {
	db := New(rawDB, time.Second)
	var affectedRows int64
	if err := db.RunInTransaction(func(tx *Tx) error {
		result, err := tx.ExecCtx(
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
			return err
		}
		affectedRows, err = result.RowsAffected()
		return err
	}); err != nil {
		log.Panic(err)
	}
	fmt.Println(affectedRows)
	// Output:
	// 2
}

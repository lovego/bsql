package scan

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:@localhost/bsql_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
}

type statusType int8

type Student struct {
	Id        int64
	Name      string
	Cities    []string
	FriendIds pq.Int64Array
	Scores    []interface{}
	Money     decimal.Decimal
	Status    statusType
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func getTestStudents() *sql.Rows {
	return getTestRows(`
select * from (values
(1, '李雷',   '{成都,上海}'::text[], '{1001,1002}'::int[], '["语文",99,"数学",100]'::JSON,
'25.04', 0, '2001-09-01 12:25:48+08'::timestamptz, NULL),
(2, '韩梅梅', '{广州,北京}'::text[], '{1001,1003}'::int[], '["语文",98,"数学",95]'::JSON,
'95.90', 0, '2001-09-01 10:25:48+08'::timestamptz, '2001-09-02 10:25:58+08'::timestamptz)
) as tmp(id, name, cities, friend_ids, scores, money, status, created_at, updated_at)
`)
}

func getTestIntValues() *sql.Rows {
	return getTestRows(`select * from (values (9), (99), (999)) as tmp`)
}

func getTestNull() *sql.Rows {
	return getTestRows(`select null`)
}

func ExampleScan_struct() {
	var row Student
	if err := Scan(getTestStudents(), &row); err != nil {
		log.Panic(err)
	}
	row.CreatedAt = row.CreatedAt.UTC()
	fmt.Printf("{%d %s %v %v %v %v %d\n  %v %v}\n",
		row.Id, row.Name, row.Cities, row.FriendIds, row.Scores, row.Money, row.Status,
		row.CreatedAt, row.UpdatedAt,
	)
	// Output:
	// {1 李雷 [成都 上海] [1001 1002] [语文 99 数学 100] 25.04 0
	//   2001-09-01 04:25:48 +0000 UTC <nil>}
}

func ExampleScan_structSlice() {
	var rows []Student
	if err := Scan(getTestStudents(), &rows); err != nil {
		log.Panic(err)
	}
	for _, row := range rows {
		row.CreatedAt = row.CreatedAt.UTC()
		if row.UpdatedAt != nil {
			t := row.UpdatedAt.UTC()
			row.UpdatedAt = &t
		}
		fmt.Printf("{%d %s %v %v %v %v %d\n  %v %v}\n",
			row.Id, row.Name, row.Cities, row.FriendIds, row.Scores, row.Money, row.Status,
			row.CreatedAt, row.UpdatedAt,
		)
	}
	// Output:
	// {1 李雷 [成都 上海] [1001 1002] [语文 99 数学 100] 25.04 0
	//   2001-09-01 04:25:48 +0000 UTC <nil>}
	// {2 韩梅梅 [广州 北京] [1001 1003] [语文 98 数学 95] 95.9 0
	//   2001-09-01 02:25:48 +0000 UTC 2001-09-02 02:25:58 +0000 UTC}
}

func ExampleScan_string() {
	var s string
	if err := Scan(getTestRows(`select 'abc'`), &s); err != nil {
		log.Panic(err)
	}
	fmt.Println(s)

	if err := Scan(getTestNull(), &s); err != nil {
		log.Panic(err)
	}
	fmt.Printf("'%s'\n", s)

	// Output:
	// abc
	// ''
}

func ExampleScan_stringPointer() {
	var pointer *string
	if err := Scan(getTestRows(`select 'abc'`), &pointer); err != nil {
		log.Panic(err)
	}
	fmt.Println(*pointer)

	if err := Scan(getTestNull(), &pointer); err != nil {
		log.Panic(err)
	}
	fmt.Println(pointer)

	var p **string
	if err := Scan(getTestRows(`select 'abc'`), &p); err != nil {
		log.Panic(err)
	}
	fmt.Println(**p)

	// Output:
	// abc
	// <nil>
	// abc
}

func ExampleScan_int() {
	var i int
	if err := Scan(getTestIntValues(), &i); err != nil {
		log.Panic(err)
	}
	fmt.Println(i)

	if err := Scan(getTestNull(), &i); err != nil {
		fmt.Println(err)
	}
	fmt.Println(i)

	// Output:
	// 9
	// 0
}

func ExampleScan_uint() {
	var i uint
	if err := Scan(getTestIntValues(), &i); err != nil {
		log.Panic(err)
	}
	fmt.Println(i)

	if err := Scan(getTestNull(), &i); err != nil {
		fmt.Println(err)
	}
	fmt.Println(i)

	// Output:
	// 9
	// 0
}

func ExampleScan_intPointer() {
	var pointer *int
	if err := Scan(getTestIntValues(), &pointer); err != nil {
		log.Panic(err)
	}
	fmt.Println(*pointer)

	if err := Scan(getTestNull(), &pointer); err != nil {
		log.Panic(err)
	}
	fmt.Println(pointer)

	var p **int
	if err := Scan(getTestIntValues(), &p); err != nil {
		log.Panic(err)
	}
	fmt.Println(**p)

	// Output:
	// 9
	// <nil>
	// 9
}

func ExampleScan_intValueOutOfRange() {
	var i int8
	if err := Scan(getTestRows(`select 128`), &i); err != nil {
		fmt.Println(strings.ReplaceAll(err.Error(), `, name "?column?"`, ``))
	}
	// Output:
	// sql: Scan error on column index 0: bsql: cannot assign int64(128) to int8: value out of range
}

func ExampleScan_intSlice() {
	var a []int
	if err := Scan(getTestIntValues(), &a); err != nil {
		log.Panic(err)
	}
	fmt.Println(a)
	// Output: [9 99 999]
}

func ExampleScan_pqInt64Array() {
	var a pq.Int64Array
	if err := Scan(getTestRows(`select '{9,99,999}'::int[]`), &a); err != nil {
		log.Panic(err)
	}
	fmt.Println(a)

	if err := Scan(getTestNull(), &a); err != nil {
		log.Panic(err)
	}
	fmt.Println(a)
	// Output:
	// [9 99 999]
	// []
}

func ExampleScan_float() {
	var f float32
	if err := Scan(getTestRows(`select 1.23::float`), &f); err != nil {
		log.Panic(err)
	}
	fmt.Println(f)
	// Output: 1.23
}

func ExampleScan_bool() {
	var b bool
	if err := Scan(getTestRows(`select true`), &b); err != nil {
		log.Panic(err)
	}
	fmt.Println(b)
	if err := Scan(getTestNull(), &b); err != nil {
		log.Panic(err)
	}
	fmt.Println(b)
	// Output:
	// true
	// false
}

func ExampleScan_time() {
	var t time.Time
	if err := Scan(getTestRows(`select '2001-09-01 12:25:48+08'::timestamptz`), &t); err != nil {
		log.Panic(err)
	}
	fmt.Println(t.UTC())
	// Output: 2001-09-01 04:25:48 +0000 UTC
}

func ExampleScan_scanner() {
	var d ***decimal.Decimal
	if err := Scan(getTestRows(`select '12.34'`), &d); err != nil {
		log.Panic(err)
	}
	fmt.Println(***d)

	if err := Scan(getTestNull(), &d); err != nil {
		log.Panic(err)
	}
	fmt.Println(d)
	// Output:
	// 12.34
	// <nil>
}

func ExampleScan_interface() {
	var s struct {
		A interface{}
	}
	if err := Scan(getTestRows(`SELECT 12 AS a`), &s); err != nil {
		log.Panic(err)
	}
	fmt.Printf("%v %v\n", s, reflect.TypeOf(s.A))

	var u uint64
	s.A = &u
	if err := Scan(getTestRows(`SELECT 12 AS a`), &s); err != nil {
		log.Panic(err)
	}
	value := *s.A.(*uint64)
	fmt.Printf("%d %v %v\n", u, value, reflect.TypeOf(value))

	// Output:
	// {12} int64
	// 12 12 uint64
}

func getTestRows(sql string) *sql.Rows {
	rows, err := db.Query(sql)
	if err != nil {
		log.Panic(err)
	}
	return rows
}

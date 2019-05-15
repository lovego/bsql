package bsql

import "fmt"

type fakeTable struct {
	Id     int64  `json:"id" sql:"default null" comment:"主键ID"`
	Name   string `json:"name" sql:"default null" comment:"名称"`
	Field1 string `json:"field1" comment:"field1"`
	Field2 string `json:"field2" comment:"field2"`
}

func ExampleTableSql_1() {
	constraints := []string{
		`UNIQUE(name)`,
		`UNIQUE(name,field1),`,
	}

	extSqls := []string{
		`
		CREATE INDEX IF NOT EXISTS order_fields_index1 ON pur.orders
		(name,field1)

		`,
		`
		CREATE INDEX IF NOT EXISTS order_fields_index2 ON pur.orders
		(name,field2);
		`,
	}
	fmt.Println(TableSql(`pur.orders`, &fakeTable{}, constraints, extSqls))

	// Output:
	// CREATE TABLE IF NOT EXISTS pur.orders (
	//   id serial8 primary key default null,
	//   name text default null,
	//   field1 text not null,
	//   field2 text not null,
	//   UNIQUE(name),
	//   UNIQUE(name,field1)
	// );
	// CREATE INDEX IF NOT EXISTS order_fields_index1 ON pur.orders
	// 		(name,field1);
	// CREATE INDEX IF NOT EXISTS order_fields_index2 ON pur.orders
	// 		(name,field2);
	// comment on column pur.orders.id is '主键ID';
	// comment on column pur.orders.name is '名称';
	// comment on column pur.orders.field1 is 'field1';
	// comment on column pur.orders.field2 is 'field2';
}

func ExampleTableSql_2() {
	fmt.Println(TableSql(`pur.orders`, &fakeTable{}, nil, nil))
	// Output:
	// CREATE TABLE IF NOT EXISTS pur.orders (
	//   id serial8 primary key default null,
	//   name text default null,
	//   field1 text not null,
	//   field2 text not null
	// );
	//
	// comment on column pur.orders.id is '主键ID';
	// comment on column pur.orders.name is '名称';
	// comment on column pur.orders.field1 is 'field1';
	// comment on column pur.orders.field2 is 'field2';
}

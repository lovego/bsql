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
		CREATE INDEX IF NOT EXISTS order_fields_index1 ON pur.orders (name,field1)

		`,
		`
		CREATE INDEX IF NOT EXISTS order_fields_index2 ON pur.orders (name,field2);
		`,
	}
	fmt.Println(Table{
		Name:        `pur.orders`,
		Desc:        `订单表`,
		Struct:      &fakeTable{},
		Constraints: constraints,
		ExtraSqls:   extSqls,
	}.Sql())

	// Output:
	// CREATE TABLE IF NOT EXISTS pur.orders (
	//   id serial8 PRIMARY KEY default null,
	//   name text default null,
	//   field1 text NOT NULL,
	//   field2 text NOT NULL,
	//   UNIQUE(name),
	//   UNIQUE(name,field1)
	// );
	// CREATE INDEX IF NOT EXISTS order_fields_index1 ON pur.orders (name,field1);
	// CREATE INDEX IF NOT EXISTS order_fields_index2 ON pur.orders (name,field2);
	// COMMENT ON TABLE pur.orders is '订单表';
	// COMMENT ON COLUMN pur.orders.id IS '主键ID';
	// COMMENT ON COLUMN pur.orders.name IS '名称';
	// COMMENT ON COLUMN pur.orders.field1 IS 'field1';
	// COMMENT ON COLUMN pur.orders.field2 IS 'field2';
}

func ExampleTableSql_2() {
	fmt.Println(Table{
		Name:        `pur.orders`,
		Desc:        `采购单`,
		Struct:      &fakeTable{},
		Constraints: nil,
		ExtraSqls:   nil,
	}.Sql())
	// Output:
	// CREATE TABLE IF NOT EXISTS pur.orders (
	//   id serial8 PRIMARY KEY default null,
	//   name text default null,
	//   field1 text NOT NULL,
	//   field2 text NOT NULL
	// );
	// COMMENT ON TABLE pur.orders is '采购单';
	// COMMENT ON COLUMN pur.orders.id IS '主键ID';
	// COMMENT ON COLUMN pur.orders.name IS '名称';
	// COMMENT ON COLUMN pur.orders.field1 IS 'field1';
	// COMMENT ON COLUMN pur.orders.field2 IS 'field2';
}

func ExampleTableSql_3() {
	fmt.Println(Table{
		Name:        `pur.orders`,
		Desc:        `采购单`,
		Struct:      &fakeTable{},
		Constraints: nil,
		Options:     []string{`with (fillfactor = 70)`},
		ExtraSqls:   nil,
	}.Sql())
	// Output:
	// CREATE TABLE IF NOT EXISTS pur.orders (
	//   id serial8 PRIMARY KEY default null,
	//   name text default null,
	//   field1 text NOT NULL,
	//   field2 text NOT NULL
	// )
	// with (fillfactor = 70);
	// COMMENT ON TABLE pur.orders is '采购单';
	// COMMENT ON COLUMN pur.orders.id IS '主键ID';
	// COMMENT ON COLUMN pur.orders.name IS '名称';
	// COMMENT ON COLUMN pur.orders.field1 IS 'field1';
	// COMMENT ON COLUMN pur.orders.field2 IS 'field2';
}

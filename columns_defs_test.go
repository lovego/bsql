package bsql

import "fmt"

func ExampleColumnsDefs() {
	fmt.Println(ColumnsDefs(Student{}))
	// Output:
	// id serial8 not null primary key,
	// name text not null,
	// friend_ids int[] not null,
	// cities jsonb not null,
	// scores jsonb not null,
	// money decimal not null,
	// status int2 not null default 0,
	// created_at timestamptz not null,
	// updated_at timestamptz not null
}

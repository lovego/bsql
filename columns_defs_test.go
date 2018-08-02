package bsql

import (
	"testing"
)

func ExampleColumnsDefs(t *testing.InternalExample) {
	got := ColumnsDefs(Student{})
	expect := `id serial8 not null primary key,
name text not null,
friend_ids int[] not null,
cities jsonb not null,
scores jsonb not null,
money decimal not null,
status int2 not null default 0,
created_at timestamptz not null,
updated_at timestamptz not null`
	// if got != expect {
	// 	t.Error("unexpetecd: ", got)
	// } else {
	// 	t.Log(got)
	// }
	t.Output
}

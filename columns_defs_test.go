package bsql

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

func ExampleColumnsDefs() {
	type Student struct {
		Id        int64
		Name      string
		FriendIds pq.Int64Array `sql:"int[]"`
		Cities    []string
		Scores    map[string]int
		Money     decimal.Decimal
		Status    **int8 `sql:"default 0"`
		CreatedAt time.Time
		UpdatedAt *time.Time
	}

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

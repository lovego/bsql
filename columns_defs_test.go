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
	// id serial8 NOT NULL PRIMARY KEY,
	// name text NOT NULL,
	// friend_ids int[] NOT NULL,
	// cities jsonb NOT NULL,
	// scores jsonb NOT NULL,
	// money decimal NOT NULL,
	// status int2 NOT NULL default 0,
	// created_at timestamptz NOT NULL,
	// updated_at timestamptz NOT NULL
}

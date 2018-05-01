package bsql

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/lovego/deep"
)

func TestRunInTransactionAndTx(t *testing.T) {
	db := getTestDB()
	defer db.db.Close()
	if err := db.RunInTransaction(func(tx *Tx) error {
		runTests(t, tx)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func runTests(t *testing.T, db DbOrTx) {
	createTable(t, db)
	var expects = testStudents()
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
	sql = fmt.Sprintf(`select name from students where scores @> %v order by id`, V(map[string]int{"英语": 97}))
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

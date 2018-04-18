package bsql

import (
	"testing"
)

func TestRunInTransactionAndTx(t *testing.T) {
	var student Student
	if err := testDb.RunInTransaction(func(testTx *Tx) error {
		if _, err := testTx.Exec(`insert into users(id,phone,account,name,status,created_at) values
	 (2,'18380461682','jack999','jack2',0,now())`); err != nil {
			return err
		}
		if err := testTx.Query(&student, `select phone from users where id = 2`); err != nil {
			t.Fatal(err)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

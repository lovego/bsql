package bsql

import (
	"testing"
)

func TestRunInTransactionAndTx(t *testing.T) {
	var user User
	var staff Staffs
	if err := testDb.RunInTransaction(func(testTx *Tx) error {
		if _, err := testTx.Exec(`insert into users(id,phone,account,name,status,created_at) values
	 (2,'18380461682','jack999','jack2',0,now())`); err != nil {
			return err
		}
		if _, err := testTx.Exec(`insert into staffs
			(id, company_id, staff_id, staff_name) values
			(1, 2, 2, 'jack')`); err != nil {
			return err
		}
		if err := testTx.Query(&user, `select phone from users where id = 2`); err != nil {
			t.Fatal(err)
		}
		if err := testTx.Query(&staff, `select staff_id from staffs where id = 1`); err != nil {
			t.Fatal(err)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if user.Phone != `18380461682` {
		t.Logf("unexpected phone: %v", user.Phone)
	}
	if staff.StaffId != 2 {
		t.Logf("unexpected staff_id: %d", staff.StaffId)
	}
}

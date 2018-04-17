package bsql

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"testing"
	"time"

	"github.com/lovego/errs"
)

var testDb *DB

type User struct {
	Id        int64
	Phone     string
	Account   string
	Name      string
	Status    int8
	CreateAt  time.Time
	UpdatedAt time.Time
}

type Staffs struct {
	Id        int64
	CompanyId int64
	StaffId   int64
	StaffName string
	Type      int8
	Status    int8
	SwUserId  int
	CreatedBy int64
	CreatedAt time.Time
	UpdatedBy int64
	UpdatedAt time.Time
}

func init() {
	db, err := sql.Open("postgres", "postgres://develop:@localhost/accounts_test?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}
	testDb = &DB{db, time.Minute}

	if _, err := db.Exec(`truncate users`); err != nil {
		log.Panic(err)
	}
	if _, err := db.Exec(`insert into users(id,phone,account,name,status,create_at) values
	 (1,'18380461680','jack666','jack',0,now())`); err != nil {
		log.Panic(err)
	}
}

func TestQuery(t *testing.T) {
	var user User
	if err := testDb.Query(&user, `select * from users where id = ?`, 1); err != nil {
		t.Fatal(errs.WithStack(err))
	}
	if user.Phone != `18380461680` {
		t.Logf("unexpected phone: %v", user.Phone)
	}
	t.Log(user)
}

func TestExec(t *testing.T) {
	if _, err := testDb.Exec(`update users set phone = '18380461689' where id = ?'`, 1); err != nil {
		t.Fatal(errs.WithStack(err))
	}
	var user User
	if err := testDb.Query(&user, `select phone from users where id = ?`, 1); err != nil {
		t.Fatal(errs.WithStack(err))
	}
	if user.Phone != `18380461689` {
		t.Logf("unexpected phone: %v", user.Phone)
	}
}

func TestRunInTransaction(t *testing.T) {
	if err := testDb.RunInTransaction(func(testTx *Tx) error {
		if _, err := testTx.Exec(`insert into users(id,phone,account,name,status,create_at) values
	 (2,'18380461682','jack999','jack2',0,now())`); err != nil {
			return errs.Trace(err)
		}
		if _, err := testTx.Exec(`insert into staffs
			(id, company_id, staff_id, staff_name, type, status, created_at) values
			(1, 2, 2, 'jack', 1, 1, now())`); err != nil {
			return errs.Trace(err)
		}
		return nil
	}); err != nil {
		t.Fatal(errs.WithStack(err))
	}
	var user User
	var staff Staffs
	if err := testDb.Query(&user, `select phone from users where id = 2`); err != nil {
		t.Fatal(errs.WithStack(err))
	}
	if user.Phone != `18380461682` {
		t.Logf("unexpected phone: %v", user.Phone)
	}

	if err := testDb.Query(&staff, `select staff_id from staffs where id = 1`); err != nil {
		t.Fatal(errs.WithStack(err))
	}
	if staff.StaffId != 2 {
		t.Logf("unexpected staff_id: %d", staff.StaffId)
	}
}

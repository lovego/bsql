package bsql

import "fmt"

func ExampleUpsertSql_1() {
	fmt.Println(UpsertSql("users", []string{
		"Phone", "Name", "Sex", "Birthday", "Status", "CreatedBy", "CreatedAt",
	}, []string{"Phone"}, []string{"Status"},
	))
	// Output:
	// INSERT INTO users (phone, name, sex, birthday, status, created_by, created_at)
	// VALUES %s
	// ON CONFLICT (phone) DO UPDATE SET
	// (         name,          sex,          birthday,          updated_by,          updated_at) =
	// (excluded.name, excluded.sex, excluded.birthday, excluded.created_by, excluded.created_at)
}

func ExampleUpsertSql_2() {
	fmt.Println(UpsertSql(
		"users", []string{
			"phone", "name", "sex", "birthday", "Status", "created_by", "created_at",
		}, []string{"phone"}, []string{"status"},
	))
	// Output:
	// INSERT INTO users (phone, name, sex, birthday, status, created_by, created_at)
	// VALUES %s
	// ON CONFLICT (phone) DO UPDATE SET
	// (         name,          sex,          birthday,          updated_by,          updated_at) =
	// (excluded.name, excluded.sex, excluded.birthday, excluded.created_by, excluded.created_at)
}

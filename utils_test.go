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

func ExampleOffsetToLineAndColumn() {
	fmt.Println(OffsetToLineAndColumn("", 0))
	fmt.Println(OffsetToLineAndColumn("", 1))

	fmt.Println(OffsetToLineAndColumn("a", 1))
	fmt.Println(OffsetToLineAndColumn("a\nb", 2))
	fmt.Println(OffsetToLineAndColumn("中文3\n\n\n789", 8))

	fmt.Println(OffsetToLineAndColumn("a\r\n", 3))
	fmt.Println(OffsetToLineAndColumn("a\r\nb", 4))

	// Output:
	// 0 0
	// 0 0
	// 1 1
	// 1 2
	// 4 2
	// 1 3
	// 2 1
}

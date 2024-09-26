package bsql

import (
	"fmt"
	"testing"
	"time"
)

func TestBuilder(t *testing.T) {
	ExampleBuilder()
}

func ExampleBuilder() {
	builder := NewBuilder()
	var sql string
	//Select
	sql = builder.Select("a.tableA").
		Alias("t").
		Fields("t.*", "tb.*", "sum(t.num) sum_num").
		LeftJoin("a.tableB", "tb", "t.id=tb.tid").
		Where("t.field1 = 1").
		OrderBy("a.id").
		GroupBy("a.name").
		Having("sum_num > 300").
		ForUpdate().
		Build()
	fmt.Println("select:")
	fmt.Println(sql)
	// print:
	// SELECT t.*,tb.*,sum(t.num) sum_num
	// FROM a.tableA AS t
	// LEFT JOIN a.tableB AS tb ON t.id=tb.tid
	// WHERE (t.field1 = 1)
	// GROUP BY a.name
	// ORDER BY a.id
	// HAVING sum_num > 300
	// FOR UPDATE

	builder.Clear()

	//Delete

	sql = builder.Delete("a.tableA").
		Alias("t").
		Equal("t.id", 1).
		Build()
	fmt.Println("delete:")
	fmt.Println(sql)
	// print:
	// DELETE FROM a.tableA AS t WHERE (t.id = 1)

	builder.Clear()

	// Insert
	now, _ := time.Parse("2006-01-02 15:04:05", "2020-05-01 00:00:00")
	builder = builder.Insert("a.tableA").
		Values([]interface{}{
			User{
				Name:      "a",
				Age:       1,
				IsAdmin:   true,
				Remark:    "aa",
				CreatedBy: 1,
				CreatedAt: &now,
			},
			User{
				Name:      "b",
				Age:       10,
				IsAdmin:   false,
				Remark:    "bb",
				CreatedBy: 1,
				CreatedAt: &now,
			},
		})
	// 如果需要限定只插入某些字段,使用Cols()方法
	// builder.Cols("Name", "Age", "IsAdmin", "CreatedBy", "CreatedAt")
	builder.Build()
	fmt.Println("insert:")
	fmt.Println(sql)
	// print:
	// INSERT INTO a.tableA(name,age,is_admin,remark,created_by,created_at)
	// VALUES
	//     ('a',1,true,'aa',1,'2020-05-01T00:00:00Z'),
	//     ('b',10,false,'bb',1,'2020-05-01T00:00:00Z')

	builder.Clear()

	// update

	builder = builder.Update("a.user").
		Set("name = 'aaa'", "age = 12").
		//Set("age = age + 1").     // Set()支持复杂的update赋值操作
		SetMap(map[string]interface{}{
			"is_admin":  true,
			"update_at": now,
		}).
		Equal("id", 1)
	fmt.Println("update:Set/setMap")
	fmt.Println(builder.Build())
	// print:
	// UPDATE a.user SET
	// name = 'aaa',age = 12,update_at = '2020-05-01T00:00:00Z',is_admin = true
	// WHERE (id = 1)

	// 也可以使用struct方式，SetStruct()方法
	// 使用SetStruct()后，Set()和SetMap()的设置将不会生效
	builder.SetStruct(User{
		Name:      "a",
		Age:       1,
		IsAdmin:   true,
		Remark:    "aa",
		CreatedBy: 1,
		CreatedAt: &now,
	})
	// 如果需要限定SetStruct()中只更新某些字段,使用Cols()方法
	// builder.Cols("Name", "Age", "IsAdmin", "CreatedBy", "CreatedAt")

	fmt.Println("update:SetStruct")
	fmt.Println(builder.Build())
	// print:
	// UPDATE a.user SET
	// (id,name,age,is_admin,remark,updated_by,updated_at) =
	// (0,'a',1,true,'aa',0,NULL)
	// WHERE (id = 1)
}

type User struct {
	Id        int64
	Name      string
	Age       int8
	IsAdmin   bool
	Remark    string     `json:"remark,omitempty" sql:"default ''" comment:"备注"`
	CreatedBy int64      `json:"createdBy,omitempty" comment:"创建者员工ID"`
	CreatedAt *time.Time `json:"createdAt,omitempty" comment:"创建时间"`
	UpdatedBy int64      `json:"updatedBy,omitempty" sql:"default 0" comment:"更新者员工ID"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" sql:"default null" comment:"更新时间"`
}

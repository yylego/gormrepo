// Package example10_test compares the GORM Preload with gormjoin on the same need: each student with the room of its class
// The two tests below give the same rows
//
// Preload, as GORM does it:
//   - the association is a struct field, named with a string, e.g. Preload("Class"), a typo shows up at run time
//   - it sends one SELECT each table, and gives nested objects, which need flattening when a flat row is wanted
//
// gormjoin:
//   - the ON condition and the select columns are written with the generated columns, a typo fails at compile time
//   - it sends one SELECT, and scans the rows straight into a flat result struct
//
// example10_test 在同一个需求上对比 GORM 的 Preload 和 gormjoin：每个学生和所在班级的教室
// 下面两个测试查出的结果一样
//
// Preload，GORM 官方的做法：
//   - 关联是结构体里的一个字段，用字符串指定，例如 Preload("Class")，写错了运行时才发现
//   - 每张表各发一条查询，得到嵌套的对象，要扁平的结果行时还得自己展开
//
// gormjoin：
//   - 连接条件和选列都用生成的列来写，写错了编译就过不去
//   - 只发一条查询，结果行直接扫进扁平的结果结构体
package example10_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/gormjoin"
	"github.com/yylego/gormrepo/gormtablerepo"
	"github.com/yylego/gormrepo/internal/examples/example10/internal/models"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var caseDB *gorm.DB

var (
	studentRepo = gormtablerepo.NewRepo(gormclass.UseTable(&models.Student{}))
	classRepo   = gormtablerepo.NewRepo(gormclass.UseTable(&models.Class{}))
)

// expected is the rows both ways give
// expected 是两种做法都应该查出的结果
var expected = []models.StudentClass{
	{StudentName: "Alice", ClassRoom: "room-101"},
	{StudentName: "Bob", ClassRoom: "room-102"},
	{StudentName: "Chloe", ClassRoom: ""}, // not in a class yet // 还没分班
}

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	done.Done(db.AutoMigrate(&models.Class{}, &models.Student{}))
	done.Done(db.Create([]*models.Class{{ID: 1, Room: "room-101"}, {ID: 2, Room: "room-102"}}).Error)
	done.Done(db.Create([]*models.Student{
		{ID: 1, Name: "Alice", ClassID: 1},
		{ID: 2, Name: "Bob", ClassID: 2},
		{ID: 3, Name: "Chloe"},
	}).Error)

	caseDB = db
	m.Run()
}

func TestPreload(t *testing.T) {
	// Two SELECTs: the students first, then the classes WHERE id IN (...)
	// 两条查询：先查学生，再按 id IN (...) 查班级
	var students []*models.StudentWithClass
	require.NoError(t, caseDB.Preload("Class").Order("id").Find(&students).Error)

	// The nested objects need flattening into rows
	// 嵌套的对象要自己展开成结果行
	rows := make([]models.StudentClass, 0, len(students))
	for _, student := range students {
		row := models.StudentClass{StudentName: student.Name}
		if student.Class != nil {
			row.ClassRoom = student.Class.Room
		}
		rows = append(rows, row)
	}
	require.Equal(t, expected, rows)
}

func TestJoin(t *testing.T) {
	// One SELECT, the rows go straight into the flat result struct
	// 一条查询，结果行直接进扁平的结果结构体
	join := gormjoin.LeftJoin(studentRepo, classRepo, func(student *models.StudentColumns, class *models.ClassColumns) []string {
		return []string{class.ID.OnEq(student.ClassID)}
	})

	var rows []models.StudentClass
	require.NoError(t, caseDB.Model(&models.Student{}).
		Joins(join.Stmt()).
		Select(join.Names[models.StudentClass](func(student *models.StudentColumns, class *models.ClassColumns, res *models.StudentClassColumns) []string {
			return []string{
				student.Name.AsName(res.StudentName),
				class.Room.AsName(res.ClassRoom),
			}
		})).
		Scopes(studentRepo.Base().NewScope(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
			return db.Order(cls.ID.Asc())
		})).
		Scan(&rows).Error)
	require.Equal(t, expected, rows)
}

// Package example03_test shows AND / OR / NOT conditions and reusable scopes
// cls.Qx wraps one condition, then AND / OR / NOT combine conditions with the parentheses in place
//
// example03_test 演示 AND / OR / NOT 组合条件和可复用的 scope
// cls.Qx 包住一个条件，再用 AND / OR / NOT 组合，括号自动加好
package example03_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/internal/examples/example03/internal/models"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var caseDB *gorm.DB

var studentRepo = gormrepo.NewRepo(gormclass.Use(&models.Student{}))

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	done.Done(db.AutoMigrate(&models.Student{}))
	done.Done(db.Create([]*models.Student{
		{Name: "Alice", Age: 18, Grade: "grade-1"},
		{Name: "Bob", Age: 19, Grade: "grade-1"},
		{Name: "Chloe", Age: 20, Grade: "grade-2"},
		{Name: "David", Age: 21, Grade: "grade-2"},
	}).Error)

	caseDB = db
	m.Run()
}

// findNames returns the names of the students matching the where condition, sorted on ID
// findNames 返回符合条件的学生名字，按 ID 排序
func findNames(t *testing.T, where func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB) []string {
	names, err := studentRepo.With(t.Context(), caseDB).Pluck(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return where(db, cls).Order(cls.ID.Asc())
	}, func(cls *models.StudentColumns) gormcnm.ColumnName[string] {
		return cls.Name
	})
	require.NoError(t, err)
	return names
}

func TestOr(t *testing.T) {
	// A simple OR: GORM's own Or
	// 简单的 OR：直接用 GORM 的 Or
	require.Equal(t, []string{"Alice", "David"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Age.Lt(19)).Or(cls.Age.Gt(20))
	}))
}

func TestAndOrGroups(t *testing.T) {
	// (grade = 'grade-1' AND age >= 19) OR name = 'Chloe'
	require.Equal(t, []string{"Bob", "Chloe"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		qx := cls.Qx(cls.Grade.Eq("grade-1")).AND(cls.Qx(cls.Age.Gte(19))).OR(cls.Qx(cls.Name.Eq("Chloe")))
		return db.Where(qx.Qs(), qx.Args()...)
	}))
}

func TestNot(t *testing.T) {
	// NOT (grade = 'grade-1')
	require.Equal(t, []string{"Chloe", "David"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		qx := cls.Qx(cls.Grade.Eq("grade-1")).NOT()
		return db.Where(qx.Qs(), qx.Args()...)
	}))
}

func TestReusableScopes(t *testing.T) {
	// Write a condition once as a scope, then add it to each find that needs it
	// 把条件写成 scope，任何查询都能加上它
	grade2 := studentRepo.NewScope(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Grade.Eq("grade-2"))
	})
	adult := studentRepo.NewScope(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Age.Gte(21))
	})

	require.Equal(t, []string{"Chloe", "David"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Scopes(grade2)
	}))
	require.Equal(t, []string{"David"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Scopes(grade2, adult)
	}))
}

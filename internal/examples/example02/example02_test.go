// Package example02_test shows the common conditions written with the generated columns, plus Count and Exist
// Each condition is type-checked: cls.Age.Gt takes an int, cls.Name.In takes a []string
//
// example02_test 演示用生成的列写常规条件，以及 Count 和 Exist
// 每个条件都有类型检查：cls.Age.Gt 收 int，cls.Name.In 收 []string
package example02_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/internal/examples/example02/internal/models"
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

	coach := "Coach Wang"
	done.Done(db.AutoMigrate(&models.Student{}))
	done.Done(db.Create([]*models.Student{
		{Name: "Alice", Age: 18, Grade: "grade-1", Coach: &coach},
		{Name: "Bob", Age: 19, Grade: "grade-1"},
		{Name: "Chloe", Age: 20, Grade: "grade-2", Coach: &coach},
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

func TestCompare(t *testing.T) {
	require.Equal(t, []string{"Bob"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("Bob"))
	}))
	require.Equal(t, []string{"Alice", "Chloe", "David"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Ne("Bob"))
	}))
	require.Equal(t, []string{"Chloe", "David"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Age.Gt(19))
	}))
	require.Equal(t, []string{"Alice", "Bob"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Age.Lte(19))
	}))
}

func TestInAndBetween(t *testing.T) {
	require.Equal(t, []string{"Alice", "David"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.In([]string{"Alice", "David"}))
	}))
	require.Equal(t, []string{"Bob", "Chloe"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Age.Between(19, 20))
	}))
}

func TestLikeAndNull(t *testing.T) {
	require.Equal(t, []string{"Alice", "Chloe"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Like("%e%"))
	}))
	require.Equal(t, []string{"Bob", "David"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Coach.IsNull())
	}))
}

func TestConditionsChained(t *testing.T) {
	// Chained Where calls are joined with AND
	// 多次 Where 用 AND 连接
	require.Equal(t, []string{"Chloe"}, findNames(t, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Grade.Eq("grade-2")).Where(cls.Coach.IsNotNull())
	}))
}

func TestCountAndExist(t *testing.T) {
	repo := studentRepo.With(t.Context(), caseDB)

	count, err := repo.Count(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Grade.Eq("grade-1"))
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), count)

	exist, err := repo.Exist(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Age.Gt(30))
	})
	require.NoError(t, err)
	require.False(t, exist)
}

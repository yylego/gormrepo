// Package example05_test shows creating and updating rows
// The new values are written with the generated columns, e.g. cls.Score.Kv(95), so the column names and the value types can not go wrong
//
// example05_test 演示新增和更新
// 新值用生成的列来写，例如 cls.Score.Kv(95)，列名和值的类型都不会写错
package example05_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/internal/examples/example05/internal/models"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var studentRepo = gormrepo.NewRepo(gormclass.Use(&models.Student{}))

// newCaseDB opens a new database in each test, since these tests change the data
// newCaseDB 给每个测试开一个新库，因为这些测试会改数据
func newCaseDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	t.Cleanup(func() { rese.F0(rese.P1(db.DB()).Close) })

	done.Done(db.AutoMigrate(&models.Student{}))
	return db
}

func TestCreate(t *testing.T) {
	repo := studentRepo.With(t.Context(), newCaseDB(t))

	// Create fills the ID back into the object
	// Create 会把 ID 回填到对象里
	alice := &models.Student{Name: "Alice", Grade: "grade-1", Score: 90}
	require.NoError(t, repo.Create(alice))
	require.NotZero(t, alice.ID)

	require.NoError(t, repo.Creates([]*models.Student{
		{Name: "Bob", Grade: "grade-1", Score: 75},
		{Name: "Chloe", Grade: "grade-2", Score: 80},
	}))

	// CreateInBatches inserts a batch of rows, batchSize rows each SQL
	// CreateInBatches 批量插入，每条 SQL 插 batchSize 行
	require.NoError(t, repo.CreateInBatches([]*models.Student{
		{Name: "David", Grade: "grade-2", Score: 60},
		{Name: "Eric", Grade: "grade-2", Score: 85},
		{Name: "Frank", Grade: "grade-2", Score: 70},
	}, 2))

	count, err := repo.Count(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db
	})
	require.NoError(t, err)
	require.Equal(t, int64(6), count)
}

func TestSave(t *testing.T) {
	repo := studentRepo.With(t.Context(), newCaseDB(t))

	// Save inserts when the ID is zero, and updates each column when the ID is set
	// ID 为零时 Save 插入，ID 有值时 Save 更新所有列
	student := &models.Student{Name: "Alice", Grade: "grade-1", Score: 90}
	require.NoError(t, repo.Save(student))

	student.Score = 95
	require.NoError(t, repo.Save(student))

	res, err := repo.First(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.ID.Eq(student.ID))
	})
	require.NoError(t, err)
	require.Equal(t, 95, res.Score)
}

func TestUpdate(t *testing.T) {
	repo := studentRepo.With(t.Context(), newCaseDB(t))
	require.NoError(t, repo.Creates([]*models.Student{
		{Name: "Alice", Grade: "grade-1", Score: 90},
		{Name: "Bob", Grade: "grade-1", Score: 75},
	}))

	// One column
	// 更新一列
	require.NoError(t, repo.Update(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("Bob"))
	}, func(cls *models.StudentColumns) (string, any) {
		return cls.Score.Kv(80)
	}))

	// An expression, score = score + 5
	// 用表达式更新，score = score + 5
	require.NoError(t, repo.Update(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("Bob"))
	}, func(cls *models.StudentColumns) (string, any) {
		return cls.Score.KeAdd(5)
	}))

	// Two columns and up
	// 更新多列
	require.NoError(t, repo.UpdatesM(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("Alice"))
	}, func(cls *models.StudentColumns) gormcnm.ColumnValueMap {
		return cls.Kw(cls.Grade.Kv("grade-2")).Kw(cls.Score.Kv(92))
	}))

	students, err := repo.FindO(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db
	}, func(cls *models.StudentColumns) gormcnm.Ordering {
		return cls.ID.Ob("asc")
	})
	require.NoError(t, err)
	require.Equal(t, "grade-2", students[0].Grade)
	require.Equal(t, 92, students[0].Score)
	require.Equal(t, 85, students[1].Score)
}

func TestUpdatesObject(t *testing.T) {
	repo := studentRepo.With(t.Context(), newCaseDB(t))
	alice := &models.Student{Name: "Alice", Grade: "grade-1", Score: 90}
	require.NoError(t, repo.Create(alice))

	// UpdatesO updates the row with the ID of the object
	// UpdatesO 按对象的主键更新这一行
	require.NoError(t, repo.UpdatesO(alice, func(cls *models.StudentColumns) gormcnm.ColumnValueMap {
		return cls.Kw(cls.Score.Kv(93))
	}))

	// UpdatesC adds a where condition on top of the ID, it changes nothing when the condition fails
	// UpdatesC 在主键之外再加条件，条件不满足时什么都不改
	require.NoError(t, repo.UpdatesC(alice, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Grade.Eq("grade-2"))
	}, func(cls *models.StudentColumns) gormcnm.ColumnValueMap {
		return cls.Kw(cls.Score.Kv(0))
	}))

	res, err := repo.First(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.ID.Eq(alice.ID))
	})
	require.NoError(t, err)
	require.Equal(t, 93, res.Score)
}

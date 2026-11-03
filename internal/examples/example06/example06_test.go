// Package example06_test shows deleting rows, soft delete and Unscoped
// The Student struct has a gorm.DeletedAt field, so a delete just sets deleted_at and the finds skip the row
// Unscoped reads the soft deleted rows too, and removes the row from the table
//
// example06_test 演示删除、软删除和 Unscoped
// 模型有 gorm.DeletedAt 字段，删除只是填上 deleted_at，查询时自动跳过这一行
// Unscoped 能读到软删除的行，也能真正把行删掉
package example06_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/internal/examples/example06/internal/models"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var studentRepo = gormrepo.NewRepo(gormclass.Use(&models.Student{}))

// newCaseDB opens a new database with three students in each test, since these tests delete rows
// newCaseDB 给每个测试开一个有三个学生的新库，因为这些测试会删数据
func newCaseDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	t.Cleanup(func() { rese.F0(rese.P1(db.DB()).Close) })

	done.Done(db.AutoMigrate(&models.Student{}))
	done.Done(db.Create([]*models.Student{
		{Name: "Alice", Grade: "grade-1"},
		{Name: "Bob", Grade: "grade-1"},
		{Name: "Chloe", Grade: "grade-2"},
	}).Error)
	return db
}

// countRows counts the rows the finds can see
// countRows 统计查询能看到的行数
func countRows(t *testing.T, repo *gormrepo.GormRepo[models.Student, *models.StudentColumns]) int64 {
	count, err := repo.Count(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db
	})
	require.NoError(t, err)
	return count
}

func TestDelete(t *testing.T) {
	repo := studentRepo.With(t.Context(), newCaseDB(t))

	alice, err := repo.First(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("Alice"))
	})
	require.NoError(t, err)

	// Delete removes the row with its ID
	// Delete 按主键删除这一行
	require.NoError(t, repo.Delete(alice))
	require.Equal(t, int64(2), countRows(t, repo))
}

func TestDeleteW(t *testing.T) {
	repo := studentRepo.With(t.Context(), newCaseDB(t))

	// DeleteW removes the rows matching the where condition
	// DeleteW 删除符合条件的行
	require.NoError(t, repo.DeleteW(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Grade.Eq("grade-1"))
	}))
	require.Equal(t, int64(1), countRows(t, repo))
}

func TestDeleteM(t *testing.T) {
	repo := studentRepo.With(t.Context(), newCaseDB(t))

	chloe, err := repo.First(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("Chloe"))
	})
	require.NoError(t, err)

	// DeleteM deletes the row just when the extra condition holds too
	// DeleteM 只有附加条件也满足时才删除
	require.NoError(t, repo.DeleteM(chloe, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Grade.Eq("grade-1"))
	}))
	require.Equal(t, int64(3), countRows(t, repo))

	require.NoError(t, repo.DeleteM(chloe, func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Grade.Eq("grade-2"))
	}))
	require.Equal(t, int64(2), countRows(t, repo))
}

func TestUnscoped(t *testing.T) {
	repo := studentRepo.With(t.Context(), newCaseDB(t))

	require.NoError(t, repo.DeleteW(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("Bob"))
	}))

	// The soft deleted row remains in the table, Unscoped reads it
	// 软删除的行还在表里，Unscoped 能读到
	require.Equal(t, int64(2), countRows(t, repo))
	require.Equal(t, int64(3), countRows(t, repo.Unscoped()))

	// Unscoped removes the row from the table
	// Unscoped 真正把这一行删掉
	require.NoError(t, repo.Unscoped().DeleteW(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("Bob"))
	}))
	require.Equal(t, int64(2), countRows(t, repo.Unscoped()))
}

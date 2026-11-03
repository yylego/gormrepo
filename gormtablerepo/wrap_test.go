package gormtablerepo_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/gormtablerepo"
	"github.com/yylego/gormrepo/gormtablerepo/internal/models"
	"github.com/yylego/must"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestRepo_Base(t *testing.T) {
	repo := gormtablerepo.NewRepo(gormclass.UseTable(&models.Student{}))
	scope := repo.Base().NewScope(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("A"))
	})

	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	done.Done(db.AutoMigrate(&models.Student{}))
	must.Done(db.Save(&models.Student{ID: 0, Name: "A", Rank: 100}).Error)
	must.Done(db.Save(&models.Student{ID: 0, Name: "B", Rank: 85}).Error)

	// The scope from Base uses the table-prefixed columns, e.g. students.name
	// Base 得到的 scope 用的是带表名前缀的列，例如 students.name
	var students []*models.Student
	require.NoError(t, db.Scopes(scope).Find(&students).Error)
	require.Len(t, students, 1)
	require.Equal(t, "A", students[0].Name)
}

func TestRepo_Repo(t *testing.T) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	done.Done(db.AutoMigrate(&models.Student{}))

	must.Done(db.Save(&models.Student{ID: 0, Name: "A", Rank: 100}).Error)
	must.Done(db.Save(&models.Student{ID: 0, Name: "B", Rank: 85}).Error)
	must.Done(db.Save(&models.Student{ID: 0, Name: "C", Rank: 90}).Error)

	repo := gormtablerepo.NewRepo(gormclass.UseTable(&models.Student{}))
	students, count, err := repo.Repo(db).FindPageAndCount(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.In([]string{"A", "B", "C"}))
	}, func(cls *models.StudentColumns) gormcnm.Ordering {
		return cls.Rank.Ob("desc")
	}, &gormrepo.Pagination{
		Limit:  3,
		Offset: 1,
	})
	require.NoError(t, err)
	require.Equal(t, int64(3), count) // Expecting 3 students in aggregate count without pagination
	require.Len(t, students, 2)       // Expecting 2 students in the pagination result
	require.Equal(t, "C", students[0].Name)
	require.Equal(t, "B", students[1].Name)
}

func TestRepo_Gorm(t *testing.T) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	done.Done(db.AutoMigrate(&models.Student{}))

	must.Done(db.Save(&models.Student{ID: 0, Name: "A", Rank: 100}).Error)
	must.Done(db.Save(&models.Student{ID: 0, Name: "B", Rank: 85}).Error)
	must.Done(db.Save(&models.Student{ID: 0, Name: "C", Rank: 90}).Error)

	repo := gormtablerepo.NewRepo(gormclass.UseTable(&models.Student{}))
	var students []*models.Student
	pagination := &gormrepo.Pagination{
		Limit:  3,
		Offset: 1,
	}
	err := repo.Gorm(db).Find(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.In([]string{"A", "B", "C"})).
			Order(cls.Rank.Desc()).
			Scopes(pagination.Scope())
	}, &students).Error
	require.NoError(t, err)
	var count int64
	err = repo.Gorm(db).Modeled().Invoke(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.In([]string{"A", "B", "C"})).Count(&count)
	}).Error
	require.NoError(t, err)
	require.Equal(t, int64(3), count) // Expecting 3 students in aggregate count without pagination
	require.Len(t, students, 2)       // Expecting 2 students in the pagination result
	require.Equal(t, "C", students[0].Name)
	require.Equal(t, "B", students[1].Name)
}

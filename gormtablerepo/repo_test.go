package gormtablerepo_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/gormtablerepo"
	"github.com/yylego/gormrepo/gormtablerepo/internal/models"
	"github.com/yylego/must"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestNewRepo(t *testing.T) {
	repo := gormtablerepo.NewRepo(gormclass.UseTable(&models.Student{}))
	require.NotNil(t, repo)
}

func TestRepo_TableName(t *testing.T) {
	repo := gormtablerepo.NewRepo(gormclass.UseTable(&models.Student{}))
	require.Equal(t, "students", repo.TableName())
}

func TestRepo_TableColumns(t *testing.T) {
	repo := gormtablerepo.NewRepo(gormclass.UseTable(&models.Student{}))
	cls := repo.TableColumns()
	require.Equal(t, "students.name", cls.Name.Name())
	require.Equal(t, "students.rank", cls.Rank.Name())
}

func TestRepo_BuildColumns(t *testing.T) {
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
	var result models.StudentScore
	res := result.Columns()
	err := repo.Repo(db).Modeled().Invoke(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Select(res.Merge(repo.BuildColumns(func(cls *models.StudentColumns) []string {
			return []string{
				cls.Name.AsName(res.Name),
				cls.Rank.AsName(res.Score),
			}
		})...)).Where(cls.Name.Eq("B")).Take(&result)
	})
	require.NoError(t, err)
	require.Equal(t, "B", result.Name)
	require.Equal(t, 85, result.Score)
}

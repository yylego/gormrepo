// Package example01_test shows the basics: define a struct, generate its columns, build a repo and find rows
// The columns of Student are generated into internal/models/models.gen.go via go generate
//
// example01_test 演示入门用法：定义模型、生成列、建 repo、查询
// Student 的列由 go generate 生成到 internal/models/models.gen.go
package example01_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/internal/examples/example01/internal/models"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var caseDB *gorm.DB

// The repo is built once and used everywhere, it holds no db, pass the db when querying via With
// repo 建一次到处用，它不带 db，查询时通过 With 传入
var studentRepo = gormrepo.NewRepo(gormclass.Use(&models.Student{}))

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	done.Done(db.AutoMigrate(&models.Student{}))
	done.Done(db.Create([]*models.Student{
		{Name: "Alice", Age: 18},
		{Name: "Bob", Age: 19},
		{Name: "Chloe", Age: 20},
	}).Error)

	caseDB = db
	m.Run()
}

func TestFirst(t *testing.T) {
	student, err := studentRepo.With(t.Context(), caseDB).First(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("Bob"))
	})
	require.NoError(t, err)
	require.Equal(t, 19, student.Age)
}

func TestFirstE(t *testing.T) {
	// FirstE tells "no such row" apart from the rest of the errors
	// FirstE 能把"没有这一行"和其他错误区分开
	student, erb := studentRepo.With(t.Context(), caseDB).FirstE(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("David"))
	})
	require.NotNil(t, erb)
	require.True(t, erb.NotExist)
	require.Nil(t, student)
}

func TestFind(t *testing.T) {
	students, err := studentRepo.With(t.Context(), caseDB).Find(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Where(cls.Age.Gte(19)).Order(cls.Age.Asc())
	})
	require.NoError(t, err)
	require.Len(t, students, 2)
	require.Equal(t, "Bob", students[0].Name)
	require.Equal(t, "Chloe", students[1].Name)
}

func TestFindAll(t *testing.T) {
	students, err := studentRepo.With(t.Context(), caseDB).FindAll()
	require.NoError(t, err)
	require.Len(t, students, 3)
}

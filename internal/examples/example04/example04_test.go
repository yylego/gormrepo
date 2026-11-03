// Package example04_test shows sorting and paging
// FindO sorts, FindN takes the first N, FindPage / FindPageAndCount page with a Pagination, FindC pages with a custom callback
// Each one takes the ordering as a required param, since the first N rows and the pages are not stable without ORDER BY
//
// example04_test 演示排序与分页
// FindO 排序，FindN 取前 N 条，FindPage / FindPageAndCount 按 Pagination 分页，FindC 用自定义回调分页
// 这几个方法都把排序作为必填参数，因为取前 N 条或分页时不排序，拿到的行不稳定
package example04_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/internal/examples/example04/internal/models"
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
		{Name: "Alice", Grade: "grade-1", Score: 90},
		{Name: "Bob", Grade: "grade-1", Score: 75},
		{Name: "Chloe", Grade: "grade-1", Score: 90},
		{Name: "David", Grade: "grade-1", Score: 60},
		{Name: "Eric", Grade: "grade-1", Score: 85},
		{Name: "Frank", Grade: "grade-2", Score: 95},
	}).Error)

	caseDB = db
	m.Run()
}

// namesOf returns the names of the students in sequence
// namesOf 按顺序返回学生的名字
func namesOf(students []*models.Student) []string {
	names := make([]string, 0, len(students))
	for _, student := range students {
		names = append(names, student.Name)
	}
	return names
}

func TestFindO(t *testing.T) {
	// Grade 1, from the high score to the low, the same score sorted on the name
	// 一年级，分数从高到低，分数相同按名字排
	students, err := studentRepo.With(t.Context(), caseDB).FindO(
		func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
			return db.Where(cls.Grade.Eq("grade-1"))
		},
		func(cls *models.StudentColumns) gormcnm.Ordering {
			return cls.Score.Ob("desc").Ob(cls.Name.Ob("asc"))
		},
	)
	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Chloe", "Eric", "Bob", "David"}, namesOf(students))
}

func TestFindN(t *testing.T) {
	// The top 2 of grade 1
	// 一年级的前 2 名
	students, err := studentRepo.With(t.Context(), caseDB).FindN(
		func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
			return db.Where(cls.Grade.Eq("grade-1"))
		},
		func(cls *models.StudentColumns) gormcnm.Ordering {
			return cls.Score.Ob("desc").Ob(cls.Name.Ob("asc"))
		},
		2,
	)
	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Chloe"}, namesOf(students))
}

func TestFindPage(t *testing.T) {
	// Page 2 with 2 rows each page
	// 每页 2 条，取第 2 页
	students, err := studentRepo.With(t.Context(), caseDB).FindPage(
		func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
			return db.Where(cls.Grade.Eq("grade-1"))
		},
		func(cls *models.StudentColumns) gormcnm.Ordering {
			return cls.Score.Ob("desc").Ob(cls.Name.Ob("asc"))
		},
		&gormrepo.Pagination{Limit: 2, Offset: 2},
	)
	require.NoError(t, err)
	require.Equal(t, []string{"Eric", "Bob"}, namesOf(students))
}

func TestFindPageAndCount(t *testing.T) {
	// Page 3 with 2 rows each page, and the count of matching rows across pages
	// 每页 2 条取第 3 页，加上不分页时符合条件的总数
	students, count, err := studentRepo.With(t.Context(), caseDB).FindPageAndCount(
		func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
			return db.Where(cls.Grade.Eq("grade-1"))
		},
		func(cls *models.StudentColumns) gormcnm.Ordering {
			return cls.Score.Ob("desc").Ob(cls.Name.Ob("asc"))
		},
		&gormrepo.Pagination{Limit: 2, Offset: 4},
	)
	require.NoError(t, err)
	require.Equal(t, []string{"David"}, namesOf(students))
	require.Equal(t, int64(5), count)
}

func TestFindC(t *testing.T) {
	// FindC fits paging that needs code: here the page is cut just when a page index is given, 0 gives each row
	// FindC 适合分页部分要写逻辑的场景：这里只有传了页码才分页，传 0 返回全部
	page := 2
	students, count, err := studentRepo.With(t.Context(), caseDB).FindC(
		func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
			return db.Where(cls.Grade.Eq("grade-1"))
		},
		func(cls *models.StudentColumns) gormcnm.Ordering {
			return cls.Score.Ob("desc").Ob(cls.Name.Ob("asc"))
		},
		func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
			// db comes with the where and the ordering on it // db 上已经带了条件和排序
			if page > 0 {
				db = db.Offset((page - 1) * 2).Limit(2)
			}
			return db
		},
	)
	require.NoError(t, err)
	require.Equal(t, []string{"Eric", "Bob"}, namesOf(students))
	require.Equal(t, int64(5), count)
}

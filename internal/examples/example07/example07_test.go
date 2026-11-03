// Package example07_test shows the finds whose result is not the table struct itself
// Pluck gives one column, Scan gives rows of any struct, ScanOne gives one value such as an aggregate
// ScanList / ScanItem give rows of a result struct with generated columns, so each alias in Select is type-checked
//
// example07_test 演示结果不是模型本身的查询
// Pluck 取一列，Scan 把结果行扫进任意结构体，ScanOne 取一个值（如聚合结果）
// ScanList / ScanItem 把结果行扫进带生成列的结果结构体，Select 里的每个别名都有类型检查
package example07_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/internal/examples/example07/internal/models"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var caseDB *gorm.DB

var scoreRepo = gormrepo.NewRepo(gormclass.Use(&models.Score{}))

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	done.Done(db.AutoMigrate(&models.Score{}))
	done.Done(db.Create([]*models.Score{
		{StudentName: "Alice", Subject: "math", Points: 90},
		{StudentName: "Bob", Subject: "math", Points: 70},
		{StudentName: "Chloe", Subject: "math", Points: 80},
		{StudentName: "Alice", Subject: "music", Points: 60},
		{StudentName: "Bob", Subject: "music", Points: 100},
	}).Error)

	caseDB = db
	m.Run()
}

func TestPluck(t *testing.T) {
	// The value type comes from the column: cls.StudentName is a string column, so the result is []string
	// 值的类型由列决定：cls.StudentName 是 string 列，结果就是 []string
	names, err := scoreRepo.With(t.Context(), caseDB).Pluck(func(db *gorm.DB, cls *models.ScoreColumns) *gorm.DB {
		return db.Where(cls.Subject.Eq("music")).Order(cls.StudentName.Asc())
	}, func(cls *models.ScoreColumns) gormcnm.ColumnName[string] {
		return cls.StudentName
	})
	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, names)
}

func TestScan(t *testing.T) {
	// Scan works with any struct, even a struct defined in the function
	// Scan 用任意结构体都行，函数里临时定义的也可以
	type nameAndPoints struct {
		StudentName string
		Points      int
	}

	rows, err := scoreRepo.With(t.Context(), caseDB).Scan[nameAndPoints](func(db *gorm.DB, cls *models.ScoreColumns) *gorm.DB {
		return db.Select(cls.StudentName.Name(), cls.Points.Name()).
			Where(cls.Subject.Eq("math")).
			Order(cls.Points.Desc())
	})
	require.NoError(t, err)
	require.Equal(t, []nameAndPoints{
		{StudentName: "Alice", Points: 90},
		{StudentName: "Chloe", Points: 80},
		{StudentName: "Bob", Points: 70},
	}, rows)
}

func TestScanOne(t *testing.T) {
	// One value, it fails on zero rows and on two rows and up, so a missing value is not a silent zero
	// 只取一个值，查不到行或查出多行都报错，查不到的值不会悄悄变成零
	high, err := scoreRepo.With(t.Context(), caseDB).ScanOne[int](func(db *gorm.DB, cls *models.ScoreColumns) *gorm.DB {
		return db.Select(cls.Points.COALESCE().Max()).Where(cls.Subject.Eq("math"))
	})
	require.NoError(t, err)
	require.Equal(t, 90, high)
}

func TestScanList(t *testing.T) {
	// The stats of each subject, each alias is a column of the result struct
	// 每门科目的统计，每个别名都是结果结构体的列
	stats, err := scoreRepo.With(t.Context(), caseDB).ScanList[models.SubjectStat](func(db *gorm.DB, cls *models.ScoreColumns, res *models.SubjectStatColumns) *gorm.DB {
		return db.Select(res.Merge(
			cls.Subject.AsName(res.Subject),
			res.CountAs(res.ScoreCount.Name()),
			cls.Points.COALESCE().SumAs(res.PointsSum.Name()),
			cls.Points.COALESCE().MaxAs(res.PointsHigh.Name()),
			cls.Points.COALESCE().AvgAs(res.PointsAvg.Name()),
		)).Group(cls.Subject.Name()).Order(res.PointsSum.Desc())
	})
	require.NoError(t, err)
	require.Equal(t, []models.SubjectStat{
		{Subject: "math", ScoreCount: 3, PointsSum: 240, PointsHigh: 90, PointsAvg: 80},
		{Subject: "music", ScoreCount: 2, PointsSum: 160, PointsHigh: 100, PointsAvg: 80},
	}, stats)
}

func TestScanItem(t *testing.T) {
	// One row of the result struct, it fails on zero rows and on two rows and up
	// 只要一行结果，查不到行或查出多行都报错
	stat, err := scoreRepo.With(t.Context(), caseDB).ScanItem[models.SubjectStat](func(db *gorm.DB, cls *models.ScoreColumns, res *models.SubjectStatColumns) *gorm.DB {
		return db.Select(res.Merge(
			cls.Subject.AsName(res.Subject),
			res.CountAs(res.ScoreCount.Name()),
			cls.Points.COALESCE().SumAs(res.PointsSum.Name()),
			cls.Points.COALESCE().MaxAs(res.PointsHigh.Name()),
			cls.Points.COALESCE().AvgAs(res.PointsAvg.Name()),
		)).Where(cls.Subject.Eq("music")).Group(cls.Subject.Name())
	})
	require.NoError(t, err)
	require.Equal(t, models.SubjectStat{Subject: "music", ScoreCount: 2, PointsSum: 160, PointsHigh: 100, PointsAvg: 80}, stat)
}

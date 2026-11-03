package gormjoin_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/gormjoin"
	"github.com/yylego/gormrepo/gormjoin/internal/models"
	"github.com/yylego/gormrepo/gormtablerepo"
	"github.com/yylego/must"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var caseDB *gorm.DB

// The repos are built once and used across the tests, the same as in business code
// repo 建一次，各个测试共用，跟业务代码里的用法一样
var (
	studentRepo = gormtablerepo.NewRepo(gormclass.UseTable(&models.Student{}))
	scoreRepo   = gormtablerepo.NewRepo(gormclass.UseTable(&models.Score{}))
	courseRepo  = gormtablerepo.NewRepo(gormclass.UseTable(&models.Course{}))
)

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.NewString())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	// student-1 has a math score, student-2 has none
	// student-1 有一条数学成绩，student-2 没有成绩
	must.Done(db.AutoMigrate(&models.Student{}, &models.Score{}, &models.Course{}))
	must.Done(db.Create([]*models.Student{{ID: 1, Name: "student-1"}, {ID: 2, Name: "student-2"}}).Error)
	must.Done(db.Create([]*models.Score{{ID: 1, StudentID: 1, Subject: "math", Points: 90}}).Error)
	must.Done(db.Create([]*models.Course{{ID: 1, Subject: "math", CoachName: "coach-a"}}).Error)

	caseDB = db
	m.Run()
}

func TestNew(t *testing.T) {
	join := gormjoin.New(studentRepo, clause.LeftJoin, scoreRepo, func(student *models.StudentColumns, score *models.ScoreColumns) []string {
		// Two conditions and up are joined with AND
		// 两个及以上的条件用 AND 连接
		return []string{score.StudentID.OnEq(student.ID), score.Points.Qs("> 0")}
	})
	require.Equal(t, "LEFT JOIN scores ON scores.student_id=students.id AND scores.points > 0", join.Stmt())
}

func TestLeftJoin(t *testing.T) {
	join := gormjoin.LeftJoin(studentRepo, scoreRepo, func(student *models.StudentColumns, score *models.ScoreColumns) []string {
		return []string{score.StudentID.OnEq(student.ID)}
	})

	t.Run("scan-rows", func(t *testing.T) {
		var rows []models.StudentScore
		require.NoError(t, caseDB.Model(&models.Student{}).Joins(join.Stmt()).
			Select(join.Names[models.StudentScore](func(student *models.StudentColumns, score *models.ScoreColumns, res *models.StudentScoreColumns) []string {
				return []string{
					student.Name.AsName(res.StudentName),
					score.Subject.AsName(res.Subject),
					score.Points.AsName(res.Points),
				}
			})).
			Scopes(studentRepo.Base().NewScope(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
				return db.Order(cls.ID.Asc())
			})).
			Scan(&rows).Error)
		require.Equal(t, []models.StudentScore{
			{StudentName: "student-1", Subject: "math", Points: 90},
			{StudentName: "student-2"}, // LEFT JOIN keeps the student with no score
		}, rows)
	})

	t.Run("count", func(t *testing.T) {
		var count int64
		require.NoError(t, caseDB.Model(&models.Student{}).Joins(join.Stmt()).
			Scopes(scoreRepo.Base().NewScope(func(db *gorm.DB, cls *models.ScoreColumns) *gorm.DB {
				return db.Where(cls.Subject.Eq("math"))
			})).
			Count(&count).Error)
		require.Equal(t, int64(1), count)
	})

	t.Run("where", func(t *testing.T) {
		var rows []models.StudentScore
		require.NoError(t, caseDB.Model(&models.Student{}).Joins(join.Stmt()).
			Select(join.Names[models.StudentScore](func(student *models.StudentColumns, score *models.ScoreColumns, res *models.StudentScoreColumns) []string {
				return []string{
					student.Name.AsName(res.StudentName),
					score.Subject.AsName(res.Subject),
					score.Points.AsName(res.Points),
				}
			})).
			Scopes(
				studentRepo.Base().NewScope(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
					return db.Where(cls.Name.In([]string{"student-1", "student-2"}))
				}),
				scoreRepo.Base().NewScope(func(db *gorm.DB, cls *models.ScoreColumns) *gorm.DB {
					return db.Where(cls.Points.Gte(60))
				}),
			).
			Scan(&rows).Error)
		require.Equal(t, []models.StudentScore{
			{StudentName: "student-1", Subject: "math", Points: 90},
		}, rows)
	})
}

func TestLeftJoin_ThreeTables(t *testing.T) {
	// three tables are two joins of two tables
	// 三张表就是两段两表连接
	studentScore := gormjoin.LeftJoin(studentRepo, scoreRepo, func(student *models.StudentColumns, score *models.ScoreColumns) []string {
		return []string{score.StudentID.OnEq(student.ID)}
	})
	scoreCourse := gormjoin.LeftJoin(scoreRepo, courseRepo, func(score *models.ScoreColumns, course *models.CourseColumns) []string {
		return []string{course.Subject.OnEq(score.Subject)}
	})

	var rows []models.StudentCourse
	require.NoError(t, caseDB.Model(&models.Student{}).
		Joins(studentScore.Stmt()).
		Joins(scoreCourse.Stmt()).
		Select(slices.Concat(
			// each join picks columns of its own two tables, the parts concat into one select
			// 每段连接只选自己两张表的列，几段拼成一个 select
			studentScore.Names[models.StudentCourse](func(student *models.StudentColumns, score *models.ScoreColumns, res *models.StudentCourseColumns) []string {
				return []string{
					student.Name.AsName(res.StudentName),
					score.Subject.AsName(res.Subject),
					score.Points.AsName(res.Points),
				}
			}),
			scoreCourse.Names[models.StudentCourse](func(score *models.ScoreColumns, course *models.CourseColumns, res *models.StudentCourseColumns) []string {
				return []string{
					course.CoachName.AsName(res.CoachName),
				}
			}),
		)).
		Scopes(studentRepo.Base().NewScope(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
			return db.Order(cls.ID.Asc())
		})).
		Scan(&rows).Error)
	require.Equal(t, []models.StudentCourse{
		{StudentName: "student-1", Subject: "math", Points: 90, CoachName: "coach-a"},
		{StudentName: "student-2"}, // LEFT JOIN keeps the student with no score
	}, rows)
}

func TestRightJoin(t *testing.T) {
	join := gormjoin.RightJoin(studentRepo, scoreRepo, func(student *models.StudentColumns, score *models.ScoreColumns) []string {
		return []string{score.StudentID.OnEq(student.ID)}
	})
	require.Equal(t, "RIGHT JOIN scores ON scores.student_id=students.id", join.Stmt())
}

func TestInnerJoin(t *testing.T) {
	join := gormjoin.InnerJoin(studentRepo, scoreRepo, func(student *models.StudentColumns, score *models.ScoreColumns) []string {
		return []string{score.StudentID.OnEq(student.ID)}
	})
	require.Equal(t, "INNER JOIN scores ON scores.student_id=students.id", join.Stmt())

	var rows []models.StudentScore
	require.NoError(t, caseDB.Model(&models.Student{}).Joins(join.Stmt()).
		Select(join.Names[models.StudentScore](func(student *models.StudentColumns, score *models.ScoreColumns, res *models.StudentScoreColumns) []string {
			return []string{
				student.Name.AsName(res.StudentName),
				score.Subject.AsName(res.Subject),
				score.Points.AsName(res.Points),
			}
		})).
		Scan(&rows).Error)
	require.Equal(t, []models.StudentScore{
		{StudentName: "student-1", Subject: "math", Points: 90}, // INNER JOIN leaves out the student with no score
	}, rows)
}

func TestJoin_Stmt(t *testing.T) {
	join := gormjoin.LeftJoin(studentRepo, scoreRepo, func(student *models.StudentColumns, score *models.ScoreColumns) []string {
		return []string{score.StudentID.OnEq(student.ID)}
	})
	require.Equal(t, "LEFT JOIN scores ON scores.student_id=students.id", join.Stmt())
}

func TestJoin_Names(t *testing.T) {
	join := gormjoin.LeftJoin(studentRepo, scoreRepo, func(student *models.StudentColumns, score *models.ScoreColumns) []string {
		return []string{score.StudentID.OnEq(student.ID)}
	})
	names := join.Names[models.StudentScore](func(student *models.StudentColumns, score *models.ScoreColumns, res *models.StudentScoreColumns) []string {
		return []string{
			student.Name.AsName(res.StudentName),
			score.Points.AsName(res.Points),
		}
	})
	require.Equal(t, []string{"students.name as student_name", "scores.points as points"}, names)
}

// Package example09_test shows joining tables with gormjoin
// gormjoin.LeftJoin takes two table repos and the ON condition, Stmt gives the JOIN statement to pass to db.Joins
// Names gives the select columns of the two tables, aliased to the result columns with type checks
// The where and the sorting go through each repo's own scope, with *gorm.DB as the link
//
// example09_test 演示用 gormjoin 联表查询
// gormjoin.LeftJoin 传两个表 repo 和连接条件，Stmt 给出传给 db.Joins 的 JOIN 语句
// Names 给出两张表要选的列，别名用结果的列，带类型检查
// 条件和排序走各 repo 自己的 scope，以 *gorm.DB 为纽带
package example09_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/gormjoin"
	"github.com/yylego/gormrepo/gormtablerepo"
	"github.com/yylego/gormrepo/internal/examples/example09/internal/models"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var caseDB *gorm.DB

// Table repos give the columns with the table name, e.g. students.id, which a join needs
// 表 repo 给出带表名的列，例如 students.id，联表需要这种列
var (
	studentRepo = gormtablerepo.NewRepo(gormclass.UseTable(&models.Student{}))
	scoreRepo   = gormtablerepo.NewRepo(gormclass.UseTable(&models.Score{}))
	courseRepo  = gormtablerepo.NewRepo(gormclass.UseTable(&models.Course{}))
)

// The join is built once and used everywhere, the same as the repos
// join 建一次到处用，跟 repo 一样
var studentScore = gormjoin.LeftJoin(studentRepo, scoreRepo, func(student *models.StudentColumns, score *models.ScoreColumns) []string {
	return []string{score.StudentID.OnEq(student.ID)}
})

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	defer rese.F0(rese.P1(db.DB()).Close)

	// Alice has two scores, Bob has one, Chloe has none
	// Alice 有两条成绩，Bob 有一条，Chloe 没有成绩
	done.Done(db.AutoMigrate(&models.Student{}, &models.Score{}, &models.Course{}))
	done.Done(db.Create([]*models.Student{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}, {ID: 3, Name: "Chloe"}}).Error)
	done.Done(db.Create([]*models.Score{
		{ID: 1, StudentID: 1, Subject: "math", Points: 90},
		{ID: 2, StudentID: 1, Subject: "music", Points: 60},
		{ID: 3, StudentID: 2, Subject: "math", Points: 70},
	}).Error)
	done.Done(db.Create([]*models.Course{
		{ID: 1, Subject: "math", CoachName: "Coach Wang"},
		{ID: 2, Subject: "music", CoachName: "Coach Li"},
	}).Error)

	caseDB = db
	m.Run()
}

// selectStudentScore picks the columns of the two tables into StudentScore
// selectStudentScore 把两张表的列选进 StudentScore
func selectStudentScore(student *models.StudentColumns, score *models.ScoreColumns, res *models.StudentScoreColumns) []string {
	return []string{
		student.Name.AsName(res.StudentName),
		score.Subject.AsName(res.Subject),
		score.Points.AsName(res.Points),
	}
}

// sortStudentScore sorts on the student, then on the score
// sortStudentScore 先按学生排，再按成绩排
var sortStudentScore = []func(db *gorm.DB) *gorm.DB{
	studentRepo.Base().NewScope(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
		return db.Order(cls.ID.Asc())
	}),
	scoreRepo.Base().NewScope(func(db *gorm.DB, cls *models.ScoreColumns) *gorm.DB {
		return db.Order(cls.ID.Asc())
	}),
}

func TestStmt(t *testing.T) {
	require.Equal(t, "LEFT JOIN scores ON scores.student_id=students.id", studentScore.Stmt())
}

func TestJoinTwoTables(t *testing.T) {
	var rows []models.StudentScore
	require.NoError(t, caseDB.Model(&models.Student{}).
		Joins(studentScore.Stmt()).
		Select(studentScore.Names[models.StudentScore](selectStudentScore)).
		Scopes(sortStudentScore...).
		Scan(&rows).Error)
	require.Equal(t, []models.StudentScore{
		{StudentName: "Alice", Subject: "math", Points: 90},
		{StudentName: "Alice", Subject: "music", Points: 60},
		{StudentName: "Bob", Subject: "math", Points: 70},
		{StudentName: "Chloe"}, // LEFT JOIN keeps the student with no score
	}, rows)
}

func TestJoinWhere(t *testing.T) {
	// Each repo writes the where of its own table
	// 每个 repo 写自己那张表的条件
	var rows []models.StudentScore
	require.NoError(t, caseDB.Model(&models.Student{}).
		Joins(studentScore.Stmt()).
		Select(studentScore.Names[models.StudentScore](selectStudentScore)).
		Scopes(
			studentRepo.Base().NewScope(func(db *gorm.DB, cls *models.StudentColumns) *gorm.DB {
				return db.Where(cls.Name.In([]string{"Alice", "Bob"}))
			}),
			scoreRepo.Base().NewScope(func(db *gorm.DB, cls *models.ScoreColumns) *gorm.DB {
				return db.Where(cls.Subject.Eq("math"))
			}),
		).
		Scopes(sortStudentScore...).
		Scan(&rows).Error)
	require.Equal(t, []models.StudentScore{
		{StudentName: "Alice", Subject: "math", Points: 90},
		{StudentName: "Bob", Subject: "math", Points: 70},
	}, rows)
}

func TestJoinCount(t *testing.T) {
	var count int64
	require.NoError(t, caseDB.Model(&models.Student{}).
		Joins(studentScore.Stmt()).
		Scopes(scoreRepo.Base().NewScope(func(db *gorm.DB, cls *models.ScoreColumns) *gorm.DB {
			return db.Where(cls.Points.Gte(70))
		})).
		Count(&count).Error)
	require.Equal(t, int64(2), count)
}

func TestJoinThreeTables(t *testing.T) {
	// Three tables are two joins of two tables, each join picks the columns of its own two tables
	// 三张表就是两段两表连接，每段只选自己两张表的列
	scoreCourse := gormjoin.LeftJoin(scoreRepo, courseRepo, func(score *models.ScoreColumns, course *models.CourseColumns) []string {
		return []string{course.Subject.OnEq(score.Subject)}
	})

	var rows []models.StudentCourse
	require.NoError(t, caseDB.Model(&models.Student{}).
		Joins(studentScore.Stmt()).
		Joins(scoreCourse.Stmt()).
		Select(slices.Concat(
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
		Scopes(sortStudentScore...).
		Scan(&rows).Error)
	require.Equal(t, []models.StudentCourse{
		{StudentName: "Alice", Subject: "math", Points: 90, CoachName: "Coach Wang"},
		{StudentName: "Alice", Subject: "music", Points: 60, CoachName: "Coach Li"},
		{StudentName: "Bob", Subject: "math", Points: 70, CoachName: "Coach Wang"},
		{StudentName: "Chloe"},
	}, rows)
}

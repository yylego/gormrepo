// Package example08_test shows using the repos in a GORM transaction
// Pass the tx to With, then each repo works in the transaction, and an error rolls back each change
//
// example08_test 演示在 GORM 事务里使用 repo
// 把 tx 传给 With，各个 repo 就在事务里工作，出错时所有改动一起回滚
package example08_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/yylego/done"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/internal/examples/example08/internal/models"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	courseRepo     = gormrepo.NewRepo(gormclass.Use(&models.Course{}))
	enrollmentRepo = gormrepo.NewRepo(gormclass.Use(&models.Enrollment{}))
)

// newCaseDB opens a new database in each test, with one course that takes 1 student
// newCaseDB 给每个测试开一个新库，里面有一门容量为 1 的课
func newCaseDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:db-%s?mode=memory&cache=shared", uuid.New().String())
	db := rese.P1(gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}))
	t.Cleanup(func() { rese.F0(rese.P1(db.DB()).Close) })

	done.Done(db.AutoMigrate(&models.Course{}, &models.Enrollment{}))
	done.Done(db.Create(&models.Course{Name: "math", Seats: 1}).Error)
	return db
}

// admitStudent writes an enrollment and adds 1 to the enrolled count in one transaction, it fails when the course has no seat left
// admitStudent 在一个事务里写选课记录并给已选人数加 1，课程满员时失败
func admitStudent(ctx context.Context, db *gorm.DB, courseName string, studentName string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		course, err := courseRepo.With(ctx, tx).First(func(db *gorm.DB, cls *models.CourseColumns) *gorm.DB {
			return db.Where(cls.Name.Eq(courseName))
		})
		if err != nil {
			return err
		}

		if err := enrollmentRepo.With(ctx, tx).Create(&models.Enrollment{CourseID: course.ID, StudentName: studentName}); err != nil {
			return err
		}

		if err := courseRepo.With(ctx, tx).Update(func(db *gorm.DB, cls *models.CourseColumns) *gorm.DB {
			return db.Where(cls.ID.Eq(course.ID))
		}, func(cls *models.CourseColumns) (string, any) {
			return cls.Enrolled.KeAdd(1)
		}); err != nil {
			return err
		}

		// The check comes behind the writes on purpose, to show the rollback undoes them
		// 故意在写完之后才检查，用来演示回滚会撤销前面的写入
		if course.Enrolled+1 > course.Seats {
			return errors.Errorf("course %s has no seat left", courseName)
		}
		return nil
	})
}

func TestTransactionCommit(t *testing.T) {
	db := newCaseDB(t)
	require.NoError(t, admitStudent(t.Context(), db, "math", "Alice"))

	course, err := courseRepo.With(t.Context(), db).First(func(db *gorm.DB, cls *models.CourseColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("math"))
	})
	require.NoError(t, err)
	require.Equal(t, 1, course.Enrolled)
}

func TestTransactionRollback(t *testing.T) {
	db := newCaseDB(t)
	require.NoError(t, admitStudent(t.Context(), db, "math", "Alice"))

	// The course has no seat left now, the enrollment of Bob and the added count are both rolled back
	// 课程已满，Bob 的选课记录和加上的人数都被回滚
	require.ErrorContains(t, admitStudent(t.Context(), db, "math", "Bob"), "has no seat left")

	course, err := courseRepo.With(t.Context(), db).First(func(db *gorm.DB, cls *models.CourseColumns) *gorm.DB {
		return db.Where(cls.Name.Eq("math"))
	})
	require.NoError(t, err)
	require.Equal(t, 1, course.Enrolled)

	exist, err := enrollmentRepo.With(t.Context(), db).Exist(func(db *gorm.DB, cls *models.EnrollmentColumns) *gorm.DB {
		return db.Where(cls.StudentName.Eq("Bob"))
	})
	require.NoError(t, err)
	require.False(t, exist)
}

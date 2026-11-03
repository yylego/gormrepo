package gormrepo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/internal/models"
	"github.com/yylego/gormrepo/internal/tests"
	"github.com/yylego/rese"
	"gorm.io/gorm"
)

// TestNewRepo tests NewRepo, which gives a Repo, the generic type alias of BaseRepo
// TestNewRepo 测试 NewRepo，它返回 Repo，即 BaseRepo 的泛型类型别名
func TestNewRepo(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewRepo(gormclass.Use(&models.Account{}))

	ctx := t.Context()

	t.Run("case-1", func(t *testing.T) {
		res, err := repo.With(ctx, db).First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-1-username"))
		})
		require.NoError(t, err)
		require.Equal(t, "demo-1-nickname", res.Nickname)
	})

	t.Run("case-2", func(t *testing.T) {
		res, err := repo.With(ctx, db).First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-2-username"))
		})
		require.NoError(t, err)
		require.Equal(t, "demo-2-nickname", res.Nickname)
	})
}

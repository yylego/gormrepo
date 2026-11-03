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

// TestNewBaseRepo tests NewBaseRepo, which holds the columns and takes the db at each use
// TestNewBaseRepo 测试 NewBaseRepo，它持有列定义，每次使用时再传入数据库连接
func TestNewBaseRepo(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{}))
	account, err := repo.With(t.Context(), db).First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq("demo-1-username"))
	})
	require.NoError(t, err)
	require.Equal(t, "demo-1-nickname", account.Nickname)
}

// TestBaseRepo_Repo tests BaseRepo.Repo method to create GormRepo with db connection
// TestBaseRepo_Repo 测试 BaseRepo.Repo 方法创建带数据库连接的 GormRepo
func TestBaseRepo_Repo(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{}))

	t.Run("demo1", func(t *testing.T) {
		res, err := repo.Repo(db).First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-1-username"))
		})
		require.NoError(t, err)
		require.Equal(t, "demo-1-nickname", res.Nickname)
	})

	t.Run("demo2", func(t *testing.T) {
		require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
			res, err := repo.Repo(tx).First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
				return db.Where(cls.Username.Eq("demo-2-username"))
			})
			require.NoError(t, err)
			require.Equal(t, "demo-2-nickname", res.Nickname)
			return nil
		}))
	})
}

// TestBaseRepo_Gorm tests BaseRepo.Gorm method to create GormWrap with db connection
// TestBaseRepo_Gorm 测试 BaseRepo.Gorm 方法创建带数据库连接的 GormWrap
func TestBaseRepo_Gorm(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{}))

	t.Run("demo1", func(t *testing.T) {
		var account models.Account
		require.NoError(t, repo.Gorm(db).First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-1-username"))
		}, &account).Error)
		require.Equal(t, "demo-1-nickname", account.Nickname)
	})

	t.Run("demo2", func(t *testing.T) {
		require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
			var account models.Account
			require.NoError(t, repo.Gorm(tx).First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
				return db.Where(cls.Username.Eq("demo-2-username"))
			}, &account).Error)
			require.Equal(t, "demo-2-nickname", account.Nickname)
			return nil
		}))
	})
}

// TestBaseRepo_With tests BaseRepo.With method to create GormRepo with context
// TestBaseRepo_With 测试 BaseRepo.With 方法创建带上下文的 GormRepo
func TestBaseRepo_With(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{}))

	ctx := t.Context()
	res, err := repo.With(ctx, db).First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq("demo-1-username"))
	})
	require.NoError(t, err)
	require.Equal(t, "demo-1-nickname", res.Nickname)
}

// TestBaseRepo_Wrap tests BaseRepo.Wrap method to create GormWrap with context
// TestBaseRepo_Wrap 测试 BaseRepo.Wrap 方法创建带上下文的 GormWrap
func TestBaseRepo_Wrap(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{}))

	ctx := t.Context()
	var account models.Account
	require.NoError(t, repo.Wrap(ctx, db).First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq("demo-1-username"))
	}, &account).Error)
	require.Equal(t, "demo-1-nickname", account.Nickname)
}

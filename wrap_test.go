package gormrepo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/internal/models"
	"github.com/yylego/gormrepo/internal/tests"
	"github.com/yylego/rese"
	"gorm.io/gorm"
)

// TestGormRepo_Gorm tests converting GormRepo to GormWrap
// TestGormRepo_Gorm 测试将 GormRepo 转换为 GormWrap
func TestGormRepo_Gorm(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{}))

	var account models.Account
	require.NoError(t, repo.Repo(db).Gorm().First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq("demo-1-username"))
	}, &account).Error)
	require.Equal(t, "demo-1-nickname", account.Nickname)
}

// TestGormWrap_Repo tests converting GormWrap to GormRepo
// TestGormWrap_Repo 测试将 GormWrap 转换为 GormRepo
func TestGormWrap_Repo(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{}))

	res, err := repo.Gorm(db).Repo().First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq("demo-1-username"))
	})
	require.NoError(t, err)
	require.Equal(t, "demo-1-nickname", res.Nickname)
}

// TestGormRepo_Modeled tests Modeled, which applies db.Model(new(MOD)) on GormRepo
// TestGormRepo_Modeled 测试在 GormRepo 上把 MOD 设为 GORM 模型
func TestGormRepo_Modeled(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{}))

	t.Run("case-1", func(t *testing.T) {
		var nickname string
		require.NoError(t, repo.Repo(db).Modeled().Invoke(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-1-username")).Select(string(cls.Nickname)).First(&nickname)
		}))
		require.Equal(t, "demo-1-nickname", nickname)
	})

	t.Run("case-2", func(t *testing.T) {
		require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
			var nickname string
			require.NoError(t, repo.Repo(tx).Modeled().Invoke(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
				return db.Where(cls.Username.Eq("demo-2-username")).Select(string(cls.Nickname)).First(&nickname)
			}))
			require.Equal(t, "demo-2-nickname", nickname)
			return nil
		}))
	})
}

// TestGormWrap_Modeled tests Modeled, which applies db.Model(new(MOD)) on GormWrap
// TestGormWrap_Modeled 测试在 GormWrap 上把 MOD 设为 GORM 模型
func TestGormWrap_Modeled(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{}))

	t.Run("case-1", func(t *testing.T) {
		var nickname string
		require.NoError(t, repo.Gorm(db).Modeled().Invoke(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-1-username")).Select(string(cls.Nickname)).First(&nickname)
		}).Error)
		require.Equal(t, "demo-1-nickname", nickname)
	})

	t.Run("case-2", func(t *testing.T) {
		require.NoError(t, db.Transaction(func(tx *gorm.DB) error {
			var nickname string
			require.NoError(t, repo.Gorm(tx).Modeled().Invoke(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
				return db.Where(cls.Username.Eq("demo-2-username")).Select(string(cls.Nickname)).First(&nickname)
			}).Error)
			require.Equal(t, "demo-2-nickname", nickname)
			return nil
		}))
	})
}

// TestGormRepo_WithContext tests setting context on GormRepo
// TestGormRepo_WithContext 测试在 GormRepo 上设置上下文
func TestGormRepo_WithContext(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{}))
	type ctxSlotType struct{}
	ctx := context.WithValue(t.Context(), ctxSlotType{}, "value-abc")

	res, err := repo.Repo(db).WithContext(ctx).First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		value := ctx.Value(ctxSlotType{})
		require.NotNil(t, value)
		rawValue, ok := value.(string)
		require.True(t, ok)
		require.Equal(t, "value-abc", rawValue)

		return db.Where(cls.Username.Eq("demo-1-username"))
	})
	require.NoError(t, err)
	require.Equal(t, "demo-1-nickname", res.Nickname)
}

// TestGormWrap_WithContext tests setting context on GormWrap
// TestGormWrap_WithContext 测试在 GormWrap 上设置上下文
func TestGormWrap_WithContext(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{}))
	type ctxSlotType struct{}
	ctx := context.WithValue(t.Context(), ctxSlotType{}, "value-abc")

	var account models.Account
	require.NoError(t, repo.Gorm(db).WithContext(ctx).First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		value := ctx.Value(ctxSlotType{})
		require.NotNil(t, value)
		rawValue, ok := value.(string)
		require.True(t, ok)
		require.Equal(t, "value-abc", rawValue)

		return db.Where(cls.Username.Eq("demo-1-username"))
	}, &account).Error)
	require.Equal(t, "demo-1-nickname", account.Nickname)
}

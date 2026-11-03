package gormrepo_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/internal/models"
	"github.com/yylego/gormrepo/internal/tests"
	"gorm.io/gorm"
)

// newAccount creates a new Account instance with the given username
// newAccount 使用给定的用户名创建新的 Account 实例
func newAccount(username string) *models.Account {
	return &models.Account{
		Model:    gorm.Model{},
		Username: username,
		Password: uuid.New().String(),
		Nickname: uuid.New().String(),
	}
}

// TestUse tests the Use function to get db, model, and columns
// TestUse 测试 Use 函数获取 db、model 和 columns
func TestUse(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))
		require.NotNil(t, repo)
	})
}

// TestUmc tests the Umc function (alias to Use)
// TestUmc 测试 Umc 函数（Use 的别名）
func TestUmc(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		repo := gormrepo.NewGormRepo(gormrepo.Umc(db, &models.Account{}))
		require.NotNil(t, repo)
	})
}

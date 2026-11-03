package gormrepo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/internal/models"
	"github.com/yylego/gormrepo/internal/tests"
	"github.com/yylego/neatjson/neatjsons"
	"github.com/yylego/rese"
	"gorm.io/gorm"
)

// TestBaseRepo_NewScope tests NewScope via matching Accounts via username
// TestBaseRepo_NewScope 测试通过用户名匹配 Account 的 NewScope 方法
func TestBaseRepo_NewScope(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{})) // Init repo with Account // 使用 Account 初始化仓储

	// Create scope to match via username="demo-1-username"
	// 创建 scope 通过 username="demo-1-username" 匹配
	scope := repo.NewScope(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq("demo-1-username"))
	})

	// Query first matching Account
	// 查询第一个匹配的 Account
	var account models.Account
	err := db.Scopes(scope).First(&account).Error
	require.NoError(t, err) // Check no error // 检查无错误

	t.Log(neatjsons.S(account)) // Log result // 打印结果

	require.Equal(t, "demo-1-nickname", account.Nickname) // Check nickname // 检查昵称
}

// TestBaseRepo_NewOrderingScope tests NewOrderingScope via sorting Accounts on username (desc)
// TestBaseRepo_NewOrderingScope 测试按 username 倒序排列 Account 的 NewOrderingScope 方法
func TestBaseRepo_NewOrderingScope(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{})) // Init repo with Account // 使用 Account 初始化仓储

	orderingScope := repo.NewOrderingScope(func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.Username.Ob("desc")
	})

	var accounts []*models.Account
	err := db.Scopes(orderingScope).Find(&accounts).Error
	require.NoError(t, err) // Check no error // 检查无错误

	t.Log(neatjsons.S(accounts)) // Log results // 打印结果

	require.Len(t, accounts, 2)
	require.Equal(t, "demo-2-username", accounts[0].Username) // Created second, comes first // 后创建的排在前面
	require.Equal(t, "demo-1-username", accounts[1].Username)
}

// TestPagination_Scope tests Pagination.Scope, which adds LIMIT and OFFSET to a find
// TestPagination_Scope 测试 Pagination.Scope，它给查询加上 LIMIT 和 OFFSET
func TestPagination_Scope(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	cls := gormclass.Cls(&models.Account{})
	page := &gormrepo.Pagination{Limit: 1, Offset: 1}

	var accounts []*models.Account
	err := db.Order(cls.ID.Asc()).Scopes(page.Scope()).Find(&accounts).Error
	require.NoError(t, err) // Check no error // 检查无错误

	require.Len(t, accounts, 1)
	require.Equal(t, "demo-2-username", accounts[0].Username) // Offset 1 skips demo-1 // Offset 1 跳过 demo-1
}

// TestBaseRepo_NewPaginateScope tests NewPaginateScope with sorting and pagination
// TestBaseRepo_NewPaginateScope 测试带排序和分页的 NewPaginateScope 方法
func TestBaseRepo_NewPaginateScope(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewBaseRepo(gormclass.Use(&models.Account{})) // Init repo with Account // 使用 Account 初始化仓储

	// Create scope to match via username in ("demo-1-username", "demo-2-username")
	// 创建 scope 通过 username 在 ("demo-1-username", "demo-2-username") 中匹配
	condScope := repo.NewScope(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"}))
	})

	// Create scope to sort via username (desc) nickname (asc) with limit=10, offset=0
	// 创建 scope 按 username (desc) nickname (asc) 排序，limit=10, offset=0
	pageScope := repo.NewPaginateScope(func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.Username.Ob("desc").Ob(cls.Nickname.Ob("asc"))
	}, &gormrepo.Pagination{
		Limit:  10,
		Offset: 0,
	})

	// Query Accounts with pagination
	// 带分页查询 Accounts
	var accounts []*models.Account
	err := db.Scopes(condScope, pageScope).Find(&accounts).Error
	require.NoError(t, err) // Check no error // 检查无错误

	t.Log(neatjsons.S(accounts)) // Log results // 打印结果

	require.Len(t, accounts, 2)
	require.Equal(t, "demo-2-username", accounts[0].Username) // Check username // 检查用户名
	require.Equal(t, "demo-1-username", accounts[1].Username) // Check username // 检查用户名
}

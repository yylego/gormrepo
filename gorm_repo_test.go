package gormrepo_test

import (
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/internal/models"
	"github.com/yylego/gormrepo/internal/tests"
	"github.com/yylego/must"
	"github.com/yylego/neatjson/neatjsons"
	"github.com/yylego/rese"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// createDemoAccounts creates the accounts table with two demo accounts
// createDemoAccounts 建 accounts 表并插入两个演示账号
func createDemoAccounts(t *testing.T, db *gorm.DB) {
	require.NoError(t, db.AutoMigrate(&models.Account{}))
	require.NoError(t, db.Create([]*models.Account{
		{Username: "demo-1-username", Password: "demo-1-password", Nickname: "demo-1-nickname"},
		{Username: "demo-2-username", Password: "demo-2-password", Nickname: "demo-2-nickname"},
	}).Error)
}

// TestNewGormRepo tests NewGormRepo, which binds the db and the columns in one repo
// TestNewGormRepo 测试 NewGormRepo，它把数据库连接和列定义绑在一个 repo 里
func TestNewGormRepo(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))
	count, err := repo.Count(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), count)
}

// TestGormRepo_First tests the First method to find a single record
// TestGormRepo_First 测试 First 方法查找单条记录
func TestGormRepo_First(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Umc(db, &models.Account{}))

	t.Run("case-1", func(t *testing.T) {
		res, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-1-username"))
		})
		require.NoError(t, err)
		require.Equal(t, "demo-1-nickname", res.Nickname)
	})

	t.Run("case-2", func(t *testing.T) {
		res, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-2-username"))
		})
		require.NoError(t, err)
		require.Equal(t, "demo-2-nickname", res.Nickname)
	})
}

// TestGormRepo_FirstE tests the FirstE method with structured error handling
// TestGormRepo_FirstE 测试带结构化错误处理的 FirstE 方法
func TestGormRepo_FirstE(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Umc(db, &models.Account{}))

	t.Run("case-1", func(t *testing.T) {
		res, erb := repo.FirstE(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-1-username"))
		})
		require.Nil(t, erb)
		require.Equal(t, "demo-1-nickname", res.Nickname)
	})

	t.Run("case-2", func(t *testing.T) {
		res, erb := repo.FirstE(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-x-username"))
		})
		require.NotNil(t, erb)
		require.ErrorIs(t, erb.Cause, gorm.ErrRecordNotFound)
		require.True(t, erb.NotExist)
		require.Nil(t, res)
	})
}

// TestGormRepo_Where tests the Where method to apply custom conditions
// TestGormRepo_Where 测试 Where 方法应用自定义条件
func TestGormRepo_Where(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Umc(db, &models.Account{}))

	t.Run("case-1", func(t *testing.T) {
		var nicknames []string
		dbx := repo.Where(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Model(&models.Account{}).
				Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"})).
				Distinct(cls.Nickname.Name())
		})
		require.NoError(t, dbx.Find(&nicknames).Error)
		t.Log(nicknames)
		require.Len(t, nicknames, 2)
		slices.Sort(nicknames)
		require.Equal(t, []string{"demo-1-nickname", "demo-2-nickname"}, nicknames)
	})

	t.Run("case-2", func(t *testing.T) {
		var nickname string
		dbx := repo.Where(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Model(&models.Account{}).
				Select("MIN(" + cls.Nickname.Name() + ")").
				Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"}))
		})
		require.NoError(t, dbx.First(&nickname).Error)
		t.Log(nickname)
		require.Equal(t, "demo-1-nickname", nickname)
	})
}

// TestGormRepo_Exist tests the Exist method to check record existence
// TestGormRepo_Exist 测试 Exist 方法检查记录是否存在
func TestGormRepo_Exist(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Umc(db, &models.Account{}))

	t.Run("case-1", func(t *testing.T) {
		exist, err := repo.Exist(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-1-username"))
		})
		require.NoError(t, err)
		require.True(t, exist)
	})

	t.Run("case-2", func(t *testing.T) {
		exist, err := repo.Exist(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq("demo-x-username"))
		})
		require.NoError(t, err)
		require.False(t, exist)
	})
}

// TestGormRepo_Find tests the Find method to retrieve multiple records
// TestGormRepo_Find 测试 Find 方法检索多条记录
func TestGormRepo_Find(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	accounts, err := repo.Find(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Like("demo-%-username")).Order(cls.Username.Asc())
	})
	require.NoError(t, err)
	require.Len(t, accounts, 2)
	require.Equal(t, "demo-1-username", accounts[0].Username)
	require.Equal(t, "demo-2-username", accounts[1].Username)
	t.Log(neatjsons.S(accounts))
}

// TestGormRepo_FindAll tests that FindAll returns each record without a where function
// TestGormRepo_FindAll 测试 FindAll 不传 where 函数就能取回全部记录
func TestGormRepo_FindAll(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	accounts, err := repo.FindAll()
	require.NoError(t, err)
	require.Len(t, accounts, 2)
}

// TestGormRepo_FindO tests the FindO method with where condition and ordering
// TestGormRepo_FindO 测试带 where 条件和排序的 FindO 方法
func TestGormRepo_FindO(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	accounts, err := repo.FindO(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Like("demo-%-username"))
	}, func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.Username.Ob("desc")
	})
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Equal(t, "demo-2-username", accounts[0].Username)
	require.Equal(t, "demo-1-username", accounts[len(accounts)-1].Username)
	t.Log(neatjsons.S(accounts))
}

// TestGormRepo_FindN tests the FindN method to retrieve limited records
// TestGormRepo_FindN 测试 FindN 方法检索有限数量的记录
func TestGormRepo_FindN(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	accounts, err := repo.FindN(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"}))
	}, func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.Username.Ob("desc")
	}, 1)
	require.NoError(t, err)
	require.Len(t, accounts, 1)
	require.Equal(t, "demo-2-username", accounts[0].Username) // Created second, comes first in desc // 后创建的在倒序里排第一
	t.Log(neatjsons.S(accounts))
}

func TestGormRepo_FindPage(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	accounts, err := repo.FindPage(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"}))
	}, func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.Username.Ob("desc").Ob(cls.Nickname.Ob("asc"))
	}, &gormrepo.Pagination{
		Offset: 0,
		Limit:  2,
	})
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 2)
	require.Equal(t, "demo-2-username", accounts[0].Username)
	require.Equal(t, "demo-2-nickname", accounts[0].Nickname)
	require.Equal(t, "demo-1-username", accounts[1].Username)
	require.Equal(t, "demo-1-nickname", accounts[1].Nickname)
	t.Log(neatjsons.S(accounts))
}

// TestGormRepo_FindPageUnlimited checks Limit -1, which GORM treats as no limit
// TestGormRepo_FindPageUnlimited 验证 Limit 为 -1 时（GORM 视为不限条数）能正常返回全部记录
func TestGormRepo_FindPageUnlimited(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	accounts, err := repo.FindPage(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"}))
	}, func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.Username.Ob("asc")
	}, &gormrepo.Pagination{
		Offset: 0,
		Limit:  -1,
	})
	require.NoError(t, err)
	require.Len(t, accounts, 2)

	accounts, err = repo.FindN(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"}))
	}, func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.Username.Ob("asc")
	}, -1)
	require.NoError(t, err)
	require.Len(t, accounts, 2)
}

// TestGormRepo_FindPageZeroLimit checks a Pagination with Limit left out gets an error, not a silent blank page
// TestGormRepo_FindPageZeroLimit 验证忘了填 Limit 的 Pagination 会报错，而不是悄悄返回空页
func TestGormRepo_FindPageZeroLimit(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	_, err := repo.FindPage(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db
	}, func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.Username.Ob("asc")
	}, &gormrepo.Pagination{Offset: 0})
	require.Error(t, err)

	_, _, err = repo.FindPageAndCount(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db
	}, func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.Username.Ob("asc")
	}, &gormrepo.Pagination{Offset: 0})
	require.Error(t, err)
}

// TestGormRepo_FindPageAndCount tests paginated search with total count
// TestGormRepo_FindPageAndCount 测试带总数的分页搜索
func TestGormRepo_FindPageAndCount(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	accounts, count, err := repo.FindPageAndCount(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"}))
	}, func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.Username.Ob("asc")
	}, &gormrepo.Pagination{
		Offset: 1,
		Limit:  2,
	})
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 1)
	require.Equal(t, "demo-2-username", accounts[0].Username)
	require.Equal(t, "demo-2-nickname", accounts[0].Nickname)
	require.Equal(t, int64(2), count)
	t.Log(neatjsons.S(accounts))
}

// TestGormRepo_FindC tests the FindC method with custom paging and count
// TestGormRepo_FindC 测试带自定义分页和计数的 FindC 方法
func TestGormRepo_FindC(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	accounts, count, err := repo.FindC(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"}))
	}, func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.Username.Ob("asc")
	}, func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Offset(1).Limit(2)
	})
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 1)
	require.Equal(t, "demo-2-username", accounts[0].Username)
	require.Equal(t, "demo-2-nickname", accounts[0].Nickname)
	require.Equal(t, int64(2), count)
	t.Log(neatjsons.S(accounts))
}

// TestGormRepo_FindC_SwappedCallbacks tests FindC gives the same page and count when the two callbacks are swapped
// TestGormRepo_FindC_SwappedCallbacks 测试 FindC 两个回调正着传、反着传，当前页和总数都一样
func TestGormRepo_FindC_SwappedCallbacks(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	matchDemo2 := func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq("demo-2-username"))
	}
	idDesc := func(cls *models.AccountColumns) gormcnm.Ordering {
		return cls.ID.Ob("desc")
	}
	pageAt := func(offset int) func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Offset(offset).Limit(1)
		}
	}

	testCases := []struct {
		name      string
		where     func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB
		paging    func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB
		wantNames []string
	}{
		{name: "correct-page-one", where: matchDemo2, paging: pageAt(0), wantNames: []string{"demo-2-username"}},
		{name: "correct-page-two", where: matchDemo2, paging: pageAt(1), wantNames: []string{}},
		{name: "swapped-page-one", where: pageAt(0), paging: matchDemo2, wantNames: []string{"demo-2-username"}},
		{name: "swapped-page-two", where: pageAt(1), paging: matchDemo2, wantNames: []string{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accounts, count, err := repo.FindC(tc.where, idDesc, tc.paging)
			require.NoError(t, err)
			names := make([]string, 0, len(accounts))
			for _, account := range accounts {
				names = append(names, account.Username)
			}
			require.Equal(t, tc.wantNames, names)
			require.Equal(t, int64(1), count)
		})
	}
}

func TestGormRepo_Count(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	count, err := repo.Count(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"}))
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), count)
}

func TestGormRepo_Pluck(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	usernames, err := repo.Pluck(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"})).Order(cls.Username.Asc())
	}, func(cls *models.AccountColumns) gormcnm.ColumnName[string] {
		return cls.Username
	})
	require.NoError(t, err)
	require.Equal(t, []string{"demo-1-username", "demo-2-username"}, usernames)

	ids, err := repo.Pluck(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq("demo-2-username"))
	}, func(cls *models.AccountColumns) gormcnm.ColumnName[uint] {
		return cls.ID
	})
	require.NoError(t, err)
	require.Len(t, ids, 1)
	require.NotZero(t, ids[0])
}

func TestGormRepo_Scan(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	type accountInfo struct {
		Username string
		Nickname string
	}

	infos, err := repo.Scan[accountInfo](func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Select(cls.Username.Name(), cls.Nickname.Name()).
			Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"})).
			Order(cls.Username.Desc())
	})
	require.NoError(t, err)
	require.Equal(t, []accountInfo{
		{Username: "demo-2-username", Nickname: "demo-2-nickname"},
		{Username: "demo-1-username", Nickname: "demo-1-nickname"},
	}, infos)
}

func TestGormRepo_ScanList(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	infos, err := repo.ScanList[models.AccountInfo](func(db *gorm.DB, cls *models.AccountColumns, res *models.AccountInfoColumns) *gorm.DB {
		return db.Select(res.Merge(
			cls.Username.AsName(res.LoginName),
			cls.Nickname.AsName(res.ShowName),
		)).
			Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"})).
			Order(res.LoginName.Desc())
	})
	require.NoError(t, err)
	require.Equal(t, []models.AccountInfo{
		{LoginName: "demo-2-username", ShowName: "demo-2-nickname"},
		{LoginName: "demo-1-username", ShowName: "demo-1-nickname"},
	}, infos)
}

func TestGormRepo_ScanItem(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	for _, tc := range []struct {
		name      string
		pattern   string
		want      models.AccountInfo
		wantFault bool
	}{
		{name: "one-row", pattern: "demo-2-username", want: models.AccountInfo{LoginName: "demo-2-username", ShowName: "demo-2-nickname"}},
		{name: "no-rows", pattern: "missing-username", wantFault: true},
		{name: "multiple-rows", pattern: "demo-%-username", wantFault: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			info, err := repo.ScanItem[models.AccountInfo](func(db *gorm.DB, cls *models.AccountColumns, res *models.AccountInfoColumns) *gorm.DB {
				return db.Select(res.Merge(
					cls.Username.AsName(res.LoginName),
					cls.Nickname.AsName(res.ShowName),
				)).Where(cls.Username.Like(tc.pattern))
			})
			if tc.wantFault {
				require.Error(t, err)
				require.Equal(t, models.AccountInfo{}, info)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, info)
		})
	}
}

// TestGormRepo_ScanOne tests scanning one aggregate value, and the faults on zero rows and on two rows
// TestGormRepo_ScanOne 测试查一个聚合值，以及查不到行、查出两行时报错
func TestGormRepo_ScanOne(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	t.Run("max-value", func(t *testing.T) {
		nickname, err := repo.ScanOne[string](func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Select(cls.Nickname.COALESCE().Func("MAX", "''"))
		})
		require.NoError(t, err)
		require.Equal(t, "demo-2-nickname", nickname)
	})

	t.Run("max-on-no-match-gives-the-default", func(t *testing.T) {
		nickname, err := repo.ScanOne[string](func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Select(cls.Nickname.COALESCE().Func("MAX", "''")).Where(cls.Username.Eq("missing-username"))
		})
		require.NoError(t, err)
		require.Equal(t, "", nickname)
	})

	t.Run("no-row-is-a-fault", func(t *testing.T) {
		_, err := repo.ScanOne[string](func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Select(cls.Nickname.Name()).Where(cls.Username.Eq("missing-username"))
		})
		require.Error(t, err)
	})

	t.Run("two-rows-is-a-fault", func(t *testing.T) {
		_, err := repo.ScanOne[string](func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Select(cls.Nickname.Name())
		})
		require.Error(t, err)
	})
}

func TestGormRepo_Update(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	must.Done(db.AutoMigrate(&models.Account{}))

	username := uuid.New().String()
	require.NoError(t, db.Save(newAccount(username)).Error)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	newNickname := uuid.New().String()
	require.NoError(t, repo.Update(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq(username))
	}, func(cls *models.AccountColumns) (string, any) {
		return cls.Nickname.Kv(newNickname)
	}))

	res, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq(username))
	})
	require.NoError(t, err)
	require.Equal(t, newNickname, res.Nickname)
}

func TestGormRepo_Updates(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	must.Done(db.AutoMigrate(&models.Account{}))

	username := uuid.New().String()
	require.NoError(t, db.Save(newAccount(username)).Error)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	newNickname := uuid.New().String()
	newPassword := uuid.New().String()
	require.NoError(t, repo.Updates(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq(username))
	}, func(cls *models.AccountColumns) map[string]any {
		return cls.
			Kw(cls.Nickname.Kv(newNickname)).
			Kw(cls.Password.Kv(newPassword)).
			AsMap() // Convert to map[string]any // 转换为 map[string]any 类型
	}))

	res, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq(username))
	})
	require.NoError(t, err)
	require.Equal(t, newNickname, res.Nickname)
	require.Equal(t, newPassword, res.Password)
}

// TestGormRepo_UpdatesM tests fluent API without AsMap() call
// TestGormRepo_UpdatesM 测试流畅的 API，无需调用 AsMap()
func TestGormRepo_UpdatesM(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	must.Done(db.AutoMigrate(&models.Account{}))

	username := uuid.New().String()
	require.NoError(t, db.Save(newAccount(username)).Error)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	newNickname := uuid.New().String()
	newPassword := uuid.New().String()
	require.NoError(t, repo.UpdatesM(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq(username))
	}, func(cls *models.AccountColumns) gormcnm.ColumnValueMap {
		return cls.
			Kw(cls.Nickname.Kv(newNickname)).
			Kw(cls.Password.Kv(newPassword))
		// No AsMap() needed! // 不需要 AsMap()！
	}))

	res, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
		return db.Where(cls.Username.Eq(username))
	})
	require.NoError(t, err)
	require.Equal(t, newNickname, res.Nickname)
	require.Equal(t, newPassword, res.Password)
}

// TestGormRepo_UpdatesO tests update using the object ID as condition
// O = Object, updates the record that has the object's ID
//
// TestGormRepo_UpdatesO 测试使用 object 主键作为条件的更新
// O = Object，通过 object 的主键定位要更新的记录
func TestGormRepo_UpdatesO(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		username := uuid.New().String()
		account := newAccount(username)
		require.NoError(t, repo.Create(account))

		newNickname := uuid.New().String()
		newPassword := uuid.New().String()
		require.NoError(t, repo.UpdatesO(account, func(cls *models.AccountColumns) gormcnm.ColumnValueMap {
			return cls.
				Kw(cls.Nickname.Kv(newNickname)).
				Kw(cls.Password.Kv(newPassword))
		}))

		res, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		})
		require.NoError(t, err)
		require.Equal(t, newNickname, res.Nickname)
		require.Equal(t, newPassword, res.Password)
	})
}

// TestGormRepo_UpdatesC tests update using combined conditions: the object ID plus where clause
// C = Combined, uses both the object ID and where conditions
//
// TestGormRepo_UpdatesC 测试使用组合条件的更新：object 主键加上 where 子句
// C = Combined，同时使用 object 主键和 where 条件进行精确定位
func TestGormRepo_UpdatesC(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		username := uuid.New().String()
		account := newAccount(username)
		require.NoError(t, repo.Create(account))

		newNickname := uuid.New().String()
		newPassword := uuid.New().String()
		require.NoError(t, repo.UpdatesC(account, func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		}, func(cls *models.AccountColumns) gormcnm.ColumnValueMap {
			return cls.
				Kw(cls.Nickname.Kv(newNickname)).
				Kw(cls.Password.Kv(newPassword))
		}))

		res, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		})
		require.NoError(t, err)
		require.Equal(t, newNickname, res.Nickname)
		require.Equal(t, newPassword, res.Password)
	})
}

func TestGormRepo_Invoke(t *testing.T) {
	db := tests.NewMemDB(t)
	defer rese.F0(rese.P1(db.DB()).Close)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Umc(db, &models.Account{}))

	t.Run("case-1", func(t *testing.T) {
		var passwords []string
		require.NoError(t, repo.Invoke(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Model(&models.Account{}).Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"})).Distinct(cls.Password.Name()).Find(&passwords)
		}))
		t.Log(passwords)
		require.Len(t, passwords, 2)
		slices.Sort(passwords)
		require.Equal(t, []string{"demo-1-password", "demo-2-password"}, passwords)
	})

	t.Run("case-2", func(t *testing.T) {
		var nickname string
		require.NoError(t, repo.Invoke(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Model(&models.Account{}).Where(cls.Username.In([]string{"demo-1-username", "demo-2-username"})).Select("MAX(" + cls.Nickname.Name() + ")").First(&nickname)
		}))
		t.Log(nickname)
		require.Equal(t, "demo-2-nickname", nickname)
	})
}

func TestGormRepo_Create(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		username := uuid.New().String()
		account := newAccount(username)
		require.NoError(t, repo.Create(account))
		require.NotZero(t, account.ID)

		res, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		})
		require.NoError(t, err)
		require.Equal(t, username, res.Username)
	})
}

// TestGormRepo_Creates tests batch insert of multiple records
// TestGormRepo_Creates 测试批量插入多条记录
func TestGormRepo_Creates(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		accounts := []*models.Account{
			newAccount(uuid.New().String()),
			newAccount(uuid.New().String()),
			newAccount(uuid.New().String()),
		}
		require.NoError(t, repo.Creates(&accounts))

		// Verify all records have IDs assigned
		// 验证所有记录都分配了 ID
		for _, account := range accounts {
			require.NotZero(t, account.ID)
		}

		// Verify all records exist in database
		// 验证所有记录都存在于数据库中
		count, err := repo.Count(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.In([]string{
				accounts[0].Username,
				accounts[1].Username,
				accounts[2].Username,
			}))
		})
		require.NoError(t, err)
		require.Equal(t, int64(3), count)
	})
}

// TestGormRepo_CreateInBatches tests batch insert with batch size control
// TestGormRepo_CreateInBatches 测试带批次大小控制的分批插入
func TestGormRepo_CreateInBatches(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		// Create 10 accounts with batch size 3
		// 创建 10 个账户，每批 3 条
		accounts := make([]*models.Account, 10)
		for i := range 10 {
			accounts[i] = newAccount(uuid.New().String())
		}
		require.NoError(t, repo.CreateInBatches(accounts, 3))

		// Verify all records have IDs assigned
		// 验证所有记录都分配了 ID
		for _, account := range accounts {
			require.NotZero(t, account.ID)
		}

		// Verify all records exist in database
		// 验证所有记录都存在于数据库中
		usernames := make([]string, len(accounts))
		for i, account := range accounts {
			usernames[i] = account.Username
		}
		count, err := repo.Count(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.In(usernames))
		})
		require.NoError(t, err)
		require.Equal(t, int64(10), count)
	})
}

func TestGormRepo_Save(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		username := uuid.New().String()
		account := newAccount(username)
		require.NoError(t, repo.Save(account))
		require.NotZero(t, account.ID)

		newNickname := uuid.New().String()
		account.Nickname = newNickname
		require.NoError(t, repo.Save(account))

		res, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		})
		require.NoError(t, err)
		require.Equal(t, newNickname, res.Nickname)
	})
}

// TestGormRepo_Clauses_Saves_Returning tests Saves with OnConflict DoNothing and Returning
// The conflicting record gets skipped, the slice becomes the records that got saved
//
// TestGormRepo_Clauses_Saves_Returning 测试 Saves 搭配 OnConflict DoNothing 与 Returning
// 冲突的记录被跳过，切片变成实际保存的那几条记录
func TestGormRepo_Clauses_Saves_Returning(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		existing := newAccount(uuid.New().String())
		require.NoError(t, repo.Create(existing))

		username2 := uuid.New().String()
		accounts := []*models.Account{
			newAccount(existing.Username),
			newAccount(username2),
		}
		require.NoError(t, repo.Clauses(
			clause.OnConflict{DoNothing: true},
			clause.Returning{},
		).Saves(&accounts))

		require.Len(t, accounts, 1)
		require.Equal(t, username2, accounts[0].Username)
		require.NotZero(t, accounts[0].ID)
	})
}

// TestGormRepo_Saves tests batch insert or update of multiple records
// TestGormRepo_Saves 测试批量插入或更新多条记录
func TestGormRepo_Saves(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		// First save - should create new records
		// 首次保存 - 应该创建新记录
		accounts := []*models.Account{
			newAccount(uuid.New().String()),
			newAccount(uuid.New().String()),
		}
		require.NoError(t, repo.Saves(&accounts))

		for _, account := range accounts {
			require.NotZero(t, account.ID)
		}

		// Update nicknames and save again - should update existing records
		// 更新 nickname 并再次保存 - 应该更新现有记录
		newNickname0 := uuid.New().String()
		newNickname1 := uuid.New().String()
		accounts[0].Nickname = newNickname0
		accounts[1].Nickname = newNickname1
		require.NoError(t, repo.Saves(&accounts))

		// Verify updates
		// 验证更新
		res0, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.ID.Eq(accounts[0].ID))
		})
		require.NoError(t, err)
		require.Equal(t, newNickname0, res0.Nickname)

		res1, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.ID.Eq(accounts[1].ID))
		})
		require.NoError(t, err)
		require.Equal(t, newNickname1, res1.Nickname)
	})
}

func TestGormRepo_Delete(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		username := uuid.New().String()
		account := newAccount(username)
		require.NoError(t, repo.Create(account))

		require.NoError(t, repo.Delete(account))

		exist, err := repo.Exist(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.ID.Eq(account.ID))
		})
		require.NoError(t, err)
		require.False(t, exist)
	})
}

func TestGormRepo_DeleteW(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		username := uuid.New().String()
		account := newAccount(username)
		require.NoError(t, repo.Create(account))

		require.NoError(t, repo.DeleteW(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		}))

		exist, err := repo.Exist(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		})
		require.NoError(t, err)
		require.False(t, exist)
	})
}

func TestGormRepo_DeleteM(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		username := uuid.New().String()
		account := newAccount(username)
		require.NoError(t, repo.Create(account))

		require.NoError(t, repo.DeleteM(account, func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		}))

		exist, err := repo.Exist(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.ID.Eq(account.ID))
		})
		require.NoError(t, err)
		require.False(t, exist)
	})
}

func TestGormRepo_Unscoped(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		username := uuid.New().String()
		account := newAccount(username)
		require.NoError(t, repo.Create(account))

		// soft delete
		require.NoError(t, repo.DeleteW(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		}))

		// default scope: not visible
		exist, err := repo.Exist(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		})
		require.NoError(t, err)
		require.False(t, exist)

		// unscoped: visible (soft-deleted)
		exist, err = repo.Unscoped().Exist(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		})
		require.NoError(t, err)
		require.True(t, exist)

		// unscoped hard delete
		require.NoError(t, repo.Unscoped().DeleteW(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		}))

		// unscoped: now gone
		exist, err = repo.Unscoped().Exist(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(username))
		})
		require.NoError(t, err)
		require.False(t, exist)
	})
}

// TestGormRepo_Clauses tests Clauses method for upsert operations
// TestGormRepo_Clauses 测试 Clauses 方法的 upsert 操作
func TestGormRepo_Clauses(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		// First create
		// 首次创建
		account1 := &models.Account{
			Username: uuid.New().String(),
			Password: uuid.New().String(),
			Nickname: uuid.New().String(),
		}
		require.NoError(t, repo.Create(account1))
		require.NotZero(t, account1.ID)

		// Upsert with Clauses - should update nickname
		// 使用 Clauses 进行 upsert - 应该更新 nickname
		account2 := &models.Account{
			Username: account1.Username,
			Password: uuid.New().String(),
			Nickname: uuid.New().String(),
		}
		cls := account2.Columns()
		require.NoError(t, repo.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: cls.Username.Name()}},
			DoUpdates: clause.AssignmentColumns([]string{cls.Nickname.Name()}),
		}).Create(account2))

		// Verify the nickname was updated
		// 验证 nickname 已被更新
		res, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(account1.Username))
		})
		require.NoError(t, err)
		require.Equal(t, account2.Username, res.Username)
		require.Equal(t, account2.Nickname, res.Nickname)
		require.Equal(t, account1.Password, res.Password)
	})
}

// TestGormRepo_Clause tests Clause method for type-safe upsert operations
// TestGormRepo_Clause 测试 Clause 方法的类型安全 upsert 操作
func TestGormRepo_Clause(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		// First create
		// 首次创建
		account1 := &models.Account{
			Username: uuid.New().String(),
			Password: uuid.New().String(),
			Nickname: uuid.New().String(),
		}
		require.NoError(t, repo.Create(account1))
		require.NotZero(t, account1.ID)

		// Upsert with Clause - type-safe column names
		// 使用 Clause 进行 upsert - 类型安全的列名
		account2 := &models.Account{
			Username: account1.Username,
			Password: uuid.New().String(),
			Nickname: uuid.New().String(),
		}
		require.NoError(t, repo.Clause(func(cls *models.AccountColumns) clause.Expression {
			return clause.OnConflict{
				Columns:   []clause.Column{{Name: cls.Username.Name()}},
				DoUpdates: clause.AssignmentColumns([]string{cls.Nickname.Name()}),
			}
		}).Create(account2))

		// Verify the nickname was updated
		// 验证 nickname 已被更新
		res, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(account1.Username))
		})
		require.NoError(t, err)
		require.Equal(t, account2.Username, res.Username)
		require.Equal(t, account2.Nickname, res.Nickname)
		require.Equal(t, account1.Password, res.Password)
	})
}

// TestGormRepo_Clauses_Creates_Returning tests Creates with OnConflict DoNothing and Returning
// The conflicting record gets skipped, the slice becomes the records that got inserted
//
// TestGormRepo_Clauses_Creates_Returning 测试 Creates 搭配 OnConflict DoNothing 与 Returning
// 冲突的记录被跳过，切片变成实际插入的那几条记录
func TestGormRepo_Clauses_Creates_Returning(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		existing := newAccount(uuid.New().String())
		require.NoError(t, repo.Create(existing))

		username2 := uuid.New().String()
		username3 := uuid.New().String()
		accounts := []*models.Account{
			newAccount(existing.Username),
			newAccount(username2),
			newAccount(username3),
		}
		require.NoError(t, repo.Clauses(
			clause.OnConflict{DoNothing: true},
			clause.Returning{},
		).Creates(&accounts))

		require.Len(t, accounts, 2)
		require.Equal(t, username2, accounts[0].Username)
		require.Equal(t, username3, accounts[1].Username)
		require.NotZero(t, accounts[0].ID)
		require.NotZero(t, accounts[1].ID)
	})
}

// TestGormRepo_Clause_Creates tests Clause + Creates for batch upsert operations
// TestGormRepo_Clause_Creates 测试 Clause + Creates 的批量 upsert 操作
func TestGormRepo_Clause_Creates(t *testing.T) {
	tests.NewDBRun(t, func(db *gorm.DB) {
		must.Done(db.AutoMigrate(&models.Account{}))

		repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

		// Prepare nicknames
		// 准备 nicknames
		preNick1 := uuid.New().String()
		preNick2 := uuid.New().String()
		preNick3 := uuid.New().String()

		// First create some accounts
		// 首先创建一些账户
		accounts := []*models.Account{
			{Username: uuid.New().String(), Password: uuid.New().String(), Nickname: preNick1},
			{Username: uuid.New().String(), Password: uuid.New().String(), Nickname: preNick2},
			{Username: uuid.New().String(), Password: uuid.New().String(), Nickname: preNick3},
		}
		require.NoError(t, repo.Creates(&accounts))

		newNick1 := uuid.New().String()
		newNick2 := uuid.New().String()
		newNick4 := uuid.New().String()

		// Batch upsert with Clause - update nicknames
		// 使用 Clause 进行批量 upsert - 更新 nicknames
		upsertAccounts := []*models.Account{
			{Username: accounts[0].Username, Password: uuid.New().String(), Nickname: newNick1},
			{Username: accounts[1].Username, Password: uuid.New().String(), Nickname: newNick2},
			{Username: uuid.New().String(), Password: uuid.New().String(), Nickname: newNick4},
		}
		require.NoError(t, repo.Clause(func(cls *models.AccountColumns) clause.Expression {
			return clause.OnConflict{
				Columns:   []clause.Column{{Name: cls.Username.Name()}},
				DoUpdates: clause.AssignmentColumns([]string{cls.Nickname.Name()}),
			}
		}).Creates(&upsertAccounts))

		// Verify updates
		// 验证更新结果
		res1, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(accounts[0].Username))
		})
		require.NoError(t, err)
		require.Equal(t, newNick1, res1.Nickname)
		require.Equal(t, accounts[0].Password, res1.Password)

		res2, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(accounts[1].Username))
		})
		require.NoError(t, err)
		require.Equal(t, newNick2, res2.Nickname)

		// Verify new insert
		// 验证新插入的记录
		res4, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(upsertAccounts[2].Username))
		})
		require.NoError(t, err)
		require.Equal(t, newNick4, res4.Nickname)

		// Verify original account not affected
		// 验证原始账户未受影响
		res3, err := repo.First(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
			return db.Where(cls.Username.Eq(accounts[2].Username))
		})
		require.NoError(t, err)
		require.Equal(t, preNick3, res3.Nickname)
	})
}

// TestGormRepo_DerivedRepoReuse tests a repo from Modeled / Unscoped / Clauses can be saved and reused,
// each use starts clean and does not take the conditions of the previous use
//
// TestGormRepo_DerivedRepoReuse 测试 Modeled / Unscoped / Clauses 得到的 repo 可以存下来多次调用，
// 每次调用都是干净的，不会带上前一次调用的条件
func TestGormRepo_DerivedRepoReuse(t *testing.T) {
	db := tests.NewMemDB(t)
	createDemoAccounts(t, db)

	repo := gormrepo.NewGormRepo(gormrepo.Use(db, &models.Account{}))

	for name, derived := range map[string]*gormrepo.GormRepo[models.Account, *models.AccountColumns]{
		"modeled":  repo.Modeled(),
		"unscoped": repo.Unscoped(),
		"clauses":  repo.Clauses(clause.Locking{Strength: "UPDATE"}),
	} {
		t.Run(name, func(t *testing.T) {
			for _, username := range []string{"demo-1-username", "demo-2-username"} {
				accounts, err := derived.Find(func(db *gorm.DB, cls *models.AccountColumns) *gorm.DB {
					return db.Where(cls.Username.Eq(username))
				})
				require.NoError(t, err)
				require.Len(t, accounts, 1)
				require.Equal(t, username, accounts[0].Username)
			}
		})
	}
}

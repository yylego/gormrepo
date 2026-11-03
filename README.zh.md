<!-- TEMPLATE (ZH) BEGIN: BADGES -->

[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/gormrepo/release.yml?branch=main&label=BUILD)](https://github.com/yylego/gormrepo/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/gormrepo)](https://pkg.go.dev/github.com/yylego/gormrepo)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/gormrepo/main.svg)](https://coveralls.io/github/yylego/gormrepo?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.27%2B-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/yylego/gormrepo.svg)](https://github.com/yylego/gormrepo/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/gormrepo)](https://goreportcard.com/report/github.com/yylego/gormrepo)

<!-- TEMPLATE (ZH) CLOSE: BADGES -->

# 🚀 gormrepo - 基于 GORM 的类型安全 Repo

**gormrepo** 用泛型 repo 包装 GORM：查询条件、排序、新值都用 gormcngen 生成的列来写，列名或值的类型写错了，编译就过不去。

> 🌟 **回调拿到的是原生的 `*gorm.DB`，回调里写的就是原生 GORM 代码**

---

## 生态系统

![GORM Type-Safe Ecosystem](https://github.com/yylego/gormcnm/raw/main/assets/gormcnm-ecosystem.svg)

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->

## 英文文档

[ENGLISH README](README.md)

<!-- TEMPLATE (ZH) CLOSE: LANGUAGE NAVIGATION -->

---

## 🔄 技术对比

| 生态系统       | Java MyBatis Plus  | Python SQLAlchemy | Go GORM 生态系统 |
| -------------- | ------------------ | ----------------- | ---------------- |
| **类型安全列** | `Example::getName` | `Example.name`    | `cls.Name.Eq()`  |
| **仓储模式**   | ✅ BaseMapper      | ✅ Session API    | ✅ GormRepo      |

---

## 🚀 快速开始

### 安装

```bash
go get github.com/yylego/gormrepo
```

### 完整使用流程

#### 1. 定义模型（支持原生字段）

```go
type Account struct {
    gorm.Model // ID、CreatedAt、UpdatedAt，以及开启软删除的 DeletedAt
    Username string `gorm:"uniqueIndex"`
    Nickname string `gorm:"index"`
    Age      int
}
```

#### 2. 自动生成列结构体（gormcngen）

```go
// 自动生成
func (c *Account) Columns() *AccountColumns {
    return c.TableColumns(gormcnm.NewPlainDecoration())
}

type AccountColumns struct {
    gormcnm.ColumnOperationClass // 嵌入的辅助方法，例如 cls.Kw(...)、cls.Merge(...)
    ID        gormcnm.ColumnName[uint]
    CreatedAt gormcnm.ColumnName[time.Time]
    UpdatedAt gormcnm.ColumnName[time.Time]
    DeletedAt gormcnm.ColumnName[gorm.DeletedAt]
    Username  gormcnm.ColumnName[string]
    Nickname  gormcnm.ColumnName[string]
    Age       gormcnm.ColumnName[int]
}

// TableColumns 按装饰给出列，例如联表时带上表名：accounts.id
func (c *Account) TableColumns(decoration gormcnm.ColumnNameDecoration) *AccountColumns {
    return &AccountColumns{
        ID:        gormcnm.Cmn(c.ID, "id", decoration),
        CreatedAt: gormcnm.Cmn(c.CreatedAt, "created_at", decoration),
        UpdatedAt: gormcnm.Cmn(c.UpdatedAt, "updated_at", decoration),
        DeletedAt: gormcnm.Cmn(c.DeletedAt, "deleted_at", decoration),
        Username:  gormcnm.Cmn(c.Username, "username", decoration),
        Nickname:  gormcnm.Cmn(c.Nickname, "nickname", decoration),
        Age:       gormcnm.Cmn(c.Age, "age", decoration),
    }
}
```

#### 3. 创建 Repo

```go
// 创建 repo，传入列定义
repo := gormrepo.NewRepo(&Account{}, (&Account{}).Columns())

// gormrepo/gormclass 简洁写法
repo := gormrepo.NewRepo(gormclass.Use(&Account{}))
```

#### 4. 查询记录

```go
// 经典 GORM
var alice Account
err := db.Where("username = ?", "alice").First(&alice).Error

// gormrepo - First
account, err := repo.With(ctx, db).First(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Username.Eq("alice"))
})

// gormrepo - FirstE (返回 ErrorOrNotExist)
account, erb := repo.With(ctx, db).FirstE(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Username.Eq("alice"))
})

// gormrepo - Find
accounts, err := repo.With(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18))
})

// gormrepo - FindN（取前 N 条）
accounts, err := repo.With(ctx, db).FindN(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.Age.Gte(18))
    },
    func(cls *AccountColumns) gormcnm.Ordering {
        return cls.Age.Ob("DESC")
    },
    10,
)

// gormrepo - FindAll（查全部记录，不带条件）
accounts, err := repo.With(ctx, db).FindAll()

// gormrepo - FindO（带排序）
accounts, err := repo.With(ctx, db).FindO(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.Age.Gte(18))
    },
    func(cls *AccountColumns) gormcnm.Ordering {
        return cls.ID.Ob("DESC")
    },
)

// gormrepo - FindPage
accounts, err := repo.With(ctx, db).FindPage(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.Age.Gte(18))
    },
    func(cls *AccountColumns) gormcnm.Ordering {
        return cls.ID.Ob("DESC")
    },
    &gormrepo.Pagination{Limit: 10, Offset: 0},
)

// gormrepo - FindPageAndCount
accounts, count, err := repo.With(ctx, db).FindPageAndCount(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.Age.Between(18, 65))
    },
    func(cls *AccountColumns) gormcnm.Ordering {
        return cls.ID.Ob("DESC")
    },
    &gormrepo.Pagination{Limit: 10, Offset: 0},
)

// gormrepo - FindC（分页写成自定义回调，同时返回不分页时的总数）
accounts, count, err := repo.With(ctx, db).FindC(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.Age.Gte(18))
    },
    func(cls *AccountColumns) gormcnm.Ordering {
        return cls.Age.Ob("DESC")
    },
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Limit(10)
    },
)

// gormrepo - Count
count, err := repo.With(ctx, db).Count(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18))
})

// gormrepo - Exist
exist, err := repo.With(ctx, db).Exist(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Username.Eq("alice"))
})

// gormrepo - Pluck（取一列的值，结果类型由列决定：这里是 []string）
names, err := repo.With(ctx, db).Pluck(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18))
}, func(cls *AccountColumns) gormcnm.ColumnName[string] {
    return cls.Username
})

// 结果行，不是表：跟模型放在一起，同样生成列
type AccountInfo struct {
    ID       uint
    Username string
}
type AgeCount struct {
    Age          int
    AccountCount int64
}

// gormrepo - Scan（把结果行扫描到结构体，适合 Select / Joins / Group）
infos, err := repo.With(ctx, db).Scan[AccountInfo](func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Select(cls.ID.Name(), cls.Username.Name()).Where(cls.Age.Gte(18))
})

// gormrepo - ScanList（返回结果列表，回调同时获得结果结构体的列）
ageCounts, err := repo.With(ctx, db).ScanList[AgeCount](func(db *gorm.DB, cls *AccountColumns, res *AgeCountColumns) *gorm.DB {
    return db.Select(res.Merge(cls.Age.Name(), res.CountAs(res.AccountCount.Name()))).
              Group(cls.Age.Name()).
              Order(res.AccountCount.Desc())
})

// gormrepo - ScanItem（返回恰好一条结果，回调同时获得结果结构体的列）
ageCount, err := repo.With(ctx, db).ScanItem[AgeCount](func(db *gorm.DB, cls *AccountColumns, res *AgeCountColumns) *gorm.DB {
    return db.Select(res.Merge(cls.Age.Name(), res.CountAs(res.AccountCount.Name()))).
              Where(cls.Age.Eq(18)).
              Group(cls.Age.Name())
})

// gormrepo - ScanOne（查一个值，如聚合结果；查不到行或查出多行都报错）
maxAge, err := repo.With(ctx, db).ScanOne[int](func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Select(cls.Age.COALESCE().Max())
})

// gormrepo - Where 条件: Eq, Ne, Gt, Lt, Gte, Lte
accounts, err := repo.With(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gt(18)).    // age > 18
              Where(cls.Age.Lt(60)).    // age < 60
              Where(cls.Age.Ne(30))     // age != 30
})

// gormrepo - Like / NotLike
accounts, err := repo.With(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Nickname.Like("%test%")).
              Where(cls.Username.NotLike("%admin%"))
})

// gormrepo - In / NotIn
accounts, err := repo.With(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.ID.In([]uint{1, 2, 3})).
              Where(cls.Age.NotIn([]int{18, 19, 20}))
})

// gormrepo - Between / NotBetween
accounts, err := repo.With(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Between(20, 40)).
              Where(cls.ID.NotBetween(100, 200))
})

// gormrepo - IsNull / IsNotNull
accounts, err := repo.With(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Nickname.IsNotNull())
})

// gormrepo - Or 条件
accounts, err := repo.With(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18)).
              Or(cls.Nickname.Eq("vip"))
})

// gormrepo - Order / Select
accounts, err := repo.With(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18)).
              Order(cls.Age.Desc()).
              Select(cls.ID.Name(), cls.Username.Name())
})
```

#### 查询操作

| 方法               | 参数                                                  | 返回值                   | 描述                             |
| ------------------ | ----------------------------------------------------- | ------------------------ | -------------------------------- |
| `First`            | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `*MOD, error`            | 查询第一个匹配记录               |
| `FirstE`           | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `*MOD, *ErrorOrNotExist` | 查询首条或返回不存在标识         |
| `Find`             | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `[]*MOD, error`          | 查询所有匹配记录                 |
| `FindAll`          | (无)                                                  | `[]*MOD, error`          | 查全部记录                       |
| `FindN`            | `where, ordering, size`                               | `[]*MOD, error`          | 取前 N 条匹配记录                |
| `FindO`            | `where, ordering`                                     | `[]*MOD, error`          | 带排序查询所有匹配记录           |
| `FindC`            | `where, ordering, paging`                             | `[]*MOD, int64, error`   | 自定义分页回调，带总数           |
| `FindPage`         | `where, ordering, pagination`                         | `[]*MOD, error`          | 分页查询                         |
| `FindPageAndCount` | `where, ordering, pagination`                         | `[]*MOD, int64, error`   | 分页查询带总数                   |
| `Count`            | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `int64, error`           | 统计匹配记录数                   |
| `Exist`            | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `bool, error`            | 检查记录是否存在                 |
| `Pluck[T]`         | `where, column`                                       | `[]T, error`             | 取某一列的值                     |
| `Scan[T]`          | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `[]T, error`             | 结果行扫描到结构体 T             |
| `ScanList[RES]`    | `where func(db *gorm.DB, cls CLS, res CLSR) *gorm.DB` | `[]RES, error`           | 返回列表，回调提供结果列         |
| `ScanItem[RES]`    | `where func(db *gorm.DB, cls CLS, res CLSR) *gorm.DB` | `RES, error`             | 返回恰好一条结果，回调提供结果列 |
| `ScanOne[T]`       | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `T, error`               | 查一个值，不是一行就报错         |

`ScanList` 和 `ScanItem` 都是 `GormRepo` 的方法，结果类型需要列定义（例如 `AgeCountColumns`），可以通过 gormcngen 生成。`ScanItem` 遇到零行或多行都会报错，不会自动取第一行。`Scan` 和 `ScanOne` 不要求结果类型具有列定义。

排序的传法：取前 N 条或分页的方法（`FindN`、`FindPage`、`FindPageAndCount`、`FindC`）都把排序作为必填参数，因为不排序的话，前 N 条是哪几条数据库不保证，翻页时数据也可能重复或漏掉。普通分页用 `FindPageAndCount`；分页部分要写逻辑时用 `FindC`（例如只有传了页码才分页），它的分页回调拿到已经带上条件和排序的 `*gorm.DB`，在上面加 `Limit` / `Offset`。`FindO` 是带排序的 `Find`。其余方法跟原生 GORM 一样，要排序就在 where 回调里写 `db.Order(cls.Age.Desc())`。

#### 5. 创建记录

```go
// 经典 GORM
err := db.Create(&Account{Username: "bob", Nickname: "Bob", Age: 25}).Error

// gormrepo
err := repo.With(ctx, db).Create(&Account{Username: "bob", Nickname: "Bob", Age: 25})

// gormrepo - CreateInBatches（插入很多条，每条 SQL 插 batchSize 条）
err := repo.With(ctx, db).CreateInBatches(accounts, 100)
```

#### 创建操作

| 方法              | 参数                     | 返回值  | 描述               |
| ----------------- | ------------------------ | ------- | ------------------ |
| `Create`          | `one *MOD`               | `error` | 创建新记录         |
| `Creates`         | `ones []*MOD`            | `error` | 批量创建记录       |
| `CreateInBatches` | `ones []*MOD, batchSize` | `error` | 分批创建记录       |
| `Save`            | `one *MOD`               | `error` | 插入或更新记录     |
| `Saves`           | `ones []*MOD`            | `error` | 批量插入或更新记录 |

#### 6. 更新记录

```go
// 经典 GORM
err := db.Model(&Account{}).Where("id = ?", 1).Updates(map[string]any{"age": 26}).Error

// gormrepo - Update（更新一列）
err := repo.With(ctx, db).Update(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.ID.Eq(1))
    },
    func(cls *AccountColumns) (string, any) {
        return cls.Age.Kv(26)
    },
)

// gormrepo - Updates（原生的 map[string]any，跟 GORM 写法一样）
err := repo.With(ctx, db).Updates(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.ID.Eq(1))
    },
    func(cls *AccountColumns) map[string]any {
        return cls.Kw(cls.Age.Kv(26)).
                  Kw(cls.Nickname.Kv("NewNick")).
                  AsMap()
    },
)

// gormrepo - UpdatesM（使用 ColumnValueMap，每个值都有类型检查，推荐写法）
err := repo.With(ctx, db).UpdatesM(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.ID.Eq(1))
    },
    func(cls *AccountColumns) gormcnm.ColumnValueMap {
        return cls.Kw(cls.Age.Kv(26)).
                  Kw(cls.Nickname.Kv("NewNick"))
    },
)

// gormrepo - UpdatesO (通过主键)
account := &Account{Model: gorm.Model{ID: 1}}
err := repo.With(ctx, db).UpdatesO(account, func(cls *AccountColumns) gormcnm.ColumnValueMap {
    return cls.Kw(cls.Age.Kv(26))
})

// gormrepo - UpdatesC (组合条件)
account := &Account{Model: gorm.Model{ID: 1}}
err := repo.With(ctx, db).UpdatesC(account,
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.Age.Lt(30))
    },
    func(cls *AccountColumns) gormcnm.ColumnValueMap {
        return cls.Kw(cls.Age.Kv(26))
    },
)
```

#### 更新操作

| 方法       | 参数                       | 返回值  | 描述                             |
| ---------- | -------------------------- | ------- | -------------------------------- |
| `Update`   | `where, newValue`          | `error` | 更新单个字段                     |
| `Updates`  | `where, newValues`         | `error` | 更新多个字段                     |
| `UpdatesM` | `where, newValues`         | `error` | 使用 ColumnValueMap 更新 (M=Map) |
| `UpdatesO` | `object, newValues`        | `error` | 通过主键更新 (O=Object)          |
| `UpdatesC` | `object, where, newValues` | `error` | 组合条件更新 (C=Combined)        |

#### 7. 删除记录

```go
// 经典 GORM
err := db.Where("id = ?", 1).Delete(&Account{}).Error

// gormrepo - Delete (通过实例)
account := &Account{Model: gorm.Model{ID: 1}}
err := repo.With(ctx, db).Delete(account)

// gormrepo - DeleteW (通过条件)
err := repo.With(ctx, db).DeleteW(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.ID.Eq(1))
})

// gormrepo - DeleteM (实例 + 条件)
account := &Account{Model: gorm.Model{ID: 1}}
err := repo.With(ctx, db).DeleteM(account, func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Lt(30))
})
```

#### 删除操作

| 方法      | 参数                                                  | 返回值  | 描述               |
| --------- | ----------------------------------------------------- | ------- | ------------------ |
| `Delete`  | `one *MOD`                                            | `error` | 通过实例删除记录   |
| `DeleteW` | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `error` | 通过条件删除记录   |
| `DeleteM` | `one *MOD, where func(db *gorm.DB, cls CLS) *gorm.DB` | `error` | 实例加条件删除记录 |

#### 8. 自定义操作

遇到未涵盖的操作，使用 `Invoke`：

```go
// Invoke - 批量更新
err := repo.With(ctx, db).Invoke(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18)).Update(cls.Nickname.Name(), "adult")
})

// Invoke - 字段自增
err := repo.With(ctx, db).Invoke(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.ID.Eq(1)).
        Update(cls.Age.Name(), gorm.Expr(cls.Age.Name()+" + ?", 1))
})

// Invoke - 复杂条件
err := repo.With(ctx, db).Invoke(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Between(18, 60)).
        Where(cls.Nickname.IsNotNull()).
        Updates(map[string]any{
            cls.Nickname.Name(): "active",
        })
})
```

#### 调用操作

| 方法     | 参数                                      | 返回值  | 描述           |
| -------- | ----------------------------------------- | ------- | -------------- |
| `Invoke` | `run func(db *gorm.DB, cls CLS) *gorm.DB` | `error` | 执行自定义操作 |

#### 9. 使用 Clauses 实现 Upsert

使用 `Clauses` 或 `Clause` 实现 upsert（冲突时插入或更新）：

```go
account := &Account{Username: "alice", Nickname: "Alice"}

// Clauses - 直接传入 clause.Expression
cls := account.Columns()
err := repo.With(ctx, db).Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: cls.Username.Name()}},
    DoUpdates: clause.AssignmentColumns([]string{cls.Nickname.Name()}),
}).Create(account)

// Clause - 类型安全，使用列定义构建
err := repo.With(ctx, db).Clause(func(cls *AccountColumns) clause.Expression {
    return clause.OnConflict{
        Columns:   []clause.Column{{Name: cls.Username.Name()}},
        DoUpdates: clause.AssignmentColumns([]string{cls.Nickname.Name()}),
    }
}).Create(account)

// 批量 upsert，使用 Clause + Creates
accounts := []*Account{
    {Username: "user1", Nickname: "nick1"},
    {Username: "user2", Nickname: "nick2"},
}
err := repo.With(ctx, db).Clause(func(cls *AccountColumns) clause.Expression {
    return clause.OnConflict{
        Columns:   []clause.Column{{Name: cls.Username.Name()}},
        DoUpdates: clause.AssignmentColumns([]string{cls.Nickname.Name()}),
    }
}).Creates(accounts)
```

#### 子句操作

| 方法      | 参数                              | 返回值      | 描述                  |
| --------- | --------------------------------- | ----------- | --------------------- |
| `Clauses` | `clauses ...clause.Expression`    | `*GormRepo` | 添加子句并返回新 repo |
| `Clause`  | `func(cls CLS) clause.Expression` | `*GormRepo` | 使用列定义构建子句    |

#### 10. Unscoped（永久删除 / 绕过软删除）

使用 `Unscoped()` 绕过软删除作用域，通过方法链返回新的 repo 实例：

```go
// 永久删除——绕过软删除，从数据库中彻底移除
err := repo.With(ctx, db).Unscoped().DeleteW(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Username.Eq(username))
})

// 查询已软删除的记录
accounts, err := repo.With(ctx, db).Unscoped().Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Username.Eq(username))
})
```

| 方法       | 参数 | 返回值      | 描述                         |
| ---------- | ---- | ----------- | ---------------------------- |
| `Unscoped` | (无) | `*GormRepo` | 通过链式调用禁用软删除作用域 |

#### 11. 可复用的 scope

把条件、排序或分页写成一次 GORM scope，任何查询都能通过 `db.Scopes` 加上它：

```go
adult := repo.NewScope(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18))
})
ageDescScope := repo.NewOrderingScope(func(cls *AccountColumns) gormcnm.Ordering {
    return cls.Age.Ob("desc")
})

accounts, err := repo.With(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Scopes(adult, ageDescScope)
})
```

| 方法               | 参数                                        | 返回值          | 描述             |
| ------------------ | ------------------------------------------- | --------------- | ---------------- |
| `NewScope`         | `where func(db *gorm.DB, cls CLS) *gorm.DB` | `ScopeFunction` | 带列定义的 scope |
| `NewOrderingScope` | `ordering func(cls CLS) gormcnm.Ordering`   | `ScopeFunction` | 排序 scope       |
| `NewPaginateScope` | `ordering, page *Pagination`                | `ScopeFunction` | 排序加分页 scope |

#### 12. 联表查询（gormjoin）

`gormjoin` 只负责连表：传两个表 repo 和连接条件，它给出 JOIN 语句和要选的列，其余交给 GORM，以 `*gorm.DB` 为纽带。表 repo 给出带表名的列，例如 `students.id`。结果结构体也用 gormcngen 生成列，起别名时直接指向结果的列，类型对不上编译就报错：

```go
type StudentScore struct { // 结果行，生成列后得到 StudentScoreColumns
    StudentName string
    Subject     string
    Points      int
}

studentRepo := gormtablerepo.NewRepo(gormclass.UseTable(&Student{}))
scoreRepo := gormtablerepo.NewRepo(gormclass.UseTable(&Score{}))

join := gormjoin.LeftJoin(studentRepo, scoreRepo, func(student *StudentColumns, score *ScoreColumns) []string {
    return []string{score.StudentID.OnEq(student.ID)}
})

var rows []StudentScore
err := db.Model(&Student{}).
    Joins(join.Stmt()).
    Select(join.Names[StudentScore](func(student *StudentColumns, score *ScoreColumns, res *StudentScoreColumns) []string {
        return []string{
            student.Name.AsName(res.StudentName),
            score.Subject.AsName(res.Subject),
            score.Points.AsName(res.Points),
        }
    })).
    Scopes(scoreRepo.Base().NewScope(func(db *gorm.DB, cls *ScoreColumns) *gorm.DB {
        return db.Where(cls.Points.Gte(60))
    })).
    Scan(&rows).Error
// SELECT students.name as student_name,scores.subject as subject,scores.points as points FROM `students`
// LEFT JOIN scores ON scores.student_id=students.id WHERE scores.points>=60
```

连三张表就是两段两表连接：每段的 `Stmt()` 各用一次 `Joins`，每段的 `Names` 用 `slices.Concat` 拼起来，见 [example09](internal/examples/example09)。软删除照 GORM 的规矩：GORM 通过 `db.Model` 给第一张表加，被连的表需要时自己写在连接条件里，例如 `score.DeletedAt.IsNull()`。[example10](internal/examples/example10) 对比了 gormjoin 和 GORM 的 `Preload`。

| 函数 / 方法                            | 参数                                      | 返回       | 说明                       |
| -------------------------------------- | ----------------------------------------- | ---------- | -------------------------- |
| `LeftJoin` / `RightJoin` / `InnerJoin` | `repo1, repo2, func(cls1, cls2) []string` | `*Join`    | 把 repo2 连到 repo1 上     |
| `New`                                  | `repo1, joinType, repo2, onFunc`          | `*Join`    | 指定连接类型               |
| `(*Join).Stmt`                         | (无)                                      | `string`   | JOIN 语句，传给 `db.Joins` |
| `(*Join).Names[RES]`                   | `func(cls1, cls2, res) []string`          | `[]string` | 要选的列，传给 `db.Select` |

---

## 📝 完整示例

查看 [examples](internal/examples) 目录获取完整集成示例。

---

## 关联项目

探索完整的 GORM 生态系统集成包：

### 核心生态

- **[gormcnm](https://github.com/yylego/gormcnm)** - GORM 基础层，提供类型安全的列操作和条件构建
- **[gormcngen](https://github.com/yylego/gormcngen)** - 使用 AST 的代码生成引擎，支持类型安全的 GORM 操作
- **[gormrepo](https://github.com/yylego/gormrepo)** - 仓储模式实现，遵循 GORM 最佳实践（本项目）
- **[gormmom](https://github.com/yylego/gormmom)** - 原生语言 GORM 标签生成引擎，支持智能列名
- **[gormzhcn](https://github.com/yylego/gormzhcn)** - 完整的 GORM 中文编程接口

每个包针对 GORM 开发的不同方面，从本地化到类型安全操作和代码生成。

---

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 许可证类型

MIT 许可证 - 详见 [LICENSE](LICENSE)。

---

## 💬 联系与反馈

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **问题报告？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **新颖思路？** 创建 issue 讨论
- 📖 **文档疑惑？** 报告问题，帮助我们完善文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，协助解决性能问题
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：面向用户的更改需要更新文档
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来贡献此项目。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) CLOSE: STANDARD PROJECT FOOTER -->

<!-- TEMPLATE (EN) BEGIN: BADGES -->

[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/gormrepo/release.yml?branch=main&label=BUILD)](https://github.com/yylego/gormrepo/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/gormrepo)](https://pkg.go.dev/github.com/yylego/gormrepo)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/gormrepo/main.svg)](https://coveralls.io/github/yylego/gormrepo?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.27%2B-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/yylego/gormrepo.svg)](https://github.com/yylego/gormrepo/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/gormrepo)](https://goreportcard.com/report/github.com/yylego/gormrepo)

<!-- TEMPLATE (EN) CLOSE: BADGES -->

# 🚀 gormrepo - Type-Safe Repos on GORM

**gormrepo** wraps GORM in generic repos: the conditions, the sorting and the new values are written with the columns that gormcngen generates, so a wrong column name / value type fails at compile time.

> 🌟 **The callbacks take the native `*gorm.DB`, so the code inside them is plain GORM code**

---

## Ecosystem

![GORM Type-Safe Ecosystem](https://github.com/yylego/gormcnm/raw/main/assets/gormcnm-ecosystem.svg)

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->

## CHINESE README

[中文说明](README.zh.md)

<!-- TEMPLATE (EN) CLOSE: LANGUAGE NAVIGATION -->

---

## 🔄 Tech Comparison

| Ecosystem             | Java MyBatis Plus  | Python SQLAlchemy | Go GORM Ecosystem |
| --------------------- | ------------------ | ----------------- | ----------------- |
| **Type-Safe Columns** | `Example::getName` | `Example.name`    | `cls.Name.Eq()`   |
| **Repo Pattern**      | ✅ BaseMapper      | ✅ Session API    | ✅ GormRepo       |

---

## 🚀 Quick Start

### Installation

```bash
go get github.com/yylego/gormrepo
```

### Complete Usage Flow

#### 1. Define The Model (Supporting Native Fields)

```go
type Account struct {
    gorm.Model // ID, CreatedAt, UpdatedAt, and DeletedAt that turns on soft delete
    Username string `gorm:"uniqueIndex"`
    Nickname string `gorm:"index"`
    Age      int
}
```

#### 2. Auto-Generate Column Structs (gormcngen)

```go
// Auto-generated
func (c *Account) Columns() *AccountColumns {
    return c.TableColumns(gormcnm.NewPlainDecoration())
}

type AccountColumns struct {
    gormcnm.ColumnOperationClass // embedded aids, e.g. cls.Kw(...) and cls.Merge(...)
    ID        gormcnm.ColumnName[uint]
    CreatedAt gormcnm.ColumnName[time.Time]
    UpdatedAt gormcnm.ColumnName[time.Time]
    DeletedAt gormcnm.ColumnName[gorm.DeletedAt]
    Username  gormcnm.ColumnName[string]
    Nickname  gormcnm.ColumnName[string]
    Age       gormcnm.ColumnName[int]
}

// TableColumns gives the columns with a decoration, e.g. the table name in joins: accounts.id
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

#### 3. Create Repo

```go
// Create repo with columns
repo := gormrepo.NewRepo(&Account{}, (&Account{}).Columns())

// Concise approach with gormrepo/gormclass
repo := gormrepo.NewRepo(gormclass.Use(&Account{}))
```

#### 4. Select Records

```go
// Classic GORM
var alice Account
err := db.Where("username = ?", "alice").First(&alice).Error

// gormrepo - First
account, err := repo.With(ctx, db).First(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Username.Eq("alice"))
})

// gormrepo - FirstE (returns ErrorOrNotExist)
account, erb := repo.With(ctx, db).FirstE(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Username.Eq("alice"))
})

// gormrepo - Find
accounts, err := repo.With(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18))
})

// gormrepo - FindN (the first N records)
accounts, err := repo.With(ctx, db).FindN(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.Age.Gte(18))
    },
    func(cls *AccountColumns) gormcnm.Ordering {
        return cls.Age.Ob("DESC")
    },
    10,
)

// gormrepo - FindAll (each record, no condition)
accounts, err := repo.With(ctx, db).FindAll()

// gormrepo - FindO (with ordering)
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

// gormrepo - NewPagination (from a page, counted from 1, and a page size; page 3 at size 20 is Limit 20, Offset 40)
page, err := gormrepo.NewPagination(3, 20)
accounts, count, err := repo.With(ctx, db).FindPageAndCount(where, ordering, page)
pageCount := page.PageCount(count) // 41 records make 3 pages

// gormrepo - FindC (the page with a custom paging callback, and the count across pages)
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

// gormrepo - Pluck (one column, result type comes from the column: []string here)
names, err := repo.With(ctx, db).Pluck(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18))
}, func(cls *AccountColumns) gormcnm.ColumnName[string] {
    return cls.Username
})

// Result rows, which are not tables: put them beside the models and generate the columns of them too
type AccountInfo struct {
    ID       uint
    Username string
}
type AgeCount struct {
    Age          int
    AccountCount int64
}

// gormrepo - Scan (rows into a result struct, fits Select / Joins / Group)
infos, err := repo.With(ctx, db).Scan[AccountInfo](func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Select(cls.ID.Name(), cls.Username.Name()).Where(cls.Age.Gte(18))
})

// gormrepo - ScanList (result rows, with result columns available in the callback)
ageCounts, err := repo.With(ctx, db).ScanList[AgeCount](func(db *gorm.DB, cls *AccountColumns, res *AgeCountColumns) *gorm.DB {
    return db.Select(res.Merge(cls.Age.Name(), res.CountAs(res.AccountCount.Name()))).
              Group(cls.Age.Name()).
              Order(res.AccountCount.Desc())
})

// gormrepo - ScanItem (just one result row, with result columns in the callback)
ageCount, err := repo.With(ctx, db).ScanItem[AgeCount](func(db *gorm.DB, cls *AccountColumns, res *AgeCountColumns) *gorm.DB {
    return db.Select(res.Merge(cls.Age.Name(), res.CountAs(res.AccountCount.Name()))).
              Where(cls.Age.Eq(18)).
              Group(cls.Age.Name())
})

// gormrepo - ScanOne (one value such as an aggregate, fails on zero rows and on two rows and up)
maxAge, err := repo.With(ctx, db).ScanOne[int](func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Select(cls.Age.COALESCE().Max())
})

// gormrepo - Where conditions: Eq, Ne, Gt, Lt, Gte, Lte
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

// gormrepo - Or conditions
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

#### Select Operations

| Method             | Parameters                                            | Returns                  | Description                                             |
| ------------------ | ----------------------------------------------------- | ------------------------ | ------------------------------------------------------- |
| `First`            | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `*MOD, error`            | Find first matching record                              |
| `FirstE`           | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `*MOD, *ErrorOrNotExist` | Find first with not-exist indication                    |
| `Find`             | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `[]*MOD, error`          | Find matching records                                   |
| `FindAll`          | (none)                                                | `[]*MOD, error`          | Find each record                                        |
| `FindN`            | `where, ordering, size`                               | `[]*MOD, error`          | Find the first N matching records                       |
| `FindO`            | `where, ordering`                                     | `[]*MOD, error`          | Find matching records with ordering                     |
| `FindC`            | `where, ordering, paging`                             | `[]*MOD, int64, error`   | Page with a custom paging callback, with record count   |
| `FindPage`         | `where, ordering, pagination`                         | `[]*MOD, error`          | Paginated search                                        |
| `FindPageAndCount` | `where, ordering, pagination`                         | `[]*MOD, int64, error`   | Paginated search with record count                      |
| `Count`            | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `int64, error`           | Count matching records                                  |
| `Exist`            | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `bool, error`            | Check if records exist                                  |
| `Pluck[T]`         | `where, column`                                       | `[]T, error`             | Values of one column                                    |
| `Scan[T]`          | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `[]T, error`             | Rows into result struct T                               |
| `ScanList[RES]`    | `where func(db *gorm.DB, cls CLS, res CLSR) *gorm.DB` | `[]RES, error`           | Result rows with result columns in the callback         |
| `ScanItem[RES]`    | `where func(db *gorm.DB, cls CLS, res CLSR) *gorm.DB` | `RES, error`             | Just one result row with result columns in the callback |
| `ScanOne[T]`       | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `T, error`               | One value, fails unless one row                         |

`ScanList` and `ScanItem` are methods on `GormRepo`. Both require result column definitions (such as `AgeCountColumns`), which can be generated with gormcngen. `ScanItem` returns an error on zero rows and on two rows and up; it does not pick the first row on its own. `Scan` and `ScanOne` do not require result column definitions.

How to pass the ordering: the methods that take the first N rows / a page (`FindN`, `FindPage`, `FindPageAndCount`, `FindC`) take it as a required param, since without ORDER BY the first N rows come in no fixed sequence, and paging can give repeats and gaps across pages. Use `FindPageAndCount` on plain paging, and `FindC` when the paging needs code (e.g. paging just when a page index is given): its paging callback gets the `*gorm.DB` with the conditions and the ordering on it, then adds `Limit` / `Offset`. `FindO` is `Find` with an ordering. The rest work as GORM does: write `db.Order(cls.Age.Desc())` in the where callback when an ordering is wanted.

#### 5. Create Records

```go
// Classic GORM
err := db.Create(&Account{Username: "bob", Nickname: "Bob", Age: 25}).Error

// gormrepo
err := repo.With(ctx, db).Create(&Account{Username: "bob", Nickname: "Bob", Age: 25})

// gormrepo - Creates (pass the address of the slice, as db.Create(&accounts) does in GORM)
accounts := []*Account{{Username: "chloe"}, {Username: "david"}}
err := repo.With(ctx, db).Creates(&accounts)

// gormrepo - CreateInBatches (a batch of records, batchSize records each SQL)
err := repo.With(ctx, db).CreateInBatches(accounts, 100)
```

#### Create Operations

| Method            | Parameters               | Returns | Description                 |
| ----------------- | ------------------------ | ------- | --------------------------- |
| `Create`          | `one *MOD`               | `error` | Create new record           |
| `Creates`         | `ones *[]*MOD`           | `error` | Batch create records        |
| `CreateInBatches` | `ones []*MOD, batchSize` | `error` | Create in batches           |
| `Save`            | `one *MOD`               | `error` | Insert/update record        |
| `Saves`           | `ones *[]*MOD`           | `error` | Batch insert/update records |

`Creates` and `Saves` take the address of the slice: with `clause.Returning{}` GORM rewrites the slice into the records the database returned, e.g. plus `OnConflict{DoNothing: true}` the slice shrinks to the inserted records. `CreateInBatches` takes the slice itself and keeps its length, so use `Creates` when the inserted records are wanted.

#### 6. Update Records

```go
// Classic GORM
err := db.Model(&Account{}).Where("id = ?", 1).Updates(map[string]any{"age": 26}).Error

// gormrepo - Update (one column)
err := repo.With(ctx, db).Update(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.ID.Eq(1))
    },
    func(cls *AccountColumns) (string, any) {
        return cls.Age.Kv(26)
    },
)

// gormrepo - Updates (a native map[string]any, as GORM does it)
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

// gormrepo - UpdatesM (using ColumnValueMap, each value type-checked, the recommended style)
err := repo.With(ctx, db).UpdatesM(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.ID.Eq(1))
    },
    func(cls *AccountColumns) gormcnm.ColumnValueMap {
        return cls.Kw(cls.Age.Kv(26)).
                  Kw(cls.Nickname.Kv("NewNick"))
    },
)

// gormrepo - UpdatesO (by primary key)
account := &Account{Model: gorm.Model{ID: 1}}
err := repo.With(ctx, db).UpdatesO(account, func(cls *AccountColumns) gormcnm.ColumnValueMap {
    return cls.Kw(cls.Age.Kv(26))
})

// gormrepo - UpdatesC (combined conditions)
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

#### Update Operations

| Method     | Parameters                 | Returns | Description                                  |
| ---------- | -------------------------- | ------- | -------------------------------------------- |
| `Update`   | `where, newValue`          | `error` | Update single field                          |
| `Updates`  | `where, newValues`         | `error` | Update multiple fields                       |
| `UpdatesM` | `where, newValues`         | `error` | Update with ColumnValueMap (M=Map)           |
| `UpdatesO` | `object, newValues`        | `error` | Update via primary key (O=Object)            |
| `UpdatesC` | `object, where, newValues` | `error` | Update with combined conditions (C=Combined) |

#### 7. Delete Records

```go
// Classic GORM
err := db.Where("id = ?", 1).Delete(&Account{}).Error

// gormrepo - Delete (by instance)
account := &Account{Model: gorm.Model{ID: 1}}
err := repo.With(ctx, db).Delete(account)

// gormrepo - DeleteW (by where conditions)
err := repo.With(ctx, db).DeleteW(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.ID.Eq(1))
})

// gormrepo - DeleteM (instance + conditions)
account := &Account{Model: gorm.Model{ID: 1}}
err := repo.With(ctx, db).DeleteM(account, func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Lt(30))
})
```

#### Delete Operations

| Method    | Parameters                                            | Returns | Description                 |
| --------- | ----------------------------------------------------- | ------- | --------------------------- |
| `Delete`  | `one *MOD`                                            | `error` | Delete record via instance  |
| `DeleteW` | `where func(db *gorm.DB, cls CLS) *gorm.DB`           | `error` | Delete via conditions       |
| `DeleteM` | `one *MOD, where func(db *gorm.DB, cls CLS) *gorm.DB` | `error` | Delete item with conditions |

#### 8. Custom Operations

When above methods do not fit, use `Invoke`, which runs the whole GORM chain in the callback, and `Where`, which returns the `*gorm.DB` with the conditions on it. The `*gorm.DB` in these two comes with no `Model` set, so use `Modeled()` first (the same as `db.Model(&Account{})`) when the chain ends in `Update` / `Updates` / `Count` / `Pluck`, which need the table:

```go
// Invoke - batch update
err := repo.With(ctx, db).Modeled().Invoke(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18)).Update(cls.Nickname.Name(), "adult")
})

// Invoke - increment field
err := repo.With(ctx, db).Modeled().Invoke(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.ID.Eq(1)).
        Update(cls.Age.Name(), gorm.Expr(cls.Age.Name()+" + ?", 1))
})

// Invoke - complex conditions
err := repo.With(ctx, db).Modeled().Invoke(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Between(18, 60)).
        Where(cls.Nickname.IsNotNull()).
        Updates(map[string]any{
            cls.Nickname.Name(): "active",
        })
})

// Where - take the *gorm.DB with the conditions, then go on with native GORM
var count int64
err := repo.With(ctx, db).Modeled().Where(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18))
}).Count(&count).Error
```

#### Invoke Operations

| Method    | Parameters                                  | Returns     | Description                                |
| --------- | ------------------------------------------- | ----------- | ------------------------------------------ |
| `Invoke`  | `run func(db *gorm.DB, cls CLS) *gorm.DB`   | `error`     | Execute custom operation                   |
| `Where`   | `where func(db *gorm.DB, cls CLS) *gorm.DB` | `*gorm.DB`  | Add the conditions, return the native db   |
| `Modeled` | (none)                                      | `*GormRepo` | New repo with `db.Model(new(MOD))` applied |

#### 9. Upsert with Clauses

Use `Clauses` / `Clause` to handle upsert (insert / update on conflict):

```go
account := &Account{Username: "alice", Nickname: "Alice"}

// Clauses - pass clause.Expression
cls := account.Columns()
err := repo.With(ctx, db).Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: cls.Username.Name()}},
    DoUpdates: clause.AssignmentColumns([]string{cls.Nickname.Name()}),
}).Create(account)

// Clause - type-safe with column definitions
err := repo.With(ctx, db).Clause(func(cls *AccountColumns) clause.Expression {
    return clause.OnConflict{
        Columns:   []clause.Column{{Name: cls.Username.Name()}},
        DoUpdates: clause.AssignmentColumns([]string{cls.Nickname.Name()}),
    }
}).Create(account)

// Batch upsert with Clause + Creates
accounts := []*Account{
    {Username: "user1", Nickname: "nick1"},
    {Username: "user2", Nickname: "nick2"},
}
err := repo.With(ctx, db).Clause(func(cls *AccountColumns) clause.Expression {
    return clause.OnConflict{
        Columns:   []clause.Column{{Name: cls.Username.Name()}},
        DoUpdates: clause.AssignmentColumns([]string{cls.Nickname.Name()}),
    }
}).Creates(&accounts)
```

#### Clause Operations

| Method    | Parameters                        | Returns     | Description                          |
| --------- | --------------------------------- | ----------- | ------------------------------------ |
| `Clauses` | `clauses ...clause.Expression`    | `*GormRepo` | Add clauses and return new repo      |
| `Clause`  | `func(cls CLS) clause.Expression` | `*GormRepo` | Build clause with column definitions |

#### 10. Unscoped (Permanent Deletion / Soft-Delete Bypass)

Use `Unscoped()` to bypass soft-delete scope. Returns a new repo instance via method chaining:

```go
// Permanent deletion — bypasses soft delete, removes from database
err := repo.With(ctx, db).Unscoped().DeleteW(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Username.Eq(username))
})

// Read soft-deleted records
accounts, err := repo.With(ctx, db).Unscoped().Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Username.Eq(username))
})
```

| Method     | Parameters | Returns     | Description                            |
| ---------- | ---------- | ----------- | -------------------------------------- |
| `Unscoped` | (none)     | `*GormRepo` | Disable soft-delete scope via chaining |

#### 11. Reusable Scopes

Write a condition, an ordering / a page once as a GORM scope, then add it to each find via `db.Scopes`:

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

| Method             | Parameters                                  | Returns         | Description                  |
| ------------------ | ------------------------------------------- | --------------- | ---------------------------- |
| `NewScope`         | `where func(db *gorm.DB, cls CLS) *gorm.DB` | `ScopeFunction` | A scope with the columns     |
| `NewOrderingScope` | `ordering func(cls CLS) gormcnm.Ordering`   | `ScopeFunction` | An ordering scope            |
| `NewPaginateScope` | `ordering, page *Pagination`                | `ScopeFunction` | An ordering and paging scope |

#### 12. Join Tables (gormjoin)

`gormjoin` does the join and nothing else: pass two table repos and the ON condition, it gives the JOIN statement and the select columns. The rest stays with GORM, with `*gorm.DB` as the link. Table repos give the columns with the table name, e.g. `students.id`. Generate columns for the result struct with gormcngen too, then each alias points at a result column and a type mismatch fails at compile time:

```go
type StudentScore struct { // result row, generate its columns to get StudentScoreColumns
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

Joining three tables is two joins of two tables: add each `Stmt()` via `Joins`, and concat each `Names` via `slices.Concat`, see [example09](internal/examples/example09). Soft delete works as in GORM: GORM adds it to the first table via `db.Model`, and the joined table needs it in the ON condition when wanted, e.g. `score.DeletedAt.IsNull()`. [example10](internal/examples/example10) compares gormjoin with the GORM `Preload`.

| Function / Method                      | Parameters                                | Returns    | Description                           |
| -------------------------------------- | ----------------------------------------- | ---------- | ------------------------------------- |
| `LeftJoin` / `RightJoin` / `InnerJoin` | `repo1, repo2, func(cls1, cls2) []string` | `*Join`    | Join repo2 onto repo1                 |
| `CrossJoin`                            | `repo1, repo2`                            | `*Join`    | CROSS JOIN, which has no ON condition |
| `New`                                  | `repo1, joinType, repo2, onFunc`          | `*Join`    | Join with the given join type         |
| `(*Join).Stmt`                         | (none)                                    | `string`   | JOIN statement to pass to `db.Joins`  |
| `(*Join).Names[RES]`                   | `func(cls1, cls2, res) []string`          | `[]string` | Select columns to pass to `db.Select` |

CROSS JOIN takes no ON condition in standard SQL (PostgreSQL rejects `CROSS JOIN ... ON`), so `CrossJoin` takes the two repos and nothing else:

```go
join := gormjoin.CrossJoin(studentRepo, scoreRepo)
// join.Stmt() gives: CROSS JOIN scores
```

#### 13. GormWrap (Returns `*gorm.DB`)

`repo.With(ctx, db)` / `repo.Repo(db)` give a `GormRepo`, whose methods return `(result, error)`. `repo.Wrap(ctx, db)` / `repo.Gorm(db)` give a `GormWrap`: the same params, and each method returns `*gorm.DB` as GORM does, so `RowsAffected` is at hand; its `First` / `Find` take the dest as GORM does. `GormWrap` keeps to the GORM methods, the extras (`FindO`, `FindC`, `Count`, `Exist`, `Pluck`, `Scan...`) are on `GormRepo`:

```go
res := repo.Wrap(ctx, db).Update(
    func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
        return db.Where(cls.Age.Gte(18))
    },
    func(cls *AccountColumns) (string, any) {
        return cls.Nickname.Kv("adult")
    },
)
// res.Error, res.RowsAffected

var accounts []*Account
err := repo.Wrap(ctx, db).Find(func(db *gorm.DB, cls *AccountColumns) *gorm.DB {
    return db.Where(cls.Age.Gte(18))
}, &accounts).Error
```

| Method                      | On                      | Returns                   | Description                      |
| --------------------------- | ----------------------- | ------------------------- | -------------------------------- |
| `With(ctx, db)`             | `Repo`                  | `*GormRepo`               | Repo with the context and the db |
| `Repo(db)`                  | `Repo`                  | `*GormRepo`               | Repo with the db                 |
| `Wrap(ctx, db)`             | `Repo`                  | `*GormWrap`               | Wrap with the context and the db |
| `Gorm(db)`                  | `Repo`                  | `*GormWrap`               | Wrap with the db                 |
| `Gorm()` / `Repo()`         | `GormRepo` / `GormWrap` | `*GormWrap` / `*GormRepo` | Convert one into the next        |
| `WithContext(ctx)`          | `GormRepo` / `GormWrap` | same type                 | Set the context                  |
| `First(where, dest *MOD)`   | `GormWrap`              | `*gorm.DB`                | First matching record into dest  |
| `Find(where, dest *[]*MOD)` | `GormWrap`              | `*gorm.DB`                | Matching records into dest       |

---

## 📝 Complete Examples

Check [examples](internal/examples) DIR with complete integration examples.

---

## Related Projects

Explore the complete GORM ecosystem with these integrated packages:

### Core Ecosystem

- **[gormcnm](https://github.com/yylego/gormcnm)** - GORM foundation providing type-safe column queries operations
- **[gormcngen](https://github.com/yylego/gormcngen)** - Code generation engine using AST, enables type-safe GORM operations
- **[gormrepo](https://github.com/yylego/gormrepo)** - Repo pattern implementation with GORM best practices (this project)
- **[gormmom](https://github.com/yylego/gormmom)** - Native language GORM tag generation engine with smart column naming
- **[gormzhcn](https://github.com/yylego/gormzhcn)** - Complete Chinese programming interface with GORM

Each package targets different aspects of GORM development, from localization to type-safe operations and code generation.

---

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 License

MIT License - see [LICENSE](LICENSE).

---

## 💬 Contact & Feedback

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Mistake reports?** Open an issue on GitHub with reproduction steps
- 💡 **Fresh ideas?** Create an issue to discuss
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) CLOSE: STANDARD PROJECT FOOTER -->

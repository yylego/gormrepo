package gormrepo

import (
	"github.com/pkg/errors"
	"github.com/yylego/gormcnm"
	"gorm.io/gorm"
)

// ScopeFunction is a type alias representing a function that modifies a GORM DB instance
// Used with db.Scopes() to set custom conditions
// See: https://github.com/go-gorm/gorm/blob/c44405a25b0fb15c20265e672b8632b8774793ca/chainable_api.go#L376
//
// ScopeFunction 是表示修改 GORM DB 实例的函数的类型别名
// 与 db.Scopes() 配合使用以设置自定义条件
// 参见：https://github.com/go-gorm/gorm/blob/c44405a25b0fb15c20265e672b8632b8774793ca/chainable_api.go#L376
type ScopeFunction = func(db *gorm.DB) *gorm.DB

// NewScope creates a GORM scope function that applies a custom where condition
// Returns a ScopeFunction that can be used with db.Scopes()
//
// NewScope 创建应用自定义 where 条件的 GORM scope 函数
// 返回可与 db.Scopes() 配合使用的 ScopeFunction
func (repo *BaseRepo[MOD, CLS]) NewScope(where func(db *gorm.DB, cls CLS) *gorm.DB) ScopeFunction {
	return func(db *gorm.DB) *gorm.DB {
		return where(db, repo.cls)
	}
}

// NewOrderingScope creates a GORM scope function that applies an ordering condition
// Returns a ScopeFunction that can be used with db.Scopes() to set the ordering
//
// NewOrderingScope 创建应用排序条件的 GORM scope 函数
// 返回可与 db.Scopes() 配合使用以设置排序的 ScopeFunction
func (repo *BaseRepo[MOD, CLS]) NewOrderingScope(ordering func(cls CLS) gormcnm.Ordering) ScopeFunction {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(ordering(repo.cls).Ox())
	}
}

// Pagination defines parameters enabling paginated queries
// Contains limit (max records) and offset (records to skip)
//
// Pagination 定义分页查询参数
// 包含 limit（最大记录数）和 offset（跳过的记录数）
type Pagination struct {
	Limit  int // Max records to retrieve // 最大检索记录数
	Offset int // Records to skip // 跳过的记录数
}

// validate returns an error on a nil Pagination, and on Limit 0, which GORM turns into LIMIT 0 and gives no rows
// Limit -1 is fine, GORM takes it as no limit
//
// validate 在 Pagination 为 nil，或 Limit 为 0 时返回错误：GORM 会把 0 变成 LIMIT 0，什么都查不到
// Limit 为 -1 是合法的，GORM 当作不限条数
func (p *Pagination) validate() error {
	if p == nil {
		return errors.New("pagination is missing")
	}
	if p.Limit == 0 {
		return errors.New("pagination limit is 0, which gives no rows, use -1 to get each row")
	}
	return nil
}

// Scope creates a GORM scope function that applies pagination parameters
// Returns a ScopeFunction that sets limit and offset during pagination
// A missing Pagination and Limit 0 put an error on the db, the same check as in FindPage
//
// Scope 创建应用分页参数的 GORM scope 函数
// 返回在分页时设置 limit 和 offset 的 ScopeFunction
// Pagination 没传或 Limit 为 0 时把错误放到 db 上，跟 FindPage 是同一个检查
func (p *Pagination) Scope() ScopeFunction {
	return func(db *gorm.DB) *gorm.DB {
		if err := p.validate(); err != nil {
			_ = db.AddError(err)
			return db
		}
		return db.Limit(p.Limit).Offset(p.Offset)
	}
}

// NewPaginateScope creates a GORM scope function that applies pagination and ordering
// Combines ordering with limit and offset in a single scope, the page goes through the same check as in Pagination.Scope
//
// NewPaginateScope 创建应用分页和排序的 GORM scope 函数
// 在单个 scope 中组合排序、limit 和 offset，分页参数的检查跟 Pagination.Scope 相同
func (repo *BaseRepo[MOD, CLS]) NewPaginateScope(ordering func(cls CLS) gormcnm.Ordering, page *Pagination) ScopeFunction {
	return func(db *gorm.DB) *gorm.DB {
		return page.Scope()(db.Order(ordering(repo.cls).Ox()))
	}
}

package gormrepo

import (
	"math"

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

// NewPagination builds a Pagination from a page, counted from 1, and a page size
// Page 1 with size 20 gives Limit 20 and Offset 0, page 3 gives Offset 40
// A page below 1, a page size below 1, and an offset too big to fit in int each give an error:
// defaults and size limits belong to the code that calls it, such as "size 20 when the client sends none, 200 at most"
//
// NewPagination 用从 1 开始的页码和每页条数构造 Pagination
// 第 1 页、每页 20 条得到 Limit 20、Offset 0，第 3 页得到 Offset 40
// 页码小于 1、每页条数小于 1、偏移大到 int 放不下，都返回错误：
// 默认值和上限由调用方自己定，比如"客户端没传就按 20 条、最多 200 条"
func NewPagination(pageNum, pageSize int) (*Pagination, error) {
	if pageNum < 1 {
		return nil, errors.Errorf("page %d is below 1", pageNum)
	}
	if pageSize < 1 {
		return nil, errors.Errorf("page size %d is below 1", pageSize)
	}
	if pageNum-1 > math.MaxInt/pageSize {
		return nil, errors.Errorf("page %d with page size %d gives an offset that does not fit in int", pageNum, pageSize)
	}
	return &Pagination{Limit: pageSize, Offset: (pageNum - 1) * pageSize}, nil
}

// PageCount counts the pages that the given record count makes at this page size, rounding up
// 41 records at Limit 20 make 3 pages, and 0 records make 0 pages
// A Limit below 1 has no page size to count with, so the result is 0
//
// PageCount 按这个每页条数算出给定的记录总数有多少页，向上取整
// 每页 20 条时 41 条记录是 3 页，0 条记录是 0 页
// Limit 小于 1 时没有每页条数可算，结果为 0
func (p *Pagination) PageCount(count int64) int64 {
	if p == nil || p.Limit < 1 || count < 1 {
		return 0
	}
	size := int64(p.Limit)
	return (count-1)/size + 1
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

package gormrepo

import (
	"github.com/pkg/errors"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormrepo/gormclass"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GormRepo is the main repo built on GormWrap, each method returns (T, error) in place of *gorm.DB
// Beyond the GORM methods it adds the extras GORM does not give, e.g. FindO / FindC / Exist / ScanOne / ScanList
// New conveniences go here, and GormWrap stays as the plain GORM counterpart
//
// GormRepo 是建在 GormWrap 之上的主仓储，每个方法返回 (结果, error)，不返回 *gorm.DB
// 在 GORM 官方方法之外，加上 GORM 没有的便利，如 FindO / FindC / Exist / ScanOne / ScanList
// 新的便利都加在这里，GormWrap 只保持跟 GORM 官方一一对应
type GormRepo[MOD any, CLS any] struct {
	db  *gorm.DB // Database connection // 数据库连接
	cls CLS      // Column definitions // 列定义
}

// NewGormRepo creates a new GormRepo instance with database connection and column definitions
// The MOD param is used to deduce the type, its value is not used
//
// NewGormRepo 使用数据库连接和列定义创建新的 GormRepo 实例
// MOD 参数用于类型推断，其值不使用
func NewGormRepo[MOD any, CLS any](db *gorm.DB, _ *MOD, cls CLS) *GormRepo[MOD, CLS] {
	return &GormRepo[MOD, CLS]{
		db:  db,
		cls: cls,
	}
}

// First finds the first record matching the where condition
// Returns the found record or error if not found or query fails
//
// First 查找符合 where 条件的第一条记录
// 返回找到的记录，如果未找到或查询失败则返回错误
func (repo *GormRepo[MOD, CLS]) First(where func(db *gorm.DB, cls CLS) *gorm.DB) (*MOD, error) {
	var result = new(MOD)
	if err := repo.Gorm().First(where, result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// FirstE finds the first record with structured error handling
// Returns ErrorOrNotExist to distinguish between not found and other errors
//
// FirstE 查找第一条记录，带有结构化错误处理
// 返回 ErrorOrNotExist 以区分记录未找到和其他错误
func (repo *GormRepo[MOD, CLS]) FirstE(where func(db *gorm.DB, cls CLS) *gorm.DB) (*MOD, *ErrorOrNotExist) {
	var result = new(MOD)
	if err := repo.Gorm().First(where, result).Error; err != nil {
		return nil, NewErrorOrNotExist(err)
	}
	return result, nil
}

// Where applies the where condition and returns the gorm.DB to enable chaining
// Use when you need custom operations not provided by GormRepo
//
// Where 应用 where 条件并返回 gorm.DB 以便链式调用
// 当需要 GormRepo 未提供的自定义操作时使用
func (repo *GormRepo[MOD, CLS]) Where(where func(db *gorm.DB, cls CLS) *gorm.DB) *gorm.DB {
	return where(repo.db, repo.cls)
}

// Exist checks if any record matches the where condition
// Returns true if at least one record exists, false otherwise
//
// Exist 检查是否存在符合 where 条件的记录
// 如果至少存在一条记录则返回 true，否则返回 false
func (repo *GormRepo[MOD, CLS]) Exist(where func(db *gorm.DB, cls CLS) *gorm.DB) (bool, error) {
	var exists bool
	if err := where(repo.db, repo.cls).Model((*MOD)(nil)).Select("1").Limit(1).Find(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

// Find retrieves all records matching the where condition
// Returns slice of records or error if query fails
//
// Find 检索所有符合 where 条件的记录
// 返回记录切片，如果查询失败则返回错误
func (repo *GormRepo[MOD, CLS]) Find(where func(db *gorm.DB, cls CLS) *gorm.DB) ([]*MOD, error) {
	var results []*MOD
	if err := repo.Gorm().Find(where, &results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// FindAll retrieves each record of the MOD table, with no condition
// Saves passing a where function that returns db unchanged
//
// FindAll 检索 MOD 表的全部记录，不带条件
// 省去传一个原样返回 db 的 where 函数
func (repo *GormRepo[MOD, CLS]) FindAll() ([]*MOD, error) {
	return repo.Find(func(db *gorm.DB, cls CLS) *gorm.DB {
		return db
	})
}

// FindO retrieves records matching the where condition with ordering
// O = Ordering, returns all matching records in specified sequence
//
// FindO 检索符合 where 条件的记录并按指定顺序排序
// O = Ordering，返回所有匹配记录并按指定顺序排列
func (repo *GormRepo[MOD, CLS]) FindO(where func(db *gorm.DB, cls CLS) *gorm.DB, ordering func(cls CLS) gormcnm.Ordering) ([]*MOD, error) {
	var results []*MOD
	db := where(repo.db, repo.cls)
	db = db.Order(ordering(repo.cls).Ox())
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// FindN retrieves the first size records matching the where condition, in the given ordering
// The ordering is a required param: without ORDER BY the first N records come in no fixed sequence
//
// FindN 按指定排序检索符合 where 条件的前 size 条记录
// 排序是必填参数：不排序的话，"前 N 条"是哪几条数据库不保证
func (repo *GormRepo[MOD, CLS]) FindN(where func(db *gorm.DB, cls CLS) *gorm.DB, ordering func(cls CLS) gormcnm.Ordering, size int) ([]*MOD, error) {
	var results []*MOD
	if err := where(repo.db, repo.cls).Order(ordering(repo.cls).Ox()).Limit(size).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// FindPage retrieves paginated records with ordering
// Uses Pagination struct to give the offset and the limit
// The ordering is a required param: paging without ORDER BY can give repeats and gaps across pages
//
// FindPage 使用排序检索分页记录
// 使用 Pagination 结构体指定偏移和限制
// 排序是必填参数：不排序就分页，翻页时数据可能重复或漏掉
func (repo *GormRepo[MOD, CLS]) FindPage(where func(db *gorm.DB, cls CLS) *gorm.DB, ordering func(cls CLS) gormcnm.Ordering, page *Pagination) ([]*MOD, error) {
	if err := page.validate(); err != nil {
		return nil, err
	}
	db := where(repo.db, repo.cls)
	db = db.Order(ordering(repo.cls).Ox())
	db = db.Limit(page.Limit).Offset(page.Offset)
	var results []*MOD
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// FindPageAndCount retrieves paginated records with ordering and returns total count
// Same as FindC with the Pagination as the paging callback, the ordering is required as in FindPage
//
// FindPageAndCount 使用排序检索分页记录并返回总数
// 等同于 FindC，分页回调用 Pagination；排序跟 FindPage 一样是必填参数
func (repo *GormRepo[MOD, CLS]) FindPageAndCount(where func(db *gorm.DB, cls CLS) *gorm.DB, ordering func(cls CLS) gormcnm.Ordering, page *Pagination) ([]*MOD, int64, error) {
	if err := page.validate(); err != nil {
		return nil, 0, err
	}
	return repo.FindC(where, ordering, func(db *gorm.DB, cls CLS) *gorm.DB {
		return db.Limit(page.Limit).Offset(page.Offset)
	})
}

// FindC retrieves the page records and the count of matching records across pages
// The ordering is a required param, and the paging callback gets the *gorm.DB with the where and the ordering on it, then adds Limit/Offset
// Use FindPageAndCount on plain paging; FindC fits paging that needs code, e.g. paging just when a page index is given
// The where and the paging build one chain: the page uses it as is, and the count cancels LIMIT/OFFSET on it (GORM drops ORDER BY when counting),
// so the two callbacks give the same result regardless of sequence; GROUP BY in them turns the count into a group count
//
// FindC 检索当前页记录，并返回去掉分页后符合条件的总数
// 排序是必填参数，分页回调拿到的是已经带上 where 和排序的 *gorm.DB，在上面加 Limit/Offset
// 普通分页用 FindPageAndCount；分页部分要写逻辑时用 FindC，例如只有传了页码才分页
// where 和分页回调先拼成同一个查询：查当前页直接用它，查总数时在它上面撤掉 LIMIT/OFFSET（ORDER BY 由 GORM 计数时自动去掉），
// 所以两个回调谁先谁后结果都一样；回调里带 GROUP BY 时，总数会变成分组数
func (repo *GormRepo[MOD, CLS]) FindC(where func(db *gorm.DB, cls CLS) *gorm.DB, ordering func(cls CLS) gormcnm.Ordering, paging func(db *gorm.DB, cls CLS) *gorm.DB) ([]*MOD, int64, error) {
	db := paging(where(repo.db.Model((*MOD)(nil)), repo.cls).Order(ordering(repo.cls).Ox()), repo.cls).Session(&gorm.Session{})
	var results []*MOD
	if err := db.Find(&results).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	if err := db.Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return results, count, nil
}

// Count returns the number of records matching the where condition
//
// Count 返回符合 where 条件的记录数量
func (repo *GormRepo[MOD, CLS]) Count(where func(db *gorm.DB, cls CLS) *gorm.DB) (int64, error) {
	var count int64
	if err := where(repo.db, repo.cls).Model((*MOD)(nil)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Pluck returns the values of one column across records matching the where condition
// The result type T comes from the column, e.g. cls.Username (ColumnName[string]) gives []string
//
// Pluck 返回符合 where 条件的记录中某一列的值
// 结果类型 T 由列决定，例如 cls.Username（ColumnName[string]）得到 []string
func (repo *GormRepo[MOD, CLS]) Pluck[T any](where func(db *gorm.DB, cls CLS) *gorm.DB, column func(cls CLS) gormcnm.ColumnName[T]) ([]T, error) {
	var values []T
	if err := where(repo.db, repo.cls).Model((*MOD)(nil)).Pluck(column(repo.cls).Name(), &values).Error; err != nil {
		return nil, err
	}
	return values, nil
}

// Scan selects from the MOD table and scans rows into a result struct T
// Fits Select / Joins / Group queries whose rows do not match MOD, e.g. repo.Scan[AccountInfo](where)
//
// Scan 在 MOD 对应的表上查询，并把结果行扫描到结构体 T
// 适合 Select / Joins / Group 等结果行与 MOD 不一致的查询，例如 repo.Scan[AccountInfo](where)
func (repo *GormRepo[MOD, CLS]) Scan[T any](where func(db *gorm.DB, cls CLS) *gorm.DB) ([]T, error) {
	var results []T
	// Rows() 在没有 Dest 时把 Model 当 Dest 用，nil 指针会报 ErrInvalidValue，因此用 new(MOD)
	if err := where(repo.db.Model(new(MOD)), repo.cls).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// ScanList scans result rows and also passes the RES columns to the callback.
// RES needs column definitions; use AsName with result columns to get type-checked aliases.
//
// ScanList 扫描结果列表，并把 RES 的列定义传给回调。
// RES 需要列定义；通过 AsName 引用结果列时，别名具有类型检查。
func (repo *GormRepo[MOD, CLS]) ScanList[RES any, PRES interface {
	*RES
	gormclass.ColumnsInterface[CLSR]
}, CLSR any](where func(db *gorm.DB, cls CLS, res CLSR) *gorm.DB) ([]RES, error) {
	return repo.Scan[RES](func(db *gorm.DB, cls CLS) *gorm.DB {
		return where(db, cls, PRES(new(RES)).Columns())
	})
}

// ScanItem scans just one result row and passes the RES columns to the callback.
// As in ScanList, RES needs column definitions. Zero rows and two rows and up both return an error.
//
// ScanItem 扫描恰好一条结果，并把 RES 的列定义传给回调。
// 与 ScanList 一样，RES 需要列定义。零行或多行均返回错误。
func (repo *GormRepo[MOD, CLS]) ScanItem[RES any, PRES interface {
	*RES
	gormclass.ColumnsInterface[CLSR]
}, CLSR any](where func(db *gorm.DB, cls CLS, res CLSR) *gorm.DB) (RES, error) {
	return repo.ScanOne[RES](func(db *gorm.DB, cls CLS) *gorm.DB {
		return where(db, cls, PRES(new(RES)).Columns())
	})
}

// ScanOne selects one value, such as an aggregate, from the MOD table and returns it
// Fails when the select gives no row, and fails on two rows and up, so a missing value does not come back as a silent zero
// Example: repo.ScanOne[int](func(db, cls) *gorm.DB { return db.Select(cls.Age.COALESCE().Max()) })
//
// ScanOne 在 MOD 对应的表上查询一个值（如聚合结果）并返回
// 查不到行、或查出多行时都报错，查不到的值不会悄悄变成零值
// 示例：repo.ScanOne[int](func(db, cls) *gorm.DB { return db.Select(cls.Age.COALESCE().Max()) })
func (repo *GormRepo[MOD, CLS]) ScanOne[T any](where func(db *gorm.DB, cls CLS) *gorm.DB) (T, error) {
	var zero T
	results, err := repo.Scan[T](where)
	if err != nil {
		return zero, err
	}
	if len(results) != 1 {
		return zero, errors.Errorf("scan one: got %d rows, want 1", len(results))
	}
	return results[0], nil
}

// Update updates a single column for records matching the where condition
// Uses newValue to give the column name and the value
//
// Update 更新符合 where 条件的记录的单个列
// 使用 newValue 指定列名和值
func (repo *GormRepo[MOD, CLS]) Update(where func(db *gorm.DB, cls CLS) *gorm.DB, newValue func(cls CLS) (string, any)) error {
	if err := repo.Gorm().Update(where, newValue).Error; err != nil {
		return err
	}
	return nil
}

// Updates updates multiple columns for records matching the where condition
// Uses newValues to give the column-value pairs
//
// Updates 更新符合 where 条件的记录的多个列
// 使用 newValues 指定列值对
func (repo *GormRepo[MOD, CLS]) Updates(where func(db *gorm.DB, cls CLS) *gorm.DB, newValues func(cls CLS) map[string]any) error {
	if err := repo.Gorm().Updates(where, newValues).Error; err != nil {
		return err
	}
	return nil
}

// UpdatesM updates multiple columns using ColumnValueMap, provides fluent API without AsMap() call
// Usage: repo.UpdatesM(where, func(cls CLS) gormcnm.ColumnValueMap { return cls.Kw(cls.Name.Kv("x"))... })
//
// UpdatesM 使用 ColumnValueMap 更新多列，提供流畅的 API，无需调用 AsMap()
// 用法：repo.UpdatesM(where, func(cls CLS) gormcnm.ColumnValueMap { return cls.Kw(cls.Name.Kv("x"))... })
func (repo *GormRepo[MOD, CLS]) UpdatesM(where func(db *gorm.DB, cls CLS) *gorm.DB, newValues func(cls CLS) gormcnm.ColumnValueMap) error {
	return repo.Updates(where, func(cls CLS) map[string]any {
		return newValues(cls).AsMap()
	})
}

// UpdatesO updates object using primary key as condition, with ColumnValueMap for values
// O = Object, the object must have valid primary key value, GORM uses it to locate the record
//
// UpdatesO 使用主键作为条件更新对象，使用 ColumnValueMap 指定更新值
// O = Object，object 必须有有效的主键值，GORM 会用它来定位要更新的记录
func (repo *GormRepo[MOD, CLS]) UpdatesO(object *MOD, newValues func(cls CLS) gormcnm.ColumnValueMap) error {
	if err := repo.db.Model(object).Updates(newValues(repo.cls).AsMap()).Error; err != nil {
		return err
	}
	return nil
}

// UpdatesC updates object using combined conditions: primary key from object plus where clause
// C = Combined, uses both object primary key and where conditions for precise targeting
//
// UpdatesC 使用组合条件更新对象：object 的主键加上 where 子句
// C = Combined，同时使用 object 主键和 where 条件进行精确定位
func (repo *GormRepo[MOD, CLS]) UpdatesC(object *MOD, where func(db *gorm.DB, cls CLS) *gorm.DB, newValues func(cls CLS) gormcnm.ColumnValueMap) error {
	if err := where(repo.db.Model(object), repo.cls).Updates(newValues(repo.cls).AsMap()).Error; err != nil {
		return err
	}
	return nil
}

// Invoke executes a custom operation using the database connection and column definitions
// Returns the custom operation's resulting error
//
// Invoke 使用数据库连接和列定义执行自定义操作
// 返回自定义操作产生的错误
func (repo *GormRepo[MOD, CLS]) Invoke(run func(db *gorm.DB, cls CLS) *gorm.DB) error {
	if err := run(repo.db, repo.cls).Error; err != nil {
		return err
	}
	return nil
}

// Create inserts a new record into the database
// The record's primary key will be populated after creation
//
// Create 向数据库插入一条新记录
// 创建后记录的主键将被填充
func (repo *GormRepo[MOD, CLS]) Create(one *MOD) error {
	if err := repo.db.Create(one).Error; err != nil {
		return err
	}
	return nil
}

// Creates inserts multiple records into the database in batch
// All records' primary keys will be populated after creation
// Takes the address of the slice, as db.Create(&users) does in GORM: with clause.Returning GORM rewrites the slice
// into the records the database returned, e.g. with OnConflict DoNothing the slice shrinks to the inserted records
//
// Creates 批量向数据库插入多条记录
// 创建后所有记录的主键将被填充
// 参数是切片的地址，跟 GORM 的 db.Create(&users) 一致：搭配 clause.Returning 时 GORM 会把切片改写成数据库返回的记录，
// 例如再配上 OnConflict DoNothing，切片就缩成实际插入的那几条
func (repo *GormRepo[MOD, CLS]) Creates(ones *[]*MOD) error {
	if err := repo.db.Create(ones).Error; err != nil {
		return err
	}
	return nil
}

// CreateInBatches inserts records in batches to reduce database and memory pressure
// Uses batchSize to control records per batch, suitable for large data insertion
// Takes the slice itself, as GORM does: GORM writes each batch back into the same positions and does not change the length,
// so with OnConflict DoNothing plus clause.Returning the skipped records remain in the slice, use Creates when the inserted records are wanted
//
// CreateInBatches 分批插入记录以减少数据库和内存压力
// 使用 batchSize 控制每批的记录数，适合大量数据插入
// 参数是切片本身，跟 GORM 一致：GORM 把每批结果写回原位置，不改切片长度，
// 所以 OnConflict DoNothing 搭配 clause.Returning 时被跳过的记录还留在切片里，要拿到实际插入的记录请用 Creates
func (repo *GormRepo[MOD, CLS]) CreateInBatches(ones []*MOD, batchSize int) error {
	if err := repo.db.CreateInBatches(ones, batchSize).Error; err != nil {
		return err
	}
	return nil
}

// Save inserts or updates a record based on primary key
// If primary key is zero value, creates new record; otherwise updates existing
//
// Save 根据主键插入或更新记录
// 如果主键是零值，创建新记录；否则更新现有记录
func (repo *GormRepo[MOD, CLS]) Save(one *MOD) error {
	if err := repo.db.Save(one).Error; err != nil {
		return err
	}
	return nil
}

// Saves inserts or updates multiple records based on primary keys
// For each record: if primary key is zero value, creates new; otherwise updates existing
// Takes the address of the slice, same as Creates
//
// Saves 根据主键批量插入或更新多条记录
// 对于每条记录：如果主键是零值，创建新记录；否则更新现有记录
// 参数是切片的地址，跟 Creates 一致
func (repo *GormRepo[MOD, CLS]) Saves(ones *[]*MOD) error {
	if err := repo.db.Save(ones).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes the given record using its primary key
// When using GORM Delete, param one cannot be nil, as GORM requires valid instance
//
// Delete 使用主键删除给定记录
// 使用 GORM Delete 时，参数 one 不能为 nil，因为 GORM 需要有效实例
func (repo *GormRepo[MOD, CLS]) Delete(one *MOD) error {
	if err := repo.db.Delete(one).Error; err != nil {
		return err
	}
	return nil
}

// DeleteW deletes records matching the where condition
// W = Where, uses where condition instead of object primary key
//
// DeleteW 删除符合 where 条件的记录
// W = Where，使用 where 条件而不是对象主键
func (repo *GormRepo[MOD, CLS]) DeleteW(where func(db *gorm.DB, cls CLS) *gorm.DB) error {
	// GORM Delete needs valid instance, not nil, otherwise gets ErrInvalidValue
	// GORM Delete 需要有效实例，不能为 nil，否则报 ErrInvalidValue 错误
	if err := where(repo.db, repo.cls).Delete(new(MOD)).Error; err != nil {
		return err
	}
	return nil
}

// DeleteM deletes the given object with additional where condition
// M = Model + Where, combines object with where condition
//
// DeleteM 删除给定对象并附加 where 条件
// M = Model + Where，组合对象和 where 条件
func (repo *GormRepo[MOD, CLS]) DeleteM(one *MOD, where func(db *gorm.DB, cls CLS) *gorm.DB) error {
	if err := where(repo.db, repo.cls).Delete(one).Error; err != nil {
		return err
	}
	return nil
}

// Unscoped disables soft-delete scope and returns a new GormRepo
// Enables permanent deletion via method chaining: repo.Unscoped().DeleteW(...)
//
// Unscoped 禁用软删除作用域并返回新的 GormRepo
// 通过方法链实现永久删除：repo.Unscoped().DeleteW(...)
func (repo *GormRepo[MOD, CLS]) Unscoped() *GormRepo[MOD, CLS] {
	return NewGormRepo[MOD, CLS](repo.db.Unscoped().Session(&gorm.Session{}), (*MOD)(nil), repo.cls)
}

// Clauses adds clauses to the database and returns a new GormRepo
// Enables upsert and other clause-based operations via method chaining
// Example: repo.Clauses(clause.OnConflict{...}).Create(&record)
//
// Clauses 向数据库添加子句并返回新的 GormRepo
// 通过方法链支持 upsert 和其他基于子句的操作
// 示例：repo.Clauses(clause.OnConflict{...}).Create(&record)
func (repo *GormRepo[MOD, CLS]) Clauses(clauses ...clause.Expression) *GormRepo[MOD, CLS] {
	return NewGormRepo[MOD, CLS](repo.db.Clauses(clauses...).Session(&gorm.Session{}), (*MOD)(nil), repo.cls)
}

// Clause adds a clause built from column definitions and returns a new GormRepo
// Enables type-safe clause construction using column names
// Example: repo.Clause(func(cls) clause.Expression { return clause.OnConflict{Columns: []clause.Column{{Name: cls.ID.Name()}}} }).Create(&record)
//
// Clause 使用列定义构建子句并返回新的 GormRepo
// 支持使用列名进行类型安全的子句构建
// 示例：repo.Clause(func(cls) clause.Expression { return clause.OnConflict{Columns: []clause.Column{{Name: cls.ID.Name()}}} }).Create(&record)
func (repo *GormRepo[MOD, CLS]) Clause(clauseFunc func(cls CLS) clause.Expression) *GormRepo[MOD, CLS] {
	return repo.Clauses(clauseFunc(repo.cls))
}

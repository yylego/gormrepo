package gormtablerepo

// Repo holds the table name and table-prefixed column definitions of one GORM struct
// Use it in join queries, where columns need the table name to avoid ambiguous names
//
// Repo 保存一个模型的表名和带表名前缀的列定义
// 用于连表查询，列需要带上表名以避免重名歧义
type Repo[MOD any, CLS any] struct {
	tableName string
	tbColumns CLS
}

// NewRepo creates a Repo with table name and table-prefixed column definitions
// The MOD param is used to deduce the type, its value is not used
//
// NewRepo 使用表名和带表名前缀的列定义创建 Repo
// MOD 参数用于类型推断，其值不使用
func NewRepo[MOD any, CLS any](_ *MOD, tableName string, tbColumns CLS) *Repo[MOD, CLS] {
	return &Repo[MOD, CLS]{
		tableName: tableName,
		tbColumns: tbColumns,
	}
}

// TableName returns the table name, the same method name as in the GORM schema.Tabler
//
// TableName 返回表名，方法名跟 GORM 的 schema.Tabler 一致
func (repo *Repo[MOD, CLS]) TableName() string {
	return repo.tableName
}

// TableColumns returns the table-prefixed column definitions
//
// TableColumns 返回带表名前缀的列定义
func (repo *Repo[MOD, CLS]) TableColumns() CLS {
	return repo.tbColumns
}

// BuildColumns builds a column name list from the table-prefixed column definitions
//
// BuildColumns 使用带表名前缀的列定义构建列名列表
func (repo *Repo[MOD, CLS]) BuildColumns(namesFunc func(cls CLS) []string) []string {
	return namesFunc(repo.tbColumns)
}

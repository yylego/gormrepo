package gormclass

import (
	"github.com/yylego/gormcnm"
	"gorm.io/gorm/schema"
)

// ColumnsInterface is the constraint on models with the generated Columns method
// ColumnsInterface 约束模型带有生成的 Columns 方法
type ColumnsInterface[CLS any] interface {
	Columns() CLS
}

// ColumnsWithTableNameInterface is ColumnsInterface plus the TableName method
// ColumnsWithTableNameInterface 在 ColumnsInterface 之上，再要求有 TableName 方法
type ColumnsWithTableNameInterface[CLS any] interface {
	schema.Tabler
	ColumnsInterface[CLS]
}

// TableColumnsInterface is the constraint on models with the generated TableColumns method, which takes a decoration such as the table name
// TableColumnsInterface 约束模型带有生成的 TableColumns 方法，它接收一个装饰（如表名前缀）
type TableColumnsInterface[CLS any] interface {
	TableColumns(gormcnm.ColumnNameDecoration) CLS
}

// TableColumnsWithTableNameInterface is TableColumnsInterface plus the TableName method
// TableColumnsWithTableNameInterface 在 TableColumnsInterface 之上，再要求有 TableName 方法
type TableColumnsWithTableNameInterface[CLS any] interface {
	schema.Tabler
	TableColumnsInterface[CLS]
}

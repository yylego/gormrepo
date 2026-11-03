// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/gormcngen

//go:build !gormcngen_generate

// Generated from: models.gen_test.go:26 -> models.TestGenerateModels
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

package models

import "github.com/yylego/gormcnm"

func (c *Student) Columns() *StudentColumns {
	return c.TableColumns(gormcnm.NewPlainDecoration())
}

type StudentColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	ID    gormcnm.ColumnName[uint]
	Name  gormcnm.ColumnName[string]
	Age   gormcnm.ColumnName[int]
	Grade gormcnm.ColumnName[string]
}

func (c *Student) TableColumns(decoration gormcnm.ColumnNameDecoration) *StudentColumns {
	return &StudentColumns{
		// Auto-generated: column mapping in table operations. DO NOT EDIT. // 自动生成：表操作的列映射。请勿编辑。
		ID:    gormcnm.Cmn(c.ID, "id", decoration),
		Name:  gormcnm.Cmn(c.Name, "name", decoration),
		Age:   gormcnm.Cmn(c.Age, "age", decoration),
		Grade: gormcnm.Cmn(c.Grade, "grade", decoration),
	}
}

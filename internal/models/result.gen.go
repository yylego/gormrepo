// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/gormcngen

//go:build !gormcngen_generate

// Generated from: result.gen_test.go:27 -> models.TestGenerateResult
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

package models

import "github.com/yylego/gormcnm"

func (c *AccountInfo) Columns() *AccountInfoColumns {
	return &AccountInfoColumns{
		// Auto-generated: column names and types mapping. DO NOT EDIT. // 自动生成：列名和类型映射。请勿编辑。
		LoginName: gormcnm.Cnm(c.LoginName, "login_name"),
		ShowName:  gormcnm.Cnm(c.ShowName, "show_name"),
	}
}

type AccountInfoColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	LoginName gormcnm.ColumnName[string]
	ShowName  gormcnm.ColumnName[string]
}

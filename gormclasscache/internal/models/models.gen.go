// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/gormcngen

//go:build !gormcngen_generate

// Generated from: models.gen_test.go:27 -> models.TestGenerateModels
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

package models

import "github.com/yylego/gormcnm"

func (c *Account) Columns() *AccountColumns {
	return c.TableColumns(gormcnm.NewPlainDecoration())
}

func (c *Account) TableColumns(decoration gormcnm.ColumnNameDecoration) *AccountColumns {
	return &AccountColumns{
		// Auto-generated: column mapping in table operations. DO NOT EDIT. // 自动生成：表操作的列映射。请勿编辑。
		ID:       gormcnm.Cmn(c.ID, "id", decoration),
		Username: gormcnm.Cmn(c.Username, "username", decoration),
		Password: gormcnm.Cmn(c.Password, "password", decoration),
	}
}

type AccountColumns struct {
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	ID       gormcnm.ColumnName[uint]
	Username gormcnm.ColumnName[string]
	Password gormcnm.ColumnName[string]
}

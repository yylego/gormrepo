// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/gormcngen

//go:build !gormcngen_generate

// Generated from: models.gen_test.go:26 -> models.TestGenerateModels
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

package models

import (
	"time"

	"github.com/yylego/gormcnm"
	"gorm.io/gorm"
)

func (c *Account) Columns() *AccountColumns {
	return c.TableColumns(gormcnm.NewPlainDecoration())
}

type AccountColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	ID        gormcnm.ColumnName[uint]
	CreatedAt gormcnm.ColumnName[time.Time]
	UpdatedAt gormcnm.ColumnName[time.Time]
	DeletedAt gormcnm.ColumnName[gorm.DeletedAt]
	Username  gormcnm.ColumnName[string]
	Password  gormcnm.ColumnName[string]
	Nickname  gormcnm.ColumnName[string]
}

func (c *Account) TableColumns(decoration gormcnm.ColumnNameDecoration) *AccountColumns {
	return &AccountColumns{
		// Auto-generated: column mapping in table operations. DO NOT EDIT. // 自动生成：表操作的列映射。请勿编辑。
		ID:        gormcnm.Cmn(c.ID, "id", decoration),
		CreatedAt: gormcnm.Cmn(c.CreatedAt, "created_at", decoration),
		UpdatedAt: gormcnm.Cmn(c.UpdatedAt, "updated_at", decoration),
		DeletedAt: gormcnm.Cmn(c.DeletedAt, "deleted_at", decoration),
		Username:  gormcnm.Cmn(c.Username, "username", decoration),
		Password:  gormcnm.Cmn(c.Password, "password", decoration),
		Nickname:  gormcnm.Cmn(c.Nickname, "nickname", decoration),
	}
}

package gormtablerepo

// The methods in this file turn the table repo into the gormrepo forms: BaseRepo, GormRepo, GormWrap
// Each one keeps the same table-prefixed cls, so the file name wrap means a change of state, not a change of kind
//
// 本文件把表 repo 转成 gormrepo 的几种形态：BaseRepo、GormRepo、GormWrap
// 转换前后用的都是同一份带表名前缀的列定义，所以文件名 wrap 表示的是状态变化，不是类型变化

import (
	"github.com/yylego/gormrepo"
	"gorm.io/gorm"
)

// Base converts to a BaseRepo that uses the table-prefixed column definitions
//
// Base 转换成使用带表名前缀列定义的 BaseRepo
func (repo *Repo[MOD, CLS]) Base() *gormrepo.BaseRepo[MOD, CLS] {
	return gormrepo.NewBaseRepo((*MOD)(nil), repo.tbColumns)
}

// Repo creates a GormRepo with the db connection
//
// Repo 使用数据库连接创建 GormRepo
func (repo *Repo[MOD, CLS]) Repo(db *gorm.DB) *gormrepo.GormRepo[MOD, CLS] {
	return repo.Base().Repo(db)
}

// Gorm creates a GormWrap with the db connection
//
// Gorm 使用数据库连接创建 GormWrap
func (repo *Repo[MOD, CLS]) Gorm(db *gorm.DB) *gormrepo.GormWrap[MOD, CLS] {
	return repo.Base().Gorm(db)
}

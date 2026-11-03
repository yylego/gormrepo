package gormrepo

// The methods in this file turn the repo from one form to the next: GormRepo <-> GormWrap, Modeled, WithContext
// Each one keeps the same db and cls, so the file name wrap means a change of state, not a change of kind
//
// 本文件放的是形态之间的转换：GormRepo 与 GormWrap 互转，以及 Modeled、WithContext
// 转换前后用的都是同一份 db 和 cls，所以文件名 wrap 表示的是状态变化，不是类型变化

import (
	"context"

	"gorm.io/gorm"
)

// Gorm converts a GormRepo to a GormWrap, sharing the same DB and CLS instances
// Returns a new GormWrap instance enabling chainable operations
//
// Gorm 将 GormRepo 转换为 GormWrap，共享相同的 DB 和 CLS 实例
// 返回新的 GormWrap 实例以支持链式操作
func (repo *GormRepo[MOD, CLS]) Gorm() *GormWrap[MOD, CLS] {
	return NewGormWrap(repo.db, (*MOD)(nil), repo.cls)
}

// Repo converts a GormWrap to a GormRepo, sharing the same DB and CLS instances
// Returns a new GormRepo instance enabling chainable operations
//
// Repo 将 GormWrap 转换为 GormRepo，共享相同的 DB 和 CLS 实例
// 返回新的 GormRepo 实例以支持链式操作
func (wrap *GormWrap[MOD, CLS]) Repo() *GormRepo[MOD, CLS] {
	return NewGormRepo(wrap.db, (*MOD)(nil), wrap.cls)
}

// Modeled returns a new GormRepo that has db.Model(new(MOD)) applied
// It takes no param since the repo knows its MOD, and pairs with Unscoped as a state of the repo
//
// Modeled 返回一个设好 MOD 作为 GORM 模型的新 GormRepo，等同于 db.Model(new(MOD))
// repo 已经知道自己的 MOD，所以不用传参数；它跟 Unscoped 一样表示 repo 的一种状态
func (repo *GormRepo[MOD, CLS]) Modeled() *GormRepo[MOD, CLS] {
	return NewGormRepo(repo.db.Model((*MOD)(nil)).Session(&gorm.Session{}), (*MOD)(nil), repo.cls)
}

// Modeled returns a new GormWrap that has db.Model(new(MOD)) applied
// It takes no param since the wrap knows its MOD, and pairs with Unscoped as a state of the wrap
//
// Modeled 返回一个设好 MOD 作为 GORM 模型的新 GormWrap，等同于 db.Model(new(MOD))
// wrap 已经知道自己的 MOD，所以不用传参数；它跟 Unscoped 一样表示 wrap 的一种状态
func (wrap *GormWrap[MOD, CLS]) Modeled() *GormWrap[MOD, CLS] {
	return NewGormWrap(wrap.db.Model((*MOD)(nil)).Session(&gorm.Session{}), (*MOD)(nil), wrap.cls)
}

// WithContext sets the context on the GormRepo DB instance
// Returns a new GormRepo instance enabling chainable operations
//
// WithContext 在 GormRepo 的 DB 实例上设置上下文
// 返回新的 GormRepo 实例以支持链式操作
func (repo *GormRepo[MOD, CLS]) WithContext(ctx context.Context) *GormRepo[MOD, CLS] {
	return NewGormRepo(repo.db.WithContext(ctx), (*MOD)(nil), repo.cls)
}

// WithContext sets the context on the GormWrap DB instance
// Returns a new GormWrap instance enabling chainable operations
//
// WithContext 在 GormWrap 的 DB 实例上设置上下文
// 返回新的 GormWrap 实例以支持链式操作
func (wrap *GormWrap[MOD, CLS]) WithContext(ctx context.Context) *GormWrap[MOD, CLS] {
	return NewGormWrap(wrap.db.WithContext(ctx), (*MOD)(nil), wrap.cls)
}

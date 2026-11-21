package gormjoin

import (
	"strings"

	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/gormrepo/gormtablerepo"
	"gorm.io/gorm/clause"
)

// Join is a join of two table repos with the ON condition set, it does the join and nothing else
// Select / Where / Scan remain with GORM and the repos, with *gorm.DB as the link
//
// Join 是设置好连接条件的两表连接，只负责连表
// Select / Where / Scan 仍由 GORM 和各 repo 去写，以 *gorm.DB 为纽带
type Join[MOD1 any, CLS1 any, MOD2 any, CLS2 any] struct {
	repo1        *gormtablerepo.Repo[MOD1, CLS1]
	joinType     clause.JoinType
	repo2        *gormtablerepo.Repo[MOD2, CLS2]
	onConditions []string
}

// New joins repo2 onto repo1 with the given join type, each type is inferred from the repos
// The columns passed to onFunc come with the table name, multiple conditions are joined with AND
// A CROSS JOIN takes no ON condition, use CrossJoin on it
//
// New 用指定的连接类型把 repo2 连到 repo1 上，类型都由两个 repo 推断
// 传给 onFunc 的列都带表名，多个条件用 AND 连接
// CROSS JOIN 不带连接条件，请用 CrossJoin
func New[MOD1 any, CLS1 any, MOD2 any, CLS2 any](repo1 *gormtablerepo.Repo[MOD1, CLS1], joinType clause.JoinType, repo2 *gormtablerepo.Repo[MOD2, CLS2], onFunc func(cls1 CLS1, cls2 CLS2) []string) *Join[MOD1, CLS1, MOD2, CLS2] {
	return &Join[MOD1, CLS1, MOD2, CLS2]{
		repo1:        repo1,
		joinType:     joinType,
		repo2:        repo2,
		onConditions: onFunc(repo1.TableColumns(), repo2.TableColumns()),
	}
}

// LeftJoin joins repo2 onto repo1 with LEFT JOIN, e.g. gormjoin.LeftJoin(studentRepo, scoreRepo, onFunc)
// LeftJoin 用 LEFT JOIN 把 repo2 连到 repo1 上，例如 gormjoin.LeftJoin(studentRepo, scoreRepo, onFunc)
func LeftJoin[MOD1 any, CLS1 any, MOD2 any, CLS2 any](repo1 *gormtablerepo.Repo[MOD1, CLS1], repo2 *gormtablerepo.Repo[MOD2, CLS2], onFunc func(cls1 CLS1, cls2 CLS2) []string) *Join[MOD1, CLS1, MOD2, CLS2] {
	return New(repo1, clause.LeftJoin, repo2, onFunc)
}

// RightJoin joins repo2 onto repo1 with RIGHT JOIN
// RightJoin 用 RIGHT JOIN 把 repo2 连到 repo1 上
func RightJoin[MOD1 any, CLS1 any, MOD2 any, CLS2 any](repo1 *gormtablerepo.Repo[MOD1, CLS1], repo2 *gormtablerepo.Repo[MOD2, CLS2], onFunc func(cls1 CLS1, cls2 CLS2) []string) *Join[MOD1, CLS1, MOD2, CLS2] {
	return New(repo1, clause.RightJoin, repo2, onFunc)
}

// InnerJoin joins repo2 onto repo1 with INNER JOIN
// InnerJoin 用 INNER JOIN 把 repo2 连到 repo1 上
func InnerJoin[MOD1 any, CLS1 any, MOD2 any, CLS2 any](repo1 *gormtablerepo.Repo[MOD1, CLS1], repo2 *gormtablerepo.Repo[MOD2, CLS2], onFunc func(cls1 CLS1, cls2 CLS2) []string) *Join[MOD1, CLS1, MOD2, CLS2] {
	return New(repo1, clause.InnerJoin, repo2, onFunc)
}

// CrossJoin joins repo2 onto repo1 with CROSS JOIN, each row of repo1 pairs with each row of repo2
// It takes no onFunc: standard SQL has no ON condition on a CROSS JOIN, and PostgreSQL rejects CROSS JOIN ... ON
//
// CrossJoin 用 CROSS JOIN 把 repo2 连到 repo1 上，repo1 的每一行跟 repo2 的每一行各配一次
// 它不收 onFunc：标准 SQL 的 CROSS JOIN 不带 ON 条件，PostgreSQL 不认 CROSS JOIN ... ON
func CrossJoin[MOD1 any, CLS1 any, MOD2 any, CLS2 any](repo1 *gormtablerepo.Repo[MOD1, CLS1], repo2 *gormtablerepo.Repo[MOD2, CLS2]) *Join[MOD1, CLS1, MOD2, CLS2] {
	return &Join[MOD1, CLS1, MOD2, CLS2]{
		repo1:    repo1,
		joinType: clause.CrossJoin,
		repo2:    repo2,
	}
}

// Stmt returns the JOIN statement to pass to db.Joins, e.g. "LEFT JOIN scores ON scores.student_id=students.id"
// With no ON condition the statement has no ON part, e.g. "CROSS JOIN scores"
//
// Stmt 返回传给 db.Joins 的 JOIN 语句，例如 "LEFT JOIN scores ON scores.student_id=students.id"
// 没有连接条件时不带 ON 部分，例如 "CROSS JOIN scores"
func (jo *Join[MOD1, CLS1, MOD2, CLS2]) Stmt() string {
	stmt := string(jo.joinType) + " JOIN " + jo.repo2.TableName()
	if len(jo.onConditions) == 0 {
		return stmt
	}
	return stmt + " ON " + strings.Join(jo.onConditions, " AND ")
}

// Names returns the column names to select from the two tables, aliased to the RES columns with type checks
// Pass it to db.Select, joins of three tables and up concat each part, e.g. db.Select(slices.Concat(join1.Names[RES](f1), join2.Names[RES](f2)))
//
// Names 返回两张表要选的列名，用 RES 的列起别名，带类型检查
// 把它传给 db.Select，三张及以上的表把每段拼起来，例如 db.Select(slices.Concat(join1.Names[RES](f1), join2.Names[RES](f2)))
func (jo *Join[MOD1, CLS1, MOD2, CLS2]) Names[RES any, PRES interface {
	*RES
	gormclass.ColumnsInterface[CLSR]
}, CLSR any](selectFunc func(cls1 CLS1, cls2 CLS2, res CLSR) []string) []string {
	return selectFunc(jo.repo1.TableColumns(), jo.repo2.TableColumns(), PRES(new(RES)).Columns())
}

// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/gormcngen

//go:build !gormcngen_generate

// Generated from: result.gen_test.go:27 -> models.TestGenerateResult
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

package models

import "github.com/yylego/gormcnm"

func (c *SubjectStat) Columns() *SubjectStatColumns {
	return &SubjectStatColumns{
		// Auto-generated: column names and types mapping. DO NOT EDIT. // 自动生成：列名和类型映射。请勿编辑。
		Subject:    gormcnm.Cnm(c.Subject, "subject"),
		ScoreCount: gormcnm.Cnm(c.ScoreCount, "score_count"),
		PointsSum:  gormcnm.Cnm(c.PointsSum, "points_sum"),
		PointsHigh: gormcnm.Cnm(c.PointsHigh, "points_high"),
		PointsAvg:  gormcnm.Cnm(c.PointsAvg, "points_avg"),
	}
}

type SubjectStatColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	Subject    gormcnm.ColumnName[string]
	ScoreCount gormcnm.ColumnName[int64]
	PointsSum  gormcnm.ColumnName[int]
	PointsHigh gormcnm.ColumnName[int]
	PointsAvg  gormcnm.ColumnName[float64]
}

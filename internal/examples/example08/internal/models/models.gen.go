// Code generated using gormcngen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/gormcngen

//go:build !gormcngen_generate

// Generated from: models.gen_test.go:27 -> models.TestGenerateModels
// ========== GORMCNGEN:DO-NOT-EDIT-MARKER:END ==========

package models

import "github.com/yylego/gormcnm"

func (c *Course) Columns() *CourseColumns {
	return c.TableColumns(gormcnm.NewPlainDecoration())
}

type CourseColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	ID       gormcnm.ColumnName[uint]
	Name     gormcnm.ColumnName[string]
	Seats    gormcnm.ColumnName[int]
	Enrolled gormcnm.ColumnName[int]
}

func (c *Enrollment) Columns() *EnrollmentColumns {
	return c.TableColumns(gormcnm.NewPlainDecoration())
}

type EnrollmentColumns struct {
	// Auto-generated: embedding operation functions to make it simple to use. DO NOT EDIT. // 自动生成：嵌入操作函数便于使用。请勿编辑。
	gormcnm.ColumnOperationClass
	// Auto-generated: column names and types in database table. DO NOT EDIT. // 自动生成：数据库表的列名和类型。请勿编辑。
	ID          gormcnm.ColumnName[uint]
	CourseID    gormcnm.ColumnName[uint]
	StudentName gormcnm.ColumnName[string]
}

func (c *Course) TableColumns(decoration gormcnm.ColumnNameDecoration) *CourseColumns {
	return &CourseColumns{
		// Auto-generated: column mapping in table operations. DO NOT EDIT. // 自动生成：表操作的列映射。请勿编辑。
		ID:       gormcnm.Cmn(c.ID, "id", decoration),
		Name:     gormcnm.Cmn(c.Name, "name", decoration),
		Seats:    gormcnm.Cmn(c.Seats, "seats", decoration),
		Enrolled: gormcnm.Cmn(c.Enrolled, "enrolled", decoration),
	}
}

func (c *Enrollment) TableColumns(decoration gormcnm.ColumnNameDecoration) *EnrollmentColumns {
	return &EnrollmentColumns{
		// Auto-generated: column mapping in table operations. DO NOT EDIT. // 自动生成：表操作的列映射。请勿编辑。
		ID:          gormcnm.Cmn(c.ID, "id", decoration),
		CourseID:    gormcnm.Cmn(c.CourseID, "course_id", decoration),
		StudentName: gormcnm.Cmn(c.StudentName, "student_name", decoration),
	}
}

package models

import "gorm.io/gorm"

// Student is a row of the students table, the gorm.DeletedAt field turns on soft delete
// Student 是 students 表的一行，gorm.DeletedAt 字段开启软删除
type Student struct {
	ID        uint
	Name      string `gorm:"unique"`
	Grade     string
	DeletedAt gorm.DeletedAt
}

func (*Student) TableName() string {
	return "students"
}

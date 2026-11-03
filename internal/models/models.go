package models

import "gorm.io/gorm"

// Account is the test struct with common fields
// Account 是带有常用字段的测试结构体
type Account struct {
	gorm.Model
	Username string `gorm:"unique;"`
	Password string `gorm:"size:255;"`
	Nickname string `gorm:"column:nickname;"`
}

// TableName returns the database table name
// TableName 返回数据库表名
func (*Account) TableName() string {
	return "accounts"
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Username string `gorm:"unique;"`
	Password string `gorm:"size:255;"`
	Nickname string `gorm:"column:nickname;"`
}

func (*Account) TableName() string {
	return "accounts"
}

type Example struct {
	ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name      string    `gorm:"unique;"`
	Age       int       `gorm:"type:int32;index:idx_example_age;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (*Example) TableName() string {
	return "examples"
}

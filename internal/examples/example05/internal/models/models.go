package models

// Student is a row of the students table
// Student 是 students 表的一行
type Student struct {
	ID    uint
	Name  string `gorm:"unique"`
	Grade string
	Score int
}

func (*Student) TableName() string {
	return "students"
}

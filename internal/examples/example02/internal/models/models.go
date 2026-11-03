package models

// Student is a row of the students table, Coach is NULL when the student has no coach
// Student 是 students 表的一行，没有教练时 Coach 为 NULL
type Student struct {
	ID    uint
	Name  string `gorm:"unique"`
	Age   int
	Grade string
	Coach *string
}

func (*Student) TableName() string {
	return "students"
}

package models

// Class is a row of the classes table
// Class 是 classes 表的一行
type Class struct {
	ID   uint
	Room string
}

func (*Class) TableName() string {
	return "classes"
}

// Student is a row of the students table, ClassID is zero when the student is not in a class yet
// Student 是 students 表的一行，还没分班时 ClassID 为零
type Student struct {
	ID      uint
	Name    string
	ClassID uint
}

func (*Student) TableName() string {
	return "students"
}

// StudentWithClass carries the Class association that Preload needs, the table remains students
// StudentWithClass 带上 Preload 需要的 Class 关联字段，表还是 students
type StudentWithClass struct {
	Student
	Class *Class `gorm:"foreignKey:ClassID"`
}

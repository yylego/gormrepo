package models

type Student struct {
	ID   uint
	Name string
}

func (*Student) TableName() string {
	return "students"
}

type Score struct {
	ID        uint
	StudentID uint
	Subject   string
	Points    int
}

func (*Score) TableName() string {
	return "scores"
}

type Course struct {
	ID        uint
	Subject   string
	CoachName string
}

func (*Course) TableName() string {
	return "courses"
}

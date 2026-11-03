package models

type Student struct {
	ID   uint
	Name string
	Rank int
}

func (*Student) TableName() string {
	return "students"
}

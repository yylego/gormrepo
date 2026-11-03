package models

// Score is a row of the scores table, one score of a student in one subject
// Score 是 scores 表的一行，一个学生在一门科目上的成绩
type Score struct {
	ID          uint
	StudentName string
	Subject     string
	Points      int
}

func (*Score) TableName() string {
	return "scores"
}

package models

// StudentScore is the result row of a find that picks the name and the rank of a student as the score
// StudentScore 是一行查询结果：学生的名字，以及作为分数的排名
type StudentScore struct {
	Name  string
	Score int
}

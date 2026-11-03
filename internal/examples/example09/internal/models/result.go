package models

// StudentScore is the result row of the students LEFT JOIN scores SELECT, it is not a table
// StudentScore 是 students LEFT JOIN scores 查询的结果行，不是表
type StudentScore struct {
	StudentName string
	Subject     string
	Points      int
}

// StudentCourse is the result row of the students, scores and courses SELECT, it is not a table
// StudentCourse 是 students、scores、courses 三表查询的结果行，不是表
type StudentCourse struct {
	StudentName string
	Subject     string
	Points      int
	CoachName   string
}

package models

// StudentClass is the result row of each student with the room of its class, it is not a table
// StudentClass 是每个学生和所在班级教室的结果行，不是表
type StudentClass struct {
	StudentName string
	ClassRoom   string
}

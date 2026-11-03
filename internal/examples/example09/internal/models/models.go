package models

// Student is a row of the students table
// Student 是 students 表的一行
type Student struct {
	ID   uint
	Name string
}

func (*Student) TableName() string {
	return "students"
}

// Score is a row of the scores table, one score of a student in one subject
// Score 是 scores 表的一行，一个学生在一门科目上的成绩
type Score struct {
	ID        uint
	StudentID uint
	Subject   string
	Points    int
}

func (*Score) TableName() string {
	return "scores"
}

// Course is a row of the courses table, the coach of each subject
// Course 是 courses 表的一行，每门科目的任课教练
type Course struct {
	ID        uint
	Subject   string
	CoachName string
}

func (*Course) TableName() string {
	return "courses"
}

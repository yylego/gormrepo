package models

// Course is a row of the courses table, Seats is the most students it takes, Enrolled counts the students enrolled in it
// Course 是 courses 表的一行，Seats 是最多能选的人数，Enrolled 是已选这门课的人数
type Course struct {
	ID       uint
	Name     string `gorm:"unique"`
	Seats    int
	Enrolled int
}

func (*Course) TableName() string {
	return "courses"
}

// Enrollment is a row of the enrollments table, one student enrolled in one course
// Enrollment 是 enrollments 表的一行，一个学生选了一门课
type Enrollment struct {
	ID          uint
	CourseID    uint
	StudentName string
}

func (*Enrollment) TableName() string {
	return "enrollments"
}

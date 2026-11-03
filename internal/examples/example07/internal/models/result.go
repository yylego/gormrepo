package models

// SubjectStat is a result row of the scores grouped on subject, it is not a table
// SubjectStat 是按科目分组统计成绩的结果行，不是表
type SubjectStat struct {
	Subject    string
	ScoreCount int64
	PointsSum  int
	PointsHigh int
	PointsAvg  float64
}

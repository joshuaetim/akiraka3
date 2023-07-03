package model

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	Questions string `json:"questions"`
	Options   string `json:"options"`
	Answers   string `json:"answers"`
	Course    uint   `json:"course"`
	Staff     uint   `json:"staff"`
	Title     string `json:"title"`
	Duration  int    `json:"duration"`
}

func (q *Quiz) PublicQuiz() *Quiz {
	quiz := &Quiz{
		Questions: q.Questions,
		Options:   q.Options,
		// Answers:   q.Answers,
		Course:   q.Course,
		Staff:    q.Staff,
		Title:    q.Title,
		Duration: q.Duration,
	}
	quiz.ID = q.ID
	quiz.CreatedAt = q.CreatedAt
	return quiz
}

func (*Quiz) PublicQuizArray(quizzes []Quiz) []Quiz {
	var res []Quiz
	for _, q := range quizzes {
		res = append(res, *q.PublicQuiz())
	}
	return res
}

// table name
func (Quiz) TableName() string {
	return "quizzes"
}

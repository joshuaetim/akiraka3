package model

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	Questions string `json:"questions"`
	Options   string `json:"options"`
	Answers   string `json:"answers"`
	Course    uint   `json:"course"`
	Staff     uint   `json:"staff"`
}

func (q *Quiz) PublicQuiz() *Quiz {
	quiz := &Quiz{
		Questions: q.Questions,
		Options:   q.Options,
		Answers:   q.Answers,
		Course:    q.Course,
		Staff:     q.Staff,
	}
	quiz.ID = q.ID
	quiz.CreatedAt = q.CreatedAt
	return quiz
}

// table name
func (Quiz) TableName() string {
	return "quizzes"
}

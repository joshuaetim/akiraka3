package model

import "gorm.io/gorm"

type Grade struct {
	gorm.Model
	Quiz  uint    `json:"quiz"`
	Score float64 `json:"score"`
	User  uint    `json:"user"`
}

func (g *Grade) PublicGrade() *Grade {
	grade := &Grade{
		Quiz:  g.Quiz,
		Score: g.Score,
		User:  g.User,
	}
	grade.ID = g.ID
	grade.CreatedAt = g.CreatedAt
	return grade
}

// table name
func (Grade) TableName() string {
	return "grades"
}

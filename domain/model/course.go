package model

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Staff       uint   `json:"staff"`
}

func (c *Course) PublicCourse() *Course {
	course := &Course{
		Title:       c.Title,
		Description: c.Description,
		Staff:       c.Staff,
	}
	course.ID = c.ID
	course.CreatedAt = c.CreatedAt
	return course
}

// table name
func (Course) TableName() string {
	return "courses"
}

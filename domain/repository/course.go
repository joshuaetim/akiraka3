package repository

import "github.com/joshuaetim/akiraka3/domain/model"

type CourseRepository interface {
	AddCourse(model.Course) (model.Course, error)
	GetCourse(uint) (model.Course, error)
	GetAllCoursesByStaff(uint) ([]model.Course, error)
	GetAllCourses() ([]model.Course, error)
	UpdateCourse(model.Course) (model.Course, error)
	DeleteCourse(model.Course) error
}

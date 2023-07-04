package repository

import "github.com/joshuaetim/frontdesk/domain/model"

type GradeRepository interface {
	AddGrade(model.Grade) error
	GetGrade(uint) (model.Grade, error)
	GetAllGrades() ([]model.Grade, error)
	GetAllByQuiz(uint) ([]model.Grade, error)
	GetByUser(uint) ([]model.Grade, error)
	GetGradesByMap(map[string]interface{}) ([]model.Grade, error)

	GetByUserAndQuiz(uint, uint) (model.Grade, error)
	UpdateGrade(model.Grade) (model.Grade, error)
	DeleteGrade(model.Grade) error
}

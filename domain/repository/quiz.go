package repository

import "github.com/joshuaetim/frontdesk/domain/model"

type QuizRepository interface {
	AddQuiz(model.Quiz) (model.Quiz, error)
	GetQuiz(uint) (model.Quiz, error)
	// get quizzes by a particualar staff id
	// GetQuizzesByStaff(uint) ([]model.Quiz, error)
	// GetQuizzesByCourse(uint) ([]model.Quiz, error)
	// get all quizzes
	GetAllQuizzes() ([]model.Quiz, error)
	UpdateQuiz(model.Quiz) (model.Quiz, error)
	DeleteQuiz(model.Quiz) error
	GetQuizzesByMap(map[string]interface{}) ([]model.Quiz, error)
}

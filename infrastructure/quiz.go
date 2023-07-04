package infrastructure

import (
	"fmt"

	"github.com/joshuaetim/frontdesk/domain/model"
	"github.com/joshuaetim/frontdesk/domain/repository"
	"gorm.io/gorm"
)

type quizRepo struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) repository.QuizRepository {
	return &quizRepo{
		db: db,
	}
}

func (r *quizRepo) AddQuiz(quiz model.Quiz) (model.Quiz, error) {
	return quiz, r.db.Create(&quiz).Error
}

func (r *quizRepo) GetQuiz(id uint) (model.Quiz, error) {
	var quiz model.Quiz
	return quiz, r.db.First(&quiz, id).Error
}

// {"course": 1, "staff": 2}
func (r *quizRepo) GetQuizzesByMap(query map[string]interface{}) ([]model.Quiz, error) {
	var queryString string
	var fields []interface{}
	var quiz []model.Quiz
	for k, v := range query {
		if queryString != "" {
			queryString = " " + queryString + " AND "
		}
		queryString = fmt.Sprintf("%s%s = ?", queryString, k)
		fields = append(fields, v)
	}
	// fields[0]
	var queryMain []interface{}
	queryMain = append(queryMain, queryString)
	queryMain = append(queryMain, fields...)

	return quiz, r.db.Find(&quiz, queryMain...).Error
}

func (r *quizRepo) GetQuizzesByStaff(staffid uint) ([]model.Quiz, error) {
	var quiz []model.Quiz
	return quiz, r.db.Find(&quiz, "staff = ?", staffid).Error
}

func (r *quizRepo) GetAllQuizzes() ([]model.Quiz, error) {
	var quizzes []model.Quiz
	return quizzes, r.db.Find(&quizzes).Error
}

func (r *quizRepo) UpdateQuiz(quiz model.Quiz) (model.Quiz, error) {
	return quiz, r.db.Model(&quiz).Updates(&quiz).Error
}

func (q *quizRepo) DeleteQuiz(quiz model.Quiz) error {
	return q.db.Unscoped().Delete(&quiz).Error
}

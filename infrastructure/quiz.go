package infrastructure

import (
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

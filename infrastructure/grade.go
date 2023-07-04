package infrastructure

import (
	"fmt"

	"github.com/joshuaetim/frontdesk/domain/model"
	"github.com/joshuaetim/frontdesk/domain/repository"
	"gorm.io/gorm"
)

type gradeRepo struct {
	db *gorm.DB
}

func NewGradeRepository(db *gorm.DB) repository.GradeRepository {
	return &gradeRepo{
		db: db,
	}
}

func (r *gradeRepo) AddGrade(grade model.Grade) error {
	return r.db.Create(&grade).Error
}

func (r *gradeRepo) GetGrade(id uint) (model.Grade, error) {
	var grade model.Grade
	return grade, r.db.First(&grade, id).Error
}

func (r *gradeRepo) GetAllGrades() ([]model.Grade, error) {
	var grades []model.Grade
	return grades, r.db.Find(&grades).Error
}

func (r *gradeRepo) GetAllByQuiz(quizId uint) ([]model.Grade, error) {
	var grades []model.Grade
	return grades, r.db.Find(&grades, "quiz = ?", quizId).Error
}

func (r *gradeRepo) GetByUserAndQuiz(user uint, quiz uint) (model.Grade, error) {
	var grade model.Grade
	return grade, r.db.First(&grade, "quiz = ? AND user = ?", quiz, user).Error
}

func (r *gradeRepo) UpdateGrade(grade model.Grade) (model.Grade, error) {
	return grade, r.db.Model(&grade).Updates(&grade).Error
}

func (r *gradeRepo) DeleteGrade(grade model.Grade) error {
	return r.db.Delete(&grade).Error
}

func (r *gradeRepo) GetByUser(id uint) ([]model.Grade, error) {
	var grades []model.Grade
	return grades, r.db.Find(&grades, "user = ?", id).Error
}

func (r *gradeRepo) GetGradesByMap(query map[string]interface{}) ([]model.Grade, error) {
	var queryString string
	var fields []interface{}
	var grade []model.Grade
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

	return grade, r.db.Find(&grade, queryMain...).Error
}

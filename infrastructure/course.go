package infrastructure

import (
	"github.com/joshuaetim/akiraka3/domain/model"
	"github.com/joshuaetim/akiraka3/domain/repository"
	"gorm.io/gorm"
)

type courseRepo struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) repository.CourseRepository {
	return &courseRepo{
		db: db,
	}
}

func (cr *courseRepo) AddCourse(course model.Course) (model.Course, error) {
	return course, cr.db.Create(&course).Error
}

func (r *courseRepo) GetCourse(id uint) (model.Course, error) {
	var course model.Course
	return course, r.db.First(&course, id).Error
}

func (r *courseRepo) GetAllCoursesByStaff(staffid uint) ([]model.Course, error) {
	var courses []model.Course
	return courses, r.db.Find(&courses, "staff = ?", staffid).Error
}

func (r *courseRepo) GetAllCourses() ([]model.Course, error) {
	var courses []model.Course
	return courses, r.db.Find(&courses).Error
}

func (r *courseRepo) UpdateCourse(course model.Course) (model.Course, error) {
	return course, r.db.Model(&course).Updates(&course).Error
}

func (r *courseRepo) DeleteCourse(course model.Course) error {
	return r.db.Delete(&course).Error
}

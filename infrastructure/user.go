package infrastructure

import (
	"github.com/joshuaetim/frontdesk/domain/model"
	"github.com/joshuaetim/frontdesk/domain/repository"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepo{
		db: db,
	}
}

var _ repository.UserRepository = &userRepo{}

func (r *userRepo) AddUser(user model.User) (model.User, error) {
	return user, r.db.Create(&user).Error
}

func (r *userRepo) GetUser(id uint) (model.User, error) {
	var user model.User
	return user, r.db.First(&user, id).Error
}

func (r *userRepo) GetByEmail(email string) (model.User, error) {
	var user model.User
	return user, r.db.First(&user, "email = ?", email).Error
}

func (r *userRepo) GetAllUser() ([]model.User, error) {
	var users []model.User
	return users, r.db.Find(&users).Error
}

func (r *userRepo) CountUsers() int {
	type Result struct {
		Total int
	}
	var result Result
	r.db.Raw("select count(*) as total from users").Scan(&result)

	return result.Total
}

func (r *userRepo) UpdateUser(user model.User) (model.User, error) {
	// check if exists
	var uModel model.User = user
	if err := r.db.First(&uModel).Error; err != nil {
		return user, err
	}
	err := r.db.Model(&user).Updates(&user).Error
	user, _ = r.GetUser(user.ID)
	return user, err
}

func (r *userRepo) DeleteUser(user model.User) error {
	// exists?
	if err := r.db.First(&user).Error; err != nil {
		return err
	}
	return r.db.Delete(&user).Error
}

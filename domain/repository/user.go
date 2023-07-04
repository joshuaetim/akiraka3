package repository

import "github.com/joshuaetim/akiraka3/domain/model"

type UserRepository interface {
	AddUser(model.User) (model.User, error)
	GetUser(uint) (model.User, error)
	GetByEmail(string) (model.User, error)
	GetAllUser() ([]model.User, error)
	UpdateUser(model.User) (model.User, error)
	DeleteUser(model.User) error
	CountUsers() int
}

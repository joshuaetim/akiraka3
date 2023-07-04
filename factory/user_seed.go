package factory

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/joshuaetim/akiraka3/domain/model"
	"github.com/joshuaetim/akiraka3/infrastructure"
	"gorm.io/gorm"
)

func SeedUser(dbInstance *gorm.DB) (model.User, error) {
	user := model.User{
		Firstname: gofakeit.FirstName(),
		Lastname:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  "$2a$10$VbEAUZR5q.M88TtfaA0ghuYDPS.qFlim/R51pSN4mFAQVdCW4jmtO", // "password"
	}

	ur := infrastructure.NewUserRepository(dbInstance)
	u, err := ur.AddUser(user)

	return u, err
}

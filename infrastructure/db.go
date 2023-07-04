package infrastructure

import (
	"fmt"
	"log"
	"os"

	"github.com/joshuaetim/akiraka3/domain/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DB() *gorm.DB {
	var dialect gorm.Dialector
	var dsn string

	switch os.Getenv("DB_DRIVER") {
	case "sqlite":
		dsn = os.Getenv("DATABASE_URL")
		if mem := os.Getenv("SQLITE_MEMORY"); mem != "" {
			dialect = sqlite.Open(mem)
		} else {
			dialect = sqlite.Open(dsn)
		}
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_DATABASE"))
		dialect = mysql.Open(dsn)
	default:
		log.Fatalf("invalid driver: %s", os.Getenv("DB_DRIVER"))
	}

	db, err := gorm.Open(dialect, &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalf("Error connecting to database (%v)(%v): %v", dialect.Name(), dsn, err)
		return nil
	}

	db.AutoMigrate(&model.User{}, &model.Course{}, &model.Quiz{}, &model.Grade{})
	return db
}

package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection() *gorm.DB {
	db, err := gorm.Open(postgres.Open(Settings.APP_DATABASE_DSN), &gorm.Config{})
	if err != nil {
		panic(errors.Join(fmt.Errorf("error when connecting to database"), err))
	}

	return db
}

var DatabaseConnection *gorm.DB = GetConnection()

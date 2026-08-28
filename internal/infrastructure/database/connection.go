package database

import (
	"chats-service/configs"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection() *gorm.DB {
	configuration := configs.GetDatabaseConfiguration()
	db, err := gorm.Open(postgres.Open(configuration.DSN.String()), &gorm.Config{})
	if err != nil {
		panic(errors.Join(fmt.Errorf("error when connecting to database"), err))
	}

	return db
}

var DatabaseConnection *gorm.DB = GetConnection()

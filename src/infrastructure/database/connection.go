package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabaseDsn() string {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		Settings.APP_DATABASE_HOST,
		Settings.APP_DATABASE_USER,
		Settings.APP_DATABASE_PASSWORD,
		Settings.APP_DATABASE_NAME,
		Settings.APP_DATABASE_PORT,
	)

	return dsn
}

func GetConnection() *gorm.DB {
	dsn := GetDatabaseDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(errors.Join(fmt.Errorf("error when connecting to database"), err))
	}

	return db
}

var DatabaseConnection *gorm.DB = GetConnection()

package database

import (
	"fmt"

	"github.com/chack-check/chats-service/settings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabaseDsn() string {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		settings.Settings.APP_DATABASE_HOST,
		settings.Settings.APP_DATABASE_USER,
		settings.Settings.APP_DATABASE_PASSWORD,
		settings.Settings.APP_DATABASE_NAME,
		settings.Settings.APP_DATABASE_PORT,
	)

	return dsn
}

func GetConnection() *gorm.DB {
	dsn := GetDatabaseDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Errorf("Error when connecting to database: %s", err))
	}

	return db
}

var DB *gorm.DB = GetConnection()

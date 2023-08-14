package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabaseDsn() string {
	host := os.Getenv("APP_DATABASE_HOST")
	port := os.Getenv("APP_DATABASE_PORT")
	user := os.Getenv("APP_DATABASE_USER")
	password := os.Getenv("APP_DATABASE_PASSWORD")
	dbname := os.Getenv("APP_DATABASE_NAME")

	if len(port) == 0 {
		port = "5432"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
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

package settings

import (
	"os"
	"strconv"
)

type BaseSettings struct {
	PORT                  int
	USERS_GRPC_HOST       string
	USERS_GRPC_PORT       int
	APP_DATABASE_HOST     string
	APP_DATABASE_PORT     int
	APP_DATABASE_USER     string
	APP_DATABASE_PASSWORD string
	APP_DATABASE_NAME     string
	SECRET_KEY            string
}

var Settings *BaseSettings = NewSettings()

func NewSettings() *BaseSettings {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8000
	}

	users_grpc_host := os.Getenv("USERS_GRPC_HOST")
	if users_grpc_host == "" {
		panic("You need to specify USERS_GRPC_HOST environment variable")
	}

	users_grpc_port, err := strconv.Atoi(os.Getenv("USERS_GRPC_PORT"))
	if err != nil {
		panic("You need to specify numeric USERS_GRPC_PORT environment variable")
	}

	db_host := os.Getenv("APP_DATABASE_HOST")
	if db_host == "" {
		panic("You need to specify APP_DATABASE_HOST environment variable")
	}

	db_port, err := strconv.Atoi(os.Getenv("APP_DATABASE_PORT"))
	if err != nil {
		db_port = 5432
	}

	db_user := os.Getenv("APP_DATABASE_USER")
	if db_user == "" {
		panic("You need to specify APP_DATABASE_USER environment variable")
	}

	db_password := os.Getenv("APP_DATABASE_PASSWORD")
	if db_password == "" {
		panic("You need to specify APP_DATABASE_PASSWORD environment variable")
	}

	db_name := os.Getenv("APP_DATABASE_NAME")
	if db_name == "" {
		panic("You need to specify APP_DATABASE_NAME environment variable")
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		panic("You need to specify SECRET_KEY environment variable")
	}

	return &BaseSettings{
		PORT:                  port,
		USERS_GRPC_HOST:       users_grpc_host,
		USERS_GRPC_PORT:       users_grpc_port,
		APP_DATABASE_HOST:     db_host,
		APP_DATABASE_PORT:     db_port,
		APP_DATABASE_USER:     db_user,
		APP_DATABASE_PASSWORD: db_password,
		APP_DATABASE_NAME:     db_name,
		SECRET_KEY:            secretKey,
	}
}

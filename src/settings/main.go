package settings

import (
	"os"
	"strconv"
)

type BaseSettings struct {
	PORT                         int
	USERS_GRPC_HOST              string
	USERS_GRPC_PORT              int
    APP_RABBIT_USER              string
    APP_RABBIT_PASSWORD          string
    APP_RABBIT_HOST              string
    APP_RABBIT_PORT              int
    APP_RABBIT_EVENTS_QUEUE_NAME string
	APP_DATABASE_HOST            string
	APP_DATABASE_PORT            int
	APP_DATABASE_USER            string
	APP_DATABASE_PASSWORD        string
	APP_DATABASE_NAME            string
    APP_ALLOW_ORIGINS            string
	SECRET_KEY                   string
}

var Settings *BaseSettings = NewSettings()

func NewSettings() *BaseSettings {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8000
	}

    rabbit_user := os.Getenv("APP_RABBIT_USER")
    if rabbit_user == "" {
        panic("You need to specify rabbitmq user")
    }

    rabbit_password := os.Getenv("APP_RABBIT_PASSWORD")
    if rabbit_password == "" {
        panic("You need to specify rabbitmq password")
    }

    rabbit_host := os.Getenv("APP_RABBIT_HOST")
    if rabbit_host == "" {
        panic("You need to specify rabbitmq host")
    }

	rabbit_port, err := strconv.Atoi(os.Getenv("APP_RABBIT_PORT"))
    if err != nil {
        panic("You need to specify rabbit port")
    }

    rabbit_queue_name := os.Getenv("APP_RABBIT_EVENTS_QUEUE_NAME")
    if rabbit_queue_name == "" {
        panic("You need to specify rabbit events queue  name")
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

    allowOrigins := os.Getenv("APP_ALLOW_ORIGINS")
    if allowOrigins == "" {
        allowOrigins = "*"
    }

	return &BaseSettings{
		PORT:                         port,
		USERS_GRPC_HOST:              users_grpc_host,
		USERS_GRPC_PORT:              users_grpc_port,
        APP_RABBIT_HOST:              rabbit_host,
        APP_RABBIT_PORT:              rabbit_port,
        APP_RABBIT_USER:              rabbit_user,
        APP_RABBIT_PASSWORD:          rabbit_password,
        APP_RABBIT_EVENTS_QUEUE_NAME: rabbit_queue_name,
		APP_DATABASE_HOST:            db_host,
		APP_DATABASE_PORT:            db_port,
		APP_DATABASE_USER:            db_user,
		APP_DATABASE_PASSWORD:        db_password,
		APP_DATABASE_NAME:            db_name,
        APP_ALLOW_ORIGINS:            allowOrigins,
		SECRET_KEY:                   secretKey,
	}
}

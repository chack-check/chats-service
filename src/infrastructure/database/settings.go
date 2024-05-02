package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type SettingsSchema struct {
	APP_DATABASE_HOST     string
	APP_DATABASE_NAME     string
	APP_DATABASE_PASSWORD string
	APP_DATABASE_PORT     int
	APP_DATABASE_USER     string
}

func InitSettings() SettingsSchema {
	host := os.Getenv("APP_DATABASE_HOST")
	if host == "" {
		panic(fmt.Errorf("you need to specify `APP_DATABASE_HOST` environment variable"))
	}

	name := os.Getenv("APP_DATABASE_NAME")
	if name == "" {
		panic(fmt.Errorf("you need to specify `APP_DATABASE_NAME` environment variable"))
	}

	password := os.Getenv("APP_DATABASE_PASSWORD")
	if password == "" {
		panic(fmt.Errorf("you need to specify `APP_DATABASE_PASSWORD` environment variable"))
	}

	port := os.Getenv("APP_DATABASE_PORT")
	if port == "" {
		panic(fmt.Errorf("you need to specify `APP_DATABASE_PORT` environment variable"))
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(errors.Join(fmt.Errorf("specify the correct number for database port"), err))
	}

	user := os.Getenv("APP_DATABASE_USER")
	if user == "" {
		panic(fmt.Errorf("you need to specify `APP_DATABASE_USER` environment variable"))
	}

	return SettingsSchema{
		APP_DATABASE_HOST:     host,
		APP_DATABASE_NAME:     name,
		APP_DATABASE_PASSWORD: password,
		APP_DATABASE_PORT:     portInt,
		APP_DATABASE_USER:     user,
	}
}

var Settings SettingsSchema = InitSettings()

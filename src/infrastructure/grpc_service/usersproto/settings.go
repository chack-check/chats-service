package usersproto

import (
	"fmt"
	"os"
	"strconv"
)

type SettingsSchema struct {
	APP_USERS_GRPC_HOST string
	APP_USERS_GRPC_PORT int
}

func InitSettings() SettingsSchema {
	host := os.Getenv("APP_USERS_GRPC_HOST")
	if host == "" {
		panic(fmt.Errorf("you need to specify `APP_USERS_GRPC_HOST` environment variable"))
	}

	port := os.Getenv("APP_USERS_GRPC_PORT")
	if port == "" {
		panic(fmt.Errorf("you need to specify `APP_USERS_GRPC_PORT` environment variable"))
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	return SettingsSchema{
		APP_USERS_GRPC_HOST: host,
		APP_USERS_GRPC_PORT: portInt,
	}
}

var Settings SettingsSchema = InitSettings()

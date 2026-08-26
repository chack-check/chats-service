package redisdb

import (
	"fmt"
	"os"
)

type SettingsSchema struct {
	APP_REDIS_URL string
}

func InitSettings() SettingsSchema {
	url := os.Getenv("APP_REDIS_URL")
	if url == "" {
		panic(fmt.Errorf("you need to specify `APP_REDIS_URL` environment variable"))
	}

	return SettingsSchema{
		APP_REDIS_URL: url,
	}
}

var Settings SettingsSchema = InitSettings()

package settings

import (
	"fmt"
	"os"
	"strconv"
)

type SettingsSchema struct {
	APP_PORT          int
	APP_SECRET_KEY    string
	APP_ALLOW_ORIGINS string
}

func InitSettings() SettingsSchema {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(fmt.Errorf("error parsing port. Please specify the correct number"))
	}

	secretKey := os.Getenv("APP_SECRET_KEY")
	if secretKey == "" {
		panic(fmt.Errorf("you need to specify `APP_SECRET_KEY` environment variable"))
	}

	allowOrigins := os.Getenv("APP_ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "*"
	}

	return SettingsSchema{
		APP_PORT:          portInt,
		APP_SECRET_KEY:    secretKey,
		APP_ALLOW_ORIGINS: allowOrigins,
	}
}

var Settings SettingsSchema = InitSettings()

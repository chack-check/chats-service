package chatsproto

import (
	"fmt"
	"os"
)

type SettingsSchema struct {
	APP_SECRET_KEY string
}

func InitSettings() SettingsSchema {
	secretKey := os.Getenv("APP_SECRET_KEY")
	if secretKey == "" {
		panic(fmt.Errorf("you need to specify `APP_SECRET_KEY` environment variable"))
	}

	return SettingsSchema{
		APP_SECRET_KEY: secretKey,
	}
}

var Settings SettingsSchema = InitSettings()

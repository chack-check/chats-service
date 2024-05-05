package database

import (
	"fmt"
	"os"
)

type SettingsSchema struct {
	APP_DATABASE_DSN string
}

func InitSettings() SettingsSchema {
	databaseDsn := os.Getenv("APP_DATABASE_DSN")
	if databaseDsn == "" {
		panic(fmt.Errorf("you need to specify `APP_DATABASE_DSN` environment variable"))
	}

	return SettingsSchema{
		APP_DATABASE_DSN: databaseDsn,
	}
}

var Settings SettingsSchema = InitSettings()

package filesservice

import (
	"fmt"
	"os"
)

type SettingsSchema struct {
	FILES_SIGNATURE_KEY string
}

func InitSettings() SettingsSchema {
	key := os.Getenv("FILES_SIGNATURE_KEY")
	if key == "" {
		panic(fmt.Errorf("you need to specify `FILES_SIGNATURE_KEY` environment variable"))
	}

	return SettingsSchema{
		FILES_SIGNATURE_KEY: key,
	}
}

var Settings SettingsSchema = InitSettings()

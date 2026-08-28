package configs

import (
	"net/url"
	"sync"

	"github.com/caarlos0/env/v11"
)

type APIConfiguration struct {
	Port                       int     `env:"APP_PORT,required,notEmpty"`
	SecretKey                  string  `env:"APP_SECRET_KEY,required,notEmpty"`
	AllowOrigins               string  `env:"APP_ALLOW_ORIGINS"`
	SavedMessagesChatAvatarURL url.URL `env:"APP_SAVED_MESSAGES_CHAT_AVATAR_URL"`
}

var (
	apiConfig APIConfiguration
	apiOnce   sync.Once
)

func GetAPIConfiguration() *APIConfiguration {
	apiOnce.Do(func() {
		if err := env.Parse(&apiConfig); err != nil {
			panic(err)
		}
	})

	return &apiConfig
}

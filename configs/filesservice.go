package configs

import (
	"sync"

	"github.com/caarlos0/env/v11"
)

type FilesServiceConfiguration struct {
	SignatureKey string `env:"FILES_SIGNATURE_KEY,required,notEmpty"`
}

var (
	filesServiceConfig FilesServiceConfiguration
	filesServiceOnce   sync.Once
)

func GetFilesServiceConfiguration() *FilesServiceConfiguration {
	filesServiceOnce.Do(func() {
		if err := env.Parse(&filesServiceConfig); err != nil {
			panic(err)
		}
	})

	return &filesServiceConfig
}

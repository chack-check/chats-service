package configs

import (
	"sync"

	"github.com/caarlos0/env/v11"
)

type DatabaseConfiguration struct {
	DSN string `env:"APP_DATABASE_DSN,required,notEmpty"`
}

var (
	databaseConfig DatabaseConfiguration
	databaseOnce   sync.Once
)

func GetDatabaseConfiguration() *DatabaseConfiguration {
	databaseOnce.Do(func() {
		if err := env.Parse(&databaseConfig); err != nil {
			panic(err)
		}
	})

	return &databaseConfig
}

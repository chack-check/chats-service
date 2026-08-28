package configs

import (
	"sync"

	"github.com/caarlos0/env/v11"
)

type GRPCUsersConfiguration struct {
	Host string `env:"APP_USERS_GRPC_HOST,required,notEmpty"`
	Port int    `env:"APP_USERS_GRPC_PORT,required,notEmpty"`
}

var (
	grpcUsersConfig GRPCUsersConfiguration
	grpcUsersOnce   sync.Once
)

func GetGRPCUsersConfiguration() *GRPCUsersConfiguration {
	grpcUsersOnce.Do(func() {
		if err := env.Parse(&grpcUsersConfig); err != nil {
			panic(err)
		}
	})

	return &grpcUsersConfig
}

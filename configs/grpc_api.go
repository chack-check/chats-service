package configs

import (
	"sync"

	"github.com/caarlos0/env/v11"
)

type GRPCAPIConfiguration struct {
	Host string `env:"APP_GRPC_HOST,required,notEmpty"`
	Port int    `env:"APP_GRPC_PORT,required,notEmpty"`
}

var (
	grpcAPIConfig GRPCAPIConfiguration
	grpcAPIOnce   sync.Once
)

func GetGRPCAPIConfiguration() *GRPCAPIConfiguration {
	grpcAPIOnce.Do(func() {
		if err := env.Parse(&grpcAPIConfig); err != nil {
			panic(err)
		}
	})

	return &grpcAPIConfig
}

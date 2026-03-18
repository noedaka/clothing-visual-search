package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerPort     string `env:"GO_PORT"`
	MLServiceAddr  string `env:"GRPC_ML_SERVICE_ADDR"`
	DataBaseURL    string `env:"DATABASE_URL"`
	MinIOAddr      string `env:"MINIO_ENDPOINT"`
	MinIOAccessKey string `env:"MINIO_ACCESS_KEY"`
	MinIOSecretKey string `env:"MINIO_SECRET_KEY"`
}

func Init() (*Config, error) {
	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	flag.Parse()

	return cfg, nil
}

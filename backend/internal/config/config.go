package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerPort    string `env:"GO_PORT"`
	MLServiceAddr string `env:"GRPC_ML_SERVICE_ADDR"`

	PostgresURL string `env:"DATABASE_URL"`

	MilvusAddr string `env:"MILVUS_ENDPOINT"`

	TopK      int     `env:"TOP_K"`
	Threshold float64 `env:"THRESHOLD"`

	MinIOAddr            string `env:"MINIO_ENDPOINT"`
	MinIOAccessKey       string `env:"MINIO_ACCESS_KEY"`
	MinIOSecretKey       string `env:"MINIO_SECRET_KEY"`
	MinIOExternalBaseURL string `env:"MINIO_EXTERNAL_URL"`
	MinIOBucket          string `env:"MINIO_BUCKET"`
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

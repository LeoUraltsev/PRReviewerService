package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLever          slog.Level    `env:"LOG_LEVEL" env-default:"INFO"`
	Host              string        `env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port              string        `env:"SERVER_PORT" env-default:"8080"`
	ReadTimeout       time.Duration `env:"SERVER_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout      time.Duration `env:"SERVER_WRITE_TIMEOUT" env-default:"5s"`
	IdleTimeout       time.Duration `env:"SERVER_IDLE_TIMEOUT" env-default:"5s"`
	ReadHeaderTimeout time.Duration `env:"SERVER_READ_HEADER_TIMEOUT" env-default:"5s"`
	ConnectionString  string        `env:"POSTGRES_PR_CONNECTION_STRING" env-default:"postgres://postgres:password@localhost:5432/database"`
	MaxRetries        int           `env:"DB_MAX_RETRIES" env-default:"3"`
	RetryInterval     time.Duration `env:"DB_RETRY_INTERVAL" env-default:"5s"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}

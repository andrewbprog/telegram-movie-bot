package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"time"
	"tlgbs/internal/postgres"
)

type Config struct {
	Database postgres.DBConfig

	AppEnv string `env:"APP_ENV" envDefault:"development"`

	TelegramToken string `env:"TELEGRAM_BOT_TOKEN" envDefault:""`

	MigrationTimeout time.Duration `env:"MIGRATION_TIMEOUT" envDefault:"30s"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	return cfg, nil
}

package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"time"
	"tlgbs/internal/postgres"
)

type KafkaConfig struct {
	Brokers []string `env:"KAFKA_BROKERS" envDefault:"kafka:9092"`
	Topic   string   `env:"KAFKA_TOPIC" envDefault:"recommendations_notifications"`
	GroupID string   `env:"KAFKA_GROUP_ID" envDefault:"=telegram_bot_consumer"`
}

type Config struct {
	Database postgres.DBConfig

	AppEnv string `env:"APP_ENV" envDefault:"development"`

	TelegramToken string `env:"TELEGRAM_BOT_TOKEN" envDefault:""`

	MigrationTimeout time.Duration `env:"MIGRATION_TIMEOUT" envDefault:"30s"`

	GatewayURL     string        `env:"GATEWAY_URL,required"`
	GatewayTimeout time.Duration `env:"GATEWAY_TIMEOUT" envDefault:"5s"`

	Kafka KafkaConfig
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	return cfg, nil
}

package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Addr string `validate:"required" env:"ADDR_TELEGRAM_FORWARDER"`
}

type Kafka struct {
	BootstrapServersValue string `validate:"required" env:"BOOTSTRAP_SERVERS_VALUE"`
	TgMessTopicName       string `validate:"required" env:"TOPIC_NAME_TELEGRAM_UPDATES"`
	SaslUsername          string `validate:"required" env:"SASL_USERNAME"`
	SaslPassword          string `validate:"required" env:"SASL_PASSWORD"`
}

type Telegram struct {
	Token string `validate:"required" env:"TOKEN"`
}

type Config struct {
	Server   Server   `envPrefix:"SERVER_"`
	Kafka    Kafka    `envPrefix:"KAFKA_"`
	Telegram Telegram `envPrefix:"TELEGRAM_"`
}

func Load() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("env parse: %w", err)
	}

	v := validator.New()
	if err := v.Struct(c); err != nil {
		return nil, fmt.Errorf("config validate: %w", err)
	}

	return &c, nil
}

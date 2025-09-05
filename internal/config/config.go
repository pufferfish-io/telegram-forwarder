package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Addr string `validate:"required" env:"ADDR" envDefault:":8080"`
}

type Kafka struct {
	BootstrapServersKey   string `validate:"required" env:"BOOTSTRAP_SERVERS_KEY" envDefault:"bootstrap.servers"`
	BootstrapServersValue string `validate:"required" env:"BOOTSTRAP_SERVERS_VALUE"`
}

type Telegram struct {
	Token string `validate:"required" env:"TOKEN"`
}

type Config struct {
	Server   Server   `envPrefix:"TG_FORWARDER_TELEGRAM_SERVER_"`
	Kafka    Kafka    `envPrefix:"TG_FORWARDER_TELEGRAM_KAFKA_"`
	Telegram Telegram `envPrefix:"TG_FORWARDER_TELEGRAM_"`
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

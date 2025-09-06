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
	BootstrapServersKey   string `validate:"required" env:"BOOTSTRAP_SERVERS_KEY"   envDefault:"bootstrap.servers"`
	BootstrapServersValue string `validate:"required" env:"BOOTSTRAP_SERVERS_VALUE"`
	TgMessTopicName       string `validate:"required" env:"TG_MESS_TOPIC_NAME"`

	// Новые поля для SASL/SCRAM
	SecurityProtocol string `validate:"required,oneof=PLAINTEXT SSL SASL_PLAINTEXT SASL_SSL" env:"SECURITY_PROTOCOL" envDefault:"SASL_PLAINTEXT"`
	SaslMechanism    string `validate:"required,oneof=PLAIN SCRAM-SHA-256 SCRAM-SHA-512 OAUTHBEARER" env:"SASL_MECHANISM" envDefault:"SCRAM-SHA-512"`
	SaslUsername     string `validate:"required" env:"SASL_USERNAME"`
	SaslPassword     string `validate:"required" env:"SASL_PASSWORD"`
}

type Telegram struct {
	Token string `validate:"required" env:"TOKEN"`
}

type Api struct {
	TgWebHookPath   string `validate:"required" env:"TG_WEB_HOOK_PATH"`
	HealthCheckPath string `validate:"required" env:"HEALTH_CHECK_PATH"`
}
type Config struct {
	Server   Server   `envPrefix:"TG_FORWARDER_SERVER_"`
	Kafka    Kafka    `envPrefix:"TG_FORWARDER_KAFKA_"`
	Telegram Telegram `envPrefix:"TG_FORWARDER_TELEGRAM_"`
	Api      Api      `envPrefix:"TG_FORWARDER_API_"`
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

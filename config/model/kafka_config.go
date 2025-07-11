package model

type KafkaConfig struct {
	BootstrapServersValue string `yaml:"bootstrap_servers_value"`
}

func (KafkaConfig) SectionName() string {
	return "kafka"
}

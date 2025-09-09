package messaging

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/xdg-go/scram"
)

func EnsureTopic(opt Option, topic string, partitions int32, replicationFactor int16) error {
	if topic == "" {
		return fmt.Errorf("topic is empty")
	}

	cfg := sarama.NewConfig()

	cfg.Net.SASL.Enable = true
	cfg.Net.SASL.User = opt.SaslUsername
	cfg.Net.SASL.Password = opt.SaslPassword
	cfg.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
		return &xdgSCRAMClient{hash: scram.SHA512}
	}

	admin, err := sarama.NewClusterAdmin([]string{opt.Broker}, cfg)
	if err != nil {
		return fmt.Errorf("kafka admin init: %w")
	}
	defer func() { _ = admin.Close() }()

	topics, err := admin.ListTopics()
	if err != nil {
		return fmt.Errorf("list topics: %w", err)
	}
	if _, exists := topics[topic]; exists {
		if opt.Logger != nil {
			opt.Logger.Info("Kafka topic already exists: %s", topic)
		}
		return nil
	}

	detail := &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
		ConfigEntries:     map[string]*string{},
	}
	if err := admin.CreateTopic(topic, detail, false); err != nil {
		if opt.Logger != nil {
			opt.Logger.Error("Kafka create topic failed: %v", err)
		}
		return fmt.Errorf("create topic: %w", err)
	}
	if opt.Logger != nil {
		opt.Logger.Info("Kafka topic created: %s partitions=%d rf=%d", topic, partitions, replicationFactor)
	}
	return nil
}

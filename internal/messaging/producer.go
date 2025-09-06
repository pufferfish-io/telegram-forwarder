package messaging

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tgforwarder/internal/logger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Option struct {
	Logger              logger.Logger
	BootstrapServersKey string // обычно "bootstrap.servers"
	Brokers             string // "176.98.176.199:9092"

	// SASL/SCRAM (или другой SASL) — опционально, если пусто, не настраиваем
	SecurityProtocol string // "SASL_PLAINTEXT" | "SASL_SSL" | "PLAINTEXT" | "SSL"
	SaslMechanism    string // "SCRAM-SHA-512" | "SCRAM-SHA-256" | "PLAIN" | "OAUTHBEARER"
	SaslUsername     string
	SaslPassword     string

	// (опц.) таймаут ожидания доставки в Send, если контекст без дедлайна
	DeliveryTimeout time.Duration
}

type KafkaProducer struct {
	log      logger.Logger
	producer *kafka.Producer
}

func NewKafkaProducer(opt Option) (*KafkaProducer, error) {
	if opt.Brokers == "" {
		return nil, errors.New("brokers is empty")
	}
	if opt.BootstrapServersKey == "" {
		opt.BootstrapServersKey = "bootstrap.servers"
	}
	if opt.DeliveryTimeout <= 0 {
		opt.DeliveryTimeout = 5 * time.Second
	}

	cfg := &kafka.ConfigMap{
		opt.BootstrapServersKey: opt.Brokers,
		// полезные дефолты на клиенте:
		// "message.timeout.ms": 30000, // сколько ждать ack от брокера
		// "linger.ms":          5,     // микро-буферизация на продюсере
		// "acks":               "all", // надежнее (можно "1" для быстрее)
	}

	// Добавим SASL, если задано
	if opt.SecurityProtocol != "" {
		_ = cfg.SetKey("security.protocol", opt.SecurityProtocol)
	}
	if opt.SaslMechanism != "" {
		_ = cfg.SetKey("sasl.mechanism", opt.SaslMechanism)
	}
	if opt.SaslUsername != "" {
		_ = cfg.SetKey("sasl.username", opt.SaslUsername)
	}
	if opt.SaslPassword != "" {
		_ = cfg.SetKey("sasl.password", opt.SaslPassword)
	}

	p, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, fmt.Errorf("kafka producer init: %w", err)
	}
	if opt.Logger != nil {
		opt.Logger.Info("Kafka Producer init: %s=%s", opt.BootstrapServersKey, opt.Brokers)
	}
	return &KafkaProducer{log: opt.Logger, producer: p}, nil
}

func (kp *KafkaProducer) Send(ctx context.Context, topic string, data []byte) error {
	if kp.producer == nil {
		return errors.New("kafka producer is nil")
	}
	delivery := make(chan kafka.Event, 1)
	defer close(delivery)

	if err := kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic, // отдельная переменная — безопасно
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}, delivery); err != nil {
		return fmt.Errorf("produce: %w", err)
	}

	// если у ctx нет дедлайна — подстрахуем своим
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	select {
	case e := <-delivery:
		m, ok := e.(*kafka.Message)
		if !ok {
			return errors.New("unexpected delivery event type")
		}
		if m.TopicPartition.Error != nil {
			return m.TopicPartition.Error
		}
		if kp.log != nil {
			kp.log.Info("Kafka delivered topic=%s partition=%d offset=%v bytes=%d",
				topic, m.TopicPartition.Partition, m.TopicPartition.Offset, len(data))
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("delivery timeout: %w", ctx.Err())
	}
}

func (kp *KafkaProducer) Close() {
	if kp.producer == nil {
		return
	}
	kp.producer.Flush(5000) // мс
	kp.producer.Close()
	if kp.log != nil {
		kp.log.Info("Kafka Producer closed")
	}
}

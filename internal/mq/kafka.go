package mq

import (
	"context"
	"my-chat/internal/config"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
}

func NewKafkaClient(cfg *config.KafkaConfig) *KafkaClient {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Addr),
		Topic:        cfg.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.Addr},
		Topic:    cfg.Topic,
		GroupID:  cfg.Group,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &KafkaClient{
		Writer: writer,
		Reader: reader,
	}
}
func (k *KafkaClient) Publish(ctx context.Context, key, value []byte) error {
	return k.Writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}
func (k *KafkaClient) Close() {
	if k.Writer != nil {
		_ = k.Writer.Close()
	}
	if k.Reader != nil {
		_ = k.Reader.Close()
	}
}

package mq

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"my-chat/internal/config"
	"time"
)

type KafkaClient struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
}

var GlobalKafka *KafkaClient

func InitKafka() {
	cfg := config.GlobalConfig.Kafka
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
	GlobalKafka = &KafkaClient{
		Writer: writer,
		Reader: reader,
	}
	log.Println("Kafka初始化成功")
}
func (k *KafkaClient) Publish(ctx context.Context, key, value []byte) error {
	return k.Writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}

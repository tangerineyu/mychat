package bootstrap

import (
	"context"
	"fmt"
	"my-chat/internal/config"
	"my-chat/internal/dao"
	"my-chat/internal/mq"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Deps struct {
	DB    *gorm.DB
	Redis *redis.Client
	Kafka *mq.KafkaClient
}

func InitDeps(cfg *config.Config) (*Deps, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil config")
	}

	db, err := dao.NewMySQL(&cfg.MySQL)
	if err != nil {
		return nil, fmt.Errorf("init mysql: %w", err)
	}

	rdb, err := dao.NewRedis(&cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("init redis: %w", err)
	}

	kafkaClient := mq.NewKafkaClient(&cfg.Kafka)

	return &Deps{DB: db, Redis: rdb, Kafka: kafkaClient}, nil
}

func CloseDeps(ctx context.Context, d *Deps, log *zap.Logger) {
	_ = ctx
	if d == nil {
		return
	}
	if d.Kafka != nil {
		d.Kafka.Close()
	}
	if d.Redis != nil {
		_ = d.Redis.Close()
	}
	// gorm DB close: gorm v2 needs sqlDB
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		} else if log != nil {
			log.Warn("get sql db failed", zap.Error(err))
		}
	}
}

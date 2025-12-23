package dao

import (
	"context"
	"log"
	"my-chat/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	c := config.GlobalConfig.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		PoolSize:     c.PoolSize,
		MinIdleConns: 10,
		DialTimeout:  5 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis 连接失败: %v", err)
	}
	log.Println("redis 连接成功")
}

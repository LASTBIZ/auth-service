package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type redisConfig struct {
	Host     string
	Port     string
	Password string
}

func NewRedisConfig(host, port, password string) *redisConfig {
	return &redisConfig{
		Host:     host,
		Port:     port,
		Password: password,
	}
}

func NewClient(ctx context.Context, cfg redisConfig) (client *redis.Client, err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       0,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

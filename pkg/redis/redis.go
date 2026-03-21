package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rzkhosroshahi/velox/config"
	"github.com/rzkhosroshahi/velox/pkg/logger"
)

func New(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	logger.Log.Info("connected to the redis!")

	return client, nil
}

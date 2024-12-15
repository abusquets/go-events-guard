package redis

import (
	"context"
	"eventsguard/internal/infrastructure/config"
	"eventsguard/internal/infrastructure/mylog"

	"github.com/go-redis/redis/v8"
)

func ConnectToRedis(config *config.AppConfig) (*redis.Client, error) {
	logger := mylog.GetLogger()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// Perform a 'ping' to ensure the connection is working
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.ErrorWithErr("Error connecting to Redis", err)
		return nil, err
	}

	return rdb, nil
}

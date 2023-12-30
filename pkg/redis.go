package pkg

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func NewDiscordBotUserSessionRedisClient() *redis.Client {
	host := os.Getenv("REDIS_DISCORD_HOST")
	port := os.Getenv("REDIS_DISCORD_PORT")
	url := fmt.Sprintf("%s:%s", host, port)

	return redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})
}

func NewRedisClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_POST")
	url := fmt.Sprintf("%s:%s", host, port)

	db := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})

	_, err := db.Ping(context.Background()).Result()
	if err != nil {
		logrus.Error("failed to ping redis server")
		return nil
	}

	return db
}

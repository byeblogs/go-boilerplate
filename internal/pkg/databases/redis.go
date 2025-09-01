package databases

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func NewRedisClient() (*redis.Client, error) {
	var (
		rdb *redis.Client
		err error
	)

	err = godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDb,
	})

	_, err = rdb.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

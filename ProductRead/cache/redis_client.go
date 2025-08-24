package cache

import (
	"context"
	"os"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClientFromEnv() (*RedisClient, error) {
	masterName := os.Getenv("REDIS_MASTER_NAME")
	sentinels := os.Getenv("REDIS_SENTINELS")
	if masterName == "" || sentinels == "" {
		return nil, nil
	}

	addrs := strings.Split(sentinels, ",")

	opt := &redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: addrs,
		DialTimeout:   5 * time.Second,
		ReadTimeout:   3 * time.Second,
		WriteTimeout:  3 * time.Second,
		PoolSize:      10,
		MinIdleConns:  2,
	}

	cli := redis.NewFailoverClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{client: cli}, nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

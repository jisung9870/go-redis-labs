package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     int
	DB       int
	Username string
	Password string
}

type RedisClient struct {
	client *redis.Client
}

func RedisConnect(cfg RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return &RedisClient{}, fmt.Errorf("failed to connect to redis server: %v", err)
	}

	return &RedisClient{client: rdb}, nil
}

func (r *RedisClient) Close() {
	r.client.Close()
}

func (r *RedisClient) LPush(key string, data string) error {
	ctx := context.Background()
	err := r.client.LPush(ctx, key, data).Err()
	if err != nil {
		return fmt.Errorf("failed to LPush value to Redis: %v", err)
	}
	return nil
}

func (r RedisClient) RPop(key string) (string, error) {
	ctx := context.Background()
	data, err := r.client.RPop(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist: %v", err)
	} else if err != nil {
		return "", fmt.Errorf("failed to RPop value from redis: %v", err)
	}
	return data, nil
}

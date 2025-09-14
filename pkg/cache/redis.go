package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"ai-assistant/internal/app/config"
	"ai-assistant/pkg/logger"
)

type RedisService struct {
	client *redis.Client
	ctx    context.Context
	logger *logger.Logger
}

func NewRedisService(cfg *config.Config) *RedisService {
	opts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		logger.New().Error("Failed to parse Redis URL:", err)
		return nil
	}

	client := redis.NewClient(opts)
	
	return &RedisService{
		client: client,
		ctx:    context.Background(),
		logger: logger.New(),
	}
}

func (r *RedisService) Connect() error {
	_, err := r.client.Ping(r.ctx).Result()
	if err != nil {
		r.logger.Error("Failed to connect to Redis:", err)
		return err
	}
	r.logger.Info("Successfully connected to Redis")
	return nil
}

func (r *RedisService) Disconnect() error {
	return r.client.Close()
}

func (r *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisService) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisService) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisService) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	return count > 0, err
}

func (r *RedisService) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

func (r *RedisService) HSet(key string, field string, value interface{}) error {
	return r.client.HSet(r.ctx, key, field, value).Err()
}

func (r *RedisService) HGet(key string, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

func (r *RedisService) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

func (r *RedisService) HDel(key string, fields ...string) error {
	return r.client.HDel(r.ctx, key, fields...).Err()
}
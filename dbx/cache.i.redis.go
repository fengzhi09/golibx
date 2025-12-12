package dbx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fengzhi09/golibx/jsonx"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存实例
func NewRedisCache(ctx context.Context, conf CacheConf) (*RedisCache, error) {
	// 解析配置
	addr, ok := conf["addr"].(string)
	if !ok {
		return nil, fmt.Errorf("redis addr not found in config")
	}

	db := 0
	if dbVal, exists := conf["db"]; exists {
		db = int(dbVal.(float64))
	}

	password := ""
	if pwVal, exists := conf["password"]; exists {
		password = pwVal.(string)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return &RedisCache{
		client: client,
	}, nil
}

// Has 检查键是否存在
func (r *RedisCache) Has(ctx context.Context, key string) bool {
	exists, err := r.client.Exists(ctx, key).Result()
	return err == nil && exists > 0
}

// Get 获取值
func (r *RedisCache) Get(ctx context.Context, key string) (jsonx.JValue, error) {
	val := r.client.Get(ctx, key).Val()
	return jsonx.GoV2JV(val), nil
}

// Set 设置值
func (r *RedisCache) Set(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, 0).Err()
}

// SetEx 设置带过期时间的值
func (r *RedisCache) SetEx(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

// Close 关闭缓存连接
func (r *RedisCache) Close(ctx context.Context) error {
	return r.client.Close()
}

package dbx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fengzhi09/golibx/jsonx"
	"sync"
	"time"

	"github.com/fengzhi09/golibx/gox"
)

// MemCache 内存缓存实现
type MemCache struct {
	lru       *gox.LRUMap[cacheItem]
	cleanSecs int
	sync.RWMutex
}

// cacheItem 缓存项
type cacheItem struct {
	value      []byte
	expiration int64
}

// NewMemCache 创建内存缓存实例
func NewMemCache(ctx context.Context, conf CacheConf) *MemCache {
	size := conf.GetOr("max", 10000).ToInt()
	cleanSec := conf.GetOr("clean_sec", 300).ToInt()
	cache := &MemCache{lru: gox.NewLRUMap[cacheItem](size), cleanSecs: cleanSec}
	// 启动清理过期项的协程
	go cache.cleanup(ctx)
	return cache
}

// Has 检查键是否存在
func (c *MemCache) Has(ctx context.Context, key string) bool {
	return c.lru.Data().Has(key)
}

// Get 获取值
func (c *MemCache) Get(ctx context.Context, key string) (jsonx.JValue, error) {
	item, exists := c.lru.Get(key)

	if !exists {
		return jsonx.JNull{}, fmt.Errorf("key %s not found", key)
	}

	// 检查是否过期
	if item.expiration > 0 && item.expiration < time.Now().UnixNano() {
		// 删除过期项
		c.lru.Del(key)
		return jsonx.JNull{}, fmt.Errorf("key %s expired", key)
	}

	return jsonx.GoV2JV(item.value), nil
}

// Set 设置值
func (c *MemCache) Set(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.lru.Set(key, cacheItem{
		value:      data,
		expiration: 0, // 不过期, 但会被自动清除策略清除
	})

	return nil
}

// SetEx 设置带过期时间的值
func (c *MemCache) SetEx(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.lru.Set(key, cacheItem{
		value:      data,
		expiration: time.Now().Add(expiration).UnixNano(),
	})

	return nil
}

// Close 关闭缓存连接
func (c *MemCache) Close(ctx context.Context) error {
	c.lru.Clear()
	return nil
}

// cleanup 清理过期项
func (c *MemCache) cleanup(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(c.cleanSecs) * time.Second)
	defer ticker.Stop()

	// 启动协程进行定时清理；优雅重启：当上下文取消时，程序终止时，停止清理协程
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.lru.DropByUpdateAt(time.Now().UnixNano() - int64(c.cleanSecs)*1000000000)
		}
	}

}

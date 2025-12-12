package dbx

// 缓存模块
import (
	"context"
	"sync"
	"time"

	"github.com/fengzhi09/golibx/jsonx"

	"github.com/fengzhi09/golibx/logx"

	"github.com/pkg/errors"
)

// ICache 缓存接口
type ICache interface {
	// Has 检查键是否存在
	Has(ctx context.Context, key string) bool
	// Get 获取值
	Get(ctx context.Context, key string) (jsonx.JValue, error)
	// Set 设置值
	Set(ctx context.Context, key string, value any) error
	// SetEx 设置带过期时间的值
	SetEx(ctx context.Context, key string, value any, expiration time.Duration) error
	// Close 关闭缓存连接
	Close(ctx context.Context) error
}

// CacheMgr 缓存管理器
type CacheMgr struct {
	cacheMap map[string]ICache
	confMap  map[string]CacheConf
	mutex    sync.RWMutex
}

var (
	cacheMgrInstance *CacheMgr
	cacheMgrOnce     sync.Once
)

// CacheX 获取缓存管理器实例
func CacheX() *CacheMgr {
	cacheMgrOnce.Do(func() {
		cacheMgrInstance = &CacheMgr{
			cacheMap: make(map[string]ICache),
			confMap:  make(map[string]CacheConf),
		}
	})
	return cacheMgrInstance
}

type CacheConf = jsonx.JObj

// Init 初始化缓存
func (cm *CacheMgr) Init(ctx context.Context, caches map[string]CacheConf) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	var err error

	for name, conf := range caches {
		cacheType := getOrDefault(&conf, "type", "mem")

		var cache ICache

		switch cacheType {
		case "mem":
			cache = NewMemCache(ctx, conf)
		case "redis":
			cache, err = NewRedisCache(ctx, conf)
		default:
			cache = NewMemCache(ctx, conf)
			err = errors.New("type not supported")
		}
		if cache == nil {
			conf.Put("mode", "mem")
			logx.WarnfM(ctx, "CacheMgr", "use mem instead; name:%v type:%s created failed: %v ", name, cacheType, err)
			cache = NewMemCache(ctx, conf)
		}

		cm.cacheMap[name] = cache
		cm.confMap[name] = conf
	}

	return nil
}

// Use 获取缓存实例
func (cm *CacheMgr) Use(name string) (ICache, error) {
	return cm.getOrCreate(name)
}

func (cm *CacheMgr) getOrCreate(name string) (ICache, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cache, exists := cm.cacheMap[name]
	if !exists {
		cm.cacheMap[name] = NewMemCache(context.Background(), jsonx.JObj{})
		cache = cm.cacheMap[name]
	}
	return cache, nil
}

// CloseAll 关闭所有缓存连接
func (cm *CacheMgr) CloseAll(ctx context.Context) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for name, cache := range cm.cacheMap {
		if err := cache.Close(ctx); err != nil {
			logx.ErrorfM(ctx, "CacheMgr", "关闭缓存 %s 失败: %v", name, err)
		}
	}

	cm.cacheMap = make(map[string]ICache)
}

// 导出的全局函数

// InitCache 初始化全局缓存
func InitCache(ctx context.Context, modules map[string]CacheConf) error {
	return CacheX().Init(ctx, modules)
}

// UseCache 获取全局缓存实例
func UseCache(name string) (ICache, error) {
	return CacheX().Use(name)
}

// CloseCaches 关闭所有全局缓存
func CloseCaches(ctx context.Context) {
	CacheX().CloseAll(ctx)
}

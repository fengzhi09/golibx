package dbx

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"
	"github.com/fengzhi09/golibx/logx"

	_ "github.com/go-sql-driver/mysql"
)

// ISQL 数据库接口
type ISQL interface {
	// Connect 建立数据库连接
	Connect(ctx context.Context) error
	// Query 执行查询
	Query(ctx context.Context, query string, args ...any) ([]*jsonx.JObj, error)
	// Close 关闭数据库连接
	Close(ctx context.Context) error
}

type DBConf = *jsonx.JObj

// getOrDefault 获取配置字符串
func getOrDefault(conf *jsonx.JObj, key, defaultValue string) string {
	return gox.IfElse(conf.Contains(key), conf.GetStr(key), defaultValue).(string)
}

func getInOrder(conf *jsonx.JObj, keys ...string) string {
	for _, key := range keys {
		if val := conf.GetStr(key); val != "" {
			return val
		}
	}
	return ""
}

// NewPSQL 创建PostgreSQL数据库实例
func NewPSQL(name string, conf DBConf) ISQL {
	return &PSql{
		name:   name,
		conf:   conf,
		dbType: "PostgreSQL",
	}
}

// NewMSQL 创建MySQL数据库实例
func NewMSQL(name string, conf DBConf) ISQL {
	return &MSql{
		name:   name,
		conf:   conf,
		dbType: "MySQL",
	}
}

// NewDorisDB 创建Doris数据库实例
func NewDoris(name string, conf DBConf) ISQL {
	return &Doris{
		name:   name,
		conf:   conf,
		dbType: "Doris",
	}
}

// DBMgr 数据库管理器
type DBMgr struct {
	dbMap map[string]ISQL
	mutex sync.RWMutex
}

var (
	dbMgrInstance *DBMgr
	dbMgrOnce     sync.Once
)

// GetDBMgrInstance 获取数据库管理器实例
func DB() *DBMgr {
	dbMgrOnce.Do(func() {
		dbMgrInstance = &DBMgr{
			dbMap: make(map[string]ISQL),
		}
	})
	return dbMgrInstance
}

// Init 初始化数据库连接
func (dm *DBMgr) Init(ctx context.Context, dbs map[string]DBConf) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	for name, conf := range dbs {
		// 从配置中获取URL或单独的连接参数
		dbType, err := dm.detectDBType(conf)
		if err != nil {
			return fmt.Errorf("failed to detect DB type for %s: %v", name, err)
		}

		var db ISQL

		switch dbType {
		case "postgresql", "postgres":
			db = NewPSQL(name, conf)
		case "mysql":
			db = NewMSQL(name, conf)
		case "doris":
			db = NewDoris(name, conf)
		default:
			return fmt.Errorf("unsupported database type: %s", dbType)
		}

		// 建立连接
		if err := db.Connect(ctx); err != nil {
			return fmt.Errorf("failed to connect to database %s: %v", name, err)
		}

		dm.dbMap[name] = db
	}

	return nil
}

// GetDB 获取数据库实例
func (dm *DBMgr) Use(name string) (ISQL, error) {
	dm.mutex.RLock()
	db, exists := dm.dbMap[name]
	dm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("database %s not found", name)
	}

	return db, nil
}

// CloseAll 关闭所有数据库连接
func (dm *DBMgr) CloseAll(ctx context.Context) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	for _, db := range dm.dbMap {
		if err := db.Close(ctx); err != nil {
			// 错误处理
			logx.Warnf(ctx, "关闭数据库 %s 连接失败: %v", db.(*PSql).name, err)
		}
	}

	dm.dbMap = make(map[string]ISQL)
}

// detectDBType 检测数据库类型
func (dm *DBMgr) detectDBType(conf DBConf) (string, error) {
	// 优先从type字段获取
	dbType := getInOrder(conf, "db_type", "type", "db_driver", "driver")
	if dbType != "" {
		return strings.ToLower(dbType), nil
	}
	url := getOrDefault(conf, "db_url", getOrDefault(conf, "pg_url", ""))
	// 从URL中检测数据库类型
	if url != "" {
		conf.Put("db_url", url)
		if strings.HasPrefix(strings.ToLower(url), "postgresql://") || strings.HasPrefix(strings.ToLower(url), "postgres://") {
			return "postgresql", nil
		} else if strings.HasPrefix(strings.ToLower(url), "mysql://") {
			return "mysql", nil
		} else if strings.HasPrefix(strings.ToLower(url), "doris://") {
			return "doris", nil
		}
	}

	return "", fmt.Errorf("无法检测数据库类型")
}

// 导出的全局函数

// InitDB 初始化全局数据库
func InitDB(ctx context.Context, dbs map[string]DBConf) error {
	return DB().Init(ctx, dbs)
}

// GetDB 获取全局数据库实例
func UseDB(name string) (ISQL, error) {
	return DB().Use(name)
}

// CloseAllDBs 关闭所有全局数据库连接
func CloseDBs(ctx context.Context) {
	DB().CloseAll(ctx)
}

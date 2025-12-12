package dbx_vec

import (
	"context"
	"fmt"
	"github.com/fengzhi09/golibx/jsonx"
	"github.com/fengzhi09/golibx/logx"
	"math"
	"sync"
)

func init() {
	ctx := context.Background()
	var _impl VecDB
	_impl = &QdrantDB{}
	logx.Debugf(ctx, "vecDBImplCheck: %T", _impl)
	_impl = &PgVecDB{}
	logx.Debugf(ctx, "vecDBImplCheck: %T", _impl)
	_impl = &MilvusDB{}
	logx.Debugf(ctx, "vecDBImplCheck: %T", _impl)
}

// FilterCondition 过滤条件
type FilterCondition struct {
	Field string `json:"field"`
	Op    string `json:"op"`
	Val   any    `json:"val"`
}

// VectorFilter 向量过滤条件
type VectorFilter struct {
	Field     string    `json:"field"`
	Val       []float64 `json:"val"`
	Threshold float64   `json:"threshold"`
}

// VecQuery 向量查询条件
type VecQuery struct {
	Table           string            `json:"table"`
	Topn            int               `json:"topn"`
	Filters         []FilterCondition `json:"filters"`
	FiltersVec      []VectorFilter    `json:"filters_vec"`
	IncludeMetadata bool              `json:"include_metadata"`
	IncludeVectors  bool              `json:"include_vectors"`
}

type VecNode struct {
	Id       string
	MetaData jsonx.JObj
	Vec      []float32
}
type MVecNode struct {
	Id       string
	MetaData jsonx.JObj
	VecMap   map[string][]float32
}
type ResNode struct {
	Id       string
	MetaData jsonx.JObj
	Vec      []float32
	Score    float32
	Rank     int
}
type DbApi interface {
	// Connect 建立数据库连接
	Connect(ctx context.Context) error
	// Close 关闭数据库连接
	Close(ctx context.Context) error
}

type MVecApi interface {
	// Upsert  插入多维向量数据
	UpsertM(ctx context.Context, nodes ...MVecNode) error
}
type VecApi interface {
	// Search 执行向量搜索
	Search(ctx context.Context, query VecQuery) ([]*ResNode, error)
	// Upsert  插入单维向量数据
	Upsert(ctx context.Context, nodes ...VecNode) error
	// Delete 删除向量数据
	// Delete(ctx context.Context,table string, dimension int, kwargs jsonx.JObj) (bool, error)
	// Insert 插入向量数据
	// Insert(ctx context.Context,table string, vectors [][]float64, ids []string, metadata []jsonx.JObj) (bool, error)
}
type TableApi interface {
	// NewTable 创建新表
	NewTable(ctx context.Context, name string, conf jsonx.JObj) error
}

// VecDB 向量数据库接口
type VecDB interface {
	DbApi
	TableApi
	VecApi
	MVecApi
}
type VecDBConf = jsonx.JObj

// VecDBBase 向量数据库基类
type VecDBBase struct {
	resName string
	resConf jsonx.JObj
}

// VecDBMgr 向量数据库管理器
type VecDBMgr struct {
	dbMap map[string]VecDB
	mutex sync.RWMutex
}

var vecDBMgrInstance *VecDBMgr
var vecDBMgrOnce sync.Once

// GetVecDBMgrInstance 获取向量数据库管理器实例
func GetVecDBMgrInstance() *VecDBMgr {
	vecDBMgrOnce.Do(func() {
		vecDBMgrInstance = &VecDBMgr{
			dbMap: make(map[string]VecDB),
		}
	})
	return vecDBMgrInstance
}

// Init 初始化向量数据库连接
func (vm *VecDBMgr) Init(ctx context.Context, dbs map[string]VecDBConf) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	for dbName, dbConf := range dbs {
		var db VecDB
		var err error

		if _, ok := dbConf["qdrant"]; ok {
			db, err = NewQdrantDB(dbName, dbConf)
			if err != nil {
				return fmt.Errorf("failed to create qdrant db: %v", err)
			}
		} else if _, ok := dbConf["pgvec"]; ok {
			db, err = NewPgVecDB(ctx, dbName, dbConf)
			if err != nil {
				return fmt.Errorf("failed to create pgvec db: %v", err)
			}
		} else if _, ok := dbConf["milvus"]; ok {
			db, err = NewMilvusDB(ctx, dbName, dbConf)
			if err != nil {
				return fmt.Errorf("failed to create qdrmilvusant db: %v", err)
			}
		} else {
			return fmt.Errorf("不支持的向量数据库配置: %v", dbConf)
		}

		// 建立连接
		if err := db.Connect(ctx); err != nil {
			return fmt.Errorf("failed to connect to vector database %s: %v", dbName, err)
		}

		vm.dbMap[dbName] = db
	}

	return nil
}

// GetDB 获取向量数据库客户端
func (vm *VecDBMgr) GetDB(dbName string) (VecDB, error) {
	vm.mutex.RLock()
	db, exists := vm.dbMap[dbName]
	vm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("未配置向量数据库: %s", dbName)
	}

	return db, nil
}

// CloseAll 关闭所有数据库连接
func (vm *VecDBMgr) CloseAll(ctx context.Context) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	for _, db := range vm.dbMap {
		_ = db.Close(ctx)
	}

	vm.dbMap = make(map[string]VecDB)
}

// 导出的全局函数

// InitVecDB 初始化全局向量数据库
func InitVecDB(ctx context.Context, dbs map[string]VecDBConf) error {
	return GetVecDBMgrInstance().Init(ctx, dbs)
}

// GetVecDB 获取全局向量数据库实例
func GetVecDB(name string) (VecDB, error) {
	return GetVecDBMgrInstance().GetDB(name)
}

// CloseAllVecDBs 关闭所有全局向量数据库连接
func CloseAllVecDBs(ctx context.Context) {
	GetVecDBMgrInstance().CloseAll(ctx)
}

// 工具函数：计算向量余弦相似度
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float64
	var normA, normB float64

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	normA = math.Sqrt(normA)
	normB = math.Sqrt(normB)

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (normA * normB)
}

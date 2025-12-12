# dbx

一个用于Golang的数据库操作库，为不同的数据库提供统一的接口。

## 特性

- **统一接口**: 为不同数据库类型提供一致的API
- **多数据库支持**: PostgreSQL、MySQL、Doris
- **向量数据库支持**: Milvus、PGVector、Qdrant（在dbx_vec子包中）
- **数据库管理器**: 轻松管理多个数据库连接
- **缓存支持**: 内存缓存和Redis缓存实现

## 使用方法

### 基本用法

```go
import (
    "context"
    "github.com/fengzhi09/golibx/dbx"
    "github.com/fengzhi09/golibx/jsonx"
)

func main() {
    ctx := context.Background()
    
    // 创建一个PostgreSQL数据库实例
    pgConf := jsonx.ParseObj([]byte(`{"db_type":"postgres","db_url":"postgres://user:pass@localhost:5432/dbname"}`))
    pgDB := dbx.NewPSQL("pgdb", pgConf)
    
    // 连接到数据库
    err := pgDB.Connect(ctx)
    if err != nil {
        panic(err)
    }
    defer pgDB.Close(ctx)
    
    // 查询数据
    results, err := pgDB.Query(ctx, "SELECT * FROM users WHERE id = $1", 1)
    if err != nil {
        panic(err)
    }
    
    // 使用结果
    for _, row := range results {
        fmt.Printf("User: %s\n", row.GetStr("name"))
    }
}
```

### 使用数据库管理器

```go
import (
    "context"
    "github.com/fengzhi09/golibx/dbx"
    "github.com/fengzhi09/golibx/jsonx"
)

func main() {
    ctx := context.Background()
    
    // 初始化多个数据库
    dbs := map[string]dbx.DBConf{
        "pgdb": jsonx.ParseObj([]byte(`{"db_type":"postgres","db_url":"postgres://user:pass@localhost:5432/dbname"}`)),
        "mydb": jsonx.ParseObj([]byte(`{"db_type":"mysql","db_url":"mysql://user:pass@localhost:3306/dbname"}`)),
    }
    
    err := dbx.InitDB(ctx, dbs)
    if err != nil {
        panic(err)
    }
    defer dbx.CloseDBs(ctx)
    
    // 获取数据库实例
    pgDB, err := dbx.UseDB("pgdb")
    if err != nil {
        panic(err)
    }
    
    // 查询数据
    results, err := pgDB.Query(ctx, "SELECT * FROM users")
    // ...
}
```

## 向量数据库支持(待完成)

`dbx_vec`子包提供了对向量数据库的支持：

- Milvus
- PGVector
- Qdrant

### 使用示例

```go
import (
    "context"
    "github.com/fengzhi09/golibx/dbx/dbx_vec"
    "github.com/fengzhi09/golibx/jsonx"
)

func main() {
    ctx := context.Background()
    
    // 创建一个Milvus向量数据库实例
    milvusConf := jsonx.ParseObj([]byte(`{"addr":"localhost:19530","collection_name":"test_collection"}`))
    milvus := dbx_vec.NewMilvus(milvusConf)
    
    // 连接到向量数据库
    err := milvus.Connect(ctx)
    if err != nil {
        panic(err)
    }
    defer milvus.Close(ctx)
    
    // 插入向量
    vectors := [][]float32{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}}
    ids, err := milvus.Insert(ctx, vectors)
    if err != nil {
        panic(err)
    }
    
    // 搜索向量
    queryVector := []float32{1.1, 2.1, 3.1}
    results, err := milvus.Search(ctx, queryVector, 2)
    if err != nil {
        panic(err)
    }
    
    // 使用搜索结果
    for _, result := range results {
        fmt.Printf("ID: %d, Score: %f\n", result.ID, result.Score)
    }
}
```

## 缓存支持

dbx包还提供了缓存支持，有两种实现：

- 内存缓存
- Redis缓存

### 使用示例

```go
import (
    "context"
    "github.com/fengzhi09/golibx/dbx"
    "github.com/fengzhi09/golibx/jsonx"
)

func main() {
    ctx := context.Background()
    
    // 创建一个内存缓存
    memCache := dbx.NewMemCache()
    
    // 设置值
    err := memCache.Set(ctx, "key", "value", 60) // 60秒后过期
    if err != nil {
        panic(err)
    }
    
    // 获取值
    value, err := memCache.Get(ctx, "key")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Value: %s\n", value)
    
    // 创建一个Redis缓存
    redisConf := jsonx.ParseObj([]byte(`{"addr":"localhost:6379","password":"","db":0}`))
    redisCache := dbx.NewRedisCache(redisConf)
    
    // 类似地使用Redis缓存
    // ...
}
```

## API参考

### ISQL接口

```go
type ISQL interface {
    Connect(ctx context.Context) error
    Query(ctx context.Context, query string, args ...any) ([]*jsonx.JObj, error)
    Close(ctx context.Context) error
}
```

### 数据库管理器函数

```go
// 初始化数据库
func InitDB(ctx context.Context, dbs map[string]DBConf) error

// 获取数据库实例
func UseDB(name string) (ISQL, error)

// 关闭所有数据库
func CloseDBs(ctx context.Context)
```

### 向量数据库接口

```go
type IVectorDB interface {
    Connect(ctx context.Context) error
    Insert(ctx context.Context, vectors [][]float32) ([]int64, error)
    Search(ctx context.Context, vector []float32, topK int) ([]*VectorResult, error)
    Close(ctx context.Context) error
}
```

### 缓存接口

```go
type ICache interface {
    Set(ctx context.Context, key, value string, expire int) error
    Get(ctx context.Context, key string) (string, error)
    Del(ctx context.Context, key string) error
    Close(ctx context.Context) error
}
```

## 许可证

MIT

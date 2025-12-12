# dbx

A database operation library for Golang, providing a unified interface for different databases.

## Features

- **Unified Interface**: Consistent API for different database types
- **Multiple Database Support**: PostgreSQL, MySQL, Doris
- **Vector Database Support**: Milvus, PGVector, Qdrant (in dbx_vec subpackage)
- **Database Manager**: Easy management of multiple database connections
- **Cache Support**: In-memory and Redis cache implementation

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/fengzhi09/golibx/dbx"
    "github.com/fengzhi09/golibx/jsonx"
)

func main() {
    ctx := context.Background()
    
    // Create a PostgreSQL database instance
    pgConf := jsonx.ParseObj([]byte(`{"db_type":"postgres","db_url":"postgres://user:pass@localhost:5432/dbname"}`))
    pgDB := dbx.NewPSQL("pgdb", pgConf)
    
    // Connect to the database
    err := pgDB.Connect(ctx)
    if err != nil {
        panic(err)
    }
    defer pgDB.Close(ctx)
    
    // Query data
    results, err := pgDB.Query(ctx, "SELECT * FROM users WHERE id = $1", 1)
    if err != nil {
        panic(err)
    }
    
    // Use the results
    for _, row := range results {
        fmt.Printf("User: %s\n", row.GetStr("name"))
    }
}
```

### Using Database Manager

```go
import (
    "context"
    "github.com/fengzhi09/golibx/dbx"
    "github.com/fengzhi09/golibx/jsonx"
)

func main() {
    ctx := context.Background()
    
    // Initialize multiple databases
    dbs := map[string]dbx.DBConf{
        "pgdb": jsonx.ParseObj([]byte(`{"db_type":"postgres","db_url":"postgres://user:pass@localhost:5432/dbname"}`)),
        "mydb": jsonx.ParseObj([]byte(`{"db_type":"mysql","db_url":"mysql://user:pass@localhost:3306/dbname"}`)),
    }
    
    err := dbx.InitDB(ctx, dbs)
    if err != nil {
        panic(err)
    }
    defer dbx.CloseDBs(ctx)
    
    // Get a database instance
    pgDB, err := dbx.UseDB("pgdb")
    if err != nil {
        panic(err)
    }
    
    // Query data
    results, err := pgDB.Query(ctx, "SELECT * FROM users")
    // ...
}
```

## Vector Database Support(DOING)

The `dbx_vec` subpackage provides support for vector databases:

- Milvus
- PGVector
- Qdrant

### Usage Example

```go
import (
    "context"
    "github.com/fengzhi09/golibx/dbx/dbx_vec"
    "github.com/fengzhi09/golibx/jsonx"
)

func main() {
    ctx := context.Background()
    
    // Create a Milvus vector database instance
    milvusConf := jsonx.ParseObj([]byte(`{"addr":"localhost:19530","collection_name":"test_collection"}`))
    milvus := dbx_vec.NewMilvus(milvusConf)
    
    // Connect to the vector database
    err := milvus.Connect(ctx)
    if err != nil {
        panic(err)
    }
    defer milvus.Close(ctx)
    
    // Insert vectors
    vectors := [][]float32{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}}
    ids, err := milvus.Insert(ctx, vectors)
    if err != nil {
        panic(err)
    }
    
    // Search vectors
    queryVector := []float32{1.1, 2.1, 3.1}
    results, err := milvus.Search(ctx, queryVector, 2)
    if err != nil {
        panic(err)
    }
    
    // Use the search results
    for _, result := range results {
        fmt.Printf("ID: %d, Score: %f\n", result.ID, result.Score)
    }
}
```

## Cache Support

The dbx package also provides cache support with two implementations:

- In-memory cache
- Redis cache

### Usage Example

```go
import (
    "context"
    "github.com/fengzhi09/golibx/dbx"
    "github.com/fengzhi09/golibx/jsonx"
)

func main() {
    ctx := context.Background()
    
    // Create an in-memory cache
    memCache := dbx.NewMemCache()
    
    // Set a value
    err := memCache.Set(ctx, "key", "value", 60) // Expires in 60 seconds
    if err != nil {
        panic(err)
    }
    
    // Get a value
    value, err := memCache.Get(ctx, "key")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Value: %s\n", value)
    
    // Create a Redis cache
    redisConf := jsonx.ParseObj([]byte(`{"addr":"localhost:6379","password":"","db":0}`))
    redisCache := dbx.NewRedisCache(redisConf)
    
    // Use Redis cache similarly
    // ...
}
```

## API Reference

### ISQL Interface

```go
type ISQL interface {
    Connect(ctx context.Context) error
    Query(ctx context.Context, query string, args ...any) ([]*jsonx.JObj, error)
    Close(ctx context.Context) error
}
```

### Database Manager Functions

```go
// Initialize databases
func InitDB(ctx context.Context, dbs map[string]DBConf) error

// Get a database instance
func UseDB(name string) (ISQL, error)

// Close all databases
func CloseDBs(ctx context.Context)
```

### Vector Database Interface

```go
type IVectorDB interface {
    Connect(ctx context.Context) error
    Insert(ctx context.Context, vectors [][]float32) ([]int64, error)
    Search(ctx context.Context, vector []float32, topK int) ([]*VectorResult, error)
    Close(ctx context.Context) error
}
```

### Cache Interface

```go
type ICache interface {
    Set(ctx context.Context, key, value string, expire int) error
    Get(ctx context.Context, key string) (string, error)
    Del(ctx context.Context, key string) error
    Close(ctx context.Context) error
}
```

## License

MIT

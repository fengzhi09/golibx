# gox

一个全面的Golang实用工具库，为常见编程任务提供广泛的辅助函数。

## 特性

- **数组操作**: 去重、子集、indexOf、in等
- **比较功能**: 支持多种类型的深度比较，包括JSON
- **转换工具**: 类型转换工具
- **文件操作**: 文件读写和操作
- **JSON处理**: 简单的JSON解析和操作
- **时间工具**: 时间格式化、解析和操作
- **字符串工具**: 随机字符串生成、正则表达式助手
- **映射工具**: LRU缓存实现
- **任务管理**: 简单的任务执行工具
- **语言助手**: 语言相关工具
- **CSV处理**: CSV文件读写
- **UTF-8支持**: UTF-8字符串工具

## 使用方法

### 数组操作

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // 数组中的唯一元素
    arr1 := []int{1, 2, 3, 2, 1}
    unique := gox.ArrUniq(arr1)
    fmt.Printf("唯一元素: %v\n", unique) // [1 2 3]
    
    // 数组的子集
    items := []string{"a", "b", "c", "d"}
    skips := []string{"b", "d"}
    subset := gox.ArrSub(items, skips)
    fmt.Printf("子集: %v\n", subset) // ["a", "c"]
    
    // 检查元素是否在数组中
    exists := gox.In("b", items...)
    fmt.Printf("存在: %v\n", exists) // true
    
    // 获取元素的索引
    index := gox.IndexOf(items, "c")
    fmt.Printf("索引: %d\n", index) // 2
}
```

### 比较功能

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // 比较两个值
    a := 10
    b := 20
    equal := gox.XEq(a, b)
    fmt.Printf("相等: %v\n", equal) // false
    
    // 比较JSON字符串
    json1 := `{"name":"test","age":10}`
    json2 := `{"age":10,"name":"test"}`
    jsonEqual := gox.CmpJSON(json1, json2)
    fmt.Printf("JSON相等: %v\n", jsonEqual) // true
    
    // 比较错误
    err1 := fmt.Errorf("error 1")
    err2 := fmt.Errorf("error 1")
    errEqual := gox.CmpErr(err1, err2)
    fmt.Printf("错误相等: %v\n", errEqual) // true
}
```

### 时间工具

```go
import (
    "fmt"
    "time"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // 格式化时间
    now := time.Now()
    formatted := gox.FormatTime(now, "2006-01-02 15:04:05")
    fmt.Printf("格式化: %s\n", formatted)
    
    // 解析时间
    parsed, err := gox.ParseTime("2023-01-01 12:00:00", "2006-01-02 15:04:05")
    if err == nil {
        fmt.Printf("解析: %v\n", parsed)
    }
    
    // 获取时间差
    later := now.Add(time.Hour)
    diff := gox.TimeDiff(now, later)
    fmt.Printf("时间差: %v\n", diff) // 1h0m0s
}
```

### 字符串工具

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // 生成随机字符串
    randStr := gox.RandStr(10)
    fmt.Printf("随机字符串: %s\n", randStr)
    
    // 检查字符串是否以指定前缀开头
    startsWith := gox.StrStartsWith("hello world", "hello")
    fmt.Printf("以hello开头: %v\n", startsWith) // true
    
    // 获取子字符串
    substr := gox.SubStr("hello world", 6, 5)
    fmt.Printf("子字符串: %s\n", substr) // "world"
}
```

### 文件操作

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // 读取文件内容
    content, err := gox.ReadFile("test.txt")
    if err == nil {
        fmt.Printf("内容: %s\n", content)
    }
    
    // 写入文件
    err = gox.WriteFile("output.txt", "Hello, World!")
    if err == nil {
        fmt.Println("文件写入成功")
    }
    
    // 检查文件是否存在
    exists := gox.FileExists("test.txt")
    fmt.Printf("文件存在: %v\n", exists)
}
```

### JSON处理

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // 解析JSON
    jsonStr := `{"name":"test","age":10}`
    var data map[string]interface{}
    err := gox.ParseJSON(jsonStr, &data)
    if err == nil {
        fmt.Printf("解析: %v\n", data)
    }
    
    // 格式化JSON
    formatted, err := gox.FormatJSON(data)
    if err == nil {
        fmt.Printf("格式化: %s\n", formatted)
    }
}
```

### 映射工具 (LRU缓存)

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // 创建容量为2的LRU缓存
    cache := gox.NewLRUCache(2)
    
    // 设置值
    cache.Put("key1", "value1")
    cache.Put("key2", "value2")
    cache.Put("key3", "value3") // 这将驱逐key1
    
    // 获取值
    val1, exists1 := cache.Get("key1")
    val2, exists2 := cache.Get("key2")
    val3, exists3 := cache.Get("key3")
    
    fmt.Printf("Key1: %v, 存在: %v\n", val1, exists1) // nil, false
    fmt.Printf("Key2: %v, 存在: %v\n", val2, exists2) // "value2", true
    fmt.Printf("Key3: %v, 存在: %v\n", val3, exists3) // "value3", true
}
```

## API参考

### 数组函数

```go
// 获取数组中的唯一元素
func ArrUniq[T comparable](arrs ...[]T) []T

// 获取数组的子集，移除指定元素
func ArrSub[T comparable](items, skips []T) []T

// 检查值是否在数组中
func In[T any](val T, arr ...T) bool

// 获取值在数组中的索引
func IndexOf[T any](a []T, b T) int
```

### 比较函数

```go
// 深度比较两个值
func XEq(a, b any) bool

// 比较两个JSON字符串
func CmpJSON(a, b string) bool

// 比较两个错误
func CmpErr(a, b error) bool
```

### 时间函数

```go
// 使用布局格式化时间
func FormatTime(t time.Time, layout string) string

// 从字符串解析时间
func ParseTime(s, layout string) (time.Time, error)

// 获取时间差
func TimeDiff(a, b time.Time) time.Duration
```

### 字符串函数

```go
// 生成随机字符串
func RandStr(length int) string

// 检查字符串是否以指定前缀开头
func StrStartsWith(s, prefix string) bool

// 获取子字符串
func SubStr(s string, start, length int) string
```

### 文件函数

```go
// 读取文件内容
func ReadFile(filePath string) (string, error)

// 将内容写入文件
func WriteFile(filePath, content string) error

// 检查文件是否存在
func FileExists(filePath string) bool
```

### 映射函数

```go
// 创建新的LRU缓存
func NewLRUCache(capacity int) *LRUCache

// 将值放入缓存
func (c *LRUCache) Put(key, value any)

// 从缓存获取值
func (c *LRUCache) Get(key any) (any, bool)
```

## 许可证

MIT

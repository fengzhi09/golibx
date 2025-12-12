# gox

A comprehensive utility library for Golang, providing a wide range of helper functions for common programming tasks.

## Features

- **Array Operations**: Unique, subset, indexOf, in, and more
- **Comparison**: Deep comparison for various types including JSON
- **Conversion**: Type conversion utilities
- **File Operations**: File reading, writing, and manipulation
- **JSON Handling**: Easy JSON parsing and manipulation
- **Time Utilities**: Time formatting, parsing, and manipulation
- **String Utilities**: Random string generation, regexp helpers
- **Map Utilities**: LRU cache implementation
- **Task Management**: Simple task execution utilities
- **Lang Helpers**: Language-related utilities
- **CSV Handling**: CSV file reading and writing
- **UTF-8 Support**: UTF-8 string utilities

## Usage

### Array Operations

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // Unique elements in array
    arr1 := []int{1, 2, 3, 2, 1}
    unique := gox.ArrUniq(arr1)
    fmt.Printf("Unique: %v\n", unique) // [1 2 3]
    
    // Subset of array
    items := []string{"a", "b", "c", "d"}
    skips := []string{"b", "d"}
    subset := gox.ArrSub(items, skips)
    fmt.Printf("Subset: %v\n", subset) // ["a", "c"]
    
    // Check if element is in array
    exists := gox.In("b", items...)
    fmt.Printf("Exists: %v\n", exists) // true
    
    // Get index of element
    index := gox.IndexOf(items, "c")
    fmt.Printf("Index: %d\n", index) // 2
}
```

### Comparison

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // Compare two values
    a := 10
    b := 20
    equal := gox.XEq(a, b)
    fmt.Printf("Equal: %v\n", equal) // false
    
    // Compare JSON strings
    json1 := `{"name":"test","age":10}`
    json2 := `{"age":10,"name":"test"}`
    jsonEqual := gox.CmpJSON(json1, json2)
    fmt.Printf("JSON Equal: %v\n", jsonEqual) // true
    
    // Compare errors
    err1 := fmt.Errorf("error 1")
    err2 := fmt.Errorf("error 1")
    errEqual := gox.CmpErr(err1, err2)
    fmt.Printf("Error Equal: %v\n", errEqual) // true
}
```

### Time Utilities

```go
import (
    "fmt"
    "time"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // Format time
    now := time.Now()
    formatted := gox.FormatTime(now, "2006-01-02 15:04:05")
    fmt.Printf("Formatted: %s\n", formatted)
    
    // Parse time
    parsed, err := gox.ParseTime("2023-01-01 12:00:00", "2006-01-02 15:04:05")
    if err == nil {
        fmt.Printf("Parsed: %v\n", parsed)
    }
    
    // Get time difference
    later := now.Add(time.Hour)
    diff := gox.TimeDiff(now, later)
    fmt.Printf("Diff: %v\n", diff) // 1h0m0s
}
```

### String Utilities

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // Generate random string
    randStr := gox.RandStr(10)
    fmt.Printf("Random: %s\n", randStr)
    
    // Check if string starts with
    startsWith := gox.StrStartsWith("hello world", "hello")
    fmt.Printf("Starts with: %v\n", startsWith) // true
    
    // Get substring
    substr := gox.SubStr("hello world", 6, 5)
    fmt.Printf("Substring: %s\n", substr) // "world"
}
```

### File Operations

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // Read file content
    content, err := gox.ReadFile("test.txt")
    if err == nil {
        fmt.Printf("Content: %s\n", content)
    }
    
    // Write to file
    err = gox.WriteFile("output.txt", "Hello, World!")
    if err == nil {
        fmt.Println("File written successfully")
    }
    
    // Check if file exists
    exists := gox.FileExists("test.txt")
    fmt.Printf("File exists: %v\n", exists)
}
```

### JSON Handling

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // Parse JSON
    jsonStr := `{"name":"test","age":10}`
    var data map[string]interface{}
    err := gox.ParseJSON(jsonStr, &data)
    if err == nil {
        fmt.Printf("Parsed: %v\n", data)
    }
    
    // Format JSON
    formatted, err := gox.FormatJSON(data)
    if err == nil {
        fmt.Printf("Formatted: %s\n", formatted)
    }
}
```

### Map Utilities (LRU Cache)

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/gox"
)

func main() {
    // Create LRU cache with capacity 2
    cache := gox.NewLRUCache(2)
    
    // Set values
    cache.Put("key1", "value1")
    cache.Put("key2", "value2")
    cache.Put("key3", "value3") // This will evict key1
    
    // Get values
    val1, exists1 := cache.Get("key1")
    val2, exists2 := cache.Get("key2")
    val3, exists3 := cache.Get("key3")
    
    fmt.Printf("Key1: %v, Exists: %v\n", val1, exists1) // nil, false
    fmt.Printf("Key2: %v, Exists: %v\n", val2, exists2) // "value2", true
    fmt.Printf("Key3: %v, Exists: %v\n", val3, exists3) // "value3", true
}
```

## API Reference

### Array Functions

```go
// Get unique elements from arrays
func ArrUniq[T comparable](arrs ...[]T) []T

// Get subset of array, removing specified elements
func ArrSub[T comparable](items, skips []T) []T

// Check if value is in array
func In[T any](val T, arr ...T) bool

// Get index of value in array
func IndexOf[T any](a []T, b T) int
```

### Comparison Functions

```go
// Deep compare two values
func XEq(a, b any) bool

// Compare two JSON strings
func CmpJSON(a, b string) bool

// Compare two errors
func CmpErr(a, b error) bool
```

### Time Functions

```go
// Format time with layout
func FormatTime(t time.Time, layout string) string

// Parse time from string
func ParseTime(s, layout string) (time.Time, error)

// Get time difference
func TimeDiff(a, b time.Time) time.Duration
```

### String Functions

```go
// Generate random string
func RandStr(length int) string

// Check if string starts with prefix
func StrStartsWith(s, prefix string) bool

// Get substring
func SubStr(s string, start, length int) string
```

### File Functions

```go
// Read file content
func ReadFile(filePath string) (string, error)

// Write content to file
func WriteFile(filePath, content string) error

// Check if file exists
func FileExists(filePath string) bool
```

### Map Functions

```go
// Create new LRU cache
func NewLRUCache(capacity int) *LRUCache

// Put value in cache
func (c *LRUCache) Put(key, value any)

// Get value from cache
func (c *LRUCache) Get(key any) (any, bool)
```

## License

MIT

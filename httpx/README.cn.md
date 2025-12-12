# httpx

一个功能强大的Golang HTTP客户端库，提供低级和高级API，支持钩子机制。

## 特性

- **双API设计**: 低级和高级API，适用于不同的使用场景
- **请求/响应钩子**: 用于请求和响应处理的自定义钩子
- **Web钩子**: 用于监控和日志记录HTTP请求的钩子
- **超时支持**: 可配置的请求超时
- **基础URL支持**: 为所有请求设置基础URL
- **头部管理**: 简单的头部配置
- **JSON支持**: 自动JSON序列化/反序列化
- **表单支持**: URL编码的表单数据支持
- **文件上传/下载**: 便捷的文件操作
- **灵活的选项**: 用于客户端配置的链式选项函数

## 使用方法

### 低级API

```go
import (
    "fmt"
    "net/url"
    "github.com/fengzhi09/golibx/httpx"
)

func main() {
    // 创建30秒超时的HTTP客户端
    client := httpx.NewHttp(30)
    
    // 配置客户端，设置基础URL和头部
    client = client.WithOpts(
        httpx.WithBaseURL("https://api.example.com"),
        httpx.WithHeaders(map[string]string{
            "Content-Type": "application/json",
            "Authorization": "Bearer token123",
        }),
    )
    
    // 添加请求钩子
    client = client.WithReqHooks(func(client httpx.Httpx, req *http.Request) error {
        fmt.Printf("请求URL: %s\n", req.URL)
        return nil
    })
    
    // 添加响应钩子
    client = client.WithRspHooks(func(client httpx.Httpx, req *http.Request, rsp *http.Response) error {
        fmt.Printf("响应状态: %s\n", rsp.Status)
        return nil
    })
    
    // 添加用于日志记录的web钩子
    client = client.WithHooks(httpx.WithLog(func(method, path, input, output string) {
        fmt.Printf("方法: %s, 路径: %s\n", method, path)
        fmt.Printf("请求: %s\n", input)
        fmt.Printf("响应: %s\n", output)
    }))
    
    // 发送GET请求
    resp, err := client.Do("GET", "/users", nil, nil)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    // 发送带有JSON体的POST请求
    body := map[string]string{
        "name": "test",
        "email": "test@example.com",
    }
    resp, err = client.Do("POST", "/users", body, nil)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    // 发送带有查询参数的请求
    params := url.Values{"page": {"1"}, "limit": {"10"}}
    resp, err = client.Do("GET", "/items", nil, params)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
}
```

### 高级API

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/httpx"
)

func main() {
    // 创建API客户端（实现取决于您的设置）
    apiClient := httpx.NewApiX() // 示例，实际实现可能有所不同
    
    // GET请求
    status, data, err := apiClient.Get("/users", nil)
    if err != nil {
        panic(err)
    }
    fmt.Printf("GET状态: %d, 数据: %v\n", status, data)
    
    // POST JSON请求
    body := map[string]string{
        "name": "test",
        "email": "test@example.com",
    }
    status, data, err = apiClient.PostJson("/users", body)
    if err != nil {
        panic(err)
    }
    fmt.Printf("POST状态: %d, 数据: %v\n", status, data)
    
    // POST表单请求
    formData := map[string]string{
        "username": "test",
        "password": "password123",
    }
    status, data, err = apiClient.PostForm("/login", formData)
    if err != nil {
        panic(err)
    }
    fmt.Printf("登录状态: %d, 数据: %v\n", status, data)
    
    // 上传文件
    file, err := os.Open("test.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    status, data, err = apiClient.Upload("/upload", file)
    if err != nil {
        panic(err)
    }
    fmt.Printf("上传状态: %d, 数据: %v\n", status, data)
    
    // 下载文件
    status, data, err = apiClient.Download("/download/file.txt", "/local/path/file.txt")
    if err != nil {
        panic(err)
    }
    fmt.Printf("下载状态: %d\n", status)
}
```

### 选项

```go
// 设置基础URL
httpx.WithBaseURL("https://api.example.com")

// 设置多个头部
httpx.WithHeaders(map[string]string{
    "Content-Type": "application/json",
})

// 设置单个头部
httpx.WithHeader("Authorization", "Bearer token123")

// 设置超时
httpx.WithTimeout(60) // 60秒

// 自定义HTTP客户端
httpx.WithClient(&http.Client{/* 自定义客户端 */})

// 指标钩子
httpx.WithMetric(func(method, path string, statusCode int, elapsedMs int64, err error) {
    // 记录指标
})

// 日志钩子
httpx.WithLog(func(method, path, input, output string) {
    // 记录请求/响应
})
```

## API参考

### 低级API

```go
// 创建新的HTTP客户端
func NewHttp(timeout int) Httpx

// Httpx接口
type Httpx interface {
    WithOpts(opts ...HttpOpt) Httpx
    WithReqHooks(hooks ...WebReqHook) Httpx
    WithRspHooks(hooks ...WebRspHook) Httpx
    WithHooks(hooks ...WebHook) Httpx
    Do(method string, path string, body any, params url.Values) (*http.Response, error)
}
```

### 高级API

```go
// ApiX接口
type ApiX interface {
    Get(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
    PostTxt(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
    PostForm(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
    PostJson(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
    Upload(path string, file *os.File, opts ...HttpOpt) (int, *jsonx.JObj, error)
    Download(urlPath string, savePath string, opts ...HttpOpt) (int, *jsonx.JObj, error)
}
```

### 选项函数

```go
// WithBaseURL为所有请求设置基础URL
func WithBaseURL(baseURL string) HttpOpt

// WithHeaders设置多个头部
func WithHeaders(headers map[string]string) HttpOpt

// WithHeader设置单个头部
func WithHeader(key, value string) HttpOpt

// WithTimeout设置请求超时
func WithTimeout(timeout int) HttpOpt

// WithClient设置自定义HTTP客户端
func WithClient(client *http.Client) HttpOpt

// WithMetric创建用于指标的web钩子
func WithMetric(writer func(method, path string, statusCode int, elapsedMs int64, err error)) WebHook

// WithLog创建用于日志的web钩子
func WithLog(writer func(method, path string, input string, output string)) WebHook
```

## 许可证

MIT

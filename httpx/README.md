# httpx

A powerful HTTP client library for Golang, providing both low-level and high-level APIs with hooks support.

## Features

- **Dual API Design**: Low-level and high-level APIs for different use cases
- **Request/Response Hooks**: Custom hooks for request and response processing
- **Web Hooks**: Hooks for monitoring and logging HTTP requests
- **Timeout Support**: Configurable request timeouts
- **Base URL Support**: Set a base URL for all requests
- **Header Management**: Easy header configuration
- **JSON Support**: Automatic JSON serialization/deserialization
- **Form Support**: URL-encoded form data support
- **File Upload/Download**: Convenient file operations
- **Flexible Options**: Chainable option functions for client configuration

## Usage

### Low-level API

```go
import (
    "fmt"
    "net/url"
    "github.com/fengzhi09/golibx/httpx"
)

func main() {
    // Create HTTP client with 30 seconds timeout
    client := httpx.NewHttp(30)
    
    // Configure client with base URL and headers
    client = client.WithOpts(
        httpx.WithBaseURL("https://api.example.com"),
        httpx.WithHeaders(map[string]string{
            "Content-Type": "application/json",
            "Authorization": "Bearer token123",
        }),
    )
    
    // Add request hook
    client = client.WithReqHooks(func(client httpx.Httpx, req *http.Request) error {
        fmt.Printf("Request URL: %s\n", req.URL)
        return nil
    })
    
    // Add response hook
    client = client.WithRspHooks(func(client httpx.Httpx, req *http.Request, rsp *http.Response) error {
        fmt.Printf("Response Status: %s\n", rsp.Status)
        return nil
    })
    
    // Add web hook for logging
    client = client.WithHooks(httpx.WithLog(func(method, path, input, output string) {
        fmt.Printf("Method: %s, Path: %s\n", method, path)
        fmt.Printf("Request: %s\n", input)
        fmt.Printf("Response: %s\n", output)
    }))
    
    // Send GET request
    resp, err := client.Do("GET", "/users", nil, nil)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    // Send POST request with JSON body
    body := map[string]string{
        "name": "test",
        "email": "test@example.com",
    }
    resp, err = client.Do("POST", "/users", body, nil)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    // Send request with query parameters
    params := url.Values{"page": {"1"}, "limit": {"10"}}
    resp, err = client.Do("GET", "/items", nil, params)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
}
```

### High-level API

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/httpx"
)

func main() {
    // Create API client (implementation depends on your setup)
    apiClient := httpx.NewApiX() // Example, actual implementation may vary
    
    // GET request
    status, data, err := apiClient.Get("/users", nil)
    if err != nil {
        panic(err)
    }
    fmt.Printf("GET Status: %d, Data: %v\n", status, data)
    
    // POST JSON request
    body := map[string]string{
        "name": "test",
        "email": "test@example.com",
    }
    status, data, err = apiClient.PostJson("/users", body)
    if err != nil {
        panic(err)
    }
    fmt.Printf("POST Status: %d, Data: %v\n", status, data)
    
    // POST form request
    formData := map[string]string{
        "username": "test",
        "password": "password123",
    }
    status, data, err = apiClient.PostForm("/login", formData)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Login Status: %d, Data: %v\n", status, data)
    
    // Upload file
    file, err := os.Open("test.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    status, data, err = apiClient.Upload("/upload", file)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Upload Status: %d, Data: %v\n", status, data)
    
    // Download file
    status, data, err = apiClient.Download("/download/file.txt", "/local/path/file.txt")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Download Status: %d\n", status)
}
```

### Options

```go
// Set base URL
httpx.WithBaseURL("https://api.example.com")

// Set headers
httpx.WithHeaders(map[string]string{
    "Content-Type": "application/json",
})

// Set single header
httpx.WithHeader("Authorization", "Bearer token123")

// Set timeout
httpx.WithTimeout(60) // 60 seconds

// Custom HTTP client
httpx.WithClient(&http.Client{/* custom client */})

// Metric hook
httpx.WithMetric(func(method, path string, statusCode int, elapsedMs int64, err error) {
    // Record metrics
})

// Log hook
httpx.WithLog(func(method, path, input, output string) {
    // Log request/response
})
```

## API Reference

### Low-level API

```go
// Create new HTTP client
func NewHttp(timeout int) Httpx

// Httpx interface
type Httpx interface {
    WithOpts(opts ...HttpOpt) Httpx
    WithReqHooks(hooks ...WebReqHook) Httpx
    WithRspHooks(hooks ...WebRspHook) Httpx
    WithHooks(hooks ...WebHook) Httpx
    Do(method string, path string, body any, params url.Values) (*http.Response, error)
}
```

### High-level API

```go
// ApiX interface
type ApiX interface {
    Get(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
    PostTxt(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
    PostForm(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
    PostJson(path string, body any, opts ...HttpOpt) (int, *jsonx.JObj, error)
    Upload(path string, file *os.File, opts ...HttpOpt) (int, *jsonx.JObj, error)
    Download(urlPath string, savePath string, opts ...HttpOpt) (int, *jsonx.JObj, error)
}
```

### Option Functions

```go
// WithBaseURL sets the base URL for all requests
func WithBaseURL(baseURL string) HttpOpt

// WithHeaders sets multiple headers
func WithHeaders(headers map[string]string) HttpOpt

// WithHeader sets a single header
func WithHeader(key, value string) HttpOpt

// WithTimeout sets the request timeout
func WithTimeout(timeout int) HttpOpt

// WithClient sets a custom HTTP client
func WithClient(client *http.Client) HttpOpt

// WithMetric creates a web hook for metrics
func WithMetric(writer func(method, path string, statusCode int, elapsedMs int64, err error)) WebHook

// WithLog creates a web hook for logging
func WithLog(writer func(method, path string, input string, output string)) WebHook
```

## License

MIT

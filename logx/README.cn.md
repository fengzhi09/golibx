# logx

一个灵活且强大的Golang日志库，支持多种日志级别、基于模块的日志记录和各种输出格式。

## 特性

- **多种日志级别**: DEBUG、INFO、WARN、ERROR、PANIC
- **基于模块的日志记录**: 日志消息可以与特定模块关联
- **灵活的输出**: 控制台和文件输出支持
- **日志旋转**: 自动日志文件旋转
- **上下文支持**: 上下文感知的日志记录
- **结构化日志**: 丰富的日志元数据，包括时间戳、级别和模块名称
- **易用性**: 用于常见日志任务的简单API
- **恢复支持**: 带有日志记录的自动panic恢复
- **可定制**: 可配置的日志格式和输出目标

## 使用方法

### 基本用法

```go
import (
    "context"
    "github.com/fengzhi09/golibx/logx"
)

func main() {
    ctx := context.Background()
    
    // 使用默认配置初始化日志器
    err := logx.Init("myapp", logx.INFO, nil)
    if err != nil {
        panic(err)
    }
    defer logx.CloseLogs(ctx)
    
    // 使用不同级别记录日志
    logx.Debugf(ctx, "调试消息: %v", 123)
    logx.Infof(ctx, "信息消息: %s", "test")
    logx.Warnf(ctx, "警告消息: %v", true)
    logx.Errorf(ctx, "错误消息: %v", fmt.Errorf("test error"))
    
    // 简单日志方法
    logx.Debug(ctx, "简单调试消息")
    logx.Info(ctx, "简单信息消息")
    logx.Warn(ctx, "简单警告消息")
    logx.Error(ctx, "简单错误消息")
}
```

### 基于模块的日志记录

```go
import (
    "context"
    "github.com/fengzhi09/golibx/logx"
)

func main() {
    ctx := context.Background()
    
    // 初始化日志器
    logx.Init("myapp", logx.INFO, nil)
    defer logx.CloseLogs(ctx)
    
    // 使用模块名称记录日志
    logx.DebugfM(ctx, "auth", "用户认证尝试: %s", "testuser")
    logx.InfofM(ctx, "db", "数据库连接已建立")
    logx.WarnfM(ctx, "http", "检测到慢请求: %s", "/api/slow")
    logx.ErrorfM(ctx, "worker", "任务失败: %v", fmt.Errorf("task error"))
    
    // 简单的模块日志方法
    logx.DebugM(ctx, "auth", "auth模块的简单调试消息")
    logx.InfoM(ctx, "db", "db模块的简单信息消息")
    logx.WarnM(ctx, "http", "http模块的简单警告消息")
    logx.ErrorM(ctx, "worker", "worker模块的简单错误消息")
}
```

### 配置

```go
import (
    "context"
    "github.com/fengzhi09/golibx/logx"
)

func main() {
    ctx := context.Background()
    
    // 配置日志旋转
    rotateConf := &logx.RotateConf{
        MaxSize:    10,  // 最大文件大小，单位MB
        MaxBackups: 5,   // 最大备份文件数量
        MaxAge:     30,  // 最大保留天数
        Compress:   true, // 压缩旧日志文件
    }
    
    // 使用配置初始化日志器
    err := logx.Init("myapp", logx.DEBUG, rotateConf)
    if err != nil {
        panic(err)
    }
    defer logx.CloseLogs(ctx)
    
    // 日志消息现在将写入控制台和带有旋转功能的文件
    logx.Infof(ctx, "日志器已使用旋转配置初始化")
}
```

### 恢复支持

```go
import (
    "context"
    "github.com/fengzhi09/golibx/logx"
)

func main() {
    ctx := context.Background()
    
    // 初始化日志器
    logx.Init("myapp", logx.INFO, nil)
    defer logx.CloseLogs(ctx)
    
    // 使用恢复中间件
    recovered := logx.Recovery(ctx, func() {
        // 这将引发panic
        panic("测试panic")
    })
    
    if recovered {
        logx.Info(ctx, "Panic已成功恢复")
    }
}
```

## 日志级别

```go
const (
    DEBUG LogLevel = iota // 调试级别
    INFO                  // 信息级别
    WARN                  // 警告级别
    ERROR                 // 错误级别
    PANIC                 // 紧急级别
)
```

## API参考

### 初始化

```go
// 初始化全局日志器
func Init(app string, level LogLevel, rotate *RotateConf) error

// 关闭所有日志器
func CloseLogs(ctx context.Context)
```

### 日志方法

```go
// 调试级别日志
func Debugf(ctx context.Context, format string, args ...any)
func Debug(ctx context.Context, format string, args ...any)

// 信息级别日志
func Infof(ctx context.Context, format string, args ...any)
func Info(ctx context.Context, message string)

// 警告级别日志
func Warnf(ctx context.Context, format string, args ...any)
func Warn(ctx context.Context, message string)

// 错误级别日志
func Errorf(ctx context.Context, format string, args ...any)
func Error(ctx context.Context, message string)
```

### 基于模块的日志方法

```go
// 带有模块的调试级别日志
func DebugfM(ctx context.Context, module string, format string, args ...any)
func DebugM(ctx context.Context, module string, message string)

// 带有模块的信息级别日志
func InfofM(ctx context.Context, module string, format string, args ...any)
func InfoM(ctx context.Context, module string, message string)

// 带有模块的警告级别日志
func WarnfM(ctx context.Context, module string, format string, args ...any)
func WarnM(ctx context.Context, module string, message string)

// 带有模块的错误级别日志
func ErrorfM(ctx context.Context, module string, format string, args ...any)
func ErrorM(ctx context.Context, module string, message string)
```

### 恢复

```go
// 从panic中恢复并记录日志
func Recovery(ctx context.Context, fn func()) bool
```

### 配置

```go
// 旋转配置

type RotateConf struct {
    Filename   string // 日志文件名
    MaxSize    int    // 最大文件大小，单位MB
    MaxBackups int    // 最大备份文件数量
    MaxAge     int    // 最大保留天数
    Compress   bool   // 压缩旧日志文件
}
```

## 许可证

MIT

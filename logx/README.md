# logx

A flexible and powerful logging library for Golang, supporting multiple log levels, module-based logging, and various output formats.

## Features

- **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR, PANIC
- **Module-based Logging**: Log messages can be associated with specific modules
- **Flexible Output**: Console and file output support
- **Log Rotation**: Automatic log file rotation
- **Context Support**: Context-aware logging
- **Structured Logging**: Rich log metadata including timestamps, levels, and module names
- **Easy to Use**: Simple API for common logging tasks
- **Recovery Support**: Automatic panic recovery with logging
- **Customizable**: Configurable log formats and output destinations

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/fengzhi09/golibx/logx"
)

func main() {
    ctx := context.Background()
    
    // Initialize logger with default configuration
    err := logx.Init("myapp", logx.INFO, nil)
    if err != nil {
        panic(err)
    }
    defer logx.CloseLogs(ctx)
    
    // Log messages with different levels
    logx.Debugf(ctx, "Debug message: %v", 123)
    logx.Infof(ctx, "Info message: %s", "test")
    logx.Warnf(ctx, "Warning message: %v", true)
    logx.Errorf(ctx, "Error message: %v", fmt.Errorf("test error"))
    
    // Simple log methods
    logx.Debug(ctx, "Simple debug message")
    logx.Info(ctx, "Simple info message")
    logx.Warn(ctx, "Simple warning message")
    logx.Error(ctx, "Simple error message")
}
```

### Module-based Logging

```go
import (
    "context"
    "github.com/fengzhi09/golibx/logx"
)

func main() {
    ctx := context.Background()
    
    // Initialize logger
    logx.Init("myapp", logx.INFO, nil)
    defer logx.CloseLogs(ctx)
    
    // Log with module names
    logx.DebugfM(ctx, "auth", "Authentication attempt for user: %s", "testuser")
    logx.InfofM(ctx, "db", "Database connection established")
    logx.WarnfM(ctx, "http", "Slow request detected: %s", "/api/slow")
    logx.ErrorfM(ctx, "worker", "Task failed: %v", fmt.Errorf("task error"))
    
    // Simple module log methods
    logx.DebugM(ctx, "auth", "Simple debug message for auth module")
    logx.InfoM(ctx, "db", "Simple info message for db module")
    logx.WarnM(ctx, "http", "Simple warning message for http module")
    logx.ErrorM(ctx, "worker", "Simple error message for worker module")
}
```

### Configuration

```go
import (
    "context"
    "github.com/fengzhi09/golibx/logx"
)

func main() {
    ctx := context.Background()
    
    // Configure log rotation
    rotateConf := &logx.RotateConf{
        MaxSize:    10,  // Maximum file size in MB
        MaxBackups: 5,   // Maximum number of backup files
        MaxAge:     30,  // Maximum age in days
        Compress:   true, // Compress old log files
    }
    
    // Initialize logger with configuration
    err := logx.Init("myapp", logx.DEBUG, rotateConf)
    if err != nil {
        panic(err)
    }
    defer logx.CloseLogs(ctx)
    
    // Log messages will now be written to both console and file with rotation
    logx.Infof(ctx, "Logger initialized with rotation configuration")
}
```

### Recovery Support

```go
import (
    "context"
    "github.com/fengzhi09/golibx/logx"
)

func main() {
    ctx := context.Background()
    
    // Initialize logger
    logx.Init("myapp", logx.INFO, nil)
    defer logx.CloseLogs(ctx)
    
    // Use recovery middleware
    recovered := logx.Recovery(ctx, func() {
        // This will panic
        panic("test panic")
    })
    
    if recovered {
        logx.Info(ctx, "Panic recovered successfully")
    }
}
```

## Log Levels

```go
const (
    DEBUG LogLevel = iota // Debug level
    INFO                  // Info level
    WARN                  // Warning level
    ERROR                 // Error level
    PANIC                 // Panic level
)
```

## API Reference

### Initialization

```go
// Initialize global logger
func Init(app string, level LogLevel, rotate *RotateConf) error

// Close all loggers
func CloseLogs(ctx context.Context)
```

### Log Methods

```go
// Debug level logging
func Debugf(ctx context.Context, format string, args ...any)
func Debug(ctx context.Context, format string, args ...any)

// Info level logging
func Infof(ctx context.Context, format string, args ...any)
func Info(ctx context.Context, message string)

// Warning level logging
func Warnf(ctx context.Context, format string, args ...any)
func Warn(ctx context.Context, message string)

// Error level logging
func Errorf(ctx context.Context, format string, args ...any)
func Error(ctx context.Context, message string)
```

### Module-based Log Methods

```go
// Debug level logging with module
func DebugfM(ctx context.Context, module string, format string, args ...any)
func DebugM(ctx context.Context, module string, message string)

// Info level logging with module
func InfofM(ctx context.Context, module string, format string, args ...any)
func InfoM(ctx context.Context, module string, message string)

// Warning level logging with module
func WarnfM(ctx context.Context, module string, format string, args ...any)
func WarnM(ctx context.Context, module string, message string)

// Error level logging with module
func ErrorfM(ctx context.Context, module string, format string, args ...any)
func ErrorM(ctx context.Context, module string, message string)
```

### Recovery

```go
// Recover from panic and log it
func Recovery(ctx context.Context, fn func()) bool
```

### Configuration

```go
// Rotate configuration

type RotateConf struct {
    Filename   string // Log file name
    MaxSize    int    // Maximum file size in MB
    MaxBackups int    // Maximum number of backup files
    MaxAge     int    // Maximum age in days
    Compress   bool   // Compress old log files
}
```

## License

MIT

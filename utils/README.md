# utils

A collection of utility functions and tools for Golang, providing various helper utilities for common programming tasks.

## Features

- **Event Bus**: In-memory event bus with synchronous and asynchronous support
- **Time Utilities**: Time manipulation functions for dates, months, and weeks
- **URL Query String Handling**: URL query parameter parsing and manipulation
- **Viper Configuration**: Enhanced Viper configuration utilities
- **Wait Group**: Custom wait group implementation
- **Chinese Text Reading**: Chinese text reading utilities

## Usage

### Event Bus

```go
import (
    "context"
    "github.com/fengzhi09/golibx/utils"
)

func main() {
    ctx := context.Background()
    
    // Create an event bus with queue size 100
    bus := utils.NewEventBus("test-bus", 100, nil)
    bus.Start()
    defer bus.Close()
    
    // Subscribe to events
    bus.Sub("sub1", func(ctx context.Context, step utils.EventStep, event *utils.MemEvent) error {
        fmt.Printf("Received event: %v, Type: %s, Data: %v\n", event.Id, event.Type, event.Data)
        return nil
    }, utils.AcceptAny())
    
    // Subscribe to specific event types
    bus.Sub("sub2", func(ctx context.Context, step utils.EventStep, event *utils.MemEvent) error {
        fmt.Printf("Received event of type %s: %v\n", event.Type, event.Data)
        return nil
    }, utils.AcceptTypes("event1", "event2"))
    
    // Publish events synchronously
    bus.PubSync(ctx, &utils.MemEvent{
        Type: "event1",
        Data: "test data 1",
    })
    
    // Publish events asynchronously
    bus.PubAsync(ctx, &utils.MemEvent{
        Type: "event2",
        Data: "test data 2",
    })
}
```

### Time Utilities

```go
import (
    "fmt"
    "time"
    "github.com/fengzhi09/golibx/utils"
)

func main() {
    now := time.Now()
    dateLayout := "2006-01-02"
    
    // Get start and end of day
    startOfDay := utils.GetStartDayTime(now, dateLayout)
    endOfDay := utils.GetEndDayTime(now, dateLayout)
    fmt.Printf("Start of day: %v\n", startOfDay)
    fmt.Printf("End of day: %v\n", endOfDay)
    
    // Get start and end of month
    startOfMonth := utils.GetMonthStart(now)
    endOfMonth := utils.GetMonthEnd(now)
    fmt.Printf("Start of month: %v\n", startOfMonth)
    fmt.Printf("End of month: %v\n", endOfMonth)
    
    // Get start and end of week
    startOfWeek := utils.GetWeekStart(now)
    endOfWeek := utils.GetWeekEnd(now)
    fmt.Printf("Start of week: %v\n", startOfWeek)
    fmt.Printf("End of week: %v\n", endOfWeek)
    
    // Get number of days in month
    daysInMonth := utils.GetMonthDays(now)
    fmt.Printf("Days in month: %d\n", daysInMonth)
    
    // Get Chinese weekday (1-7, 1 = Monday)
    weekdayCn := utils.GetWeekDayCn(now)
    fmt.Printf("Chinese weekday: %d\n", weekdayCn)
}
```

### URL Query String Handling

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/utils"
)

func main() {
    // Parse URL query string
    queryStr := "name=test&age=25&tags=tag1&tags=tag2"
    params, err := utils.ParseQueryString(queryStr)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Parsed params: %v\n", params)
    fmt.Printf("Name: %s\n", params.Get("name"))
    fmt.Printf("Age: %s\n", params.Get("age"))
    fmt.Printf("Tags: %v\n", params["tags"])
    
    // Build URL query string
    newParams := map[string][]string{
        "name": []string{"newtest"},
        "age":  []string{"30"},
        "tags": []string{"newtag1", "newtag2"},
    }
    newQueryStr := utils.BuildQueryString(newParams)
    fmt.Printf("Built query string: %s\n", newQueryStr)
}
```

### Viper Configuration

```go
import (
    "github.com/fengzhi09/golibx/utils"
    "github.com/spf13/viper"
)

func main() {
    // Create a viper instance with default configuration
    v := utils.NewViper("app", []string{"./config"}, "yaml")
    
    // Read configuration
    err := v.ReadInConfig()
    if err != nil {
        panic(err)
    }
    
    // Get configuration values
    appName := v.GetString("app.name")
    appPort := v.GetInt("app.port")
    
    fmt.Printf("App name: %s\n", appName)
    fmt.Printf("App port: %d\n", appPort)
}
```

## API Reference

### Event Bus

```go
// Create a new event bus
func NewEventBus(name string, queueSize int64, debug Debugger, steps ...EventStep) *MemEventBus

// Start the event bus
func (re *MemEventBus) Start()

// Close the event bus
func (re *MemEventBus) Close()

// Subscribe to events
func (re *MemEventBus) Sub(subName string, handle StepHandle, accept Acceptor)

// Unsubscribe from events
func (re *MemEventBus) Del(subName string)

// Publish event synchronously
func (re *MemEventBus) PubSync(ctx context.Context, event *MemEvent)

// Publish event asynchronously
func (re *MemEventBus) PubAsync(ctx context.Context, event *MemEvent)

// Event acceptance functions
func AcceptAny() Acceptor
func AcceptTypes(types ...EventType) Acceptor
func AcceptPattern(patterns ...EventType) Acceptor
```

### Time Utilities

```go
// Get start of day
func GetStartDayTime(start time.Time, dateLayout string) time.Time

// Get end of day
func GetEndDayTime(start time.Time, dateLayout string) time.Time

// Get start of month
func GetMonthStart(start time.Time) time.Time

// Get end of month
func GetMonthEnd(start time.Time) time.Time

// Get start of week
func GetWeekStart(start time.Time) time.Time

// Get end of week
func GetWeekEnd(start time.Time) time.Time

// Get number of days in month
func GetMonthDays(start time.Time) int

// Get Chinese weekday (1-7, 1 = Monday)
func GetWeekDayCn(start time.Time) int
```

## License

MIT

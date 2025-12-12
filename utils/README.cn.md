# utils

一个Golang实用工具函数和工具的集合，为常见编程任务提供各种辅助工具。

## 特性

- **事件总线**: 支持同步和异步的内存事件总线
- **时间工具**: 日期、月份和周的时间操作函数
- **URL查询字符串处理**: URL查询参数解析和操作
- **Viper配置**: 增强的Viper配置工具
- **等待组**: 自定义等待组实现
- **中文文本读取**: 中文文本读取工具

## 使用方法

### 事件总线

```go
import (
    "context"
    "github.com/fengzhi09/golibx/utils"
)

func main() {
    ctx := context.Background()
    
    // 创建一个队列大小为100的事件总线
    bus := utils.NewEventBus("test-bus", 100, nil)
    bus.Start()
    defer bus.Close()
    
    // 订阅所有事件
    bus.Sub("sub1", func(ctx context.Context, step utils.EventStep, event *utils.MemEvent) error {
        fmt.Printf("收到事件: %v, 类型: %s, 数据: %v\n", event.Id, event.Type, event.Data)
        return nil
    }, utils.AcceptAny())
    
    // 订阅特定类型的事件
    bus.Sub("sub2", func(ctx context.Context, step utils.EventStep, event *utils.MemEvent) error {
        fmt.Printf("收到类型为 %s 的事件: %v\n", event.Type, event.Data)
        return nil
    }, utils.AcceptTypes("event1", "event2"))
    
    // 同步发布事件
    bus.PubSync(ctx, &utils.MemEvent{
        Type: "event1",
        Data: "测试数据 1",
    })
    
    // 异步发布事件
    bus.PubAsync(ctx, &utils.MemEvent{
        Type: "event2",
        Data: "测试数据 2",
    })
}
```

### 时间工具

```go
import (
    "fmt"
    "time"
    "github.com/fengzhi09/golibx/utils"
)

func main() {
    now := time.Now()
    dateLayout := "2006-01-02"
    
    // 获取一天的开始和结束时间
    startOfDay := utils.GetStartDayTime(now, dateLayout)
    endOfDay := utils.GetEndDayTime(now, dateLayout)
    fmt.Printf("一天的开始: %v\n", startOfDay)
    fmt.Printf("一天的结束: %v\n", endOfDay)
    
    // 获取一个月的开始和结束时间
    startOfMonth := utils.GetMonthStart(now)
    endOfMonth := utils.GetMonthEnd(now)
    fmt.Printf("一个月的开始: %v\n", startOfMonth)
    fmt.Printf("一个月的结束: %v\n", endOfMonth)
    
    // 获取一周的开始和结束时间
    startOfWeek := utils.GetWeekStart(now)
    endOfWeek := utils.GetWeekEnd(now)
    fmt.Printf("一周的开始: %v\n", startOfWeek)
    fmt.Printf("一周的结束: %v\n", endOfWeek)
    
    // 获取一个月的天数
    daysInMonth := utils.GetMonthDays(now)
    fmt.Printf("一个月的天数: %d\n", daysInMonth)
    
    // 获取中文星期几 (1-7, 1 = 星期一)
    weekdayCn := utils.GetWeekDayCn(now)
    fmt.Printf("中文星期几: %d\n", weekdayCn)
}
```

### URL查询字符串处理

```go
import (
    "fmt"
    "github.com/fengzhi09/golibx/utils"
)

func main() {
    // 解析URL查询字符串
    queryStr := "name=test&age=25&tags=tag1&tags=tag2"
    params, err := utils.ParseQueryString(queryStr)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("解析后的参数: %v\n", params)
    fmt.Printf("Name: %s\n", params.Get("name"))
    fmt.Printf("Age: %s\n", params.Get("age"))
    fmt.Printf("Tags: %v\n", params["tags"])
    
    // 构建URL查询字符串
    newParams := map[string][]string{
        "name": []string{"newtest"},
        "age":  []string{"30"},
        "tags": []string{"newtag1", "newtag2"},
    }
    newQueryStr := utils.BuildQueryString(newParams)
    fmt.Printf("构建的查询字符串: %s\n", newQueryStr)
}
```

### Viper配置

```go
import (
    "github.com/fengzhi09/golibx/utils"
    "github.com/spf13/viper"
)

func main() {
    // 创建一个带有默认配置的viper实例
    v := utils.NewViper("app", []string{"./config"}, "yaml")
    
    // 读取配置
    err := v.ReadInConfig()
    if err != nil {
        panic(err)
    }
    
    // 获取配置值
    appName := v.GetString("app.name")
    appPort := v.GetInt("app.port")
    
    fmt.Printf("应用名称: %s\n", appName)
    fmt.Printf("应用端口: %d\n", appPort)
}
```

## API参考

### 事件总线

```go
// 创建一个新的事件总线
func NewEventBus(name string, queueSize int64, debug Debugger, steps ...EventStep) *MemEventBus

// 启动事件总线
func (re *MemEventBus) Start()

// 关闭事件总线
func (re *MemEventBus) Close()

// 订阅事件
func (re *MemEventBus) Sub(subName string, handle StepHandle, accept Acceptor)

// 取消订阅事件
func (re *MemEventBus) Del(subName string)

// 同步发布事件
func (re *MemEventBus) PubSync(ctx context.Context, event *MemEvent)

// 异步发布事件
func (re *MemEventBus) PubAsync(ctx context.Context, event *MemEvent)

// 事件接受函数
func AcceptAny() Acceptor
func AcceptTypes(types ...EventType) Acceptor
func AcceptPattern(patterns ...EventType) Acceptor
```

### 时间工具

```go
// 获取一天的开始时间
func GetStartDayTime(start time.Time, dateLayout string) time.Time

// 获取一天的结束时间
func GetEndDayTime(start time.Time, dateLayout string) time.Time

// 获取一个月的开始时间
func GetMonthStart(start time.Time) time.Time

// 获取一个月的结束时间
func GetMonthEnd(start time.Time) time.Time

// 获取一周的开始时间
func GetWeekStart(start time.Time) time.Time

// 获取一周的结束时间
func GetWeekEnd(start time.Time) time.Time

// 获取一个月的天数
func GetMonthDays(start time.Time) int

// 获取中文星期几 (1-7, 1 = 星期一)
func GetWeekDayCn(start time.Time) int
```

## 许可证

MIT

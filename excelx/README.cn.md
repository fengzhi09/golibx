# excelx

一个功能强大的Golang Excel文件处理库，支持XLSX、XLS和CSV格式。

## 特性

- **多格式支持**: XLSX、XLS和CSV文件
- **统一API**: 不同文件格式的一致接口
- **流处理**: 高效的基于流的读写
- **易用性**: 用于常见Excel操作的简单API
- **文件转换**: 不同Excel格式之间的转换
- **基于时间的行拆分**: 根据时间列拆分行

## 使用方法

### 基本用法

```go
import (
    "context"
    "fmt"
    "os"
    "github.com/fengzhi09/golibx/excelx"
)

func main() {
    ctx := context.Background()
    
    // 打开Excel文件
    file, err := os.Open("example.xlsx")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    fileInfo, err := file.Stat()
    if err != nil {
        panic(err)
    }
    
    // 读取Excel文件
    xlsxFile, err := excelx.AsXlsx("example.xlsx", file, fileInfo.Size())
    if err != nil {
        panic(err)
    }
    defer xlsxFile.Release(ctx)
    
    // 获取所有工作表
    sheets := xlsxFile.GetSheetList()
    fmt.Printf("工作表: %v\n", sheets)
    
    // 从第一个工作表获取行
    rows, err := xlsxFile.GetRows(sheets[0])
    if err != nil {
        panic(err)
    }
    
    // 打印行
    for _, row := range rows {
        fmt.Printf("行: %v\n", row)
    }
}
```

### 读取不同文件格式

```go
import (
    "context"
    "fmt"
    "os"
    "github.com/fengzhi09/golibx/excelx"
)

func main() {
    ctx := context.Background()
    
    // 读取XLSX文件
    xlsxFile, _ := os.Open("example.xlsx")
    xlsxInfo, _ := xlsxFile.Stat()
    xlsx, _ := excelx.ReadXlsx(xlsxFile, xlsxInfo.Size())
    defer xlsx.Release(ctx)
    
    // 读取XLS文件
    xlsFile, _ := os.Open("example.xls")
    xlsInfo, _ := xlsFile.Stat()
    xls, _ := excelx.ReadXls(xlsFile, xlsInfo.Size())
    defer xls.Release(ctx)
    
    // 读取CSV文件
    csvFile, _ := os.Open("example.csv")
    csvInfo, _ := csvFile.Stat()
    csv, _ := excelx.ReadCsv(csvFile, csvInfo.Size())
    defer csv.Release(ctx)
    
    fmt.Println("所有文件读取成功")
}
```

### 按时间拆分行

```go
import (
    "context"
    "fmt"
    "os"
    "github.com/fengzhi09/golibx/excelx"
)

func main() {
    // 打开Excel文件
    file, _ := os.Open("time_data.xlsx")
    defer file.Close()
    
    fileInfo, _ := file.Stat()
    xlsxFile, _ := excelx.ReadXlsx(file, fileInfo.Size())
    
    // 按时间列拆分行
    columns, data := excelx.SplitRowsByTime(xlsxFile, []string{"date"}, 0)
    
    // 打印结果
    fmt.Printf("列: %v\n", columns)
    for date, rows := range data {
        fmt.Printf("日期: %s, 行数: %d\n", date, len(rows))
    }
}
```

## API参考

### 文件读取函数

```go
// 读取XLSX文件
func ReadXlsx(file IFile, size int64) (*XlsxFile, error)

// 读取XLS文件
func ReadXls(file IFile, size int64) (*XlsxFile, error)

// 读取CSV文件
func ReadCsv(file IFile, size int64) (*XlsxFile, error)

// 自动检测文件类型并读取
func AsXlsx(fileName string, file IFile, size int64) (*XlsxFile, error)
```

### XlsxFile方法

```go
// 释放资源
func (xf *XlsxFile) Release(ctx context.Context)

// 获取工作表列表
func (xf *XlsxFile) GetSheetList() []string

// 从工作表获取行
func (xf *XlsxFile) GetRows(sheetName string) ([][]string, error)
```

### 实用函数

```go
// 按时间列拆分行
func SplitRowsByTime(xlsxFile *XlsxFile, splitColNames []string, baserow int) ([]string, map[string][][]string)

// 将字符串切片转换为接口切片
func StringS2InterfaceS(src []string) []any
```

## 许可证

MIT

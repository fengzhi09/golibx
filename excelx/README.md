# excelx

A powerful Excel file processing library for Golang, supporting XLSX, XLS, and CSV formats.

## Features

- **Multi-format Support**: XLSX, XLS, and CSV files
- **Unified API**: Consistent interface for different file formats
- **Stream Processing**: Efficient stream-based reading and writing
- **Easy-to-use**: Simple API for common Excel operations
- **File Conversion**: Convert between different Excel formats
- **Time-based Row Splitting**: Split rows based on time columns

## Usage

### Basic Usage

```go
import (
    "context"
    "fmt"
    "os"
    "github.com/fengzhi09/golibx/excelx"
)

func main() {
    ctx := context.Background()
    
    // Open an Excel file
    file, err := os.Open("example.xlsx")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    fileInfo, err := file.Stat()
    if err != nil {
        panic(err)
    }
    
    // Read the Excel file
    xlsxFile, err := excelx.AsXlsx("example.xlsx", file, fileInfo.Size())
    if err != nil {
        panic(err)
    }
    defer xlsxFile.Release(ctx)
    
    // Get all sheets
    sheets := xlsxFile.GetSheetList()
    fmt.Printf("Sheets: %v\n", sheets)
    
    // Get rows from the first sheet
    rows, err := xlsxFile.GetRows(sheets[0])
    if err != nil {
        panic(err)
    }
    
    // Print rows
    for _, row := range rows {
        fmt.Printf("Row: %v\n", row)
    }
}
```

### Reading Different File Formats

```go
import (
    "context"
    "fmt"
    "os"
    "github.com/fengzhi09/golibx/excelx"
)

func main() {
    ctx := context.Background()
    
    // Read XLSX file
    xlsxFile, _ := os.Open("example.xlsx")
    xlsxInfo, _ := xlsxFile.Stat()
    xlsx, _ := excelx.ReadXlsx(xlsxFile, xlsxInfo.Size())
    defer xlsx.Release(ctx)
    
    // Read XLS file
    xlsFile, _ := os.Open("example.xls")
    xlsInfo, _ := xlsFile.Stat()
    xls, _ := excelx.ReadXls(xlsFile, xlsInfo.Size())
    defer xls.Release(ctx)
    
    // Read CSV file
    csvFile, _ := os.Open("example.csv")
    csvInfo, _ := csvFile.Stat()
    csv, _ := excelx.ReadCsv(csvFile, csvInfo.Size())
    defer csv.Release(ctx)
    
    fmt.Println("All files read successfully")
}
```

### Splitting Rows by Time

```go
import (
    "context"
    "fmt"
    "os"
    "github.com/fengzhi09/golibx/excelx"
)

func main() {
    // Open an Excel file
    file, _ := os.Open("time_data.xlsx")
    defer file.Close()
    
    fileInfo, _ := file.Stat()
    xlsxFile, _ := excelx.ReadXlsx(file, fileInfo.Size())
    
    // Split rows by time column
    columns, data := excelx.SplitRowsByTime(xlsxFile, []string{"date"}, 0)
    
    // Print results
    fmt.Printf("Columns: %v\n", columns)
    for date, rows := range data {
        fmt.Printf("Date: %s, Rows: %d\n", date, len(rows))
    }
}
```

## API Reference

### File Reading Functions

```go
// Read XLSX file
func ReadXlsx(file IFile, size int64) (*XlsxFile, error)

// Read XLS file  
func ReadXls(file IFile, size int64) (*XlsxFile, error)

// Read CSV file
func ReadCsv(file IFile, size int64) (*XlsxFile, error)

// Auto-detect file type and read
func AsXlsx(fileName string, file IFile, size int64) (*XlsxFile, error)
```

### XlsxFile Methods

```go
// Release resources
func (xf *XlsxFile) Release(ctx context.Context)

// Get sheet list
func (xf *XlsxFile) GetSheetList() []string

// Get rows from a sheet
func (xf *XlsxFile) GetRows(sheetName string) ([][]string, error)
```

### Utility Functions

```go
// Split rows by time columns
func SplitRowsByTime(xlsxFile *XlsxFile, splitColNames []string, baserow int) ([]string, map[string][][]string)

// Convert string slice to interface slice
func StringS2InterfaceS(src []string) []any
```

## License

MIT

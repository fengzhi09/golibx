package excelx

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"
	"github.com/fengzhi09/golibx/logx"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

var acceptedTypes = map[string]bool{
	".xlsx": true,
	".xls":  true,
	".csv":  true,
}

var (
	maxFileSize int64 = 300 << 20

	ErrFileParam    = errors.New("file param err")
	ErrFileOverSize = errors.New("file over size")
	ErrFileType     = errors.New("unsupported file type")
)

func GetFormFile(c *gin.Context, field string) (multipart.File, int64, error) {
	fileHeader, err := c.FormFile(field)
	if err != nil {
		logx.Warnf(context.Background(), "get form file failed: %s", err)
		return nil, 0, ErrFileParam
	}

	if fileHeader.Size > maxFileSize {
		logx.Warnf(context.Background(), "file oversize: %d", fileHeader.Size)
		return nil, 0, ErrFileOverSize

	}

	ext := path.Ext(fileHeader.Filename)
	if !acceptedTypes[ext] {
		logx.Warnf(context.Background(), "unsupported file type: %s", fileHeader.Filename)
		return nil, 0, ErrFileType

	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, 0, err
	}
	return file, fileHeader.Size, nil
}

func ExportXlsxFile(c *gin.Context, fileName string, file *XlsxFile) error {
	fileBuffer := bytes.NewBuffer(nil)
	err := file.Write(fileBuffer)
	if err != nil {
		return err
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet; charset=UTF-8")
	c.Header("Accept-Length", strconv.Itoa(fileBuffer.Len()))
	_, err = c.Writer.Write(fileBuffer.Bytes())
	return err
}

func ReadExcel(ctx context.Context, fileName string) (*XlsxFile, error) {
	fInfo := gox.NewFilePath(fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		cErr := file.Close()
		if cErr != nil {
			logx.Errorf(ctx, "文件未正确关闭，请查看:%v", cErr.Error())
		}
	}(file)
	fileInfo, _ := file.Stat()
	fType := strings.ToLower(fInfo.Type)
	switch fType {
	case "xlsx":
		return ReadXlsx(file, fileInfo.Size())
	case "xls":
		return ReadXls(file, fileInfo.Size())
	case "csv":
		return ReadCsv(file, fileInfo.Size())
	default:
		return nil, fmt.Errorf("不支持的文件类型:%s", fType)

	}
}

func Excel2Records(ctx context.Context, fileName string) ([]*Record, error) {
	file, err := ReadExcel(ctx, fileName)
	if err != nil {
		return nil, err
	}
	return GetRecords(file, 0), nil
}

func Record2Json(ctx context.Context, record *Record) (*jsonx.JObj, error) {
	raw := &jsonx.JObj{}
	var err error = nil
	emptyCnt := 0
	if len(record.Columns) == 0 {
		return raw, fmt.Errorf("空行")
	}
	for col, cell := range record.Data {
		if cell == "" {
			emptyCnt++
		}
		if col >= len(record.Columns) {
			return raw, fmt.Errorf("太多列")
		}
		column := record.Columns[col]
		column = strings.ToLower(column)
		raw.PutStr(column, cell)
	}
	if emptyCnt == len(record.Columns) {
		return raw, fmt.Errorf("空行")
	}
	return raw, err
}

func Record2Str(ctx context.Context, record *Record) string {
	raw, err := Record2Json(ctx, record)
	if err != nil {
		return err.Error()
	}
	return raw.String()
}

func Records2Excel(ctx context.Context, records []*Record, fileName string) error {
	excelFile := excelize.NewFile()
	columnMapper := map[string]int{}
	var columns []string
	for _, record := range records {
		for idx, column := range record.Columns {
			if _, ok := columnMapper[column]; ok {
				continue
			}
			columnMapper[column] = idx
			columns = append(columns, column)
		}
	}
	streamSheet1, err := excelFile.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}

	err = streamSheet1.SetRow("A1", StringS2InterfaceS(columns))
	if err != nil {
		return err
	}
	for idx, record := range records {
		row := idx + 2
		rowErr := streamSheet1.SetRow(fmt.Sprintf("A%v", row), StringS2InterfaceS(record.Data))
		if rowErr != nil {
			logx.Warnf(ctx, "写入行 %v,失败: %v", row, rowErr)
			continue
		}
	}
	err = streamSheet1.Flush()
	if err != nil {
		return err
	}
	return excelFile.SaveAs(fileName + ".xlsx")
}

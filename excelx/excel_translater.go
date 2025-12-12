package excelx

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/logx"
	"github.com/fengzhi09/golibx/utils"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2" // https://github.com/qax-os/excelize
)

type IFile interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type (
	excelReader func(file IFile, size int64) (*XlsxFile, error)
	fileType    = string
)

var (
	typeXlsx = "xlsx"
	typeXls  = "xls"
	typeCsv  = "csv"
)
var Readers = map[fileType]excelReader{typeXlsx: ReadXlsx, typeXls: ReadXls, typeCsv: ReadCsv}

func StringS2InterfaceS(src []string) []any {
	ret := make([]any, len(src))
	for i, v := range src {
		ret[i] = v
	}
	return ret
}

type XlsxFile struct {
	*excelize.File
	Raw     IFile
	cleaner func() error
}

func (xf *XlsxFile) Release(ctx context.Context) {
	if xf != nil && xf.Raw != nil {
		gox.SafeClose(ctx, "xlsx_file.raw", xf.Raw.Close)
	}
	if xf != nil && xf.File != nil {
		gox.SafeClose(ctx, "xlsx_file.excel", xf.File.Close)
	}
	if xf != nil && xf.cleaner != nil {
		err := xf.cleaner()
		if err != nil {
			logx.Errorf(ctx, "xlsx_file.cleaner failed, err:%v", err)
		}
	}
}

func AsXlsx(fileName string, file IFile, size int64) (*XlsxFile, error) {
	nameParts := strings.Split(fileName, ".")
	fileType_ := nameParts[len(nameParts)-1]
	if reader, hit := Readers[fileType_]; hit {
		return reader(file, size)
	}
	return nil, fmt.Errorf("not support file type:%v", fileType_)
}

func ReadCsv(file IFile, size int64) (*XlsxFile, error) {
	reader := csv.NewReader(utils.GbkReader(file))
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	csvData, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("err:%v", err)
	}
	ctx, xlsxFile := context.Background(), excelize.NewFile()
	path := fmt.Sprintf("%v-csv.xlsx", gox.NewOID().Hex())
	defer gox.CloseRes(ctx, path, xlsxFile)
	if err = csvFlushToFile(csvData, xlsxFile); err != nil {
		return nil, err
	}
	if err = xlsxFile.SaveAs(path); err != nil {
		return nil, err
	}
	return readXlsxWithDel(path)
}

func csvFlushToFile(records [][]string, xlsxFile *excelize.File) error {
	streamWriter, err := xlsxFile.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}
	for rowNum, row := range records {
		axis := fmt.Sprintf("A%d", rowNum+1)
		err = streamWriter.SetRow(axis, StringS2InterfaceS(row))
		if err != nil {
			return err
		}
	}
	if err = streamWriter.Flush(); err != nil {
		return fmt.Errorf("结束流式写入失败: %v", err)
	}
	return nil
}

func ReadXls(file IFile, size int64) (*XlsxFile, error) {
	xlFile, err := xls.OpenReader(file, "utf-8")
	if err != nil {
		return nil, fmt.Errorf("open_err:%v", err)
	}
	xlsSheet1 := xlFile.GetSheet(0)
	ctx, xlsxFile := context.Background(), excelize.NewFile()
	path := fmt.Sprintf("%v-xls.xlsx", gox.NewOID().Hex())
	defer gox.CloseRes(ctx, path, xlsxFile)
	if err = xlsFlushToFile(xlsSheet1, xlsxFile); err != nil {
		return nil, err
	}
	if err = xlsxFile.SaveAs(path); err != nil {
		return nil, err
	}
	return readXlsxWithDel(path)
}

func xlsFlushToFile(sheet *xls.WorkSheet, xlsxFile *excelize.File) error {
	streamWriter, err := xlsxFile.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}
	for i := 0; i < int(sheet.MaxRow); i++ {
		axis, xlsData := fmt.Sprintf("A%d", i+1), sheet.Row(i)
		err = copyRowFromXls(xlsData, axis, streamWriter)
		if err != nil {
			return err
		}
	}
	if err = streamWriter.Flush(); err != nil {
		return fmt.Errorf("结束流式写入失败: %v", err)
	}
	return err
}

func readXlsxWithDel(path string) (*XlsxFile, error) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	return &XlsxFile{file, nil, func() error { return os.Remove(path) }}, nil
}

func copyRowFromXls(rawData *xls.Row, axis string, streamWriter *excelize.StreamWriter) error {
	columnCnt := rawData.LastCol()
	targetRow := make([]any, 0)
	for j := 0; j < columnCnt; j++ {
		targetRow = append(targetRow, rawData.Col(j))
	}
	return streamWriter.SetRow(axis, targetRow)
}

func ReadXlsx(file IFile, size int64) (*XlsxFile, error) {
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	return &XlsxFile{xlsx, file, nil}, nil
}

func SplitRowsByTime(xlsxFile *XlsxFile, splitColNames []string, baserow int) ([]string, map[string][][]string) {
	defer xlsxFile.Release(context.Background())
	allColumns := make([]string, 0)
	timeColIndex := make([]int, 0)
	result := make(map[string][][]string)
	matrixXlsx, _ := xlsxFile.GetRows(xlsxFile.GetSheetList()[0])
	for t, row := range matrixXlsx {
		if t < baserow {
			continue
		} else if t == baserow {
			for colNum, cell := range row {
				cache := strings.Trim(cell, " ")
				allColumns = append(allColumns, cache)
				if gox.In(cache, splitColNames...) {
					timeColIndex = append(timeColIndex, colNum)
				}
			}
		} else {
			if len(timeColIndex) == 0 {
				return nil, nil
			}
			RowSplitString := ""
			for _, colIndex := range timeColIndex {
				RowSplitString += strings.Trim(row[colIndex], " ")
			}
			result[RowSplitString] = append(result[RowSplitString], row)
		}
	}
	return allColumns, result
}

func writeExcel(data [][]string, colNames []string, targetFile *excelize.File) error {
	if len(data) == 0 {
		return nil
	}
	streamWriter, err := targetFile.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}
	streamWriter.SetRow("A1", StringS2InterfaceS(colNames))
	for rowNum, copyRow := range data {
		axis := fmt.Sprintf("A%v", rowNum+2)
		streamWriter.SetRow(axis, StringS2InterfaceS(copyRow))
	}
	err = streamWriter.Flush()
	return err
}

func preTrim(src []string) []string {
	target := make([]string, 0)
	for _, value := range src {
		target = append(target, strings.TrimSpace(value))
	}
	return target
}

func simplyExecute(path, baseColumns, savePath, targetFileName string) error {
	columns := strings.Split(baseColumns, ",")

	file, _ := os.Open(path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("文件未正确关闭，请查看:%v", err.Error())
		}
	}(file)
	fileInfo, _ := file.Stat()
	xlsxFile, err := AsXlsx(path, file, fileInfo.Size())
	if err != nil {
		return err
	}

	if xlsxFile == nil {
		return fmt.Errorf("生成的xlsx读取失败，已结束")
	} else {
		colName, data := SplitRowsByTime(xlsxFile, columns, 0)
		targetFile := excelize.NewFile()
		for dataTime, dt := range data {
			err = writeExcel(dt, preTrim(colName), targetFile)
			if err != nil {
				return fmt.Errorf("写入xlsx内存失败:%v", err)
			}
			// 使用filepath.Join来确保跨平台兼容性
			outputPath := filepath.Join(savePath, dataTime+"_"+targetFileName+".xlsx")
			err = targetFile.SaveAs(outputPath)
			if err != nil {
				return fmt.Errorf("写入磁盘失败:%v", err)
			}
		}
		return nil
	}
}

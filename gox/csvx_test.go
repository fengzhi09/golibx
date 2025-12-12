package gox

import (
	"context"
	"os"
	"strings"
	"testing"
)

// 测试toCell函数对不同类型值的处理
func TestToCell(t *testing.T) {
	// 测试字符串类型
	if result := toCell("test"); result != "'test'" {
		t.Errorf("toCell(\"test\") = %v, want %v", result, "'test'")
	}

	// 测试nil值
	if result := toCell(nil); result != "" {
		t.Errorf("toCell(nil) = %v, want %v", result, "")
	}

	// 测试数字类型
	if result := toCell(42); result != "42" {
		t.Errorf("toCell(42) = %v, want %v", result, "42")
	}

	// 测试布尔类型
	if result := toCell(true); result != "true" {
		t.Errorf("toCell(true) = %v, want %v", result, "true")
	}
}

// 测试Csvx的基本功能
func TestCsvxBasic(t *testing.T) {
	ctx := context.Background()
	// 创建临时文件
	tempFile := "test_csv_basic.csv"
	defer os.Remove(tempFile) // 清理临时文件

	// 定义列
	columns := []string{"id", "name", "value"}

	// 测试NewCsv和Reset功能
	csv, err := NewCsv(ctx, tempFile, columns)
	if err != nil {
		t.Fatalf("NewCsv failed: %v", err)
	}

	// 验证文件是否创建并包含列标题
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	contentStr := string(content)
	expectedHeader := "'id','name','value'\n"
	if !strings.HasPrefix(contentStr, expectedHeader) {
		t.Errorf("File header mismatch. Got: %v, Want: %v", contentStr, expectedHeader)
	}

	// 测试Append功能
	row1 := []string{"1", "test1", "100"}
	if err := csv.Append(row1); err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	row2 := []string{"2", "test2", "200"}
	if err := csv.Append(row2); err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// 重新读取文件内容验证
	content, err = os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read temp file after append: %v", err)
	}

	contentStr = string(content)
	// 检查是否包含我们写入的数据
	if !strings.Contains(contentStr, "'1','test1','100'") {
		t.Errorf("File does not contain expected row1 data")
	}
	if !strings.Contains(contentStr, "'2','test2','200'") {
		t.Errorf("File does not contain expected row2 data")
	}

	// 测试Reset功能
	if err := csv.Reset(); err != nil {
		t.Fatalf("Reset failed: %v", err)
	}

	// 验证Reset后只包含列标题
	content, err = os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read temp file after reset: %v", err)
	}

	contentStr = string(content)
	if contentStr != expectedHeader {
		t.Errorf("After Reset, content mismatch. Got: %v, Want: %v", contentStr, expectedHeader)
	}

	// 测试Close功能
	if err := csv.Close(); err != nil {
		t.Errorf("CloseWriters failed: %v", err)
	}
}

// 测试不同数量的列与数据的情况
func TestCsvxDifferentColumnCount(t *testing.T) {
	ctx := context.Background()
	// 创建临时文件
	tempFile := "test_csv_different_cols.csv"
	defer os.Remove(tempFile) // 清理临时文件

	// 定义列
	columns := []string{"col1", "col2"}

	// 创建CSV实例
	csv, err := NewCsv(ctx, tempFile, columns)
	if err != nil {
		t.Fatalf("NewCsv failed: %v", err)
	}

	// 测试数据列数少于定义列数
	lessData := []string{"only_one"}
	if err := csv.Append(lessData); err != nil {
		t.Fatalf("Append with less data failed: %v", err)
	}

	// 测试数据列数多于定义列数
	moreData := []string{"data1", "data2", "extra_data"}
	if err := csv.Append(moreData); err != nil {
		t.Fatalf("Append with more data failed: %v", err)
	}

	// 验证文件内容包含写入的数据
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "'only_one'") {
		t.Errorf("File does not contain less data row")
	}
	if !strings.Contains(contentStr, "'data1','data2','extra_data'") {
		t.Errorf("File does not contain more data row")
	}
}

// 测试特殊字符处理
func TestCsvxSpecialCharacters(t *testing.T) {
	ctx := context.Background()
	// 创建临时文件
	tempFile := "test_csv_special_chars.csv"
	defer os.Remove(tempFile) // 清理临时文件

	// 定义列
	columns := []string{"text"}

	// 创建CSV实例
	csv, err := NewCsv(ctx, tempFile, columns)
	if err != nil {
		t.Fatalf("NewCsv failed: %v", err)
	}

	// 测试包含逗号的文本
	commaText := []string{"text,with,commas"}
	if err := csv.Append(commaText); err != nil {
		t.Fatalf("Append with comma text failed: %v", err)
	}

	// 测试包含引号的文本
	quoteText := []string{"text with \"quotes\""}
	if err := csv.Append(quoteText); err != nil {
		t.Fatalf("Append with quote text failed: %v", err)
	}

	// 验证文件内容
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "'text,with,commas'") {
		t.Errorf("File does not contain comma text")
	}
	if !strings.Contains(contentStr, "'text with \"quotes\"'") {
		t.Errorf("File does not contain quote text")
	}
}

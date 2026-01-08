package gox

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathExists(t *testing.T) {
	// 测试存在的文件
	_, err := os.Create("test_exist.txt")
	assert.NoError(t, err)
	defer os.Remove("test_exist.txt")

	assert.True(t, PathExists("test_exist.txt"))
	assert.False(t, PathExists("non_exist_file.txt"))
}

func TestWorkDir(t *testing.T) {
	expected, err := os.Getwd()
	assert.NoError(t, err)
	expected = strings.Replace(expected, "\\", "/", -1)

	actual := WorkDir()
	assert.Equal(t, expected, actual)
}

func TestCloseRes(t *testing.T) {
	ctx := context.Background()
	file, err := os.Create("test_close.txt")
	assert.NoError(t, err)
	defer os.Remove("test_close.txt")

	// 测试正常关闭
	CloseRes(ctx, "test_file", file)

	// 测试nil资源
	CloseRes(ctx, "nil_resource", nil)
}

func TestSafeClose(t *testing.T) {
	ctx := context.Background()
	// 测试正常关闭函数
	closed := false
	closer := func() error {
		closed = true
		return nil
	}
	SafeClose(ctx, "test_closer", closer)
	assert.True(t, closed)

	// 测试返回错误的关闭函数
	errClosed := false
	errCloser := func() error {
		errClosed = true
		return os.ErrClosed
	}
	SafeClose(ctx, "error_closer", errCloser)
	assert.True(t, errClosed)

	// 测试nil关闭函数
	SafeClose(ctx, "nil_closer", nil)
}

func TestReadAll(t *testing.T) {
	// 创建测试文件
	testContent := "Hello, World!\nThis is a test file."
	_ = os.WriteFile("test_readall.txt", []byte(testContent), 0o644)
	defer os.Remove("test_readall.txt")

	// 测试读取文件
	content, err := ReadAll("test_readall.txt")
	assert.NoError(t, err)
	assert.Equal(t, []byte(testContent), content)

	// 测试读取不存在的文件
	_, err = ReadAll("non_exist_file.txt")
	assert.Error(t, err)
}

func TestReadString(t *testing.T) {
	// 创建测试文件
	testContent := "Hello, String!"
	_ = os.WriteFile("test_readstring.txt", []byte(testContent), 0o644)
	defer os.Remove("test_readstring.txt")

	// 测试读取文件
	content := ReadString("test_readstring.txt")
	assert.Equal(t, testContent, content)

	// 测试读取不存在的文件（应返回空字符串）
	emptyContent := ReadString("non_exist_file.txt")
	assert.Equal(t, "", emptyContent)
}

func TestReadLines(t *testing.T) {
	// 创建测试文件
	testContent := "Line 1\nLine 2\nLine 3"
	_ = os.WriteFile("test_readlines.txt", []byte(testContent), 0o644)
	defer os.Remove("test_readlines.txt")

	// 测试按行读取
	lines := ReadLines("test_readlines.txt")
	expected := []string{"Line 1", "Line 2", "Line 3"}
	assert.Equal(t, expected, lines)

	// 测试读取不存在的文件
	emptyLines := ReadLines("non_exist_file.txt")
	assert.Equal(t, []string{""}, emptyLines)
}

func TestReadJson(t *testing.T) {
	ctx := context.Background()
	// 创建测试JSON文件
	testJson := `{"name":"test","age":30}`
	_ = os.WriteFile("test_readjson.json", []byte(testJson), 0o644)
	defer os.Remove("test_readjson.json")

	// 测试读取JSON
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var result TestStruct
	err := ReadJson(ctx, "test_readjson.json", &result)
	assert.NoError(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 30, result.Age)

	// 测试读取不存在的文件
	var emptyResult TestStruct
	err = ReadJson(ctx, "non_exist_file.json", &emptyResult)
	assert.Error(t, err)

	// 测试读取无效的JSON
	_ = os.WriteFile("test_invalid_json.json", []byte("invalid json"), 0o644)
	defer os.Remove("test_invalid_json.json")
	var invalidResult TestStruct
	err = ReadJson(ctx, "test_invalid_json.json", &invalidResult)
	assert.Error(t, err)
}

func TestParseJson(t *testing.T) {
	// 测试解析有效的JSON
	testJson := `{"name":"parse","value":42}`
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	var result TestStruct
	err := ParseJson(testJson, &result)
	assert.NoError(t, err)
	assert.Equal(t, "parse", result.Name)
	assert.Equal(t, 42, result.Value)

	// 测试解析带有换行符的JSON
	testJsonWithNewlines := `{"name":"parse\r","value":42\n}`
	var resultWithNewlines TestStruct
	err = ParseJson(testJsonWithNewlines, &resultWithNewlines)
	assert.NoError(t, err)
	assert.Equal(t, "parse", resultWithNewlines.Name) // 换行符应该被移除
	assert.Equal(t, 42, resultWithNewlines.Value)

	// 测试解析无效的JSON
	var invalidResult TestStruct
	err = ParseJson("invalid json", &invalidResult)
	assert.Error(t, err)
}

func TestNewFilePath(t *testing.T) {
	// 测试普通文件路径
	fp1 := NewFilePath("test/file.txt")
	assert.Equal(t, "test", fp1.Dir)
	assert.Equal(t, "file", fp1.Name)
	assert.Equal(t, "txt", fp1.Type)
	assert.False(t, fp1.IsDir)

	// 测试带反斜杠的路径
	fp2 := NewFilePath("test\\file.txt")
	assert.Equal(t, "test", fp2.Dir)
	assert.Equal(t, "file", fp2.Name)
	assert.Equal(t, "txt", fp2.Type)
	assert.False(t, fp2.IsDir)

	// 测试没有扩展名的文件
	fp3 := NewFilePath("test/file")
	assert.Equal(t, "test", fp3.Dir)
	assert.Equal(t, "file", fp3.Name)
	assert.Equal(t, "", fp3.Type)
	assert.True(t, fp3.IsDir) // 没有扩展名，应该被认为是目录

	// 测试根目录下的文件
	fp4 := NewFilePath("file.txt")
	assert.Equal(t, ".", fp4.Dir)
	assert.Equal(t, "file", fp4.Name)
	assert.Equal(t, "txt", fp4.Type)
	assert.False(t, fp4.IsDir)

	// 测试多级目录
	fp5 := NewFilePath("test/subdir/file.txt")
	assert.Equal(t, "test/subdir", fp5.Dir)
	assert.Equal(t, "file", fp5.Name)
	assert.Equal(t, "txt", fp5.Type)
	assert.False(t, fp5.IsDir)
}

func TestMkDir(t *testing.T) {
	// 测试创建单级目录
	path1 := MkDir("test_mkdir")
	defer os.Remove("test_mkdir")
	assert.Equal(t, "test_mkdir", path1)
	exists := PathExists("test_mkdir")
	assert.True(t, exists)

	// 测试创建多级目录
	path2 := MkDir("test_mkdir/subdir/child")
	defer os.RemoveAll("test_mkdir")
	assert.Equal(t, "test_mkdir/subdir/child", path2)
	exists = PathExists("test_mkdir/subdir/child")
	assert.True(t, exists)

	// 测试带反斜杠的路径
	path3 := MkDir("test_mkdir\\backslash")
	assert.Equal(t, "test_mkdir/backslash", path3) // 应该将反斜杠转换为正斜杠
	exists = PathExists("test_mkdir/backslash")
	assert.True(t, exists)
}

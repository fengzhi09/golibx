package gox

import (
	"context"
	"io"
	"os"
	"reflect"
	"strings"
)

func MkDir(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	_ = os.MkdirAll(path, 0o750)
	return path
}

func ClearDir(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	_ = os.RemoveAll(path)
	_ = os.MkdirAll(path, 0o750)
	return path
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func WorkDir() string {
	pwd, _ := os.Getwd()
	return strings.ReplaceAll(pwd, "\\", "/")
}

func CloseRes(ctx context.Context, name string, res io.Closer) {
	if res != nil && !reflect.ValueOf(res).IsNil() {
		SafeClose(ctx, name, res.Close)
	}
}

func SafeClose(ctx context.Context, name string, closer func() error) error {
	if closer != nil {
		err := closer()
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadAll(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func ReadString(path string) string {
	data, _ := ReadAll(path)
	return string(data)
}

func ReadLines(path string) []string {
	data := ReadString(path)
	return strings.Split(data, "\n")
}

func ReadJson(ctx context.Context, path string, ret any) error {
	data, err := ReadAll(path)
	if err != nil {
		return err
	}
	return Unmarshal(data, ret)
}

func ParseJson(content string, ret any) error {
	content = strings.ReplaceAll(strings.ReplaceAll(content, "\\r", ""), "\\n", "")
	return Unmarshal([]byte(content), ret)
}

type FilePath struct {
	Dir   string
	Name  string
	Type  string
	IsDir bool
}

func (fp *FilePath) Path() string {
	return fp.Dir + "/" + fp.Name + "." + fp.Type
}

func (fp *FilePath) Exists() bool {
	return PathExists(fp.Path())
}

func NewFilePath(input string) *FilePath {
	input = strings.ReplaceAll(strings.ReplaceAll(input, "//", "/"), "\\", "/")

	// 处理目录分隔符
	pSep := strings.LastIndex(input, "/")
	var dir, file string
	if pSep >= 0 {
		dir, file = MkDir(input[:pSep]), input[pSep+1:]
	} else {
		dir, file = ".", input
	}

	// 处理文件扩展名
	fSep := strings.LastIndex(file, ".")
	var name, typ string
	if fSep >= 0 {
		name, typ = file[:fSep], file[fSep+1:]
	} else {
		name, typ = file, ""
	}

	return &FilePath{dir, name, typ, typ == ""}
}

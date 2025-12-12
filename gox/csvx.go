package gox

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// 优化toCell函数，参考Python版本
func toCell(val interface{}) string {
	switch v := val.(type) {
	case map[string]interface{}:
		return "'" + UnsafeMarshalString(v) + "'"
	case string:
		return "'" + v + "'"
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

type Csvx struct {
	file    string
	columns []string
}

func NewCsv(ctx context.Context, file string, columns []string) (*Csvx, error) {
	c := &Csvx{file: file, columns: columns}
	return c, c.Reset()
}

func (c *Csvx) Append(values []string) error {
	f, err := os.OpenFile(c.file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	// 格式化值
	formatted := make([]string, len(values))
	for i, v := range values {
		formatted[i] = toCell(v)
	}

	// 写入行数据
	line := strings.Join(formatted, ",") + "\n"
	_, err = f.WriteString(line)
	return err
}

func (c *Csvx) Reset() error {
	f, err := os.Create(c.file)
	if err != nil {
		return err
	}
	defer f.Close()
	// 写入列标题
	_, err = f.WriteString("'" + strings.Join(c.columns, "','") + "'\n")
	return err
}

func (c *Csvx) Close() error {
	return nil
}

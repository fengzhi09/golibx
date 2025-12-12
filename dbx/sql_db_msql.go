package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/fengzhi09/golibx/jsonx"

	_ "github.com/go-sql-driver/mysql"
)

// MSql 数据库基类（用于MySQL和Doris）
type MSql struct {
	name   string
	conf   DBConf
	db     *sql.DB
	dbType string
}

// Connect 建立数据库连接
func (m *MSql) Connect(ctx context.Context) error {
	// 具体实现由子类提供
	// 构建连接字符串
	connStr := m.ConnString(m.conf)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %v", err)
	}

	m.db = db
	return nil
}

// ConnString 构建MySQL连接字符串
func (m *MSql) ConnString(conf DBConf) string {
	url := getInOrder(conf, "db_url", "url", "my_url")
	if url == "" {
		// 从单独的配置项构建
		host := getOrDefault(conf, "host", "localhost")
		port := getOrDefault(conf, "port", "3306")
		dbname := getOrDefault(conf, "dbname", "mysql")
		user := getOrDefault(conf, "user", "root")
		password := getOrDefault(conf, "password", "")
		url = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, password, host, port, dbname)
	}

	return strings.ReplaceAll(url, m.dbType+"://", "")
}

// Query 执行MySQL查询
func (m *MSql) Query(ctx context.Context, query string, args ...any) ([]*jsonx.JObj, error) {
	if m.db == nil {
		if err := m.Connect(ctx); err != nil {
			return nil, fmt.Errorf("failed to connect: %v", err)
		}
	}
	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%v query error: %v", m.dbType, err)
	}
	defer rows.Close()
	results := make([]*jsonx.JObj, 0)
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("%v get columns error: %v", m.dbType, err)
	}
	values := make([]any, len(columns))
	scanArgs := make([]any, len(columns))
	for i := range columns {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("%v scan row error: %v", m.dbType, err)
		}
		row := make(jsonx.JObj)
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, &row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v rows error: %v", m.dbType, err)
	}

	return results, nil
}

// Close 关闭数据库连接（用于MySQL和Doris）
func (m *MSql) Close(ctx context.Context) error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

type Doris = MSql

package dbx

import (
	"context"
	"fmt"

	"github.com/fengzhi09/golibx/jsonx"

	"github.com/jackc/pgx/v5"
)

// PSql PostgreSQL数据库基类（使用pgx）
type PSql struct {
	name   string
	conf   DBConf
	db     *pgx.Conn
	dbType string
}

// ConnString 构建PostgreSQL连接字符串
func (p *PSql) ConnString(conf DBConf) string {
	url := getInOrder(conf, "db_url", "url", "pg_url", "driver")
	if url != "" {
		return url
	}

	// 从单独的配置项构建
	host := getOrDefault(conf, "host", "localhost")
	port := getOrDefault(conf, "port", "5432")
	dbname := getOrDefault(conf, "dbname", "postgres")
	user := getOrDefault(conf, "user", "postgres")
	password := getOrDefault(conf, "password", "")

	return fmt.Sprintf("postgres://%s:%s@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)
}

// Connect 建立PostgreSQL基础连接
func (p *PSql) Connect(ctx context.Context) error {
	// 由子类实现具体连接逻辑
	// 构建连接字符串
	connStr := p.ConnString(p.conf)
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// 测试连接
	if err := db.Ping(ctx); err != nil {
		db.Close(ctx)
		return fmt.Errorf("failed to ping database: %v", err)
	}

	p.db = db
	return nil
}

// Close 关闭PostgreSQL连接
func (p *PSql) Close(ctx context.Context) error {
	if p.db != nil {
		return p.db.Close(ctx)
	}
	return nil
}

// Query 执行PostgreSQL查询
func (p *PSql) Query(ctx context.Context, query string, args ...any) ([]*jsonx.JObj, error) {
	if p.db == nil {
		if err := p.Connect(ctx); err != nil {
			return nil, fmt.Errorf("failed to connect: %v", err)
		}
	}
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%v query error: %v", p.dbType, err)
	}
	defer rows.Close()
	results := make([]*jsonx.JObj, 0)
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = fd.Name
	}
	values := make([]any, len(columns))
	scanArgs := make([]any, len(columns))
	for i := range columns {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("%v scan row error: %v", p.dbType, err)
		}
		row := make(jsonx.JObj)
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, &row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v rows error: %v", p.dbType, err)
	}

	return results, nil
}

package dbx_vec

import (
	"context"
	"fmt"
	"github.com/fengzhi09/golibx/jsonx"

	pgx "github.com/jackc/pgx/v5"
)

// PgVecDB PostgreSQL向量数据库实现
type PgVecDB struct {
	VecDBBase
	db *pgx.Conn
}

// Close implements VecDB.
func (p *PgVecDB) Close(ctx context.Context) error {
	panic("unimplemented")
}

func (p *PgVecDB) NewTable(ctx context.Context, name string, conf jsonx.JObj) error {
	//TODO implement me
	panic("implement me")
}

func (p *PgVecDB) Search(ctx context.Context, query VecQuery) ([]*ResNode, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PgVecDB) Upsert(ctx context.Context, nodes ...VecNode) error {
	//TODO implement me
	panic("implement me")
}

func (p *PgVecDB) UpsertM(ctx context.Context, nodes ...MVecNode) error {
	//TODO implement me
	panic("implement me")
}

// NewPgVecDB 创建PostgreSQL向量数据库实例
func NewPgVecDB(ctx context.Context, resName string, resConf map[string]any) (VecDB, error) {
	if resConf == nil {
		resConf = jsonx.NewObj(map[string]any{})
	}
	return &PgVecDB{
		VecDBBase: VecDBBase{
			resName: resName,
			resConf: resConf,
		},
	}, nil
}

// 辅助方法
func (p *PgVecDB) buildConnString() string {
	conf := p.resConf
	// 从配置中获取连接参数
	host := "localhost"
	if val, ok := conf["host"].(string); ok {
		host = val
	}

	port := "5432"
	if val, ok := conf["port"].(string); ok {
		port = val
	} else if val, ok := conf["port"].(float64); ok {
		port = fmt.Sprintf("%d", int(val))
	}

	dbname := "postgres"
	if val, ok := conf["dbname"].(string); ok {
		dbname = val
	}

	user := "postgres"
	if val, ok := conf["user"].(string); ok {
		user = val
	}

	password := ""
	if val, ok := conf["password"].(string); ok {
		password = val
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

// Connect 建立数据库连接
func (p *PgVecDB) Connect(ctx context.Context) error {
	// 构建连接字符串
	connStr := p.buildConnString()

	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	// 测试连接
	if err := db.Ping(ctx); err != nil {
		db.Close(ctx)
		return fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}

	p.db = db
	return nil
}

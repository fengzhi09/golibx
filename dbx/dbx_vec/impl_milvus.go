package dbx_vec

import (
	"context"
	"fmt"
	"github.com/fengzhi09/golibx/jsonx"

	milvus "github.com/milvus-io/milvus-sdk-go/v2/client"
)

// MilvusDB Milvus向量数据库实现
type MilvusDB struct {
	VecDBBase
	client any // 实际应该使用Milvus Go SDK客户端
}

// Close implements VecDB.
func (m *MilvusDB) Close(ctx context.Context) error {
	panic("unimplemented")
}

// NewMilvusDB 创建Milvus数据库实例
func NewMilvusDB(ctx context.Context, resName string, resConf VecDBConf) (VecDB, error) {
	if resConf == nil {
		resConf = jsonx.NewObj(map[string]any{})
	}
	return &MilvusDB{
		VecDBBase: VecDBBase{
			resName: resName,
			resConf: resConf,
		},
	}, nil
}

// Connect 建立数据库连接
func (m *MilvusDB) Connect(ctx context.Context) error {
	client, err := milvus.NewClient(ctx, milvus.Config{
		Address:    m.resConf.GetStr("address"),
		APIKey:     m.resConf.GetStr("api_key"),
		Identifier: m.resConf.GetStr("app_name"),
	})
	if err != nil {
		return fmt.Errorf("failed to create milvus client: %v", err)
	}
	m.client = client
	return nil
}

func (m *MilvusDB) NewTable(ctx context.Context, name string, conf jsonx.JObj) error {
	//TODO implement me
	panic("implement me")
}

func (m *MilvusDB) Search(ctx context.Context, query VecQuery) ([]*ResNode, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MilvusDB) Upsert(ctx context.Context, nodes ...VecNode) error {
	//TODO implement me
	panic("implement me")
}

func (m *MilvusDB) UpsertM(ctx context.Context, nodes ...MVecNode) error {
	//TODO implement me
	panic("implement me")
}

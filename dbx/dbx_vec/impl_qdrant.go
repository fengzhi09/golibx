package dbx_vec

import (
	"context"
	"fmt"

	"github.com/fengzhi09/golibx/gox"
	"github.com/fengzhi09/golibx/jsonx"

	"github.com/qdrant/go-client/qdrant"
)

// QdrantDB Qdrant向量数据库实现
type QdrantDB struct {
	VecDBBase
	client *qdrant.Client
}

// NewQdrantDB 创建Qdrant数据库实例
func NewQdrantDB(resName string, resConf VecDBConf) (VecDB, error) {
	if resConf == nil {
		resConf = jsonx.NewObj(map[string]any{})
	}
	if _, ok := resConf["use_tls"]; ok {
		resConf["use_tls"] = true
	}
	if _, ok := resConf["pool_size"]; ok {
		resConf["pool_size"] = 10
	}
	return &QdrantDB{
		VecDBBase: VecDBBase{
			resName: resName,
			resConf: resConf,
		},
		client: nil,
	}, nil
}

// Connect 建立Qdrant数据库连接（使用官方Go客户端）
func (q *QdrantDB) Connect(ctx context.Context) (err error) {
	config := &qdrant.Config{}
	if q.resConf.Contains("api_key") {
		config.APIKey = q.resConf.GetStr("api_key")
	} else {
		config.Host = q.resConf.GetStr("host")
		config.Port = q.resConf.GetInt("port")
	}
	config.UseTLS = q.resConf.GetBool("use_tls")
	config.PoolSize = uint(q.resConf.GetInt("pool_size"))
	q.client, err = qdrant.NewClient(config)
	return err
}

func (q *QdrantDB) Close(ctx context.Context) error {
	if q.client != nil {
		err := q.client.Close()
		q.client = nil
		return err
	}
	return nil
}

func (q *QdrantDB) ensure(ctx context.Context) {
	if q.client == nil {
		if err := q.Connect(ctx); err != nil {
			panic(err)
		}
	}
}

// Search 执行向量搜索（使用官方Go客户端）
func (q *QdrantDB) Search(ctx context.Context, query VecQuery) ([]*ResNode, error) {
	q.ensure(ctx)
	// 提取查询向量
	if len(query.FiltersVec) == 0 || len(query.FiltersVec[0].Val) == 0 {
		return nil, fmt.Errorf("no query vector provided")
	}
	musts, skips := []*qdrant.Condition{}, []*qdrant.Condition{}
	for _, filter := range query.Filters {
		if gox.In(filter.Op, "=", "==", "eq") {
			musts = append(musts, qdrant.NewMatch(filter.Field, gox.AsStr(filter.Val)))
		} else if gox.In(filter.Op, "!=", "ne", "skip") {
			skips = append(skips, qdrant.NewMatch(filter.Field, gox.AsStr(filter.Val)))
		}
	}
	res, err := q.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: query.Table, Filter: &qdrant.Filter{Must: musts, MustNot: skips},
	})
	nodes := make([]*ResNode, 0, len(res))
	for _, point := range res {
		nodes = append(nodes, &ResNode{
			Id:       point.Id.String(),
			MetaData: AsJson(point.Payload),
			Vec:      point.Vectors.GetVector().Data,
			Score:    point.Score,
			Rank:     int(point.OrderValue.GetInt()),
		})
	}

	return nodes, err
}

func AsJson(payload map[string]*qdrant.Value) jsonx.JObj {
	meta := jsonx.JObj{}
	for k, v := range payload {
		if val := QdAsJson(v); val != nil {
			meta.Put(k, val)
		}
	}
	return meta
}

func QdAsJson(v *qdrant.Value) jsonx.JValue {
	if v == nil {
		return nil
	}
	if _, ok := v.GetKind().(*qdrant.Value_NullValue); ok {
		return nil
	} else if x, ok := v.GetKind().(*qdrant.Value_BoolValue); ok {
		return jsonx.NewJBool(x.BoolValue)
	} else if x, ok := v.GetKind().(*qdrant.Value_IntegerValue); ok {
		return jsonx.NewJInt(x.IntegerValue)
	} else if x, ok := v.GetKind().(*qdrant.Value_DoubleValue); ok {
		return jsonx.NewJNum(x.DoubleValue)
	} else if x, ok := v.GetKind().(*qdrant.Value_StringValue); ok {
		return jsonx.NewJStr(x.StringValue)
	} else if x, ok := v.GetKind().(*qdrant.Value_ListValue); ok {
		vl := jsonx.JArr{}
		for _, v := range x.ListValue.Values {
			vl = append(vl, QdAsJson(v))
		}
		return vl
	} else if x, ok := v.GetKind().(*qdrant.Value_StructValue); ok {
		vm := jsonx.JObj{}
		for k, v2 := range x.StructValue.Fields {
			vm[k] = QdAsJson(v2)
		}
		return vm
	}
	return jsonx.NewJStr(v.String())
}

func (q *QdrantDB) NewTable(ctx context.Context, name string, conf jsonx.JObj) error {
	q.ensure(ctx)
	return q.client.CreateCollection(ctx,
		&qdrant.CreateCollection{
			CollectionName: name,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     uint64(conf.GetLong("size")),
				Distance: qdrant.Distance_Cosine,
			}),
		})
}

func (q *QdrantDB) Upsert(ctx context.Context, nodes ...VecNode) error {
	q.ensure(ctx)
	pts := make([]*qdrant.PointStruct, 0, len(nodes))
	for _, node := range nodes {
		pts = append(pts, &qdrant.PointStruct{
			Id:      qdrant.NewID(node.Id),
			Vectors: qdrant.NewVectorsDense(node.Vec),
			Payload: qdrant.NewValueMap(node.MetaData),
		})
	}
	opInfo, err := q.client.Upsert(ctx,
		&qdrant.UpsertPoints{
			CollectionName: nodes[0].Id,
			Points:         pts,
		})
	if opInfo.Status == qdrant.UpdateStatus_ClockRejected && err == nil {
		err = fmt.Errorf("update rejected due to an outdated clock")
	}
	return err
}

func (q *QdrantDB) UpsertM(ctx context.Context, nodes ...MVecNode) error {
	// TODO implement me
	panic("implement me")
}

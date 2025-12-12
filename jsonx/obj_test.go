package jsonx

import (
	"reflect"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
)

func TestParseObjAndMapEq(t *testing.T) {
	r1 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
	r2 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629741c8f5310f5558527f41","629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{"send_from":["manu_send"]},"seek_type":"flag.send_from.msg"}]}}}}`
	var m1, m2 map[string]interface{}

	err := Unmarshal([]byte(r1), &m1)
	assert.NoError(t, err)

	err = Unmarshal([]byte(r2), &m2)
	assert.NoError(t, err)

	assert.False(t, MapEq(m1, m2))

	r2 = `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
	err = sonic.Unmarshal([]byte(r2), &m2)
	assert.NoError(t, err)
	assert.True(t, MapEq(m1, m2))
}

func TestReflectDeepEq(t *testing.T) {
	r1 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
	r2 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629741c8f5310f5558527f41","629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{"send_from":["manu_send"]},"seek_type":"flag.send_from.msg"}]}}}}`
	var m1, m2 map[string]any

	err := Unmarshal([]byte(r1), &m1)
	assert.NoError(t, err)

	err = Unmarshal([]byte(r2), &m2)
	assert.NoError(t, err)

	assert.False(t, reflect.DeepEqual(m1, m2))

	r2 = `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
	err = sonic.Unmarshal([]byte(r2), &m2)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(m1, m2))
}

func BenchmarkReflectDeepEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r1 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
		r2 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629741c8f5310f5558527f41","629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{"send_from":["manu_send"]},"seek_type":"flag.send_from.msg"}]}}}}`
		var m1, m2 map[string]any

		err := Unmarshal([]byte(r1), &m1)
		assert.NoError(b, err)

		err = Unmarshal([]byte(r2), &m2)
		assert.NoError(b, err)

		assert.False(b, reflect.DeepEqual(m1, m2))

		r2 = `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
		err = Unmarshal([]byte(r2), &m2)
		assert.NoError(b, err)
		assert.True(b, reflect.DeepEqual(m1, m2))
	}
}

func BenchmarkMapEq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r1 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
		r2 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629741c8f5310f5558527f41","629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{"send_from":["manu_send"]},"seek_type":"flag.send_from.msg"}]}}}}`
		var m1, m2 map[string]any

		err := sonic.Unmarshal([]byte(r1), &m1)
		assert.NoError(b, err)

		err = sonic.Unmarshal([]byte(r2), &m2)
		assert.NoError(b, err)

		assert.False(b, MapEq(m1, m2))

		r2 = `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
		err = sonic.Unmarshal([]byte(r2), &m2)
		assert.NoError(b, err)
		assert.True(b, MapEq(m1, m2))
	}
}

// go test -bench=. -run=none -benchtime=3s ./libs/jsonx/
// goos: darwin
// goarch: amd64
// pkg: github.com/fengzhi09/golibx/jsonx
// cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
// BenchmarkReflectDeepEqual-16               71871             49481 ns/op
// BenchmarkMapEq-16                          26368            138189 ns/op
// PASS
// ok      github.com/fengzhi09/golibx/jsonx 10.415s

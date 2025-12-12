package gox

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXEq(t *testing.T) {
	r1 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
	r2 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629741c8f5310f5558527f41","629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{"send_from":["manu_send"]},"seek_type":"flag.send_from.msg"}]}}}}`
	r3 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
	m1, m2, m3 := map[string]any{}, map[string]any{}, map[string]any{}
	_ = Unmarshal([]byte(r1), &m1)
	_ = Unmarshal([]byte(r2), &m2)
	_ = Unmarshal([]byte(r3), &m3)
	type testCase_ struct {
		name string
		eq   func(a, b any, tx *testing.T) bool
		e12  bool
		e13  bool
	}
	tests := []*testCase_{
		{"json-my", func(a, b any, tx *testing.T) bool {
			err := jsonDiff(a, b)
			if err != nil {
				tx.Logf("got err : %v", err)
			}
			return err == nil
		}, false, true},
		{"reflect-my", func(a, b any, tx *testing.T) bool {
			err := reflectDiff(reflect.ValueOf(a), reflect.ValueOf(b))
			if err != nil {
				tx.Logf("got err : %v", err)
			}
			return err == nil
		}, false, true},
		{"reflect-deep", func(a, b any, tx *testing.T) bool {
			ok := reflect.DeepEqual(a, b)
			if !ok {
				tx.Log("got err : not eq")
			}
			return ok
		}, false, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Logf("%v starts", "eq(m1, m2)")
			assert.Equal(tt, test.e12, test.eq(m1, m2, tt), "eq(m1, m2)")
			tt.Logf("%v starts", "eq(m1, m3)")
			assert.Equal(tt, test.e13, test.eq(m1, m3, tt), "eq(m1, m3)")
		})
	}
}

func BenchmarkXEq(t *testing.B) {
	r1 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
	r2 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629741c8f5310f5558527f41","629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{"send_from":["manu_send"]},"seek_type":"flag.send_from.msg"}]}}}}`
	r3 := `{"company_id":"612c53cb7250e1e5140faded","create_time":"2022-06-01T02:14:58.72Z","description":"","id":"6296cba2cff3ca464345a28b","name":"无无无2位","platform":"","qc_norm_ids":["629611ac43f804bd918a3301","62415ac5e191a6bfada59c0d","629741c8f5310f5558527f41"],"rule_category":3,"rule_type":0,"settings":{"cal_op":1,"check_sale_status":[],"check_step":-1,"description":"","name":"无无无2位","score":1,"xact":{"tag_mode":"trigger"},"xrule":{"filter":{"sub_conds":[]},"trigger":{"logic_op":"and","sub_conds":[{"seek_conf":{},"seek_type":"flag.send_from.msg"}]}}}}`
	m1, m2, m3 := map[string]any{}, map[string]any{}, map[string]any{}
	_ = Unmarshal([]byte(r1), &m1)
	_ = Unmarshal([]byte(r2), &m2)
	_ = Unmarshal([]byte(r3), &m3)
	type testCase_ struct {
		name string
		eq   func(a, b any, tx *testing.B) bool
		e12  bool
		e13  bool
	}
	tests := []*testCase_{
		{"json-my", func(a, b any, tx *testing.B) bool {
			err := jsonDiff(a, b)
			return err == nil
		}, false, true},
		{"reflect-my", func(a, b any, tx *testing.B) bool {
			err := reflectDiff(reflect.ValueOf(a), reflect.ValueOf(b))
			return err == nil
		}, false, true},
		{"reflect-deep", func(a, b any, tx *testing.B) bool {
			ok := reflect.DeepEqual(a, b)
			return ok
		}, false, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.B) {
			for i := 0; i < tt.N; i++ {
				assert.Equal(tt, test.e12, test.eq(m1, m2, tt), "eq(m1, m2)")
				assert.Equal(tt, test.e13, test.eq(m1, m3, tt), "eq(m1, m3)")
			}
		})
	}
}

/*
-test.v -test.paniconexit0 -test.bench . -test.run ^$ -test.benchmem -test.benchtime 1s
goos: windows
goarch: amd64
pkg: github.com/fengzhi09/golibx/gox
cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
BenchmarkXEq
BenchmarkXEq/json-my
BenchmarkXEq/json-my-8         	   83456	     13455 ns/op	    4581 B/op	     134 allocs/op
BenchmarkXEq/reflect-my
BenchmarkXEq/reflect-my-8      	   15098	     79761 ns/op	   36294 B/op	     937 allocs/op
BenchmarkXEq/reflect-deep
BenchmarkXEq/reflect-deep-8    	   58610	     20520 ns/op	    8644 B/op	     127 allocs/op
PASS
*/

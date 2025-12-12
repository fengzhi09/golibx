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

func TestObj_CreationMethods(t *testing.T) {
	// 测试NewObj创建对象
	mapData := map[string]any{"key1": "value1", "key2": 123, "key3": 3.14, "key4": true}
	obj := NewObj(mapData)
	assert.Equal(t, 4, obj.Size())
	assert.Equal(t, "value1", obj.GetStr("key1"))
	assert.Equal(t, 123, obj.GetInt("key2"))
	assert.Equal(t, 3.14, obj.GetDouble("key3"))
	assert.Equal(t, true, obj.GetBool("key4"))

	// 测试NewObjPtr创建对象指针
	objPtr := NewObjPtr(mapData)
	assert.NotNil(t, objPtr)
	assert.Equal(t, 4, objPtr.Size())
	assert.Equal(t, "value1", objPtr.GetStr("key1"))

	// 测试ParseObj和ParseObjPtr解析JSON
	jsonStr := `{"name": "test", "age": 25}`
	parsedObj := ParseObj([]byte(jsonStr))
	assert.Equal(t, 2, parsedObj.Size())
	assert.Equal(t, "test", parsedObj.GetStr("name"))
	assert.Equal(t, 25, parsedObj.GetInt("age"))

	parsedObjPtr := ParseObjPtr([]byte(jsonStr))
	assert.NotNil(t, parsedObjPtr)
	assert.Equal(t, 2, parsedObjPtr.Size())
	assert.Equal(t, "test", parsedObjPtr.GetStr("name"))
}

func TestObj_ToMapKeysValues(t *testing.T) {
	// 测试ToMap方法
	obj := NewObj(map[string]any{"key1": "value1", "key2": 123})
	m := obj.ToMap()
	// ToMap返回的是map[string]any，其中的值可能是JStr、JInt等类型
	// 我们只需要验证键的存在性和值的类型
	assert.Contains(t, m, "key1")
	assert.Contains(t, m, "key2")
	assert.IsType(t, JStr(""), m["key1"])
	assert.IsType(t, JInt(0), m["key2"])

	// 测试Keys方法
	keys := obj.Keys()
	assert.Equal(t, 2, len(keys))
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")

	// 测试Values方法
	values := obj.Values()
	assert.Equal(t, 2, len(values))
	for _, val := range values {
		assert.NotNil(t, val)
	}
}

func TestObj_SizeIsEmptyContains(t *testing.T) {
	// 测试Size和IsEmpty方法
	emptyObj := NewObj(map[string]any{})
	assert.Equal(t, 0, emptyObj.Size())
	assert.True(t, emptyObj.IsEmpty())

	nonEmptyObj := NewObj(map[string]any{"key": "value"})
	assert.Equal(t, 1, nonEmptyObj.Size())
	assert.False(t, nonEmptyObj.IsEmpty())

	// 测试Contains方法
	assert.True(t, nonEmptyObj.Contains("key"))
	assert.False(t, nonEmptyObj.Contains("non_existent_key"))
}

func TestObj_ForeachMethods(t *testing.T) {
	// 测试Foreach方法
	obj := NewObj(map[string]any{"key1": "value1", "key2": 123, "key3": 3.14})
	count := 0
	keysSeen := make(map[string]bool)

	obj.Foreach(func(key string, val JValue) bool {
		count++
		keysSeen[key] = true
		return true // 继续遍历
	})

	assert.Equal(t, 3, count)
	assert.True(t, keysSeen["key1"])
	assert.True(t, keysSeen["key2"])
	assert.True(t, keysSeen["key3"])

	// 测试Foreach中途停止
	count = 0
	obj.Foreach(func(key string, val JValue) bool {
		count++
		return count < 2 // 只遍历前2个元素
	})
	assert.Equal(t, 2, count)
}

func TestObj_MergeMethods(t *testing.T) {
	// 测试Merge方法
	obj1 := NewObj(map[string]any{"key1": "value1", "key2": 123})
	obj2 := NewObj(map[string]any{"key2": 456, "key3": "value3"})

	obj1.Merge(obj2)
	assert.Equal(t, 3, obj1.Size())
	assert.Equal(t, "value1", obj1.GetStr("key1"))
	assert.Equal(t, 456, obj1.GetInt("key2"))      // 应该被覆盖
	assert.Equal(t, "value3", obj1.GetStr("key3")) // 应该被添加

	// 测试MergePtr方法
	obj3 := NewObj(map[string]any{"key1": "value1"})
	obj4 := NewObj(map[string]any{"key2": "value2"})
	obj3.MergePtr(&obj4)
	assert.Equal(t, 2, obj3.Size())
	assert.Equal(t, "value1", obj3.GetStr("key1"))
	assert.Equal(t, "value2", obj3.GetStr("key2"))
}

func TestObj_GetMethods(t *testing.T) {
	obj := NewObj(map[string]any{"string_key": "value", "int_key": 123, "double_key": 3.14, "bool_key": true})

	// 测试GetVal和GetVal2方法
	val, exists := obj.GetVal2("string_key")
	assert.True(t, exists)
	assert.Equal(t, "value", val.String())

	val = obj.GetVal("non_existent_key")
	assert.True(t, val.IsNull())

	// 测试GetValIgnore方法（忽略大小写）
	obj.Put("CASE_SENSITIVE_KEY", "test_value")
	val = obj.GetValIgnore("case_sensitive_key")
	assert.Equal(t, "test_value", val.String())

	val = obj.GetValIgnore("CASE_SENSITIVE_KEY")
	assert.Equal(t, "test_value", val.String())

	val = obj.GetValIgnore("CASE_SENSITIVE_KEY")
	assert.Equal(t, "test_value", val.String())

	// 测试各种Get*方法
	assert.Equal(t, "value", obj.GetStr("string_key"))
	assert.Equal(t, 123, obj.GetInt("int_key"))
	assert.Equal(t, int64(123), obj.GetLong("int_key"))
	assert.Equal(t, 3.14, obj.GetDouble("double_key"))
	assert.Equal(t, float32(3.14), obj.GetFloat("double_key"))
	assert.Equal(t, true, obj.GetBool("bool_key"))

	// 测试GetOrDefault方法
	assert.Equal(t, "value", obj.GetOrDefault("string_key", "default"))
	assert.Equal(t, "default", obj.GetOrDefault("non_existent", "default"))

	// 测试GetOr方法
	val = obj.GetOr("string_key", "default")
	assert.Equal(t, "value", val.String())

	val = obj.GetOr("non_existent", "default")
	assert.Equal(t, "default", val.String())
}

func TestObj_PutMethods(t *testing.T) {
	obj := NewObj(map[string]any{})

	// 测试各种Put*方法
	obj.Put("put_key", "put_value")
	assert.Equal(t, "put_value", obj.GetStr("put_key"))

	obj.PutStr("str_key", "test_string")
	assert.Equal(t, "test_string", obj.GetStr("str_key"))

	obj.PutInt("int_key", 123)
	assert.Equal(t, 123, obj.GetInt("int_key"))

	obj.PutLong("long_key", 456789)
	assert.Equal(t, int64(456789), obj.GetLong("long_key"))

	obj.PutDouble("double_key", 3.14159)
	assert.Equal(t, 3.14159, obj.GetDouble("double_key"))

	obj.PutBool("bool_key", true)
	assert.Equal(t, true, obj.GetBool("bool_key"))
}

func TestObj_SnapClone(t *testing.T) {
	obj := NewObj(map[string]any{"key1": "value1", "key2": 123, "key3": 3.14})

	// 测试Snap方法
	snapped := obj.Snap("key1", "key3")
	assert.Equal(t, 2, snapped.Size())
	assert.Equal(t, "value1", snapped.GetStr("key1"))
	assert.Equal(t, 3.14, snapped.GetDouble("key3"))
	assert.False(t, snapped.Contains("key2"))

	// 测试Clone方法
	cloned := obj.Clone()
	assert.Equal(t, obj.Size(), cloned.Size())
	assert.Equal(t, obj.GetStr("key1"), cloned.GetStr("key1"))
	assert.Equal(t, obj.GetInt("key2"), cloned.GetInt("key2"))
	assert.Equal(t, obj.GetDouble("key3"), cloned.GetDouble("key3"))
}

func TestObj_ToLineToCSVToCells(t *testing.T) {
	obj := NewObj(map[string]any{"name": "test", "age": 25, "city": "beijing"})

	// 测试ToCells方法
	cells := obj.ToCells([]string{"name", "age", "city"})
	assert.Equal(t, []string{"test", "25", "beijing"}, cells)

	// 测试ToLine方法
	line := obj.ToLine(",", []string{"name", "age", "city"})
	assert.Equal(t, "name=test,age=25,city=beijing", line)

	// 测试ToCSV方法
	csv := obj.ToCSV([]string{"name", "age", "city"})
	assert.Equal(t, "test,25,beijing", csv)
}

func TestObj_ToObjToObjPtrToArr(t *testing.T) {
	obj := NewObj(map[string]any{"key": "value"})

	// 测试ToObj和ToObjPtr方法
	sameObj := obj.ToObj()
	assert.Equal(t, obj, sameObj)

	objPtr := obj.ToObjPtr()
	assert.NotNil(t, objPtr)
	assert.Equal(t, obj, *objPtr)

	// 测试ToArr方法（应该返回空数组）
	emptyArr := obj.ToArr()
	assert.Equal(t, JArr{}, emptyArr)
	assert.Equal(t, 0, emptyArr.Size())
}

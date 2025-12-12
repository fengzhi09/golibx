package jsonx

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fengzhi09/golibx/gox"
)

type JObj map[string]any

func (re JObj) Type() JType       { return JOBJ }
func (re JObj) IsNull() bool      { return false }
func (re JObj) ToInt() int        { return 0 }
func (re JObj) ToTime() time.Time { return AsTime(nil) }
func (re JObj) ToLong() int64     { return 0 }
func (re JObj) ToDouble() float64 { return 0 }
func (re JObj) ToFloat() float32  { return 0 }
func (re JObj) ToBool() bool      { return re.Size() > 0 }
func (re JObj) String() string    { return re.ToJDoc().String() }
func (re JObj) Pretty() string    { return re.ToJDoc().String() }
func (re JObj) ToJDoc() JDoc      { return DumpJObj(re) }
func (re JObj) ToObj() JObj       { return re }
func (re JObj) ToObjPtr() *JObj   { return &re }
func (re JObj) ToArr() JArr       { return JArr{} }
func (re JObj) ToJVal() JValue    { return re }
func (re JObj) ToGVal() any       { return re }

func (re JObj) ToMap() map[string]any { return re }
func (re JObj) Keys() []string {
	var keys []string
	for key := range re {
		keys = append(keys, key)
	}
	return keys
}

func (re JObj) Values() []JValue {
	var values []JValue
	for _, v := range re {
		values = append(values, v.(JValue))
	}
	return values
}
func (re JObj) Size() int                { return len(re) }
func (re JObj) IsEmpty() bool            { return re.Size() == 0 }
func (re JObj) Contains(key string) bool { _, h := re.GetVal2(key); return h }
func (re JObj) Foreach(act func(string, JValue) bool) {
	if re.IsEmpty() {
		return
	}
	for key := range re {
		if !act(key, re.GetVal(key)) {
			break
		}
	}
}

func (re JObj) ForeachObj(act func(string, *JObj) bool) {
	if re.IsEmpty() {
		return
	}
	for key := range re {
		if !act(key, re.GetObjPtr(key)) {
			break
		}
	}
}

func (re JObj) Merge(oth JObj) {
	if oth.IsEmpty() {
		return
	}
	oth.Foreach(func(key string, val JValue) bool {
		re.PutVal(key, val)
		return true
	})
}

func (re JObj) MergePtr(oth *JObj) {
	if oth.IsEmpty() {
		return
	}
	oth.Foreach(func(key string, val JValue) bool {
		re.PutVal(key, val)
		return true
	})
}

func (re JObj) GetVal2(key string) (JValue, bool) {
	v, h := (re)[key]
	if !h || v == nil {
		return NewJNull(), false
	}
	return GoV2JV(v), true
}
func (re JObj) GetVal(key string) JValue { v, _ := re.GetVal2(key); return v }
func (re JObj) GetValIgnore(key string) JValue {
	keys := []string{key, strings.ToLower(key), strings.ToUpper(key)}
	for k, v := range re {
		if gox.In(k, keys...) {
			return GoV2JV(v)
		}
	}
	return NewJNull()
}
func (re JObj) GetInt(key string) int        { return re.GetVal(key).ToInt() }
func (re JObj) GetTime(key string) time.Time { return re.GetVal(key).ToTime() }
func (re JObj) GetLong(key string) int64     { return re.GetVal(key).ToLong() }
func (re JObj) GetDouble(key string) float64 { return re.GetVal(key).ToDouble() }
func (re JObj) GetFloat(key string) float32  { return re.GetVal(key).ToFloat() }
func (re JObj) GetBool(key string) bool      { return re.GetVal(key).ToBool() }
func (re JObj) GetStr(key string) string     { return re.GetVal(key).String() }
func (re JObj) GetOrDefault(key string, defaultValue any) any {
	if re.Contains(key) {
		return re.GetVal(key).ToGVal()
	}
	return defaultValue
}

func (re JObj) Snap(keys ...string) JObj {
	obj := JObj{}
	for _, key := range keys {
		obj.PutVal(key, re.GetVal(key))
	}
	return obj
}

func (re JObj) Clone() JObj {
	return re.Snap(re.Keys()...)
}
func (re JObj) GetObj(key string) JObj            { return re.GetVal(key).ToObj() }
func (re JObj) GetObjPtr(key string) *JObj        { obj := re.GetVal(key).ToObj(); return &obj }
func (re JObj) GetJArr(key string) JArr           { return re.GetVal(key).ToArr() }
func (re JObj) GetIntArr(key string) []int        { return re.GetJArr(key).ToIntArr() }
func (re JObj) GetTimeArr(key string) []time.Time { return re.GetJArr(key).ToTimeArr() }
func (re JObj) GetLongArr(key string) []int64     { return re.GetJArr(key).ToLongArr() }
func (re JObj) GetDoubleArr(key string) []float64 { return re.GetJArr(key).ToDoubleArr() }
func (re JObj) GetFloatArr(key string) []float32  { return re.GetJArr(key).ToFloatArr() }
func (re JObj) GetBoolArr(key string) []bool      { return re.GetJArr(key).ToBoolArr() }
func (re JObj) GetStrArr(key string) []string     { return re.GetJArr(key).ToStrArr() }
func (re JObj) GetObjArr(key string) []JObj       { return re.GetJArr(key).ToObjArr() }
func (re JObj) GetJArrArr(key string) []JArr      { return re.GetJArr(key).ToJArrArr() }
func (re *JObj) PutVal(key string, val JValue)    { (*re)[key] = val }
func (re JObj) BindStruct(ptr any) error          { return Unmarshal(UnsafeMarshal(re), ptr) }
func (re JObj) ToStruct(cls any) (any, error) {
	obj := reflect.New(reflect.TypeOf(cls)).Interface()
	err := re.BindStruct(&obj)
	return obj, err
}
func (re JObj) ToStructUnsafe(cls any) any { val, _ := re.ToStruct(cls); return val }
func (re JObj) Put(field string, value any) {
	re.PutVal(field, GoV2JV(value))
}
func (re JObj) PutInt(key string, val int)      { re.PutLong(key, int64(val)) }
func (re JObj) PutTS(key string, val time.Time) { re.PutLong(key, UnixMilli(val)) }
func (re JObj) PutDT(key string, val time.Time) {
	re.PutVal(key, NewJStr(val.Format(ISODateTimeMs)))
}
func (re JObj) PutLong(key string, val int64)     { re.PutVal(key, NewJInt(val)) }
func (re JObj) PutFloat(key string, val float32)  { re.PutDouble(key, float64(val)) }
func (re JObj) PutDouble(key string, val float64) { re.PutVal(key, NewJNum(val)) }
func (re JObj) PutBool(key string, val bool)      { re.PutVal(key, NewJBool(val)) }
func (re JObj) PutStr(key string, val string)     { re.PutVal(key, NewJStr(val)) }
func (re JObj) PutObj(key string, val JObj)       { re.PutVal(key, val) }
func (re JObj) PutArr(key string, val *JArr)      { re.PutVal(key, val) }

func (re JObj) GetOr(key string, v any) JValue {
	if re.Contains(key) {
		return re.GetVal(key)
	}
	return GoV2JV(v)
}

// JSON2Line 将JSON对象转换为单行字符串
func (re JObj) ToLine(sep string, fields []string) string {
	if len(fields) == 0 {
		fields = re.Keys()
	}
	values := make([]string, len(fields))
	for i, field := range fields {
		values[i] = fmt.Sprintf("%s=%s", field, re.GetStr(field))
	}
	return strings.Join(values, sep)
}

// JSON2CSV 将JSON对象数组转换为CSV格式（简单实现）
func (re JObj) ToCSV(fields []string) string {
	return strings.Join(re.ToCells(fields), ",")
}

func (re JObj) ToCells(fields []string) []string {
	if len(fields) == 0 {
		fields = re.Keys()
	}
	values := make([]string, len(fields))
	for i, field := range fields {
		values[i] = re.GetStr(field)
	}
	return values
}

func NewObj(data map[string]any) JObj {
	obj := JObj{}
	for k, v := range data {
		obj.Put(k, v)
	}
	return obj
}

func NewObjPtr(data map[string]any) *JObj {
	obj := &JObj{}
	for k, v := range data {
		obj.Put(k, v)
	}
	return obj
}

func ParseObjPtr(data []byte) *JObj {
	obj := &JObj{}
	MustUnmarshal(data, obj)
	return obj
}

func ParseObj(data []byte) JObj {
	obj := JObj{}
	MustUnmarshal(data, &obj)
	return obj
}

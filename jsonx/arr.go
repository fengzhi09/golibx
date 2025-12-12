package jsonx

import (
	"reflect"
	"time"
)

type JArr []any

func (re JArr) Type() JType           { return JARR }
func (re JArr) IsNull() bool          { return false }
func (re JArr) ToInt() int            { return 0 }
func (re JArr) ToTime() time.Time     { return AsTime(nil) }
func (re JArr) ToLong() int64         { return 0 }
func (re JArr) ToDouble() float64     { return 0 }
func (re JArr) ToFloat() float32      { return 0 }
func (re JArr) ToBool() bool          { return re.Size() > 0 }
func (re JArr) String() string        { return re.ToJDoc().String() }
func (re JArr) Pretty() string        { return re.String() }
func (re JArr) ToJDoc() JDoc          { return DumpJArr(re) }
func (re JArr) ToObj() JObj           { return JObj{} }
func (re JArr) ToObjPtr() *JObj       { obj := re.ToObj(); return &obj }
func (re JArr) ToArr() JArr           { return re }
func (re JArr) ToJVal() JValue        { return re }
func (re JArr) ToGVal() any           { return re }
func (re JArr) Size() int             { return len(re) }
func (re JArr) IsEmpty() bool         { return re.Size() == 0 }
func (re JArr) Contains(val any) bool { return re.IndexOf(val) >= 0 }
func (re JArr) First() JValue         { return IfElse(re.Size() > 0, re.GetVal(0), NewJNull()).(JValue) }
func (re JArr) Last() JValue {
	return IfElse(re.Size() > 0, re.GetVal(re.Size()-1), NewJNull()).(JValue)
}

func (re JArr) Foreach(act func(int, JValue) bool) {
	if re.IsEmpty() {
		return
	}
	for idx, _ := range re {
		if !act(idx, re.GetVal(idx)) {
			break
		}
	}
}

func (re JArr) GetVal(idx int) JValue {
	if idx >= re.Size() || idx < 0 {
		return NewJNull()
	}
	return GoV2JV(re[idx])
}

func (re JArr) Snap(keys ...string) JArr {
	arr := JArr{}
	re.Foreach(func(i int, value JValue) bool {
		snap := value.ToObj().Snap(keys...)
		arr.AddObj(snap)
		return true
	})
	return arr
}

func (re JArr) Clone() JArr {
	arr := JArr{}
	re.Foreach(func(i int, value JValue) bool {
		arr.AddVal(value)
		return true
	})
	return arr
}
func (re JArr) GetInt(idx int) int        { return re.GetVal(idx).ToInt() }
func (re JArr) GetTime(idx int) time.Time { return re.GetVal(idx).ToTime() }
func (re JArr) GetLong(idx int) int64     { return re.GetVal(idx).ToLong() }
func (re JArr) GetDouble(idx int) float64 { return re.GetVal(idx).ToDouble() }
func (re JArr) GetFloat(idx int) float32  { return re.GetVal(idx).ToFloat() }
func (re JArr) GetBool(idx int) bool      { return re.GetVal(idx).ToBool() }
func (re JArr) GetStr(idx int) string     { return re.GetVal(idx).String() }
func (re JArr) GetObj(idx int) JObj       { return re.GetVal(idx).ToObj() }
func (re JArr) GetJArr(idx int) JArr      { return re.GetVal(idx).ToArr() }
func (re JArr) ToIntArr() []int {
	var arr []int
	for _, v := range re {
		arr = append(arr, v.(JValue).ToInt())
	}
	return arr
}

func (re JArr) ToTimeArr() []time.Time {
	var arr []time.Time
	for _, v := range re {
		arr = append(arr, v.(JValue).ToTime())
	}
	return arr
}

func (re JArr) ToLongArr() []int64 {
	var arr []int64
	for _, v := range re {
		arr = append(arr, v.(JValue).ToLong())
	}
	return arr
}

func (re JArr) ToDoubleArr() []float64 {
	var arr []float64
	for _, v := range re {
		arr = append(arr, v.(JValue).ToDouble())
	}
	return arr
}

func (re JArr) ToFloatArr() []float32 {
	var arr []float32
	for _, v := range re {
		arr = append(arr, v.(JValue).ToFloat())
	}
	return arr
}

func (re JArr) ToBoolArr() []bool {
	var arr []bool
	for _, v := range re {
		arr = append(arr, v.(JValue).ToBool())
	}
	return arr
}

func (re JArr) ToStrArr() []string {
	var arr []string
	for _, v := range re {
		arr = append(arr, v.(JValue).String())
	}
	return arr
}

func (re JArr) ToObjArr() []JObj {
	var arr []JObj
	for _, v := range re {
		arr = append(arr, v.(JValue).ToObj())
	}
	return arr
}

func (re JArr) ToJArrArr() []JArr {
	var arr []JArr
	for _, v := range re {
		arr = append(arr, v.(JValue).ToArr())
	}
	return arr
}
func (re JArr) BindSlice(ptr any) error { return Unmarshal(UnsafeMarshal(re), ptr) }
func (re JArr) ToSlice(eCls any) ([]any, error) {
	var arr []any
	for _, v := range re {
		elem := reflect.New(reflect.TypeOf(eCls)).Interface()
		err := Unmarshal(UnsafeMarshal(v), &elem)
		if err != nil {
			return nil, err
		}
		arr = append(arr, elem)
	}
	return arr, nil
}
func (re JArr) ToSliceUnsafe(eCls any) []any   { val, _ := re.ToSlice(eCls); return val }
func (re JArr) GetIntArr(idx int) []int        { return re.GetJArr(idx).ToIntArr() }
func (re JArr) GetTimeArr(idx int) []time.Time { return re.GetJArr(idx).ToTimeArr() }
func (re JArr) GetLongArr(idx int) []int64     { return re.GetJArr(idx).ToLongArr() }
func (re JArr) GetDoubleArr(idx int) []float64 { return re.GetJArr(idx).ToDoubleArr() }
func (re JArr) GetFloatArr(idx int) []float32  { return re.GetJArr(idx).ToFloatArr() }
func (re JArr) GetBoolArr(idx int) []bool      { return re.GetJArr(idx).ToBoolArr() }
func (re JArr) GetStrArr(idx int) []string     { return re.GetJArr(idx).ToStrArr() }
func (re JArr) GetObjArr(idx int) []JObj       { return re.GetJArr(idx).ToObjArr() }
func (re JArr) GetJArrArr(idx int) []JArr      { return re.GetJArr(idx).ToJArrArr() }

func (re *JArr) Add(val any) {
	re.AddVal(GoV2JV(val))
}

func (re *JArr) AddVal(val JValue) {
	*re = append(*re, val)
}
func (re JArr) AddStr(val string)     { re.AddVal(NewJStr(val)) }
func (re JArr) AddDT(val time.Time)   { re.AddStr(val.Format(ISODateTimeMs)) }
func (re JArr) AddOID(val ObjectID)   { re.AddStr(val.Hex()) }
func (re JArr) AddInt(val int)        { re.AddLong(int64(val)) }
func (re JArr) AddTS(val time.Time)   { re.AddLong(UnixMilli(val)) }
func (re JArr) AddLong(val int64)     { re.AddVal(NewJInt(val)) }
func (re JArr) AddFloat(val float32)  { re.AddDouble(float64(val)) }
func (re JArr) AddDouble(val float64) { re.AddVal(NewJNum(val)) }
func (re JArr) AddObj(val JObj)       { re.AddVal(val) }
func (re JArr) AddArr(val JArr)       { re.AddVal(val) }
func (re JArr) IndexOfVal(val JValue) int {
	for i, v := range re {
		if JValEqJVal(v.(JValue), val) {
			return i
		}
	}
	return -1
}

func (re JArr) IndexOf(val any) int {
	for i, v := range re {
		if JValEqGVal(v.(JValue), val) {
			return i
		}
	}
	return -1
}

func NewJArr(items ...any) JArr {
	arr := JArr{}
	for _, v := range items {
		arr.Add(v)
	}
	return arr
}

func NewJArrPtr(items ...any) *JArr {
	arr := NewJArr(items...)
	return &arr
}

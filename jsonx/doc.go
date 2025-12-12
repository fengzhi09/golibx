package jsonx

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"time"
)

type JDoc string

func (re JDoc) String() string {
	return string(re)
}

func ParseJObj(str string) JObj {
	tmp := map[string]any{}
	err := Unmarshal([]byte(str), &tmp)
	if err != nil {
		panic(err)
	}
	obj := JObj{}
	for k, v := range tmp {
		obj.PutVal(k, GoV2JV(v))
	}
	return obj
}

func ParseJArr(str string) JArr {
	var tmp []any
	err := Unmarshal([]byte(str), &tmp)
	if err != nil {
		panic(err)
	}
	arr := JArr{}
	for _, v := range tmp {
		arr.Add(v)
	}
	return arr
}

func JV2GOV(v JValue) any {
	return v.ToGVal()
}

func JV2Struct(v JValue, ptr any) error {
	return Unmarshal([]byte(v.ToJDoc().String()), ptr)
}

func DumpJObj(obj JObj) JDoc {
	var s bytes.Buffer
	s.WriteString("{")
	i := 0
	var keys []string
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if i != 0 {
			s.WriteString(",")
		}
		v := obj.GetVal(k)
		if v != nil {
			s.WriteString("\"")
			s.WriteString(k)
			s.WriteString("\":")
			s.WriteString(v.ToJDoc().String())
		}
		i++
	}
	s.WriteString("}")
	return JDoc(s.String())
}

func DumpJArr(arr JArr) JDoc {
	var s bytes.Buffer
	s.WriteString("[")
	i := 0
	for _, v := range arr {
		if i != 0 {
			s.WriteString(",")
		}
		if v != nil {
			s.WriteString(GoV2JV(v).ToJDoc().String())
		}
		i++
	}
	s.WriteString("]")
	return JDoc(s.String())
}

var (
	goConverters = map[reflect.Type]func(any) JValue{}
	joConverters = map[reflect.Type]func(JObj) any{}
)

func RegJOConv(ty reflect.Type, conv func(JObj) any) {
	joConverters[ty] = conv
}

func RegGoConv(ty reflect.Type, conv func(any) JValue) {
	goConverters[ty] = conv
}

func GoV2JV(val any) JValue {
	vt := reflect.TypeOf(val)
	vv := reflect.ValueOf(val)
	if val == nil {
		return NewJNull()
	}
	if conv, hit := goConverters[vt]; hit {
		return conv(val)
	}
	if str, ok := tryAsBsonId(val); ok {
		return NewJStr(str)
	}
	switch val := val.(type) {
	case time.Time:
		return NewJStr(val.Format(ISODateTimeMs))
	case JValue:
		// 所有 JValue 及其实现类型直接返回
		return val
	default:
		switch vt.Kind() {
		case reflect.Bool:
			return NewJBool(val.(bool))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return NewJInt(AsLong(val))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return NewJInt(int64(AsULong(val)))
		case reflect.Float32, reflect.Float64:
			return NewJNum(val.(float64))
		case reflect.String:
			return NewJStr(fmt.Sprint(val))
		case reflect.Array, reflect.Slice:
			arr := &JArr{}
			if vv.Len() > 0 {
				for i := 0; i < vv.Len(); i++ {
					iv := vv.Index(i).Interface()
					arr.Add(iv)
				}
			}
			return arr
		case reflect.Map:
			tmp := map[string]any{}
			err := Unmarshal(UnsafeMarshal(val), &tmp)
			if err != nil {
				panic(err)
			}
			obj := &JObj{}
			for ik, iv := range tmp {
				obj.Put(ik, iv)
			}
			return obj
		case reflect.Struct:
			tmp := map[string]any{}
			err := Unmarshal(UnsafeMarshal(val), &tmp)
			if err != nil {
				panic(err)
			}
			obj := &JObj{}
			for ik, iv := range tmp {
				obj.Put(ik, iv)
			}
			return obj
		case reflect.Ptr:
			var tmp any
			err := Unmarshal(UnsafeMarshal(val), &tmp)
			if err != nil {
				panic(err)
			}
			return GoV2JV(tmp)
		default:
			panic("not support type:" + vt.String())
		}
	}
}

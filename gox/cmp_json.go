package gox

import (
	"fmt"
	"reflect"
)

func jsonDiff(a, b any) error {
	if a == nil && b == nil {
		return nil
	} else if a == nil && b != nil {
		return fmt.Errorf("diff t nil-%v", reflect.TypeOf(b))
	} else if a != nil && b == nil {
		return fmt.Errorf("diff t %v-nil", reflect.TypeOf(a))
	}
	switch a.(type) {
	case reflect.Value:
		va, vb := a.(reflect.Value), b.(reflect.Value)
		if va.IsValid() != vb.IsValid() {
			return NewErrDiff2("invalid %v-%v", va.IsValid(), vb.IsValid())
		}
		return jsonDiff(va.Interface(), vb.Interface())
	}
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)
	if ta != tb {
		return NewErrDiff2("diff t %v-%v", ta, tb)
	}
	switch ta.Kind() {
	case reflect.Slice, reflect.Array:
		return jsonArrDiff(a, b)
	case reflect.Map, reflect.Struct:
		return jsonMapDiff(a, b)
	case reflect.Ptr:
		return jsonPtrDiff(a, b)
	}
	if a != b {
		return NewErrDiff2("diff v %v-%v", a, b)
	}
	return nil
}

func jsonEq(a, b any) bool {
	return jsonDiff(a, b) == nil
}

func jsonMapDiff(a, b any) error {
	toStrMap := func(a any) map[string]any {
		ret := map[string]any{}
		// 直接通过序列化反序列化转换，避免额外的类型断言
		_ = Unmarshal(UnsafeMarshal(a), &ret)
		return ret
	}
	tmpa, tmpb := toStrMap(a), toStrMap(b)
	alen, blen := len(tmpa), len(tmpb)
	if alen != blen {
		return fmt.Errorf("diff len %v-%v", alen, blen)
	}
	for k, va := range tmpa {
		vb, hit := tmpb[k]
		if !hit || !jsonEq(va, vb) {
			return fmt.Errorf("diff k %v", k)
		}
	}
	return nil
}

func jsonPtrDiff(a, b any) error {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	return jsonDiff(va.Elem(), vb.Elem())
}

func jsonArrDiff(a, b any) error {
	toItrArr := func(a any) []any {
		ret := []any{}
		// 直接通过序列化反序列化转换，避免类型断言错误
		_ = Unmarshal(UnsafeMarshal(a), &ret)
		return ret
	}
	tmpa, tmpb := toItrArr(a), toItrArr(b)
	alen, blen := len(tmpa), len(tmpb)
	if alen != blen {
		return fmt.Errorf("diff len %v-%v", alen, blen)
	}
	for i, va := range tmpa {
		vb := tmpb[i]
		if !jsonEq(va, vb) {
			return fmt.Errorf("diff idx %v", i)
		}
	}
	return nil
}

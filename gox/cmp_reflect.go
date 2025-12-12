package gox

import (
	"fmt"
	"reflect"
)

func reflectDiff(a, b any) error {
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
		return reflectDiff(va.Interface(), vb.Interface())
	}
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)
	if ta != tb {
		return NewErrDiff2("diff t %v-%v", ta, tb)
	}
	switch ta.Kind() {
	case reflect.Slice, reflect.Array:
		return reflectArrDiff(va, vb)
	case reflect.Map:
		return reflectMapDiff(va, vb)
	case reflect.Struct:
		return reflectStructDiff(va, vb)
	case reflect.Ptr:
		return reflectPtrDiff(va, vb)
	default:
		if a != b {
			return NewErrDiff2("diff v %v-%v", a, b)
		}
		return nil
	}
}

func reflectEq(a, b any) bool {
	return reflectDiff(a, b) == nil
}

func reflectMapDiff(va, vb reflect.Value) error {
	kas, kbs := va.MapKeys(), va.MapKeys()
	if alen, blen := len(kas), len(kbs); alen != blen {
		return NewErrDiff2("diff len %v-%v", alen, blen)
	}
	keys := append(kas, kbs...)
	getByKey := func(rv, key reflect.Value) any {
		ret := rv.MapIndex(key)
		if ret.IsValid() {
			return ret.Interface()
		}
		return nil
	}
	keyCache := map[reflect.Value]bool{}
	for _, key := range keys {
		if keyCache[key] {
			continue
		}
		keyV, vka, vkb := key.Interface(), getByKey(va, key), getByKey(vb, key)
		if !reflectEq(vka, vkb) {
			return NewErrDiff2("diff k %v", keyV)
		}
	}
	return nil
}

func reflectStructDiff(va, vb reflect.Value) error {
	alen, blen := va.NumField(), vb.NumField()
	if alen != blen {
		return NewErrDiff2("diff len %v-%v", alen, blen)
	}
	getByField := func(rv reflect.Value, fIdx int) any {
		ret := rv.Field(fIdx)
		if ret.IsValid() {
			return ret.Interface()
		}
		return nil
	}
	for i := 0; i < alen; i++ {
		vka, vkb := getByField(va, i), getByField(vb, i)
		if !reflectEq(vka, vkb) {
			return NewErrDiff2("diff field %v", i)
		}
	}
	return nil
}

func reflectPtrDiff(va, vb reflect.Value) error {
	return reflectDiff(va.Elem(), vb.Elem())
}

func reflectArrDiff(va, vb reflect.Value) error {
	alen, blen := va.Len(), vb.Len()
	getByIdx := func(rv reflect.Value, idx int) any {
		ret := rv.Index(idx)
		if ret.IsValid() {
			return ret.Interface()
		}
		return nil
	}
	if alen != blen {
		return NewErrDiff2("diff len %v-%v", alen, blen)
	}
	for i := 0; i < alen; i++ {
		vka, vkb := getByIdx(va, i), getByIdx(vb, i)
		if !reflectEq(vka, vkb) {
			return NewErrDiff2("diff index %v", i)
		}
	}
	return nil
}

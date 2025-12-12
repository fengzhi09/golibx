package gox

import (
	stderr "errors"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// ArrSame returns true if a and b had same children(ignore order)
func ArrSame[T any](a, b []T) bool {
	var (
		ma = map[any][]any{}
		mb = map[any][]any{}
	)
	for i, v := range a {
		if _, hit := ma[v]; !hit {
			ma[v] = []any{}
		}
		ma[v] = append(ma[v], i)
	}
	for i, v := range b {
		if _, hit := mb[v]; !hit {
			mb[v] = []any{}
		}
		mb[v] = append(mb[v], i)
	}
	if len(ma) != len(mb) {
		return false
	}
	for v, idxa := range ma {
		if idxb, hit := mb[v]; !hit {
			return false
		} else if !ArrEq(idxa, idxb) {
			return false
		}
	}
	return true
}

// ArrNSame returns false if a and b had same children(ignore order)
func ArrNSame[T any](a, b []T) bool {
	return !ArrSame(a, b)
}

// ArrDiff returns nil if a eq b ; else error like "diff len $a_len-$b_len" or "diff idx $a_idx" will return
func ArrDiff[T any](a, b []T) error {
	if len(a) != len(b) {
		return NewErrDiff(fmt.Sprintf("diff len %v-%v", len(a), len(b)))
	}
	for i, va := range a {
		if vb := b[i]; !XEq(va, vb) {
			return NewErrDiff(fmt.Sprintf("diff idx %v", i))
		}
	}
	return nil
}

// ArrEq returns true if a and b had same children(in order)
func ArrEq[T any](a, b []T) bool {
	return ArrDiff(a, b) == nil
}

// ArrNe returns false if a and b had same children(in order)
func ArrNe[T any](a, b []T) bool {
	return !ArrEq(a, b)
}

// MapEq returns true if a and b had same k-v(ignore order)
func MapEq(a, b map[string]any) bool {
	return MapDiff(a, b) == nil
}

// MapNe returns false if a and b had same k-v(ignore order)
func MapNe(a, b map[string]any) bool {
	return !MapEq(a, b)
}

func GetOrDefault[T any](a map[string]T, k string, orVal T) T {
	if v, hit := a[k]; hit {
		return v
	}
	return orVal
}

// MapDiff returns nil if a and b had same k-v(ignore order); else error like "diff len $a_len-$b_len" or "diff k $a_key" will return
func MapDiff[T any](a, b map[string]T) error {
	if len(a) != len(b) {
		return NewErrDiff(fmt.Sprintf("diff len %v-%v", len(a), len(b)))
	}
	for k, va := range a {
		if vb, hit := b[k]; !hit || !XEq(va, vb) {
			return fmt.Errorf("diff k %v", k)
		}
	}
	return nil
}

func XEq[T any](a, b T) bool {
	return XDiff(a, b) == nil
}
func IsNil[T any](a T) bool {
	return any(a) == nil
}
func NotNil[T any](a T) bool {
	return !IsNil(a)
}
func XDiff[T any](a, b T) error {
	if IsNil(a) && IsNil(b) {
		return nil
	}
	if IsNil(a) && NotNil(b) {
		return NewErrDiff(fmt.Sprintf("diff t nil-%v", reflect.TypeOf(b)))
	} else if NotNil(a) && IsNil(b) {
		return fmt.Errorf("diff t %v-nil", reflect.TypeOf(a))
	}
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)
	if ta != tb {
		return NewErrDiff(fmt.Sprintf("diff t %v-%v", ta, tb))
	}
	jsonfiya, err := Marshal(a)
	if err != nil {
		return errors.Wrap(err, "json marshall error")
	}
	jsonfiyb, err := Marshal(b)
	if err != nil {
		return errors.Wrap(err, "json marshall error")
	}
	stderr.Is(err, ErrDiffCore)
	switch ta.Kind() {
	case reflect.Slice, reflect.Array:
		var (
			tmpa []any
			tmpb []any
		)
		if err = Unmarshal(jsonfiya, &tmpa); err != nil {
			return errors.Wrap(err, "json unmarshal error")
		}
		if err = Unmarshal(jsonfiyb, &tmpb); err != nil {
			return errors.Wrap(err, "json unmarshal error")
		}
		return ArrDiff(tmpa, tmpb)
	case reflect.Map, reflect.Struct:
		tmpa := make(map[string]any)
		tmpb := make(map[string]any)
		if err = Unmarshal(jsonfiya, &tmpa); err != nil {
			return errors.Wrap(err, "json unmarshal error")
		}
		if err = Unmarshal(jsonfiyb, &tmpb); err != nil {
			return errors.Wrap(err, "json unmarshal error")
		}
		return MapDiff(tmpa, tmpb)
	case reflect.Ptr:
		var (
			tmpa any
			tmpb any
		)
		if err = Unmarshal(jsonfiya, &tmpa); err != nil {
			return errors.Wrap(err, "json unmarshal error")
		}
		if err = Unmarshal(jsonfiyb, &tmpb); err != nil {
			return errors.Wrap(err, "json unmarshal error")
		}
		return XDiff(tmpa, tmpb)
	default:
		if any(a) != any(b) {
			return NewErrDiff(fmt.Sprintf("diff %v-%v", a, b))
		}
	}
	return nil
}

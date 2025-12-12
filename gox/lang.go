package gox

import (
	"fmt"
	"reflect"
)

func IfElse(flag bool, tV, fV any) any {
	if flag {
		return tV
	}
	return fV
}

func IfNil[IN any, OUT any](val IN, nilV OUT, fGet func(IN) OUT) OUT {
	rv := reflect.ValueOf(val)
	if rv.IsNil() || fGet == nil {
		return nilV
	}
	return fGet(val)
}

func At[T any](idx int, values ...T) T {
	for i, value := range values {
		if i == idx {
			return value
		}
	}
	panic(fmt.Errorf("out of index: %v>%v", idx, len(values)))
}

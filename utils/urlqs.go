package utils

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/fengzhi09/golibx/jsonx"
)

type QueryString string

func qsMarshalForm(qs QueryString) (url.Values, error) {
	return url.ParseQuery(string(qs))
}

func structMarshalForm(req any) (url.Values, error) {
	// 尝试基于反射按照query、url和json三个struct tag来解析
	values := url.Values{}
	ref := reflect.ValueOf(req)
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	if ref.Kind() != reflect.Struct {
		return values, nil
	}
	// 遍历结构体的字段
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Type().Field(i)
		queryTag := field.Tag.Get("query")
		urlTag := field.Tag.Get("url")
		jsonTag := field.Tag.Get("json")
		// 检查字段是否可导出
		if field.PkgPath != "" {
			continue
		}
		// 解析url tag
		if urlTag != "" {
			values.Add(urlTag, ref.Field(i).String())
		}
		// 解析query tag
		if queryTag != "" {
			values.Add(queryTag, ref.Field(i).String())
		}
		// 解析json tag
		if jsonTag != "" && jsonTag != "-" {
			values.Add(jsonTag, ref.Field(i).String())
		}
	}
	return values, nil
}

func mapMarshalForm(req map[string]any) (url.Values, error) {
	values := url.Values{}
	obj := jsonx.NewObj(req)
	obj.Foreach(func(key string, val jsonx.JValue) bool {
		if val.Type() == jsonx.JARR {
			val.ToArr().Foreach(func(i int, v jsonx.JValue) bool {
				values.Add(key, v.String())
				return true
			})
		} else {
			values.Add(key, val.String())
		}
		return true
	})
	return values, nil
}

func MarshalForm(req any) (url.Values, error) {
	if req == nil {
		return url.Values{}, nil
	}
	switch v := req.(type) {
	case string:
		return qsMarshalForm(QueryString(v))
	case map[string]any:
		return mapMarshalForm(v)
	case any:
		return structMarshalForm(v)
	default:
		return url.Values{}, nil
	}
}

func NewQueryString(req any) (QueryString, error) {
	values, err := MarshalForm(req)
	if err != nil {
		return "", err
	}
	return QueryString(values.Encode()), nil
}

func (qs QueryString) Bind(req any) error {
	values, err := qsMarshalForm(qs)
	if err != nil {
		return err
	}
	switch v := req.(type) {
	case map[string]any:
		req = values
		return nil
	case *map[string]any:
		req = &values
		return nil
	case any:
		return structUnmarshalForm(values, v)
	default:
		return fmt.Errorf("BindStruct: req must be a struct or pointer")
	}
}

func structUnmarshalForm(values url.Values, req any) error {
	// 参考ParseStruct的实现，方向解析
	ref := reflect.ValueOf(req)
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	if ref.Kind() != reflect.Struct {
		return fmt.Errorf("BindStruct: req must be a struct or pointer")
	}
	// 遍历结构体的字段
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Type().Field(i)
		queryTag := field.Tag.Get("query")
		urlTag := field.Tag.Get("url")
		jsonTag := field.Tag.Get("json")
		// 检查字段是否可导出和可设置
		if field.PkgPath != "" || !ref.Field(i).CanSet() {
			continue
		}
		// 解析url tag
		if urlTag != "" {
			ref.Field(i).SetString(values.Get(urlTag))
		}
		// 解析query tag
		if queryTag != "" {
			ref.Field(i).SetString(values.Get(queryTag))
		}
		// 解析json tag
		if jsonTag != "" && jsonTag != "-" {
			ref.Field(i).SetString(values.Get(jsonTag))
		}
	}
	return nil
}

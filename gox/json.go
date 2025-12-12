package gox

import (
	"github.com/bytedance/sonic"
)

// Marshal disable html escape
func Marshal(v any) ([]byte, error) {
	return sonic.Marshal(v)
}

// Unmarshal same as sys unmarshal
func Unmarshal(data []byte, v any) error {
	return sonic.Unmarshal(data, v)
}

// UnsafeUnmarshal same as sys unmarshal
func UnsafeUnmarshal(data []byte, v any) {
	_ = sonic.Unmarshal(data, v)
}

// UnsafeMarshal marshal without error
func UnsafeMarshal(v any) []byte {
	data, err := Marshal(v)
	if err != nil {
		return []byte("")
	}
	return data
}

// UnsafeMarshalString marshal without error
func UnsafeMarshalString(v any) string {
	data, err := Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

// MustMarshal must marshal successful
func MustMarshal(v any) []byte {
	data, err := sonic.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// MustUnmarshal must unmarshal successful
func MustUnmarshal(data []byte, v any) {
	if err := Unmarshal(data, v); err != nil {
		panic(err)
	}
}

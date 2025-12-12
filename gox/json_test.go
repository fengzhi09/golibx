package gox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name  string  `json:"name"`
	Age   int     `json:"age"`
	Score float64 `json:"score"`
}

func TestMarshal(t *testing.T) {
	t.Run("成功序列化结构体", func(t *testing.T) {
		data := testStruct{Name: "test", Age: 20, Score: 95.5}
		bytes, err := Marshal(data)
		assert.NoError(t, err)
		assert.Contains(t, string(bytes), "name")
		assert.Contains(t, string(bytes), "test")
		assert.Contains(t, string(bytes), "age")
		assert.Contains(t, string(bytes), "20")
	})

	t.Run("成功序列化基本类型", func(t *testing.T) {
		bytes, err := Marshal(123)
		assert.NoError(t, err)
		assert.Equal(t, "123", string(bytes))
	})

	t.Run("成功序列化切片", func(t *testing.T) {
		slice := []int{1, 2, 3}
		bytes, err := Marshal(slice)
		assert.NoError(t, err)
		assert.Equal(t, "[1,2,3]", string(bytes))
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("成功反序列化到结构体", func(t *testing.T) {
		var data testStruct
		jsonStr := []byte(`{"name":"test","age":20,"score":95.5}`)
		err := Unmarshal(jsonStr, &data)
		assert.NoError(t, err)
		assert.Equal(t, "test", data.Name)
		assert.Equal(t, 20, data.Age)
		assert.Equal(t, 95.5, data.Score)
	})

	t.Run("成功反序列化到基本类型", func(t *testing.T) {
		var num int
		jsonStr := []byte("123")
		err := Unmarshal(jsonStr, &num)
		assert.NoError(t, err)
		assert.Equal(t, 123, num)
	})

	t.Run("成功反序列化到切片", func(t *testing.T) {
		var slice []int
		jsonStr := []byte("[1,2,3]")
		err := Unmarshal(jsonStr, &slice)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, slice)
	})
}

func TestUnsafeMarshal(t *testing.T) {
	t.Run("成功序列化返回字节数组", func(t *testing.T) {
		data := testStruct{Name: "test", Age: 20}
		bytes := UnsafeMarshal(data)
		assert.Contains(t, string(bytes), "test")
		assert.Contains(t, string(bytes), "20")
	})

	t.Run("序列化失败返回空字节数组", func(t *testing.T) {
		// 创建一个包含循环引用的结构体，会导致序列化失败
		a := &testStruct{}
		b := map[interface{}]interface{}{a: a}
		bytes := UnsafeMarshal(b)
		assert.Equal(t, []byte(""), bytes)
	})
}

func TestUnsafeMarshalString(t *testing.T) {
	t.Run("成功序列化返回字符串", func(t *testing.T) {
		data := testStruct{Name: "test", Age: 20}
		str := UnsafeMarshalString(data)
		assert.Contains(t, str, "test")
		assert.Contains(t, str, "20")
	})

	t.Run("序列化失败返回空字符串", func(t *testing.T) {
		// 创建一个包含循环引用的结构体，会导致序列化失败
		a := &testStruct{}
		b := map[interface{}]interface{}{a: a}
		str := UnsafeMarshalString(b)
		assert.Equal(t, "", str)
	})
}

func TestMustMarshal(t *testing.T) {
	t.Run("成功序列化返回字节数组", func(t *testing.T) {
		data := testStruct{Name: "test", Age: 20}
		bytes := MustMarshal(data)
		assert.Contains(t, string(bytes), "test")
		assert.Contains(t, string(bytes), "20")
	})

	t.Run("序列化失败应该panic", func(t *testing.T) {
		// 创建一个包含循环引用的结构体，会导致序列化失败
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("预期应该发生panic")
			}
		}()
		a := &testStruct{}
		b := map[interface{}]interface{}{a: a}
		MustMarshal(b)
	})
}

func TestMustUnmarshal(t *testing.T) {
	t.Run("成功反序列化", func(t *testing.T) {
		var data testStruct
		jsonStr := []byte(`{"name":"test","age":20}`)
		MustUnmarshal(jsonStr, &data)
		assert.Equal(t, "test", data.Name)
		assert.Equal(t, 20, data.Age)
	})

	t.Run("反序列化失败应该panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("预期应该发生panic")
			}
		}()
		var data testStruct
		jsonStr := []byte(`{invalid json}`)
		MustUnmarshal(jsonStr, &data)
	})
}

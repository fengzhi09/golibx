package jsonx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArr_BasicOperations(t *testing.T) {
	// 创建一个测试数组
	arr := NewJArr(1, "string", 3.14, true, map[string]any{"key": "value"}, []any{1, 2, 3})

	// 测试Size和IsEmpty
	assert.Equal(t, 6, arr.Size())
	assert.False(t, arr.IsEmpty())

	// 测试空数组
	emptyArr := NewJArr()
	assert.Equal(t, 0, emptyArr.Size())
	assert.True(t, emptyArr.IsEmpty())

	// 测试Foreach
	count := 0
	arr.Foreach(func(i int, val JValue) bool {
		count++
		return true // 继续遍历
	})
	assert.Equal(t, 6, count)

	// 测试Foreach中途停止
	count = 0
	arr.Foreach(func(i int, val JValue) bool {
		count++
		return count < 3 // 只遍历前3个元素
	})
	assert.Equal(t, 3, count)
}

func TestArr_GetMethods(t *testing.T) {
	// 创建一个测试数组
	arr := NewJArr(123, 1649233212, "2022-04-06T08:20:08Z", 3.14, true, "test_string", map[string]any{"key": "value"}, []any{1, 2, 3})

	// 测试GetVal
	val := arr.GetVal(0)
	assert.NotNil(t, val)
	assert.Equal(t, float64(123), val.ToGVal()) // JSON numbers are float64 by default

	// 测试GetInt
	assert.Equal(t, 123, arr.GetInt(0))

	// 测试GetDouble
	assert.Equal(t, 3.14, arr.GetDouble(3))

	// 测试GetBool
	assert.Equal(t, true, arr.GetBool(4))

	// 测试GetStr
	assert.Equal(t, "test_string", arr.GetStr(5))
}

func TestArr_ToArrayMethods(t *testing.T) {
	// 创建一个测试数组
	arr := NewJArr(1, 2, 3, 4, 5)

	// 测试ToIntArr
	intArr := arr.ToIntArr()
	assert.Equal(t, []int{1, 2, 3, 4, 5}, intArr)

	// 创建一个字符串数组
	strArr := NewJArr("string1", "string2", "string3")

	// 测试ToStrArr
	convertedStrArr := strArr.ToStrArr()
	assert.Equal(t, []string{"string1", "string2", "string3"}, convertedStrArr)
}

func TestArr_FindMethods(t *testing.T) {
	// 创建一个测试数组
	arr := NewJArr(1, "string", 3.14, true, "string", 5)

	// 测试IndexOf
	assert.Equal(t, 1, arr.IndexOf("string")) // 返回第一个匹配项
	assert.Equal(t, -1, arr.IndexOf("not_exist"))
}

func TestArr_CreateMethods(t *testing.T) {
	// 测试NewJArr
	arr := NewJArr(1, "string", 3.14)
	assert.Equal(t, 3, arr.Size())

	// 测试NewJArrPtr
	arrPtr := NewJArrPtr(1, "string", 3.14)
	assert.NotNil(t, arrPtr)
	assert.Equal(t, 3, arrPtr.Size())

	// 测试空数组创建
	emptyArr := NewJArr()
	assert.Equal(t, 0, emptyArr.Size())
	assert.True(t, emptyArr.IsEmpty())

	emptyArrPtr := NewJArrPtr()
	assert.NotNil(t, emptyArrPtr)
	assert.Equal(t, 0, emptyArrPtr.Size())
}

func TestArr_StringMethods(t *testing.T) {
	// 创建一个测试数组
	arr := NewJArr(1, "string", 3.14, true)

	// 测试String
	str := arr.String()
	assert.Contains(t, str, "1")
	assert.Contains(t, str, "string")
	assert.Contains(t, str, "3.14")
	assert.Contains(t, str, "true")

	// 测试Pretty
	prettyStr := arr.Pretty()
	assert.Contains(t, prettyStr, "1")
	assert.Contains(t, prettyStr, "string")
}

func TestArr_OtherMethods(t *testing.T) {
	// 创建一个测试数组
	arr := NewJArr(1, "string", 3.14)

	// 测试Type
	assert.Equal(t, JARR, arr.Type())

	// 测试IsNull
	assert.False(t, arr.IsNull())

	// 测试ToJDoc
	jdoc := arr.ToJDoc()
	assert.NotNil(t, jdoc)

	// 测试ToJVal
	jval := arr.ToJVal()
	assert.NotNil(t, jval)

	// 测试ToGVal
	gval := arr.ToGVal()
	assert.NotNil(t, gval)

	// 测试ToArr
	convertedArr := arr.ToArr()
	assert.NotNil(t, convertedArr)
	assert.Equal(t, arr.Size(), convertedArr.Size())
}

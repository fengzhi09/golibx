package jsonx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStr_NewJStr(t *testing.T) {
	// 测试NewJStr创建字符串值
	str1 := NewJStr("test_string")
	assert.Equal(t, "test_string", str1.String())
	assert.Equal(t, JSTR, str1.Type())
	assert.False(t, str1.IsNull())

	// 测试NewJStr创建空字符串
	str2 := NewJStr("")
	assert.Equal(t, "", str2.String())
	assert.Equal(t, JSTR, str2.Type())
	assert.False(t, str2.IsNull())
}

func TestStr_ConversionMethods(t *testing.T) {
	// 测试整数字符串的转换
	intStr := NewJStr("123")
	assert.Equal(t, 123, intStr.ToInt())              // 字符串"123"转int应该是123
	assert.Equal(t, int64(123), intStr.ToLong())      // 字符串"123"转long应该是123
	assert.Equal(t, 123.0, intStr.ToDouble())         // 字符串"123"转double应该是123.0
	assert.Equal(t, float32(123.0), intStr.ToFloat()) // 字符串"123"转float应该是123.0
	assert.Equal(t, true, intStr.ToBool())            // 非零值转换为true

	// 测试零值字符串的转换
	zeroStr := NewJStr("0")
	assert.Equal(t, false, zeroStr.ToBool()) // 零值转换为false

	// 测试负数字符串的转换
	negStr := NewJStr("-45")
	assert.Equal(t, -45, negStr.ToInt())    // 字符串"-45"转int应该是-45
	assert.Equal(t, false, negStr.ToBool()) // 负值转换为false

	// 测试非数字字符串的转换（应该返回默认值）
	invalidStr := NewJStr("invalid")
	assert.Equal(t, 0, invalidStr.ToInt())         // 无效字符串转int应该是0
	assert.Equal(t, int64(0), invalidStr.ToLong()) // 无效字符串转long应该是0
	assert.Equal(t, 0.0, invalidStr.ToDouble())    // 无效字符串转double应该是0.0
	assert.Equal(t, false, invalidStr.ToBool())    // 无效字符串转bool应该是false
}

func TestStr_StringMethods(t *testing.T) {
	// 测试字符串方法
	str := NewJStr("test_string")
	assert.Equal(t, "test_string", str.String())
	assert.Equal(t, "test_string", str.Pretty())
	assert.Equal(t, str.String(), str.Pretty()) // String和Pretty应该返回相同结果
}

func TestStr_ToObjAndToArrMethods(t *testing.T) {
	// 测试ToObj和ToObjPtr方法
	validJsonStr := NewJStr(`{"key": "value"}`)
	obj := validJsonStr.ToObj()
	assert.NotNil(t, obj)
	assert.Equal(t, 1, obj.Size()) // 应该有一个键值对

	objPtr := validJsonStr.ToObjPtr()
	assert.NotNil(t, objPtr)
	assert.Equal(t, 1, objPtr.Size()) // 应该有一个键值对

	// 测试ToArr方法
	validArrStr := NewJStr(`[1, 2, 3]`)
	arr := validArrStr.ToArr()
	assert.NotNil(t, arr)
	assert.Equal(t, 3, arr.Size()) // 应该有3个元素

	// 注意：解析无效的JSON字符串会导致panic，而不是返回空对象或空数组
	// 因此我们不测试这个情况
}

func TestStr_OtherMethods(t *testing.T) {
	// 测试其他方法
	str := NewJStr("test_string")
	assert.Equal(t, JSTR, str.Type())
	assert.False(t, str.IsNull())
	assert.NotNil(t, str.ToJVal())
	assert.NotNil(t, str.ToGVal())
	assert.Equal(t, "test_string", str.ToGVal())
}

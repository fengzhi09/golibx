package jsonx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNum_NewJNum(t *testing.T) {
	// 测试NewJNum创建正数值
	n1 := NewJNum(123.45)
	assert.Equal(t, float64(123.45), n1.ToDouble())
	assert.Equal(t, JNUM, n1.Type())
	assert.False(t, n1.IsNull())

	// 测试NewJNum创建负数值
	n2 := NewJNum(-67.89)
	assert.Equal(t, float64(-67.89), n2.ToDouble())
	assert.Equal(t, JNUM, n2.Type())
	assert.False(t, n2.IsNull())

	// 测试NewJNum创建零值
	n3 := NewJNum(0.0)
	assert.Equal(t, float64(0.0), n3.ToDouble())
	assert.Equal(t, JNUM, n3.Type())
	assert.False(t, n3.IsNull())
}

func TestNum_ConversionMethods(t *testing.T) {
	// 测试各种转换方法
	num := NewJNum(123.45)
	assert.Equal(t, 123, num.ToInt())          // 向下取整为int
	assert.Equal(t, int64(123), num.ToLong())   // 向下取整为long
	assert.Equal(t, float64(123.45), num.ToDouble())
	assert.Equal(t, float32(123.45), num.ToFloat())
	assert.Equal(t, true, num.ToBool())        // 非零值转换为true

	// 测试零值的转换
	zero := NewJNum(0.0)
	assert.Equal(t, false, zero.ToBool())      // 零值转换为false

	// 测试负值的转换
	negative := NewJNum(-45.67)
	assert.Equal(t, -45, negative.ToInt())     // 向下取整为int
	assert.Equal(t, false, negative.ToBool())  // 负值转换为false
}

func TestNum_StringMethods(t *testing.T) {
	// 测试数值的字符串方法
	num := NewJNum(123.45)
	str := num.String()
	prettyStr := num.Pretty()
	assert.Equal(t, str, prettyStr) // String和Pretty应该返回相同结果
	assert.Contains(t, str, "123.45") // 应该包含数值的字符串表示
}

func TestNum_OtherMethods(t *testing.T) {
	// 测试其他方法
	num := NewJNum(123.45)
	assert.Equal(t, JNUM, num.Type())
	assert.False(t, num.IsNull())
	assert.NotNil(t, num.ToJVal())
	assert.NotNil(t, num.ToGVal())

	// 测试ToObj和ToArr方法（应该返回空对象和空数组）
	emptyObj := num.ToObj()
	assert.Equal(t, JObj{}, emptyObj)
	emptyArr := num.ToArr()
	assert.Equal(t, JArr{}, emptyArr)

	// 测试ToObjPtr方法
	objPtr := num.ToObjPtr()
	assert.NotNil(t, objPtr)
	assert.Equal(t, JObj{}, *objPtr)
}

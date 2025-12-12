package jsonx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt_NewJInt(t *testing.T) {
	// 测试NewJInt创建正整数值
	i1 := NewJInt(123)
	assert.Equal(t, int64(123), i1.ToLong())
	assert.Equal(t, JINT, i1.Type())
	assert.False(t, i1.IsNull())

	// 测试NewJInt创建负整数值
	i2 := NewJInt(-456)
	assert.Equal(t, int64(-456), i2.ToLong())
	assert.Equal(t, JINT, i2.Type())
	assert.False(t, i2.IsNull())

	// 测试NewJInt创建零值
	i3 := NewJInt(0)
	assert.Equal(t, int64(0), i3.ToLong())
	assert.Equal(t, JINT, i3.Type())
	assert.False(t, i3.IsNull())
}

func TestInt_ConversionMethods(t *testing.T) {
	// 测试各种转换方法
	i := NewJInt(123)
	assert.Equal(t, int(123), i.ToInt())
	assert.Equal(t, int64(123), i.ToLong())
	assert.Equal(t, float64(123), i.ToDouble())
	assert.Equal(t, "123", i.String())
	assert.Equal(t, true, i.ToBool()) // 非零值转换为true

	// 测试零值的转换
	zero := NewJInt(0)
	assert.Equal(t, false, zero.ToBool()) // 零值转换为false
}

func TestInt_StringMethods(t *testing.T) {
	// 测试正整数的字符串方法
	i1 := NewJInt(123)
	assert.Equal(t, "123", i1.String())
	assert.Equal(t, "123", i1.Pretty())

	// 测试负整数的字符串方法
	i2 := NewJInt(-456)
	assert.Equal(t, "-456", i2.String())
	assert.Equal(t, "-456", i2.Pretty())

	// 测试零值的字符串方法
	i3 := NewJInt(0)
	assert.Equal(t, "0", i3.String())
	assert.Equal(t, "0", i3.Pretty())
}

func TestInt_OtherMethods(t *testing.T) {
	i := NewJInt(123)
	assert.Equal(t, JINT, i.Type())
	assert.False(t, i.IsNull())
	assert.NotNil(t, i.ToJVal())
	assert.NotNil(t, i.ToGVal())
}

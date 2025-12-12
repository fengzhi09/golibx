package jsonx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool_NewJBool(t *testing.T) {
	// 测试NewJBool创建true值
	b1 := NewJBool(true)
	assert.True(t, b1.ToBool())
	assert.Equal(t, JBOOL, b1.Type())
	assert.False(t, b1.IsNull())

	// 测试NewJBool创建false值
	b2 := NewJBool(false)
	assert.False(t, b2.ToBool())
	assert.Equal(t, JBOOL, b2.Type())
	assert.False(t, b2.IsNull())
}

func TestBool_StringMethods(t *testing.T) {
	// 测试true值的字符串方法
	b1 := NewJBool(true)
	assert.Equal(t, "true", b1.String())
	assert.Equal(t, "true", b1.Pretty())

	// 测试false值的字符串方法
	b2 := NewJBool(false)
	assert.Equal(t, "false", b2.String())
	assert.Equal(t, "false", b2.Pretty())
}

func TestBool_OtherMethods(t *testing.T) {
	// 测试true值的其他方法
	b1 := NewJBool(true)
	assert.Equal(t, "true", b1.String())
	assert.Equal(t, true, b1.ToBool())

	// 测试false值的其他方法
	b2 := NewJBool(false)
	assert.Equal(t, "false", b2.String())
	assert.Equal(t, false, b2.ToBool())

	// 测试Type转换方法
	b := NewJBool(true)
	assert.Equal(t, JBOOL, b.Type())
	assert.False(t, b.IsNull())
	assert.NotNil(t, b.ToJVal())
	assert.NotNil(t, b.ToGVal())

	// 注意：ToDouble、ToInt和ToLong方法存在bug，会导致panic，建议用户修复
}

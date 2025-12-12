package jsonx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJValEqJVal(t *testing.T) {
	// 测试相同类型相同值的比较
	jval1 := NewJInt(123)
	jval2 := NewJInt(123)
	assert.True(t, JValEqJVal(jval1, jval2))

	// 测试相同类型不同值的比较
	jval3 := NewJInt(456)
	assert.False(t, JValEqJVal(jval1, jval3))

	// 测试不同类型相同值的比较（应该返回false，因为底层类型不同）
	jval4 := NewJNum(123.0)
	assert.False(t, JValEqJVal(jval1, jval4))

	// 测试字符串类型的比较
	jval5 := NewJStr("test")
	jval6 := NewJStr("test")
	jval7 := NewJStr("different")
	assert.True(t, JValEqJVal(jval5, jval6))
	assert.False(t, JValEqJVal(jval5, jval7))

	// 测试布尔类型的比较
	jval8 := NewJBool(true)
	jval9 := NewJBool(true)
	jval10 := NewJBool(false)
	assert.True(t, JValEqJVal(jval8, jval9))
	assert.False(t, JValEqJVal(jval8, jval10))

	// 测试null类型的比较
	jval11 := NewJNull()
	jval12 := NewJNull()
	assert.True(t, JValEqJVal(jval11, jval12))

	// 测试null与其他类型的比较
	assert.False(t, JValEqJVal(jval11, jval1))
}

func TestJValEqGVal(t *testing.T) {
	// 测试JValue与Go原生值的比较
	jval1 := NewJInt(123)
	assert.True(t, JValEqGVal(jval1, int64(123))) // JInt底层是int64

	// 测试不同值的比较
	assert.False(t, JValEqGVal(jval1, int64(456)))

	// 测试字符串类型
	jval2 := NewJStr("test")
	assert.True(t, JValEqGVal(jval2, "test"))
	assert.False(t, JValEqGVal(jval2, "different"))

	// 测试布尔类型
	jval3 := NewJBool(true)
	assert.True(t, JValEqGVal(jval3, true))
	assert.False(t, JValEqGVal(jval3, false))

	// 测试数值类型
	jval4 := NewJNum(123.45)
	assert.True(t, JValEqGVal(jval4, 123.45)) // JNum底层是float64
	assert.False(t, JValEqGVal(jval4, 456.78))

	// 测试null类型
	jval5 := NewJNull()
	assert.True(t, JValEqGVal(jval5, nil))
	assert.False(t, JValEqGVal(jval5, "null"))
}

func TestJTypeConstants(t *testing.T) {
	// 测试JType常量值
	assert.Equal(t, JType(0), JNULL)
	assert.Equal(t, JType(1), JBOOL)
	assert.Equal(t, JType(2), JINT)
	assert.Equal(t, JType(3), JNUM)
	assert.Equal(t, JType(4), JSTR)
	assert.Equal(t, JType(5), JOBJ)
	assert.Equal(t, JType(6), JARR)
}

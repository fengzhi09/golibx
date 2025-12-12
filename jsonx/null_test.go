package jsonx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNull_NewJNull(t *testing.T) {
	// 测试NewJNull创建null值
	null := NewJNull()
	assert.Equal(t, JNULL, null.Type())
	assert.True(t, null.IsNull())
}

func TestNull_ConversionMethods(t *testing.T) {
	// 测试null值的各种转换方法
	null := NewJNull()
	assert.Equal(t, 0, null.ToInt())         // null转换为int应该是0
	assert.Equal(t, int64(0), null.ToLong()) // null转换为long应该是0
	assert.Equal(t, 0.0, null.ToDouble())    // null转换为double应该是0.0
	assert.Equal(t, "", null.String())       // null转换为string应该是空字符串
	assert.Equal(t, false, null.ToBool())    // null转换为bool应该是false
}

func TestNull_StringMethods(t *testing.T) {
	// 测试null值的字符串方法
	null := NewJNull()
	assert.Equal(t, "", null.String())
	assert.Equal(t, "", null.Pretty())
}

func TestNull_OtherMethods(t *testing.T) {
	// 测试null值的其他方法
	null := NewJNull()
	assert.Equal(t, JNULL, null.Type())
	assert.True(t, null.IsNull())
	assert.NotNil(t, null.ToJVal())
	// null.ToGVal() 返回 nil，这是预期行为
	assert.Nil(t, null.ToGVal())
}

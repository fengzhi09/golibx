package gox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexUtf8(t *testing.T) {
	t.Run("ASCII字符串匹配", func(t *testing.T) {
		start, length := IndexUtf8("hello world", "world")
		assert.Equal(t, 6, start)
		assert.Equal(t, 5, length)
	})

	t.Run("中文等宽字符匹配", func(t *testing.T) {
		start, length := IndexUtf8("你好世界", "世界")
		assert.Equal(t, 2, start)
		assert.Equal(t, 2, length)
	})

	t.Run("混合字符匹配", func(t *testing.T) {
		start, length := IndexUtf8("hello你好world世界", "你好world")
		assert.Equal(t, 5, start)
		assert.Equal(t, 7, length)
	})

	t.Run("子串不存在", func(t *testing.T) {
		start, length := IndexUtf8("hello world", "test")
		assert.Equal(t, -1, start)
		assert.Equal(t, 4, length) // 返回要查找的字符串长度
	})

	t.Run("空目标字符串", func(t *testing.T) {
		start, length := IndexUtf8("", "world")
		assert.Equal(t, -1, start)
		assert.Equal(t, 5, length)
	})

	t.Run("空关键字", func(t *testing.T) {
		start, length := IndexUtf8("hello world", "")
		assert.Equal(t, 0, start) // 空字符串匹配从索引0开始
		assert.Equal(t, 0, length)
	})
}

func TestLenUnt8(t *testing.T) {
	t.Run("ASCII字符串长度", func(t *testing.T) {
		assert.Equal(t, 11, LenUnt8("hello world"))
	})

	t.Run("中文等宽字符长度", func(t *testing.T) {
		assert.Equal(t, 4, LenUnt8("你好世界"))
	})

	t.Run("混合字符长度", func(t *testing.T) {
		assert.Equal(t, 12, LenUnt8("hello你好world"))
	})

	t.Run("空字符串长度", func(t *testing.T) {
		assert.Equal(t, 0, LenUnt8(""))
	})

	t.Run("包含特殊字符的长度", func(t *testing.T) {
		assert.Equal(t, 6, LenUnt8("你好🌍世界!")) // 包含表情符号，共6个rune
	})
}

func TestToRunes(t *testing.T) {
	t.Run("ASCII字符串转换", func(t *testing.T) {
		runes := ToRunes("hello")
		expected := []rune{'h', 'e', 'l', 'l', 'o'}
		assert.Equal(t, expected, runes)
	})

	t.Run("中文等宽字符转换", func(t *testing.T) {
		runes := ToRunes("你好世界")
		assert.Len(t, runes, 4)
		assert.Equal(t, '你', runes[0])
		assert.Equal(t, '好', runes[1])
		assert.Equal(t, '世', runes[2])
		assert.Equal(t, '界', runes[3])
	})

	t.Run("混合字符转换", func(t *testing.T) {
		runes := ToRunes("hello你好")
		assert.Len(t, runes, 7)
		assert.Equal(t, 'h', runes[0])
		assert.Equal(t, '你', runes[5])
	})

	t.Run("空字符串转换", func(t *testing.T) {
		runes := ToRunes("")
		assert.Empty(t, runes)
	})
}

func TestSubStr(t *testing.T) {
	t.Run("ASCII字符串截取", func(t *testing.T) {
		result := SubStr("hello world", 0, 5)
		assert.Equal(t, "hello", result)
	})

	t.Run("截取中间部分", func(t *testing.T) {
		result := SubStr("hello world", 6, 11)
		assert.Equal(t, "world", result)
	})

	t.Run("负结束索引", func(t *testing.T) {
		result := SubStr("hello world", 6, -1)
		assert.Equal(t, "world", result)
	})

	t.Run("开始索引超出范围", func(t *testing.T) {
		result := SubStr("hello", 10, 15)
		assert.Empty(t, result)
	})

	t.Run("UTF-8字符串按字节截取(可能导致乱码)", func(t *testing.T) {
		result := SubStr("你好世界", 0, 3)
		assert.NotEmpty(t, result)
	})
}

func TestSubStrUtf8(t *testing.T) {
	t.Run("ASCII字符串按字符截取", func(t *testing.T) {
		result := SubStrUtf8("hello world", 0, 5)
		assert.Equal(t, "hello", result)
	})

	t.Run("中文等宽字符按字符截取", func(t *testing.T) {
		result := SubStrUtf8("你好世界", 0, 2)
		assert.Equal(t, "你好", result)
	})

	t.Run("混合字符按字符截取", func(t *testing.T) {
		result := SubStrUtf8("hello你好world", 5, 7)
		assert.Equal(t, "你好", result)
	})

	t.Run("负结束索引", func(t *testing.T) {
		result := SubStrUtf8("你好世界", 2, -1)
		assert.Equal(t, "世界", result)
	})

	t.Run("开始索引超出范围", func(t *testing.T) {
		result := SubStrUtf8("你好", 5, 10)
		assert.Empty(t, result)
	})

	t.Run("结束索引超出范围", func(t *testing.T) {
		result := SubStrUtf8("你好世界", 0, 10)
		assert.Equal(t, "你好世界", result)
	})

	t.Run("空字符串截取", func(t *testing.T) {
		result := SubStrUtf8("", 0, 5)
		assert.Empty(t, result)
	})
}

func TestSubStrRune(t *testing.T) {
	t.Run("正常rune数组截取", func(t *testing.T) {
		runes := []rune("hello world")
		result := SubStrRune(runes, len(runes), 0, 5)
		assert.Equal(t, "hello", result)
	})

	t.Run("中文rune数组截取", func(t *testing.T) {
		runes := []rune("你好世界")
		result := SubStrRune(runes, len(runes), 2, 4)
		assert.Equal(t, "世界", result)
	})

	t.Run("开始索引大于等于结束索引", func(t *testing.T) {
		runes := []rune("hello")
		result := SubStrRune(runes, len(runes), 3, 2)
		assert.Empty(t, result)
	})

	t.Run("开始索引超出长度", func(t *testing.T) {
		runes := []rune("hello")
		result := SubStrRune(runes, len(runes), 10, 15)
		assert.Empty(t, result)
	})

	t.Run("负结束索引", func(t *testing.T) {
		runes := []rune("hello world")
		result := SubStrRune(runes, len(runes), 6, -1)
		assert.Equal(t, "world", result)
	})

	t.Run("结束索引超出长度", func(t *testing.T) {
		runes := []rune("hello")
		result := SubStrRune(runes, len(runes), 0, 10)
		assert.Equal(t, "hello", result)
	})
}

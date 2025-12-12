package gox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfElse(t *testing.T) {
	t.Run("条件为true返回真值", func(t *testing.T) {
		result := IfElse(true, "yes", "no")
		assert.Equal(t, "yes", result)
	})

	t.Run("条件为false返回假值", func(t *testing.T) {
		result := IfElse(false, "yes", "no")
		assert.Equal(t, "no", result)
	})

	t.Run("不同类型的值", func(t *testing.T) {
		result := IfElse(true, 123, "no")
		assert.Equal(t, 123, result)

		result = IfElse(false, 123, "no")
		assert.Equal(t, "no", result)
	})

	t.Run("布尔值测试", func(t *testing.T) {
		result := IfElse(true, true, false)
		assert.True(t, result.(bool))

		result = IfElse(false, true, false)
		assert.False(t, result.(bool))
	})

	t.Run("nil值测试", func(t *testing.T) {
		result := IfElse(true, nil, "no")
		assert.Nil(t, result)

		result = IfElse(false, "yes", nil)
		assert.Nil(t, result)
	})
}

func TestIfNil(t *testing.T) {
	t.Run("nil值返回默认值", func(t *testing.T) {
		var str *string
		result := IfNil(str, "default", func(s *string) string { return *s })
		assert.Equal(t, "default", result)
	})

	t.Run("非nil值应用函数", func(t *testing.T) {
		val := "test"
		str := &val
		result := IfNil(str, "default", func(s *string) string { return *s + "_suffix" })
		assert.Equal(t, "test_suffix", result)
	})

	t.Run("nil函数返回默认值", func(t *testing.T) {
		val := "test"
		str := &val
		result := IfNil(str, "default", nil)
		assert.Equal(t, "default", result)
	})

	t.Run("nil切片处理", func(t *testing.T) {
		var slice []int
		result := IfNil(slice, 0, func(s []int) int { return len(s) })
		assert.Equal(t, 0, result)
	})

	t.Run("非nil空切片处理", func(t *testing.T) {
		slice := []int{}
		result := IfNil(slice, 100, func(s []int) int { return len(s) })
		assert.Equal(t, 0, result)
	})
}

func TestAt(t *testing.T) {
	t.Run("获取有效索引的值", func(t *testing.T) {
		result := At(0, 1, 2, 3, 4, 5)
		assert.Equal(t, 1, result)

		result = At(2, 1, 2, 3, 4, 5)
		assert.Equal(t, 3, result)

		result = At(4, 1, 2, 3, 4, 5)
		assert.Equal(t, 5, result)
	})

	t.Run("获取字符串切片的值", func(t *testing.T) {
		result := At(1, "a", "b", "c")
		assert.Equal(t, "b", result)
	})

	t.Run("索引超出范围应该panic", func(t *testing.T) {
		assert.Panics(t, func() {
			At(5, 1, 2, 3, 4, 5)
		})

		assert.Panics(t, func() {
			At(-1, 1, 2, 3)
		})
	})

	t.Run("布尔值切片测试", func(t *testing.T) {
		result := At(1, true, false, true)
		assert.False(t, result)
	})
}

package gox

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandStrN(t *testing.T) {
	t.Run("生成指定长度的随机字符串", func(t *testing.T) {
		length := 10
		result := RandStrN(length)
		assert.Equal(t, length, len(result))
	})

	t.Run("生成不同长度的随机字符串", func(t *testing.T) {
		assert.Equal(t, 5, len(RandStrN(5)))
		assert.Equal(t, 15, len(RandStrN(15)))
		assert.Equal(t, 20, len(RandStrN(20)))
	})

	t.Run("生成的字符串只包含指定字符集", func(t *testing.T) {
		result := RandStrN(100) // 生成足够长的字符串以增加测试覆盖率
		for _, char := range result {
			assert.True(t, strings.ContainsRune(letterBytes, char), "字符 %c 不在允许的字符集中", char)
		}
	})

	t.Run("生成多次字符串应该不同", func(t *testing.T) {
		// 由于是随机的，理论上有极小概率相同，但在实际测试中几乎不可能
		s1 := RandStrN(20)
		s2 := RandStrN(20)
		s3 := RandStrN(20)

		// 检查至少有两个不同的结果
		different := false
		if s1 != s2 || s2 != s3 || s1 != s3 {
			different = true
		}
		assert.True(t, different, "多次生成的随机字符串应该不同")
	})

	t.Run("生成空字符串", func(t *testing.T) {
		result := RandStrN(0)
		assert.Empty(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("生成较长的随机字符串", func(t *testing.T) {
		length := 1000
		result := RandStrN(length)
		assert.Equal(t, length, len(result))

		// 检查是否包含大小写字母和数字
		hasLower := false
		hasUpper := false
		hasDigit := false

		for _, char := range result {
			if 'a' <= char && char <= 'z' {
				hasLower = true
			} else if 'A' <= char && char <= 'Z' {
				hasUpper = true
			} else if '0' <= char && char <= '9' {
				hasDigit = true
			}
		}
		assert.True(t, hasLower, "生成的随机字符串应包含小写字母")
		assert.True(t, hasUpper, "生成的随机字符串应包含大写字母")
		assert.True(t, hasDigit, "生成的随机字符串应包含数字")

		// 由于是随机的，不能保证每次都包含所有类型，但在1000长度下概率很高
		// 如果这个测试偶尔失败，可以考虑放宽条件或增加长度
	})
}

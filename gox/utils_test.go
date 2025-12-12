package gox

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOID(t *testing.T) {
	t.Run("生成有效的ObjectId", func(t *testing.T) {
		id := NewOID()
		assert.NotEmpty(t, id)

		// 验证生成的是有效的bson.ObjectId类型
		assert.IsType(t, ObjectID{}, id)
	})

	t.Run("多次生成的ObjectId应该不同", func(t *testing.T) {
		id1 := NewOID()
		id2 := NewOID()
		assert.NotEqual(t, id1, id2, "多次生成的ObjectId应该不同")
	})

	t.Run("生成的ObjectId应该有效", func(t *testing.T) {
		id := NewOID()
		hex := id.Hex()
		// ObjectId的十六进制表示应该是24个字符长
		assert.Equal(t, 24, len(hex))

		// 验证可以从十六进制字符串解析回ObjectId
		parsedId, err := AsOID(hex)
		assert.NoError(t, err)
		assert.Equal(t, id, parsedId)
	})
}

func TestNewOIDHex(t *testing.T) {
	t.Run("生成有效的ObjectId十六进制字符串", func(t *testing.T) {
		hex := NewOIDHex()
		assert.NotEmpty(t, hex)

		// 验证长度为24个字符
		assert.Equal(t, 24, len(hex))

		// 验证只包含有效的十六进制字符（小写）
		for _, char := range hex {
			assert.True(t, (char >= '0' && char <= '9') || (char >= 'a' && char <= 'f'),
				"字符 %c 不是有效的十六进制字符", char)
		}
	})

	t.Run("多次生成的ObjectId十六进制字符串应该不同", func(t *testing.T) {
		hex1 := NewOIDHex()
		hex2 := NewOIDHex()
		hex3 := NewOIDHex()

		// 检查至少有两个不同的结果
		different := false
		if hex1 != hex2 || hex2 != hex3 || hex1 != hex3 {
			different = true
		}
		assert.True(t, different, "多次生成的ObjectId十六进制字符串应该不同")
	})

	t.Run("生成的十六进制字符串应该可以解析为ObjectId", func(t *testing.T) {
		hex := NewOIDHex()

		// 尝试从十六进制字符串解析ObjectId
		id, err := AsOID(hex)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		// 验证解析后的ObjectId的十六进制表示与原字符串相同
		assert.Equal(t, hex, id.Hex())
	})

	t.Run("生成的十六进制字符串应该是小写的", func(t *testing.T) {
		hex := NewOIDHex()
		assert.Equal(t, hex, strings.ToLower(hex))
	})
}

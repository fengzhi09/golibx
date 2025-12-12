package gox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqArr(t *testing.T) {
	t.Run("单个数组去重", func(t *testing.T) {
		arr := []int{1, 2, 2, 3, 3, 3}
		result := ArrUniq(arr)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("多个数组去重合并", func(t *testing.T) {
		arr1 := []int{1, 2, 3}
		arr2 := []int{2, 3, 4}
		arr3 := []int{3, 4, 5}
		result := ArrUniq(arr1, arr2, arr3)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("空数组去重", func(t *testing.T) {
		result := ArrUniq([]int{})
		assert.Empty(t, result)
	})

	t.Run("字符串数组去重", func(t *testing.T) {
		arr := []string{"a", "b", "b", "c"}
		result := ArrUniq(arr)
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("布尔数组去重", func(t *testing.T) {
		arr := []bool{true, false, true, false}
		result := ArrUniq(arr)
		assert.Equal(t, []bool{true, false}, result)
	})
}

func TestIn(t *testing.T) {
	t.Run("元素在数组中", func(t *testing.T) {
		assert.True(t, In(2, 1, 2, 3, 4, 5))
	})

	t.Run("元素不在数组中", func(t *testing.T) {
		assert.False(t, In(6, 1, 2, 3, 4, 5))
	})

	t.Run("空数组查询", func(t *testing.T) {
		assert.False(t, In(1))
	})

	t.Run("字符串元素查询", func(t *testing.T) {
		assert.True(t, In("b", "a", "b", "c"))
		assert.False(t, In("d", "a", "b", "c"))
	})
}

func TestIndexOf(t *testing.T) {
	t.Run("元素在数组中返回正确索引", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		assert.Equal(t, 2, IndexOf(arr, 3))
	})

	t.Run("元素不在数组中返回-1", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		assert.Equal(t, -1, IndexOf(arr, 6))
	})

	t.Run("空数组查询返回-1", func(t *testing.T) {
		assert.Equal(t, -1, IndexOf([]int{}, 1))
	})

	t.Run("重复元素返回第一个索引", func(t *testing.T) {
		arr := []int{1, 2, 3, 2, 5}
		assert.Equal(t, 1, IndexOf(arr, 2))
	})

	t.Run("字符串数组索引查询", func(t *testing.T) {
		arr := []string{"a", "b", "c", "d"}
		assert.Equal(t, 2, IndexOf(arr, "c"))
	})
}

func TestArrStarts(t *testing.T) {
	t.Run("数组以子数组开头", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{1, 2, 3}
		assert.True(t, ArrStarts(a, b))
	})

	t.Run("数组不以子数组开头", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{2, 3, 4}
		assert.False(t, ArrStarts(a, b))
	})

	t.Run("子数组长度大于数组", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{1, 2, 3, 4, 5}
		assert.False(t, ArrStarts(a, b))
	})

	t.Run("空数组检查", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{}
		assert.True(t, ArrStarts(a, b)) // 空数组是任何数组的前缀
	})

	t.Run("字符串数组前缀检查", func(t *testing.T) {
		a := []string{"a", "b", "c", "d"}
		b := []string{"a", "b"}
		assert.True(t, ArrStarts(a, b))
	})
}

func TestArrEnds(t *testing.T) {
	t.Run("数组以子数组结尾", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{3, 4, 5}
		assert.True(t, ArrEnds(a, b))
	})

	t.Run("数组不以子数组结尾", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{1, 2, 3}
		assert.False(t, ArrEnds(a, b))
	})

	t.Run("子数组长度大于数组", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{1, 2, 3, 4, 5}
		assert.False(t, ArrEnds(a, b))
	})

	t.Run("空数组检查", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{}
		assert.True(t, ArrEnds(a, b)) // 空数组是任何数组的后缀
	})

	t.Run("字符串数组后缀检查", func(t *testing.T) {
		a := []string{"a", "b", "c", "d"}
		b := []string{"c", "d"}
		assert.True(t, ArrEnds(a, b))
	})
}

func TestArrIn(t *testing.T) {
	t.Run("子数组在数组中", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{2, 3, 4}
		assert.True(t, ArrIn(a, b))
	})

	t.Run("子数组不在数组中", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{2, 4, 5}
		assert.False(t, ArrIn(a, b))
	})

	t.Run("子数组长度大于数组", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{1, 2, 3, 4, 5}
		assert.False(t, ArrIn(a, b))
	})

	t.Run("空子数组", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{}
		assert.True(t, ArrIn(a, b)) // 空数组是任何数组的子数组
	})
}

func TestArrIndexOf(t *testing.T) {
	t.Run("子数组在数组中返回正确索引", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{2, 3, 4}
		assert.Equal(t, 1, ArrIndexOf(a, b))
	})

	t.Run("子数组在数组开头", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{1, 2, 3}
		assert.Equal(t, 0, ArrIndexOf(a, b))
	})

	t.Run("子数组在数组结尾", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{3, 4, 5}
		assert.Equal(t, 2, ArrIndexOf(a, b))
	})

	t.Run("子数组不在数组中返回-1", func(t *testing.T) {
		a := []int{1, 2, 3, 4, 5}
		b := []int{2, 4, 5}
		assert.Equal(t, -1, ArrIndexOf(a, b))
	})

	t.Run("子数组长度大于数组返回-1", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{1, 2, 3, 4, 5}
		assert.Equal(t, -1, ArrIndexOf(a, b))
	})

	t.Run("空子数组返回0", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{}
		assert.Equal(t, 0, ArrIndexOf(a, b))
	})

	t.Run("字符串数组子数组索引", func(t *testing.T) {
		a := []string{"a", "b", "c", "d", "e"}
		b := []string{"b", "c", "d"}
		assert.Equal(t, 1, ArrIndexOf(a, b))
	})
}

func TestWordIn(t *testing.T) {
	t.Run("字符串包含单词返回索引", func(t *testing.T) {
		content := "hello world this is a test"
		words := []string{"world", "test", "hello"}
		assert.Equal(t, 0, WordIn(content, words)) // "world"在words数组中的索引是0
	})

	t.Run("字符串包含多个单词返回第一个匹配的索引", func(t *testing.T) {
		content := "hello world test"
		words := []string{"test", "world", "hello"}
		assert.Equal(t, 0, WordIn(content, words)) // "test"是第一个匹配的单词，在words数组中的索引是0
	})

	t.Run("字符串不包含任何单词返回-1", func(t *testing.T) {
		content := "hello world"
		words := []string{"test", "example"}
		assert.Equal(t, -1, WordIn(content, words))
	})

	t.Run("空单词列表返回-1", func(t *testing.T) {
		content := "hello world"
		words := []string{}
		assert.Equal(t, -1, WordIn(content, words))
	})

	t.Run("子字符串匹配", func(t *testing.T) {
		content := "abcdefg"
		words := []string{"bcd", "def"}
		assert.Equal(t, 0, WordIn(content, words)) // "bcd"在words数组中的索引是0
	})
}

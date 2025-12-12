package gox

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试SortedMap的基本功能
func TestSortedMapBasic(t *testing.T) {
	// 创建一个按值排序的SortedMap
	sm := NewSortedMap(func(a, b *SortNode[int]) int {
		return b.Data - a.Data
	})

	// 测试Set和Get
	sm.Set("one", 1)
	sm.Set("two", 2)
	sm.Set("zero", 0)

	val, exists := sm.Get("one")
	assert.True(t, exists)
	assert.Equal(t, 1, val)

	val, exists = sm.Get("unknown")
	assert.False(t, exists)
	assert.Equal(t, 0, val) // 默认零值

	// 测试排序
	keys := sm.GetKeys()
	assert.Equal(t, []string{"zero", "one", "two"}, keys)

	// 测试修改现有值
	sm.Set("one", 3)
	keys = sm.GetKeys()
	assert.Equal(t, []string{"zero", "two", "one"}, keys) // 重新排序

	// 测试删除
	sm.Del("zero")
	keys = sm.GetKeys()
	assert.Equal(t, []string{"two", "one"}, keys)

	// 测试删除不存在的键
	sm.Del("unknown")
	keys = sm.GetKeys()
	assert.Equal(t, []string{"two", "one"}, keys) // 不变

	// 测试批量删除
	sm.Del("one", "two")
	keys = sm.GetKeys()
	assert.Empty(t, keys)
}

// 测试SortedMap的基本功能
func TestSortedMapOrder(t *testing.T) {
	// 创建一个按值排序的SortedMap
	sm := NewSortedMap(func(a, b *SortNode[int]) int {
		return a.Data - b.Data
	})

	// 测试Set和Get
	sm.Set("one", 1)
	sm.Set("two", 2)
	sm.Set("zero", 0)

	strs := []string{"two", "one", "zero"}
	fmt.Println("v-dec", strs)
	assert.Equal(t, strs, sm.GetKeys())
	sort.SliceStable(strs, func(i, j int) bool { return strings.Compare(strs[i], strs[j]) > 0 })
	fmt.Println("k-dec", strs)
	assert.Equal(t, strs, sm.SortOnly(func(a, b *SortNode[int]) int { return strings.Compare(a.Key, b.Key) }))

	sort.SliceStable(strs, func(i, j int) bool { return strings.Compare(strs[i], strs[j]) < 0 })
	fmt.Println("k-inc", strs)
	assert.Equal(t, strs, sm.SortOnly(func(a, b *SortNode[int]) int { return strings.Compare(a.Key, b.Key) * -1 }))
}

// 测试LRUCache的功能
func TestLRUCache(t *testing.T) {
	// 创建容量为2的LRU缓存
	sleep := func() {
		time.Sleep(time.Millisecond * 1)
	}
	// 测试Set和Get
	t.Run("get-key1+key2", func(tt *testing.T) {
		cache := NewLRUMap[string](2)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		sleep()
		assert.Nil(tt, ArrDiff([]string{"key2", "key1"}, cache.GetKeys()))
		assert.Nil(tt, MapDiff(
			map[string]string{"key1": "value1", "key2": "value2"}, cache.Items()))
	})

	t.Run("get-key2", func(tt *testing.T) {
		cache := NewLRUMap[string](2)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		_, _ = cache.Get("key2")
		sleep()
		assert.Nil(tt, ArrDiff([]string{"key2", "key1"}, cache.GetKeys()))
		assert.Nil(tt, MapDiff(
			map[string]string{"key1": "value1", "key2": "value2"}, cache.Items()))
	})

	t.Run("get-key1", func(tt *testing.T) {
		cache := NewLRUMap[string](2)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		_, _ = cache.Get("key2")
		_, _ = cache.Get("key1")
		sleep()
		assert.Nil(tt, ArrDiff([]string{"key1", "key2"}, cache.GetKeys()))
		assert.Nil(tt, MapDiff(
			map[string]string{"key1": "value1", "key2": "value2"}, cache.Items()))
	})

	t.Run("get-key3-not-exist=>keys", func(tt *testing.T) {
		cache := NewLRUMap[string](2)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		_, _ = cache.Get("key2")
		_, _ = cache.Get("key1")
		_, _ = cache.Get("key3")
		sleep()
		assert.Nil(tt, ArrDiff([]string{"key1", "key2"}, cache.GetKeys()))
		assert.Nil(tt, MapDiff(
			map[string]string{"key1": "value1", "key2": "value2"}, cache.Items()))
	})

	t.Run("set-key3", func(tt *testing.T) {
		cache := NewLRUMap[string](2)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		_, _ = cache.Get("key2")
		_, _ = cache.Get("key1")
		_, _ = cache.Get("key3")
		sleep()
		cache.Set("key3", "value3")
		assert.Nil(tt, ArrDiff([]string{"key3", "key1"}, cache.GetKeys()))
		assert.Nil(tt, MapDiff(
			map[string]string{"key3": "value3", "key1": "value1"}, cache.Items()))
	})

	t.Run("edit-key1", func(tt *testing.T) {
		cache := NewLRUMap[string](2)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		_, _ = cache.Get("key2")
		_, _ = cache.Get("key1")
		_, _ = cache.Get("key3")
		cache.Set("key3", "value3")
		cache.Set("key1", "updated1")
		_, _ = cache.Get("key1")
		sleep()
		assert.Nil(tt, ArrDiff([]string{"key1", "key3"}, cache.GetKeys()))
		assert.Nil(tt, MapDiff(
			map[string]string{"key3": "value3", "key1": "updated1"}, cache.Items()))
	})

	t.Run("set-key4", func(tt *testing.T) {
		cache := NewLRUMap[string](2)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		_, _ = cache.Get("key2")
		_, _ = cache.Get("key1")
		_, _ = cache.Get("key3")
		cache.Set("key3", "value3")
		cache.Set("key1", "updated1")
		_, _ = cache.Get("key1")
		_, _ = cache.Get("key3")
		cache.Set("key4", "value4")
		sleep()
		assert.Equal(tt, []string{"key4", "key3"}, cache.GetKeys())
		assert.Nil(tt, MapDiff(
			map[string]string{"key4": "value4", "key3": "value3"}, cache.Items()))
	})
}

// 测试并发安全性（简单测试）
func TestSortedMapConcurrent(t *testing.T) {
	// 创建一个按值排序的SortedMap
	sm := NewSortedMap[int](func(a, b *SortNode[int]) int {
		return a.Data - b.Data
	})

	// 并发写入
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			sm.Set("key"+string(rune(n)), n)
		}(i)
	}
	wg.Wait()

	// 并发读取
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sm.GetKeys()
		}()
	}
	wg.Wait()

	// 只要没有panic，测试就算通过
	assert.True(t, true)
}

// 测试LRUCache在边界情况下的行为
func TestLRUCacheEdgeCases(t *testing.T) {
	// 测试容量为0的情况
	cache := NewLRUMap[string](0)
	cache.Set("key", "value")
	cache.DropBySize()
	val, exists := cache.Get("key")
	assert.False(t, exists)
	assert.Empty(t, val)

	// 测试容量为1的情况
	cache = NewLRUMap[string](1)
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	val, exists = cache.Get("key1")
	assert.False(t, exists)
	val, exists = cache.Get("key2")
	assert.True(t, exists)
	assert.Equal(t, "value2", val)
}

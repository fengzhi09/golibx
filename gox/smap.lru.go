package gox

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

type lruItem[T any] struct {
	Data     T
	UpdateAt int64
	VisitAt  int64
	ActIdx   int64
}

type LRUMap[T any] struct {
	sync.Mutex
	sm     *SortedMap[*lruItem[T]]
	size   int
	actIdx int64
}

func NewLRUMap[T any](size int) *LRUMap[T] {
	// 创建一个按访问时间(纳秒级)降序排序的SortedMap（最新访问的在前）
	sm := NewSortedMap(
		func(a, b *SortNode[*lruItem[T]]) int {
			if a.Data.VisitAt == b.Data.VisitAt {
				return int(a.Data.ActIdx - b.Data.ActIdx)
			}
			return int(a.Data.VisitAt - b.Data.VisitAt)
		},
	)
	return &LRUMap[T]{sm: sm, size: size, actIdx: 0}
}

func (lc *LRUMap[T]) Data() *SortedMap[*lruItem[T]] {
	return lc.sm
}

// GetKeys 获取当前缓存的所有键
func (lc *LRUMap[T]) GetKeys() []string {
	return lc.sm.GetKeys()
}

// Has 检测缓存项是否存在, 不会影响lru计算
func (lc *LRUMap[T]) Has(key string) bool {
	// 获取项
	_, exists := lc.sm.Get(key)
	return exists
}

// Get 获取缓存项, 会影响lru计算
func (lc *LRUMap[T]) Get(key string) (T, bool) {
	// before := lc.metas()
	// defer func() { fmt.Printf("Get[%v] before:%v => after: %v\n\n", key, before, lc.metas()) }()
	item, exists := lc.sm.Get(key)
	if exists {
		item.VisitAt = time.Now().UnixNano()
		item.ActIdx = lc.actInc()
		lc.sm.ensureSort()
		return item.Data, exists
	}
	return *new(T), exists
}

// Items 获取所有缓存项, 不会影响lru计算
func (lc *LRUMap[T]) Items() map[string]T {
	items := lc.sm.GetItems()
	result := make(map[string]T, len(items))
	for key, item := range items {
		result[key] = item.Data
	}
	return result
}

// Set 设置缓存项, 会影响lru计算
func (lc *LRUMap[T]) Set(key string, data T) {
	// before := lc.metas()
	// defer func() { fmt.Printf("Set[%v] before:%v => after: %v\n\n", key, before, lc.metas()) }()
	currentTime := time.Now().UnixNano() // 纳秒级时间戳, 毫秒级 会导致排序不准确
	item := &lruItem[T]{
		Data:     data,
		UpdateAt: currentTime,
		VisitAt:  currentTime,
		ActIdx:   lc.actInc(),
	}
	lc.sm.Set(key, item)
	// 自动清理超出容量的项
	lc.DropBySize()
}

// VisitTsMap 获取所有访问时间(纳秒级)
func (lc *LRUMap[T]) VisitTsMap() map[string]int64 {
	items := lc.sm.GetItems()
	result := make(map[string]int64, len(items))
	for key, item := range items {
		result[key] = item.VisitAt
	}
	return result
}

// EditTsMap 获取所有编辑时间(纳秒级)
func (lc *LRUMap[T]) EditTsMap() map[string]int64 {
	items := lc.sm.GetItems()
	result := make(map[string]int64, len(items))
	for key, item := range items {
		result[key] = item.UpdateAt
	}
	return result
}

// DropBySize 手动清理超出容量的项
func (lc *LRUMap[T]) DropBySize() {
	var toDel []string
	keys := lc.sm.SortKeys()

	if len(lc.sm.keys) <= lc.size {
		return
	}
	before, caller := lc.metas(), CallerDep(2)
	defer func() { fmt.Printf("Drop[%v] before:%v => after: %v\n\n", caller, before, lc.metas()) }()

	// 删除最久未访问的项
	// 由于按VisitAt降序排列（最新访问的在前），所以前几个元素是最新访问的，我们只需要删除后几个元素即可
	toDel = keys[lc.size:]
	lc.sm.Del(toDel...)
}

// DropByUpdateAt 按更新时间(纳秒级)删除
func (lc *LRUMap[T]) DropByUpdateAt(nanoTs int64) {
	lc.dropByTime(nanoTs, false)
}

// DropByVisit 按访问时间(纳秒级)删除
func (lc *LRUMap[T]) DropByVisit(nanoTs int64) {
	lc.dropByTime(nanoTs, true)
}

func (lc *LRUMap[T]) Del(keys ...string) {
	lc.sm.Del(keys...)
	lc.sm.ensureSort()
}

// dropByTime 内部清理方法(无锁)，按时间(纳秒级)删除
func (lc *LRUMap[T]) dropByTime(nanoTs int64, byVisit bool) {
	items := lc.sm.GetItems()
	toDel := make([]string, 0)
	for key, item := range items {
		itemTs := IfElse(byVisit, item.VisitAt, item.UpdateAt).(int64)
		if byVisit && itemTs <= 0 || itemTs > nanoTs {
			continue
		}
		toDel = append(toDel, key)
	}
	lc.sm.Del(toDel...)
}

func (lc *LRUMap[T]) actInc() int64 {
	lc.Lock()
	defer lc.Unlock()
	lc.actIdx++
	if lc.actIdx > int64(lc.size*1000) {
		lc.clear()
	}
	return lc.actIdx
}

// metas 获取所有项的元数据
func (lc *LRUMap[T]) metas() string {
	lc.Lock()
	defer lc.Unlock()
	keys := lc.sm.SortKeys()
	result := make([]string, len(keys))
	for _, key := range keys {
		item := lc.sm.data[key]
		result = append(result,
			fmt.Sprintf("K=%v,I=%v,V=%v,U=%v,D=%v", key, item.ActIdx, item.VisitAt, item.UpdateAt, item.Data))
	}
	return strings.Join(result, "\n")
}

func (lc *LRUMap[T]) Clear() {
	lc.Lock()
	defer lc.Unlock()
	lc.clear()
}
func (lc *LRUMap[T]) clear() {
	lc.sm.Clear()
	lc.actIdx = 0
	// 强制GC
	go func() {
		time.Sleep(time.Millisecond * 100)
		runtime.GC()
	}()
}

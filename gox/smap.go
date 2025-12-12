package gox

import (
	"sort"
	"sync"
)

type SortedMap[T any] struct {
	sync.RWMutex
	keys   []string
	data   map[string]T
	sorter func(a, b *SortNode[T]) int
}

type SortNode[T any] struct {
	Idx  int
	Key  string
	Data T
}

func NewSortedMap[T any](sorter func(a, b *SortNode[T]) int) *SortedMap[T] {
	sm := &SortedMap[T]{
		keys: make([]string, 0),
		data: make(map[string]T),
	}
	if sorter == nil {
		sorter = func(a, b *SortNode[T]) int {
			return a.Idx - b.Idx
		}
	}
	sm.sorter = sorter
	return sm
}

func (sm *SortedMap[T]) Has(key string) bool {
	sm.RLock()
	defer sm.RUnlock()
	_, exists := sm.data[key]
	return exists
}

func (sm *SortedMap[T]) Len() int {
	sm.RLock()
	defer sm.RUnlock()
	return len(sm.keys)
}

func (sm *SortedMap[T]) Clear() {
	sm.Lock()
	defer sm.Unlock()
	sm.keys = make([]string, 0)
	sm.data = make(map[string]T)
}

func (sm *SortedMap[T]) NodeAt(i int) *SortNode[T] {
	return &SortNode[T]{Idx: i, Key: sm.keys[i], Data: sm.data[sm.keys[i]]}
}

// SortOnly 仅按当前sorter排序，不改变数据
func (sm *SortedMap[T]) SortOnly(sorter func(a, b *SortNode[T]) int) []string {
	if sorter == nil {
		sorter = sm.sorter
	}
	// fmt.Printf("caller:%v SortOnly before:%v\n", CallerSkip("smap.go"), sm.keys)

	// 创建一个临时节点数组并排序
	nodes := make([]*SortNode[T], len(sm.keys))
	for i, key := range sm.keys {
		nodes[i] = &SortNode[T]{Idx: i, Key: key, Data: sm.data[key]}
	}

	// 对节点数组排序
	sort.SliceStable(nodes, func(i, j int) bool {
		return sorter(nodes[i], nodes[j]) > 0
	})

	// 从排序后的节点中提取键
	sortedKeys := make([]string, len(nodes))
	for i, node := range nodes {
		sortedKeys[i] = node.Key
	}

	// fmt.Printf("caller:%v SortOnly after:%v\n", CallerSkip("smap.go"), sortedKeys)
	return sortedKeys
}

// SetSorter 设置排序函数
func (sm *SortedMap[T]) SetSorter(sorter func(a, b *SortNode[T]) int) {
	sm.Lock()
	defer sm.Unlock()
	sm.sorter = sorter
	sm.ensureSort()
}

func (sm *SortedMap[T]) At(i int) (string, T) {
	sm.RLock()
	defer sm.RUnlock()
	return sm.keys[i], sm.data[sm.keys[i]]
}

func (sm *SortedMap[T]) GetOrSet(key string, defaultValue T) (string, T) {
	sm.RLock()
	defer sm.RUnlock()
	if item, exists := sm.data[key]; exists {
		return key, item
	}
	return key, defaultValue
}

func (sm *SortedMap[T]) Get(key string) (T, bool) {
	sm.RLock()
	defer sm.RUnlock()
	item, exists := sm.data[key]
	return item, exists
}

func (sm *SortedMap[T]) Set(key string, item T) {
	sm.Lock()
	defer sm.Unlock()
	if _, exists := sm.data[key]; !exists {
		sm.keys = append(sm.keys, key)
	}
	sm.data[key] = item
	sm.ensureSort()
}

func (sm *SortedMap[T]) Del(toDel ...string) {
	sm.Lock()
	defer sm.Unlock()
	if len(toDel) == 0 {
		return
	}
	dels, keeps, keys := map[string]bool{}, map[string]T{}, []string{}
	for _, key := range toDel {
		dels[key] = true
	}
	for _, key := range sm.keys {
		if ok, hit := dels[key]; hit || ok {
			continue
		}
		keys = append(keys, key)
		keeps[key] = sm.data[key]
	}
	sm.keys = keys
	sm.data = keeps
	sm.ensureSort()
}

func (sm *SortedMap[T]) SortKeys() []string {
	sm.RLock()
	defer sm.RUnlock()
	sm.ensureSort()
	return sm.keys
}

func (sm *SortedMap[T]) GetKeys() []string {
	sm.Lock()
	defer sm.Unlock()
	return sm.keys
}

func (sm *SortedMap[T]) ensureSort() {
	sm.keys = sm.SortOnly(sm.sorter)
}

func (sm *SortedMap[T]) GetItems() map[string]T {
	sm.RLock()
	defer sm.RUnlock()
	sm.ensureSort()
	return sm.data
}

// stringpool.go - 完整替换
package fastlog

import (
	"container/list"
	"sync"
)

type StringPool struct {
	pool    map[string]*poolEntry
	lruList *list.List
	maxSize int
	mu      sync.RWMutex
}

type poolEntry struct {
	value   *string
	element *list.Element
}

func NewStringPool(maxSize int) *StringPool {
	if maxSize <= 0 {
		maxSize = 10000 // 默认最大容量
	}
	return &StringPool{
		pool:    make(map[string]*poolEntry),
		lruList: list.New(),
		maxSize: maxSize,
	}
}

func (sp *StringPool) Intern(s string) *string {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	// 如果已存在，移到前面
	if entry, exists := sp.pool[s]; exists {
		sp.lruList.MoveToFront(entry.element)
		return entry.value
	}

	// 如果达到最大容量，删除最久未使用的
	if len(sp.pool) >= sp.maxSize {
		sp.evictLRU()
	}

	// 添加新字符串
	interned := &s
	element := sp.lruList.PushFront(s)
	sp.pool[s] = &poolEntry{
		value:   interned,
		element: element,
	}

	return interned
}

func (sp *StringPool) evictLRU() {
	if sp.lruList.Len() == 0 {
		return
	}

	// 删除最后一个元素
	element := sp.lruList.Back()
	if element != nil {
		key := element.Value.(string)
		delete(sp.pool, key)
		sp.lruList.Remove(element)
	}
}

func (sp *StringPool) Size() int {
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	return len(sp.pool)
}

func (sp *StringPool) Clear() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.pool = make(map[string]*poolEntry)
	sp.lruList.Init()
}

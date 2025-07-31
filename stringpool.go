package fastlog

import (
	"sync"
	"sync/atomic"
)

// StringPool 字符串池
type StringPool struct {
	mu    sync.RWMutex       // 读写锁
	pool  map[string]*string // 字符串池
	stats StringPoolStats    // 统计信息
}

// StringPoolStats 字符串池统计信息
type StringPoolStats struct {
	Hits   int64 // 命中次数
	Misses int64 // 未命中次数
	Size   int64 // 池大小
}

// 获取或创建字符串
func (sp *StringPool) Intern(s string) *string {
	// 先尝试读锁
	sp.mu.RLock()
	if interned, exists := sp.pool[s]; exists {
		atomic.AddInt64(&sp.stats.Hits, 1)
		sp.mu.RUnlock()
		return interned
	}
	sp.mu.RUnlock()

	// 写锁创建新字符串
	sp.mu.Lock()
	defer sp.mu.Unlock()

	// 双重检查
	if interned, exists := sp.pool[s]; exists {
		atomic.AddInt64(&sp.stats.Hits, 1)
		return interned
	}

	// 创建新的字符串指针
	interned := &s
	sp.pool[s] = interned
	atomic.AddInt64(&sp.stats.Misses, 1)
	atomic.AddInt64(&sp.stats.Size, 1)

	return interned
}

// NewStringPool 创建一个新的字符串池
func NewStringPool(initialCapacity int) *StringPool {
	return &StringPool{
		pool:  make(map[string]*string, initialCapacity),
		stats: StringPoolStats{},
	}
}

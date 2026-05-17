package fastlog

import (
	"sync/atomic"
	"time"
)

const (
	samplerLevels  = 6    // 日志级别数: DEBUG ~ PANIC
	samplerBuckets = 4096 // 每个级别下的桶数
)

// samplerCounter 采样计数器
type samplerCounter struct {
	resetAt atomic.Int64  // 窗口过期时间 (UnixNano)
	counter atomic.Uint64 // 当前窗口内的计数
}

// Sampler 日志采样器
//
// 使用固定桶 + atomic 计数器实现, 无锁设计。
// 相同 level 和 message 的日志会被哈希到同一个桶, 在时间窗口内按规则放行或抑制。
type Sampler struct {
	tick       time.Duration                                 // 时间窗口
	initial    int                                           // 窗口内前 N 条放行
	thereafter int                                           // 之后每 M 条放行 1 条
	counters   [samplerLevels][samplerBuckets]samplerCounter // 每个级别下的每个桶的计数器
}

// NewSampler 创建日志采样器
//
// 参数:
//   - tick: 时间窗口, 如 10*time.Second。如果 <= 0, 默认使用 1 秒
//   - initial: 窗口内前 N 条放行。如果 < 0, 默认使用 1
//   - thereafter: 之后每 M 条放行 1 条, 0 表示不再放行。如果 < 0, 默认使用 10
//
// 示例:
//
//	// 每 10 秒, 前 3 条放行, 之后每 10 条放行 1 条
//	sampler := fastlog.NewSampler(10*time.Second, 3, 10)
func NewSampler(tick time.Duration, initial, thereafter int) *Sampler {
	if tick <= 0 {
		tick = time.Second
	}
	if initial < 0 {
		initial = 1
	}
	if thereafter < 0 {
		thereafter = 10
	}
	return &Sampler{
		tick:       tick,
		initial:    initial,
		thereafter: thereafter,
	}
}

// DefaultSampler 创建默认日志采样器
//
// 默认参数: 窗口 1 秒, 前 3 条放行, 之后每 10 条放行 1 条。
// 适合大多数场景直接使用, 无需额外配置。
//
// 示例:
//
//	sampler := fastlog.DefaultSampler()
func DefaultSampler() *Sampler {
	return NewSampler(time.Second, 3, 10)
}

// Allow 判断是否放行这条日志
//
// 参数:
//   - level: 日志级别
//   - msg: 日志消息
//
// 返回:
//   - bool: true 放行, false 抑制
func (s *Sampler) Allow(level Level, msg string) bool {
	// level 转索引 (DEBUG=1 → 0, INFO=2 → 1, ...)
	i := int(level) - 1
	if i < 0 || i >= samplerLevels {
		return true
	}

	// message 哈希到桶
	j := fnv32a(msg) % samplerBuckets
	c := &s.counters[i][j]

	now := time.Now()
	tn := now.UnixNano()

	// 如果窗口已过期, 重置计数器
	if tn > c.resetAt.Load() {
		c.counter.Store(1)
		c.resetAt.Store(tn + s.tick.Nanoseconds())
		return true
	}

	// 窗口内计数递增
	n := c.counter.Add(1)

	// 前 N 条放行
	if uint64(s.initial) >= n {
		return true
	}

	// 之后每 M 条放行 1 条 (thereafter=0 表示不再放行)
	if s.thereafter > 0 && (n-uint64(s.initial))%uint64(s.thereafter) == 0 {
		return true
	}

	return false
}

// fnv32a FNV-1a 哈希函数, 无内存分配
func fnv32a(s string) uint32 {
	const (
		offset32 = 2166136261
		prime32  = 16777619
	)
	hash := uint32(offset32)
	for i := 0; i < len(s); i++ {
		hash ^= uint32(s[i])
		hash *= prime32
	}
	return hash
}

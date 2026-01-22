package pool

import (
	"strings"
	"sync"
)

// 全局默认字符串构建器池实例
//
// 说明:
//   - 该实例用于在全局范围内管理字符串构建器对象，避免频繁创建和销毁对象导致的性能问题。
//   - 初始容量为256，最大回收容量为32KB。
var defStrPool = NewStrPool(256, 32*1024)

// GetStr 从默认字符串池获取默认容量的字符串构建器
//
// 返回值:
//   - *strings.Builder: 容量至少为默认容量的字符串构建器
func GetStr() *strings.Builder { return defStrPool.Get() }

// GetStrCap 从默认字符串池获取指定容量的字符串构建器
//
// 参数:
//   - cap: 字符串构建器初始容量
//
// 返回值:
//   - *strings.Builder: 容量至少为capacity的字符串构建器
func GetStrCap(cap int) *strings.Builder {
	return defStrPool.GetCap(cap)
}

// PutStr 将字符串构建器归还到默认字符串池
//
// 参数:
//   - buf: 要归还的字符串构建器
//
// 说明:
//   - 该函数将字符串构建器归还到对象池，以便后续复用。
func PutStr(buf *strings.Builder) { defStrPool.Put(buf) }

// WithStr 使用默认容量的字符串构建器执行函数，自动管理获取和归还
//
// 参数:
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 使用示例:
//
//	result := pool.WithStr(func(buf *strings.Builder) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func WithStr(fn func(*strings.Builder)) string { return defStrPool.With(fn) }

// WithStrCap 使用指定容量的字符串构建器执行函数，自动管理获取和归还
//
// 参数:
//   - cap: 字符串构建器初始容量
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 使用示例:
//
//	result := pool.WithStrCap(64, func(buf *strings.Builder) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func WithStrCap(cap int, fn func(*strings.Builder)) string {
	return defStrPool.WithCap(cap, fn)
}

// StrPool 字符串构建器对象池，支持自定义配置
type StrPool struct {
	pool   sync.Pool // 字符串构建器对象池
	maxCap int       // 最大回收构建器容量
	defCap int       // 默认构建器容量
}

// NewStrPool 创建新的字符串构建器对象池
//
// 参数:
//   - defCap: 默认字符串构建器容量
//   - maxCap: 最大回收构建器容量，超过此容量的构建器不会被回收
//
// 返回值:
//   - *StrPool: 字符串构建器对象池实例
func NewStrPool(defCap, maxCap int) *StrPool {
	if defCap <= 0 {
		defCap = 256 // 默认256字节
	}
	if maxCap <= 0 {
		maxCap = 32 * 1024 // 默认32KB
	}

	return &StrPool{
		maxCap: maxCap,
		defCap: defCap,
		pool: sync.Pool{
			New: func() any {
				return new(strings.Builder)
			},
		},
	}
}

// Get 获取默认容量的字符串构建器
//
// 返回:
//   - *strings.Builder: 容量至少为默认容量的字符串构建器
//
// 说明:
//   - 返回的字符串构建器已经重置为空状态，可以直接使用
//   - 底层容量可能大于默认容量，来自对象池的复用构建器
func (sp *StrPool) Get() *strings.Builder {
	return sp.GetCap(sp.defCap)
}

// GetCap 获取指定容量的字符串构建器
//
// 参数:
//   - cap: 需要的字符串构建器容量
//
// 返回:
//   - *strings.Builder: 容量至少为capacity的字符串构建器
//
// 说明:
//   - 返回的字符串构建器已经重置为空状态，可以直接使用
//   - 底层容量可能大于capacity，来自对象池的复用构建器
//   - 如果capacity <= 0, 返回默认容量的构建器
func (sp *StrPool) GetCap(cap int) *strings.Builder {
	if cap <= 0 {
		cap = sp.defCap
	}
	builder, ok := sp.pool.Get().(*strings.Builder)
	if !ok {
		panic("string pool: invalid builder type")
	}

	// 先Reset确保内容完全干净(注意: Reset()会将buf设为nil，容量变为0)
	builder.Reset()

	// Reset后重新分配所需容量
	// 由于Reset()后容量为0，直接Grow(cap)即可
	if cap > 0 {
		builder.Grow(cap)
	}

	return builder
}

// Put 归还字符串构建器到对象池
//
// 参数:
//   - b: 要归还的字符串构建器
//
// 说明:
//   - nil构建器不会被回收
//   - 容量不超过maxCap的构建器直接重置后归还
//   - 容量超过maxCap的构建器会创建一个新的小容量构建器进行归还（智能缩容）
func (sp *StrPool) Put(buf *strings.Builder) {
	if buf == nil || buf.Cap() > sp.maxCap {
		return // 为nil或容量过大不处理, 交给gc回收
	}
	sp.pool.Put(buf)
}

// With 使用默认容量的字符串构建器执行函数，自动管理获取和归还
//
// 参数:
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 说明:
//   - 自动从对象池获取默认容量的字符串构建器
//   - 执行用户提供的函数
//   - 获取构建的字符串结果
//   - 自动归还字符串构建器到对象池
//   - 即使函数发生panic也会正确归还资源
func (sp *StrPool) With(fn func(*strings.Builder)) string {
	buf := sp.Get()
	defer sp.Put(buf)
	fn(buf)
	return buf.String()
}

// WithCap 使用指定容量的字符串构建器执行函数，自动管理获取和归还
//
// 参数:
//   - cap: 字符串构建器初始容量
//   - fn: 使用字符串构建器的函数
//
// 返回值:
//   - string: 函数执行后构建的字符串结果
//
// 说明:
//   - 自动从对象池获取指定容量的字符串构建器
//   - 执行用户提供的函数
//   - 获取构建的字符串结果
//   - 自动归还字符串构建器到对象池
//   - 即使函数发生panic也会正确归还资源
func (sp *StrPool) WithCap(cap int, fn func(*strings.Builder)) string {
	buf := sp.GetCap(cap)
	defer sp.Put(buf)
	fn(buf)
	return buf.String()
}

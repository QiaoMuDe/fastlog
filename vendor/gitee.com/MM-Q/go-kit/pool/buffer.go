package pool

import (
	"bytes"
	"sync"
)

// 全局默认缓冲区池实例, 初始容量为256, 最大容量为32KB
var defBufPool = NewBufPool(256, 32*1024)

// GetBuf 从默认缓冲区池获取默认容量的字节缓冲区
//
// 返回值:
//   - *bytes.Buffer: 容量至少为默认容量的字节缓冲区
func GetBuf() *bytes.Buffer { return defBufPool.Get() }

// GetBufCap 从默认缓冲区池获取指定容量的字节缓冲区
//
// 参数:
//   - cap: 缓冲区初始容量
//
// 返回值:
//   - *bytes.Buffer: 容量至少为capacity的字节缓冲区
func GetBufCap(cap int) *bytes.Buffer { return defBufPool.GetCap(cap) }

// PutBuf 将字节缓冲区归还到默认缓冲区池
//
// 参数:
//   - buf: 要归还的字节缓冲区
//
// 说明:
//   - 该函数将字节缓冲区归还到对象池，以便后续复用。
func PutBuf(buf *bytes.Buffer) { defBufPool.Put(buf) }

// WithBuf 使用默认容量的字节缓冲区执行函数，自动管理获取和归还
//
// 参数:
//   - fn: 使用字节缓冲区的函数
//
// 返回值:
//   - []byte: 函数执行后缓冲区的字节数据副本
//
// 使用示例:
//
//	data := pool.WithBuf(func(buf *bytes.Buffer) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func WithBuf(fn func(*bytes.Buffer)) []byte { return defBufPool.With(fn) }

// WithBufCap 使用指定容量的字节缓冲区执行函数，自动管理获取和归还
//
// 参数:
//   - cap: 字节缓冲区初始容量
//   - fn: 使用字节缓冲区的函数
//
// 返回值:
//   - []byte: 函数执行后缓冲区的字节数据副本
//
// 使用示例:
//
//	data := pool.WithBufCap(1024, func(buf *bytes.Buffer) {
//	    buf.WriteString("Hello")
//	    buf.WriteByte(' ')
//	    buf.WriteString("World")
//	})
func WithBufCap(cap int, fn func(*bytes.Buffer)) []byte {
	return defBufPool.WithCap(cap, fn)
}

// BufPool 字节缓冲区对象池，支持自定义配置
type BufPool struct {
	pool   sync.Pool // 字节缓冲区对象池
	maxCap int       // 最大回收缓冲区容量
	defCap int       // 默认缓冲区容量
}

// NewBufPool 创建新的字节缓冲区对象池
//
// 参数:
//   - defCap: 默认字节缓冲区容量
//   - maxCap: 最大回收缓冲区容量，超过此容量的缓冲区不会被回收
//
// 返回值:
//   - *BufPool: 字节缓冲区对象池实例
func NewBufPool(defCap, maxCap int) *BufPool {
	if defCap <= 0 {
		defCap = 256 // 默认256字节
	}
	if maxCap <= 0 {
		maxCap = 32 * 1024 // 默认32KB
	}

	return &BufPool{
		maxCap: maxCap,
		defCap: defCap,
		pool: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

// Get 获取默认容量的字节缓冲区
//
// 返回:
//   - *bytes.Buffer: 容量至少为默认容量的字节缓冲区
//
// 说明:
//   - 返回的字节缓冲区已经重置为空状态，可以直接使用
//   - 底层容量可能大于默认容量，来自对象池的复用缓冲区
func (bp *BufPool) Get() *bytes.Buffer { return bp.GetCap(bp.defCap) }

// GetCap 获取指定容量的字节缓冲区
//
// 参数:
//   - cap: 需要的字节缓冲区容量
//
// 返回:
//   - *bytes.Buffer: 容量至少为capacity的字节缓冲区
//
// 说明:
//   - 返回的字节缓冲区已经重置为空状态，可以直接使用
//   - 底层容量可能大于capacity，来自对象池的复用缓冲区
//   - 如果capacity <= 0, 返回默认容量的缓冲区
func (bp *BufPool) GetCap(cap int) *bytes.Buffer {
	if cap <= 0 {
		cap = bp.defCap
	}

	buf, ok := bp.pool.Get().(*bytes.Buffer)
	if !ok {
		// 一旦触发说明代码契约被破坏了,直接panic比静默继续更安全
		panic("buffer pool: unexpected type")
	}

	// 如果当前容量不足，扩容到所需容量
	if buf.Cap() < cap {
		// 直接Grow(cap)，让Buffer内部决定如何分配足够的容量
		// 注意：Grow(n)确保还能写入n个字节，但实际分配可能更多
		buf.Grow(cap)
	}

	// 重置缓冲区状态
	buf.Reset()

	return buf
}

// Put 归还字节缓冲区到对象池
//
// 参数:
//   - buf: 要归还的字节缓冲区
func (bp *BufPool) Put(buf *bytes.Buffer) {
	if buf == nil || buf.Cap() > bp.maxCap {
		return // 为nil或容量过大不处理, 交给gc回收
	}
	buf.Reset()
	bp.pool.Put(buf)
}

// With 使用默认容量的字节缓冲区执行函数，自动管理获取和归还
//
// 参数:
//   - fn: 使用字节缓冲区的函数
//
// 返回值:
//   - []byte: 函数执行后缓冲区的字节数据副本
//
// 说明:
//   - 自动从对象池获取默认容量的字节缓冲区
//   - 执行用户提供的函数
//   - 获取缓冲区字节数据的副本
//   - 自动归还字节缓冲区到对象池
//   - 即使函数发生panic也会正确归还资源
func (bp *BufPool) With(fn func(*bytes.Buffer)) []byte {
	buf := bp.Get()
	defer bp.Put(buf)
	fn(buf)
	return append([]byte(nil), buf.Bytes()...) // 一次性拷贝
}

// WithCap 使用指定容量的字节缓冲区执行函数，自动管理获取和归还
//
// 参数:
//   - cap: 字节缓冲区初始容量
//   - fn: 使用字节缓冲区的函数
//
// 返回值:
//   - []byte: 函数执行后缓冲区的字节数据副本
//
// 说明:
//   - 自动从对象池获取指定容量的字节缓冲区
//   - 执行用户提供的函数
//   - 获取缓冲区字节数据的副本
//   - 自动归还字节缓冲区到对象池
//   - 即使函数发生panic也会正确归还资源
func (bp *BufPool) WithCap(cap int, fn func(*bytes.Buffer)) []byte {
	buf := bp.GetCap(cap)
	defer bp.Put(buf)
	fn(buf)
	return append([]byte(nil), buf.Bytes()...) // 一次性拷贝
}

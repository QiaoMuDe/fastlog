// Package pool 提供高性能对象池管理功能, 通过复用对象优化内存使用。
//
// 该包实现了基于 sync.Pool 的多种对象池, 用于减少频繁的内存分配和回收。
// 通过复用对象, 可以显著提升应用程序的性能, 特别是在高并发场景下。
//
// 主要功能：
//   - 字节切片对象池管理
//   - 动态容量对象获取
//   - 自动内存回收控制
//   - 防止内存泄漏的容量限制
//   - 支持多种对象类型的池化
//
// 性能优化：
//   - 使用 sync.Pool 减少 GC 压力
//   - 支持不同容量的对象需求
//   - 自动限制大对象回收
//   - 预热机制提升冷启动性能
package pool

import "sync"

// 全局默认对象池实例, 默认容量为256, 最大容量为32KB
var defBytePool = NewBytePool(256, 32*1024)

// GetByte 从默认字节池获取默认容量的缓冲区
//
// 返回值:
//   - []byte: 长度为默认容量, 容量至少为默认容量的缓冲区
func GetByte() []byte { return defBytePool.Get() }

// GetByteCap 从默认字节池获取指定容量的缓冲区
//
// 参数:
//   - size: 缓冲区容量
//
// 返回值:
//   - []byte: 长度为capacity, 容量至少为capacity的缓冲区
func GetByteCap(size int) []byte { return defBytePool.GetCap(size) }

// PutByte 将缓冲区归还到默认字节池
//
// 参数:
//   - buf: 要归还的缓冲区
//
// 说明:
//   - 该函数将缓冲区归还到对象池, 以便后续复用。
func PutByte(buf []byte) { defBytePool.Put(buf) }

// GetByteEmpty 从默认字节池获取空缓冲区
//
// 参数:
//   - size: 指定容量要求
//
// 返回值:
//   - []byte: 长度为0但容量至少为capacity的缓冲区切片
func GetByteEmpty(size int) []byte { return defBytePool.GetEmpty(size) }

// // WithByte 使用默认容量的字节缓冲区执行函数，自动管理获取和归还
// //
// // 参数:
// //   - fn: 使用字节缓冲区的函数
// //
// // 返回值:
// //   - []byte: 函数执行后缓冲区的字节数据副本
// //
// // 使用示例:
// //
// //	data := pool.WithByte(func(buf []byte) {
// //	    buf = append(buf, "Hello"...)
// //	    buf = append(buf, ' ')
// //	    buf = append(buf, "World"...)
// //	})
// func WithByte(fn func(*[]byte)) []byte { return defBytePool.With(fn) }

// // WithByteCap 使用指定容量的字节缓冲区执行函数，自动管理获取和归还
// //
// // 参数:
// //   - size: 字节缓冲区初始容量
// //   - fn: 使用字节缓冲区的函数
// //
// // 返回值:
// //   - []byte: 函数执行后缓冲区的字节数据副本
// //
// // 使用示例:
// //
// //	data := pool.WithByteCap(1024, func(buf []byte) {
// //	    buf = append(buf, "Hello"...)
// //	    buf = append(buf, ' ')
// //	    buf = append(buf, "World"...)
// //	})
// func WithByteCap(size int, fn func(*[]byte)) []byte {
// 	return defBytePool.WithCap(size, fn)
// }

// BytePool 字节切片对象池, 支持自定义配置
type BytePool struct {
	pool   sync.Pool // 缓冲区对象池
	maxCap int       // 最大回收缓冲区容量
	defCap int       // 默认缓冲区容量
}

// NewBytePool 创建新的字节切片对象池
//
// 参数:
//   - defCap: 默认缓冲区容量
//   - maxCap: 最大回收缓冲区容量, 超过此容量的缓冲区不会被回收
//
// 返回值:
//   - *BytePool: 字节切片对象池实例
func NewBytePool(defCap, maxCap int) *BytePool {
	if defCap <= 0 {
		defCap = 256 // 默认256字节
	}
	if maxCap <= 0 {
		maxCap = 32 * 1024 // 默认32KB
	}

	return &BytePool{
		maxCap: maxCap,
		defCap: defCap,
		pool: sync.Pool{
			New: func() any {
				return make([]byte, 0, defCap)
			},
		},
	}
}

// Get 获取默认容量的缓冲区
//
// 返回:
//   - []byte: 长度为默认容量, 容量至少为默认容量的缓冲区切片
//
// 说明:
//   - 返回的缓冲区长度等于默认容量, 可以直接使用
//   - 底层容量可能大于默认容量, 来自对象池的复用缓冲区
func (bp *BytePool) Get() []byte {
	return bp.GetCap(bp.defCap)
}

// GetCap 获取指定容量的缓冲区
//
// 参数:
//   - size: 需要的缓冲区容量
//
// 返回:
//   - []byte: 长度为capacity, 容量至少为capacity的缓冲区切片
//
// 说明:
//   - 返回的缓冲区长度等于请求的capacity, 可以直接使用
//   - 底层容量可能大于capacity, 来自对象池的复用缓冲区
//   - 如果capacity <= 0, 使用默认容量
func (bp *BytePool) GetCap(size int) []byte {
	if size <= 0 {
		size = bp.defCap
	}

	buf, ok := bp.pool.Get().([]byte)
	if !ok {
		// 类型断言失败, 创建新的
		return make([]byte, size)
	}

	// 容量足够，返回
	if cap(buf) >= size {
		return buf[:size] // 返回长度为capacity的切片
	}

	// 容量不足，创建新的
	bp.pool.Put(buf) //nolint:all
	return make([]byte, size)
}

// Put 归还缓冲区到对象池
//
// 参数:
//   - buf: 要归还的缓冲区
func (bp *BytePool) Put(buf []byte) {
	if buf == nil || cap(buf) > bp.maxCap {
		return // 不回收空指针或容量超过最大回收容量
	}

	bp.pool.Put(buf[:0]) //nolint:all
}

// GetEmpty 获取指定容量的空缓冲区
//
// 参数:
//   - size: 指定容量要求
//
// 返回:
//   - []byte: 长度为0但容量至少为capacity的缓冲区切片
//
// 说明:
//   - 适用于需要使用append操作逐步构建数据的场景
//   - 避免频繁的内存重新分配
//   - 如果capacity <= 0, 使用默认容量
func (bp *BytePool) GetEmpty(size int) []byte {
	if size <= 0 {
		size = bp.defCap
	}

	buf, ok := bp.pool.Get().([]byte)
	if !ok {
		// 类型断言失败, 创建新的
		return make([]byte, 0, size)
	}

	// 缓冲区容量足够, 返回空切片
	if cap(buf) >= size {
		return buf[:0]
	}

	// 缓冲区容量不足, 创建新的
	bp.pool.Put(buf) //nolint:all
	return make([]byte, 0, size)
}

// // With 使用默认容量的字节缓冲区执行函数，自动管理获取和归还
// //
// // 参数:
// //   - fn: 使用字节缓冲区的函数
// //
// // 返回值:
// //   - []byte: 函数执行后缓冲区的字节数据副本
// //
// // 说明:
// //   - 自动从对象池获取默认容量的字节缓冲区
// //   - 执行用户提供的函数
// //   - 获取缓冲区字节数据的副本
// //   - 自动归还字节缓冲区到对象池
// //   - 即使函数发生panic也会正确归还资源
// func (bp *BytePool) With(fn func(*[]byte)) []byte {
// 	buf := bp.GetEmpty(0)
// 	defer bp.Put(buf)
// 	fn(&buf)
// 	return append([]byte(nil), buf...) // 一次性拷贝
// }

// // WithCap 使用指定容量的字节缓冲区执行函数，自动管理获取和归还
// //
// // 参数:
// //   - size: 字节缓冲区初始容量
// //   - fn: 使用字节缓冲区的函数
// //
// // 返回值:
// //   - []byte: 函数执行后缓冲区的字节数据副本
// //
// // 说明:
// //   - 自动从对象池获取指定容量的字节缓冲区
// //   - 执行用户提供的函数
// //   - 获取缓冲区字节数据的副本
// //   - 自动归还字节缓冲区到对象池
// //   - 即使函数发生panic也会正确归还资源
// func (bp *BytePool) WithCap(size int, fn func(*[]byte)) []byte {
// 	buf := bp.GetEmpty(size)
// 	defer bp.Put(buf)
// 	fn(&buf)
// 	return append([]byte(nil), buf...)
// }

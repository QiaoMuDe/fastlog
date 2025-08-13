/*
smart_buffer_pool.go - 智能分层缓冲区池实现
针对文件和控制台输出分别优化，使用90%阈值触发智能切换，
实现高效的内存管理和性能优化。
*/
package fastlog

import (
	"bytes"
	"sync"
)

// smartTieredBufferPool 智能分层缓冲区池
// 针对文件和控制台输出分别优化，使用90%阈值触发智能切换
type smartTieredBufferPool struct {
	// 文件缓冲区池（大容量，适合批量I/O操作）
	fileSmall  sync.Pool // 32KB - 小文件批量
	fileMedium sync.Pool // 256KB - 中等文件批量
	fileLarge  sync.Pool // 1MB - 大文件批量

	// 控制台缓冲区池（小容量，适合实时显示）
	consoleSmall  sync.Pool // 8KB - 小控制台批量
	consoleMedium sync.Pool // 32KB - 中等控制台批量
	consoleLarge  sync.Pool // 64KB - 大控制台批量
}

// 全局智能缓冲区池实例
var globalSmartBufferPool = newSmartTieredBufferPool()

// newSmartTieredBufferPool 创建新的智能分层缓冲区池
//
// 返回值：
//   - *smartTieredBufferPool: 智能分层缓冲区池实例
func newSmartTieredBufferPool() *smartTieredBufferPool {
	return &smartTieredBufferPool{
		// 文件缓冲区池初始化（基于types.go中的配置）
		fileSmall: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, fileSmallBufferCapacity)) // 32KB
			},
		},
		fileMedium: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, fileMediumBufferCapacity)) // 256KB
			},
		},
		fileLarge: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, fileLargeBufferCapacity)) // 1MB
			},
		},

		// 控制台缓冲区池初始化（更小的容量）
		consoleSmall: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, consoleSmallBufferCapacity)) // 8KB
			},
		},
		consoleMedium: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, consoleMediumBufferCapacity)) // 32KB
			},
		},
		consoleLarge: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, consoleLargeBufferCapacity)) // 64KB
			},
		},
	}
}

// GetFileBuffer 获取文件缓冲区（大容量，适合批量I/O）
//
// 参数：
//   - estimatedSize: 预估的数据大小（字节）
//
// 返回值：
//   - *bytes.Buffer: 合适大小的文件缓冲区
func (stp *smartTieredBufferPool) GetFileBuffer(estimatedSize int) *bytes.Buffer {
	switch {
	//  🎯 关键逻辑：90%阈值触发切换
	case estimatedSize <= fileSmallThreshold: // <= 28.8KB
		if buffer, ok := stp.fileSmall.Get().(*bytes.Buffer); ok {
			return buffer
		}
		return bytes.NewBuffer(make([]byte, 0, fileSmallBufferCapacity))

	// 🎯 关键逻辑：90%阈值触发切换
	case estimatedSize <= fileMediumThreshold: // <= 230.4KB
		if buffer, ok := stp.fileMedium.Get().(*bytes.Buffer); ok {
			return buffer
		}
		return bytes.NewBuffer(make([]byte, 0, fileMediumBufferCapacity))

	// 默认: 默认封顶缓冲区
	default: // > 230.4KB
		if buffer, ok := stp.fileLarge.Get().(*bytes.Buffer); ok {
			return buffer
		}
		return bytes.NewBuffer(make([]byte, 0, fileLargeBufferCapacity))
	}
}

// GetConsoleBuffer 获取控制台缓冲区（小容量，适合实时显示）
//
// 参数：
//   - estimatedSize: 预估的数据大小（字节）
//
// 返回值：
//   - *bytes.Buffer: 合适大小的控制台缓冲区
func (stp *smartTieredBufferPool) GetConsoleBuffer(estimatedSize int) *bytes.Buffer {
	switch {
	// 🎯 关键逻辑：90%阈值触发切换
	case estimatedSize <= consoleSmallThreshold: // <= 7.2KB
		if buffer, ok := stp.consoleSmall.Get().(*bytes.Buffer); ok {
			return buffer
		}
		return bytes.NewBuffer(make([]byte, 0, consoleSmallBufferCapacity))

	// 🎯 关键逻辑：90%阈值触发切换
	case estimatedSize <= consoleMediumThreshold: // <= 28.8KB
		if buffer, ok := stp.consoleMedium.Get().(*bytes.Buffer); ok {
			return buffer
		}
		return bytes.NewBuffer(make([]byte, 0, consoleMediumBufferCapacity))

	// 默认: 默认封顶缓冲区
	default: // > 28.8KB
		if buffer, ok := stp.consoleLarge.Get().(*bytes.Buffer); ok {
			return buffer
		}
		return bytes.NewBuffer(make([]byte, 0, consoleLargeBufferCapacity))
	}
}

// CheckAndUpgradeFileBuffer 检查文件缓冲区是否需要升级
// 当缓冲区使用量达到90%阈值时，自动切换到更大的缓冲区
//
// 参数：
//   - buffer: 当前使用的缓冲区
//   - newDataSize: 即将写入的数据大小
//
// 返回值：
//   - *bytes.Buffer: 升级后的缓冲区（可能是原缓冲区或新缓冲区）
func (stp *smartTieredBufferPool) CheckAndUpgradeFileBuffer(buffer *bytes.Buffer, newDataSize int) *bytes.Buffer {
	if buffer == nil {
		return stp.GetFileBuffer(newDataSize)
	}

	currentLen := buffer.Len()                // 当前缓冲区已使用长度
	currentCap := buffer.Cap()                // 当前缓冲区总容量
	afterWriteLen := currentLen + newDataSize // 新写入的数据长度

	// 🎯 关键逻辑：90%阈值触发切换
	switch {
	case currentCap <= fileSmallBufferCapacity && afterWriteLen > fileSmallThreshold:
		// 小文件缓冲区达到90%，切换到中等缓冲区
		newBuffer := stp.GetFileBuffer(fileMediumBufferCapacity)
		newBuffer.Write(buffer.Bytes()) // 复制现有数据
		stp.PutFileBuffer(buffer)       // 归还旧缓冲区
		return newBuffer

	case currentCap <= fileMediumBufferCapacity && afterWriteLen > fileMediumThreshold:
		// 中等文件缓冲区达到90%，切换到大缓冲区
		newBuffer := stp.GetFileBuffer(fileLargeBufferCapacity)
		newBuffer.Write(buffer.Bytes()) // 复制现有数据
		stp.PutFileBuffer(buffer)       // 归还旧缓冲区
		return newBuffer

	case currentCap <= fileLargeBufferCapacity && afterWriteLen > fileLargeThreshold:
		// 大文件缓冲区达到90%，创建超大缓冲区（不放入池中）
		newBuffer := bytes.NewBuffer(make([]byte, 0, currentCap*2)) // 扩容2倍
		newBuffer.Write(buffer.Bytes())                             // 复制现有数据
		stp.PutFileBuffer(buffer)                                   // 归还旧缓冲区到池中
		return newBuffer                                            // 新缓冲区不进池，GC时自动回收
	}

	return buffer // 未达到阈值，继续使用当前缓冲区
}

// CheckAndUpgradeConsoleBuffer 检查控制台缓冲区是否需要升级
// 当缓冲区使用量达到90%阈值时，自动切换到更大的缓冲区
//
// 参数：
//   - buffer: 当前使用的缓冲区
//   - newDataSize: 即将写入的数据大小
//
// 返回值：
//   - *bytes.Buffer: 升级后的缓冲区（可能是原缓冲区或新缓冲区）
func (stp *smartTieredBufferPool) CheckAndUpgradeConsoleBuffer(buffer *bytes.Buffer, newDataSize int) *bytes.Buffer {
	if buffer == nil {
		return stp.GetConsoleBuffer(newDataSize)
	}

	currentLen := buffer.Len()
	currentCap := buffer.Cap()
	afterWriteLen := currentLen + newDataSize

	// 🎯 关键逻辑：90%阈值触发切换
	switch {
	case currentCap <= consoleSmallBufferCapacity && afterWriteLen > consoleSmallThreshold:
		// 小控制台缓冲区达到90%，切换到中等缓冲区
		newBuffer := stp.GetConsoleBuffer(consoleMediumBufferCapacity)
		newBuffer.Write(buffer.Bytes()) // 复制现有数据
		stp.PutConsoleBuffer(buffer)    // 归还旧缓冲区
		return newBuffer

	case currentCap <= consoleMediumBufferCapacity && afterWriteLen > consoleMediumThreshold:
		// 中等控制台缓冲区达到90%，切换到大缓冲区
		newBuffer := stp.GetConsoleBuffer(consoleLargeBufferCapacity)
		newBuffer.Write(buffer.Bytes()) // 复制现有数据
		stp.PutConsoleBuffer(buffer)    // 归还旧缓冲区
		return newBuffer

	case currentCap <= consoleLargeBufferCapacity && afterWriteLen > consoleLargeThreshold:
		// 大控制台缓冲区达到90%，创建超大缓冲区（不放入池中）
		newBuffer := bytes.NewBuffer(make([]byte, 0, currentCap*2)) // 扩容2倍
		newBuffer.Write(buffer.Bytes())                             // 复制现有数据
		stp.PutConsoleBuffer(buffer)                                // 归还旧缓冲区到池中
		return newBuffer                                            // 新缓冲区不进池，GC时自动回收
	}

	return buffer // 未达到阈值，继续使用当前缓冲区
}

// PutFileBuffer 归还文件缓冲区到对应的池中
// 根据缓冲区的实际容量重新分类到合适的池
//
// 参数：
//   - buffer: 要归还的文件缓冲区
func (stp *smartTieredBufferPool) PutFileBuffer(buffer *bytes.Buffer) {
	if buffer == nil {
		return
	}

	buffer.Reset() // 清空内容但保留容量

	// 根据实际容量重新分类
	switch cap := buffer.Cap(); {
	case cap <= fileSmallBufferCapacity: // <= 32KB
		stp.fileSmall.Put(buffer)
	case cap <= fileMediumBufferCapacity: // <= 256KB
		stp.fileMedium.Put(buffer)
	case cap <= fileLargeBufferCapacity: // <= 1MB
		stp.fileLarge.Put(buffer)
	default:
		// 🗑️ 超大缓冲区不放入池中，让GC回收
		// 这样避免池中积累过大的缓冲区
	}
}

// PutConsoleBuffer 归还控制台缓冲区到对应的池中
// 根据缓冲区的实际容量重新分类到合适的池
//
// 参数：
//   - buffer: 要归还的控制台缓冲区
func (stp *smartTieredBufferPool) PutConsoleBuffer(buffer *bytes.Buffer) {
	if buffer == nil {
		return
	}

	buffer.Reset() // 清空内容但保留容量

	// 根据实际容量重新分类
	switch cap := buffer.Cap(); {
	case cap <= consoleSmallBufferCapacity: // <= 8KB
		stp.consoleSmall.Put(buffer)
	case cap <= consoleMediumBufferCapacity: // <= 32KB
		stp.consoleMedium.Put(buffer)
	case cap <= consoleLargeBufferCapacity: // <= 64KB
		stp.consoleLarge.Put(buffer)
	default:
		// 🗑️ 超大缓冲区不放入池中，让GC回收
		// 这样避免池中积累过大的缓冲区
	}
}

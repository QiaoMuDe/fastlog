/*
buffered_writer.go - 带缓冲批量写入器
实现简洁高效的批量写入优化，通过双重条件触发减少系统调用开销。
*/
package fastlog

import (
	"bytes"
	"io"
	"sync"

	"gitee.com/MM-Q/logrotatex"
)

// BufferedWriter 带缓冲批量写入器
// 内部嵌入日志切割器，提供批量写入功能
type BufferedWriter struct {
	logger *logrotatex.LogRotateX // 嵌入的日志切割器
	buffer *bytes.Buffer          // 缓冲区
	mutex  sync.Mutex             // 保护缓冲区和状态

	// 刷新条件
	maxBufferSize int // 最大缓冲区大小（字节）
	maxLogCount   int // 最大日志条数
	currentCount  int // 当前日志条数

	closed bool // 是否已关闭
}

// BufferedWriterConfig 缓冲写入器配置
type BufferedWriterConfig struct {
	MaxBufferSize int // 最大缓冲区大小，默认64KB
	MaxLogCount   int // 最大日志条数，默认100条
}

// DefaultBufferedWriterConfig 默认缓冲写入器配置
func DefaultBufferedWriterConfig() *BufferedWriterConfig {
	return &BufferedWriterConfig{
		MaxBufferSize: 64 * 1024, // 64KB缓冲区
		MaxLogCount:   100,       // 100条日志
	}
}

// NewBufferedWriter 创建新的带缓冲批量写入器
func NewBufferedWriter(logger *logrotatex.LogRotateX, config *BufferedWriterConfig) *BufferedWriter {
	if config == nil {
		config = DefaultBufferedWriterConfig()
	}

	return &BufferedWriter{
		logger:        logger,
		buffer:        bytes.NewBuffer(make([]byte, 0, config.MaxBufferSize)),
		maxBufferSize: config.MaxBufferSize,
		maxLogCount:   config.MaxLogCount,
	}
}

// Write 实现 io.Writer 接口
// 将数据写入缓冲区，达到刷新条件时自动批量写入
func (bw *BufferedWriter) Write(p []byte) (n int, err error) {
	if bw.closed {
		return 0, io.ErrClosedPipe
	}

	bw.mutex.Lock()
	defer bw.mutex.Unlock()

	// 1. 写入缓冲区
	n, err = bw.buffer.Write(p)
	if err != nil {
		return n, err
	}

	// 2. 增加日志计数
	bw.currentCount++

	// 3. 检查是否需要刷新（双重条件触发）
	if bw.shouldFlush() {
		return n, bw.flushLocked()
	}

	return n, nil
}

// shouldFlush 检查是否应该刷新缓冲区
// 双重条件：缓冲区大小 OR 日志条数
func (bw *BufferedWriter) shouldFlush() bool {
	return bw.buffer.Len() >= bw.maxBufferSize ||
		bw.currentCount >= bw.maxLogCount
}

// flushLocked 刷新缓冲区（需要持有锁）
func (bw *BufferedWriter) flushLocked() error {
	if bw.buffer.Len() == 0 {
		return nil
	}

	// 一次性写入所有数据到日志切割器
	_, err := bw.logger.Write(bw.buffer.Bytes())
	if err != nil {
		return err
	}

	// 重置缓冲区和计数器
	bw.buffer.Reset()
	bw.currentCount = 0
	return nil
}

// Flush 手动刷新缓冲区
func (bw *BufferedWriter) Flush() error {
	bw.mutex.Lock()
	defer bw.mutex.Unlock()
	return bw.flushLocked()
}

// Close 关闭缓冲写入器
func (bw *BufferedWriter) Close() error {
	bw.mutex.Lock()
	defer bw.mutex.Unlock()

	if bw.closed {
		return nil
	}

	bw.closed = true

	// 关闭前最后一次刷新，确保数据不丢失
	err := bw.flushLocked()

	// 关闭日志切割器
	if bw.logger != nil {
		if closeErr := bw.logger.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}

	return err
}

// GetStats 获取缓冲写入器统计信息
func (bw *BufferedWriter) GetStats() (bufferSize, logCount, maxBufferSize, maxLogCount int) {
	bw.mutex.Lock()
	defer bw.mutex.Unlock()

	return bw.buffer.Len(), bw.currentCount, bw.maxBufferSize, bw.maxLogCount
}

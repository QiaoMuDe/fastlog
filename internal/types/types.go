package types

import (
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// 文件大小配置常量
	DefaultMaxFileSize = 10 // 默认最大文件大小（MB）

	// 默认文件写入器配置
	DefaultMaxBufferSize = 64 * 1024       // 默认最大缓冲区大小（64KB）
	DefaultMaxWriteCount = 500             // 默认最大写入次数（500次）
	DefaultFlushInterval = 1 * time.Second // 默认最大刷新间隔（1秒）
)

// 优化的时间戳缓存结构，使用原子操作 + 读写锁的混合方案
// 读取时使用原子操作快速检查，只在必要时使用读写锁
type rwTimestampCache struct {
	lastSecond   int64        // 原子操作的秒数，用于快速检查
	cachedString string       // 缓存的时间戳字符串
	mu           sync.RWMutex // 读写锁，读多写少场景的最佳选择
}

// 全局时间戳缓存实例
var globalRWCache = &rwTimestampCache{}

// GetCachedTimestamp 获取缓存的时间戳，读写锁优化版本
//
// 性能特点：
//   - 快路径：原子操作检查 + 读锁保护
//   - 慢路径：写锁保护更新操作
//   - 多读者并发，单写者独占
//   - 无unsafe操作，完全内存安全
//
// 返回值：
//   - string: 格式化的时间戳字符串 "2006-01-02 15:04:05"
func GetCachedTimestamp() string {
	now := time.Now()           // 获取当前完整时间对象
	currentSecond := now.Unix() // 提取Unix时间戳的秒数部分

	// 🚀 快路径：原子操作快速检查
	if atomic.LoadInt64(&globalRWCache.lastSecond) == currentSecond {
		// 使用读锁保护字符串读取，允许多个goroutine并发读取
		globalRWCache.mu.RLock()
		result := globalRWCache.cachedString
		globalRWCache.mu.RUnlock()
		return result // 大多数情况走这里，性能很好
	}

	// 慢路径：需要更新缓存
	globalRWCache.mu.Lock()
	defer globalRWCache.mu.Unlock()

	// 双重检查：在等待写锁期间，可能其他goroutine已经更新了
	if atomic.LoadInt64(&globalRWCache.lastSecond) == currentSecond {
		return globalRWCache.cachedString
	}

	// 执行更新
	// 先更新字符串，再原子更新秒数（确保一致性）
	newTimestamp := now.Format("2006-01-02 15:04:05")
	globalRWCache.cachedString = newTimestamp
	atomic.StoreInt64(&globalRWCache.lastSecond, currentSecond)

	return newTimestamp
}

// 文件名缓存，用于缓存 filepath.Base() 的结果，减少重复的字符串处理开销
// key: 完整文件路径，value: 文件名（不含路径）
var fileNameCache = sync.Map{}

// NeedsFileInfo 判断日志格式是否需要文件信息
//
// 参数：
//   - format: 日志格式类型
//
// 返回值：
//   - bool: true表示需要文件信息，false表示不需要
func NeedsFileInfo(format LogFormatType) bool {
	_, exists := FileInfoRequiredFormats[format]
	return exists
}

// GetCallerInfo 获取调用者的信息（优化版本，使用文件名缓存）
//
// 参数：
//   - skip: 跳过的调用层数（通常设置为1或2, 具体取决于调用链的深度）
//
// 返回值：
//   - fileName: 调用者的文件名（不包含路径）
//   - functionName: 调用者的函数名
//   - line: 调用者的行号
//   - ok: 是否成功获取到调用者信息
func GetCallerInfo(skip int) (fileName string, functionName string, line uint16, ok bool) {
	// 获取调用者信息, 跳过指定的调用层数
	pc, file, lineInt, ok := runtime.Caller(skip)
	if !ok {
		line = 0
		return
	}

	// 行号转换和边界检查
	if lineInt >= 0 && lineInt <= 65535 {
		line = uint16(lineInt)
	} else {
		line = 0 // 超出范围使用默认值
	}

	// 优化：使用缓存获取文件名，避免重复的 filepath.Base() 调用
	// 尝试从缓存中获取文件名
	if cached, exists := fileNameCache.Load(file); exists {
		// 缓存命中：直接使用缓存的文件名（性能提升5-10倍）
		fileName = cached.(string)
	} else {
		// 缓存未命中：计算文件名并存储到缓存中
		fileName = filepath.Base(file)      // 执行字符串处理："/path/to/file.go" -> "file.go"
		fileNameCache.Store(file, fileName) // 存储到缓存，供后续调用复用
	}

	// 获取函数名（保持原有逻辑）
	function := runtime.FuncForPC(pc)
	if function != nil {
		functionName = function.Name()
	} else {
		functionName = "???"
	}

	return
}

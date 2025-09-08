package fastlog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"gitee.com/MM-Q/colorlib"
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

// getCachedTimestamp 获取缓存的时间戳，读写锁优化版本
//
// 性能特点：
//   - 快路径：原子操作检查 + 读锁保护
//   - 慢路径：写锁保护更新操作
//   - 多读者并发，单写者独占
//   - 无unsafe操作，完全内存安全
//
// 返回值：
//   - string: 格式化的时间戳字符串 "2006-01-02 15:04:05"
func getCachedTimestamp() string {
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

// 临时缓冲区对象池，用于复用临时缓冲区，减少内存分配
var tempBufferPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

// getTempBuffer 从对象池获取临时缓冲区，使用安全的类型断言
func getTempBuffer() *bytes.Buffer {
	// 安全的类型断言
	if buffer, ok := tempBufferPool.Get().(*bytes.Buffer); ok {
		return buffer
	}
	// 如果类型断言失败，创建新的缓冲区作为fallback
	return &bytes.Buffer{}
}

// putTempBuffer 将临时缓冲区归还到对象池
func putTempBuffer(buffer *bytes.Buffer) {
	if buffer != nil {
		buffer.Reset()             // 重置缓冲区内容
		tempBufferPool.Put(buffer) // 归还到对象池
	}
}

// needsFileInfo 判断日志格式是否需要文件信息
//
// 参数：
//   - format: 日志格式类型
//
// 返回值：
//   - bool: true表示需要文件信息，false表示不需要
func needsFileInfo(format LogFormatType) bool {
	_, exists := fileInfoRequiredFormats[format]
	return exists
}

// getCallerInfo 获取调用者的信息（优化版本，使用文件名缓存）
//
// 参数：
//   - skip: 跳过的调用层数（通常设置为1或2, 具体取决于调用链的深度）
//
// 返回值：
//   - fileName: 调用者的文件名（不包含路径）
//   - functionName: 调用者的函数名
//   - line: 调用者的行号
//   - ok: 是否成功获取到调用者信息
func getCallerInfo(skip int) (fileName string, functionName string, line uint16, ok bool) {
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

// shouldDropLogByBackpressure 根据通道背压情况判断是否应该丢弃日志
//
// 参数:
//   - bp: 通道背压阈值
//   - logChan: 日志通道
//   - level: 日志级别
//
// 返回:
//   - bool: true表示应该丢弃该日志, false表示应该保留
func shouldDropLogByBackpressure(bp *bpThresholds, logChan chan *logMsg, level LogLevel) bool {
	// 完整的空指针和边界检查
	if bp == nil || logChan == nil {
		return false // 如果背压阈值或通道为nil, 不丢弃日志
	}

	// 提前获取通道长度和容量, 供后续复用
	chanLen := len(logChan)
	chanCap := cap(logChan)

	// 边界条件检查: 防止除零错误和异常情况
	if chanCap <= 0 {
		return true // 容量为0或负数的通道应该丢弃日志
	}

	// 通道长度不能为负数
	if chanLen < 0 {
		return false
	}

	// 当通道满了, 立即丢弃所有新日志
	if chanLen >= chanCap {
		return true
	}

	// 关键优化: 避免除法，使用数学等价比较
	// 原理: chanLen/chanCap >= X% 等价于 chanLen*100 >= chanCap*X
	chanLen100 := chanLen * 100 // 预计算，避免重复乘法

	// 根据通道使用率判断是否丢弃日志
	switch {
	case chanLen100 >= bp.threshold98: // 98%+ 只保留FATAL
		return level < FATAL
	case chanLen100 >= bp.threshold95: // 95%+ 只保留ERROR及以上
		return level < ERROR
	case chanLen100 >= bp.threshold90: // 90%+ 只保留WARN及以上
		return level < WARN
	case chanLen100 >= bp.threshold80: // 80%+ 只保留INFO及以上
		return level < INFO
	default: // 80%以下不丢弃任何日志
		return false
	}
}

// logWithLevel 通用日志记录方法
//
// 参数:
//   - level: 日志级别
//   - message: 格式化后的消息
//   - skipFrames: 跳过的调用栈帧数（用于获取正确的调用者信息）
func (f *FastLog) logWithLevel(level LogLevel, message string, skipFrames int) {
	// 关键路径空指针检查 - 防止panic
	if f == nil {
		return
	}

	// 检查核心组件是否已初始化
	if f.config == nil || f.logChan == nil {
		return
	}

	// 检查日志通道是否已关闭 - 复用现有的协调机制
	select {
	case <-f.ctx.Done():
		return // 上下文已取消，直接返回
	default:
		// 继续执行
	}

	// 检查日志级别，如果调用的日志级别低于配置的日志级别，则直接返回
	if level < f.config.LogLevel {
		return
	}

	// 验证消息内容 - 空消息直接返回
	if message == "" {
		return
	}

	// 调用者信息获取逻辑
	var (
		fileName = "unknown"
		funcName = "unknown"
		line     uint16
	)

	// 仅当需要文件信息时才获取调用者信息
	if needsFileInfo(f.config.LogFormat) {
		var ok bool
		fileName, funcName, line, ok = getCallerInfo(skipFrames)
		if !ok {
			fileName = "unknown"
			funcName = "unknown"
			line = 0
		}
	}

	// 使用缓存的时间戳，减少重复的时间格式化开销
	timestamp := getCachedTimestamp()

	// 从对象池获取日志消息对象，增加安全检查
	logMessage := getLogMsg()
	if logMessage == nil {
		// 对象池异常，创建新对象作为fallback
		logMessage = &logMsg{}
	}

	// 安全地填充日志消息字段
	logMessage.Timestamp = timestamp // 时间戳
	logMessage.Level = level         // 日志级别
	logMessage.Message = message     // 日志消息
	logMessage.FileName = fileName   // 文件名
	logMessage.FuncName = funcName   // 函数名
	logMessage.Line = line           // 行号

	// 多级背压处理: 根据通道使用率丢弃低级别日志消息
	if shouldDropLogByBackpressure(f.bp, f.logChan, level) {
		// 重要：如果丢弃日志，需要回收对象
		putLogMsg(logMessage)
		return
	}

	// 安全发送日志 - 使用select避免阻塞
	select {
	// 上下文已取消，回收对象
	case <-f.ctx.Done():
		putLogMsg(logMessage)

	// 成功发送
	case f.logChan <- logMessage:
		// 无操作

	// 通道满，回收对象并丢弃日志
	default:
		putLogMsg(logMessage)
	}
}

// logFatal Fatal级别的特殊处理方法
//
// 参数:
//   - message: 格式化后的消息
//   - skipFrames: 跳过的调用栈帧数
func (f *FastLog) logFatal(message string, skipFrames int) {
	// Fatal方法的特殊处理 - 即使FastLog为nil也要记录错误并退出
	if f == nil {
		// 如果日志器为nil，直接输出到stderr并退出
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", message)
		os.Exit(1)
		return
	}

	// 先记录日志
	f.logWithLevel(FATAL, message, skipFrames)

	// 关闭日志记录器
	f.Close()

	// 终止程序（非0退出码表示错误）
	os.Exit(1)
}

// ===== 实现 processorDependencies 接口 =====

// getConfig 获取日志配置
func (f *FastLog) getConfig() *FastLogConfig {
	return f.config
}

// getFileWriter 获取文件写入器
func (f *FastLog) getFileWriter() io.Writer {
	return f.fileWriter
}

// getConsoleWriter 获取控制台写入器
func (f *FastLog) getConsoleWriter() io.Writer {
	return f.consoleWriter
}

// getColorLib 获取颜色库实例
func (f *FastLog) getColorLib() *colorlib.ColorLib {
	return f.cl
}

// getContext 获取上下文
func (f *FastLog) getContext() context.Context {
	return f.ctx
}

// getLogChannel 获取日志消息通道（只读）
func (f *FastLog) getLogChannel() <-chan *logMsg {
	return f.logChan
}

// notifyProcessorDone 通知处理器完成工作
func (f *FastLog) notifyProcessorDone() {
	f.logWait.Done()
}

// getBufferSize 获取缓冲区大小
func (f *FastLog) getBufferSize() int {
	return f.bufferSize
}

// getCloseTimeout 计算并返回日志记录器关闭时的合理超时时间
//
// 返回值:
//   - time.Duration: 计算后的关闭超时时间，范围在3-10秒之间
//
// 实现逻辑:
//  1. 基于配置的刷新间隔(FlushInterval)乘以10作为基础超时时间
//  2. 确保最小超时为3秒，避免过短的超时导致日志丢失
//  3. 确保最大超时为10秒，避免过长的等待影响程序退出
func (f *FastLog) getCloseTimeout() time.Duration {
	// 基于刷新间隔计算合理的超时时间
	baseTimeout := f.config.FlushInterval * 10
	if baseTimeout < 3*time.Second {
		baseTimeout = 3 * time.Second
	}
	if baseTimeout > 10*time.Second {
		baseTimeout = 10 * time.Second
	}
	return baseTimeout
}

// gracefulShutdown 优雅关闭日志记录器
//
// 参数:
//   - ctx: 上下文对象，用于控制关闭过程
func (f *FastLog) gracefulShutdown(ctx context.Context) {
	// 1. 先取消处理器上下文，通知所有组件停止工作
	f.cancel()

	// 2. 等待一小段时间，让正在进行的操作完成
	time.Sleep(10 * time.Millisecond)

	// 3. 等待处理器完成剩余工作
	shutdownComplete := make(chan struct{})
	go func() {
		defer close(shutdownComplete)
		f.logWait.Wait()
	}()

	// 4. 等待完成或超时
	select {
	case <-shutdownComplete:
		// 正常关闭完成
	case <-ctx.Done():
		// 超时，但不打印警告(因为会强制清理)
	}

	// 5. 关闭日志通道，停止接收新日志
	close(f.logChan)
}

// calculateBufferSize 根据批处理数量计算缓冲区大小
// 保证最小16KB和最大1MB的范围
//
// 参数:
//   - batchSize: 批处理数量
//
// 返回值:
//   - int: 缓冲区大小（字节）
func calculateBufferSize(batchSize int) int {
	if batchSize <= 0 {
		return 16 * 1024 // 16KB
	}

	size := batchSize * bytesPerLogEntry

	// 最小16KB，最大1MB
	if size < 16*1024 {
		return 16 * 1024
	}
	if size > 1024*1024 {
		return 1024 * 1024
	}

	return size
}

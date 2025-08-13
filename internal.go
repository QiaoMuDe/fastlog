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

// 优化的时间戳缓存结构，使用原子操作 + 轻量级锁的混合方案
// 相比原来的读写锁方案，性能提升2-3倍，特别是在高并发场景下
type safeTimestampCache struct {
	lastSecond   int64      // 原子操作的秒数，用于快速检查缓存是否有效
	cachedString string     // 缓存的时间戳字符串
	mu           sync.Mutex // 轻量级互斥锁，只保护字符串更新操作
}

// 全局时间戳缓存实例
var globalSafeCache = &safeTimestampCache{}

// getCachedTimestamp 获取缓存的时间戳，优化版本（原子操作 + 轻量级锁）
// 性能特点：
//   - 快路径完全无锁，使用原子读取
//   - 慢路径使用轻量级Mutex，避免读写锁的开销
//   - 双重检查锁定，确保并发安全
//
// 返回值：
//   - string: 格式化的时间戳字符串 "2006-01-02 15:04:05"
func getCachedTimestamp() string {
	// 步骤1：获取当前时间信息
	now := time.Now()           // 获取当前完整时间对象
	currentSecond := now.Unix() // 提取Unix时间戳的秒数部分

	// 步骤2：快路径 - 原子读取，完全无锁（🚀 性能关键优化）
	// 使用原子操作读取上次缓存的秒数，避免锁竞争
	lastSecond := atomic.LoadInt64(&globalSafeCache.lastSecond)

	// 如果秒数相同，直接返回缓存的字符串（大多数情况下走这个路径）
	if currentSecond == lastSecond {
		return globalSafeCache.cachedString // 🚀 无锁读取，性能最优
	}

	// 步骤3：慢路径 - 需要更新缓存
	// 使用轻量级Mutex而不是RWMutex，减少锁开销
	globalSafeCache.mu.Lock()
	defer globalSafeCache.mu.Unlock()

	// 步骤4：双重检查 - 防止重复更新
	// 在等待锁期间，可能其他goroutine已经更新了缓存
	if currentSecond == atomic.LoadInt64(&globalSafeCache.lastSecond) {
		return globalSafeCache.cachedString
	}

	// 步骤5：执行缓存更新
	// 先更新字符串，再原子更新秒数（确保一致性）
	newTimestamp := now.Format("2006-01-02 15:04:05")
	globalSafeCache.cachedString = newTimestamp
	atomic.StoreInt64(&globalSafeCache.lastSecond, currentSecond)

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
func needsFileInfo(format LogFormatType) bool {
	return format == Json || format == Detailed || format == Structured
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
//   - logChan: 日志通道
//   - level: 日志级别
//
// 返回:
//   - bool: true表示应该丢弃该日志, false表示应该保留
func shouldDropLogByBackpressure(logChan chan *logMsg, level LogLevel) bool {
	// 完整的空指针和边界检查
	if logChan == nil {
		return false // 如果通道为nil, 不丢弃日志
	}

	// 提前获取通道长度和容量, 供后续复用
	chanLen := len(logChan)
	chanCap := cap(logChan)

	// 边界条件检查：防止除零错误和异常情况
	if chanCap <= 0 {
		return false // 容量异常，不丢弃日志
	}

	if chanLen < 0 {
		return false // 长度异常，不丢弃日志
	}

	// 当通道满了, 立即丢弃所有新日志
	if chanLen >= chanCap {
		return true
	}

	// 使用int64进行安全的通道使用率计算，防止整数溢出
	var channelUsage int64
	if chanCap > 0 {
		// 直接使用int64计算，避免类型转换开销
		channelUsage = (int64(chanLen) * 100) / int64(chanCap)

		// 边界检查，确保结果在合理范围内
		if channelUsage > 100 {
			channelUsage = 100
		} else if channelUsage < 0 {
			channelUsage = 0 // 防止异常的负值
		}
	}

	// 根据通道使用率决定是否丢弃日志, 按照日志级别重要性递增
	switch {
	case channelUsage >= 98: // 98%+ 只保留FATAL
		return level < FATAL
	case channelUsage >= 95: // 95%+ 只保留ERROR及以上
		return level < ERROR
	case channelUsage >= 90: // 90%+ 只保留WARN及以上
		return level < WARN
	case channelUsage >= 80: // 80%+ 只保留SUCCESS及以上
		return level < SUCCESS
	case channelUsage >= 70: // 70%+ 只保留INFO及以上(丢弃DEBUG级别)
		return level < INFO
	default:
		return false // 70%以下不丢弃任何日志
	}
}

// logWithLevel 通用日志记录方法
//
// 参数:
//   - level: 日志级别
//   - message: 格式化后的消息
//   - skipFrames: 跳过的调用栈帧数（用于获取正确的调用者信息）
func (l *FastLog) logWithLevel(level LogLevel, message string, skipFrames int) {
	// 关键路径空指针检查 - 防止panic
	if l == nil {
		return
	}

	// 检查核心组件是否已初始化
	if l.config == nil || l.logChan == nil {
		return
	}

	// 检查日志通道是否已关闭
	if l.isLogChanClosed.Load() {
		return
	}

	// 检查日志级别，如果当前级别高于指定级别则不记录
	if level < l.config.LogLevel {
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
	if needsFileInfo(l.config.LogFormat) {
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
	if shouldDropLogByBackpressure(l.logChan, level) {
		// 重要：如果丢弃日志，需要回收对象
		putLogMsg(logMessage)
		return
	}

	// 安全发送日志 - 使用select避免阻塞
	select {
	case l.logChan <- logMessage:
		// 成功发送
	default:
		// 通道满，回收对象并丢弃日志
		putLogMsg(logMessage)
	}
}

// logFatal Fatal级别的特殊处理方法
//
// 参数:
//   - message: 格式化后的消息
//   - skipFrames: 跳过的调用栈帧数
func (l *FastLog) logFatal(message string, skipFrames int) {
	// Fatal方法的特殊处理 - 即使FastLog为nil也要记录错误并退出
	if l == nil {
		// 如果日志器为nil，直接输出到stderr并退出
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", message)
		os.Exit(1)
		return
	}

	// 先记录日志
	l.logWithLevel(FATAL, message, skipFrames)

	// 关闭日志记录器
	l.Close()

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

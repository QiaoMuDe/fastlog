/*
tools.go - 工具函数集合
提供路径检查、调用者信息获取、协程ID获取、日志格式化和颜色添加等辅助功能。
*/
package fastlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"gitee.com/MM-Q/colorlib"
)

// checkPath 检查给定路径的信息
//
// 参数：
//   - path: 要检查的路径
//
// 返回值：
//   - PathInfo: 路径信息
//   - error: 错误信息
func checkPath(path string) (PathInfo, error) {
	// 创建一个 PathInfo 结构体
	var info PathInfo

	// 清理路径, 确保没有多余的斜杠
	path = filepath.Clean(path)

	// 设置路径
	info.Path = path

	// 使用 os.Stat 获取文件状态
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果路径不存在, 则直接返回
			info.Exists = false
			return info, fmt.Errorf("路径 '%s' 不存在, 请检查路径是否正确: %s", path, err)
		} else {
			return info, fmt.Errorf("无法访问路径 '%s': %s", path, err)
		}
	}

	// 路径存在, 填充信息
	info.Exists = true                // 标记路径存在
	info.IsFile = !fileInfo.IsDir()   // 通过取反判断是否为文件, 因为 IsDir 返回 false 表示是文件
	info.IsDir = fileInfo.IsDir()     // 直接使用 IsDir 方法判断是否为目录
	info.Size = fileInfo.Size()       // 获取文件大小
	info.Mode = fileInfo.Mode()       // 获取文件权限
	info.ModTime = fileInfo.ModTime() // 获取文件的最后修改时间

	// 返回路径信息结构体
	return info, nil
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

// logLevelToString 将 LogLevel 转换为对应的字符串（不带填充，用于JSON序列化）
//
// 参数：
//   - level: 要转换的日志级别
//
// 返回值：
//   - string: 对应的日志级别字符串, 如果 level 无效, 则返回 "UNKNOWN"
func logLevelToString(level LogLevel) string {
	// 使用预构建的映射表进行O(1)查询（不带填充，适用于JSON）
	if str, exists := logLevelStringMap[level]; exists {
		return str
	}
	return "UNKNOWN"
}

// logLevelToPaddedString 将 LogLevel 转换为带填充的字符串（用于文本格式化）
//
// 参数：
//   - level: 要转换的日志级别
//
// 返回值：
//   - string: 对应的带填充的日志级别字符串（7个字符），如果 level 无效, 则返回 "UNKNOWN"
func logLevelToPaddedString(level LogLevel) string {
	// 使用预构建的带填充映射表进行O(1)查询（适用于文本格式）
	if str, exists := logLevelPaddedStringMap[level]; exists {
		return str
	}
	return "UNKNOWN"
}

// addColorToBuffer 直接将带颜色的消息写入缓冲区，避免创建中间字符串（零拷贝优化版本）
//
// 参数：
//   - buffer: 目标缓冲区
//   - cl: 颜色库实例
//   - level: 日志级别
//   - sourceBuffer: 源缓冲区（包含原始消息内容）
func addColorToBuffer(buffer *bytes.Buffer, cl *colorlib.ColorLib, level LogLevel, sourceBuffer *bytes.Buffer) {
	// 完整的空指针和参数检查
	if buffer == nil || cl == nil || sourceBuffer == nil {
		return
	}

	// 检查源缓冲区是否为空
	if sourceBuffer.Len() == 0 {
		return
	}

	// 获取源缓冲区的内容（避免String()调用的内存分配）
	sourceBytes := sourceBuffer.Bytes()
	sourceString := string(sourceBytes) // 这里仍需要一次转换，但比多次String()调用更高效

	// 根据日志级别添加颜色并直接写入目标缓冲区
	switch level {
	case INFO:
		buffer.WriteString(cl.Sblue(sourceString)) // Blue
	case WARN:
		buffer.WriteString(cl.Syellow(sourceString)) // Yellow
	case ERROR:
		buffer.WriteString(cl.Sred(sourceString)) // Red
	case SUCCESS:
		buffer.WriteString(cl.Sgreen(sourceString)) // Green
	case DEBUG:
		buffer.WriteString(cl.Spurple(sourceString)) // Purple
	case FATAL:
		buffer.WriteString(cl.Sred(sourceString)) // Red
	default:
		// 如果没有匹配到日志级别，直接写入原始内容
		buffer.Write(sourceBytes)
	}
}

// formatLogDirectlyToBuffer 直接将日志消息格式化到缓冲区，避免创建中间字符串（零拷贝优化）
//
// 参数：
//   - buffer: 目标缓冲区
//   - config: 日志配置
//   - logMsg: 日志消息
//   - withColor: 是否添加颜色（用于控制台输出）
//   - colorLib: 颜色库实例（当withColor为true时使用）
func formatLogDirectlyToBuffer(buffer *bytes.Buffer, config *FastLogConfig, logMsg *logMsg, withColor bool, colorLib *colorlib.ColorLib) {
	// 检查参数有效性
	if buffer == nil || config == nil || logMsg == nil || colorLib == nil {
		return
	}

	// 如果时间戳为空，使用缓存的时间戳
	if logMsg.Timestamp == "" {
		logMsg.Timestamp = getCachedTimestamp()
	}

	// 检查关键字段是否为空，设置默认值
	if logMsg.Message == "" {
		return // 消息为空直接返回
	}
	if logMsg.FileName == "" {
		logMsg.FileName = "unknown-file"
	}
	if logMsg.FuncName == "" {
		logMsg.FuncName = "unknown-func"
	}

	// 文本格式处理：先格式化到临时缓冲区，然后根据需要添加颜色
	tempBuffer := getTempBuffer()
	defer putTempBuffer(tempBuffer)

	// 根据日志格式格式化到临时缓冲区
	switch config.LogFormat {
	// JSON格式
	case Json:
		// 序列化为JSON并直接写入缓冲区
		if jsonBytes, err := json.Marshal(logMsg); err == nil {
			tempBuffer.Write(jsonBytes)
		} else {
			// JSON序列化失败时的降级处理
			fmt.Fprintf(tempBuffer,
				logFormatMap[Json],
				logMsg.Timestamp, logLevelToString(logMsg.Level), "unknown", "unknown", 0,
				fmt.Sprintf("原始消息序列化失败: %v | 原始内容: %s", err, logMsg.Message),
			)
		}

	// 详细格式
	case Detailed:
		tempBuffer.WriteString(logMsg.Timestamp) // 时间戳
		tempBuffer.WriteString(" | ")

		levelStr := logLevelToPaddedString(logMsg.Level) // 使用预填充的日志级别字符串
		tempBuffer.WriteString(levelStr)

		tempBuffer.WriteString(" | ")
		tempBuffer.WriteString(logMsg.FileName) // 文件信息
		tempBuffer.WriteByte(':')
		tempBuffer.WriteString(logMsg.FuncName) // 函数
		tempBuffer.WriteByte(':')
		tempBuffer.WriteString(strconv.Itoa(int(logMsg.Line))) // 行号
		tempBuffer.WriteString(" - ")
		tempBuffer.WriteString(logMsg.Message) // 消息

	// 简约格式
	case Simple:
		tempBuffer.WriteString(logMsg.Timestamp) // 时间戳
		tempBuffer.WriteString(" | ")

		levelStr := logLevelToPaddedString(logMsg.Level) // 使用预填充的日志级别字符串
		tempBuffer.WriteString(levelStr)

		tempBuffer.WriteString(" | ")
		tempBuffer.WriteString(logMsg.Message) // 消息

	// 结构化格式
	case Structured:
		tempBuffer.WriteString("T:") // 时间戳
		tempBuffer.WriteString(logMsg.Timestamp)
		tempBuffer.WriteString("|L:") // 日志级别

		levelStr := logLevelToPaddedString(logMsg.Level) // 使用预填充的日志级别字符串
		tempBuffer.WriteString(levelStr)

		tempBuffer.WriteString("|F:") // 文件信息
		tempBuffer.WriteString(logMsg.FileName)
		tempBuffer.WriteByte(':')
		tempBuffer.WriteString(logMsg.FuncName)
		tempBuffer.WriteByte(':')
		tempBuffer.WriteString(strconv.Itoa(int(logMsg.Line)))
		tempBuffer.WriteString("|M:") // 消息
		tempBuffer.WriteString(logMsg.Message)

	// 自定义格式
	case Custom:
		tempBuffer.WriteString(logMsg.Message)

	// 默认情况
	default:
		tempBuffer.WriteString("无法识别的日志格式选项: ")
		fmt.Fprintf(tempBuffer, "%v", config.LogFormat)
	}

	// 根据withColor参数决定是否添加颜色
	if withColor {
		// 使用零拷贝版本：直接将带颜色的内容写入目标缓冲区(控制台)
		addColorToBuffer(buffer, colorLib, logMsg.Level, tempBuffer)
	} else {
		// 直接写入原始内容(文件)
		buffer.Write(tempBuffer.Bytes())
	}
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

// 时间戳缓存结构，用于缓存秒级时间戳，减少重复的时间格式化开销
type timestampCache struct {
	lastSecond      int64        // 上次缓存的秒数（Unix时间戳）
	cachedTimestamp string       // 缓存的格式化时间戳字符串
	mu              sync.RWMutex // 读写锁，保证并发安全
}

// 全局时间戳缓存实例
var globalTimestampCache = &timestampCache{}

// getCachedTimestamp 获取缓存的时间戳，如果缓存过期则更新
// 使用双重检查锁定模式，在高并发场景下提供最佳性能
//
// 返回值：
//   - string: 格式化的时间戳字符串 "2006-01-02 15:04:05"
func getCachedTimestamp() string {
	// 步骤1：获取当前时间信息
	// time.Now() 是系统调用，相对昂贵（约200-300ns）
	now := time.Now()           // 获取当前完整时间对象，包含纳秒精度
	currentSecond := now.Unix() // 提取Unix时间戳的秒数部分，用于缓存键比较

	// 步骤2：快路径 - 尝试读取缓存（高并发优化）
	// 使用读锁允许多个goroutine同时读取缓存，提高并发性能
	// 大多数情况下（同一秒内的后续调用）会走这个快路径
	globalTimestampCache.mu.RLock() // 获取读锁，允许并发读取
	if currentSecond == globalTimestampCache.lastSecond {
		// 缓存命中：当前秒数与缓存的秒数相同
		// 这意味着我们可以复用之前格式化好的时间戳字符串
		// 性能提升：避免了昂贵的Format()调用（约200-300ns -> 20-30ns）
		cached := globalTimestampCache.cachedTimestamp // 复制缓存的时间戳字符串
		globalTimestampCache.mu.RUnlock()              // 释放读锁
		return cached                                  // 直接返回缓存结果，性能最优
	}
	globalTimestampCache.mu.RUnlock() // 缓存未命中，释放读锁准备进入慢路径

	// 步骤3：慢路径 - 更新缓存（双重检查锁定模式）
	// 当缓存未命中时（通常是秒数发生变化），需要更新缓存
	// 使用写锁确保同一时间只有一个goroutine能更新缓存
	globalTimestampCache.mu.Lock()         // 获取写锁，独占访问缓存
	defer globalTimestampCache.mu.Unlock() // 确保函数退出时释放写锁

	// 步骤4：双重检查 - 防止重复更新（重要的并发安全措施）
	// 在等待写锁期间，可能其他goroutine已经更新了缓存
	// 再次检查避免重复的Format()调用，提高效率
	if currentSecond == globalTimestampCache.lastSecond {
		// 其他goroutine已经更新了缓存，直接返回最新的缓存值
		// 这种情况在高并发场景下经常发生
		return globalTimestampCache.cachedTimestamp
	}

	// 步骤5：执行缓存更新 - 格式化新的时间戳
	// 只有在确实需要更新时才执行昂贵的Format()操作
	// 使用之前获取的now对象，保证时间一致性
	globalTimestampCache.lastSecond = currentSecond                          // 更新缓存的秒数键
	globalTimestampCache.cachedTimestamp = now.Format("2006-01-02 15:04:05") // 格式化时间戳并缓存

	// 返回新格式化的时间戳
	// 这个结果将被后续同一秒内的调用复用，实现性能优化
	return globalTimestampCache.cachedTimestamp
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

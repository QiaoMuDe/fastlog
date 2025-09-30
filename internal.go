/*
internal.go - FastLog内部实现文件
包含日志系统的核心内部功能实现，包括时间戳缓存、调用者信息获取、背压控制、
日志消息处理和接口实现等，为FastLog提供高性能的底层支持。
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
	"sync/atomic"
	"time"

	"gitee.com/MM-Q/go-kit/pool"
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

// logFatal Fatal级别的特殊处理方法

// 参数:
//   - message: 格式化后的消息
func (f *FastLog) logFatal(message string) {
	// Fatal方法的特殊处理 - 即使FastLog为nil也要记录错误并退出
	if f == nil {
		// 如果日志器为nil，直接输出到stderr并退出
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", message)
		os.Exit(1)
		return
	}

	// 先记录日志
	f.processLog(FATAL, message)

	// 关闭日志记录器
	f.Close()

	// 终止程序（非0退出码表示错误）
	os.Exit(1)
}

// processLog 内部用于处理日志消息的方法
//
// 参数:
//   - level: 日志级别
//   - msg: 日志消息
func (f *FastLog) processLog(level LogLevel, msg string) {
	// 检查核心组件是否已初始化
	if f == nil || f.config == nil || msg == "" {
		return
	}

	// 检查日志级别，如果调用的日志级别低于配置的日志级别，则直接返回
	if level < f.config.LogLevel {
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
		fileName, funcName, line, ok = getCallerInfo(3)
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
	defer putLogMsg(logMessage) // 确保在函数返回时回收对象

	// 安全地填充日志消息字段
	logMessage.Timestamp = timestamp // 时间戳
	logMessage.Level = level         // 日志级别
	logMessage.Message = msg         // 日志消息
	logMessage.FileName = fileName   // 文件名
	logMessage.FuncName = funcName   // 函数名
	logMessage.Line = line           // 行号

	// 获取缓冲区
	buf := pool.GetBuf()
	defer pool.PutBuf(buf)

	// 根据日志格式格式化到缓冲区
	f.formatLogToBuffer(buf, logMessage)

	// 控制台输出 - 直接使用 colorlib 打印
	if f.config.OutputToConsole {
		srcString := buf.String()

		// 直接调用 colorlib 的打印方法（自带换行）
		switch logMessage.Level {
		case INFO:
			f.cl.Blue(srcString)
		case WARN:
			f.cl.Yellow(srcString)
		case ERROR:
			f.cl.Red(srcString)
		case DEBUG:
			f.cl.Magenta(srcString)
		case FATAL:
			f.cl.Red(srcString)
		default:
			fmt.Println(srcString) // 默认打印
		}
	}

	// 将缓冲区中的日志消息写入日志文件
	if f.config.OutputToFile && f.fileWriter != nil {
		buf.WriteString("\n") // 添加换行符，确保每条日志单独一行
		if _, err := f.fileWriter.Write(buf.Bytes()); err != nil {
			fmt.Printf("Error writing to log file: %v\n", err)
		}
	}

}

// formatLogToBuffer 将日志消息格式化到缓冲区，避免创建中间字符串（零拷贝优化）
//
// 参数:
//   - buf: 目标缓冲区
//   - logmsg: 日志消息
func (f *FastLog) formatLogToBuffer(buf *bytes.Buffer, logmsg *logMsg) {
	// 检查参数有效性
	if buf == nil || logmsg == nil {
		return
	}

	// 如果时间戳为空，使用缓存的时间戳
	if logmsg.Timestamp == "" {
		logmsg.Timestamp = getCachedTimestamp()
	}

	// 检查关键字段是否为空，设置默认值
	if logmsg.Message == "" {
		return // 消息为空直接返回
	}
	if logmsg.FileName == "" {
		logmsg.FileName = "unknown-file"
	}
	if logmsg.FuncName == "" {
		logmsg.FuncName = "unknown-func"
	}

	// 根据日志格式直接格式化到目标缓冲区
	switch f.config.LogFormat {
	// JSON格式
	case Json:
		// 序列化为JSON并直接写入缓冲区
		if jsonBytes, err := json.Marshal(logmsg); err == nil {
			buf.Write(jsonBytes)
		} else {
			// JSON序列化失败时的降级处理
			fmt.Fprintf(buf,
				logFormatMap[Json],
				logmsg.Timestamp, logLevelToString(logmsg.Level), "unknown", "unknown", 0,
				fmt.Sprintf("Failed to serialize original message: %v | Original content: %s", err, logmsg.Message),
			)
		}

	// JsonSimple格式（无文件信息）
	case JsonSimple:
		// 序列化为JSON并直接写入缓冲区
		if jsonBytes, err := json.Marshal(simpleLogMsg{
			Timestamp: logmsg.Timestamp,
			Level:     logmsg.Level,
			Message:   logmsg.Message,
		}); err == nil {
			buf.Write(jsonBytes)
		} else {
			// JSON序列化失败时的降级处理
			fmt.Fprintf(buf, logFormatMap[JsonSimple],
				logmsg.Timestamp, logLevelToString(logmsg.Level), fmt.Sprintf("Failed to serialize: %v | Original: %s", err, logmsg.Message))
		}

	// 详细格式
	case Detailed:
		buf.WriteString(logmsg.Timestamp) // 时间戳
		buf.WriteString(" | ")
		levelStr := logLevelToPaddedString(logmsg.Level) // 使用预填充的日志级别字符串
		buf.WriteString(levelStr)
		buf.WriteString(" | ")
		buf.WriteString(logmsg.FileName) // 文件信息
		buf.WriteByte(':')
		buf.WriteString(logmsg.FuncName) // 函数
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(int(logmsg.Line))) // 行号
		buf.WriteString(" - ")
		buf.WriteString(logmsg.Message) // 消息

	// 简约格式
	case Simple:
		buf.WriteString(logmsg.Timestamp) // 时间戳
		buf.WriteString(" | ")
		levelStr := logLevelToPaddedString(logmsg.Level) // 使用预填充的日志级别字符串
		buf.WriteString(levelStr)
		buf.WriteString(" | ")
		buf.WriteString(logmsg.Message) // 消息

	// 结构化格式
	case Structured:
		buf.WriteString("T:") // 时间戳
		buf.WriteString(logmsg.Timestamp)
		buf.WriteString("|L:")                           // 日志级别
		levelStr := logLevelToPaddedString(logmsg.Level) // 使用预填充的日志级别字符串
		buf.WriteString(levelStr)
		buf.WriteString("|F:") // 文件信息
		buf.WriteString(logmsg.FileName)
		buf.WriteByte(':')
		buf.WriteString(logmsg.FuncName)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(int(logmsg.Line)))
		buf.WriteString("|M:") // 消息
		buf.WriteString(logmsg.Message)

	// 基础结构化格式(无文件信息)
	case BasicStructured:
		buf.WriteString("T:") // 时间戳
		buf.WriteString(logmsg.Timestamp)
		buf.WriteString("|L:")                           // 日志级别
		levelStr := logLevelToPaddedString(logmsg.Level) // 使用预填充的日志级别字符串
		buf.WriteString(levelStr)
		buf.WriteString("|M:") // 消息
		buf.WriteString(logmsg.Message)

	// 简单时间格式
	case SimpleTimestamp:
		buf.WriteString(logmsg.Timestamp) // 时间戳
		buf.WriteString(" ")
		levelStr := logLevelToPaddedString(logmsg.Level) // 使用预填充的日志级别字符串
		buf.WriteString(levelStr)                        // 日志级别
		buf.WriteString(" ")
		buf.WriteString(logmsg.Message) // 消息

	// 自定义格式
	case Custom:
		buf.WriteString(logmsg.Message)

	// 默认情况
	default:
		buf.WriteString("Unrecognized log format option: ")
		fmt.Fprintf(buf, "%v", f.config.LogFormat)
	}
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

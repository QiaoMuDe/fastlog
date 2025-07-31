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
	"strings"
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

	// 清理路径，确保没有多余的斜杠
	path = filepath.Clean(path)

	// 设置路径
	info.Path = path

	// 使用 os.Stat 获取文件状态
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果路径不存在, 则直接返回
			info.Exists = false
			return info, fmt.Errorf("路径 '%s' 不存在，请检查路径是否正确: %s", path, err)
		} else {
			return info, fmt.Errorf("无法访问路径 '%s': %s", path, err)
		}
	}

	// 路径存在，填充信息
	info.Exists = true                // 标记路径存在
	info.IsFile = !fileInfo.IsDir()   // 通过取反判断是否为文件，因为 IsDir 返回 false 表示是文件
	info.IsDir = fileInfo.IsDir()     // 直接使用 IsDir 方法判断是否为目录
	info.Size = fileInfo.Size()       // 获取文件大小
	info.Mode = fileInfo.Mode()       // 获取文件权限
	info.ModTime = fileInfo.ModTime() // 获取文件的最后修改时间

	// 返回路径信息结构体
	return info, nil
}

// getCallerInfo 获取调用者的信息
//
// 参数：
//   - skip: 跳过的调用层数（通常设置为1或2，具体取决于调用链的深度）
//
// 返回值：
//   - fileName: 调用者的文件名（不包含路径）
//   - functionName: 调用者的函数名
//   - line: 调用者的行号
//   - ok: 是否成功获取到调用者信息
func getCallerInfo(skip int) (fileName string, functionName string, line uint16, ok bool) {
	// 获取调用者信息，跳过指定的调用层数
	pc, file, lineInt, ok := runtime.Caller(skip)
	if !ok {
		line = 0
		return
	}

	// 在这里做一次转换和边界检查
	if lineInt >= 0 && lineInt <= 65535 {
		line = uint16(lineInt)
	} else {
		line = 0 // 超出范围使用默认值
	}

	// 获取文件名（只保留文件名，不包含路径）
	fileName = filepath.Base(file)

	// 获取函数名
	function := runtime.FuncForPC(pc)
	if function != nil {
		functionName = function.Name()
	} else {
		functionName = "???"
	}

	return
}

// getGoroutineID 获取当前 Goroutine 的 ID
//
// 返回值：
//   - uint32: 当前 Goroutine 的 ID
func getGoroutineID() uint32 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := bytes.Fields(buf[:n])[1]
	id, _ := strconv.ParseInt(string(idField), 10, 64)

	// 安全转换
	if id >= 0 && id <= 0xFFFFFFFF {
		return uint32(id)
	}
	return 0
}

// logLevelToString 将 LogLevel 转换为对应的字符串，并以大写形式返回
//
// 参数：
//   - level: 要转换的日志级别
//
// 返回值：
//   - string: 对应的日志级别字符串，如果 level 无效，则返回 "UNKNOWN"
func logLevelToString(level LogLevel) string {
	// 使用 switch 语句根据日志级别返回对应的字符串
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case SUCCESS:
		return "SUCCESS"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	case NONE:
		return "NONE"
	default:
		return "UNKNOWN"
	}
}

// addColor 根据日志级别添加颜色
//
// 参数：
//   - f: FastLog 实例
//   - l: 日志消息
//   - s: 原始字符串
//
// 返回值:
//   - string: 带有颜色的字符串
func addColor(f *FastLog, l *logMessage, s string) string {
	// 添加空指针检查
	if f == nil || l == nil || f.cl == nil {
		return s // 如果任何参数为nil，返回原始字符串
	}

	// 根据匹配到的日志级别添加颜色
	switch l.Level {
	case INFO:
		return f.cl.Sblue(s) // Blue
	case WARN:
		return f.cl.Syellow(s) // Yellow
	case ERROR:
		return f.cl.Sred(s) // Red
	case SUCCESS:
		return f.cl.Sgreen(s) // Green
	case DEBUG:
		return f.cl.Spurple(s) // Purple
	case FATAL:
		return f.cl.Sred(s) // Red
	default:
		return s // 如果没有匹配到日志级别，返回原始字符串
	}
}

// formatLog 格式化日志消息（优化版本，使用 strings.Builder 提升性能）
//
// 参数：
//   - f: FastLog 实例
//   - l: 日志消息
//
// 返回值:
//   - string: 格式化后的日志消息
func formatLog(f *FastLog, l *logMessage) string {
	if f == nil || l == nil {
		return "" // 如果 FastLog 或 logMessage 为 nil，返回空字符串
	}

	// 预先格式化公共部分，避免重复计算
	levelStr := logLevelToString(l.Level)

	// 根据日志格式选项，格式化日志消息
	switch f.config.LogFormat {
	// Json格式 - 保持使用 fmt.Sprintf（JSON格式复杂，解析开销可接受）
	case Json:
		// 构建json数据
		logData := logMessage{
			Timestamp:   l.Timestamp,   // 格式化时间
			Level:       l.Level,       // 格式化日志级别
			FileName:    l.FileName,    // 文件名
			FuncName:    l.FuncName,    // 函数名
			Line:        l.Line,        // 行号
			GoroutineID: l.GoroutineID, // 协程ID
			Message:     l.Message,     // 日志消息
		}

		// 编码json
		jsonBytes, err := json.Marshal(logData)
		if err != nil {
			// 使用字符串池处理json编码失败的情况
			errorMsg := fmt.Sprintf("原始消息序列化失败: %v | 原始内容: %s", err, safeDeref(l.Message))
			logData.Message = f.stringPool.Intern(errorMsg)

			// 再次尝试序列化，如果还失败就使用最基本的格式
			if fallbackBytes, fallbackErr := json.Marshal(logData); fallbackErr == nil {
				return string(fallbackBytes)
			} else {
				// 最后的兜底方案：手动构建JSON字符串
				return fmt.Sprintf(
					logFormatMap[Json],
					l.Timestamp, levelStr, "unknown", "unknown", 0, l.GoroutineID, logData.Message,
				)
			}
		}
		return string(jsonBytes)

	// 详细格式 - 使用 strings.Builder 优化
	case Detailed:
		// 动态计算容量: 80 + 消息长度 + 文件名长度 + 函数名长度
		estimatedSize := 80 + len(safeDeref(l.Message)) + len(safeDeref(l.FileName)) + len(safeDeref(l.FuncName))

		var builder strings.Builder
		builder.Grow(estimatedSize)

		builder.WriteString(safeDeref(l.Timestamp))
		builder.WriteString(" | ")

		// 格式化日志级别，左对齐7个字符
		builder.WriteString(levelStr)
		for i := len(levelStr); i < 7; i++ {
			builder.WriteByte(' ')
		}

		builder.WriteString(" | ")
		builder.WriteString(safeDeref(l.FileName))
		builder.WriteByte(':')
		builder.WriteString(safeDeref(l.FuncName))
		builder.WriteByte(':')
		builder.WriteString(strconv.Itoa(int(l.Line)))
		builder.WriteString(" - ")
		builder.WriteString(safeDeref(l.Message))

		return builder.String()

	// 括号格式 - 使用 strings.Builder 优化
	case Bracket:
		// 动态计算容量: 80 + 消息长度 + 文件名长度 + 函数名长度
		estimatedSize := 80 + len(safeDeref(l.Message)) + len(safeDeref(l.FileName)) + len(safeDeref(l.FuncName))

		var builder strings.Builder
		builder.Grow(estimatedSize)

		builder.WriteByte('[')
		builder.WriteString(levelStr)
		builder.WriteString("] ")
		builder.WriteString(safeDeref(l.Message))

		return builder.String()

	// 协程格式 - 使用 strings.Builder 优化
	case Threaded:
		// 动态计算容量: 80 + 消息长度 + 文件名长度 + 函数名长度
		estimatedSize := 80 + len(safeDeref(l.Message)) + len(safeDeref(l.FileName)) + len(safeDeref(l.FuncName))

		var builder strings.Builder
		builder.Grow(estimatedSize)

		builder.WriteString(safeDeref(l.Timestamp))
		builder.WriteString(" | ")

		// 格式化日志级别，左对齐7个字符
		builder.WriteString(levelStr)
		for i := len(levelStr); i < 7; i++ {
			builder.WriteByte(' ')
		}

		builder.WriteString(" | [thread=\"")
		builder.WriteString(strconv.FormatUint(uint64(l.GoroutineID), 10))
		builder.WriteString("\"] ")
		builder.WriteString(safeDeref(l.Message))

		return builder.String()

	// 简约格式 - 使用 strings.Builder 优化
	case Simple:
		// 动态计算容量: 80 + 消息长度 + 文件名长度 + 函数名长度
		estimatedSize := 80 + len(safeDeref(l.Message)) + len(safeDeref(l.FileName)) + len(safeDeref(l.FuncName))

		var builder strings.Builder
		builder.Grow(estimatedSize)

		builder.WriteString(safeDeref(l.Timestamp))
		builder.WriteString(" | ")

		// 格式化日志级别，左对齐7个字符
		builder.WriteString(levelStr)
		for i := len(levelStr); i < 7; i++ {
			builder.WriteByte(' ')
		}

		builder.WriteString(" | ")
		builder.WriteString(safeDeref(l.Message))

		return builder.String()

	// 自定义格式
	case Custom:
		return safeDeref(l.Message)

	// 无法识别的日志格式选项
	default:
		return fmt.Sprintf("无法识别的日志格式选项: %v", f.config.LogFormat)
	}
}

// shouldDropLogByBackpressure 根据通道背压情况判断是否应该丢弃日志
//
// 参数:
//   - logChan: 日志通道
//   - level: 日志级别
//
// 返回:
//   - bool: true表示应该丢弃该日志，false表示应该保留
func shouldDropLogByBackpressure(logChan chan *logMessage, level LogLevel) bool {
	// 添加空指针检查
	if logChan == nil {
		return false // 如果通道为nil，不丢弃日志
	}

	// 提前获取通道长度和容量，供后续复用
	chanLen := len(logChan)
	chanCap := cap(logChan)

	// 当通道满了，立即丢弃所有新日志
	if chanLen >= chanCap {
		return true
	}

	// 使用整数计算通道使用率(扩大100倍避免浮点运算)
	// 例如: chanLen=60, chanCap=100 => channelUsage=60
	channelUsage := chanLen * 100 / chanCap

	// 根据通道使用率决定是否丢弃日志，按照日志级别重要性递增
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

// safeDeref 安全地解引用字符串指针，如果指针为nil则返回空字符串
//
// 参数:
//   - s: 字符串指针
//
// 返回:
//   - 解引用后的字符串，如果指针为nil则返回空字符串
func safeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

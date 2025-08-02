/*
tools.go - 工具函数集合
提供路径检查、调用者信息获取、协程ID获取、日志格式化和颜色添加等辅助功能。
*/
package fastlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

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

// getCallerInfo 获取调用者的信息
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

	// 在这里做一次转换和边界检查
	if lineInt >= 0 && lineInt <= 65535 {
		line = uint16(lineInt)
	} else {
		line = 0 // 超出范围使用默认值
	}

	// 获取文件名（只保留文件名, 不包含路径）
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

// logLevelToString 将 LogLevel 转换为对应的字符串, 并以大写形式返回
//
// 参数：
//   - level: 要转换的日志级别
//
// 返回值：
//   - string: 对应的日志级别字符串, 如果 level 无效, 则返回 "UNKNOWN"
func logLevelToString(level LogLevel) string {
	// 使用预构建的映射表进行O(1)查询
	if str, exists := logLevelStringMap[level]; exists {
		return str
	}
	return "UNKNOWN"
}

// addColorToMessage 根据日志级别为消息添加颜色（纯函数版本）
//
// 参数：
//   - cl: 颜色库实例
//   - level: 日志级别
//   - message: 原始消息字符串
//
// 返回值:
//   - string: 带有颜色的字符串
func addColorToMessage(cl *colorlib.ColorLib, level LogLevel, message string) string {
	// 完整的空指针和参数检查
	if cl == nil {
		return message
	}

	// 检查消息是否为空
	if message == "" {
		return message
	}

	// 根据日志级别添加颜色
	switch level {
	case INFO:
		return cl.Sblue(message) // Blue
	case WARN:
		return cl.Syellow(message) // Yellow
	case ERROR:
		return cl.Sred(message) // Red
	case SUCCESS:
		return cl.Sgreen(message) // Green
	case DEBUG:
		return cl.Spurple(message) // Purple
	case FATAL:
		return cl.Sred(message) // Red
	default:
		return message // 如果没有匹配到日志级别, 返回原始字符串
	}
}

// addColor 根据日志级别添加颜色（兼容性包装器）
//
// 参数：
//   - f: FastLog 实例
//   - l: 日志消息
//   - s: 原始字符串
//
// 返回值:
//   - string: 带有颜色的字符串
func addColor(f *FastLog, l *logMsg, s string) string {
	// 添加空指针检查
	if f == nil || l == nil || f.cl == nil {
		return s
	}

	// 调用新的纯函数版本
	return addColorToMessage(f.cl, l.Level, s)
}

// formatLogMessage 格式化日志消息（纯函数版本，优化版本, 使用 strings.Builder 提升性能）
//
// 参数：
//   - config: 日志配置
//   - logMsg: 日志消息
//
// 返回值:
//   - string: 格式化后的日志消息
func formatLogMessage(config *FastLogConfig, logMsg *logMsg) string {
	// 完整的空指针检查
	if config == nil {
		return ""
	}
	if logMsg == nil {
		return ""
	}

	// 检查关键字段是否为空
	if logMsg.Message == "" {
		return ""
	}
	if logMsg.Timestamp == "" {
		logMsg.Timestamp = "unknown-time"
	}
	if logMsg.FileName == "" {
		logMsg.FileName = "unknown-file"
	}
	if logMsg.FuncName == "" {
		logMsg.FuncName = "unknown-func"
	}

	// 预先格式化公共部分, 避免重复计算
	levelStr := logLevelToString(logMsg.Level)

	// 根据日志格式选项, 格式化日志消息
	switch config.LogFormat {
	// Json格式 - 保持使用 fmt.Sprintf（JSON格式复杂, 解析开销可接受）
	case Json:
		// 直接序列化传入的logMsg结构体
		jsonBytes, err := json.Marshal(logMsg)
		if err != nil {
			// JSON编码失败时的兜底方案：手动构建JSON字符串
			return fmt.Sprintf(
				logFormatMap[Json],
				logMsg.Timestamp, levelStr, "unknown", "unknown", 0,
				fmt.Sprintf("原始消息序列化失败: %v | 原始内容: %s", err, logMsg.Message),
			)
		}
		return string(jsonBytes)

	// 详细格式 - 使用 strings.Builder 优化
	case Detailed:
		// 动态计算容量: 80 + 消息长度 + 文件名长度 + 函数名长度
		estimatedSize := 80 + len(logMsg.Message) + len(logMsg.FileName) + len(logMsg.FuncName)

		var builder strings.Builder
		builder.Grow(estimatedSize)

		builder.WriteString(logMsg.Timestamp)
		builder.WriteString(" | ")

		// 格式化日志级别, 左对齐7个字符
		builder.WriteString(levelStr)
		for i := len(levelStr); i < 7; i++ {
			builder.WriteByte(' ')
		}

		builder.WriteString(" | ")
		builder.WriteString(logMsg.FileName)
		builder.WriteByte(':')
		builder.WriteString(logMsg.FuncName)
		builder.WriteByte(':')
		builder.WriteString(strconv.Itoa(int(logMsg.Line)))
		builder.WriteString(" - ")
		builder.WriteString(logMsg.Message)

		return builder.String()

	// 简约格式 - 使用 strings.Builder 优化
	case Simple:
		// 动态计算容量: 80 + 消息长度
		estimatedSize := 80 + len(logMsg.Message)

		var builder strings.Builder
		builder.Grow(estimatedSize)

		builder.WriteString(logMsg.Timestamp)
		builder.WriteString(" | ")

		// 格式化日志级别, 左对齐7个字符
		builder.WriteString(levelStr)
		for i := len(levelStr); i < 7; i++ {
			builder.WriteByte(' ')
		}

		builder.WriteString(" | ")
		builder.WriteString(logMsg.Message)

		return builder.String()

	// 结构化格式 - 使用 strings.Builder 优化
	case Structured:
		estimatedSize := 100 + len(logMsg.Message) + len(logMsg.FileName) + len(logMsg.FuncName)

		var builder strings.Builder
		builder.Grow(estimatedSize)

		builder.WriteString("T:") // 时间戳
		builder.WriteString(logMsg.Timestamp)
		builder.WriteString("|L:") // 格式化日志级别, 左对齐7个字符

		// 格式化日志级别, 左对齐7个字符
		builder.WriteString(levelStr)
		for i := len(levelStr); i < 7; i++ {
			builder.WriteByte(' ')
		}

		builder.WriteString("|F:") // 文件信息
		builder.WriteString(logMsg.FileName)
		builder.WriteByte(':')
		builder.WriteString(logMsg.FuncName)
		builder.WriteByte(':')
		builder.WriteString(strconv.Itoa(int(logMsg.Line)))
		builder.WriteString("|M:") // 消息
		builder.WriteString(logMsg.Message)

		return builder.String()

	// 自定义格式
	case Custom:
		return logMsg.Message

	// 无法识别的日志格式选项
	default:
		return fmt.Sprintf("无法识别的日志格式选项: %v", config.LogFormat)
	}
}

// formatLog 格式化日志消息（兼容性包装器）
//
// 参数：
//   - f: FastLog 实例
//   - l: 日志消息
//
// 返回值:
//   - string: 格式化后的日志消息
func formatLog(f *FastLog, l *logMsg) string {
	if f == nil || l == nil {
		return ""
	}

	// 调用新的纯函数版本
	return formatLogMessage(f.config, l)
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

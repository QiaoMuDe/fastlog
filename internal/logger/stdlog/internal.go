package stdlog

import (
	"bytes"
	"fmt"
	"os"

	"gitee.com/MM-Q/fastlog/internal/types"
	"gitee.com/MM-Q/go-kit/pool"
	"gitee.com/MM-Q/go-kit/utils"
)

// logFatal Fatal级别的特殊处理方法
//
// 参数:
//   - message: 格式化后的消息
func (s *StdLog) logFatal(message string) {
	// Fatal方法的特殊处理 - 即使StdLog为nil也要记录错误并退出
	if s == nil {
		// 如果日志器为nil，直接输出到stderr并退出
		fmt.Printf("FATAL: %s\n", message)
		os.Exit(1)
		return
	}

	// 先记录日志
	s.processLog(types.FATAL, message)

	// 关闭日志记录器
	if err := s.Close(); err != nil {
		// 如果关闭失败，记录到stderr
		fmt.Printf("FATAL: failed to close logger: %v\n", err)
	}

	// 终止程序（非0退出码表示错误）
	os.Exit(1)
}

// processLog 内部用于处理日志消息的方法
//
// 参数:
//   - level: 日志级别
//   - msg: 日志消息
func (s *StdLog) processLog(level types.LogLevel, msg string) {
	// 检查核心组件是否已初始化
	if s == nil || s.cfg == nil || msg == "" {
		return
	}

	// 检查日志级别，使用位运算判断是否应该记录该级别的日志
	if !types.ShouldLog(level, s.cfg.LogLevel) {
		return
	}

	// 使用缓存的时间戳，减少重复的时间格式化开销
	timestamp := types.GetCachedTimestamp()

	// 从对象池获取日志消息对象，增加安全检查
	logMessage := getLogMsg()
	defer putLogMsg(logMessage) // 确保在函数返回时回收对象

	// 安全地填充日志消息字段
	logMessage.Timestamp = timestamp // 时间戳
	logMessage.Level = level         // 日志级别
	logMessage.Message = msg         // 日志消息

	// 仅当需要文件信息时才获取调用者信息
	if s.cfg.CallerInfo {
		logMessage.Caller = types.GetCallerInfo(3)
	}

	// 获取预分配的缓冲区，避免动态内存分配
	buf := pool.GetBufCap(256)
	defer pool.PutBuf(buf)

	// 根据日志格式格式化到缓冲区
	s.formatLogToBuffer(buf, logMessage)

	// 控制台输出 - 直接使用 colorlib 打印
	if s.cfg.OutputToConsole {
		// 直接调用 colorlib 的打印方法（自带换行）
		switch logMessage.Level {
		case types.INFO:
			s.cl.Blue(buf.String())
		case types.WARN:
			s.cl.Yellow(buf.String())
		case types.ERROR:
			s.cl.Red(buf.String())
		case types.DEBUG:
			s.cl.Magenta(buf.String())
		case types.FATAL:
			s.cl.Red(buf.String())
		default:
			fmt.Println(buf.String()) // 默认打印
		}
	}

	// 将缓冲区中的日志消息写入日志文件
	if s.cfg.OutputToFile && s.fileWriter != nil {
		buf.WriteString("\n") // 添加换行符，确保每条日志单独一行
		if _, err := s.fileWriter.Write(buf.Bytes()); err != nil {
			fmt.Printf("Error writing to log file: %v\n", err)
		}
	}
}

// formatLogToBuffer 将日志消息格式化到缓冲区，避免创建中间字符串（零拷贝优化）
//
// 参数:
//   - buf: 目标缓冲区
//   - logmsg: 日志消息
func (s *StdLog) formatLogToBuffer(buf *bytes.Buffer, logmsg *logMsg) {
	// 检查参数有效性
	if buf == nil || logmsg == nil {
		return
	}

	// 检查关键字段是否为空，设置默认值
	if logmsg.Message == "" {
		return // 消息为空直接返回
	}

	// 如果时间戳为空，使用缓存的时间戳
	if logmsg.Timestamp == "" {
		logmsg.Timestamp = types.GetCachedTimestamp()
	}

	// 根据日志格式直接格式化到目标缓冲区
	switch s.cfg.LogFormat {
	case types.Json: // Json格式
		buf.Write([]byte(`{"time":"`))
		buf.WriteString(logmsg.Timestamp)
		buf.Write([]byte(`","level":"`))
		buf.WriteString(logmsg.Level.String())
		if s.cfg.CallerInfo {
			// 仅当需要文件信息时才添加caller字段
			buf.Write([]byte(`","caller":"`))
			buf.Write(logmsg.Caller)
		}
		buf.Write([]byte(`","msg":"`))
		buf.WriteString(utils.QuoteString(logmsg.Message))
		buf.Write([]byte(`"}`))

	case types.Def: // 默认格式
		buf.WriteString(logmsg.Timestamp) // 时间戳
		buf.WriteString(" | ")
		buf.WriteString(types.LogLevelToPaddedString(logmsg.Level))
		buf.WriteString(" | ")
		if s.cfg.CallerInfo {
			buf.Write(logmsg.Caller) // 调用者信息
			buf.WriteString(" - ")
		}
		buf.WriteString(utils.QuoteString(logmsg.Message))

	case types.Structured: // 结构化格式
		buf.WriteString("T:") // 时间戳
		buf.WriteString(logmsg.Timestamp)
		buf.WriteString("|L:") // 日志级别
		buf.WriteString(types.LogLevelToPaddedString(logmsg.Level))
		if s.cfg.CallerInfo {
			buf.WriteString("|C:") // 调用者信息
			buf.Write(logmsg.Caller)
		}
		buf.WriteString("|M:") // 消息
		buf.WriteString(utils.QuoteString(logmsg.Message))

	case types.Timestamp: // 时间格式
		buf.WriteString(logmsg.Timestamp) // 时间戳
		buf.WriteString(" ")
		buf.WriteString(types.LogLevelToPaddedString(logmsg.Level)) // 日志级别
		buf.WriteString(" ")
		buf.WriteString(utils.QuoteString(logmsg.Message))

	case types.Custom: // 自定义格式
		buf.WriteString(utils.QuoteString(logmsg.Message))

	default: // 未识别的日志格式选项
		buf.WriteString(logmsg.Timestamp) // 时间戳
		buf.WriteString(" | ")
		buf.WriteString(types.LogLevelToPaddedString(logmsg.Level))
		buf.WriteString(" | ")
		if s.cfg.CallerInfo {
			buf.Write(logmsg.Caller) // 调用者信息
			buf.WriteString(" - ")
		}
		buf.WriteString(utils.QuoteString(logmsg.Message))
	}
}

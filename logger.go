/*
logger.go - 日志记录方法实现
提供不同级别日志的记录方法（带占位符和不带占位符），
实现日志级别过滤和调用者信息获取功能。
*/
package fastlog

import (
	"fmt"
	"os"
)

/*====== 辅助函数 ======*/

// needsFileInfo 判断日志格式是否需要文件信息
func needsFileInfo(format LogFormatType) bool {
	return format == Json || format == Detailed || format == Structured
}

/*====== 内部通用方法 ======*/

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

/* ====== 不带占位符方法 ======*/

// Info 记录信息级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Info(v ...any) {
	// 公共API入口参数验证
	if l == nil {
		return
	}
	l.logWithLevel(INFO, fmt.Sprint(v...), 3)
}

// Debug 记录调试级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Debug(v ...any) {
	// 公共API入口参数验证
	if l == nil {
		return
	}
	l.logWithLevel(DEBUG, fmt.Sprint(v...), 3)
}

// Warn 记录警告级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Warn(v ...any) {
	// 公共API入口参数验证
	if l == nil {
		return
	}
	l.logWithLevel(WARN, fmt.Sprint(v...), 3)
}

// Error 记录错误级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Error(v ...any) {
	// 公共API入口参数验证
	if l == nil {
		return
	}
	l.logWithLevel(ERROR, fmt.Sprint(v...), 3)
}

// Success 记录成功级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Success(v ...any) {
	// 公共API入口参数验证
	if l == nil {
		return
	}
	l.logWithLevel(SUCCESS, fmt.Sprint(v...), 3)
}

// Fatal 记录致命级别的日志，不支持占位符，发送后关闭日志记录器
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Fatal(v ...any) {
	// 公共API入口参数验证
	if l == nil {
		return
	}
	l.logFatal(fmt.Sprint(v...), 3)
}

/*====== 占位符方法 ======*/

// Infof 记录信息级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Infof(format string, v ...any) {
	// 公共API入口参数验证
	if l == nil || format == "" {
		return
	}
	l.logWithLevel(INFO, fmt.Sprintf(format, v...), 3)
}

// Debugf 记录调试级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Debugf(format string, v ...any) {
	// 公共API入口参数验证
	if l == nil || format == "" {
		return
	}
	l.logWithLevel(DEBUG, fmt.Sprintf(format, v...), 3)
}

// Warnf 记录警告级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Warnf(format string, v ...any) {
	// 公共API入口参数验证
	if l == nil || format == "" {
		return
	}
	l.logWithLevel(WARN, fmt.Sprintf(format, v...), 3)
}

// Errorf 记录错误级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Errorf(format string, v ...any) {
	// 公共API入口参数验证
	if l == nil || format == "" {
		return
	}
	l.logWithLevel(ERROR, fmt.Sprintf(format, v...), 3)
}

// Successf 记录成功级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Successf(format string, v ...any) {
	// 公共API入口参数验证
	if l == nil || format == "" {
		return
	}
	l.logWithLevel(SUCCESS, fmt.Sprintf(format, v...), 3)
}

// Fatalf 记录致命级别的日志，支持占位符，发送后关闭日志记录器
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Fatalf(format string, v ...any) {
	// 公共API入口参数验证
	if l == nil || format == "" {
		return
	}
	l.logFatal(fmt.Sprintf(format, v...), 3)
}

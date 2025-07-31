/*
logger.go - 日志记录方法实现
提供不同级别日志的记录方法（带占位符和不带占位符），
实现日志级别过滤和调用者信息获取功能。
*/
package fastlog

import (
	"fmt"
	"os"
	"time"
)

/*====== 内部通用方法 ======*/

// logWithLevel 通用日志记录方法
//
// 参数:
//   - level: 日志级别
//   - message: 格式化后的消息
//   - skipFrames: 跳过的调用栈帧数（用于获取正确的调用者信息）
func (l *FastLog) logWithLevel(level LogLevel, message string, skipFrames int) {
	// 检查日志级别，如果当前级别高于指定级别则不记录
	if level < l.config.LogLevel {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(skipFrames)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 直接获取当前时间，避免不必要的转换
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 从对象池获取日志消息对象
	logMessage := getLogMsg()

	// 使用字符串池
	logMessage.Timestamp = l.stringPool.Intern(timestamp) // 时间戳
	logMessage.Level = level                              // 日志级别
	logMessage.Message = l.stringPool.Intern(message)     // 日志消息
	logMessage.FileName = l.stringPool.Intern(filename)   // 文件名
	logMessage.FuncName = l.stringPool.Intern(funcName)   // 函数名
	logMessage.Line = line                                // 行号

	// 多级背压处理: 根据通道使用率丢弃低级别日志消息
	if shouldDropLogByBackpressure(l.logChan, level) {
		// 重要：如果丢弃日志，需要回收对象
		putLogMsg(logMessage)
		return
	}

	// 发送日志
	l.logChan <- logMessage
}

// logFatal Fatal级别的特殊处理方法
//
// 参数:
//   - message: 格式化后的消息
//   - skipFrames: 跳过的调用栈帧数
func (l *FastLog) logFatal(message string, skipFrames int) {
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
	l.logWithLevel(INFO, fmt.Sprint(v...), 3)
}

// Debug 记录调试级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Debug(v ...any) {
	l.logWithLevel(DEBUG, fmt.Sprint(v...), 3)
}

// Warn 记录警告级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Warn(v ...any) {
	l.logWithLevel(WARN, fmt.Sprint(v...), 3)
}

// Error 记录错误级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Error(v ...any) {
	l.logWithLevel(ERROR, fmt.Sprint(v...), 3)
}

// Success 记录成功级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Success(v ...any) {
	l.logWithLevel(SUCCESS, fmt.Sprint(v...), 3)
}

// Fatal 记录致命级别的日志，不支持占位符，发送后关闭日志记录器
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Fatal(v ...any) {
	l.logFatal(fmt.Sprint(v...), 3)
}

/*====== 占位符方法 ======*/

// Infof 记录信息级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Infof(format string, v ...any) {
	l.logWithLevel(INFO, fmt.Sprintf(format, v...), 3)
}

// Debugf 记录调试级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Debugf(format string, v ...any) {
	l.logWithLevel(DEBUG, fmt.Sprintf(format, v...), 3)
}

// Warnf 记录警告级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Warnf(format string, v ...any) {
	l.logWithLevel(WARN, fmt.Sprintf(format, v...), 3)
}

// Errorf 记录错误级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Errorf(format string, v ...any) {
	l.logWithLevel(ERROR, fmt.Sprintf(format, v...), 3)
}

// Successf 记录成功级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Successf(format string, v ...any) {
	l.logWithLevel(SUCCESS, fmt.Sprintf(format, v...), 3)
}

// Fatalf 记录致命级别的日志，支持占位符，发送后关闭日志记录器
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (l *FastLog) Fatalf(format string, v ...any) {
	l.logFatal(fmt.Sprintf(format, v...), 3)
}

/*
logger.go - 日志记录方法实现
提供不同级别日志的记录方法（带占位符和不带占位符），
实现日志级别过滤和调用者信息获取功能。
*/
package fastlog

import (
	"fmt"
)

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

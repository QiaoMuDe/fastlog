// 用于存放fastlog的方法
package fastlog

import (
	"fmt"
	"time"
)

// Info 记录信息级别的日志，不支持占位符
func (l *FastLog) Info(v ...any) {
	// 检查日志级别，如果小于等于 Info 级别，则不记录日志。
	if INFO < l.config.LogLevel {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,        // 时间戳
		level:       INFO,             // 日志级别
		message:     fmt.Sprint(v...), // 日志消息
		fileName:    filename,         // 文件名
		funcName:    funcName,         // 函数名
		line:        line,             // 行号
		goroutineID: getGoroutineID(), // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}

// debug 记录调试级别的日志，不支持占位符
func (l *FastLog) Debug(v ...any) {
	// 检查日志级别，如果小于等于 Debug 级别，则不记录日志。
	if DEBUG < l.config.LogLevel {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,        // 时间戳
		level:       DEBUG,            // 日志级别
		message:     fmt.Sprint(v...), // 日志消息
		fileName:    filename,         // 文件名
		funcName:    funcName,         // 函数名
		line:        line,             // 行号
		goroutineID: getGoroutineID(), // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}

// Warn 记录警告级别的日志，不支持占位符
func (l *FastLog) Warn(v ...any) {
	// 检查日志级别，如果小于等于 Warn 级别，则不记录日志。
	if WARN < l.config.LogLevel {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,        // 时间戳
		level:       WARN,             // 日志级别
		message:     fmt.Sprint(v...), // 日志消息
		fileName:    filename,         // 文件名
		funcName:    funcName,         // 函数名
		line:        line,             // 行号
		goroutineID: getGoroutineID(), // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}

// Error 记录错误级别的日志，不支持占位符
func (l *FastLog) Error(v ...any) {
	// 检查日志级别，如果小于等于 Error 级别，则不记录日志。
	if ERROR < l.config.LogLevel {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,        // 时间戳
		level:       ERROR,            // 日志级别
		message:     fmt.Sprint(v...), // 日志消息
		fileName:    filename,         // 文件名
		funcName:    funcName,         // 函数名
		line:        line,             // 行号
		goroutineID: getGoroutineID(), // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}

// Success 记录成功级别的日志，不支持占位符
func (l *FastLog) Success(v ...any) {
	// 检查日志级别，如果小于等于 Success 级别，则不记录日志。
	if SUCCESS < l.config.LogLevel {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,        // 时间戳
		level:       SUCCESS,          // 日志级别
		message:     fmt.Sprint(v...), // 日志消息
		fileName:    filename,         // 文件名
		funcName:    funcName,         // 函数名
		line:        line,             // 行号
		goroutineID: getGoroutineID(), // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}

// Infof 记录信息级别的日志，支持占位符，格式化
func (l *FastLog) Infof(format string, v ...any) {
	// 检查日志级别，如果小于等于 Info 级别，则不记录日志。
	if INFO < l.config.LogLevel {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,                 // 时间戳
		level:       INFO,                      // 日志级别
		message:     fmt.Sprintf(format, v...), // 日志消息
		fileName:    filename,                  // 文件名
		funcName:    funcName,                  // 函数名
		line:        line,                      // 行号
		goroutineID: getGoroutineID(),          // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}

// Debugf 记录调试级别的日志，支持占位符，格式化
func (l *FastLog) Debugf(format string, v ...any) {
	// 检查日志级别，如果小于等于 Debug 级别，则不记录日志。
	if DEBUG < l.config.LogLevel {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,                 // 时间戳
		level:       DEBUG,                     // 日志级别
		message:     fmt.Sprintf(format, v...), // 日志消息
		fileName:    filename,                  // 文件名
		funcName:    funcName,                  // 函数名
		line:        line,                      // 行号
		goroutineID: getGoroutineID(),          // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg
}

// Warnf 记录警告级别的日志，支持占位符，格式化
func (l *FastLog) Warnf(format string, v ...any) {
	// 检查日志级别，如果小于等于 Warn 级别，则不记录日志。
	if WARN < l.config.LogLevel {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,                 // 时间戳
		level:       WARN,                      // 日志级别
		message:     fmt.Sprintf(format, v...), // 日志消息
		fileName:    filename,                  // 文件名
		funcName:    funcName,                  // 函数名
		line:        line,                      // 行号
		goroutineID: getGoroutineID(),          // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}

// Errorf 记录错误级别的日志，支持占位符，格式化
func (l *FastLog) Errorf(format string, v ...any) {
	// 检查日志级别，如果小于等于 Error 级别，则不记录日志。
	if ERROR < l.config.LogLevel {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,                 // 时间戳
		level:       ERROR,                     // 日志级别
		message:     fmt.Sprintf(format, v...), // 日志消息
		fileName:    filename,                  // 文件名
		funcName:    funcName,                  // 函数名
		line:        line,                      // 行号
		goroutineID: getGoroutineID(),          // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}

// Successf 记录成功级别的日志，支持占位符，格式化
func (l *FastLog) Successf(format string, v ...any) {
	// 检查日志级别，如果小于等于 Success 级别，则不记录日志。
	if SUCCESS < l.config.LogLevel {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,                 // 时间戳
		level:       SUCCESS,                   // 日志级别
		message:     fmt.Sprintf(format, v...), // 日志消息
		fileName:    filename,                  // 文件名
		funcName:    funcName,                  // 函数名
		line:        line,                      // 行号
		goroutineID: getGoroutineID(),          // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}

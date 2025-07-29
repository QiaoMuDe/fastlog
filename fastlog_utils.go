// 用于存放fastlog的方法
package fastlog

import (
	"fmt"
	"os"
	"time"
)

// Info 记录信息级别的日志，不支持占位符
func (l *FastLog) Info(v ...any) {
	// 检查日志级别，如果小于等于 Info 级别，则不记录日志。
	if INFO < l.config.GetLogLevel() {
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
	if DEBUG < l.config.GetLogLevel() {
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
	if WARN < l.config.GetLogLevel() {
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
	if ERROR < l.config.GetLogLevel() {
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
	if SUCCESS < l.config.GetLogLevel() {
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
	if INFO < l.config.GetLogLevel() {
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
	if DEBUG < l.config.GetLogLevel() {
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
	if WARN < l.config.GetLogLevel() {
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
	if ERROR < l.config.GetLogLevel() {
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
	if SUCCESS < l.config.GetLogLevel() {
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

// Fatal 记录致命级别的日志，不支持占位符，发送后关闭日志记录器
func (l *FastLog) Fatal(v ...any) {
	// 检查日志级别，如果当前级别高于Fatal则不记录
	if FATAL < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,        // 时间戳
		level:       FATAL,            // 日志级别
		message:     fmt.Sprint(v...), // 日志消息
		fileName:    filename,         // 文件名
		funcName:    funcName,         // 函数名
		line:        line,             // 行号
		goroutineID: getGoroutineID(), // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

	// 等待日志处理完成 - 更精确的同步方式
	// 创建一个完成信号通道
	done := make(chan struct{})

	// 启动一个goroutine检查日志是否处理完成
	go func() {
		// 等待通道为空（所有日志都被处理）
		for len(l.logChan) > 0 {
			// 休眠10毫秒
			time.Sleep(10 * time.Millisecond)
		}
		// 再等待一个刷新周期确保写入文件
		time.Sleep(l.config.GetFlushInterval())
		close(done)
	}()

	// 等待完成信号，但设置超时避免无限等待
	select {
	case <-done:
		// 日志处理完成
	case <-time.After(3 * time.Second):
		// 超时保护，避免无限等待
	}

	// 发送完成后关闭日志记录器
	_ = l.Close()

	// 终止程序（非0退出码表示错误）
	os.Exit(1)
}

// Fatalf 记录致命级别的日志，支持占位符，发送后关闭日志记录器
func (l *FastLog) Fatalf(format string, v ...any) {
	// 检查日志级别，如果当前级别高于Fatal则不记录
	if FATAL < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,                 // 时间戳
		level:       FATAL,                     // 日志级别
		message:     fmt.Sprintf(format, v...), // 日志消息
		fileName:    filename,                  // 文件名
		funcName:    funcName,                  // 函数名
		line:        line,                      // 行号
		goroutineID: getGoroutineID(),          // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

	// 等待日志处理完成 - 更精确的同步方式
	// 创建一个完成信号通道
	done := make(chan struct{})

	// 启动一个goroutine检查日志是否处理完成
	go func() {
		// 等待通道为空（所有日志都被处理）
		for len(l.logChan) > 0 {
			// 休眠10毫秒
			time.Sleep(10 * time.Millisecond)
		}
		// 再等待一个刷新周期确保写入文件
		time.Sleep(l.config.GetFlushInterval())
		close(done)
	}()

	// 等待完成信号，但设置超时避免无限等待
	select {
	case <-done:
		// 日志处理完成
	case <-time.After(3 * time.Second):
		// 超时保护，避免无限等待
	}

	// 发送完成后关闭日志记录器
	_ = l.Close()

	// 终止程序（非0退出码表示错误）
	os.Exit(1)
}

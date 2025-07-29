// 定义全局常量变量或结构体
package fastlog

import (
	"time"
)

// 定义缓冲区相关常量
const (
	// 1 KB
	kb = 1024

	// 8KB 初始容量
	initialBufferCapacity = 8 * kb

	// 256KB 最大容量
	maxBufferCapacity = 256 * kb

	// 缓冲区90%阈值
	flushThreshold = maxBufferCapacity * 9 / 10
)

// 日志格式选项
type LogFormatType int

// 日志格式选项
const (
	Detailed LogFormatType = iota // 详细格式
	Bracket                       // 方括号格式
	Json                          // json格式
	Threaded                      // 协程格式
	Simple                        // 简约格式
	Custom                        // 自定义格式
)

// 日志级别枚举
type LogLevel int

// 定义日志级别
const (
	DEBUG   LogLevel = 10  // 调试级别
	INFO    LogLevel = 20  // 信息级别
	SUCCESS LogLevel = 30  // 成功级别
	WARN    LogLevel = 40  // 警告级别
	ERROR   LogLevel = 50  // 错误级别
	FATAL   LogLevel = 60  // 致命级别
	NONE    LogLevel = 999 // 无日志级别
)

// 定义一个接口, 声明对外暴露的方法
type FastLogInterface interface {
	Close() // 关闭日志记录器

	Info(v ...any)    // 记录信息级别的日志，不支持占位符
	Warn(v ...any)    // 记录警告级别的日志，不支持占位符
	Error(v ...any)   // 记录错误级别的日志，不支持占位符
	Success(v ...any) // 记录成功级别的日志，不支持占位符
	Debug(v ...any)   // 记录调试级别的日志，不支持占位符
	Fatal(v ...any)   // 记录致命级别的日志，不支持占位符(调用后程序会退出)

	Infof(format string, v ...any)    // 记录信息级别的日志，支持占位符，格式化
	Warnf(format string, v ...any)    // 记录警告级别的日志，支持占位符，格式化
	Errorf(format string, v ...any)   // 记录错误级别的日志，支持占位符，格式化
	Successf(format string, v ...any) // 记录成功级别的日志，支持占位符，格式化
	Debugf(format string, v ...any)   // 记录调试级别的日志，支持占位符，格式化
	Fatalf(format string, v ...any)   // 记录致命级别的日志，支持占位符，格式化(调用后程序会退出)
}

// 定义一个结构体，表示日志消息的元数据
type logMessage struct {
	timestamp   time.Time // 日志时间
	level       LogLevel  // 日志级别
	message     string    // 日志消息
	funcName    string    // 调用函数名
	fileName    string    // 文件名
	line        int       // 行号
	goroutineID int64     // 协程ID
}

// 定义日志格式
var logFormatMap = map[LogFormatType]string{
	Json:     `{"time":"%s","level":"%s","file":"%s","function":"%s","line":"%d", "thread":"%d","message":"%s"}`, // Json格式
	Detailed: `%s | %-7s | %s:%s:%d - %s`,                                                                        // 详细格式
	Bracket:  `[%s] %s`,                                                                                          // 方括号格式
	Threaded: `%s | %-7s | [thread="%d"] %s`,                                                                     // 协程格式
	Simple:   `%s | %-7s | %s`,                                                                                   // 简约格式                                                                                                // 自定义格式
}

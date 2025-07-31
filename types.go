/*
types.go - 日志系统核心类型定义
定义FastLog的核心数据结构、常量和枚举类型，包括日志级别、日志格式、路径信息和日志消息结构体等。
*/
package fastlog

import (
	"os"
	"time"
)

// PathInfo 是一个结构体，用于封装路径的信息
type PathInfo struct {
	Path    string      // 路径
	Exists  bool        // 是否存在
	IsFile  bool        // 是否为文件
	IsDir   bool        // 是否为目录
	Size    int64       // 文件大小（字节）
	Mode    os.FileMode // 文件权限
	ModTime time.Time   // 文件修改时间
}

// 定义缓冲区相关常量
const (
	// 1 KB
	kb = 1024

	// 文件缓冲区配置（更大的缓冲区用于文件写入）
	fileInitialBufferCapacity = 32 * kb                        // 32KB 文件缓冲区初始容量
	fileMaxBufferCapacity     = 1024 * kb                      // 1MB 文件缓冲区最大容量
	fileFlushThreshold        = fileMaxBufferCapacity * 9 / 10 // 文件缓冲区90%阈值

	// 控制台缓冲区配置（较小的缓冲区用于控制台输出）
	consoleInitialBufferCapacity = 8 * kb                            // 8KB 控制台缓冲区初始容量
	consoleMaxBufferCapacity     = 64 * kb                           // 64KB 控制台缓冲区最大容量
	consoleFlushThreshold        = consoleMaxBufferCapacity * 9 / 10 // 控制台缓冲区90%阈值

	// 默认批量处理大小
	defaultBatchSize = 1000
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
type LogLevel uint8

// 将日志级别转换为字符串
func (l LogLevel) MarshalJSON() ([]byte, error) {
	return []byte(`"` + logLevelToString(l) + `"`), nil
}

// 定义日志级别
const (
	DEBUG   LogLevel = 10  // 调试级别
	INFO    LogLevel = 20  // 信息级别
	SUCCESS LogLevel = 30  // 成功级别
	WARN    LogLevel = 40  // 警告级别
	ERROR   LogLevel = 50  // 错误级别
	FATAL   LogLevel = 60  // 致命级别
	NONE    LogLevel = 255 // 无日志级别
)

// logMessage 结构体用于封装日志消息
type logMessage struct {
	Timestamp   string   `json:"time"`     // 预格式化的时间字符串
	Level       LogLevel `json:"level"`    // 日志级别
	Message     string   `json:"message"`  // 日志消息
	FuncName    string   `json:"function"` // 调用函数名
	FileName    string   `json:"file"`     // 文件名
	Line        uint16   `json:"line"`     // 行号
	GoroutineID uint32   `json:"thread"`   // 协程ID
}

// 定义日志格式
var logFormatMap = map[LogFormatType]string{
	Json:     `{"time":"%s","level":"%s","file":"%s","function":"%s","line":"%d", "thread":"%d","message":"%s"}`, // Json格式
	Detailed: `%s | %-7s | %s:%s:%d - %s`,                                                                        // 详细格式
	Bracket:  `[%s] %s`,                                                                                          // 方括号格式
	Threaded: `%s | %-7s | [thread="%d"] %s`,                                                                     // 协程格式
	Simple:   `%s | %-7s | %s`,                                                                                   // 简约格式                                                                                                // 自定义格式
}

// 文件名验证相关常量
var (
	// Windows和Unix系统中文件名不能包含的字符
	invalidFileChars = []string{
		"<", ">", ":", "\"", "|", "?", "*",
		"\x00", "\x01", "\x02", "\x03", "\x04", "\x05", "\x06", "\x07",
		"\x08", "\x09", "\x0a", "\x0b", "\x0c", "\x0d", "\x0e", "\x0f",
		"\x10", "\x11", "\x12", "\x13", "\x14", "\x15", "\x16", "\x17",
		"\x18", "\x19", "\x1a", "\x1b", "\x1c", "\x1d", "\x1e", "\x1f",
	}
	maxFileNameLength = 255  // 大多数文件系统的文件名长度限制
	maxPathLength     = 4096 // 大多数系统的路径长度限制
)

/*
types.go - 日志系统核心类型定义
定义FastLog的核心数据结构、常量和枚举类型，包括日志级别、日志格式、路径信息和日志消息结构体等。
*/
package fastlog

import (
	"os"
	"sync"
	"time"
)

// 预构建的日志级别到字符串的映射表（不带填充，用于JSON序列化）
var logLevelStringMap = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	SUCCESS: "SUCCESS",
	WARN:    "WARN",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
	NONE:    "NONE",
}

// 预构建的日志级别到字符串的映射表（带填充，用于文本格式化）
var logLevelPaddedStringMap = map[LogLevel]string{
	DEBUG:   "DEBUG  ", // 7个字符(预填充空格)
	INFO:    "INFO   ", // 7个字符
	SUCCESS: "SUCCESS", // 7个字符
	WARN:    "WARN   ", // 7个字符
	ERROR:   "ERROR  ", // 7个字符
	FATAL:   "FATAL  ", // 7个字符
	NONE:    "NONE   ", // 7个字符
}

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

// logMsg 结构体用于封装日志消息
type logMsg struct {
	Timestamp string   `json:"time"`     // 预格式化的时间字符串
	Level     LogLevel `json:"level"`    // 日志级别
	FileName  string   `json:"file"`     // 文件名
	FuncName  string   `json:"function"` // 调用函数名
	Line      uint16   `json:"line"`     // 行号
	Message   string   `json:"message"`  // 日志消息
}

// logMsgPool 是一个日志消息对象池
var logMsgPool = sync.Pool{
	New: func() interface{} {
		return &logMsg{}
	},
}

// getLogMsg 获取日志消息对象，使用安全的类型断言
//
// 返回：
//   - *logMsg: 日志消息对象指针，保证非nil
//   - 注意：返回的对象总是可以安全地传递给putLogMsg
func getLogMsg() *logMsg {
	// 尝试从对象池获取对象
	if poolObj := logMsgPool.Get(); poolObj != nil {
		if msg, ok := poolObj.(*logMsg); ok && msg != nil {
			// 重置对象状态，确保干净的对象
			*msg = logMsg{}
			return msg
		}
	}

	// 对象池异常或类型断言失败时的fallback
	// 创建新对象确保函数永远不返回nil
	return &logMsg{}
}

// putLogMsg 归还日志消息对象
//
// 参数：
//   - msg: 日志消息对象指针
//   - 注意：该函数可以安全地处理任何来源的logMsg对象，
//     包括从getLogMsg获取的对象和通过new/&logMsg{}创建的对象
func putLogMsg(msg *logMsg) {
	// 安全检查：防止空指针
	if msg == nil {
		return
	}

	// 使用零值重置，确保完全清理所有字段
	// 这种方式比逐个字段清理更安全，不会遗漏任何字段
	*msg = logMsg{}

	// 归还对象到池中
	// 注意：sync.Pool设计上允许放入任何类型匹配的对象，
	// 无论该对象是否最初从池中获取，都不会导致错误
	logMsgPool.Put(msg)
}

// 日志格式选项
type LogFormatType int

// 日志格式选项
const (
	Detailed   LogFormatType = iota // 详细格式
	Json                            // json格式
	Simple                          // 简约格式
	Structured                      // 结构化格式
	Custom                          // 自定义格式
	//ExtendedStructured                      // 可扩展结构化格式
)

// 定义日志格式
var logFormatMap = map[LogFormatType]string{
	Json:       `{"time":"%s","level":"%s","file":"%s","function":"%s","line":"%d","message":"%s"}`, // Json格式
	Detailed:   `%s | %-7s | %s:%s:%d - %s`,                                                         // 详细格式
	Simple:     `%s | %-7s | %s`,                                                                    // 简约格式                                                                                                // 自定义格式
	Structured: `T:%s|L:%-7s|F:%s:%s:%d|M:%s`,                                                       // 结构化格式
	// ExtendedStructured: `T:%s|L:%-7s|%s|M:%s`,                                                                         // 可扩展结构化格式
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

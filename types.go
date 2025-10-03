/*
types.go - 日志系统核心类型定义
定义FastLog的核心数据结构、常量和枚举类型，包括日志级别、日志格式、路径信息和日志消息结构体等。
*/
package fastlog

import (
	"sync"
	"time"
)

// 预构建的日志级别到字符串的映射表（不带填充，用于JSON序列化）
var logLevelStringMap = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
	NONE:  "NONE",
}

// 预构建的日志级别到字符串的映射表（带填充，用于文本格式化）
var logLevelPaddedStringMap = map[LogLevel]string{
	DEBUG: "DEBUG", // 5个字符(预填充空格)
	INFO:  "INFO ", // 5个字符
	WARN:  "WARN ", // 5个字符
	ERROR: "ERROR", // 5个字符
	FATAL: "FATAL", // 5个字符
	NONE:  "NONE ", // 5个字符
}

const (
	// 文件大小配置常量
	defaultMaxFileSize = 10 // 默认最大文件大小（MB）

	// 默认文件写入器配置
	defaultMaxBufferSize = 64 * 1024       // 默认最大缓冲区大小（64KB）
	defaultMaxWriteCount = 500             // 默认最大写入次数（500次）
	defaultFlushInterval = 1 * time.Second // 默认最大刷新间隔（1秒）
)

// 日志级别枚举
type LogLevel uint8

// 将日志级别转换为字符串
func (l LogLevel) MarshalJSON() ([]byte, error) {
	return []byte(`"` + logLevelToString(l) + `"`), nil
}

// logLevelToString 将 LogLevel 转换为对应的字符串（不带填充，用于JSON序列化）
//
// 参数：
//   - level: 要转换的日志级别
//
// 返回值：
//   - string: 对应的日志级别字符串, 如果 level 无效, 则返回 "UNKNOWN"
func logLevelToString(level LogLevel) string {
	// 使用预构建的映射表进行O(1)查询(不带填充，适用于JSON)
	if str, exists := logLevelStringMap[level]; exists {
		return str
	}
	return "UNKNOWN"
}

// 定义日志级别
const (
	DEBUG LogLevel = 10  // 调试级别
	INFO  LogLevel = 20  // 信息级别
	WARN  LogLevel = 30  // 警告级别
	ERROR LogLevel = 40  // 错误级别
	FATAL LogLevel = 50  // 致命级别
	NONE  LogLevel = 255 // 无日志级别
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
	// 尝试从对象池获取对象并进行类型断言
	if msg, ok := logMsgPool.Get().(*logMsg); ok {
		return msg
	}

	// 创建新的对象
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
	logMsgPool.Put(msg)
}

// simpleLogMsg 简化的日志消息结构（用于JsonSimple格式）
type simpleLogMsg struct {
	Timestamp string   `json:"time"`    // 预格式化的时间字符串
	Level     LogLevel `json:"level"`   // 日志级别
	Message   string   `json:"message"` // 日志消息
}

// 日志格式选项
type LogFormatType int

// 日志格式选项
const (
	// 详细格式
	Detailed LogFormatType = iota

	// json格式
	Json

	// json简化格式(无文件信息)
	JsonSimple

	// 简约格式(无文件信息)
	Simple

	// 结构化格式
	Structured

	// 基础结构化格式(无文件信息)
	BasicStructured

	// 简单时间格式(无文件信息)
	SimpleTimestamp

	// 自定义格式(无文件信息)
	Custom
)

// fileInfoRequiredFormats 需要处理文件信息的日志格式集合
// 如果日志格式需要文件信息(文件名、函数名、行号)，则将其添加到该集合中
var fileInfoRequiredFormats = map[LogFormatType]struct{}{
	Json:       {},
	Detailed:   {},
	Structured: {},
}

// 定义日志格式
var logFormatMap = map[LogFormatType]string{
	Json:            `{"time":"%s","level":"%s","file":"%s","function":"%s","line":"%d","message":"%s"}`, // Json格式
	JsonSimple:      `{"time":"%s","level":"%s","message":"%s"}`,                                         // Json简化格式(无文件信息)
	Detailed:        `%s | %-6s | %s:%s:%d - %s`,                                                         // 详细格式
	Simple:          `%s | %-6s | %s`,                                                                    // 简约格式                                                                                                // 自定义格式
	Structured:      `T:%s|L:%-6s|F:%s:%s:%d|M:%s`,                                                       // 结构化格式
	BasicStructured: `T:%s|L:%-6s|M:%s`,                                                                  // 基础结构化格式(无文件信息)
	SimpleTimestamp: `%s %s %s`,                                                                          // 简单时间格式
}

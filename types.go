/*
types.go - 日志系统核心类型定义
定义FastLog的核心数据结构、常量和枚举类型，包括日志级别、日志格式、路径信息和日志消息结构体等。
*/
package fastlog

import (
	"context"
	"io"
	"sync"
	"time"

	"gitee.com/MM-Q/colorlib"
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
	DEBUG: "DEBUG ", // 6个字符(预填充空格)
	INFO:  "INFO  ", // 6个字符
	WARN:  "WARN  ", // 6个字符
	ERROR: "ERROR ", // 6个字符
	FATAL: "FATAL ", // 6个字符
	NONE:  "NONE  ", // 6个字符
}

const (
	// 默认批量处理大小
	defaultBatchSize = 1000

	// 每条日志的估算字节数常量
	bytesPerLogEntry = 256

	// 通道大小配置常量
	defaultChanSize = 10000 // 默认通道大小
	largeChanSize   = 20000 // 大通道大小（生产/静默模式）
	smallChanSize   = 5000  // 小通道大小（控制台模式）

	// 刷新间隔配置常量
	fastFlushInterval   = 100 * time.Millisecond  // 快速刷新（开发/控制台模式）
	normalFlushInterval = 500 * time.Millisecond  // 正常刷新（默认/文件模式）
	slowFlushInterval   = 1000 * time.Millisecond // 慢速刷新（生产/静默模式）

	// 文件大小配置常量
	defaultMaxFileSize = 10 // 默认最大文件大小（MB）

	// 文件保留配置常量
	developmentMaxAge     = 7   // 开发模式保留天数
	developmentMaxBackups = 20  // 开发模式保留文件数
	productionMaxAge      = 30  // 生产模式保留天数
	productionMaxBackups  = 50  // 生产模式保留文件数
	fileMaxAge            = 14  // 文件模式保留天数
	fileMaxBackups        = 30  // 文件模式保留文件数
	silentMaxAge          = 30  // 静默模式保留天数
	silentMaxBackups      = 100 // 静默模式保留文件数

	// 系统资源限制常量
	maxChanSize          = 1000000                // 最大通道大小限制
	maxSingleFileSize    = 10000                  // 最大单文件大小限制（MB）
	minFlushInterval     = time.Microsecond       // 最小刷新间隔
	maxFlushInterval     = 30 * time.Second       // 最大刷新间隔
	normalMinFlush       = 10 * time.Millisecond  // 正常最小刷新间隔
	maxBatchSize         = 5000                   // 最大批处理大小
	chanSizeLimit        = 100000                 // 通道大小上限
	maxRetentionDays     = 3650                   // 最大保留天数（10年）
	maxRetentionFiles    = 1000                   // 最大保留文件数
	minRetentionDays     = 7                      // 最小保留天数
	minRetentionFiles    = 5                      // 最小保留文件数
	performanceThreshold = 50000                  // 性能阈值
	performanceFlushMin  = 100 * time.Millisecond // 性能模式最小刷新间隔
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
	Detailed:        `%s | %-7s | %s:%s:%d - %s`,                                                         // 详细格式
	Simple:          `%s | %-7s | %s`,                                                                    // 简约格式                                                                                                // 自定义格式
	Structured:      `T:%s|L:%-7s|F:%s:%s:%d|M:%s`,                                                       // 结构化格式
	BasicStructured: `T:%s|L:%-7s|M:%s`,                                                                  // 基础结构化格式(无文件信息)
	SimpleTimestamp: `%s %s %s`,                                                                          // 简单时间格式
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

	// 文件处理字符串常量
	defaultLogDir      = "logs"       // 默认日志目录名
	defaultLogFileName = "app.log"    // 默认日志文件名
	truncatedSuffix    = "_truncated" // 路径截断后缀
	charReplacement    = "_"          // 非法字符替换符
	truncateReserve    = 10           // 路径截断预留长度
)

// processorDependencies 定义处理器所需的最小依赖接口
// 通过接口隔离原则，processor 只能访问必要的功能，避免持有完整的 FastLog 引用
type processorDependencies interface {
	// getConfig 获取日志配置
	getConfig() *FastLogConfig

	// getFileWriter 获取文件写入器
	getFileWriter() io.Writer

	// getConsoleWriter 获取控制台写入器
	getConsoleWriter() io.Writer

	// getColorLib 获取颜色库实例
	getColorLib() *colorlib.ColorLib

	// getContext 获取上下文，用于控制处理器生命周期
	getContext() context.Context

	// getLogChannel 获取日志消息通道
	getLogChannel() <-chan *logMsg

	// notifyProcessorDone 通知处理器完成工作
	notifyProcessorDone()

	// getBufferSize 获取缓冲区大小
	getBufferSize() int
}

// WriterPair 写入器对，用于批量传递写入器
type WriterPair struct {
	FileWriter    io.Writer
	ConsoleWriter io.Writer
}

// ProcessorConfig 处理器配置结构
type ProcessorConfig struct {
	BatchSize     int           // 批量处理大小
	FlushInterval time.Duration // 刷新间隔
}

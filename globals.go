// 定义全局常量变量或结构体
package fastlog

import (
	"bytes"
	"context"
	"io"
	"strings"
	"sync"
	"time"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/logrotatex"
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
	None    LogLevel = 999 // 无日志级别
)

// 日志记录器
type FastLog struct {
	/*  私有属性 内部使用无需修改  */
	logFilePath    string             // 日志文件路径  内部拼接的 [logDirName+logFileName]
	logChan        chan *logMessage   // 日志通道  用于异步写入日志文件
	logWait        sync.WaitGroup     // 等待组 用于等待所有goroutine完成
	fileWriter     io.Writer          // 文件写入器
	fileMu         sync.Mutex         // 文件锁 用于保护文件缓冲区的写入操作
	consoleMu      sync.Mutex         // 控制台锁 用于保护控制台缓冲区的写入操作
	consoleWriter  io.Writer          // 控制台写入器
	startOnce      sync.Once          // 用于确保日志处理器只启动一次
	fileBuffer     *bytes.Buffer      // 文件缓冲区 用于存储待写入的日志消息
	consoleBuffer  *bytes.Buffer      // 控制台缓冲区 用于存储待写入的日志消息
	fileBuilder    strings.Builder    // 文件构建器 用于构建待写入的日志消息
	consoleBuilder strings.Builder    // 控制台构建器 用于构建待写入的日志消息
	ctx            context.Context    // 控制刷新器的上下文
	cancel         context.CancelFunc // 控制刷新器的取消函数
	cl             *colorlib.ColorLib // 提供终端颜色输出的库

	/* logrotatex 日志文件切割 */
	logGer *logrotatex.LogRotateX // 日志文件切割器

	/* 嵌入的配置结构体 */
	config *FastLogConfig // 配置结构体
}

// 定义一个配置结构体，用于配置日志记录器
type FastLogConfig struct {
	LogDirName     string        // 日志目录路径
	LogFileName    string        // 日志文件名
	PrintToConsole bool          // 是否将日志输出到控制台
	ConsoleOnly    bool          // 是否仅输出到控制台
	FlushInterval  time.Duration // 刷新间隔，单位为秒
	LogLevel       LogLevel      // 日志级别
	ChanIntSize    int           // 通道大小 默认1000
	LogFormat      LogFormatType // 日志格式选项
	MaxBufferSize  int           // 最大缓冲区大小, 单位为MB, 默认1MB
	NoColor        bool          // 是否禁用终端颜色
	NoBold         bool          // 是否禁用终端字体加粗
	MaxLogFileSize int           // 最大日志文件大小, 单位为MB, 默认10MB
	MaxLogAge      int           // 最大日志文件保留天数, 默认为0, 表示不做限制
	MaxLogBackups  int           // 最大日志文件保留数量, 默认为0, 表示不做限制
	IsLocalTime    bool          // 是否使用本地时间 默认使用UTC时间
	EnableCompress bool          // 是否启用日志文件压缩 默认不启用
}

// 定义一个接口, 声明对外暴露的方法
type FastLogInterface interface {
	Info(v ...any)                    // 记录信息级别的日志，不支持占位符
	Warn(v ...any)                    // 记录警告级别的日志，不支持占位符
	Error(v ...any)                   // 记录错误级别的日志，不支持占位符
	Success(v ...any)                 // 记录成功级别的日志，不支持占位符
	Debug(v ...any)                   // 记录调试级别的日志，不支持占位符
	Close()                           // 关闭日志记录器
	Infof(format string, v ...any)    // 记录信息级别的日志，支持占位符，格式化
	Warnf(format string, v ...any)    // 记录警告级别的日志，支持占位符，格式化
	Errorf(format string, v ...any)   // 记录错误级别的日志，支持占位符，格式化
	Successf(format string, v ...any) // 记录成功级别的日志，支持占位符，格式化
	Debugf(format string, v ...any)   // 记录调试级别的日志，支持占位符，格式化
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

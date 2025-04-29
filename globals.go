// 定义全局常量变量或结构体
package fastlog

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"gitee.com/MM-Q/colorlib"
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
	logFile        *os.File           // 日志文件句柄
	logFilePath    string             // 日志文件路径  内部拼接的 [logDirName+logFileName] 绝对路径
	logChan        chan *logMessage   // 日志通道  用于异步写入日志文件
	logWait        sync.WaitGroup     // 等待组 用于等待所有goroutine完成
	fileMu         sync.Mutex         // 文件锁 用于保护文件写入操作的并发安全
	consoleMu      sync.Mutex         // 控制台锁 用于保护控制台写入操作的并发安全
	fileWriter     io.Writer          // 文件写入器
	consoleWriter  io.Writer          // 控制台写入器
	startOnce      sync.Once          // 用于确保日志处理器只启动一次
	fileBuffer     *bytes.Buffer      // 文件缓冲区 用于存储待写入的日志消息
	consoleBuffer  *bytes.Buffer      // 控制台缓冲区 用于存储待写入的日志消息
	flushInterval  time.Duration      // 刷新间隔，单位为秒
	fileBuilder    strings.Builder    // 文件构建器 用于构建待写入的日志消息
	consoleBuilder strings.Builder    // 控制台构建器 用于构建待写入的日志消息
	ctx            context.Context    // 控制刷新器的上下文
	cancel         context.CancelFunc // 控制刷新器的取消函数

	/*  测试属性  */
	maxLogFileSize              int64         // 单个日志文件的最大大小，单位为字节
	maxLogFileHour              time.Duration // 日志文件的最大保留小时数
	rotationCheckIntervalSecond time.Duration // 定时检查日志轮转的间隔时间(秒)
	currentLogFileSize          int64         // 当前日志文件的大小，单位为字节
	rotationEnabled             bool          // 是否启用日志轮转功能

	/*  公共属性 可以通过属性自定义配置  */
	logDirName     string        // 日志目录路径
	logFileName    string        // 日志文件名
	printToConsole bool          // 是否将日志输出到控制台
	consoleOnly    bool          // 是否仅输出到控制台
	logLevel       LogLevel      // 日志级别
	chanIntSize    int           // 通道大小 默认1000
	logFormat      LogFormatType // 日志格式选项
	maxBufferSize  int           // 最大缓冲区大小，单位为MB，默认为1MB
}

// 获取一个新的ColorLib实例
var CL = colorlib.NewColorLib()

// 定义一个配置结构体，用于配置日志记录器
type FastLogConfig struct {
	logDirName                  string        // 日志目录路径
	LogFileName                 string        // 日志文件名
	logFilePath                 string        // 日志文件路径  内部拼接的 [logDirName+logFileName]
	PrintToConsole              bool          // 是否将日志输出到控制台
	ConsoleOnly                 bool          // 是否仅输出到控制台
	LogLevel                    LogLevel      // 日志级别
	ChanIntSize                 int           // 通道大小 默认1000
	LogFormat                   LogFormatType // 日志格式选项
	MaxBufferSize               int           // 最大缓冲区大小
	MaxLogFileSize              int           // 单个日志文件的最大大小，单位为MB
	MaxLogFileHour              int           // 日志文件的最大保留小时数
	RotationCheckIntervalSecond int           // 定时检查日志轮转的间隔时间(秒)
	RotationEnabled             bool          // 是否启用日志轮转功能，默认为false
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

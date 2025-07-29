// 定义全局常量变量或结构体
package fastlog

import (
	"bytes"
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/logrotatex"
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

// 日志记录器
type FastLog struct {
	// 日志文件路径  内部拼接的 [logDirName+logFileName]
	logFilePath string
	// 日志通道  用于异步写入日志文件
	logChan chan *logMessage
	// 等待组 用于等待所有goroutine完成
	logWait sync.WaitGroup
	// 文件写入器
	fileWriter io.Writer
	// 文件锁 用于保护文件缓冲区的写入操作
	fileMu sync.Mutex
	// 控制台锁 用于保护控制台缓冲区的写入操作
	consoleMu sync.Mutex
	// 控制台写入器
	consoleWriter io.Writer
	// 用于确保日志处理器只启动一次
	startOnce sync.Once
	// 控制刷新器的上下文
	ctx context.Context
	// 控制刷新器的取消函数
	cancel context.CancelFunc
	// 提供终端颜色输出的库
	cl *colorlib.ColorLib
	// 用于确保日志处理器只关闭一次
	closeOnce sync.Once

	// 文件双缓冲区
	fileBuffers [2]*bytes.Buffer
	// 当前使用的文件缓冲区索引
	fileBufferIdx atomic.Int32
	// 控制台双缓冲区
	consoleBuffers [2]*bytes.Buffer
	// 当前使用的控制台缓冲区索引
	consoleBufferIdx atomic.Int32
	// 文件缓冲区锁
	fileBufferMu sync.Mutex
	// 控制台缓冲区锁
	consoleBufferMu sync.Mutex
	// 用于控制缓冲区刷新的锁
	flushLock sync.Mutex

	// 用于控制关闭过程的锁
	closeLock sync.Mutex

	// logrotatex 日志文件切割
	logGer *logrotatex.LogRotateX

	// 嵌入的配置结构体
	config *FastLogConfig
}

// 定义一个配置结构体，用于配置日志记录器
type FastLogConfig struct {
	logDirName     string        // 日志目录路径
	logFileName    string        // 日志文件名
	printToConsole bool          // 是否将日志输出到控制台
	consoleOnly    bool          // 是否仅输出到控制台
	flushInterval  time.Duration // 刷新间隔，单位为time.Duration
	logLevel       LogLevel      // 日志级别
	chanIntSize    int           // 通道大小 默认10000
	logFormat      LogFormatType // 日志格式选项
	maxBufferSize  int           // 最大缓冲区大小, 单位为MB, 默认1MB
	noColor        bool          // 是否禁用终端颜色
	noBold         bool          // 是否禁用终端字体加粗
	maxLogFileSize int           // 最大日志文件大小, 单位为MB, 默认5MB
	maxLogAge      int           // 最大日志文件保留天数, 默认为0, 表示不做限制
	maxLogBackups  int           // 最大日志文件保留数量, 默认为0, 表示不做限制
	isLocalTime    bool          // 是否使用本地时间 默认使用UTC时间
	enableCompress bool          // 是否启用日志文件压缩 默认不启用
	setMu          sync.RWMutex  // 用于保护配置的锁 (读写锁优化并发读性能)
}

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

// FastLogConfigurer 定义日志配置器接口，包含所有配置项的设置和获取方法
type FastLogConfigurer interface {
	// SetLogDirName 设置日志目录路径
	SetLogDirName(dirName string)
	// GetLogDirName 获取日志目录路径
	GetLogDirName() string

	// SetLogFileName 设置日志文件名
	SetLogFileName(fileName string)
	// GetLogFileName 获取日志文件名
	GetLogFileName() string

	// SetPrintToConsole 设置是否将日志输出到控制台
	SetPrintToConsole(print bool)
	// GetPrintToConsole 获取是否将日志输出到控制台的状态
	GetPrintToConsole() bool

	// SetConsoleOnly 设置是否仅输出到控制台
	SetConsoleOnly(only bool)
	// GetConsoleOnly 获取是否仅输出到控制台的状态
	GetConsoleOnly() bool

	// SetFlushInterval 设置刷新间隔
	SetFlushInterval(interval time.Duration)
	// GetFlushInterval 获取刷新间隔
	GetFlushInterval() time.Duration

	// SetLogLevel 设置日志级别
	SetLogLevel(level LogLevel)
	// GetLogLevel 获取日志级别
	GetLogLevel() LogLevel

	// SetChanIntSize 设置通道大小
	SetChanIntSize(size int)
	// GetChanIntSize 获取通道大小
	GetChanIntSize() int

	// SetLogFormat 设置日志格式选项
	SetLogFormat(format LogFormatType)
	// GetLogFormat 获取日志格式选项
	GetLogFormat() LogFormatType

	// SetNoColor 设置是否禁用终端颜色
	SetNoColor(noColor bool)
	// GetNoColor 获取是否禁用终端颜色的状态
	GetNoColor() bool

	// SetNoBold 设置是否禁用终端字体加粗
	SetNoBold(noBold bool)
	// GetNoBold 获取是否禁用终端字体加粗的状态
	GetNoBold() bool

	// SetMaxLogFileSize 设置最大日志文件大小(MB)
	SetMaxLogFileSize(size int)
	// GetMaxLogFileSize 获取最大日志文件大小(MB)
	GetMaxLogFileSize() int

	// SetMaxLogAge 设置最大日志文件保留天数
	SetMaxLogAge(age int)
	// GetMaxLogAge 获取最大日志文件保留天数
	GetMaxLogAge() int

	// SetMaxLogBackups 设置最大日志文件保留数量
	SetMaxLogBackups(backups int)
	// GetMaxLogBackups 获取最大日志文件保留数量
	GetMaxLogBackups() int

	// SetIsLocalTime 设置是否使用本地时间
	SetIsLocalTime(local bool)
	// GetIsLocalTime 获取是否使用本地时间的状态
	GetIsLocalTime() bool

	// SetEnableCompress 设置是否启用日志文件压缩
	SetEnableCompress(compress bool)
	// GetEnableCompress 获取是否启用日志文件压缩的状态
	GetEnableCompress() bool
}

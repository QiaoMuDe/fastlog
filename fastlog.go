package fastlog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"path/filepath"

	"gitee.com/MM-Q/colorlib"
)

// 日志级别枚举
type LogLevel int

// 定义一个原子计数器
var counter int32 = 0

type LogFormatType int

const (
	// 日志格式选项
	Detailed LogFormatType = iota // 详细格式
	Bracket                       // 方括号格式
	Json                          // json格式
	Threaded                      // 协程格式
)

// 定义日志级别
const (
	Debug LogLevel = iota
	Info
	Warn
	Error
	Success
	None // 用于表示不记录任何日志
)

// 日志记录器
type Logger struct {
	logDirName        string         // 日志目录
	logFileName       string         // 日志文件名
	logFile           *os.File       // 日志文件句柄
	logFilePath       string         // 日志文件路径  内部拼接的 [logDirName+logFileName]
	logger            *log.Logger    // 底层日志记录器
	printToConsole    bool           // 是否将日志输出到控制台
	consoleOnly       bool           // 是否仅输出到控制台
	logLevel          LogLevel       // 日志级别
	fileMu            sync.Mutex     // 文件写入的互斥锁
	consoleMu         sync.Mutex     // 控制台写入的互斥锁
	logChan           chan string    // 日志通道
	stopChan          chan string    // 停止通道
	wg                sync.WaitGroup // 等待组，用于等待日志通道中的日志被处理
	fileBuffer        *bytes.Buffer  // 文件缓冲区
	consoleBuffer     *bytes.Buffer  // 控制台缓冲区
	ticker            *time.Ticker   // 定时器
	fileWriter        io.Writer      // 文件写入器
	consoleWriter     io.Writer      // 控制台写入器
	chanIntSize       int            // 通道大小
	bufferKbSize      int            // 缓冲区大小
	closeOnce         sync.Once      // 确保通道只被关闭一次
	logFormat         LogFormatType  // 日志格式选项 [Json(json格式)|Bracket(方括号格式)|Detailed(详细格式)|Threaded(协程格式)]
	enableLogRotation bool           // 是否启用日志切割 [true|false] 默认false
	logRetentionDays  int            // 日志保留天数 默认7天 单位[天]
	logMaxSize        string         // 日志文件最大大小 默认3MB 单位[MB|GB]
	logSizeBytes      int64          // 日志文件大小字节数 用于比较大小, 无需提供或设置修改
	logRetentionCount int            // 日志文件保留数量 默认3 单位[个]
	rotationInterval  int            // 日志轮转的间隔时间 默认10分钟 单位[分钟]
	enableCompression bool           // 是否启用日志压缩 [true|false]
	compressionFormat string         // 日志压缩格式 ["zip"|"gzip"|"bz2"|"xz"|"tar"|"tgz"] 默认zip
}

// 日志配置
type LoggerConfig struct {
	LogDirName        string        // 日志目录名称
	LogFileName       string        // 日志文件名称
	LogPath           string        // 日志文件路径
	PrintToConsole    bool          // 是否将日志输出到控制台
	ConsoleOnly       bool          // 是否仅输出到控制台
	LogLevel          LogLevel      // 日志级别
	ChanIntSize       int           // 通道大小
	BufferKbSize      int           // 缓冲区大小
	LogFormat         LogFormatType // 日志格式选项 [Json(json格式)|Bracket(括号格式)|Detailed(详细格式)|Threaded(协程格式)]
	EnableLogRotation bool          // 是否启用日志切割 [true|false] 默认false
	LogRetentionDays  int           // 日志保留天数 默认7天 单位[天]
	LogMaxSize        string        // 日志文件最大大小 默认3MB 单位[MB|GB]
	LogRetentionCount int           // 日志文件保留数量 默认3 单位[个]
	EnableCompression bool          // 是否启用日志压缩 [true|false]
	RotationInterval  int           // 日志轮转的间隔时间 默认10分钟 单位[分钟]
	CompressionFormat string        // 日志压缩格式 [zip|gzip|bz2|xz|tar|tgz] 默认zip
}

// 全局编译正则表达式 用于匹配日志文件中的时间戳
var timestampRegex = regexp.MustCompile(`_(\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})\.log`)

// 定义一个接口, 声明对外暴露的方法
type LoggerInterface interface {
	Info(v ...interface{})                    // 记录信息级别的日志，不支持占位符，需要自己拼接
	Warn(v ...interface{})                    // 记录警告级别的日志，不支持占位符，需要自己拼接
	Error(v ...interface{})                   // 记录错误级别的日志，不支持占位符，需要自己拼接
	Success(v ...interface{})                 // 记录成功级别的日志，不支持占位符，需要自己拼接
	Debug(v ...interface{})                   // 记录调试级别的日志，不支持占位符，需要自己拼接
	Close()                                   // 关闭日志记录器
	Infof(format string, v ...interface{})    // 记录信息级别的日志，支持占位符，格式化
	Warnf(format string, v ...interface{})    // 记录警告级别的日志，支持占位符，格式化
	Errorf(format string, v ...interface{})   // 记录错误级别的日志，支持占位符，格式化
	Successf(format string, v ...interface{}) // 记录成功级别的日志，支持占位符，格式化
	Debugf(format string, v ...interface{})   // 记录调试级别的日志，支持占位符，格式化
}

// 创建一个colorlib的实例
var color = colorlib.NewColorLib()

// 创建一个新的日志配置器
func NewConfig(logDirName string, logFileName string) LoggerConfig {
	// 生成日志文件路径
	logPath := filepath.Join(logDirName, logFileName)
	return LoggerConfig{
		LogDirName:        logDirName,  // 日志目录 必须提供
		LogFileName:       logFileName, // 日志文件名 必须提供
		LogPath:           logPath,     // 内部自己拼接的路径,无需提供
		PrintToConsole:    true,        // 是否将日志输出到控制台, 默认值为true
		ConsoleOnly:       false,       // 是否仅输出到控制台, 默认值为false
		LogLevel:          Info,        // 日志过滤级别, 默认值为Info
		ChanIntSize:       1000,        // 日志通道大小, 默认1000
		BufferKbSize:      1024,        // 缓冲区大小, 默认1MB
		LogFormat:         Detailed,    // 日志格式选项 [Json(json格式)|Bracket(括号格式)|Detailed(详细格式)|Threaded(协程格式)] 默认详细格式
		EnableLogRotation: false,       // 是否启用日志切割 [true|false] 默认false
		LogRetentionDays:  7,           // 日志保留天数 默认7天 单位[天]
		LogMaxSize:        "3MB",       // 日志文件最大大小 默认3MB 单位[MB|GB]
		LogRetentionCount: 3,           // 日志文件保留数量 默认3 单位[个]
		RotationInterval:  10,          // 日志轮转的间隔时间 默认10分钟 单位[分钟]
		EnableCompression: false,       // 是否启用日志压缩 [true|false] 默认false
		CompressionFormat: "zip",       // 日志压缩格式 [zip|gzip|bz2|xz|tar|tgz] 默认zip
	}
}

// 创建一个新的日志记录器
func NewLogger(cfg LoggerConfig) (*Logger, error) {
	// 声明一些配置变量
	var (
		outputToConsole     bool          // 是否将日志输出到控制台
		outputOnlyToConsole bool          // 是否仅输出到控制台
		consoleBuffer       *bytes.Buffer // 初始化控制台缓冲区
		fileBuffer          *bytes.Buffer // 初始化日志文件缓冲区
		logFile             *os.File      // 日志文件句柄
		fileWriter          io.Writer     // 文件写入器
	)

	// 检查是否仅输出到控制台
	if cfg.ConsoleOnly {
		// 如果仅输出到控制台，设置 outputToConsole 为 true，outputOnlyToConsole 为 true
		outputToConsole = true
		outputOnlyToConsole = true
	} else {
		// 如果不是仅输出到控制台：
		// 1. 根据配置决定是否输出到控制台（由 PrintToConsole 控制）
		// 2. 默认情况下，日志会输出到文件
		outputToConsole = cfg.PrintToConsole
		outputOnlyToConsole = false
	}

	// 初始化缓冲区大小 默认1MB
	bufferSize := cfg.BufferKbSize * 1024

	// 如果需要输出到控制台，初始化控制台缓冲区
	if outputToConsole {
		consoleBuffer = bytes.NewBuffer(make([]byte, bufferSize))
		consoleBuffer.Reset()
	}

	// 如果不是仅输出到控制台，初始化文件缓冲区和打开日志文件
	if !outputOnlyToConsole {
		// 检查日志目录是否存在
		if err := ensureLogDirExists(cfg.LogDirName); err != nil {
			return nil, err
		}

		// 打开日志文件
		logFile, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil, fmt.Errorf("无法打开日志文件: %w", err)
		}

		// 初始化文件缓冲区
		fileBuffer = bytes.NewBuffer(make([]byte, bufferSize))
		fileBuffer.Reset()

		// 设置文件写入器
		fileWriter = logFile
	} else {
		// 如果仅输出到控制台，文件写入器为 nil
		fileWriter = nil
	}

	// 创建底层日志记录器
	logger := log.New(fileWriter, "", 0)

	// 创建 Logger 实例
	lg := &Logger{
		logDirName:        cfg.LogDirName,        // 保存日志目录名称
		logFileName:       cfg.LogFileName,       // 保存日志文件名称
		logFile:           logFile,               // 保存日志文件句柄
		logFilePath:       cfg.LogPath,           // 保存日志文件路径
		logger:            logger,                // 保存底层日志记录器
		printToConsole:    outputToConsole,       // 保存是否将日志输出到控制台
		consoleOnly:       outputOnlyToConsole,   // 保存是否仅输出到控制台
		logLevel:          cfg.LogLevel,          // 保存日志级别
		fileWriter:        fileWriter,            // 保存文件写入器
		consoleWriter:     os.Stdout,             // 保存控制台写入器
		chanIntSize:       cfg.ChanIntSize,       // 保存通道大小
		bufferKbSize:      cfg.BufferKbSize,      // 保存缓冲区大小
		fileBuffer:        fileBuffer,            // 保存预分配内存后的文件缓冲区
		consoleBuffer:     consoleBuffer,         // 保存预分配内存后的控制台缓冲区
		logFormat:         cfg.LogFormat,         // 保存日志格式选项 [Json(json格式)|Bracket(括号格式)|Detailed(详细格式)|Threaded(协程格式)]
		enableLogRotation: cfg.EnableLogRotation, // 保存是否启用日志切割 [true|false] 默认false,
		logRetentionDays:  cfg.LogRetentionDays,  // 保存日志保留天数 默认7天 单位[天]
		logMaxSize:        cfg.LogMaxSize,        // 保存日志文件最大大小 默认3MB 单位[MB|GB]
		logRetentionCount: cfg.LogRetentionCount, // 保存日志文件保留数量 默认3 单位[个]
		rotationInterval:  cfg.RotationInterval,  // 保存日志轮转的间隔时间 默认10分钟 单位[分钟]
		enableCompression: cfg.EnableCompression, // 保存是否启用日志压缩 [true|false] 默认false
		compressionFormat: cfg.CompressionFormat, // 保存日志压缩格式 [zip|gzip|bz2|xz|tar|tgz] 默认zip
	}

	// 验证日志配置
	_, err := validateLoggerConfig(lg)
	if err != nil {
		return nil, err
	}

	// 初始化日志通道, 容量为1000
	lg.logChan = make(chan string, lg.chanIntSize)

	// 初始化停止通道
	lg.stopChan = make(chan string, 10)

	// 启动异步日志记录 使用 sync.Once 确保只执行一次
	onefunc := sync.OnceFunc(func() {
		lg.startAsyncLogging()
	})
	onefunc() // 执行一次

	// 检查是否启用日志切割
	if lg.enableLogRotation {
		// 启动日志轮转主协程 使用 sync.Once 确保只执行一次
		onefunc := sync.OnceFunc(func() {
			lg.startLogRotate()
		})
		onefunc() // 执行一次
	}

	// 返回 Logger 实例
	return lg, nil
}

// validateLoggerConfig 检查 Logger 结构体的配置是否合理
func validateLoggerConfig(logger *Logger) (bool, error) {
	// 检查日志格式是否支持
	supportedFormats := map[string]LogFormatType{
		"json":     Json,
		"bracket":  Bracket,
		"detailed": Detailed,
		"threaded": Threaded,
	}
	userFormat := strings.ToLower(logger.logFormat.String()) // 转换为小写
	if _, ok := supportedFormats[userFormat]; !ok {          // 获取日志格式
		// 提取支持的格式列表
		var supportedList []string
		for key := range supportedFormats { // 遍历支持的格式
			supportedList = append(supportedList, key)
		}
		return false, fmt.Errorf("不支持的日志格式: %s, 支持的格式: %s", userFormat, strings.Join(supportedList, ", "))
	}

	// 检查是否仅输出到控制台
	if logger.consoleOnly {
		return true, nil // 仅输出到控制台，无需进一步检查
	}

	// 检查日志文件名是否为空
	if logger.logFileName == "" {
		return false, errors.New("日志文件名不能为空")
	}

	// 检查日志文件名是否以 .log 结尾
	if !strings.HasSuffix(logger.logFileName, ".log") {
		return false, errors.New("日志文件名必须以 .log 结尾")
	}

	// 检查日志保留天数是否合理
	if logger.logRetentionDays <= 0 {
		logger.logRetentionDays = 7             // 默认值
		logger.Warn("日志保留天数不能小于等于0, 已调整为默认值7天") // 打印警告信息
	}

	// 检查日志文件最大大小是否合理
	if logger.logMaxSize == "" {
		logger.logMaxSize = "3MB" // 默认值
		logger.Warn("日志文件最大大小未设置, 已调整为默认值3MB")
	}

	// 保存一下日志文件最大大小用于在日志中打印
	printLogSize := logger.logMaxSize

	// 解析日志大小字符串并转换为字节
	sizeInBytes, err := parseSize(logger.logMaxSize)
	if err != nil {
		return false, fmt.Errorf("日志文件最大大小设置不合理，支持格式如 '3MB' 或 '1GB'")
	}
	logger.logSizeBytes = sizeInBytes // 保存解析后的日志大小

	// 检查日志文件保留数量是否合理
	if logger.logRetentionCount <= 0 {
		logger.logRetentionCount = 3 // 默认值
		logger.Warn("日志文件保留数量设置不合理, 已调整为默认值3个")
	}

	// 检查日志轮转的间隔时间是否合理
	if logger.rotationInterval <= 0 {
		logger.rotationInterval = 10 // 默认值
		logger.Warn("日志轮转的间隔时间不能小于等于0, 已调整为默认值10分钟")
	}

	// 检查日志压缩格式是否支持
	if logger.enableCompression { // 如果启用压缩
		if logger.compressionFormat == "" {
			logger.compressionFormat = "zip" // 默认值
			logger.Warn("未设置压缩格式, 已调整为默认值zip")
		} else {
			supportedFormats := []string{"zip", "gzip", "bz2", "xz", "tar", "tgz"}
			if !strings.Contains(strings.Join(supportedFormats, ","), logger.compressionFormat) {
				return false, fmt.Errorf("不支持的压缩格式：%s，支持的格式：%s", logger.compressionFormat, strings.Join(supportedFormats, ","))
			}
		}
	}

	// 配置检查通过，记录日志
	logger.Info("fastlog 配置检查通过!")
	logger.Info(fmt.Sprintf("日志目录: %s", logger.logDirName))
	logger.Info(fmt.Sprintf("日志文件名: %s", logger.logFileName))
	logger.Info(fmt.Sprintf("日志级别: %d", logger.logLevel))
	logger.Info(fmt.Sprintf("日志格式: %s", logger.logFormat.String()))
	logger.Info(fmt.Sprintf("是否启用日志轮转: %t", logger.enableLogRotation))
	logger.Info(fmt.Sprintf("日志保留天数: %d", logger.logRetentionDays))
	logger.Info(fmt.Sprintf("日志文件最大大小: %s", printLogSize)) // 打印日志文件最大大小
	logger.Info(fmt.Sprintf("日志文件保留数量: %d", logger.logRetentionCount))
	logger.Info(fmt.Sprintf("日志轮转间隔时间: %d 分钟", logger.rotationInterval))
	logger.Info(fmt.Sprintf("是否启用压缩: %t", logger.enableCompression))
	logger.Info(fmt.Sprintf("压缩格式: %s", logger.compressionFormat))

	// 返回 true 表示配置检查通过
	return true, nil
}

// parseSize 函数用于解析表示日志大小的字符串，并返回以字节为单位的整数大小。
// 如果解析失败，则返回错误。
func parseSize(sizeStr string) (int64, error) {
	// 将输入的字符串转换为大写，以便统一处理单位。
	sizeStr = strings.ToUpper(sizeStr)
	// 提取字符串的最后两个字符作为单位（例如 "MB" 或 "GB"）。
	unit := sizeStr[len(sizeStr)-2:]
	// 尝试将字符串的前部分转换为整数，表示大小。
	// 这里使用 strconv.ParseInt 函数，并指定基数为 10 和 64 位整数。
	size, err := strconv.ParseInt(sizeStr[:len(sizeStr)-2], 10, 64)
	if err != nil {
		// 如果转换失败，返回错误，并使用 fmt.Errorf 包装原始错误信息。
		return 0, fmt.Errorf("无法解析日志大小: %w", err)
	}

	// 根据单位进行不同的转换。
	switch unit {
	case "MB":
		return size * 1024 * 1024, nil // 转换为字节
	case "GB":
		return size * 1024 * 1024 * 1024, nil // 转换为字节
	default:
		return 0, fmt.Errorf("不支持的日志大小单位: %s", unit)
	}
}

// ensureLogDirExists 确保日志目录存在，如果不存在则创建
func ensureLogDirExists(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建日志目录失败: %w", err)
		}
		fmt.Printf("日志目录不存在，已自动创建：%s\n", dir)
	} else if err != nil {
		return fmt.Errorf("检查日志目录失败: %w", err)
	}
	return nil
}

// 获取当前 Goroutine 的 ID
func getGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := bytes.Fields(buf[:n])[1]
	id, _ := strconv.ParseInt(string(idField), 10, 64)
	return id
}

// 核心日志记录函数
func (l *Logger) logWithLevel(level string, v ...interface{}) {
	// 将字符串日志级别转换为 LogLevel
	var logLevelValue LogLevel
	switch level {
	case "DEBUG":
		logLevelValue = Debug
	case "INFO":
		logLevelValue = Info
	case "WARN":
		logLevelValue = Warn
	case "ERROR":
		logLevelValue = Error
	case "SUCCESS":
		logLevelValue = Success
	default:
		logLevelValue = None
	}

	// 检查日志级别是否满足输出条件, 如果不满足则拦截并返回
	if logLevelValue < l.logLevel {
		return
	}

	// 获取调用者的信息 文件名、函数名、行号
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	// 获取文件名（只保留文件名，不包含路径）
	fileName := filepath.Base(file)

	// 获取函数名
	pc, _, _, _ := runtime.Caller(2)
	functionName := runtime.FuncForPC(pc).Name()

	// 根据日志格式选项设置日志格式
	var logMsg string
	switch l.logFormat {
	// Json格式
	case Json:
		logMsg = fmt.Sprintf(
			`{"time":"%s","level":"%s","file":"%s","function":"%s","line":"%d", "thread":"%d","message":"%s"}`,
			time.Now().Format("2006-01-02 15:04:05"), level, fileName, functionName, line, getGoroutineID(), fmt.Sprint(v...),
		)
	// 详细格式
	case Detailed:
		// 按照指定格式输出日志，使用%-7s让日志级别左对齐且宽度为7个字符
		logMsg = fmt.Sprintf(
			"%s | %-7s | %s:%s:%d - %s",
			time.Now().Format("2006-01-02 15:04:05"), level, fileName, functionName, line, fmt.Sprint(v...),
		)
	// 括号格式
	case Bracket:
		logMsg = fmt.Sprintf("[%s] %s", level, fmt.Sprint(v...))
	// 协程格式
	case Threaded:
		logMsg = fmt.Sprintf(`[%s] | %-7s | [thread="%d"] %s`, time.Now().Format("2006-01-02 15:04:05"), level, getGoroutineID(), fmt.Sprint(v...))
	// 无法识别的日志格式选项
	default:
		l.logger.Printf("%s", v...)
	}

	// 写入日志通道
	l.logChan <- logMsg

	// 增加计数器
	atomic.AddInt32(&counter, 1)
}

// 启动异步日志记录---日志通道消费者/缓冲区生产者
func (l *Logger) startAsyncLogging() {
	// 计算缓冲区大小
	BufferSize := l.bufferKbSize * 1024
	// 等待组计数器加1
	l.wg.Add(1)
	// 启动一个单独的 goroutine 来处理日志写入
	go func() {
		// 确保在函数退出时等待组计数器减1
		defer l.wg.Done()
		// 使用 for range 循环接收日志消息
		for {
			select {
			case <-l.stopChan:
				// 收到停止信号，退出循环
				return
			// 接收日志消息
			case logMsg := <-l.logChan:
				// 只有当需要输出到文件时才追加日志消息到文件缓冲区
				// 如果 consoleOnly 为 true，则不输出到文件
				if !l.consoleOnly && l.fileBuffer != nil {
					// 追加日志消息到文件缓冲区
					_, err := l.fileBuffer.WriteString(logMsg + "\n")
					if err != nil {
						l.logger.Printf("消费日志消息到文件缓冲区失败: %v", err)
						return
					}
				}

				// 如果需要输出到控制台，将日志消息追加到控制台缓冲区
				// consoleOnly 为 true 或 printToConsole 为 true 时，输出到控制台
				if l.printToConsole && l.consoleBuffer != nil {
					// 为日志消息添加颜色
					coloredLogMsg := l.addColor(logMsg)
					_, err := l.consoleBuffer.WriteString(coloredLogMsg + "\n")
					if err != nil {
						l.logger.Printf("消费日志消息到控制台缓冲区失败: %v", err)
						return
					}
				}

				// 检查文件缓冲区大小，如果达到阈值则立即进行写入
				if l.fileBuffer != nil && l.fileBuffer.Len() >= BufferSize {
					if !l.consoleOnly {
						l.flushBuffers()
					}
				}

				// 每次处理完一条日志到缓冲区后，减少计数器
				atomic.AddInt32(&counter, -1)
			}
		}
	}()

	// 启动定时消费缓冲区内容到文件和控制台 使用 sync.Once 确保只执行一次
	onefunc := sync.OnceFunc(func() {
		l.startBufferedLogging()
	})
	onefunc() // 执行一次
}

// 定时消费缓冲区内容到文件和控制台---缓冲区消费者
func (l *Logger) startBufferedLogging() {
	// 初始化定时器, 每1秒检查一次缓冲区
	l.ticker = time.NewTicker(1 * time.Second)
	// 等待组计数器加1
	l.wg.Add(1)
	go func() {
		// 确保在函数退出时等待组计数器减1
		defer l.wg.Done()
		// 使用 select 语句同时监听定时器通道和停止信号通道
		for {
			select {
			case <-l.ticker.C:
				// 定时检查并消费缓冲区内容到文件和控制台
				l.flushBuffers()
			case <-l.stopChan:
				// 收到停止信号，停止定时器
				l.ticker.Stop()
				// 再次检查缓冲区，确保所有日志都已处理
				l.flushBuffers()
				return
			}
		}
	}()
}

// 消费缓冲区内容到文件和控制台
func (l *Logger) flushBuffers() {
	// 如果文件缓冲区中有内容, 且不是仅输出到控制台
	if !l.consoleOnly && l.fileBuffer != nil {
		if l.fileBuffer.Len() > 0 {
			// 加锁，确保并发安全
			l.fileMu.Lock()
			// 写入文件缓冲区的内容到文件
			_, err := l.fileWriter.Write(l.fileBuffer.Bytes())
			if err != nil {
				l.logger.Printf("写入文件缓冲区内容到文件失败: %v", err)
			}
			// 清空文件缓冲区
			l.fileBuffer.Reset()
			// 解锁
			l.fileMu.Unlock()
		}
	}

	// 如果需要输出到控制台且控制台缓冲区中有内容
	if l.printToConsole && l.consoleBuffer != nil {
		if l.consoleBuffer.Len() > 0 {
			// 加锁，确保并发安全
			l.consoleMu.Lock()
			// 写入控制台缓冲区的内容到控制台
			_, err := l.consoleWriter.Write(l.consoleBuffer.Bytes())
			if err != nil {
				l.logger.Printf("写入控制台缓冲区内容到控制台失败: %v", err)
			}
			// 清空控制台缓冲区
			l.consoleBuffer.Reset()
			// 解锁
			l.consoleMu.Unlock()
		}
	}
}

// Info 级别的日志
func (l *Logger) Info(v ...interface{}) {
	l.logWithLevel("INFO", v...)
}

// Warn 级别的日志
func (l *Logger) Warn(v ...interface{}) {
	l.logWithLevel("WARN", v...)
}

// Error 级别的日志
func (l *Logger) Error(v ...interface{}) {
	l.logWithLevel("ERROR", v...)
}

// Success 级别的日志
func (l *Logger) Success(v ...interface{}) {
	l.logWithLevel("SUCCESS", v...)
}

// Debug 级别的日志
func (l *Logger) Debug(v ...interface{}) {
	l.logWithLevel("DEBUG", v...)
}

// 支持格式化的 Info 级别的日志
func (l *Logger) Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logWithLevel("INFO", msg)
}

// 支持格式化的 Warn 级别的日志
func (l *Logger) Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logWithLevel("WARN", msg)
}

// 支持格式化的 Error 级别的日志
func (l *Logger) Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logWithLevel("ERROR", msg)
}

// 支持格式化的 Success 级别的日志
func (l *Logger) Successf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logWithLevel("SUCCESS", msg)
}

// 支持格式化的 Debug 级别的日志
func (l *Logger) Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logWithLevel("DEBUG", msg)
}

// 关闭日志文件
func (l *Logger) Close() {
	// 循环检查计数器是否消费到缓存区
	for atomic.LoadInt32(&counter) > 0 {
		// 等待日志处理完成
		continue
	}
	l.closeOnce.Do(func() {
		// 发送停止信号
		close(l.stopChan)
		// 等待所有日志处理完成
		l.wg.Wait()
		// 确保所有缓冲区内容都被写入
		l.flushBuffers()
		// 关闭日志通道
		close(l.logChan)
		// 如果日志文件存在，则关闭日志文件
		if l.fileWriter != nil {
			l.logFile.Close()
		}
	})
}

// String 方法，将 LogFormatType 转换为字符串
func (lft LogFormatType) String() string {
	switch lft {
	case Json:
		return "json"
	case Bracket:
		return "bracket"
	case Detailed:
		return "detailed"
	case Threaded:
		return "threaded"
	default:
		return "unknown"
	}
}

// 为控制台输出添加颜色
func (w *Logger) addColor(s string) string {
	// 使用正则表达式精确匹配日志级别，确保匹配的是独立的单词
	re := regexp.MustCompile(`\b(INFO|WARN|ERROR|SUCCESS|DEBUG)\b`)
	match := re.FindString(s)

	// 根据匹配到的日志级别添加颜色
	switch match {
	case "INFO":
		return color.Sblue(s) // Blue
	case "WARN":
		return color.Syellow(s) // Yellow
	case "ERROR":
		return color.Sred(s) // Red
	case "SUCCESS":
		return color.Sgreen(s) // Green
	case "DEBUG":
		return color.Spurple(s) // Purple
	default:
		return s // 如果没有匹配到日志级别，返回原始字符串
	}
}

// 日志轮转主协程 启动日志切割工作协程
func (l *Logger) startLogRotate() {
	// 定义一个定时器 每隔10分钟检查一次
	ticker := time.NewTicker(time.Duration(l.rotationInterval) * time.Minute)

	// 启动一个单独的 goroutine 来处理日志轮转
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for {
			select {
			case <-ticker.C: // 定时器触发，触发日志切割工作协程
				l.logRotateWorker()
			case <-l.stopChan:
				// 收到停止信号，停止定时器
				ticker.Stop()
				return // 退出 goroutine
			}
		}
	}()
}

// 日志切割工作方法---检查
func (l *Logger) logRotateWorker() {
	// 检查日志文件是否存在
	if _, err := os.Stat(l.logFilePath); os.IsNotExist(err) {
		l.Errorf("日志文件不存在: %s", l.logFilePath)
		return
	}

	// 检查日志大小是否超过限制
	fileInfo, err := os.Stat(l.logFilePath)
	if err != nil {
		l.Errorf("获取日志文件信息失败, 日志轮转返回: %v", err)
		return
	}
	if fileInfo.Size() >= l.logSizeBytes { // 检查日志大小是否超过限制
		err := l.logRotate() // 调用日志切割函数
		if err != nil {
			l.Errorf("日志切割失败, 日志轮转返回: %v", err)
			return
		}
	}

	// 启动日志清理协程
	l.Info("开始清理历史日志...")
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		err := l.logFileClean() // 调用日志清理函数
		if err != nil {
			l.Errorf("日志清理失败, 日志轮转返回: %v", err)
			return
		}
	}()

}

// 日志切割工作---切割
func (l *Logger) logRotate() error {
	// 先获取日志文件锁，堵塞写入
	l.fileMu.Lock()
	defer l.fileMu.Unlock()

	// 检查旧日志文件是否存在
	if _, err := os.Stat(l.logFilePath); os.IsNotExist(err) {
		// 如果文件不存在，记录错误并尝试重新创建日志文件
		l.Errorf("旧日志文件 %s 不存在，尝试重新创建...", l.logFilePath)
		newLogFile, err := os.OpenFile(l.logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return fmt.Errorf("重新创建日志文件失败: %w", err)
		}
		l.logFile = newLogFile   // 更新日志文件句柄
		l.fileWriter = l.logFile // 更新日志文件写入器句柄
		return nil
	} else if err != nil {
		return fmt.Errorf("检查旧日志文件失败: %w", err)
	}

	// 获取当前时间
	currentTime := time.Now().Format("2006-01-02_15-04-05")

	// 创建新的日志文件名
	newLogFileName := fmt.Sprintf("%s_%s.log", l.logFileName, currentTime)
	newLogFilePath := filepath.Join(l.logDirName, newLogFileName)

	// 关闭旧的日志文件
	l.logFile.Close()

	// 重命名旧的日志文件 (问题点稍后研究)
	if err := os.Rename(l.logFilePath, newLogFilePath); err != nil {
		return fmt.Errorf("重命名旧日志文件失败: %w", err)
	}

	// 创建新的日志文件
	newLogFile, err := os.OpenFile(l.logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("创建新日志文件失败: %w", err)
	}

	// 更新日志文件句柄
	l.logFile = newLogFile

	// 更新日志文件写入器句柄
	l.fileWriter = l.logFile

	// 可选：记录切割操作到日志
	l.Successf("日志文件已切割: %s -> %s", l.logFilePath, newLogFilePath)

	return nil
}

// logFileClean 清理超出保留天数或保留数量的日志文件
func (l *Logger) logFileClean() error {
	files, err := os.ReadDir(l.logDirName)
	if err != nil {
		return fmt.Errorf("无法读取日志目录: %w", err)
	}

	// 存储日志文件信息
	var logFiles []struct {
		Path    string
		ModTime time.Time
	}

	// 过滤出日志文件并按时间排序
	for _, file := range files {
		if file.IsDir() { // 跳过目录
			continue
		}
		if filepath.Ext(file.Name()) == ".log" {
			filePath := filepath.Join(l.logDirName, file.Name()) // 拼接后的格式为 "logDirName/fileName"
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				l.Errorf("无法获取文件信息: %s, 错误: %v", filePath, err)
				continue
			}

			// 如果是当前正在写入的日志文件，跳过
			if filePath == l.logFilePath { // 跳过当前正在写入的日志文件
				continue
			}

			// 添加到日志文件列表
			logFiles = append(logFiles, struct {
				Path    string
				ModTime time.Time
			}{Path: filePath, ModTime: fileInfo.ModTime()}) // 记录文件路径和修改时间
		}
	}

	// 按文件名中的时间戳排序，如果文件名中没有时间戳，则按修改时间排序
	sort.Slice(logFiles, func(i, j int) bool {
		timestampI := l.extractTimestamp(filepath.Base(logFiles[i].Path)) // 提取文件名中的时间戳
		timestampJ := l.extractTimestamp(filepath.Base(logFiles[j].Path)) // 提取文件名中的时间戳

		if timestampI.IsZero() { // 如果文件名中没有时间戳，则按修改时间排序
			return logFiles[i].ModTime.Before(logFiles[j].ModTime)
		}
		if timestampJ.IsZero() { // 如果文件名中没有时间戳，则按修改时间排序
			return logFiles[i].ModTime.Before(logFiles[j].ModTime)
		}
		return timestampI.Before(timestampJ) // 按时间戳排序
	})

	// 计算保留天数前的时间
	cutoffTime := time.Now().AddDate(0, 0, -l.logRetentionDays)

	// 清理超出保留天数的日志文件
	for _, file := range logFiles {
		if file.ModTime.Before(cutoffTime) { // 如果文件修改时间早于保留天数前的时间
			if err := os.Remove(file.Path); err != nil {
				l.Errorf("清理日志文件失败: %s, 错误: %v", file.Path, err)
			} else {
				l.Successf("已清理超出保留天数的日志文件: %s", file.Path)
			}
		}
	}

	// 清理超出保留数量的日志文件
	if len(logFiles) > l.logRetentionCount {
		for i := 0; i < len(logFiles)-l.logRetentionCount; i++ {
			if err := os.Remove(logFiles[i].Path); err != nil {
				l.Errorf("清理日志文件失败: %s, 错误: %v", logFiles[i].Path, err)
			} else {
				l.Successf("已清理超出保留数量的日志文件: %s", logFiles[i].Path)
			}
		}
	}

	return nil
}

// 提取文件名中的时间戳
func (l *Logger) extractTimestamp(fileName string) time.Time {
	match := timestampRegex.FindStringSubmatch(fileName) // 匹配文件名中的时间戳
	if len(match) > 1 {
		timestamp, err := time.Parse("2006-01-02_15-04-05", match[1]) // 解析时间戳
		if err != nil {
			// 记录错误日志
			l.Errorf("解析时间戳失败: %s, 错误: %v", fileName, err)
			return time.Time{}
		}
		return timestamp
	}
	return time.Time{}
}

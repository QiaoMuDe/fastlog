package fastlog

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
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
	Json     LogFormatType = iota // Json格式
	Bracket                       // 方括号格式
	Detailed                      // 详细格式
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
	logDirName     string         // 日志目录
	logFileName    string         // 日志文件名
	logFile        *os.File       // 日志文件句柄
	logger         *log.Logger    // 底层日志记录器
	printToConsole bool           // 是否将日志输出到控制台
	consoleOnly    bool           // 是否仅输出到控制台
	logLevel       LogLevel       // 日志级别
	fileMu         sync.Mutex     // 文件写入的互斥锁
	consoleMu      sync.Mutex     // 控制台写入的互斥锁
	logChan        chan string    // 日志通道
	stopChan       chan string    // 停止通道
	wg             sync.WaitGroup // 等待组，用于等待日志通道中的日志被处理
	fileBuffer     *bytes.Buffer  // 文件缓冲区
	consoleBuffer  *bytes.Buffer  // 控制台缓冲区
	ticker         *time.Ticker   // 定时器
	fileWriter     io.Writer      // 文件写入器
	consoleWriter  io.Writer      // 控制台写入器
	chanIntSize    int            // 通道大小
	bufferKbSize   int            // 缓冲区大小
	closeOnce      sync.Once      // 确保通道只被关闭一次
	logFormat      LogFormatType  // 日志格式选项 [Json(json格式)|Bracket(方括号格式)|Detailed(详细格式)|Threaded(协程格式)]
}

// 日志配置
type LoggerConfig struct {
	LogDirName     string        // 日志目录名称
	LogFileName    string        // 日志文件名称
	LogPath        string        // 日志文件路径
	PrintToConsole bool          // 是否将日志输出到控制台
	ConsoleOnly    bool          // 是否仅输出到控制台
	LogLevel       LogLevel      // 日志级别
	ChanIntSize    int           // 通道大小
	BufferKbSize   int           // 缓冲区大小
	LogFormat      LogFormatType // 日志格式选项 [Json(json格式)|Bracket(括号格式)|Detailed(详细格式)|Threaded(协程格式)]
}

// 定义一个接口
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

// 创建一个新的日志配置器
func NewConfig(logDirName string, logFileName string) LoggerConfig {
	// 生成日志文件路径
	logPath := filepath.Join(logDirName, logFileName)
	return LoggerConfig{
		LogDirName:     logDirName,  // 日志目录 必须提供
		LogFileName:    logFileName, // 日志文件名 必须提供
		LogPath:        logPath,     // 内部自己拼接的路径,无需提供
		PrintToConsole: true,        // 是否将日志输出到控制台, 默认值为true
		ConsoleOnly:    false,       // 是否仅输出到控制台, 默认值为false
		LogLevel:       Info,        // 日志过滤级别, 默认值为Info
		ChanIntSize:    1000,        // 日志通道大小, 默认1000
		BufferKbSize:   1024,        // 缓冲区大小, 默认1MB
		LogFormat:      Detailed,    // 日志格式选项 [Json(json格式)|Bracket(括号格式)|Detailed(详细格式)|Threaded(协程格式)] 默认详细格式
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
		// 检查日志目录是否存在，如果不存在则创建
		if err := ensureLogDirExists(cfg.LogDirName); err != nil {
			return nil, err
		}

		// 打开日志文件
		logFile, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
		logDirName:     cfg.LogDirName,      // 保存日志目录名称
		logFileName:    cfg.LogFileName,     // 保存日志文件名称
		logFile:        logFile,             // 保存日志文件句柄
		logger:         logger,              // 保存底层日志记录器
		printToConsole: outputToConsole,     // 保存是否将日志输出到控制台
		consoleOnly:    outputOnlyToConsole, // 保存是否仅输出到控制台
		logLevel:       cfg.LogLevel,        // 保存日志级别
		fileWriter:     fileWriter,          // 保存文件写入器
		consoleWriter:  os.Stdout,           // 保存控制台写入器
		chanIntSize:    cfg.ChanIntSize,     // 保存通道大小
		bufferKbSize:   cfg.BufferKbSize,    // 保存缓冲区大小
		fileBuffer:     fileBuffer,          // 保存预分配内存后的文件缓冲区
		consoleBuffer:  consoleBuffer,       // 保存预分配内存后的控制台缓冲区
		logFormat:      cfg.LogFormat,       // 保存日志格式选项 [Json(json格式)|Bracket(括号格式)|Detailed(详细格式)|Threaded(协程格式)]
	}

	// 初始化日志通道, 容量为1000
	lg.logChan = make(chan string, lg.chanIntSize)

	// 初始化停止通道
	lg.stopChan = make(chan string, 10)

	// 启动异步日志记录
	lg.startAsyncLogging()

	// 返回 Logger 实例
	return lg, nil
}

// ensureLogDirExists 确保日志目录存在，如果不存在则创建
func ensureLogDirExists(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建日志目录失败: %w", err)
		}
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

// 启动异步日志记录
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

	// 启动定时器，定期将缓冲区的内容写入文件和控制台
	l.startBufferedLogging()
}

// 定时消费缓冲区内容到文件和控制台
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

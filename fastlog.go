package gitee.com/MM-Q/fastlog.git

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"path/filepath"
)

// 日志级别枚举
type LogLevel int

// 定义一个原子计数器
var counter int32 = 0

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
}

// 日志配置
type LoggerConfig struct {
	LogDirName     string   // 日志目录名称
	LogFileName    string   // 日志文件名称
	LogPath        string   // 日志文件路径
	PrintToConsole bool     // 是否将日志输出到控制台
	LogLevel       LogLevel // 日志级别
	ChanIntSize    int      // 通道大小
	BufferKbSize   int      // 缓冲区大小
}

// 定义一个接口
type LoggerInterface interface {
	Info(v ...interface{})                         // 记录信息级别的日志
	Warn(v ...interface{})                         // 记录警告级别的日志
	Error(v ...interface{})                        // 记录错误级别的日志
	Success(v ...interface{})                      // 记录成功级别的日志
	Debug(v ...interface{})                        // 记录调试级别的日志
	Close()                                        // 关闭日志记录器
	DefaultConfig(logFilePath string) LoggerConfig // 创建一个日志配置器
	NewLogger(cfg LoggerConfig) (*Logger, error)   // 创建一个日志记录器
}

// 默认配置
func DefaultConfig(logDirName string, logFileName string) LoggerConfig {
	// 生成日志文件路径
	logPath := filepath.Join(logDirName, logFileName)
	return LoggerConfig{
		LogDirName:     "logs",    // 必须提供
		LogFileName:    "app.log", // 必须提供
		LogPath:        logPath,   // 无需修改
		PrintToConsole: true,      // 默认值
		LogLevel:       Info,      // 默认值
		ChanIntSize:    1000,      // 默认1000
		BufferKbSize:   1024,      // 默认1MB
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
		return fmt.Sprintf("\033[34m%s\033[0m", s) // Blue
	case "WARN":
		return fmt.Sprintf("\033[33m%s\033[0m", s) // Yellow
	case "ERROR":
		return fmt.Sprintf("\033[31m%s\033[0m", s) // Red
	case "SUCCESS":
		return fmt.Sprintf("\033[32m%s\033[0m", s) // Green
	case "DEBUG":
		return fmt.Sprintf("\033[35m%s\033[0m", s) // Purple
	default:
		return s // 如果没有匹配到日志级别，返回原始字符串
	}
}

// 创建一个新的日志记录器
func NewLogger(cfg LoggerConfig) (*Logger, error) {
	// 检查日志目录是否存在，如果不存在则创建
	if _, err := os.Stat(cfg.LogDirName); os.IsNotExist(err) {
		if err := os.Mkdir(cfg.LogDirName, 0755); err != nil {
			fmt.Println("创建日志目录失败: ", err)
			os.Exit(1)
		}
	}

	// 打开日志文件
	logF, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// 初始化 fileWriter写入器 和 consoleWriter写入器
	fileWriter := logF
	consoleWriter := os.Stdout

	// 创建一个自定义的日志记录器，不设置默认的前缀和标志
	logger := log.New(fileWriter, "", 0)

	// 初始化缓冲区大小 默认1MB
	bufferSize := cfg.BufferKbSize * 1024

	// 使用预分配内存的方式初始化文件缓冲区
	fileBuffer := bytes.NewBuffer(make([]byte, bufferSize))
	// 初始化缓冲区就先清空, 防止上次的日志残留

	fileBuffer.Reset()
	// 初始化控制台缓冲区
	var consoleBuffer *bytes.Buffer
	// 如果需要输出到控制台，使用预分配内存的方式初始化控制台缓冲区
	if cfg.PrintToConsole {
		consoleBuffer = bytes.NewBuffer(make([]byte, bufferSize))
		// 初始化缓冲区就先清空, 防止上次的日志残留
		consoleBuffer.Reset()
	}

	// 创建 Logger 实例
	lg := &Logger{
		logDirName:     cfg.LogDirName,     // 保存日志目录名称
		logFileName:    cfg.LogFileName,    // 保存日志文件名称
		logFile:        logF,               // 保存日志文件句柄
		logger:         logger,             // 保存底层日志记录器
		printToConsole: cfg.PrintToConsole, // 保存是否将日志输出到控制台
		logLevel:       cfg.LogLevel,       // 保存日志级别
		fileWriter:     fileWriter,         // 保存文件写入器
		consoleWriter:  consoleWriter,      // 保存控制台写入器
		chanIntSize:    cfg.ChanIntSize,    // 保存通道大小
		bufferKbSize:   cfg.BufferKbSize,   // 保存缓冲区大小
		fileBuffer:     fileBuffer,         // 保存预分配内存后的文件缓冲区
		consoleBuffer:  consoleBuffer,      // 保存预分配内存后的控制台缓冲区
	}

	// 初始化日志通道, 容量为1000
	lg.logChan = make(chan string, lg.chanIntSize)

	// 初始化停止通道
	lg.stopChan = make(chan string, 10)

	// 启动异步日志记录
	lg.StartAsyncLogging()

	// 返回 Logger 实例
	return lg, nil
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

	// 检查日志级别是否满足输出条件
	if logLevelValue < l.logLevel {
		return
	}

	// 获取调用者的信息
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

	// 按照指定格式输出日志，使用%-7s让日志级别左对齐且宽度为7个字符
	logMsg := fmt.Sprintf("%s | %-7s | %s:%s:%d - %s", time.Now().Format("2006-01-02 15:04:05"), level, fileName, functionName, line, fmt.Sprint(v...))

	// 写入日志通道
	l.logChan <- logMsg
	// 增加计数器
	atomic.AddInt32(&counter, 1)
}

// 启动异步日志记录
func (l *Logger) StartAsyncLogging() {
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
				// 将日志消息追加到文件缓冲区
				l.fileBuffer.WriteString(logMsg + "\n")
				// 如果需要输出到控制台，将日志消息追加到控制台缓冲区
				if l.printToConsole {
					coloredLogMsg := l.addColor(logMsg)
					l.consoleBuffer.WriteString(coloredLogMsg + "\n")
				}
				// 检查缓冲区大小，如果达到阈值则立即进行写入
				if l.fileBuffer.Len() >= BufferSize {
					l.flushBuffers()
				}
				// 减少计数器
				atomic.AddInt32(&counter, -1)
			}
		}
	}()

	// 启动定时器，定期将缓冲区的内容写入文件和控制台
	l.StartBufferedLogging()
}

// 定时消费缓冲区内容到文件和控制台
func (l *Logger) StartBufferedLogging() {
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
	// 如果文件缓冲区中有内容
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

	// 如果需要输出到控制台且控制台缓冲区中有内容
	if l.printToConsole && l.consoleBuffer.Len() > 0 {
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
		// 关闭日志文件
		l.logFile.Close()
	})
}

package fastlog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// NewFastLogConfig 创建一个新的FastLogConfig实例，用于配置日志记录器。
// 参数:
//   - logDirName: 日志目录名称，默认为"applogs"。
//   - logFileName: 日志文件名称，默认为"app.log"。
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
func NewFastLogConfig(logDirPath string, logFileName string) *FastLogConfig {
	// 如果日志目录名称为空，则使用默认值"logs"。
	if logDirPath == "" {
		logDirPath = "logs"
	}
	// 如果日志文件名称为空，则使用默认值"app.log"。
	if logFileName == "" {
		logFileName = "app.log"
	}
	// 合并日志目录和日志文件名称，生成日志文件路径。
	logFilePath := filepath.Join(logDirPath, logFileName)
	return &FastLogConfig{
		logDirPath:     logDirPath,  // 日志目录名称
		LogFileName:    logFileName, // 日志文件名称
		logFilePath:    logFilePath, // 日志文件路径
		PrintToConsole: true,        // 是否将日志输出到控制台
		ConsoleOnly:    false,       // 是否仅输出到控制台
		LogLevel:       INFO,        // 日志级别 默认INFO
		ChanIntSize:    1000,        // 通道大小 默认1000
		BufferKbSize:   1024,        // 缓冲区大小 默认1024 单位KB
		LogFormat:      Detailed,    // 日志格式选项
	}
}

// NewFastLog 创建一个新的FastLog实例，用于记录日志。
// 参数:
//   - config: 一个指向FastLogConfig实例的指针，用于配置日志记录器。
//
// 返回值:
//   - *FastLog: 一个指向FastLog实例的指针。
//   - error: 如果创建日志记录器失败，则返回一个错误。
func NewFastLog(config *FastLogConfig) (*FastLog, error) {
	// 声明一些配置变量
	var (
		consoleBuffer *bytes.Buffer // 初始化控制台缓冲区
		fileBuffer    *bytes.Buffer // 初始化日志文件缓冲区
		logFile       *os.File      // 日志文件句柄
		fileWriter    io.Writer     // 文件写入器
		consoleWriter io.Writer     // 控制台写入器
		err           error         // 错误变量
	)

	// 初始化缓冲区大小 默认1MB
	bufferSize := config.BufferKbSize * 1024

	// 如果允许将日志输出到控制台，或者仅输出到控制台，则初始化控制台缓冲区和写入器。
	if config.ConsoleOnly || config.PrintToConsole {
		consoleBuffer = bytes.NewBuffer(make([]byte, bufferSize)) // 创建一个新的缓冲区, 初始大小为 bufferSize
		consoleBuffer.Reset()                                     // 重置缓冲区
		consoleWriter = os.Stdout                                 // 控制台写入器
	} else {
		consoleWriter = io.Discard // 不输出到控制台，直接丢弃
	}

	// 如果不是仅输出到控制台，则初始化日志文件缓冲区和写入器。
	if !config.ConsoleOnly {
		// 检查日志目录是否存在，如果不存在则创建。
		if _, err := checkPath(config.logDirPath); err != nil {
			if err := os.MkdirAll(config.logDirPath, 0644); err != nil {
				return nil, fmt.Errorf("创建日志目录失败: %s", err)
			}
		}

		// 打开日志文件，如果文件不存在则创建。
		logFile, err = os.OpenFile(config.logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("打开日志文件失败: %s", err)
		}

		// 初始化文件缓冲区
		fileBuffer = bytes.NewBuffer(make([]byte, bufferSize))
		fileBuffer.Reset()

		// 初始化文件写入器
		fileWriter = logFile
	} else {
		fileWriter = io.Discard // 仅输出到控制台，不输出到文件
	}

	// 创建一个上下文，用于控制日志记录器的生命周期。
	ctx, cancel := context.WithCancel(context.Background())

	// 创建一个新的FastLog实例，将配置和缓冲区赋值给实例。
	fastLog := &FastLog{
		logFile:        logFile,               // 日志文件句柄
		fileWriter:     fileWriter,            // 文件写入器,
		consoleWriter:  consoleWriter,         // 控制台写入器,
		logFilePath:    config.logFilePath,    // 日志文件路径
		fileBuffer:     fileBuffer,            // 日志文件缓冲区
		consoleBuffer:  consoleBuffer,         // 控制台缓冲区
		logDirPath:     config.logDirPath,     // 日志目录路径
		logFileName:    config.LogFileName,    // 日志文件名
		printToConsole: config.PrintToConsole, // 是否将日志输出到控制台
		consoleOnly:    config.ConsoleOnly,    // 是否仅输出到控制台
		logLevel:       config.LogLevel,       // 日志级别
		chanIntSize:    config.ChanIntSize,    // 通道大小 默认1000
		bufferKbSize:   config.BufferKbSize,   // 缓冲区大小 默认1024 单位KB
		logFormat:      config.LogFormat,      // 日志格式选项
		logCtx:         ctx,                   // 日志上下文,
		stopChan:       cancel,                // 日志上下文取消函数,
		isWriter:       false,                 // 是否在写入日志
	}

	// 初始化日志通道
	fastLog.logChan = make(chan *logMessage, fastLog.chanIntSize)

	// 使用sync.Once确保只执行一次
	var ProcessLogsOnce sync.Once
	ProcessLogsOnce.Do(func() {
		// 启动日志处理器
		ProcessLogs(fastLog)
	}) // 执行一次

	// 使用sync.Once确保只执行一次
	var startBufferedLoggingOnce sync.Once
	startBufferedLoggingOnce.Do(func() {
		// 启动日志缓冲器
		startBufferedLogging(fastLog)
	}) // 执行一次

	// 返回FastLog实例和nil错误
	return fastLog, nil
}

// ProcessLogs 日志处理器，用于处理日志消息。
// 读取通道中的日志消息并将其处理为指定的日志格式，然后写入到缓冲区中。
func ProcessLogs(f *FastLog) {
	f.logWait.Add(1) // 增加等待组中的计数器。
	// 启动一个goroutine，用于处理日志消息。
	go func() {
		defer func() {
			// 减少等待组中的计数器。
			f.logWait.Done()

			// 捕获panic
			if r := recover(); r != nil {
				panic(r)
			}
		}()

		// 循环监听
		for {
			select {
			case logMsg := <-f.logChan:
				// 发送缓冲区内容到文件和控制台
				if err := sendBufferToFile(f, logMsg); err != nil {
					panic(err)
				}

			case <-f.logCtx.Done(): // 如果收到停止信号，则退出循环
				// 最后循环一次，确保通道中的所有日志消息都被处理
				for logMsg := range f.logChan {
					if err := sendBufferToFile(f, logMsg); err != nil {
						panic(err)
					}
				}
				return // 退出goroutine

			default:
				// 如果通道中没有日志消息，则等待一段时间后继续循环。
				time.Sleep(time.Millisecond * 100) // 等待100毫秒
			}
		}
	}()
}

// 日志缓冲器，用于定时调用日志写入器，将缓冲区中的日志写入到文件和控制台。
func startBufferedLogging(f *FastLog) {
	// 初始化定时器, 每1秒检查一次缓冲区
	ticker := time.NewTicker(1 * time.Second)
	// 等待组计数器加1
	f.logWait.Add(1)
	go func() {
		// 确保在函数退出时等待组计数器减1
		defer func() {
			f.logWait.Done()

			// 捕获panic
			if r := recover(); r != nil {
				panic(r)
			}
		}()
		// 使用 select 语句同时监听定时器通道和停止信号通道
		for {
			select {
			case <-ticker.C:
				// 定时检查并消费缓冲区内容到文件和控制台
				if err := logBuffersWriter(f); err != nil {
					panic(err)
				}
			case <-f.logCtx.Done():
				// 收到停止信号，停止定时器
				ticker.Stop()
				// 再次检查缓冲区，确保所有日志都已处理
				if err := logBuffersWriter(f); err != nil {
					panic(err)
				}
				return
			}
		}
	}()
}

// 日志写入器，用于将缓冲区中的日志写入到文件和控制台。
func logBuffersWriter(f *FastLog) error {
	// 加锁，确保并发安全
	f.fileMu.Lock()
	defer f.fileMu.Unlock()

	// 如果正在写入日志, 则直接返回
	if f.isWriter {
		return nil
	}
	// 设置正在写入日志
	f.isWriter = true

	// 返回时重置写入标志
	defer func() {
		f.isWriter = false
	}()

	// 如果文件缓冲区中有内容, 且不是仅输出到控制台
	if !f.consoleOnly && f.fileBuffer != nil {
		if f.fileBuffer.Len() > 0 {

			// 写入文件缓冲区的内容到文件
			_, err := f.fileWriter.Write(f.fileBuffer.Bytes())
			if err != nil {
				return fmt.Errorf("写入文件缓冲区内容到文件失败: %v", err)
			}
			// 清空文件缓冲区
			f.fileBuffer.Reset()
		}
	}

	// 加锁，确保并发安全
	f.consoleMu.Lock()
	defer f.consoleMu.Unlock()

	// 如果需要输出到控制台且控制台缓冲区中有内容
	if f.printToConsole && f.consoleBuffer != nil {
		if f.consoleBuffer.Len() > 0 {
			// 写入控制台缓冲区的内容到控制台
			_, err := f.consoleWriter.Write(f.consoleBuffer.Bytes())
			if err != nil {
				return fmt.Errorf("写入控制台缓冲区内容到控制台失败: %v", err)
			}
			// 清空控制台缓冲区
			f.consoleBuffer.Reset()
		}
	}

	return nil
}

// Close 关闭FastLog实例，并等待所有日志处理完成。
func (f *FastLog) Close() {
	// 打印关闭日志记录器的信息
	f.Info("开始关闭日志记录器...")

	// 确保只关闭一次
	var closeOnce sync.Once
	closeOnce.Do(func() {
		// 关闭文件句柄, 如果输出到日志文件
		defer func() {
			if !f.consoleOnly && f.logFile != nil {
				f.logFile.Close()
			}
		}()

		// 关闭日志通道
		close(f.logChan)
		// 调用日志上下文取消函数，发送停止信号
		f.stopChan()

		// 等待所有日志处理完成
		f.logWait.Wait()

		// 清空缓冲区
		if f.fileBuffer != nil {
			f.fileBuffer.Reset()
		}
		if f.consoleBuffer != nil {
			f.consoleBuffer.Reset()
		}
	}) // 执行一次
}

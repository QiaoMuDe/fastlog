package fastlog

import (
	"bytes"
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
		LogFormat:      Detailed,    // 日志格式选项
		MaxBufferSize:  1,           // 最大缓冲区大小 默认1MB，单位为MB
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
		logFile       *os.File  // 日志文件句柄
		fileWriter    io.Writer // 文件写入器
		consoleWriter io.Writer // 控制台写入器
		err           error     // 错误变量
	)

	// 如果允许将日志输出到控制台，或者仅输出到控制台，则初始化控制台缓冲区和写入器。
	if config.ConsoleOnly || config.PrintToConsole {
		consoleWriter = os.Stdout // 控制台写入器
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

		// 初始化文件写入器
		fileWriter = logFile
	} else {
		fileWriter = io.Discard // 仅输出到控制台，不输出到文件
	}

	// 创建一个新的FastLog实例，将配置和缓冲区赋值给实例。
	f := &FastLog{
		logFile:        logFile,                            // 日志文件句柄
		fileWriter:     fileWriter,                         // 文件写入器,
		consoleWriter:  consoleWriter,                      // 控制台写入器,
		logFilePath:    config.logFilePath,                 // 日志文件路径
		logDirPath:     config.logDirPath,                  // 日志目录路径
		logFileName:    config.LogFileName,                 // 日志文件名
		printToConsole: config.PrintToConsole,              // 是否将日志输出到控制台
		consoleOnly:    config.ConsoleOnly,                 // 是否仅输出到控制台
		logLevel:       config.LogLevel,                    // 日志级别
		chanIntSize:    config.ChanIntSize,                 // 通道大小 默认1000
		logFormat:      config.LogFormat,                   // 日志格式选项
		maxBufferSize:  config.MaxBufferSize * 1024 * 1024, // 最大缓冲区大小 默认1MB
		flushInterval:  1 * time.Second,                    // 刷新间隔，单位为秒
	}

	// 初始化日志通道
	f.logChan = make(chan *logMessage, f.chanIntSize)

	// 初始化文件缓冲区
	f.fileBuffer = bytes.NewBuffer(make([]byte, f.maxBufferSize))
	f.fileBuffer.Reset() // 重置缓冲区，清空内容

	// 初始化控制台缓冲区
	f.consoleBuffer = bytes.NewBuffer(make([]byte, f.maxBufferSize))
	f.consoleBuffer.Reset() // 重置缓冲区，清空内容

	// 创建定时刷新任务
	f.flushTicker = time.NewTicker(f.flushInterval)

	f.logWait.Add(1) // 增加等待组中的计数器。
	// 使用 sync.Once 确保日志处理器只启动一次
	f.startOnce.Do(func() {
		go processLogs(f) // 启动日志处理器
		go flushBuffer(f) // 启动定时刷新缓冲区
	})

	// 返回FastLog实例和nil错误
	return f, nil
}

// processLogs 日志处理器，用于处理日志消息。
// 读取通道中的日志消息并将其处理为指定的日志格式，然后写入到缓冲区中。
func processLogs(f *FastLog) {
	defer func() {
		// 减少等待组中的计数器。
		f.logWait.Done()

		// 捕获panic
		if r := recover(); r != nil {
			panic(fmt.Sprintf("日志处理器发生panic: %v", r))
		}
	}()

	// 计算最大缓冲区大小的80%
	maxBufferSize := int(float64(f.maxBufferSize) * 0.8)

	// 初始化控制台字符串构建器
	for rawMsg := range f.logChan {
		// 检查缓冲区大小是否达到80%
		if f.fileBuffer.Len() >= maxBufferSize || f.consoleBuffer.Len() >= maxBufferSize {
			f.flushBufferNow()
		}

		// 格式化日志消息
		formattedLog := formatLog(f, rawMsg)

		// 如果不是仅输出到控制台，则将日志消息写入到日志文件缓冲区中。
		if !f.consoleOnly {
			// 复制格式化后的日志消息到文件日志变量中
			fileLog := formattedLog

			// 给格式化后的日志消息添加换行符
			f.fileBuilder.WriteString(fileLog) // 将格式化后的日志消息写入到构建器中
			f.fileBuilder.WriteString("\n")    // 添加换行符
			fileLog = f.fileBuilder.String()   // 将构建器中的内容赋值给 formattedLog

			// 将格式化后的日志消息写入到日志文件缓冲区中
			if _, err := f.fileBuffer.WriteString(fileLog); err != nil {
				panic(fmt.Errorf("写入文件缓冲区失败: %v", err))
			}

			f.fileBuilder.Reset() // 重置构建器
		}

		// 如果允许将日志输出到控制台，或者仅输出到控制台
		if f.consoleOnly || f.printToConsole {
			// 复制格式化后的日志消息到控制台日志变量中
			consoleLog := formattedLog

			// 调用addColor方法给日志消息添加颜色
			consoleLog = addColor(formattedLog)

			// 给格式化后的日志消息添加换行符
			f.consoleBuilder.WriteString(consoleLog) // 将格式化后的日志消息写入到构建器中
			f.consoleBuilder.WriteString("\n")       // 添加换行符
			consoleLog = f.consoleBuilder.String()   // 将构建器中的内容赋值给 formattedLog

			// 将格式化后的日志消息写入到控制台缓冲区中
			if _, err := f.consoleBuffer.WriteString(consoleLog); err != nil {
				panic(fmt.Errorf("写入控制台缓冲区失败: %v", err))
			}

			f.consoleBuilder.Reset() // 重置构建器
		}
	}
}

// flushBuffer 定时刷新缓冲区
func flushBuffer(f *FastLog) {
	// 使用定时器，每隔flushInterval时间触发一次
	for range f.flushTicker.C {
		f.flushBufferNow()
	}
}

// flushBufferNow 立即刷新缓冲区
func (f *FastLog) flushBufferNow() {
	// 写入文件
	if f.fileBuffer.Len() > 0 {
		// 获取文件写入写入锁，确保线程安全
		f.fileMu.Lock()
		defer f.fileMu.Unlock()

		// 获取文件缓冲区的内容
		bufferContent := f.fileBuffer.String()

		// 将文件缓冲区的内容写入到文件中
		if _, err := f.fileWriter.Write([]byte(bufferContent)); err != nil {
			fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
		}

		// 重置缓冲区
		f.fileBuffer.Reset()
	}

	// 写入控制台
	if f.consoleBuffer.Len() > 0 {
		// 获取控制台写入写入锁，确保线程安全
		f.consoleMu.Lock()
		defer f.consoleMu.Unlock()

		// 获取控制台缓冲区的内容
		bufferContent := f.consoleBuffer.String()

		// 将控制台缓冲区的内容写入到控制台
		if _, err := f.consoleWriter.Write([]byte(bufferContent)); err != nil {
			fmt.Fprintf(os.Stderr, "写入控制台失败: %v\n", err)
		}

		// 重置缓冲区
		f.consoleBuffer.Reset()
	}
}

// Close 关闭FastLog实例，并等待所有日志处理完成。
func (f *FastLog) Close() {
	// 打印关闭日志记录器的信息
	f.Info("关闭日志记录器...")

	// 确保只关闭一次
	var closeOnce sync.Once
	closeOnce.Do(func() {
		// 关闭日志通道
		close(f.logChan)

		// 等待所有日志处理完成
		f.logWait.Wait()

		// 停止定时刷新任务
		f.flushTicker.Stop()

		// 刷新剩余的日志
		f.flushBufferNow()

		// 如果不是仅输出到控制台，同时日志文件句柄不为nil，则关闭日志文件。
		if !f.consoleOnly && f.logFile != nil {
			f.logFile.Close()
		}

	}) // 执行一次
}

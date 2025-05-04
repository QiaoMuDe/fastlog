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

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/logrotatex"
)

// NewFastLogConfig 创建一个新的FastLogConfig实例，用于配置日志记录器。
// 参数:
//   - logDirName: 日志目录名称，默认为"applogs"。
//   - logFileName: 日志文件名称，默认为"app.log"。
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
func NewFastLogConfig(logDirName string, logFileName string) *FastLogConfig {
	// 如果日志目录名称为空，则使用默认值"logs"。
	if logDirName == "" {
		logDirName = "logs"
	}

	// 如果日志文件名称为空，则使用默认值"app.log"。
	if logFileName == "" {
		logFileName = "app.log"
	}

	// 合并日志目录和日志文件名称，生成日志文件路径。
	logFilePath := filepath.Join(logDirName, logFileName)

	// 返回一个新的FastLogConfig实例
	return &FastLogConfig{
		logDirName:     logDirName,  // 日志目录名称
		LogFileName:    logFileName, // 日志文件名称
		logFilePath:    logFilePath, // 日志文件路径 = 工作目录 + 日志目录名称 + 日志文件名称
		PrintToConsole: true,        // 是否将日志输出到控制台
		ConsoleOnly:    false,       // 是否仅输出到控制台
		LogLevel:       INFO,        // 日志级别 默认INFO
		ChanIntSize:    1000,        // 通道大小 默认1000
		LogFormat:      Detailed,    // 日志格式选项
		MaxBufferSize:  1,           // 最大缓冲区大小 默认1MB，单位为MB
		MaxLogFileSize: 10,          // 最大日志文件大小，单位为MB, 默认10MB
		MaxLogAge:      0,           // 最大日志文件保留天数, 默认为0, 表示不做限制
		MaxLogBackups:  0,           // 最大日志文件保留数量, 默认为0, 表示不做限制
		IsLocalTime:    false,       // 是否使用本地时间 默认使用UTC时间
		EnableCompress: false,       // 是否启用日志文件压缩 默认不启用
		NoColor:        false,       // 是否禁用终端颜色
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
		fileWriter    io.Writer // 文件写入器
		consoleWriter io.Writer // 控制台写入器
	)

	// 如果允许将日志输出到控制台，或者仅输出到控制台，则初始化控制台写入器。
	if config.ConsoleOnly || config.PrintToConsole {
		consoleWriter = os.Stdout // 控制台写入器
	} else {
		consoleWriter = io.Discard // 不输出到控制台，直接丢弃
	}

	// 如果不是仅输出到控制台，则初始化日志文件写入器。
	var logger *logrotatex.LogRotateX
	if !config.ConsoleOnly {
		// 检查日志目录是否存在，如果不存在则创建。
		if _, checkPathErr := checkPath(config.logDirName); checkPathErr != nil {
			if mkdirErr := os.MkdirAll(config.logDirName, 0644); mkdirErr != nil {
				return nil, fmt.Errorf("创建日志目录失败: %s", mkdirErr)
			}
		}

		// 初始化日志文件切割器
		logger = &logrotatex.LogRotateX{
			Filename:   config.logFilePath,    // 日志文件路径,
			MaxSize:    config.MaxLogFileSize, // 最大日志文件大小，单位为MB
			MaxAge:     config.MaxLogAge,      // 最大日志文件保留天数
			MaxBackups: config.MaxLogBackups,  // 最大日志文件保留数量
			LocalTime:  config.IsLocalTime,    // 是否使用本地时间
			Compress:   config.EnableCompress, // 是否启用日志文件压缩
		}

		// 初始化文件写入器
		fileWriter = logger
	} else {
		fileWriter = io.Discard // 仅输出到控制台，不输出到文件
	}

	// 创建一个新的FastLog实例，将配置和缓冲区赋值给实例。
	f := &FastLog{
		logGer:         logger,                             // 日志文件切割器
		fileWriter:     fileWriter,                         // 文件写入器,
		consoleWriter:  consoleWriter,                      // 控制台写入器,
		logFilePath:    config.logFilePath,                 // 日志文件路径
		logDirName:     config.logDirName,                  // 日志目录路径
		logFileName:    config.LogFileName,                 // 日志文件名
		printToConsole: config.PrintToConsole,              // 是否将日志输出到控制台
		consoleOnly:    config.ConsoleOnly,                 // 是否仅输出到控制台
		logLevel:       config.LogLevel,                    // 日志级别
		chanIntSize:    config.ChanIntSize,                 // 通道大小 默认1000
		logFormat:      config.LogFormat,                   // 日志格式选项
		maxBufferSize:  config.MaxBufferSize * 1024 * 1024, // 最大缓冲区大小 默认1MB
		flushInterval:  1 * time.Second,                    // 刷新间隔，单位为秒
		noColor:        config.NoColor,                     // 是否禁用终端颜色
		cl:             colorlib.NewColorLib(),             // 颜色库实例
	}

	// 根据noColor的值，设置颜色库的颜色选项
	if f.noColor {
		f.cl.NoColor = true // 设置颜色库的颜色选项为禁用
	}

	// 初始化日志通道
	f.logChan = make(chan *logMessage, f.chanIntSize)

	// 初始化文件缓冲区
	f.fileBuffer = bytes.NewBuffer(make([]byte, f.maxBufferSize))
	f.fileBuffer.Reset() // 重置缓冲区，清空内容

	// 初始化控制台缓冲区
	f.consoleBuffer = bytes.NewBuffer(make([]byte, f.maxBufferSize))
	f.consoleBuffer.Reset() // 重置缓冲区，清空内容

	// 创建 context 用于控制协程退出
	f.ctx, f.cancel = context.WithCancel(context.Background())

	// 使用 sync.Once 确保日志处理器只启动一次
	f.startOnce.Do(func() {
		go f.processLogs() // 启动日志处理器
		go f.flushBuffer() // 启动定时刷新缓冲区
	})

	// 返回FastLog实例和nil错误
	return f, nil
}

// processLogs 日志处理器，用于处理日志消息。
// 读取通道中的日志消息并将其处理为指定的日志格式，然后写入到缓冲区中。
func (f *FastLog) processLogs() {
	f.logWait.Add(1) // 增加等待组中的计数器

	// 创建一个goroutine，用于处理日志消息
	go func() {
		defer func() {
			// 减少等待组中的计数器。
			f.logWait.Done()

			// 捕获panic
			if r := recover(); r != nil {
				f.Errorf("日志处理器发生panic: %v", r)
			}
		}()

		// 计算最大缓冲区大小的80%
		maxBufferSize := int(float64(f.maxBufferSize) * 0.8)

		// 初始化控制台字符串构建器
		for {
			select {
			// 监听通道，如果通道关闭，则退出循环
			case <-f.ctx.Done():
				// 处理通道中剩余的日志消息
				for rawMsg := range f.logChan {
					f.handleLog(rawMsg, maxBufferSize)
				}
				return
			// 监听通道，如果通道中有日志消息，则处理日志消息
			case rawMsg, ok := <-f.logChan:
				if !ok {
					return
				}

				// 处理日志消息
				f.handleLog(rawMsg, maxBufferSize)
			}
		}
	}()
}

// handleLog 处理单条日志消息的逻辑
// 格式化日志消息，然后写入到缓冲区中。
// 如果缓冲区大小达到80%，则立即刷新缓冲区。
// 参数:
//   - rawMsg: 日志消息
//   - maxBufferSize: 最大缓冲区大小
//
// 返回值:
//   - 无
func (f *FastLog) handleLog(rawMsg *logMessage, maxBufferSize int) {
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
		f.fileBuilder.WriteString(fileLog)
		f.fileBuilder.WriteString("\n")
		fileLog = f.fileBuilder.String()

		// 写入文件锁
		f.fileMu.Lock()

		// 将格式化后的日志消息写入到文件缓冲区中
		if _, err := f.fileBuffer.WriteString(fileLog); err != nil {
			f.Errorf("写入文件缓冲区失败: %v", err)
		}

		// 释放文件锁
		f.fileMu.Unlock()

		// 重置构建器，以便下次使用
		f.fileBuilder.Reset()
	}

	// 如果允许将日志输出到控制台，或者仅输出到控制台
	if f.consoleOnly || f.printToConsole {
		// 调用addColor方法给日志消息添加颜色
		consoleLog := addColor(f, rawMsg, formattedLog)

		// 给格式化后的日志消息添加换行符
		f.consoleBuilder.WriteString(consoleLog)
		f.consoleBuilder.WriteString("\n")
		consoleLog = f.consoleBuilder.String()

		// 写入控制台锁
		f.consoleMu.Lock()

		// 将格式化后的日志消息写入到控制台缓冲区中
		if _, err := f.consoleBuffer.WriteString(consoleLog); err != nil {
			f.Errorf("写入控制台缓冲区失败: %v", err)
		}

		// 释放控制台锁
		f.consoleMu.Unlock()

		// 重置构建器，以便下次使用
		f.consoleBuilder.Reset()
	}
}

// flushBuffer 定时刷新缓冲区
func (f *FastLog) flushBuffer() {
	// 新增一个等待组，用于等待刷新缓冲区的协程完成
	f.logWait.Add(1)

	// 定义一个定时器，用于定时刷新缓冲区
	ticker := time.NewTicker(f.flushInterval)

	// 创建一个goroutine，用于定时刷新缓冲区
	go func() {
		defer func() {
			// 减少等待组中的计数器。
			f.logWait.Done()

			// 关闭定时器
			if ticker != nil {
				ticker.Stop()
			}

			// 捕获panic
			if r := recover(); r != nil {
				f.Errorf("刷新缓冲区发生panic: %v", r)
			}
		}()

		// 循环监听定时器
		for {
			select {
			case <-f.ctx.Done():
				return // 当 context 被取消时，退出协程
			case <-ticker.C:
				f.flushBufferNow() // 刷新缓冲区
			}
		}
	}()
}

// flushBufferNow 立即刷新缓冲区
func (f *FastLog) flushBufferNow() {

	// 如果不是仅输出到控制台，则刷新文件缓冲区
	if !f.consoleOnly {
		// 检查文件缓冲区大小是否大于0
		if f.fileBuffer.Len() > 0 {
			// 获取文件锁
			f.fileMu.Lock()

			// 获取文件缓冲区的内容
			bufferContent := f.fileBuffer.String()

			// 重置缓冲区
			f.fileBuffer.Reset()

			// 释放文件锁
			f.fileMu.Unlock()

			// 将文件缓冲区的内容写入到文件中
			if _, err := f.fileWriter.Write([]byte(bufferContent)); err != nil {
				f.Errorf("写入文件失败: %v", err)
			}
		}
	}

	// 如果允许将日志输出到控制台，或者仅输出到控制台
	if f.printToConsole || f.consoleOnly {
		// 检查控制台缓冲区大小是否大于0
		if f.consoleBuffer.Len() > 0 {
			// 获取控制台写入锁
			f.consoleMu.Lock()

			// 获取控制台缓冲区的内容
			bufferContent := f.consoleBuffer.String()

			// 重置缓冲区
			f.consoleBuffer.Reset()

			// 释放控制台写入锁
			f.consoleMu.Unlock()

			// 将控制台缓冲区的内容写入到控制台
			if _, err := f.consoleWriter.Write([]byte(bufferContent)); err != nil {
				f.Errorf("写入控制台失败: %v", err)
			}
		}
	}
}

// Close 关闭FastLog实例，并等待所有日志处理完成。
func (f *FastLog) Close() error {
	// 打印关闭日志记录器的信息
	f.Info("关闭日志记录器...")

	// 确保只关闭一次
	var closeOnce sync.Once
	var closeErr error
	closeOnce.Do(func() {
		// 关闭日志通道
		close(f.logChan)

		// 关闭协程
		f.cancel()

		// 等待所有日志处理完成
		f.logWait.Wait()

		// 刷新剩余的日志
		f.flushBufferNow()

		// 如果不是仅输出到控制台，同时日志文件句柄不为nil，则关闭日志文件。
		if !f.consoleOnly && f.logGer != nil {
			if err := f.logGer.Close(); err != nil {
				closeErr = fmt.Errorf("关闭日志文件失败: %v", err)
			}
		}

	}) // 执行一次

	// 检查是否有错误发生
	if closeErr != nil {
		return closeErr
	}

	return nil
}

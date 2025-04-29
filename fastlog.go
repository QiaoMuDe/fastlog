package fastlog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
		logDirName:                  logDirName,  // 日志目录名称
		LogFileName:                 logFileName, // 日志文件名称
		logFilePath:                 logFilePath, // 日志文件路径 = 工作目录 + 日志目录名称 + 日志文件名称
		PrintToConsole:              true,        // 是否将日志输出到控制台
		ConsoleOnly:                 false,       // 是否仅输出到控制台
		LogLevel:                    INFO,        // 日志级别 默认INFO
		ChanIntSize:                 1000,        // 通道大小 默认1000
		LogFormat:                   Detailed,    // 日志格式选项
		MaxBufferSize:               1,           // 最大缓冲区大小 默认1MB，单位为MB
		MaxLogFileSize:              1,           // 单个日志文件的最大大小，默认1MB
		MaxLogFileHour:              72,          // 单个日志文件的最大保存时间，默认72小时
		RotationCheckIntervalSecond: 3600,        // 定时检查日志轮转的间隔时间(秒) 默认1小时,
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
		if _, checkPathErr := checkPath(config.logDirName); checkPathErr != nil {
			if mkdirErr := os.MkdirAll(config.logDirName, 0644); mkdirErr != nil {
				return nil, fmt.Errorf("创建日志目录失败: %s", mkdirErr)
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
		logDirName:     config.logDirName,                  // 日志目录路径
		logFileName:    config.LogFileName,                 // 日志文件名
		printToConsole: config.PrintToConsole,              // 是否将日志输出到控制台
		consoleOnly:    config.ConsoleOnly,                 // 是否仅输出到控制台
		logLevel:       config.LogLevel,                    // 日志级别
		chanIntSize:    config.ChanIntSize,                 // 通道大小 默认1000
		logFormat:      config.LogFormat,                   // 日志格式选项
		maxBufferSize:  config.MaxBufferSize * 1024 * 1024, // 最大缓冲区大小 默认1MB
		flushInterval:  1 * time.Second,                    // 刷新间隔，单位为秒

		maxLogFileSize:              config.MaxLogFileSize * 1024 * 1024,               // 单个日志文件的最大大小，单位为字节
		maxLogFileHour:              time.Duration(config.MaxLogFileHour),              // 单个日志文件的最大保存时间，单位为小时,
		rotationCheckIntervalSecond: time.Duration(config.RotationCheckIntervalSecond), // 定时检查日志轮转的间隔时间(秒) 默认1小时,
		currentLogFileSize:          0,                                                 // 当前日志文件的大小，单位为字节
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
		// 启动日志轮转检查协程
		go f.rotateLogsPeriodically()
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

		// 将格式化后的日志消息写入到文件缓冲区中
		if _, err := f.fileBuffer.WriteString(fileLog); err != nil {
			f.Errorf("写入文件缓冲区失败: %v", err)
		}
		f.fileBuilder.Reset()
	}

	// 如果允许将日志输出到控制台，或者仅输出到控制台
	if f.consoleOnly || f.printToConsole {
		// 调用addColor方法给日志消息添加颜色
		consoleLog := addColor(rawMsg, formattedLog)

		// 给格式化后的日志消息添加换行符
		f.consoleBuilder.WriteString(consoleLog)
		f.consoleBuilder.WriteString("\n")
		consoleLog = f.consoleBuilder.String()

		// 将格式化后的日志消息写入到控制台缓冲区中
		if _, err := f.consoleBuffer.WriteString(consoleLog); err != nil {
			f.Errorf("写入控制台缓冲区失败: %v", err)
		}
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
	// 写入文件
	if f.fileBuffer.Len() > 0 {
		// 获取文件写入写入锁，确保线程安全
		f.fileMu.Lock()
		defer f.fileMu.Unlock()

		// 获取文件缓冲区的内容
		bufferContent := f.fileBuffer.String()

		// 将文件缓冲区的内容写入到文件中
		if _, err := f.fileWriter.Write([]byte(bufferContent)); err != nil {
			f.Errorf("写入文件失败: %v", err)
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
			f.Errorf("写入控制台失败: %v", err)
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

		// 关闭协程
		f.cancel()

		// 等待所有日志处理完成
		f.logWait.Wait()

		// 刷新剩余的日志
		f.flushBufferNow()

		// 如果不是仅输出到控制台，同时日志文件句柄不为nil，则关闭日志文件。
		if !f.consoleOnly && f.logFile != nil {
			f.logFile.Close()
		}

	}) // 执行一次
}

// rotateLogsPeriodically 定时检查日志文件是否需要轮转和清理
func (f *FastLog) rotateLogsPeriodically() {
	// 新增一个等待组，用于等待刷新缓冲区的协程完成
	f.logWait.Add(1)

	go func() {
		defer func() {
			// 减少等待组中的计数器。
			f.logWait.Done()

			// 捕获panic
			if r := recover(); r != nil {
				f.Errorf("日志轮转发生panic: %v", r)
			}
		}()

		// 创建一个定时器，用于定时检查日志文件是否需要轮转和清理
		ticker := time.NewTicker(time.Duration(f.rotationCheckIntervalSecond))
		defer ticker.Stop()

		// 检查是否为控制台输出
		if f.consoleOnly {
			f.Errorf("当前仅为控制台输出, 无需进行日志轮转")
			return
		}

		// 循环监听定时器
		for {
			select {
			case <-f.ctx.Done():
				return
			case <-ticker.C:
				f.checkAndRotateLogs() // 检查并轮转日志文件
				f.cleanupOldLogs()     // 清理旧的日志文件
			}
		}
	}()
}

// checkAndRotateLogs 检查当前日志文件大小是否超过最大限制，如果超过则进行日志轮转
func (f *FastLog) checkAndRotateLogs() {
	f.fileMu.Lock()
	defer f.fileMu.Unlock()

	// 获取当前日志文件的大小
	fileInfo, err := f.logFile.Stat()
	if err != nil {
		f.Errorf("获取日志文件信息失败: %v", err)
		return
	}
	f.currentLogFileSize = fileInfo.Size()

	// 检查是否需要轮转
	if f.currentLogFileSize >= f.maxLogFileSize {
		f.rotateLogFile() // 进行日志轮转切割
	}
}

// cleanupOldLogs 清理超过最大保留天数的日志文件
func (f *FastLog) cleanupOldLogs() {
	f.fileMu.Lock()
	defer f.fileMu.Unlock()

	// 切换到日志目录
	if err := os.Chdir(f.logDirName); err != nil {
		f.Errorf("切换到日志目录失败: %v", err)
		return
	}

	// 通过正则表达式获取日志文件列表
	files, err := filepath.Glob(`^.*_.*\.log$`)
	if err != nil {
		f.Errorf("获取日志文件列表失败: %v", err)
		return
	}

	// 获取当前时间
	now := time.Now()

	// 遍历日志目录下的所有文件
	for _, file := range files {
		// 如果是目录，则跳过
		fileInfo, statErr := os.Stat(file)
		if statErr == nil && fileInfo.IsDir() {
			// 跳过目录
			continue
		} else if statErr != nil {
			// 处理文件不存在的情况
			f.Errorf("获取文件信息失败: %v", statErr)
			continue
		}

		// 获取文件的修改时间
		modTime := fileInfo.ModTime()

		// 检查文件是否超过最大保留天数
		if now.Sub(modTime) > f.maxLogFileHour {
			// 删除超过最大保留天数的文件
			if err := os.Remove(file); err != nil {
				f.Errorf("删除日志文件失败: %v", err)
			}
		}
	}
}

// rotateLogFile 进行日志轮转切割
func (f *FastLog) rotateLogFile() {
	// 检查是否仅为控制台输出
	if f.consoleOnly {
		f.Errorf("当前仅为控制台输出, 无需进行日志轮转")
		return
	}

	// 获取文件写入写入锁，确保线程安全
	f.fileMu.Lock()
	defer f.fileMu.Unlock()

	// 关闭当前日志文件
	if f.logFile != nil {
		f.logFile.Close()
	}

	// 检查父目录是否存在，如果不存在则创建
	if _, err := os.Stat(f.logDirName); os.IsNotExist(err) {
		if err := os.MkdirAll(f.logDirName, 0644); err != nil {
			f.Errorf("创建日志目录失败: %s", err)
			return
		}
	}

	// 检查日志文件名是否为.log结尾，如果是则提取文件名部分
	var logFileName string
	if filepath.Ext(f.logFileName) == ".log" {
		// 通过.分割字符串，获取文件名部分
		logFileName = strings.TrimSuffix(f.logFileName, ".log")
	} else {
		// 直接使用文件名
		logFileName = f.logFileName
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102150405")
	backupFileName := fmt.Sprintf("%s_%s.log", logFileName, timestamp)

	// 重命名当前日志文件为备份文件
	if _, err := os.Stat(f.logFilePath); err != nil {
		f.Errorf("日志文件不存在: %s,日志轮转失败", f.logFilePath)
		return
	}
	if err := os.Rename(f.logFilePath, filepath.Join(f.logDirName, backupFileName)); err != nil {
		f.Errorf("重命名日志文件失败: %s,日志轮转失败", err)
		return
	}

	// 创建新的日志文件
	logFile, err := os.OpenFile(f.logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		f.Errorf("创建新的日志文件失败: %s,日志轮转失败", err)
		return
	}

	// 更新日志文件句柄和当前日志文件大小
	f.logFile = logFile
	f.fileWriter = logFile
	f.currentLogFileSize = 0

	// 打印日志轮转成功信息
	f.Successf("日志文件轮转成功: %s -> %s", f.logFileName, backupFileName)
}

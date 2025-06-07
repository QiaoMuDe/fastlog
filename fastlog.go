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

// NewFastLogConfig 创建一个新的FastLogConfig实例, 用于配置日志记录器。
// 参数:
//   - logDirName: 日志目录名称, 默认为"applogs"。
//   - logFileName: 日志文件名称, 默认为"app.log"。
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
func NewFastLogConfig(logDirName string, logFileName string) *FastLogConfig {
	// 如果日志目录名称为空, 则使用默认值"logs"。
	if logDirName == "" {
		logDirName = "logs"
	}

	// 如果日志文件名称为空, 则使用默认值"app.log"。
	if logFileName == "" {
		logFileName = "app.log"
	}

	// 返回一个新的FastLogConfig实例
	return &FastLogConfig{
		LogDirName:     logDirName,             // 日志目录名称
		LogFileName:    logFileName,            // 日志文件名称
		PrintToConsole: true,                   // 是否将日志输出到控制台
		ConsoleOnly:    false,                  // 是否仅输出到控制台
		LogLevel:       INFO,                   // 日志级别 默认INFO
		ChanIntSize:    10000,                  // 通道大小 增加到10000
		FlushInterval:  500 * time.Millisecond, // 刷新间隔 缩短到500毫秒
		LogFormat:      Detailed,               // 日志格式选项
		MaxBufferSize:  1 * 1024 * 1024,        // 最大缓冲区大小 默认1MB, 单位为MB
		MaxLogFileSize: 5,                      // 最大日志文件大小, 单位为MB, 默认5MB
		MaxLogAge:      0,                      // 最大日志文件保留天数, 默认为0, 表示不做限制
		MaxLogBackups:  0,                      // 最大日志文件保留数量, 默认为0, 表示不做限制
		IsLocalTime:    false,                  // 是否使用本地时间 默认使用UTC时间
		EnableCompress: false,                  // 是否启用日志文件压缩 默认不启用
		NoColor:        false,                  // 是否禁用终端颜色
		NoBold:         false,                  // 是否禁用终端字体加粗
	}
}

// NewFastLog 创建一个新的FastLog实例, 用于记录日志。
// 参数:
//   - config: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。
//
// 返回值:
//   - *FastLog: 一个指向FastLog实例的指针。
//   - error: 如果创建日志记录器失败, 则返回一个错误。
func NewFastLog(config *FastLogConfig) (*FastLog, error) {
	// 检查配置结构体是否为nil
	if config == nil {
		return nil, fmt.Errorf("FastLogConfig 为 nil")
	}

	// 声明一些配置变量
	var (
		fileWriter    io.Writer // 文件写入器
		consoleWriter io.Writer // 控制台写入器
	)

	// 如果允许将日志输出到控制台, 或者仅输出到控制台, 则初始化控制台写入器。
	if config.ConsoleOnly || config.PrintToConsole {
		consoleWriter = os.Stdout // 控制台写入器
	} else {
		consoleWriter = io.Discard // 不输出到控制台, 直接丢弃
	}

	// 拼接日志文件路径
	var logFilePath string
	// 如果日志目录名称和日志文件名称都不为空, 并且不是仅输出到控制台, 则拼接日志文件路径。
	if config.LogDirName != "" || config.LogFileName != "" && !config.ConsoleOnly {
		logFilePath = filepath.Join(config.LogDirName, config.LogFileName)
	}

	// 如果不是仅输出到控制台, 则初始化日志文件写入器。
	var logger *logrotatex.LogRotateX
	if !config.ConsoleOnly {
		// 检查日志目录是否存在, 如果不存在则创建。
		if _, checkPathErr := checkPath(config.LogDirName); checkPathErr != nil {
			if mkdirErr := os.MkdirAll(config.LogDirName, 0644); mkdirErr != nil {
				return nil, fmt.Errorf("创建日志目录失败: %s", mkdirErr)
			}
		}

		// 初始化日志文件切割器
		logger = &logrotatex.LogRotateX{
			Filename:   logFilePath,           // 日志文件路径,
			MaxSize:    config.MaxLogFileSize, // 最大日志文件大小, 单位为MB
			MaxAge:     config.MaxLogAge,      // 最大日志文件保留天数
			MaxBackups: config.MaxLogBackups,  // 最大日志文件保留数量
			LocalTime:  config.IsLocalTime,    // 是否使用本地时间
			Compress:   config.EnableCompress, // 是否启用日志文件压缩
		}

		// 初始化文件写入器
		fileWriter = logger
	} else {
		fileWriter = io.Discard // 仅输出到控制台, 不输出到文件
	}

	// 初始化双缓冲区
	fileBuffers := [2]*bytes.Buffer{
		bytes.NewBuffer(make([]byte, config.MaxBufferSize)),
		bytes.NewBuffer(make([]byte, config.MaxBufferSize)),
	}
	consoleBuffers := [2]*bytes.Buffer{
		bytes.NewBuffer(make([]byte, config.MaxBufferSize)),
		bytes.NewBuffer(make([]byte, config.MaxBufferSize)),
	}

	// 清空缓冲区
	fileBuffers[0].Reset()
	fileBuffers[1].Reset()
	consoleBuffers[0].Reset()
	consoleBuffers[1].Reset()

	// 创建一个新的FastLog实例, 将配置和缓冲区赋值给实例。
	f := &FastLog{
		logGer:         logger,                                     // 日志文件切割器
		fileWriter:     fileWriter,                                 // 文件写入器,
		consoleWriter:  consoleWriter,                              // 控制台写入器,
		logFilePath:    logFilePath,                                // 日志文件路径
		cl:             colorlib.NewColorLib(),                     // 颜色库实例
		config:         config,                                     // 配置结构体
		logChan:        make(chan *logMessage, config.ChanIntSize), // 日志消息通道
		fileBuffers:    fileBuffers,                                // 文件缓冲
		consoleBuffers: consoleBuffers,                             // 控制台缓冲
	}

	// 根据noColor的值, 设置颜色库的颜色选项
	if f.config.NoColor {
		f.cl.NoColor.Store(true) // 设置颜色库的颜色选项为禁用
	}

	// 根据noBold的值, 设置颜色库的字体加粗选项
	if f.config.NoBold {
		f.cl.NoBold.Store(true) // 设置颜色库的字体加粗选项为禁用
	}

	// 设置缓冲区索引为0
	f.fileBufferIdx.Store(0)
	f.consoleBufferIdx.Store(0)

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

// processLogs 日志处理器, 用于处理日志消息。
// 读取通道中的日志消息并将其处理为指定的日志格式, 然后写入到缓冲区中。
func (f *FastLog) processLogs() {
	f.logWait.Add(1) // 增加等待组中的计数器

	// 创建一个goroutine, 用于处理日志消息
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
		maxBufferSize := int(float64(f.config.MaxBufferSize) * 0.8)

		// 初始化控制台字符串构建器
		for {
			select {
			// 监听通道, 如果通道关闭, 则退出循环
			case <-f.ctx.Done():
				// 检查通道是否为空
				if len(f.logChan) != 0 {
					// 处理通道中剩余的日志消息
					for rawMsg := range f.logChan {
						f.handleLog(rawMsg, maxBufferSize)
					}
				}
				return
			// 监听通道, 如果通道中有日志消息, 则处理日志消息
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
// 格式化日志消息, 然后写入到缓冲区中。
// 如果缓冲区大小达到80%, 则立即刷新缓冲区。
// 参数:
//   - rawMsg: 日志消息
//   - maxBufferSize: 最大缓冲区大小
//
// 返回值:
//   - 无
func (f *FastLog) handleLog(rawMsg *logMessage, maxBufferSize int) {
	// 获取文件缓冲区的索引和文件缓冲区的大小
	f.fileBufferMu.Lock()
	currentFileBufferIdx := f.fileBufferIdx.Load()
	currentFileBufferSize := f.fileBuffers[currentFileBufferIdx].Len()
	f.fileBufferMu.Unlock()

	// 获取控制台缓冲区的索引和控制台缓冲区的大小
	f.consoleBufferMu.Lock()
	currentConsoleBufferIdx := f.consoleBufferIdx.Load()
	currentConsoleBufferSize := f.consoleBuffers[currentConsoleBufferIdx].Len()
	f.consoleBufferMu.Unlock()

	// 检查缓冲区大小是否达到80%
	if currentFileBufferSize >= maxBufferSize || currentConsoleBufferSize >= maxBufferSize {
		f.flushBufferNow()
	}

	// 格式化日志消息
	formattedLog := formatLog(f, rawMsg)

	// 写入文件缓冲区
	if !f.config.ConsoleOnly {
		f.fileBufferMu.Lock()
		defer f.fileBufferMu.Unlock()
		// 写入文件缓冲区
		f.fileBuffers[currentFileBufferIdx].WriteString(formattedLog + "\n")
	}

	// 写入控制台缓冲区
	if f.config.PrintToConsole || f.config.ConsoleOnly {
		f.consoleBufferMu.Lock()
		defer f.consoleBufferMu.Unlock()

		// 渲染颜色
		consoleLog := addColor(f, rawMsg, formattedLog)

		// 写入控制台缓冲区
		f.consoleBuffers[currentConsoleBufferIdx].WriteString(consoleLog + "\n")
	}
}

// flushBuffer 定时刷新缓冲区
func (f *FastLog) flushBuffer() {
	// 新增一个等待组, 用于等待刷新缓冲区的协程完成
	f.logWait.Add(1)

	// 定义一个定时器, 用于定时刷新缓冲区
	ticker := time.NewTicker(f.config.FlushInterval)

	// 创建一个goroutine, 用于定时刷新缓冲区
	go func() {
		defer func() {
			// 捕获panic
			if r := recover(); r != nil {
				f.Errorf("刷新缓冲区发生panic: %v", r)
			}

			// 减少等待组中的计数器。
			f.logWait.Done()

			// 关闭定时器
			if ticker != nil {
				ticker.Stop()
			}
		}()

		// 循环监听定时器
		for {
			select {
			case <-f.ctx.Done():
				return // 当 context 被取消时, 退出协程
			case <-ticker.C:
				f.flushBufferNow() // 刷新缓冲区
			}
		}
	}()
}

// flushBufferNow 立即刷新缓冲区
func (f *FastLog) flushBufferNow() {
	f.flushLock.Lock() // 加锁, 防止并发刷新
	defer f.flushLock.Unlock()

	// 判断是否需要刷新文件缓冲区
	if !f.config.ConsoleOnly {
		f.fileBufferMu.Lock()
		defer f.fileBufferMu.Unlock()

		// 获取当前缓冲区索引
		currentIdx := f.fileBufferIdx.Load()

		// 检查当前缓冲区是否有内容
		if f.fileBuffers[currentIdx].Len() > 0 {
			// 切换到另一个缓冲区
			newIdx := 1 - currentIdx // 0 -> 1, 1 -> 0
			f.fileBufferIdx.Store(newIdx)

			// 将文件缓冲区中的内容写入文件
			f.fileMu.Lock()
			defer f.fileMu.Unlock()
			if _, err := f.fileWriter.Write(f.fileBuffers[currentIdx].Bytes()); err != nil {
				f.Errorf("写入文件失败: %v", err)
			}

			// 重置当前缓冲区
			f.fileBuffers[currentIdx].Reset()
		}
	}

	// 判断是否需要刷新控制台缓冲区
	if f.config.PrintToConsole || f.config.ConsoleOnly {
		f.consoleBufferMu.Lock()
		defer f.consoleBufferMu.Unlock()

		// 获取当前缓冲区索引
		currentIdx := f.consoleBufferIdx.Load()

		// 检查当前缓冲区是否有内容
		if f.consoleBuffers[currentIdx].Len() > 0 {
			// 切换到另一个缓冲区
			newIdx := 1 - currentIdx // 0  -> 1, 1 -> 0
			f.consoleBufferIdx.Store(newIdx)

			// 将控制台缓冲区中的内容写入控制台
			f.consoleMu.Lock()
			defer f.consoleMu.Unlock()
			if _, err := f.consoleWriter.Write(f.consoleBuffers[currentIdx].Bytes()); err != nil {
				f.Errorf("写入控制台失败: %v", err)
			}

			// 重置当前缓冲区
			f.consoleBuffers[currentIdx].Reset()
		}
	}
}

// Close 关闭FastLog实例, 并等待所有日志处理完成。
func (f *FastLog) Close() error {
	f.closeLock.Lock()
	defer f.closeLock.Unlock()

	// 打印关闭日志记录器的信息
	f.Info("stop logging...")

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

		// 刷新剩余的日志 缓冲区1
		f.flushBufferNow()

		// 刷新剩余的日志 缓冲区2
		f.flushBufferNow()

		// 如果不是仅输出到控制台, 同时日志文件句柄不为nil, 则关闭日志文件。
		if !f.config.ConsoleOnly && f.logGer != nil {
			f.fileMu.Lock()
			defer f.fileMu.Unlock()
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

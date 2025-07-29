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

// 简写实现
var (
	NewFlog = NewFastLog
	NewFcfg = NewFastLogConfig
)

// NewFastLogConfig 创建一个新的FastLogConfig实例, 用于配置日志记录器。
//
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
		logDirName:     logDirName,             // 日志目录名称
		logFileName:    logFileName,            // 日志文件名称
		printToConsole: true,                   // 是否将日志输出到控制台
		consoleOnly:    false,                  // 是否仅输出到控制台
		logLevel:       INFO,                   // 日志级别 默认INFO
		chanIntSize:    10000,                  // 通道大小 增加到10000
		flushInterval:  500 * time.Millisecond, // 刷新间隔 缩短到500毫秒
		logFormat:      Detailed,               // 日志格式选项
		maxLogFileSize: 5,                      // 最大日志文件大小, 单位为MB, 默认5MB
		maxLogAge:      0,                      // 最大日志文件保留天数, 默认为0, 表示不做限制
		maxLogBackups:  0,                      // 最大日志文件保留数量, 默认为0, 表示不做限制
		isLocalTime:    false,                  // 是否使用本地时间 默认使用UTC时间
		enableCompress: false,                  // 是否启用日志文件压缩 默认不启用
		noColor:        false,                  // 是否禁用终端颜色
		noBold:         false,                  // 是否禁用终端字体加粗
	}
}

// NewFastLog 创建一个新的FastLog实例, 用于记录日志。
//
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
	if config.GetConsoleOnly() || config.GetPrintToConsole() {
		consoleWriter = os.Stdout // 控制台写入器
	} else {
		consoleWriter = io.Discard // 不输出到控制台, 直接丢弃
	}

	// 拼接日志文件路径
	var logFilePath string
	// 如果日志目录名称和日志文件名称都不为空, 并且不是仅输出到控制台, 则拼接日志文件路径。
	if config.GetLogDirName() != "" || config.GetLogFileName() != "" && !config.GetConsoleOnly() {
		logFilePath = filepath.Join(config.GetLogDirName(), config.GetLogFileName())
	}

	// 如果不是仅输出到控制台, 则初始化日志文件写入器。
	var logger *logrotatex.LogRotateX
	if !config.GetConsoleOnly() {
		// 检查日志目录是否存在, 如果不存在则创建。
		if _, checkPathErr := checkPath(config.GetLogDirName()); checkPathErr != nil {
			if mkdirErr := os.MkdirAll(config.GetLogDirName(), 0755); mkdirErr != nil {
				return nil, fmt.Errorf("创建日志目录失败: %s", mkdirErr)
			}
		}

		// 初始化日志文件切割器
		logger = &logrotatex.LogRotateX{
			Filename:   logFilePath,                // 日志文件路径,
			MaxSize:    config.GetMaxLogFileSize(), // 最大日志文件大小, 单位为MB
			MaxAge:     config.GetMaxLogAge(),      // 最大日志文件保留天数
			MaxBackups: config.GetMaxLogBackups(),  // 最大日志文件保留数量
			LocalTime:  config.GetIsLocalTime(),    // 是否使用本地时间
			Compress:   config.GetEnableCompress(), // 是否启用日志文件压缩
		}

		// 初始化文件写入器
		fileWriter = logger
	} else {
		fileWriter = io.Discard // 仅输出到控制台, 不输出到文件
	}

	// 创建双缓冲区
	fileBuffers := [2]*bytes.Buffer{
		bytes.NewBuffer(nil),
		bytes.NewBuffer(nil),
	}
	consoleBuffers := [2]*bytes.Buffer{
		bytes.NewBuffer(nil),
		bytes.NewBuffer(nil),
	}

	// 预留合理的初始容量，避免频繁扩容
	for i := range fileBuffers {
		fileBuffers[i].Grow(initialBufferCapacity)
		consoleBuffers[i].Grow(initialBufferCapacity)
	}

	// 创建一个新的FastLog实例, 将配置和缓冲区赋值给实例。
	f := &FastLog{
		logGer:         logger,                                          // 日志文件切割器
		fileWriter:     fileWriter,                                      // 文件写入器,
		consoleWriter:  consoleWriter,                                   // 控制台写入器,
		logFilePath:    logFilePath,                                     // 日志文件路径
		cl:             colorlib.NewColorLib(),                          // 颜色库实例
		config:         config,                                          // 配置结构体
		logChan:        make(chan *logMessage, config.GetChanIntSize()), // 日志消息通道
		fileBuffers:    fileBuffers,                                     // 文件缓冲
		consoleBuffers: consoleBuffers,                                  // 控制台缓冲
		closeOnce:      sync.Once{},                                     // 用于在结束时确保只执行一次
	}

	// 根据noColor的值, 设置颜色库的颜色选项
	if f.config.GetNoColor() {
		f.cl.NoColor.Store(true) // 设置颜色库的颜色选项为禁用
	}

	// 根据noBold的值, 设置颜色库的字体加粗选项
	if f.config.GetNoBold() {
		f.cl.NoBold.Store(true) // 设置颜色库的字体加粗选项为禁用
	}

	// 设置缓冲区索引为0
	f.fileBufferIdx.Store(0)
	f.consoleBufferIdx.Store(0)

	// 创建 context 用于控制协程退出
	f.ctx, f.cancel = context.WithCancel(context.Background())

	// 使用 sync.Once 确保日志处理器只启动一次
	var startErr error
	f.startOnce.Do(func() {
		// 启动日志处理器和刷新器
		defer func() {
			if r := recover(); r != nil {
				startErr = fmt.Errorf("failed to start log processor: %v", r)
			}
		}()

		go f.processLogs() // 启动日志处理器
		go f.flushBuffer() // 启动定时刷新缓冲区
	})

	// 检查启动是否成功
	if startErr != nil {
		return nil, startErr
	}

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

		// 初始化控制台字符串构建器
		for {
			select {
			// 监听通道, 如果通道关闭, 则退出循环
			case <-f.ctx.Done():
				// 检查通道是否为空
				if len(f.logChan) != 0 {
					// 处理通道中剩余的日志消息
					for rawMsg := range f.logChan {
						f.handleLog(rawMsg)
					}
				}
				return
			// 监听通道, 如果通道中有日志消息, 则处理日志消息
			case rawMsg, ok := <-f.logChan:
				if !ok {
					return
				}

				// 处理日志消息
				f.handleLog(rawMsg)
			}
		}
	}()
}

// handleLog 处理单条日志消息的逻辑
// 格式化日志消息, 然后写入到缓冲区中。
// 如果缓冲区大小达到90%, 则立即刷新缓冲区。

// 参数:
//   - rawMsg: 日志消息

// 返回值:
//   - 无
// func (f *FastLog) handleLog(rawMsg *logMessage) {
// 	// 获取文件缓冲区的索引和文件缓冲区的大小
// 	f.fileBufferMu.Lock()
// 	currentFileBufferIdx := f.fileBufferIdx.Load()
// 	currentFileBufferSize := f.fileBuffers[currentFileBufferIdx].Len()
// 	f.fileBufferMu.Unlock()

// 	// 获取控制台缓冲区的索引和控制台缓冲区的大小
// 	f.consoleBufferMu.Lock()
// 	currentConsoleBufferIdx := f.consoleBufferIdx.Load()
// 	currentConsoleBufferSize := f.consoleBuffers[currentConsoleBufferIdx].Len()
// 	f.consoleBufferMu.Unlock()

// 	// 检查缓冲区大小是否达到90%
// 	if currentFileBufferSize >= flushThreshold || currentConsoleBufferSize >= flushThreshold {
// 		f.flushBufferNow()
// 	}

// 	// 格式化日志消息
// 	formattedLog := formatLog(f, rawMsg)

// 	// 写入文件缓冲区
// 	if !f.config.GetConsoleOnly() {
// 		f.fileBufferMu.Lock() // 锁定文件缓冲区

// 		// 写入文件缓冲区
// 		if f.fileBuffers[currentFileBufferIdx] != nil {
// 			f.fileBuffers[currentFileBufferIdx].WriteString(formattedLog + "\n")
// 		}

// 		f.fileBufferMu.Unlock() // 解锁文件缓冲区
// 	}

// 	// 写入控制台缓冲区
// 	if f.config.GetPrintToConsole() || f.config.GetConsoleOnly() {
// 		f.consoleBufferMu.Lock() // 锁定控制台缓冲区

// 		// 渲染颜色
// 		consoleLog := addColor(f, rawMsg, formattedLog)

// 		// 写入控制台缓冲区
// 		if f.consoleBuffers[currentConsoleBufferIdx] != nil {
// 			f.consoleBuffers[currentConsoleBufferIdx].WriteString(consoleLog + "\n")
// 		}

//			f.consoleBufferMu.Unlock() // 解锁控制台缓冲区
//		}
//	}
func (f *FastLog) handleLog(rawMsg *logMessage) {
	// 格式化日志消息
	formattedLog := formatLog(f, rawMsg)

	// 局部变量记录是否需要刷新
	var needFileFlush bool
	var needConsoleFlush bool

	// 处理文件缓冲区 - 在同一个锁内完成检查和写入
	if !f.config.GetConsoleOnly() {
		f.fileBufferMu.Lock()

		// 记录当前文件缓冲区的索引和大小
		currentFileBufferIdx := f.fileBufferIdx.Load()
		currentFileBufferSize := f.fileBuffers[currentFileBufferIdx].Len()

		// 检查是否需要刷新（写入前检查）
		needFileFlush = currentFileBufferSize >= flushThreshold

		// 写入文件缓冲区
		if f.fileBuffers[currentFileBufferIdx] != nil {
			f.fileBuffers[currentFileBufferIdx].WriteString(formattedLog + "\n")
		}
		f.fileBufferMu.Unlock()
	}

	// 处理控制台缓冲区 - 同样的逻辑
	if f.config.GetPrintToConsole() || f.config.GetConsoleOnly() {
		f.consoleBufferMu.Lock()

		// 记录当前的控制台缓冲区的索引和大小
		currentConsoleBufferIdx := f.consoleBufferIdx.Load()
		currentConsoleBufferSize := f.consoleBuffers[currentConsoleBufferIdx].Len()

		// 检查是否需要刷新（写入前检查）
		needConsoleFlush = currentConsoleBufferSize >= flushThreshold

		// 渲染颜色
		consoleLog := addColor(f, rawMsg, formattedLog)

		// 写入控制台缓冲区
		if f.consoleBuffers[currentConsoleBufferIdx] != nil {
			f.consoleBuffers[currentConsoleBufferIdx].WriteString(consoleLog + "\n")
		}
		f.consoleBufferMu.Unlock()
	}

	// 统一判断是否需要刷新（在所有写入完成后）
	if needFileFlush || needConsoleFlush {
		f.flushBufferNow()
	}
}

// flushBuffer 定时刷新缓冲区
func (f *FastLog) flushBuffer() {
	// 新增一个等待组, 用于等待刷新缓冲区的协程完成
	f.logWait.Add(1)

	// 定义一个定时器, 用于定时刷新缓冲区
	ticker := time.NewTicker(f.config.GetFlushInterval())

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
// 刷新缓冲区, 并将控制台缓冲区的内容输出到控制台
// 将控制台缓冲区的内容输出到控制台后, 将控制台缓冲区清空。
func (f *FastLog) flushBufferNow() {
	f.flushLock.Lock() // 加锁, 防止并发刷新
	defer f.flushLock.Unlock()

	// 判断是否需要刷新文件缓冲区
	if !f.config.GetConsoleOnly() {
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
	if f.config.GetPrintToConsole() || f.config.GetConsoleOnly() {
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
//
// 返回值:
//   - 关闭过程中可能发生的错误
func (f *FastLog) Close() (closeErr error) {
	f.closeLock.Lock()
	defer f.closeLock.Unlock()

	// 确保只关闭一次
	f.closeOnce.Do(func() {
		// 等待日志处理完成 - 确保所有日志都被消费
		done := make(chan struct{})

		// 启动一个goroutine检查日志是否处理完成
		go func() {
			// 等待通道为空（所有日志都被处理）
			for len(f.logChan) > 0 {
				// 休眠10毫秒
				time.Sleep(10 * time.Millisecond)
			}
			// 再等待一个刷新周期确保写入文件
			time.Sleep(f.config.GetFlushInterval())
			close(done)
		}()

		// 等待完成信号，但设置超时避免无限等待
		select {
		case <-done:
			// 日志处理完成
		case <-time.After(3 * time.Second):
			// 超时保护，避免无限等待
		}

		// 打印关闭日志记录器的信息
		f.Info("stop logging...")

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

		// 清理缓冲区资源
		f.cleanupBuffers()

		// 如果不是仅输出到控制台, 同时日志文件句柄不为nil, 则关闭日志文件。
		if !f.config.GetConsoleOnly() && f.logGer != nil {
			f.fileMu.Lock()
			defer f.fileMu.Unlock()
			if err := f.logGer.Close(); err != nil {
				closeErr = fmt.Errorf("关闭日志文件失败: %v", err)
			}
		}

	}) // 执行一次

	// 返回关闭错误
	if closeErr != nil {
		return closeErr
	}
	return nil
}

// cleanupBuffers 清理缓冲区资源，释放内存
func (f *FastLog) cleanupBuffers() {
	// 清理文件缓冲区
	f.fileBufferMu.Lock()
	for i := range f.fileBuffers {
		if f.fileBuffers[i] != nil {
			f.fileBuffers[i].Reset()
			f.fileBuffers[i] = nil // 帮助GC回收
		}
	}
	f.fileBufferMu.Unlock()

	// 清理控制台缓冲区
	f.consoleBufferMu.Lock()
	for i := range f.consoleBuffers {
		if f.consoleBuffers[i] != nil {
			f.consoleBuffers[i].Reset()
			f.consoleBuffers[i] = nil // 帮助GC回收
		}
	}
	f.consoleBufferMu.Unlock()
}

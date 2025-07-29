// fastlog.go 存放fastlog包的实现
package fastlog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/logrotatex"
)

// 简写实现
var (
	NewFlog = NewFastLog
	NewFcfg = NewFastLogConfig
)

// FastLog 日志记录器
type FastLog struct {
	/* 私有属性 */
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
	// 用于控制关闭过程的锁
	closeLock sync.Mutex
	// logrotatex 日志文件切割
	logGer *logrotatex.LogRotateX

	/* 双缓冲区属性 */
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

	// 嵌入的配置结构体
	config *FastLogConfig
}

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

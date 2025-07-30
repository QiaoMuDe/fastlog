// fastlog.go - FastLog日志记录器核心实现
// 提供日志记录器的创建、初始化、日志写入及关闭等核心功能，
// 集成配置管理、缓冲区管理和日志处理流程。
package fastlog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
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

// FastLog 日志记录器
type FastLog struct {
	// 日志通道  用于异步写入日志文件
	logChan chan *logMessage

	// 等待组 用于等待所有goroutine完成
	logWait sync.WaitGroup

	// 文件写入器
	fileWriter io.Writer

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

	// logrotatex 日志文件切割
	logGer *logrotatex.LogRotateX

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
	// 返回一个新的FastLogConfig实例
	return &FastLogConfig{
		LogDirName:      logDirName,             // 日志目录名称
		LogFileName:     logFileName,            // 日志文件名称
		OutputToConsole: true,                   // 是否将日志输出到控制台
		OutputToFile:    true,                   // 是否将日志输出到文件
		LogLevel:        INFO,                   // 日志级别 默认INFO
		ChanIntSize:     10000,                  // 通道大小 增加到10000
		FlushInterval:   500 * time.Millisecond, // 刷新间隔 缩短到500毫秒
		LogFormat:       Detailed,               // 日志格式选项
		MaxLogFileSize:  5,                      // 最大日志文件大小, 单位为MB, 默认5MB
		MaxLogAge:       0,                      // 最大日志文件保留天数, 默认为0, 表示不做限制
		MaxLogBackups:   0,                      // 最大日志文件保留数量, 默认为0, 表示不做限制
		IsLocalTime:     false,                  // 是否使用本地时间 默认使用UTC时间
		EnableCompress:  false,                  // 是否启用日志文件压缩 默认不启用
		NoColor:         false,                  // 是否禁用终端颜色
		NoBold:          false,                  // 是否禁用终端字体加粗
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
		panic("FastLogConfig 不能为 nil")
	}

	// 克隆配置结构体防止原配置被意外修改
	cfg := &FastLogConfig{
		LogDirName:      config.LogDirName,      // 日志目录名称
		LogFileName:     config.LogFileName,     // 日志文件名称
		OutputToConsole: config.OutputToConsole, // 是否将日志输出到控制台
		OutputToFile:    config.OutputToFile,    // 是否将日志输出到文件
		LogLevel:        config.LogLevel,        // 日志级别
		ChanIntSize:     config.ChanIntSize,     // 通道大小
		FlushInterval:   config.FlushInterval,   // 刷新间隔
		LogFormat:       config.LogFormat,       // 日志格式
		MaxLogFileSize:  config.MaxLogFileSize,  // 最大日志文件大小, 单位为MB
		MaxLogAge:       config.MaxLogAge,       // 最大日志文件保留天数(单位为天)
		MaxLogBackups:   config.MaxLogBackups,   // 最大日志文件保留数量(默认为0, 表示不清理)
		IsLocalTime:     config.IsLocalTime,     // 是否使用本地时间
		EnableCompress:  config.EnableCompress,  // 是否启用日志文件压缩
		NoColor:         config.NoColor,         // 是否禁用终端颜色
		NoBold:          config.NoBold,          // 是否禁用终端字体加粗
	}

	// 最终配置修正 - 修正所有不合理的值
	cfg.fixFinalConfig()

	// 初始化写入器
	var fileWriter io.Writer    // 文件写入器
	var consoleWriter io.Writer // 控制台写入器

	// 如果允许将日志输出到控制台, 则初始化控制台写入器。
	if cfg.OutputToConsole {
		consoleWriter = os.Stdout // 控制台写入器
	} else {
		consoleWriter = io.Discard // 不允许将日志输出到控制台, 则直接丢弃
	}

	// 如果允许将日志输出到文件, 则初始化日志文件写入器。
	var logger *logrotatex.LogRotateX
	if cfg.OutputToFile {
		// 拼接日志文件路径
		logFilePath := filepath.Join(cfg.LogDirName, cfg.LogFileName)

		// 检查日志目录是否存在, 如果不存在则创建。
		if _, checkPathErr := checkPath(cfg.LogDirName); checkPathErr != nil {
			if mkdirErr := os.MkdirAll(cfg.LogDirName, 0755); mkdirErr != nil {
				return nil, fmt.Errorf("创建日志目录失败: %s", mkdirErr)
			}
		}

		// 初始化日志文件切割器
		logger = &logrotatex.LogRotateX{
			Filename:   logFilePath,        // 日志文件路径,
			MaxSize:    cfg.MaxLogFileSize, // 最大日志文件大小, 单位为MB
			MaxAge:     cfg.MaxLogAge,      // 最大日志文件保留天数
			MaxBackups: cfg.MaxLogBackups,  // 最大日志文件保留数量
			LocalTime:  cfg.IsLocalTime,    // 是否使用本地时间
			Compress:   cfg.EnableCompress, // 是否启用日志文件压缩
		}

		fileWriter = logger // 初始化文件写入器
	} else {
		fileWriter = io.Discard // 不允许将日志输出到文件, 则直接丢弃
	}

	// 创建 context 用于控制协程退出
	ctx, cancel := context.WithCancel(context.Background())

	// 创建一个新的FastLog实例, 将配置和缓冲区赋值给实例。
	f := &FastLog{
		logGer:        logger,                                  // 日志文件切割器
		fileWriter:    fileWriter,                              // 文件写入器, 用于将日志写入文件
		consoleWriter: consoleWriter,                           // 控制台写入器, 用于将日志写入控制台
		cl:            colorlib.NewColorLib(),                  // 颜色库实例, 用于在终端中显示颜色
		config:        cfg,                                     // 配置结构体
		logChan:       make(chan *logMessage, cfg.ChanIntSize), // 日志消息通道
		closeOnce:     sync.Once{},                             // 用于在结束时确保只执行一次
		cancel:        cancel,                                  // 用于取消上下文的函数
		ctx:           ctx,                                     // 上下文, 用于控制协程退出
	}

	// 根据noColor的值, 设置颜色库的颜色选项
	if f.config.NoColor {
		f.cl.NoColor.Store(true) // 设置颜色库的颜色选项为禁用
	}

	// 根据noBold的值, 设置颜色库的字体加粗选项
	if f.config.NoBold {
		f.cl.NoBold.Store(true) // 设置颜色库的字体加粗选项为禁用
	}

	// 使用 sync.Once 确保日志处理器只启动一次
	var startErr error
	f.startOnce.Do(func() {
		// 启动日志处理器和刷新器
		defer func() {
			if r := recover(); r != nil {
				startErr = fmt.Errorf("failed to start log processor: %v", r)
			}
		}()

		// 创建处理器
		processor := &processor{
			f:             f,                     // 日志记录器
			batchSize:     defaultBatchSize,      // 批量处理大小
			bufferSize:    initialBufferCapacity, // 缓冲区大小
			flushInterval: cfg.FlushInterval,     // 刷新间隔
			fileBuffer:    bytes.NewBuffer(nil),  // 文件缓冲区
			consoleBuffer: bytes.NewBuffer(nil),  // 控制台缓冲区
		}

		// 预分配缓冲区以减少内存分配
		processor.fileBuffer.Grow(processor.bufferSize)
		processor.consoleBuffer.Grow(processor.bufferSize)

		// 启动处理器
		f.logWait.Add(1)
		go processor.singleThreadProcessor()
	})

	// 检查启动是否成功
	if startErr != nil {
		return nil, startErr
	}

	// 返回FastLog实例和nil错误
	return f, nil
}

// Close 关闭 FastLog 实例
func (f *FastLog) Close() {
	// 确保只关闭一次
	f.closeOnce.Do(func() {
		// 打印关闭日志记录器的信息
		f.Info("stop logging...")

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
			time.Sleep(f.config.FlushInterval)
			close(done)
		}()

		// 等待完成信号，但设置超时避免无限等待
		select {
		case <-done:
			// 日志处理完成
		case <-time.After(3 * time.Second):
			// 超时保护，避免无限等待
		}

		// 关闭通道
		close(f.logChan)

		// 取消上下文
		f.cancel()

		// 等待日志处理器退出
		f.logWait.Wait()

		// 如果启用了文件输出，则关闭文件
		if f.config.OutputToFile {
			if f.logGer != nil {
				if closeErr := f.logGer.Close(); closeErr != nil {
					f.cl.PrintErrf("关闭日志文件失败: %v\nstack: %s\n", closeErr, debug.Stack())
				}
			}
		}
	})
}

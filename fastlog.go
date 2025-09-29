/*
fastlog.go - FastLog日志记录器核心实现
提供日志记录器的创建、初始化、日志写入及关闭等核心功能，
集成配置管理、缓冲区管理和日志处理流程。
*/
package fastlog

import (
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

// 为了提供更简洁的API调用方式，定义以下函数别名:
// 这样用户可以使用更短的函数名来创建日志实例和配置
var (
	// New 是 NewFastLog 的简写别名
	//
	// 用法:
	//  - logger := fastlog.New(config)
	//
	// 等价于:
	//  - logger := fastlog.NewFastLog(config)
	New = NewFastLog

	// NewCfg 是 NewFastLogConfig 的简写别名
	//
	// 用法:
	//  - config := fastlog.NewCfg()
	//
	// 等价于:
	//  - config := fastlog.NewFastLogConfig()
	NewCfg = NewFastLogConfig
)

// FastLog 日志记录器
type FastLog struct {
	fileWriter    io.Writer              // 文件写入器
	consoleWriter io.Writer              // 控制台写入器
	ctx           context.Context        // 用于控制全局日志记录器的生命周期
	cancel        context.CancelFunc     // 用于取消全局日志记录器的生命周期
	cl            *colorlib.ColorLib     // 提供终端颜色输出的库
	closeOnce     sync.Once              // 用于确保日志处理器只关闭一次
	startOnce     sync.Once              // 用于确保日志处理器只启动一次
	logChan       chan *logMsg           // 日志通道  用于异步写入日志文件
	logWait       sync.WaitGroup         // 等待组 用于等待所有goroutine完成
	logger        *logrotatex.LogRotateX // logrotatex 日志文件切割
	config        *FastLogConfig         // 嵌入的配置结构体
	bp            *bpThresholds          // 预计算背压阈值
	bufferSize    int                    // 缓冲区大小
}

// bpThresholds 预计算的背压阈值, 避免运行时频繁计算
type bpThresholds struct {
	threshold80 int // 80% 阈值
	threshold90 int // 90% 阈值
	threshold95 int // 95% 阈值
	threshold98 int // 98% 阈值
}

// NewFastLog 创建一个新的FastLog实例, 用于记录日志。
//
// 参数:
//   - config: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。
//
// 返回值:
//   - *FastLog: 一个指向FastLog实例的指针。
func NewFastLog(config *FastLogConfig) *FastLog {
	// 检查配置结构体是否为nil
	if config == nil {
		panic("FastLogConfig cannot be nil")
	}

	// 最终配置修正: 直接在原始配置上修正所有不合理的值
	config.validateConfig()

	// 克隆配置结构体防止原配置被意外修改
	cfg := &FastLogConfig{
		LogDirName:          config.LogDirName,          // 日志目录名称
		LogFileName:         config.LogFileName,         // 日志文件名称
		OutputToConsole:     config.OutputToConsole,     // 是否将日志输出到控制台
		OutputToFile:        config.OutputToFile,        // 是否将日志输出到文件
		LogLevel:            config.LogLevel,            // 日志级别
		ChanIntSize:         config.ChanIntSize,         // 通道大小
		FlushInterval:       config.FlushInterval,       // 刷新间隔
		LogFormat:           config.LogFormat,           // 日志格式
		MaxSize:             config.MaxSize,             // 最大日志文件大小, 单位为MB
		MaxAge:              config.MaxAge,              // 最大日志文件保留天数(单位为天)
		MaxFiles:            config.MaxFiles,            // 最大日志文件保留数量(默认为0, 表示不清理)
		LocalTime:           config.LocalTime,           // 是否使用本地时间
		Compress:            config.Compress,            // 是否启用日志文件压缩
		Color:               config.Color,               // 是否启用终端颜色
		Bold:                config.Bold,                // 是否启用终端字体加粗
		BatchSize:           config.BatchSize,           // 批处理数量
		DisableBackpressure: config.DisableBackpressure, // 是否禁用背压控制
	}

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

		// 初始化日志文件切割器
		logger = logrotatex.New(logFilePath) // 初始化日志文件切割器
		logger.MaxSize = cfg.MaxSize         // 最大日志文件大小, 单位为MB
		logger.MaxAge = cfg.MaxAge           // 最大日志文件保留天数
		logger.MaxFiles = cfg.MaxFiles       // 最大日志文件保留数量
		logger.Compress = cfg.Compress       // 是否启用日志文件压缩
		logger.LocalTime = cfg.LocalTime     // 是否使用本地时间

		fileWriter = logger // 初始化文件写入器
	} else {
		fileWriter = io.Discard // 不允许将日志输出到文件, 则直接丢弃
	}

	// 创建 context 用于控制协程退出
	ctx, cancel := context.WithCancel(context.Background())

	// 创建一个新的FastLog实例, 将配置和缓冲区赋值给实例。
	f := &FastLog{
		logger:        logger,                              // 日志文件切割器
		fileWriter:    fileWriter,                          // 文件写入器, 用于将日志写入文件
		consoleWriter: consoleWriter,                       // 控制台写入器, 用于将日志写入控制台
		cl:            colorlib.NewColorLib(),              // 颜色库实例, 用于在终端中显示颜色
		config:        cfg,                                 // 配置结构体
		logChan:       make(chan *logMsg, cfg.ChanIntSize), // 日志消息通道
		closeOnce:     sync.Once{},                         // 用于在结束时确保只执行一次
		cancel:        cancel,                              // 用于控制全局日志记录器的生命周期
		ctx:           ctx,                                 // 用于控制全局日志记录器的生命周期
	}

	// 根据Color的值, 设置颜色库的颜色选项
	if f.config.Color {
		f.cl.SetColor(true)
	} else {
		f.cl.SetColor(false)
	}

	// 根据Bold的值, 设置颜色库的字体加粗选项
	if f.config.Bold {
		f.cl.SetBold(true)
	} else {
		f.cl.SetBold(false)
	}

	// 预计算背压阈值
	logChanCap := cap(f.logChan) // 日志通道容量
	f.bp = &bpThresholds{
		threshold80: logChanCap * 80, // 80%阈值
		threshold90: logChanCap * 90, // 90%阈值
		threshold95: logChanCap * 95, // 95%阈值
		threshold98: logChanCap * 98, // 98%阈值
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

		// 创建处理器(使用依赖注入避免循环依赖)
		processor := newProcessor(
			f,                 // 传入FastLog作为依赖接口
			defaultBatchSize,  // 批量处理大小
			cfg.FlushInterval, // 刷新间隔
		)

		// 启动处理器(智能缓冲区池已在newProcessor中初始化)
		f.logWait.Add(1)
		go processor.singleThreadProcessor()
	})

	// 检查启动是否成功
	if startErr != nil {
		panic(startErr)
	}

	// 返回FastLog实例
	return f
}

// Close 安全关闭日志记录器
func (f *FastLog) Close() {
	// 使用 sync.Once 确保关闭操作只执行一次
	f.closeOnce.Do(func() {
		// 记录关闭日志
		f.Info("stop logging...")
		time.Sleep(10 * time.Millisecond)

		// 获取合适的超时时间
		closeTimeout := f.getCloseTimeout()

		// 创建关闭上下文
		closeCtx, closeCancel := context.WithTimeout(context.Background(), closeTimeout)
		defer closeCancel()

		// 优雅关闭: 先通知各组件关闭, 再等待处理完成
		f.gracefulShutdown(closeCtx)

		// 如果启用了文件写入器，则尝试关闭它。
		if f.config.OutputToFile && f.logger != nil {
			if err := f.logger.Close(); err != nil {
				f.cl.PrintErrorf("Failed to close log file: %v\n", err)
			}
		}
	})
}

/* ====== 不带占位符方法 ======*/

// Info 记录信息级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FastLog) Info(v ...any) {
	// 公共API入口参数验证
	if f == nil {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	// 调用logWithLevel方法记录日志
	f.logWithLevel(INFO, fmt.Sprint(v...), 3)
}

// Debug 记录调试级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FastLog) Debug(v ...any) {
	// 公共API入口参数验证
	if f == nil {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.logWithLevel(DEBUG, fmt.Sprint(v...), 3)
}

// Warn 记录警告级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FastLog) Warn(v ...any) {
	// 公共API入口参数验证
	if f == nil {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.logWithLevel(WARN, fmt.Sprint(v...), 3)
}

// Error 记录错误级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FastLog) Error(v ...any) {
	// 公共API入口参数验证
	if f == nil {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.logWithLevel(ERROR, fmt.Sprint(v...), 3)
}

// Fatal 记录致命级别的日志，不支持占位符，发送后关闭日志记录器
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FastLog) Fatal(v ...any) {
	// 公共API入口参数验证
	if f == nil {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.logFatal(fmt.Sprint(v...), 3)
}

/*====== 占位符方法 ======*/

// Infof 记录信息级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FastLog) Infof(format string, v ...any) {
	// 公共API入口参数验证
	if f == nil || format == "" {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.logWithLevel(INFO, fmt.Sprintf(format, v...), 3)
}

// Debugf 记录调试级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FastLog) Debugf(format string, v ...any) {
	// 公共API入口参数验证
	if f == nil || format == "" {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.logWithLevel(DEBUG, fmt.Sprintf(format, v...), 3)
}

// Warnf 记录警告级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FastLog) Warnf(format string, v ...any) {
	// 公共API入口参数验证
	if f == nil || format == "" {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.logWithLevel(WARN, fmt.Sprintf(format, v...), 3)
}

// Errorf 记录错误级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FastLog) Errorf(format string, v ...any) {
	// 公共API入口参数验证
	if f == nil || format == "" {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.logWithLevel(ERROR, fmt.Sprintf(format, v...), 3)
}

// Fatalf 记录致命级别的日志，支持占位符，发送后关闭日志记录器
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FastLog) Fatalf(format string, v ...any) {
	// 公共API入口参数验证
	if f == nil || format == "" {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.logFatal(fmt.Sprintf(format, v...), 3)
}

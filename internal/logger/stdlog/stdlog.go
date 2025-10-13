package stdlog

import (
	"path/filepath"
	"sync/atomic"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fastlog/internal/config"
	"gitee.com/MM-Q/logrotatex"
)

// StdLog 日志记录器
type StdLog struct {
	fileWriter *logrotatex.BufferedWriter // 带缓冲的文件写入器
	cl         *colorlib.ColorLib         // 提供终端颜色输出的库
	cfg        *config.FastLogConfig      // 嵌入的配置结构体
	closed     atomic.Bool                // 标记日志处理器是否已关闭
}

// NewStdLog 创建一个新的StdLog实例, 用于记录日志。
//
// 参数:
//   - config: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。
//
// 返回值:
//   - *StdLog: 一个指向StdLog实例的指针。
func NewStdLog(cfg *config.FastLogConfig) *StdLog {
	// 检查配置结构体是否为nil
	if cfg == nil {
		panic("FastLogConfig cannot be nil")
	}

	// 最终验证
	cfg.ValidateConfig()

	// 克隆配置结构体防止原配置被意外修改
	cfg = &config.FastLogConfig{
		LogDirName:      cfg.LogDirName,      // 日志目录名称
		LogFileName:     cfg.LogFileName,     // 日志文件名称
		OutputToConsole: cfg.OutputToConsole, // 是否将日志输出到控制台
		OutputToFile:    cfg.OutputToFile,    // 是否将日志输出到文件
		LogLevel:        cfg.LogLevel,        // 日志级别
		LogFormat:       cfg.LogFormat,       // 日志格式
		MaxSize:         cfg.MaxSize,         // 最大日志文件大小, 单位为MB
		MaxAge:          cfg.MaxAge,          // 最大日志文件保留天数(单位为天)
		MaxFiles:        cfg.MaxFiles,        // 最大日志文件保留数量(默认为0, 表示不清理)
		LocalTime:       cfg.LocalTime,       // 是否使用本地时间
		Compress:        cfg.Compress,        // 是否启用日志文件压缩
		Color:           cfg.Color,           // 是否启用终端颜色
		Bold:            cfg.Bold,            // 是否启用终端字体加粗
		FlushInterval:   cfg.FlushInterval,   // 刷新间隔, 单位为秒, 默认为0, 表示不做限制
		MaxBufferSize:   cfg.MaxBufferSize,   // 缓冲区最大容量, 单位为字节
		MaxWriteCount:   cfg.MaxWriteCount,   // 最大写入次数, 默认为0, 表示不做限制
		Async:           cfg.Async,           // 是否异步清理日志, 默认同步清理
	}

	// 初始化写入器
	var fileWriter *logrotatex.BufferedWriter

	// 如果允许将日志输出到文件, 则初始化带缓冲的文件写入器
	if cfg.OutputToFile {
		// 拼接日志文件路径
		logFilePath := filepath.Join(cfg.LogDirName, cfg.LogFileName)

		// 初始化日志文件切割器
		logger := logrotatex.NewLogRotateX(logFilePath) // 初始化日志文件切割器
		logger.MaxSize = cfg.MaxSize                    // 最大日志文件大小, 单位为MB
		logger.MaxAge = cfg.MaxAge                      // 最大日志文件保留天数
		logger.MaxFiles = cfg.MaxFiles                  // 最大日志文件保留数量
		logger.Compress = cfg.Compress                  // 是否启用日志文件压缩
		logger.LocalTime = cfg.LocalTime                // 是否使用本地时间
		logger.Async = cfg.Async                        // 是否异步清理日志

		// 初始化缓冲区配置
		bufCfg := logrotatex.DefBufCfg()
		bufCfg.FlushInterval = cfg.FlushInterval // 刷新间隔, 单位为秒, 默认为0, 表示不做限制
		bufCfg.MaxBufferSize = cfg.MaxBufferSize // 缓冲区最大容量, 单位为字节
		bufCfg.MaxWriteCount = cfg.MaxWriteCount // 最大写入次数, 默认为0, 表示不做限制

		// 创建带缓冲的批量写入器，嵌入日志切割器
		fileWriter = logrotatex.NewBufferedWriter(logger, bufCfg)
	}

	// 创建一个新的StdLog实例, 将配置和缓冲区赋值给实例。
	f := &StdLog{
		fileWriter: fileWriter,             // 带缓冲的文件写入器
		cl:         colorlib.NewColorLib(), // 颜色库实例, 用于在终端中显示颜色
		cfg:        cfg,                    // 配置结构体
		closed:     atomic.Bool{},          // 标记日志处理器是否已关闭
	}

	f.cl.SetColor(f.cfg.Color) // 设置颜色库的颜色选项
	f.cl.SetBold(f.cfg.Bold)   // 设置颜色库的字体加粗选项
	f.closed.Store(false)      // 初始化日志处理器为未关闭状态

	// 返回StdLog实例
	return f
}

// Close 关闭日志记录器
//
// 返回值:
//   - error: 如果关闭过程中发生错误, 返回错误信息; 否则返回nil
func (f *StdLog) Close() error {
	if f == nil {
		return nil
	}

	// 记录关闭日志
	f.Info("stop logging...")

	// 确保日志处理器只关闭一次 (原子操作)
	if !f.closed.CompareAndSwap(false, true) {
		return nil
	}

	// 如果启用了文件写入器，则尝试关闭它。
	if f.cfg.OutputToFile && f.fileWriter != nil {
		if err := f.fileWriter.Close(); err != nil {
			return err
		}
	}

	return nil
}

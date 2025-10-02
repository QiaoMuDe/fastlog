/*
fastlog.go - FastLog日志记录器核心实现
提供日志记录器的创建、初始化、日志写入及关闭等核心功能，
集成配置管理、缓冲区管理和日志处理流程。
*/
package fastlog

import (
	"fmt"
	"path/filepath"
	"sync/atomic"

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
	fileWriter *logrotatex.BufferedWriter // 带缓冲的文件写入器
	cl         *colorlib.ColorLib         // 提供终端颜色输出的库
	config     *FastLogConfig             // 嵌入的配置结构体
	closed     atomic.Bool                // 标记日志处理器是否已关闭
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
		LogDirName:      config.LogDirName,      // 日志目录名称
		LogFileName:     config.LogFileName,     // 日志文件名称
		OutputToConsole: config.OutputToConsole, // 是否将日志输出到控制台
		OutputToFile:    config.OutputToFile,    // 是否将日志输出到文件
		LogLevel:        config.LogLevel,        // 日志级别
		LogFormat:       config.LogFormat,       // 日志格式
		MaxSize:         config.MaxSize,         // 最大日志文件大小, 单位为MB
		MaxAge:          config.MaxAge,          // 最大日志文件保留天数(单位为天)
		MaxFiles:        config.MaxFiles,        // 最大日志文件保留数量(默认为0, 表示不清理)
		LocalTime:       config.LocalTime,       // 是否使用本地时间
		Compress:        config.Compress,        // 是否启用日志文件压缩
		Color:           config.Color,           // 是否启用终端颜色
		Bold:            config.Bold,            // 是否启用终端字体加粗
		FlushInterval:   config.FlushInterval,   // 刷新间隔, 单位为秒, 默认为0, 表示不做限制
		MaxBufferSize:   config.MaxBufferSize,   // 缓冲区最大容量, 单位为字节
		MaxWriteCount:   config.MaxWriteCount,   // 最大写入次数, 默认为0, 表示不做限制
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

		// 初始化缓冲区配置
		bufCfg := logrotatex.DefBufCfg()
		bufCfg.FlushInterval = cfg.FlushInterval // 刷新间隔, 单位为秒, 默认为0, 表示不做限制
		bufCfg.MaxBufferSize = cfg.MaxBufferSize // 缓冲区最大容量, 单位为字节
		bufCfg.MaxWriteCount = cfg.MaxWriteCount // 最大写入次数, 默认为0, 表示不做限制

		// 创建带缓冲的批量写入器，嵌入日志切割器
		fileWriter = logrotatex.NewBufferedWriter(logger, bufCfg)
	}

	// 创建一个新的FastLog实例, 将配置和缓冲区赋值给实例。
	f := &FastLog{
		fileWriter: fileWriter,             // 带缓冲的文件写入器
		cl:         colorlib.NewColorLib(), // 颜色库实例, 用于在终端中显示颜色
		config:     cfg,                    // 配置结构体
		closed:     atomic.Bool{},          // 标记日志处理器是否已关闭
	}

	f.cl.SetColor(f.config.Color) // 设置颜色库的颜色选项
	f.cl.SetBold(f.config.Bold)   // 设置颜色库的字体加粗选项
	f.closed.Store(false)         // 初始化日志处理器为未关闭状态

	// 返回FastLog实例
	return f
}

// Close 安全关闭日志记录器
func (f *FastLog) Close() {
	if f == nil {
		return
	}

	// 确保日志处理器只关闭一次
	if f.closed.Load() {
		return
	}

	// 记录关闭日志
	f.Info("stop logging...")

	// 如果启用了文件写入器，则尝试关闭它。
	if f.config.OutputToFile && f.fileWriter != nil {
		if err := f.fileWriter.Close(); err != nil {
			fmt.Println(err)
		}
	}
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

	// 检查是否已关闭日志记录器
	if f.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	// 调用processLog方法记录日志
	f.processLog(INFO, fmt.Sprint(v...))
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

	// 检查是否已关闭日志记录器
	if f.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.processLog(DEBUG, fmt.Sprint(v...))
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

	// 检查是否已关闭日志记录器
	if f.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.processLog(WARN, fmt.Sprint(v...))
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

	// 检查是否已关闭日志记录器
	if f.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.processLog(ERROR, fmt.Sprint(v...))
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

	// 检查是否已关闭日志记录器
	if f.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.logFatal(fmt.Sprint(v...))
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

	// 检查是否已关闭日志记录器
	if f.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.processLog(INFO, fmt.Sprintf(format, v...))
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

	// 检查是否已关闭日志记录器
	if f.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.processLog(DEBUG, fmt.Sprintf(format, v...))
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

	// 检查是否已关闭日志记录器
	if f.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.processLog(WARN, fmt.Sprintf(format, v...))
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

	// 检查是否已关闭日志记录器
	if f.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.processLog(ERROR, fmt.Sprintf(format, v...))
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

	// 检查是否已关闭日志记录器
	if f.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	f.logFatal(fmt.Sprintf(format, v...))
}

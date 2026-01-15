package flog

import (
	"fmt"
	"os"
	"sync/atomic"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fastlog/internal/config"
	"gitee.com/MM-Q/fastlog/internal/types"
	"gitee.com/MM-Q/logrotatex"
)

// FLog 是一个高性能的日志记录器, 支持键值对风格的使用和标准库fmt类似的使用,
// 同时提供了丰富的配置选项, 如日志级别、输出格式、日志轮转等。
type FLog struct {
	fileWriter *logrotatex.BufferedWriter // 带缓冲的文件写入器
	cl         *colorlib.ColorLib         // 提供终端颜色输出的库
	cfg        *config.FastLogConfig      // 嵌入的配置结构体
	closed     atomic.Bool                // 标记日志处理器是否已关闭
}

// NewFLog 创建一个新的FLog实例, 用于记录日志。
//
// 参数:
//   - cfg: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。
//
// 返回值:
//   - *FLog: 一个指向FLog实例的指针。
func NewFLog(cfg *config.FastLogConfig) *FLog {
	// 检查配置结构体是否为nil
	if cfg == nil {
		panic("fastlog: FastLogConfig cannot be nil")
	}

	// 最终验证
	cfg.ValidateConfig()

	// 克隆配置结构体防止原配置被意外修改
	cfg = cfg.Clone()

	// 初始化写入器
	fileWriter := config.CreateBufferedWriter(cfg)

	// 创建一个新的Fastlog实例, 将配置和缓冲区赋值给实例。
	f := &FLog{
		fileWriter: fileWriter,             // 带缓冲的文件写入器
		cl:         colorlib.NewColorLib(), // 颜色库实例, 用于在终端中显示颜色
		cfg:        cfg,                    // 配置结构体
		closed:     atomic.Bool{},          // 标记日志处理器是否已关闭
	}

	// 配置设置
	f.cl.SetColor(f.cfg.Color) // 设置颜色库的颜色选项
	f.cl.SetBold(f.cfg.Bold)   // 设置颜色库的字体加粗选项
	f.closed.Store(false)      // 初始化日志处理器为未关闭状态

	// 返回Fastlog实例
	return f
}

// Close 关闭日志处理器
//
// 返回值：
//   - error: 如果关闭过程中发生错误, 返回错误信息; 否则返回 nil。
func (f *FLog) Close() error {
	if f == nil || f.cfg == nil {
		return fmt.Errorf("fastlog: cannot close nil logger")
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

/* ====== 不带占位符方法 ======*/

// Info 记录信息级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FLog) Info(v ...any) {
	// 公共API入口参数验证
	if f == nil {
		return
	}

	// 调用processLog方法记录日志
	f.handleLog(types.INFO_Mask, fmt.Sprint(v...))
}

// Debug 记录调试级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FLog) Debug(v ...any) {
	// 公共API入口参数验证
	if f == nil {
		return
	}

	f.handleLog(types.DEBUG_Mask, fmt.Sprint(v...))
}

// Warn 记录警告级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FLog) Warn(v ...any) {
	// 公共API入口参数验证
	if f == nil {
		return
	}

	f.handleLog(types.WARN_Mask, fmt.Sprint(v...))
}

// Error 记录错误级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FLog) Error(v ...any) {
	// 公共API入口参数验证
	if f == nil {
		return
	}

	f.handleLog(types.ERROR_Mask, fmt.Sprint(v...))
}

// Fatal 记录致命级别的日志，不支持占位符，发送后关闭日志记录器
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FLog) Fatal(v ...any) {
	// 公共API入口参数验证
	if f == nil {
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
func (f *FLog) Infof(format string, v ...any) {
	// 公共API入口参数验证
	if f == nil || format == "" {
		return
	}

	f.handleLog(types.INFO_Mask, fmt.Sprintf(format, v...))
}

// Debugf 记录调试级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FLog) Debugf(format string, v ...any) {
	// 公共API入口参数验证
	if f == nil || format == "" {
		return
	}

	f.handleLog(types.DEBUG_Mask, fmt.Sprintf(format, v...))
}

// Warnf 记录警告级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FLog) Warnf(format string, v ...any) {
	// 公共API入口参数验证
	if f == nil || format == "" {
		return
	}

	f.handleLog(types.WARN_Mask, fmt.Sprintf(format, v...))
}

// Errorf 记录错误级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FLog) Errorf(format string, v ...any) {
	// 公共API入口参数验证
	if f == nil || format == "" {
		return
	}

	f.handleLog(types.ERROR_Mask, fmt.Sprintf(format, v...))
}

// Fatalf 记录致命级别的日志，支持占位符，发送后关闭日志记录器
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *FLog) Fatalf(format string, v ...any) {
	// 公共API入口参数验证
	if f == nil || format == "" {
		return
	}

	f.logFatal(fmt.Sprintf(format, v...))
}

// ====== 键值对方法 ======

// InfoFields 记录Info级别的键值对日志
//
// 参数：
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *FLog) InfoFields(msg string, fields ...*Field) {
	if f == nil {
		return
	}

	f.handleLog(types.INFO_Mask, msg, fields...)
}

// WarnFields 记录Warn级别的键值对日志
//
// 参数：
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *FLog) WarnFields(msg string, fields ...*Field) {
	if f == nil {
		return
	}

	f.handleLog(types.WARN_Mask, msg, fields...)
}

// ErrorFields 记录Error级别的键值对日志
//
// 参数：
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *FLog) ErrorFields(msg string, fields ...*Field) {
	if f == nil {
		return
	}

	f.handleLog(types.ERROR_Mask, msg, fields...)
}

// DebugFields 记录Debug级别的键值对日志
//
// 参数：
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *FLog) DebugFields(msg string, fields ...*Field) {
	if f == nil {
		return
	}

	f.handleLog(types.DEBUG_Mask, msg, fields...)
}

// FatalFields 记录Fatal级别的键值对日志并触发程序退出
//
// 参数：
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *FLog) FatalFields(msg string, fields ...*Field) {
	if f == nil {
		return
	}

	f.handleLog(types.FATAL_Mask, msg, fields...)

	// 关闭日志处理器
	if err := f.Close(); err != nil {
		fmt.Printf("fastlog: failed to close logger: %v\n", err)
	}

	// 退出程序
	os.Exit(1)
}

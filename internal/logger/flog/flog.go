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

// Flog 日志记录器
type Flog struct {
	fileWriter *logrotatex.BufferedWriter // 带缓冲的文件写入器
	cl         *colorlib.ColorLib         // 提供终端颜色输出的库
	cfg        *config.FastLogConfig      // 嵌入的配置结构体
	closed     atomic.Bool                // 标记日志处理器是否已关闭
}

// NewFlog 创建一个新的Flog实例, 用于记录日志。
//
// 参数:
//   - config: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。
//
// 返回值:
//   - *Flog: 一个指向Flog实例的指针。
func NewFlog(cfg *config.FastLogConfig) *Flog {
	// 检查配置结构体是否为nil
	if cfg == nil {
		panic("FastLogConfig cannot be nil")
	}

	// 最终验证
	cfg.ValidateConfig()

	// 克隆配置结构体防止原配置被意外修改
	cfg = cfg.Clone()

	// 初始化写入器
	fileWriter := config.CreateBufferedWriter(cfg)

	// 创建一个新的Flog实例, 将配置和缓冲区赋值给实例。
	f := &Flog{
		fileWriter: fileWriter,             // 带缓冲的文件写入器
		cl:         colorlib.NewColorLib(), // 颜色库实例, 用于在终端中显示颜色
		cfg:        cfg,                    // 配置结构体
		closed:     atomic.Bool{},          // 标记日志处理器是否已关闭
	}

	// 配置设置
	f.cl.SetColor(f.cfg.Color) // 设置颜色库的颜色选项
	f.cl.SetBold(f.cfg.Bold)   // 设置颜色库的字体加粗选项
	f.closed.Store(false)      // 初始化日志处理器为未关闭状态

	// 返回Flog实例
	return f
}

// Close 关闭日志处理器
//
// 返回值：
//   - error: 如果关闭过程中发生错误, 返回错误信息; 否则返回 nil。
func (f *Flog) Close() error {
	if f == nil || f.cfg == nil {
		return fmt.Errorf("fastlog: cannot close nil logger")
	}

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

// handleLog 处理日志记录
//
// 参数：
//   - level: 日志级别。
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *Flog) handleLog(level types.LogLevel, msg string, fields ...*Field) {
	if f == nil || f.cfg == nil {
		return
	}

	// 检查日志处理器是否已关闭
	if f.closed.Load() {
		return
	}

	// 检查日志级别，使用位运算判断是否应该记录该级别的日志
	if !types.ShouldLog(level, f.cfg.LogLevel) {
		return
	}

	// 创建日志条目
	e := NewEntry(f.cfg.CallerInfo, level, msg, fields...)
	defer putEntry(e) // 确保在函数返回前归还Entry实例到对象池

	// 构建日志条目
	log := buildLog(f.cfg, e)

	// 写入到终端
	if f.cfg.OutputToConsole {
		switch level {
		case types.INFO_Mask:
			f.cl.Blue(string(log))
		case types.WARN_Mask:
			f.cl.Yellow(string(log))
		case types.ERROR_Mask:
			f.cl.Red(string(log))
		case types.DEBUG_Mask:
			f.cl.Magenta(string(log))
		case types.FATAL_Mask:
			f.cl.Red(string(log))
		default:
			// 对于未知级别，使用默认颜色输出
			f.cl.White(string(log))
		}
	}

	// 写入到文件
	if f.cfg.OutputToFile && f.fileWriter != nil {
		// 确保日志以换行符结尾
		if len(log) == 0 || log[len(log)-1] != '\n' {
			log = append(log, '\n')
		}
		if _, err := f.fileWriter.Write(log); err != nil {
			fmt.Printf("fastlog: failed to write log: %v\n", err)
		}
	}
}

// Info 记录Info级别的日志
//
// 参数：
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *Flog) Info(msg string, fields ...*Field) {
	f.handleLog(types.INFO_Mask, msg, fields...)
}

// Warn 记录Warn级别的日志
//
// 参数：
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *Flog) Warn(msg string, fields ...*Field) {
	f.handleLog(types.WARN_Mask, msg, fields...)
}

// Error 记录Error级别的日志
//
// 参数：
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *Flog) Error(msg string, fields ...*Field) {
	f.handleLog(types.ERROR_Mask, msg, fields...)
}

// Debug 记录Debug级别的日志
//
// 参数：
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *Flog) Debug(msg string, fields ...*Field) {
	f.handleLog(types.DEBUG_Mask, msg, fields...)
}

// Fatal 记录Fatal级别的日志并触发程序退出
//
// 参数：
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *Flog) Fatal(msg string, fields ...*Field) {
	f.handleLog(types.FATAL_Mask, msg, fields...)

	// 关闭日志处理器
	if err := f.Close(); err != nil {
		fmt.Printf("fastlog: failed to close logger: %v\n", err)
	}

	// 退出程序
	os.Exit(1)
}

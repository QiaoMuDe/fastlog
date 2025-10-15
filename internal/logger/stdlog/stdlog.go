package stdlog

import (
	"fmt"
	"sync/atomic"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fastlog/internal/config"
	"gitee.com/MM-Q/fastlog/internal/types"
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
	cfg = cfg.Clone()

	// 初始化写入器
	fileWriter := config.CreateBufferedWriter(cfg)

	// 创建一个新的StdLog实例, 将配置和缓冲区赋值给实例。
	s := &StdLog{
		fileWriter: fileWriter,             // 带缓冲的文件写入器
		cl:         colorlib.NewColorLib(), // 颜色库实例, 用于在终端中显示颜色
		cfg:        cfg,                    // 配置结构体
		closed:     atomic.Bool{},          // 标记日志处理器是否已关闭
	}

	s.cl.SetColor(s.cfg.Color) // 设置颜色库的颜色选项
	s.cl.SetBold(s.cfg.Bold)   // 设置颜色库的字体加粗选项
	s.closed.Store(false)      // 初始化日志处理器为未关闭状态

	// 返回StdLog实例
	return s
}

// Close 关闭日志记录器
//
// 返回值:
//   - error: 如果关闭过程中发生错误, 返回错误信息; 否则返回nil
func (s *StdLog) Close() error {
	if s == nil {
		return nil
	}

	// 记录关闭日志
	s.Info("stop logging...")

	// 确保日志处理器只关闭一次 (原子操作)
	if !s.closed.CompareAndSwap(false, true) {
		return nil
	}

	// 如果启用了文件写入器，则尝试关闭它。
	if s.cfg.OutputToFile && s.fileWriter != nil {
		if err := s.fileWriter.Close(); err != nil {
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
func (s *StdLog) Info(v ...any) {
	// 公共API入口参数验证
	if s == nil {
		return
	}

	// 检查是否已关闭日志记录器
	if s.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	// 调用processLog方法记录日志
	s.processLog(types.INFO, fmt.Sprint(v...))
}

// Debug 记录调试级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (s *StdLog) Debug(v ...any) {
	// 公共API入口参数验证
	if s == nil {
		return
	}

	// 检查是否已关闭日志记录器
	if s.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	s.processLog(types.DEBUG, fmt.Sprint(v...))
}

// Warn 记录警告级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (s *StdLog) Warn(v ...any) {
	// 公共API入口参数验证
	if s == nil {
		return
	}

	// 检查是否已关闭日志记录器
	if s.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	s.processLog(types.WARN, fmt.Sprint(v...))
}

// Error 记录错误级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (s *StdLog) Error(v ...any) {
	// 公共API入口参数验证
	if s == nil {
		return
	}

	// 检查是否已关闭日志记录器
	if s.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	s.processLog(types.ERROR, fmt.Sprint(v...))
}

// Fatal 记录致命级别的日志，不支持占位符，发送后关闭日志记录器
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (s *StdLog) Fatal(v ...any) {
	// 公共API入口参数验证
	if s == nil {
		return
	}

	// 检查是否已关闭日志记录器
	if s.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	s.logFatal(fmt.Sprint(v...))
}

/*====== 占位符方法 ======*/

// Infof 记录信息级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (s *StdLog) Infof(format string, v ...any) {
	// 公共API入口参数验证
	if s == nil || format == "" {
		return
	}

	// 检查是否已关闭日志记录器
	if s.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	s.processLog(types.INFO, fmt.Sprintf(format, v...))
}

// Debugf 记录调试级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (s *StdLog) Debugf(format string, v ...any) {
	// 公共API入口参数验证
	if s == nil || format == "" {
		return
	}

	// 检查是否已关闭日志记录器
	if s.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	s.processLog(types.DEBUG, fmt.Sprintf(format, v...))
}

// Warnf 记录警告级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (s *StdLog) Warnf(format string, v ...any) {
	// 公共API入口参数验证
	if s == nil || format == "" {
		return
	}

	// 检查是否已关闭日志记录器
	if s.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	s.processLog(types.WARN, fmt.Sprintf(format, v...))
}

// Errorf 记录错误级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (s *StdLog) Errorf(format string, v ...any) {
	// 公共API入口参数验证
	if s == nil || format == "" {
		return
	}

	// 检查是否已关闭日志记录器
	if s.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	s.processLog(types.ERROR, fmt.Sprintf(format, v...))
}

// Fatalf 记录致命级别的日志，支持占位符，发送后关闭日志记录器
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (s *StdLog) Fatalf(format string, v ...any) {
	// 公共API入口参数验证
	if s == nil || format == "" {
		return
	}

	// 检查是否已关闭日志记录器
	if s.closed.Load() {
		return
	}

	// 检查参数是否为空
	if len(v) == 0 {
		return
	}

	s.logFatal(fmt.Sprintf(format, v...))
}

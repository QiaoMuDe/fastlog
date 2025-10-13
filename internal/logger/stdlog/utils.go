package stdlog

import (
	"fmt"

	"gitee.com/MM-Q/fastlog/internal/types"
)

/* ====== 不带占位符方法 ======*/

// Info 记录信息级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *StdLog) Info(v ...any) {
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
	f.processLog(types.INFO, fmt.Sprint(v...))
}

// Debug 记录调试级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *StdLog) Debug(v ...any) {
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

	f.processLog(types.DEBUG, fmt.Sprint(v...))
}

// Warn 记录警告级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *StdLog) Warn(v ...any) {
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

	f.processLog(types.WARN, fmt.Sprint(v...))
}

// Error 记录错误级别的日志，不支持占位符
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *StdLog) Error(v ...any) {
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

	f.processLog(types.ERROR, fmt.Sprint(v...))
}

// Fatal 记录致命级别的日志，不支持占位符，发送后关闭日志记录器
//
// 参数:
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *StdLog) Fatal(v ...any) {
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
func (f *StdLog) Infof(format string, v ...any) {
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

	f.processLog(types.INFO, fmt.Sprintf(format, v...))
}

// Debugf 记录调试级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *StdLog) Debugf(format string, v ...any) {
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

	f.processLog(types.DEBUG, fmt.Sprintf(format, v...))
}

// Warnf 记录警告级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *StdLog) Warnf(format string, v ...any) {
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

	f.processLog(types.WARN, fmt.Sprintf(format, v...))
}

// Errorf 记录错误级别的日志，支持占位符，格式化
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *StdLog) Errorf(format string, v ...any) {
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

	f.processLog(types.ERROR, fmt.Sprintf(format, v...))
}

// Fatalf 记录致命级别的日志，支持占位符，发送后关闭日志记录器
//
// 参数:
//   - format: 格式字符串
//   - v: 可变参数，可以是任意类型，会被转换为字符串
func (f *StdLog) Fatalf(format string, v ...any) {
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

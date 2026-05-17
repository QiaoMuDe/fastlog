package fastlog

import (
	"bytes"
	"errors"
	"io"
	"os"

	"gitee.com/MM-Q/color"
)

// ConsoleWriter 控制台写入器
type ConsoleWriter struct {
	w io.Writer
}

// Write 写入数据到控制台
//
// 参数:
//   - p: 要写入的字节数据
//
// 返回:
//   - int: 写入的字节数
//   - error: 写入过程中的错误
func (c *ConsoleWriter) Write(p []byte) (n int, err error) {
	return c.w.Write(p)
}

// Close 关闭控制台写入器
//
// 返回:
//   - error: 始终返回 nil
func (c *ConsoleWriter) Close() error {
	return nil
}

// ColorWriter 彩色控制台写入器
//
// 通过扫描字节流中的日志级别关键字自动着色输出。
// 将 NoColor 设为 true 可禁用颜色输出, 恢复原始文本。
type ColorWriter struct {
	w       io.Writer
	NoColor bool // 设为 true 禁用颜色, false 启用颜色
}

// NewColorWriter 创建彩色控制台写入器, 默认写入 os.Stdout
//
// 参数:
//   - noColor: 设为 true 禁用颜色输出
//
// 返回:
//   - *ColorWriter: 彩色写入器实例
func NewColorWriter(noColor bool) *ColorWriter {
	return &ColorWriter{w: os.Stdout, NoColor: noColor}
}

// Write 写入数据到控制台, 自动根据日志级别着色
//
// 参数:
//   - p: 要写入的字节数据
//
// 返回:
//   - int: 写入的字节数
//   - error: 写入过程中的错误
func (c *ColorWriter) Write(p []byte) (n int, err error) {
	if c.NoColor {
		return c.w.Write(p)
	}

	// 检测日志级别
	level := c.detectLevel(p)

	// 根据级别着色
	var colored string
	switch level {
	case DEBUG:
		colored = color.SCyan(string(p))

	case INFO:
		colored = color.SBlue(string(p))

	case WARN:
		colored = color.SYellow(string(p))

	case ERROR:
		colored = color.SRed(string(p))

	case FATAL:
		colored = color.New(color.FgRed, color.Bold).Sprint(string(p))

	case PANIC:
		colored = color.New(color.FgMagenta, color.Bold).Sprint(string(p))

	default:
		return c.w.Write(p)
	}

	return c.w.Write([]byte(colored))
}

// Close 关闭写入器
//
// 返回:
//   - error: 始终返回 nil
func (c *ColorWriter) Close() error {
	return nil
}

// detectLevel 从字节流中检测日志级别
//
// 参数:
//   - p: 字节流数据
//
// 返回:
//   - Level: 检测到的日志级别, 未识别时返回 INFO
func (c *ColorWriter) detectLevel(p []byte) Level {
	// 从高优先级到低优先级匹配, 避免误判
	for level := PANIC; level >= DEBUG; level-- {
		if bytes.Contains(p, []byte(level.String())) {
			return level
		}
	}
	return INFO
}

// MultiWriter 多路写入器, 同时将日志写入多个输出目标
type MultiWriter struct {
	writers []io.WriteCloser
}

// NewMultiWriter 创建多路写入器
//
// 参数:
//   - writers: 多个写入器
//
// 返回:
//   - *MultiWriter: 多路写入器实例
func NewMultiWriter(writers ...io.WriteCloser) *MultiWriter {
	return &MultiWriter{writers: writers}
}

// Write 写入数据到所有输出目标
//
// 参数:
//   - p: 要写入的字节数据
//
// 返回:
//   - int: 写入的字节数
//   - error: 写入过程中的错误
func (m *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range m.writers {
		_, err = w.Write(p)
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}

// Close 关闭所有输出目标
//
// 返回:
//   - error: 关闭过程中的错误, 多个错误会合并返回
func (m *MultiWriter) Close() error {
	var errs []error
	for _, w := range m.writers {
		if err := w.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Package colorlib 提供了扩展的日志级别输出功能。
// 该文件实现了带有级别标识的消息输出方法，包括调试、信息、成功、警告和错误等级别。
// 每个级别都有对应的颜色和前缀标识，便于在终端中快速识别不同类型的消息。
package colorlib

import (
	"fmt"

	"gitee.com/MM-Q/colorlib/internal/color"
)

const (
	levelDebug = iota // 调试级别
	levelInfo         // 信息级别
	levelOk           // 成功级别
	levelWarn         // 警告级别
	levelError        // 错误级别
)

var (
	// levelMap 存储了每个级别的名称与颜色信息
	levelMap = map[int]string{
		levelDebug: "debug: ", // 调试
		levelInfo:  "info: ",  // 信息
		levelOk:    "ok: ",    // 成功
		levelWarn:  "warn: ",  // 警告
		levelError: "error: ", // 错误
	}
)

// ==================================================================
// 级别消息提示方法
// ==================================================================

// PrintOk 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) PrintOk(msg ...any) {
	if len(msg) == 0 {
		return
	}

	c.promptMsg(levelOk, color.Green, true, "%s", fmt.Sprint(msg...))
}

// PrintError 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) PrintError(msg ...any) {
	if len(msg) == 0 {
		return
	}

	c.promptMsg(levelError, color.Red, true, "%s", fmt.Sprint(msg...))
}

// PrintWarn 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) PrintWarn(msg ...any) {
	if len(msg) == 0 {
		return
	}

	c.promptMsg(levelWarn, color.Yellow, true, "%s", fmt.Sprint(msg...))
}

// PrintInfo 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) PrintInfo(msg ...any) {
	if len(msg) == 0 {
		return
	}

	c.promptMsg(levelInfo, color.Blue, true, "%s", fmt.Sprint(msg...))
}

// PrintDebug 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) PrintDebug(msg ...any) {
	if len(msg) == 0 {
		return
	}

	c.promptMsg(levelDebug, color.Magenta, true, "%s", fmt.Sprint(msg...))
}

// PrintOkf 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) PrintOkf(format string, a ...any) {
	// 调用 promptMsg 方法，传入格式化后的字符串
	c.promptMsg(levelOk, color.Green, false, format, a...)
}

// PrintErrorf 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) PrintErrorf(format string, a ...any) {
	// 调用 promptMsg 方法，传入格式化后的字符串
	c.promptMsg(levelError, color.Red, false, format, a...)
}

// PrintWarnf 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) PrintWarnf(format string, a ...any) {
	// 调用 promptMsg 方法，传入格式化后的字符串
	c.promptMsg(levelWarn, color.Yellow, false, format, a...)
}

// PrintInfof 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) PrintInfof(format string, a ...any) {
	// 调用 promptMsg 方法，传入格式化后的字符串
	c.promptMsg(levelInfo, color.Blue, false, format, a...)
}

// PrintDebugf 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) PrintDebugf(format string, a ...any) {
	// 调用 promptMsg 方法，传入格式化后的字符串
	c.promptMsg(levelDebug, color.Magenta, false, format, a...)
}

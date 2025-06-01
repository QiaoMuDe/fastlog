package colorlib

import (
	"fmt"
)

// Bluef 方法用于将传入的参数以蓝色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Bluef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("blue", formattedMsg)
}

// Greenf 方法用于将传入的参数以绿色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Greenf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("green", formattedMsg)
}

// Redf 方法用于将传入的参数以红色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Redf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("red", formattedMsg)
}

// Yellowf 方法用于将传入的参数以黄色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Yellowf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("yellow", formattedMsg)
}

// Purplef 方法用于将传入的参数以紫色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Purplef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	c.printWithColor("purple", formattedMsg)
}

// Blackf 方法用于将传入的参数以黑色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Blackf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("black", formattedMsg)
}

// Cyanf 方法用于将传入的参数以青色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Cyanf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("cyan", formattedMsg)
}

// Whitef 方法用于将传入的参数以白色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Whitef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("white", formattedMsg)
}

// Grayf 方法用于将传入的参数以灰色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Grayf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("gray", formattedMsg)
}

// PrintSuccessf 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（带占位符）。
func (c *ColorLib) PrintSuccessf(format string, a ...any) {
	c.PromptMsg("success", "green", format, a...)
}

// PrintErrorf 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（带占位符）。
func (c *ColorLib) PrintErrorf(format string, a ...any) {
	c.PromptMsg("error", "red", format, a...)
}

// PrintWarningf 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（带占位符）。
func (c *ColorLib) PrintWarningf(format string, a ...any) {
	c.PromptMsg("warning", "yellow", format, a...)
}

// PrintInfof 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（带占位符）。
func (c *ColorLib) PrintInfof(format string, a ...any) {
	c.PromptMsg("info", "blue", format, a...)
}

// PrintDebugf 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（带占位符）。
func (c *ColorLib) PrintDebugf(format string, a ...any) {
	c.PromptMsg("debug", "purple", format, a...)
}

// Lredf 方法用于将传入的参数以亮红色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lredf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lred", formattedMsg)
}

// Lgreenf 方法用于将传入的参数以亮绿色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lgreenf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lgreen", formattedMsg)
}

// Lyellowf 方法用于将传入的参数以亮黄色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lyellowf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lyellow", formattedMsg)
}

// Lbluef 方法用于将传入的参数以亮蓝色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lbluef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lblue", formattedMsg)
}

// Lgreenf 方法用于将传入的参数以亮绿色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lpurplef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lpurple", formattedMsg)
}

// Lcyanf 方法用于将传入的参数以亮青色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lcyanf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lcyan", formattedMsg)
}

// Lwhitef 方法用于将传入的参数以亮白色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lwhitef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lwhite", formattedMsg)
}

// PrintOkf 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（带占位符）。
func (c *ColorLib) PrintOkf(format string, a ...any) {
	// 调用 PMsg 方法，传入格式化后的字符串
	c.PMsg("ok", "green", format, a...)
}

// PrintErrf 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（带占位符）。
func (c *ColorLib) PrintErrf(format string, a ...any) {
	// 调用 PMsg 方法，传入格式化后的字符串
	c.PMsg("err", "red", format, a...)
}

// PrintWarnf 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（带占位符）。
func (c *ColorLib) PrintWarnf(format string, a ...any) {
	// 调用 PMsg 方法，传入格式化后的字符串
	c.PMsg("warn", "yellow", format, a...)
}

// PrintInff 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（带占位符）。
func (c *ColorLib) PrintInff(format string, a ...any) {
	// 调用 PMsg 方法，传入格式化后的字符串
	c.PMsg("inf", "blue", format, a...)
}

// PrintDbgf 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（带占位符）。
func (c *ColorLib) PrintDbgf(format string, a ...any) {
	// 调用 PMsg 方法，传入格式化后的字符串
	c.PMsg("dbg", "purple", format, a...)
}

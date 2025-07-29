package colorlib

import (
	"fmt"
)

// PrintColorf 方法根据颜色代码常量打印对应颜色的文本
func (c *ColorLib) PrintColorf(code int, format string, a ...any) {
	// 根据颜色代码获取颜色名称
	var color string
	switch code {
	case Black:
		color = "black"
	case Red:
		color = "red"
	case Green:
		color = "green"
	case Yellow:
		color = "yellow"
	case Blue:
		color = "blue"
	case Purple:
		color = "purple"
	case Cyan:
		color = "cyan"
	case White:
		color = "white"
	case Gray:
		color = "gray"
	case Lred:
		color = "lred"
	case Lgreen:
		color = "lgreen"
	case Lyellow:
		color = "lyellow"
	case Lblue:
		color = "lblue"
	case Lpurple:
		color = "lpurple"
	case Lcyan:
		color = "lcyan"
	case Lwhite:
		color = "lwhite"
	default:
		fmt.Println("Invalid color code:", code)
		return
	}

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor(color, fmt.Sprintf(format, a...))
}

// Bluef 方法用于将传入的参数以蓝色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Bluef(format string, a ...any) {
	c.PrintColorf(Blue, format, a...)
}

// Greenf 方法用于将传入的参数以绿色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Greenf(format string, a ...any) {
	c.PrintColorf(Green, format, a...)
}

// Redf 方法用于将传入的参数以红色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Redf(format string, a ...any) {
	c.PrintColorf(Red, format, a...)
}

// Yellowf 方法用于将传入的参数以黄色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Yellowf(format string, a ...any) {
	c.PrintColorf(Yellow, format, a...)
}

// Purplef 方法用于将传入的参数以紫色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Purplef(format string, a ...any) {
	c.PrintColorf(Purple, format, a...)
}

// Blackf 方法用于将传入的参数以黑色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Blackf(format string, a ...any) {
	c.PrintColorf(Black, format, a...)
}

// Cyanf 方法用于将传入的参数以青色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Cyanf(format string, a ...any) {
	c.PrintColorf(Cyan, format, a...)
}

// Whitef 方法用于将传入的参数以白色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Whitef(format string, a ...any) {
	c.PrintColorf(White, format, a...)
}

// Grayf 方法用于将传入的参数以灰色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Grayf(format string, a ...any) {
	c.PrintColorf(Gray, format, a...)
}

// PrintSuccessf 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（带占位符）。
func (c *ColorLib) PrintSuccessf(format string, a ...any) {
	c.promptMsg("success", "green", format, a...)
}

// PrintErrorf 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（带占位符）。
func (c *ColorLib) PrintErrorf(format string, a ...any) {
	c.promptMsg("error", "red", format, a...)
}

// PrintWarningf 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（带占位符）。
func (c *ColorLib) PrintWarningf(format string, a ...any) {
	c.promptMsg("warning", "yellow", format, a...)
}

// PrintInfof 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（带占位符）。
func (c *ColorLib) PrintInfof(format string, a ...any) {
	c.promptMsg("info", "blue", format, a...)
}

// PrintDebugf 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（带占位符）。
func (c *ColorLib) PrintDebugf(format string, a ...any) {
	c.promptMsg("debug", "purple", format, a...)
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

// Lgreenf 方法用于将传入的参数以亮紫色文本形式打印到控制台（带占位符）。
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
	// 调用 promptMsg 方法，传入格式化后的字符串
	c.promptMsg("ok", "green", format, a...)
}

// PrintErrf 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（带占位符）。
func (c *ColorLib) PrintErrf(format string, a ...any) {
	// 调用 promptMsg 方法，传入格式化后的字符串
	c.promptMsg("err", "red", format, a...)
}

// PrintWarnf 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（带占位符）。
func (c *ColorLib) PrintWarnf(format string, a ...any) {
	// 调用 promptMsg 方法，传入格式化后的字符串
	c.promptMsg("warn", "yellow", format, a...)
}

// PrintInff 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（带占位符）。
func (c *ColorLib) PrintInff(format string, a ...any) {
	// 调用 promptMsg 方法，传入格式化后的字符串
	c.promptMsg("inf", "blue", format, a...)
}

// PrintDbgf 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（带占位符）。
func (c *ColorLib) PrintDbgf(format string, a ...any) {
	// 调用 promptMsg 方法，传入格式化后的字符串
	c.promptMsg("dbg", "purple", format, a...)
}

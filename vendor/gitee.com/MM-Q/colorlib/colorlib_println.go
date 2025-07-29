package colorlib

import "fmt"

// PrintColorln 方法根据颜色代码常量打印对应颜色的文本
func (c *ColorLib) PrintColorln(code int, msg ...any) {
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

	// 直接拼接参数并添加换行符
	c.printWithColor(color, fmt.Sprint(msg...)+"\n")
}

// Blue 方法用于将传入的参数以蓝色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Blue(msg ...any) {
	c.PrintColorln(Blue, msg...)
}

// Green 方法用于将传入的参数以绿色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Green(msg ...any) {
	c.PrintColorln(Green, msg...)
}

// Red 方法用于将传入的参数以红色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Red(msg ...any) {
	c.PrintColorln(Red, msg...)
}

// Yellow 方法用于将传入的参数以黄色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Yellow(msg ...any) {
	c.PrintColorln(Yellow, msg...)
}

// Purple 方法用于将传入的参数以紫色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Purple(msg ...any) {
	c.PrintColorln(Purple, msg...)
}

// Black 方法用于将传入的参数以黑色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Black(msg ...any) {
	c.PrintColorln(Black, msg...)
}

// Cyan 方法用于将传入的参数以青色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Cyan(msg ...any) {
	c.PrintColorln(Cyan, msg...)
}

// White 方法用于将传入的参数以白色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) White(msg ...any) {
	c.PrintColorln(White, msg...)
}

// Gray 方法用于将传入的参数以灰色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Gray(msg ...any) {
	c.PrintColorln(Gray, msg...)
}

// PrintSuccess 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（不带占位符）。
func (c *ColorLib) PrintSuccess(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("success", "green", "%s", "\n")
		return
	}

	// 直接拼接参数并添加换行符
	c.promptMsg("success", "green", "%s", fmt.Sprint(msg...)+"\n")
}

// PrintError 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（不带占位符）。
func (c *ColorLib) PrintError(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("error", "red", "%s", "\n")
		return
	}

	// 直接拼接参数并添加换行符
	c.promptMsg("error", "red", "%s", fmt.Sprint(msg...)+"\n")
}

// PrintWarning 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（不带占位符）。
func (c *ColorLib) PrintWarning(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("warning", "yellow", "%s", "\n")
		return
	}

	// 直接拼接参数并添加换行符
	c.promptMsg("warning", "yellow", "%s", fmt.Sprint(msg...)+"\n")
}

// PrintInfo 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（不带占位符）。
func (c *ColorLib) PrintInfo(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("info", "blue", "%s", "\n")
		return
	}

	// 直接拼接参数并添加换行符
	c.promptMsg("info", "blue", "%s", fmt.Sprint(msg...)+"\n")
}

// PrintDebug 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（不带占位符）。
func (c *ColorLib) PrintDebug(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("debug", "purple", "%s", "\n")
		return
	}

	// 直接拼接参数并添加换行符
	c.promptMsg("debug", "purple", "%s", fmt.Sprint(msg...)+"\n")
}

// Lred 方法用于将传入的参数以亮红色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lred(msg ...any) {
	// 直接拼接参数并添加换行符
	c.printWithColor("Lred", fmt.Sprint(msg...)+"\n")
}

// Lgreen 方法用于将传入的参数以亮绿色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lgreen(msg ...any) {
	// 直接拼接参数并添加换行符
	c.printWithColor("lgreen", fmt.Sprint(msg...)+"\n")
}

// Lyellow 方法用于将传入的参数以亮黄色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lyellow(msg ...any) {
	// 直接拼接参数并添加换行符
	c.printWithColor("lyellow", fmt.Sprint(msg...)+"\n")
}

// Lblue 方法用于将传入的参数以亮蓝色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lblue(msg ...any) {
	// 直接拼接参数并添加换行符
	c.printWithColor("lblue", fmt.Sprint(msg...)+"\n")
}

// Lgreen 方法用于将传入的参数以亮紫色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lpurple(msg ...any) {
	// 直接拼接参数并添加换行符
	c.printWithColor("lpurple", fmt.Sprint(msg...)+"\n")
}

// Lcyan 方法用于将传入的参数以亮青色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lcyan(msg ...any) {
	// 直接拼接参数并添加换行符
	c.printWithColor("lcyan", fmt.Sprint(msg...)+"\n")
}

// Lwhite 方法用于将传入的参数以亮白色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lwhite(msg ...any) {
	// 直接拼接参数并添加换行符
	c.printWithColor("lwhite", fmt.Sprint(msg...)+"\n")
}

// PrintOk 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（不带占位符）。
func (c *ColorLib) PrintOk(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("ok", "green", "%s", "\n")
		return
	}

	c.promptMsg("ok", "green", "%s", fmt.Sprint(msg...)+"\n")
}

// PrintErr 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（不带占位符）。
func (c *ColorLib) PrintErr(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("err", "red", "%s", "\n")
		return
	}

	c.promptMsg("err", "red", "%s", fmt.Sprint(msg...)+"\n")
}

// PrintWarn 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（不带占位符）。
func (c *ColorLib) PrintWarn(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("warn", "yellow", "%s", "\n")
		return
	}

	c.promptMsg("warn", "yellow", "%s", fmt.Sprint(msg...)+"\n")
}

// PrintInf 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（不带占位符）。
func (c *ColorLib) PrintInf(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("inf", "blue", "%s", "\n")
		return
	}

	c.promptMsg("inf", "blue", "%s", fmt.Sprint(msg...)+"\n")
}

// PrintDbg 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（不带占位符）。
func (c *ColorLib) PrintDbg(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("dbg", "purple", "%s", "\n")
		return
	}

	c.promptMsg("dbg", "purple", "%s", fmt.Sprint(msg...)+"\n")
}

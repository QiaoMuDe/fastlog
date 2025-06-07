package colorlib

import "fmt"

// Blue 方法用于将传入的参数以蓝色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Blue(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("blue", combinedMsg)
}

// Green 方法用于将传入的参数以绿色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Green(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("green", combinedMsg)
}

// Red 方法用于将传入的参数以红色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Red(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("red", combinedMsg)
}

// Yellow 方法用于将传入的参数以黄色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Yellow(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("yellow", combinedMsg)
}

// Purple 方法用于将传入的参数以紫色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Purple(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}
	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("purple", combinedMsg)
}

// Black 方法用于将传入的参数以黑色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Black(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("black", combinedMsg)
}

// Cyan 方法用于将传入的参数以青色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Cyan(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("cyan", combinedMsg)
}

// White 方法用于将传入的参数以白色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) White(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("white", combinedMsg)
}

// Gray 方法用于将传入的参数以灰色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Gray(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("gray", combinedMsg)
}

// PrintSuccess 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（不带占位符）。
func (c *ColorLib) PrintSuccess(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("success", "green", "%s", "\n")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.promptMsg("success", "green", "%s", combinedMsg)
}

// PrintError 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（不带占位符）。
func (c *ColorLib) PrintError(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("error", "red", "%s", "\n")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.promptMsg("error", "red", "%s", combinedMsg)
}

// PrintWarning 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（不带占位符）。
func (c *ColorLib) PrintWarning(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("warning", "yellow", "%s", "\n")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.promptMsg("warning", "yellow", "%s", combinedMsg)
}

// PrintInfo 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（不带占位符）。
func (c *ColorLib) PrintInfo(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("info", "blue", "%s", "\n")
		return
	}

	// 使用 fmt.Sprint 将 msg 中的所有元素拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.promptMsg("info", "blue", "%s", combinedMsg)
}

// PrintDebug 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（不带占位符）。
func (c *ColorLib) PrintDebug(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("debug", "purple", "%s", "\n")
		return
	}

	// 使用 fmt.Sprint 将 msg 中的所有元素拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.promptMsg("debug", "purple", "%s", combinedMsg)
}

// Lred 方法用于将传入的参数以亮红色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lred(msg ...any) {
	// 检查传入的参数数量，如果为0，则直接返回
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	// 调用 ColorLib 类型的 printWithColor 方法，传入颜色 "lred" 和拼接后的字符串
	c.printWithColor("lred", combinedMsg)
}

// Lgreen 方法用于将传入的参数以亮绿色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lgreen(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("lgreen", combinedMsg)
}

// Lyellow 方法用于将传入的参数以亮黄色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lyellow(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("lyellow", combinedMsg)
}

// Lblue 方法用于将传入的参数以亮蓝色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lblue(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("lblue", combinedMsg)
}

// Lgreen 方法用于将传入的参数以亮紫色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lpurple(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("lpurple", combinedMsg)
}

// Lcyan 方法用于将传入的参数以亮青色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lcyan(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("lcyan", combinedMsg)
}

// Lwhite 方法用于将传入的参数以亮白色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lwhite(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.printWithColor("lwhite", combinedMsg)
}

// PrintOk 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（不带占位符）。
func (c *ColorLib) PrintOk(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("ok", "green", "%s", "\n")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.promptMsg("ok", "green", "%s", combinedMsg)
}

// PrintErr 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（不带占位符）。
func (c *ColorLib) PrintErr(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("err", "red", "%s", "\n")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.promptMsg("err", "red", "%s", combinedMsg)
}

// PrintWarn 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（不带占位符）。
func (c *ColorLib) PrintWarn(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("warn", "yellow", "%s", "\n")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.promptMsg("warn", "yellow", "%s", combinedMsg)
}

// PrintInf 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（不带占位符）。
func (c *ColorLib) PrintInf(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("inf", "blue", "%s", "\n")
		return
	}

	// 使用 fmt.Sprint 将 msg 中的所有元素拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.promptMsg("inf", "blue", "%s", combinedMsg)
}

// PrintDbg 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（不带占位符）。
func (c *ColorLib) PrintDbg(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.promptMsg("dbg", "purple", "%s", "\n")
		return
	}

	// 使用 fmt.Sprint 将 msg 中的所有元素拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	combinedMsg += "\n"
	c.promptMsg("dbg", "purple", "%s", combinedMsg)
}

package colorlib

import "fmt"

// SColorf 方法根据颜色代码常量打印对应颜色的文本
func (c *ColorLib) SColorf(code int, format string, a ...any) string {
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
		return fmt.Sprintf("Invalid color code: %d", code)
	}

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor(color, fmt.Sprintf(format, a...))
}

// Sbluef 方法用于将传入的参数以蓝色文本形式返回（带占位符）。
func (c *ColorLib) Sbluef(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("blue", fmt.Sprintf(format, a...))
}

// Sgreenf 方法用于将传入的参数以绿色文本形式返回（带占位符）。
func (c *ColorLib) Sgreenf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("green", fmt.Sprintf(format, a...))
}

// Sredf 方法用于将传入的参数以红色文本形式返回（带占位符）。
func (c *ColorLib) Sredf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("red", fmt.Sprintf(format, a...))
}

// Syellowf 方法用于将传入的参数以黄色文本形式返回（带占位符）。
func (c *ColorLib) Syellowf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("yellow", fmt.Sprintf(format, a...))
}

// Spurplef 方法用于将传入的参数以紫色文本形式返回（带占位符）。
func (c *ColorLib) Spurplef(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("purple", fmt.Sprintf(format, a...))
}

// Sblackf 方法用于将传入的参数以黑色文本形式返回（带占位符）。
func (c *ColorLib) Sblackf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("black", fmt.Sprintf(format, a...))
}

// Scyanf 方法用于将传入的参数以青色文本形式返回（带占位符）。
func (c *ColorLib) Scyanf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("cyan", fmt.Sprintf(format, a...))
}

// Swhitef 方法用于将传入的参数以白色文本形式返回（带占位符）。
func (c *ColorLib) Swhitef(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("white", fmt.Sprintf(format, a...))
}

// Sgrayf 方法用于将传入的参数以灰色文本形式返回（带占位符）。
func (c *ColorLib) Sgrayf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("gray", fmt.Sprintf(format, a...))
}

// Slredf 方法用于将传入的参数以亮红色文本形式返回（带占位符）。
func (c *ColorLib) Slredf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lred", fmt.Sprintf(format, a...))
}

// Slgreenf 方法用于将传入的参数以亮绿色文本形式返回（带占位符）。
func (c *ColorLib) Slgreenf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lgreen", fmt.Sprintf(format, a...))
}

// Slyellowf 方法用于将传入的参数以亮黄色文本形式返回（带占位符）。
func (c *ColorLib) Slyellowf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lyellow", fmt.Sprintf(format, a...))
}

// Slbluef 方法用于将传入的参数以亮蓝色文本形式返回（带占位符）。
func (c *ColorLib) Slbluef(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lblue", fmt.Sprintf(format, a...))
}

// Slgreenf 方法用于将传入的参数以亮绿色文本形式返回（带占位符）。
func (c *ColorLib) Slpurplef(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lpurple", fmt.Sprintf(format, a...))
}

// Slcyanf 方法用于将传入的参数以亮青色文本形式返回（带占位符）。
func (c *ColorLib) Slcyanf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lcyan", fmt.Sprintf(format, a...))
}

// Slwhitef 方法用于将传入的参数以亮白色文本形式返回（带占位符）。
func (c *ColorLib) Slwhitef(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lwhite", fmt.Sprintf(format, a...))
}

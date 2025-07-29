package colorlib

import "fmt"

// SColor 方法根据颜色代码常量打印对应颜色的文本
func (c *ColorLib) SColor(code int, msg ...any) string {
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

	// 调用 returnWithColor 方法，传入拼接后的字符串
	return c.returnWithColor(color, fmt.Sprint(msg...))
}

// Sblue 方法用于将传入的参数以蓝色文本形式返回（不带占位符）。
func (c *ColorLib) Sblue(msg ...any) string {
	return c.returnWithColor("blue", fmt.Sprint(msg...))
}

// Sgreen 方法用于将传入的参数以绿色文本形式返回（不带占位符）。
func (c *ColorLib) Sgreen(msg ...any) string {
	return c.returnWithColor("green", fmt.Sprint(msg...))
}

// Sred 方法用于将传入的参数以红色文本形式返回（不带占位符）。
func (c *ColorLib) Sred(msg ...any) string {
	return c.returnWithColor("red", fmt.Sprint(msg...))
}

// Syellow 方法用于将传入的参数以黄色文本形式返回（不带占位符）。
func (c *ColorLib) Syellow(msg ...any) string {
	return c.returnWithColor("yellow", fmt.Sprint(msg...))
}

// Spurple 方法用于将传入的参数以紫色文本形式返回（不带占位符）。
func (c *ColorLib) Spurple(msg ...any) string {
	return c.returnWithColor("purple", fmt.Sprint(msg...))
}

// Sblack 方法用于将传入的参数以黑色文本形式返回（不带占位符）。
func (c *ColorLib) Sblack(msg ...any) string {
	return c.returnWithColor("black", fmt.Sprint(msg...))
}

// Scyan 方法用于将传入的参数以青色文本形式返回（不带占位符）。
func (c *ColorLib) Scyan(msg ...any) string {
	return c.returnWithColor("cyan", fmt.Sprint(msg...))
}

// Swhite 方法用于将传入的参数以白色文本形式返回（不带占位符）。
func (c *ColorLib) Swhite(msg ...any) string {
	return c.returnWithColor("white", fmt.Sprint(msg...))
}

// Sgray 方法用于将传入的参数以灰色文本形式返回（不带占位符）。
func (c *ColorLib) Sgray(msg ...any) string {
	return c.returnWithColor("gray", fmt.Sprint(msg...))
}

// Slred 方法用于将传入的参数以亮红色文本形式返回（不带占位符）。
func (c *ColorLib) Slred(msg ...any) string {
	return c.returnWithColor("lred", fmt.Sprint(msg...))
}

// Slgreen 方法用于将传入的参数以亮绿色文本形式返回（不带占位符）。
func (c *ColorLib) Slgreen(msg ...any) string {
	return c.returnWithColor("lgreen", fmt.Sprint(msg...))
}

// Slyellow 方法用于将传入的参数以亮黄色文本形式返回（不带占位符）。
func (c *ColorLib) Slyellow(msg ...any) string {
	return c.returnWithColor("lyellow", fmt.Sprint(msg...))
}

// Slblue 方法用于将传入的参数以亮蓝色文本形式返回（不带占位符）。
func (c *ColorLib) Slblue(msg ...any) string {
	return c.returnWithColor("lblue", fmt.Sprint(msg...))
}

// Slgreen 方法用于将传入的参数以亮绿色文本形式返回（不带占位符）。
func (c *ColorLib) Slpurple(msg ...any) string {
	return c.returnWithColor("lpurple", fmt.Sprint(msg...))
}

// Slcyan 方法用于将传入的参数以亮青色文本形式返回（不带占位符）。
func (c *ColorLib) Slcyan(msg ...any) string {
	return c.returnWithColor("lcyan", fmt.Sprint(msg...))
}

// Slwhite 方法用于将传入的参数以亮白色文本形式返回（不带占位符）。
func (c *ColorLib) Slwhite(msg ...any) string {
	return c.returnWithColor("lwhite", fmt.Sprint(msg...))
}

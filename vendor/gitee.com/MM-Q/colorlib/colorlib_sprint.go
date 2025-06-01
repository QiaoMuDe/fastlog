package colorlib

import "fmt"

// Sblue 方法用于将传入的参数以蓝色文本形式返回（不带占位符）。
func (c *ColorLib) Sblue(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("blue", combinedMsg)
}

// Sgreen 方法用于将传入的参数以绿色文本形式返回（不带占位符）。
func (c *ColorLib) Sgreen(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("green", combinedMsg)
}

// Sred 方法用于将传入的参数以红色文本形式返回（不带占位符）。
func (c *ColorLib) Sred(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("red", combinedMsg)
}

// Syellow 方法用于将传入的参数以黄色文本形式返回（不带占位符）。
func (c *ColorLib) Syellow(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("yellow", combinedMsg)
}

// Spurple 方法用于将传入的参数以紫色文本形式返回（不带占位符）。
func (c *ColorLib) Spurple(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("purple", combinedMsg)
}

// Sblack 方法用于将传入的参数以黑色文本形式返回（不带占位符）。
func (c *ColorLib) Sblack(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("black", combinedMsg)
}

// Scyan 方法用于将传入的参数以青色文本形式返回（不带占位符）。
func (c *ColorLib) Scyan(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("cyan", combinedMsg)
}

// Swhite 方法用于将传入的参数以白色文本形式返回（不带占位符）。
func (c *ColorLib) Swhite(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("white", combinedMsg)
}

// Sgray 方法用于将传入的参数以灰色文本形式返回（不带占位符）。
func (c *ColorLib) Sgray(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("gray", combinedMsg)
}

// Slred 方法用于将传入的参数以亮红色文本形式返回（不带占位符）。
func (c *ColorLib) Slred(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lred", combinedMsg)
}

// Slgreen 方法用于将传入的参数以亮绿色文本形式返回（不带占位符）。
func (c *ColorLib) Slgreen(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}
	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lgreen", combinedMsg)
}

// Slblue 方法用于将传入的参数以亮蓝色文本形式返回（不带占位符）。
func (c *ColorLib) Slyellow(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lyellow", combinedMsg)
}

// Slblue 方法用于将传入的参数以亮蓝色文本形式返回（不带占位符）。
func (c *ColorLib) Slblue(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lblue", combinedMsg)
}

// Slgreen 方法用于将传入的参数以亮绿色文本形式返回（不带占位符）。
func (c *ColorLib) Slpurple(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lpurple", combinedMsg)
}

// Slcyan 方法用于将传入的参数以亮青色文本形式返回（不带占位符）。
func (c *ColorLib) Slcyan(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lcyan", combinedMsg)
}

// Slwhite 方法用于将传入的参数以亮白色文本形式返回（不带占位符）。
func (c *ColorLib) Slwhite(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lwhite", combinedMsg)
}

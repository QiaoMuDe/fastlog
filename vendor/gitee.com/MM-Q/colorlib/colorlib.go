package colorlib

import (
	"fmt"
	"strings"
)

// ColorLib 结构体用于管理颜色输出和日志级别映射。
type ColorLib struct {
	levelMap     map[string]string // LevelMap 是一个映射，用于将日志级别映射到对应的前缀,// 日志级别映射到对应的前缀, 后面留空, 方便后面拼接提示内容
	colorMap     map[string]int    // colorMap 是一个映射，用于将颜色名称映射到对应的 ANSI 颜色代码。
	NoColor      bool              // NoColor 控制是否禁用颜色输出
	formatBuffer strings.Builder   // formatBuffer 用于构建格式化后的字符串。
	NoBold       bool              // NoBold 控制是否禁用字体加粗
}

const (
	reset   = 0  // Reset 重置所有属性
	Black   = 30 // Black 黑色
	Red     = 31 // Red 红色
	Green   = 32 // Green 绿色
	Yellow  = 33 // Yellow 黄色
	Blue    = 34 // Blue 蓝色
	Purple  = 35 // Purple 紫色
	Cyan    = 36 // Cyan 青色
	White   = 37 // White 白色
	Gray    = 90 // Gray 灰色
	Lred    = 91 // Lred 亮红色
	Lgreen  = 92 // Lgreen 亮绿色
	Lyellow = 93 // Lyellow 亮黄色
	Lblue   = 94 // Lblue 亮蓝色
	Lpurple = 95 // Lpurple 亮紫色
	Lcyan   = 96 // Lcyan 亮青色
	Lwhite  = 97 // Lwhite 亮白色
)

// 定义日志级别和颜色的映射
var (
	colorMap = map[string]int{
		"black":   Black,
		"red":     Red,
		"green":   Green,
		"yellow":  Yellow,
		"blue":    Blue,
		"purple":  Purple,
		"cyan":    Cyan,
		"white":   White,
		"gray":    Gray,
		"lred":    Lred,
		"lgreen":  Lgreen,
		"lyellow": Lyellow,
		"lblue":   Lblue,
		"lpurple": Lpurple,
		"lcyan":   Lcyan,
		"lwhite":  Lwhite,
	}

	levelMap = map[string]string{
		"success": "[Success] ", // 成功信息级别的前缀
		"error":   "[Error] ",   // 错误信息级别的前缀
		"warning": "[Warning] ", // 警告信息级别的前缀
		"info":    "[Info] ",    // 信息信息级别的前缀
		"debug":   "[Debug] ",   // 调试信息级别的前缀
		"ok":      "ok: ",       // 简写形式
		"err":     "err: ",      // 简写形式
		"inf":     "info: ",     // 简写形式
		"dbg":     "debug: ",    // 简写形式
		"warn":    "warn: ",     // 简写形式
	}
)

// ColorLibInterface 是一个接口，定义了一组方法，用于打印和返回带有颜色的文本。
type ColorLibInterface interface {
	// 需要占位符的方法(自带换行符)
	Bluef(format string, a ...any)           // 打印蓝色信息到控制台（带占位符）
	Greenf(format string, a ...any)          // 打印绿色信息到控制台（带占位符）
	Redf(format string, a ...any)            // 打印红色信息到控制台（带占位符）
	Yellowf(format string, a ...any)         // 打印黄色信息到控制台（带占位符）
	Purplef(format string, a ...any)         // 打印紫色信息到控制台（带占位符）
	Sbluef(format string, a ...any) string   // 返回构造后的蓝色字符串（带占位符）
	Sgreenf(format string, a ...any) string  // 返回构造后的绿色字符串（带占位符）
	Sredf(format string, a ...any) string    // 返回构造后的红色字符串（带占位符）
	Syellowf(format string, a ...any) string // 返回构造后的黄色字符串（带占位符）
	Spurplef(format string, a ...any) string // 返回构造后的紫色字符串（带占位符）
	PrintSuccessf(format string, a ...any)   // 打印成功信息到控制台（带占位符）
	PrintErrorf(format string, a ...any)     // 打印错误信息到控制台（带占位符）
	PrintWarningf(format string, a ...any)   // 打印警告信息到控制台（带占位符）
	PrintInfof(format string, a ...any)      // 打印信息到控制台（带占位符）
	PrintDebugf(format string, a ...any)     // 打印调试信息到控制台（带占位符）

	// 直接打印信息, 无需占位符
	Blue(msg ...any)           // 打印蓝色信息到控制台, 无需占位符
	Green(msg ...any)          // 打印绿色信息到控制台, 无需占位符
	Red(msg ...any)            // 打印红色信息到控制台, 无需占位符
	Yellow(msg ...any)         // 打印黄色信息到控制台, 无需占位符
	Purple(msg ...any)         // 打印紫色信息到控制台, 无需占位符
	Sblue(msg ...any) string   // 返回构造后的蓝色字符串, 无需占位符
	Sgreen(msg ...any) string  // 返回构造后的绿色字符串, 无需占位符
	Sred(msg ...any) string    // 返回构造后的红色字符串, 无需占位符
	Syellow(msg ...any) string // 返回构造后的黄色字符串, 无需占位符
	Spurple(msg ...any) string // 返回构造后的紫色字符串, 无需占位符
	PrintSuccess(msg ...any)   // 打印成功信息到控制台, 无需占位符
	PrintError(msg ...any)     // 打印错误信息到控制台, 无需占位符
	PrintWarning(msg ...any)   // 打印警告信息到控制台, 无需占位符
	PrintInfo(msg ...any)      // 打印信息到控制台, 无需占位符
	PrintDebug(msg ...any)     // 打印调试信息到控制台, 无需占位符

	// 新增扩展颜色的方法
	Black(msg ...any)                         // 打印黑色信息到控制台, 无需占位符
	Blackf(format string, a ...any)           // 打印黑色信息到控制台（带占位符）
	Sblack(msg ...any) string                 // 返回构造后的黑色字符串, 无需占位符
	Sblackf(format string, a ...any) string   // 返回构造后的黑色字符串（带占位符）
	Cyan(msg ...any)                          // 打印青色信息到控制台, 无需占位符
	Cyanf(format string, a ...any)            // 打印青色信息到控制台（带占位符）
	Scyan(msg ...any) string                  // 返回构造后的青色字符串, 无需占位符
	Scyanf(format string, a ...any) string    // 返回构造后的青色字符串（带占位符）
	White(msg ...any)                         // 打印白色信息到控制台, 无需占位符
	Whitef(format string, a ...any)           // 打印白色信息到控制台（带占位符）
	Swhite(msg ...any) string                 // 返回构造后的白色字符串, 无需占位符
	Swhitef(format string, a ...any) string   // 返回构造后的白色字符串（带占位符）
	Gray(msg ...any)                          // 打印灰色信息到控制台, 无需占位符
	Grayf(format string, a ...any)            // 打印灰色信息到控制台（带占位符）
	Sgray(msg ...any) string                  // 返回构造后的灰色字符串, 无需占位符
	Sgrayf(format string, a ...any) string    // 返回构造后的灰色字符串（带占位符）
	Lred(msg ...any)                          // 打印亮红色信息到控制台, 无需占位符
	Lredf(format string, a ...any)            // 打印亮红色信息到控制台（带占位符）
	Slred(msg ...any) string                  // 返回构造后的亮红色字符串, 无需占位符
	Slredf(format string, a ...any) string    // 返回构造后的亮红色字符串（带占位符）
	Lgreen(msg ...any)                        // 打印亮绿色信息到控制台, 无需占位符
	Lgreenf(format string, a ...any)          // 打印亮绿色信息到控制台（带占位符）
	Slgreen(msg ...any) string                // 返回构造后的亮绿色字符串, 无需占位符
	Slgreenf(format string, a ...any) string  // 返回构造后的亮绿色字符串（带占位符）
	Lyellow(msg ...any)                       // 打印亮黄色信息到控制台, 无需占位符
	Lyellowf(format string, a ...any)         // 打印亮黄色信息到控制台（带占位符）
	Slyellow(msg ...any) string               // 返回构造后的亮黄色字符串, 无需占位符
	Slyellowf(format string, a ...any) string // 返回构造后的亮黄色字符串（带占位符）
	Lblue(msg ...any)                         // 打印亮蓝色信息到控制台, 无需占位符
	Lbluef(format string, a ...any)           // 打印亮蓝色信息到控制台（带占位符）
	Slblue(msg ...any) string                 // 返回构造后的亮蓝色字符串, 无需占位符
	Slbluef(format string, a ...any) string   // 返回构造后的亮蓝色字符串（带占位符）
	Lpurple(msg ...any)                       // 打印亮紫色信息到控制台, 无需占位符
	Lpurplef(format string, a ...any)         // 打印亮紫色信息到控制台（带占位符）
	Slpurple(msg ...any) string               // 返回构造后的亮紫色字符串, 无需占位符
	Slpurplef(format string, a ...any) string // 返回构造后的亮紫色字符串（带占位符）
	Lcyan(msg ...any)                         // 打印亮青色信息到控制台, 无需占位符
	Lcyanf(format string, a ...any)           // 打印亮青色信息到控制台（带占位符）
	Slcyan(msg ...any) string                 // 返回构造后的亮青色字符串, 无需占位符
	Slcyanf(format string, a ...any) string   // 返回构造后的亮青色字符串（带占位符）
	Lwhite(msg ...any)                        // 打印亮白色信息到控制台, 无需占位符
	Lwhitef(format string, a ...any)          // 打印亮白色信息到控制台（带占位符）
	Slwhite(msg ...any) string                // 返回构造后的亮白色字符串, 无需占位符
	Slwhitef(format string, a ...any) string  // 返回构造后的亮白色字符串（带占位符）

	// 新增简洁版的方法, 无需占位符
	PrintOk(msg ...any)   // 打印成功信息到控制台, 无需占位符
	PrintErr(msg ...any)  // 打印错误信息到控制台, 无需占位符
	PrintInf(msg ...any)  // 打印信息到控制台, 无需占位符
	PrintDbg(msg ...any)  // 打印调试信息到控制台, 无需占位符
	PrintWarn(msg ...any) // 打印警告信息到控制台, 无需占位符

	// 新增简洁版的方法, 带占位符
	PrintOkf(format string, a ...any)   // 打印成功信息到控制台（带占位符）
	PrintErrf(format string, a ...any)  // 打印错误信息到控制台（带占位符）
	PrintInff(format string, a ...any)  // 打印信息到控制台（带占位符）
	PrintDbgf(format string, a ...any)  // 打印调试信息到控制台（带占位符）
	PrintWarnf(format string, a ...any) // 打印警告信息到控制台（带占位符）
}

// NewColorLib 函数用于创建一个新的 ColorLib 实例
func NewColorLib() *ColorLib {
	// 创建一个新的 ColorLib 实例
	cl := &ColorLib{
		levelMap: make(map[string]string),
		colorMap: make(map[string]int),
	}

	// 初始化颜色映射
	for k, v := range colorMap {
		cl.colorMap[k] = v
	}

	// 初始化日志级别映射
	for k, v := range levelMap {
		cl.levelMap[k] = v
	}

	return cl
}

// printWithColor 方法用于将传入的参数以指定颜色文本形式打印到控制台。
func (c *ColorLib) printWithColor(color string, msg ...any) {
	// 检查是否禁用颜色输出
	if c.NoColor {
		fmt.Print(msg...)
		return
	}

	// 获取颜色代码
	code, ok := c.colorMap[color]
	if !ok {
		fmt.Println("Invalid color:", color)
		return
	}

	// 清理缓冲区
	c.formatBuffer.Reset()

	// 检查是否禁用粗体输出
	if c.NoBold {
		// 写入前缀
		c.formatBuffer.WriteString(fmt.Sprintf("\033[%dm", code))
	} else {
		// 写入前缀
		c.formatBuffer.WriteString(fmt.Sprintf("\033[1;%dm", code))
	}

	// 写入消息
	if len(msg) > 0 {
		c.formatBuffer.WriteString(fmt.Sprint(msg...)) // 拼接消息内容
	} else {
		c.formatBuffer.WriteString(" ") // 如果没有消息，添加一个空格，避免完全空白的输出
	}

	// 写入颜色重置代码
	c.formatBuffer.WriteString(fmt.Sprintf("\033[%dm", reset))

	// 使用 fmt.Print 根据外部调用选择性添加换行符
	fmt.Print(c.formatBuffer.String())

	// 重置缓冲区
	c.formatBuffer.Reset()
}

// returnWithColor 方法用于将传入的参数以指定颜色文本形式返回。
func (c *ColorLib) returnWithColor(color string, msg ...any) string {
	// 检查是否禁用颜色输出
	if c.NoColor {
		return fmt.Sprint(msg...)
	}

	// 获取颜色代码
	code, ok := c.colorMap[color]
	if !ok {
		return fmt.Sprintf("Invalid color: %s", color)
	}

	// 检查 msg 是否为空
	if len(msg) == 0 {
		if c.NoBold {
			return fmt.Sprintf("\033[%dm\033[%dm", code, reset) // 返回空字符串，但带有颜色代码
		} else {
			return fmt.Sprintf("\033[1;%dm\033[%dm", code, reset) // 返回空字符串，但带有颜色代码
		}
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)

	// 清理缓冲区
	c.formatBuffer.Reset()

	// 写入前缀
	if c.NoBold {
		c.formatBuffer.WriteString(fmt.Sprintf("\033[%dm", code))
	} else {
		c.formatBuffer.WriteString(fmt.Sprintf("\033[1;%dm", code)) // 添加颜色代码，并加粗
	}

	// 写入消息
	c.formatBuffer.WriteString(combinedMsg) // 拼接消息内容

	// 写入颜色重置代码
	c.formatBuffer.WriteString(fmt.Sprintf("\033[%dm", reset))

	// 获取最终字符串
	result := c.formatBuffer.String()

	// 重置缓冲区
	c.formatBuffer.Reset()

	return result
}

// PromptMsg 方法用于打印带有颜色和前缀的消息。
func (c *ColorLib) PromptMsg(level, color, format string, a ...any) {
	// 获取指定级别对应的前缀
	prefix, ok := c.levelMap[level]
	if !ok {
		fmt.Println("Invalid level:", level)
		return
	}

	// 清理缓冲区
	c.formatBuffer.Reset()

	// 写入前缀
	c.formatBuffer.WriteString(prefix)

	// 如果没有参数，直接打印前缀
	if len(a) == 0 {
		if c.NoColor {
			fmt.Print(c.formatBuffer.String())
			c.formatBuffer.Reset()
		} else {
			c.printWithColor(color, c.formatBuffer.String())
			c.formatBuffer.Reset()
		}
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprintf(format, a...)

	// 写入消息
	c.formatBuffer.WriteString(combinedMsg)

	// 打印最终消息
	if c.NoColor {
		fmt.Print(c.formatBuffer.String())
	} else {
		c.printWithColor(color, c.formatBuffer.String())
	}

	// 重置缓冲区
	c.formatBuffer.Reset()
}

// PMsg 方法用于打印带有颜色和前缀的消息。
func (c *ColorLib) PMsg(level, color, format string, a ...any) {
	// 获取指定级别对应的前缀
	prefix, ok := c.levelMap[level]
	if !ok {
		fmt.Println("Invalid level:", level)
		return
	}

	// 清理缓冲区
	c.formatBuffer.Reset()

	// 写入前缀
	c.formatBuffer.WriteString(prefix)

	// 如果没有参数，直接打印前缀
	if len(a) == 0 {
		if c.NoColor {
			fmt.Print(c.formatBuffer.String())
			c.formatBuffer.Reset()
		} else {
			c.printWithColor(color, c.formatBuffer.String())
			c.formatBuffer.Reset()
		}
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprintf(format, a...)

	// 写入消息
	c.formatBuffer.WriteString(combinedMsg)

	// 打印最终消息
	if c.NoColor {
		fmt.Print(c.formatBuffer.String())
	} else {
		c.printWithColor(color, c.formatBuffer.String())
	}

	// 重置缓冲区
	c.formatBuffer.Reset()
}

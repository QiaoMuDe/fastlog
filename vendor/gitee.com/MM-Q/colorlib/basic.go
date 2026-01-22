// Package colorlib 提供了基础的颜色输出功能，包括格式化输出、直接打印和字符串返回等方法。
// 该文件实现了 ColorLib 结构体的基础颜色方法，支持标准颜色和亮色系列的文本输出。
package colorlib

import (
	"fmt"

	"gitee.com/MM-Q/colorlib/internal/color"
)

// ==================================================================
// 颜色格式化方法
// ==================================================================

// Bluef 方法用于将传入的参数以蓝色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) Bluef(format string, a ...any) {
	c.printWithColor(color.Blue, fmt.Sprintf(format, a...), false)
}

// Greenf 方法用于将传入的参数以绿色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) Greenf(format string, a ...any) {
	c.printWithColor(color.Green, fmt.Sprintf(format, a...), false)
}

// Redf 方法用于将传入的参数以红色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) Redf(format string, a ...any) {
	c.printWithColor(color.Red, fmt.Sprintf(format, a...), false)
}

// Yellowf 方法用于将传入的参数以黄色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) Yellowf(format string, a ...any) {
	c.printWithColor(color.Yellow, fmt.Sprintf(format, a...), false)
}

// Magentaf 方法用于将传入的参数以品红色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) Magentaf(format string, a ...any) {
	c.printWithColor(color.Magenta, fmt.Sprintf(format, a...), false)
}

// Blackf 方法用于将传入的参数以黑色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) Blackf(format string, a ...any) {
	c.printWithColor(color.Black, fmt.Sprintf(format, a...), false)
}

// Cyanf 方法用于将传入的参数以青色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) Cyanf(format string, a ...any) {
	c.printWithColor(color.Cyan, fmt.Sprintf(format, a...), false)
}

// Whitef 方法用于将传入的参数以白色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) Whitef(format string, a ...any) {
	c.printWithColor(color.White, fmt.Sprintf(format, a...), false)
}

// Grayf 方法用于将传入的参数以灰色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) Grayf(format string, a ...any) {
	c.printWithColor(color.Gray, fmt.Sprintf(format, a...), false)
}

// BrightRedf 方法用于将传入的参数以亮红色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) BrightRedf(format string, a ...any) {
	c.printWithColor(color.BrightRed, fmt.Sprintf(format, a...), false)
}

// BrightGreenf 方法用于将传入的参数以亮绿色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) BrightGreenf(format string, a ...any) {
	c.printWithColor(color.BrightGreen, fmt.Sprintf(format, a...), false)
}

// BrightYellowf 方法用于将传入的参数以亮黄色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) BrightYellowf(format string, a ...any) {
	c.printWithColor(color.BrightYellow, fmt.Sprintf(format, a...), false)
}

// BrightBluef 方法用于将传入的参数以亮蓝色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) BrightBluef(format string, a ...any) {
	c.printWithColor(color.BrightBlue, fmt.Sprintf(format, a...), false)
}

// BrightMagentaf 方法用于将传入的参数以亮品红色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) BrightMagentaf(format string, a ...any) {
	c.printWithColor(color.BrightMagenta, fmt.Sprintf(format, a...), false)
}

// BrightCyanf 方法用于将传入的参数以亮青色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) BrightCyanf(format string, a ...any) {
	c.printWithColor(color.BrightCyan, fmt.Sprintf(format, a...), false)
}

// BrightWhitef 方法用于将传入的参数以亮白色文本形式打印到控制台
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
func (c *ColorLib) BrightWhitef(format string, a ...any) {
	c.printWithColor(color.BrightWhite, fmt.Sprintf(format, a...), false)
}

// ==================================================================
// 颜色打印方法
// ==================================================================

// Blue 方法用于将传入的参数以蓝色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) Blue(msg ...any) {
	c.printWithColor(color.Blue, fmt.Sprint(msg...), true)
}

// Green 方法用于将传入的参数以绿色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) Green(msg ...any) {
	c.printWithColor(color.Green, fmt.Sprint(msg...), true)
}

// Red 方法用于将传入的参数以红色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) Red(msg ...any) {
	c.printWithColor(color.Red, fmt.Sprint(msg...), true)
}

// Yellow 方法用于将传入的参数以黄色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) Yellow(msg ...any) {
	c.printWithColor(color.Yellow, fmt.Sprint(msg...), true)
}

// Magenta 方法用于将传入的参数以品红色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) Magenta(msg ...any) {
	c.printWithColor(color.Magenta, fmt.Sprint(msg...), true)
}

// Black 方法用于将传入的参数以黑色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) Black(msg ...any) {
	c.printWithColor(color.Black, fmt.Sprint(msg...), true)
}

// Cyan 方法用于将传入的参数以青色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) Cyan(msg ...any) {
	c.printWithColor(color.Cyan, fmt.Sprint(msg...), true)
}

// White 方法用于将传入的参数以白色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) White(msg ...any) {
	c.printWithColor(color.White, fmt.Sprint(msg...), true)
}

// Gray 方法用于将传入的参数以灰色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) Gray(msg ...any) {
	c.printWithColor(color.Gray, fmt.Sprint(msg...), true)
}

// BrightRed 方法用于将传入的参数以亮红色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) BrightRed(msg ...any) {
	c.printWithColor(color.BrightRed, fmt.Sprint(msg...), true)
}

// BrightGreen 方法用于将传入的参数以亮绿色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) BrightGreen(msg ...any) {
	c.printWithColor(color.BrightGreen, fmt.Sprint(msg...), true)
}

// BrightYellow 方法用于将传入的参数以亮黄色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) BrightYellow(msg ...any) {
	c.printWithColor(color.BrightYellow, fmt.Sprint(msg...), true)
}

// BrightBlue 方法用于将传入的参数以亮蓝色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) BrightBlue(msg ...any) {
	c.printWithColor(color.BrightBlue, fmt.Sprint(msg...), true)
}

// BrightMagenta 方法用于将传入的参数以亮品红色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) BrightMagenta(msg ...any) {
	c.printWithColor(color.BrightMagenta, fmt.Sprint(msg...), true)
}

// BrightCyan 方法用于将传入的参数以亮青色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) BrightCyan(msg ...any) {
	c.printWithColor(color.BrightCyan, fmt.Sprint(msg...), true)
}

// BrightWhite 方法用于将传入的参数以亮白色文本形式打印到控制台
//
// 参数:
//   - msg: 可变参数，要打印的消息内容
func (c *ColorLib) BrightWhite(msg ...any) {
	c.printWithColor(color.BrightWhite, fmt.Sprint(msg...), true)
}

// ==================================================================
// 颜色返回方法
// ==================================================================

// Sblue 方法用于将传入的参数以蓝色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有蓝色格式的字符串
func (c *ColorLib) Sblue(msg ...any) string {
	return c.returnWithColor(color.Blue, fmt.Sprint(msg...))
}

// Sgreen 方法用于将传入的参数以绿色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有绿色格式的字符串
func (c *ColorLib) Sgreen(msg ...any) string {
	return c.returnWithColor(color.Green, fmt.Sprint(msg...))
}

// Sred 方法用于将传入的参数以红色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有红色格式的字符串
func (c *ColorLib) Sred(msg ...any) string {
	return c.returnWithColor(color.Red, fmt.Sprint(msg...))
}

// Syellow 方法用于将传入的参数以黄色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有黄色格式的字符串
func (c *ColorLib) Syellow(msg ...any) string {
	return c.returnWithColor(color.Yellow, fmt.Sprint(msg...))
}

// Smagenta 方法用于将传入的参数以品红色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有品红色格式的字符串
func (c *ColorLib) Smagenta(msg ...any) string {
	return c.returnWithColor(color.Magenta, fmt.Sprint(msg...))
}

// Sblack 方法用于将传入的参数以黑色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有黑色格式的字符串
func (c *ColorLib) Sblack(msg ...any) string {
	return c.returnWithColor(color.Black, fmt.Sprint(msg...))
}

// Scyan 方法用于将传入的参数以青色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有青色格式的字符串
func (c *ColorLib) Scyan(msg ...any) string {
	return c.returnWithColor(color.Cyan, fmt.Sprint(msg...))
}

// Swhite 方法用于将传入的参数以白色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有白色格式的字符串
func (c *ColorLib) Swhite(msg ...any) string {
	return c.returnWithColor(color.White, fmt.Sprint(msg...))
}

// Sgray 方法用于将传入的参数以灰色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有灰色格式的字符串
func (c *ColorLib) Sgray(msg ...any) string {
	return c.returnWithColor(color.Gray, fmt.Sprint(msg...))
}

// SbrightRed 方法用于将传入的参数以亮红色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有亮红色格式的字符串
func (c *ColorLib) SbrightRed(msg ...any) string {
	return c.returnWithColor(color.BrightRed, fmt.Sprint(msg...))
}

// SbrightGreen 方法用于将传入的参数以亮绿色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有亮绿色格式的字符串
func (c *ColorLib) SbrightGreen(msg ...any) string {
	return c.returnWithColor(color.BrightGreen, fmt.Sprint(msg...))
}

// SbrightYellow 方法用于将传入的参数以亮黄色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有亮黄色格式的字符串
func (c *ColorLib) SbrightYellow(msg ...any) string {
	return c.returnWithColor(color.BrightYellow, fmt.Sprint(msg...))
}

// SbrightBlue 方法用于将传入的参数以亮蓝色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有亮蓝色格式的字符串
func (c *ColorLib) SbrightBlue(msg ...any) string {
	return c.returnWithColor(color.BrightBlue, fmt.Sprint(msg...))
}

// SbrightMagenta 方法用于将传入的参数以亮品红色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有亮品红色格式的字符串
func (c *ColorLib) SbrightMagenta(msg ...any) string {
	return c.returnWithColor(color.BrightMagenta, fmt.Sprint(msg...))
}

// SbrightCyan 方法用于将传入的参数以亮青色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有亮青色格式的字符串
func (c *ColorLib) SbrightCyan(msg ...any) string {
	return c.returnWithColor(color.BrightCyan, fmt.Sprint(msg...))
}

// SbrightWhite 方法用于将传入的参数以亮白色文本形式返回
//
// 参数:
//   - msg: 可变参数，要处理的消息内容
//
// 返回值:
//   - string: 带有亮白色格式的字符串
func (c *ColorLib) SbrightWhite(msg ...any) string {
	return c.returnWithColor(color.BrightWhite, fmt.Sprint(msg...))
}

// ==================================================================
// 颜色格式化方法
// ==================================================================

// Sbluef 方法用于将传入的参数以蓝色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有蓝色格式的字符串
func (c *ColorLib) Sbluef(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor(color.Blue, fmt.Sprintf(format, a...))
}

// Sgreenf 方法用于将传入的参数以绿色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有绿色格式的字符串
func (c *ColorLib) Sgreenf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor(color.Green, fmt.Sprintf(format, a...))
}

// Sredf 方法用于将传入的参数以红色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有红色格式的字符串
func (c *ColorLib) Sredf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor(color.Red, fmt.Sprintf(format, a...))
}

// Syellowf 方法用于将传入的参数以黄色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有黄色格式的字符串
func (c *ColorLib) Syellowf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor(color.Yellow, fmt.Sprintf(format, a...))
}

// Smagentaf 方法用于将传入的参数以品红色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有品红色格式的字符串
func (c *ColorLib) Smagentaf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor(color.Magenta, fmt.Sprintf(format, a...))
}

// Sblackf 方法用于将传入的参数以黑色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有黑色格式的字符串
func (c *ColorLib) Sblackf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor(color.Black, fmt.Sprintf(format, a...))
}

// Scyanf 方法用于将传入的参数以青色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有青色格式的字符串
func (c *ColorLib) Scyanf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor(color.Cyan, fmt.Sprintf(format, a...))
}

// Swhitef 方法用于将传入的参数以白色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有白色格式的字符串
func (c *ColorLib) Swhitef(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor(color.White, fmt.Sprintf(format, a...))
}

// Sgrayf 方法用于将传入的参数以灰色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有灰色格式的字符串
func (c *ColorLib) Sgrayf(format string, a ...any) string {
	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor(color.Gray, fmt.Sprintf(format, a...))
}

// SbrightRedf 方法用于将传入的参数以亮红色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有亮红色格式的字符串
func (c *ColorLib) SbrightRedf(format string, a ...any) string {
	return c.returnWithColor(color.BrightRed, fmt.Sprintf(format, a...))
}

// SbrightGreenf 方法用于将传入的参数以亮绿色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有亮绿色格式的字符串
func (c *ColorLib) SbrightGreenf(format string, a ...any) string {
	return c.returnWithColor(color.BrightGreen, fmt.Sprintf(format, a...))
}

// SbrightYellowf 方法用于将传入的参数以亮黄色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有亮黄色格式的字符串
func (c *ColorLib) SbrightYellowf(format string, a ...any) string {
	return c.returnWithColor(color.BrightYellow, fmt.Sprintf(format, a...))
}

// SbrightBluef 方法用于将传入的参数以亮蓝色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有亮蓝色格式的字符串
func (c *ColorLib) SbrightBluef(format string, a ...any) string {
	return c.returnWithColor(color.BrightBlue, fmt.Sprintf(format, a...))
}

// SbrightMagentaf 方法用于将传入的参数以亮品红色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有亮品红色格式的字符串
func (c *ColorLib) SbrightMagentaf(format string, a ...any) string {
	return c.returnWithColor(color.BrightMagenta, fmt.Sprintf(format, a...))
}

// SbrightCyanf 方法用于将传入的参数以亮青色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有亮青色格式的字符串
func (c *ColorLib) SbrightCyanf(format string, a ...any) string {
	return c.returnWithColor(color.BrightCyan, fmt.Sprintf(format, a...))
}

// SbrightWhitef 方法用于将传入的参数以亮白色文本形式返回
//
// 参数:
//   - format: 格式化字符串，用于指定输出的格式
//   - a: 可变参数，用于填充格式化字符串中的占位符
//
// 返回值:
//   - string: 带有亮白色格式的字符串
func (c *ColorLib) SbrightWhitef(format string, a ...any) string {
	return c.returnWithColor(color.BrightWhite, fmt.Sprintf(format, a...))
}

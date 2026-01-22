// Package colorlib 提供了内部核心功能实现。
// 该文件实现了 ColorLib 结构体的核心方法，包括颜色输出、ANSI 序列构建、
// 样式处理和消息格式化等功能，是整个颜色库的核心实现文件。
package colorlib

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/go-kit/pool"
)

// printWithColor 方法用于将传入的参数以指定颜色文本形式打印到控制台
//
// 参数:
//   - colorCode: 颜色代码
//   - msg: 要打印的消息
//   - needLineFeed: 是否需要换行
func (c *ColorLib) printWithColor(colorCode int, msg string, needLineFeed bool) {
	// 提前检查消息是否为空
	if len(msg) == 0 {
		return // 直接返回，什么都不打印
	}

	// 检查是否禁用颜色输出
	if !c.configMgr.GetColor() {
		// 如果需要换行
		if needLineFeed {
			_, _ = fmt.Fprintln(c.writer, msg)
			return
		}

		// 如果不需要换行
		_, _ = fmt.Fprint(c.writer, msg)
		return
	}

	// 检查颜色代码是否有效
	if ok := c.colorMgr.IsColorCode(colorCode); !ok {
		_, _ = fmt.Fprintln(c.writer, "Invalid color:", colorCode)
		return
	}

	// 直接构建字符串: 预估容量：\033[ + 最多4个样式代码(每个2位) + 分号 + 颜色代码(2位) + m + 消息 + \033[0m
	result := pool.WithStrCap(len(msg)+32, func(builder *strings.Builder) {
		// 使用辅助方法构建ANSI序列
		c.buildAnsiSequence(builder, colorCode)

		// 写入消息
		_, _ = builder.WriteString(msg)

		// 写入重置代码
		_, _ = builder.WriteString("\033[0m")
	})

	// 输出结果到指定的writer
	if needLineFeed {
		// 如果需要换行
		_, _ = fmt.Fprintln(c.writer, result)
		return
	}

	// 如果不需要换行
	_, _ = fmt.Fprint(c.writer, result)
}

// returnWithColor 方法用于将传入的参数以指定颜色文本形式返回
//
// 参数:
//   - colorCode: 指定颜色代码
//   - msg: 要处理的消息
//
// 返回值:
//   - string: 带颜色的字符串
func (c *ColorLib) returnWithColor(colorCode int, msg string) string {
	// 检查 msg 是否为空
	if len(msg) == 0 {
		return "" // 直接返回空字符串，简单明了
	}

	// 检查是否禁用颜色输出
	if !c.configMgr.GetColor() {
		return fmt.Sprint(msg)
	}

	// 检查颜色代码是否有效
	if ok := c.colorMgr.IsColorCode(colorCode); !ok {
		return fmt.Sprintf("Invalid color: %d", colorCode)
	}

	// 使用字符串构建器构建带颜色的字符串
	return pool.WithStrCap(len(msg)+32, func(builder *strings.Builder) {
		// 使用辅助方法构建ANSI序列
		c.buildAnsiSequence(builder, colorCode)

		// 写入消息内容
		_, _ = builder.WriteString(msg)

		// 写入重置代码
		_, _ = builder.WriteString("\033[0m")
	})

}

// buildAnsiSequence 方法用于构建ANSI控制序列（使用缓存优化）
//
// 参数:
//   - builder: 字符串构建器
//   - colorCode: 颜色代码
func (c *ColorLib) buildAnsiSequence(builder *strings.Builder, colorCode int) {
	// 使用缓存获取ANSI序列，显著提升性能
	ansiSeq := c.ansiCache.GetANSI(
		colorCode,
		c.configMgr.GetBold(),
		c.configMgr.GetUnderline(),
		c.configMgr.GetBlink(),
	)
	_, _ = builder.WriteString(ansiSeq)
}

// promptMsg 方法用于打印带有指定前缀的消息
//
// 参数:
//   - level: 日志级别
//   - colorCode: 颜色代码
//   - needLineFeed: 是否需要换行
//   - format: 格式化字符串
//   - a: 格式化参数
func (c *ColorLib) promptMsg(level int, colorCode int, needLineFeed bool, format string, a ...any) {
	// 检查a是否为空
	if len(a) == 0 {
		return
	}

	// 获取指定级别对应的前缀
	prefix, ok := levelMap[level]
	if !ok {
		_, _ = fmt.Fprintln(c.writer, "Invalid level:", level)
		return
	}

	// 使用字符串构建器构建消息
	message := pool.WithStrCap(len(prefix)+len(format)+32, func(builder *strings.Builder) {
		// 写入前缀
		_, _ = builder.WriteString(prefix)

		// 如果有参数，格式化并写入消息
		if len(a) > 0 {
			combinedMsg := fmt.Sprintf(format, a...)
			_, _ = builder.WriteString(combinedMsg)
		}
	})

	// 使用颜色打印
	c.printWithColor(colorCode, message, needLineFeed)
}

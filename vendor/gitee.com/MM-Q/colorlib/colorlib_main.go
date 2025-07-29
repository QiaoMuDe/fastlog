package colorlib

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
)

// GetCL 是一个线程安全用于获取全局唯一的 ColorLib 实例的函数
func GetCL() *ColorLib {
	once.Do(func() {
		CL = NewColorLib()
	})
	return CL
}

// NewColorLib 函数用于创建一个新的 ColorLib 实例
func NewColorLib() *ColorLib {
	// 创建一个新的 ColorLib 实例
	cl := &ColorLib{
		levelMap: sync.Map{}, // 使用 sync.Map 来存储日志级别映射
		colorMap: sync.Map{}, // 使用 sync.Map 来存储颜色映射
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			}, // 使用 sync.Pool 来管理缓冲区
		},
	}

	// 初始化是否禁用颜色
	cl.NoColor.Store(false)

	// 初始化是否禁用加粗
	cl.NoBold.Store(false)

	// 初始化是否下划线
	cl.Underline.Store(false)

	// 初始化是否闪烁
	cl.Blink.Store(false)

	// 初始化颜色映射
	for k, v := range colorMap {
		cl.colorMap.Store(k, v)
	}

	// 初始化日志级别映射
	for k, v := range levelMap {
		cl.levelMap.Store(k, v)
	}

	return cl
}

// SetNoColor 设置是否禁用颜色输出,并返回ColorLib实例以支持链式调用
func (c *ColorLib) SetNoColor(enable bool) *ColorLib {
	c.NoColor.Store(enable)
	return c
}

// SetNoBold 设置是否禁用字体加粗,并返回ColorLib实例以支持链式调用
func (c *ColorLib) SetNoBold(enable bool) *ColorLib {
	c.NoBold.Store(enable)
	return c
}

// SetUnderline 设置是否启用下划线,并返回ColorLib实例以支持链式调用
func (c *ColorLib) SetUnderline(enable bool) *ColorLib {
	c.Underline.Store(enable)
	return c
}

// SetBlink 设置是否启用闪烁效果,并返回ColorLib实例以支持链式调用
func (c *ColorLib) SetBlink(enable bool) *ColorLib {
	c.Blink.Store(enable)
	return c
}

// printWithColor 方法用于将传入的参数以指定颜色文本形式打印到控制台。
func (c *ColorLib) printWithColor(color string, msg ...any) {
	// 检查是否禁用颜色输出
	if c.NoColor.Load() {
		fmt.Print(msg...)
		return
	}

	// 获取颜色代码
	code, ok := c.colorMap.Load(color)
	if !ok {
		fmt.Println("Invalid color:", color)
		return
	}

	// 从对象池获取缓冲区
	buffer := c.bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buffer.Reset()           // 清空缓冲区
		c.bufferPool.Put(buffer) // 将缓冲区放回对象池
	}()

	// 构建ANSI控制序列
	var ansiCodes []string
	if !c.NoBold.Load() { // 检查是否加粗
		ansiCodes = append(ansiCodes, "1")
	}
	if c.Underline.Load() { // 检查是否下划线
		ansiCodes = append(ansiCodes, "4")
	}
	if c.Blink.Load() { // 检查是否闪烁
		ansiCodes = append(ansiCodes, "5")
	}

	// 添加颜色代码
	ansiCodes = append(ansiCodes, fmt.Sprintf("%d", code))

	// 写入前缀
	buffer.WriteString(fmt.Sprintf("\033[%sm", strings.Join(ansiCodes, ";")))

	// 写入消息
	if len(msg) > 0 {
		buffer.WriteString(fmt.Sprint(msg...)) // 拼接消息内容
	} else {
		buffer.WriteString(" ") // 如果没有消息,添加一个空格,避免完全空白的输出
	}

	// 写入颜色重置代码
	buffer.WriteString("\033[0m")

	// 使用 fmt.Print 根据外部调用选择性添加换行符
	fmt.Print(buffer.String())
}

// returnWithColor 方法用于将传入的参数以指定颜色文本形式返回。
func (c *ColorLib) returnWithColor(color string, msg ...any) string {
	// 检查是否禁用颜色输出
	if c.NoColor.Load() {
		return fmt.Sprint(msg...)
	}

	// 获取颜色代码
	code, ok := c.colorMap.Load(color)
	if !ok {
		return fmt.Sprintf("Invalid color: %s", color)
	}

	// 检查 msg 是否为空
	if len(msg) == 0 {
		var ansiCodes []string
		if !c.NoBold.Load() {
			ansiCodes = append(ansiCodes, "1")
		}
		if c.Underline.Load() {
			ansiCodes = append(ansiCodes, "4")
		}
		if c.Blink.Load() {
			ansiCodes = append(ansiCodes, "5")
		}
		ansiCodes = append(ansiCodes, fmt.Sprintf("%d", code))
		return fmt.Sprintf("\033[%sm\033[0m", strings.Join(ansiCodes, ";"))
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)

	// 从对象池获取缓冲区
	buffer := c.bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buffer.Reset()           // 清空缓冲区
		c.bufferPool.Put(buffer) // 将缓冲区放回对象池
	}()

	// 构建ANSI控制序列
	var ansiCodes []string
	if !c.NoBold.Load() {
		ansiCodes = append(ansiCodes, "1")
	}
	if c.Underline.Load() {
		ansiCodes = append(ansiCodes, "4")
	}
	if c.Blink.Load() {
		ansiCodes = append(ansiCodes, "5")
	}
	ansiCodes = append(ansiCodes, fmt.Sprintf("%d", code))

	// 写入前缀
	buffer.WriteString(fmt.Sprintf("\033[%sm", strings.Join(ansiCodes, ";")))

	// 写入消息
	buffer.WriteString(combinedMsg) // 拼接消息内容

	// 写入颜色重置代码
	buffer.WriteString("\033[0m")

	// 获取最终字符串
	result := buffer.String()

	return result
}

// promptMsg 方法用于打印带有指定前缀的消息。
func (c *ColorLib) promptMsg(level, color, format string, a ...any) {
	// 获取指定级别对应的前缀
	prefix, ok := c.levelMap.Load(level)
	if !ok {
		fmt.Println("Invalid level:", level)
		return
	}

	// 从对象池获取缓冲区
	buffer := c.bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buffer.Reset()           // 清空缓冲区
		c.bufferPool.Put(buffer) // 将缓冲区放回对象池
	}()

	// 写入前缀
	// 由于 prefix 是 interface{} 类型,需要进行类型断言才能作为字符串传递给 WriteString 方法
	if prefixStr, ok := prefix.(string); ok {
		buffer.WriteString(prefixStr)
	} else {
		// 转为字符串
		buffer.WriteString(fmt.Sprintf("%v", prefix))
	}

	// 如果没有参数,直接打印前缀
	if len(a) == 0 {
		if c.NoColor.Load() {
			fmt.Print(buffer.String())
		} else {
			c.printWithColor(color, buffer.String())
		}
		return
	}

	// 使用 fmt.Sprintf 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprintf(format, a...)

	// 写入消息
	buffer.WriteString(combinedMsg)

	// 打印最终消息
	if c.NoColor.Load() {
		fmt.Print(buffer.String())
	} else {
		c.printWithColor(color, buffer.String())
	}
}

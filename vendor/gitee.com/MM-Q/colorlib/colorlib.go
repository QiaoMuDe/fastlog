// Package colorlib 是一个功能强大的 Go 语言终端颜色输出库。
// 它提供了丰富的颜色输出功能，包括基础颜色、亮色、样式设置（粗体、下划线、闪烁）等。
// 支持链式调用、自定义输出接口、线程安全的全局实例等特性。
// 主要用于在终端中输出带有颜色和样式的文本，提升命令行程序的用户体验。
package colorlib

import (
	"io"
	"os"
	"sync"

	"gitee.com/MM-Q/colorlib/internal/cache"
	"gitee.com/MM-Q/colorlib/internal/color"
	"gitee.com/MM-Q/colorlib/internal/style"
)

var (
	clOnce sync.Once // 确保全局实例只被初始化一次
	cl     *ColorLib // 全局实例可导入直接使用
)

// ColorLib 结构体用于管理颜色输出和日志级别映射。
type ColorLib struct {
	// 样式配置管理器
	configMgr *style.StyleConfig

	// 颜色管理器
	colorMgr *color.ColorManager

	// ANSI序列缓存
	ansiCache *cache.ANSICache

	// 输出接口（不可变，初始化后不可修改）
	writer io.Writer
}

// ===========================================================
// 构造函数
// ===========================================================

// GetCL 是一个线程安全用于获取全局唯一的 ColorLib 实例的函数
//
// 返回值:
//   - *ColorLib: 全局唯一的 ColorLib 实例指针
func GetCL() *ColorLib {
	clOnce.Do(func() {
		cl = NewColorLib()
	})
	return cl
}

// New  创建一个新的 ColorLib 实例(NewColorLib 的别名)
//
// 返回值:
//   - *ColorLib: 新创建的 ColorLib 实例指针
var New = NewColorLib

// NewColorLib 函数用于创建一个新的 ColorLib 实例（默认输出到标准输出）
//
// 返回值:
//   - *ColorLib: 新创建的 ColorLib 实例指针
func NewColorLib() *ColorLib {
	return NewColorLibWithWriter(os.Stdout)
}

// NewColorLibWithWriter 创建一个指定输出接口的 ColorLib 实例
//
// 参数:
//   - writer: 输出接口，如 os.Stdout, os.Stderr, 文件等
//
// 返回值:
//   - *ColorLib: 新创建的 ColorLib 实例指针
func NewColorLibWithWriter(writer io.Writer) *ColorLib {
	if writer == nil {
		writer = os.Stdout // 防御性编程
	}

	// 创建一个新的 ColorLib 实例
	cl := &ColorLib{
		configMgr: style.NewStyleConfig(),  // 初始化样式配置
		colorMgr:  color.NewColorManager(), // 初始化颜色管理器
		ansiCache: cache.NewANSICache(),    // 初始化ANSI序列缓存
		writer:    writer,                  // 设置输出接口（不可变）
	}

	return cl
}

// WithWriter 创建一个使用指定输出接口的新实例（不可变设计）
//
// 参数:
//   - w: 输出接口
//
// 返回值:
//   - *ColorLib: 新的ColorLib实例
func (c *ColorLib) WithWriter(w io.Writer) *ColorLib {
	if w == nil {
		w = os.Stdout
	}
	return NewColorLibWithWriter(w)
}

// ===========================================================
// 样式设置方法
// ===========================================================

// SetColor 设置是否启用颜色输出
//
// 参数:
//   - enabled: 是否启用颜色输出（true - 启用，false - 禁用）
func (c *ColorLib) SetColor(enabled bool) {
	c.configMgr.SetColor(enabled)
}

// SetBold 设置是否启用粗体输出
//
// 参数:
//   - enabled: 是否启用粗体输出（true - 启用，false - 禁用）
func (c *ColorLib) SetBold(enabled bool) {
	c.configMgr.SetBold(enabled)
}

// SetUnderline 设置是否启用下划线输出
//
// 参数:
//   - enabled: 是否启用下划线输出（true - 启用，false - 禁用）
func (c *ColorLib) SetUnderline(enabled bool) {
	c.configMgr.SetUnderline(enabled)
}

// SetBlink 设置是否启用闪烁输出
//
// 参数:
//   - enabled: 是否启用闪烁输出（true - 启用，false - 禁用）
func (c *ColorLib) SetBlink(enabled bool) {
	c.configMgr.SetBlink(enabled)
}

// WithColor 设置是否启用颜色输出（链式调用）
//
// 参数:
//   - enabled: 是否启用颜色输出（true - 启用，false - 禁用）
//
// 返回:
//   - *ColorLib: 当前 ColorLib 对象
func (c *ColorLib) WithColor(enabled bool) *ColorLib {
	c.configMgr.SetColor(enabled)
	return c
}

// WithBold 设置是否启用粗体输出（链式调用）
//
// 参数:
//   - enabled: 是否启用粗体输出（true - 启用，false - 禁用）
//
// 返回:
//   - *ColorLib: 当前 ColorLib 对象
func (c *ColorLib) WithBold(enabled bool) *ColorLib {
	c.configMgr.SetBold(enabled)
	return c
}

// WithUnderline 设置是否启用下划线输出（链式调用）
//
// 参数:
//   - enabled: 是否启用下划线输出（true - 启用，false - 禁用）
//
// 返回:
//   - *ColorLib: 当前 ColorLib 对象
func (c *ColorLib) WithUnderline(enabled bool) *ColorLib {
	c.configMgr.SetUnderline(enabled)
	return c
}

// WithBlink 启用闪烁效果（链式调用）
//
// 参数:
//   - enabled: 是否启用闪烁效果（true - 启用，false - 禁用）
//
// 返回:
//   - *ColorLib: 当前 ColorLib 对象
func (c *ColorLib) WithBlink(enabled bool) *ColorLib {
	c.configMgr.SetBlink(enabled)
	return c
}

// ===========================================================
// 样式获取方法
// ===========================================================

// GetColor 获取颜色启用状态
//
// 返回值:
//   - bool: 颜色启用状态
func (c *ColorLib) GetColor() bool {
	return c.configMgr.GetColor()
}

// GetBold 获取粗体启用状态
//
// 返回值:
//   - bool: 粗体启用状态
func (c *ColorLib) GetBold() bool {
	return c.configMgr.GetBold()
}

// GetUnderline 获取下划线启用状态
//
// 返回值:
//   - bool: 下划线启用状态
func (c *ColorLib) GetUnderline() bool {
	return c.configMgr.GetUnderline()
}

// GetBlink 获取闪烁启用状态
//
// 返回值:
//   - bool: 闪烁启用状态
func (c *ColorLib) GetBlink() bool {
	return c.configMgr.GetBlink()
}

// ===========================================================
// 缓存统计方法
// ===========================================================

// GetCacheSize 获取当前缓存大小
//
// 返回值:
//   - int: 缓存中的条目数量
func (c *ColorLib) GetCacheSize() int {
	return c.ansiCache.GetCacheSize()
}

// ClearCache 清空ANSI序列缓存
func (c *ColorLib) ClearCache() {
	c.ansiCache.Clear()
}

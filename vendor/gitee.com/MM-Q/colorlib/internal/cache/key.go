// Package cache 提供了智能缓存键和ANSI序列缓存功能。
// 该文件实现了基于位运算的高效缓存键设计，用于缓存ANSI颜色序列，
// 支持颜色代码、粗体、下划线、闪烁等样式的组合缓存。
package cache

import "fmt"

// CacheKey 智能缓存键类型，使用32位整数压缩存储样式信息
type CacheKey uint32

const (
	// 样式位掩码定义
	BoldMask      = 1 << 15 // 0x8000 - 粗体标志位
	UnderlineMask = 1 << 14 // 0x4000 - 下划线标志位
	BlinkMask     = 1 << 13 // 0x2000 - 闪烁标志位

	// 颜色位掩码
	ColorMask = 0xFF // 0x00FF - 低8位用于颜色代码

	// 保留位用于未来扩展
	// 第12-8位: 保留给未来的样式扩展
	// 第31-16位: 保留给未来功能扩展
)

// BuildCacheKey 构建智能缓存键
//
// 参数:
//   - colorCode: ANSI颜色代码 (30-37, 90-97等)
//   - bold: 是否启用粗体
//   - underline: 是否启用下划线
//   - blink: 是否启用闪烁
//
// 返回值:
//   - CacheKey: 压缩后的缓存键
func BuildCacheKey(colorCode int, bold, underline, blink bool) CacheKey {
	// 颜色代码占用低8位，确保不超出范围
	key := CacheKey(colorCode & ColorMask)

	// 如果是粗体，则设置BoldMask标志
	if bold {
		key |= BoldMask
	}

	// 如果是下划线，则设置UnderlineMask标志
	if underline {
		key |= UnderlineMask
	}

	// 如果是闪烁，则设置BlinkMask标志
	if blink {
		key |= BlinkMask
	}

	return key
}

// Parse 从缓存键解析出原始的样式信息
//
// 返回值:
//   - colorCode: 颜色代码
//   - bold: 是否启用粗体
//   - underline: 是否启用下划线
//   - blink: 是否启用闪烁
func (key CacheKey) Parse() (colorCode int, bold, underline, blink bool) {
	// 解析缓存键
	colorCode = int(key & ColorMask)

	// 解析加粗标志
	bold = (key & BoldMask) != 0

	// 解析下划线标志
	underline = (key & UnderlineMask) != 0

	// 解析闪烁标志
	blink = (key & BlinkMask) != 0
	return
}

// String 返回缓存键的字符串表示，用于调试
//
// 返回值:
//   - string: 格式化的字符串表示
func (key CacheKey) String() string {
	// 解析缓存键
	color, bold, underline, blink := key.Parse()

	// 返回格式化字符串
	return fmt.Sprintf("CacheKey{color:%d, bold:%v, underline:%v, blink:%v}",
		color, bold, underline, blink)
}

// IsValid 检查缓存键是否有效
//
// 返回值:
//   - bool: 缓存键是否有效
func (key CacheKey) IsValid() bool {
	colorCode, _, _, _ := key.Parse()
	// 检查颜色代码是否在有效范围内
	return isValidColorCode(colorCode)
}

// isValidColorCode 检查颜色代码是否有效
//
// 参数:
//   - code: 颜色代码
//
// 返回值:
//   - bool: 颜色代码是否有效
func isValidColorCode(code int) bool {
	// 标准颜色 (30-37)
	if code >= 30 && code <= 37 {
		return true
	}
	// 亮色 (90-97)
	if code >= 90 && code <= 97 {
		return true
	}
	// 其他可能的有效颜色代码可以在这里扩展
	return false
}

// Hash 返回缓存键的哈希值（实际上就是键本身）
//
// 返回值:
//   - uint32: 哈希值
func (key CacheKey) Hash() uint32 {
	return uint32(key)
}

// Equals 比较两个缓存键是否相等
//
// 参数:
//   - other: 另一个缓存键
//
// 返回值:
//   - bool: 是否相等
func (key CacheKey) Equals(other CacheKey) bool {
	return key == other
}

// Package style 提供了样式配置管理功能。
// 该文件实现了 StyleConfig 结构体，用于统一管理颜色、粗体、下划线、闪烁等样式配置，
// 使用原子操作确保线程安全的样式状态管理。
package style

import "sync/atomic"

// StyleConfig 统一管理所有样式配置
type StyleConfig struct {
	color     atomic.Bool // 是否启用颜色
	bold      atomic.Bool // 是否启用加粗
	underline atomic.Bool // 是否启用下划线
	blink     atomic.Bool // 是否启用闪烁
}

// NewStyleConfig 创建一个新的样式配置实例
//
// 返回值：
//   - *StyleConfig：新创建的样式配置实例
func NewStyleConfig() *StyleConfig {
	style := &StyleConfig{}      // 创建样式配置实例
	style.color.Store(true)      // 启用颜色
	style.bold.Store(true)       // 启用加粗
	style.underline.Store(false) // 启用下划线
	style.blink.Store(false)     // 启用闪烁
	return style
}

// Package style 提供了样式配置的工具函数。
// 该文件实现了 StyleConfig 结构体的方法，包括样式状态的设置和获取功能，
// 支持颜色、粗体、下划线、闪烁等样式属性的线程安全操作。
package style

// SetColor 设置颜色启用状态
//
// 参数:
//   - enable: 布尔值，表示是否启用颜色
func (s *StyleConfig) SetColor(enable bool) {
	s.color.Store(enable)
}

// SetBold 设置加粗启用状态
//
// 参数:
//   - enable: 布尔值，表示是否启用加粗
func (s *StyleConfig) SetBold(enable bool) {
	s.bold.Store(enable)
}

// SetUnderline 设置下划线启用状态
//
// 参数:
//   - enable: 布尔值，表示是否启用下划线
func (s *StyleConfig) SetUnderline(enable bool) {
	s.underline.Store(enable)
}

// SetBlink 设置闪烁启用状态
//
// 参数:
//   - enable: 布尔值，表示是否启用闪烁
func (s *StyleConfig) SetBlink(enable bool) {
	s.blink.Store(enable)
}

// GetColor 获取颜色启用状态
//
// 返回值:
//   - bool: 颜色启用状态
func (s *StyleConfig) GetColor() bool {
	return s.color.Load()
}

// GetBold 获取加粗启用状态
//
// 返回值:
//   - bool: 加粗启用状态
func (s *StyleConfig) GetBold() bool {
	return s.bold.Load()
}

// GetUnderline 获取下划线启用状态
//
// 返回值:
//   - bool: 下划线启用状态
func (s *StyleConfig) GetUnderline() bool {
	return s.underline.Load()
}

// GetBlink 获取闪烁启用状态
//
// 返回值:
//   - bool: 闪烁启用状态
func (s *StyleConfig) GetBlink() bool {
	return s.blink.Load()
}

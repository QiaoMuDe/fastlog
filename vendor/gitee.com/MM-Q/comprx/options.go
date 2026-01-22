// Package comprx 提供压缩和解压缩操作的配置选项。
//
// 该文件定义了 Options 结构体和相关的配置方法，用于控制压缩和解压缩操作的行为。
// 支持压缩等级设置、进度条显示、文件过滤、路径验证等功能的配置。
//
// 主要类型：
//   - Options: 压缩/解压配置选项结构体
//
// 主要功能：
//   - 提供默认配置选项
//   - 支持链式配置方法
//   - 提供各种预设配置选项
package comprx

// Options 压缩/解压配置选项
type Options struct {
	CompressionLevel      CompressionLevel // 压缩等级
	OverwriteExisting     bool             // 是否覆盖已存在的文件
	ProgressEnabled       bool             // 是否启用进度显示
	ProgressStyle         ProgressStyle    // 进度条样式
	DisablePathValidation bool             // 是否禁用路径验证
	Filter                FilterOptions    // 过滤选项
}

// DefaultOptions 返回默认配置选项
//
// 返回:
//   - Options: 默认配置选项
//
// 默认配置:
//   - CompressionLevel: 默认压缩等级
//   - OverwriteExisting: false (不覆盖已存在文件)
//   - ProgressEnabled: false (不显示进度)
//   - ProgressStyle: 文本样式
//   - DisablePathValidation: false (启用路径验证)
func DefaultOptions() Options {
	return Options{
		CompressionLevel:      CompressionLevelDefault,
		OverwriteExisting:     false,
		ProgressEnabled:       false,
		ProgressStyle:         ProgressStyleText,
		DisablePathValidation: false,
	}
}

// ProgressOptions 返回带进度显示的配置选项
//
// 参数:
//   - style: 进度条样式
//
// 返回:
//   - Options: 带进度显示的配置选项
func ProgressOptions(style ProgressStyle) Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = style
	return opts
}

// TextProgressOptions 返回文本样式进度条配置选项
//
// 返回:
//   - Options: 文本样式进度条配置选项
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", TextProgressOptions())
func TextProgressOptions() Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = ProgressStyleText
	return opts
}

// UnicodeProgressOptions 返回Unicode样式进度条配置选项
//
// 返回:
//   - Options: Unicode样式进度条配置选项
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", UnicodeProgressOptions())
func UnicodeProgressOptions() Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = ProgressStyleUnicode
	return opts
}

// ASCIIProgressOptions 返回ASCII样式进度条配置选项
//
// 返回:
//   - Options: ASCII样式进度条配置选项
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", ASCIIProgressOptions())
func ASCIIProgressOptions() Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = ProgressStyleASCII
	return opts
}

// DefaultProgressOptions 返回默认样式进度条配置选项
//
// 返回:
//   - Options: 默认样式进度条配置选项
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", DefaultProgressOptions())
func DefaultProgressOptions() Options {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	opts.ProgressStyle = ProgressStyleDefault
	return opts
}

// ForceOptions 返回强制模式配置选项
//
// 返回:
//   - Options: 强制模式配置选项
//
// 配置特点:
//   - OverwriteExisting: true (覆盖已存在文件)
//   - DisablePathValidation: true (禁用路径验证)
//   - ProgressEnabled: false (关闭进度条)
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", ForceOptions())
func ForceOptions() Options {
	opts := DefaultOptions()
	opts.OverwriteExisting = true
	opts.DisablePathValidation = true
	opts.ProgressEnabled = false
	return opts
}

// NoCompressionOptions 返回禁用压缩且启用进度条的配置选项
//
// 返回:
//   - Options: 禁用压缩且启用进度条的配置选项
//
// 配置特点:
//   - CompressionLevel: 无压缩 (存储模式)
//   - ProgressEnabled: true (启用进度条)
//   - ProgressStyle: 文本样式
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", NoCompressionOptions())
func NoCompressionOptions() Options {
	opts := DefaultOptions()
	opts.CompressionLevel = CompressionLevelNone
	opts.ProgressEnabled = true
	opts.ProgressStyle = ProgressStyleText
	return opts
}

// NoCompressionProgressOptions 返回禁用压缩且启用指定样式进度条的配置选项
//
// 参数:
//   - style: 进度条样式
//
// 返回:
//   - Options: 禁用压缩且启用指定样式进度条的配置选项
//
// 配置特点:
//   - CompressionLevel: 无压缩 (存储模式)
//   - ProgressEnabled: true (启用进度条)
//   - ProgressStyle: 指定样式
//
// 使用示例:
//
//	err := PackOptions("output.zip", "input_dir", NoCompressionProgressOptions(ProgressStyleUnicode))
func NoCompressionProgressOptions(style ProgressStyle) Options {
	opts := DefaultOptions()
	opts.CompressionLevel = CompressionLevelNone
	opts.ProgressEnabled = true
	opts.ProgressStyle = style
	return opts
}

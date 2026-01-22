// Package comprx 提供了一个统一的压缩和解压缩库，支持多种压缩格式。
//
// 该包提供了简单易用的 API 来处理 ZIP、TAR、GZIP、BZIP2、ZLIB 和 TGZ 格式的压缩文件。
// 支持进度条显示、文件过滤、并发安全操作等高级功能。
//
// 主要功能：
//   - 压缩和解压缩多种格式的文件
//   - 支持进度条显示
//   - 文件过滤功能
//   - 线程安全操作
//   - 灵活的配置选项
//
// 基本使用示例：
//
//	// 简单压缩
//	err := comprx.Pack("output.zip", "input_dir")
//
//	// 简单解压
//	err := comprx.Unpack("archive.zip", "output_dir")
//
//	// 带进度条的压缩
//	err := comprx.PackProgress("output.zip", "input_dir")
package comprx

import (
	"fmt"

	"gitee.com/MM-Q/comprx/internal/core"
)

// ==============================================
// 简单便捷函数 - 线程安全版本
// ==============================================

// Pack 压缩文件或目录(禁用进度条) - 线程安全
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	err := Pack("output.zip", "input_dir")
func Pack(dst string, src string) error {
	return PackOptions(dst, src, DefaultOptions())
}

// Unpack 解压文件(禁用进度条) - 线程安全
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	err := Unpack("archive.zip", "output_dir")
func Unpack(src string, dst string) error {
	return UnpackOptions(src, dst, DefaultOptions())
}

// PackProgress 压缩文件或目录(启用进度条) - 线程安全
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	err := PackWithProgress("output.zip", "input_dir")
func PackProgress(dst string, src string) error {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	return PackOptions(dst, src, opts)
}

// UnpackProgress 解压文件(启用进度条) - 线程安全
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	err := UnpackWithProgress("archive.zip", "output_dir")
func UnpackProgress(src string, dst string) error {
	opts := DefaultOptions()
	opts.ProgressEnabled = true
	return UnpackOptions(src, dst, opts)
}

// ==============================================
// 指定内容解压便捷函数 - 线程安全版本
// ==============================================

// UnpackFile 解压指定文件名 - 线程安全
//
// 参数:
//   - archivePath: 压缩包路径
//   - fileName: 要解压的文件名
//   - outputDir: 输出目录路径
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	err := UnpackFile("archive.zip", "config.json", "output/")
func UnpackFile(archivePath string, fileName string, outputDir string) error {
	// 参数验证
	if archivePath == "" {
		return fmt.Errorf("压缩包路径不能为空")
	}
	if fileName == "" {
		return fmt.Errorf("文件名不能为空")
	}
	if outputDir == "" {
		return fmt.Errorf("输出目录不能为空")
	}

	// 创建过滤器选项
	opts := DefaultOptions()
	opts.Filter = FilterOptions{
		Include: []string{fileName},
		Exclude: []string{},
		MaxSize: 0,
		MinSize: 0,
	}

	return UnpackOptions(archivePath, outputDir, opts)
}

// UnpackDir 解压指定目录 - 线程安全
//
// 参数:
//   - archivePath: 压缩包路径
//   - dirName: 要解压的目录名
//   - outputDir: 输出目录路径
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	err := UnpackDir("archive.zip", "src", "output/")
func UnpackDir(archivePath string, dirName string, outputDir string) error {
	// 参数验证
	if archivePath == "" {
		return fmt.Errorf("压缩包路径不能为空")
	}
	if dirName == "" {
		return fmt.Errorf("目录名不能为空")
	}
	if outputDir == "" {
		return fmt.Errorf("输出目录不能为空")
	}

	// 创建过滤器选项
	opts := DefaultOptions()
	opts.Filter = FilterOptions{
		Include: []string{
			dirName,        // 匹配目录本身
			dirName + "/*", // 匹配目录下所有内容
		},
		Exclude: []string{},
		MaxSize: 0,
		MinSize: 0,
	}

	return UnpackOptions(archivePath, outputDir, opts)
}

// UnpackMatch 解压匹配关键字的文件 - 线程安全
//
// 参数:
//   - archivePath: 压缩包路径
//   - keyword: 匹配关键字
//   - outputDir: 输出目录路径
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	err := UnpackMatch("archive.zip", "test", "output/")
func UnpackMatch(archivePath string, keyword string, outputDir string) error {
	// 参数验证
	if archivePath == "" {
		return fmt.Errorf("压缩包路径不能为空")
	}
	if keyword == "" {
		return fmt.Errorf("关键字不能为空")
	}
	if outputDir == "" {
		return fmt.Errorf("输出目录不能为空")
	}

	// 创建过滤器选项
	opts := DefaultOptions()
	opts.Filter = FilterOptions{
		Include: []string{"*" + keyword + "*"},
		Exclude: []string{},
		MaxSize: 0,
		MinSize: 0,
	}

	return UnpackOptions(archivePath, outputDir, opts)
}

// ==============================================
// 配置化便捷函数 - 线程安全版本
// ==============================================

// PackOptions 使用指定配置压缩文件或目录 - 线程安全
//
// 参数:
//   - dst: 目标文件路径
//   - src: 源文件路径
//   - opts: 配置选项
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	opts := Options{
//	    CompressionLevel: config.CompressionLevelBest,
//	    OverwriteExisting: true,
//	    ProgressEnabled: true,
//	    ProgressStyle: ProgressStyleUnicode,
//	}
//	err := PackOptions("output.zip", "input_dir", opts)
func PackOptions(dst string, src string, opts Options) error {
	comprx := core.New()

	// 验证压缩等级
	if !opts.CompressionLevel.IsValid() {
		return fmt.Errorf("无效的压缩等级: %s，有效范围: -2 到 9", opts.CompressionLevel.String())
	}
	comprx.Config.CompressionLevel = opts.CompressionLevel
	comprx.Config.OverwriteExisting = opts.OverwriteExisting
	comprx.Config.Progress.Enabled = opts.ProgressEnabled

	// 验证进度条样式
	if !opts.ProgressStyle.IsValid() {
		return fmt.Errorf("invalid progress style: %v", opts.ProgressStyle)
	}
	comprx.Config.Progress.BarStyle = opts.ProgressStyle
	comprx.Config.DisablePathValidation = opts.DisablePathValidation

	// 验证过滤器选项
	if err := opts.Filter.Validate(); err != nil {
		return err
	}
	comprx.Config.Filter = &FilterOptions{
		Include: opts.Filter.Include,
		Exclude: opts.Filter.Exclude,
		MaxSize: opts.Filter.MaxSize,
		MinSize: opts.Filter.MinSize,
	}

	return comprx.Pack(dst, src)
}

// UnpackOptions 使用指定配置解压文件 - 线程安全
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//   - opts: 配置选项
//
// 返回:
//   - error: 错误信息
//
// 使用示例:
//
//	opts := Options{
//	    OverwriteExisting: true,
//	    ProgressEnabled: true,
//	    ProgressStyle: ProgressStyleASCII,
//	}
//	err := UnpackOptions("archive.zip", "output_dir", opts)
func UnpackOptions(src string, dst string, opts Options) error {
	comprx := core.New()

	// 设置配置（解压时不需要验证压缩等级）
	comprx.Config.OverwriteExisting = opts.OverwriteExisting
	comprx.Config.Progress.Enabled = opts.ProgressEnabled

	// 验证进度条样式
	if !opts.ProgressStyle.IsValid() {
		return fmt.Errorf("invalid progress style: %v", opts.ProgressStyle)
	}
	comprx.Config.Progress.BarStyle = opts.ProgressStyle
	comprx.Config.DisablePathValidation = opts.DisablePathValidation

	// 验证并设置过滤器
	if err := opts.Filter.Validate(); err != nil {
		return err
	}
	comprx.Config.Filter = &FilterOptions{
		Include: opts.Filter.Include,
		Exclude: opts.Filter.Exclude,
		MaxSize: opts.Filter.MaxSize,
		MinSize: opts.Filter.MinSize,
	}

	return comprx.Unpack(src, dst)
}

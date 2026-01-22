// Package progress 提供压缩和解压缩操作的进度显示功能。
//
// 该包实现了多种样式的进度显示，包括文本模式和进度条模式。
// 支持压缩和解压缩过程中的实时进度反馈，提供了统一的进度管理接口。
//
// 主要类型：
//   - Progress: 进度显示器结构体
//
// 主要功能：
//   - 支持多种进度条样式（文本、ASCII、Unicode、默认）
//   - 提供压缩和解压缩的进度显示
//   - 支持文件扫描进度显示
//   - 提供带进度的数据复制功能
//   - 自动管理进度条生命周期
//
// 支持的进度样式：
//   - 文本模式：显示操作文本信息
//   - ASCII模式：使用ASCII字符的进度条
//   - Unicode模式：使用Unicode字符的精美进度条
//   - 默认模式：使用库默认样式的进度条
//
// 使用示例：
//
//	// 创建进度显示器
//	progress := progress.New()
//
//	// 开始进度显示
//	err := progress.Start(totalSize, "archive.zip", "正在解压...")
//
//	// 带进度的数据复制
//	written, err := progress.CopyBuffer(dst, src, buffer)
//
//	// 关闭进度显示
//	err := progress.Close()
package progress

import (
	"fmt"
	"io"
	"path/filepath"

	"gitee.com/MM-Q/comprx/types"
	"github.com/schollz/progressbar/v3"
)

// 操作标签常量 - 确保冒号对齐
const (
	labelArchive     = "Archive:    " // 表示操作压缩包
	labelInflating   = "  inflating:" // 表示操作解压文件
	labelCreating    = "   creating:" // 表示操作创建目录
	labelAdding      = "     adding:" // 表示操作添加文件
	labelStoring     = "    storing:" // 表示操作存储目录
	labelCompressing = "compressing:" // 表示操作压缩文件
)

// Progress 控制台进度显示器
type Progress struct {
	Enabled  bool                // 是否启用进度显示
	BarStyle types.ProgressStyle // 进度条样式

	// 当前进度条相关字段 //
	totalSize   int64                    // 总大小
	currentBar  *progressbar.ProgressBar // 当前进度条实例
	isActive    bool                     // 是否有活跃的进度操作
	description string                   // 操作描述
}

// New 创建进度显示器
//
// 返回:
//   - *Progress: 简单进度显示器
func New() *Progress {
	return &Progress{
		Enabled:  true,                    // 是否启用进度显示
		BarStyle: types.ProgressStyleText, // 进度条样式
	}
}

// Start 开始进度显示，创建进度条
//
// 参数:
//   - totalSize: 总数据大小
//   - archivePath: 压缩包路径（用于文本模式显示）
//   - description: 操作描述（用于进度条模式显示，如"正在解压 xxx.zip..."）
//
// 返回:
//   - error: 初始化错误
func (s *Progress) Start(totalSize int64, archivePath, description string) error {
	// 未启用 或 已存在活跃进度条 则返回
	if !s.Enabled || s.isActive {
		return nil
	}

	// 文本模式 或 总大小小于等于0 则显示Archive信息(文本模式)
	if s.BarStyle == types.ProgressStyleText || totalSize <= 0 {
		ext := filepath.Ext(archivePath)
		if ext == "" {
			s.Archive(archivePath)
			return nil
		}

		// 如果是bz2、gz、bzip2格式的压缩文件，则显示压缩文件信息
		if ext == types.CompressTypeBz2.String() || ext == types.CompressTypeGz.String() || ext == types.CompressTypeBzip2.String() {
			s.Compressing(archivePath)
			return nil
		}

		// 如果是其他格式的压缩包，则显示压缩包信息
		s.Archive(archivePath)
		return nil
	}

	// 仅在进度条模式下初始化进度条
	s.totalSize = totalSize     // 总数据大小
	s.description = description // 操作描述
	s.isActive = true           // 有活跃的进度操作

	// 创建底层进度条
	if bar := s.newProgressBar(totalSize, description); bar != nil {
		s.currentBar = bar // 保存进度条实例
		return nil
	}

	return nil
}

// StartScan 开始扫描进度显示，创建进度条
//
// 参数:
//   - description: 操作描述（如"正在计算文件大小..."）
//
// 返回:
//   - bar: 进度条实例
func (s *Progress) StartScan(description string) *progressbar.ProgressBar {
	return s.newProgressBar(-1, description)
}

// CloseBar 通用进度条关闭方法
//
// 参数:
//   - bar: 进度条实例
//
// 返回:
//   - error: 清理错误
func (s *Progress) CloseBar(bar *progressbar.ProgressBar) error {
	if bar == nil {
		return nil
	}

	// 完成进度条
	if err := bar.Finish(); err != nil {
		return err
	}

	// 关闭进度条
	if err := bar.Close(); err != nil {
		return err
	}

	return nil
}

// CopyBuffer 带进度显示的数据复制
//
// 参数:
//   - dst: 目标写入器
//   - src: 源读取器
//   - buf: 缓冲区
//
// 返回:
//   - written: 写入的字节数
//   - err: 错误信息
//
// 使用示例:
//
//	written, err := cfg.Progress.CopyBuffer(fileWriter, zipReader, buffer, "file.txt")
func (s *Progress) CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	if dst == nil || src == nil {
		return 0, fmt.Errorf("dst 或 src 不能为 nil")
	}

	// 进度条未启用 或 未开始 直接使用标准库copybuffer复制
	if !s.Enabled || !s.isActive {
		return io.CopyBuffer(dst, src, buf)
	}

	// 文字模式也使用标准库copybuffer复制
	if s.BarStyle == types.ProgressStyleText {
		return io.CopyBuffer(dst, src, buf)
	}

	// 如果进度条写入器未空则直接使用标准库copybuffer复制
	if s.currentBar == nil {
		return io.CopyBuffer(dst, src, buf)
	}

	// 使用多写入器同时写入文件和更新进度条
	multiWriter := io.MultiWriter(dst, s.currentBar)
	written, err = io.CopyBuffer(multiWriter, src, buf)

	return written, err
}

// Close 关闭进度显示，清理资源
//
// 返回:
//   - error: 清理错误
//
// 使用示例:
//   - err := cfg.Progress.Close()
func (s *Progress) Close() error {
	// 没有使用的进度条 或 未启用 或 进度条实例为空 则直接返回
	if !s.isActive || !s.Enabled || s.currentBar == nil {
		return nil
	}

	// 完成进度条
	if err := s.currentBar.Finish(); err != nil {
		return err
	}

	// 关闭进度条
	if err := s.currentBar.Close(); err != nil {
		return err
	}

	// 重置进度条实例
	s.currentBar = nil

	// 重置状态
	s.totalSize = 0
	s.isActive = false
	s.description = ""

	return nil
}

// toTheme 将进度条样式转换为主题
//
// 参数:
//   - style: 进度条样式
//
// 返回:
//   - progressbar.Theme: 主题
func toTheme(style types.ProgressStyle) progressbar.Theme {
	switch style {
	// 使用 unicode 样式
	case types.ProgressStyleUnicode:
		return progressbar.ThemeUnicode

	/// 使用 ASCII 样式
	case types.ProgressStyleASCII:
		return progressbar.ThemeASCII

	/// 使用默认样式
	case types.ProgressStyleDefault:
		return progressbar.ThemeDefault

	// 如果指定文本样式则转为外部库的ASCII样式
	case types.ProgressStyleText:
		return progressbar.ThemeASCII

	default:
		// 默认使用ASCII样式，兼容性最好
		return progressbar.ThemeASCII
	}
}

// newProgressBar 创建一个进度条
//
// 参数:
//   - total: 进度条总大小
//   - description: 进度条描述信息
//
// 返回:
//   - *progressbar.ProgressBar: 进度条指针
//
// 进度条样式:
//   - types.ProgressStyleUnicode: Unicode样式进度条 - 使用Unicode字符绘制精美进度条
//   - types.ProgressStyleASCII: ASCII样式进度条 - 使用基础ASCII字符绘制兼容性最好的进度条
//   - types.ProgressStyleDefault: 默认样式进度条 - 使用外部库默认主题
func (s *Progress) newProgressBar(total int64, description string) *progressbar.ProgressBar {
	// 转换进度条样式为主题
	theme := toTheme(s.BarStyle)

	return progressbar.NewOptions64(
		total,                             // 进度条总大小
		progressbar.OptionClearOnFinish(), // 完成后清除进度条
		progressbar.OptionSetDescription(description), // 进度条描述信息
		progressbar.OptionSetElapsedTime(true),        // 显示已用时间
		progressbar.OptionSetPredictTime(true),        // 显示预计剩余时间
		progressbar.OptionSetRenderBlankState(true),   // 在进度条完成之前显示空白状态
		progressbar.OptionShowBytes(true),             // 显示进度条传输的字节
		progressbar.OptionShowCount(),                 // 显示当前进度的总和
		//progressbar.OptionShowElapsedTimeOnFinish(),        // 完成后显示已用时间
		progressbar.OptionSetTheme(theme), // 设置进度条主题样式
		progressbar.OptionFullWidth(),     // 设置进度条全宽显示
	)
}

// IsEnabled 检查是否启用
//
// 返回:
//   - bool: 是否启用
func (s *Progress) IsEnabled() bool {
	return s.Enabled
}

// Archive 显示压缩文件信息
//
// 参数:
//   - archivePath: 压缩文件路径
func (s *Progress) Archive(archivePath string) {
	// 如果不启用进度显示, 则直接返回
	// 如果进度条样式不是文本样式, 则直接返回
	if !s.Enabled || s.BarStyle != types.ProgressStyleText {
		return
	}
	fmt.Printf("%s %s\n", labelArchive, filepath.Base(archivePath))
}

// Compressing 显示压缩文件信息
//
// 参数:
//   - filePath: 文件路径
func (s *Progress) Compressing(filePath string) {
	if !s.Enabled || s.BarStyle != types.ProgressStyleText {
		return
	}
	fmt.Printf("%s %s\n", labelCompressing, filepath.Base(filePath))
}

// ======================================================
// 解压进度
// ======================================================

// Inflating 显示解压文件
//
// 参数:
//   - filePath: 文件路径
func (s *Progress) Inflating(filePath string) {
	if !s.Enabled || s.BarStyle != types.ProgressStyleText {
		return
	}
	fmt.Printf("%s %s\n", labelInflating, filePath)
}

// Creating 显示创建目录
//
// 参数:
//   - dirPath: 目录路径
func (s *Progress) Creating(dirPath string) {
	if !s.Enabled || s.BarStyle != types.ProgressStyleText {
		return
	}
	fmt.Printf("%s %s\n", labelCreating, dirPath)
}

// ======================================================
// 压缩进度
// ======================================================

// Adding 显示添加文件
//
// 参数:
//   - filePath: 文件路径
func (s *Progress) Adding(filePath string) {
	if !s.Enabled || s.BarStyle != types.ProgressStyleText {
		return
	}
	fmt.Printf("%s %s\n", labelAdding, filePath)
}

// Storing 显示存储目录
//
// 参数:
//   - dirPath: 目录路径
func (s *Progress) Storing(dirPath string) {
	if !s.Enabled || s.BarStyle != types.ProgressStyleText {
		return
	}
	fmt.Printf("%s %s\n", labelStoring, dirPath)
}

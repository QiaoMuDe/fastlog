// Package progress 提供源文件大小计算和进度显示的实用工具函数。
//
// 该文件实现了在压缩操作前计算源文件总大小的功能，用于初始化进度条的总进度。
// 支持单个文件和目录的递归大小计算，并在计算过程中显示扫描进度。
//
// 主要功能：
//   - 计算源路径中所有普通文件的总大小
//   - 在计算过程中显示扫描进度
//   - 支持文件过滤器，跳过不需要的文件
//   - 自动区分文件和目录处理
//   - 只在进度条模式下执行计算
//
// 性能优化：
//   - 文本模式下跳过大小计算，直接返回0
//   - 支持过滤器提前跳过不需要的文件和目录
//   - 实时更新扫描进度条
//   - 错误容忍，遇到错误继续遍历
//
// 使用示例：
//
//	// 计算源文件总大小并显示进度
//	totalSize := progress.CalculateSourceTotalSizeWithProgress(
//	    srcPath,
//	    progressObj,
//	    "正在分析内容...",
//	    filterOptions
//	)
package progress

import (
	"io/fs"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/types"
)

// CalculateSourceTotalSizeWithProgress 计算源路径中所有普通文件的总大小并显示进度
//
// 参数:
//   - srcPath: 源路径（文件或目录）
//   - progress: 进度显示对象
//   - scanMessage: 扫描时显示的消息，如 "正在分析内容..."
//   - filter: 文件过滤器，用于跳过不需要的文件
//
// 返回值:
//   - int64: 普通文件的总大小（字节）
//
// 功能:
//   - 只在进度条模式下计算总大小，文本模式返回 0
//   - 显示扫描进度条并实时更新
//   - 支持单个文件和目录的大小计算
//   - 只计算普通文件，忽略目录、符号链接等特殊文件
//   - 应用过滤器跳过不需要处理的文件
func CalculateSourceTotalSizeWithProgress(srcPath string, progress *Progress, scanMessage string, filter *types.FilterOptions) int64 {
	// 只在进度条模式下计算总大小
	if !progress.Enabled || progress.BarStyle == types.ProgressStyleText {
		return 0
	}

	// 开始扫描进度显示
	bar := progress.StartScan(scanMessage)
	defer func() {
		_ = progress.CloseBar(bar)
	}()

	var totalSize int64

	// 检查是文件还是目录
	info, err := os.Stat(srcPath)
	if err != nil {
		return 0
	}

	// 单个文件处理
	if info.Mode().IsRegular() {
		// 检查是否应该跳过
		if filter != nil && filter.ShouldSkipByParams(srcPath, info.Size(), false) {
			return 0 // 文件被过滤器跳过
		}
		totalSize = info.Size()
		_ = bar.Add64(totalSize)
		return totalSize
	}

	// 目录处理
	if info.IsDir() {
		// 遍历目录下所有文件处理
		_ = filepath.WalkDir(srcPath, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return nil // 忽略错误，继续遍历
			}

			// 获取文件信息用于过滤检查
			fileInfo, err := entry.Info()
			if err != nil {
				return nil // 忽略错误，继续遍历
			}

			// 应用过滤器检查
			if filter != nil && filter.ShouldSkipByParams(path, fileInfo.Size(), fileInfo.IsDir()) {
				if fileInfo.IsDir() {
					return filepath.SkipDir // 跳过整个目录
				}
				return nil // 跳过文件
			}

			// 只计算普通文件的大小
			if entry.Type().IsRegular() {
				fileSize := fileInfo.Size()
				totalSize += fileSize   // 累加文件大小
				_ = bar.Add64(fileSize) // 实时更新进度条
			}

			return nil
		})
		return totalSize
	}

	// 其他类型文件（符号链接、设备文件等）不计算大小
	return 0
}

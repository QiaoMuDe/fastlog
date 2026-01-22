// Package comprx 提供压缩包内容列表和信息查看功能。
//
// 该文件提供了查看压缩包内容的各种方法，包括列出文件信息、打印压缩包信息等。
// 支持多种压缩格式，提供简洁和详细两种显示样式，支持文件过滤和数量限制。
//
// 主要功能：
//   - 列出压缩包内的文件信息
//   - 打印压缩包基本信息
//   - 支持文件名模式匹配
//   - 支持限制显示文件数量
//   - 提供简洁和详细两种显示样式
//
// 使用示例：
//
//	// 列出压缩包内所有文件
//	info, err := comprx.List("archive.zip")
//
//	// 打印压缩包信息（简洁样式）
//	err := comprx.PrintLs("archive.zip")
//
//	// 打印匹配模式的文件（详细样式）
//	err := comprx.PrintLlMatch("archive.zip", "*.go")
package comprx

import (
	"gitee.com/MM-Q/comprx/internal/core"
	"gitee.com/MM-Q/comprx/internal/utils"
)

// ==============================================
// 压缩包信息获取方法
// ==============================================

// List 列出压缩包的所有文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - *ArchiveInfo: 压缩包信息
//   - error: 错误信息
func List(archivePath string) (*ArchiveInfo, error) {
	return core.List(archivePath)
}

// ListLimit 列出指定数量的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制返回的文件数量
//
// 返回:
//   - *ArchiveInfo: 压缩包信息
//   - error: 错误信息
func ListLimit(archivePath string, limit int) (*ArchiveInfo, error) {
	return core.ListLimit(archivePath, limit)
}

// ListMatch 列出匹配指定模式的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//
// 返回:
//   - *ArchiveInfo: 压缩包信息
//   - error: 错误信息
func ListMatch(archivePath string, pattern string) (*ArchiveInfo, error) {
	return core.ListMatch(archivePath, pattern)
}

// ==============================================
// 打印压缩包本身信息
// ==============================================

// PrintArchiveInfo 打印压缩包本身的基本信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - error: 错误信息
func PrintArchiveInfo(archivePath string) error {
	archiveInfo, err := core.List(archivePath)
	if err != nil {
		return err
	}

	utils.PrintArchiveSummary(archiveInfo)
	return nil
}

// ==============================================
// 打印压缩包内文件信息
// ==============================================

// PrintFiles 打印压缩包内所有文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintFiles(archivePath string, detailed bool) error {
	archiveInfo, err := core.List(archivePath)
	if err != nil {
		return err
	}

	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// PrintFilesLimit 打印压缩包内指定数量的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintFilesLimit(archivePath string, limit int, detailed bool) error {
	archiveInfo, err := core.ListLimit(archivePath, limit)
	if err != nil {
		return err
	}

	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// PrintFilesMatch 打印压缩包内匹配指定模式的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintFilesMatch(archivePath string, pattern string, detailed bool) error {
	archiveInfo, err := core.ListMatch(archivePath, pattern)
	if err != nil {
		return err
	}

	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// ==============================================
// 便捷函数 - 简洁样式
// ==============================================

// PrintLs 打印压缩包内所有文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - error: 错误信息
func PrintLs(archivePath string) error {
	return PrintFiles(archivePath, false)
}

// PrintLsLimit 打印压缩包内指定数量的文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//
// 返回:
//   - error: 错误信息
func PrintLsLimit(archivePath string, limit int) error {
	return PrintFilesLimit(archivePath, limit, false)
}

// PrintLsMatch 打印压缩包内匹配指定模式的文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//
// 返回:
//   - error: 错误信息
func PrintLsMatch(archivePath string, pattern string) error {
	return PrintFilesMatch(archivePath, pattern, false)
}

// ==============================================
// 便捷函数 - 详细样式
// ==============================================

// PrintLl 打印压缩包内所有文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - error: 错误信息
func PrintLl(archivePath string) error {
	return PrintFiles(archivePath, true)
}

// PrintLlLimit 打印压缩包内指定数量的文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//
// 返回:
//   - error: 错误信息
func PrintLlLimit(archivePath string, limit int) error {
	return PrintFilesLimit(archivePath, limit, true)
}

// PrintLlMatch 打印压缩包内匹配指定模式的文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//
// 返回:
//   - error: 错误信息
func PrintLlMatch(archivePath string, pattern string) error {
	return PrintFilesMatch(archivePath, pattern, true)
}

// ==============================================
// 打印压缩包信息+文件信息
// ==============================================

// PrintArchiveAndFiles 打印压缩包信息和所有文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintArchiveAndFiles(archivePath string, detailed bool) error {
	archiveInfo, err := core.List(archivePath)
	if err != nil {
		return err
	}

	utils.PrintArchiveSummary(archiveInfo)
	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// PrintArchiveAndFilesLimit 打印压缩包信息和指定数量的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintArchiveAndFilesLimit(archivePath string, limit int, detailed bool) error {
	archiveInfo, err := core.ListLimit(archivePath, limit)
	if err != nil {
		return err
	}

	utils.PrintArchiveSummary(archiveInfo)
	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// PrintArchiveAndFilesMatch 打印压缩包信息和匹配指定模式的文件信息
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//   - detailed: true=详细样式, false=简洁样式(默认)
//
// 返回:
//   - error: 错误信息
func PrintArchiveAndFilesMatch(archivePath string, pattern string, detailed bool) error {
	archiveInfo, err := core.ListMatch(archivePath, pattern)
	if err != nil {
		return err
	}

	utils.PrintArchiveSummary(archiveInfo)
	utils.PrintFileList(archiveInfo.Files, detailed)
	return nil
}

// ==============================================
// 便捷函数 - 压缩包信息+文件信息（简洁样式）
// ==============================================

// PrintInfo 打印压缩包信息和所有文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - error: 错误信息
func PrintInfo(archivePath string) error {
	return PrintArchiveAndFiles(archivePath, false)
}

// PrintInfoLimit 打印压缩包信息和指定数量的文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//
// 返回:
//   - error: 错误信息
func PrintInfoLimit(archivePath string, limit int) error {
	return PrintArchiveAndFilesLimit(archivePath, limit, false)
}

// PrintInfoMatch 打印压缩包信息和匹配指定模式的文件信息（简洁样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//
// 返回:
//   - error: 错误信息
func PrintInfoMatch(archivePath string, pattern string) error {
	return PrintArchiveAndFilesMatch(archivePath, pattern, false)
}

// ==============================================
// 便捷函数 - 压缩包信息+文件信息（详细样式）
// ==============================================

// PrintInfoDetailed 打印压缩包信息和所有文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//
// 返回:
//   - error: 错误信息
func PrintInfoDetailed(archivePath string) error {
	return PrintArchiveAndFiles(archivePath, true)
}

// PrintInfoDetailedLimit 打印压缩包信息和指定数量的文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - limit: 限制打印的文件数量
//
// 返回:
//   - error: 错误信息
func PrintInfoDetailedLimit(archivePath string, limit int) error {
	return PrintArchiveAndFilesLimit(archivePath, limit, true)
}

// PrintInfoDetailedMatch 打印压缩包信息和匹配指定模式的文件信息（详细样式）
//
// 参数:
//   - archivePath: 压缩包文件路径
//   - pattern: 文件名匹配模式 (支持通配符 * 和 ?)
//
// 返回:
//   - error: 错误信息
func PrintInfoDetailedMatch(archivePath string, pattern string) error {
	return PrintArchiveAndFilesMatch(archivePath, pattern, true)
}

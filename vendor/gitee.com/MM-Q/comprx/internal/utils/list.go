// Package utils 提供压缩包文件列表处理和格式化显示的实用工具函数。
//
// 该文件实现了压缩包内容显示相关的工具函数，包括文件大小格式化、文件权限显示、
// 文件列表打印、模式匹配过滤等功能。这些函数被压缩包列表功能广泛使用。
//
// 主要功能：
//   - 文件大小格式化显示（B、KB、MB、GB等）
//   - 文件权限格式化显示
//   - 文件名模式匹配（支持通配符）
//   - 压缩包摘要信息打印
//   - 文件列表格式化打印
//   - 文件列表过滤和限制
//
// 显示格式：
//   - 简洁模式：仅显示文件名
//   - 详细模式：显示权限、大小、时间等完整信息
//   - 支持符号链接目标显示
//   - 自动计算压缩率
//
// 使用示例：
//
//	// 格式化文件大小
//	sizeStr := utils.FormatFileSize(1024*1024) // "1.0 MB"
//
//	// 打印压缩包摘要
//	utils.PrintArchiveSummary(archiveInfo)
//
//	// 打印文件列表（详细模式）
//	utils.PrintFileList(files, true)
//
//	// 模式匹配过滤
//	filtered := utils.FilterFilesByPattern(files, "*.go")
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/types"
)

// 文件大小格式化相关常量
const (
	// SizeUnit 文件大小计算单位 (1024字节)
	SizeUnit = 1024
	// SizeUnitStr 文件大小单位字符串 (KB, MB, GB, TB, PB, EB)
	SizeUnitStr = "KMGTPE"
)

// 切片预分配相关常量
const (
	// DefaultFileCapacity 默认文件列表初始容量
	// 适用于 TAR/TGZ 等无法预先知道文件数量的格式
	DefaultFileCapacity = 256
)

// 文件处理相关常量
const (
	// DefaultBufferSize 默认缓冲区大小 (32KB)
	// 用于读取压缩文件内容时的缓冲区
	DefaultBufferSize = 32 * 1024

	// DefaultFileMode 默认文件权限 (0644)
	// 用于不保存文件权限的压缩格式 (如 GZIP, BZ2)
	DefaultFileMode = 0644
)

// 文件扩展名相关常量
const (
	// DecompressedSuffix 解压缩文件的默认后缀
	DecompressedSuffix = ".decompressed"
)

// FormatFileSize 格式化文件大小显示
//
// 参数:
//   - size: 文件大小
//
// 返回:
//   - string: 格式化后的文件大小字符串
func FormatFileSize(size int64) string {
	if size < SizeUnit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(SizeUnit), 0
	for n := size / SizeUnit; n >= SizeUnit; n /= SizeUnit {
		div *= SizeUnit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), SizeUnitStr[exp])
}

// FormatFileMode 格式化文件权限显示
//
// 参数:
//   - mode: 文件权限
//
// 返回:
//   - string: 格式化后的文件权限字符串
func FormatFileMode(mode os.FileMode) string {
	return mode.String()
}

// MatchPattern 文件名模式匹配
// 支持简单的通配符匹配: * 和 ?
//
// 参数:
//   - name: 文件名
//   - pattern: 模式字符串
//
// 返回:
//   - bool: 是否匹配成功
func MatchPattern(name, pattern string) bool {
	if pattern == "" {
		return true
	}

	// 使用 filepath.Match 进行通配符匹配
	matched, err := filepath.Match(pattern, name)
	if err != nil {
		// 如果模式无效，尝试简单的字符串包含匹配
		return strings.Contains(strings.ToLower(name), strings.ToLower(pattern))
	}

	// 如果没有匹配到文件名，尝试匹配路径的任意部分
	if !matched {
		parts := strings.Split(name, "/")
		for _, part := range parts {
			if matched, _ := filepath.Match(pattern, part); matched {
				return true
			}
		}
		// 最后尝试字符串包含匹配
		return strings.Contains(strings.ToLower(name), strings.ToLower(pattern))
	}

	return matched
}

// PrintFileInfo 格式化打印单个文件信息
//
// 参数:
//   - info: 文件信息
//   - showDetails: 是否显示详细信息
func PrintFileInfo(info types.FileInfo, showDetails bool) {
	if showDetails {
		// 详细模式：显示权限、大小、时间等
		modeStr := FormatFileMode(info.Mode)
		sizeStr := FormatFileSize(info.Size)
		timeStr := info.ModTime.Format("2006-01-02 15:04:05")

		if info.IsSymlink {
			fmt.Printf("%s %8s %s %s -> %s\n", modeStr, sizeStr, timeStr, info.Name, info.LinkTarget)
		} else {
			fmt.Printf("%s %8s %s %s\n", modeStr, sizeStr, timeStr, info.Name)
		}
	} else {
		// 简单模式：只显示文件名
		if info.IsSymlink {
			fmt.Printf("%s -> %s\n", info.Name, info.LinkTarget)
		} else {
			fmt.Printf("%s\n", info.Name)
		}
	}
}

// PrintArchiveSummary 打印压缩包摘要信息
//
// 参数:
//   - archiveInfo: 压缩包信息
func PrintArchiveSummary(archiveInfo *types.ArchiveInfo) {
	fmt.Println(strings.Repeat("-", 50))                            // 分隔线
	fmt.Printf("压缩包类型: %s\n", archiveInfo.Type)                     // 压缩包类型
	fmt.Printf("文件总数: %d\n", archiveInfo.TotalFiles)                // 文件总数
	fmt.Printf("原始大小: %s\n", FormatFileSize(archiveInfo.TotalSize)) // 原始大小

	// 如果有压缩大小，则显示压缩大小和压缩率
	if archiveInfo.CompressedSize > 0 {
		fmt.Printf("压缩大小: %s\n", FormatFileSize(archiveInfo.CompressedSize)) // 压缩大小

		// 计算压缩率，避免负值
		if archiveInfo.CompressedSize >= archiveInfo.TotalSize {
			// 当压缩大小大于等于原始大小时，显示0.0%（适用于TAR等归档格式）
			fmt.Printf("压缩率: 0.0%%\n")
		} else {
			// 正常压缩率计算
			ratio := (1.0 - float64(archiveInfo.CompressedSize)/float64(archiveInfo.TotalSize)) * 100
			fmt.Printf("压缩率: %.1f%%\n", ratio)
		}
	}
	fmt.Println(strings.Repeat("-", 50)) // 分隔线
}

// PrintFileList 打印文件列表
//
// 参数:
//   - files: 文件列表
//   - showDetails: 是否显示详细信息
func PrintFileList(files []types.FileInfo, showDetails bool) {
	// 遍历文件列表并打印
	for _, file := range files {
		PrintFileInfo(file, showDetails)
	}
}

// FilterFilesByPattern 根据模式过滤文件列表
//
// 参数:
//   - files: 文件列表
//   - pattern: 模式字符串
//
// 返回:
//   - []types.FileInfo: 过滤后的文件列表
func FilterFilesByPattern(files []types.FileInfo, pattern string) []types.FileInfo {
	if pattern == "" {
		return files
	}

	var filtered []types.FileInfo
	for _, file := range files {
		if MatchPattern(file.Name, pattern) {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

// LimitFiles 限制文件列表数量
//
// 参数:
//   - files: 文件列表
//   - limit: 限制数量
//
// 返回:
//   - []types.FileInfo: 限制后的文件列表
func LimitFiles(files []types.FileInfo, limit int) []types.FileInfo {
	if limit <= 0 || limit >= len(files) {
		return files
	}
	return files[:limit]
}

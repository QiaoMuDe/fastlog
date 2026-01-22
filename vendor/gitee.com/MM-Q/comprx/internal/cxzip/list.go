// Package cxzip 提供 ZIP 格式的压缩包内容列表功能实现。
//
// 该包实现了 ZIP 格式压缩包的文件信息获取功能，包括基本列表、限制数量列表和模式匹配列表。
// ZIP 格式支持多种文件类型，包括普通文件、目录和符号链接，并提供完整的压缩信息。
//
// 主要功能：
//   - ZIP 压缩包完整文件列表获取
//   - 限制数量的文件列表获取
//   - 模式匹配的文件列表过滤
//   - 多种文件类型支持（普通文件、目录、符号链接）
//   - 完整的压缩信息统计
//
// 文件类型支持：
//   - 普通文件：完整的文件信息和压缩信息
//   - 目录：目录标识和权限信息
//   - 符号链接：链接目标路径读取和保存
//
// 压缩信息：
//   - 原始文件大小
//   - 压缩后大小
//   - 压缩率计算
//   - 文件修改时间
//   - 文件权限模式
//
// 性能优化：
//   - 限制模式下的容量预分配
//   - 符号链接目标的安全读取
//   - 错误容忍的链接目标处理
//
// 使用示例：
//
//	// 获取 ZIP 文件完整列表
//	info, err := cxzip.ListZip("archive.zip")
//
//	// 获取前 10 个文件信息
//	info, err := cxzip.ListZipLimit("archive.zip", 10)
//
//	// 获取匹配 *.go 模式的文件
//	info, err := cxzip.ListZipMatch("archive.zip", "*.go")
package cxzip

import (
	"archive/zip"
	"fmt"
	"os"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
)

// ListZip 获取ZIP压缩包的所有文件信息
func ListZip(archivePath string) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "ZIP文件路径")
	if err != nil {
		return nil, err
	}

	// 打开ZIP文件
	reader, err := zip.OpenReader(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开ZIP文件失败: %w", err)
	}
	defer func() { _ = reader.Close() }()

	// 获取压缩包文件信息
	stat, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("获取ZIP文件信息失败: %w", err)
	}

	// 根据文件名检测压缩格式类型
	compressType, err := utils.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 创建 ArchiveInfo 结构体
	archiveInfo := &types.ArchiveInfo{
		Type:           compressType,
		TotalFiles:     len(reader.File),
		CompressedSize: stat.Size(),
		Files:          make([]types.FileInfo, 0, len(reader.File)),
	}

	// 遍历ZIP文件中的每个条目
	for _, file := range reader.File {
		fileInfo := types.FileInfo{
			Name:           file.Name,
			Size:           int64(file.UncompressedSize64),
			CompressedSize: int64(file.CompressedSize64),
			ModTime:        file.Modified,
			Mode:           file.Mode(),
			IsDir:          file.Mode().IsDir(),
			IsSymlink:      file.Mode()&os.ModeSymlink != 0,
		}

		// 如果是符号链接，读取链接目标
		if fileInfo.IsSymlink {
			if target, err := readSymlinkTarget(file); err == nil {
				fileInfo.LinkTarget = target
			}
		}

		archiveInfo.Files = append(archiveInfo.Files, fileInfo)
		archiveInfo.TotalSize += fileInfo.Size
	}

	return archiveInfo, nil
}

// ListZipLimit 获取ZIP压缩包指定数量的文件信息
func ListZipLimit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "ZIP文件路径")
	if err != nil {
		return nil, err
	}

	// 打开ZIP文件
	reader, err := zip.OpenReader(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开ZIP文件失败: %w", err)
	}
	defer func() { _ = reader.Close() }()

	// 获取压缩包文件信息
	stat, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("获取ZIP文件信息失败: %w", err)
	}

	// 计算实际需要处理的文件数量
	totalFiles := len(reader.File)
	maxFiles := totalFiles
	if limit > 0 && limit < totalFiles {
		maxFiles = limit
	}

	// 根据文件名检测压缩格式类型
	compressType, err := utils.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 创建 ArchiveInfo 结构体
	archiveInfo := &types.ArchiveInfo{
		Type:           compressType,
		TotalFiles:     maxFiles, // 注意：这里是实际返回的文件数
		CompressedSize: stat.Size(),
		Files:          make([]types.FileInfo, 0, maxFiles), // 优化容量分配
	}

	// 只遍历前 maxFiles 个文件
	for i := 0; i < maxFiles; i++ {
		file := reader.File[i]
		fileInfo := types.FileInfo{
			Name:           file.Name,
			Size:           int64(file.UncompressedSize64),
			CompressedSize: int64(file.CompressedSize64),
			ModTime:        file.Modified,
			Mode:           file.Mode(),
			IsDir:          file.Mode().IsDir(),
			IsSymlink:      file.Mode()&os.ModeSymlink != 0,
		}

		// 如果是符号链接，读取链接目标
		if fileInfo.IsSymlink {
			if target, err := readSymlinkTarget(file); err == nil {
				fileInfo.LinkTarget = target
			}
		}

		archiveInfo.Files = append(archiveInfo.Files, fileInfo)
		archiveInfo.TotalSize += fileInfo.Size
	}

	return archiveInfo, nil
}

// ListZipMatch 获取ZIP压缩包中匹配指定模式的文件信息
func ListZipMatch(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListZip(archivePath)
	if err != nil {
		return nil, err
	}

	archiveInfo.Files = utils.FilterFilesByPattern(archiveInfo.Files, pattern)
	archiveInfo.TotalFiles = len(archiveInfo.Files)

	// 重新计算总大小
	var totalSize int64
	for _, file := range archiveInfo.Files {
		totalSize += file.Size
	}
	archiveInfo.TotalSize = totalSize

	return archiveInfo, nil
}

// readSymlinkTarget 读取ZIP文件中符号链接的目标
func readSymlinkTarget(file *zip.File) (string, error) {
	reader, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = reader.Close() }()

	// 读取符号链接目标
	target := make([]byte, file.UncompressedSize64)
	n, err := reader.Read(target)
	if err != nil {
		return "", err
	}

	return string(target[:n]), nil
}

// Package cxtar 提供 TAR 格式的压缩包内容列表功能实现。
//
// 该包实现了 TAR 格式压缩包的文件信息获取功能，包括基本列表、限制数量列表和模式匹配列表。
// TAR 是归档格式，不进行压缩，因此压缩大小等于原始大小。
//
// 主要功能：
//   - TAR 压缩包完整文件列表获取
//   - 限制数量的文件列表获取
//   - 模式匹配的文件列表过滤
//   - 多种文件类型支持（普通文件、目录、符号链接、硬链接）
//   - 文件元数据完整保存
//
// 文件类型支持：
//   - 普通文件：完整的文件信息
//   - 目录：目录标识和权限信息
//   - 符号链接：链接目标路径保存
//   - 硬链接：链接目标路径保存
//
// 元数据信息：
//   - 文件名和路径
//   - 文件大小（原始大小和压缩大小相同）
//   - 修改时间
//   - 文件权限模式
//   - 文件类型标识
//   - 符号链接目标
//
// 使用示例：
//
//	// 获取 TAR 文件完整列表
//	info, err := cxtar.ListTar("archive.tar")
//
//	// 获取前 10 个文件信息
//	info, err := cxtar.ListTarLimit("archive.tar", 10)
//
//	// 获取匹配 *.go 模式的文件
//	info, err := cxtar.ListTarMatch("archive.tar", "*.go")
package cxtar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
)

// ListTar 获取TAR压缩包的所有文件信息
func ListTar(archivePath string) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "TAR文件路径")
	if err != nil {
		return nil, err
	}

	// 打开TAR文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开TAR文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 获取压缩包文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取TAR文件信息失败: %w", err)
	}

	// 创建TAR读取器
	tarReader := tar.NewReader(file)

	// 根据文件名检测压缩格式类型
	compressType, err := utils.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 创建TAR文件信息
	archiveInfo := &types.ArchiveInfo{
		Type:           compressType,
		CompressedSize: stat.Size(),
		Files:          make([]types.FileInfo, 0, utils.DefaultFileCapacity),
	}

	// 遍历TAR文件中的每个条目
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取TAR条目失败: %w", err)
		}

		fileInfo := types.FileInfo{
			Name:           header.Name,
			Size:           header.Size,
			CompressedSize: header.Size, // TAR禁用压缩，压缩大小等于原始大小
			ModTime:        header.ModTime,
			Mode:           header.FileInfo().Mode(),
			IsDir:          header.FileInfo().IsDir(),
			IsSymlink:      header.Typeflag == tar.TypeSymlink || header.Typeflag == tar.TypeLink,
		}

		// 如果是符号链接，设置链接目标
		if fileInfo.IsSymlink {
			fileInfo.LinkTarget = header.Linkname
		}

		archiveInfo.Files = append(archiveInfo.Files, fileInfo)
		archiveInfo.TotalSize += fileInfo.Size
		archiveInfo.TotalFiles++
	}

	return archiveInfo, nil
}

// ListTarLimit 获取TAR压缩包指定数量的文件信息
func ListTarLimit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "TAR文件路径")
	if err != nil {
		return nil, err
	}

	// 打开TAR文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开TAR文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 获取压缩包文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取TAR文件信息失败: %w", err)
	}

	// 创建TAR读取器
	tarReader := tar.NewReader(file)

	// 根据文件名检测压缩格式类型
	compressType, err := utils.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 创建TAR文件信息
	archiveInfo := &types.ArchiveInfo{
		Type:           compressType,
		CompressedSize: stat.Size(),
		Files:          make([]types.FileInfo, 0, utils.DefaultFileCapacity),
	}

	// 遍历TAR文件中的每个条目，但限制数量
	count := 0
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取TAR条目失败: %w", err)
		}

		// 达到限制数量就提前退出
		if limit > 0 && count >= limit {
			break
		}

		fileInfo := types.FileInfo{
			Name:           header.Name,
			Size:           header.Size,
			CompressedSize: header.Size, // TAR禁用压缩，压缩大小等于原始大小
			ModTime:        header.ModTime,
			Mode:           header.FileInfo().Mode(),
			IsDir:          header.FileInfo().IsDir(),
			IsSymlink:      header.Typeflag == tar.TypeSymlink || header.Typeflag == tar.TypeLink,
		}

		// 如果是符号链接，设置链接目标
		if fileInfo.IsSymlink {
			fileInfo.LinkTarget = header.Linkname
		}

		archiveInfo.Files = append(archiveInfo.Files, fileInfo)
		archiveInfo.TotalSize += fileInfo.Size
		count++
	}

	archiveInfo.TotalFiles = count
	return archiveInfo, nil
}

// ListTarMatch 获取TAR压缩包中匹配指定模式的文件信息
func ListTarMatch(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListTar(archivePath)
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

// Package cxbzip2 提供 BZIP2 格式的压缩包内容列表功能实现。
//
// 该包实现了 BZIP2 格式压缩包的文件信息获取功能，包括基本列表、限制数量列表和模式匹配列表。
// 由于 BZIP2 是单文件压缩格式，所有列表操作都针对单个文件进行处理。
//
// 主要功能：
//   - BZIP2 压缩包文件信息获取
//   - 原始文件名智能推导
//   - 原始文件大小计算
//   - 模式匹配过滤
//   - 压缩率计算
//
// 特殊处理：
//   - 智能推导原始文件名（去除 .bz2 或 .bzip2 后缀）
//   - 文件名缺失时使用默认后缀
//   - 通过完整读取计算原始文件大小
//   - 使用压缩文件修改时间（BZIP2 不保存原始时间）
//   - 使用默认文件权限（BZIP2 不保存权限信息）
//
// 格式限制：
//   - BZIP2 格式不保存文件元数据
//   - 无原始文件名信息
//   - 无修改时间信息
//   - 无文件权限信息
//
// 性能优化：
//   - 使用缓冲区池减少内存分配
//   - 高效的文件大小计算方法
//   - 错误容忍的大小估算机制
//
// 使用示例：
//
//	// 获取 BZIP2 文件信息
//	info, err := cxbzip2.ListBz2("archive.bz2")
//
//	// 获取匹配模式的文件信息
//	info, err := cxbzip2.ListBz2Match("archive.bzip2", "*.txt")
//
//	// 限制返回文件数量（对 BZIP2 无实际效果）
//	info, err := cxbzip2.ListBz2Limit("archive.bz2", 10)
package cxbzip2

import (
	"compress/bzip2"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
	"gitee.com/MM-Q/go-kit/pool"
)

// ListBz2 获取BZ2压缩包的文件信息
func ListBz2(archivePath string) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "BZ2文件路径")
	if err != nil {
		return nil, err
	}

	// 打开BZ2文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开BZ2文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 获取压缩包文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取BZ2文件信息失败: %w", err)
	}

	// 创建BZ2读取器
	bz2Reader := bzip2.NewReader(file)

	// 获取原始文件名
	baseName := filepath.Base(absPath)
	var originalName string
	if ext := filepath.Ext(baseName); ext == ".bz2" {
		originalName = baseName[:len(baseName)-len(ext)]
	} else if ext == ".bzip2" {
		originalName = baseName[:len(baseName)-len(ext)]
	} else {
		originalName = baseName + utils.DecompressedSuffix
	}

	// BZ2是单文件压缩，需要读取整个文件来获取原始大小
	// 使用io.CopyBuffer配合io.Discard，既高效又准确
	buffer := pool.GetByteCap(utils.DefaultBufferSize)
	defer pool.PutByte(buffer)

	originalSize, err := io.CopyBuffer(io.Discard, bz2Reader, buffer)
	if err != nil {
		// 如果读取失败，使用压缩文件大小作为估算
		originalSize = stat.Size()
	}

	// 创建BZ2文件信息
	fileInfo := types.FileInfo{
		Name:           originalName,
		Size:           originalSize,
		CompressedSize: stat.Size(),
		ModTime:        stat.ModTime(),
		Mode:           utils.DefaultFileMode, // BZ2不保存文件权限，使用默认权限
		IsDir:          false,
		IsSymlink:      false,
	}

	// 根据文件名检测压缩格式类型
	compressType, err := utils.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 创建BZ2文件信息
	archiveInfo := &types.ArchiveInfo{
		Type:           compressType,
		TotalFiles:     1,
		TotalSize:      originalSize,
		CompressedSize: stat.Size(),
		Files:          []types.FileInfo{fileInfo},
	}

	return archiveInfo, nil
}

// ListBz2Limit 获取BZ2压缩包指定数量的文件信息
func ListBz2Limit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListBz2(archivePath)
	if err != nil {
		return nil, err
	}

	// BZ2只有一个文件，limit不影响结果
	return archiveInfo, nil
}

// ListBz2Match 获取BZ2压缩包中匹配指定模式的文件信息
func ListBz2Match(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListBz2(archivePath)
	if err != nil {
		return nil, err
	}

	// 检查单个文件是否匹配模式
	if len(archiveInfo.Files) > 0 && utils.MatchPattern(archiveInfo.Files[0].Name, pattern) {
		return archiveInfo, nil
	}

	// 如果不匹配，返回空列表
	archiveInfo.Files = []types.FileInfo{}
	archiveInfo.TotalFiles = 0
	archiveInfo.TotalSize = 0

	return archiveInfo, nil
}

// Package cxgzip 提供 GZIP 格式的压缩包内容列表功能实现。
//
// 该包实现了 GZIP 格式压缩包的文件信息获取功能，包括基本列表、限制数量列表和模式匹配列表。
// 由于 GZIP 是单文件压缩格式，所有列表操作都针对单个文件进行处理。
//
// 主要功能：
//   - GZIP 压缩包文件信息获取
//   - 原始文件名和大小计算
//   - 文件修改时间获取
//   - 模式匹配过滤
//   - 压缩率计算
//
// 特殊处理：
//   - 自动从 GZIP 文件头获取原始文件名
//   - 文件名缺失时智能推导（去除 .gz 后缀）
//   - 通过完整读取计算原始文件大小
//   - 使用默认文件权限（GZIP 不保存权限信息）
//
// 性能优化：
//   - 使用缓冲区池减少内存分配
//   - 高效的文件大小计算方法
//   - 错误容忍的大小估算机制
//
// 使用示例：
//
//	// 获取 GZIP 文件信息
//	info, err := cxgzip.ListGzip("archive.gz")
//
//	// 获取匹配模式的文件信息
//	info, err := cxgzip.ListGzipMatch("archive.gz", "*.txt")
//
//	// 限制返回文件数量（对 GZIP 无实际效果）
//	info, err := cxgzip.ListGzipLimit("archive.gz", 10)
package cxgzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
	"gitee.com/MM-Q/go-kit/pool"
)

// ListGzip 获取GZIP压缩包的文件信息
func ListGzip(archivePath string) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "GZIP文件路径")
	if err != nil {
		return nil, err
	}

	// 根据文件名检测压缩格式类型
	compressType, err := utils.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 打开GZIP文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开GZIP文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 获取压缩包文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取GZIP文件信息失败: %w", err)
	}

	// 创建GZIP读取器
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("创建GZIP读取器失败: %w", err)
	}
	defer func() { _ = gzipReader.Close() }()

	// 获取原始文件名
	originalName := gzipReader.Name
	if originalName == "" {
		// 如果GZIP头中没有文件名，从压缩包文件名推导
		baseName := filepath.Base(absPath)
		if ext := filepath.Ext(baseName); ext == ".gz" {
			//originalName = baseName[:len(baseName)-len(ext)]
			// 去除.gz后缀
			originalName = strings.TrimSuffix(baseName, ".gz")
		} else {
			originalName = baseName + utils.DecompressedSuffix
		}
	}

	// GZIP是单文件压缩，需要读取整个文件来获取原始大小
	// 使用io.CopyBuffer配合io.Discard，既高效又准确
	buffer := pool.GetByteCap(utils.DefaultBufferSize)
	defer pool.PutByte(buffer)

	originalSize, err := io.CopyBuffer(io.Discard, gzipReader, buffer)
	if err != nil {
		// 如果读取失败，使用压缩文件大小作为估算
		originalSize = stat.Size()
	}

	// 创建FileInfo
	fileInfo := types.FileInfo{
		Name:           originalName,
		Size:           originalSize,
		CompressedSize: stat.Size(),
		ModTime:        gzipReader.ModTime,
		Mode:           utils.DefaultFileMode, // GZIP不保存文件权限，使用默认权限
		IsDir:          false,
		IsSymlink:      false,
	}

	// 创建ArchiveInfo
	archiveInfo := &types.ArchiveInfo{
		Type:           compressType,               // 类型
		TotalFiles:     1,                          // 文件数量
		TotalSize:      originalSize,               // 原始文件大小
		CompressedSize: stat.Size(),                // 压缩文件大小
		Files:          []types.FileInfo{fileInfo}, // 文件列表
	}

	return archiveInfo, nil
}

// ListGzipLimit 获取GZIP压缩包指定数量的文件信息
func ListGzipLimit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListGzip(archivePath)
	if err != nil {
		return nil, err
	}

	// GZIP只有一个文件，limit不影响结果
	return archiveInfo, nil
}

// ListGzipMatch 获取GZIP压缩包中匹配指定模式的文件信息
func ListGzipMatch(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListGzip(archivePath)
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

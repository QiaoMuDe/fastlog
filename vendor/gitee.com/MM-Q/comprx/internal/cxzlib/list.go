// Package cxzlib 提供 ZLIB 格式的压缩包内容列表功能实现。
//
// 该包实现了 ZLIB 格式压缩包的文件信息获取功能，包括基本列表、限制数量列表和模式匹配列表。
// 由于 ZLIB 是单文件压缩格式，所有列表操作都针对单个文件进行处理。
//
// 主要功能：
//   - ZLIB 压缩包文件信息获取
//   - 原始文件名智能推导
//   - 原始文件大小计算
//   - 模式匹配过滤
//   - 压缩率计算
//
// 特殊处理：
//   - 智能推导原始文件名（去除 .zlib 后缀）
//   - 文件名缺失时使用默认后缀
//   - 通过完整读取计算原始文件大小
//   - 使用压缩文件修改时间（ZLIB 不保存原始时间）
//   - 使用默认文件权限（ZLIB 不保存权限信息）
//
// 格式限制：
//   - ZLIB 格式不保存文件元数据
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
//	// 获取 ZLIB 文件信息
//	info, err := cxzlib.ListZlib("archive.zlib")
//
//	// 获取匹配模式的文件信息
//	info, err := cxzlib.ListZlibMatch("archive.zlib", "*.txt")
//
//	// 限制返回文件数量（对 ZLIB 无实际效果）
//	info, err := cxzlib.ListZlibLimit("archive.zlib", 10)
package cxzlib

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
	"gitee.com/MM-Q/go-kit/pool"
)

// ListZlib 获取ZLIB压缩包的文件信息
func ListZlib(archivePath string) (*types.ArchiveInfo, error) {
	// 确保路径为绝对路径
	absPath, err := utils.EnsureAbsPath(archivePath, "ZLIB文件路径")
	if err != nil {
		return nil, err
	}

	// 根据文件名检测压缩格式类型
	compressType, err := utils.DetectCompressFormat(absPath)
	if err != nil {
		return nil, fmt.Errorf("检测压缩格式失败: %w", err)
	}

	// 打开ZLIB文件
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开ZLIB文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	// 获取压缩包文件信息
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取ZLIB文件信息失败: %w", err)
	}

	// 创建ZLIB读取器
	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("创建ZLIB读取器失败: %w", err)
	}
	defer func() { _ = zlibReader.Close() }()

	// 获取原始文件名（ZLIB格式没有文件名信息，从压缩包文件名推导）
	baseName := filepath.Base(absPath)
	var originalName string
	if ext := filepath.Ext(baseName); ext == ".zlib" {
		// 去除.zlib后缀
		originalName = strings.TrimSuffix(baseName, ".zlib")
	} else {
		originalName = baseName + utils.DecompressedSuffix
	}

	// ZLIB是单文件压缩，需要读取整个文件来获取原始大小
	// 使用io.CopyBuffer配合io.Discard，既高效又准确
	buffer := pool.GetByteCap(utils.DefaultBufferSize)
	defer pool.PutByte(buffer)

	originalSize, err := io.CopyBuffer(io.Discard, zlibReader, buffer)
	if err != nil {
		// 如果读取失败，使用压缩文件大小作为估算
		originalSize = stat.Size()
	}

	// 创建FileInfo
	fileInfo := types.FileInfo{
		Name:           originalName,
		Size:           originalSize,
		CompressedSize: stat.Size(),
		ModTime:        stat.ModTime(),        // ZLIB不保存修改时间，使用压缩文件的修改时间
		Mode:           utils.DefaultFileMode, // ZLIB不保存文件权限，使用默认权限
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

// ListZlibLimit 获取ZLIB压缩包指定数量的文件信息
func ListZlibLimit(archivePath string, limit int) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListZlib(archivePath)
	if err != nil {
		return nil, err
	}

	// ZLIB只有一个文件，limit不影响结果
	return archiveInfo, nil
}

// ListZlibMatch 获取ZLIB压缩包中匹配指定模式的文件信息
func ListZlibMatch(archivePath string, pattern string) (*types.ArchiveInfo, error) {
	archiveInfo, err := ListZlib(archivePath)
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

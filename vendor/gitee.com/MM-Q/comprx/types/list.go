// Package types 定义了压缩包文件信息和压缩包整体信息的数据结构。
//
// 该文件提供了用于表示压缩包内容的核心数据结构，包括单个文件的详细信息
// 和整个压缩包的统计信息。这些结构体用于压缩包内容列表功能。
//
// 主要类型：
//   - FileInfo: 压缩包内单个文件的详细信息
//   - ArchiveInfo: 压缩包的整体信息和文件列表
//
// 主要功能：
//   - 存储文件的基本属性（名称、大小、时间等）
//   - 记录压缩相关信息（原始大小、压缩后大小）
//   - 支持符号链接信息
//   - 提供压缩包统计信息
//
// 使用示例：
//
//	// 获取压缩包信息
//	info, err := comprx.List("archive.zip")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 访问压缩包统计信息
//	fmt.Printf("文件总数: %d\n", info.TotalFiles)
//	fmt.Printf("原始大小: %d 字节\n", info.TotalSize)
//
//	// 遍历文件列表
//	for _, file := range info.Files {
//	    fmt.Printf("文件: %s, 大小: %d\n", file.Name, file.Size)
//	}
package types

import (
	"os"
	"time"
)

// FileInfo 压缩包内文件信息
type FileInfo struct {
	Name           string      // 文件名/路径
	Size           int64       // 原始大小
	CompressedSize int64       // 压缩后大小
	ModTime        time.Time   // 修改时间
	Mode           os.FileMode // 文件权限
	IsDir          bool        // 是否为目录
	IsSymlink      bool        // 是否为符号链接
	LinkTarget     string      // 符号链接目标(如果是符号链接)
}

// ArchiveInfo 压缩包整体信息
type ArchiveInfo struct {
	Type           CompressType // 压缩包类型
	TotalFiles     int          // 总文件数
	TotalSize      int64        // 总原始大小
	CompressedSize int64        // 总压缩大小
	Files          []FileInfo   // 文件列表
}

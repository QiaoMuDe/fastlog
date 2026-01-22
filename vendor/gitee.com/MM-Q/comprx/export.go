// Package comprx 提供核心类型和常量的导出定义。
//
// 该文件将 types 包中的核心类型和常量重新导出到主包中，为用户提供更简洁统一的 API。
// 用户可以直接通过 comprx 包访问所有核心类型，而无需额外导入 types 包。
//
// 导出的类型包括：
//   - FilterOptions: 文件过滤配置选项
//   - CompressionLevel: 压缩等级类型及相关常量
//   - CompressType: 压缩格式类型及相关常量
//   - FileInfo: 压缩包内文件信息
//   - ArchiveInfo: 压缩包整体信息
//   - ProgressStyle: 进度条样式类型及相关常量
//
// 所有导出的类型都是原始类型的别名，与原始类型完全兼容，可以互换使用。
// 这种设计既提供了便利的用户体验，又保持了与原有 types 包的向后兼容性。
//
// 使用示例：
//
//	// 使用导出的类型和常量
//	filter := &comprx.FilterOptions{
//	    Include: []string{"*.go"},
//	    MaxSize: 10 * 1024 * 1024,
//	}
//
//	opts := comprx.Options{
//	    CompressionLevel: comprx.CompressionLevelBest,
//	    ProgressStyle:    comprx.ProgressStyleUnicode,
//	    Filter:          filter,
//	}
//
//	err := comprx.PackOptions("output.zip", "src/", opts)
package comprx

import "gitee.com/MM-Q/comprx/types"

// ==============================================
// 导出类型
// ==============================================

// FilterOptions 过滤配置选项
//
// 用于指定压缩时或解压时的文件过滤条件：
//   - Include: 包含模式列表，支持 glob 语法，只有匹配的文件才会被处理
//   - Exclude: 排除模式列表，支持 glob 语法，匹配的文件会被跳过
//   - MaxSize: 最大文件大小限制(字节)，0 表示无限制，超过此大小的文件会被跳过
//   - MinSize: 最小文件大小限制(字节)，默认为 0，小于此大小的文件会被跳过
type FilterOptions = types.FilterOptions

// CompressionLevel 压缩等级类型
//
// 支持的压缩等级：
//   - CompressionLevelDefault: 默认压缩等级(支持该等级的类型: zip, tgz, tar.gz, zlib, gz)
//   - CompressionLevelNone: 禁用压缩(支持该等级的类型: zip, tgz, tar.gz, zlib, gz)
//   - CompressionLevelFast: 快速压缩(支持该等级的类型: tgz, tar.gz, zlib, gz)
//   - CompressionLevelBest: 最佳压缩(支持该等级的类型: tgz, tar.gz, zlib, gz)
//   - CompressionLevelHuffmanOnly: 仅使用Huffman编码(支持该等级的类型: tgz, tar.gz, zlib, gz)
type CompressionLevel = types.CompressionLevel

// CompressType 压缩格式类型
//
// 支持的压缩格式：
//   - CompressTypeZip: zip 压缩格式
//   - CompressTypeTar: tar 压缩格式
//   - CompressTypeTgz: tgz 压缩格式
//   - CompressTypeTarGz: tar.gz 压缩格式
//   - CompressTypeGz: gz 压缩格式
//   - CompressTypeBz2: bz2 压缩格式
//   - CompressTypeBzip2: bzip2 压缩格式
//   - CompressTypeZlib: zlib 压缩格式
type CompressType = types.CompressType

// FileInfo 压缩包内文件信息
//
// 包含文件的详细信息：
//   - Name: 文件名/路径
//   - Size: 原始大小
//   - CompressedSize: 压缩后大小
//   - ModTime: 修改时间
//   - Mode: 文件权限
//   - IsDir: 是否为目录
//   - IsSymlink: 是否为符号链接
//   - LinkTarget: 符号链接目标(如果是符号链接)
type FileInfo = types.FileInfo

// ArchiveInfo 压缩包整体信息
//
// 包含压缩包的统计信息和文件列表：
//   - Type: 压缩包类型
//   - TotalFiles: 总文件数
//   - TotalSize: 总原始大小
//   - CompressedSize: 总压缩大小
//   - Files: 文件列表
type ArchiveInfo = types.ArchiveInfo

// ==============================================
// 压缩等级常量
// ==============================================

const (
	// CompressionLevelDefault 默认压缩等级
	CompressionLevelDefault = types.CompressionLevelDefault

	// CompressionLevelNone 禁用压缩
	CompressionLevelNone = types.CompressionLevelNone

	// CompressionLevelFast 快速压缩
	CompressionLevelFast = types.CompressionLevelFast

	// CompressionLevelBest 最佳压缩
	CompressionLevelBest = types.CompressionLevelBest

	// CompressionLevelHuffmanOnly 仅使用Huffman编码
	CompressionLevelHuffmanOnly = types.CompressionLevelHuffmanOnly
)

// ==============================================
// 压缩格式常量
// ==============================================

const (
	// CompressTypeZip zip 压缩格式
	CompressTypeZip = types.CompressTypeZip

	// CompressTypeTar tar 压缩格式
	CompressTypeTar = types.CompressTypeTar

	// CompressTypeTgz tgz 压缩格式
	CompressTypeTgz = types.CompressTypeTgz

	// CompressTypeTarGz tar.gz 压缩格式
	CompressTypeTarGz = types.CompressTypeTarGz

	// CompressTypeGz gz 压缩格式
	CompressTypeGz = types.CompressTypeGz

	// CompressTypeBz2 bz2 压缩格式
	CompressTypeBz2 = types.CompressTypeBz2

	// CompressTypeBzip2 bzip2 压缩格式
	CompressTypeBzip2 = types.CompressTypeBzip2

	// CompressTypeZlib zlib 压缩格式
	CompressTypeZlib = types.CompressTypeZlib
)

// ProgressStyle 进度条样式类型
//
// 进度条样式类型定义：
//   - ProgressStyleText: 文本样式进度条 - 使用文字描述进度
//   - ProgressStyleDefault: 默认进度条样式 - progress库的默认进度条样式
//   - ProgressStyleUnicode: Unicode样式进度条 - 使用Unicode字符绘制精美进度条
//   - ProgressStyleASCII: ASCII样式进度条 - 使用基础ASCII字符绘制兼容性最好的进度条
type ProgressStyle = types.ProgressStyle

// ==============================================
// 进度条样式常量
// ==============================================

const (
	// ProgressStyleText 文本样式进度条 - 使用文字描述进度
	ProgressStyleText = types.ProgressStyleText

	// ProgressStyleDefault 默认进度条样式 - progress库的默认进度条样式
	ProgressStyleDefault = types.ProgressStyleDefault

	// ProgressStyleUnicode Unicode样式进度条 - 使用Unicode字符绘制精美进度条
	// 示例: ████████████░░░░░░░░ 60%
	ProgressStyleUnicode = types.ProgressStyleUnicode

	// ProgressStyleASCII ASCII样式进度条 - 使用基础ASCII字符绘制兼容性最好的进度条
	// 示例: [##########          ] 50%
	ProgressStyleASCII = types.ProgressStyleASCII
)

// Package config 提供压缩器的配置管理功能。
//
// 该包定义了压缩器的核心配置结构体，包括压缩等级、文件覆盖策略、
// 进度显示、路径验证和文件过滤等配置选项。提供了配置的创建和
// 压缩等级转换等实用功能。
//
// 主要类型：
//   - Config: 压缩器配置结构体
//
// 主要功能：
//   - 创建默认配置
//   - 压缩等级转换
//   - 进度显示配置
//   - 文件过滤配置
//   - 路径验证配置
//
// 使用示例：
//
//	// 创建默认配置
//	cfg := config.New()
//
//	// 设置压缩等级
//	cfg.CompressionLevel = types.CompressionLevelBest
//
//	// 启用文件覆盖
//	cfg.OverwriteExisting = true
package config

import (
	"compress/gzip"

	"gitee.com/MM-Q/comprx/internal/progress"
	"gitee.com/MM-Q/comprx/types"
)

// Config 压缩器配置
type Config struct {
	CompressionLevel      types.CompressionLevel // 压缩等级
	OverwriteExisting     bool                   // 是否覆盖已存在的文件
	Progress              *progress.Progress     // 进度显示
	DisablePathValidation bool                   // 是否禁用路径验证
	Filter                *types.FilterOptions   // 文件过滤配置
}

// New 创建新的压缩器配置
func New() *Config {
	return &Config{
		CompressionLevel:      types.CompressionLevelDefault, // 默认压缩等级
		OverwriteExisting:     false,                         // 默认不覆盖已存在文件
		Progress:              progress.New(),                // 创建进度显示
		DisablePathValidation: false,                         // 默认启用路径验证
		Filter:                nil,                           // 初始化空过滤器(不启用过滤时为nil)
	}
}

// GetCompressionLevel 根据配置返回对应的压缩等级
//
// 参数:
//   - level: types.CompressionLevel - 压缩等级
//
// 返回值:
//   - int - 压缩等级
func GetCompressionLevel(level types.CompressionLevel) int {
	switch level {
	case types.CompressionLevelNone: // 不进行压缩(禁用压缩)
		return gzip.NoCompression

	case types.CompressionLevelFast: // 快速压缩(压缩速度最快)
		return gzip.BestSpeed

	case types.CompressionLevelBest: // 最佳压缩(压缩率最高)
		return gzip.BestCompression

	case types.CompressionLevelHuffmanOnly: // 只使用哈夫曼编码(仅对文本文件有效)
		return gzip.HuffmanOnly

	default: // 默认的压缩等级(在压缩速度和压缩率之间取得平衡)
		return gzip.DefaultCompression
	}
}

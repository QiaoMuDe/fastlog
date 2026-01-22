// Package core 提供压缩库的核心功能实现。
//
// 该包实现了压缩库的主要业务逻辑，包括压缩和解压缩操作的统一接口。
// 支持多种压缩格式的自动检测和处理，提供了配置化的压缩解压缩功能。
//
// 主要类型：
//   - Comprx: 压缩器核心结构体
//
// 主要功能：
//   - 统一的压缩和解压缩接口
//   - 自动压缩格式检测
//   - 支持多种压缩格式（ZIP、TAR、TGZ、GZIP、BZIP2、ZLIB）
//   - 配置化的压缩参数
//   - 文件和目录的智能处理
//
// 支持的压缩格式：
//   - ZIP: .zip
//   - TAR: .tar
//   - TGZ: .tgz, .tar.gz
//   - GZIP: .gz
//   - BZIP2: .bz2, .bzip2（仅解压）
//   - ZLIB: .zlib
//
// 使用示例：
//
//	// 创建压缩器实例
//	comprx := core.New()
//
//	// 压缩文件
//	err := comprx.Pack("output.zip", "input_dir")
//
//	// 解压文件
//	err := comprx.Unpack("archive.zip", "output_dir")
package core

import (
	"fmt"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/internal/config"
	"gitee.com/MM-Q/comprx/internal/cxbzip2"
	"gitee.com/MM-Q/comprx/internal/cxgzip"
	"gitee.com/MM-Q/comprx/internal/cxtar"
	"gitee.com/MM-Q/comprx/internal/cxtgz"
	"gitee.com/MM-Q/comprx/internal/cxzip"
	"gitee.com/MM-Q/comprx/internal/cxzlib"
	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
)

// Comprx 压缩器
type Comprx struct {
	Config *config.Config // 压缩器配置
}

// ==============================================
// 构造函数
// ==============================================

// New 创建压缩器实例(NewComprx的别名)
//
// 返回:
//   - *Comprx: 压缩器实例
var New = NewComprx

// NewComprx 创建压缩器实例
//
// 返回:
//   - *Comprx: 压缩器实例
func NewComprx() *Comprx {
	return &Comprx{
		Config: config.New(),
	}
}

// ==============================================
// 打包压缩方法
// ==============================================

// Pack 压缩文件或目录
//
// 参数:
//   - dst: 目标压缩文件路径
//   - src: 源文件或目录路径
//
// 返回:
//   - error: 错误信息
func (c *Comprx) Pack(dst string, src string) error {
	// 检查参数
	if src == "" || dst == "" {
		return fmt.Errorf("源路径或目标路径不能为空")
	}

	// 规范化路径
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	// 智能检测压缩文件格式
	compressType, err := utils.DetectCompressFormat(dst)
	if err != nil {
		return fmt.Errorf("检测压缩格式失败: %v", err)
	}

	// 检查是否为.bz2格式的压缩文件，暂不支持
	if compressType == types.CompressTypeBz2 || compressType == types.CompressTypeBzip2 {
		return fmt.Errorf("暂不支持 %s 和 %s 格式的压缩文件", types.CompressTypeBz2.String(), types.CompressTypeBzip2.String())
	}

	// 检查目标文件是否存在
	if utils.Exists(dst) {
		if !c.Config.OverwriteExisting {
			return fmt.Errorf("文件 %s 已存在，如需覆盖请设置 OverwriteExisting 为 true", dst)
		}
	}

	// 检查目标目录是否存在, 不存在则创建
	targetDir := filepath.Dir(dst)
	if err := utils.EnsureDir(targetDir); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 根据压缩格式进行打包
	switch compressType {
	case types.CompressTypeZip: // Zip
		return cxzip.Zip(dst, src, c.Config)

	case types.CompressTypeTar: // Tar
		return cxtar.Tar(dst, src, c.Config)

	case types.CompressTypeTgz, types.CompressTypeTarGz: // Tar.gz 或 .tgz
		return cxtgz.Tgz(dst, src, c.Config)

	case types.CompressTypeGz: // Gz
		return cxgzip.Gzip(dst, src, c.Config)

	case types.CompressTypeZlib: // Zlib
		return cxzlib.Zlib(dst, src, c.Config)

	default:
		return fmt.Errorf("不支持的压缩格式: %s", compressType)
	}
}

// ==============================================
// 解压方法
// ==============================================

// Unpack 解压文件
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标目录路径
//
// 返回:
//   - error: 错误信息
func (c *Comprx) Unpack(src string, dst string) error {
	// 检查源路径是否为空
	if src == "" {
		return fmt.Errorf("源路径不能为空")
	}

	// 规范化路径
	src = filepath.Clean(src)
	if dst != "" {
		dst = filepath.Clean(dst)
	}

	// 智能检测压缩文件格式
	compressType, err := utils.DetectCompressFormat(src)
	if err != nil {
		return fmt.Errorf("检测压缩格式失败: %v", err)
	}

	// 检查源文件是否存在
	if !utils.Exists(src) {
		return fmt.Errorf("源文件 %s 不存在", src)
	}

	// 当目标目录为空时，自动生成目标目录, 如: /path/to/file.tar.gz -> /path/to/file
	if dst == "" {
		baseName := filepath.Base(src)
		baseName = strings.TrimSuffix(baseName, ".tar.gz")
		baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
		dst = filepath.Join(filepath.Dir(src), baseName)
	}

	// 检查目标目录是否存在, 不存在则创建
	if err := utils.EnsureDir(dst); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 根据压缩格式进行解压
	switch compressType {
	case types.CompressTypeZip: // Zip
		return cxzip.Unzip(src, dst, c.Config)

	case types.CompressTypeTar: // Tar
		return cxtar.Untar(src, dst, c.Config)

	case types.CompressTypeTgz, types.CompressTypeTarGz: // Tgz, TarGz
		return cxtgz.Untgz(src, dst, c.Config)

	case types.CompressTypeGz: // Gzip
		return cxgzip.Ungzip(src, dst, c.Config)

	case types.CompressTypeBz2, types.CompressTypeBzip2: // Bz2, Bzip2
		return cxbzip2.Unbz2(src, dst, c.Config)

	case types.CompressTypeZlib: // Zlib
		return cxzlib.Unzlib(src, dst, c.Config)

	default:
		return fmt.Errorf("不支持的压缩格式: %s", compressType)
	}
}

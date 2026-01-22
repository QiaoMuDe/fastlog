// Package cxzlib 提供 ZLIB 格式的解压缩功能实现。
//
// 该包实现了 ZLIB 格式的单文件解压缩操作，支持进度显示和智能目标路径处理。
// ZLIB 格式是一种无文件头的压缩格式，不保存原始文件名和修改时间信息。
//
// 主要功能：
//   - ZLIB 格式单文件解压缩
//   - 进度显示支持
//   - 智能目标路径处理
//   - 文件覆盖控制
//   - 高效的缓冲区管理
//
// 智能处理：
//   - 目标为目录时自动生成文件名
//   - 自动去除 .zlib 扩展名作为文件名
//   - 自动创建目标文件的父目录
//
// 解压特性：
//   - 使用 DEFLATE 解压算法
//   - 包含 Adler-32 校验和验证
//   - 流式解压，内存占用低
//   - 支持大文件解压
//
// 性能优化：
//   - 智能缓冲区大小选择
//   - 进度条模式下的大小预计算
//   - 错误容忍的大小估算机制
//
// 使用示例：
//
//	// 创建配置
//	cfg := config.New()
//	cfg.OverwriteExisting = true
//
//	// 解压文件到指定路径
//	err := cxzlib.Unzlib("archive.zlib", "output.txt", cfg)
//
//	// 解压文件到目录（自动生成文件名）
//	err := cxzlib.Unzlib("archive.zlib", "output_dir/", cfg)
package cxzlib

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/internal/config"
	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
	"gitee.com/MM-Q/go-kit/pool"
)

// calculateZlibTotalSize 计算ZLIB文件的解压后大小
//
// 参数:
//   - zlibFilePath: ZLIB文件路径
//   - cfg: 解压配置
//
// 返回值:
//   - int64: 解压后的文件大小（字节）
func calculateZlibTotalSize(zlibFilePath string, cfg *config.Config) int64 {
	// 只在进度条模式下计算总大小
	if !cfg.Progress.Enabled || cfg.Progress.BarStyle == types.ProgressStyleText {
		return 0
	}

	// 开始扫描进度显示
	bar := cfg.Progress.StartScan("正在分析内容...")
	defer func() {
		_ = cfg.Progress.CloseBar(bar)
	}()

	// 打开ZLIB文件进行扫描
	zlibFile, err := os.Open(zlibFilePath)
	if err != nil {
		return 0
	}
	defer func() { _ = zlibFile.Close() }()

	// 创建ZLIB读取器
	zlibReader, err := zlib.NewReader(zlibFile)
	if err != nil {
		return 0
	}
	defer func() { _ = zlibReader.Close() }()

	// 由于ZLIB是流式压缩，我们需要读取整个文件来计算大小
	// 使用进度条作为写入器，直接通过io.CopyBuffer计算总大小
	buffer := pool.GetByteCap(utils.DefaultBufferSize) // 32KB缓冲区
	defer pool.PutByte(buffer)

	totalSize, err := io.CopyBuffer(bar, zlibReader, buffer)
	if err != nil {
		return 0 // 如果出错，返回0表示无法计算大小
	}

	return totalSize
}

// Unzlib 解压缩 ZLIB 文件
//
// 参数:
//   - zlibFilePath: 要解压缩的 ZLIB 文件路径
//   - targetPath: 解压缩后的目标文件路径
//   - config: 解压缩配置
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Unzlib(zlibFilePath string, targetPath string, config *config.Config) error {
	// 在进度条模式下计算总大小
	totalSize := calculateZlibTotalSize(zlibFilePath, config)

	// 开始进度显示
	if err := config.Progress.Start(totalSize, zlibFilePath, fmt.Sprintf("正在解压'%s'...", filepath.Base(zlibFilePath))); err != nil {
		return fmt.Errorf("开始进度显示失败: %w", err)
	}
	defer func() {
		_ = config.Progress.Close()
	}()

	// 打开 ZLIB 文件（同时检查文件是否存在）
	zlibFile, err := os.Open(zlibFilePath)
	if err != nil {
		return fmt.Errorf("打开 ZLIB 文件失败: %w", err)
	}
	defer func() { _ = zlibFile.Close() }()

	// 获取ZLIB文件信息用于预验证
	zlibInfo, err := zlibFile.Stat()
	if err != nil {
		return fmt.Errorf("获取ZLIB文件信息失败: %w", err)
	}

	// 创建 ZLIB 读取器
	zlibReader, err := zlib.NewReader(zlibFile)
	if err != nil {
		return fmt.Errorf("创建 ZLIB 读取器失败: %w", err)
	}
	defer func() { _ = zlibReader.Close() }()

	// 检查目标路径状态，处理目录情况和覆盖检查
	if targetStat, statErr := os.Stat(targetPath); statErr == nil {
		if targetStat.IsDir() {
			// 目标是目录，生成文件名（去掉.zlib扩展名）
			baseName := filepath.Base(zlibFilePath)
			baseName = strings.TrimSuffix(baseName, ".zlib")
			targetPath = filepath.Join(targetPath, baseName)

			// 重新检查生成的目标文件是否存在
			if _, statErr := os.Stat(targetPath); statErr == nil && !config.OverwriteExisting {
				return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
			}
		} else {
			// 目标是文件，检查是否允许覆盖
			if !config.OverwriteExisting {
				return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
			}
		}
	}

	// 检查目标文件的父目录是否存在，如果不存在则创建
	parentDir := filepath.Dir(targetPath)
	if mkdirErr := utils.EnsureDir(parentDir); mkdirErr != nil {
		return fmt.Errorf("创建目标文件父目录失败: %w", mkdirErr)
	}

	// 创建目标文件
	targetFile, createErr := os.Create(targetPath)
	if createErr != nil {
		return fmt.Errorf("创建目标文件失败: %w", createErr)
	}
	defer func() { _ = targetFile.Close() }()

	// 使用之前获取的zlibInfo来估算缓冲区大小
	bufferSize := pool.CalculateBufferSize(zlibInfo.Size())
	buffer := pool.GetByteCap(bufferSize)
	defer pool.PutByte(buffer)

	// 打印解压缩进度
	config.Progress.Inflating(targetPath)

	// 解压缩文件内容
	if _, err := config.Progress.CopyBuffer(targetFile, zlibReader, buffer); err != nil {
		return fmt.Errorf("解压缩文件失败: %w", err)
	}

	return nil
}

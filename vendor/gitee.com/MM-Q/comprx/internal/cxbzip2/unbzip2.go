// Package cxbzip2 提供 BZIP2 格式的解压缩功能实现。
//
// 该包实现了 BZIP2 格式的单文件解压缩操作，支持进度显示和智能目标路径处理。
// BZIP2 格式是一种高压缩比的压缩格式，不保存原始文件名和修改时间信息。
//
// 主要功能：
//   - BZIP2 格式单文件解压缩
//   - 进度显示支持
//   - 智能目标路径处理
//   - 文件覆盖控制
//   - 路径安全验证
//
// 智能处理：
//   - 目标为目录时自动生成文件名
//   - 自动去除 .bz2 和 .bzip2 扩展名作为文件名
//   - 自动创建目标文件的父目录
//   - 路径安全验证防止路径遍历攻击
//
// 解压特性：
//   - 使用 Burrows-Wheeler 变换算法
//   - 高压缩比，适合文本文件
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
//	err := cxbzip2.Unbz2("archive.bz2", "output.txt", cfg)
//
//	// 解压文件到目录（自动生成文件名）
//	err := cxbzip2.Unbz2("archive.bzip2", "output_dir/", cfg)
package cxbzip2

import (
	"compress/bzip2"
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

// calculateBzip2TotalSize 计算BZIP2文件的解压后大小
//
// 参数:
//   - bz2FilePath: BZIP2文件路径
//   - cfg: 解压配置
//
// 返回值:
//   - int64: 解压后的文件大小（字节）
func calculateBzip2TotalSize(bz2FilePath string, cfg *config.Config) int64 {
	// 只在进度条模式下计算总大小
	if !cfg.Progress.Enabled || cfg.Progress.BarStyle == types.ProgressStyleText {
		return 0
	}

	// 开始扫描进度显示
	bar := cfg.Progress.StartScan("正在分析内容...")
	defer func() {
		_ = cfg.Progress.CloseBar(bar)
	}()

	// 打开BZIP2文件进行扫描
	bz2File, err := os.Open(bz2FilePath)
	if err != nil {
		return 0
	}
	defer func() { _ = bz2File.Close() }()

	// 创建BZIP2读取器
	bz2Reader := bzip2.NewReader(bz2File)

	// 由于BZIP2是流式压缩，我们需要读取整个文件来计算大小
	// 使用进度条作为写入器，直接通过io.CopyBuffer计算总大小
	buffer := pool.GetByteCap(utils.DefaultBufferSize)
	defer pool.PutByte(buffer)

	totalSize, err := io.CopyBuffer(bar, bz2Reader, buffer)
	if err != nil {
		return 0 // 如果出错，返回0表示无法计算大小
	}

	return totalSize
}

// Unbz2 解压缩 BZIP2 文件
//
// 参数:
//   - bz2FilePath: 要解压缩的 BZIP2 文件路径
//   - targetPath: 解压缩后的目标文件路径
//   - cfg: 解压缩配置
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Unbz2(bz2FilePath string, targetPath string, cfg *config.Config) error {
	// 在进度条模式下计算总大小
	totalSize := calculateBzip2TotalSize(bz2FilePath, cfg)

	// 开始进度显示
	if err := cfg.Progress.Start(totalSize, bz2FilePath, fmt.Sprintf("正在解压'%s'...", filepath.Base(bz2FilePath))); err != nil {
		return fmt.Errorf("开始进度显示失败: %w", err)
	}
	defer func() {
		_ = cfg.Progress.Close()
	}()

	// 打开 BZIP2 文件（同时检查文件是否存在）
	bz2File, err := os.Open(bz2FilePath)
	if err != nil {
		return fmt.Errorf("打开 BZIP2 文件失败: %w", err)
	}
	defer func() { _ = bz2File.Close() }()

	// 获取BZIP2文件信息
	bz2Info, err := bz2File.Stat()
	if err != nil {
		return fmt.Errorf("获取BZIP2文件信息失败: %w", err)
	}

	// 创建 BZIP2 读取器
	bz2Reader := bzip2.NewReader(bz2File)

	// 检查目标路径状态，处理目录情况和覆盖检查
	if targetStat, err := os.Stat(targetPath); err == nil {
		if targetStat.IsDir() {
			// 目标是目录，生成文件名
			baseName := filepath.Base(bz2FilePath)
			baseName = strings.TrimSuffix(baseName, ".bz2")
			baseName = strings.TrimSuffix(baseName, ".bzip2")

			// 添加安全验证
			validatedPath, err := utils.ValidatePathSimple(targetPath, baseName, cfg.DisablePathValidation)
			if err != nil {
				return fmt.Errorf("BZIP2文件名包含不安全的路径: %w", err)
			}
			targetPath = validatedPath

			// 重新检查生成的目标文件是否存在
			if _, err := os.Stat(targetPath); err == nil && !cfg.OverwriteExisting {
				return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
			}
		} else {
			// 目标是文件，检查是否允许覆盖
			if !cfg.OverwriteExisting {
				return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
			}
		}
	}

	// 检查目标文件的父目录是否存在，如果不存在则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("创建目标文件父目录失败: %w", err)
	}

	// 创建目标文件
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer func() { _ = targetFile.Close() }()

	// 获取缓冲区大小并创建缓冲区
	bufferSize := pool.CalculateBufferSize(bz2Info.Size())
	buffer := pool.GetByteCap(bufferSize)
	defer pool.PutByte(buffer)

	// 打印解压缩进度
	cfg.Progress.Inflating(targetPath)

	// 解压缩文件内容到目标文件
	if _, err := cfg.Progress.CopyBuffer(targetFile, bz2Reader, buffer); err != nil {
		return fmt.Errorf("解压缩文件失败: %w", err)
	}

	return nil
}

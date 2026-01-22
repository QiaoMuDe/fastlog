// Package cxgzip 提供 GZIP 格式的解压缩功能实现。
//
// 该包实现了 GZIP 格式的单文件解压缩操作，支持进度显示、文件元数据恢复和路径安全验证。
// 能够处理 GZIP 文件头中的原始文件名和修改时间信息。
//
// 主要功能：
//   - GZIP 格式单文件解压缩
//   - 进度显示支持
//   - 文件元数据恢复（文件名、修改时间）
//   - 路径安全验证
//   - 文件覆盖控制
//   - 智能目标路径处理
//
// 安全特性：
//   - 路径遍历攻击防护
//   - GZIP 文件头文件名验证
//   - 可配置的路径验证开关
//
// 智能处理：
//   - 自动从 GZIP 文件头获取原始文件名
//   - 目标为目录时自动生成文件名
//   - 自动去除 .gz 扩展名作为备选文件名
//
// 使用示例：
//
//	// 创建配置
//	cfg := config.New()
//	cfg.OverwriteExisting = true
//
//	// 解压文件到指定路径
//	err := cxgzip.Ungzip("archive.gz", "output.txt", cfg)
//
//	// 解压文件到目录（自动生成文件名）
//	err := cxgzip.Ungzip("archive.gz", "output_dir/", cfg)
package cxgzip

import (
	"compress/gzip"
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

// calculateGzipTotalSize 计算GZIP文件的解压后大小
//
// 参数:
//   - gzipFilePath: GZIP文件路径
//   - cfg: 解压配置
//
// 返回值:
//   - int64: 解压后的文件大小（字节）
func calculateGzipTotalSize(gzipFilePath string, cfg *config.Config) int64 {
	// 只在进度条模式下计算总大小
	if !cfg.Progress.Enabled || cfg.Progress.BarStyle == types.ProgressStyleText {
		return 0
	}

	// 开始扫描进度显示
	bar := cfg.Progress.StartScan("正在分析内容...")
	defer func() {
		_ = cfg.Progress.CloseBar(bar)
	}()

	// 打开GZIP文件进行扫描
	gzipFile, err := os.Open(gzipFilePath)
	if err != nil {
		return 0
	}
	defer func() { _ = gzipFile.Close() }()

	// 创建GZIP读取器
	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return 0
	}
	defer func() { _ = gzipReader.Close() }()

	// 由于GZIP是流式压缩，我们需要读取整个文件来计算大小
	// 使用进度条作为写入器，直接通过io.CopyBuffer计算总大小
	buffer := pool.GetByteCap(utils.DefaultBufferSize) // 32KB缓冲区
	defer pool.PutByte(buffer)

	totalSize, err := io.CopyBuffer(bar, gzipReader, buffer)
	if err != nil {
		return 0 // 如果出错，返回0表示无法计算大小
	}

	return totalSize
}

// Ungzip 解压缩 GZIP 文件
//
// 参数:
//   - gzipFilePath: 要解压缩的 GZIP 文件路径
//   - targetPath: 解压缩后的目标文件路径
//   - config: 解压缩配置
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Ungzip(gzipFilePath string, targetPath string, config *config.Config) error {
	// 在进度条模式下计算总大小
	totalSize := calculateGzipTotalSize(gzipFilePath, config)

	// 开始进度显示
	if err := config.Progress.Start(totalSize, gzipFilePath, fmt.Sprintf("正在解压'%s'...", filepath.Base(gzipFilePath))); err != nil {
		return fmt.Errorf("开始进度显示失败: %w", err)
	}
	defer func() {
		_ = config.Progress.Close()
	}()

	// 打开 GZIP 文件（同时检查文件是否存在）
	gzipFile, err := os.Open(gzipFilePath)
	if err != nil {
		return fmt.Errorf("打开 GZIP 文件失败: %w", err)
	}
	defer func() { _ = gzipFile.Close() }()

	// 获取GZIP文件信息用于预验证
	gzipInfo, err := gzipFile.Stat()
	if err != nil {
		return fmt.Errorf("获取GZIP文件信息失败: %w", err)
	}

	// 创建 GZIP 读取器
	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return fmt.Errorf("创建 GZIP 读取器失败: %w", err)
	}
	defer func() { _ = gzipReader.Close() }()

	// 检查目标路径状态，处理目录情况和覆盖检查
	if targetStat, statErr := os.Stat(targetPath); statErr == nil {
		if targetStat.IsDir() {
			// 目标是目录，生成文件名
			if gzipReader.Name != "" {
				// 直接验证 GZIP 头部的文件名，并与目标目录合并
				validatedPath, validateErr := utils.ValidatePathSimple(targetPath, gzipReader.Name, config.DisablePathValidation)
				if validateErr != nil {
					return fmt.Errorf("GZIP文件头包含不安全的文件名: %w", validateErr)
				}
				targetPath = validatedPath
			} else {
				// 如果GZIP文件头中没有原始文件名，则去掉.gz扩展名
				baseName := filepath.Base(gzipFilePath)
				baseName = strings.TrimSuffix(baseName, ".gz")
				targetPath = filepath.Join(targetPath, baseName)
			}

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

	// 使用之前获取的gzipInfo来估算缓冲区大小
	// 获取缓冲区大小并创建缓冲区
	bufferSize := pool.CalculateBufferSize(gzipInfo.Size())
	buffer := pool.GetByteCap(bufferSize)
	defer pool.PutByte(buffer)

	// 打印解压缩进度
	config.Progress.Inflating(targetPath)

	// 解压缩文件内容
	if _, err := config.Progress.CopyBuffer(targetFile, gzipReader, buffer); err != nil {
		return fmt.Errorf("解压缩文件失败: %w", err)
	}

	// 如果GZIP文件头中有修改时间信息，则设置目标文件的修改时间
	if !gzipReader.ModTime.IsZero() {
		if err := os.Chtimes(targetPath, gzipReader.ModTime, gzipReader.ModTime); err != nil {
			// 设置时间失败不是致命错误，只记录警告
			fmt.Printf("警告: 设置文件修改时间失败: %v\n", err)
		}
	}

	return nil
}

// Package cxzip 提供 ZIP 格式的解压缩功能实现。
//
// 该包实现了 ZIP 格式的解压缩操作，支持普通文件、目录、符号链接的解压。
// 提供了进度显示、文件过滤、路径安全验证和配置化的解压缩功能。
//
// 主要功能：
//   - ZIP 格式解压缩
//   - 支持普通文件、目录、符号链接解压
//   - 文件过滤功能
//   - 进度显示支持
//   - 路径安全验证
//   - 文件覆盖控制
//
// 安全特性：
//   - 路径遍历攻击防护
//   - 安全的文件路径验证
//   - 可配置的路径验证开关
//
// 支持的文件类型：
//   - 普通文件：完整解压文件内容
//   - 目录：创建目录结构
//   - 符号链接：重建符号链接
//   - 空文件：创建空文件
//
// 使用示例：
//
//	// 创建配置
//	cfg := config.New()
//	cfg.OverwriteExisting = true
//
//	// 解压文件
//	err := cxzip.Unzip("archive.zip", "output_dir", cfg)
package cxzip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/internal/config"
	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
	"gitee.com/MM-Q/go-kit/pool"
)

// Unzip 解压缩 ZIP 文件到指定目录
//
// 参数:
//   - zipFilePath: 要解压缩的 ZIP 文件路径
//   - targetDir: 解压缩后的目标目录路径
//   - cfg: 解压缩配置
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Unzip(zipFilePath string, targetDir string, cfg *config.Config) error {
	// 打开 ZIP 文件
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("打开 ZIP 文件失败: %w", err)
	}
	defer func() { _ = zipReader.Close() }()

	// 在进度条模式下计算总大小
	totalSize := calculateZipTotalSize(zipReader, cfg)

	// 开始进度显示
	if err := cfg.Progress.Start(totalSize, zipFilePath, fmt.Sprintf("正在解压'%s'...", filepath.Base(zipFilePath))); err != nil {
		return fmt.Errorf("开始进度显示失败: %w", err)
	}
	defer func() {
		_ = cfg.Progress.Close()
	}()

	// 检查目标目录是否存在, 如果不存在, 则创建
	if err := utils.EnsureDir(targetDir); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 遍历 ZIP 文件中的每个文件或目录
	for _, file := range zipReader.File {
		// 应用过滤器检查
		if cfg.Filter != nil {
			// 使用通用的过滤方法，传入文件路径、大小和是否为目录
			if cfg.Filter.ShouldSkipByParams(file.Name, int64(file.UncompressedSize64), file.Mode().IsDir()) {
				continue // 跳过此文件
			}
		}

		// 安全的路径验证和拼接
		targetPath, err := utils.ValidatePathSimple(targetDir, file.Name, cfg.DisablePathValidation)
		if err != nil {
			return fmt.Errorf("处理文件 '%s' 时路径验证失败: %w", file.Name, err)
		}

		// 获取文件的模式
		mode := file.Mode()

		// 使用 switch 语句处理不同类型的文件
		switch {
		// 处理目录
		case mode.IsDir():
			cfg.Progress.Creating(targetPath) // 更新进度
			if err := extractDirectory(targetPath, file.Name); err != nil {
				return err
			}

		// 处理软链接
		case mode&os.ModeSymlink != 0:
			cfg.Progress.Inflating(targetPath) // 更新进度
			if err := extractSymlink(file, targetPath); err != nil {
				return err
			}

		// 处理普通文件
		default:
			cfg.Progress.Inflating(targetPath) // 更新进度
			if err := extractRegularFileWithWriter(file, targetPath, mode, cfg); err != nil {
				return err
			}
		}
	}

	return nil
}

// calculateZipTotalSize 计算ZIP文件中所有普通文件的总大小
//
// 参数:
//   - zipReader: ZIP文件读取器
//   - cfg: 解压配置
//
// 返回值:
//   - int64: 普通文件的总大小（字节）
func calculateZipTotalSize(zipReader *zip.ReadCloser, cfg *config.Config) int64 {
	var totalSize int64

	// 只在进度条模式下计算总大小
	if !cfg.Progress.Enabled || cfg.Progress.BarStyle == types.ProgressStyleText {
		return 0
	}

	// 开始扫描进度显示
	bar := cfg.Progress.StartScan("正在分析内容...")
	defer func() {
		_ = cfg.Progress.CloseBar(bar)
	}()

	// 遍历ZIP文件中的所有条目
	for _, file := range zipReader.File {
		// 应用过滤器检查
		if cfg.Filter != nil {
			if cfg.Filter.ShouldSkipByParams(file.Name, int64(file.UncompressedSize64), file.Mode().IsDir()) {
				continue // 跳过被过滤的文件
			}
		}

		// 只计算普通文件的大小，跳过目录和软链接
		if file.Mode().IsRegular() {
			totalSize += int64(file.UncompressedSize64)   // 累加普通文件大小
			_ = bar.Add64(int64(file.UncompressedSize64)) // 更新进度条
		}
	}

	return totalSize
}

// extractDirectory 处理目录解压
//
// 参数:
//   - targetPath: 目标路径
//   - fileName: 文件名（用于错误信息）
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractDirectory(targetPath, fileName string) error {
	if err := utils.EnsureDir(targetPath); err != nil {
		return fmt.Errorf("处理目录 '%s' 时出错 - 创建目录失败: %w", fileName, err)
	}
	return nil
}

// extractSymlink 处理软链接解压
//
// 参数:
//   - file: ZIP文件条目
//   - targetPath: 目标路径
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractSymlink(file *zip.File, targetPath string) error {
	zipFileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 打开 ZIP 文件中的软链接失败: %w", file.Name, err)
	}
	defer func() { _ = zipFileReader.Close() }()

	// 使用 io.ReadAll 读取完整的软链接目标路径
	targetBytes, err := io.ReadAll(zipFileReader)
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 读取软链接目标失败: %w", file.Name, err)
	}
	target := string(targetBytes) // 软链接的目标

	// 检查软链接的父目录是否存在，如果不存在，则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建软链接父目录失败: %w", file.Name, err)
	}

	// 创建软链接
	if err := os.Symlink(target, targetPath); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建软链接失败: %w", file.Name, err)
	}

	return nil
}

// extractRegularFileWithWriter 处理普通文件解压的通用实现
//
// 参数:
//   - file: ZIP文件条目
//   - targetPath: 目标路径
//   - mode: 文件模式
//   - cfg: 解压配置
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractRegularFileWithWriter(file *zip.File, targetPath string, mode os.FileMode, cfg *config.Config) error {
	// 检查目标文件是否已存在
	if _, err := os.Stat(targetPath); err == nil {
		// 文件已存在，检查是否允许覆盖
		if !cfg.OverwriteExisting {
			return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
		}
	}

	// 检查file的父目录是否存在, 如果不存在, 则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建文件父目录失败: %w", file.Name, err)
	}

	// 获取文件的大小
	fileSize := file.UncompressedSize64

	// 如果文件大小为0，只创建空文件，不进行读写操作
	if fileSize == 0 {
		// 创建空文件
		emptyFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
		if err != nil {
			return fmt.Errorf("处理文件 '%s' 时出错 - 创建空文件失败: %w", file.Name, err)
		}
		defer func() { _ = emptyFile.Close() }()
		return nil
	}

	// 创建文件
	fileWriter, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建文件失败: %w", file.Name, err)
	}
	defer func() { _ = fileWriter.Close() }()

	// 打开 ZIP 文件中的文件
	zipFileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 打开 zip 文件中的文件失败: %w", file.Name, err)
	}
	defer func() { _ = zipFileReader.Close() }()

	// 获取对应文件大小的缓冲区
	bufferSize := pool.CalculateBufferSize(int64(fileSize))
	buffer := pool.GetByteCap(bufferSize)
	defer pool.PutByte(buffer)

	// 将文件内容写入目标文件
	if _, err := cfg.Progress.CopyBuffer(fileWriter, zipFileReader, buffer); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 写入文件失败: %w", file.Name, err)
	}

	return nil
}

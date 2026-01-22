// Package cxtgz 提供 TGZ (tar.gz) 格式的解压缩功能实现。
//
// 该包实现了 TGZ 格式的文件和目录解压缩操作，TGZ 是 TAR 归档格式与 GZIP 压缩的组合。
// 支持多种文件类型的处理，包括普通文件、目录、符号链接和硬链接，提供完整的进度显示、
// 路径安全验证和文件过滤功能。
//
// 主要功能：
//   - TGZ 格式文件和目录解压缩
//   - 支持多种文件类型（普通文件、目录、符号链接、硬链接）
//   - 进度显示支持
//   - 路径安全验证
//   - 文件过滤功能
//   - 文件覆盖控制
//
// 解压流程：
//  1. 打开 TGZ 文件
//  2. 创建 GZIP 解压缩流
//  3. 在 GZIP 流上创建 TAR 读取流
//  4. 按 TAR 格式解析并解压文件
//
// 文件类型支持：
//   - 普通文件：完整内容解压
//   - 目录：创建目录结构
//   - 符号链接：恢复链接关系
//   - 硬链接：创建硬链接
//   - 其他类型：跳过处理并提示
//
// 安全特性：
//   - 路径遍历攻击防护
//   - 文件路径验证
//   - 可配置的路径验证开关
//   - 文件覆盖保护
//
// 性能优化：
//   - 智能缓冲区大小选择
//   - 空文件特殊处理
//   - 进度条模式下的大小预计算
//
// 使用示例：
//
//	// 创建配置
//	cfg := config.New()
//	cfg.OverwriteExisting = true
//
//	// 解压 TGZ 文件
//	err := cxtgz.Untgz("archive.tar.gz", "output_dir", cfg)
package cxtgz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/internal/config"
	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
	"gitee.com/MM-Q/go-kit/pool"
)

// Untgz 解压缩 TGZ 文件到指定目录
//
// 参数:
//   - tgzFilePath: 要解压缩的 TGZ 文件路径
//   - targetDir: 解压缩后的目标目录路径
//   - cfg: 解压缩配置
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Untgz(tgzFilePath string, targetDir string, cfg *config.Config) error {
	// 在进度条模式下计算总大小
	totalSize := calculateTgzTotalSize(tgzFilePath, cfg)

	// 打开 TGZ 文件
	tgzFile, err := os.Open(tgzFilePath)
	if err != nil {
		return fmt.Errorf("打开 TGZ 文件失败: %w", err)
	}
	defer func() { _ = tgzFile.Close() }()

	// 创建 GZIP 读取器
	gzipReader, err := gzip.NewReader(tgzFile)
	if err != nil {
		return fmt.Errorf("创建 GZIP 读取器失败: %w", err)
	}
	defer func() { _ = gzipReader.Close() }()

	// 创建 TAR 读取器
	tarReader := tar.NewReader(gzipReader)

	// 开始进度显示
	if err := cfg.Progress.Start(totalSize, tgzFilePath, fmt.Sprintf("正在解压'%s'...", filepath.Base(tgzFilePath))); err != nil {
		return fmt.Errorf("开始进度显示失败: %w", err)
	}
	defer func() {
		_ = cfg.Progress.Close()
	}()

	// 检查目标目录是否存在, 如果不存在, 则创建
	if err := utils.EnsureDir(targetDir); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 遍历 TAR 文件中的每个文件或目录
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 到达文件末尾
		}
		if err != nil {
			return fmt.Errorf("读取 TAR 文件头失败: %w", err)
		}

		// 检查是否应该跳过此文件
		if cfg.Filter != nil {
			// 应用过滤器检查
			isDir := header.Typeflag == tar.TypeDir
			if cfg.Filter.ShouldSkipByParams(header.Name, header.Size, isDir) {
				continue // 跳过此文件
			}
		}

		// 安全的路径验证和拼接
		targetPath, err := utils.ValidatePathSimple(targetDir, header.Name, cfg.DisablePathValidation)
		if err != nil {
			return fmt.Errorf("处理文件 '%s' 时路径验证失败: %w", header.Name, err)
		}

		// 使用 switch 语句处理不同类型的文件
		switch header.Typeflag {
		case tar.TypeDir: // 处理目录
			cfg.Progress.Creating(targetPath) // 显示进度
			if err := extractDirectory(targetPath, header.Name); err != nil {
				return err
			}

		case tar.TypeReg: // 处理普通文件
			cfg.Progress.Inflating(targetPath) // 显示进度
			if err := extractRegularFile(tarReader, targetPath, header, cfg); err != nil {
				return err
			}

		case tar.TypeSymlink: // 处理符号链接
			cfg.Progress.Inflating(targetPath) // 显示进度
			if err := extractSymlink(header, targetPath); err != nil {
				return err
			}

		case tar.TypeLink: // 处理硬链接
			cfg.Progress.Inflating(targetPath) // 显示进度
			if err := extractHardlink(header, targetPath, targetDir); err != nil {
				return err
			}

		default:
			// 对于其他类型的文件，我们跳过处理
			fmt.Printf("跳过不支持的文件类型: %s (类型: %c)\n", header.Name, header.Typeflag)
		}
	}

	return nil
}

// calculateTgzTotalSize 计算TGZ文件中所有普通文件的总大小
//
// 参数:
//   - tgzFilePath: TGZ文件路径
//   - cfg: 解压配置
//
// 返回值:
//   - int64: 普通文件的总大小（字节）
func calculateTgzTotalSize(tgzFilePath string, cfg *config.Config) int64 {
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

	// 打开 TGZ 文件进行扫描
	tgzFile, err := os.Open(tgzFilePath)
	if err != nil {
		return 0
	}
	defer func() { _ = tgzFile.Close() }()

	// 创建 GZIP 读取器
	gzipReader, err := gzip.NewReader(tgzFile)
	if err != nil {
		return 0
	}
	defer func() { _ = gzipReader.Close() }()

	// 创建 TAR 读取器
	tarReader := tar.NewReader(gzipReader)

	// 遍历TAR文件中的所有条目
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			break // 出错时停止扫描
		}

		// 检查是否应该跳过此文件
		if cfg.Filter != nil {
			// 应用过滤器检查
			isDir := header.Typeflag == tar.TypeDir
			if cfg.Filter.ShouldSkipByParams(header.Name, header.Size, isDir) {
				continue
			}
		}

		// 只计算普通文件的大小
		if header.Typeflag == tar.TypeReg {
			totalSize += header.Size   // 累加普通文件大小
			_ = bar.Add64(header.Size) // 更新进度条
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

// extractRegularFile 处理普通文件解压
//
// 参数:
//   - tarReader: TAR读取器
//   - targetPath: 目标路径
//   - header: TAR文件头
//   - cfg: 解压配置
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractRegularFile(tarReader *tar.Reader, targetPath string, header *tar.Header, cfg *config.Config) error {
	// 检查目标文件是否已存在
	if _, err := os.Stat(targetPath); err == nil {
		// 文件已存在，检查是否允许覆盖
		if !cfg.OverwriteExisting {
			return fmt.Errorf("目标文件已存在且不允许覆盖: %s", targetPath)
		}
	}

	// 检查文件的父目录是否存在, 如果不存在, 则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建文件父目录失败: %w", header.Name, err)
	}

	// 获取文件的大小
	fileSize := header.Size

	// 如果文件大小为0，只创建空文件，不进行读写操作
	if fileSize == 0 {
		// 创建空文件
		emptyFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return fmt.Errorf("处理文件 '%s' 时出错 - 创建空文件失败: %w", header.Name, err)
		}
		defer func() { _ = emptyFile.Close() }()
		return nil
	}

	// 创建文件
	fileWriter, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建文件失败: %w", header.Name, err)
	}
	defer func() { _ = fileWriter.Close() }()

	// 获取缓冲区大小并创建缓冲区
	bufferSize := pool.CalculateBufferSize(fileSize)
	buffer := pool.GetByteCap(bufferSize)
	defer pool.PutByte(buffer)

	// 将文件内容写入目标文件
	if _, err := cfg.Progress.CopyBuffer(fileWriter, tarReader, buffer); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 写入文件失败: %w", header.Name, err)
	}

	return nil
}

// extractSymlink 处理软链接解压
//
// 参数:
//   - header: TAR文件头
//   - targetPath: 目标路径
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractSymlink(header *tar.Header, targetPath string) error {
	// 检查软链接的父目录是否存在，如果不存在，则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建软链接父目录失败: %w", header.Name, err)
	}

	// 创建软链接
	if err := os.Symlink(header.Linkname, targetPath); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建软链接失败: %w", header.Name, err)
	}

	return nil
}

// extractHardlink 处理硬链接解压
//
// 参数:
//   - header: TAR文件头
//   - targetPath: 目标路径
//   - targetDir: 目标目录
//
// 返回值:
//   - error: 操作过程中遇到的错误
func extractHardlink(header *tar.Header, targetPath, targetDir string) error {
	// 检查硬链接的父目录是否存在，如果不存在，则创建
	parentDir := filepath.Dir(targetPath)
	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf("处理硬链接 '%s' 时出错 - 创建硬链接父目录失败: %w", header.Name, err)
	}

	// 获取硬链接的源文件路径
	linkSourcePath := filepath.Join(targetDir, header.Linkname)

	// 创建硬链接
	if err := os.Link(linkSourcePath, targetPath); err != nil {
		return fmt.Errorf("处理硬链接 '%s' 时出错 - 创建硬链接失败: %w", header.Name, err)
	}

	return nil
}

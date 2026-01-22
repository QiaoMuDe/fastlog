// Package cxtar 提供 TAR 格式的归档功能实现。
//
// 该包实现了 TAR 格式的文件和目录归档操作，支持多种文件类型的处理，包括普通文件、
// 目录、符号链接和特殊文件。提供完整的进度显示和文件过滤功能。
//
// 主要功能：
//   - TAR 格式文件和目录归档
//   - 支持多种文件类型（普通文件、目录、符号链接、特殊文件）
//   - 进度显示支持
//   - 文件过滤功能
//   - 文件覆盖控制
//   - 相对路径处理
//
// 文件类型支持：
//   - 普通文件：完整内容复制
//   - 目录：创建目录条目
//   - 符号链接：保存链接目标
//   - 特殊文件：保存文件元数据
//
// 路径处理：
//   - 自动转换为 TAR 标准路径格式（正斜杠）
//   - 保留目录结构的相对路径
//   - 支持单文件和目录归档
//
// 使用示例：
//
//	// 创建配置
//	cfg := config.New()
//	cfg.OverwriteExisting = true
//
//	// 归档目录
//	err := cxtar.Tar("archive.tar", "source_dir", cfg)
//
//	// 归档单个文件
//	err := cxtar.Tar("file.tar", "single_file.txt", cfg)
package cxtar

import (
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/internal/config"
	"gitee.com/MM-Q/comprx/internal/progress"
	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/go-kit/pool"
)

// Tar 函数用于创建TAR归档文件
//
// 参数:
//   - dst: 生成的TAR文件路径
//   - src: 需要归档的源路径
//   - cfg: 压缩配置指针
//
// 返回值:
//   - error: 操作过程中遇到的错误
func Tar(dst string, src string, cfg *config.Config) error {
	// 确保路径为绝对路径
	var absErr error
	if dst, absErr = utils.EnsureAbsPath(dst, "TAR文件路径"); absErr != nil {
		return absErr
	}
	if src, absErr = utils.EnsureAbsPath(src, "源路径"); absErr != nil {
		return absErr
	}

	// 检查目标文件是否已存在
	if _, err := os.Stat(dst); err == nil {
		// 文件已存在，检查是否允许覆盖
		if !cfg.OverwriteExisting {
			return fmt.Errorf("目标文件已存在且不允许覆盖: %s", dst)
		}
	}

	// 确保目标目录存在
	if err := utils.EnsureDir(filepath.Dir(dst)); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 在进度条模式下计算源文件总大小
	totalSize := progress.CalculateSourceTotalSizeWithProgress(src, cfg.Progress, "正在分析内容...", cfg.Filter)

	// 开始进度显示
	if err := cfg.Progress.Start(totalSize, dst, fmt.Sprintf("正在压缩'%s'...", filepath.Base(dst))); err != nil {
		return fmt.Errorf("开始进度显示失败: %w", err)
	}
	defer func() {
		_ = cfg.Progress.Close()
	}()

	// 创建 TAR 文件
	tarFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建 TAR 文件失败: %w", err)
	}
	defer func() { _ = tarFile.Close() }()

	// 创建 TAR 写入器
	tarWriter := tar.NewWriter(tarFile)
	defer func() { _ = tarWriter.Close() }()

	// 检查源路径是文件还是目录
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源路径信息失败: %w", err)
	}

	// 根据源路径类型处理
	var tarErr error
	if srcInfo.IsDir() {
		// 遍历目录并添加文件到 TAR 包
		tarErr = walkDirectoryForTar(src, tarWriter, cfg)
	} else {
		// 单文件处理逻辑 - 检查是否应该跳过
		if cfg.Filter != nil && cfg.Filter.ShouldSkipByParams(src, srcInfo.Size(), srcInfo.IsDir()) {
			// 文件被过滤器跳过，直接返回成功
			return nil
		}
		cfg.Progress.Adding(src)
		tarErr = processRegularFile(tarWriter, src, filepath.Base(src), srcInfo, cfg)
	}

	// 检查是否有错误发生
	if tarErr != nil {
		return fmt.Errorf("打包目录到 TAR 失败: %w", tarErr)
	}

	return nil
}

// processDirectory 处理目录
//
// 参数:
//   - tarWriter: *tar.Writer - TAR 文件写入器
//   - headerName: string - TAR 文件中的目录名
//   - info: os.FileInfo - 目录信息
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processDirectory(tarWriter *tar.Writer, headerName string, info os.FileInfo) error {
	// 创建目录文件头
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("处理目录 '%s' 时出错 - 创建 TAR 文件头失败: %w", headerName, err)
	}
	// 设置目录名
	header.Name = headerName + "/" // 目录名后添加斜杠

	// 写入目录文件头
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("处理目录 '%s' 时出错 - 写入 TAR 目录头失败: %w", headerName, err)
	}
	return nil
}

// processSymlink 处理软链接
//
// 参数:
//   - tarWriter: *tar.Writer - TAR 文件写入器
//   - path: string - 软链接路径
//   - headerName: string - TAR 文件中的软链接名
//   - info: os.FileInfo - 文件信息
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processSymlink(tarWriter *tar.Writer, path, headerName string, info os.FileInfo) error {
	// 读取软链接目标
	target, err := os.Readlink(path)
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 读取软链接目标失败: %w", path, err)
	}

	// 创建软链接文件头
	header, err := tar.FileInfoHeader(info, target)
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建 TAR 文件头失败: %w", path, err)
	}
	header.Name = headerName

	// 写入软链接文件头
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 写入 TAR 软链接头失败: %w", path, err)
	}
	return nil
}

// processSpecialFile 处理特殊文件类型
//
// 参数:
//   - tarWriter: *tar.Writer - TAR 文件写入器
//   - headerName: string - TAR 文件中的特殊文件名
//   - info: os.FileInfo - 文件信息
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processSpecialFile(tarWriter *tar.Writer, headerName string, info os.FileInfo) error {
	// 创建 TAR 文件头
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("处理特殊文件 '%s' 时出错 - 创建 TAR 文件头失败: %w", headerName, err)
	}
	header.Name = headerName

	// 写入 TAR 文件头
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("处理特殊文件 '%s' 时出错 - 写入 TAR 特殊文件头失败: %w", headerName, err)
	}
	return nil
}

// walkDirectoryForTar 遍历目录并处理文件到TAR包
//
// 参数:
//   - src: string - 源目录路径
//   - tarWriter: *tar.Writer - TAR 文件写入器
//   - cfg: *config.Config - 配置
//
// 返回值:
//   - error - 操作过程中遇到的错误
func walkDirectoryForTar(src string, tarWriter *tar.Writer, cfg *config.Config) error {
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			// 如果不存在则忽略
			if os.IsNotExist(err) {
				return nil
			}
			// 其他错误
			return fmt.Errorf("遍历路径 '%s' 时出错: %w", path, err)
		}

		// 获取文件信息用于过滤检查
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("处理路径 '%s' 时出错 - 获取文件信息失败: %w", path, err)
		}

		// 应用过滤器检查
		if cfg.Filter != nil && cfg.Filter.ShouldSkipByParams(path, info.Size(), info.IsDir()) {
			if info.IsDir() {
				return filepath.SkipDir // 跳过整个目录
			}
			return nil // 跳过文件
		}

		// 获取相对路径，保留顶层目录
		headerName, err := filepath.Rel(filepath.Dir(src), path)
		if err != nil {
			return fmt.Errorf("处理路径 '%s' 时出错 - 获取相对路径失败: %w", path, err)
		}

		// 替换路径分隔符为正斜杠(TAR 文件格式要求)
		headerName = filepath.ToSlash(headerName)

		// 根据文件类型处理
		switch {
		// 处理普通文件
		case entry.Type().IsRegular():
			cfg.Progress.Adding(headerName) // 更新进度
			return processRegularFile(tarWriter, path, headerName, info, cfg)

		// 处理目录
		case entry.IsDir():
			cfg.Progress.Storing(headerName) // 更新进度
			return processDirectory(tarWriter, headerName, info)

		// 处理符号链接
		case entry.Type()&os.ModeSymlink != 0:
			cfg.Progress.Adding(headerName) // 更新进度
			return processSymlink(tarWriter, path, headerName, info)

		// 处理特殊文件
		default:
			cfg.Progress.Adding(headerName) // 更新进度
			return processSpecialFile(tarWriter, headerName, info)
		}
	})
}

// processRegularFile 处理普通文件
//
// 参数:
//   - tarWriter: *tar.Writer - TAR 文件写入器
//   - path: string - 源路径
//   - headerName: string - TAR 文件中的文件名
//   - info: os.FileInfo - 文件信息
//   - cfg: *config.Config - 配置
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processRegularFile(tarWriter *tar.Writer, path, headerName string, info os.FileInfo, cfg *config.Config) error {
	// 创建文件头
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建 TAR 文件头失败: %w", path, err)
	}
	header.Name = headerName // 设置文件名

	// 写入文件头
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 写入 TAR 文件头失败: %w", path, err)
	}

	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 打开文件失败: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	// 获取文件大小
	fileSize := info.Size()

	// 获取缓冲区大小并创建缓冲区
	bufferSize := pool.CalculateBufferSize(fileSize)
	buffer := pool.GetByteCap(bufferSize)
	defer pool.PutByte(buffer)

	// 复制文件内容到TAR写入器
	if _, err := cfg.Progress.CopyBuffer(tarWriter, file, buffer); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 写入 TAR 文件失败: %w", path, err)
	}

	return nil
}

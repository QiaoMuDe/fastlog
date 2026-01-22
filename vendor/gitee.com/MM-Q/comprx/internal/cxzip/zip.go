// Package cxzip 提供 ZIP 格式的压缩功能实现。
//
// 该包实现了 ZIP 格式的压缩操作，支持文件和目录的压缩，
// 包括普通文件、目录、符号链接和特殊文件的处理。
// 提供了进度显示、文件过滤和配置化的压缩功能。
//
// 主要功能：
//   - ZIP 格式压缩
//   - 支持文件和目录压缩
//   - 符号链接和特殊文件处理
//   - 文件过滤功能
//   - 进度显示支持
//   - 可配置的压缩等级
//
// 支持的文件类型：
//   - 普通文件：使用配置的压缩方法
//   - 目录：创建目录条目
//   - 符号链接：保存链接目标
//   - 特殊文件：创建占位符条目
//
// 使用示例：
//
//	// 创建配置
//	cfg := config.New()
//	cfg.CompressionLevel = types.CompressionLevelBest
//
//	// 压缩文件
//	err := cxzip.Zip("output.zip", "input_dir", cfg)
package cxzip

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gitee.com/MM-Q/comprx/internal/config"
	"gitee.com/MM-Q/comprx/internal/progress"
	"gitee.com/MM-Q/comprx/internal/utils"
	"gitee.com/MM-Q/comprx/types"
	"gitee.com/MM-Q/go-kit/pool"
)

// Zip 函数用于创建ZIP压缩文件
//
// 参数:
//   - dst: 生成的ZIP文件路径
//   - src: 需要压缩的源路径
//   - cfg: 压缩配置指针
//
// 返回值:
//   - error: 操作过程中遇到的错误
func Zip(dst string, src string, cfg *config.Config) error {
	// 确保路径为绝对路径
	var absErr error
	if dst, absErr = utils.EnsureAbsPath(dst, "ZIP文件路径"); absErr != nil {
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

	// 创建 ZIP 文件
	zipFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建 ZIP 文件失败: %w", err)
	}
	defer func() { _ = zipFile.Close() }()

	// 创建 ZIP 写入器（使用带验证的写入器）
	zipWriter := zip.NewWriter(zipFile)
	defer func() { _ = zipWriter.Close() }()

	// 检查源路径是文件还是目录
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源路径信息失败: %w", err)
	}

	// 根据源路径类型处理
	var zipErr error
	if srcInfo.IsDir() {
		// 遍历目录并添加文件到 ZIP 包
		zipErr = walkDirectoryForZip(src, zipWriter, cfg)
	} else {
		// 单文件处理逻辑 - 检查是否应该跳过
		if cfg.Filter != nil && cfg.Filter.ShouldSkipByParams(src, srcInfo.Size(), srcInfo.IsDir()) {
			// 文件被过滤器跳过，直接返回成功
			return nil
		}
		cfg.Progress.Adding(src)
		zipErr = processRegularFile(zipWriter, src, filepath.Base(src), srcInfo, cfg)
	}

	// 检查是否有错误发生
	if zipErr != nil {
		return fmt.Errorf("打包目录到 ZIP 失败: %w", zipErr)
	}

	return nil
}

// processRegularFile 处理普通文件
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - path: string - 文件路径
//   - headerName: string - ZIP 文件中的文件名
//   - info: os.FileInfo - 文件信息
//   - cfg: 压缩配置
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processRegularFile(zipWriter *zip.Writer, path, headerName string, info os.FileInfo, cfg *config.Config) error {
	// 创建文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建 ZIP 文件头失败: %w", path, err)
	}
	header.Name = headerName                  // 设置文件名
	header.Method = getCompressionMethod(cfg) // 使用配置的压缩方法

	// 创建 ZIP 写入器
	fileWriter, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 创建 ZIP 写入器失败: %w", path, err)
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

	// 复制文件内容到ZIP写入器
	if _, err := cfg.Progress.CopyBuffer(fileWriter, file, buffer); err != nil {
		return fmt.Errorf("处理文件 '%s' 时出错 - 写入 ZIP 文件失败: %w", path, err)
	}

	return nil
}

// processDirectory 处理目录
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - headerName: string - ZIP 文件中的目录名
//   - info: os.FileInfo - 目录信息
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processDirectory(zipWriter *zip.Writer, headerName string, info os.FileInfo) error {
	// 创建目录文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("处理目录 '%s' 时出错 - 创建 ZIP 文件头失败: %w", headerName, err)
	}
	// 设置目录名
	header.Name = headerName + "/" // 目录名后添加斜杠
	header.Method = zip.Store      // 使用禁用压缩的方法

	// 创建目录文件头
	if _, err := zipWriter.CreateHeader(header); err != nil {
		return fmt.Errorf("处理目录 '%s' 时出错 - 创建 ZIP 目录失败: %w", headerName, err)
	}
	return nil
}

// processSymlink 处理软链接
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - path: string - 软链接路径
//   - headerName: string - ZIP 文件中的软链接名
//   - mode: fs.FileMode - 文件模式
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processSymlink(zipWriter *zip.Writer, path, headerName string, mode fs.FileMode) error {
	// 读取软链接目标
	target, err := os.Readlink(path)
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 读取软链接目标失败: %w", path, err)
	}

	// 创建软链接文件头
	header := &zip.FileHeader{
		Name:   headerName,
		Method: zip.Store,
	}
	header.SetMode(mode)

	// 创建软链接文件
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 创建 ZIP 软链接失败: %w", path, err)
	}
	if _, err := writer.Write([]byte(target)); err != nil {
		return fmt.Errorf("处理软链接 '%s' 时出错 - 写入软链接目标失败: %w", path, err)
	}
	return nil
}

// processSpecialFile 处理特殊文件类型
//
// 参数:
//   - zipWriter: *zip.Writer - ZIP 文件写入器
//   - headerName: string - ZIP 文件中的特殊文件名
//   - mode: fs.FileMode - 文件模式
//
// 返回值:
//   - error - 操作过程中遇到的错误
func processSpecialFile(zipWriter *zip.Writer, headerName string, mode fs.FileMode) error {
	// 创建 ZIP 文件头
	header := &zip.FileHeader{
		Name:   headerName,
		Method: zip.Store,
	}
	header.SetMode(mode)

	// 创建 ZIP 文件写入器
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("处理特殊文件 '%s' 时出错 - 创建 ZIP 特殊文件失败: %w", headerName, err)
	}
	if _, err := writer.Write([]byte{}); err != nil {
		return fmt.Errorf("处理特殊文件 '%s' 时出错 - 写入特殊文件失败: %w", headerName, err)
	}
	return nil
}

// getCompressionMethod 根据配置返回对应的压缩方法
func getCompressionMethod(cfg *config.Config) uint16 {
	if cfg.CompressionLevel == types.CompressionLevelNone {
		return zip.Store // 禁用压缩，只存储
	}

	return zip.Deflate // 使用默认的压缩方法
}

// walkDirectoryForZip 遍历目录并处理文件到ZIP包
//
// 参数:
//   - src: 源目录路径
//   - zipWriter: ZIP写入器
//   - cfg: 压缩配置
//
// 返回值:
//   - error: 遍历过程中发生的错误
func walkDirectoryForZip(src string, zipWriter *zip.Writer, cfg *config.Config) error {
	return filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
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

		// 替换路径分隔符为正斜杠(ZIP 文件格式要求)
		headerName = filepath.ToSlash(headerName)

		// 根据文件类型处理
		switch {
		case entry.Type().IsRegular(): // 处理普通文件
			cfg.Progress.Adding(headerName) // 显示进度
			return processRegularFile(zipWriter, path, headerName, info, cfg)

		case entry.IsDir(): // 处理目录
			cfg.Progress.Storing(headerName) // 显示进度
			return processDirectory(zipWriter, headerName, info)

		case entry.Type()&fs.ModeSymlink != 0: // 处理符号链接
			cfg.Progress.Adding(headerName) // 显示进度
			return processSymlink(zipWriter, path, headerName, entry.Type())

		default: // 处理特殊文件
			cfg.Progress.Adding(headerName) // 显示进度
			return processSpecialFile(zipWriter, headerName, entry.Type())
		}
	})
}

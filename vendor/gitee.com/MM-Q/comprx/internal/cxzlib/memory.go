// Package cxzlib 提供 ZLIB 格式的内存压缩和流式压缩功能实现。
//
// 该包实现了 ZLIB 格式的内存中压缩和解压缩操作，以及流式压缩功能。
// 支持字节数组、字符串和流式数据的压缩与解压缩，提供了高性能的内存管理。
//
// 主要功能：
//   - ZLIB 内存压缩：字节数组和字符串的压缩解压
//   - ZLIB 流式压缩：支持 io.Reader 和 io.Writer 接口
//   - 支持自定义压缩等级
//   - 优化的内存分配策略
//   - 完善的错误处理和资源管理
//
// 压缩特性：
//   - 使用 DEFLATE 压缩算法
//   - 包含 Adler-32 校验和
//   - 比 GZIP 格式更紧凑（无文件头信息）
//   - 支持多种压缩等级
//
// 性能优化：
//   - 预分配缓冲区减少内存重分配
//   - 智能估算压缩后大小
//   - 直接字节操作避免额外拷贝
//   - 自动资源清理防止内存泄漏
//
// 使用示例：
//
//	// 压缩字节数据
//	compressed, err := cxzlib.CompressBytes(data, types.CompressionLevelBest)
//
//	// 解压字节数据
//	decompressed, err := cxzlib.DecompressBytes(compressed)
//
//	// 压缩字符串
//	compressed, err := cxzlib.CompressString("hello world", types.CompressionLevelDefault)
//
//	// 流式压缩
//	err := cxzlib.CompressStream(dst, src, types.CompressionLevelFast)
package cxzlib

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"

	"gitee.com/MM-Q/comprx/internal/config"
	"gitee.com/MM-Q/comprx/types"
)

// ================================ 内存压缩API ================================

// CompressBytes 压缩字节数据到内存
//
// 参数:
//   - data: 要压缩的字节数据
//   - level: 压缩级别
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
func CompressBytes(data []byte, level types.CompressionLevel) (result []byte, err error) {
	// 参数验证 - 更精确的nil检查
	if data == nil {
		return nil, fmt.Errorf("输入数据不能为nil")
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("输入数据不能为空")
	}

	// 创建内存缓冲区 - 预分配容量减少重分配
	// 预分配原大小的50%
	estimatedSize := len(data) / 2
	if estimatedSize < 64 {
		estimatedSize = 64 // 最小64字节
	}
	buf := bytes.NewBuffer(make([]byte, 0, estimatedSize))

	// 创建zlib写入器
	writer, err := zlib.NewWriterLevel(buf, config.GetCompressionLevel(level))
	if err != nil {
		return nil, fmt.Errorf("创建zlib写入器失败: %w", err)
	}

	// 直接写入数据，无需额外缓冲区
	if _, err = writer.Write(data); err != nil {
		_ = writer.Close() // 确保资源清理
		return nil, fmt.Errorf("压缩数据失败: %w", err)
	}

	// 关闭写入器确保数据完整写入
	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("完成压缩失败: %w", err)
	}

	return buf.Bytes(), nil
}

// DecompressBytes 从内存解压字节数据
//
// 参数:
//   - compressedData: 压缩的字节数据
//
// 返回:
//   - []byte: 解压后的数据
//   - error: 错误信息
func DecompressBytes(compressedData []byte) (result []byte, err error) {
	// 参数验证 - 更精确的nil检查
	if compressedData == nil {
		return nil, fmt.Errorf("压缩数据不能为nil")
	}
	if len(compressedData) == 0 {
		return nil, fmt.Errorf("压缩数据不能为空")
	}

	// 创建字节读取器
	reader := bytes.NewReader(compressedData)

	// 创建zlib读取器
	zlibReader, err := zlib.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("创建zlib读取器失败: %w", err)
	}

	// 预分配解压缓冲区 - 解压通常是压缩数据的2-3倍
	estimatedSize := len(compressedData) * 2
	if estimatedSize < 128 {
		estimatedSize = 128 // 最小128字节
	}
	buf := bytes.NewBuffer(make([]byte, 0, estimatedSize))

	// 直接读取解压数据，无需额外缓冲区
	if _, err = io.Copy(buf, zlibReader); err != nil {
		_ = zlibReader.Close() // 确保资源清理
		return nil, fmt.Errorf("解压数据失败: %w", err)
	}

	// 关闭读取器
	if err = zlibReader.Close(); err != nil {
		return nil, fmt.Errorf("关闭zlib读取器失败: %w", err)
	}

	return buf.Bytes(), nil
}

// CompressString 压缩字符串到内存
//
// 参数:
//   - text: 要压缩的字符串
//   - level: 压缩级别
//
// 返回:
//   - []byte: 压缩后的数据
//   - error: 错误信息
func CompressString(text string, level types.CompressionLevel) ([]byte, error) {
	// 快速失败判断
	if text == "" {
		return nil, fmt.Errorf("输入字符串不能为空")
	}

	// 直接复用CompressBytes
	return CompressBytes([]byte(text), level)
}

// DecompressString 从内存解压为字符串
//
// 参数:
//   - compressedData: 压缩的字节数据
//
// 返回:
//   - string: 解压后的字符串
//   - error: 错误信息
func DecompressString(compressedData []byte) (string, error) {
	// 快速失败判断
	if compressedData == nil {
		return "", fmt.Errorf("压缩数据不能为nil")
	}
	if len(compressedData) == 0 {
		return "", fmt.Errorf("压缩数据不能为空")
	}

	// 先解压为字节
	decompressed, err := DecompressBytes(compressedData)
	if err != nil {
		return "", err
	}

	// 转换为字符串
	return string(decompressed), nil
}

// ==================== 流式压缩API ====================

// CompressStream 流式压缩数据
//
// 参数:
//   - dst: 目标写入器
//   - src: 源读取器
//   - level: 压缩级别
//
// 返回:
//   - error: 错误信息
func CompressStream(dst io.Writer, src io.Reader, level types.CompressionLevel) (err error) {
	// 1. 参数验证
	if dst == nil {
		err = fmt.Errorf("目标写入器不能为nil")
		return
	}
	if src == nil {
		err = fmt.Errorf("源读取器不能为nil")
		return
	}

	// 2. 创建zlib写入器
	writer, createErr := zlib.NewWriterLevel(dst, config.GetCompressionLevel(level))
	if createErr != nil {
		err = fmt.Errorf("创建zlib写入器失败: %w", createErr)
		return
	}
	defer func() {
		if closeErr := writer.Close(); closeErr != nil && err == nil {
			// 只有在没有其他错误时才设置关闭错误
			err = fmt.Errorf("关闭zlib写入器失败: %w", closeErr)
		}
	}()

	// 3. 流式复制数据
	if _, copyErr := io.Copy(writer, src); copyErr != nil {
		err = fmt.Errorf("压缩数据失败: %w", copyErr)
		return
	}

	// 4. 确保数据完整写入
	if closeErr := writer.Close(); closeErr != nil {
		err = fmt.Errorf("完成压缩失败: %w", closeErr)
		return
	}

	return
}

// DecompressStream 流式解压数据
//
// 参数:
//   - dst: 目标写入器
//   - src: 源读取器（压缩数据）
//
// 返回:
//   - error: 错误信息
func DecompressStream(dst io.Writer, src io.Reader) (err error) {
	// 1. 参数验证
	if dst == nil {
		err = fmt.Errorf("目标写入器不能为nil")
		return
	}
	if src == nil {
		err = fmt.Errorf("源读取器不能为nil")
		return
	}

	// 2. 创建zlib读取器
	reader, createErr := zlib.NewReader(src)
	if createErr != nil {
		err = fmt.Errorf("创建zlib读取器失败: %w", createErr)
		return
	}
	defer func() {
		if closeErr := reader.Close(); closeErr != nil && err == nil {
			// 只有在没有其他错误时才设置关闭错误
			err = fmt.Errorf("关闭zlib读取器失败: %w", closeErr)
		}
	}()

	// 3. 流式复制数据
	if _, copyErr := io.Copy(dst, reader); copyErr != nil {
		err = fmt.Errorf("解压数据失败: %w", copyErr)
		return
	}

	return
}

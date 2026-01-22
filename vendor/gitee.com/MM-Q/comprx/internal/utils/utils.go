// Package utils 提供压缩库使用的通用工具函数。
//
// 该包包含了文件系统操作、路径处理、缓冲区管理等实用工具函数。
// 这些函数被压缩库的各个模块广泛使用，提供了统一的基础功能。
//
// 主要功能：
//   - 文件和目录存在性检查
//   - 目录创建和确保
//   - 动态缓冲区大小计算
//   - 路径安全验证和转换
//   - 绝对路径处理
//
// 安全特性：
//   - 路径遍历攻击防护
//   - 绝对路径检测
//   - UNC路径和协议前缀检测
//   - Windows特殊路径处理
//
// 使用示例：
//
//	// 检查文件是否存在
//	if utils.Exists("file.txt") {
//	    // 文件存在
//	}
//
//	// 确保目录存在
//	err := utils.EnsureDir("output/dir")
//
//	// 获取动态缓冲区大小
//	bufSize := utils.GetBufferSize(fileSize)
//
//	// 验证路径安全性
//	safePath, err := utils.ValidatePathSimple(targetDir, filePath, false)
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/MM-Q/comprx/types"
)

// Exists 检查指定路径的文件或目录是否存在
//
// 参数：
//   - path: 要检查的路径
//
// 返回值：
//   - bool: 如果文件或目录存在，则返回true，否则返回false
func Exists(path string) bool {
	// 使用os.Stat尝试获取文件信息
	_, err := os.Stat(path)

	// 如果没有错误，说明文件/目录存在
	if err == nil {
		return true
	}

	// 如果错误是文件不存在，则返回false
	if os.IsNotExist(err) {
		return false
	}

	// 其他错误情况（如权限问题等）也视为不存在
	// 根据实际需求，也可以选择返回错误
	return false
}

// EnsureDir 检查指定路径的目录是否存在，不存在则创建
//
// 参数：
//   - path: 要检查的目录路径
//
// 返回值：
//   - error: 如果创建目录成功，则返回nil，否则返回错误信息
func EnsureDir(path string) error {
	// 检查目录是否存在
	_, err := os.Stat(path)
	if err == nil {
		// 目录存在，返回nil
		return nil
	}

	// 检查错误是否为目录不存在
	if os.IsNotExist(err) {
		// 创建目录，使用0755权限，并递归创建父目录
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		return nil
	}

	// 其他错误（如权限问题等）
	return err
}

// EnsureAbsPath 确保路径为绝对路径，如果不是则转换为绝对路径
//
// 参数:
//   - path: 待检查的路径
//   - pathType: 路径类型描述（用于错误信息）
//
// 返回值:
//   - string: 绝对路径
//   - error: 转换过程中的错误
func EnsureAbsPath(path, pathType string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("转换 %s 为绝对路径失败: %w", pathType, err)
	}
	return absPath, nil
}

// ValidatePathSimple 验证文件路径是否安全，防止路径遍历攻击
//
// 参数：
//   - targetDir: 目标目录
//   - filePath: 要验证的文件路径
//   - skipValidation: 是否跳过安全验证（警告：仅在处理可信数据时使用）
//
// 返回值：
//   - string: 安全的文件路径
//   - error: 如果路径不安全，则返回错误信息
func ValidatePathSimple(targetDir, filePath string, skipValidation bool) (string, error) {
	// 如果跳过验证，直接拼接返回
	if skipValidation {
		return filepath.Join(targetDir, filePath), nil
	}

	// === 第一阶段：检查原始路径中的危险模式 ===
	// 注意：这些检查必须在filepath.Clean()之前进行，因为Clean会改变路径格式

	// 检查路径遍历攻击（最常见的攻击方式）
	if strings.Contains(filePath, "..") {
		return "", fmt.Errorf("不安全的路径: %s", filePath)
	}

	// 检查协议前缀攻击（如 file:// 等）
	if strings.Contains(filePath, "://") {
		return "", fmt.Errorf("不安全的路径: %s", filePath)
	}

	// === 第二阶段：清理路径并进行进一步检查 ===
	cleanPath := filepath.Clean(filePath)

	// 检查绝对路径 - Unix风格
	if strings.HasPrefix(cleanPath, "/") {
		return "", fmt.Errorf("不安全的路径: %s", filePath)
	}

	// 检查绝对路径 - Windows风格
	if len(cleanPath) >= 2 && cleanPath[1] == ':' {
		return "", fmt.Errorf("不安全的路径: %s", filePath)
	}

	// 检查UNC路径
	if strings.HasPrefix(cleanPath, "\\\\") || strings.HasPrefix(cleanPath, "//") {
		return "", fmt.Errorf("不安全的路径: %s", filePath)
	}

	// 检查Windows特殊路径前缀
	if strings.HasPrefix(cleanPath, "\\\\?\\") || strings.HasPrefix(cleanPath, "//?/") {
		return "", fmt.Errorf("不安全的路径: %s", filePath)
	}

	// 双重检查：确保Clean后没有残留的上级目录引用
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("不安全的路径: %s", filePath)
	}

	// === 第三阶段：构建最终安全路径 ===
	finalPath := filepath.Join(targetDir, cleanPath)
	return finalPath, nil
}

// DetectCompressFormat 智能检测压缩文件格式
//
// 参数:
//   - filename: 文件名
//
// 返回:
//   - types.CompressType: 检测到的压缩格式
//   - error: 错误信息
func DetectCompressFormat(filename string) (types.CompressType, error) {
	// 转换为小写进行处理
	lowerFilename := strings.ToLower(filename)

	// 处理.tar.gz特殊情况
	if strings.HasSuffix(lowerFilename, ".tar.gz") {
		return types.CompressTypeTarGz, nil
	}

	// 获取文件扩展名并转换为小写
	ext := strings.ToLower(filepath.Ext(filename))
	if !types.IsSupportedCompressType(ext) {
		return "", fmt.Errorf("不支持的压缩文件格式: %s, 支持的格式: %v", ext, types.SupportedCompressTypes())
	}

	return types.CompressType(ext), nil
}

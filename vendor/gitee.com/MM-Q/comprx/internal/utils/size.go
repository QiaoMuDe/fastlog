// Package utils 提供文件和目录大小计算的实用工具函数。
//
// 该文件实现了文件和目录大小的计算功能，支持单个文件大小获取和目录递归大小计算。
// 提供了安全版本和详细版本两种接口，满足不同的错误处理需求。
//
// 主要功能：
//   - 获取单个文件的大小
//   - 递归计算目录的总大小
//   - 提供安全版本（出错返回0）和详细版本（返回错误信息）
//   - 自动忽略符号链接等特殊文件
//   - 详细的错误分类和处理
//
// 错误处理：
//   - 文件不存在错误
//   - 权限不足错误
//   - 其他系统错误
//   - 遍历过程中的动态错误处理
//
// 使用示例：
//
//	// 安全版本，出错时返回0
//	size := utils.GetSizeOrZero("./mydir")
//
//	// 详细版本，返回错误信息
//	size, err := utils.GetSize("./myfile.txt")
//	if err != nil {
//	    log.Printf("获取大小失败: %v", err)
//	}
package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetSizeOrZero 获取文件或目录的大小，出错时返回 0
//
// 参数:
//   - path: 文件或目录路径
//
// 返回:
//   - int64: 文件或目录的总大小（字节），出错时返回 0
//
// 功能:
//   - 如果是文件，返回文件大小
//   - 如果是目录，返回目录中所有普通文件的总大小
//   - 忽略符号链接等特殊文件
//   - 发生任何错误时返回 0，不抛出异常
//
// 注意:
//   - 此函数为 GetSize 的安全版本，适用于不需要错误处理的场景
//   - 如需详细错误信息，请使用 GetSize 函数
func GetSizeOrZero(path string) int64 {
	if size, err := GetSize(path); err == nil {
		return size
	}
	return 0
}

// GetSize 获取文件或目录的大小(字节)
//
// 参数:
//   - path: 文件或目录路径
//
// 返回:
//   - int64: 文件或目录的总大小(字节)
//   - error: 错误信息
//
// 注意:
//   - 如果是文件，返回文件大小
//   - 如果是目录，返回目录中所有文件的总大小
//   - 如果路径不存在，返回错误
//   - 只计算普通文件的大小，忽略符号链接等特殊文件
func GetSize(path string) (int64, error) {
	// 获取路径信息
	info, err := os.Stat(path)
	if err != nil {
		// 判断错误类型，返回精准的错误信息
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("路径不存在: %s", path)
		}
		if os.IsPermission(err) {
			return 0, fmt.Errorf("访问路径 '%s' 时权限不足: %w", path, err)
		}
		// 其他错误
		return 0, fmt.Errorf("获取路径 '%s' 信息失败: %w", path, err)
	}

	// 如果是普通文件，直接返回文件大小
	if info.Mode().IsRegular() {
		return info.Size(), nil
	}

	// 如果不是目录，提前返回 0(符号链接等特殊文件)
	if !info.IsDir() {
		return 0, nil
	}

	// 如果是目录，遍历计算总大小
	var totalSize int64
	walkDirErr := filepath.WalkDir(path, func(walkPath string, entry os.DirEntry, err error) error {
		if err != nil {
			// 判断错误类型
			if os.IsNotExist(err) {
				// 文件不存在，忽略并继续遍历
				return nil
			}
			if os.IsPermission(err) {
				// 权限错误，返回具体错误信息
				return fmt.Errorf("访问路径 '%s' 时权限不足: %w", walkPath, err)
			}
			// 其他错误，返回通用错误信息
			return fmt.Errorf("访问路径 '%s' 时出错: %w", walkPath, err)
		}

		// 只计算普通文件的大小
		if entry.Type().IsRegular() {
			if info, err := entry.Info(); err == nil {
				// 累加文件大小
				totalSize += info.Size()
			} else {
				// 判断获取文件信息时的错误类型
				if os.IsNotExist(err) {
					// 文件在遍历过程中被删除，忽略
					return nil
				}
				if os.IsPermission(err) {
					// 权限错误
					return fmt.Errorf("获取文件 '%s' 信息时权限不足: %w", walkPath, err)
				}
				// 其他错误
				return fmt.Errorf("获取文件 '%s' 信息时出错: %w", walkPath, err)
			}
		}

		return nil
	})

	if walkDirErr != nil {
		return 0, fmt.Errorf("遍历目录失败: %w", walkDirErr)
	}

	return totalSize, nil
}

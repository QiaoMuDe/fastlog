// Package comprx 提供文件和目录大小计算功能。
//
// 该文件提供了计算文件或目录大小的实用函数。
// 支持单个文件大小获取和目录递归大小计算，提供安全和详细两种版本。
//
// 主要功能：
//   - 获取单个文件的大小
//   - 递归计算目录的总大小
//   - 提供安全版本（出错返回0）和详细版本（返回错误信息）
//   - 自动忽略符号链接等特殊文件
//
// 使用示例：
//
//	// 安全版本，出错时返回0
//	size := comprx.GetSizeOrZero("./mydir")
//
//	// 详细版本，返回错误信息
//	size, err := comprx.GetSize("./myfile.txt")
package comprx

import "gitee.com/MM-Q/comprx/internal/utils"

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
	return utils.GetSizeOrZero(path)
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
	return utils.GetSize(path)
}

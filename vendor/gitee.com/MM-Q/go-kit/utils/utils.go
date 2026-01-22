package utils

import (
	"strconv"
)

const (
	// 使用位运算常量，1024 = 1 << 10
	_KB = 1 << 10 // 1024
	_MB = 1 << 20 // 1048576
	_GB = 1 << 30 // 1073741824
	_TB = 1 << 40 // 1099511627776
	_PB = 1 << 50 // 1125899906842624
)

// 预定义单位数组，避免每次函数调用时重新创建
var units = [6]string{"B", "KB", "MB", "GB", "TB", "PB"}

// FormatBytes 将字节数转换为人类可读的带单位的字符串
// 用于将字节数格式化为易读的存储单位格式，支持B到PB的转换
//
// 参数:
//   - bytes: 字节数（int64类型）
//
// 返回:
//   - string: 格式化后的字符串，如 "1.23 KB", "456.78 MB", "2.34 GB" 等
func FormatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	// 处理负数
	if bytes < 0 {
		return "-" + FormatBytes(-bytes)
	}

	// 使用条件判断代替循环，提高性能
	switch {
	case bytes < _KB:
		return strconv.FormatInt(bytes, 10) + " B"
	case bytes < _MB:
		return formatWithUnit(bytes, _KB, 0)
	case bytes < _GB:
		return formatWithUnit(bytes, _MB, 1)
	case bytes < _TB:
		return formatWithUnit(bytes, _GB, 2)
	case bytes < _PB:
		return formatWithUnit(bytes, _TB, 3)
	default:
		return formatWithUnit(bytes, _PB, 4)
	}
}

// formatWithUnit 格式化字节数为指定单位
// 用于将字节数按指定除数转换为对应单位的格式化字符串
//
// 参数:
//   - bytes: 字节数（int64类型）
//   - divisor: 除数，用于计算单位
//   - unitIndex: 单位索引，对应units数组中的位置
//
// 返回:
//   - string: 格式化后的字符串，保留两位小数
func formatWithUnit(bytes, divisor int64, unitIndex int) string {
	// 计算整数部分和小数部分
	quotient := bytes / divisor
	remainder := bytes % divisor

	// 计算两位小数（乘以100后除以divisor再取整）
	decimal := (remainder * 100) / divisor

	// 构建结果字符串
	if decimal == 0 {
		return strconv.FormatInt(quotient, 10) + " " + units[unitIndex+1]
	}

	// 格式化小数部分，确保两位数显示
	var decimalStr string
	if decimal < 10 {
		decimalStr = "0" + strconv.FormatInt(decimal, 10)
	} else {
		decimalStr = strconv.FormatInt(decimal, 10)
	}

	return strconv.FormatInt(quotient, 10) + "." + decimalStr + " " + units[unitIndex+1]
}

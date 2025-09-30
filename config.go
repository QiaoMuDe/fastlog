/*
config.go - 日志配置管理模块
定义日志配置结构体及配置项的设置与获取方法，负责管理FastLog的所有可配置参数。
*/
package fastlog

import (
	"fmt"
	"path/filepath"
	"strings"
)

// FastLogConfig 定义一个配置结构体, 用于配置日志记录器
type FastLogConfig struct {
	LogDirName      string        // 日志目录路径
	LogFileName     string        // 日志文件名
	OutputToConsole bool          // 是否将日志输出到控制台
	OutputToFile    bool          // 是否将日志输出到文件
	LogLevel        LogLevel      // 日志级别
	LogFormat       LogFormatType // 日志格式选项
	Color           bool          // 是否启用终端颜色
	Bold            bool          // 是否启用终端字体加粗
	MaxSize         int           // 最大日志文件大小, 单位为MB, 默认10MB
	MaxAge          int           // 最大日志文件保留天数, 默认为0, 表示不做限制
	MaxFiles        int           // 最大日志文件保留数量, 默认为0, 表示不做限制
	LocalTime       bool          // 是否使用本地时间 默认使用UTC时间
	Compress        bool          // 是否启用日志文件压缩 默认不启用
}

// NewFastLogConfig 创建一个新的FastLogConfig实例, 用于配置日志记录器。
//
// 参数:
//   - logDirName: 日志目录名称, 默认为"applogs"。
//   - logFileName: 日志文件名称, 默认为"app.log"。
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
func NewFastLogConfig(logDirName string, logFileName string) *FastLogConfig {
	// 返回一个新的FastLogConfig实例
	return &FastLogConfig{
		LogDirName:      logDirName,         // 日志目录名称
		LogFileName:     logFileName,        // 日志文件名称
		OutputToConsole: true,               // 是否将日志输出到控制台
		OutputToFile:    true,               // 是否将日志输出到文件
		LogLevel:        INFO,               // 日志级别 默认INFO
		LogFormat:       Simple,             // 日志格式选项
		MaxSize:         defaultMaxFileSize, // 最大日志文件大小, 单位为MB, 默认10MB
		MaxAge:          0,                  // 最大日志文件保留天数, 默认为0, 表示不做限制
		MaxFiles:        0,                  // 最大日志文件保留数量, 默认为0, 表示不做限制
		LocalTime:       true,               // 是否使用本地时间 默认使用本地时间
		Compress:        false,              // 是否启用日志文件压缩 默认不启用
		Color:           true,               // 是否启用终端颜色
		Bold:            true,               // 是否启用终端字体加粗
	}
}

// ========================================================================
// 内部辅助函数
// ========================================================================

// validateConfig 验证配置并设置默认值
// 发现任何不合理的配置值都会panic，确保调用者提供正确配置
func (c *FastLogConfig) validateConfig() {
	// 配置对象不能为nil
	if c == nil {
		panic("FastLogConfig cannot be nil")
	}

	// ========================================================================
	// 第一步：设置所有默认值
	// ========================================================================

	// 设置日志级别默认值
	if c.LogLevel == 0 {
		c.LogLevel = INFO
	}

	// 设置日志格式默认值
	if c.LogFormat == 0 {
		c.LogFormat = Simple
	}

	// 设置文件大小默认值
	if c.MaxSize == 0 {
		c.MaxSize = defaultMaxFileSize
	}

	// 设置文件相关默认值（仅在启用文件输出时）
	if c.OutputToFile {
		if strings.TrimSpace(c.LogDirName) == "" {
			c.LogDirName = defaultLogDir
		}
		if strings.TrimSpace(c.LogFileName) == "" {
			c.LogFileName = defaultLogFileName
		}
	}

	// ========================================================================
	// 第二步：验证所有配置值
	// ========================================================================

	// 必须启用至少一种输出方式
	if !c.OutputToConsole && !c.OutputToFile {
		panic("at least one output method must be enabled: OutputToConsole or OutputToFile")
	}

	// 验证日志级别
	if c.LogLevel < DEBUG || c.LogLevel > NONE {
		panic(fmt.Sprintf("invalid LogLevel %d, must be %d-%d", c.LogLevel, DEBUG, NONE))
	}

	// 验证日志格式
	if c.LogFormat < Detailed || c.LogFormat > Custom {
		panic(fmt.Sprintf("invalid LogFormat %d, must be %d-%d", c.LogFormat, Detailed, Custom))
	}

	// 验证文件大小
	if c.MaxSize < 0 {
		panic("MaxSize cannot be negative")
	}
	if c.MaxSize > maxSingleFileSize {
		panic(fmt.Sprintf("MaxSize %d exceeds maximum %d MB", c.MaxSize, maxSingleFileSize))
	}

	// 验证保留天数
	if c.MaxAge < 0 {
		panic("MaxAge cannot be negative")
	}
	if c.MaxAge > maxRetentionDays {
		panic(fmt.Sprintf("MaxAge %d exceeds maximum %d days", c.MaxAge, maxRetentionDays))
	}

	// 验证保留文件数
	if c.MaxFiles < 0 {
		panic("MaxFiles cannot be negative")
	}
	if c.MaxFiles > maxRetentionFiles {
		panic(fmt.Sprintf("MaxFiles %d exceeds maximum %d files", c.MaxFiles, maxRetentionFiles))
	}

	// 验证文件输出相关配置
	if c.OutputToFile {
		// 验证路径安全性
		if strings.Contains(c.LogDirName, "..") {
			panic("LogDirName contains path traversal '..'")
		}
		if strings.Contains(c.LogFileName, "..") {
			panic("LogFileName contains path traversal '..'")
		}

		// 清理文件名
		c.LogDirName = cleanFileName(c.LogDirName)
		c.LogFileName = cleanFileName(c.LogFileName)
	}
}

// cleanFileName 清理文件名中的非法字符和格式问题
//
// 参数:
//   - filename: 原始文件名（可能包含路径）
//
// 返回:
//   - string: 清理后的文件名
func cleanFileName(filename string) string {
	// 处理空字符串
	if strings.TrimSpace(filename) == "" {
		return defaultLogFileName
	}

	// 1. 使用 filepath.Clean 进行路径规范化
	// 这会自动处理 "./"、多余的分隔符等问题
	cleaned := filepath.Clean(filename)

	// 2. 安全检查：移除上级目录引用
	if strings.Contains(cleaned, "..") {
		// 如果包含 ".."，只保留文件名部分
		cleaned = filepath.Base(cleaned)
	}

	// 3. 分离目录和文件名
	dir := filepath.Dir(cleaned)
	actualFileName := filepath.Base(cleaned)

	// 4. 清理文件名中的非法字符
	for _, char := range invalidFileChars {
		actualFileName = strings.ReplaceAll(actualFileName, char, charReplacement)
	}

	// 5. 移除文件名开头/结尾的点和空格
	actualFileName = strings.Trim(actualFileName, ". ")

	// 6. 处理文件名长度限制
	if len(actualFileName) > maxFileNameLength {
		ext := filepath.Ext(actualFileName)
		maxNameLen := maxFileNameLength - len(ext)
		if maxNameLen > 0 {
			actualFileName = actualFileName[:maxNameLen] + ext
		} else {
			actualFileName = actualFileName[:maxFileNameLength]
		}
	}

	// 7. 确保文件名不为空
	if actualFileName == "" {
		actualFileName = defaultLogFileName
	}

	// 8. 重新组合路径
	if dir == "." {
		cleaned = actualFileName
	} else {
		cleaned = filepath.Join(dir, actualFileName)
	}

	// 9. 最终路径长度检查
	if len(cleaned) > maxPathLength {
		cleaned = cleaned[:maxPathLength-truncateReserve] + truncatedSuffix
	}

	return cleaned
}

/*
config.go - 日志配置管理模块
定义日志配置结构体及配置项的设置与获取方法，负责管理FastLog的所有可配置参数。
*/
package fastlog

import (
	"path/filepath"
	"strings"
	"time"
)

// FastLogConfig 定义一个配置结构体, 用于配置日志记录器
type FastLogConfig struct {
	LogDirName      string        // 日志目录路径
	LogFileName     string        // 日志文件名
	OutputToConsole bool          // 是否将日志输出到控制台
	OutputToFile    bool          // 是否将日志输出到文件
	FlushInterval   time.Duration // 刷新间隔, 单位为time.Duration
	LogLevel        LogLevel      // 日志级别
	ChanIntSize     int           // 通道大小 默认10000
	LogFormat       LogFormatType // 日志格式选项
	NoColor         bool          // 是否禁用终端颜色
	NoBold          bool          // 是否禁用终端字体加粗
	MaxLogFileSize  int           // 最大日志文件大小, 单位为MB, 默认10MB
	MaxLogAge       int           // 最大日志文件保留天数, 默认为0, 表示不做限制
	MaxLogBackups   int           // 最大日志文件保留数量, 默认为0, 表示不做限制
	IsLocalTime     bool          // 是否使用本地时间 默认使用UTC时间
	EnableCompress  bool          // 是否启用日志文件压缩 默认不启用
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
		LogDirName:      logDirName,             // 日志目录名称
		LogFileName:     logFileName,            // 日志文件名称
		OutputToConsole: true,                   // 是否将日志输出到控制台
		OutputToFile:    true,                   // 是否将日志输出到文件
		LogLevel:        INFO,                   // 日志级别 默认INFO
		ChanIntSize:     10000,                  // 通道大小 增加到10000
		FlushInterval:   500 * time.Millisecond, // 刷新间隔 缩短到500毫秒
		LogFormat:       Detailed,               // 日志格式选项
		MaxLogFileSize:  10,                     // 最大日志文件大小, 单位为MB, 默认10MB
		MaxLogAge:       0,                      // 最大日志文件保留天数, 默认为0, 表示不做限制
		MaxLogBackups:   0,                      // 最大日志文件保留数量, 默认为0, 表示不做限制
		IsLocalTime:     false,                  // 是否使用本地时间 默认使用UTC时间
		EnableCompress:  false,                  // 是否启用日志文件压缩 默认不启用
		NoColor:         false,                  // 是否禁用终端颜色
		NoBold:          false,                  // 是否禁用终端字体加粗
	}
}

// fixFinalConfig 最终配置修正函数 - 在NewFastLog开始时调用
// 负责修正所有不合理的配置值, 确保系统稳定运行
func (c *FastLogConfig) fixFinalConfig() {
	// 判断是否为文件输出模式
	outputToFile := c.OutputToFile

	// 只在文件输出模式下修正
	if outputToFile {
		// 1. 修正基本字符串字段
		if c.LogDirName == "" {
			c.LogDirName = "logs"
		}

		if c.LogFileName == "" {
			c.LogFileName = "app.log"
		}

		// 2. 清理文件名中的非法字符
		originalDir := c.LogDirName
		cleanedDir := cleanFileName(originalDir)
		if originalDir != cleanedDir {
			c.LogDirName = cleanedDir
		}

		originalFile := c.LogFileName
		cleanedFile := cleanFileName(originalFile)
		if originalFile != cleanedFile {
			c.LogFileName = cleanedFile
		}
	}

	// 3. 修正通道大小
	if c.ChanIntSize <= 0 {
		c.ChanIntSize = 10000
	} else if c.ChanIntSize > 100000 {
		c.ChanIntSize = 100000
	}

	// 4. 修正刷新间隔
	if c.FlushInterval <= 0 {
		c.FlushInterval = 500 * time.Millisecond
	} else if c.FlushInterval < 10*time.Millisecond {
		c.FlushInterval = 10 * time.Millisecond
	} else if c.FlushInterval > 30*time.Second {
		c.FlushInterval = 30 * time.Second
	}

	// 5. 修正最大文件大小
	if c.MaxLogFileSize <= 0 {
		c.MaxLogFileSize = 10
	} else if c.MaxLogFileSize > 1000 {
		c.MaxLogFileSize = 1000
	}

	// 6. 修正最大保留天数
	if c.MaxLogAge < 0 {
		c.MaxLogAge = 0
	} else if c.MaxLogAge > 3650 {
		c.MaxLogAge = 3650
	}

	// 7. 修正最大备份数量
	if c.MaxLogBackups < 0 {
		c.MaxLogBackups = 0
	} else if c.MaxLogBackups > 1000 {
		c.MaxLogBackups = 1000
	}

	// 8. 修正日志级别
	if c.LogLevel < DEBUG || c.LogLevel > NONE {
		c.LogLevel = INFO
	}

	// 9. 修正日志格式
	if c.LogFormat < Detailed || c.LogFormat > Custom {
		c.LogFormat = Detailed
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
		return "default.log"
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
		actualFileName = strings.ReplaceAll(actualFileName, char, "_")
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
		actualFileName = "app.log"
	}

	// 8. 重新组合路径
	if dir == "." {
		cleaned = actualFileName
	} else {
		cleaned = filepath.Join(dir, actualFileName)
	}

	// 9. 最终路径长度检查
	if len(cleaned) > maxPathLength {
		cleaned = cleaned[:maxPathLength-10] + "_truncated"
	}

	return cleaned
}

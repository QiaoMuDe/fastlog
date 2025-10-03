/*
config.go - 日志配置管理模块
定义日志配置结构体及配置项的设置与获取方法，负责管理FastLog的所有可配置参数。
*/
package fastlog

import (
	"fmt"
	"strings"
	"time"
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
	MaxBufferSize   int           // 缓冲区大小, 单位为字节, 默认64KB
	MaxWriteCount   int           // 最大写入次数, 默认500次
	FlushInterval   time.Duration // 刷新间隔, 默认1秒
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
		LogDirName:      logDirName,           // 日志目录名称
		LogFileName:     logFileName,          // 日志文件名称
		OutputToConsole: true,                 // 是否将日志输出到控制台
		OutputToFile:    true,                 // 是否将日志输出到文件
		LogLevel:        INFO,                 // 日志级别 默认INFO
		LogFormat:       Simple,               // 日志格式选项
		MaxSize:         defaultMaxFileSize,   // 最大日志文件大小, 单位为MB, 默认10MB
		MaxAge:          0,                    // 最大日志文件保留天数, 默认为0, 表示不做限制
		MaxFiles:        0,                    // 最大日志文件保留数量, 默认为0, 表示不做限制
		LocalTime:       true,                 // 是否使用本地时间 默认使用本地时间
		Compress:        false,                // 是否启用日志文件压缩 默认不启用
		Color:           true,                 // 是否启用终端颜色
		Bold:            true,                 // 是否启用终端字体加粗
		MaxBufferSize:   defaultMaxBufferSize, // 缓冲区大小, 单位为字节, 默认64KB
		MaxWriteCount:   defaultMaxWriteCount, // 最大写入次数, 默认500次
		FlushInterval:   defaultFlushInterval, // 刷新间隔, 默认1秒
	}
}

// DevConfig 创建一个开发环境下的FastLogConfig实例
//
// 参数:
//   - logDirName: 日志目录名
//   - logFileName: 日志文件名
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
//
// 特性:
//   - 启用详细日志格式
//   - 设置日志级别为DEBUG
//   - 设置最大日志文件保留数量为5
//   - 设置最大日志文件保留天数为7天
func DevConfig(logDirName string, logFileName string) *FastLogConfig {
	// 创建一个新的FastLogConfig实例
	cfg := NewFastLogConfig(logDirName, logFileName)
	cfg.LogLevel = DEBUG     // 设置日志级别为DEBUG
	cfg.LogFormat = Detailed // 设置日志格式为详细格式
	cfg.MaxFiles = 5         // 设置最大日志文件保留数量为5
	cfg.MaxAge = 7           // 设置最大日志文件保留天数为7天
	return cfg
}

// ProdConfig 创建一个生产环境下的FastLogConfig实例
//
// 参数:
//   - logDirName: 日志目录名
//   - logFileName: 日志文件名
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
//
// 特性:
//   - 启用日志文件压缩
//   - 禁用控制台输出
//   - 设置最大日志文件保留天数为30天
//   - 设置最大日志文件保留数量为24个
func ProdConfig(logDirName string, logFileName string) *FastLogConfig {
	cfg := NewFastLogConfig(logDirName, logFileName)
	cfg.MaxAge = 30             // 设置最大日志文件保留天数为30天
	cfg.MaxFiles = 24           // 设置最大日志文件保留数量为24个
	cfg.Compress = true         // 启用日志文件压缩
	cfg.OutputToConsole = false // 禁用控制台输出
	return cfg
}

// ConsoleConfig 创建一个控制台环境下的FastLogConfig实例
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
//
// 特性:
//   - 禁用文件输出
//   - 设置日志级别为DEBUG
func ConsoleConfig() *FastLogConfig {
	cfg := NewFastLogConfig("", "")
	cfg.OutputToFile = false // 禁用文件输出
	cfg.LogLevel = DEBUG     // 设置日志级别为DEBUG
	return cfg
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

	// 仅进行只读校验，不做任何字段修改

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

	// 验证保留天数
	if c.MaxAge < 0 {
		panic("MaxAge cannot be negative")
	}

	// 验证保留文件数
	if c.MaxFiles < 0 {
		panic("MaxFiles cannot be negative")
	}

	// 文件输出相关校验（只校验不更改）
	if c.OutputToFile {
		// 必填项
		if strings.TrimSpace(c.LogDirName) == "" {
			panic("LogDirName cannot be empty when OutputToFile is enabled")
		}
		if strings.TrimSpace(c.LogFileName) == "" {
			panic("LogFileName cannot be empty when OutputToFile is enabled")
		}
		// 路径穿越检测
		if strings.Contains(c.LogDirName, "..") {
			panic("LogDirName contains path traversal '..'")
		}
		if strings.Contains(c.LogFileName, "..") {
			panic("LogFileName contains path traversal '..'")
		}
	}
}

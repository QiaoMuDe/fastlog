package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"gitee.com/MM-Q/fastlog/internal/types"
	"gitee.com/MM-Q/logrotatex"
)

// FastLogConfig 定义一个配置结构体, 用于配置日志记录器
type FastLogConfig struct {
	LogDirName      string              // 日志目录路径
	LogFileName     string              // 日志文件名
	OutputToConsole bool                // 是否将日志输出到控制台
	OutputToFile    bool                // 是否将日志输出到文件
	LogLevel        types.LogLevel      // 日志级别
	LogFormat       types.LogFormatType // 日志格式选项
	Color           bool                // 是否启用终端颜色
	Bold            bool                // 是否启用终端字体加粗
	MaxSize         int                 // 最大日志文件大小, 单位为MB, 默认10MB
	MaxAge          int                 // 最大日志文件保留天数, 默认为0, 表示不做限制
	MaxFiles        int                 // 最大日志文件保留数量, 默认为0, 表示不做限制
	LocalTime       bool                // 是否使用本地时间 默认使用UTC时间
	Compress        bool                // 是否启用日志文件压缩 默认不启用
	MaxBufferSize   int                 // 缓冲区大小, 单位为字节, 默认64KB
	MaxWriteCount   int                 // 最大写入次数, 默认500次
	FlushInterval   time.Duration       // 刷新间隔, 默认1秒
	Async           bool                // 是否异步清理日志, 默认同步清理
	CallerInfo      bool                // 是否获取调用者信息, 默认不获取
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
		LogDirName:      logDirName,                 // 日志目录名称
		LogFileName:     logFileName,                // 日志文件名称
		OutputToConsole: true,                       // 是否将日志输出到控制台
		OutputToFile:    true,                       // 是否将日志输出到文件
		LogLevel:        types.INFO,                 // 日志级别 默认INFO
		LogFormat:       types.Def,                  // 日志格式选项
		MaxSize:         types.DefaultMaxFileSize,   // 最大日志文件大小, 单位为MB, 默认10MB
		MaxAge:          0,                          // 最大日志文件保留天数, 默认为0, 表示不做限制
		MaxFiles:        0,                          // 最大日志文件保留数量, 默认为0, 表示不做限制
		LocalTime:       true,                       // 是否使用本地时间 默认使用本地时间
		Compress:        false,                      // 是否启用日志文件压缩 默认不启用
		Color:           true,                       // 是否启用终端颜色
		Bold:            true,                       // 是否启用终端字体加粗
		MaxBufferSize:   types.DefaultMaxBufferSize, // 缓冲区大小, 单位为字节, 默认64KB
		MaxWriteCount:   types.DefaultMaxWriteCount, // 最大写入次数, 默认500次
		FlushInterval:   types.DefaultFlushInterval, // 刷新间隔, 默认1秒
		Async:           false,                      // 是否异步清理日志, 默认同步清理
		CallerInfo:      false,                      // 是否获取调用者信息, 默认不获取
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
	cfg.LogLevel = types.DEBUG // 设置日志级别为DEBUG
	cfg.LogFormat = types.Def  // 设置日志格式为默认格式
	cfg.MaxFiles = 5           // 设置最大日志文件保留数量为5
	cfg.MaxAge = 7             // 设置最大日志文件保留天数为7天
	cfg.CallerInfo = true      // 启用调用者信息
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
	cfg.LogFormat = types.Json  // 设置日志格式为json格式
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
	cfg.OutputToFile = false   // 禁用文件输出
	cfg.LogLevel = types.DEBUG // 设置日志级别为DEBUG
	cfg.CallerInfo = true      // 启用调用者信息
	cfg.LogFormat = types.Def  // 设置日志格式为默认格式
	return cfg
}

// ========================================================================
// 内部辅助函数
// ========================================================================

// ValidateConfig 验证配置并设置默认值
// 发现任何不合理的配置值都会panic，确保调用者提供正确配置
func (c *FastLogConfig) ValidateConfig() {
	// 配置对象不能为nil
	if c == nil {
		panic("FastLogConfig cannot be nil")
	}

	// 必须启用至少一种输出方式
	if !c.OutputToConsole && !c.OutputToFile {
		panic("at least one output method must be enabled: OutputToConsole or OutputToFile")
	}

	// 验证日志格式
	if c.LogFormat < types.Def || c.LogFormat > types.Custom {
		panic(fmt.Sprintf("invalid LogFormat %d, must be %d-%d", c.LogFormat, types.Def, types.Custom))
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

// Clone 创建一个FastLogConfig的副本
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针，该实例是当前实例的副本。
func (c *FastLogConfig) Clone() *FastLogConfig {
	return &FastLogConfig{
		LogDirName:      c.LogDirName,
		LogFileName:     c.LogFileName,
		LogLevel:        c.LogLevel,
		LogFormat:       c.LogFormat,
		Color:           c.Color,
		Bold:            c.Bold,
		OutputToFile:    c.OutputToFile,
		OutputToConsole: c.OutputToConsole,
		MaxSize:         c.MaxSize,
		MaxAge:          c.MaxAge,
		MaxFiles:        c.MaxFiles,
		Compress:        c.Compress,
		Async:           c.Async,
		LocalTime:       c.LocalTime,
		MaxBufferSize:   c.MaxBufferSize,
		MaxWriteCount:   c.MaxWriteCount,
		FlushInterval:   c.FlushInterval,
		CallerInfo:      c.CallerInfo,
	}
}

// CreateBufferedWriter 根据配置创建一个带缓冲的文件写入器
//
// 参数:
//   - cfg: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。
//
// 返回值:
//   - *logrotatex.BufferedWriter: 一个指向BufferedWriter实例的指针。
func CreateBufferedWriter(cfg *FastLogConfig) *logrotatex.BufferedWriter {
	// 如果不允许将日志输出到文件, 则返回nil
	if !cfg.OutputToFile {
		return nil
	}

	// 拼接日志文件路径
	logFilePath := filepath.Join(cfg.LogDirName, cfg.LogFileName)

	// 初始化日志文件切割器
	logger := logrotatex.NewLogRotateX(logFilePath) // 初始化日志文件切割器
	logger.MaxSize = cfg.MaxSize                    // 最大日志文件大小, 单位为MB
	logger.MaxAge = cfg.MaxAge                      // 最大日志文件保留天数
	logger.MaxFiles = cfg.MaxFiles                  // 最大日志文件保留数量
	logger.Compress = cfg.Compress                  // 是否启用日志文件压缩
	logger.LocalTime = cfg.LocalTime                // 是否使用本地时间
	logger.Async = cfg.Async                        // 是否异步清理日志

	// 初始化缓冲区配置
	bufCfg := logrotatex.DefBufCfg()
	bufCfg.FlushInterval = cfg.FlushInterval // 刷新间隔, 单位为秒, 默认为0, 表示不做限制
	bufCfg.MaxBufferSize = cfg.MaxBufferSize // 缓冲区最大容量, 单位为字节
	bufCfg.MaxWriteCount = cfg.MaxWriteCount // 最大写入次数, 默认为0, 表示不做限制

	// 创建带缓冲的批量写入器，嵌入日志切割器
	return logrotatex.NewBufferedWriter(logger, bufCfg)
}

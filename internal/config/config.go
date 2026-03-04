package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitee.com/MM-Q/comprx"
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
	MaxBufferSize   int                 // 缓冲区大小, 单位为字节, 默认256KB
	FlushInterval   time.Duration       // 刷新间隔, 默认1秒, 最低为500毫秒
	Async           bool                // 是否异步清理日志, 默认同步清理
	CallerInfo      bool                // 是否获取调用者信息, 默认不获取

	// DateDirLayout 决定是否启用按日期目录存放轮转后的日志。
	// true: 轮转后的日志存放在 YYYY-MM-DD/ 目录下
	// false: 轮转后的日志存放在当前目录下 (默认)
	DateDirLayout bool `json:"datedirlayout" yaml:"datedirlayout"`

	// RotateByDay 决定是否启用按天轮转。
	// true: 每天自动轮转一次 (跨天时触发)
	// false: 只按文件大小轮转 (默认)
	RotateByDay bool `json:"rotatebyday" yaml:"rotatebyday"`

	// CompressType 压缩类型, 默认为: comprx.CompressTypeZip
	//
	// 支持的压缩格式：
	//   - comprx.CompressTypeZip: zip 压缩格式
	//   - comprx.CompressTypeTar: tar 压缩格式
	//   - comprx.CompressTypeTgz: tgz 压缩格式
	//   - comprx.CompressTypeTarGz: tar.gz 压缩格式
	//   - comprx.CompressTypeGz: gz 压缩格式
	//   - comprx.CompressTypeBz2: bz2 压缩格式
	//   - comprx.CompressTypeBzip2: bzip2 压缩格式
	//   - comprx.CompressTypeZlib: zlib 压缩格式
	CompressType comprx.CompressType `json:"compress_type" yaml:"compress_type"`

	// BufferedWrite 是否使用带缓冲的批量写入器
	//  - true: 使用带缓冲的批量写入器（默认，高性能）
	//  - false: 使用普通文件句柄直接写入（低延迟）
	BufferedWrite bool `json:"buffered_write" yaml:"buffered_write"`
}

// Default 返回一个默认的FastLogConfig实例, 用于配置日志记录器。
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
//
// 默认配置:
//   - 日志目录: "logs"
//   - 日志文件名: "app.log"
//   - 日志级别: INFO
//   - 日志格式: Def
//   - 最大日志文件大小: 10MB
//   - 最大日志文件保留天数: 0 (不做限制)
//   - 最大日志文件保留数量: 0 (不做限制)
//   - 是否使用本地时间: true
//   - 是否启用日志文件压缩: false
//   - 是否启用终端颜色: true
//   - 是否启用终端字体加粗: true
//   - 缓冲区大小: 256KB
//   - 刷新间隔: 1秒
//   - 是否异步清理日志: false
//   - 是否获取调用者信息: false
//   - 是否将日志输出到控制台: true
//   - 是否将日志输出到文件: true
//   - 是否启用按日期目录存放轮转后的日志: true
//   - 是否启用按天轮转: true
//   - 压缩类型: comprx.CompressTypeZip
//   - 是否使用带缓冲的批量写入器: true (默认)
func Default() *FastLogConfig {
	return NewFastLogConfig("logs", "app.log")
}

// NewFastLogConfig 创建一个新的FastLogConfig实例, 用于配置日志记录器。
//
// 参数:
//   - logDirName: 日志目录名称
//   - logFileName: 日志文件名称
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
//
// 默认配置:
//   - 日志级别: INFO
//   - 日志格式: Def
//   - 最大日志文件大小: 10MB
//   - 最大日志文件保留天数: 0 (不做限制)
//   - 最大日志文件保留数量: 0 (不做限制)
//   - 是否使用本地时间: true
//   - 是否启用日志文件压缩: false
//   - 是否启用终端颜色: true
//   - 是否启用终端字体加粗: true
//   - 缓冲区大小: 256KB
//   - 刷新间隔: 1秒
//   - 是否异步清理日志: false
//   - 是否获取调用者信息: false
//   - 是否将日志输出到控制台: true
//   - 是否将日志输出到文件: true
//   - 是否启用按日期目录存放轮转后的日志: true
//   - 是否启用按天轮转: true
//   - 压缩类型: comprx.CompressTypeZip
//   - 是否使用带缓冲的批量写入器: true (默认)
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
		MaxBufferSize:   types.DefaultMaxBufferSize, // 缓冲区大小, 单位为字节, 默认256KB
		FlushInterval:   types.DefaultFlushInterval, // 刷新间隔, 默认1秒
		Async:           false,                      // 是否异步清理日志, 默认同步清理
		CallerInfo:      false,                      // 是否获取调用者信息, 默认不获取
		DateDirLayout:   true,                       // 是否启用按日期目录存放轮转后的日志, 默认启用
		RotateByDay:     true,                       // 是否启用按天轮转, 默认启用
		CompressType:    comprx.CompressTypeZip,     // 压缩类型, 默认zip压缩格式
		BufferedWrite:   true,                       // 是否使用带缓冲的批量写入器, 默认使用带缓冲的批量写入器
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
		panic("maxSize cannot be negative")
	}

	// 验证保留天数
	if c.MaxAge < 0 {
		panic("maxAge cannot be negative")
	}

	// 验证保留文件数
	if c.MaxFiles < 0 {
		panic("maxFiles cannot be negative")
	}

	// 文件输出相关校验（只校验不更改）
	if c.OutputToFile {
		// 必填项校验
		if strings.TrimSpace(c.LogDirName) == "" {
			panic("logDirName cannot be empty when OutputToFile is enabled")
		}
		if strings.TrimSpace(c.LogFileName) == "" {
			panic("logFileName cannot be empty when OutputToFile is enabled")
		}
		// 路径穿越检测
		if strings.Contains(c.LogDirName, "..") {
			panic("logDirName contains path traversal '..'")
		}
		if strings.Contains(c.LogFileName, "..") {
			panic("logFileName contains path traversal '..'")
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
		FlushInterval:   c.FlushInterval,
		CallerInfo:      c.CallerInfo,
		DateDirLayout:   c.DateDirLayout,
		RotateByDay:     c.RotateByDay,
		CompressType:    c.CompressType,
		BufferedWrite:   c.BufferedWrite,
	}
}

// CreateWriter 根据配置创建一个文件写入器
//
// 参数:
//   - cfg: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。
//
// 返回值:
//   - io.WriteCloser: 一个指向WriteCloser实例的指针，用于写入日志到文件。
func CreateWriter(cfg *FastLogConfig) io.WriteCloser {
	// 如果不允许将日志输出到文件, 则返回nil
	if !cfg.OutputToFile {
		return nil
	}

	// 拼接日志文件路径
	logFilePath := filepath.Join(cfg.LogDirName, cfg.LogFileName)

	// 如果不使用缓冲写入，则直接打开文件
	if !cfg.BufferedWrite {
		// 确保日志目录存在
		if err := os.MkdirAll(cfg.LogDirName, 0755); err != nil {
			panic(fmt.Sprintf("failed to create log directory: %v", err))
		}

		// 打开日志文件
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(fmt.Sprintf("failed to open log file: %v", err))
		}
		return file
	}

	// 默认使用带有缓冲的批量写入器
	logger := &logrotatex.LogRotateX{
		LogFilePath:   logFilePath,       // 日志文件路径
		MaxSize:       cfg.MaxSize,       // 最大日志文件大小, 单位为MB
		MaxAge:        cfg.MaxAge,        // 最大日志文件保留天数
		MaxFiles:      cfg.MaxFiles,      // 最大日志文件保留数量
		Compress:      cfg.Compress,      // 是否启用日志文件压缩
		LocalTime:     cfg.LocalTime,     // 是否使用本地时间
		Async:         cfg.Async,         // 是否异步清理日志
		DateDirLayout: cfg.DateDirLayout, // 是否启用按日期目录存放轮转后的日志
		RotateByDay:   cfg.RotateByDay,   // 是否启用按天轮转
		CompressType:  cfg.CompressType,  // 压缩类型
	}

	// 初始化缓冲区配置
	bufCfg := &logrotatex.BufCfg{
		FlushInterval: cfg.FlushInterval, // 刷新间隔, 单位为秒, 默认为1秒, 最低为500毫秒
		MaxBufferSize: cfg.MaxBufferSize, // 缓冲区最大容量, 单位为字节
	}

	// 创建带缓冲的批量写入器，嵌入日志切割器
	return logrotatex.NewBufferedWriter(logger, bufCfg)
}

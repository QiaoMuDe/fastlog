package fastlog

import (
	"gitee.com/MM-Q/fastlog/internal/config"
	"gitee.com/MM-Q/fastlog/internal/logger/flog"
	"gitee.com/MM-Q/fastlog/internal/types"
)

// LogLevel 定义为位掩码类型，每一位代表一个日志级别
//
// 支持的日志级别:
//   - DEBUG: 调试级别
//   - INFO: 信息级别
//   - WARN: 警告级别
//   - ERROR: 错误级别
//   - FATAL: 致命错误级别
//   - NONE: 不记录任何日志级别
type LogLevel = types.LogLevel

// LogLevel 日志级别
const (
	DEBUG = types.DEBUG // 调试级别
	INFO  = types.INFO  // 信息级别
	WARN  = types.WARN  // 警告级别
	ERROR = types.ERROR // 错误级别
	FATAL = types.FATAL // 致命错误级别
	NONE  = types.NONE  // 不记录任何日志级别
)

// LogFormatType 日志格式选项
//
// 格式:
//   - Detailed: 详细格式
//   - Json: json格式
//   - Timestamp: 时间格式
//   - KVFmt: 键值对格式
//   - LogFmt: logfmt格式
//   - Custom: 自定义格式
type LogFormatType = types.LogFormatType

const (
	Def       = types.Def       // 默认格式
	Json      = types.Json      // json格式
	Timestamp = types.Timestamp // 时间格式
	KVFmt     = types.KVFmt     // 键值对格式
	LogFmt    = types.LogFmt    // 日志格式
	Custom    = types.Custom    // 自定义格式
)

// FastLogConfig 定义一个配置结构体, 用于配置日志记录器
type FastLogConfig = config.FastLogConfig

// NewFastLogConfig 创建一个新的FastLogConfig实例, 用于配置日志记录器。
//
// 参数:
//   - logDirName: 日志目录名称, 默认为"applogs"。
//   - logFileName: 日志文件名称, 默认为"app.log"。
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
var NewFastLogConfig = config.NewFastLogConfig

// NewCfg 是 NewFastLogConfig 的简写, 用于创建一个新的FastLogConfig实例。
//
// 参数:
//   - logDirName: 日志目录名称, 默认为"applogs"。
//   - logFileName: 日志文件名称, 默认为"app.log"。
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
var NewCfg = config.NewFastLogConfig

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
var DevConfig = config.DevConfig

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
var ProdConfig = config.ProdConfig

// ConsoleConfig 创建一个控制台环境下的FastLogConfig实例
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
//
// 特性:
//   - 禁用文件输出
//   - 设置日志级别为DEBUG
var ConsoleConfig = config.ConsoleConfig

// FLog 是一个高性能的日志记录器, 支持键值对风格的使用和标准库fmt类似的使用,
// 同时提供了丰富的配置选项, 如日志级别、输出格式、日志轮转等。
type FLog = flog.FLog

// NewFLog 创建一个新的FLog实例, 用于记录日志。
//
// 参数:
//   - config: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。
//
// 返回值:
//   - *FLog: 一个指向FLog实例的指针。
var NewFLog = flog.NewFLog

// String 添加字符串字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为字符串。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var String = flog.String

// Int 添加整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var Int = flog.Int

// Int64 添加64位整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为64位整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var Int64 = flog.Int64

// Float64 添加64位浮点数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为64位浮点数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var Float64 = flog.Float64

// Bool 添加布尔字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为布尔值。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var Bool = flog.Bool

// Time 添加时间字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为time.Time类型。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var Time = flog.Time

// Duration 添加持续时间字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为time.Duration类型。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var Duration = flog.Duration

// Uint64 添加64位无符号整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为64位无符号整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var Uint64 = flog.Uint64

// Uint32 添加32位无符号整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为32位无符号整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var Uint32 = flog.Uint32

// Uint16 添加16位无符号整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为16位无符号整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var Uint16 = flog.Uint16

// Uint8 添加8位无符号整数字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为8位无符号整数。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var Uint8 = flog.Uint8

// Error 添加错误字段
//
// 参数:
//   - key: 字段的键名, 不能为空字符串。
//   - value: 字段值, 必须为error类型。
//
// 返回值:
//   - *Field: 一个指向 Field 实例的指针。
var Error = flog.Error

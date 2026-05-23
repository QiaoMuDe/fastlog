package fastlog

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"gitee.com/MM-Q/comprx"
	"gitee.com/MM-Q/logrotatex"
)

// DefaultTimeFormat 默认时间格式
const DefaultTimeFormat = time.RFC3339

// NewConfig 创建一个默认配置实例
//
// 默认配置同时输出到终端和文件, 适用于开发和测试环境。
//
// 默认配置详情:
//   - Level: INFO - 日志级别为信息级别
//   - Formatter: Def{} - 使用默认格式输出
//   - Caller: false - 不记录调用者信息
//   - Fields: []Field{} - 无预设字段
//   - SamplerTick: 10s - 采样窗口为10秒
//   - SamplerInitial: 3 - 前3条日志放行
//   - SamplerThereafter: 10 - 之后每10条放行1条
//   - LevelRouter: false - 不启用级别路由
//   - OutputConsole: true - 输出到终端
//   - NoColor: false - 启用彩色输出
//   - OutputFile: true - 输出到文件
//   - LogPath: 参数指定 - 日志文件路径由参数指定
//   - Async: false - 不启用异步清理
//   - MaxSize: 20MB - 单文件最大20MB
//   - MaxFiles: 24 - 保留24个历史文件
//   - MaxAge: 7天 - 保留7天日志
//   - Compress: false - 不压缩历史日志
//   - CompressType: Gz - 压缩类型为gzip
//   - LocalTime: true - 使用本地时间命名
//   - DateDirLayout: true - 按日期目录存放
//   - RotateByDay: true - 按天轮转文件
//   - MaxBufferSize: 256KB - 缓冲区256KB
//   - SyncInterval: 1s - 每秒同步
//   - TimeFormat: RFC3339 - 时间格式为RFC3339
//
// 参数:
//   - logPath: 日志文件路径
//
// 返回:
//   - *Config: 配置实例, 所有字段都设置了默认值
func NewConfig(logPath string) *Config {
	return &Config{
		// 基础日志配置
		Level:             INFO,                     // 日志级别, 零值默认 INFO
		Formatter:         Def{},                    // 日志格式化器, 零值默认 Def
		Caller:            false,                    // 是否记录调用者信息 (文件:函数:行号)
		Fields:            []Field{},                // 预设字段
		SamplerTick:       DefaultSamplerTick,       // 每 DefaultSamplerTick 一个窗口
		SamplerInitial:    DefaultSamplerInitial,    // 前 DefaultSamplerInitial 条放行
		SamplerThereafter: DefaultSamplerThereafter, // 之后每 DefaultSamplerThereafter 条放行 1 条
		LevelRouter:       false,                    // 是否启用级别路由, 零值默认 false

		// 终端输出配置
		OutputConsole: true,  // 是否输出到终端 (彩色自动检测)
		NoColor:       false, // 设为 true 时禁用终端彩色输出, 仅当 OutputConsole=true 时生效

		// 文件输出配置
		OutputFile:    true,                  // 是否输出到文件
		LogPath:       logPath,               // 日志文件路径
		Async:         false,                 // 是否启用异步清理 (单协程、合并触发) , 关闭文件时后台清理旧日志
		MaxSize:       20,                    // 单文件最大大小 (MB) , 超过后触发轮转。零值默认 10MB。
		MaxFiles:      24,                    // 保留的历史日志文件数, 超过的旧文件自动删除。零值表示不限制。
		MaxAge:        7,                     // 保留天数, 超过指定天数的旧文件自动清理。零值表示不限制。
		Compress:      false,                 // 是否压缩历史日志文件
		CompressType:  comprx.CompressTypeGz, // 压缩类型, 默认: comprx.CompressTypeGz
		LocalTime:     true,                  // 是否使用本地时间命名轮转文件
		DateDirLayout: true,                  // 是否按日期目录存放轮转文件
		RotateByDay:   true,                  // 是否按天轮转文件

		// 缓冲写入配置
		MaxBufferSize: 256 * 1024,      // 缓冲区最大大小 (字节) , 零值默认 256KB。
		SyncInterval:  1 * time.Second, // 同步间隔, 零值默认 1秒。

		// 时间格式
		TimeFormat: DefaultTimeFormat, // 时间格式, 默认: time.RFC3339
	}
}

// Default 创建一个默认配置实例
//
// 默认配置同时输出到终端和文件, 适用于开发和测试环境。
//
// 默认配置详情:
//   - Level: INFO - 日志级别为信息级别
//   - Formatter: Def{} - 使用默认格式输出
//   - Caller: false - 不记录调用者信息
//   - Fields: []Field{} - 无预设字段
//   - SamplerTick: 10s - 采样窗口为10秒
//   - SamplerInitial: 3 - 前3条日志放行
//   - SamplerThereafter: 10 - 之后每10条放行1条
//   - LevelRouter: false - 不启用级别路由
//   - OutputConsole: true - 输出到终端
//   - NoColor: false - 启用彩色输出
//   - OutputFile: true - 输出到文件
//   - LogPath: 参数指定 - 日志文件路径由参数指定
//   - Async: false - 不启用异步清理
//   - MaxSize: 20MB - 单文件最大20MB
//   - MaxFiles: 24 - 保留24个历史文件
//   - MaxAge: 7天 - 保留7天日志
//   - Compress: false - 不压缩历史日志
//   - CompressType: Gz - 压缩类型为gzip
//   - LocalTime: true - 使用本地时间命名
//   - DateDirLayout: true - 按日期目录存放
//   - RotateByDay: true - 按天轮转文件
//   - MaxBufferSize: 256KB - 缓冲区256KB
//   - SyncInterval: 1s - 每秒同步
//   - TimeFormat: RFC3339 - 时间格式为RFC3339
//
// 返回:
//   - *Config: 配置实例, 所有字段都设置了默认值
func Default() *Config {
	return NewConfig("logs/app.log")
}

// Dev 创建开发环境配置
//
// 基于 NewConfig 修改以下配置:
//   - Level: DEBUG (更多信息便于调试)
//   - Caller: true (开启调用者信息, 方便定位代码)
//   - SamplerTick: 0 (关闭采样, 保留所有日志)
//   - MaxSize: 10MB (小文件, 便于查看)
//   - Compress: false (不压缩, 快速写入)
//   - RotateByDay: false (不按天轮转)
//
// 参数:
//   - logPath: 日志文件路径
//
// 返回:
//   - *Config: 开发环境配置实例
func Dev(logPath string) *Config {
	cfg := NewConfig(logPath)
	cfg.Level = DEBUG
	cfg.Caller = true
	cfg.SamplerTick = 0
	cfg.SamplerInitial = 0
	cfg.SamplerThereafter = 0
	cfg.MaxSize = 10
	cfg.Compress = false
	cfg.RotateByDay = false
	return cfg
}

// Prod 创建生产环境配置
//
// 基于 NewConfig 修改以下配置:
//   - Level: WARN (只记录警告及以上, 减少日志量)
//   - OutputConsole: false (只输出到文件, 无终端开销)
//   - Async: true (异步清理, 性能更好)
//   - MaxSize: 100MB (大文件, 减少轮转频率)
//   - MaxFiles: 14 (保留2周)
//   - MaxAge: 14 (保留14天)
//   - Compress: true (开启压缩, 节省磁盘)
//   - LevelRouter: true (启用级别路由, 便于快速定位错误)
//
// 参数:
//   - logPath: 日志文件路径
//
// 返回:
//   - *Config: 生产环境配置实例
func Prod(logPath string) *Config {
	cfg := NewConfig(logPath)
	cfg.Level = WARN
	cfg.OutputConsole = false
	cfg.Async = true
	cfg.MaxSize = 100
	cfg.MaxFiles = 14
	cfg.MaxAge = 14
	cfg.Compress = true
	cfg.LevelRouter = true
	return cfg
}

// Console 创建纯控制台配置
//
// 基于 NewConfig 修改以下配置:
//   - Level: DEBUG (保留所有信息)
//   - OutputFile: false (不输出到文件)
//   - SamplerTick: 0 (关闭采样)
//
// 返回:
//   - *Config: 纯控制台配置实例
func Console() *Config {
	cfg := NewConfig("")
	cfg.Level = DEBUG
	cfg.OutputFile = false
	cfg.LogPath = ""
	cfg.SamplerTick = 0
	cfg.SamplerInitial = 0
	cfg.SamplerThereafter = 0
	return cfg
}

// Docker 创建容器生产环境配置
//
// 适用于 Docker、Kubernetes 等容器环境。
// 基于 NewConfig 修改以下配置:
//   - Level: WARN (只记录警告及以上)
//   - Formatter: JSON{} (结构化日志, 方便收集系统解析)
//   - OutputConsole: true (输出到 stdout, 容器标准做法)
//   - OutputFile: false (不写入文件, 由容器收集)
//   - LogPath: "" (无文件路径)
//   - SamplerTick: DefaultSamplerTick
//   - SamplerInitial: DefaultSamplerInitial
//   - SamplerThereafter: DefaultSamplerThereafter
//
// 返回:
//   - *Config: 容器生产环境配置实例
func Docker() *Config {
	cfg := NewConfig("")
	cfg.Level = WARN
	cfg.Formatter = JSON{}
	cfg.OutputConsole = true
	cfg.OutputFile = false
	cfg.LogPath = ""
	return cfg
}

// Config 日志记录器配置
//
// OutputConsole 和 OutputFile 可同时启用, 日志会同时写入终端和文件。
// 两者必须设置一个输出, 否则会报错。
type Config struct {
	// ======== 基础日志配置 ========

	// Level 日志级别, 零值默认 INFO
	Level Level

	// Formatter 日志格式化器, 零值默认 Def
	Formatter Formatter

	// Caller 是否记录调用者信息 (文件:函数:行号)
	Caller bool

	// Fields 预设字段, 每条日志都会自动携带这些字段
	Fields []Field

	// SamplerTick 采样时间窗口, 零值表示不启用采样
	// 例如 10*time.Second 表示每 10 秒为一个采样窗口
	SamplerTick time.Duration

	// SamplerInitial 每个窗口内前 N 条日志放行
	SamplerInitial int

	// SamplerThereafter 之后每 M 条放行 1 条, 0 表示不再放行
	SamplerThereafter int

	// ======== 终端输出配置 ========

	// OutputConsole 是否输出到终端 (彩色自动检测)
	OutputConsole bool

	// NoColor 设为 true 时禁用终端彩色输出, 仅当 OutputConsole=true 时生效
	NoColor bool

	// ======== 文件输出配置 ========

	// OutputFile 是否输出到文件
	OutputFile bool

	// LogPath 日志文件路径
	LogPath string

	// Async 是否启用异步清理 (单协程、合并触发) , 关闭文件时后台清理旧日志
	Async bool

	// MaxSize 单文件最大大小 (MB) , 超过后触发轮转。零值默认 10MB。
	MaxSize int

	// MaxFiles 保留的历史日志文件数, 超过的旧文件自动删除。零值表示不限制。
	MaxFiles int

	// MaxAge 保留天数, 超过指定天数的旧文件自动清理。零值表示不限制。
	MaxAge int

	// Compress 是否压缩历史日志文件
	Compress bool

	// CompressType 压缩类型, 默认: comprx.CompressTypeGz
	//
	// 支持的压缩格式:
	//   - comprx.CompressTypeZip
	//   - comprx.CompressTypeTar
	//   - comprx.CompressTypeTgz
	//   - comprx.CompressTypeTarGz
	//   - comprx.CompressTypeGz
	//   - comprx.CompressTypeBz2
	//   - comprx.CompressTypeBzip2
	//   - comprx.CompressTypeZlib
	CompressType comprx.CompressType

	// LocalTime 是否使用本地时间命名轮转文件
	LocalTime bool

	// DateDirLayout 是否按日期目录存放轮转文件
	DateDirLayout bool

	// RotateByDay 是否按天轮转
	RotateByDay bool

	// ======== 缓冲写入配置 ========

	// MaxBufferSize 缓冲区大小 (字节) , 零值默认 256KB
	MaxBufferSize int

	// SyncInterval 自动同步间隔, 零值默认 1 秒
	SyncInterval time.Duration

	// TimeFormat 时间格式
	// 默认 time.RFC3339，支持 Go time 包所有格式常量
	// 常用值: time.RFC3339, time.DateTime, time.TimeOnly
	TimeFormat string

	// LevelRouter 启用级别路由
	// 为 true 时，自动在 LogPath 同级目录创建 {LEVEL}.log 文件
	// 例: LogPath="logs/app.log" → 创建 logs/DEBUG.log, logs/INFO.log 等
	// 注意：启用后，每条日志会同时写入主文件和对应级别专属文件
	LevelRouter bool
}

// NewSampler 根据配置创建采样器
//
// 当 SamplerTick > 0 时创建采样器, 否则返回 nil 表示不启用采样。
func (c *Config) NewSampler() *Sampler {
	if c.SamplerTick <= 0 {
		return nil
	}
	return NewSampler(c.SamplerTick, c.SamplerInitial, c.SamplerThereafter)
}

// Clone 克隆配置
//
// 返回配置的深拷贝副本, 与原始配置完全独立互不干扰。
// Fields 切片会独立复制。
func (c *Config) Clone() *Config {
	clone := *c
	if len(c.Fields) > 0 {
		clone.Fields = make([]Field, len(c.Fields))
		copy(clone.Fields, c.Fields)
	}
	return &clone
}

// NewWriter 根据配置创建日志写入器
//
// 返回日志写入器, 用于将日志写入终端或文件。
func (c *Config) NewWriter() io.WriteCloser {
	// 根据配置创建对应的写入器
	switch {
	// 情况1: 仅终端输出
	case !c.OutputFile && c.OutputConsole:
		return NewColorWriter(c.NoColor)

	// 情况2: 仅文件输出
	case c.OutputFile && !c.OutputConsole:
		return c.newFileWriter()

	// 情况3: 同时输出到文件和终端
	case c.OutputFile && c.OutputConsole:
		fileWriter := c.newFileWriter()
		consoleWriter := NewColorWriter(c.NoColor)
		return NewMultiWriter(fileWriter, consoleWriter)

	// 情况4: 未设置任何输出 (理论上不会走到这里, 因为 Validate 已检查)
	default:
		return nil
	}
}

// newFileWriter 创建文件写入器 (内部辅助方法)
func (c *Config) newFileWriter() io.WriteCloser {
	// 创建日志切割器
	logger := &logrotatex.LogRotateX{
		LogFilePath:   c.LogPath,       // 日志文件路径
		MaxSize:       c.MaxSize,       // 最大日志文件大小, 单位为MB
		MaxAge:        c.MaxAge,        // 最大日志文件保留天数
		MaxFiles:      c.MaxFiles,      // 最大日志文件保留数量
		Compress:      c.Compress,      // 是否启用日志文件压缩
		LocalTime:     c.LocalTime,     // 是否使用本地时间
		Async:         c.Async,         // 是否异步清理日志
		DateDirLayout: c.DateDirLayout, // 是否启用按日期目录存放轮转后的日志
		RotateByDay:   c.RotateByDay,   // 是否启用按天轮转
		CompressType:  c.CompressType,  // 压缩类型
	}

	// 配置缓冲写入器
	bufCfg := &logrotatex.BufCfg{
		SyncInterval:  c.SyncInterval,  // 同步间隔
		MaxBufferSize: c.MaxBufferSize, // 缓冲区最大容量
	}

	// 返回带缓冲的批量写入器, 嵌入日志切割器
	return logrotatex.NewBufferedWriter(logger, bufCfg)
}

// Validate 验证配置是否有效
//
// 返回:
//   - error: 验证通过时返回 nil, 否则返回错误信息
func (c *Config) Validate() error {
	// 如果未设置输出, 返回错误
	if !c.OutputFile && !c.OutputConsole {
		return errors.New("output must be set")
	}

	// 验证采样器配置
	if c.SamplerTick > 0 {
		// 如果启用了采样, SamplerInitial 必须 >= 0 (零值表示不放行)
		if c.SamplerInitial < 0 {
			return errors.New("sampler initial must be >= 0")
		}
		// SamplerThereafter 必须 >= 0 (零值表示之后不再放行)
		if c.SamplerThereafter < 0 {
			return errors.New("sampler thereafter must be >= 0")
		}
	}

	// 验证文件输出配置
	if c.OutputFile {
		// LogPath 不能为空
		if c.LogPath == "" {
			return errors.New("log path must be set when output file is enabled")
		}
		// MaxSize 不能为负数
		if c.MaxSize < 0 {
			return errors.New("max size must be >= 0")
		}
		// MaxFiles 不能为负数
		if c.MaxFiles < 0 {
			return errors.New("max files must be >= 0")
		}
		// MaxAge 不能为负数
		if c.MaxAge < 0 {
			return errors.New("max age must be >= 0")
		}
	}

	// 验证级别路由配置
	if c.LevelRouter {
		// 必须设置文件输出
		if !c.OutputFile || c.LogPath == "" {
			return errors.New("level router requires file output and log path")
		}
		// 检查路径冲突：LogPath 不能与任何级别文件冲突
		dir := filepath.Dir(c.LogPath)
		for _, lvl := range AllLevels() {
			lvlPath := filepath.Join(dir, lvl.String()+".log")
			if lvlPath == c.LogPath {
				return fmt.Errorf("log path %s conflicts with level file path", c.LogPath)
			}
		}
	}

	// 验证缓冲写入配置
	if c.MaxBufferSize < 0 {
		return errors.New("max buffer size must be >= 0")
	}
	// 如果设置了缓冲区大小, 不能小于 64KB
	if c.MaxBufferSize > 0 && c.MaxBufferSize < 64*1024 {
		return errors.New("max buffer size must be >= 64KB")
	}
	if c.SyncInterval < 0 {
		return errors.New("sync interval must be >= 0")
	}
	// 如果设置了同步间隔, 不能小于 500ms
	if c.SyncInterval > 0 && c.SyncInterval < 500*time.Millisecond {
		return errors.New("sync interval must be >= 500ms")
	}

	return nil
}

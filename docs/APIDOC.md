# package fastlog

```
package fastlog // import "gitee.com/MM-Q/fastlog"
```

## CONSTANTS

```go
const (
	LevelNameDebug = "DEBUG"
	LevelNameInfo  = "INFO"
	LevelNameWarn  = "WARN"
	LevelNameError = "ERROR"
	LevelNameFatal = "FATAL"
	LevelNamePanic = "PANIC"
)
```

日志级别名称常量

```go
const (
	DefaultSamplerTick       = 10 * time.Second // 默认采样时间窗口
	DefaultSamplerInitial    = 3                // 默认窗口内前 N 条放行
	DefaultSamplerThereafter = 10               // 默认之后每 M 条放行 1 条
)
```

```go
const DefaultTimeFormat = time.DateTime // 2006-01-02 15:04:05
```

DefaultTimeFormat 默认时间格式

## VARIABLES

```go
var EntryPool = sync.Pool{
	New: func() interface{} {
		return &Entry{
			Fields: make([]Field, 0, 8),
		}
	},
}
```

EntryPool 日志条目池, 用于减少内存分配

## FUNCTIONS

### Close

```go
func Close() error
```

Close 关闭全局默认日志记录器

如果全局日志记录器未初始化（未调用 L()）, 返回 nil。

### LogRequest

```go
func LogRequest(log *Logger, next http.Handler) http.Handler
```

LogRequest 日志中间件，用于记录HTTP请求日志

**参数:**
- `log`: 日志实例
- `next`: 下一个处理器

**返回:**
- `http.Handler`: 中间件处理后的处理器

### PutEntry

```go
func PutEntry(e *Entry)
```

PutEntry 将日志条目放回池中

**参数:**
- `e`: 要放回池的日志条目

### Sync

```go
func Sync() error
```

Sync 同步全局默认日志记录器的缓冲区数据到存储

如果全局日志记录器未初始化（未调用 L()）或写入器不支持同步, 返回 nil。

## TYPES

### ColorWriter

```go
type ColorWriter struct {
	NoColor bool // 设为 true 禁用颜色, false 启用颜色
	// Has unexported fields.
}
```

ColorWriter 彩色控制台写入器

通过扫描字节流中的日志级别关键字自动着色输出。 将 NoColor 设为 true 可禁用颜色输出, 恢复原始文本。

#### func NewColorWriter

```go
func NewColorWriter(noColor bool) *ColorWriter
```

NewColorWriter 创建彩色控制台写入器, 默认写入 os.Stdout

**参数:**
- `noColor`: 设为 true 禁用颜色输出

**返回:**
- `*ColorWriter`: 彩色写入器实例

#### func (c *ColorWriter) Close

```go
func (c *ColorWriter) Close() error
```

Close 关闭写入器

**返回:**
- `error`: 始终返回 nil

#### func (c *ColorWriter) Write

```go
func (c *ColorWriter) Write(p []byte) (n int, err error)
```

Write 写入数据到控制台, 自动根据日志级别着色

**参数:**
- `p`: 要写入的字节数据

**返回:**
- `int`: 写入的字节数
- `error`: 写入过程中的错误

### Compact

```go
type Compact struct{}
```

Compact 极简格式 格式: [I] 2025-01-15 10:30:45 用户登录成功 特点: 级别首字母 + 时间戳，简洁易读，适合容器环境
时间格式遵循 Config.TimeFormat，默认 time.DateTime (2006-01-02 15:04:05)

#### func (f Compact) Format

```go
func (f Compact) Format(entry *Entry) ([]byte, error)
```

Format 实现 Compact 格式

**参数:**
- `entry`: 日志条目

**返回:**
- `[]byte`: 格式化后的字节数组
- `error`: 如果格式化失败

### Config

```go
type Config struct {

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

	// OutputConsole 是否输出到终端 (彩色自动检测)
	OutputConsole bool

	// NoColor 设为 true 时禁用终端彩色输出, 仅当 OutputConsole=true 时生效
	NoColor bool

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

	// MaxBufferSize 缓冲区大小 (字节) , 零值默认 256KB
	MaxBufferSize int

	// SyncInterval 自动同步间隔, 零值默认 1 秒
	SyncInterval time.Duration

	// TimeFormat 时间格式
	// 默认 time.DateTime (2006-01-02 15:04:05)，支持 Go time 包所有格式常量
	// 常用值: time.DateTime, time.RFC3339, time.TimeOnly
	TimeFormat string

	// LevelRouter 启用级别路由
	// 为 true 时，自动在 LogPath 同级目录创建 {LEVEL}.log 文件
	// 例: LogPath="logs/app.log" → 创建 logs/DEBUG.log, logs/INFO.log 等
	// 注意：启用后，每条日志会同时写入主文件和对应级别专属文件
	LevelRouter bool

	// BufferEnabled 是否启用缓冲写入
	// true:  使用 BufferedWriter（默认），提升写入性能
	// false: 直接使用 LogRotateX，无缓冲，立即落盘
	//
	// 使用建议：
	//   - 开发环境：设为 false，立即看到日志，便于调试
	//   - 生产环境：设为 true（默认），提升写入性能
	//   - 高可靠性场景：设为 false，避免数据丢失风险
	BufferEnabled bool
}
```

Config 日志记录器配置

OutputConsole 和 OutputFile 可同时启用, 日志会同时写入终端和文件。 两者必须设置一个输出, 否则会报错。

#### func Console

```go
func Console() *Config
```

Console 创建纯控制台配置

基于 NewConfig 修改以下配置:
- `Level`: DEBUG (保留所有信息)
- `OutputFile`: false (不输出到文件)
- `SamplerTick`: 0 (关闭采样)

**返回:**
- `*Config`: 纯控制台配置实例

#### func Default

```go
func Default() *Config
```

Default 创建一个默认配置实例

默认配置同时输出到终端和文件, 适用于开发和测试环境。

默认配置详情:
- `Level`: INFO - 日志级别为信息级别
- `Formatter`: Def{} - 使用默认格式输出
- `Caller`: false - 不记录调用者信息
- `Fields`: []Field{} - 无预设字段
- `SamplerTick`: 10s - 采样窗口为10秒
- `SamplerInitial`: 3 - 前3条日志放行
- `SamplerThereafter`: 10 - 之后每10条放行1条
- `LevelRouter`: false - 不启用级别路由
- `OutputConsole`: true - 输出到终端
- `NoColor`: false - 启用彩色输出
- `OutputFile`: true - 输出到文件
- `LogPath`: 参数指定 - 日志文件路径由参数指定
- `Async`: false - 不启用异步清理
- `MaxSize`: 20MB - 单文件最大20MB
- `MaxFiles`: 24 - 保留24个历史文件
- `MaxAge`: 7天 - 保留7天日志
- `Compress`: false - 不压缩历史日志
- `CompressType`: Gz - 压缩类型为gzip
- `LocalTime`: true - 使用本地时间命名
- `DateDirLayout`: true - 按日期目录存放
- `RotateByDay`: true - 按天轮转文件
- `BufferEnabled`: true - 启用缓冲写入
- `MaxBufferSize`: 256KB - 缓冲区256KB
- `SyncInterval`: 1s - 每秒同步
- `TimeFormat`: 2006-01-02 15:04:05 - 时间格式为 DateTime

**返回:**
- `*Config`: 配置实例, 所有字段都设置了默认值

#### func Dev

```go
func Dev(logPath string) *Config
```

Dev 创建开发环境配置

基于 NewConfig 修改以下配置:
- `Level`: DEBUG (更多信息便于调试)
- `Caller`: true (开启调用者信息, 方便定位代码)
- `SamplerTick`: 0 (关闭采样, 保留所有日志)
- `MaxSize`: 10MB (小文件, 便于查看)
- `Compress`: false (不压缩, 快速写入)
- `RotateByDay`: false (不按天轮转)
- `BufferEnabled`: false (禁用缓冲, 立即写入, 便于调试)

**参数:**
- `logPath`: 日志文件路径

**返回:**
- `*Config`: 开发环境配置实例

#### func Docker

```go
func Docker() *Config
```

Docker 创建容器生产环境配置

适用于 Docker、Kubernetes 等容器环境。 基于 NewConfig 修改以下配置:
- `Level`: WARN (只记录警告及以上)
- `Formatter`: JSON{} (结构化日志, 方便收集系统解析)
- `OutputConsole`: true (输出到 stdout, 容器标准做法)
- `OutputFile`: false (不写入文件, 由容器收集)
- `LogPath`: "" (无文件路径)
- `SamplerTick`: DefaultSamplerTick
- `SamplerInitial`: DefaultSamplerInitial
- `SamplerThereafter`: DefaultSamplerThereafter

**返回:**
- `*Config`: 容器生产环境配置实例

#### func NewConfig

```go
func NewConfig(logPath string) *Config
```

NewConfig 创建一个默认配置实例

默认配置同时输出到终端和文件, 适用于开发和测试环境。

默认配置详情:
- `Level`: INFO - 日志级别为信息级别
- `Formatter`: Def{} - 使用默认格式输出
- `Caller`: false - 不记录调用者信息
- `Fields`: []Field{} - 无预设字段
- `SamplerTick`: 10s - 采样窗口为10秒
- `SamplerInitial`: 3 - 前3条日志放行
- `SamplerThereafter`: 10 - 之后每10条放行1条
- `LevelRouter`: false - 不启用级别路由
- `OutputConsole`: true - 输出到终端
- `NoColor`: false - 启用彩色输出
- `OutputFile`: true - 输出到文件
- `LogPath`: 参数指定 - 日志文件路径由参数指定
- `Async`: false - 不启用异步清理
- `MaxSize`: 20MB - 单文件最大20MB
- `MaxFiles`: 24 - 保留24个历史文件
- `MaxAge`: 7天 - 保留7天日志
- `Compress`: false - 不压缩历史日志
- `CompressType`: Gz - 压缩类型为gzip
- `LocalTime`: true - 使用本地时间命名
- `DateDirLayout`: true - 按日期目录存放
- `RotateByDay`: true - 按天轮转文件
- `BufferEnabled`: true - 启用缓冲写入
- `MaxBufferSize`: 256KB - 缓冲区256KB
- `SyncInterval`: 1s - 每秒同步
- `TimeFormat`: 2006-01-02 15:04:05 - 时间格式为 DateTime

**参数:**
- `logPath`: 日志文件路径

**返回:**
- `*Config`: 配置实例, 所有字段都设置了默认值

#### func Prod

```go
func Prod(logPath string) *Config
```

Prod 创建生产环境配置

基于 NewConfig 修改以下配置:
- `Level`: WARN (只记录警告及以上, 减少日志量)
- `OutputConsole`: false (只输出到文件, 无终端开销)
- `Async`: true (异步清理, 性能更好)
- `MaxSize`: 100MB (大文件, 减少轮转频率)
- `MaxFiles`: 14 (保留2周)
- `MaxAge`: 14 (保留14天)
- `Compress`: true (开启压缩, 节省磁盘)
- `LevelRouter`: true (启用级别路由, 便于快速定位错误)
- `BufferEnabled`: true (启用缓冲, 提升写入性能)

**参数:**
- `logPath`: 日志文件路径

**返回:**
- `*Config`: 生产环境配置实例

#### func (c *Config) Clone

```go
func (c *Config) Clone() *Config
```

Clone 克隆配置

返回配置的深拷贝副本, 与原始配置完全独立互不干扰。 Fields 切片会独立复制。

#### func (c *Config) NewSampler

```go
func (c *Config) NewSampler() *Sampler
```

NewSampler 根据配置创建采样器

当 SamplerTick > 0 时创建采样器, 否则返回 nil 表示不启用采样。

#### func (c *Config) NewWriter

```go
func (c *Config) NewWriter() io.WriteCloser
```

NewWriter 根据配置创建日志写入器

返回日志写入器, 用于将日志写入终端或文件。

#### func (c *Config) Validate

```go
func (c *Config) Validate() error
```

Validate 验证配置是否有效

**返回:**
- `error`: 验证通过时返回 nil, 否则返回错误信息

### ConsoleWriter

```go
type ConsoleWriter struct {
	// Has unexported fields.
}
```

ConsoleWriter 控制台写入器

#### func (c *ConsoleWriter) Close

```go
func (c *ConsoleWriter) Close() error
```

Close 关闭控制台写入器

**返回:**
- `error`: 始终返回 nil

#### func (c *ConsoleWriter) Write

```go
func (c *ConsoleWriter) Write(p []byte) (n int, err error)
```

Write 写入数据到控制台

**参数:**
- `p`: 要写入的字节数据

**返回:**
- `int`: 写入的字节数
- `error`: 写入过程中的错误

### Def

```go
type Def struct{}
```

Def 默认格式 格式: 2025-01-15T10:30:45 | INFO | main.go:main:15 - 用户登录成功

#### func (f Def) Format

```go
func (f Def) Format(entry *Entry) ([]byte, error)
```

Format 实现默认格式化器

**参数:**
- `entry`: 日志条目

**返回:**
- `[]byte`: 格式化后的字节数组
- `error`: 如果格式化失败

### Entry

```go
type Entry struct {
	Time       time.Time // 时间戳
	Level      Level     // 日志级别
	Message    string    // 日志消息
	Caller     string    // 调用者信息: file.go:func:line
	Fields     []Field   // 键值对字段
	TimeFormat string    // 时间格式, 从 Config.TimeFormat 传递
}
```

Entry 表示一条日志记录

#### func GetEntry

```go
func GetEntry() *Entry
```

GetEntry 从池中获取日志条目

**返回:**
- `*Entry`: 日志条目实例
- `error`: 如果池为空, 返回错误

### Field

```go
type Field struct {
	// Has unexported fields.
}
```

Field 表示一个键值对字段, 包含所有可能的类型

#### func Any

```go
func Any(key string, val interface{}) Field
```

Any 创建一个任意类型字段

**参数:**
- `key`: 字段键
- `val`: 任意类型的值

**返回:**
- `Field`: 字段实例

#### func Bool

```go
func Bool(key string, val bool) Field
```

Bool 创建一个 bool 字段

**参数:**
- `key`: 字段键
- `val`: 布尔值

**返回:**
- `Field`: 字段实例

#### func Duration

```go
func Duration(key string, val time.Duration) Field
```

Duration 创建一个 time.Duration 字段

**参数:**
- `key`: 字段键
- `val`: 时间持续值

**返回:**
- `Field`: 字段实例

#### func Err

```go
func Err(key string, err error) Field
```

Err 创建一个自定义键名的 error 字段

**参数:**
- `key`: 字段键名
- `err`: 错误值, nil 时返回 "<nil>" 字符串

**返回:**
- `Field`: 字段实例

**示例:**

```go
logger.Infow("操作失败",
    fastlog.Err("db_error", dbErr),
    fastlog.Err("cache_error", cacheErr),
)
```

#### func Error

```go
func Error(err error) Field
```

Error 创建一个 error 字段

**参数:**
- `err`: 错误值, nil 时返回 "<nil>" 字符串

**返回:**
- `Field`: 字段实例, 键名为 "error"

#### func Float64

```go
func Float64(key string, val float64) Field
```

Float64 创建一个 float64 字段

**参数:**
- `key`: 字段键
- `val`: 浮点数值

**返回:**
- `Field`: 字段实例

#### func Int

```go
func Int(key string, val int) Field
```

Int 创建一个 int 字段

**参数:**
- `key`: 字段键
- `val`: 整数值

**返回:**
- `Field`: 字段实例

#### func Int64

```go
func Int64(key string, val int64) Field
```

Int64 创建一个 int64 字段

**参数:**
- `key`: 字段键
- `val`: 64位整数值

**返回:**
- `Field`: 字段实例

#### func Stack

```go
func Stack() Field
```

Stack 创建一个堆栈字段

**返回:**
- `Field`: 字段实例, 键名为 "stack", 值为当前堆栈信息

#### func String

```go
func String(key, val string) Field
```

String 创建一个字符串字段

**参数:**
- `key`: 字段键
- `val`: 字符串值

**返回:**
- `Field`: 字段实例

#### func Time

```go
func Time(key string, val time.Time) Field
```

Time 创建一个 time.Time 字段

**参数:**
- `key`: 字段键
- `val`: 时间值

**返回:**
- `Field`: 字段实例

#### func Uint

```go
func Uint(key string, val uint) Field
```

Uint 创建一个 uint 字段

**参数:**
- `key`: 字段键
- `val`: 无符号整数值

**返回:**
- `Field`: 字段实例

#### func Uint64

```go
func Uint64(key string, val uint64) Field
```

Uint64 创建一个 uint64 字段

**参数:**
- `key`: 字段键
- `val`: 64位无符号整数值

**返回:**
- `Field`: 字段实例

#### func (f Field) Format

```go
func (f Field) Format() string
```

Format 将字段格式化为 key=value 形式

**返回:**
- `string`: 格式化后的字段字符串，格式为 "key=value"

**示例:**

```go
field := fastlog.String("user", "admin")
result := field.Format() // 返回 "user=admin"
```

#### func (f Field) Key

```go
func (f Field) Key() string
```

Key 返回字段键

**返回:**
- `string`: 字段键名

#### func (f Field) Type

```go
func (f Field) Type() FieldType
```

Type 返回字段类型

**返回:**
- `FieldType`: 字段类型

#### func (f Field) Value

```go
func (f Field) Value() string
```

Value 将字段值转换为字符串返回

根据字段类型将值格式化为字符串:
- `StringType`/`ErrorType`: 直接返回字符串值
- `IntType`/`Int64Type`: 转为 10 进制字符串
- `UintType`/`Uint64Type`: 转为 10 进制无符号字符串
- `Float64Type`: 转为浮点数字符串
- `BoolType`: 转为 "true" 或 "false"
- `TimeType`: 转为 DateTime 格式时间字符串 (2006-01-02 15:04:05)
- `DurationType`: 转为持续时间字符串 (如 "1h30m")
- `AnyType`: 使用 fmt.Sprintf("%v") 格式化
- 其他类型: 返回空字符串

**返回:**
- `string`: 字段值的字符串表示

### FieldType

```go
type FieldType int8
```

FieldType 表示字段类型

```go
const (
	UnknownType  FieldType = iota // 未知类型
	StringType                    // 字符串类型
	IntType                       // 整数类型
	Int64Type                     // 64位整数类型
	UintType                      // 无符号整数类型
	Uint64Type                    // 64位无符号整数类型
	Float64Type                   // 浮点数类型
	BoolType                      // 布尔类型
	TimeType                      // 时间类型
	DurationType                  // 时间持续类型
	ErrorType                     // 错误类型
	AnyType                       // any类型
)
```

字段类型常量

### Formatter

```go
type Formatter interface {
	// Format 将日志条目格式化为字节数组
	Format(entry *Entry) ([]byte, error)
}
```

Formatter 定义日志格式化器接口

### JSON

```go
type JSON struct{}
```

JSON JSON 格式

#### func (f JSON) Format

```go
func (f JSON) Format(entry *Entry) ([]byte, error)
```

Format 实现 JSON 格式

**参数:**
- `entry`: 日志条目

**返回:**
- `[]byte`: 格式化后的字节数组
- `error`: 如果格式化失败

### KV

```go
type KV struct{}
```

KV 键值对格式 格式: time=2025-01-15T10:30:45 level=INFO message=用户登录成功

#### func (f KV) Format

```go
func (f KV) Format(entry *Entry) ([]byte, error)
```

Format 实现键值对格式

**参数:**
- `entry`: 日志条目

**返回:**
- `[]byte`: 格式化后的字节数组
- `error`: 如果格式化失败

### Level

```go
type Level int32
```

Level 表示日志级别

```go
const (
	DEBUG Level = iota + 1 // 调试级别 (1)
	INFO                   // 信息级别 (2)
	WARN                   // 警告级别 (3)
	ERROR                  // 错误级别 (4)
	FATAL                  // 致命级别 (5)
	PANIC                  // 恐慌级别 (6)
)
```

日志级别常量

#### func AllLevels

```go
func AllLevels() []Level
```

AllLevels 返回所有日志级别

**返回:**
- `[]Level`: 包含所有日志级别的切片

#### func ParseLevel

```go
func ParseLevel(s string) (Level, error)
```

ParseLevel 从字符串解析日志级别

**参数:**
- `s`: 要解析的字符串

**返回:**
- `Level`: 解析后的日志级别
- `error`: 如果解析失败

#### func (l Level) Enabled

```go
func (l Level) Enabled(lvl Level) bool
```

Enabled 检查是否启用该级别 (lvl >= l 时启用)

**参数:**
- `lvl`: 要检查的级别

**返回:**
- `bool`: 是否启用该级别

#### func (l Level) String

```go
func (l Level) String() string
```

String 返回级别的字符串表示

### Logger

```go
type Logger struct {
	// Has unexported fields.
}
```

Logger 日志记录器

Logger 是 FastLog 的核心日志记录器, 提供日志记录、级别控制、采样等功能。 支持 6 种日志级别: DEBUG, INFO, WARN,
ERROR, FATAL, PANIC。 支持三种调用方式: 标准日志 (Info)、格式化日志 (Infof)、结构化日志 (Infow)。

必须通过 fastlog.New(cfg) 构造函数创建, 切勿直接声明空结构体使用。 直接声明 Logger{} 会导致内部
writer/sampler 等关键字段为 nil, 引发 panic。

使用示例:

```go
cfg := fastlog.NewConfig("logs/app.log")
logger := fastlog.New(cfg)
defer func() { _ = logger.Close() }()
logger.Info("服务启动成功")
```

#### func L

```go
func L() *Logger
```

L 返回全局默认日志记录器

全局日志记录器在第一次调用时创建, 使用 Console() 配置（DEBUG 级别, 纯控制台输出）。 适合快速使用和调试, 无需手动创建
Logger 实例:

```go
fastlog.L().Info("服务启动")
fastlog.L().Errorw("连接失败", fastlog.Err(err))
```

#### func New

```go
func New(cfg *Config) *Logger
```

New 创建一个新的日志记录器

**参数:**
- `cfg`: 日志配置, 零值时使用默认配置

**返回:**
- `*Logger`: 新的日志记录器实例

注意: 如果配置验证失败, 会触发 panic 以便快速发现问题

#### func (l *Logger) Close

```go
func (l *Logger) Close() error
```

Close 关闭日志记录器

**返回:**
- `error`: 关闭过程中的错误

#### func (l *Logger) Debug

```go
func (l *Logger) Debug(msg string)
```

Debug 记录调试日志

**参数:**
- `msg`: 日志消息

#### func (l *Logger) Debugf

```go
func (l *Logger) Debugf(format string, args ...interface{})
```

Debugf 记录格式化的调试日志

**参数:**
- `format`: 格式化字符串
- `args`: 格式化参数

#### func (l *Logger) Debugw

```go
func (l *Logger) Debugw(msg string, fields ...Field)
```

Debugw 记录带字段的调试日志

**参数:**
- `msg`: 日志消息
- `fields`: 日志字段

#### func (l *Logger) Error

```go
func (l *Logger) Error(msg string)
```

Error 记录错误日志

**参数:**
- `msg`: 日志消息

#### func (l *Logger) Errorf

```go
func (l *Logger) Errorf(format string, args ...interface{})
```

Errorf 记录格式化的错误日志

**参数:**
- `format`: 格式化字符串
- `args`: 格式化参数

#### func (l *Logger) Errorw

```go
func (l *Logger) Errorw(msg string, fields ...Field)
```

Errorw 记录带字段的错误日志

**参数:**
- `msg`: 日志消息
- `fields`: 日志字段

#### func (l *Logger) Fatal

```go
func (l *Logger) Fatal(msg string)
```

Fatal 记录致命日志并退出程序

**参数:**
- `msg`: 日志消息

#### func (l *Logger) Fatalf

```go
func (l *Logger) Fatalf(format string, args ...interface{})
```

Fatalf 记录格式化的致命日志并退出程序

**参数:**
- `format`: 格式化字符串
- `args`: 格式化参数

#### func (l *Logger) Fatalw

```go
func (l *Logger) Fatalw(msg string, fields ...Field)
```

Fatalw 记录带字段的致命日志并退出程序

**参数:**
- `msg`: 日志消息
- `fields`: 日志字段

#### func (l *Logger) Info

```go
func (l *Logger) Info(msg string)
```

Info 记录信息日志

**参数:**
- `msg`: 日志消息

#### func (l *Logger) Infof

```go
func (l *Logger) Infof(format string, args ...interface{})
```

Infof 记录格式化的信息日志

**参数:**
- `format`: 格式化字符串
- `args`: 格式化参数

#### func (l *Logger) Infow

```go
func (l *Logger) Infow(msg string, fields ...Field)
```

Infow 记录带字段的信息日志

**参数:**
- `msg`: 日志消息
- `fields`: 日志字段

#### func (l *Logger) Level

```go
func (l *Logger) Level() Level
```

Level 返回当前运行时日志级别

**返回:**
- `Level`: 当前日志级别

#### func (l *Logger) Panic

```go
func (l *Logger) Panic(msg string)
```

Panic 记录恐慌日志并触发 panic

**参数:**
- `msg`: 日志消息

#### func (l *Logger) Panicf

```go
func (l *Logger) Panicf(format string, args ...interface{})
```

Panicf 记录格式化的恐慌日志并触发 panic

**参数:**
- `format`: 格式化字符串
- `args`: 格式化参数

#### func (l *Logger) Panicw

```go
func (l *Logger) Panicw(msg string, fields ...Field)
```

Panicw 记录带字段的恐慌日志并触发 panic

**参数:**
- `msg`: 日志消息
- `fields`: 日志字段

#### func (l *Logger) SetLevel

```go
func (l *Logger) SetLevel(level Level)
```

SetLevel 运行时动态修改日志级别, 立即生效

**参数:**
- `level`: 新的日志级别

#### func (l *Logger) Sync

```go
func (l *Logger) Sync() error
```

Sync 同步日志到存储

**返回:**
- `error`: 同步过程中的错误, 如果写入器不支持同步则返回 nil

#### func (l *Logger) Warn

```go
func (l *Logger) Warn(msg string)
```

Warn 记录警告日志

**参数:**
- `msg`: 日志消息

#### func (l *Logger) Warnf

```go
func (l *Logger) Warnf(format string, args ...interface{})
```

Warnf 记录格式化的警告日志

**参数:**
- `format`: 格式化字符串
- `args`: 格式化参数

#### func (l *Logger) Warnw

```go
func (l *Logger) Warnw(msg string, fields ...Field)
```

Warnw 记录带字段的警告日志

**参数:**
- `msg`: 日志消息
- `fields`: 日志字段

### MultiWriter

```go
type MultiWriter struct {
	// Has unexported fields.
}
```

MultiWriter 多路写入器, 同时将日志写入多个输出目标

#### func NewMultiWriter

```go
func NewMultiWriter(writers ...io.WriteCloser) *MultiWriter
```

NewMultiWriter 创建多路写入器

**参数:**
- `writers`: 多个写入器

**返回:**
- `*MultiWriter`: 多路写入器实例

#### func (m *MultiWriter) Close

```go
func (m *MultiWriter) Close() error
```

Close 关闭所有输出目标

**返回:**
- `error`: 关闭过程中的错误, 多个错误会合并返回

#### func (m *MultiWriter) Write

```go
func (m *MultiWriter) Write(p []byte) (n int, err error)
```

Write 写入数据到所有输出目标

**参数:**
- `p`: 要写入的字节数据

**返回:**
- `int`: 写入的字节数
- `error`: 写入过程中的错误

### Sampler

```go
type Sampler struct {
	// Has unexported fields.
}
```

Sampler 日志采样器

使用固定桶 + atomic 计数器实现, 无锁设计。 相同 level 和 message 的日志会被哈希到同一个桶, 在时间窗口内按规则放行或抑制。

#### func DefaultSampler

```go
func DefaultSampler() *Sampler
```

DefaultSampler 创建默认日志采样器

默认参数: 窗口 DefaultSamplerTick, 前 DefaultSamplerInitial 条放行, 之后每
DefaultSamplerThereafter 条放行 1 条。 适合大多数场景直接使用, 无需额外配置。

**示例:**

```go
sampler := fastlog.DefaultSampler()
```

#### func NewSampler

```go
func NewSampler(tick time.Duration, initial, thereafter int) *Sampler
```

NewSampler 创建日志采样器

**参数:**
- `tick`: 时间窗口, 如 10*time.Second。如果 <= 0, 默认使用 DefaultSamplerTick
- `initial`: 窗口内前 N 条放行。如果 < 0, 默认使用 DefaultSamplerInitial
- `thereafter`: 之后每 M 条放行 1 条, 0 表示不再放行。如果 < 0, 默认使用
  DefaultSamplerThereafter

**示例:**

```go
// 每 10 秒, 前 3 条放行, 之后每 10 条放行 1 条
sampler := fastlog.NewSampler(10*time.Second, 3, 10)
```

#### func (s *Sampler) Allow

```go
func (s *Sampler) Allow(level Level, msg string) bool
```

Allow 判断是否放行这条日志

**参数:**
- `level`: 日志级别
- `msg`: 日志消息

**返回:**
- `bool`: true 放行, false 抑制

### Simple

```go
type Simple struct{}
```

Simple 简单格式 格式: 2025-01-15T10:30:45 INFO 用户登录成功

#### func (f Simple) Format

```go
func (f Simple) Format(entry *Entry) ([]byte, error)
```

Format 实现简单格式

**参数:**
- `entry`: 日志条目

**返回:**
- `[]byte`: 格式化后的字节数组
- `error`: 如果格式化失败

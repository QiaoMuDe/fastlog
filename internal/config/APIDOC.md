# Package config

Package config (import "gitee.com/MM-Q/fastlog/internal/config")

## Functions

### CreateWriter

CreateWriter 根据配置创建文件写入器（支持带缓冲和普通文件写入）

```go
func CreateWriter(cfg *FastLogConfig) io.WriteCloser
```

**参数:**
- cfg: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。

**返回值:**
- io.WriteCloser: 文件写入器接口，同时支持写入和关闭操作。

**说明:**
- 当 BufferedWrite 为 true 时，返回带缓冲的批量写入器，支持日志轮转和压缩功能
- 当 BufferedWrite 为 false 时，返回普通文件句柄，直接写入文件，不提供日志轮转功能
- 如果 OutputToFile 为 false，则返回 nil

### Default

Default 返回一个默认的FastLogConfig实例, 用于配置日志记录器。

```go
func Default() *FastLogConfig
```

**返回值:**
- *FastLogConfig: 一个指向FastLogConfig实例的指针。

**默认配置:**
- 日志目录: "logs"
- 日志文件名: "app.log"
- 日志级别: INFO
- 日志格式: Def
- 最大日志文件大小: 10MB
- 最大日志文件保留天数: 0 (不做限制)
- 最大日志文件保留数量: 0 (不做限制)
- 是否使用本地时间: true
- 是否启用日志文件压缩: false
- 是否启用终端颜色: true
- 是否启用终端字体加粗: true
- 缓冲区大小: 256KB
- 刷新间隔: 1秒
- 是否异步清理日志: false
- 是否获取调用者信息: false
- 是否将日志输出到控制台: true
- 是否将日志输出到文件: true
- 是否启用按日期目录存放轮转后的日志: true
- 是否启用按天轮转: true
- 压缩类型: comprx.CompressTypeZip
- 是否使用带缓冲的批量写入器: true (默认)

## Types

### FastLogConfig

FastLogConfig 定义一个配置结构体, 用于配置日志记录器

```go
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
```

#### ConsoleConfig

ConsoleConfig 创建一个控制台环境下的FastLogConfig实例

```go
func ConsoleConfig() *FastLogConfig
```

**返回值:**
- *FastLogConfig: 一个指向FastLogConfig实例的指针。

**特性:**
- 禁用文件输出
- 设置日志级别为DEBUG
- 启用调用者信息
- 设置日志格式为默认格式

#### Default

Default 返回一个默认的FastLogConfig实例, 用于配置日志记录器。

```go
func Default() *FastLogConfig
```

**返回值:**
- *FastLogConfig: 一个指向FastLogConfig实例的指针。

**默认配置:**
- 日志目录: "logs"
- 日志文件名: "app.log"
- 日志级别: INFO
- 日志格式: Def
- 最大日志文件大小: 10MB
- 最大日志文件保留天数: 0 (不做限制)
- 最大日志文件保留数量: 0 (不做限制)
- 是否使用本地时间: true
- 是否启用日志文件压缩: false
- 是否启用终端颜色: true
- 是否启用终端字体加粗: true
- 缓冲区大小: 256KB
- 刷新间隔: 1秒
- 是否异步清理日志: false
- 是否获取调用者信息: false
- 是否将日志输出到控制台: true
- 是否将日志输出到文件: true
- 是否启用按日期目录存放轮转后的日志: true
- 是否启用按天轮转: true
- 压缩类型: comprx.CompressTypeZip
- 是否使用带缓冲的批量写入器: true (默认)

#### DevConfig

DevConfig 创建一个开发环境下的FastLogConfig实例

```go
func DevConfig(logDirName string, logFileName string) *FastLogConfig
```

**参数:**
- logDirName: 日志目录名
- logFileName: 日志文件名

**返回值:**
- *FastLogConfig: 一个指向FastLogConfig实例的指针。

**特性:**
- 启用详细日志格式
- 设置日志级别为DEBUG
- 设置最大日志文件保留数量为5
- 设置最大日志文件保留天数为7天
- 启用调用者信息

#### NewFastLogConfig

NewFastLogConfig 创建一个新的FastLogConfig实例, 用于配置日志记录器。

```go
func NewFastLogConfig(logDirName string, logFileName string) *FastLogConfig
```

**参数:**
- logDirName: 日志目录名称
- logFileName: 日志文件名称

**返回值:**
- *FastLogConfig: 一个指向FastLogConfig实例的指针。

**默认配置:**
- 日志级别: INFO
- 日志格式: Def
- 最大日志文件大小: 10MB
- 最大日志文件保留天数: 0 (不做限制)
- 最大日志文件保留数量: 0 (不做限制)
- 是否使用本地时间: true
- 是否启用日志文件压缩: false
- 是否启用终端颜色: true
- 是否启用终端字体加粗: true
- 缓冲区大小: 256KB
- 刷新间隔: 1秒
- 是否异步清理日志: false
- 是否获取调用者信息: false
- 是否将日志输出到控制台: true
- 是否将日志输出到文件: true
- 是否启用按日期目录存放轮转后的日志: true
- 是否启用按天轮转: true
- 压缩类型: comprx.CompressTypeZip
- 是否使用带缓冲的批量写入器: true (默认)

#### ProdConfig

ProdConfig 创建一个生产环境下的FastLogConfig实例

```go
func ProdConfig(logDirName string, logFileName string) *FastLogConfig
```

**参数:**
- logDirName: 日志目录名
- logFileName: 日志文件名

**返回值:**
- *FastLogConfig: 一个指向FastLogConfig实例的指针。

**特性:**
- 启用日志文件压缩
- 禁用控制台输出
- 设置最大日志文件保留天数为30天
- 设置最大日志文件保留数量为24个
- 设置日志格式为json格式

#### Clone

Clone 创建一个FastLogConfig的副本

```go
func (c *FastLogConfig) Clone() *FastLogConfig
```

**返回值:**
- *FastLogConfig: 一个指向FastLogConfig实例的指针，该实例是当前实例的副本。

#### ValidateConfig

ValidateConfig 验证配置并设置默认值 发现任何不合理的配置值都会panic，确保调用者提供正确配置

```go
func (c *FastLogConfig) ValidateConfig()
```
# Package config

Package config (import "gitee.com/MM-Q/fastlog/internal/config")

## Functions

### CreateBufferedWriter

CreateBufferedWriter 根据配置创建一个带缓冲的文件写入器

```go
func CreateBufferedWriter(cfg *FastLogConfig) *logrotatex.BufferedWriter
```

**参数:**
- cfg: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。

**返回值:**
- *logrotatex.BufferedWriter: 一个指向BufferedWriter实例的指针。

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
	MaxBufferSize   int                 // 缓冲区大小, 单位为字节, 默认64KB
	MaxWriteCount   int                 // 最大写入次数, 默认500次
	FlushInterval   time.Duration       // 刷新间隔, 默认1秒
	Async           bool                // 是否异步清理日志, 默认同步清理
	CallerInfo      bool                // 是否获取调用者信息, 默认不获取
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

#### NewFastLogConfig

NewFastLogConfig 创建一个新的FastLogConfig实例, 用于配置日志记录器。

```go
func NewFastLogConfig(logDirName string, logFileName string) *FastLogConfig
```

**参数:**
- logDirName: 日志目录名称, 默认为"applogs"。
- logFileName: 日志文件名称, 默认为"app.log"。

**返回值:**
- *FastLogConfig: 一个指向FastLogConfig实例的指针。

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
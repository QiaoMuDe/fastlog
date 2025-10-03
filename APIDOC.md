# Package fastlog

`fastlog` 是一个高性能的日志库，提供日志记录器的创建、初始化、日志写入及关闭等核心功能，集成配置管理、缓冲区管理和日志处理流程。

## 文件说明

- **config.go**：日志配置管理模块，定义日志配置结构体及配置项的设置与获取方法，负责管理 `FastLog` 的所有可配置参数。
- **fastlog.go**：`FastLog` 日志记录器核心实现，提供日志记录器的创建、初始化、日志写入及关闭等核心功能，集成配置管理、缓冲区管理和日志处理流程。
- **internal.go**：`FastLog` 内部实现文件，包含日志系统的核心内部功能实现，包括时间戳缓存、调用者信息获取、背压控制、日志消息处理和接口实现等，为 `FastLog` 提供高性能的底层支持。
- **types.go**：日志系统核心类型定义，定义 `FastLog` 的核心数据结构、常量和枚举类型，包括日志级别、日志格式、路径信息和日志消息结构体等。

## VARIABLES

```go
var (
	// New 是 NewFastLog 的简写别名
	//
	// 用法:
	//  - logger := fastlog.New(config)
	//
	// 等价于:
	//  - logger := fastlog.NewFastLog(config)
	New = NewFastLog

	// NewCfg 是 NewFastLogConfig 的简写别名
	//
	// 用法:
	//  - config := fastlog.NewCfg()
	//
	// 等价于:
	//  - config := fastlog.NewFastLogConfig()
	NewCfg = NewFastLogConfig
)
```

为了提供更简洁的 API 调用方式，定义以下函数别名：这样用户可以使用更短的函数名来创建日志实例和配置。

## TYPES

### FastLog

`FastLog` 日志记录器

```go
type FastLog struct {
	// Has unexported fields.
}
```

#### NewFastLog

创建一个新的 `FastLog` 实例，用于记录日志。

```go
func NewFastLog(config *FastLogConfig) *FastLog
```

- **参数**：
  - `config`：一个指向 `FastLogConfig` 实例的指针，用于配置日志记录器。
- **返回值**：
  - `*FastLog`：一个指向 `FastLog` 实例的指针。

#### Close

关闭日志记录器。

```go
func (f *FastLog) Close() error
```

- **返回值**：
  - `error`：如果关闭过程中发生错误，返回错误信息；否则返回 `nil`。

#### Debug

记录调试级别的日志，不支持占位符。

```go
func (f *FastLog) Debug(v ...any)
```

- **参数**：
  - `v`：可变参数，可以是任意类型，会被转换为字符串。

#### Debugf

记录调试级别的日志，支持占位符，格式化。

```go
func (f *FastLog) Debugf(format string, v ...any)
```

- **参数**：
  - `format`：格式字符串。
  - `v`：可变参数，可以是任意类型，会被转换为字符串。

#### Error

记录错误级别的日志，不支持占位符。

```go
func (f *FastLog) Error(v ...any)
```

- **参数**：
  - `v`：可变参数，可以是任意类型，会被转换为字符串。

#### Errorf

记录错误级别的日志，支持占位符，格式化。

```go
func (f *FastLog) Errorf(format string, v ...any)
```

- **参数**：
  - `format`：格式字符串。
  - `v`：可变参数，可以是任意类型，会被转换为字符串。

#### Fatal

记录致命级别的日志，不支持占位符，发送后关闭日志记录器。

```go
func (f *FastLog) Fatal(v ...any)
```

- **参数**：
  - `v`：可变参数，可以是任意类型，会被转换为字符串。

#### Fatalf

记录致命级别的日志，支持占位符，发送后关闭日志记录器。

```go
func (f *FastLog) Fatalf(format string, v ...any)
```

- **参数**：
  - `format`：格式字符串。
  - `v`：可变参数，可以是任意类型，会被转换为字符串。

#### Info

记录信息级别的日志，不支持占位符。

```go
func (f *FastLog) Info(v ...any)
```

- **参数**：
  - `v`：可变参数，可以是任意类型，会被转换为字符串。

#### Infof

记录信息级别的日志，支持占位符，格式化。

```go
func (f *FastLog) Infof(format string, v ...any)
```

- **参数**：
  - `format`：格式字符串。
  - `v`：可变参数，可以是任意类型，会被转换为字符串。

#### Warn

记录警告级别的日志，不支持占位符。

```go
func (f *FastLog) Warn(v ...any)
```

- **参数**：
  - `v`：可变参数，可以是任意类型，会被转换为字符串。

#### Warnf

记录警告级别的日志，支持占位符，格式化。

```go
func (f *FastLog) Warnf(format string, v ...any)
```

- **参数**：
  - `format`：格式字符串。
  - `v`：可变参数，可以是任意类型，会被转换为字符串。

### FastLogConfig

定义一个配置结构体，用于配置日志记录器。

```go
type FastLogConfig struct {
	LogDirName      string        // 日志目录路径
	LogFileName     string        // 日志文件名
	OutputToConsole bool          // 是否将日志输出到控制台
	OutputToFile    bool          // 是否将日志输出到文件
	LogLevel        LogLevel      // 日志级别
	LogFormat       LogFormatType // 日志格式选项
	Color           bool          // 是否启用终端颜色
	Bold            bool          // 是否启用终端字体加粗
	MaxSize         int           // 最大日志文件大小，单位为MB，默认10MB
	MaxAge          int           // 最大日志文件保留天数，默认为0，表示不做限制
	MaxFiles        int           // 最大日志文件保留数量，默认为0，表示不做限制
	LocalTime       bool          // 是否使用本地时间，默认使用UTC时间
	Compress        bool          // 是否启用日志文件压缩，默认不启用
	MaxBufferSize   int           // 缓冲区大小，单位为字节，默认64KB
	MaxWriteCount   int           // 最大写入次数，默认500次
	FlushInterval   time.Duration // 刷新间隔，默认1秒
}
```

#### ConsoleConfig

创建一个控制台环境下的 `FastLogConfig` 实例。

```go
func ConsoleConfig() *FastLogConfig
```

- **返回值**：
  - `*FastLogConfig`：一个指向 `FastLogConfig` 实例的指针。
- **特性**：
  - 禁用文件输出。
  - 设置日志级别为 `DEBUG`。

#### DevConfig

创建一个开发环境下的 `FastLogConfig` 实例。

```go
func DevConfig(logDirName string, logFileName string) *FastLogConfig
```

- **参数**：
  - `logDirName`：日志目录名。
  - `logFileName`：日志文件名。
- **返回值**：
  - `*FastLogConfig`：一个指向 `FastLogConfig` 实例的指针。
- **特性**：
  - 启用详细日志格式。
  - 设置日志级别为 `DEBUG`。
  - 设置最大日志文件保留数量为5。
  - 设置最大日志文件保留天数为7天。

#### NewFastLogConfig

创建一个新的 `FastLogConfig` 实例，用于配置日志记录器。

```go
func NewFastLogConfig(logDirName string, logFileName string) *FastLogConfig
```

- **参数**：
  - `logDirName`：日志目录名称，默认为 `"applogs"`。
  - `logFileName`：日志文件名称，默认为 `"app.log"`。
- **返回值**：
  - `*FastLogConfig`：一个指向 `FastLogConfig` 实例的指针。

#### ProdConfig

创建一个生产环境下的 `FastLogConfig` 实例。

```go
func ProdConfig(logDirName string, logFileName string) *FastLogConfig
```

- **参数**：
  - `logDirName`：日志目录名。
  - `logFileName`：日志文件名。
- **返回值**：
  - `*FastLogConfig`：一个指向`FastLogConfig` 实例的指针。
- **特性**：
  - 启用日志文件压缩。
  - 禁用控制台输出。
  - 设置最大日志文件保留天数为30天。
  - 设置最大日志文件保留数量为24个。

### LogFormatType

日志格式选项。

```go
type LogFormatType int
```

```go
const (
	// 详细格式
	Detailed LogFormatType = iota

	// json格式
	Json

	// json简化格式（无文件信息）
	JsonSimple

	// 简约格式（无文件信息）
	Simple

	// 结构化格式
	Structured

	// 基础结构化格式（无文件信息）
	BasicStructured

	// 简单时间格式（无文件信息）
	SimpleTimestamp

	// 自定义格式（无文件信息）
	Custom
)
```

### LogLevel

日志级别枚举。

```go
type LogLevel uint8
```

```go
const (
	DEBUG LogLevel = 10  // 调试级别
	INFO  LogLevel = 20  // 信息级别
	WARN  LogLevel = 30  // 警告级别
	ERROR LogLevel = 40  // 错误级别
	FATAL LogLevel = 50  // 致命级别
	NONE  LogLevel = 255 // 无日志级别
)
```

#### MarshalJSON

将日志级别转换为字符串。

```go
func (l LogLevel) MarshalJSON() ([]byte, error)
```

- **返回值**：
  - `[]byte`：日志级别的 JSON 字符串。
  - `error`：如果发生错误，返回错误信息。
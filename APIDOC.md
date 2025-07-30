# Package fastlog

Package `fastlog` (import: `"gitee.com/MM-Q/fastlog"`)

## Module Overview

- `config.go` - 日志配置管理模块  
  定义日志配置结构体及配置项的设置与获取方法，负责管理FastLog的所有可配置参数。

- `fastlog.go` - FastLog日志记录器核心实现  
  提供日志记录器的创建、初始化、日志写入及关闭等核心功能，集成配置管理、缓冲区管理和日志处理流程。

- `logger.go` - 日志记录方法实现  
  提供不同级别日志的记录方法（带占位符和不带占位符），实现日志级别过滤和调用者信息获取功能。

- `processor.go` - 单线程日志处理器实现  
  负责从日志通道接收消息、批量缓存，并根据批次大小或时间间隔触发处理，实现日志的批量格式化和输出。

- `tools.go` - 工具函数集合  
  提供路径检查、调用者信息获取、协程ID获取、日志格式化和颜色添加等辅助功能。

- `types.go` - 日志系统核心类型定义  
  定义FastLog的核心数据结构、常量和枚举类型，包括日志级别、日志格式、路径信息和日志消息结构体等。

## Variables

```go
var (
	NewFlog = NewFastLog
	NewFcfg = NewFastLogConfig
)
```

简写实现

## Types

### FastLog

```go
type FastLog struct {
	// Has unexported fields.
}
```

FastLog 日志记录器

#### NewFastLog

```go
func NewFastLog(config *FastLogConfig) (*FastLog, error)
```

NewFastLog 创建一个新的FastLog实例, 用于记录日志。

参数:
- config: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。

返回值:
- *FastLog: 一个指向FastLog实例的指针。
- error: 如果创建日志记录器失败, 则返回一个错误。

#### Close

```go
func (f *FastLog) Close()
```

Close 关闭 FastLog 实例

#### Debug

```go
func (l *FastLog) Debug(v ...any)
```

Debug 记录调试级别的日志，不支持占位符

参数:
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Debugf

```go
func (l *FastLog) Debugf(format string, v ...any)
```

Debugf 记录调试级别的日志，支持占位符，格式化

参数:
- format: 格式字符串
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Error

```go
func (l *FastLog) Error(v ...any)
```

Error 记录错误级别的日志，不支持占位符

参数:
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Errorf

```go
func (l *FastLog) Errorf(format string, v ...any)
```

Errorf 记录错误级别的日志，支持占位符，格式化

参数:
- format: 格式字符串
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Fatal

```go
func (l *FastLog) Fatal(v ...any)
```

Fatal 记录致命级别的日志，不支持占位符，发送后关闭日志记录器

参数:
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Fatalf

```go
func (l *FastLog) Fatalf(format string, v ...any)
```

Fatalf 记录致命级别的日志，支持占位符，发送后关闭日志记录器

参数:
- format: 格式字符串
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Info

```go
func (l *FastLog) Info(v ...any)
```

Info 记录信息级别的日志，不支持占位符

参数:
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Infof

```go
func (l *FastLog) Infof(format string, v ...any)
```

Infof 记录信息级别的日志，支持占位符，格式化

参数:
- format: 格式字符串
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Success

```go
func (l *FastLog) Success(v ...any)
```

Success 记录成功级别的日志，不支持占位符

参数:
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Successf

```go
func (l *FastLog) Successf(format string, v ...any)
```

Successf 记录成功级别的日志，支持占位符，格式化

参数:
- format: 格式字符串
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Warn

```go
func (l *FastLog) Warn(v ...any)
```

Warn 记录警告级别的日志，不支持占位符

参数:
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Warnf

```go
func (l *FastLog) Warnf(format string, v ...any)
```

Warnf 记录警告级别的日志，支持占位符，格式化

参数:
- format: 格式字符串
- v: 可变参数，可以是任意类型，会被转换为字符串

---

### FastLogConfig

```go
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
	MaxLogFileSize  int           // 最大日志文件大小, 单位为MB, 默认5MB
	MaxLogAge       int           // 最大日志文件保留天数, 默认为0, 表示不做限制
	MaxLogBackups   int           // 最大日志文件保留数量, 默认为0, 表示不做限制
	IsLocalTime     bool          // 是否使用本地时间 默认使用UTC时间
	EnableCompress  bool          // 是否启用日志文件压缩 默认不启用
}
```

FastLogConfig 定义一个配置结构体, 用于配置日志记录器

#### NewFastLogConfig

```go
func NewFastLogConfig(logDirName string, logFileName string) *FastLogConfig
```

NewFastLogConfig 创建一个新的FastLogConfig实例, 用于配置日志记录器。

参数:
- logDirName: 日志目录名称, 默认为"applogs"。
- logFileName: 日志文件名称, 默认为"app.log"。

返回值:
- *FastLogConfig: 一个指向FastLogConfig实例的指针。

---

### LogFormatType

```go
type LogFormatType int
```

日志格式选项

#### Constants

```go
const (
	Detailed LogFormatType = iota // 详细格式
	Bracket                       // 方括号格式
	Json                          // json格式
	Threaded                      // 协程格式
	Simple                        // 简约格式
	Custom                        // 自定义格式
)
```

日志格式选项

---

### LogLevel

```go
type LogLevel int
```

日志级别枚举

#### Constants

```go
const (
	DEBUG   LogLevel = 10  // 调试级别
	INFO    LogLevel = 20  // 信息级别
	SUCCESS LogLevel = 30  // 成功级别
	WARN    LogLevel = 40  // 警告级别
	ERROR   LogLevel = 50  // 错误级别
	FATAL   LogLevel = 60  // 致命级别
	NONE    LogLevel = 999 // 无日志级别
)
```

定义日志级别

---

### PathInfo

```go
type PathInfo struct {
	Path    string      // 路径
	Exists  bool        // 是否存在
	IsFile  bool        // 是否为文件
	IsDir   bool        // 是否为目录
	Size    int64       // 文件大小（字节）
	Mode    os.FileMode // 文件权限
	ModTime time.Time   // 文件修改时间
}
```

PathInfo 是一个结构体，用于封装路径的信息
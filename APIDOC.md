# Package fastlog

Package fastlog 提供了一个高性能的日志记录器，适用于快速记录日志。

导入路径：gitee.com/MM-Q/fastlog

## FastLog

日志记录器

```go
type FastLog struct {
	// Has unexported fields.
}
```

## NewFastLog

创建一个新的FastLog实例，用于记录日志。

```go
func NewFastLog(config *FastLogConfig) (*FastLog, error)
```

参数：

- `config`：一个指向FastLogConfig实例的指针，用于配置日记录志器。

返回值：

- `*FastLog`：一个指向FastLog实例的指针。
- `error`：如果创建日志记录器失败，则返回一个错误。

## Close

关闭FastLog实例，并等待所有日志处理完成。

```go
func (f *FastLog) Close() error
```

## Debug

记录调试级别的日志，不支持占位符。

```go
func (l *FastLog) Debug(v ...any)
```

## Debugf

记录调试级别的日志，支持占位符，格式化。

```go
func (l *FastLog) Debugf(format string, v ...any)
```

## Error

记录错误级别的日志，不支持占位符。

```go
func (l *FastLog) Error(v ...any)
```

## Errorf

记录错误级别的日志，支持占位符，格式化。

```go
func (l *FastLog) Errorf(format string, v ...any)
```

## Fatal

记录致命级别的日志，不支持占位符，发送后关闭日志记录器。

```go
func (l *FastLog) Fatal(v ...any)
```

## Fatalf

记录致命级别的日志，支持占位符，发送后关闭日志记录器。

```go
func (l *FastLog) Fatalf(format string, v ...any)
```

## Info

记录信息级别的日志，不支持占位符。

```go
func (l *FastLog) Info(v ...any)
```

## Infof

记录信息级别的日志，支持占位符，格式化。

```go
func (l *FastLog) Infof(format string, v ...any)
```

## Success

记录成功级别的日志，不支持占位符。

```go
func (l *FastLog) Success(v ...any)
```

## Successf

记录成功级别的日志，支持占位符，格式化。

```go
func (l *FastLog) Successf(format string, v ...any)
```

## Warn

记录警告级别的日志，不支持占位符。

```go
func (l *FastLog) Warn(v ...any)
```

## Warnf

记录警告级别的日志，支持占位符，格式化。

```go
func (l *FastLog) Warnf(format string, v ...any)
```

## FastLogConfig

定义一个配置结构体，用于配置日志记录器。

```go
type FastLogConfig struct {
	// Has unexported fields.
}
```

## NewFastLogConfig

创建一个新的FastLogConfig实例，用于配置日志记录器。

```go
func NewFastLogConfig(logDirName string, logFileName string) *FastLogConfig
```

参数：

- `logDirName`：日志目录名称，默认为"applogs"。
- `logFileName`：日志文件名称，默认为"app.log"。

返回值：

- `*FastLogConfig`：一个指向FastLogConfig实例的指针。

## FastLogConfigurer

定义日志配置器接口，包含所有配置项的设置和获取方法。

```go
type FastLogConfigurer interface {
	// SetLogDirName 设置日志目录路径
	SetLogDirName(dirName string)
	// GetLogDirName 获取日志目录路径
	GetLogDirName() string

	// SetLogFileName 设置日志文件名
	SetLogFileName(fileName string)
	// GetLogFileName 获取日志文件名
	GetLogFileName() string

	// SetPrintToConsole 设置是否将日志输出到控制台
	SetPrintToConsole(print bool)
	// GetPrintToConsole 获取是否将日志输出到控制台的状态
	GetPrintToConsole() bool

	// SetConsoleOnly 设置是否仅输出到控制台
	SetConsoleOnly(only bool)
	// GetConsoleOnly 获取是否仅输出到控制台的状态
	GetConsoleOnly() bool

	// SetFlushInterval 设置刷新间隔
	SetFlushInterval(interval time.Duration)
	// GetFlushInterval 获取刷新间隔
	GetFlushInterval() time.Duration

	// SetLogLevel 设置日志级别
	SetLogLevel(level LogLevel)
	// GetLogLevel 获取日志级别
	GetLogLevel() LogLevel

	// SetChanIntSize 设置通道大小
	SetChanIntSize(size int)
	// GetChanIntSize 获取通道大小
	GetChanIntSize() int

	// SetLogFormat 设置日志格式选项
	SetLogFormat(format LogFormatType)
	// GetLogFormat 获取日志格式选项
	GetLogFormat() LogFormatType

	// SetMaxBufferSize 设置最大缓冲区大小(MB)
	SetMaxBufferSize(size int)
	// GetMaxBufferSize 获取最大缓冲区大小(MB)
	GetMaxBufferSize() int

	// SetNoColor 设置是否禁用终端颜色
	SetNoColor(noColor bool)
	// GetNoColor 获取是否禁用终端颜色的状态
	GetNoColor() bool

	// SetNoBold 设置是否禁用终端字体加粗
	SetNoBold(noBold bool)
	// GetNoBold 获取是否禁用终端字体加粗的状态
	GetNoBold() bool

	// SetMaxLogFileSize 设置最大日志文件大小(MB)
	SetMaxLogFileSize(size int)
	// GetMaxLogFileSize 获取最大日志文件大小(MB)
	GetMaxLogFileSize() int

	// SetMaxLogAge 设置最大日志文件保留天数
	SetMaxLogAge(age int)
	// GetMaxLogAge 获取最大日志文件保留天数
	GetMaxLogAge() int

	// SetMaxLogBackups 设置最大日志文件保留数量
	SetMaxLogBackups(backups int)
	// GetMaxLogBackups 获取最大日志文件保留数量
	GetMaxLogBackups() int

	// SetIsLocalTime 设置是否使用本地时间
	SetIsLocalTime(local bool)
	// GetIsLocalTime 获取是否使用本地时间的状态
	GetIsLocalTime() bool

	// SetEnableCompress 设置是否启用日志文件压缩
	SetEnableCompress(compress bool)
	// GetEnableCompress 获取是否启用日志文件压缩的状态
	GetEnableCompress() bool
}
```

## FastLogInterface

定义一个接口，声明对外暴露的方法。

```go
type FastLogInterface interface {
	Close() // 关闭日志记录器

	Info(v ...any)    // 记录信息级别的日志，不支持占位符
	Warn(v ...any)    // 记录警告级别的日志，不支持占位符
	Error(v ...any)   // 记录错误级别的日志，不支持占位符
	Success(v ...any) // 记录成功级别的日志，不支持占位符
	Debug(v ...any)   // 录调试级别的日志，不支持占位符
	Fatal(v ...any)   // 记录致命级别的日志，不支持占位符

	Infof(format string, v ...any)    // 记录信息级别的日志，支持占位符，格式化
	Warnf(format string, v ...any)    // 记录警告级别的日志，支持占位符，格式化
	Errorf(format string, v ...any)   // 记录错误级别的日志，支持占位符，格式化
	Successf(format string, v ...any) // 记录成功级别的日志，支持占位符，格式化
	Debugf(format string, v ...any)   // 记录调试级别的日志，支持占位符，格式化
	Fatalf(format string, v ...any)   // 记录致命级别的日志，支持占位符，格式化
}
```

## LogFormatType

日志格式选项。

```go
type LogFormatType int
```

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

## LogLevel

日志级别枚举。

```go
type LogLevel int
```

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

## PathInfo

PathInfo 是一个结构体，用于封装路径的信息。

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

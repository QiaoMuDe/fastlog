# FastLog

`fastlog` 是一个高性能、灵活的日志记录库，旨在为 Go 语言项目提供高效、可配置的日志记录功能。它支持多种日志格式、日志级别过滤、异步日志记录、缓冲区写入、颜色输出到控制台等功能。

## 特点

- 支持多种日志格式（`Json(json格式)`、`Bracket(括号格式)`、`Detailed(详细格式)`、`Threaded(协程格式)`）。
- 支持分级日志（Debug、Info、Warn、Error、Success）。
- 异步日志记录，使用通道和缓冲区提高性能。
- 支持将日志输出到文件和控制台，并且可以仅输出到控制台。
- 支持日志缓冲区，减少 I/O 操作。
- 控制台输出支持渲染颜色，方便区分不同级别的日志。
- 支持日志级别过滤，可以根据需要记录不同级别的日志。
- 提供了格式化和非格式化的日志记录方法。
- 支持日志轮转，支持按日期和大小分割日志文件。
- 支持日志轮转的时候日志压缩，支持的格式有：`tgz`, `tar`, `gz`, `zip`。

## 日志级别与颜色映射

| 日志级别 | 颜色 |
| -------- | ---- |
| DEBUG    | 紫色 |
| INFO     | 蓝色 |
| WARN     | 黄色 |
| ERROR    | 红色 |
| SUCCESS  | 绿色 |

## 结构体方法和函数

### 函数

| 函数名称  | 参数类型                                  | 返回值类型           | 说明                                                 |
| --------- | ----------------------------------------- | -------------------- | ---------------------------------------------------- |
| NewConfig | `logDirName string, logFileName string` | `*FastLogConfig` | 创建一个日志配置器，日志目录和日志文件名为必需参数。 |
| NewFastLog | `cfg *FastLogConfig`              | `(*FastLog, error)` | 根据配置创建一个新的日志记录器。                     |

### 结构体

以下是将 `FastLog` 结构体的字段及其说明：

> 私有属性：模块内部使用无需修改

| 属性名称      | 类型           | 说明                                                      |
| ------------- | -------------- | --------------------------------------------------------- |
| logFile       | *os.File       | 日志文件句柄，用于操作日志文件。                          |
| logFilePath   | string         | 日志文件路径，由 `logDirName` 和 `logFileName` 拼接而成。 |
| logger        | *log.Logger    | 底层日志记录器，用于实现日志的基本功能。                  |
| fileMu        | sync.Mutex     | 文件写入的互斥锁，用于保证文件写入操作的线程安全。        |
| consoleMu     | sync.Mutex     | 控制台写入的互斥锁，用于保证控制台输出的线程安全。        |
| rmLogMu       | sync.Mutex     | 删除日志的互斥锁，用于保证删除日志操作的线程安全。        |
| isCleaning    | bool           | 标志变量，表示是否已经在清理日志文件。                    |
| logChan       | chan string    | 日志通道，用于传递日志消息。                              |
| stopChan      | chan string    | 停止通道，用于控制日志记录器的停止操作。                  |
| wg            | sync.WaitGroup | 等待组，用于等待日志通道中的日志被处理。                  |
| wgz           | sync.WaitGroup | 等待组，用于等待日志文件被压缩。                          |
| fileBuffer    | *bytes.Buffer  | 文件缓冲区，用于暂存日志数据。                            |
| consoleBuffer | *bytes.Buffer  | 控制台缓冲区，用于暂存控制台输出数据。                    |
| ticker        | *time.Ticker   | 定时器，用于触发日志轮转等定时操作。                      |
| fileWriter    | io.Writer      | 文件写入器，用于将日志写入文件。                          |
| consoleWriter | io.Writer      | 控制台写入器，用于将日志输出到控制台。                    |
| closeOnce     | sync.Once      | 确保通道只被关闭一次。                                    |
| logSizeBytes  | int64          | 日志文件大小字节数，用于比较大小。                        |

> 公共属性：对外暴露的属性，可以通过cfg配置结构体自定义

| 属性名称          | 类型          | 说明                                          |
| ----------------- | ------------- | --------------------------------------------- |
| logDirName        | string        | 日志目录，用于指定日志文件存储的目录。        |
| logFileName       | string        | 日志文件名。                                  |
| printToConsole    | bool          | 是否将日志输出到控制台。                      |
| consoleOnly       | bool          | 是否仅输出到控制台。                          |
| logLevel          | LogLevel      | 日志级别，用于控制日志的输出级别。            |
| chanIntSize       | int           | 通道大小，默认为 1000。                       |
| bufferKbSize      | int           | 缓冲区大小，默认为 1024 KB。                  |
| logFormat         | LogFormatType | 日志格式选项，如 Json、Bracket、Detailed 等。 |
| enableLogRotation | bool          | 是否启用日志切割。                            |
| logRetentionDays  | int           | 日志保留天数，默认为 7 天。                   |
| logMaxSize        | string        | 日志文件最大大小，默认为 3MB。                |
| logRetentionCount | int           | 日志文件保留数量，默认为 3 个。               |
| rotationInterval  | int           | 日志轮转的间隔时间，默认为 10 分钟。          |
| enableCompression | bool          | 是否启用日志压缩。                            |
| compressionFormat | string        | 日志压缩格式，如 zip、gz、tar ，tgz等。       |

以下是将 `FastLogConfig` 结构体的字段及其说明：

> 配置项属性

| 属性名称          | 类型          | 说明                                                         |
| ----------------- | ------------- | ------------------------------------------------------------ |
| LogDirName        | string        | 日志目录名称，指定日志文件存储的目录。                       |
| LogFileName       | string        | 日志文件名称。                                               |
| LogPath           | string        | 日志文件路径，由 `LogDirName` 和 `LogFileName` 拼接而成，无需手动提供。 |
| PrintToConsole    | bool          | 是否将日志输出到控制台。                                     |
| ConsoleOnly       | bool          | 是否仅输出到控制台。                                         |
| LogLevel          | LogLevel      | 日志级别，用于控制日志的输出级别。默认是：Info                           |
| ChanIntSize       | int           | 通道大小，用于设置日志通道的容量。 默认是：1000                          |
| BufferKbSize      | int           | 缓冲区大小，单位为 KB。                                      |
| LogFormat         | LogFormatType | 日志格式选项，如 Json、Bracket、Detailed、Threaded 等。      |
| EnableLogRotation | bool          | 是否启用日志切割功能。                                       |
| LogRetentionDays  | int           | 日志保留天数，默认为 7 天。                                  |
| LogMaxSize        | string        | 日志文件最大大小，默认为 3MB，单位可以是 KB、MB 或 GB。      |
| LogRetentionCount | int           | 日志文件保留数量，默认为 3 个。                              |
| EnableCompression | bool          | 是否启用日志压缩功能。                                       |
| RotationInterval  | int           | 日志轮转的间隔时间，默认为 10 分钟，单位为秒。               |
| CompressionFormat | string        | 日志压缩格式，支持 zip、gz、tar、tgz 等，默认为 zip。        |

> 说明

- **LogPath** 是由 `LogDirName` 和 `LogFileName` 拼接而成的，通常不需要用户手动提供，而是由程序内部生成。
- **EnableLogRotation**、**EnableCompression** 等布尔类型的配置项，用于控制日志功能的启用或禁用。
- **LogMaxSize** 和 **RotationInterval** 等配置项，提供了灵活的自定义选项，以满足不同场景下的日志管理需求。


### 方法

以下是将 `FastLogInterface` 中的方法及其说明：

| 方法名称 | 参数类型                            | 说明                                             |
| -------- | ----------------------------------- | ------------------------------------------------ |
| Info     | `v ...interface{}`                | 记录信息级别的日志，不支持占位符，需要自己拼接。 |
| Warn     | `v ...interface{}`                | 记录警告级别的日志，不支持占位符，需要自己拼接。 |
| Error    | `v ...interface{}`                | 记录错误级别的日志，不支持占位符，需要自己拼接。 |
| Success  | `v ...interface{}`                | 记录成功级别的日志，不支持占位符，需要自己拼接。 |
| Debug    | `v ...interface{}`                | 记录调试级别的日志，不支持占位符，需要自己拼接。 |
| Close    | 无                                  | 关闭日志记录器。                                 |
| Infof    | `format string, v ...interface{}` | 记录信息级别的日志，支持占位符，格式化。         |
| Warnf    | `format string, v ...interface{}` | 记录警告级别的日志，支持占位符，格式化。         |
| Errorf   | `format string, v ...interface{}` | 记录错误级别的日志，支持占位符，格式化。         |
| Successf | `format string, v ...interface{}` | 记录成功级别的日志，支持占位符，格式化。         |
| Debugf   | `format string, v ...interface{}` | 记录调试级别的日志，支持占位符，格式化。         |

## 日志格式选项

根据 `logFormat` 的值，日志格式不同：

- **Json**:

  ```json
  {"time":"2006-01-02 15:04:05","level":"INFO","file":"filename","function":"funcName","line":123, "thread":"1234","message":"log message"}
  ```
- **Bracket**:

  ```bash
  [INFO] log message
  ```
- **Detailed**:

  ```bash
  2006-01-02 15:04:05 | INFO     | filename:funcName:123 - log message
  ```
- **Threaded**:

  ```bash
  2006-01-02 15:04:05 | INFO     | [thread="1234"] log message
  ```

## 使用方法

### 下载与引入

```bash
# 确保在自己项目路径下，并且存在go.mog文件，不存在则 go init 项目名 创建
go get gitee.com/MM-Q/fastlog
```

在代码中引入：

```go
import "gitee.com/MM-Q/fastlog"
```

### 示例代码

```go
package main

import (
	"gitee.com/MM-Q/fastlog"
	"fmt"
)

func main() {
	// 创建日志配置
	cfg := fastlog.NewConfig("logs", "app.log")

	// 设置日志级别
	cfg.LogLevel = fastlog.Debug

	// 设置日志格式
	cfg.LogFormat = fastlog.Json

	// 设置日志轮转
	cfg.EnableLogRotation = true
	cfg.LogRetentionDays = 7
	cfg.LogMaxSize = "5MB"
	cfg.RotationInterval = 600 // 10分钟

	// 设置日志压缩
	cfg.EnableCompression = true
	cfg.CompressionFormat = "zip"

	// 设置控制台输出
	cfg.PrintToConsole = true

	// 创建日志记录器
	logger, err := fastlog.NewFastLog(cfg)
	if err != nil {
		fmt.Println("创建日志记录器失败:", err)
	}
	defer logger.Close() // 记得在程序结束时关闭日志记录器

	// 记录不同级别的日志
	logger.Info("This is an info message")
	logger.Warn("This is a warn message")
	logger.Error("This is an error message")
	logger.Success("This is a success message")
	logger.Debug("This is a debug message")

	// 使用格式化日志
	logger.Infof("This is an info message with format: %s", "formatted text")

	// 使用调试日志
	logger.Debug("This is a debug message with format: %s", "formatted text")
}
```

## 更多自定义配置

### 自定义级别

```go
cfg.LogLevel = fastlog.Debug
```

### 自定义格式

```go
cfg.LogFormat = fastlog.Json
```

### 自定义轮转

```go
cfg.EnableLogRotation = true
cfg.LogRetentionDays = 7
cfg.LogMaxSize = "5MB"
cfg.RotationInterval = 600 // 10分钟
```

### 自定义压缩

```go
cfg.EnableCompression = true
cfg.CompressionFormat = "zip"
```

### 自定义控制台输出

```go
cfg.PrintToConsole = true
```

### 自定义仅输出到控制台

```go
cfg.ConsoleOnly = true
```

## 注意事项

- 日志文件路径：日志文件路径由 `logDirName` 和 `logFileName` 拼接而成，无需手动提供。
- 日志轮转：日志轮转功能会在指定的时间间隔内自动创建新的日志文件，并删除旧的日志文件。
- 日志压缩：日志压缩功能会在日志轮转时。

## 贡献

欢迎大佬们贡献代码，提问题，提建议，一起共同交流，实现日志简易化。
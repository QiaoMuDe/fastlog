# FastLog

`fastlog` 是一个高性能、灵活的日志记录库，旨在为 Go 语言项目提供高效、可配置的日志记录功能。它支持多种日志格式、日志级别过滤、异步日志记录、缓冲区写入、颜色输出到控制台等功能。

## 特点

- 支持多种日志格式（`Json(json格式)`、`Bracket(括号格式)`、`Detailed(详细格式)`、`Threaded(协程格式)`）。
- 支持分级日志（Debug、Info、Warn、Error、Success）。
- 异步日志记录，使用通道和缓冲区提高性能。
- 支持将日志输出到文件和控制台，并且可以仅输出到控制台。
- 支持日志缓冲区，减少 I/O 操作。
- 控制台输出带有颜色，方便区分不同级别的日志。
- 支持日志级别过滤，可以根据需要记录不同级别的日志。
- 提供了格式化和非格式化的日志记录方法。

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

| 函数名称  | 参数类型                                | 返回值类型         | 说明                                                 |
| --------- | --------------------------------------- | ------------------ | ---------------------------------------------------- |
| NewConfig | `logDirName string, logFileName string` | `LoggerConfig`     | 创建一个日志配置器，日志目录和日志文件名为必需参数。 |
| NewLogger | `cfg LoggerConfig`                      | `(*Logger, error)` | 根据配置创建一个新的日志记录器。                     |

### 结构体

以下是将 `Logger` 结构体的字段及其说明：

| 字段名称       | 类型             | 说明                                                      |
| -------------- | ---------------- | --------------------------------------------------------- |
| logDirName     | `string`         | 日志目录路径。                                            |
| logFileName    | `string`         | 日志文件名。                                              |
| logFile        | `*os.File`       | 日志文件句柄，用于操作日志文件。                          |
| logger         | `*log.Logger`    | 底层日志记录器，用于基础的日志操作。                      |
| printToConsole | `bool`           | 是否将日志输出到控制台。                                  |
| consoleOnly    | `bool`           | 是否仅将日志输出到控制台，而不写入文件。                  |
| logLevel       | `LogLevel`       | 当前日志级别，用于过滤日志输出。                          |
| fileMu         | `sync.Mutex`     | 文件写入的互斥锁，用于保证文件写入的线程安全。            |
| consoleMu      | `sync.Mutex`     | 控制台写入的互斥锁，用于保证控制台输出的线程安全。        |
| logChan        | `chan string`    | 日志通道，用于异步传递日志消息。                          |
| stopChan       | `chan string`    | 停止通道，用于发送停止信号以关闭日志记录器。              |
| wg             | `sync.WaitGroup` | 等待组，用于等待日志通道中的日志被处理完成。              |
| fileBuffer     | `*bytes.Buffer`  | 文件缓冲区，用于暂存日志内容以减少文件 I/O 操作。         |
| consoleBuffer  | `*bytes.Buffer`  | 控制台缓冲区，用于暂存日志内容以减少控制台输出操作。      |
| ticker         | `*time.Ticker`   | 定时器，用于定期刷新缓冲区内容到文件或控制台。            |
| fileWriter     | `io.Writer`      | 文件写入器，用于将日志内容写入文件。                      |
| consoleWriter  | `io.Writer`      | 控制台写入器，用于将日志内容输出到控制台。                |
| chanIntSize    | `int`            | 日志通道的大小，控制通道的缓存能力。                      |
| bufferKbSize   | `int`            | 缓冲区大小（单位：KB），控制文件和控制台缓冲区的大小。    |
| closeOnce      | `sync.Once`      | 确保通道只被关闭一次，防止重复关闭。                      |
| logFormat      | `LogFormatType`  | 日志格式选项，支持 JSON、详细格式、括号格式、协程格式等。 |

以下是将 `LoggerConfig` 结构体的字段及其说明：

| 字段名称       | 类型            | 说明                                                         |
| -------------- | --------------- | ------------------------------------------------------------ |
| LogDirName     | `string`        | 日志目录名称，指定日志文件存储的目录路径。                   |
| LogFileName    | `string`        | 日志文件名称，指定日志文件的文件名。                         |
| LogPath        | `string`        | 日志文件路径，由 `LogDirName` 和 `LogFileName` 拼接生成。    |
| PrintToConsole | `bool`          | 是否将日志输出到控制台，默认为 `true`。                      |
| ConsoleOnly    | `bool`          | 是否仅将日志输出到控制台，而不写入文件，默认为 `false`。     |
| LogLevel       | `LogLevel`      | 日志级别，用于过滤日志输出，支持 `Debug`、`Info`、`Warn`、`Error`、`Success` 等。 |
| ChanIntSize    | `int`           | 日志通道的大小，控制通道的缓存能力，默认为 `1000`。          |
| BufferKbSize   | `int`           | 缓冲区大小（单位：KB），控制文件和控制台缓冲区的大小，默认为 `1024`（1MB）。 |
| LogFormat      | `LogFormatType` | 日志格式选项，支持以下格式：<br> - `Json`：JSON 格式<br> - `Bracket`：括号格式<br> - `Detailed`：详细格式<br> - `Threaded`：协程格式 |

### 方法

以下是将 `LoggerInterface` 中的方法及其说明：

| 方法名称 | 参数类型                          | 说明                                             |
| -------- | --------------------------------- | ------------------------------------------------ |
| Info     | `v ...interface{}`                | 记录信息级别的日志，不支持占位符，需要自己拼接。 |
| Warn     | `v ...interface{}`                | 记录警告级别的日志，不支持占位符，需要自己拼接。 |
| Error    | `v ...interface{}`                | 记录错误级别的日志，不支持占位符，需要自己拼接。 |
| Success  | `v ...interface{}`                | 记录成功级别的日志，不支持占位符，需要自己拼接。 |
| Debug    | `v ...interface{}`                | 记录调试级别的日志，不支持占位符，需要自己拼接。 |
| Close    | 无                                | 关闭日志记录器。                                 |
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
  [2006-01-02 15:04:05] | INFO     | [thread="1234"] log message
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
)

func main() {
	// 创建日志配置
	config := fastlog.NewConfig("logs", "app.log")

	// 创建日志记录器
	logger, err := fastlog.NewLogger(config)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// 记录不同级别的日志
	logger.Info("This is an info message")
	logger.Warn("This is a warn message")
	logger.Error("This is an error message")
	logger.Success("This is a success message")
	logger.Debug("This is a debug message")

	// 使用格式化日志
	logger.Infof("This is an info message with format: %s", "formatted text")
}
```

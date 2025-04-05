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

| 函数名称         | 参数类型                                  | 返回值类型            | 说明                                                 |
| ---------------- | ----------------------------------------- | --------------------- | ---------------------------------------------------- |
| NewFastLogConfig | `logDirPath string, logFileName string` | `*FastLogConfig`    | 创建一个日志配置器，日志目录和日志文件名为必需参数。 |
| NewFastLog       | `cfg *FastLogConfig`                    | `(*FastLog, error)` | 根据配置创建一个新的日志记录器。                     |

### 结构体

以下是将 `FastLogConfig` 结构体的字段及其说明：

> 配置项属性

| 属性名称       | 类型          | 说明                                                    |
| -------------- | ------------- | ------------------------------------------------------- |
| LogDirPath     | string        | 日志目录名称，指定日志文件存储的目录。                  |
| LogFileName    | string        | 日志文件名称。                                          |
| PrintToConsole | bool          | 是否将日志输出到控制台。                                |
| ConsoleOnly    | bool          | 是否仅输出到控制台。                                    |
| LogLevel       | LogLevel      | 日志级别，用于控制日志的输出级别。默认INFO              |
| ChanIntSize    | int           | 通道大小，用于设置日志通道的容量。 默认是：1000         |
| LogFormat      | LogFormatType | 日志格式选项，如 Json、Bracket、Detailed、Threaded 等。 |
| MaxBufferSize  | int           | 最大缓冲区大小，单位为 MB。                             |

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

## 注意事项

- 日志文件路径：日志文件路径由 `logDirName` 和 `logFileName` 拼接而成，无需手动提供。
- 日志级别：日志级别用于控制日志的输出级别，默认是 `INFO`。
- 日志格式：日志格式选项，如 `Json`、`Bracket`、`Detailed`、`Threaded` 等。
- 缓冲区大小：缓冲区大小用于控制日志写入的缓冲区大小，单位为 KB。

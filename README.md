# fastlog 包使用文档

## 包功能介绍

`fastlog` 包是一个高性能、灵活的日志记录库，旨在为 Go 语言应用程序提供快速、可靠的日志记录功能。它支持多种日志级别，允许将日志输出到文件和控制台，并具备以下特点：

- **异步日志记录**：通过内部通道和 goroutine 实现异步写入，减少对主程序的阻塞，提高日志记录效率。
- **缓冲区机制**：设置文件和控制台缓冲区，当缓冲区达到一定大小或定时触发时，批量写入日志，减少 I/O 操作次数。
- **日志级别控制**：提供 Debug、Info、Warn、Error、Success 五种日志级别，以及 None 级别用于关闭所有日志记录，方便根据需要筛选日志输出。
- **控制台彩色输出**：为不同级别的日志在控制台输出时添加对应的颜色，便于区分日志类型。
- **灵活配置**：通过 `LoggerConfig` 结构体提供丰富的配置选项，如日志目录、文件名、是否输出到控制台、日志级别、通道大小、缓冲区大小等，方便用户根据项目需求进行定制。

## 接口对外暴露的方法和函数

### LoggerConfig 结构体

用于配置日志记录器，包含以下字段：

| 字段名称         | 说明                                                         | 默认值    |
| :--------------- | :----------------------------------------------------------- | :-------- |
| `LogDirName`     | 日志目录名称                                                 | "logs"    |
| `LogFileName`    | 日志文件名称                                                 | "app.log" |
| `LogPath`        | 日志文件路径，由 `LogDirName` 和 `LogFileName` 拼接生成，无需手动修改 | 无        |
| `PrintToConsole` | 是否将日志输出到控制台                                       | true      |
| `LogLevel`       | 日志级别，可选值为 Debug、Info、Warn、Error、Success、None   | Info      |
| `ChanIntSize`    | 日志通道大小                                                 | 1000      |
| `BufferKbSize`   | 缓冲区大小（KB）                                             | 1MB       |

### Logger 结构体

日志记录器的核心结构体，包含以下方法：

| 方法名称                    | 说明                                             |
| :-------------------------- | :----------------------------------------------- |
| `Info(v ...interface{})`    | 记录信息级别的日志                               |
| `Warn(v ...interface{})`    | 记录警告级别的日志                               |
| `Error(v ...interface{})`   | 记录错误级别的日志                               |
| `Success(v ...interface{})` | 记录成功级别的日志                               |
| `Debug(v ...interface{})`   | 记录调试级别的日志                               |
| `Close()`                   | 关闭日志记录器，确保所有日志被正确写入并释放资源 |

### 其他函数

| 函数名称                                                     | 说明                                                         |
| :----------------------------------------------------------- | :----------------------------------------------------------- |
| `DefaultConfig(logDirName string, logFileName string) LoggerConfig` | 创建一个带有默认配置的日志配置器，需提供日志目录名称和日志文件名称 |
| `NewLogger(cfg LoggerConfig) (*Logger, error)`               | 根据提供的 `LoggerConfig` 创建一个新的日志记录器实例         |

## 使用流程演示

### 1. 引入包

在你的 Go 项目中，首先需要引入 `fastlog` 包：

```go
import "gitee.com/MM-Q/fastlog"
```

### 2. 创建日志配置器

使用 `DefaultConfig` 函数创建一个带有默认配置的日志配置器，然后根据需要修改配置：

```go
cfg := fastlog.DefaultConfig("logs", "app.log")
cfg.PrintToConsole = true  // 确保日志输出到控制台
cfg.LogLevel = fastlog.Debug  // 设置日志级别为 Debug
```

### 3. 创建日志记录器

通过 `NewLogger` 函数根据配置创建日志记录器：

```go
logger, err := fastlog.NewLogger(cfg)
if err != nil {
    fmt.Println("创建日志记录器失败: ", err)
    os.Exit(1)
}
```

### 4. 记录日志

使用日志记录器的各个方法记录不同级别的日志：

```go
logger.Debug("这是一条调试信息")
logger.Info("这是一条普通信息")
logger.Warn("这是一条警告信息")
logger.Error("这是一条错误信息")
logger.Success("这是一条成功信息")
```

### 5. 关闭日志记录器

在程序结束前，调用 `Close` 方法关闭日志记录器，确保所有日志被正确写入：

```go
logger.Close()
```

## 优质开源项目
强烈推荐好大哥全栈开发项目
- gitee：https://gitee.com/pixelmax/gin-vue-admin
- github：https://github.com/flipped-aurora/gin-vue-admin.git
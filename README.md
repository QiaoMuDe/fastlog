# FastLog - 高性能Go日志库

![License](https://img.shields.io/badge/license-GPL-blue.svg)

`FastLog`是一个高性能、异步、可扩展的Go语言日志库，支持多种日志格式和级别，提供文件和控制台双输出。

## 功能特性

- 🚀 异步非阻塞日志记录
- 📁 支持文件和控制台双输出
- 🎨 内置彩色日志输出
- 🔍 多级别日志(DEBUG/INFO/SUCCESS/WARN/ERROR)
- 📝 多种日志格式(JSON/详细/方括号/协程)
- ⏱️ 定时自动刷新缓冲区
- 🛡️ 线程安全设计

## 安装与引入

```bash
# 确保在自己项目路径下，并且存在go.mog文件，不存在则 go init 项目名 创建
go get gitee.com/MM-Q/fastlog

# 引入
import "gitee.com/MM-Q/fastlog"
```

## 快速开始

```go
package main

import "gitee.com/MM-Q/fastlog"

func main() {
    // 初始化配置
    config := &fastlog.FastLogConfig{
        LogDirPath:     "./logs",
        LogFileName:    "app.log",
        PrintToConsole: true,
        LogLevel:       fastlog.DEBUG,
        LogFormat:      fastlog.Detailed,
        MaxBufferSize:  1, // MB
    }

    // 创建日志实例
    logger, err := fastlog.NewFastLog(config)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    // 记录日志
    logger.Info("这是一条信息日志")
    logger.Debugf("调试信息: %s", "value")
    logger.Error("发生了一个错误")
}
```
## 配置选项

| 属性名称 | 类型 | 说明 | 默认值 |
|----------------|---------------|-------------------------------------------------------|---------| 
| LogDirPath | string | 日志目录路径 | 必填 | 
| LogFileName | string | 日志文件名 | 必填 | 
| PrintToConsole | bool | 是否输出到控制台 | false | 
| ConsoleOnly | bool | 是否仅输出到控制台 | false | | LogLevel | LogLevel | 日志级别(DEBUG/INFO/SUCCESS/WARN/ERROR/None) | INFO | 
| ChanIntSize | int | 日志通道缓冲区大小 | 1000 | 
| LogFormat | LogFormatType | 日志格式(Json/Bracket/Detailed/Threaded) | Detailed| 
| MaxBufferSize | int | 最大缓冲区大小(MB) | 1 |

## 日志级别

| 级别 | 值 | 说明 | 
|---------|-----|------------| 
| DEBUG | 10 | 调试信息 | 
| INFO | 20 | 普通信息 | 
| SUCCESS | 30 | 成功信息 | 
| WARN | 40 | 警告信息 | 
| ERROR | 50 | 错误信息 | 
| None | 999 | 不记录任何日志 |

## 日志格式

1. **Detailed** (默认)

```ini
2023-01-01 12:00:00 | INFO    | main.go:main:10 - 日志信息
```

2. **Json**

```json
{"time":"2023-01-01 12:00:00","level":"INFO","file":"main.go","function":"main","line":"10","thread":"1","message":"日志信息"}
```

3. **Bracket**

```ini
[INFO] 日志信息
```

4. **Threaded**

```ini
2023-01-01 12:00:00 | INFO    | [thread="1"] 日志信息
```

## 性能优化

- 异步处理：所有日志操作通过channel异步处理
- 缓冲区：使用内存缓冲区减少IO操作
- 批量写入：定时刷新缓冲区

## 函数

| 函数名称         | 参数类型                                | 返回值类型          | 说明                                                 |
| ---------------- | --------------------------------------- | ------------------- | ---------------------------------------------------- |
| NewFastLogConfig | `logDirPath string, logFileName string` | `*FastLogConfig`    | 创建一个日志配置器，日志目录和日志文件名为必需参数。 |
| NewFastLog       | `cfg *FastLogConfig`                    | `(*FastLog, error)` | 根据配置创建一个新的日志记录器。                     |

## 方法

以下是将 `FastLogInterface` 中的方法及其说明：

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

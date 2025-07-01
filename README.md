# FastLog - 高性能 Go 日志库

![License](https://img.shields.io/badge/license-GPL-blue.svg)

`FastLog`是一个高性能、异步、可扩展的 Go 语言日志库，支持多种日志格式和级别，提供文件和控制台双输出。

## DeepWiKi介绍

- 项目地址：[https://deepwiki.com/QiaoMuDe/fastlog](https://deepwiki.com/QiaoMuDe/fastlog)

## 托管平台

- Gitee：[https://gitee.com/MM-Q/fastlog](https://gitee.com/MM-Q/fastlog)
- Github：[https://github.com/QiaoMuDe/fastlog](https://github.com/QiaoMuDe/fastlog)

## 功能特性

- 🚀 异步非阻塞日志记录
- 📁 支持文件和控制台双输出
- 🎨 内置彩色日志输出
- 🔍 多级别日志(DEBUG/INFO/SUCCESS/WARN/ERROR)
- 📝 多种日志格式(JSON/详细/方括号/协程)
- ⏱️ 定时自动刷新缓冲区
- 🛡️ 线程安全设计
- 🔄 自动日志轮转功能
- 🔧 可配置的日志切割策略

## 安装与引入

```bash
# 确保在自己项目路径下，并且存在go.mog文件，不存在则 go init 项目名 创建
go get gitee.com/MM-Q/fastlog

# 引入 
import "gitee.com/MM-Q/fastlog"
```

## 快速开始

### 完整配置示例

```go
package main

import "gitee.com/MM-Q/fastlog"

func main() {
    // 完整结构体配置
    config1 := &fastlog.FastLogConfig{
        LogDirName:     "logs",              // 日志目录名
        LogFileName:    "app.log",           // 日志文件名
        PrintToConsole: true,                // 是否将日志输出到控制台
        ConsoleOnly:    false,               // 是否仅输出到控制台
        FlushInterval:  1 * time.Second,     // 刷新间隔
        LogLevel:       fastlog.DEBUG,       // 日志级别
        ChanIntSize:    1000,                // 通道大小
        LogFormat:      fastlog.Detailed,    // 日志格式选项
        MaxBufferSize:  1 * 1024 * 1024,     // 最大缓冲区大小(MB)
        NoColor:        false,               // 是否禁用终端颜色
        NoBold:        false,                // 是否禁用终端字体加粗
        MaxLogFileSize: 1,                   // 单个日志文件最大大小(MB)
        MaxLogAge:      30,                  // 日志文件保留天数
        MaxLogBackups:  10,                  // 日志文件保留数量
        IsLocalTime:    true,                // 是否使用本地时间
        EnableCompress: false,               // 是否启用压缩
    }

    // 创建日志实例
    logger, err := fastlog.NewFastLog(config1)
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

### 简化配置示例

```go
package main

import "gitee.com/MM-Q/fastlog"

func main() {
    // 仅指定日志目录和文件名的简化配置
    config2 := fastlog.NewFastLogConfig("custom_logs", "custom.log")

    // 创建日志实例
    logger, err := fastlog.NewFastLog(config2)
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

### 默认配置示例

```go
package main

import "gitee.com/MM-Q/fastlog"

func main() {
    // 完全使用默认配置
    config3 := fastlog.NewFastLogConfig("", "")

    // 创建日志实例
    logger, err := fastlog.NewFastLog(config3)
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

## 日志轮转功能

FastLog 提供了自动日志轮转功能，当满足以下条件时会自动创建新的日志文件：

1. 当前日志文件大小超过`MaxLogFileSize`设置的值

轮转后的日志文件会以原文件名加上时间的形式保存，例如：`app-2025-05-01T00-16-27.372.log`

## NoColor功能

FastLog支持通过设置`NoColor`属性为`true`来全局禁用颜色输出。当`NoColor`为`true`时，所有日志输出将直接显示原始文本，不添加任何颜色代码。

### 使用方法

```go
// 创建日志配置
cfg := fastlog.NewFastLogConfig("logs", "nocolor.log")
cfg.NoColor = true // 禁用终端颜色

// 创建日志实例
logger, err := fastlog.NewFastLog(cfg)
if err != nil {
    panic(err)
}
defer logger.Close()

// 此时所有日志输出将不会有颜色
logger.Info("这是一条无颜色信息日志")
logger.Error("这是一条无颜色错误日志")
```

### 使用场景
- 当终端不支持ANSI颜色代码时
- 需要将输出重定向到文件时
- 其他需要禁用颜色的场景

## 日志级别

| 级别    | 值  | 说明           |
| ------- | --- | -------------- |
| DEBUG   | 10  | 调试信息       |
| INFO    | 20  | 普通信息       |
| SUCCESS | 30  | 成功信息       |
| WARN    | 40  | 警告信息       |
| ERROR   | 50  | 错误信息       |
| None    | 999 | 不记录任何日志 |

## 日志格式

FastLog 支持以下几种日志格式：

| 格式名称 | 说明                                              |
| -------- | ------------------------------------------------- |
| Json     | 以 JSON 格式输出日志                              |
| Bracket  | 以方括号格式输出日志                              |
| Detailed | 详细格式，包含时间、级别、文件、函数、行号等信息  |
| Threaded | 包含线程 ID 的详细格式                            |
| Simple   | 简单格式，仅包含时间、级别和消息                  |
| Custom   | 自定义格式，通过类似于 fmt.Printf()格式进行自定义 |

1. **Detailed** (默认)

```ini
2023-01-01 12:00:00 | INFO    | main.go:main:10 - 日志信息
```

2. **Json**

```json
{
  "time": "2023-01-01 12:00:00",
  "level": "INFO",
  "file": "main.go",
  "function": "main",
  "line": "10",
  "thread": "1",
  "message": "日志信息"
}
```

3. **Bracket**

```ini
[INFO] 日志信息
```

4. **Threaded**

```ini
2023-01-01 12:00:00 | INFO    | [thread="1"] 日志信息
```

5. **Simple**

```ini
2023-01-01 12:00:00 | INFO    | 日志信息
```

6. **Custom** (自定义格式)

```ini
当配置为Custom时，请使用自定义格式，通过类似于fmt.Printf()格式进行自定义，然后传递给日志方法。
```

## 性能优化

- 异步处理：所有日志操作通过 channel 异步处理
- 缓冲区：使用内存缓冲区减少 IO 操作
- 批量写入：定时刷新缓冲区
- 内存优化：减少内存分配次数
- 并发控制：优化锁粒度提升并发性能

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

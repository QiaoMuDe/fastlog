# FastLog - 高性能 Go 日志库

![License](https://img.shields.io/badge/license-GPL-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.24.4-blue.svg)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/QiaoMuDe/fastlog)

FastLog 是一个高性能、异步、灵活格式的 Go 语言日志库，支持多种日志格式和级别，提供文件和控制台双输出。采用单线程处理器架构，通过批量处理和智能背压控制实现极致性能。

## 🌟 核心特性

- 🚀 **异步非阻塞日志记录** - 基于 channel 的异步处理架构
- 📁 **双输出支持** - 同时支持文件和控制台输出
- 🎨 **彩色日志输出** - 内置终端颜色支持，可配置禁用
- 🔍 **多级别日志** - 支持 DEBUG/INFO/SUCCESS/WARN/ERROR/FATAL 六个级别
- 📝 **多种日志格式** - Detailed/Json/Simple/Structured/Custom 五种格式
- ⏱️ **智能缓冲管理** - 定时刷新+阈值触发的双重缓冲策略
- 🛡️ **线程安全设计** - 单线程处理器确保数据一致性
- 🔄 **自动日志轮转** - 基于文件大小、时间和数量的轮转策略
- 🔧 **智能背压控制** - 根据通道使用率自动丢弃低优先级日志
- ⚡ **批量处理优化** - 批量格式化和写入，减少 IO 开销
- 🎯 **内存优化** - 预分配缓冲区，减少 GC 压力
- 🔒 **配置验证** - 自动修正不合理配置，确保系统稳定

## 📦 安装与引入

```bash
# 确保在项目路径下存在 go.mod 文件
go get gitee.com/MM-Q/fastlog
```

```go
import "gitee.com/MM-Q/fastlog"
```

## 🤠 获取性能测试报告
```bash
go test -v -run TestConcurrentFastLog
```

## 🚀 快速开始

### 基础使用示例

```go
package main

import "gitee.com/MM-Q/fastlog"

func main() {
    // 创建日志配置
    config := fastlog.NewFastLogConfig("logs", "app.log")

    // 创建日志实例
    logger, err := fastlog.NewFastLog(config)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    // 记录不同级别的日志
    logger.Debug("这是调试信息")
    logger.Info("这是一条信息日志")
    logger.Success("操作成功完成")
    logger.Warn("这是一个警告")
    logger.Error("发生了一个错误")
  
    // 使用格式化方法
    logger.Infof("用户 %s 登录成功，IP: %s", "张三", "192.168.1.1")
    logger.Errorf("数据库连接失败，重试次数: %d", 3)
}
```

### 简化创建方式

```go
package main

import "gitee.com/MM-Q/fastlog"

func main() {
    // 使用简写函数创建
    config := fastlog.NewCfg("logs", "app.log")  // NewCfg 是 NewFastLogConfig 的简写
    logger, err := fastlog.New(config)           // New 是 NewFastLog 的简写
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    logger.Info("使用简写函数创建的日志")
}
```

## ⚙️ 高级配置

### 完整配置示例

```go
package main

import (
    "time"
    "gitee.com/MM-Q/fastlog"
)

func main() {
    // 创建配置
    config := fastlog.NewFastLogConfig("logs", "app.log")
  
    // 配置日志级别和格式
    config.LogLevel = fastlog.DEBUG                      // 设置日志级别
    config.LogFormat = fastlog.Json                      // 设置 JSON 格式
    
    // 配置缓冲和刷新
    config.FlushInterval = 500 * time.Millisecond        // 设置刷新间隔
    config.ChanIntSize = 20000                           // 设置通道大小
    
    // 配置文件轮转
    config.MaxLogFileSize = 50                           // 最大文件大小 50MB
    config.MaxLogAge = 30                                // 保留 30 天
    config.MaxLogBackups = 10                            // 最多 10 个备份
    config.IsLocalTime = true                            // 使用本地时间
    config.EnableCompress = true                         // 启用压缩
    
    // 配置输出选项
    config.OutputToConsole = true                        // 启用控制台输出
    config.OutputToFile = true                           // 启用文件输出
    config.NoColor = false                               // 启用颜色
    config.NoBold = false                                // 启用加粗
  
    // 创建日志实例
    logger, err := fastlog.NewFastLog(config)
    if err != nil {
        panic(err)
    }
    defer logger.Close()
  
    // 记录各种级别的日志
    logger.Debug("调试信息：程序启动")
    logger.Info("应用程序已启动")
    logger.Success("数据库连接成功")
    logger.Warn("配置文件使用默认值")
    logger.Error("网络连接超时")
}
```

### 环境配置示例

```go
// 推荐：根据环境使用不同的配置
func createLogger(env string) (*fastlog.FastLog, error) {
    config := fastlog.NewFastLogConfig("logs", "app.log")
  
    switch env {
    case "development":
        config.LogLevel = fastlog.DEBUG
        config.LogFormat = fastlog.Detailed
        config.OutputToConsole = true
    case "production":
        config.LogLevel = fastlog.INFO
        config.LogFormat = fastlog.Json
        config.OutputToConsole = false
        config.MaxLogFileSize = 100
        config.MaxLogAge = 7
    }
  
    return fastlog.NewFastLog(config)
}
```

## 📝 日志格式

FastLog 支持五种不同的日志格式：

| 格式名称 | 枚举值 | 说明 |
|---------|--------|------|
| Detailed | `fastlog.Detailed` | 详细格式，包含时间、级别、文件、函数、行号等完整信息（默认） |
| Json | `fastlog.Json` | JSON 格式输出，便于日志分析和处理 |
| Simple | `fastlog.Simple` | 简约格式，仅包含时间、级别和消息 |
| Structured | `fastlog.Structured` | 结构化格式，使用分隔符组织信息 |
| Custom | `fastlog.Custom` | 自定义格式，直接输出原始消息 |

### 格式示例

#### 1. Detailed 格式（默认）
```
2025-01-15 10:30:45 | INFO    | main.go:main:15 - 用户登录成功
2025-01-15 10:30:46 | ERROR   | database.go:Connect:23 - 数据库连接失败
```

#### 2. JSON 格式
```json
{"time":"2025-01-15 10:30:45","level":"INFO","file":"main.go","function":"main","line":15,"message":"用户登录成功"}
{"time":"2025-01-15 10:30:46","level":"ERROR","file":"database.go","function":"Connect","line":23,"message":"数据库连接失败"}
```

#### 3. Simple 格式
```
2025-01-15 10:30:45 | INFO    | 用户登录成功
2025-01-15 10:30:46 | ERROR   | 数据库连接失败
```

#### 4. Structured 格式
```
T:2025-01-15 10:30:45|L:INFO   |F:main.go:main:15|M:用户登录成功
T:2025-01-15 10:30:46|L:ERROR  |F:database.go:Connect:23|M:数据库连接失败
```

#### 5. Custom 格式
```go
// 使用 Custom 格式时，直接输出传入的消息内容
logger.Info("自定义格式的日志消息")  // 输出: 自定义格式的日志消息
```

## 🎯 核心特性详解

### 智能背压控制

FastLog 实现了智能背压控制机制，根据日志通道的使用率自动调整日志处理策略：

- **98%+ 使用率**: 只保留 FATAL 级别日志
- **95%+ 使用率**: 只保留 ERROR 及以上级别日志
- **90%+ 使用率**: 只保留 WARN 及以上级别日志
- **80%+ 使用率**: 只保留 SUCCESS 及以上级别日志
- **70%+ 使用率**: 丢弃 DEBUG 级别日志
- **正常情况**: 处理所有级别日志

这确保了在高负载情况下系统的稳定性，避免日志积压导致的内存问题。

### 批量处理优化

- **批量大小**: 默认 1000 条日志为一批
- **缓冲区管理**: 文件缓冲区 32KB 初始容量，控制台缓冲区 8KB 初始容量
- **智能刷新**: 90% 阈值触发或定时刷新（默认 500ms）
- **内存优化**: 预分配缓冲区，减少 GC 压力

### 日志轮转功能

FastLog 基于 `logrotatex` 提供强大的日志轮转功能：

```go
config := fastlog.NewFastLogConfig("logs", "app.log")
config.MaxLogFileSize = 10      // 文件大小超过 10MB 时轮转
config.MaxLogAge = 7            // 保留 7 天的日志文件
config.MaxLogBackups = 5        // 最多保留 5 个备份文件
config.IsLocalTime = true       // 使用本地时间命名
config.EnableCompress = true    // 启用压缩功能
```

轮转后的日志文件命名格式：`app-2025-01-15T10-30-45.123.log`

## 📊 API 参考

### 核心函数

| 函数名称 | 参数 | 返回值 | 说明 |
|---------|------|--------|------|
| `NewFastLogConfig` | `logDirPath string, logFileName string` | `*FastLogConfig` | 创建日志配置实例 |
| `NewFastLog` | `cfg *FastLogConfig` | `(*FastLog, error)` | 根据配置创建日志记录器实例 |
| `New` | `cfg *FastLogConfig` | `(*FastLog, error)` | `NewFastLog` 的简写形式 |
| `NewCfg` | `logDirPath string, logFileName string` | `*FastLogConfig` | `NewFastLogConfig` 的简写形式 |

### 配置字段

| 字段名称 | 类型 | 默认值 | 说明 |
|---------|------|--------|------|
| `LogDirName` | `string` | - | 日志目录路径 |
| `LogFileName` | `string` | - | 日志文件名 |
| `OutputToConsole` | `bool` | `true` | 是否输出到控制台 |
| `OutputToFile` | `bool` | `true` | 是否输出到文件 |
| `FlushInterval` | `time.Duration` | `500ms` | 缓冲区刷新间隔 |
| `LogLevel` | `LogLevel` | `INFO` | 日志级别 |
| `ChanIntSize` | `int` | `10000` | 通道大小 |
| `LogFormat` | `LogFormatType` | `Detailed` | 日志格式 |
| `NoColor` | `bool` | `false` | 是否禁用颜色输出 |
| `NoBold` | `bool` | `false` | 是否禁用字体加粗 |
| `MaxLogFileSize` | `int` | `10` | 最大日志文件大小(MB) |
| `MaxLogAge` | `int` | `0` | 日志文件最大保留天数(0表示不限制) |
| `MaxLogBackups` | `int` | `0` | 最大备份文件数量(0表示不限制) |
| `IsLocalTime` | `bool` | `true` | 是否使用本地时间 |
| `EnableCompress` | `bool` | `false` | 是否启用日志压缩 |

### 日志记录方法

#### 基础日志方法（不支持格式化）

| 方法名称 | 参数 | 日志级别 | 说明 |
|---------|------|----------|------|
| `Debug` | `v ...any` | DEBUG | 记录调试信息 |
| `Info` | `v ...any` | INFO | 记录一般信息 |
| `Success` | `v ...any` | SUCCESS | 记录成功信息 |
| `Warn` | `v ...any` | WARN | 记录警告信息 |
| `Error` | `v ...any` | ERROR | 记录错误信息 |
| `Fatal` | `v ...any` | FATAL | 记录致命错误，程序将退出 |

#### 格式化日志方法（支持占位符）

| 方法名称 | 参数 | 日志级别 | 说明 |
|---------|------|----------|------|
| `Debugf` | `format string, v ...any` | DEBUG | 格式化记录调试信息 |
| `Infof` | `format string, v ...any` | INFO | 格式化记录一般信息 |
| `Successf` | `format string, v ...any` | SUCCESS | 格式化记录成功信息 |
| `Warnf` | `format string, v ...any` | WARN | 格式化记录警告信息 |
| `Errorf` | `format string, v ...any` | ERROR | 格式化记录错误信息 |
| `Fatalf` | `format string, v ...any` | FATAL | 格式化记录致命错误，程序将退出 |

#### 控制方法

| 方法名称 | 参数 | 说明 |
|---------|------|------|
| `Close` | 无 | 优雅关闭日志记录器，确保所有日志写入完成 |

### 日志级别常量

| 常量名称 | 数值 | 说明 |
|---------|------|------|
| `DEBUG` | 10 | 调试级别，用于开发调试 |
| `INFO` | 20 | 信息级别，记录一般信息 |
| `SUCCESS` | 30 | 成功级别，记录成功操作 |
| `WARN` | 40 | 警告级别，记录警告信息 |
| `ERROR` | 50 | 错误级别，记录错误信息 |
| `FATAL` | 60 | 致命级别，记录致命错误 |
| `NONE` | 255 | 禁用所有日志输出 |

## 💡 最佳实践

### 1. 日志级别使用建议

- **DEBUG**: 仅在开发和调试时使用，包含详细的程序执行信息
- **INFO**: 记录程序的正常运行信息，如启动、关闭、重要操作等
- **SUCCESS**: 记录成功完成的重要操作，如用户登录、订单完成等
- **WARN**: 记录警告信息，程序可以继续运行但需要注意
- **ERROR**: 记录错误信息，程序遇到问题但可以恢复
- **FATAL**: 记录致命错误，程序无法继续运行

### 2. 性能优化建议

```go
// 推荐：在高并发场景下适当调整配置
config := fastlog.NewFastLogConfig("logs", "app.log")
config.FlushInterval = 1 * time.Second    // 增加刷新间隔减少 IO
config.LogLevel = fastlog.INFO            // 生产环境避免 DEBUG 日志
config.MaxLogFileSize = 100               // 适当增大文件大小减少轮转频率

// 推荐：使用格式化方法避免不必要的字符串拼接
logger.Infof("用户 %s 执行操作 %s", username, action)  // 推荐
// logger.Info("用户 " + username + " 执行操作 " + action)  // 不推荐
```

### 3. 错误处理建议

```go
// 推荐：在关键位置使用 defer 确保日志记录器正确关闭
func main() {
    logger, err := fastlog.NewFastLog(config)
    if err != nil {
        panic(err)
    }
    defer func() {
        logger.Info("程序正在关闭...")
        logger.Close()
    }()
  
    // 业务逻辑...
}
```

## 🔧 依赖信息

FastLog 依赖以下外部库：

```go
require (
    gitee.com/MM-Q/colorlib v1.2.4      // 终端颜色输出库
    gitee.com/MM-Q/logrotatex v1.0.0    // 日志轮转库
)
```

- **colorlib**: 提供终端颜色输出功能，支持多种颜色和样式
- **logrotatex**: 提供日志文件轮转功能，支持按大小、时间和数量轮转

## 📋 版本要求

- **Go 版本**: 1.24.4 或更高版本
- **操作系统**: 支持 Windows、Linux、macOS

## ❓ 常见问题

### Q: 如何在生产环境中使用？

A: 建议在生产环境中：
- 设置日志级别为 INFO 或更高
- 使用 JSON 格式便于日志分析
- 禁用控制台输出，仅输出到文件
- 配置合适的日志轮转策略

### Q: 如何处理高并发场景？

A: FastLog 内置了智能背压控制：
- 自动根据负载丢弃低优先级日志
- 批量处理减少系统调用
- 异步处理避免阻塞业务逻辑

### Q: 日志文件过大怎么办？

A: 配置日志轮转参数：
```go
config.MaxLogFileSize = 50      // 50MB 轮转
config.MaxLogAge = 7            // 保留 7 天
config.MaxLogBackups = 10       // 最多 10 个备份
config.EnableCompress = true    // 启用压缩
```

### Q: 如何自定义日志格式？

A: 使用 Custom 格式：
```go
config.LogFormat = fastlog.Custom
logger.Info("自定义格式的消息")  // 直接输出原始消息
```

## 📈 性能测试

### 高并发性能测试报告

基于最新优化后的测试数据（22核CPU，10个Goroutine并发）：

```
============================================================
           FastLog 高并发性能测试报告
============================================================
📊 测试基本信息:
   开始时间: 2025-08-02 22:09:16.670
   结束时间: 2025-08-02 22:09:19.207
   测试耗时: 2.537s (2537.37ms)
   Goroutine数量: 10

📝 日志处理统计:
   预期生成: 300.0万条日志
   实际生成: 300.0万条日志
   文件写入: 4.1万条有效日志
   成功率: 1.37%
   实际吞吐量: 118.2万条/秒

💾 内存使用统计:
   开始内存: 337.7 KB
   结束内存: 1.2 MB
   峰值内存: 33.1 MB
   内存增长: +912.1 KB
   总分配: 874.0 MB
   系统内存: 57.0 MB
   GC次数: 71次
   GC暂停时间: 38.84ms

⚡ 性能评估:
   平均每条日志内存开销: 0.31 bytes
   平均每条日志处理时间: 0.85 μs

🖥️  系统资源:
   CPU核心数: 22
   最大并发Goroutine: 10
============================================================
```

### 性能指标总结

| 指标类型 | 数值 | 说明 |
|---------|------|------|
| **吞吐量** | 118.2万条/秒 | 10个Goroutine并发下的实际处理能力 |
| **处理延迟** | 0.85μs | 平均每条日志的处理时间 |
| **内存效率** | 0.31 bytes/条 | 平均每条日志的内存开销 |
| **智能背压** | 1.37% | 高负载下的日志保留率（智能丢弃98.63%） |
| **GC友好** | 71次GC | 2.5秒内触发71次垃圾回收 |
| **内存峰值** | 33.1MB | 处理300万条日志的最大内存占用 |

### 测试环境说明

- **测试方法**: `go test -v -run TestConcurrentFastLog`
- **并发级别**: 10个Goroutine
- **测试数据量**: 300万条日志
- **智能背压**: 启用（高负载下自动丢弃低优先级日志）
- **文件输出**: 启用（实际写入4.1万条有效日志）

### 性能特点

1. **超高吞吐量**: 在10个Goroutine并发下达到118.2万条/秒的处理能力
2. **超低延迟**: 平均0.85微秒的单条日志处理时间
3. **极致内存效率**: 每条日志仅占用0.31字节内存开销
4. **智能背压**: 高负载下自动保留重要日志，丢弃冗余信息
5. **优化的GC管理**: 通过零拷贝和对象池优化减少垃圾回收压力

### 🚀 性能优化成果

经过全面的性能优化，FastLog 在以下方面取得了显著提升：

- **吞吐量提升**: 从68.8万条/秒提升到118.2万条/秒（+71.8%）
- **延迟降低**: 从1.45μs降低到0.85μs（-41.4%）
- **内存效率**: 从7.61 bytes/条优化到0.31 bytes/条（-95.9%）
- **处理能力**: 单次测试处理300万条日志（vs 之前90万条）

#### 核心优化技术

1. **时间戳缓存机制**: 减少重复的时间格式化开销
2. **文件名缓存优化**: 避免重复的 `filepath.Base()` 调用
3. **临时缓冲区对象池**: 复用缓冲区，减少内存分配
4. **零拷贝写入优化**: 直接格式化到目标缓冲区
5. **零拷贝颜色处理**: 直接缓冲区颜色操作
6. **日志级别字符串预分配**: 消除动态填充开销
7. **静态检查优化**: 使用 `fmt.Fprintf` 替代低效的字符串操作

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 开发规范

- 遵循 Go 代码规范
- 添加必要的单元测试
- 更新相关文档
- 确保所有测试通过

## 📄 许可证

本项目采用 GPL 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系方式

- 项目主页：[https://deepwiki.com/QiaoMuDe/fastlog](https://deepwiki.com/QiaoMuDe/fastlog)
- Gitee：[https://gitee.com/MM-Q/fastlog](https://gitee.com/MM-Q/fastlog)
- GitHub：[https://github.com/QiaoMuDe/fastlog](https://github.com/QiaoMuDe/fastlog)

## 📚 相关文档

- [API 文档](./APIDOC.md) - 详细的 API 参考文档

---

**FastLog** - 让日志记录更简单、更高效！ 🚀
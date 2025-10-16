<div align="center">

# FastLog - 高性能 Go 日志库

![License](https://img.shields.io/badge/license-GPL-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.25.0-blue.svg)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/QiaoMuDe/fastlog)

</div>

FastLog 是一个高性能的 Go 日志库，面向生产可用与易用性设计。当前实现为同步写入，提供严格的配置校验、丰富的日志级别与格式、灵活的输出目标，以及完善的日志轮转支持。配套提供开发/生产/终端三种便捷模式构造函数，便于快速在不同场景落地。

## 🌟 核心特性

### ⚙️ 写入与配置
- 同步写入：简洁直接，便于理解与调试
- 异步清理：支持 Async 控制是否异步清理轮转后的旧日志（默认同步）
- 严格配置校验：validateConfig 对不合理配置直接 panic，防止隐性错误
- 刷新与大小控制：支持 FlushInterval、MaxSize/MaxAge/MaxFiles 等参数
- 本地时间与压缩：支持本地时间命名与压缩轮转

### 📊 日志能力
- 日志级别：DEBUG / INFO / WARN / ERROR / FATAL
- 输出格式：Def / Json / Structured / Timestamp / Custom
- 输出目标：文件与控制台可独立开启/关闭
- 颜色与样式：可配置 Color/Bold 提升终端可读性

### 🔧 开发友好设计
- 便捷模式构造函数：
  - DevConfig：开发模式（详细格式、DEBUG、短期保留）
  - ProdConfig：生产模式（压缩、禁用控制台、长期保留）
  - ConsoleConfig：终端模式（仅控制台、DEBUG、简洁时间戳）
- 完整测试：包含高并发性能测试与 Fatal/Fatalf 子进程行为验证
- 简洁 API：通过 NewFastLogConfig + NewStdLog 组合使用，支持格式化方法（Infof/Debugf 等）

## 📦 安装与引入

```bash
# 确保在项目路径下存在 go.mod 文件
go get gitee.com/MM-Q/fastlog
```

```go
import "gitee.com/MM-Q/fastlog"
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
    logger := fastlog.NewStdLog(config)
    defer logger.Close()

    // 记录不同级别的日志
    logger.Debug("这是调试信息")
    logger.Info("这是一条信息日志")
    logger.Warn("这是一个警告")
    logger.Error("发生了一个错误")
  
    // 使用格式化方法
    logger.Infof("用户 %s 登录成功，IP: %s", "张三", "192.168.1.1")
    logger.Errorf("数据库连接失败，重试次数: %d", 3)
}
```

### 预设配置模式（推荐）

FastLog 提供了 3 种便捷模式构造函数，覆盖常见使用场景：

```go
package main

import "gitee.com/MM-Q/fastlog"

func main() {
    // 开发模式：文件+控制台、详细格式、彩色加粗、快速刷新
    devCfg := fastlog.DevConfig("logs", "dev.log")
    logger := fastlog.NewStdLog(devCfg)
    defer logger.Close()

    // 生产模式：仅文件、结构化格式、压缩、长期保留
    // prodCfg := fastlog.ProdConfig("logs", "prod.log")
    
    // 终端模式：仅控制台、时间戳简洁格式、彩色加粗
    // consoleCfg := fastlog.ConsoleConfig()

    // 记录日志
    logger.Info("应用程序启动")
    logger.Debugf("用户ID: %d, 操作: %s", 12345, "登录")
    logger.Warn("内存使用率较高: 85%")
    logger.Error("数据库连接失败")
}
```

#### 可用的预设配置模式

- `DevConfig(logDir, logFile)` - 开发模式：文件+控制台、Detailed、DEBUG、彩色加粗、FlushInterval≈200ms、短期保留（示例：MaxFiles=5 / MaxAge=7）
- `ProdConfig(logDir, logFile)` - 生产模式：仅文件、Structured、INFO、无装饰、压缩、MaxSize=100MB、FlushInterval≈1s、长期保留（30天 / 24个）
- `ConsoleConfig()` - 终端模式：仅控制台、Timestamp、DEBUG、彩色加粗、FlushInterval≈500ms（不写文件）

### 简化创建方式

```go
package main

import "gitee.com/MM-Q/fastlog"

func main() {
    // 使用简写函数创建
    config := fastlog.NewCfg("logs", "app.log")  // NewCfg 是 NewFastLogConfig 的简写
    logger := fastlog.New(config)           // New 是 NewStdLog 的简写
    defer logger.Close()

    logger.Info("使用简写函数创建的日志")
}
```

## 📝 日志格式

FastLog 支持4种不同的日志格式：

| 格式名称 | 枚举值 | 说明 |
|---------|--------|------|
| Def | `fastlog.Def` | 默认格式，包含时间、级别、文件、函数、行号等完整信息 |
| Json | `fastlog.Json` | JSON 格式输出，便于日志分析和处理 |
| Timestamp | `fastlog.Timestamp` | 时间格式 |
| KVfmt | `fastlog.KVfmt` | 键值对格式 |
| LogFmt | `fastlog.LogFmt` | 日志格式 |
| Custom | `fastlog.Custom` | 自定义格式，直接输出原始消息 |

### 格式示例

#### 1. Def 格式
```
2025-01-15T10:30:45 | INFO    | main.go:main:15 - 用户登录成功
2025-01-15T10:30:46 | ERROR   | database.go:Connect:23 - 数据库连接失败
```

#### 2. JSON 格式
```json
{"time":"2025-01-15T10:30:45","level":"INFO","caller":"main.go:main:15","message":"用户登录成功"}
{"time":"2025-01-15T10:30:46","level":"ERROR","caller":"database.go:Connect:23","message":"数据库连接失败"}
```

#### 3. Timestamp 格式
```
2025-01-15T10:30:45 INFO  用户登录成功
2025-01-15T10:30:46 ERROR 数据库连接失败
```

#### 4. KVfmt 键值对格式
```
time=2025-01-15T10:30:45 level=INFO message=用户登录成功
time=2025-01-15T10:30:46 level=ERROR message=数据库连接失败
```

#### 5. LogFmt 日志格式
```
2025-01-15T10:30:45 [INFO ] 用户登录成功 [username=张三, age=30]
2025-01-15T10:30:46 [ERROR] database.go:Connect:23 数据库连接失败
```

#### 6. Custom 格式
```go
// 使用 Custom 格式时，直接输出传入的消息内容
logger.Info("自定义格式的日志消息")  // 输出: 自定义格式的日志消息
```

## 🎯 核心特性详解

### 同步写入与关闭
- 当前实现为同步写入，调用日志方法会直接格式化并写出到目标（文件/控制台）。
- 优雅关闭：Close 返回错误用于传递关闭阶段的异常；在测试中对 Close 的错误均进行显式处理，保证资源释放。

### 配置校验与安全
- validateConfig 对关键字段进行严格校验：输出目标至少启用一个、级别/格式范围合法、大小/天数/数量不越界。
- 路径安全：对日志目录与文件名进行路径穿越检测（禁止包含“..”）。
- 可选项：支持本地时间命名与压缩轮转，便于生产环境长期保留。

### 日志轮转功能
FastLog 基于轮转参数支持常见的日志管理策略：

```go
config := fastlog.NewFastLogConfig("logs", "app.log")
config.MaxSize = 10             // 文件大小超过 10MB 时轮转
config.MaxAge = 7               // 保留 7 天的日志文件
config.MaxFiles = 5             // 最多保留 5 个备份文件
config.LocalTime = true         // 使用本地时间命名
config.Compress = true          // 启用压缩功能
```

轮转后的日志文件命名格式示例：`app_202501010301.log`

## 💡 最佳实践

### 1. 日志级别使用建议

- **DEBUG**: 仅在开发和调试时使用，包含详细的程序执行信息
- **INFO**: 记录程序的正常运行信息，如启动、关闭、重要操作等
- **WARN**: 记录警告信息，程序可以继续运行但需要注意
- **ERROR**: 记录错误信息，程序遇到问题但可以恢复
- **FATAL**: 记录致命错误，程序无法继续运行

### 2. 配置与性能建议

```go
// 推荐：根据场景选择合适的预设配置
// 生产环境（仅文件、结构化、长期保留、压缩）
logger := fastlog.NewStdLog(fastlog.ProdConfig("logs", "app.log"))

// 开发环境（文件+控制台、详细信息、彩色加粗、快速刷新）
logger := fastlog.NewStdLog(fastlog.DevConfig("logs", "debug.log"))

// 终端环境（仅控制台、时间戳简洁格式、彩色加粗）
logger := fastlog.NewStdLog(fastlog.ConsoleConfig())

// 自定义配置（在预设不满足需求时使用）
config := fastlog.NewFastLogConfig("logs", "app.log")
config.FlushInterval = 1 * time.Second    // 增加刷新间隔减少 IO（根据需求调整）
config.LogLevel = fastlog.INFO            // 生产环境避免 DEBUG 噪音
config.MaxSize = 100                      // 增大文件大小减少轮转频率

// 使用格式化方法，避免不必要的字符串拼接
logger.Infof("用户 %s 执行操作 %s", username, action)  // 推荐
// logger.Info("用户 " + username + " 执行操作 " + action)  // 不推荐
```

### 3. 错误处理建议

```go
// 推荐：在关键位置使用 defer 并显式处理 Close 的错误（满足 errcheck）
func main() {
    logger := fastlog.NewStdLog(config)
    defer func() {
        logger.Info("程序正在关闭...")
        if err := logger.Close(); err != nil {
            // 生产环境可记录到标准错误或告警系统
            // 此处简单打印或忽略均可根据需要处理
            _ = err
        }
    }()
  
    // 业务逻辑...
}
```

## 🔧 依赖信息

FastLog 依赖以下外部库（与 go.mod 保持一致）：

```go
require (
    gitee.com/MM-Q/colorlib v1.3.2      // 终端颜色与样式
    gitee.com/MM-Q/logrotatex v1.1.2    // 日志文件轮转（大小/时间/数量、压缩、本地时间）
    gitee.com/MM-Q/go-kit v0.0.8        // 常用工具集
)
```

- colorlib：提供终端颜色输出与样式支持
- logrotatex：提供日志文件轮转功能，支持按大小、时间和数量轮转，以及压缩与本地时间
- go-kit：提供通用工具函数（如字符串/配置等）

## 📋 版本要求

- Go 版本：1.25.0 或更高版本
- 操作系统：支持 Windows、Linux、macOS

## ❓ 常见问题

### Q: 如何在生产环境中使用？

A: 建议在生产环境中：
- 设置日志级别为 INFO 或更高（避免 DEBUG 噪音）
- 使用 Structured 或 JSON 格式便于日志分析
- 禁用控制台输出，仅输出到文件（ProdConfig 默认如此）
- 配置合适的日志轮转策略（MaxSize/MaxAge/MaxFiles，建议启用压缩）

### Q: 如何处理高并发场景？

A: 当前实现为同步写入，建议：
- 使用 `ProdConfig("logs", "app.log")`，开启压缩与结构化格式
- 将日志级别设为 INFO 或更高，减少无关日志
- 适度增大 `FlushInterval`（如 1s）与 `MaxSize`（如 100MB），降低 I/O 频率
- 采用结构化/JSON 格式，便于后续日志聚合与分析

### Q: 日志文件过大怎么办？

A: 配置日志轮转参数：
```go
config.MaxSize = 50      // 50MB 轮转（可根据场景调整为 100MB）
config.MaxAge = 7        // 保留 7 天
config.MaxFiles = 10     // 最多 10 个备份
config.Compress = true   // 启用压缩
```

### Q: 如何自定义日志格式？

A: 使用 Custom 格式：
```go
config.LogFormat = fastlog.Custom
logger.Info("自定义格式的消息")  // 直接输出原始消息
```

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

<div align="center">

**FastLog** - 让日志记录更简单、更高效！ 🚀

</div>
<div align="center">

# FastLog

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev) [![License](https://img.shields.io/badge/License-MIT-green)](#许可证) [![Stars](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgitee.com%2Fapi%2Fv5%2Frepos%2FMM-Q%2Ffastlog&query=%24.stargazers_count&label=Stars&color=red&suffix=%20%E2%AD%90)](https://gitee.com/MM-Q/fastlog) [![Go Reference](https://img.shields.io/badge/Reference-pkg.go.dev-00ADD8)](https://pkg.go.dev/gitee.com/MM-Q/fastlog)

</div>

**FastLog** 是一个轻量级、高性能的 Go 日志库，专注于简洁 API 与生产可用性的平衡。
受 [zap](https://github.com/uber-go/zap) 启发，采用零分配字段设计、对象池复用、无锁采样等优化手段，同时提供友好的三级 API 和场景化的配置模式。

> 项目地址：[gitee.com/MM-Q/fastlog](https://gitee.com/MM-Q/fastlog)

---

## 核心特性

| 特性 | 说明 |
|------|------|
| 🚀 **高性能** | 零分配 Field 结构体、对象池复用 Entry、无锁原子采样、避免反射 |
| 🎨 **彩色输出** | 基于 [color](https://gitee.com/MM-Q/color) 库，自动检测日志级别着色，支持 `NoColor` 开关，兼容 Windows |
| 📋 **三级 API** | 标准日志 `Info()`、格式化日志 `Infof()`、结构化日志 `Infow()` |
| 🔧 **Config 配置** | 场景化配置函数，开箱即用，支持自定义调整 |
| ⏰ **时间格式可配置** | 通过 `TimeFormat` 自定义时间格式，默认 RFC3339，`DefaultTimeFormat` 常量统一管理 |
| 📝 **多格式支持** | 内置 5 种格式：Def、JSON、Simple、KV、Compact，支持自定义 |
| 🧩 **结构化字段** | 12 种字段类型，类型安全，零装箱分配 |
| 🎯 **日志采样** | 固定桶 + atomic 无锁设计，参考 zap，有效防洪 |
| 🔌 **多路输出** | `MultiWriter` 同时输出到多个目标 |
| 🧪 **场景化配置** | `NewConfig()`、`Dev()`、`Prod()`、`Console()`、`Docker()` 覆盖常见场景 |
| 🔒 **线程安全** | `sync.Mutex` 保证写入安全 |
| 📦 **一站式集成** | 基于 [logrotatex](https://gitee.com/MM-Q/logrotatex) 实现日志轮转、缓冲写入，[comprx](https://gitee.com/MM-Q/comprx) 实现压缩，用户无感知 |
| 🎚️ **动态级别** | 运行时通过 `SetLevel()` 调整日志级别，无需重启，基于 `atomic.Int32` 无锁实现 |
| 🗂️ **级别路由** | 通过 `LevelRouter` 启用，自动按级别分发到专属文件（如 ERROR.log），便于快速定位问题 |
| 💾 **缓冲控制** | 通过 `BufferEnabled` 控制是否启用缓冲写入，开发环境立即落盘，生产环境批量写入 |

---

## 安装

```bash
go get gitee.com/MM-Q/fastlog
```

然后在代码中引入：

```go
import "gitee.com/MM-Q/fastlog"
```

---

## 快速开始

### 最简单的用法

```go
package main

import "gitee.com/MM-Q/fastlog"

func main() {
    logger := fastlog.New(fastlog.Default())
    defer func() { _ = logger.Close() }()

    logger.Info("应用启动成功")
    logger.Warnf("磁盘使用率: %.1f%%", 85.3)
    logger.Errorw("数据库连接失败", fastlog.String("db", "mysql"))
}
```

### 开发调试模式

```go
logger := fastlog.New(fastlog.Dev("logs/dev.log"))
defer func() { _ = logger.Close() }()

logger.Debug("进入函数 processOrder")           // 彩色输出
logger.Info("订单处理完成")
logger.Error("数据库连接超时")                    // 带调用者信息
```

### 生产环境模式

```go
logger := fastlog.New(fastlog.Prod("/var/log/app.log"))
defer func() { _ = logger.Close() }()

logger.Info("服务启动")                          // 只写入文件，WARN 及以上级别
logger.Warn("连接数接近上限")
```

### 容器/Docker 模式

```go
logger := fastlog.New(fastlog.Docker())
defer func() { _ = logger.Close() }()

logger.Info("容器启动")                          // JSON 格式输出到 stdout
```

---

## 使用指南

### 场景化配置函数

FastLog 提供了 6 种场景化配置函数，覆盖最常见的使用场景：

| 函数 | 级别 | 输出 | 格式 | 调用者 | 采样 | 轮转 | 级别路由 | 缓冲 | 适用场景 |
|------|------|------|------|--------|------|------|----------|------|----------|
| `NewConfig(path)` | INFO | 终端+文件 | Def | ❌ | ✅ | ✅ | ❌ | ✅ | 通用默认 |
| `Default()` | INFO | 终端+文件 | Def | ❌ | ✅ | ✅ | ❌ | ✅ | 快速上手 |
| `Dev(path)` | DEBUG | 终端+文件 | Def | ✅ | ❌ | ❌ | ❌ | ❌ | 开发调试 |
| `Prod(path)` | WARN | 仅文件 | Def | ❌ | ✅ | ✅ | ✅ | ✅ | 生产环境 |
| `Console()` | DEBUG | 仅终端 | Def | ❌ | ❌ | ❌ | ❌ | N/A | 本地调试 |
| `Docker()` | WARN | 仅终端 | JSON | ❌ | ✅ | ❌ | ❌ | N/A | 容器/K8s |

```go
// 开发环境：DEBUG 级别、彩色输出、调用者信息
logger := fastlog.New(fastlog.Dev("logs/dev.log"))

// 生产环境：WARN 级别、文件输出、压缩、异步清理
logger := fastlog.New(fastlog.Prod("/var/log/app.log"))

// 容器环境：JSON 格式、stdout 输出
logger := fastlog.New(fastlog.Docker())

// 纯控制台：本地调试
logger := fastlog.New(fastlog.Console())
```

### 自定义配置

基于场景化配置进行修改：

```go
// 基于生产环境配置，调整为 ERROR 级别
cfg := fastlog.Prod("/var/log/app.log")
cfg.Level = fastlog.ERROR
cfg.Compress = false  // 关闭压缩

logger := fastlog.New(cfg)
defer func() { _ = logger.Close() }()
```

### 完整自定义配置

```go
cfg := &fastlog.Config{
    Level:         fastlog.INFO,
    Formatter:     fastlog.JSON{},
    Caller:        true,
    OutputConsole: true,
    OutputFile:    true,
    LogPath:       "logs/app.log",
    MaxSize:       100,
    MaxFiles:      7,
    Compress:      true,
}

logger := fastlog.New(cfg)
defer func() { _ = logger.Close() }()
```

### 结构化字段

```go
logger.Infow("用户登录",
    fastlog.String("username", "alice"),
    fastlog.Int("attempts", 3),
    fastlog.Bool("vip", true),
    fastlog.Duration("elapsed", 45*time.Millisecond),
    fastlog.Error(err),
    fastlog.Any("metadata", map[string]string{"ip": "192.168.1.1"}),
)
```

**三级 API 对应关系：**

| 级别 | 标准 | 格式化 | 结构化 |
|------|------|--------|--------|
| DEBUG | `Debug(msg)` | `Debugf(fmt, args...)` | `Debugw(msg, fields...)` |
| INFO | `Info(msg)` | `Infof(fmt, args...)` | `Infow(msg, fields...)` |
| WARN | `Warn(msg)` | `Warnf(fmt, args...)` | `Warnw(msg, fields...)` |
| ERROR | `Error(msg)` | `Errorf(fmt, args...)` | `Errorw(msg, fields...)` |
| FATAL | `Fatal(msg)` | `Fatalf(fmt, args...)` | `Fatalw(msg, fields...)` |
| PANIC | `Panic(msg)` | `Panicf(fmt, args...)` | `Panicw(msg, fields...)` |

### 多种格式输出

```go
// Def 格式（默认）
logger := fastlog.New(fastlog.NewConfig("logs/app.log"))
// 输出: 2025-01-15T10:30:45 | INFO    | main.go:main:15 - 用户登录成功

// JSON 格式
cfg := fastlog.NewConfig("logs/app.log")
cfg.Formatter = fastlog.JSON{}
logger := fastlog.New(cfg)
// 输出: {"time":"2025-01-15T10:30:45","level":"INFO","message":"用户登录成功"}

// Simple 格式
cfg.Formatter = fastlog.Simple{}
// 输出: 2025-01-15T10:30:45 INFO 用户登录成功

// KV 格式
cfg.Formatter = fastlog.KV{}
// 输出: time=2025-01-15T10:30:45 level=INFO message=用户登录成功

// Compact 格式（时间格式遵循 TimeFormat，默认 RFC3339）
cfg.Formatter = fastlog.Compact{}
// 输出: [I] 2025-01-15T10:30:45Z 用户登录成功 | username=alice count=42
```

### 动态设置日志级别

运行时动态调整日志级别，无需重启程序，立即生效：

```go
logger := fastlog.New(fastlog.Prod("/var/log/app.log"))

// 初始 WARN 级别
logger.Info("这条不会输出")  // 被过滤

// 动态调整为 DEBUG 级别
logger.SetLevel(fastlog.DEBUG)

// 后续日志立即使用新级别
logger.Info("现在可以输出了")  // 输出
logger.Debug("调试信息")       // 输出

// 获取当前级别
currentLevel := logger.Level()
fmt.Println(currentLevel)  // DEBUG
```

**使用场景：**
- 生产环境出问题，临时开启 DEBUG 排查，无需重启服务
- 根据系统负载动态调整日志详细程度
- 通过 HTTP API 热更新日志级别

### 缓冲控制

通过 `BufferEnabled` 控制是否启用缓冲写入：

```go
// 默认配置启用缓冲（批量写入，性能更好）
logger := fastlog.New(fastlog.NewConfig("logs/app.log"))

// 开发环境自动禁用缓冲（Dev 函数内部设置 BufferEnabled = false）
logger := fastlog.New(fastlog.Dev("logs/dev.log"))

// 手动禁用缓冲（立即落盘，数据更安全）
cfg := fastlog.NewConfig("logs/app.log")
cfg.BufferEnabled = false
logger := fastlog.New(cfg)
```

**缓冲控制对比：**

| 模式 | BufferEnabled | 特点 | 适用场景 |
|------|---------------|------|----------|
| 缓冲写入 | `true` | 批量写入磁盘，性能更好，有延迟 | 生产环境 |
| 直接写入 | `false` | 立即落盘，数据安全，无延迟 | 开发调试、高可靠性场景 |

### 彩色输出

```go
// 默认配置已启用彩色输出
logger := fastlog.New(fastlog.NewConfig("logs/app.log"))

// 禁用彩色输出
cfg := fastlog.NewConfig("logs/app.log")
cfg.NoColor = true
logger := fastlog.New(cfg)
```

**颜色映射：**

| 级别 | 颜色 | 效果 |
|------|------|------|
| DEBUG | 青色 | `FgCyan + Bold` |
| INFO | 蓝色 | `FgBlue + Bold` |
| WARN | 黄色 | `FgYellow + Bold` |
| ERROR | 红色加粗 | `FgRed + Bold` |
| FATAL | 红色加粗 | `FgRed + Bold` |
| PANIC | 紫色加粗 | `FgMagenta + Bold` |

### 级别路由

通过 `LevelRouter` 启用级别路由功能，日志会同时写入全量文件和级别专属文件：

```go
// 启用级别路由
cfg := fastlog.NewConfig("logs/app.log")
cfg.LevelRouter = true  // 启用级别路由

logger := fastlog.New(cfg)
defer func() { _ = logger.Close() }()

// 这些日志会同时写入 app.log 和对应的级别文件
logger.Debug("调试信息")  // 写入 app.log + DEBUG.log
logger.Info("业务信息")   // 写入 app.log + INFO.log
logger.Warn("警告信息")   // 写入 app.log + WARN.log
logger.Error("错误信息")  // 写入 app.log + ERROR.log
```

**生成的日志文件：**

| 文件 | 说明 |
|------|------|
| `app.log` | 全量日志，包含所有级别 |
| `DEBUG.log` | 仅包含 DEBUG 级别 |
| `INFO.log` | 仅包含 INFO 级别 |
| `WARN.log` | 仅包含 WARN 级别 |
| `ERROR.log` | 仅包含 ERROR 级别 |
| `FATAL.log` | 仅包含 FATAL 级别 |
| `PANIC.log` | 仅包含 PANIC 级别 |

**使用场景：**
- 快速定位错误：直接查看 ERROR.log，无需在大文件中搜索
- 分离关注点：运维关注 ERROR.log，开发关注 DEBUG.log
- 监控集成：ERROR.log 可直接接入错误监控系统

### 日志采样

```go
// 使用默认采样配置（10秒窗口，前3条放行，之后每10条放行1条）
logger := fastlog.New(fastlog.NewConfig("logs/app.log"))

// 自定义采样配置
cfg := fastlog.NewConfig("logs/app.log")
cfg.SamplerTick = 30 * time.Second      // 30秒窗口
cfg.SamplerInitial = 5                  // 前5条放行
cfg.SamplerThereafter = 20              // 之后每20条放行1条
logger := fastlog.New(cfg)

for i := 0; i < 100; i++ {
    logger.Error("数据库连接超时")     // 大量重复日志仅部分写入
}
```

### 多路输出

```go
// 同时输出到终端和文件（NewConfig 默认行为）
logger := fastlog.New(fastlog.NewConfig("logs/app.log"))

// 手动创建多路输出
cfg := &fastlog.Config{
    Level:         fastlog.INFO,
    Formatter:     fastlog.Def{},
    OutputConsole: true,
    OutputFile:    true,
    LogPath:       "logs/app.log",
}
logger := fastlog.New(cfg)

logger.Info("这条日志同时输出到控制台和文件")
```

---

## 测试

```bash
# 运行所有测试
go test ./...

# 运行并查看详细输出
go test -v ./...

# 运行特定测试
go test -v -run TestSampler ./...

# 检查代码覆盖率
go test -cover ./...
```

当前测试覆盖全部核心文件，按优先级从底层到上层覆盖：

| 优先级 | 被测组件 | 关键测点 |
|--------|----------|----------|
| P0 | 日志级别 | String / Enabled / ParseLevel / AllLevels |
| P0 | 字段系统 | 构造函数 / 取值方法 / 类型不匹配 / nil 边界 |
| P1 | 日志采样 | 放行规则 / 窗口重置 / 独立计数 / 边界值 |
| P1 | 格式化器 | 5 种格式输出模板 / JSON 合法性 / 空 Entry |
| P2 | 写入器 | ColorWriter / MultiWriter / 错误处理 |
| P3 | Logger | 配置生效 / 级别过滤 / 采样集成 / 生命周期 |

---

## 许可证

本项目采用 **MIT 许可证**。详情请参见 [LICENSE](LICENSE) 文件。

---

## 贡献指南

欢迎贡献！请遵循以下流程：

1. **Fork** 本仓库
2. 创建特性分支：`git checkout -b feat/your-feature`
3. 提交改动：`git commit -m "feat: 添加 XXX 功能"`
4. 推送分支：`git push origin feat/your-feature`
5. 创建 **Pull Request**

### 开发规范

- 提交信息遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范
- 新增功能需附带对应的单元测试
- 运行 `go vet ./...` 和 `go test ./...` 确保无错误

---

## 相关链接

- **仓库地址**: [gitee.com/MM-Q/fastlog](https://gitee.com/MM-Q/fastlog)
- **Go 官方**: [go.dev](https://go.dev)
- **依赖库**:
  - [color](https://gitee.com/MM-Q/color) — 终端彩色输出
  - [go-json](https://github.com/goccy/go-json) — 高性能 JSON 序列化
  - [logrotatex](https://gitee.com/MM-Q/logrotatex) — 日志轮转和缓冲写入
  - [comprx](https://gitee.com/MM-Q/comprx) — 日志压缩
- **参考项目**: [zap](https://github.com/uber-go/zap)

---

<p align="center">
  Made with ❤️ by <a href="https://gitee.com/MM-Q">M乔木</a>
  <br>
  <a href="https://gitee.com/MM-Q/fastlog">gitee.com/MM-Q/fastlog</a>
</p>

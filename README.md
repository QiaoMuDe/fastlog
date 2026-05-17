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
| 📝 **多格式支持** | 内置 5 种格式：Def、JSON、Timestamp、KV、LogFmt，支持自定义 |
| 🧩 **结构化字段** | 12 种字段类型，类型安全，零装箱分配 |
| 🎯 **日志采样** | 固定桶 + atomic 无锁设计，参考 zap，有效防洪 |
| 🔌 **多路输出** | `MultiWriter` 同时输出到多个目标 |
| 🧪 **场景化配置** | `NewConfig()`、`Dev()`、`Prod()`、`Console()`、`Docker()` 覆盖常见场景 |
| 🔒 **线程安全** | `sync.Mutex` 保证写入安全 |
| 📦 **一站式集成** | 基于 [logrotatex](https://gitee.com/MM-Q/logrotatex) 实现日志轮转、缓冲写入，[comprx](https://gitee.com/MM-Q/comprx) 实现压缩，用户无感知 |

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

| 函数 | 级别 | 输出 | 格式 | 调用者 | 采样 | 轮转 | 适用场景 |
|------|------|------|------|--------|------|------|----------|
| `NewConfig(path)` | INFO | 终端+文件 | Def | ❌ | ✅ | ✅ | 通用默认 |
| `Default()` | INFO | 终端+文件 | Def | ❌ | ✅ | ✅ | 快速上手 |
| `Dev(path)` | DEBUG | 终端+文件 | Def | ✅ | ❌ | ❌ | 开发调试 |
| `Prod(path)` | WARN | 仅文件 | Def | ❌ | ✅ | ✅ | 生产环境 |
| `Console()` | DEBUG | 仅终端 | Def | ❌ | ❌ | ❌ | 本地调试 |
| `Docker()` | WARN | 仅终端 | JSON | ❌ | ✅ | ❌ | 容器/K8s |

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

### 多种格式

```go
// Def 格式（默认）
logger := fastlog.New(fastlog.NewConfig("logs/app.log"))
// 输出: 2025-01-15T10:30:45 | INFO    | main.go:main:15 - 用户登录成功

// JSON 格式
cfg := fastlog.NewConfig("logs/app.log")
cfg.Formatter = fastlog.JSON{}
logger := fastlog.New(cfg)
// 输出: {"time":"2025-01-15T10:30:45","level":"INFO","message":"用户登录成功"}

// Timestamp 格式
cfg.Formatter = fastlog.Timestamp{}
// 输出: 2025-01-15T10:30:45 INFO 用户登录成功

// KV 格式
cfg.Formatter = fastlog.KV{}
// 输出: time=2025-01-15T10:30:45 level=INFO message=用户登录成功

// LogFmt 格式
cfg.Formatter = fastlog.LogFmt{}
// 输出: 2025-01-15T10:30:45 [INFO ] 用户登录成功 [username=alice, count=42]
```

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
| DEBUG | 青色 | `SCyan` |
| INFO | 蓝色 | `SBlue` |
| WARN | 黄色 | `SYellow` |
| ERROR | 红色 | `SRed` |
| FATAL | 红色加粗 | `FgRed + Bold` |
| PANIC | 紫色加粗 | `FgMagenta + Bold` |

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

## API 文档

### Logger 创建

| 函数 | 说明 |
|------|------|
| `New(cfg *Config) *Logger` | 创建 Logger，传入配置实例 |
| `Default() *Config` | 默认配置（INFO + 双输出 + 采样） |
| `NewConfig(path) *Config` | 通用配置（双输出） |
| `Dev(path) *Config` | 开发配置（DEBUG + Caller） |
| `Prod(path) *Config` | 生产配置（WARN + 仅文件 + 压缩） |
| `Console() *Config` | 控制台配置（DEBUG + 仅终端） |
| `Docker() *Config` | 容器配置（WARN + JSON + stdout） |

### 日志记录方法

每条日志方法对应三种变体：

| 级别 | 标准 | 格式化 | 结构化 |
|------|------|--------|--------|
| DEBUG | `Debug(msg)` | `Debugf(fmt, args...)` | `Debugw(msg, fields...)` |
| INFO | `Info(msg)` | `Infof(fmt, args...)` | `Infow(msg, fields...)` |
| WARN | `Warn(msg)` | `Warnf(fmt, args...)` | `Warnw(msg, fields...)` |
| ERROR | `Error(msg)` | `Errorf(fmt, args...)` | `Errorw(msg, fields...)` |
| FATAL | `Fatal(msg)` | `Fatalf(fmt, args...)` | `Fatalw(msg, fields...)` |
| PANIC | `Panic(msg)` | `Panicf(fmt, args...)` | `Panicw(msg, fields...)` |

**生命周期方法：**

| 方法 | 说明 |
|------|------|
| `Sync() error` | 刷新日志到存储 |
| `Close() error` | 关闭 Logger，释放资源 |

### Config 配置项

| 配置项 | 类型 | 说明 |
|--------|------|------|
| `Level` | `Level` | 日志级别 |
| `Formatter` | `Formatter` | 格式化器 |
| `Caller` | `bool` | 是否记录调用者信息 |
| `Fields` | `[]Field` | 预设字段 |
| `SamplerTick` | `time.Duration` | 采样时间窗口 |
| `SamplerInitial` | `int` | 采样初始放行数 |
| `SamplerThereafter` | `int` | 采样后续放行间隔 |
| `OutputConsole` | `bool` | 是否输出到终端 |
| `NoColor` | `bool` | 是否禁用彩色输出 |
| `OutputFile` | `bool` | 是否输出到文件 |
| `LogPath` | `string` | 日志文件路径 |
| `Async` | `bool` | 是否异步清理 |
| `MaxSize` | `int` | 单文件最大大小（MB） |
| `MaxFiles` | `int` | 保留的历史文件数 |
| `MaxAge` | `int` | 保留天数 |
| `Compress` | `bool` | 是否压缩历史文件 |
| `CompressType` | `comprx.CompressType` | 压缩类型 |
| `LocalTime` | `bool` | 是否使用本地时间 |
| `DateDirLayout` | `bool` | 是否按日期目录存放 |
| `RotateByDay` | `bool` | 是否按天轮转 |
| `MaxBufferSize` | `int` | 缓冲区大小（字节） |
| `SyncInterval` | `time.Duration` | 自动同步间隔 |

### 自定义格式示例

```go
// 自定义格式：实现 Formatter 接口即可
type MyFormatter struct{}

func (f MyFormatter) Format(entry *fastlog.Entry) ([]byte, error) {
    return []byte(entry.Level.String() + ": " + entry.Message + "\n"), nil
}

cfg := fastlog.NewConfig("logs/app.log")
cfg.Formatter = MyFormatter{}
logger := fastlog.New(cfg)
logger.Info("使用自定义格式")     // 输出: INFO: 使用自定义格式
```

---

## 日志格式

| 格式 | 结构体 | 输出示例 |
|------|--------|----------|
| 默认 | `Def` | `2025-01-15T10:30:45 \| INFO \| main.go:main:15 - 消息` |
| JSON | `JSON` | `{"time":"...","level":"INFO","message":"..."}` |
| 时间戳 | `Timestamp` | `2025-01-15T10:30:45 INFO 消息` |
| 键值对 | `KV` | `time=... level=INFO message=...` |
| LogFmt | `LogFmt` | `2025-01-15T10:30:45 [INFO ] 消息 [key=val]` |

---

## 字段类型

| 构造函数 | Go 类型 | 取值方法 |
|----------|---------|----------|
| `fastlog.String(key, val)` | `string` | `Field.String()` |
| `fastlog.Int(key, val)` | `int` | `Field.Int()` |
| `fastlog.Int64(key, val)` | `int64` | `Field.Int64()` |
| `fastlog.Uint(key, val)` | `uint` | `Field.Uint()` |
| `fastlog.Uint64(key, val)` | `uint64` | `Field.Uint64()` |
| `fastlog.Float64(key, val)` | `float64` | `Field.Float64()` |
| `fastlog.Bool(key, val)` | `bool` | `Field.Bool()` |
| `fastlog.Time(key, val)` | `time.Time` | `Field.Time()` |
| `fastlog.Duration(key, val)` | `time.Duration` | `File.DurationVal()` |
| `fastlog.Error(err)` | `error` | 键名固定为 `"error"` |
| `fastlog.Err(key, err)` | `error` | 自定义键名 |
| `fastlog.Any(key, val)` | `interface{}` | `Field.Interface` |
| `fastlog.Stack()` | — | 当前堆栈信息 |

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

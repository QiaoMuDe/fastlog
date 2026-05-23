---
name: fastlog-skill
description: 为 FastLog Go 日志库生成代码的 Skill。当用户需要使用 fastlog 记录日志、创建 Logger 实例、配置日志环境（开发/生产/容器）、使用结构化字段、启用高级功能（采样、级别路由、动态级别调整）时，必须使用此 Skill。适用于所有涉及 fastlog 的代码生成场景，包括新建项目、为现有代码添加日志、配置不同环境的日志行为等。
---

# FastLog 代码生成 Skill

## 概述

FastLog 是一个轻量级、高性能的 Go 日志库，提供三级 API（标准/格式化/结构化）、
6 种场景化配置、5 种输出格式、日志采样、级别路由等高级功能。

## 核心概念

### 日志级别（6个）
- `DEBUG` - 调试信息
- `INFO` - 一般信息
- `WARN` - 警告
- `ERROR` - 错误
- `FATAL` - 致命错误（记录后退出程序）
- `PANIC` - 恐慌（记录后触发 panic）

### 三级 API
每个级别都有三种调用方式：
- **标准**: `logger.Info("消息")`
- **格式化**: `logger.Infof("用户: %s", name)`
- **结构化**: `logger.Infow("用户登录", fastlog.String("user", name), fastlog.Int("age", 25))`

### 场景化配置函数
| 函数 | 适用场景 | 特点 |
|------|----------|------|
| `Default()` | 快速上手 | INFO 级别，终端+文件 |
| `Dev(path)` | 开发调试 | DEBUG 级别，彩色输出，调用者信息 |
| `Prod(path)` | 生产环境 | WARN 级别，仅文件，启用级别路由 |
| `Console()` | 本地调试 | DEBUG 级别，仅终端 |
| `Docker()` | 容器/K8s | WARN 级别，JSON 格式，stdout |
| `NewConfig(path)` | 自定义 | 通用默认，可进一步调整 |

### 结构化字段类型
```go
fastlog.String(key, value)      // 字符串
fastlog.Int(key, value)         // 整数
fastlog.Int64(key, value)       // 64位整数
fastlog.Uint(key, value)        // 无符号整数
fastlog.Uint64(key, value)      // 64位无符号整数
fastlog.Float64(key, value)     // 浮点数
fastlog.Bool(key, value)        // 布尔值
fastlog.Time(key, value)        // 时间
fastlog.Duration(key, value)    // 时长
fastlog.Error(err)              // 错误（key 固定为 "error"）
fastlog.Any(key, value)         // 任意类型
fastlog.Namespace(name)         // 命名空间（用于分组）
```

### 高级功能

#### 1. 日志采样
```go
cfg := fastlog.NewConfig("logs/app.log")
cfg.SamplerTick = 10 * time.Second      // 采样窗口
cfg.SamplerInitial = 3                  // 前3条放行
cfg.SamplerThereafter = 10              // 之后每10条放行1条
```

#### 2. 级别路由（按级别分文件）
```go
cfg := fastlog.NewConfig("logs/app.log")
cfg.LevelRouter = true  // 启用后，日志同时写入 app.log 和 {LEVEL}.log
// 生成: app.log, DEBUG.log, INFO.log, WARN.log, ERROR.log, FATAL.log, PANIC.log
```

#### 3. 动态级别调整
```go
logger.SetLevel(fastlog.DEBUG)  // 运行时调整级别，立即生效
currentLevel := logger.Level()  // 获取当前级别
```

#### 4. 多格式输出
```go
cfg.Formatter = fastlog.Def{}      // 默认格式
cfg.Formatter = fastlog.JSON{}     // JSON 格式
cfg.Formatter = fastlog.Simple{}   // 简单格式
cfg.Formatter = fastlog.KV{}       // Key=Value 格式
cfg.Formatter = fastlog.Compact{}  // 紧凑格式
```

## 代码生成规则

### 1. 基础使用模板（推荐）

```go
package main

import (
    "gitee.com/MM-Q/fastlog"
)

func main() {
    // 创建 Logger - 使用默认配置（推荐）
    logger := fastlog.New(fastlog.Default())
    defer func() { _ = logger.Close() }()

    // 记录日志
    logger.Info("应用启动成功")
}
```

### 2. 不同环境的完整示例

#### 开发环境示例
```go
package main

import (
    "gitee.com/MM-Q/fastlog"
)

func main() {
    // 开发环境：DEBUG 级别、彩色输出、带调用者信息
    logger := fastlog.New(fastlog.Dev("logs/dev.log"))
    defer func() { _ = logger.Close() }()

    logger.Debug("开始初始化数据库连接")
    logger.Infow("数据库连接成功",
        fastlog.String("host", "localhost"),
        fastlog.Int("port", 3306),
    )
}
```

#### 生产环境示例
```go
package main

import (
    "gitee.com/MM-Q/fastlog"
)

func main() {
    // 生产环境：WARN 级别、仅文件、级别路由、压缩
    logger := fastlog.New(fastlog.Prod("/var/log/myapp/app.log"))
    defer func() { _ = logger.Close() }()

    // INFO 不会输出（WARN 级别以上才输出）
    logger.Info("这条不会记录")
    
    // WARN 和 ERROR 会输出到 app.log 和对应的级别文件
    logger.Warn("连接池接近上限")
    logger.Errorw("数据库查询失败",
        fastlog.String("sql", "SELECT * FROM users"),
        fastlog.Error(err),
    )
    // 同时生成：app.log、WARN.log、ERROR.log
}
```

#### 容器环境示例
```go
package main

import (
    "gitee.com/MM-Q/fastlog"
)

func main() {
    // 容器环境：JSON 格式输出到 stdout
    logger := fastlog.New(fastlog.Docker())
    defer func() { _ = logger.Close() }()

    logger.Info("容器启动成功")
    // 输出: {"time":"2025-01-15T10:30:45Z","level":"INFO","message":"容器启动成功"}
}
```

### 3. 必须遵循的规范

- **必须**使用 `defer func() { _ = logger.Close() }()` 关闭 Logger
- **必须**通过 `fastlog.New(cfg)` 创建 Logger，不要直接声明结构体
- 使用结构化日志（`Infow` 等）时，**必须**使用 fastlog 提供的字段构造函数
- 生产环境**推荐**使用 `Prod()` 配置，自动启用级别路由
- 需要动态调整级别时，**必须**使用 `atomic` 安全的 `SetLevel()` 方法

### 4. 根据场景选择配置（优先级顺序）

**优先使用场景化快捷函数**，它们已经针对特定环境优化好了配置：

#### 快速上手（推荐用于原型和测试）
```go
// 使用默认配置，同时输出到终端和文件
logger := fastlog.New(fastlog.Default())
// 等价于 fastlog.New(fastlog.NewConfig("logs/app.log"))
```

#### 开发环境（推荐用于本地开发）
```go
// DEBUG 级别、彩色输出、调用者信息、无采样
logger := fastlog.New(fastlog.Dev("logs/dev.log"))
```

#### 生产环境（推荐用于线上服务）
```go
// WARN 级别、仅文件输出、启用级别路由、异步清理、压缩
logger := fastlog.New(fastlog.Prod("/var/log/app.log"))
```

#### 容器环境（推荐用于 Docker/K8s）
```go
// WARN 级别、JSON 格式、仅输出到 stdout
logger := fastlog.New(fastlog.Docker())
```

#### 纯控制台（推荐用于临时调试）
```go
// DEBUG 级别、仅终端输出、无文件
logger := fastlog.New(fastlog.Console())
```

#### 自定义配置（当快捷函数不满足需求时使用）
```go
// 基于 NewConfig 修改特定配置
cfg := fastlog.NewConfig("logs/app.log")
cfg.Level = fastlog.DEBUG
cfg.Caller = true
cfg.LevelRouter = true
// ... 其他自定义配置

logger := fastlog.New(cfg)
```

**选择建议：**
- 不确定用什么？→ 先用 `Default()`
- 本地开发调试？→ 用 `Dev()`
- 线上生产环境？→ 用 `Prod()`
- 跑在容器里？→ 用 `Docker()`
- 需要特殊配置？→ 基于 `NewConfig()` 自定义

### 5. 结构化日志最佳实践

```go
// 好的做法：使用具体类型字段
logger.Infow("用户登录",
    fastlog.String("username", user.Name),
    fastlog.Int("user_id", user.ID),
    fastlog.Duration("elapsed", time.Since(start)),
    fastlog.Bool("success", true),
)

// 避免：在消息中拼接变量
logger.Infof("用户 %s 登录成功，耗时 %v", user.Name, elapsed)

// 避免：使用 Any 代替具体类型
logger.Infow("用户登录", fastlog.Any("user", user))  // 性能较差
```

### 6. 错误处理日志

```go
if err != nil {
    logger.Errorw("操作失败",
        fastlog.String("operation", "database_query"),
        fastlog.Error(err),
        fastlog.String("sql", query),
    )
    return err
}
```

### 7. 采样配置（高频率日志）

```go
// 对于可能产生大量重复日志的操作，启用采样
cfg := fastlog.NewConfig("logs/app.log")
cfg.SamplerTick = 10 * time.Second

logger := fastlog.New(cfg)

// 大量重复调用，只有部分会被记录
for i := 0; i < 10000; i++ {
    logger.Debug("循环处理")  // 受采样控制
}
```

## 完整示例

### 示例 1: Web 服务日志

```go
package main

import (
    "net/http"
    "time"
    
    "gitee.com/MM-Q/fastlog"
)

func main() {
    // 生产环境配置
    logger := fastlog.New(fastlog.Prod("logs/server.log"))
    defer func() { _ = logger.Close() }()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // 处理请求...
        
        logger.Infow("请求处理完成",
            fastlog.String("method", r.Method),
            fastlog.String("path", r.URL.Path),
            fastlog.String("ip", r.RemoteAddr),
            fastlog.Duration("elapsed", time.Since(start)),
        )
    })

    logger.Info("服务器启动，监听 :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        logger.Fatalw("服务器启动失败", fastlog.Error(err))
    }
}
```

### 示例 2: 动态调整日志级别

```go
package main

import (
    "net/http"
    
    "gitee.com/MM-Q/fastlog"
)

func main() {
    cfg := fastlog.NewConfig("logs/app.log")
    cfg.Level = fastlog.WARN  // 初始 WARN 级别
    
    logger := fastlog.New(cfg)
    defer func() { _ = logger.Close() }()

    // HTTP 接口动态调整级别
    http.HandleFunc("/setlevel", func(w http.ResponseWriter, r *http.Request) {
        level := r.URL.Query().Get("level")
        switch level {
        case "debug":
            logger.SetLevel(fastlog.DEBUG)
        case "info":
            logger.SetLevel(fastlog.INFO)
        case "warn":
            logger.SetLevel(fastlog.WARN)
        case "error":
            logger.SetLevel(fastlog.ERROR)
        }
        w.Write([]byte("Level set to: " + logger.Level().String()))
    })

    http.ListenAndServe(":8080", nil)
}
```

### 示例 3: 级别路由功能

```go
package main

import "gitee.com/MM-Q/fastlog"

func main() {
    cfg := fastlog.NewConfig("logs/app.log")
    cfg.LevelRouter = true  // 启用级别路由
    cfg.Level = fastlog.DEBUG
    
    logger := fastlog.New(cfg)
    defer func() { _ = logger.Close() }()

    // 这些日志会分别写入不同文件
    logger.Debug("调试信息")  // app.log + DEBUG.log
    logger.Info("业务信息")   // app.log + INFO.log
    logger.Warn("警告信息")   // app.log + WARN.log
    logger.Error("错误信息")  // app.log + ERROR.log
    
    // ERROR.log 只包含错误日志，便于快速定位问题
}
```

## 输出格式

### Def 格式（默认）
```
2025-01-15T10:30:45+08:00 | INFO    | main.go:main:15 - 用户登录成功
```

### JSON 格式
```json
{"time":"2025-01-15T10:30:45+08:00","level":"INFO","message":"用户登录成功","username":"alice"}
```

### Simple 格式
```
2025-01-15T10:30:45 INFO 用户登录成功
```

### KV 格式
```
time=2025-01-15T10:30:45+08:00 level=INFO message=用户登录成功 username=alice
```

### Compact 格式
```
[I] 2025-01-15T10:30:45Z 用户登录成功 | username=alice
```

## 注意事项

1. **不要**重复使用已关闭的 Logger
2. **不要**在多个 goroutine 中共享 Logger 后关闭它（确保 defer 在主流程）
3. **不要**直接修改 Config 后继续使用同一个 Logger，需要重新创建
4. Fatal 和 Panic 方法在记录日志后会退出程序或触发 panic，**不要**在这之后执行清理代码
5. 使用 `Any` 字段类型会有反射开销，性能敏感场景应使用具体类型

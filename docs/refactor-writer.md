# 一站式重构方案

> 改动方向：从"纯日志引擎"改为"一站式日志库"，内置文件轮转和缓冲写入

---

## 一、核心思路

**去掉 `WithWriter`，用 `Config` 结构体声明输出方式。**

```
❌ 旧：用户注入写入器 → Logger 搞不清该不该关
   logger := fastlog.New(fastlog.WithWriter(someWriter))

✅ 新：用户声明需求 → Logger 自己造自己管
   logger := fastlog.New(fastlog.Config{
       Console: true,
       LogPath: "logs/app.log",
   })
   defer logger.Close()  // 内部统一关闭所有资源
```

Logger 内部管理一个 `writers []io.WriteCloser` 切片，`log()` 写入所有 writer，`Close()` 关闭所有 writer。

---

## 二、Config 结构体设计

```go
// Config 日志记录器配置
//
// 零值表示使用默认值或关闭该功能。
// 设置 Console=true 启用终端输出，设置 LogPath 启用文件输出。
// 两者可同时启用，日志会同时写入终端和文件。
type Config struct {
    // ======== 基础日志配置 ========

    // Level 日志级别，零值默认 INFO
    Level Level

    // Formatter 日志格式化器，零值默认 Def
    Formatter Formatter

    // Caller 是否记录调用者信息（文件:函数:行号）
    Caller bool

    // Fields 预设字段，每条日志都会带上
    Fields []Field

    // Sampler 日志采样器，nil 表示不采样
    Sampler *Sampler

    // ======== 终端输出配置 ========

    // Console 是否输出到终端（彩色自动检测）
    Console bool

    // NoColor 设为 true 禁用终端彩色输出
    NoColor bool

    // ======== 文件输出配置 ========

    // LogPath 日志文件路径，非空时自动启用文件输出
    // 内部自动创建日志轮转器 + 缓冲写入器
    LogPath string

    // MaxSize 单文件最大大小（MB），超过后轮转。零值默认 10MB
    MaxSize int

    // MaxFiles 保留的历史文件数，零值表示不限制
    MaxFiles int

    // MaxAge 保留天数，超过的旧文件自动清理，零值表示不限制
    MaxAge int

    // Compress 是否压缩历史日志文件
    Compress bool

    // LocalTime 是否使用本地时间命名轮转文件，默认 true
    LocalTime bool

    // DateDirLayout 是否按日期目录存放轮转文件
    DateDirLayout bool

    // RotateByDay 是否按天轮转
    RotateByDay bool

    // ======== 缓冲写入配置 ========

    // MaxBufferSize 缓冲区大小（字节），零值默认 256KB
    MaxBufferSize int

    // FlushInterval 自动同步间隔，零值默认 1 秒，最小 500ms
    FlushInterval time.Duration
}
```

---

## 三、Logger 结构体设计

```go
type Logger struct {
    // 配置（保存用户传入的配置，Close 后不再使用）
    config Config

    // 内部写入器列表（由 Logger 统一管理生命周期）
    writers []io.WriteCloser

    // 日志核心
    level     Level
    formatter Formatter
    caller    bool
    fields    []Field
    sampler   *Sampler

    // 线程安全
    mu   sync.Mutex
    once sync.Once
}
```

### init() 逻辑

```go
func (l *Logger) init() {
    l.once.Do(func() {
        // === 1. 基础配置默认值 ===
        if l.config.Level == 0 {
            l.level = INFO
        } else {
            l.level = l.config.Level
        }

        if l.config.Formatter == nil {
            l.formatter = Def{}
        } else {
            l.formatter = l.config.Formatter
        }

        l.caller = l.config.Caller
        l.fields = l.config.Fields
        l.sampler = l.config.Sampler

        // === 2. 构建写入器列表 ===

        // 终端输出
        if l.config.Console {
            l.writers = append(l.writers, NewColorWriter(l.config.NoColor))
        }

        // 文件输出
        if l.config.LogPath != "" {
            l.writers = append(l.writers, newFileWriter(l.config))
        }

        // 默认：没有配置任何写入器时，输出到终端
        if len(l.writers) == 0 {
            l.writers = append(l.writers, NewColorWriter(false))
        }
    })
}
```

### newFileWriter 内部实现

```go
func newFileWriter(cfg Config) io.WriteCloser {
    rotator := &logrotatex.LogRotateX{
        LogFilePath:   cfg.LogPath,
        MaxSize:       cfg.MaxSize,
        MaxFiles:      cfg.MaxFiles,
        MaxAge:        cfg.MaxAge,
        Compress:      cfg.Compress,
        LocalTime:     cfg.LocalTime,
        DateDirLayout: cfg.DateDirLayout,
        RotateByDay:   cfg.RotateByDay,
    }

    bufCfg := &logrotatex.BufCfg{
        MaxBufferSize: cfg.MaxBufferSize,
        FlushInterval: cfg.FlushInterval,
    }

    return logrotatex.NewBufferedWriter(rotator, bufCfg)
}
```

### log() 和 Close()

```go
func (l *Logger) log(level Level, msg string, fields []Field) {
    l.init()

    if !l.level.Enabled(level) {
        return
    }

    if l.sampler != nil && !l.sampler.Allow(level, msg) {
        return
    }

    entry := GetEntry()
    defer PutEntry(entry)

    entry.Time = time.Now()
    entry.Level = level
    entry.Message = msg
    entry.Fields = append(entry.Fields[:0], l.fields...)
    entry.Fields = append(entry.Fields, fields...)

    if l.caller {
        entry.Caller = getCaller(callerSkip)
    }

    data, err := l.formatter.Format(entry)
    if err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "format error: %v\n", err)
        return
    }

    l.mu.Lock()
    for _, w := range l.writers {
        _, err = w.Write(data)
        if err != nil {
            _, _ = fmt.Fprintf(os.Stderr, "write error: %v\n", err)
        }
    }
    l.mu.Unlock()
}

// Close 关闭日志记录器，释放所有内部资源（文件句柄、定时器协程等）
//
// Close 调用后不应再写入日志。
func (l *Logger) Close() error {
    l.init()
    var errs []error
    for _, w := range l.writers {
        if err := w.Close(); err != nil {
            errs = append(errs, err)
        }
    }
    l.writers = nil
    if len(errs) > 0 {
        return fmt.Errorf("close errors: %v", errs)
    }
    return nil
}
```

---

## 四、Option 设计

提供 3 个构造方式，覆盖所有场景：

### 4.1 `New(cfg Config)` — 标准构造，Config 集中配置

```go
logger := fastlog.New(fastlog.Config{
    Console:  true,                   // 终端彩色
    LogPath:  "logs/app.log",         // 文件轮转
    MaxSize:  100,
    MaxFiles: 10,
    MaxAge:   30,
    Compress: true,
    Level:    fastlog.INFO,
    Formatter: fastlog.JSON{},
    Sampler:  fastlog.DefaultSampler(),
})
defer logger.Close()
```

### 4.2 `Default()` — 终端输出，带采样

```go
logger := fastlog.Default()
defer logger.Close()
```

等价于：

```go
fastlog.New(fastlog.Config{Console: true})
```

### 4.3 `Dev()` — 开发模式，终端彩色+DEBUG+调用者

```go
logger := fastlog.Dev()
defer logger.Close()
```

等价于：

```go
fastlog.New(fastlog.Config{
    Console: true,
    Level:   fastlog.DEBUG,
    Caller:  true,
})
```

---

## 五、用户使用示例

### 终端 + 文件同时输出

```go
logger := fastlog.New(fastlog.Config{
    Console: true,              // 终端彩色
    LogPath: "logs/app.log",    // 文件轮转
    MaxSize: 100,
})
defer logger.Close()

logger.Info("这条日志终端和文件都能看到")
```

### 仅文件输出（生产环境）

```go
logger := fastlog.New(fastlog.Config{
    LogPath:   "logs/app.log",
    MaxSize:   100,
    MaxFiles:  10,
    MaxAge:    30,
    Compress:  true,
    Formatter: fastlog.JSON{},
})
defer logger.Close()
logger.Info("仅写入文件，终端不输出")
```

### 仅终端输出（开发环境）

```go
logger := fastlog.New(fastlog.Config{
    Console: true,
    Caller:  true,
})
defer logger.Close()
```

---

## 六、移除的内容

| 移除项 | 原因 |
|--------|------|
| `WithWriter(w io.WriteCloser)` | 外部写入器生命周期无法管理 |
| `WithLevel` / `WithFormatter` / `WithCaller` 等 Option | 合并到 Config 结构体 |
| `ConsoleWriter` 导出类型 | 不再需要外部包装 |
| `MultiWriter` 导出类型 | 内部用 writers 切片管理，不对外暴露 |
| 外部的 `Nop()` 函数 | 无实际用途 |
| `EntryPool` 导出 | 改为内部使用 |
| `GetEntry()` / `PutEntry()` 导出 | 同上 |

---

## 七、不修改的内容

| 保留项 | 原因 |
|--------|------|
| `Level` 类型和常量 | 核心类型，无变化 |
| `Field` 字段系统和构造函数 | 核心功能，无变化 |
| 5 种 `Formatter`（Def/JSON/Timestamp/KV/LogFmt） | 核心功能，无变化 |
| `Sampler` + `DefaultSampler` | 核心功能，无变化 |
| `ColorWriter` | 内部使用，也可以保留导出作为独立写入器 |
| 三级 API（Info/Infof/Infow） | 核心 API，无变化 |
| `Sync()` 方法 | 保留，委托给支持 Sync 的写入器 |

---

## 八、改动汇总

| 文件 | 改动程度 | 说明 |
|------|----------|------|
| `fastlog.go` | ✅ 少量改动 | Config 结构体定义放这里，或者新建文件 |
| `logger.go` | 🔴 大量改动 | 去掉 Option 函数，Logger 结构体扩写 writers，init/log/Close 重写 |
| `field.go` | ✅ 无变化 | 不动 |
| `formatter.go` | ✅ 无变化 | 不动 |
| `writer.go` | 🔴 大量改动 | 去掉 ConsoleWriter/MultiWriter 导出，保留 ColorWriter |
| `sampler.go` | ✅ 无变化 | 不动 |
| `go.mod` | ✅ 加 logrotatex | 加回依赖 |
| 测试文件 | 🔴 大量改动 | 适配新 API |
| 示例文件 | ✅ 少量改动 | 改用新 Config API |

---

## 九、三种方案对比

| 维度 | 当前（纯引擎） | 方案 A（WithWriter + FileLogger） | 方案 B（Config 一站式） |
|------|---------------|----------------------------------|------------------------|
| **API 风格** | `WithXxx` 函数式 | `WithXxx` + `FileLogger()` | **Config 结构体** |
| **写入器生命周期** | ❌ 模糊 | ⚠️ 委托给用户 | ✅ **Logger 全权管理** |
| **文件日志** | ❌ 不自带 | ✅ `FileLogger()` | ✅ `Config.LogPath` |
| **终端+文件同时** | ❌ MultiWriter 手动 | ❌ 手动组合 | ✅ **Config 字段声明** |
| **自定义写入器** | ✅ `WithWriter` | ✅ `WithWriter` | ❌ 不支持 |
| **Option 数量** | 6 个 | 6 个 + 文件相关 | **Config 一个结构体** |
| **学习成本** | 中（要理解写入器接口） | 中 | **低（一个结构体搞定）** |


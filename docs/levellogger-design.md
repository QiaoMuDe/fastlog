# LevelLogger 级别路由日志记录器设计方案

> 设计目标：完全复用现有 Config，cfg.LogPath 作为全量日志路径，级别文件自动创建在其同级目录
> 设计日期：2026-05-21
> 方案版本：方式 A（复用 Config.LogPath）

---

## 一、概述

### 1.1 背景

目前 FastLog 的 `Logger` 只能将所有日志输出到单一目标（或同时输出到终端+文件）。但在生产环境中，常见的需求是按级别将日志拆分到不同文件：

```
INFO+  → logs/app.log        # 全量日志（用户指定的路径）
ERROR  → logs/ERROR.log      # 错误日志单独存放
WARN   → logs/WARN.log       # 警告日志单独存放
```

当前用户需手动创建多个 Logger 实例并自行判断用哪个，容易出现漏报或误用。

### 1.2 目标

1. **完全复用 Config**：不新增配置结构体，不新增构造函数参数
2. **自动推导目录**：从 `cfg.LogPath` 提取目录，级别文件创建在同级目录
3. **透明体验**：方法签名与 Logger 完全一致，可无缝替换
4. **零侵入**：不修改现有 Logger/Config/Formatter 核心代码

### 1.3 非目标

- 不改变 Logger 内部实现
- 不涉及运行时增减级别路由
- 不涉及采样器共享（子 Logger 默认关闭采样）
- 不涉及 With 方法/Context 集成
- 不涉及自动堆栈跟踪/错误类型识别

---

## 二、架构设计

### 2.1 核心原理

```
LevelLogger (level=INFO)       ← 调用者只跟这一个交互
    │
    ├── all.Logger             ← 所有 >= INFO 的日志 → cfg.LogPath (用户指定)
    │
    ├── INFO.Logger            ← INFO 级别 → {dir}/INFO.log
    ├── WARN.Logger            ← WARN 级别 → {dir}/WARN.log
    └── ERROR.Logger           ← ERROR 级别 → {dir}/ERROR.log
                                ← 低于 LevelLogger.level 的级别不创建专属文件
```

### 2.2 文件自动生成逻辑

接收 `*Config`，从 `cfg.LogPath` 提取目录，从 `cfg.Level` 推导需要创建哪些级别文件：

```
NewLevelLogger(fastlog.Dev("logs/app.log"))
    cfg.LogPath = "logs/app.log"
    cfg.Level = DEBUG
    ↓
    提取目录: "logs/"
    ↓
    创建: logs/app.log        ← 所有 >= DEBUG 的日志（用户指定路径）
    创建: logs/DEBUG.log      ← 专属
    创建: logs/INFO.log       ← 专属
    创建: logs/WARN.log       ← 专属
    创建: logs/ERROR.log      ← 专属
    创建: logs/FATAL.log      ← 专属
    创建: logs/PANIC.log      ← 专属

NewLevelLogger(fastlog.Prod("/var/log/myapp/error.log"))
    cfg.LogPath = "/var/log/myapp/error.log"
    cfg.Level = WARN
    ↓
    提取目录: "/var/log/myapp/"
    ↓
    创建: /var/log/myapp/error.log   ← 所有 >= WARN 的日志（用户指定路径）
    创建: /var/log/myapp/WARN.log    ← 专属
    创建: /var/log/myapp/ERROR.log   ← 专属
    创建: /var/log/myapp/FATAL.log   ← 专属
    创建: /var/log/myapp/PANIC.log   ← 专属
    (INFO, DEBUG 不创建专属文件，只走 error.log)
```

### 2.3 数据流

```
LevelLogger.Info("msg")
    │
    ├── 级别检查 (INFO >= LevelLogger.level?)
    │   └── 否 → 丢弃
    │
    ├── all.Logger.log(INFO, msg)     → cfg.LogPath (始终写入)
    │
    └── INFO.Logger 存在? → log(INFO, msg) → {dir}/INFO.log
    └── WARN.Logger  存在? → 不匹配，跳过
    └── ERROR.Logger 存在? → 不匹配，跳过

LevelLogger.Error("msg", field)
    │
    ├── 级别检查 (ERROR >= LevelLogger.level?) → 通过
    │
    ├── all.Logger.log(ERROR, msg, field) → cfg.LogPath
    │
    ├── INFO.Logger  存在? → 不匹配，跳过
    ├── WARN.Logger   存在? → 不匹配，跳过
    └── ERROR.Logger 存在? → log(ERROR, msg, field) → {dir}/ERROR.log
```

---

## 三、API 设计

### 3.1 设计理念

**零新增参数**，完全复用现有 `*Config`。构造函数只接收 Config：

```go
// 普通 Logger 的创建方式
cfg := fastlog.Dev("logs/app.log")
log := fastlog.New(cfg)

// LevelLogger 的创建方式 — 完全相同的 Config
ll := fastlog.NewLevelLogger(fastlog.Dev("logs/app.log"))
// 内部自动:
//   logs/app.log      ← 继承 Dev 配置，作为全量日志（用户指定路径）
//   logs/DEBUG.log    ← 专属
//   logs/INFO.log     ← 专属
//   ...
```

### 3.2 构造函数

```go
// NewLevelLogger 创建级别路由日志记录器
//
// 参数:
//   - cfg: 常规日志配置，cfg.LogPath 作为全量日志路径，
//          其所在目录用于存放各级别专属日志文件,
//          cfg.Level 决定创建哪些级别文件
//
// 返回:
//   - *LevelLogger: 级别路由日志记录器
//
// 注意:
//   - cfg 为 nil 时 panic
//   - cfg.LogPath 为空时 panic（无法推导目录）
//
// 示例:
//
//	// 开发环境：DEBUG+ 全量 + 各级别专属
//	ll := fastlog.NewLevelLogger(fastlog.Dev("logs/app.log"))
//	defer ll.Close()
//
//	ll.Debug("SQL 查询")  // → logs/app.log + logs/DEBUG.log
//	ll.Info("用户登录")    // → logs/app.log + logs/INFO.log
//	ll.Error("连接失败")   // → logs/app.log + logs/ERROR.log
func NewLevelLogger(cfg *Config) *LevelLogger
```

### 3.3 方法集（与 Logger 完全一致）

```go
// 6 × 标准日志
func (ll *LevelLogger) Debug(msg string)
func (ll *LevelLogger) Info(msg string)
func (ll *LevelLogger) Warn(msg string)
func (ll *LevelLogger) Error(msg string)
func (ll *LevelLogger) Fatal(msg string)
func (ll *LevelLogger) Panic(msg string)

// 6 × 格式化日志
func (ll *LevelLogger) Debugf(format string, args ...interface{})
func (ll *LevelLogger) Infof(format string, args ...interface{})
func (ll *LevelLogger) Warnf(format string, args ...interface{})
func (ll *LevelLogger) Errorf(format string, args ...interface{})
func (ll *LevelLogger) Fatalf(format string, args ...interface{})
func (ll *LevelLogger) Panicf(format string, args ...interface{})

// 6 × 结构化日志
func (ll *LevelLogger) Debugw(msg string, fields ...Field)
func (ll *LevelLogger) Infow(msg string, fields ...Field)
func (ll *LevelLogger) Warnw(msg string, fields ...Field)
func (ll *LevelLogger) Errorw(msg string, fields ...Field)
func (ll *LevelLogger) Fatalw(msg string, fields ...Field)
func (ll *LevelLogger) Panicw(msg string, fields ...Field)

// 管理方法
func (ll *LevelLogger) SetLevel(level Level)
func (ll *LevelLogger) Level() Level
func (ll *LevelLogger) Sync() error
func (ll *LevelLogger) Close() error
```

---

## 四、内部实现

### 4.1 LevelLogger 结构体

```go
type LevelLogger struct {
	level  atomic.Int32      // 运行时日志级别
	all    *Logger           // 全量日志记录器（写入 cfg.LogPath）
	levels map[Level]*Logger // 级别专属日志记录器
}
```

### 4.2 核心 log 方法

```go
// log 内部路由核心
//
// 参数:
//   - level: 日志级别
//   - msg: 日志消息
//   - fields: 结构化字段
func (ll *LevelLogger) log(level Level, msg string, fields []Field) {
    // 级别检查 — 由 LevelLogger 统一控制
    if !Level(ll.level.Load()).Enabled(level) {
        return
    }

    // 写入全量文件（始终执行）
    ll.all.log(level, msg, fields)

    // 写入级别专属文件（如果配置了）
    if specific, ok := ll.levels[level]; ok {
        specific.log(level, msg, fields)
    }
}
```

### 4.3 各方法实现

所有 18 个日志方法统一委托给 `log()`：

```go
func (ll *LevelLogger) Debug(msg string)              { ll.log(DEBUG, msg, nil) }
func (ll *LevelLogger) Debugf(f string, a ...any)     { ll.log(DEBUG, fmt.Sprintf(f, a...), nil) }
func (ll *LevelLogger) Debugw(msg string, f ...Field)  { ll.log(DEBUG, msg, f) }

func (ll *LevelLogger) Info(msg string)               { ll.log(INFO, msg, nil) }
func (ll *LevelLogger) Infof(f string, a ...any)      { ll.log(INFO, fmt.Sprintf(f, a...), nil) }
func (ll *LevelLogger) Infow(msg string, f ...Field)   { ll.log(INFO, msg, f) }

// ... Warn, Error 同理 ...
```

### 4.4 Fatal/Panic 特殊处理

Fatal 和 Panic 不能直接委托给子 Logger 的 Fatal()/Panic() 方法，因为子 Logger 内部会调 `os.Exit(1)` / `panic()`，导致后续的子 Logger 无法执行。

```go
func (ll *LevelLogger) Fatal(msg string) {
    ll.log(FATAL, msg, nil)  // 写入 all + FATAL（如果有）
    _ = ll.Sync()            // 确保全部刷入磁盘
    os.Exit(1)               // 只执行一次
}

func (ll *LevelLogger) Fatalf(format string, args ...interface{}) {
    ll.log(FATAL, fmt.Sprintf(format, args...), nil)
    _ = ll.Sync()
    os.Exit(1)
}

func (ll *LevelLogger) Fatalw(msg string, fields ...Field) {
    ll.log(FATAL, msg, fields)
    _ = ll.Sync()
    os.Exit(1)
}
```

Panic 同理：`ll.log(PANIC, ...)` + `ll.Sync()` + `panic(msg)`。

关键：LevelLogger 的 Fatal 方法直接调用 `all.log()` 和 `levels[FATAL].log()`（内部 log 方法），而非 `all.Fatal()`（公开方法）。由 LevelLogger 自己调 `os.Exit(1)`。

### 4.5 构造函数内部实现

```go
func NewLevelLogger(cfg *Config) *LevelLogger {
    if cfg == nil {
        panic("fastlog: config is required for LevelLogger")
    }
    if cfg.LogPath == "" {
        panic("fastlog: LogPath is required for LevelLogger")
    }

    // 从 cfg.LogPath 提取目录
    // 例: "/var/log/app/error.log" → "/var/log/app/"
    // 例: "logs/app.log" → "logs/"
    dir := filepath.Dir(cfg.LogPath)

    // 从 cfg.Level 推导需要创建哪些级别文件
    // 例: cfg.Level=INFO → 创建 INFO + WARN + ERROR + FATAL + PANIC
    level := cfg.Level
    if level == 0 {
        level = INFO
    }

    // 创建 all Logger（路径: cfg.LogPath，用户指定）
    allCfg := cfg.Clone()
    allCfg.OutputConsole = false  // 级别路由模式默认不输出到终端
    all := New(allCfg)

    // 为 >= level 的每个级别创建专属 Logger（路径: dir/LEVEL.log）
    levels := make(map[Level]*Logger)
    allLevels := AllLevels()
    for _, lvl := range allLevels {
        if lvl < level {
            continue  // 低于设置级别的级别不创建专属文件
        }
        lvlCfg := cfg.Clone()
        lvlCfg.LogPath = filepath.Join(dir, lvl.String()+".log")
        lvlCfg.OutputConsole = false
        lvlCfg.SamplerTick = 0  // 级别专属不启用采样
        lvlCfg.Async = false    // 简化，不异步
        levels[lvl] = New(lvlCfg)
    }

    ll := &LevelLogger{
        all:    all,
        levels: levels,
    }
    ll.level.Store(int32(level))

    return ll
}
```

### 4.6 管理方法

```go
func (ll *LevelLogger) SetLevel(level Level) {
    ll.level.Store(int32(level))
}

func (ll *LevelLogger) Level() Level {
    return Level(ll.level.Load())
}

func (ll *LevelLogger) Sync() error {
    var errs []error
    if err := ll.all.Sync(); err != nil {
        errs = append(errs, err)
    }
    for _, l := range ll.levels {
        if err := l.Sync(); err != nil {
            errs = append(errs, err)
        }
    }
    return errors.Join(errs...)
}

func (ll *LevelLogger) Close() error {
    var errs []error
    if err := ll.all.Close(); err != nil {
        errs = append(errs, err)
    }
    for _, l := range ll.levels {
        if err := l.Close(); err != nil {
            errs = append(errs, err)
        }
    }
    return errors.Join(errs...)
}
```

---

## 五、边界情况与安全

### 5.1 并发安全

- `level` 字段使用 `atomic.Int32`，与 Logger 一致
- 各子 Logger 内部有各自的 `sync.Mutex`，互不干扰
- `levels` map 仅在构造函数中赋值，后续只读，无需锁

### 5.2 退化行为

| 场景 | 行为 |
|------|------|
| cfg 为 nil | panic |
| cfg.LogPath 为空 | panic |
| cfg.Level 为 0 | 默认 INFO，创建 INFO+ 级别文件 |
| 子 Logger 写入失败 | 各自输出错误到 stderr（与 Logger 行为一致）|
| Sync/Close 部分失败 | 收集所有错误，errors.Join 合并返回 |

### 5.3 文件名冲突

文件名固定为 `{LEVEL}.log`（如 `INFO.log`、`ERROR.log`），均为大写。

如果用户指定的 `cfg.LogPath` 恰好是 `{DIR}/INFO.log` 这种形式，则：
- 全量日志写入 `{DIR}/INFO.log`
- INFO 级别的专属文件也是 `{DIR}/INFO.log`

这会导致 INFO 级别日志被写两次到同一文件（但内容相同，只是冗余）。为避免这种情况，构造函数可以检测并 panic：

```go
// 检测冲突
for _, lvl := range allLevels {
    lvlPath := filepath.Join(dir, lvl.String()+".log")
    if lvlPath == cfg.LogPath {
        panic(fmt.Sprintf("fastlog: LogPath %s conflicts with level file path", cfg.LogPath))
    }
}
```

---

## 六、改动清单

### 6.1 新增文件

| 文件 | 内容 | 预估行数 |
|------|------|---------|
| `level_logger.go` | LevelLogger 完整实现 | ~200 行 |
| `level_logger_test.go` | LevelLogger 单元测试 | ~300 行 |

### 6.2 无需改动的文件

- `config.go` — 复用现有 Config
- `logger.go` — 不修改（同包内直接调用 log() 方法，无需改 callerSkip）
- `formatter.go` — 不修改
- `field.go` — 不修改
- `writer.go` — 不修改
- `sampler.go` — 不修改

> 注意：`level_logger.go` 与 `logger.go` 同属 `package fastlog`，可直接调用未导出的 `Logger.log()` 方法，**不需要改动 callerSkip 或暴露新方法**。因为 LevelLogger 的 Info/Warn 等方法直接调 `ll.all.log(level, msg, fields)`，绕过了 Logger 的 Info/Warn 中间层，调用栈不变。

---

## 七、使用示例

### 7.1 开发环境（DEBUG+ 全量 + 各级别专属）

```go
ll := fastlog.NewLevelLogger(fastlog.Dev("logs/app.log"))
defer ll.Close()
defer ll.Sync()

ll.Debug("数据库查询 SQL")    // → logs/app.log + logs/DEBUG.log
ll.Info("用户登录成功")        // → logs/app.log + logs/INFO.log
ll.Warn("磁盘使用率 85%")      // → logs/app.log + logs/WARN.log
ll.Error("连接超时")           // → logs/app.log + logs/ERROR.log

// 产生的文件:
//   logs/app.log      — 所有 >= DEBUG 的日志（用户指定路径）
//   logs/DEBUG.log    — 仅 DEBUG 级别
//   logs/INFO.log     — 仅 INFO 级别
//   logs/WARN.log     — 仅 WARN 级别
//   logs/ERROR.log    — 仅 ERROR 级别
//   logs/FATAL.log    — 仅 FATAL 级别
//   logs/PANIC.log    — 仅 PANIC 级别
```

### 7.2 生产环境（WARN+，只自动创建 WARN/ERROR/FATAL/PANIC）

```go
ll := fastlog.NewLevelLogger(fastlog.Prod("/var/log/myapp/app.log"))

ll.Info("用户登录")      // → /var/log/myapp/app.log only（无 INFO.log）
ll.Warn("磁盘告警")      // → /var/log/myapp/app.log + /var/log/myapp/WARN.log
ll.Error("连接失败")     // → /var/log/myapp/app.log + /var/log/myapp/ERROR.log

// 产生的文件:
//   /var/log/myapp/app.log    — 所有 >= WARN 的日志（用户指定路径）
//   /var/log/myapp/WARN.log   — 仅 WARN 级别
//   /var/log/myapp/ERROR.log  — 仅 ERROR 级别
//   /var/log/myapp/FATAL.log  — 仅 FATAL 级别
//   /var/log/myapp/PANIC.log  — 仅 PANIC 级别
```

### 7.3 配合普通 Logger 混用

```go
// 普通 Logger — 全局系统日志
sysLog := fastlog.New(fastlog.NewConfig("logs/system.log"))

// 级别路由 Logger — 业务日志按级别分文件
bizLog := fastlog.NewLevelLogger(fastlog.Dev("logs/biz/app.log"))

sysLog.Info("服务启动")
bizLog.Info("用户下单")
bizLog.Error("支付失败")
```

---

## 八、方案对比

### 8.1 三种方案对比

| 维度 | 方式 A（当前） | 方式 B（原设计） | 方式 C（混合） |
|------|---------------|-----------------|---------------|
| 构造函数 | `NewLevelLogger(cfg *Config)` | `NewLevelLogger(dir string, cfg *Config)` | `NewLevelLogger(dir string, cfg *Config)` |
| 全量日志路径 | `cfg.LogPath`（用户指定） | `dir/all.log` | 优先 `cfg.LogPath`，否则 `dir/all.log` |
| 级别文件位置 | `filepath.Dir(cfg.LogPath)` | `dir` | `dir` 或 `filepath.Dir(cfg.LogPath)` |
| 参数数量 | 1 个 | 2 个 | 2 个 |
| 学习成本 | 低（完全复用 Config） | 中（新增 dir 参数） | 高（逻辑复杂） |
| 灵活性 | 中（路径由 LogPath 决定） | 高（dir 独立指定） | 高（两者结合） |
| 代码复杂度 | 低 | 低 | 中 |

### 8.2 选择方式 A 的理由

1. **完全复用现有 Config**：不引入新概念，用户学习成本最低
2. **语义清晰**：`cfg.LogPath` 就是全量日志路径，直观易懂
3. **与 Logger 创建方式一致**：都是 `NewXXX(cfg)` 单参数模式
4. **实现简单**：无需处理 dir 和 LogPath 的优先级关系

---

## 九、测试计划

### 9.1 单元测试

| 测试组 | 用例 | 说明 |
|--------|------|------|
| 基础路由 | INFO → app.log + INFO.log | 级别匹配时正确路由到专属文件 |
| 路由过滤 | 专属文件不存在时只走 app.log | 如 WARN 没有专属文件，只写 app.log |
| 级别过滤 | SetLevel(WARN) 后 INFO 不记录 | 级别过滤正常 |
| 多级路由 | ERROR → app.log + ERROR；WARN → app.log + WARN | 多个级别各自独立路由 |
| 路径冲突 | LogPath 与级别文件冲突时 panic | 如 LogPath="logs/INFO.log" |
| Fatal | 写入 app.log + 专属后退出 | 不重复退出 |
| Panic | 写入 app.log + 专属后 panic | 不重复 panic |
| Sync | 遍历所有子 Logger Sync | 错误合并 |
| Close | 遍历所有子 Logger Close | 错误合并 |
| 并发 | 多 goroutine 同时写入 | 不 panic、不丢日志 |

### 9.2 集成测试

- 配合不同场景配置（Dev/Prod/Console/Docker）
- 配合不同 Formatter（Def/JSON）
- 配合 Caller 开关
- 配合日志轮转验证子文件轮转

---

## 十、讨论过程回顾（决策记录）

### 10.1 方案演进

| 版本 | 方案 | 决策 |
|------|------|------|
| 初版 | 新增 `LevelLoggerConfig` 结构体 | ❌ 否决 — 引入新类型，增加学习成本 |
| 二版 | Config 加 `Routes` map 字段 | ❌ 否决 — 侵入 Config，非路由用户也看到此字段 |
| 三版（原设计） | `NewLevelLogger(dir string, cfg *Config)` | ❌ 否决 — 新增 dir 参数，与 Logger 创建方式不一致 |
| 终版（当前） | `NewLevelLogger(cfg *Config)` 复用 LogPath | ✅ 采纳 — 完全复用 Config，学习成本最低 |

### 10.2 调用栈（callerSkip）处理

LevelLogger 直接调用子 Logger 的 `log()` 方法（未导出），而非 `Info()` 等公开方法：

```
调用栈:
    0: getCaller
    1: Logger.log          ← 子 Logger 内部
    2: LevelLogger.log     ← 路由层
    3: 用户代码              ← 正确位置！
```

**不需要修改 Logger 的 callerSkip**，因为：
- `level_logger.go` 和 `logger.go` 同属 `package fastlog`
- LevelLogger 直接调 `all.log(level, msg, fields)` 绕过 Info/Warn 中间层
- 调用栈仅比普通 Logger 深 1 层，但 LevelLogger.log 也是内部方法，用户代码帧位置不变

---

> **设计完成（方式 A 版本）**

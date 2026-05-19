# LevelLogger 级别路由日志记录器设计方案

> 设计目标：接收目录路径 + 常规 Config，内部自动创建级别文件，调用者零感知
> 设计日期：2026-05-19

---

## 一、概述

### 1.1 背景

目前 FastLog 的 `Logger` 只能将所有日志输出到单一目标（或同时输出到终端+文件）。但在生产环境中，常见的需求是按级别将日志拆分到不同文件：

```
INFO+  → logs/all.log        # 全部日志
ERROR  → logs/ERROR.log      # 错误日志单独存放
WARN   → logs/WARN.log       # 警告日志单独存放
```

当前用户需手动创建多个 Logger 实例并自行判断用哪个，容易出现漏报或误用。

### 1.2 目标

1. **目录参数 + 复用 Config**：接收一个目录路径和常规 `*Config`，不新增配置结构体
2. **内部自动创建**：根据 `cfg.Level` 决定创建哪些级别文件，调用者无需手动指定
3. **透明体验**：方法签名与 Logger 完全一致，可无缝替换
4. **零侵入**：不修改现有 Logger/Config/Formatter 核心代码

### 1.3 非目标

- 不改变 Logger 内部实现
- 不涉及运行时增减级别路由
- 不涉及采样器共享（子 Logger 默认关闭采样）
- 不涉及 With 方法/Context 集成（经讨论认为必要性低，Config.Fields 已覆盖）
- 不涉及自动堆栈跟踪/错误类型识别（经讨论认为属于业务逻辑范畴）

---

## 二、架构设计

### 2.1 核心原理

```
LevelLogger (level=INFO)       ← 调用者只跟这一个交互
    │
    ├── all.Logger             ← 所有 >= INFO 的日志 → logs/all.log
    │
    ├── INFO.Logger            ← INFO 级别 → logs/INFO.log
    ├── WARN.Logger            ← WARN 级别 → logs/WARN.log
    └── ERROR.Logger           ← ERROR 级别 → logs/ERROR.log
                                ← 低于 LevelLogger.level 的级别不创建专属文件
```

### 2.2 文件自动生成逻辑

接收目录路径和 `*Config`，从 `cfg.Level` 推导需要创建哪些级别文件：

```
NewLevelLogger("logs", fastlog.Dev(""))
    cfg.Level = DEBUG
    ↓
    创建: logs/all.log        ← 所有 >= DEBUG 的日志
    创建: logs/DEBUG.log      ← 专属
    创建: logs/INFO.log       ← 专属
    创建: logs/WARN.log       ← 专属
    创建: logs/ERROR.log      ← 专属
    创建: logs/FATAL.log      ← 专属
    创建: logs/PANIC.log      ← 专属

NewLevelLogger("logs", fastlog.Prod(""))
    cfg.Level = WARN
    ↓
    创建: logs/all.log        ← 所有 >= WARN 的日志
    创建: logs/WARN.log       ← 专属
    创建: logs/ERROR.log      ← 专属
    创建: logs/FATAL.log      ← 专属
    创建: logs/PANIC.log      ← 专属
    (INFO, DEBUG 不创建专属文件，只走 all.log)
```

### 2.3 数据流

```
LevelLogger.Info("msg")
    │
    ├── 级别检查 (INFO >= LevelLogger.level?)
    │   └── 否 → 丢弃
    │
    ├── all.Logger.log(INFO, msg)     → all.log (始终写入)
    │
    └── INFO.Logger 存在? → log(INFO, msg) → INFO.log
    └── WARN.Logger  存在? → 不匹配，跳过
    └── ERROR.Logger 存在? → 不匹配，跳过

LevelLogger.Error("msg", field)
    │
    ├── 级别检查 (ERROR >= LevelLogger.level?) → 通过
    │
    ├── all.Logger.log(ERROR, msg, field) → all.log
    │
    ├── INFO.Logger  存在? → 不匹配，跳过
    ├── WARN.Logger   存在? → 不匹配，跳过
    └── ERROR.Logger 存在? → log(ERROR, msg, field) → ERROR.log
```

---

## 三、API 设计

### 3.1 设计理念

**不新增配置结构体**，直接复用现有 `*Config`。构造函数接收目录路径和 Config 两参数：

```go
// 对标普通 Logger 的创建方式
cfg := fastlog.Dev("logs/app.log")
log := fastlog.New(cfg)

// LevelLogger 的创建方式 — 目录 + Config
ll := fastlog.NewLevelLogger("logs", fastlog.Dev(""))
// 内部自动:
//   logs/all.log      ← 继承 Dev 配置 (DEBUG, Caller, 10MB, 不压缩)
//   logs/DEBUG.log    ← 同上
//   logs/INFO.log     ← 同上
//   ...
```

### 3.2 构造函数

```go
// NewLevelLogger 创建级别路由日志记录器
//
// 参数:
//   - dir: 日志目录路径，内部自动在该目录下创建 all.log 和各级别 .log 文件
//   - cfg: 常规日志配置，继承其 Formatter/Caller/TimeFormat/轮转等所有设置,
//          cfg.Level 决定创建哪些级别文件
//
// 返回:
//   - *LevelLogger: 级别路由日志记录器
//
// 注意:
//   - dir 为 "" 时 panic
//   - cfg 为 nil 时使用默认配置
//
// 示例:
//
//	// 开发环境：DEBUG+ 全量 + 各级别专属
//	ll := fastlog.NewLevelLogger("logs", fastlog.Dev(""))
//	defer ll.Close()
//
//	ll.Debug("SQL 查询")  // → logs/all.log + logs/DEBUG.log
//	ll.Info("用户登录")    // → logs/all.log + logs/INFO.log
//	ll.Error("连接失败")   // → logs/all.log + logs/ERROR.log
func NewLevelLogger(dir string, cfg *Config) *LevelLogger
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
	all    *Logger           // 全量日志记录器
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
func NewLevelLogger(dir string, cfg *Config) *LevelLogger {
    if dir == "" {
        panic("fastlog: directory is required for LevelLogger")
    }
    if cfg == nil {
        cfg = NewConfig("")
    }

    // 从 cfg.Level 推导需要创建哪些级别文件
    // 例: cfg.Level=INFO → 创建 all + INFO + WARN + ERROR + FATAL + PANIC
    level := cfg.Level
    if level == 0 {
        level = INFO
    }

    // 创建 all Logger（路径: dir/all.log）
    allCfg := cfg.Clone()
    allCfg.LogPath = filepath.Join(dir, "all.log")
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
| dir 为空 | panic |
| cfg 为 nil | 使用默认配置 (NewConfig("")) |
| cfg.Level 为 0 | 默认 INFO，创建 INFO+ 级别文件 |
| 子 Logger 写入失败 | 各自输出错误到 stderr（与 Logger 行为一致）|
| Sync/Close 部分失败 | 收集所有错误，errors.Join 合并返回 |

### 5.3 文件名冲突

文件名固定为 `{LEVEL}.log`（如 `INFO.log`、`ERROR.log`），均为大写。避免与 all.log 冲突。

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
ll := fastlog.NewLevelLogger("logs", fastlog.Dev(""))
defer ll.Close()
defer ll.Sync()

ll.Debug("数据库查询 SQL")    // → logs/all.log + logs/DEBUG.log
ll.Info("用户登录成功")        // → logs/all.log + logs/INFO.log
ll.Warn("磁盘使用率 85%")      // → logs/all.log + logs/WARN.log
ll.Error("连接超时")           // → logs/all.log + logs/ERROR.log

// 产生的文件:
//   logs/all.log      — 所有 >= DEBUG 的日志
//   logs/DEBUG.log    — 仅 DEBUG 级别
//   logs/INFO.log     — 仅 INFO 级别
//   logs/WARN.log     — 仅 WARN 级别
//   logs/ERROR.log    — 仅 ERROR 级别
//   logs/FATAL.log    — 仅 FATAL 级别
//   logs/PANIC.log    — 仅 PANIC 级别
```

### 7.2 生产环境（WARN+，只自动创建 WARN/ERROR/FATAL/PANIC）

```go
ll := fastlog.NewLevelLogger("logs", fastlog.Prod(""))

ll.Info("用户登录")      // → logs/all.log only（无 INFO.log）
ll.Warn("磁盘告警")      // → logs/all.log + logs/WARN.log
ll.Error("连接失败")     // → logs/all.log + logs/ERROR.log

// 产生的文件:
//   logs/all.log      — 所有 >= WARN 的日志
//   logs/WARN.log     — 仅 WARN 级别
//   logs/ERROR.log    — 仅 ERROR 级别
//   logs/FATAL.log    — 仅 FATAL 级别
//   logs/PANIC.log    — 仅 PANIC 级别
```

### 7.4 配合普通 Logger 混用

```go
// 普通 Logger — 全局系统日志
sysLog := fastlog.New(fastlog.NewConfig("logs/system.log"))

// 级别路由 Logger — 业务日志按级别分文件
bizLog := fastlog.NewLevelLogger("logs/biz", fastlog.Dev(""))

sysLog.Info("服务启动")
bizLog.Info("用户下单")
bizLog.Error("支付失败")
```

---

## 八、讨论过程回顾（决策记录）

### 8.1 不实现的功能及原因

| 功能 | 讨论结论 | 原因 |
|------|---------|------|
| With 方法（预绑定字段） | ❌ 不实现 | Config.Fields 已覆盖，函数传参已够用 |
| Context 集成 | ❌ 不实现 | `ctx.Value()` 一行能搞定，不值得引入提取器机制 |
| 请求级 Logger | ❌ 不实现 | With + Context 的延伸，随前两者一起否决 |
| 自动堆栈跟踪 | ❌ 不实现 | Caller 已覆盖 90%，Stack() 按需调用，自动调用开销大 |
| 自定义堆栈深度 | ❌ 不实现 | 依赖自动堆栈跟踪，无独立价值 |
| 错误类型识别 | ❌ 不实现 | 职责错位，业务逻辑应在应用层判断 |
| 延迟格式化 | ✅ 已实现 | Formatter.Format() 延迟已实现，文档标记为已完成 |
| 按级别采样 | ✅ 核心已实现 | 独立计数器空间，级别过滤已覆盖 |
| 按内容采样 | ✅ 核心已实现 | FNV-1a 哈希定桶，不同内容独立采样 |
| 自适应采样 | ❌ 不实现 | 复杂度高收益低，固定采样已够用 |

### 8.2 配置结构体设计演进

设计方案的演化过程：

1. **初版**：新增 `LevelLoggerConfig` 结构体（AllPath + Levels map + Level）
   - 问题：引入新类型，增加学习成本
2. **二版**：复用 `*Config`，Config 加 `Routes` map 字段
   - 问题：侵入 Config，非路由用户也看到此字段
3. **终版（当前）**：**不复用 Config.LogPath**，改为接收目录 + Config
   - 构造函数：`NewLevelLogger(dir string, cfg *Config)`
   - 文件自动生成规则：`dir/{LEVEL}.log`
   - 级别范围推导规则：从 `cfg.Level` 自动计算

### 8.3 调用栈（callerSkip）处理

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

## 九、测试计划

### 9.1 单元测试

| 测试组 | 用例 | 说明 |
|--------|------|------|
| 基础路由 | INFO → all + INFO.log | 级别匹配时正确路由到专属文件 |
| 路由过滤 | 专属文件不存在时只走 all | 如 WARN 没有专属文件，只写 all.log |
| 级别过滤 | SetLevel(WARN) 后 INFO 不记录 | 级别过滤正常 |
| 多级路由 | ERROR → all + ERROR；WARN → all + WARN | 多个级别各自独立路由 |
| Fatal | 写入 all + 专属后退出 | 不重复退出 |
| Panic | 写入 all + 专属后 panic | 不重复 panic |
| Sync | 遍历所有子 Logger Sync | 错误合并 |
| Close | 遍历所有子 Logger Close | 错误合并 |
| 并发 | 多 goroutine 同时写入 | 不 panic、不丢日志 |

### 9.2 集成测试

- 配合不同场景配置（Dev/Prod/Console/Docker）
- 配合不同 Formatter（Def/JSON）
- 配合 Caller 开关
- 配合日志轮转验证子文件轮转

---

> **设计完成**

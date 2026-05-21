## 一、核心思想

LevelLogger 是一个**代理模式**的实现，它本身不直接写日志，而是**将日志分发到多个子 Logger**：

```
用户调用 LevelLogger.Info("msg")
           ↓
    LevelLogger 根据级别做路由
           ↓
    ├── 全量日志 Logger → 写入 cfg.LogPath（所有级别）
    └── 专属日志 Logger → 写入 {dir}/{LEVEL}.log（仅当前级别）
```

---

## 二、内部结构

```go
type LevelLogger struct {
    level  atomic.Int32      // 运行时日志级别（动态可调）
    all    *Logger           // 全量日志记录器
    levels map[Level]*Logger // 级别专属记录器 map[DEBUG/INFO/WARN/ERROR/FATAL/PANIC]*Logger
}
```

### 子 Logger 创建规则

| 配置 | 创建哪些子 Logger |
|------|------------------|
| `cfg.Level = DEBUG` | all + DEBUG + INFO + WARN + ERROR + FATAL + PANIC（7个） |
| `cfg.Level = INFO` | all + INFO + WARN + ERROR + FATAL + PANIC（6个） |
| `cfg.Level = WARN` | all + WARN + ERROR + FATAL + PANIC（5个） |

**关键点**：`cfg.Level` 决定了**从哪个级别开始**创建专属文件，低于该级别的日志只写入全量文件。

---

## 三、工作流程详解

### 3.1 构造函数流程

```
NewLevelLogger(cfg)
    │
    ├── 校验: cfg != nil, cfg.LogPath != ""
    │
    ├── 提取目录: dir = filepath.Dir(cfg.LogPath)
    │   例: "logs/app.log" → "logs/"
    │
    ├── 确定基准级别: level = cfg.Level (默认为 INFO)
    │
    ├── 创建 all Logger
    │   └── 路径: cfg.LogPath（用户指定的全量日志路径）
    │   └── OutputConsole = false（级别路由模式默认不输出终端）
    │
    └── 循环创建级别专属 Logger
        for lvl in [DEBUG, INFO, WARN, ERROR, FATAL, PANIC]:
            if lvl >= level:
                创建 Logger，路径: dir/lvl.String()+".log"
                例: "logs/INFO.log", "logs/ERROR.log"
                关闭采样、关闭异步（简化处理）
                存入 levels[lvl]
```

### 3.2 日志记录流程

以 `ll.Info("用户登录")` 为例：

```
LevelLogger.Info("用户登录")
    │
    └── 调用内部 ll.log(INFO, "用户登录", nil)
        │
        ├── 级别检查: INFO >= ll.level?
        │   └── 否 → 直接返回（日志被过滤）
        │
        ├── 写入全量: ll.all.log(INFO, "用户登录", nil)
        │   └── 写入 cfg.LogPath（如 logs/app.log）
        │
        └── 写入专属: 查找 ll.levels[INFO]
            └── 存在 → ll.levels[INFO].log(INFO, "用户登录", nil)
                └── 写入 logs/INFO.log
            └── 不存在 → 跳过（该级别没有专属文件）
```

再以 `ll.Error("连接失败")` 为例：

```
LevelLogger.Error("连接失败")
    │
    └── ll.log(ERROR, "连接失败", nil)
        │
        ├── 级别检查: ERROR >= ll.level? → 通过
        │
        ├── 写入全量: ll.all.log(ERROR, "连接失败", nil)
        │   └── 写入 logs/app.log
        │
        └── 写入专属: ll.levels[ERROR].log(ERROR, "连接失败", nil)
            └── 写入 logs/ERROR.log
```

### 3.3 数据流图示

```
┌─────────────────┐
│   用户代码       │
│  ll.Info("msg") │
└────────┬────────┘
         ▼
┌─────────────────┐
│  LevelLogger    │
│  1. 级别检查     │───级别不够──→ 丢弃
│  2. 路由分发     │
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌───────┐  ┌────────┐
│  all  │  │ levels │
│Logger │  │[INFO]  │
└───┬───┘  └────┬───┘
    │           │
    ▼           ▼
┌────────┐  ┌────────┐
│app.log │  │INFO.log│
│(全量)   │  │(专属)  │
└────────┘  └────────┘
```

---

## 四、关键设计要点

### 4.1 为什么直接调用 `log()` 而不是 `Info()`？

LevelLogger 内部直接调用子 Logger 的未导出方法 `log()`，而不是公开方法 `Info()`：

```go
// 正确做法（当前设计）
func (ll *LevelLogger) Info(msg string) {
    ll.log(INFO, msg, nil)  // 调用内部 log 方法
}

// 错误做法（会导致问题）
func (ll *LevelLogger) Info(msg string) {
    ll.all.Info(msg)  // 调用公开方法
    ll.levels[INFO].Info(msg)
}
```

**原因**：
- `Logger.Info()` 内部会再次检查级别、获取调用者信息、处理采样等
- 如果 LevelLogger 调用 `all.Info()`，调用栈会多一层，导致 `callerSkip` 计算错误
- 直接调用 `log()` 可以复用 LevelLogger 已经完成的级别检查和字段准备

### 4.2 Fatal/Panic 特殊处理

Fatal 和 Panic 不能简单委托给子 Logger：

```go
func (ll *LevelLogger) Fatal(msg string) {
    ll.log(FATAL, msg, nil)  // 写入 all + FATAL
    _ = ll.Sync()             // 确保刷盘
    os.Exit(1)                // 只退出一次！
}
```

**为什么不能调 `all.Fatal()`？**
- `Logger.Fatal()` 内部会调 `os.Exit(1)`
- 如果先调 `all.Fatal()`，程序直接退出，`levels[FATAL]` 根本没机会写
- 所以 LevelLogger 自己调 `os.Exit(1)`，确保所有子 Logger 都写入后再退出

### 4.3 动态级别调整

```go
func (ll *LevelLogger) SetLevel(level Level) {
    ll.level.Store(int32(level))  // 原子操作，线程安全
}
```

- 使用 `atomic.Int32` 实现无锁级别切换
- 调整的是 LevelLogger 的级别，**不影响**子 Logger 的级别
- 子 Logger 的级别始终等于创建时的 `cfg.Level`

---

## 五、文件生成示例

### 场景 1：开发环境

```go
ll := fastlog.NewLevelLogger(fastlog.Dev("logs/app.log"))
// cfg.Level = DEBUG
```

生成文件：
```
logs/
├── app.log      ← 全量日志（用户指定）
├── DEBUG.log    ← DEBUG 专属
├── INFO.log     ← INFO 专属
├── WARN.log     ← WARN 专属
├── ERROR.log    ← ERROR 专属
├── FATAL.log    ← FATAL 专属
└── PANIC.log    ← PANIC 专属
```

### 场景 2：生产环境

```go
ll := fastlog.NewLevelLogger(fastlog.Prod("logs/app.log"))
// cfg.Level = WARN
```

生成文件：
```
logs/
├── app.log      ← 全量日志（WARN+）
├── WARN.log     ← WARN 专属
├── ERROR.log    ← ERROR 专属
├── FATAL.log    ← FATAL 专属
└── PANIC.log    ← PANIC 专属
// DEBUG.log 和 INFO.log 不会创建！
```

---

## 六、边界情况处理

| 场景 | 处理方式 |
|------|---------|
| cfg 为 nil | panic |
| cfg.LogPath 为空 | panic（无法推导目录）|
| cfg.LogPath = "logs/INFO.log" | panic（与级别文件冲突）|
| 子 Logger 写入失败 | 输出到 stderr（与 Logger 行为一致）|
| Sync/Close 部分失败 | `errors.Join` 合并所有错误返回 |
| 并发写入 | 各子 Logger 独立加锁，互不干扰 |

---

## 七、与普通 Logger 的关系

```
┌─────────────────────────────────────┐
│           LevelLogger               │
│  （代理层：负责路由、级别控制）        │
└─────────────────────────────────────┘
              │
    ┌─────────┼─────────┐
    ▼         ▼         ▼
┌───────┐ ┌───────┐ ┌───────┐
│  all  │ │levels │ │levels │
│Logger │ │[INFO] │ │[ERROR]│
└───┬───┘ └───┬───┘ └───┬───┘
    │         │         │
    ▼         ▼         ▼
┌───────┐ ┌───────┐ ┌───────┐
│app.log│ │INFO.log│ │ERROR.log│
└───────┘ └───────┘ └───────┘
```

**LevelLogger 是对 Logger 的包装**，本身不处理文件写入、格式化等细节，完全复用 Logger 的能力。
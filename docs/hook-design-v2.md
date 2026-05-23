# FastLog Hook 机制设计方案 V2（内部机制版）

> 设计目标：Hook 作为内部机制，用户只需一个开关字段启用级别路由

---

## 一、设计概述

### 1.1 核心思想

**Hook 完全内部化**，用户无感知。通过 Config 的一个字段控制是否启用级别路由。

```
用户配置：cfg.LevelRouter = true
    ↓
内部自动创建：为 >= cfg.Level 的每个级别创建 Hook
    ↓
日志输出时：主文件 + 级别专属文件同时写入
```

### 1.2 与 V1 方案对比

| 特性 | V1（暴露 Hook） | V2（内部 Hook） |
|------|-----------------|-----------------|
| 用户接口 | `AddHook()` | `cfg.LevelRouter = true` |
| 使用复杂度 | 需手动添加 Hook | 一键开启 |
| 灵活性 | 高（自由组合） | 中（固定模式） |
| 代码复杂度 | 中 | 低 |
| 适用场景 | 高级用户 | 普通用户 |

---

## 二、Config 配置

### 2.1 新增字段

```go
// config.go

type Config struct {
    // ... 现有字段
    
    // LevelRouter 启用级别路由
    // 为 true 时，自动在 LogPath 同级目录创建 {LEVEL}.log 文件
    // 例：LogPath="logs/app.log" → 创建 logs/DEBUG.log, logs/INFO.log 等
    LevelRouter bool
}
```

### 2.2 使用示例

```go
// 基础用法（一行开启）
cfg := fastlog.NewConfig("logs/app.log")
cfg.Level = fastlog.DEBUG
cfg.LevelRouter = true  // 启用级别路由

logger := fastlog.New(cfg)
defer logger.Close()

// 效果：
// logger.Debug("...") → logs/app.log + logs/DEBUG.log
// logger.Info("...")  → logs/app.log + logs/INFO.log
// logger.Error("...") → logs/app.log + logs/ERROR.log
```

---

## 三、内部实现

### 3.1 Hook 接口（内部使用，不暴露）

```go
// hook.go（内部文件）
package fastlog

// hook 日志钩子接口（内部使用，小写不导出）
type hook interface {
    Fire(entry *Entry, data []byte) error
    Levels() []Level
    Close() error
}

// levelHook 按级别分发的钩子（内部使用）
type levelHook struct {
    level  Level
    writer io.WriteCloser
}

func (h *levelHook) Fire(entry *Entry, data []byte) error {
    if entry.Level != h.level {
        return nil
    }
    _, err := h.writer.Write(data)
    return err
}

func (h *levelHook) Levels() []Level {
    return []Level{h.level}
}

func (h *levelHook) Close() error {
    if h.writer != nil {
        return h.writer.Close()
    }
    return nil
}
```

### 3.2 Logger 结构体改造

```go
// logger.go

type Logger struct {
    config  *Config
    writer  io.WriteCloser
    sampler *Sampler
    mu      sync.Mutex
    level   atomic.Int32
    hooks   []hook  // 内部使用，小写不导出
}
```

### 3.3 New 函数改造

```go
func New(cfg *Config) *Logger {
    // ... 原有逻辑：验证、克隆、默认值
    
    // 创建主写入器
    writer := config.NewWriter()
    if writer == nil {
        writer = &ConsoleWriter{w: os.Stdout}
    }
    
    l := &Logger{
        config: config,
        writer: writer,
        // ... 其他字段
    }
    
    // 如果启用 LevelRouter，自动创建 Hooks
    if cfg.LevelRouter && cfg.LogPath != "" {
        l.initLevelHooks(cfg)
    }
    
    return l
}

// initLevelHooks 初始化级别路由 Hooks（内部方法）
func (l *Logger) initLevelHooks(cfg *Config) {
    dir := filepath.Dir(cfg.LogPath)
    
    // 为 >= cfg.Level 的每个级别创建专属文件
    for _, lvl := range AllLevels() {
        if lvl < cfg.Level {
            continue
        }
        
        // 创建级别专属配置
        lvlCfg := cfg.Clone()
        lvlCfg.LogPath = filepath.Join(dir, lvl.String()+".log")
        lvlCfg.OutputConsole = false
        lvlCfg.LevelRouter = false  // 防止递归
        
        // 创建写入器
        writer := lvlCfg.NewWriter()
        if writer != nil {
            l.hooks = append(l.hooks, &levelHook{
                level:  lvl,
                writer: writer,
            })
        }
    }
}
```

### 3.4 log 方法改造

```go
func (l *Logger) log(level Level, msg string, fields []Field) {
    // ... 原有逻辑：检查级别、采样、构建 Entry
    
    // 记录调用者（正确获取）
    if l.config.Caller {
        entry.Caller = getCaller(callerSkip)
    }
    
    // 格式化
    data, err := l.config.Formatter.Format(entry)
    if err != nil {
        return
    }
    
    // 写入主文件 + 执行内部 Hooks
    l.mu.Lock()
    _, _ = l.writer.Write(data)
    
    // 执行内部 Hooks（级别路由）
    for _, h := range l.hooks {
        for _, lvl := range h.Levels() {
            if lvl == level {
                _ = h.Fire(entry, data)
                break
            }
        }
    }
    l.mu.Unlock()
}
```

### 3.5 Close 方法改造

```go
func (l *Logger) Close() error {
    var errs []error
    
    // 关闭主写入器
    if err := l.writer.Close(); err != nil {
        errs = append(errs, err)
    }
    
    // 关闭所有 Hooks
    for _, h := range l.hooks {
        if err := h.Close(); err != nil {
            errs = append(errs, err)
        }
    }
    
    return joinErrors(errs)
}
```

---

## 四、使用示例

### 4.1 基础用法

```go
package main

import (
    "gitee.com/MM-Q/fastlog"
)

func main() {
    // 创建配置，启用级别路由
    cfg := fastlog.NewConfig("logs/app.log")
    cfg.Level = fastlog.INFO
    cfg.LevelRouter = true  // 一键开启
    
    logger := fastlog.New(cfg)
    defer logger.Close()
    
    // 正常使用
    logger.Debug("debug msg")  // 被过滤（Level=INFO）
    logger.Info("info msg")    // → app.log + INFO.log
    logger.Error("error msg")  // → app.log + ERROR.log
}
```

### 4.2 与场景配置结合

```go
// 开发环境（自动启用级别路由）
cfg := fastlog.Dev("logs/app.log")
cfg.LevelRouter = true

// 生产环境（只记录 ERROR+）
cfg := fastlog.Prod("logs/app.log")
cfg.LevelRouter = true  // 只创建 ERROR.log, FATAL.log, PANIC.log
```

### 4.3 与现有代码兼容

```go
// 原有代码（不变）
cfg := fastlog.NewConfig("logs/app.log")
logger := fastlog.New(cfg)

// 新增功能（加一行）
cfg.LevelRouter = true
```

---

## 五、文件变更清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `config.go` | 修改 | 添加 `LevelRouter bool` 字段 |
| `logger.go` | 修改 | 添加 `hooks []hook` 字段，改造 `New()`/`log()`/`Close()` |
| `hook.go` | 新增 | 内部 Hook 接口和 levelHook 实现（约 40 行，不导出） |
| `level_logger.go` | 删除 | 被替代 |
| `level_logger_test.go` | 删除 | 被替代 |

**总代码量**：约 80 行（比 LevelLogger 减少 70%）

---

## 六、优势总结

### 6.1 极致简洁

| 操作 | LevelLogger | V2 方案 |
|------|-------------|---------|
| 启用级别路由 | `NewLevelLogger(cfg)` | `cfg.LevelRouter = true` |
| 学习成本 | 需了解新类型 | 只需了解新字段 |
| API 数量 | 2 个（Logger + LevelLogger） | 1 个（Logger） |

### 6.2 完全兼容

- 不改动现有 Logger API
- 不使用时行为完全一致
- 零学习成本升级

### 6.3 Caller 正确

- 在 Logger 内部获取 Caller
- 指向用户代码，无转发层问题

### 6.4 自动管理

- 自动创建级别文件
- 自动关闭所有资源
- 用户无感知

---

## 七、注意事项

### 7.1 路径冲突

如果 `cfg.LogPath` 与级别文件路径冲突（如 `logs/INFO.log`），应该：
- 方案 A：panic（快速失败）
- 方案 B：自动调整（如改为 `logs/INFO_1.log`）

**推荐方案 A**，明确报错。

### 7.2 资源占用

启用 `LevelRouter` 后会打开多个文件：
- 1 个主文件
- N 个级别文件（N = 6 - cfg.Level + 1）

**建议**：生产环境如果 `Level=ERROR`，只打开 4 个文件（app.log + ERROR + FATAL + PANIC），可接受。

### 7.3 性能影响

- 无 Hook 时：零开销
- 有 Hook 时：顺序写入，锁内执行

**建议**：高并发场景谨慎使用，或配合 Async 模式。

---

## 八、实施步骤

1. **修改 `config.go`**：添加 `LevelRouter bool` 字段
2. **创建 `hook.go`**：内部 Hook 接口和实现（不导出）
3. **修改 `logger.go`**：
   - 添加 `hooks` 字段
   - 改造 `New()` 添加 `initLevelHooks()` 调用
   - 改造 `log()` 添加 Hooks 执行
   - 改造 `Close()` 添加 Hooks 关闭
4. **删除 `level_logger.go` 和 `level_logger_test.go`**
5. **创建 `logger_levelrouter_test.go`**：测试级别路由功能
6. **运行测试**：确保所有测试通过

---

## 九、示例代码对比

### 使用前（LevelLogger）

```go
// 需要引入新类型
cfg := fastlog.NewConfig("logs/app.log")
logger := fastlog.NewLevelLogger(cfg)  // 新类型
defer logger.Close()

logger.Info("msg")
```

### 使用后（V2 方案）

```go
// 只需加一个字段
cfg := fastlog.NewConfig("logs/app.log")
cfg.LevelRouter = true  // 一行开启

logger := fastlog.New(cfg)  // 原类型
defer logger.Close()

logger.Info("msg")
```

---

> 设计方案 V2 完成，更简洁、更易用、完全内部化。

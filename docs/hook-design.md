# FastLog Hook 机制设计方案

> 设计目标：用 Hook 机制替代 LevelLogger，实现更简洁的级别路由功能

---

## 一、设计概述

### 1.1 核心思想

在现有 Logger 基础上增加 **Hook 机制**，允许用户在日志输出时执行额外的写入操作。

```
用户调用 logger.Error("msg")
    ↓
Logger.log() 内部
    ├── 1. 构建 Entry（Caller 正确指向用户代码）
    ├── 2. 格式化日志
    ├── 3. 写入主文件（原有逻辑，不变）
    └── 4. 执行 Hooks（新增）
            └── LevelHook.Fire(entry, data)
                    └── 级别匹配 → 写入专属文件
```

### 1.2 与 LevelLogger 对比

| 特性 | LevelLogger | Hook 方案 |
|------|-------------|-----------|
| 实现方式 | 包装多个 Logger | 在 Logger 内部增加 Hooks |
| Caller 信息 | ❌ 指向 level_logger.go | ✅ 指向用户代码 |
| 代码复杂度 | 高（240 行，独立文件） | 低（约 70 行，分散到现有文件） |
| 使用方式 | `NewLevelLogger(cfg)` | `New(cfg) + AddHook()` |
| 灵活性 | 固定模式 | 自由组合多个 Hook |
| 向后兼容 | 需维护两个 API | 完全兼容，Hooks 可选 |

---

## 二、接口设计

### 2.1 Hook 接口

```go
// hook.go
package fastlog

// Hook 日志钩子接口
// 用于在日志输出时执行额外操作
type Hook interface {
    // Fire 日志触发时调用
    // 参数:
    //   - entry: 日志条目，包含完整信息
    //   - data: 格式化后的日志数据
    // 返回:
    //   - error: 执行过程中的错误
    Fire(entry *Entry, data []byte) error
    
    // Levels 返回关心的日志级别列表
    // 只有这些级别的日志会触发 Fire
    Levels() []Level
}
```

### 2.2 LevelHook 实现

```go
// level_hook.go
package fastlog

import "io"

// LevelHook 按级别分发日志的钩子
type LevelHook struct {
    level  Level          // 关心的级别
    writer io.WriteCloser // 写入目标
}

// NewLevelHook 创建级别钩子
func NewLevelHook(level Level, writer io.WriteCloser) *LevelHook {
    return &LevelHook{
        level:  level,
        writer: writer,
    }
}

// Fire 执行钩子
func (h *LevelHook) Fire(entry *Entry, data []byte) error {
    if entry.Level != h.level {
        return nil
    }
    _, err := h.writer.Write(data)
    return err
}

// Levels 返回关心的级别
func (h *LevelHook) Levels() []Level {
    return []Level{h.level}
}

// Close 关闭写入器
func (h *LevelHook) Close() error {
    if h.writer != nil {
        return h.writer.Close()
    }
    return nil
}
```

### 2.3 便捷函数

```go
// NewLevelFileHook 快速创建文件级别钩子
func NewLevelFileHook(level Level, path string) (*LevelHook, error) {
    cfg := NewConfig(path)
    writer := cfg.NewWriter()
    if writer == nil {
        return nil, fmt.Errorf("failed to create writer for %s", path)
    }
    return NewLevelHook(level, writer), nil
}
```

---

## 三、Logger 改造

### 3.1 结构体修改

```go
// logger.go

type Logger struct {
    config  *Config
    writer  io.WriteCloser
    sampler *Sampler
    mu      sync.Mutex
    level   atomic.Int32
    hooks   []Hook  // 新增：钩子列表
}
```

### 3.2 新增方法

```go
// AddHook 添加日志钩子
// 线程安全，可在运行时添加
func (l *Logger) AddHook(hook Hook) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.hooks = append(l.hooks, hook)
}
```

### 3.3 log 方法改造

```go
func (l *Logger) log(level Level, msg string, fields []Field) {
    // ... 原有逻辑：检查级别、采样、构建 Entry
    
    // 记录调用者（正确获取，指向用户代码）
    if l.config.Caller {
        entry.Caller = getCaller(callerSkip)
    }
    
    // 格式化
    data, err := l.config.Formatter.Format(entry)
    if err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "format error: %v\n", err)
        return
    }
    
    // 写入主文件 + 执行 Hooks
    l.mu.Lock()
    _, err = l.writer.Write(data)
    if err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "write error: %v\n", err)
    }
    
    // 执行 Hooks（新增逻辑）
    for _, hook := range l.hooks {
        for _, lvl := range hook.Levels() {
            if lvl == level {
                _ = hook.Fire(entry, data)  // 忽略错误，避免影响主流程
                break
            }
        }
    }
    l.mu.Unlock()
}
```

### 3.4 Close 方法改造

```go
func (l *Logger) Close() error {
    var errs []error
    
    // 关闭主写入器
    if err := l.writer.Close(); err != nil {
        errs = append(errs, err)
    }
    
    // 关闭所有 Hook 的写入器
    for _, hook := range l.hooks {
        if closer, ok := hook.(interface{ Close() error }); ok {
            if err := closer.Close(); err != nil {
                errs = append(errs, err)
            }
        }
    }
    
    return joinErrors(errs)
}
```

---

## 四、使用示例

### 4.1 基础用法：ERROR 专属文件

```go
package main

import (
    "gitee.com/MM-Q/fastlog"
)

func main() {
    // 1. 创建主 Logger（全量日志）
    cfg := fastlog.NewConfig("logs/app.log")
    cfg.Level = fastlog.DEBUG
    logger := fastlog.New(cfg)
    defer logger.Close()
    
    // 2. 添加 ERROR 专属钩子
    errorHook, err := fastlog.NewLevelFileHook(fastlog.ERROR, "logs/ERROR.log")
    if err != nil {
        panic(err)
    }
    logger.AddHook(errorHook)
    
    // 3. 正常使用
    logger.Info("服务启动成功")      // 只写入 app.log
    logger.Error("数据库连接失败")   // 写入 app.log + ERROR.log
}
```

### 4.2 多级别专属文件

```go
func main() {
    logger := fastlog.New(fastlog.NewConfig("logs/app.log"))
    defer logger.Close()
    
    // 为每个级别创建专属文件
    levels := []fastlog.Level{
        fastlog.DEBUG,
        fastlog.INFO,
        fastlog.WARN,
        fastlog.ERROR,
    }
    
    for _, lvl := range levels {
        path := fmt.Sprintf("logs/%s.log", lvl.String())
        hook, err := fastlog.NewLevelFileHook(lvl, path)
        if err != nil {
            continue
        }
        logger.AddHook(hook)
    }
    
    // 现在：
    // logger.Debug("...") → app.log + DEBUG.log
    // logger.Info("...")  → app.log + INFO.log
    // logger.Error("...") → app.log + ERROR.log
}
```

### 4.3 自定义 Hook（扩展功能）

```go
// AlertHook 错误告警钩子
type AlertHook struct {
    webhook string
}

func (h *AlertHook) Fire(entry *Entry, data []byte) error {
    if entry.Level >= ERROR {
        // 发送告警通知
        sendAlert(h.webhook, string(data))
    }
    return nil
}

func (h *AlertHook) Levels() []Level {
    return []Level{ERROR, FATAL, PANIC}
}

// 使用
logger.AddHook(&AlertHook{webhook: "https://alert.example.com"})
```

---

## 五、文件变更清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `hook.go` | 新增 | Hook 接口定义（约 15 行） |
| `level_hook.go` | 新增 | LevelHook 实现（约 50 行） |
| `logger.go` | 修改 | 添加 hooks 字段、AddHook 方法、改造 log/Close（约 30 行） |
| `level_logger.go` | 删除 | 被 Hook 方案替代 |
| `level_logger_test.go` | 删除 | 被新测试替代 |

**总代码量**：约 100 行（比 LevelLogger 的 240 行减少 60%）

---

## 六、优势总结

### 6.1 相比 LevelLogger

1. **Caller 正确**：在 Logger 内部获取，无转发层
2. **代码简洁**：100 行 vs 240 行
3. **维护简单**：分散到现有文件，不新增独立模块
4. **灵活性高**：用户可自由组合多个 Hook

### 6.2 向后兼容

- 不使用 Hook 时，行为和原来完全一致
- 现有代码无需修改
- 新增功能完全可选

### 6.3 可扩展性

除了级别路由，Hook 还可用于：
- 错误告警（ERROR 时发送通知）
- 日志统计（记录日志数量）
- 日志过滤（丢弃特定日志）
- 日志转换（修改日志内容）

---

## 七、实施步骤

1. **创建 `hook.go`**：定义 Hook 接口
2. **创建 `level_hook.go`**：实现 LevelHook
3. **修改 `logger.go`**：添加 hooks 支持和执行逻辑
4. **删除 `level_logger.go` 和 `level_logger_test.go`**
5. **创建 `hook_test.go`**：编写单元测试
6. **运行测试**：确保所有测试通过

---

## 八、风险评估

| 风险 | 等级 | 缓解措施 |
|------|------|----------|
| Hook 执行影响性能 | 低 | 无 Hook 时零开销，有 Hook 时顺序执行 |
| Hook 错误影响主流程 | 低 | 忽略 Hook 错误，只记录主文件错误 |
| 并发安全问题 | 低 | 在 mu.Lock 保护下执行 Hooks |

---

> 设计方案完成，等待评审确认后实施。

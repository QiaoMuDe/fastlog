# 动态级别调整方案

> 基于 `atomic.Int32` 实现运行时无锁切换日志级别
> 
> **更新**: `Level` 类型从 `int8` 改为 `int32`，与 `atomic.Int32` 完美配合，避免类型转换

---

## 一、目标

运行时无需重启即可修改日志级别，类似 `zap.AtomicLevel` / `slog.LevelVar`。

```go
logger.SetLevel(fastlog.DEBUG)  // 立即生效，所有后续写入使用新级别
```

---

## 二、总体设计

### 2.1 核心思路

把级别检查从 `Config.Level`（只读）改为 `Logger.level`（atomic.Int32，可写），在 `log()` 的检查路径上增加一次 `atomic.Load`。

```
改前: Logger.log() → l.config.Level.Enabled(level)  → 只读，永不改变
改后: Logger.log() → Level(l.level.Load()).Enabled(level) → 可读可写，SetLevel 修改
```

### 2.2 关键设计决策：`Level` 改为 `int32`

| 原设计 | 新设计 | 原因 |
|--------|--------|------|
| `type Level int8` | `type Level int32` | 与 `atomic.Int32` 类型一致，避免转换 |
| `l.level.Store(int64(level))` | `l.level.Store(int32(level))` | 无需转换，直接存储 |
| `Level(l.level.Load())` | `Level(l.level.Load())` | 名义转换（int32→int32），零开销 |

**优势**：
- 零类型转换开销（`int32` ↔ `int32`）
- 代码更清晰，无 `int64` 中间层
- 内存占用相同（都是 4 字节对齐）

### 2.3 涉及的文件和改动量

| 文件 | 改动 | 行数 |
|------|------|------|
| `logger.go` | `Level` 类型 `int8` → `int32` | ±1 |
| `logger.go` | Logger 结构体新增 `level atomic.Int32` | +1 |
| `logger.go` | `log()` 中一行判断改动 | ±1 |
| `logger.go` | 新增 `SetLevel()` / `Level()` 2 个方法 | +6 |
| `logger.go` | `New()` 结尾加一行 `l.level.Store(...)` | +1 |
| `logger.go` | import 加 `sync/atomic` | +1 |
| **总计** | | **~11 行** |

---

## 三、具体变更

### 3.1 Level 类型定义（修改）

```go
// Level 表示日志级别 (int32 方便与 atomic.Int32 配合)
type Level int32
```

### 3.2 Logger 结构体（修改）

```go
type Logger struct {
	config  *Config        // 日志配置（克隆后的，不可变）
	writer  io.WriteCloser // 日志写入器
	sampler *Sampler       // 日志采样器, nil 表示不启用采样
	mu      sync.Mutex     // 日志记录器的互斥锁
	level   atomic.Int32   // 运行时日志级别, 支持动态调整, 初始化时从 config.Level 设置
}
```

`atomic.Int32` 的选型理由：
- `Level` 是 `int32`（4 字节），`atomic.Int32` 完美匹配
- x86 上 `atomic.Load` = 普通 `MOV`，零额外开销（x86 TSO 保证）
- `atomic.Int32.Store` 编译为 `XCHG`，比普通写慢 ~20ns——但 `SetLevel` 调用频率极低，可忽略

### 3.3 New() 中初始化（修改）

```go
func New(cfg *Config) *Logger {
	// ... 现有逻辑不变 ...

	l := &Logger{
		config:  config,
		writer:  writer,
		sampler: sampler,
	}

	// 以 Config.Level 作为运行时级别的初始值
	l.level.Store(int32(config.Level))

	return l
}
```

### 3.4 log() 中的检查（修改）

```go
func (l *Logger) log(level Level, msg string, fields []Field) {
	// 改前: if !l.config.Level.Enabled(level) {
	// 改后:
	if !Level(l.level.Load()).Enabled(level) {
		return
	}
	// ... 后续不变 ...
}
```

**注意**：`Level(l.level.Load())` 是名义类型转换（`int32` → `Level`，底层都是 `int32`），编译后无实际指令。

### 3.5 新增方法（新增）

```go
// SetLevel 运行时动态修改日志级别, 立即生效
//
// 参数:
//   - level: 新的日志级别
func (l *Logger) SetLevel(level Level) {
	l.level.Store(int32(level))
}

// Level 返回当前运行时日志级别
//
// 返回:
//   - Level: 当前日志级别
func (l *Logger) Level() Level {
	return Level(l.level.Load())
}
```

---

## 四、使用示例

### 4.1 基本用法

```go
logger := fastlog.New(fastlog.Prod("app.log"))

// 生产环境初始 WARN，只记录警告及以上
logger.Info("这条不会输出")       // 被压制
logger.Warn("这条会输出")         // 输出

// 出问题时，运行时调整为 DEBUG
logger.SetLevel(fastlog.DEBUG)

// 后续写入立即使用新级别，无需重启
logger.Info("现在可以看到 INFO 了")   // 输出
logger.Debug("DEBUG 也能看到")       // 输出
```

### 4.2 HTTP 热更新（基于此方案的扩展思路）

```go
// HTTP handler: 读取当前级别
http.HandleFunc("/log/level", func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, _ = fmt.Fprintf(w, "%s", logger.Level())
	case http.MethodPut:
		body, _ := io.ReadAll(r.Body)
		level, err := fastlog.ParseLevel(string(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.SetLevel(level)  // 立即生效
		_, _ = fmt.Fprintf(w, "level set to %s", level)
	}
})
```

### 4.3 持久化到文件（基于此方案的扩展思路）

```go
// 启动时从文件恢复
func (l *Logger) initLevel(defaultLevel Level) {
	data, err := os.ReadFile(l.levelFilePath())
	if err != nil {
		l.level.Store(int32(defaultLevel))
		return
	}
	level, err := ParseLevel(strings.TrimSpace(string(data)))
	if err != nil {
		l.level.Store(int32(defaultLevel))
		return
	}
	l.level.Store(int32(level))
}

// SetLevel + 持久化
func (l *Logger) SetLevel(level Level) {
	l.level.Store(int32(level))
	_ = os.WriteFile(l.levelFilePath(), []byte(level.String()), 0644)
}

// 级别持久化文件路径: 日志文件同目录下的 .level 文件
func (l *Logger) levelFilePath() string {
	return l.config.LogPath + ".level"
}
```

---

## 五、性能分析

### 5.1 对比

| 操作 | 当前（读 Config） | 改后（atomic.Load） | 差异 |
|------|-----------------|--------------------|------|
| 级别检查 | `MOV` 1 条指令 | `MOV` 1 条指令 | **0** |
| CPU 栅栏 | 无 | 无（x86 TSO） | **0** |
| L1 延迟 | ~1ns | ~1ns | **0** |
| 写路径 | 不支持 | `XCHG` ~25ns | +25ns（但极少调用） |

### 5.2 类型转换开销对比

| 场景 | 原方案 (Level=int8) | 新方案 (Level=int32) | 差异 |
|------|--------------------|---------------------|------|
| Store | `int8` → `int64` → `int32` | `int32` → `int32` | **省 1 次转换** |
| Load | `int32` → `int64` → `int8` | `int32` → `int32` | **省 1 次转换** |
| 实际指令 | 2 条 MOV | 0 条（名义转换） | **更优** |

### 5.3 在整条日志路径中的占比

```
log() 完整路径:
  Level.Enabled()     ← ~1ns  (atomic.Load)
  sampler.Allow()     ← ~20ns (哈希+atomic)
  GetEntry()          ← ~10ns (sync.Pool)
  time.Now()          ← ~100ns
  Format()            ← ~500ns(字符串拼接/JSON序列化)
  mu.Lock()           ← ~30ns (无竞争)
  writer.Write()      ← ~1μs+ (系统调用)
```

级别检查占总路径占比：**< 0.1%**

### 5.4 ARM64/MIPS 等其他架构

在非 x86 架构上，`atomic.Load` 可能需要 `DMB`（Data Memory Barrier）指令，比普通读多 ~5ns。但相对于整条日志路径的微秒级耗时，仍然可以忽略。

---

## 六、与 zap / slog 的设计对比

| 库 | 方案 | 级别后门 | 持久化 |
|----|------|---------|--------|
| **zap** | `AtomicLevel`（atomic 包装器） | 无 | 支持 JSON 序列化 |
| **slog** | `LevelVar`（atomic 包装器） | `slog.SetDefault` | 无 |
| **FastLog（本方案）** | `Logger.level atomic.Int32` | `logger.SetLevel()` | 可在此基础上扩展 |

和 zap/slog 的核心区别：
- **zap/slog 的级别包装器是独立类型**，可以在多个 Logger 间共享
- **本方案把级别直接放在 Logger 上**，每个 Logger 独立，API 更简洁

如果需要"多 Logger 共享级别"或"从外部控制级别"，可以在不破坏本方案的情况下扩展：

```go
// 后续扩展：共享级别提供者
type LevelProvider interface {
	Level() Level
	SetLevel(Level)
}

// Logger 可以接受外部提供者
func (l *Logger) SetLevelProvider(p LevelProvider) { ... }
```

---

## 七、不涉及的内容

本方案**不包含**以下功能（后续按需实现）：

| 功能 | 说明 |
|------|------|
| HTTP API | 对外暴露 HTTP 端点修改级别（在 example 中实现） |
| 持久化到文件 | 将当前级别存到 `.level` 文件，重启后恢复（在 example 中实现） |
| 多 Logger 共享级别 | `LevelProvider` 接口（如需再添加） |
| 级别变更回调 | 级别变化时触发钩子（极少需要） |

---

## 八、变更总结

| 项目 | 原状态 | 新状态 | 影响 |
|------|--------|--------|------|
| `Level` 类型 | `int8` | `int32` | 内存对齐从 1 字节→4 字节，无实际影响 |
| Logger.level | 无 | `atomic.Int32` | 新增字段 |
| SetLevel() | 无 | 有 | 新增方法 |
| Level() | 无 | 有 | 新增方法 |
| 类型转换 | `int8` ↔ `int64` | 无（都是 `int32`） | 更简洁 |
| 性能 | 基准 | 相同（x86） | 零损失 |

---

> **结论**: 约 11 行代码的极小改动，`Level` 改为 `int32` 后与 `atomic.Int32` 完美配合，零运行时性能损失（x86 上），即可获得运行时动态调整级别能力。

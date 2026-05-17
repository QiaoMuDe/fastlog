# fastlog 日志采样功能设计方案

## 背景

高并发下，同一位置的相同日志可能在短时间内重复输出数千次，导致磁盘 I/O 飙升、日志存储暴增、监控告警轰炸。采样功能通过**时间窗口 + 内容去重**，确保相同日志在窗口内按规则控制输出量，被抑制的日志只计数不写入。

## 效果示例

```
// 采样配置：每 10 秒，前 3 条放行，之后每 10 条放行 1 条
ERROR | login.go:Login:42 - 数据库连接超时     ← 第1条  放行
ERROR | login.go:Login:42 - 数据库连接超时     ← 第2条  放行
ERROR | login.go:Login:42 - 数据库连接超时     ← 第3条  放行
...（第4~12条被抑制，不写入任何内容）...
ERROR | login.go:Login:42 - 数据库连接超时     ← 第13条 放行（3 + 1*10）
```

## 设计决策

### Key 的选取

用 `level + message` 作为判重依据：

- **level**: 同一消息不同级别视为不同日志，各自独立计数
- **message**: 日志消息内容，通过 FNV 哈希映射到固定桶

**不包含 caller**：caller 信息（文件名:行号）虽然能精确区分调用位置，但实践中同一类日志（如"数据库连接超时"）出现在不同位置时，我们希望合并计数而非各自独立采样——合并才能更有效地削峰。zap 也只用了 `level + message`。

**不包含 fields**：fields 经常携带动态值（user_id、req_id、latency），几乎每条都不同，包含进去会导致采样完全失效。

### 数据结构：固定桶 + atomic 计数器

使用固定大小的二维数组代替 map，消除锁竞争：

```
counters[6][4096]samplerCounter
  │       │
  │       └── message 的 FNV 哈希 % 4096，决定使用哪个桶
  └── level (DEBUG=1 ~ PANIC=6)
```

- 4096 个桶足以将不同消息均匀分布，哈希冲突概率极低
- 每个桶是一个 `samplerCounter`，内含 `atomic.Int64` 和 `atomic.Uint64`，全原子操作
- 零 GC 压力、零锁竞争

### 集成位置

在 `log()` 方法的 `Formatter.Format()` 之前插入采样检查，此时还能拿到结构化的 `level` 和 `msg` 原始值。

## 接口设计

### 使用方式

```go
// 创建采样器：每 10 秒一个窗口，前 3 条放行，之后每 10 条放行 1 条
sampler := fastlog.NewSampler(10*time.Second, 3, 10)

logger := fastlog.New(
    fastlog.WithLevel(fastlog.DEBUG),
    fastlog.WithSampler(sampler),
)

// 即使被调用 10000 次，每秒最多输出数条
for i := 0; i < 10000; i++ {
    logger.Errorw("数据库连接超时", fastlog.String("db", "mysql"))
}
```

### 参数说明

`NewSampler(tick time.Duration, initial int, thereafter int)`

| 参数 | 含义 | 示例 |
|------|------|------|
| `tick` | 时间窗口 | `10*time.Second` |
| `initial` | 窗口内前 N 条放行 | `3` |
| `thereafter` | 之后每 M 条放行 1 条 | `10` |

### 配置为不采样

`sampler = nil`（零值）表示不启用采样，不影响现有代码。

## 核心代码

```go
package fastlog

import (
    "sync/atomic"
    "time"
)

const (
    samplerLevels    = 6   // DEBUG ~ PANIC
    samplerBuckets   = 4096 // 每个级别下的桶数
)

// samplerCounter 采样计数器
type samplerCounter struct {
    resetAt atomic.Int64  // 窗口过期时间（UnixNano）
    counter atomic.Uint64 // 当前窗口内的计数
}

// Sampler 日志采样器
type Sampler struct {
    tick       time.Duration
    initial    int
    thereafter int
    counters   [samplerLevels][samplerBuckets]samplerCounter
}

// NewSampler 创建日志采样器
//
// 参数:
//   - tick: 时间窗口
//   - initial: 窗口内前 N 条放行
//   - thereafter: 之后每 M 条放行 1 条
func NewSampler(tick time.Duration, initial, thereafter int) *Sampler {
    return &Sampler{
        tick:       tick,
        initial:    initial,
        thereafter: thereafter,
    }
}

// Allow 判断是否放行这条日志
func (s *Sampler) Allow(level Level, msg string) bool {
    // level 转索引（DEBUG=1 → 0, INFO=2 → 1, ...）
    i := int(level) - 1
    if i < 0 || i >= samplerLevels {
        return true
    }

    // message 哈希到桶
    j := fnv32a(msg) % samplerBuckets
    c := &s.counters[i][j]

    now := time.Now()
    tn := now.UnixNano()

    // 如果窗口已过期，重置
    if tn > c.resetAt.Load() {
        c.counter.Store(1)
        c.resetAt.Store(tn + s.tick.Nanoseconds())
        return true
    }

    // 窗口内计数递增
    n := c.counter.Add(1)

    // 前 N 条放行
    if n <= uint64(s.initial) {
        return true
    }

    // 之后每 M 条放行 1 条
    if s.thereafter > 0 && (n-uint64(s.initial))%uint64(s.thereafter) == 0 {
        return true
    }

    return false
}

// fnv32a FNV-1a 哈希，无内存分配
func fnv32a(s string) uint32 {
    const (
        offset32 = 2166136261
        prime32  = 16777619
    )
    hash := uint32(offset32)
    for i := 0; i < len(s); i++ {
        hash ^= uint32(s[i])
        hash *= prime32
    }
    return hash
}
```

## Logger 集成

```go
type Logger struct {
    level     Level
    writer    io.WriteCloser
    formatter Formatter
    caller    bool
    fields    []Field
    sampler   *Sampler  // nil 表示不启用采样

    mu   sync.Mutex
    once sync.Once
}

func WithSampler(s *Sampler) Option {
    return func(l *Logger) { l.sampler = s }
}
```

`log()` 中的改动：

```go
func (l *Logger) log(level Level, msg string, fields []Field) {
    l.init()
    if !l.level.Enabled(level) {
        return
    }

    // 采样检查
    if l.sampler != nil && !l.sampler.Allow(level, msg) {
        return
    }

    entry := GetEntry()
    defer PutEntry(entry)
    // ... 后续不变
}
```

## 改动量评估

| 改动 | 说明 |
|------|------|
| 新增 `sampler.go` | 约 90 行 |
| Logger 结构体加字段 | 1 行 `sampler *Sampler` |
| 新增 `WithSampler` | 1 个 Option 函数 |
| `log()` 中插采样检查 | 3 行 |

零侵入：`sampler = nil`（零值）表示不启用，现有代码无需任何修改。

## 参考

- zap 采样实现：`zapcore/sampler.go` — 使用 `level + message` 哈希到 4096 个固定桶，atomic 计数器
- zerolog 采样：`BasicSampler` / `RandomSampler` / `BurstSampler` — 不基于内容去重，纯计数/概率控制

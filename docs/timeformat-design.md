# 可配置时间格式设计方案

> 设计目标：支持用户自定义日志时间格式，默认 RFC3339，改动最小化
> 设计日期：2026-05-18

---

## 一、概述

### 1.1 背景

目前 5 种格式化器均硬编码 `time.RFC3339` 作为时间格式：
```go
entry.Time.Format(time.RFC3339)
```

用户无法自定义时间格式，如需其他格式需要自己实现 Formatter 接口。

### 1.2 目标

1. Config 新增 `TimeFormat` 字段，默认 `time.RFC3339`
2. Entry 新增 `TimeFormat` 字段，从 Config 传递到 Formatter
3. 5 种内置 Formatter 改用 `entry.TimeFormat`
4. 零分配影响、零破坏性变更

---

## 二、改动明细

### 2.1 Config 新增字段

**文件**：`config.go`

```go
// Config 日志配置结构体
type Config struct {
    // ... 现有字段 ...

    // TimeFormat 时间格式
    // 默认 time.RFC3339，支持 Go time 包所有格式常量
    // 常用值: time.RFC3339, time.DateTime, time.TimeOnly
    TimeFormat string
}
```

**默认值设置**：在各场景配置函数中设置默认值

```go
// newDefaultConfig 创建默认配置
func newDefaultConfig() *Config {
    return &Config{
        Level:             INFO,
        Formatter:         Def{},
        Caller:            false,
        TimeFormat:        time.RFC3339,     // 新增
        OutputConsole:     true,
        NoColor:           false,
        OutputFile:        false,
        LogPath:           "",
        MaxBufferSize:     0,
        SyncInterval:      0,
    }
}
```

**Validate 验证**：

```go
func (c *Config) Validate() error {
    // ... 现有验证 ...

    // 验证时间格式
    if c.TimeFormat == "" {
        return errors.New("TimeFormat 不能为空")
    }

    // 尝试验证时间格式是否合法
    // 注意：Go 的时间格式是布局格式，无法完全静态验证
    // 只能检查是否为空

    return nil
}
```

### 2.2 Entry 新增字段

**文件**：`logger.go`

```go
// Entry 表示一条日志记录
type Entry struct {
    Time       time.Time   // 日志时间
    Level      Level       // 日志级别
    Message    string      // 日志消息
    Fields     []Field     // 结构化字段
    Caller     string      // 调用者信息
    TimeFormat string      // 时间格式 (新增)
}
```

**Pool 清空**：在 `PutEntry` 中重置新字段

```go
func (l *Logger) PutEntry(entry *Entry) {
    entry.Time = time.Time{}
    entry.Level = 0
    entry.Message = ""
    entry.Fields = nil
    entry.Caller = ""
    entry.TimeFormat = ""  // 新增：重置
    l.entryPool.Put(entry)
}
```

### 2.3 Logger.log 赋值

**文件**：`logger.go`

```go
func (l *Logger) log(lvl Level, msg string, fields []Field) {
    // ... 级别检查、采样检查 ...

    entry := l.entryPool.Get().(*Entry)
    entry.Time = time.Now()
    entry.Level = lvl
    entry.Message = msg
    entry.Fields = fields
    entry.Caller = ""
    entry.TimeFormat = l.config.TimeFormat  // 新增：从配置赋值

    // ... 调用者信息处理 ...
}
```

### 2.4 5 种 Formatter 修改

**文件**：`formatter.go`

每个格式化器中，将 `entry.Time.Format(time.RFC3339)` 替换为 `entry.Time.Format(entry.TimeFormat)`。

| 格式化器 | 替换位置 | 行数约 |
|----------|---------|--------|
| `Def.Format` | `entry.Time.Format(time.RFC3339)` | 1 行 |
| `JSON.Format` | `entry.Time.Format(time.RFC3339)` | 1 行 |
| `Simple.Format` | `entry.Time.Format(time.RFC3339)` | 1 行 |
| `KV.Format` | `entry.Time.Format(time.RFC3339)` | 1 行 |
| `Compact.Format` | `entry.Time.Format(time.RFC3339)` | 1 行 |

### 2.5 测试更新

**文件**：`config_test.go` - 新增 TimeFormat 验证

```go
func TestConfigTimeFormat(t *testing.T) {
    t.Run("default time format", func(t *testing.T) {
        cfg := Default()
        if cfg.TimeFormat != time.RFC3339 {
            t.Errorf("TimeFormat = %q, want %q", cfg.TimeFormat, time.RFC3339)
        }
    })

    t.Run("custom time format", func(t *testing.T) {
        cfg := NewConfig("logs/app.log")
        cfg.TimeFormat = time.DateTime
        if cfg.TimeFormat != time.DateTime {
            t.Errorf("TimeFormat = %q, want %q", cfg.TimeFormat, time.DateTime)
        }
    })

    t.Run("empty time format validation", func(t *testing.T) {
        cfg := Console()
        cfg.TimeFormat = ""
        err := cfg.Validate()
        if err == nil {
            t.Error("expected error for empty TimeFormat")
        }
    })
}
```

**文件**：`formatter_test.go` - 新增各种格式的验证

```go
func TestDefFormatWithCustomTimeFormat(t *testing.T) {
    entry := &Entry{
        Time:       time.Date(2026, 1, 15, 10, 30, 45, 0, time.UTC),
        Level:      INFO,
        Message:    "测试",
        TimeFormat: time.DateTime,  // 2026-01-15 10:30:45
    }

    b, _ := Def{}.Format(entry)
    output := string(b)

    if !strings.Contains(output, "2026-01-15 10:30:45") {
        t.Errorf("Def format with DateTime should contain '2026-01-15 10:30:45', got: %s", output)
    }
}

func TestJSONFormatWithCustomTimeFormat(t *testing.T) {
    entry := &Entry{
        Time:       time.Date(2026, 1, 15, 10, 30, 45, 0, time.UTC),
        Level:      INFO,
        Message:    "测试",
        TimeFormat: time.RFC822,  // 15 Jan 26 10:30 UTC
    }

    b, _ := JSON{}.Format(entry)
    output := string(b)

    if !strings.Contains(output, "15 Jan 26 10:30 UTC") {
        t.Errorf("JSON format with RFC822 should contain '15 Jan 26 10:30 UTC', got: %s", output)
    }
}

func TestSimpleFormatWithCustomTimeFormat(t *testing.T) {
    entry := &Entry{
        Time:       time.Date(2026, 6, 18, 15, 30, 0, 0, time.UTC),
        Level:      INFO,
        Message:    "test",
        TimeFormat: "2006-01-02",  // 只显示日期
    }

    b, _ := Simple{}.Format(entry)
    output := string(b)

    if !strings.HasPrefix(output, "2026-06-18") {
        t.Errorf("Simple format with custom format should start with '2026-06-18', got: %s", output)
    }
}

func TestKVFormatWithCustomTimeFormat(t *testing.T) {
    entry := &Entry{
        Time:       time.Date(2026, 6, 18, 15, 30, 0, 0, time.UTC),
        Level:      INFO,
        Message:    "test",
        TimeFormat: time.UnixDate,
    }

    b, _ := KV{}.Format(entry)
    output := string(b)

    if !strings.Contains(output, "time=") {
        t.Errorf("KV format should contain 'time=', got: %s", output)
    }
}
```

---

## 三、数据流

```go
// 1. 用户配置
cfg := fastlog.NewConfig("logs/app.log")
cfg.TimeFormat = time.DateTime

// 2. 创建 Logger
logger := fastlog.New(cfg)

// 3. 记录日志时
logger.Info("你好")
//     ↓
//     Logger.log()
//     ↓
//     entry.TimeFormat = l.config.TimeFormat  // "2006-01-02 15:04:05"
//     ↓
//     formatter.Format(entry)
//     ↓
//     entry.Time.Format(entry.TimeFormat)
//     ↓
//     2026-05-18 15:30:00 | INFO | main.go:main:15 - 你好
```

---

## 四、性能分析

| 操作 | 之前 | 之后 | 差异 |
|------|------|------|------|
| Entry 大小 | ~84 字节 | ~100 字节 | +16 字节 |
| Pool 复用 | ✅ 是 | ✅ 是 | 一样 |
| 堆分配 | ✅ 无额外分配 | ✅ 无额外分配 | 一样 |
| `time.Format` | 常量参数 | 变量参数 | 无差异 |

**结论**：零性能影响。

---

## 五、向后兼容性

| 维度 | 评估 | 说明 |
|------|------|------|
| API 兼容 | ✅ 完全兼容 | 新增字段，默认值等于旧行为 |
| Formatter 接口 | ✅ 完全兼容 | 不改接口签名 |
| 现有配置 | ✅ 完全兼容 | 默认 `time.RFC3339` |
| 自定义 Formatter | ✅ 完全兼容 | 不强制使用 `TimeFormat` |

---

## 六、实现步骤

### Step 1：Config 加字段
- `config.go`：`TimeFormat string` + 默认值 + Validate

### Step 2：Entry 加字段
- `logger.go`：`TimeFormat string` + PutEntry 重置

### Step 3：Logger.log 赋值
- `logger.go`：`entry.TimeFormat = l.config.TimeFormat`

### Step 4：Formatter 替换
- `formatter.go`：5 处 `time.RFC3339` → `entry.TimeFormat`

### Step 5：测试
- `config_test.go`：默认值 + Validate
- `formatter_test.go`：各格式自定义时间

---

> 方案完成，等待审阅

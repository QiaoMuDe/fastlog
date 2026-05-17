# fastlog - Go 语言高性能日志库设计方案

## 一、项目概述

fastlog 是一个专为 Go 语言设计的高性能、易扩展的日志库，提供简洁的 API 设计和丰富的功能特性。

### 核心特性
- 简洁易用的 API 设计
- 高性能与可靠性
- 支持多种日志级别和格式化方法
- 灵活的初始化方式
- 可扩展的 Writer（使用标准库 io.Writer）
- 调用者信息记录
- 多种日志格式支持

---

## 二、目录结构

```
fastlog/
├── fastlog.go              # 主入口文件，提供便捷函数
├── logger.go            # Logger 核心实现
├── level.go             # 日志级别定义
├── field.go             # 字段类型定义
├── formatter.go         # 格式化器接口 + 所有实现
├── writer.go            # Writer 接口 + 所有实现（暂不实现）
├── options.go           # 配置选项
├── utils.go             # 工具函数
└── example/
    └── main.go          # 使用示例
```

---

## 三、核心接口设计

### 3.1 Logger 接口

```go
// Logger 定义日志记录器接口
type Logger interface {
    // 标准日志方法
    Debug(msg string)
    Info(msg string)
    Warn(msg string)
    Error(msg string)
    Fatal(msg string)
    Panic(msg string)

    // 格式化日志方法
    Debugf(format string, args ...interface{})
    Infof(format string, args ...interface{})
    Warnf(format string, args ...interface{})
    Errorf(format string, args ...interface{})
    Fatalf(format string, args ...interface{})
    Panicf(format string, args ...interface{})

    // 键值对日志方法
    Debugw(msg string, fields ...Field)
    Infow(msg string, fields ...Field)
    Warnw(msg string, fields ...Field)
    Errorw(msg string, fields ...Field)
    Fatalw(msg string, fields ...Field)
    Panicw(msg string, fields ...Field)

    // 配置方法
    With(fields ...Field) Logger
    WithLevel(level Level) Logger
    WithWriter(w Writer) Logger
    WithFormatter(f Formatter) Logger
    WithCaller(enabled bool) Logger

    // 同步方法
    Sync() error
}
```

### 3.2 Writer 接口

直接使用 Go 标准库的 `io.WriteCloser` 组合接口：

```go
import "io"

// Logger 使用标准库的 io.WriteCloser 接口
type Logger struct {
    Writer io.WriteCloser  // 标准库 io.WriteCloser 接口
    // ... 其他字段
}

// 写入日志
func (l *Logger) write(p []byte) error {
    _, err := l.Writer.Write(p)
    return err
}

// 关闭 Writer
func (l *Logger) Close() error {
    return l.Writer.Close()
}

// 同步（如果底层支持）
func (l *Logger) Sync() error {
    if syncer, ok := l.Writer.(interface{ Sync() error }); ok {
        return syncer.Sync()
    }
    return nil
}
```

**内置 Writer 实现**：
- `ConsoleWriter` - 包装 `os.Stdout`/`os.Stderr`，Close() 为空操作
- `FileWriter` - 包装 `*os.File`，支持 Close() 和 Sync()
- `MultiWriter` - 包装多个 `io.WriteCloser`，依次关闭

### 3.3 Formatter 接口

```go
// Formatter 定义日志格式化器接口
type Formatter interface {
    // Format 将日志条目格式化为字节数组
    Format(entry *Entry) ([]byte, error)
}

// Entry 表示一条日志记录
type Entry struct {
    Time    time.Time
    Level   Level
    Message string
    Caller  string        // 调用者信息: file.go:func:line
    Fields  []Field       // 键值对字段
}
```

### 3.4 Field 类型

使用包含所有类型字段的结构体设计，避免 `interface{}` 的类型断言开销：

```go
// FieldType 表示字段类型
 type FieldType int8
 
 const (
     UnknownType FieldType = iota
     StringType
     IntType
     Int64Type
     UintType
     Uint64Type
     Float64Type
     BoolType
     TimeType
     DurationType
     ErrorType
     AnyType
 )
 
 // Field 表示一个键值对字段，包含所有可能的类型
 type Field struct {
     Key       string
     Type      FieldType
     StringVal string
     IntVal    int64
     UintVal   uint64
     FloatVal  float64
     BoolVal   bool
     TimeVal   time.Time
     Duration  time.Duration
     Interface interface{}
 }
 
 // 字段构造函数 - 每个函数只设置相关字段
 func String(key, val string) Field {
     return Field{Key: key, Type: StringType, StringVal: val}
 }
 
 func Int(key string, val int) Field {
     return Field{Key: key, Type: IntType, IntVal: int64(val)}
 }
 
 func Int64(key string, val int64) Field {
     return Field{Key: key, Type: Int64Type, IntVal: val}
 }
 
 func Uint(key string, val uint) Field {
     return Field{Key: key, Type: UintType, UintVal: uint64(val)}
 }
 
 func Uint64(key string, val uint64) Field {
     return Field{Key: key, Type: Uint64Type, UintVal: val}
 }
 
 func Float64(key string, val float64) Field {
     return Field{Key: key, Type: Float64Type, FloatVal: val}
 }
 
 func Bool(key string, val bool) Field {
     return Field{Key: key, Type: BoolType, BoolVal: val}
 }
 
 func Time(key string, val time.Time) Field {
     return Field{Key: key, Type: TimeType, TimeVal: val}
 }
 
 func Duration(key string, val time.Duration) Field {
     return Field{Key: key, Type: DurationType, Duration: val}
 }
 
 func Error(err error) Field {
     if err == nil {
         return Field{Key: "error", Type: StringType, StringVal: "<nil>"}
     }
     return Field{Key: "error", Type: ErrorType, StringVal: err.Error()}
 }
 
 func Any(key string, val interface{}) Field {
     return Field{Key: key, Type: AnyType, Interface: val}
 }
 
 func Stack() Field {
     return Field{Key: "stack", Type: StringType, StringVal: captureStack()}
 }
 
 // String 返回 string 值，类型不匹配返回空字符串
 func (f Field) String() string {
     if f.Type == StringType || f.Type == ErrorType {
         return f.StringVal
     }
     return ""
 }
 
 // Int 返回 int 值，类型不匹配返回 0
 func (f Field) Int() int {
     if f.Type == IntType {
         return int(f.IntVal)
     }
     return 0
 }
 
 // Int64 返回 int64 值，类型不匹配返回 0
 func (f Field) Int64() int64 {
     if f.Type == IntType || f.Type == Int64Type {
         return f.IntVal
     }
     return 0
 }
 
 // Uint 返回 uint 值，类型不匹配返回 0
 func (f Field) Uint() uint {
     if f.Type == UintType {
         return uint(f.UintVal)
     }
     return 0
 }
 
 // Uint64 返回 uint64 值，类型不匹配返回 0
 func (f Field) Uint64() uint64 {
     if f.Type == UintType || f.Type == Uint64Type {
         return f.UintVal
     }
     return 0
 }
 
 // Float64 返回 float64 值，类型不匹配返回 0
 func (f Field) Float64() float64 {
     if f.Type == Float64Type {
         return f.FloatVal
     }
     return 0
 }
 
 // Bool 返回 bool 值，类型不匹配返回 false
 func (f Field) Bool() bool {
     if f.Type == BoolType {
         return f.BoolVal
     }
     return false
 }
 
 // Time 返回 time.Time 值，类型不匹配返回零值
 func (f Field) Time() time.Time {
     if f.Type == TimeType {
         return f.TimeVal
     }
     return time.Time{}
 }
 
 // Duration 返回 time.Duration 值，类型不匹配返回 0
 func (f Field) Duration() time.Duration {
     if f.Type == DurationType {
         return f.Duration
     }
     return 0
 }
 ```
 
 **泛型方法使用示例**：
 ```go
 field := Int("count", 42)
 
 // 获取 int 值
 val, ok := Value[int](field)
 // val = 42, ok = true
 
 // 获取 int64 值（自动转换）
 val64, ok := Value[int64](field)
 // val64 = 42, ok = true
 
 // 获取 string 值（类型不匹配）
 str, ok := Value[string](field)
 // str = "", ok = false
 
 // 直接使用 MustValue
 count := MustValue[int](field)  // count = 42
 ```

**设计优势**：
- **零分配**：结构体在栈上分配，无堆分配
- **无类型断言**：根据 `Type` 字段直接访问对应值
- **高性能**：格式化时无需反射或类型 switch
- **类型安全**：编译期保证类型正确

---

## 四、日志级别设计

```go
// Level 表示日志级别
type Level int8

const (
    DEBUG Level = iota - 1  // 调试级别
    INFO                     // 信息级别
    WARN                     // 警告级别
    ERROR                    // 错误级别
    FATAL                    // 致命级别
    PANIC                    // 恐慌级别
)

// 级别方法
func (l Level) String() string           // 返回级别字符串
func (l Level) Enabled(lvl Level) bool   // 检查是否启用该级别
func ParseLevel(s string) (Level, error) // 从字符串解析级别
```

**级别优先级**: DEBUG < INFO < WARN < ERROR < FATAL < PANIC

---

## 五、格式化器实现

### 5.1 Def 格式（默认格式）

```
2025-01-15T10:30:45 | INFO    | main.go:main:15 - 用户登录成功
2025-01-15T10:30:46 | ERROR   | database.go:Connect:23 - 数据库连接失败
```

**特点**:
- 时间戳使用 ISO8601 格式
- 级别右对齐，固定 8 字符宽度
- 包含调用者信息
- 简洁直观

### 5.2 JSON 格式

```json
{"time":"2025-01-15T10:30:45","level":"INFO","caller":"main.go:main:15","message":"用户登录成功"}
{"time":"2025-01-15T10:30:46","level":"ERROR","caller":"database.go:Connect:23","message":"数据库连接失败"}
```

**特点**:
- 结构化数据，便于机器解析
- 支持嵌套字段
- 适合日志收集系统

### 5.3 Timestamp 格式

```
2025-01-15T10:30:45 INFO  用户登录成功
2025-01-15T10:30:46 ERROR 数据库连接失败
```

**特点**:
- 极简格式
- 仅包含时间、级别和消息
- 适合开发调试

### 5.4 KVFmt 键值对格式

```
time=2025-01-15T10:30:45 level=INFO message=用户登录成功
time=2025-01-15T10:30:46 level=ERROR message=数据库连接失败
```

**特点**:
- 纯键值对形式
- 易于 grep 和解析
- 适合命令行工具

### 5.5 LogFmt 格式

```
2025-01-15T10:30:45 [INFO ] 用户登录成功 [username=张三, age=30]
2025-01-15T10:30:46 [ERROR] database.go:Connect:23 数据库连接失败
```

**特点**:
- 人类可读性好
- 字段使用方括号包裹
- 适合查看和分析

---

## 六、初始化方式

### 6.1 结构体字面量方式

```go
logger := &fastlog.Logger{
    Level:     fastlog.InfoLevel,
    Writer:    fastlog.NewConsoleWriter(),
    Formatter: fastlog.NewDefFormatter(),
    Caller:    true,
}
```

### 6.2 构造函数方式

```go
// 使用默认配置
logger := fastlog.New()

// 使用自定义配置
logger := fastlog.New(
    fastlog.WithLevel(fastlog.DEBUG),
    fastlog.WithWriter(os.Stdout),  // 纯日志引擎，不内置文件写入器
    fastlog.WithFormatter(fastlog.JSON{}),
    fastlog.WithCaller(true),
)
```

### 6.3 全局默认实例

```go
// 直接使用包级别函数
fastlog.Info("应用启动")
fastlog.Infof("监听端口: %d", 8080)
fastlog.Infow("用户登录", fastlog.String("user", "admin"))

// 自定义全局实例
fastlog.SetDefault(fastlog.New(
    fastlog.WithLevel(fastlog.DebugLevel),
))
```

---

## 七、Writer 实现

> **注意**：此部分设计暂不实现，待后续讨论确定。

### 7.1 控制台写入器

暂不实现。

### 7.2 文件写入器

暂不实现。

### 7.3 多路写入器

暂不实现。

### 7.4 自定义 Writer

暂不实现。

---

## 八、使用示例

### 8.1 基础使用

```go
package main

import (
    "github.com/yourname/fastlog"
)

func main() {
    // 使用默认配置
    logger := fastlog.New()
    defer logger.Sync()

    // 标准日志
    logger.Info("应用启动成功")
    logger.Error("发生错误")

    // 格式化日志
    logger.Infof("当前版本: %s", "v1.0.0")
    logger.Errorf("操作失败: %v", err)

    // 键值对日志
    logger.Infow("用户登录",
        fastlog.String("username", "admin"),
        fastlog.Int("user_id", 123),
        fastlog.String("ip", "192.168.1.1"),
    )
}
```

### 8.2 高级配置

```go
package main

import (
    "github.com/yourname/fastlog"
)

func main() {
    // 使用 io.WriteCloser 接口注入外部写入器
    // 例如与 logrotatex 库配合实现文件轮转：
    //
    // import "gitee.com/MM-Q/logrotatex"
    // rotator := logrotatex.NewLogRotateX("logs/app.log")
    // bw := logrotatex.NewBufferedWriter(rotator, nil)
    //
    // logger := fastlog.New(
    //     fastlog.WithWriter(bw),
    //     fastlog.WithFormatter(fastlog.JSON{}),
    // )

    // 创建多路写入器（同时输出到控制台和文件）
    multiWriter := fastlog.NewMultiWriter(
        fastlog.NewColorWriter(false),
        fileWriter,  // 外部文件写入器
    )

    // 创建 Logger
    logger := fastlog.New(
        fastlog.WithLevel(fastlog.DEBUG),
        fastlog.WithWriter(multiWriter),
        fastlog.WithFormatter(fastlog.JSON{}),
        fastlog.WithCaller(true),
    )
    defer logger.Sync()

    // 使用子 Logger（继承字段）
    userLogger := logger.With(
        fastlog.String("module", "user"),
        fastlog.String("service", "auth"),
    )

    userLogger.Infow("用户认证成功",
        fastlog.String("username", "admin"),
    )
}
```

### 8.3 使用全局实例

```go
package main

import (
    "github.com/yourname/fastlog"
)

func init() {
    // 配置全局 Logger
    fastlog.SetDefault(fastlog.New(
        fastlog.WithLevel(fastlog.InfoLevel),
        fastlog.WithFormatter(fastlog.NewDefFormatter()),
    ))
}

func main() {
    // 直接使用包级别函数
    fastlog.Info("服务启动")
    
    doSomething()
}

func doSomething() {
    fastlog.Debug("调试信息")  // 不会输出，因为级别是 Info
    fastlog.Info("处理请求")
    fastlog.Infow("请求完成",
        fastlog.Int("status", 200),
        fastlog.Duration("latency", 100*time.Millisecond),
    )
}
```

---

## 九、性能优化策略

### 9.1 零分配优化

```go
// 使用 sync.Pool 复用对象
var entryPool = sync.Pool{
    New: func() interface{} {
        return &Entry{
            Fields: make([]Field, 0, 8),
        }
    },
}
```

### 9.2 异步写入

```go
// 使用缓冲通道实现异步写入
type AsyncWriter struct {
    writer Writer
    queue  chan []byte
    wg     sync.WaitGroup
}
```

### 9.3 级别检查前置

```go
func (l *logger) Info(msg string) {
    if !l.level.Enabled(InfoLevel) {
        return  // 快速返回，避免不必要的处理
    }
    // ... 处理日志
}
```

---

## 十、线程安全

- Logger 实例是线程安全的，可以在多个 goroutine 中并发使用
- Writer 实现需要保证线程安全或使用互斥锁
- 使用 sync.Mutex 保护共享状态

---

## 十一、错误处理

```go
// 日志写入错误处理
type ErrorHandler func(err error)

// Logger 配置中添加错误处理器
logger := fastlog.New(
    fastlog.WithErrorHandler(func(err error) {
        // 处理日志写入错误，如发送到监控系统
        metrics.Increment("log.write.errors")
    }),
)
```

---

## 十二、扩展性设计

### 12.1 Hook 机制

```go
// Hook 接口允许在日志处理过程中插入自定义逻辑
type Hook interface {
    Levels() []Level
    Fire(entry *Entry) error
}

// 添加 Hook
logger.AddHook(hook)
```

### 12.2 采样器

```go
// 日志采样，防止日志风暴
type Sampler interface {
    Sample(entry *Entry) bool
}

logger := fastlog.New(
    fastlog.WithSampler(fastlog.BurstSampler(100, time.Second)),
)
```

---

## 十三、API 设计原则

1. **简洁性**: API 命名清晰，符合 Go 语言习惯
2. **一致性**: 类似功能使用相似的命名模式
3. **可发现性**: 通过 IDE 自动补全即可发现功能
4. **向后兼容**: 版本升级保持 API 兼容

---

## 十四、总结

fastlog 日志库设计方案遵循以下原则：

1. **高性能**: 零分配、异步写入、级别前置检查
2. **易用性**: 简洁的 API、多种初始化方式
3. **可扩展**: Writer/Formatter 接口、Hook 机制
4. **灵活性**: 多种格式、配置选项丰富
5. **可靠性**: 线程安全、错误处理、日志采样

该设计参考了 zap、logrus 等优秀日志库的优点，同时针对 Go 语言特性进行了优化。

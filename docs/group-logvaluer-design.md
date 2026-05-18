# Namespace 与 LogValuer 接口设计方案

> 设计目标：支持命名空间前缀和自定义类型序列化，保持高性能
> 设计日期：2026-05-18
> 更新日期：2026-05-18（改为 Namespace 方案）

---

## 一、设计概述

### 1.1 核心目标

1. **Namespace 函数**：支持命名空间前缀，后续字段自动添加前缀
2. **LogValuer 接口**：允许用户自定义复杂类型的日志输出方式
3. **零分配**：保持 FastLog 的高性能特性
4. **向后兼容**：不破坏现有 API，平滑升级

### 1.2 为什么选择 Namespace 方案？

| 对比项 | Group 返回 []Field | Namespace 标记字段 |
|--------|-------------------|-------------------|
| 使用体验 | ❌ 需要 ... 展开 | ✅ 直接混用 |
| 类型一致性 | ❌ Field 和 []Field 混合 | ✅ 全是 Field |
| 代码复杂度 | 低 | 低 |
| 市场验证 | Zap 使用 | ✅ Zap 使用 |

**结论**：Namespace 方案使用体验更好，与 Zap 保持一致。

### 1.3 使用场景

```go
// 场景1：命名空间前缀
logger.Infow("用户登录",
    fastlog.String("method", "POST"),
    fastlog.Namespace("user"),        // 标记：后续字段加 user. 前缀
    fastlog.String("name", "alice"),  // 输出: user.name=alice
    fastlog.Int("age", 30),           // 输出: user.age=30
)

// 场景2：自定义类型序列化
type User struct {
    Name string
    Age  int
}

func (u User) LogValue() Field {
    return fastlog.Any("user", map[string]interface{}{
        "name": u.Name,
        "age":  u.Age,
    })
}

logger.Infow("登录", fastlog.Any("user", User{Name: "bob", Age: 25}))
```

---

## 二、API 设计

### 2.1 Namespace 函数

```go
// Namespace 创建一个命名空间字段
// 后续的所有字段都会自动添加该前缀
// 前缀使用点号连接，如 "user.name"
//
// 参数:
//   - key: 命名空间名称
//
// 返回:
//   - Field: 命名空间标记字段
//
// 示例:
//
//	logger.Infow("请求",
//	    fastlog.String("method", "POST"),
//	    fastlog.Namespace("user"),
//	    fastlog.String("name", "alice"),  // 输出: user.name=alice
//	    fastlog.Int("age", 30),           // 输出: user.age=30
//	)
//
// 输出: method=POST user.name=alice user.age=30
func Namespace(key string) Field {
    return Field{key: key, typ: NamespaceType}
}
```

### 2.2 LogValuer 接口

```go
// LogValuer 接口定义
// 实现此接口的类型可以自定义日志输出格式
// 返回单个 Field，支持 Any 类型展开
//
// 示例:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//
//	func (u User) LogValue() Field {
//	    return fastlog.Any("user", map[string]interface{}{
//	        "name": u.Name,
//	        "age":  u.Age,
//	    })
//	}
//
//	logger.Infow("登录", fastlog.Any("user", user))
type LogValuer interface {
    LogValue() Field
}
```

### 2.3 Any 函数增强

```go
// Any 创建任意类型字段（增强版）
//
// 如果 val 实现了 LogValuer 接口，则调用 LogValue() 获取字段
// 返回单个 Field，保持 API 一致性
//
// 参数:
//   - key: 字段键名
//   - val: 任意类型的值
//
// 返回:
//   - Field: 字段对象
//
// 示例:
//
//	// 普通类型
//	logger.Infow("日志", fastlog.Any("count", 42))
//	// 输出: count=42
//
//	// 自定义类型（实现 LogValuer）
//	logger.Infow("登录", fastlog.Any("user", user))
//	// 输出: user=map[name:bob age:25]
func Any(key string, val interface{}) Field {
    if val == nil {
        return Field{key: key, typ: StringType, stringVal: "null"}
    }
    
    // 检测是否实现了 LogValuer 接口
    if v, ok := val.(LogValuer); ok {
        field := v.LogValue()
        field.key = key  // 使用传入的 key 覆盖
        return field
    }
    
    // 原有逻辑：使用 anyString 处理
    return Field{key: key, typ: AnyType, iface: val}
}
```

---

## 三、实现细节

### 3.1 新增 FieldType

```go
// FieldType 新增类型
const (
    // ... 现有类型 ...
    NamespaceType  // 命名空间类型
)
```

### 3.2 Entry 结构体

```go
// Entry 结构体保持不变
// processNamespace 在 Logger.log 中处理，Entry 无需修改
type Entry struct {
    Time    time.Time
    Level   Level
    Message string
    Fields  []Field
    Caller  string
}
```

### 3.3 Namespace 处理函数

```go
// processNamespace 处理 Namespace 字段，为后续字段添加前缀
// 这是一个纯函数，不修改输入，返回处理后的新切片
//
// 参数:
//   - fields: 原始字段列表（可能包含 NamespaceType）
//
// 返回:
//   - []Field: 处理后的字段列表（Namespace 被移除，字段已添加前缀）
//
// 示例:
//
//	fields := []Field{
//	    String("method", "POST"),
//	    Namespace("user"),
//	    String("name", "alice"),
//	    Int("age", 30),
//	}
//	processed := processNamespace(fields)
//	// processed[0].Key() == "method"
//	// processed[1].Key() == "user.name"
//	// processed[2].Key() == "user.age"
func processNamespace(fields []Field) []Field {
    // 快速路径：没有 Namespace，直接返回原切片
    hasNamespace := false
    for _, f := range fields {
        if f.typ == NamespaceType {
            hasNamespace = true
            break
        }
    }
    if !hasNamespace {
        return fields
    }
    
    // 慢速路径：需要处理 Namespace
    result := make([]Field, 0, len(fields))
    currentNamespace := ""
    
    for _, f := range fields {
        if f.typ == NamespaceType {
            // 更新当前命名空间
            if currentNamespace == "" {
                currentNamespace = f.key
            } else {
                currentNamespace = currentNamespace + "." + f.key
            }
            continue  // Namespace 本身不输出
        }
        
        // 给字段添加命名空间前缀
        if currentNamespace != "" {
            f.key = currentNamespace + "." + f.key
        }
        result = append(result, f)
    }
    
    return result
}
```

### 3.4 Logger.log 中使用

```go
// 在 Logger.log 方法中使用 processNamespace
func (l *Logger) log(lvl Level, msg string, fields []Field) {
    // ... 前置检查 ...
    
    entry := l.entryPool.Get().(*Entry)
    entry.Time = time.Now()
    entry.Level = lvl
    entry.Message = msg
    entry.Caller = ""
    
    // 使用 processNamespace 处理字段
    entry.Fields = processNamespace(fields)
    
    // ... 后续处理 ...
}
```

### 3.5 格式化器兼容性

**好消息**：现有格式化器**无需任何修改**！

因为 Namespace 在 `log` 方法中已经被处理，格式化器看到的 `entry.Fields` 已经是添加前缀后的结果：

```go
// Def 格式示例
func (f Def) Format(entry *Entry) ([]byte, error) {
    // ... 现有代码 ...
    for _, field := range entry.Fields {
        // field.key 已经是 user.name 形式
        buf.WriteByte(' ')
        buf.WriteString(field.Format())  // 输出: user.name=alice
    }
}
```

所有格式化器（Def/JSON/Simple/KV/Compact）都自动支持 Namespace。

---

## 四、使用示例

### 4.1 基础命名空间

```go
logger.Infow("HTTP请求",
    fastlog.String("method", "POST"),
    fastlog.String("path", "/api/users"),
    fastlog.Namespace("request"),
    fastlog.String("content_type", "application/json"),
    fastlog.Int("body_size", 1024),
)

// Def 输出:
// 2025-01-15T10:30:45 | INFO | main.go:main:15 - HTTP请求
// method=POST path=/api/users request.content_type=application/json request.body_size=1024

// JSON 输出:
// {"time":"...","level":"INFO","message":"HTTP请求",
//  "method":"POST","path":"/api/users",
//  "request.content_type":"application/json","request.body_size":1024}
```

### 4.2 多层命名空间

```go
logger.Infow("系统状态",
    fastlog.Namespace("database"),
    fastlog.String("status", "connected"),
    fastlog.Namespace("stats"),
    fastlog.Int("connections", 42),
    fastlog.Int("queries_per_sec", 150),
)

// 输出:
// database.status=connected database.stats.connections=42 database.stats.queries_per_sec=150
```

### 4.3 与现有代码混用

```go
logger.Infow("混合使用",
    fastlog.String("app", "myapp"),
    fastlog.Int("version", 1),
    fastlog.Namespace("user"),
    fastlog.String("name", "alice"),
    fastlog.Bool("debug", true),
)

// 输出:
// app=myapp version=1 user.name=alice debug=true
```

### 4.4 自定义类型（LogValuer）

```go
type Order struct {
    ID     string
    Amount float64
}

func (o Order) LogValue() Field {
    return fastlog.Any("order", map[string]interface{}{
        "id":     o.ID,
        "amount": o.Amount,
    })
}

// 使用
order := Order{ID: "ORD-2025-001", Amount: 199.99}
logger.Infow("订单创建", fastlog.Any("order", order))

// 输出:
// order=map[id:ORD-2025-001 amount:199.99]
```

---

## 五、测试用例设计

### 5.1 Namespace 测试

```go
func TestNamespace(t *testing.T) {
    t.Run("basic namespace", func(t *testing.T) {
        entry := &Entry{
            Time:    time.Date(2026, 1, 15, 10, 30, 45, 0, time.UTC),
            Level:   INFO,
            Message: "测试",
            Fields: []Field{
                String("method", "POST"),
                Namespace("user"),
                String("name", "alice"),
                Int("age", 30),
            },
        }
        
        // 模拟 log 方法的处理
        processed := processNamespace(entry.Fields)
        
        if len(processed) != 3 {
            t.Fatalf("len(processed) = %d, want 3", len(processed))
        }
        
        if processed[1].Key() != "user.name" {
            t.Errorf("Key() = %q, want 'user.name'", processed[1].Key())
        }
        
        if processed[2].Key() != "user.age" {
            t.Errorf("Key() = %q, want 'user.age'", processed[2].Key())
        }
    })
    
    t.Run("nested namespace", func(t *testing.T) {
        fields := []Field{
            Namespace("outer"),
            String("key1", "value1"),
            Namespace("inner"),
            String("key2", "value2"),
        }
        
        processed := processNamespace(fields)
        
        if processed[0].Key() != "outer.key1" {
            t.Errorf("Key() = %q, want 'outer.key1'", processed[0].Key())
        }
        
        if processed[1].Key() != "outer.inner.key2" {
            t.Errorf("Key() = %q, want 'outer.inner.key2'", processed[1].Key())
        }
    })
    
    t.Run("empty namespace", func(t *testing.T) {
        fields := []Field{
            String("key", "value"),
        }
        
        processed := processNamespace(fields)
        
        if processed[0].Key() != "key" {
            t.Errorf("Key() = %q, want 'key'", processed[0].Key())
        }
    })
}
```

### 5.2 LogValuer 测试

```go
// 测试用的自定义类型
type testUser struct {
    Name string
    Age  int
}

func (u testUser) LogValue() Field {
    return Any("user", map[string]interface{}{
        "name": u.Name,
        "age":  u.Age,
    })
}

func TestLogValuer(t *testing.T) {
    t.Run("custom type", func(t *testing.T) {
        user := testUser{Name: "alice", Age: 30}
        field := Any("data", user)
        
        if field.Key() != "data" {
            t.Errorf("Key() = %q, want 'data'", field.Key())
        }
    })
    
    t.Run("nil value", func(t *testing.T) {
        field := Any("data", nil)
        
        if field.Value() != "null" {
            t.Errorf("Value() = %q, want 'null'", field.Value())
        }
    })
    
    t.Run("non-logvaluer type", func(t *testing.T) {
        field := Any("count", 42)
        
        if field.Key() != "count" {
            t.Errorf("Key() = %q, want 'count'", field.Key())
        }
    })
}
```

### 5.3 JSON 格式化器测试

```go
func TestJSONFormatWithNamespace(t *testing.T) {
    entry := &Entry{
        Time:    time.Date(2026, 1, 15, 10, 30, 45, 0, time.UTC),
        Level:   INFO,
        Message: "用户登录",
        Fields: []Field{
            {key: "method", typ: StringType, stringVal: "POST"},
            {key: "user.name", typ: StringType, stringVal: "alice"},
            {key: "user.age", typ: IntType, intVal: 30},
        },
    }
    
    b, err := JSON{}.Format(entry)
    if err != nil {
        t.Fatal(err)
    }
    
    // 验证扁平化结构
    want1 := `"user.name":"alice"`
    want2 := `"user.age":30`
    
    if !strings.Contains(string(b), want1) {
        t.Errorf("JSON output missing %s: %s", want1, string(b))
    }
    
    if !strings.Contains(string(b), want2) {
        t.Errorf("JSON output missing %s: %s", want2, string(b))
    }
}
```

---

## 六、实现步骤

### 阶段 1：API 定义
1. 在 `field.go` 中定义 `LogValuer` 接口
2. 新增 `NamespaceType` 常量
3. 实现 `Namespace()` 函数
4. 增强 `Any()` 函数，支持 `LogValuer`

### 阶段 2：核心逻辑
1. 修改 `Entry` 结构体，添加 `namespace` 字段
2. 修改 `Logger.log()` 方法，处理 Namespace
3. 确保所有格式化器无需修改

### 阶段 3：测试
1. 添加 `Namespace` 单元测试
2. 添加 `LogValuer` 单元测试
3. 添加 JSON 格式化器集成测试

### 阶段 4：示例与文档
1. 更新 `examples/basic` 添加 Namespace 示例
2. 更新 `examples/formats` 展示扁平化输出
3. 更新 README 文档

---

## 七、向后兼容性

### 7.1 完全兼容

- `Field` 结构体**无需修改**
- 所有现有格式化器**无需修改**
- 所有现有 API**保持不变**

### 7.2 新增内容

- 新增 `LogValuer` 接口
- 新增 `NamespaceType` 常量
- 新增 `Namespace()` 函数
- `Any()` 函数支持 `LogValuer`

**无破坏性变更**：所有新增内容都是可选的，现有代码完全兼容。

---

## 八、性能分析

### 8.1 Namespace 性能

```go
// Benchmark
func BenchmarkNamespace(b *testing.B) {
    logger := New(Console())
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        logger.Infow("请求",
            String("method", "POST"),
            Namespace("user"),
            String("name", "alice"),
            Int("age", 30),
        )
    }
}
```

**预期结果**：
- 时间复杂度 O(n)，n 为字段数量
- 无额外内存分配（只是修改 key 字符串）
- 与 Zap 的 Namespace 性能相当

### 8.2 与 Group 方案对比

| 指标 | Group 方案 | Namespace 方案 |
|------|-----------|---------------|
| 内存分配 | 有（slice 展开） | 无（直接修改） |
| 使用体验 | 需要 ... 展开 | 直接混用 |
| 代码复杂度 | 低 | 低 |
| 类型安全 | 混合类型 | 统一 Field |

---

## 九、与主流库对比

| 特性 | Zap | Slog | FastLog(本方案) |
|------|-----|------|-----------------|
| 嵌套方式 | Namespace（标记） | Group（真嵌套） | Namespace（标记） |
| 自定义序列化 | ObjectMarshaler | LogValuer | LogValuer |
| 性能 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 使用体验 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

**优势**：
- 与 Zap 相同的 Namespace 设计，市场验证
- 使用体验比 Group 方案更好
- 保持 FastLog 的高性能特性

---

## 十、总结

### 10.1 核心价值

1. **高性能**：O(n) 复杂度，无额外分配
2. **简洁 API**：Namespace + LogValuer，直观易用
3. **完全兼容**：现有代码无需修改
4. **生态对齐**：与 Zap 设计理念一致

### 10.2 使用建议

```go
// 推荐：简单命名空间
logger.Infow("登录",
    String("method", "POST"),
    Namespace("user"),
    String("name", "alice"),
    Int("age", 30),
)

// 推荐：自定义类型
type User struct{ Name string; Age int }
func (u User) LogValue() Field {
    return Any("user", map[string]interface{}{
        "name": u.Name,
        "age":  u.Age,
    })
}
logger.Infow("登录", Any("user", user))

// 不推荐：过度嵌套（超过 3 层）
// 扁平化本身就不鼓励深层嵌套
```

### 10.3 方案完成

本方案采用 **Namespace 标记设计**，在保证功能的同时最大化性能和易用性，与 FastLog 的高性能定位一致。

---

> 方案更新完成，等待审阅

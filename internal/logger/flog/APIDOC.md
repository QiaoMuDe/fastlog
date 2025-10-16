# Package flog

Package flog (import "gitee.com/MM-Q/fastlog/internal/logger/flog")

## Types

### Entry

Entry 定义了日志条目的结构体

```go
type Entry struct {
	// Has unexported fields.
}
```

#### NewEntry
NewEntry 创建一个新的日志条目（使用对象池优化）

```go
func NewEntry(needFileInfo bool, level types.LogLevel, msg string, fields ...*Field) *Entry
```

**参数：**
- needFileInfo: 是否需要文件信息。
- level: 日志级别。
- msg: 日志消息。
- *fields: 日志字段

**返回值：**
- *Entry: 一个指向 Entry 实例的指针。

---

### FLog

FLog 是一个高性能的日志记录器，支持键值对风格的使用和标准库fmt类似的使用，同时提供了丰富的配置选项，如日志级别、输出格式、日志轮转等。

```go
type FLog struct {
	// Has unexported fields.
}
```

#### NewFLog
NewFLog 创建一个新的FLog实例，用于记录日志。

```go
func NewFLog(cfg *config.FastLogConfig) *FLog
```

**参数:**
- cfg: 一个指向FastLogConfig实例的指针，用于配置日志记录器。

**返回值:**
- *FLog: 一个指向FLog实例的指针。

#### Close
Close 关闭日志处理器

```go
func (f *FLog) Close() error
```

**返回值：**
- error: 如果关闭过程中发生错误，返回错误信息；否则返回 nil。

#### Debug
Debug 记录调试级别的日志，不支持占位符

```go
func (f *FLog) Debug(v ...any)
```

**参数:**
- v: 可变参数，可以是任意类型，会被转换为字符串

#### DebugF
DebugF 记录Debug级别的键值对日志

```go
func (f *FLog) DebugF(msg string, fields ...*Field)
```

**参数：**
- msg: 日志消息。
- fields: 日志字段，可变参数。

#### Debugf
Debugf 记录调试级别的日志，支持占位符，格式化

```go
func (f *FLog) Debugf(format string, v ...any)
```

**参数:**
- format: 格式字符串
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Error
Error 记录错误级别的日志，不支持占位符

```go
func (f *FLog) Error(v ...any)
```

**参数:**
- v: 可变参数，可以是任意类型，会被转换为字符串

#### ErrorF
ErrorF 记录Error级别的键值对日志

```go
func (f *FLog) ErrorF(msg string, fields ...*Field)
```

**参数：**
- msg: 日志消息。
- fields: 日志字段，可变参数。

#### Errorf
Errorf 记录错误级别的日志，支持占位符，格式化

```go
func (f *FLog) Errorf(format string, v ...any)
```

**参数:**
- format: 格式字符串
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Fatal
Fatal 记录致命级别的日志，不支持占位符，发送后关闭日志记录器

```go
func (f *FLog) Fatal(v ...any)
```

**参数:**
- v: 可变参数，可以是任意类型，会被转换为字符串

#### FatalF
FatalF 记录Fatal级别的键值对日志并触发程序退出

```go
func (f *FLog) FatalF(msg string, fields ...*Field)
```

**参数：**
- msg: 日志消息。
- fields: 日志字段，可变参数。

#### Fatalf
Fatalf 记录致命级别的日志，支持占位符，发送后关闭日志记录器

```go
func (f *FLog) Fatalf(format string, v ...any)
```

**参数:**
- format: 格式字符串
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Info
Info 记录信息级别的日志，不支持占位符

```go
func (f *FLog) Info(v ...any)
```

**参数:**
- v: 可变参数，可以是任意类型，会被转换为字符串

#### InfoF
InfoF 记录Info级别的键值对日志

```go
func (f *FLog) InfoF(msg string, fields ...*Field)
```

**参数：**
- msg: 日志消息。
- fields: 日志字段，可变参数。

#### Infof
Infof 记录信息级别的日志，支持占位符，格式化

```go
func (f *FLog) Infof(format string, v ...any)
```

**参数:**
- format: 格式字符串
- v: 可变参数，可以是任意类型，会被转换为字符串

#### Warn
Warn 记录警告级别的日志，不支持占位符

```go
func (f *FLog) Warn(v ...any)
```

**参数:**
- v: 可变参数，可以是任意类型，会被转换为字符串

#### WarnF
WarnF 记录Warn级别的键值对日志

```go
func (f *FLog) WarnF(msg string, fields ...*Field)
```

**参数：**
- msg: 日志消息。
- fields: 日志字段，可变参数。

#### Warnf
Warnf 记录警告级别的日志，支持占位符，格式化

```go
func (f *FLog) Warnf(format string, v ...any)
```

**参数:**
- format: 格式字符串
- v: 可变参数，可以是任意类型，会被转换为字符串

---

### Field

Field 定义了日志字段的结构体

```go
type Field struct {
	// Has unexported fields.
}
```

#### Bool
Bool 添加布尔字段

```go
func Bool(key string, value bool) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为布尔值。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### Duration
Duration 添加持续时间字段

```go
func Duration(key string, value time.Duration) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为time.Duration类型。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### Error
Error 添加错误字段

```go
func Error(key string, value error) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为error类型。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### Float64
Float64 添加64位浮点数字段

```go
func Float64(key string, value float64) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为64位浮点数。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### Int
Int 添加整数字段

```go
func Int(key string, value int) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为整数。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### Int64
Int64 添加64位整数字段

```go
func Int64(key string, value int64) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为64位整数。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### String
String 添加字符串字段

```go
func String(key string, value string) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为字符串。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### Time
Time 添加时间字段

```go
func Time(key string, value time.Time) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为time.Time类型。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### Uint16
Uint16 添加16位无符号整数字段

```go
func Uint16(key string, value uint16) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为16位无符号整数。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### Uint32
Uint32 添加32位无符号整数字段

```go
func Uint32(key string, value uint32) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为32位无符号整数。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### Uint64
Uint64 添加64位无符号整数字段

```go
func Uint64(key string, value uint64) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为64位无符号整数。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### Uint8
Uint8 添加8位无符号整数字段

```go
func Uint8(key string, value uint8) *Field
```

**参数:**
- key: 字段的键名，不能为空字符串。
- value: 字段值，必须为8位无符号整数。

**返回值:**
- *Field: 一个指向 Field 实例的指针。

#### Key
Key 获取字段键名

```go
func (f *Field) Key() string
```

**返回值:**
- string: 字段键名。

#### Type
Type 获取字段类型

```go
func (f *Field) Type() FieldType
```

**返回值:**
- FieldType: 字段类型。

#### Value
Value 获取字段值

```go
func (f *Field) Value() string
```

**返回值:**
- string: 字段值的字符串表示。

---

### FieldType

FieldType 定义了日志字段的类型(13个)

```go
type FieldType uint8
```

#### 常量定义
```go
const (
	StringType   FieldType = iota // 字符串字段
	IntType                       // 整数字段
	Int64Type                     // 64位整数字段
	Float64Type                   // 64位浮点数字段
	BoolType                      // 布尔字段
	TimeType                      // 时间字段
	DurationType                  // 间隔时间字段
	Uint64Type                    // 64位无符号整数字段
	Uint32Type                    // 32位无符号整数字段
	Uint16Type                    // 16位无符号整数字段
	Uint8Type                     // 8位无符号整数字段
	ErrorType                     // 错误字段
	UnknownType                   // 未知字段
)
```

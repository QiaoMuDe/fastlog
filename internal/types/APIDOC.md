# Package types

Package types (import "gitee.com/MM-Q/fastlog/internal/types")

## Constants

| 常量名                | 类型          | 值/说明                                                                 |
|-----------------------|---------------|-------------------------------------------------------------------------|
| `DefaultMaxFileSize`  | `int`         | 10，默认最大文件大小（单位：MB）                                        |
| `DefaultTimeFormat`   | `string`      | "2006-01-02T15:04:05"，默认时间格式                                    |
| `DefaultCallerDepth`  | `int`         | 3，默认调用信息层数（0表示当前调用，1表示调用者，依此类推）             |
| `DefaultMaxBufferSize`| `int`         | 64 * 1024，默认最大缓冲区大小（单位：字节，即64KB）                    |
| `DefaultMaxWriteCount`| `int`         | 500，默认最大写入次数                                                  |
| `DefaultFlushInterval`| `time.Duration`| 1 * time.Second，默认最大刷新间隔（即1秒）                             |

---

## Variables

### LogLevelPaddedStringMap
预构建的日志级别到字符串的映射表（带填充，用于文本格式化）

```go
var LogLevelPaddedStringMap = map[LogLevel]string{
	DEBUG_Mask: "DEBUG",
	INFO_Mask:  "INFO ",
	WARN_Mask:  "WARN ",
	ERROR_Mask: "ERROR",
	FATAL_Mask: "FATAL",
	NONE_Mask:  "NONE ",
	DEBUG:      "DEBUG",
	INFO:       "INFO ",
	WARN:       "WARN ",
	ERROR:      "ERROR",
}
```

### LogLevelStringMap
预构建的日志级别到字符串的映射表（不带填充，用于JSON序列化）

```go
var LogLevelStringMap = map[LogLevel]string{
	DEBUG_Mask: "DEBUG",
	INFO_Mask:  "INFO",
	WARN_Mask:  "WARN",
	ERROR_Mask: "ERROR",
	FATAL_Mask: "FATAL",
	NONE_Mask:  "NONE",
	DEBUG:      "DEBUG",
	INFO:       "INFO",
	WARN:       "WARN",
	ERROR:      "ERROR",
}
```

---

## Functions

### GetCachedTimestamp
GetCachedTimestamp 获取缓存的时间戳，读写锁优化版本

```go
func GetCachedTimestamp() string
```

**性能特点：**
- 快路径：原子操作检查 + 读锁保护
- 慢路径：写锁保护更新操作
- 多读者并发，单写者独占
- 无unsafe操作，完全内存安全

**返回值：**
- string: 格式化的时间戳字符串 "2006-01-02 15:04:05"

### GetCallerInfo
GetCallerInfo 获取调用者的信息（优化版本，使用文件名缓存）

```go
func GetCallerInfo(skip int) []byte
```

**参数：**
- skip: 跳过的调用层数（通常设置为1或2，具体取决于调用链的深度）

**返回值：**
- []byte: 调用者的信息，格式为 "file:function:line"

### LogLevelToPaddedString
LogLevelToPaddedString 将 LogLevel 转换为带填充的字符串（用于文本格式化）

```go
func LogLevelToPaddedString(level LogLevel) string
```

**参数：**
- level: 要转换的日志级别

**返回值：**
- string: 对应的带填充的日志级别字符串（7个字符），如果 level 无效，则返回 "UNKNOWN"

### ShouldLog
ShouldLog 检查是否应该记录指定级别的日志，使用位运算优化日志级别比较，提高判断性能

```go
func ShouldLog(logLevel, minLevel LogLevel) bool
```

**参数：**
- logLevel: 当前要记录的日志级别（基本级别，如DEBUG_Mask）
- minLevel: 配置的最低记录级别（组合级别，如DEBUG, INFO等）

**返回值：**
- bool: 如果当前级别应该记录，则返回 true；否则返回 false

---

## Types

### LogFormatType
LogFormatType 日志格式选项

```go
type LogFormatType int
```

**格式说明：**
- Detailed: 详细格式
- Json: json格式
- Timestamp: 时间格式
- KVFmt: 键值对格式
- LogFmt: logfmt格式
- Custom: 自定义格式

#### LogFormatType 常量定义
```go
const (
	Def       LogFormatType = iota // 默认格式
	Json                           // json格式
	Timestamp                      // 时间格式
	KVFmt                          // KVFmt 键值对格式
	LogFmt                         // logfmt格式
	Custom                         // 自定义格式
)
```

#### String
String 将 LogFormatType 转换为对应的字符串

```go
func (f LogFormatType) String() string
```

### LogLevel
LogLevel 定义为位掩码类型，每一位代表一个日志级别

```go
type LogLevel uint8
```

**支持的日志级别：**
- DEBUG: 调试级别
- INFO: 信息级别
- WARN: 警告级别
- ERROR: 错误级别
- FATAL: 致命错误级别
- NONE: 不记录任何日志

#### LogLevel 常量定义
```go
const (
	DEBUG_Mask LogLevel = 1 << iota // 1 表示Debug级别
	INFO_Mask                       // 2  表示Info级别
	WARN_Mask                       // 4  表示Warn级别
	ERROR_Mask                      // 8  表示Error级别
	FATAL_Mask                      // 16 表示Fatal级别
	NONE_Mask  LogLevel = 0         // 0 表示不启用任何日志

	// 预定义的日志级别组合
	DEBUG LogLevel = DEBUG_Mask | INFO_Mask | WARN_Mask | ERROR_Mask | FATAL_Mask // Debug及以上级别
	INFO  LogLevel = INFO_Mask | WARN_Mask | ERROR_Mask | FATAL_Mask              // Info及以上级别
	WARN  LogLevel = WARN_Mask | ERROR_Mask | FATAL_Mask                          // Warn及以上级别
	ERROR LogLevel = ERROR_Mask | FATAL_Mask                                      // Error及以上级别
	FATAL LogLevel = FATAL_Mask                                                   // Fatal级别
	NONE  LogLevel = NONE_Mask
)
```

#### MarshalJSON
MarshalJSON 将 LogLevel 转换为 JSON 字符串

```go
func (l LogLevel) MarshalJSON() ([]byte, error)
```

**返回值：**
- []byte: 包含日志级别的 JSON 字符串（带双引号）
- error: 如果转换过程中发生错误，返回非 nil 错误；否则返回 nil

#### String
String 将 LogLevel 转换为字符串

```go
func (l LogLevel) String() string
```

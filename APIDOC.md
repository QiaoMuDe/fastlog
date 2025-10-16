# Package fastlog

导入路径：`import "gitee.com/MM-Q/fastlog"`


## 1. 常量 (CONSTANTS)

### 1.1 日志级别常量
定义日志记录的不同级别，用于控制日志输出的详细程度。

```go
const (
    DEBUG = types.DEBUG // 调试级别
    INFO  = types.INFO  // 信息级别
    WARN  = types.WARN  // 警告级别
    ERROR = types.ERROR // 错误级别
    FATAL = types.FATAL // 致命错误级别
    NONE  = types.NONE  // 不记录任何日志级别
)
```

### 1.2 日志格式常量
定义日志输出的不同格式，满足不同场景下的日志展示需求。

```go
const (
    Def       = types.Def       // 默认格式
    Json      = types.Json      // json格式
    Timestamp = types.Timestamp // 时间格式
    KVFmt     = types.KVFmt     // 键值对格式
    LogFmt    = types.LogFmt    // 日志格式
    Custom    = types.Custom    // 自定义格式
)
```


## 2. 变量 (VARIABLES)

### 2.1 Bool
添加布尔类型的日志字段。

```go
var Bool = flog.Bool
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为布尔值。
- **返回值**：`*Field`，一个指向 Field 实例的指针。

---

### 2.2 ConsoleConfig
创建控制台环境专用的 `FastLogConfig` 实例。

```go
var ConsoleConfig = config.ConsoleConfig
```

- **返回值**：`*FastLogConfig`，一个指向 FastLogConfig 实例的指针。
- **特性**：
  - 禁用文件输出。
  - 设置日志级别为 DEBUG。

---

### 2.3 DevConfig
创建开发环境专用的 `FastLogConfig` 实例。

```go
var DevConfig = config.DevConfig
```

- **参数**：
  - logDirName：日志目录名。
  - logFileName：日志文件名。
- **返回值**：`*FastLogConfig`，一个指向 FastLogConfig 实例的指针。
- **特性**：
  - 启用详细日志格式。
  - 设置日志级别为 DEBUG。
  - 最大日志文件保留数量为 5 个。
  - 最大日志文件保留天数为 7 天。

---

### 2.4 Duration
添加 `time.Duration` 类型的日志字段，用于记录时间间隔。

```go
var Duration = flog.Duration
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为 `time.Duration` 类型。
- **返回值**：`*Field`，一个指向 Field 实例的指针。

---

### 2.5 Error
添加 `error` 类型的日志字段，用于记录错误信息。

```go
var Error = flog.Error
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为 `error` 类型。
- **返回值**：`*Field`，一个指向 Field 实例的指针。

---

### 2.6 Float64
添加 64 位浮点类型的日志字段。

```go
var Float64 = flog.Float64
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为 64 位浮点数。
- **返回值**：`*Field`，一个指向 Field 实例的指针。

---

### 2.7 Int
添加整数类型的日志字段。

```go
var Int = flog.Int
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为整数。
- **返回值**：`*Field`，一个指向 Field 实例的指针。

---

### 2.8 Int64
添加 64 位整数类型的日志字段。

```go
var Int64 = flog.Int64
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为 64 位整数。
- **返回值**：`*Field`，一个指向 Field 实例的指针。

---

### 2.9 NewCfg
`NewFastLogConfig` 的简写，用于快速创建 `FastLogConfig` 实例。

```go
var NewCfg = config.NewFastLogConfig
```

- **参数**：
  - logDirName：日志目录名称，默认为 "applogs"。
  - logFileName：日志文件名称，默认为 "app.log"。
- **返回值**：`*FastLogConfig`，一个指向 FastLogConfig 实例的指针。

---

### 2.10 NewFLog
创建日志记录器核心实例 `FLog`，用于实际日志输出。

```go
var NewFLog = flog.NewFLog
```

- **参数**：
  - config：指向 `FastLogConfig` 实例的指针，用于配置日志记录器。
- **返回值**：`*FLog`，一个指向 FLog 实例的指针。

---

### 2.11 NewFastLogConfig
创建标准的 `FastLogConfig` 实例，用于配置日志记录器的基础参数。

```go
var NewFastLogConfig = config.NewFastLogConfig
```

- **参数**：
  - logDirName：日志目录名称，默认为 "applogs"。
  - logFileName：日志文件名称，默认为 "app.log"。
- **返回值**：`*FastLogConfig`，一个指向 FastLogConfig 实例的指针。

---

### 2.12 ProdConfig
创建生产环境专用的 `FastLogConfig` 实例，侧重日志存储效率和归档策略。

```go
var ProdConfig = config.ProdConfig
```

- **参数**：
  - logDirName：日志目录名。
  - logFileName：日志文件名。
- **返回值**：`*FastLogConfig`，一个指向 FastLogConfig 实例的指针。
- **特性**：
  - 启用日志文件压缩。
  - 禁用控制台输出。
  - 最大日志文件保留天数为 30 天。
  - 最大日志文件保留数量为 24 个。

---

### 2.13 String
添加字符串类型的日志字段。

```go
var String = flog.String
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为字符串。
- **返回值**：`*Field`，一个指向 Field 实例的指针。

---

### 2.14 Time
添加 `time.Time` 类型的日志字段，用于记录具体时间点。

```go
var Time = flog.Time
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为 `time.Time` 类型。
- **返回值**：`*Field`，一个指向 Field 实例的指针。

---

### 2.15 Uint16
添加 16 位无符号整数类型的日志字段。

```go
var Uint16 = flog.Uint16
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为 16 位无符号整数。
- **返回值**：`*Field`，一个指向 Field 实例的指针。

---

### 2.16 Uint32
添加 32 位无符号整数类型的日志字段。

```go
var Uint32 = flog.Uint32
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为 32 位无符号整数。
- **返回值**：`*Field`，一个指向 Field 实例的指针。

---

### 2.17 Uint64
添加 64 位无符号整数类型的日志字段。

```go
var Uint64 = flog.Uint64
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为 64 位无符号整数。
- **返回值**：`*Field`，一个指向 Field 实例的指针。

---

### 2.18 Uint8
添加 8 位无符号整数类型的日志字段。

```go
var Uint8 = flog.Uint8
```

- **参数**：
  - key：字段的键名，不能为空字符串。
  - value：字段值，必须为 8 位无符号整数。
- **返回值**：`*Field`，一个指向 Field 实例的指针。


## 3. 类型 (TYPES)

### 3.1 FLog
日志记录器核心类型，支持键值对风格和标准 `fmt` 风格的日志输出，提供丰富的配置选项。

```go
type FLog = flog.FLog
```

- **核心能力**：
  - 支持日志级别控制（DEBUG/INFO/WARN 等）。
  - 支持多种输出格式（Json/KVFmt/LogFmt 等）。
  - 支持日志轮转（按文件数量、按保留天数）。

---

### 3.2 FastLogConfig
日志记录器的配置结构体，用于定义日志的存储、格式、级别等基础参数。

```go
type FastLogConfig = config.FastLogConfig
```

- **关联变量**：`ConsoleConfig`、`DevConfig`、`ProdConfig`、`NewFastLogConfig` 等均用于创建该类型实例。

---

### 3.3 LogFormatType
日志输出格式的枚举类型，定义了所有支持的日志格式选项。

```go
type LogFormatType = types.LogFormatType
```

- **支持的格式**：
  - Detailed：详细格式（包含完整上下文信息）。
  - Json：JSON 格式（便于机器解析）。
  - Timestamp：时间格式（以时间为核心的简化格式）。
  - KVFmt：键值对格式（`key=value` 形式，可读性强）。
  - LogFmt：标准 LogFmt 格式（兼容主流日志收集工具）。
  - Custom：自定义格式（允许用户自定义输出模板）。

---

### 3.4 LogLevel
日志级别的枚举类型，采用位掩码设计，支持多级别组合。

```go
type LogLevel = types.LogLevel
```

- **支持的级别**：
  - DEBUG：调试级别（开发阶段调试信息）。
  - INFO：信息级别（正常业务流程信息）。
  - WARN：警告级别（非致命的异常情况）。
  - ERROR：错误级别（业务逻辑错误，需关注）。
  - FATAL：致命错误级别（程序无法继续运行，会终止）。
  - NONE：无日志级别（不输出任何日志）。

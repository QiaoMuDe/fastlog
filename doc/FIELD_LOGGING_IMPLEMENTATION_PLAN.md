# 字段日志功能实现方案

## 需求分析

实现一个可以按字段添加内容的日志功能，支持：
1. 通过 `AddField(key, value)` 方法添加字段
2. 支持链式调用
3. 支持多次调用添加多个字段
4. 最终通过 `Log(message)` 方法输出日志

## 设计思路

在现有日志库基础上扩展，不破坏现有API，保持向后兼容性。

## 核心实现

### 1. 添加FieldLogger结构体

```go
// FieldLogger 字段日志构建器
type FieldLogger struct {
    logger *FastLog       // 关联的FastLog实例
    level  LogLevel       // 日志级别
    fields map[string]interface{} // 字段存储
}
```

### 2. 在FastLog中添加创建FieldLogger的方法

```go
// WithFields 创建一个字段日志构建器
func (f *FastLog) WithFields(level LogLevel) *FieldLogger {
    return &FieldLogger{
        logger: f,
        level:  level,
        fields: make(map[string]interface{}),
    }
}
```

### 3. 实现FieldLogger的方法

```go
// AddField 添加字段，支持链式调用
func (fl *FieldLogger) AddField(key string, value interface{}) *FieldLogger {
    fl.fields[key] = value
    return fl
}

// Log 输出日志
func (fl *FieldLogger) Log(message string) {
    if fl.logger == nil || fl.fields == nil {
        return
    }
    
    // 格式化字段信息
    formattedFields := fl.formatFields()
    
    // 调用现有日志处理流程
    fl.logger.processLog(fl.level, message+" "+formattedFields)
}
```

### 4. 字段格式化方法

提供多种字段格式化选项：

```go
// formatFields 格式化字段为键值对字符串
func (fl *FieldLogger) formatFields() string {
    if len(fl.fields) == 0 {
        return ""
    }
    
    var sb strings.Builder
    for key, value := range fl.fields {
        if sb.Len() > 0 {
            sb.WriteString(" ")
        }
        sb.WriteString(key)
        sb.WriteString("=")
        sb.WriteString(fmt.Sprintf("%v", value))
    }
    
    return sb.String()
}
```

## 使用示例

```go
// 创建日志实例
logger := fastlog.NewFastLog(config)

// 使用字段日志
logger.WithFields(fastlog.INFO).
    AddField("user_id", 12345).
    AddField("action", "login").
    AddField("ip", "192.168.1.1").
    Log("User login event")

// 输出示例：
// 2025-01-15 10:30:45 | INFO    | User login event user_id=12345 action=login ip=192.168.1.1
```

## 文件修改计划

1. **fastlog.go** - 添加WithFields方法
2. **internal.go** - 添加FieldLogger结构体及其实现
3. **types.go** - (可选)如果需要在其他地方引用FieldLogger，可以在这里定义

## 实现步骤

1. 在internal.go中实现FieldLogger结构体及相关方法
2. 在fastlog.go中为FastLog添加WithFields方法
3. 添加测试用例验证功能
4. 更新文档说明

## 格式化选项

支持多种字段输出格式：

1. **键值对格式**：`key1=value1 key2=value2`
2. **JSON格式**：`{"key1":"value1","key2":"value2"}`
3. **结构化格式**：`key1="value1" key2="value2"`

可以通过配置选择不同的格式化方式。
# BufferEnabled 字段设计方案

## 背景

当前 `newFileWriter()` 方法**强制**使用 `logrotatex.NewBufferedWriter` 包装 `LogRotateX`，
没有提供禁用缓冲的选项。在某些场景下（如开发调试、高可靠性要求），用户希望直接写入文件，
避免缓冲带来的数据丢失风险或延迟问题。

## 目标

- 提供灵活的控制机制，允许用户选择是否启用缓冲写入
- 保持向后兼容，默认行为不变
- 场景化配置合理：开发环境禁用，生产环境启用

## 方案概述

新增 `BufferEnabled bool` 字段控制是否启用缓冲写入：
- `true`（默认）：使用 `BufferedWriter`（LogRotateX + BufCfg）
- `false`：直接使用 `LogRotateX`，无缓冲，立即落盘

## 详细设计

### 1. Config 结构体修改

```go
type Config struct {
    // ... 其他现有字段
    
    // BufferEnabled 是否启用缓冲写入
    // true:  使用 BufferedWriter（默认）
    // false: 直接使用 LogRotateX，无缓冲，立即落盘
    // 
    // 使用建议：
    //   - 开发环境：设为 false，立即看到日志，便于调试
    //   - 生产环境：设为 true（默认），提升写入性能
    //   - 高可靠性场景：设为 false，避免数据丢失风险
    BufferEnabled bool
    
    // ... 其他现有字段
}
```

### 2. 默认值设置

#### NewConfig() - 通用默认配置
```go
func NewConfig(logPath string) *Config {
    return &Config{
        // ... 其他配置
        BufferEnabled: true,  // 默认启用缓冲，保持向后兼容
        MaxBufferSize: 256 * 1024,
        SyncInterval:  1 * time.Second,
        // ...
    }
}
```

#### Dev() - 开发环境配置
```go
func Dev(logPath string) *Config {
    cfg := NewConfig(logPath)
    cfg.BufferEnabled = false  // 开发环境禁用缓冲，立即写入
    // ... 其他配置
    return cfg
}
```

#### Prod() - 生产环境配置
```go
func Prod(logPath string) *Config {
    cfg := NewConfig(logPath)
    cfg.BufferEnabled = true   // 生产环境启用缓冲（已是默认值，显式设置更清晰）
    // ... 其他配置
    return cfg
}
```

#### Console() / Docker() - 控制台/容器配置
```go
// 这两个配置 OutputFile = false，不涉及文件写入
// BufferEnabled 不影响它们的行为
```

### 3. newFileWriter() 方法修改

```go
func (c *Config) newFileWriter() io.WriteCloser {
    // 创建日志切割器（核心写入器）
    logger := &logrotatex.LogRotateX{
        LogFilePath:   c.LogPath,
        MaxSize:       c.MaxSize,
        MaxAge:        c.MaxAge,
        MaxFiles:      c.MaxFiles,
        Compress:      c.Compress,
        LocalTime:     c.LocalTime,
        Async:         c.Async,
        DateDirLayout: c.DateDirLayout,
        RotateByDay:   c.RotateByDay,
        CompressType:  c.CompressType,
    }

    // 如果禁用缓冲，直接返回 LogRotateX
    if !c.BufferEnabled {
        return logger
    }

    // 启用缓冲，包装 BufferedWriter
    bufCfg := &logrotatex.BufCfg{
        SyncInterval:  c.SyncInterval,
        MaxBufferSize: c.MaxBufferSize,
    }
    return logrotatex.NewBufferedWriter(logger, bufCfg)
}
```

### 4. 场景化配置表格更新

| 函数 | BufferEnabled | 说明 |
|------|---------------|------|
| `NewConfig(path)` | ✅ true | 默认启用缓冲 |
| `Default()` | ✅ true | 默认启用缓冲 |
| `Dev(path)` | ❌ false | 开发环境禁用，立即写入 |
| `Prod(path)` | ✅ true | 生产环境启用，性能优先 |
| `Console()` | N/A | 无文件输出，不涉及 |
| `Docker()` | N/A | 无文件输出，不涉及 |

## 使用示例

### 示例 1: 开发环境（禁用缓冲）

```go
package main

import (
    "gitee.com/MM-Q/fastlog"
)

func main() {
    // Dev() 自动禁用缓冲
    logger := fastlog.New(fastlog.Dev("logs/dev.log"))
    defer func() { _ = logger.Close() }()

    logger.Info("这条日志会立即写入文件")
    // 无需等待 SyncInterval，立即可在文件中看到
}
```

### 示例 2: 生产环境（启用缓冲）

```go
package main

import (
    "gitee.com/MM-Q/fastlog"
)

func main() {
    // Prod() 启用缓冲（默认）
    logger := fastlog.New(fastlog.Prod("/var/log/app.log"))
    defer func() { _ = logger.Close() }()

    logger.Info("这条日志会先进入缓冲区")
    // 缓冲区满或到达 SyncInterval 才会写入磁盘
}
```

### 示例 3: 自定义配置（显式控制）

```go
package main

import (
    "gitee.com/MM-Q/fastlog"
)

func main() {
    cfg := fastlog.NewConfig("logs/app.log")
    
    // 显式禁用缓冲（高可靠性场景）
    cfg.BufferEnabled = false
    
    // 或者基于 NewConfig 修改其他配置
    cfg.Level = fastlog.DEBUG
    cfg.Caller = true
    
    logger := fastlog.New(cfg)
    defer func() { _ = logger.Close() }()

    logger.Info("无缓冲，立即落盘")
}
```

## 影响分析

### 向后兼容性
- ✅ **完全兼容**：NewConfig() 默认 `BufferEnabled: true`，现有代码行为不变
- ✅ **API 不变**：只是新增字段，没有修改现有函数签名

### 性能影响
- 启用缓冲：批量写入，性能更好，但有延迟
- 禁用缓冲：立即落盘，延迟低，但 I/O 更频繁

### 数据安全性
- 启用缓冲：程序崩溃时可能丢失缓冲区数据（最多 MaxBufferSize）
- 禁用缓冲：每条日志立即落盘，数据安全性更高

## 实现步骤

1. **config.go**
   - [ ] Config 结构体添加 BufferEnabled 字段
   - [ ] NewConfig() 设置默认值 true
   - [ ] Dev() 设置 BufferEnabled = false
   - [ ] Prod() 显式设置 BufferEnabled = true（可选，为了清晰）
   - [ ] newFileWriter() 添加条件判断

2. **README.md**
   - [ ] 核心特性表格添加缓冲控制说明
   - [ ] 场景化配置表格添加 BufferEnabled 列
   - [ ] 使用指南添加缓冲控制章节

3. **fastlog-skill/SKILL.md**
   - [ ] 更新场景化配置说明
   - [ ] 添加 BufferEnabled 使用示例

4. **测试**
   - [ ] 验证 Dev() 禁用缓冲
   - [ ] 验证 Prod() 启用缓冲
   - [ ] 验证自定义配置生效

## 决策记录

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 控制方式 | `BufferEnabled bool` | 语义清晰，避免歧义 |
| 默认值 | `true`（启用缓冲） | 保持向后兼容 |
| Dev() 配置 | `false`（禁用缓冲） | 开发环境需要立即看到日志 |
| Prod() 配置 | `true`（启用缓冲） | 生产环境性能优先 |

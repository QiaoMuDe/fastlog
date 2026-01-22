# Cache Package API Documentation

## 概述

`cache` 包提供了智能缓存键和ANSI序列缓存功能，用于优化ColorLib的性能。通过预缓存常用的ANSI颜色序列组合，显著减少字符串构建开销。

## 核心组件

### CacheKey

智能缓存键，使用32位整数压缩存储样式信息。

#### 位布局设计
```
31-16位: 保留给未来功能扩展
15位:    粗体标志位 (Bold)
14位:    下划线标志位 (Underline)  
13位:    闪烁标志位 (Blink)
12-8位:  保留给未来样式扩展
7-0位:   颜色代码 (Color Code)
```

#### 主要方法

- `BuildCacheKey(colorCode int, bold, underline, blink bool) CacheKey`
- `Parse() (colorCode int, bold, underline, blink bool)`
- `String() string`
- `IsValid() bool`

### ANSICache

ANSI序列缓存管理器，提供高效的序列缓存功能。

#### 主要方法

- `NewANSICache() *ANSICache`
- `GetANSI(colorCode int, bold, underline, blink bool) string`
- `GetCacheSize() int`
- `Clear()`

## 使用示例

### 基础用法

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/colorlib/internal/cache"
)

func main() {
    // 创建缓存管理器
    ansiCache := cache.NewANSICache()
    
    // 获取ANSI序列
    redBold := ansiCache.GetANSI(31, true, false, false)
    fmt.Printf("红色粗体序列: %s\n", redBold)
    
    fmt.Printf("缓存大小: %d\n", ansiCache.GetCacheSize())
}
```

### 缓存键操作

```go
// 构建缓存键
key := cache.BuildCacheKey(31, true, false, false)

// 解析缓存键
color, bold, underline, blink := key.Parse()
fmt.Printf("颜色: %d, 粗体: %v, 下划线: %v, 闪烁: %v\n", 
           color, bold, underline, blink)

// 检查有效性
if key.IsValid() {
    fmt.Println("缓存键有效")
}
```

## 性能特性

### 内存效率
- 缓存键仅占用4字节
- 相比结构体键减少87.5%内存占用
- 支持高效的哈希查找

### 查找性能
- O(1) 哈希查找时间复杂度
- 位运算比较，单CPU指令完成

### 预热策略
- 自动预热25+常用颜色组合
- 基于实际使用模式优化
- 支持日志级别常用组合

## 扩展性设计

### 位掩码常量
```go
const (
    BoldMask      = 1 << 15 // 粗体
    UnderlineMask = 1 << 14 // 下划线
    BlinkMask     = 1 << 13 // 闪烁
    ColorMask     = 0xFF    // 颜色
)
```

### 未来扩展预留
- 第12-8位: 预留给新样式 (斜体、删除线等)
- 第31-16位: 预留给新功能 (背景色、RGB等)

## 线程安全

- 使用 `sync.Map` 实现无锁缓存
- 支持高并发读写操作

## 最佳实践

1. **单例使用**: 建议在ColorLib中使用单个ANSICache实例
2. **预热优化**: 根据应用场景自定义预热组合
3. **内存管理**: 在长期运行的应用中考虑定期清理
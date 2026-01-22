# ColorLib - Go 语言高性能彩色终端输出库

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.24-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/gitee.com/MM-Q/colorlib)](https://goreportcard.com/report/gitee.com/MM-Q/colorlib)

## 📖 简介

`ColorLib` 是一个功能强大、高性能的 Go 语言终端颜色输出库。它提供了丰富的颜色输出功能，支持标准颜色和亮色系列，具备样式设置（粗体、下划线、闪烁）、链式调用、自定义输出接口、线程安全等特性，专为提升命令行程序的用户体验而设计。

## ✨ 核心特性

- 🎨 **丰富的颜色支持**：16种标准颜色 + 7种亮色，满足各种场景需求
- 🔗 **链式调用**：支持流畅的链式API，代码更简洁优雅
- 🎯 **多种输出方式**：直接打印、格式化输出、字符串返回
- 🏷️ **日志级别支持**：内置 Debug、Info、Ok、Warn、Error 等级别
- 🎭 **样式控制**：支持粗体、下划线、闪烁等文本效果
- 🔒 **线程安全**：全局实例和并发操作完全安全
- 📝 **自定义输出**：支持输出到文件、缓冲区等任意 io.Writer
- ⚡ **高性能**：内置对象池和智能缓存，减少内存分配
- 🧪 **完整测试**：99+ 测试用例，覆盖率 > 90%

## 🚀 快速开始

### 安装

```bash
go get gitee.com/MM-Q/colorlib
```

### 基础用法

```go
package main

import (
    "gitee.com/MM-Q/colorlib"
)

func main() {
    // 使用全局实例（推荐）
    cl := colorlib.GetCL()
    
    // 基础颜色输出
    cl.Red("这是红色文本")
    cl.Green("这是绿色文本")
    cl.Blue("这是蓝色文本")
    
    // 格式化输出
    cl.Yellowf("用户 %s 登录成功，时间：%s\n", "张三", "2024-01-01")
    
    // 日志级别输出
    cl.PrintInfo("系统启动中...")
    cl.PrintOk("启动成功！")
    cl.PrintWarn("内存使用率较高")
    cl.PrintError("连接数据库失败")
    
    // 返回带颜色的字符串
    coloredMsg := cl.Sgreen("成功处理 100 条记录")
    fmt.Println(coloredMsg)
}
```

## 🎨 颜色支持

### 标准颜色

| 颜色名称 | 方法名 | 颜色代码 | 示例 |
|---------|--------|----------|------|
| 黑色 | `Black()` | 30 | `cl.Black("黑色文本")` |
| 红色 | `Red()` | 31 | `cl.Red("红色文本")` |
| 绿色 | `Green()` | 32 | `cl.Green("绿色文本")` |
| 黄色 | `Yellow()` | 33 | `cl.Yellow("黄色文本")` |
| 蓝色 | `Blue()` | 34 | `cl.Blue("蓝色文本")` |
| 品红色 | `Magenta()` | 35 | `cl.Magenta("品红色文本")` |
| 青色 | `Cyan()` | 36 | `cl.Cyan("青色文本")` |
| 白色 | `White()` | 37 | `cl.White("白色文本")` |
| 灰色 | `Gray()` | 90 | `cl.Gray("灰色文本")` |

### 亮色系列

| 颜色名称 | 方法名 | 颜色代码 | 示例 |
|---------|--------|----------|------|
| 亮红色 | `BrightRed()` | 91 | `cl.BrightRed("亮红色文本")` |
| 亮绿色 | `BrightGreen()` | 92 | `cl.BrightGreen("亮绿色文本")` |
| 亮黄色 | `BrightYellow()` | 93 | `cl.BrightYellow("亮黄色文本")` |
| 亮蓝色 | `BrightBlue()` | 94 | `cl.BrightBlue("亮蓝色文本")` |
| 亮品红色 | `BrightMagenta()` | 95 | `cl.BrightMagenta("亮品红色文本")` |
| 亮青色 | `BrightCyan()` | 96 | `cl.BrightCyan("亮青色文本")` |
| 亮白色 | `BrightWhite()` | 97 | `cl.BrightWhite("亮白色文本")` |

## 🏷️ 日志级别

| 级别 | 方法名 | 前缀 | 颜色 | 使用场景 |
|------|--------|------|------|----------|
| Debug | `PrintDebug()` | `debug: ` | 品红色 | 调试信息 |
| Info | `PrintInfo()` | `info: ` | 蓝色 | 一般信息 |
| Ok | `PrintOk()` | `ok: ` | 绿色 | 成功操作 |
| Warn | `PrintWarn()` | `warn: ` | 黄色 | 警告信息 |
| Error | `PrintError()` | `error: ` | 红色 | 错误信息 |

```go
cl := colorlib.GetCL()

cl.PrintDebug("调试信息：变量值为", value)
cl.PrintInfo("正在处理用户请求...")
cl.PrintOk("数据保存成功")
cl.PrintWarn("磁盘空间不足")
cl.PrintError("网络连接超时")

// 格式化版本
cl.PrintDebugf("用户ID: %d, 状态: %s", userID, status)
cl.PrintInfof("处理进度: %d%%", progress)
```

## 🎭 样式控制

### 基础样式设置

```go
cl := colorlib.NewColorLib()

// 设置样式
cl.SetColor(true)      // 启用颜色
cl.SetBold(true)       // 启用粗体
cl.SetUnderline(true)  // 启用下划线
cl.SetBlink(true)      // 启用闪烁

cl.Red("带样式的红色文本")
```

### 链式调用

```go
cl := colorlib.NewColorLib()

// 链式设置样式
cl.WithColor(true).
   WithBold(true).
   WithUnderline(true).
   Red("链式调用的红色粗体下划线文本")
```

### 禁用颜色输出

```go
cl := colorlib.NewColorLib()
cl.SetColor(false)  // 禁用颜色，适用于日志文件输出
cl.Red("这将显示为普通文本")
```

## 📤 输出方式

### 1. 直接打印（带换行）

```go
cl.Red("直接打印红色文本")           // 输出后自动换行
cl.Green("直接打印绿色文本")         // 输出后自动换行
```

### 2. 格式化打印（不换行）

```go
cl.Redf("用户: %s", username)       // 格式化输出，不自动换行
cl.Greenf("状态: %s", status)       // 需要手动添加 \n
```

### 3. 返回字符串

```go
redText := cl.Sred("红色字符串")     // 返回带颜色的字符串
greenText := cl.Sgreen("绿色字符串") // 返回带颜色的字符串
fmt.Println(redText, greenText)
```

### 4. 格式化返回字符串

```go
coloredMsg := cl.Sredf("错误代码: %d", errorCode)
log.Println(coloredMsg)  // 可以传递给其他日志库
```

## 🔧 高级用法

### 自定义输出接口

```go
// 输出到文件
file, _ := os.Create("colored_log.txt")
defer file.Close()
cl := colorlib.NewColorLibWithWriter(file)
cl.Red("这将写入文件")

// 输出到缓冲区
var buf bytes.Buffer
cl := colorlib.NewColorLibWithWriter(&buf)
cl.Green("这将写入缓冲区")
fmt.Println(buf.String())

// 使用 WithWriter 创建新实例
cl1 := colorlib.GetCL()
cl2 := cl1.WithWriter(os.Stderr)  // 输出到标准错误
cl2.Red("错误信息")
```

### 全局实例 vs 自定义实例

```go
// 方式1: 使用全局实例（推荐）
cl := colorlib.GetCL()  // 线程安全的单例
cl.Red("使用全局实例")

// 方式2: 创建新实例
cl := colorlib.NewColorLib()  // 或者 colorlib.New()
cl.Red("使用新实例")

// 方式3: 指定输出接口
cl := colorlib.NewColorLibWithWriter(os.Stderr)
cl.Red("输出到标准错误")
```

### 并发使用

```go
cl := colorlib.GetCL()  // 全局实例是线程安全的

// 在多个 goroutine 中安全使用
go func() {
    cl.Red("Goroutine 1")
}()

go func() {
    cl.Green("Goroutine 2")
}()
```

## 📋 完整 API 参考

### 构造函数

| 函数名 | 描述 |
|--------|------|
| `GetCL()` | 获取全局单例实例（线程安全） |
| `NewColorLib()` | 创建新实例（输出到 stdout） |
| `New()` | `NewColorLib()` 的别名 |
| `NewColorLibWithWriter(io.Writer)` | 创建指定输出接口的实例 |

### 样式设置

| 方法名 | 描述 |
|--------|------|
| `SetColor(bool)` | 设置是否启用颜色 |
| `SetBold(bool)` | 设置是否启用粗体 |
| `SetUnderline(bool)` | 设置是否启用下划线 |
| `SetBlink(bool)` | 设置是否启用闪烁 |
| `WithColor(bool)` | 链式设置颜色（返回自身） |
| `WithBold(bool)` | 链式设置粗体（返回自身） |
| `WithUnderline(bool)` | 链式设置下划线（返回自身） |
| `WithBlink(bool)` | 链式设置闪烁（返回自身） |
| `WithWriter(io.Writer)` | 创建新的输出接口实例 |

### 颜色方法命名规则

| 前缀/后缀 | 说明 | 示例 |
|-----------|------|------|
| 无前缀 | 直接打印（带换行） | `Red("text")` |
| `f` 后缀 | 格式化打印（不换行） | `Redf("user: %s", name)` |
| `S` 前缀 | 返回字符串 | `Sred("text")` |
| `S` + `f` | 格式化返回字符串 | `Sredf("user: %s", name)` |
| `Bright` 前缀 | 亮色版本 | `BrightRed("text")` |

## 🎯 使用场景

### CLI 工具

```go
func main() {
    cl := colorlib.GetCL()
    
    cl.PrintInfo("正在初始化...")
    
    if err := initialize(); err != nil {
        cl.PrintError("初始化失败:", err)
        os.Exit(1)
    }
    
    cl.PrintOk("初始化完成")
    cl.Green("欢迎使用 MyTool v1.0.0")
}
```

### 日志系统

```go
type Logger struct {
    cl *colorlib.ColorLib
}

func NewLogger() *Logger {
    return &Logger{cl: colorlib.GetCL()}
}

func (l *Logger) Info(msg string) {
    l.cl.PrintInfo(msg)
}

func (l *Logger) Error(msg string) {
    l.cl.PrintError(msg)
}
```

### 测试输出

```go
func TestSomething(t *testing.T) {
    cl := colorlib.GetCL()
    
    cl.PrintInfo("开始测试...")
    
    if result := doSomething(); result {
        cl.PrintOk("测试通过")
    } else {
        cl.PrintError("测试失败")
        t.Fail()
    }
}
```

## 🔧 配置建议

### 生产环境

```go
cl := colorlib.NewColorLib()

// 根据环境变量决定是否启用颜色
if os.Getenv("NO_COLOR") != "" {
    cl.SetColor(false)
}

// 输出到日志文件时禁用颜色
if isLogFile {
    cl.SetColor(false)
}
```

### 开发环境

```go
cl := colorlib.GetCL()
cl.WithBold(true).WithUnderline(true)  // 开发时使用更明显的样式
```

## 📊 性能特性

- **对象池技术**：内置 `strings.Builder` 对象池，减少内存分配
- **智能缓存**：ANSI 序列缓存，避免重复构建
- **零拷贝**：优化的字符串操作，减少不必要的内存复制
- **并发安全**：使用原子操作，无锁设计

## 🧪 测试

运行所有测试：

```bash
go test ./...
```

运行竞态检测：

```bash
go test -race ./...
```

查看测试覆盖率：

```bash
go test -cover ./...
```

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 联系方式

- 项目地址：[https://gitee.com/MM-Q/colorlib](https://gitee.com/MM-Q/colorlib)
- 问题反馈：[Issues](https://gitee.com/MM-Q/colorlib/issues)

---

⭐ 如果这个项目对你有帮助，请给个 Star 支持一下！
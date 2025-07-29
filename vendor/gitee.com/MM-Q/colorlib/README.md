# ColorLib - Go 语言的彩色终端输出库

## 用途和特点

`ColorLib` 是一个用于在 Go 语言中实现彩色终端输出的库。它提供了丰富的接口，用于打印和返回带有颜色的文本，并支持自定义颜色和日志级别。该库的主要特点包括：

- 支持16种颜色输出，包括标准色和亮色
- 提供带占位符和不带占位符的打印方法
- 支持日志级别前缀，方便打印带有提示信息的消息
- 简洁易用的接口，方便开发者快速集成
- 支持全局实例和自定义实例
- 可控制颜色、加粗、下划线、闪烁等文本效果
- 支持禁用颜色输出（NoColor模式）
- 内置并发安全测试，可通过 `-race` 参数进行竞态检测

## 定义的颜色常量

以下是库中定义的颜色及其对应的 ANSI 颜色代码：

`PS`:带 `l`开头的为亮色!

| 颜色常量名 | 颜色名称 |
| ---------- | -------- |
| Black      | 黑色     |
| Red        | 红色     |
| Green      | 绿色     |
| Yellow     | 黄色     |
| Blue       | 蓝色     |
| Purple     | 紫色     |
| Cyan       | 青色     |
| White      | 白色     |
| Gray       | 灰色     |
| Lred       | 亮红色   |
| Lgreen     | 亮绿色   |
| Lyellow    | 亮黄色   |
| Lblue      | 亮蓝色   |
| Lpurple    | 亮紫色   |
| Lcyan      | 亮青色   |
| Lwhite     | 亮白色   |

## 提示信息级别和名称

以下是库中定义的提示级别及其对应的前缀：

| 提示级别 | 前缀名称  |
| :------- | :-------- |
| success  | [Success] |
| error    | [Error]   |
| warning  | [Warning] |
| info     | [Info]    |
| debug    | [Debug]   |
| ok       | ok:       |
| err      | err:      |
| warn     | warn:     |
| inf      | info:     |
| dbg      | debug:    |

## 内置的函数
### 创建实例函数

- `NewColorLib()`：创建一个新的 `ColorLib` 实例
- `GetCL()`：线程安全地获取全局 `ColorLib` 实例 (CL)

### 全局实例

ColorLib 提供了一个预初始化的全局实例 `CL`，可以直接使用而无需调用 `NewColorLib()`。可以通过 `GetCL()` 函数安全地获取这个实例。

## ColorLib 结构体
`ColorLib` 结构体用于管理颜色输出和文本效果。

| 字段名称  | 字段类型       | 字段描述                     |
| :-------- | :------------- | :--------------------------- |
| NoColor   | atomic.Bool    | 原子操作控制是否禁用颜色输出 |
| NoBold    | atomic.Bool    | 原子操作控制是否禁用加粗效果 |
| Underline | atomic.Bool    | 原子操作控制是否启用下划线   |
| Blink     | atomic.Bool    | 原子操作控制是否启用闪烁效果 |

> 所有字段均采用atomic原子类型，保证并发安全

`ColorLib` 结构体实现了以下方法：

### 终端打印（不支持占位符）自带换行符

| 方法名称                     | 描述                                 |
| ---------------------------- | ------------------------------------ |
| `Blue(msg ...any)`         | 打印蓝色信息到控制台（不带占位符）   |
| `Green(msg ...any)`        | 打印绿色信息到控制台（不带占位符）   |
| `Red(msg ...any)`          | 打印红色信息到控制台（不带占位符）   |
| `Yellow(msg ...any)`       | 打印黄色信息到控制台（不带占位符）   |
| `Purple(msg ...any)`       | 打印紫色信息到控制台（不带占位符）   |
| `PrintSuccess(msg ...any)` | 打印成功信息到控制台（不带占位符）   |
| `PrintError(msg ...any)`   | 打印错误信息到控制台（不带占位符）   |
| `PrintWarning(msg ...any)` | 打印警告信息到控制台（不带占位符）   |
| `PrintInfo(msg ...any)`    | 打印信息到控制台（不带占位符）       |
| `PrintDebug(msg ...any)`   | 打印调试信息到控制台（不带占位符）   |
| `Black(msg ...any)`        | 打印黑色信息到控制台（不带占位符）   |
| `Cyan(msg ...any)`         | 打印青色信息到控制台（不带占位符）   |
| `White(msg ...any)`        | 打印白色信息到控制台（不带占位符）   |
| `Gray(msg ...any)`         | 打印灰色信息到控制台（不带占位符）   |
| `Lred(msg ...any)`         | 打印亮红色信息到控制台（不带占位符） |
| `Lgreen(msg ...any)`       | 打印亮绿色信息到控制台（不带占位符） |
| `Lyellow(msg ...any)`      | 打印亮黄色信息到控制台（不带占位符） |
| `Lblue(msg ...any)`        | 打印亮蓝色信息到控制台（不带占位符） |
| `Lpurple(msg ...any)`      | 打印亮紫色信息到控制台（不带占位符） |
| `Lcyan(msg ...any)`        | 打印亮青色信息到控制台（不带占位符） |
| `Lwhite(msg ...any)`       | 打印亮白色信息到控制台（不带占位符） |
| `PrintDbg(msg...any)`      | 打印调试信息到控制台（不带占位符）   |
| `PrintInf(msg...any)`      | 打印信息到控制台（不带占位符）       |
| `PrintWarn(msg...any)`     | 打印警告信息到控制台（不带占位符）   |
| `PrintErr(msg...any)`      | 打印错误信息到控制台（不带占位符）   |
| `PrintOk(msg...any)`       | 打印成功信息到控制台（不带占位符）   |

### 终端打印（支持占位符）

| 方法名称                                   | 描述                               |
| ------------------------------------------ | ---------------------------------- |
| `Bluef(format string, a ...any)`         | 打印蓝色信息到控制台（带占位符）   |
| `Greenf(format string, a ...any)`        | 打印绿色信息到控制台（带占位符）   |
| `Redf(format string, a ...any)`          | 打印红色信息到控制台（带占位符）   |
| `Yellowf(format string, a ...any)`       | 打印黄色信息到控制台（带占位符）   |
| `Purplef(format string, a ...any)`       | 打印紫色信息到控制台（带占位符）   |
| `PrintSuccessf(format string, a ...any)` | 打印成功信息到控制台（带占位符）   |
| `PrintErrorf(format string, a ...any)`   | 打印错误信息到控制台（带占位符）   |
| `PrintWarningf(format string, a ...any)` | 打印警告信息到控制台（带占位符）   |
| `PrintInfof(format string, a ...any)`    | 打印信息到控制台（带占位符）       |
| `PrintDebugf(format string, a ...any)`   | 打印调试信息到控制台（带占位符）   |
| `Blackf(format string, a ...any)`        | 打印黑色信息到控制台（带占位符）   |
| `Cyanf(format string, a ...any)`         | 打印青色信息到控制台（带占位符）   |
| `Whitef(format string, a ...any)`        | 打印白色信息到控制台（带占位符）   |
| `Grayf(format string, a ...any)`         | 打印灰色信息到控制台（带占位符）   |
| `Lredf(format string, a ...any)`         | 打印亮红色信息到控制台（带占位符） |
| `Lgreenf(format string, a ...any)`       | 打印亮绿色信息到控制台（带占位符） |
| `Lyellowf(format string, a ...any)`      | 打印亮黄色信息到控制台（带占位符） |
| `Lbluef(format string, a ...any)`        | 打印亮蓝色信息到控制台（带占位符） |
| `Lpurplef(format string, a ...any)`      | 打印亮紫色信息到控制台（带占位符） |
| `Lcyanf(format string, a ...any)`        | 打印亮青色信息到控制台（带占位符） |
| `Lwhitef(format string, a ...any)`       | 打印亮白色信息到控制台（带占位符） |
| `PrintDbgf(format string, a...any)`      | 打印调试信息到控制台（带占位符）   |
| `PrintInff(format string, a...any)`      | 打印信息到控制台（带占位符）       |
| `PrintWarnf(format string, a...any)`     | 打印警告信息到控制台（带占位符）   |
| `PrintErrf(format string, a...any)`      | 打印错误信息到控制台（带占位符）   |
| `PrintOkf(format string, a...any)`       | 打印成功信息到控制台（带占位符）   |

### 返回构造字符串（不支持占位符）

| 方法名称                 | 描述                                   |
| ------------------------ | -------------------------------------- |
| `Sblue(msg ...any)`    | 返回构造后的蓝色字符串（不带占位符）   |
| `Sgreen(msg ...any)`   | 返回构造后的绿色字符串（不带占位符）   |
| `Sred(msg ...any)`     | 返回构造后的红色字符串（不带占位符）   |
| `Syellow(msg ...any)`  | 返回构造后的黄色字符串（不带占位符）   |
| `Spurple(msg ...any)`  | 返回构造后的紫色字符串（不带占位符）   |
| `Sblack(msg ...any)`   | 返回构造后的黑色字符串（不带占位符）   |
| `Scyan(msg ...any)`    | 返回构造后的青色字符串（不带占位符）   |
| `Swhite(msg ...any)`   | 返回构造后的白色字符串（不带占位符）   |
| `Sgray(msg ...any)`    | 返回构造后的灰色字符串（不带占位符）   |
| `Slred(msg ...any)`    | 返回构造后的亮红色字符串（不带占位符） |
| `Slgreen(msg ...any)`  | 返回构造后的亮绿色字符串（不带占位符） |
| `Slyellow(msg ...any)` | 返回构造后的亮黄色字符串（不带占位符） |
| `Slblue(msg ...any)`   | 返回构造后的亮蓝色字符串（不带占位符） |
| `Slpurple(msg ...any)` | 返回构造后的亮紫色字符串（不带占位符） |
| `Slcyan(msg ...any)`   | 返回构造后的亮青色字符串（不带占位符） |
| `Slwhite(msg ...any)`  | 返回构造后的亮白色字符串（不带占位符） |

### 返回构造字符串（支持占位符）

| 方法名称                               | 描述                                 |
| -------------------------------------- | ------------------------------------ |
| `Sbluef(format string, a ...any)`    | 返回构造后的蓝色字符串（带占位符）   |
| `Sgreenf(format string, a ...any)`   | 返回构造后的绿色字符串（带占位符）   |
| `Sredf(format string, a ...any)`     | 返回构造后的红色字符串（带占位符）   |
| `Syellowf(format string, a ...any)`  | 返回构造后的黄色字符串（带占位符）   |
| `Spurplef(format string, a ...any)`  | 返回构造后的紫色字符串（带占位符）   |
| `Sblackf(format string, a ...any)`   | 返回构造后的黑色字符串（带占位符）   |
| `Scyanf(format string, a ...any)`    | 返回构造后的青色字符串（带占位符）   |
| `Swhitef(format string, a ...any)`   | 返回构造后的白色字符串（带占位符）   |
| `Sgrayf(format string, a ...any)`    | 返回构造后的灰色字符串（带占位符）   |
| `Slredf(format string, a ...any)`    | 返回构造后的亮红色字符串（带占位符） |
| `Slgreenf(format string, a ...any)`  | 返回构造后的亮绿色字符串（带占位符） |
| `Slyellowf(format string, a ...any)` | 返回构造后的亮黄色字符串（带占位符） |
| `Slbluef(format string, a ...any)`   | 返回构造后的亮蓝色字符串（带占位符） |
| `Slpurplef(format string, a ...any)` | 返回构造后的亮紫色字符串（带占位符） |
| `Slcyanf(format string, a ...any)`   | 返回构造后的亮青色字符串（带占位符） |
| `Slwhitef(format string, a ...any)`  | 返回构造后的亮白色字符串（带占位符） |

### 通用颜色方法

| 方法名称                                   | 描述                                   |
| ------------------------------------------ | -------------------------------------- |
| `PrintColorf(code int, format string, a ...any)`    | 打印通用颜色信息到控制台（带占位符）   |
| `PrintColor(code int, msg ...any)`                  | 打印通用颜色信息到控制台（不带占位符） |
| `Scolorf(code int, format string, a ...any)` string | 返回构造后的通用颜色字符串（带占位符） |
| `Scolor(code int, msg ...any)` string               | 返回构造后的通用颜色字符串（不带占位符） |

## 文本效果控制
`ColorLib` 提供了多种文本效果控制选项，可以灵活调整输出样式。

### NoColor - 禁用颜色输出
```go
cl := NewColorLib()
cl.NoColor.Store(true)  // 原子操作禁用颜色
cl.NoBold.Store(false)  // 启用加粗
cl.Underline.Store(true) // 原子操作启用下划线

// 输出无颜色但有下划线的文本
cl.Red("这条消息将显示为无颜色但有下划线")
```

### NoBold - 禁用加粗效果
默认情况下文本会加粗显示，设置 `NoBold` 为 `true` 可以禁用加粗效果。

### Underline - 启用下划线
设置 `Underline` 为 `true` 可以为输出文本添加下划线效果。

### Blink - 启用闪烁效果
设置 `Blink` 为 `true` 可以让输出文本闪烁显示（部分终端可能不支持）。

### 使用场景

- 当终端不支持某些效果时
- 需要将输出重定向到文件时
- 需要特殊强调某些文本时
- 其他需要调整文本显示效果的场景

## 下载和使用

### 下载

通过 Go 模块管理工具下载 `ColorLib`：

```bash
go get gitee.com/MM-Q/colorlib
```

## 下载和使用
### 引入和使用

在您的 Go 代码中引入 `ColorLib`：

```go
package main

import (
	"gitee.com/MM-Q/colorlib"
)

func main() {
	// 使用全局实例CL（无需初始化）
	colorlib.CL.PrintDebug("这是一条来自全局实例的调试消息")
	colorlib.CL.PrintError("这是一条来自全局实例的错误消息")
	colorlib.CL.Blue("这是一条来自全局实例的蓝色消息")
	
	// 或者创建新的实例
	cl := colorlib.NewColorLib()

	// 打印带有颜色的文本
	cl.Blue("这是一条蓝色的消息")
	cl.Greenf("这是一条绿色的消息：%s\n", "Hello, ColorLib!")

	// 返回带有颜色的字符串
	coloredString := cl.Sred("这是一条红色的字符串")
	fmt.Println(coloredString)

	// 打印带有日志级别的消息
	cl.PrintSuccess("操作成功！")
	cl.PrintError("发生了一个错误")
	cl.PrintWarning("请注意：这是一个警告")
	cl.PrintInfo("这是一条普通信息")
}
```

### 全局实例

ColorLib 提供了一个预初始化的全局实例 `CL`，可以直接使用而无需调用 `NewColorLib()`。

```go
// 使用全局实例
colorlib.CL.PrintSuccess("操作成功！")
colorlib.CL.PrintError("发生错误")
colorlib.CL.Blue("蓝色文本")

// 返回带颜色的字符串
msg := colorlib.CL.Sgreen("绿色字符串")
fmt.Println(msg)
```

## 常用用法

以下是 `ColorLib` 的一些常用用法示例：

### 打印彩色文本

```go
cl := colorlib.NewColorLib()
cl.Blue("蓝色文本")
cl.Greenf("绿色文本：%s\n", "带占位符")
```

### 返回彩色字符串

```go
coloredString := cl.Spurple("紫色字符串")
fmt.Println(coloredString)
```

### 打印日志级别消息

```go
cl.PrintSuccess("操作成功！")
cl.PrintError("发生错误：参数无效")
cl.PrintWarning("警告：磁盘空间不足")
cl.PrintInfo("正在处理数据...")
cl.PrintDebug("正在测试...")
```

### 打印简洁版终端提示信息

```go
cl.PrintOk("操作成功")
cl.PrintErr("发生错误")
cl.PrintWarn("警告")
cl.PrintInf("信息")
cl.PrintDbg("调试信息")
```

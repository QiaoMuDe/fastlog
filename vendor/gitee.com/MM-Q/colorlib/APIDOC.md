# Package colorlib

Package colorlib 提供了彩色控制台输出和日志级别映射的功能。

```go
package colorlib // import "gitee.com/MM-Q/colorlib"
```

## Constants

```go
const (
	Black   = 30 // Black 黑色
	Red     = 31 // Red 红色
	Green   = 32 // Green 绿色
	Yellow  = 33 // Yellow 黄色
	Blue    = 34 // Blue 蓝色
	Purple  = 35 // Purple 紫色
	Cyan    = 36 // Cyan 青色
	White   = 37 // White 白色
	Gray    = 90 // Gray 灰色
	Lred    = 91 // Lred 亮红色
	Lgreen  = 92 // Lgreen 亮绿色
	Lyellow = 93 // Lyellow 亮黄色
	Lblue   = 94 // Lblue 亮蓝色
	Lpurple = 95 // Lpurple 亮紫色
	Lcyan   = 96 // Lcyan 亮青色
	Lwhite  = 97 // Lwhite 亮白色
)
```

## Types

### ColorLib

```go
type ColorLib struct {
	NoColor   atomic.Bool // NoColor 控制是否禁用颜色输出
	NoBold    atomic.Bool // NoBold 控制是否禁用字体加粗
	Underline atomic.Bool // Underline 控制是否启用下划线
	Blink     atomic.Bool // Blink 控制是否启用闪烁效果
	// Has unexported fields.
}
```

ColorLib 结构体用于管理颜色输出和日志级别映射。

### ColorLibInterface

```go
type ColorLibInterface interface {
	// 需要占位符的方法(自带换行符)
	Bluef(format string, a ...any)           // 打印蓝色信息到控制台（带占位符）
	Greenf(format string, a ...any)          // 打印绿色信息到控制台（带占位符）
	Redf(format string, a ...any)            // 打印红色信息到控制台（带占位符）
	Yellowf(format string, a ...any)         // 打印黄色信息到控制台（带占位符）
	Purplef(format string, a ...any)         // 打印紫色信息到控制台（带占位符）
	Sbluef(format string, a ...any) string   // 返回构造后的蓝色字符串（带占位符）
	Sgreenf(format string, a ...any) string  // 返回构造后的绿色字符串（带占位符）
	Sredf(format string, a ...any) string    // 返回构造后的红色字符串（带占位符）
	Syellowf(format string, a ...any) string // 返回构造后的黄色字符串（带占位符）
	Spurplef(format string, a ...any) string // 返回构造后的紫色字符串（带占位符）
	PrintSuccessf(format string, a ...any)   // 打印成功信息到控制台（带占位符）
	PrintErrorf(format string, a ...any)     // 打印错误信息到控制台（带占位符）
	PrintWarningf(format string, a ...any)   // 打印警告信息到控制台（带占位符）
	PrintInfof(format string, a ...any)      // 打印信息到控制台（带占位符）
	PrintDebugf(format string, a ...any)     // 打印调试信息到控制台（带占位符）

	// 直接打印信息, 无需占位符
	Blue(msg ...any)           // 打印蓝色信息到控制台, 无需占位符
	Green(msg ...any)          // 打印绿色信息到控制台, 无需占位符
	Red(msg ...any)            // 打印红色信息到控制台, 无需占位符
	Yellow(msg ...any)         // 打印黄色信息到控制台, 无需占位符
	Purple(msg ...any)         // 打印紫色信息到控制台, 无需占位符
	Sblue(msg ...any) string   // 返回构造后的蓝色字符串, 无需占位符
	Sgreen(msg ...any) string  // 返回构造后的绿色字符串, 无需占位符
	Sred(msg ...any) string    // 返回构造后的红色字符串, 无需占位符
	Syellow(msg ...any) string // 返回构造后的黄色字符串, 无需占位符
	Spurple(msg ...any) string // 返回构造后的紫色字符串, 无需占位符
	PrintSuccess(msg ...any)   // 打印成功信息到控制台, 无需占位符
	PrintError(msg ...any)     // 打印错误信息到控制台, 无需占位符
	PrintWarning(msg ...any)   // 打印警告信息到控制台, 无需占位符
	PrintInfo(msg ...any)      // 打印信息到控制台, 无需占位符
	PrintDebug(msg ...any)     // 打印调试信息到控制台, 无需占位符

	// 扩展颜色的方法
	Black(msg ...any)                         // 打印黑色信息到控制台, 无需占位符
	Blackf(format string, a ...any)           // 打印黑色信息到控制台（带占位符）
	Sblack(msg ...any) string                 // 返回构造后的黑色字符串, 无需占位符
	Sblackf(format string, a ...any) string   // 返回构造后的黑色字符串（带占位符）
	Cyan(msg ...any)                          // 打印青色信息到控制台, 无需占位符
	Cyanf(format string, a ...any)            // 打印青色信息到控制台（带占位符）
	Scyan(msg ...any) string                  // 返回构造后的青色字符串, 无需占位符
	Scyanf(format string, a ...any) string    // 返回构造后的青色字符串（带占位符）
	White(msg ...any)                         // 打印白色信息到控制台, 无需占位符
	Whitef(format string, a ...any)           // 打印白色信息到控制台（带占位符）
	Swhite(msg ...any) string                 // 返回构造后的白色字符串, 无需占位符
	Swhitef(format string, a ...any) string   // 返回构造后的白色字符串（带占位符）
	Gray(msg ...any)                          // 打印灰色信息到控制台, 无需占位符
	Grayf(format string, a ...any)            // 打印灰色信息到控制台（带占位符）
	Sgray(msg ...any) string                  // 返回构造后的灰色字符串, 无需占位符
	Sgrayf(format string, a ...any) string    // 返回构造后的灰色字符串（带占位符）
	Lred(msg ...any)                          // 打印亮红色信息到控制台, 无需占位符
	Lredf(format string, a ...any)            // 打印亮红色信息到控制台（带占位符）
	Slred(msg ...any) string                  // 返回构造后的亮红色字符串, 无需占位符
	Slredf(format string, a ...any) string    // 返回构造后的亮红色字符串（带占位符）
	Lgreen(msg ...any)                        // 打印亮绿色信息到控制台, 无需占位符
	Lgreenf(format string, a ...any)          // 打印亮绿色信息到控制台（带占位符）
	Slgreen(msg ...any) string                // 返回构造后的亮绿色字符串, 无需占位符
	Slgreenf(format string, a ...any) string  // 返回构造后的亮绿色字符串（带占位符）
	Lyellow(msg ...any)                       // 打印亮黄色信息到控制台, 无需占位符
	Lyellowf(format string, a ...any)         // 打印亮黄色信息到控制台（带占位符）
	Slyellow(msg ...any) string               // 返回构造后的亮黄色字符串, 无需占位符
	Slyellowf(format string, a ...any) string // 返回构造后的亮黄色字符串（带占位符）
	Lblue(msg ...any)                         // 打印亮蓝色信息到控制台, 无需占位符
	Lbluef(format string, a ...any)           // 打印亮蓝色信息到控制台（带占位符）
	Slblue(msg ...any) string                 // 返回构造后的亮蓝色字符串, 无需占位符
	Slbluef(format string, a ...any) string   // 返回构造后的亮蓝色字符串（带占位符）
	Lpurple(msg ...any)                       // 打印亮紫色信息到控制台, 无需占位符
	Lpurplef(format string, a ...any)         // 打印亮紫色信息到控制台（带占位符）
	Slpurple(msg ...any) string               // 返回构造后的亮紫色字符串, 无需占位符
	Slpurplef(format string, a ...any) string // 返回构造后的亮紫色字符串（带占位符）
	Lcyan(msg ...any)                         // 打印亮青色信息到控制台, 无需占位符
	Lcyanf(format string, a ...any)           // 打印亮青色信息到控制台（带占位符）
	Slcyan(msg ...any) string                 // 返回构造后的亮青色字符串, 无需占位符
	Slcyanf(format string, a ...any) string   // 返回构造后的亮青色字符串（带占位符）
	Lwhite(msg ...any)                        // 打印亮白色信息到控制台, 无需占位符
	Lwhitef(format string, a ...any)          // 打印亮白色信息到控制台（带占位符）
	Slwhite(msg ...any) string                // 返回构造后的亮白色字符串, 无需占位符
	Slwhitef(format string, a ...any) string  // 返回构造后的亮白色字符串（带占位符）

	// 简洁版的方法, 无需占位符
	PrintOk(msg ...any)   // 打印成功信息到控制台, 无需占位符
	PrintErr(msg ...any)  // 打印错误信息到控制台, 无需占位符
	PrintInf(msg ...any)  // 打印信息到控制台, 无需占位符
	PrintDbg(msg ...any)  // 打印调试信息到控制台, 无需占位符
	PrintWarn(msg ...any) // 打印警告信息到控制台, 无需占位符

	// 简洁版的方法, 带占位符
	PrintOkf(format string, a ...any)   // 打印成功信息到控制台（带占位符）
	PrintErrf(format string, a ...any)  // 打印错误信息到控制台（带占位符）
	PrintInff(format string, a ...any)  // 打印信息到控制台（带占位符）
	PrintDbgf(format string, a ...any)  // 打印调试信息到控制台（带占位符）
	PrintWarnf(format string, a ...any) // 打印警告信息到控制台（带占位符）

	// 通用颜色方法
	PrintColorf(code int, format string, a ...any)    // 打印通用颜色信息到控制台（带占位符）
	PrintColor(code int, msg ...any)                  // 打印通用颜色信息到控制台, 无需占位符
	Scolorf(code int, format string, a ...any) string // 返回构造后的通用颜色字符串（带占位符）
	Scolor(code int, msg ...any) string               // 返回构造后的通用颜色字符串, 无需占位符

	// 控制颜色输出和字体样式方法
	SetNoColor(enable bool) *ColorLib   // 设置是否禁用颜色输出
	SetNoBold(enable bool) *ColorLib    // 设置是否禁用字体加粗
	SetUnderline(enable bool) *ColorLib // 设置是否启用下划线
	SetBlink(enable bool) *ColorLib     // 设置是否启用闪烁
}
```

ColorLibInterface 是一个接口，定义了一组方法，用于打印和返回带有颜色的文本。

## Variables

```go
var (
	CL *ColorLib // 全局实例可导入直接使用
)
```

## Functions

```go
func GetCL() *ColorLib
```

GetCL 是一个线程安全用于获取全局唯一的 ColorLib 实例的函数。

```go
func NewColorLib() *ColorLib
```

NewColorLib 函数用于创建一个新的 ColorLib 实例。

## Methods

### Black

```go
func (c *ColorLib) Black(msg ...any)
```

Black 方法用于将传入的参数以黑色文本形式打印到控制台（不带占位符）。

### Blackf

```go
func (c *ColorLib) Blackf(format string, a ...any)
```

Blackf 方法用于将传入的参数以黑色文本形式打印到控制台（带占位符）。

### Blue

```go
func (c *ColorLib) Blue(msg ...any)
```

Blue 方法用于将传入的参数以蓝色文本形式打印到控制台（不带占位符）。

### Bluef

```go
func (c *ColorLib) Bluef(format string, a ...any)
```

Bluef 方法用于将传入的参数以蓝色文本形式打印到控制台（带占位符）。

### Cyan

```go
func (c *ColorLib) Cyan(msg ...any)
```

Cyan 方法用于将传入的参数以青色文本形式打印到控制台（不带占位符）。

### Cyanf

```go
func (c *ColorLib) Cyanf(format string, a ...any)
```

Cyanf 方法用于将传入的参数以青色文本形式打印到控制台（带占位符）。

### Gray

```go
func (c *ColorLib) Gray(msg ...any)
```

Gray 方法用于将传入的参数以灰色文本形式打印到控制台（不带占位符）。

### Grayf

```go
func (c *ColorLib) Grayf(format string, a ...any)
```

Grayf 方法用于将传入的参数以灰色文本形式打印到控制台（带占位符）。

### Green

```go
func (c *ColorLib) Green(msg ...any)
```

Green 方法用于将传入的参数以绿色文本形式打印到控制台（不带占位符）。

### Greenf

```go
func (c *ColorLib) Greenf(format string, a ...any)
```

Greenf 方法用于将传入的参数以绿色文本形式打印到控制台（带占位符）。

### Lblue

```go
func (c *ColorLib) Lblue(msg ...any)
```

Lblue 方法用于将传入的参数以亮蓝色文本形式打印到控制台（不带占位符）。

### Lbluef

```go
func (c *ColorLib) Lbluef(format string, a ...any)
```

Lbluef 方法用于将传入的参数以亮蓝色文本形式打印到控制台（带占位符）。

### Lcyan

```go
func (c *ColorLib) Lcyan(msg ...any)
```

Lcyan 方法用于将传入的参数以亮青色文本形式打印到控制台（不带占位符）。

### Lcyanf

```go
func (c *ColorLib) Lcyanf(format string, a ...any)
```

Lcyanf 方法用于将传入的参数以亮青色文本形式打印到控制台（带占位符）。

### Lgreen

```go
func (c *ColorLib) Lgreen(msg ...any)
```

Lgreen 方法用于将传入的参数以亮绿色文本形式打印到控制台（不带占位符）。

### Lgreenf

```go
func (c *ColorLib) Lgreenf(format string, a ...any)
```

Lgreenf 方法用于将传入的参数以亮绿色文本形式打印到控制台（带占位符）。

### Lpurple

```go
func (c *ColorLib) Lpurple(msg ...any)
```

Lpurple 方法用于将传入的参数以亮紫色文本形式打印到控制台（不带占位符）。

### Lpurplef

```go
func (c *ColorLib) Lpurplef(format string, a ...any)
```

Lpurplef 方法用于将传入的参数以亮紫色文本形式打印到控制台（带占位符）。

### Lred

```go
func (c *ColorLib) Lred(msg ...any)
```

Lred 方法用于将传入的参数以亮红色文本形式打印到控制台（不带占位符）。

### Lredf

```go
func (c *ColorLib) Lredf(format string, a ...any)
```

Lredf 方法用于将传入的参数以亮红色文本形式打印到控制台（带占位符）。

### Lwhite

```go
func (c *ColorLib) Lwhite(msg ...any)
```

Lwhite 方法用于将传入的参数以亮白色文本形式打印到控制台（不带占位符）。

### Lwhitef

```go
func (c *ColorLib) Lwhitef(format string, a ...any)
```

Lwhitef 方法用于将传入的参数以亮白色文本形式打印到控制台（带占位符）。

### Lyellow

```go
func (c *ColorLib) Lyellow(msg ...any)
```

Lyellow 方法用于将传入的参数以亮黄色文本形式打印到控制台（不带占位符）。

### Lyellowf

```go
func (c *ColorLib) Lyellowf(format string, a ...any)
```

Lyellowf 方法用于将传入的参数以亮黄色文本形式打印到控制台（带占位符）。

### PrintColorf

```go
func (c *ColorLib) PrintColorf(code int, format string, a ...any)
```

PrintColorf 方法根据颜色代码常量打印对应颜色的文本。

### PrintColorln

```go
func (c *ColorLib) PrintColorln(code int, msg ...any)
```

PrintColorln 方法根据颜色代码常量打印对应颜色的文本。

### PrintDbg

```go
func (c *ColorLib) PrintDbg(msg ...any)
```

PrintDbg 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（不带占位符）。

### PrintDbgf

```go
func (c *ColorLib) PrintDbgf(format string, a ...any)
```

PrintDbgf 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（带占位符）。

### PrintDebug

```go
func (c *ColorLib) PrintDebug(msg ...any)
```

PrintDebug 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（不带占位符）。

### PrintDebugf

```go
func (c *ColorLib) PrintDebugf(format string, a ...any)
```

PrintDebugf 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（带占位符）。

### PrintErr

```go
func (c *ColorLib) PrintErr(msg ...any)
```

PrintErr 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（不带占位符）。

### PrintErrf

```go
func (c *ColorLib) PrintErrf(format string, a ...any)
```

PrintErrf 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（带占位符）。

### PrintError

```go
func (c *ColorLib) PrintError(msg ...any)
```

PrintError 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（不带占位符）。

### PrintErrorf

```go
func (c *ColorLib) PrintErrorf(format string, a ...any)
```

PrintErrorf 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（带占位符）。

### PrintInf

```go
func (c *ColorLib) PrintInf(msg ...any)
```

PrintInf 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（不带占位符）。

### PrintInff

```go
func (c *ColorLib) PrintInff(format string, a ...any)
```

PrintInff 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（带占位符）。

### PrintInfo

```go
func (c *ColorLib) PrintInfo(msg ...any)
```

PrintInfo 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（不带占位符）。

### PrintInfof

```go
func (c *ColorLib) PrintInfof(format string, a ...any)
```

PrintInfof 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（带占位符）。

### PrintOk

```go
func (c *ColorLib) PrintOk(msg ...any)
```

PrintOk 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（不带占位符）。

### PrintOkf

```go
func (c *ColorLib) PrintOkf(format string, a ...any)
```

PrintOkf 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（带占位符）。

### PrintSuccess

```go
func (c *ColorLib) PrintSuccess(msg ...any)
```

PrintSuccess 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（不带占位符）。

### PrintSuccessf

```go
func (c *ColorLib) PrintSuccessf(format string, a ...any)
```

PrintSuccessf 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（带占位符）。

### PrintWarn

```go
func (c *ColorLib) PrintWarn(msg ...any)
```

PrintWarn 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（不带占位符）。

### PrintWarnf

```go
func (c *ColorLib) PrintWarnf(format string, a ...any)
```

PrintWarnf 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（带占位符）。

### PrintWarning

```go
func (c *ColorLib) PrintWarning(msg ...any)
```

PrintWarning 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（不带占位符）。

### PrintWarningf

```go
func (c *ColorLib) PrintWarningf(format string, a ...any)
```

PrintWarningf 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（带占位符）。

### Purple

```go
func (c *ColorLib) Purple(msg ...any)
```

Purple 方法用于将传入的参数以紫色文本形式打印到控制台（不带占位符）。

### Purplef

```go
func (c *ColorLib) Purplef(format string, a ...any)
```

Purplef 方法用于将传入的参数以紫色文本形式打印到控制台（带占位符）。

### Red

```go
func (c *ColorLib) Red(msg ...any)
```

Red 方法用于将传入的参数以红色文本形式打印到控制台（不带占位符）。

### Redf

```go
func (c *ColorLib) Redf(format string, a ...any)
```

Redf 方法用于将传入的参数以红色文本形式打印到控制台（带占位符）。

### SColor

```go
func (c *ColorLib) SColor(code int, msg ...any) string
```

SColor 方法根据颜色代码常量打印对应颜色的文本。

### SColorf

```go
func (c *ColorLib) SColorf(code int, format string, a ...any) string
```

SColorf 方法根据颜色代码常量打印对应颜色的文本。

### Sblack

```go
func (c *ColorLib) Sblack(msg ...any) string
```

Sblack 方法用于将传入的参数以黑色文本形式返回（不带占位符）。

### Sblackf

```go
func (c *ColorLib) Sblackf(format string, a ...any) string
```

Sblackf 方法用于将传入的参数以黑色文本形式返回（带占位符）。

### Sblue

```go
func (c *ColorLib) Sblue(msg ...any) string
```

Sblue 方法用于将传入的参数以蓝色文本形式返回（不带占位符）。

### Sbluef

```go
func (c *ColorLib) Sbluef(format string, a ...any) string
```

Sbluef 方法用于将传入的参数以蓝色文本形式返回（带占位符）。

### Scyan

```go
func (c *ColorLib) Scyan(msg ...any) string
```

Scyan 方法用于将传入的参数以青色文本形式返回（不带占位符）。

### Scyanf

```go
func (c *ColorLib) Scyanf(format string, a ...any) string
```

Scyanf 方法用于将传入的参数以青色文本形式返回（带占位符）。

### SetBlink

```go
func (c *ColorLib) SetBlink(enable bool) *ColorLib
```

SetBlink 设置是否启用闪烁效果,并返回ColorLib实例以支持链式调用。

### SetNoBold

```go
func (c *ColorLib) SetNoBold(enable bool) *ColorLib
```

SetNoBold 设置是否禁用字体加粗,并返回ColorLib实例以支持链式调用。

### SetNoColor

```go
func (c *ColorLib) SetNoColor(enable bool) *ColorLib
```

SetNoColor 设置是否禁用颜色输出,并返回ColorLib实例以支持链式调用。

### SetUnderline

```go
func (c *ColorLib) SetUnderline(enable bool) *ColorLib
```

SetUnderline 设置是否启用下划线,并返回ColorLib实例以支持链式调用。

### Sgray

```go
func (c *ColorLib) Sgray(msg ...any) string
```

Sgray 方法用于将传入的参数以灰色文本形式返回（不带占位符）。

### Sgrayf

```go
func (c *ColorLib) Sgrayf(format string, a ...any) string
```

Sgrayf 方法用于将传入的参数以灰色文本形式返回（带占位符）。

### Sgreen

```go
func (c *ColorLib) Sgreen(msg ...any) string
```

Sgreen 方法用于将传入的参数以绿色文本形式返回（不带占位符）。

### Sgreenf

```go
func (c *ColorLib) Sgreenf(format string, a ...any) string
```

Sgreenf 方法用于将传入的参数以绿色文本形式返回（带占位符）。

### Slblue

```go
func (c *ColorLib) Slblue(msg ...any) string
```

Slblue 方法用于将传入的参数以亮蓝色文本形式返回（不带占位符）。

### Slbluef

```go
func (c *ColorLib) Slbluef(format string, a ...any) string
```

Slbluef 方法用于将传入的参数以亮蓝色文本形式返回（带占位符）。

### Slcyan

```go
func (c *ColorLib) Slcyan(msg ...any) string
```

Slcyan 方法用于将传入的参数以亮青色文本形式返回（不带占位符）。

### Slcyanf

```go
func (c *ColorLib) Slcyanf(format string, a ...any) string
```

Slcyanf 方法用于将传入的参数以亮青色文本形式返回（带占位符）。

### Slgreen

```go
func (c *ColorLib) Slgreen(msg ...any) string
```

Slgreen 方法用于将传入的参数以亮绿色文本形式返回（不带占位符）。

### Slgreenf

```go
func (c *ColorLib) Slgreenf(format string, a ...any) string
```

Slgreenf 方法用于将传入的参数以亮绿色文本形式返回（带占位符）。

### Slpurple

```go
func (c *ColorLib) Slgreen(msg ...any) string
```

Slpurple 方法用于将传入的参数以亮紫色文本形式返回（不带占位符）。

### Slpurplef

```go
func (c *ColorLib) Slgreenf(format string, a ...any) string
```

Slpurplef 方法用于将传入的参数以亮紫色文本形式返回（带占位符）。

### Slred

```go
func (c *ColorLib) Slred(msg ...any) string
```

Slred 方法用于将传入的参数以亮红色文本形式返回（不带占位符）。

### Slredf

```go
func (c *ColorLib) Slredf(format string, a ...any) string
```

Slredf 方法用于将传入的参数以亮红色文本形式返回（带占位符）。

### Slwhite

```go
func (c *ColorLib) Slwhite(msg ...any) string
```

Slwhite 方法用于将传入的参数以亮白色文本形式返回（不带占位符）。

### Slwhitef

```go
func (c *ColorLib) Slwhitef(format string, a ...any) string
```

Slwhitef 方法用于将传入的参数以亮白色文本形式返回（带占位符）。

### Slyellow

```go
func (c *ColorLib) Slyellow(msg ...any) string
```

Slyellow 方法用于将传入的参数以亮黄色文本形式返回（不带占位符）。

### Slyellowf

```go
func (c *ColorLib) Slyellowf(format string, a ...any) string
```

Slyellowf 方法用于将传入的参数以亮黄色文本形式返回（带占位符）。

### Spurple

```go
func (c *ColorLib) Spurple(msg ...any) string
```

Spurple 方法用于将传入的参数以紫色文本形式返回（不带占位符）。

### Spurplef

```go
func (c *ColorLib) Spurplef(format string, a ...any) string
```

Spurplef 方法用于将传入的参数以紫色文本形式返回（带占位符）。

### Sred

```go
func (c *ColorLib) Sred(msg ...any) string
```

Sred 方法用于将传入的参数以红色文本形式返回（不带占位符）。

### Sredf

```go
func (c *ColorLib) Sredf(format string, a ...any) string
```

Sredf 方法用于将传入的参数以红色文本形式返回（带占位符）。

### Swhite

```go
func (c *ColorLib) Swhite(msg ...any) string
```

Swhite 方法用于将传入的参数以白色文本形式返回（不带占位符）。

### Swhitef

```go
func (c *ColorLib) Swhitef(format string, a ...any) string
```

Swhitef 方法用于将传入的参数以白色文本形式返回（带占位符）。

### Syellow

```go
func (c *ColorLib) Syellow(msg ...any) string
```

Syellow 方法用于将传入的参数以黄色文本形式返回（不带占位符）。

### Syellowf

```go
func (c *ColorLib) Syellowf(format string, a ...any) string
```

Syellowf 方法用于将传入的参数以黄色文本形式返回（带占位符）。

### White

```go
func (c *ColorLib) White(msg ...any)
```

White 方法用于将传入的参数以白色文本形式打印到控制台（不带占位符）。

### Whitef

```go
func (c *ColorLib) Whitef(format string, a ...any)
```

Whitef 方法用于将传入的参数以白色文本形式打印到控制台（带占位符）。

### Yellow

```go
func (c *ColorLib) Yellow(msg ...any)
```

Yellow 方法用于将传入的参数以黄色文本形式打印到控制台（不带占位符）。

### Yellowf

```go
func (c *ColorLib) Yellowf(format string, a ...any)
```

Yellowf 方法用于将传入的参数以黄色文本形式打印到控制台（带占位符）。
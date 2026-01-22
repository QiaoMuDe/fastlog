# Package colorlib

Package colorlib 提供了基础的颜色输出功能，包括格式化输出、直接打印和字符串返回等方法。该文件实现了 ColorLib 结构体的基础颜色方法，支持标准颜色和亮色系列的文本输出。

Package colorlib 是一个功能强大的 Go 语言终端颜色输出库。它提供了丰富的颜色输出功能，包括基础颜色、亮色、样式设置（粗体、下划线、闪烁）等。支持链式调用、自定义输出接口、线程安全的全局实例等特性。主要用于在终端中输出带有颜色和样式的文本，提升命令行程序的用户体验。

Package colorlib 提供了扩展的日志级别输出功能。该文件实现了带有级别标识的消息输出方法，包括调试、信息、成功、警告和错误等级别。每个级别都有对应的颜色和前缀标识，便于在终端中快速识别不同类型的消息。

Package colorlib 提供了内部核心功能实现。该文件实现了 ColorLib 结构体的核心方法，包括颜色输出、ANSI 序列构建、样式处理和消息格式化等功能，是整个颜色库的核心实现文件。

## VARIABLES

```go
var New = NewColorLib
```

- **New**: 创建一个新的 ColorLib 实例（NewColorLib 的别名）。
  - 返回值:
    - `*ColorLib`: 新创建的 ColorLib 实例指针

## TYPES

### ColorLib

```go
type ColorLib struct {
	// Has unexported fields.
}
```

ColorLib 结构体用于管理颜色输出和日志级别映射。

### Functions

#### GetCL

```go
func GetCL() *ColorLib
```

GetCL 是一个线程安全用于获取全局唯一的 ColorLib 实例的函数。

- 返回值:
  - `*ColorLib`: 全局唯一的 ColorLib 实例指针

#### NewColorLib

```go
func NewColorLib() *ColorLib
```

NewColorLib 函数用于创建一个新的 ColorLib 实例（默认输出到标准输出）。

- 返回值:
  - `*ColorLib`: 新创建的 ColorLib 实例指针

#### NewColorLibWithWriter

```go
func NewColorLibWithWriter(writer io.Writer) *ColorLib
```

NewColorLibWithWriter 创建一个指定输出接口的 ColorLib 实例。

- 参数:
  - `writer`: 输出接口，如 os.Stdout, os.Stderr, 文件等
- 返回值:
  - `*ColorLib`: 新创建的 ColorLib 实例指针

### Methods

#### Black

```go
func (c *ColorLib) Black(msg ...any)
```

Black 方法用于将传入的参数以黑色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### Blackf

```go
func (c *ColorLib) Blackf(format string, a ...any)
```

Blackf 方法用于将传入的参数以黑色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### Blue

```go
func (c *ColorLib) Blue(msg ...any)
```

Blue 方法用于将传入的参数以蓝色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### Bluef

```go
func (c *ColorLib) Bluef(format string, a ...any)
```

Bluef 方法用于将传入的参数以蓝色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### BrightBlue

```go
func (c *ColorLib) BrightBlue(msg ...any)
```

BrightBlue 方法用于将传入的参数以亮蓝色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### BrightBluef

```go
func (c *ColorLib) BrightBluef(format string, a ...any)
```

BrightBluef 方法用于将传入的参数以亮蓝色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### BrightCyan

```go
func (c *ColorLib) BrightCyan(msg ...any)
```

BrightCyan 方法用于将传入的参数以亮青色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### BrightCyanf

```go
func (c *ColorLib) BrightCyanf(format string, a ...any)
```

BrightCyanf 方法用于将传入的参数以亮青色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### BrightGreen

```go
func (c *ColorLib) BrightGreen(msg ...any)
```

BrightGreen 方法用于将传入的参数以亮绿色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### BrightGreenf

```go
func (c *ColorLib) BrightGreenf(format string, a ...any)
```

BrightGreenf 方法用于将传入的参数以亮绿色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### BrightMagenta

```go
func (c *ColorLib) BrightMagenta(msg ...any)
```

BrightMagenta 方法用于将传入的参数以亮品红色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### BrightMagentaf

```go
func (c *ColorLib) BrightMagentaf(format string, a ...any)
```

BrightMagentaf 方法用于将传入的参数以亮品红色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### BrightRed

```go
func (c *ColorLib) BrightRed(msg ...any)
```

BrightRed 方法用于将传入的参数以亮红色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### BrightRedf

```go
func (c *ColorLib) BrightRedf(format string, a ...any)
```

BrightRedf 方法用于将传入的参数以亮红色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### BrightWhite

```go
func (c *ColorLib) BrightWhite(msg ...any)
```

BrightWhite 方法用于将传入的参数以亮白色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### BrightWhitef

```go
func (c *ColorLib) BrightWhitef(format string, a ...any)
```

BrightWhitef 方法用于将传入的参数以亮白色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### BrightYellow

```go
func (c *ColorLib) BrightYellow(msg ...any)
```

BrightYellow 方法用于将传入的参数以亮黄色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### BrightYellowf

```go
func (c *ColorLib) BrightYellowf(format string, a ...any)
```

BrightYellowf 方法用于将传入的参数以亮黄色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### Cyan

```go
func (c *ColorLib) Cyan(msg ...any)
```

Cyan 方法用于将传入的参数以青色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### Cyanf

```go
func (c *ColorLib) Cyanf(format string, a ...any)
```

Cyanf 方法用于将传入的参数以青色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### Gray

```go
func (c *ColorLib) Gray(msg ...any)
```

Gray 方法用于将传入的参数以灰色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### Grayf

```go
func (c *ColorLib) Grayf(format string, a ...any)
```

Grayf 方法用于将传入的参数以灰色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### Green

```go
func (c *ColorLib) Green(msg ...any)
```

Green 方法用于将传入的参数以绿色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### Greenf

```go
func (c *ColorLib) Greenf(format string, a ...any)
```

Greenf 方法用于将传入的参数以绿色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### Magenta

```go
func (c *ColorLib) Magenta(msg ...any)
```

Magenta 方法用于将传入的参数以品红色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### Magentaf

```go
func (c *ColorLib) Magentaf(format string, a ...any)
```

Magentaf 方法用于将传入的参数以品红色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### PrintDebug

```go
func (c *ColorLib) PrintDebug(msg ...any)
```

PrintDebug 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### PrintDebugf

```go
func (c *ColorLib) PrintDebugf(format string, a ...any)
```

PrintDebugf 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### PrintError

```go
func (c *ColorLib) PrintError(msg ...any)
```

PrintError 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### PrintErrorf

```go
func (c *ColorLib) PrintErrorf(format string, a ...any)
```

PrintErrorf 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### PrintInfo

```go
func (c *ColorLib) PrintInfo(msg ...any)
```

PrintInfo 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### PrintInfof

```go
func (c *ColorLib) PrintInfof(format string, a ...any)
```

PrintInfof 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### PrintOk

```go
func (c *ColorLib) PrintOk(msg ...any)
```

PrintOk 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### PrintOkf

```go
func (c *ColorLib) PrintOkf(format string, a ...any)
```

PrintOkf 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### PrintWarn

```go
func (c *ColorLib) PrintWarn(msg ...any)
```

PrintWarn 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### PrintWarnf

```go
func (c *ColorLib) PrintWarnf(format string, a ...any)
```

PrintWarnf 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### Red

```go
func (c *ColorLib) Red(msg ...any)
```

Red 方法用于将传入的参数以红色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### Redf

```go
func (c *ColorLib) Redf(format string, a ...any)
```

Redf 方法用于将传入的参数以红色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### Sblack

```go
func (c *ColorLib) Sblack(msg ...any) string
```

Sblack 方法用于将传入的参数以黑色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有黑色格式的字符串

#### Sblackf

```go
func (c *ColorLib) Sblackf(format string, a ...any) string
```

Sblackf 方法用于将传入的参数以黑色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有黑色格式的字符串

#### Sblue

```go
func (c *ColorLib) Sblue(msg ...any) string
```

Sblue 方法用于将传入的参数以蓝色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有蓝色格式的字符串

#### Sbluef

```go
func (c *ColorLib) Sbluef(format string, a ...any) string
```

Sbluef 方法用于将传入的参数以蓝色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有蓝色格式的字符串

#### SbrightBlue

```go
func (c *ColorLib) SbrightBlue(msg ...any) string
```

SbrightBlue 方法用于将传入的参数以亮蓝色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有亮蓝色格式的字符串

#### SbrightBluef

```go
func (c *ColorLib) SbrightBluef(format string, a ...any) string
```

SbrightBluef 方法用于将传入的参数以亮蓝色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有亮蓝色格式的字符串

#### SbrightCyan

```go
func (c *ColorLib) SbrightCyan(msg ...any) string
```

SbrightCyan 方法用于将传入的参数以亮青色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有亮青色格式的字符串

#### SbrightCyanf

```go
func (c *ColorLib) SbrightCyanf(format string, a ...any) string
```

SbrightCyanf 方法用于将传入的参数以亮青色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有亮青色格式的字符串

#### SbrightGreen

```go
func (c *ColorLib) SbrightGreen(msg ...any) string
```

SbrightGreen 方法用于将传入的参数以亮绿色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有亮绿色格式的字符串

#### SbrightGreenf

```go
func (c *ColorLib) SbrightGreenf(format string, a ...any) string
```

SbrightGreenf 方法用于将传入的参数以亮绿色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有亮绿色格式的字符串

#### SbrightMagenta

```go
func (c *ColorLib) SbrightMagenta(msg ...any) string
```

SbrightMagenta 方法用于将传入的参数以亮品红色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有亮品红色格式的字符串

#### SbrightMagentaf

```go
func (c *ColorLib) SbrightMagentaf(format string, a ...any) string
```

SbrightMagentaf 方法用于将传入的参数以亮品红色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有亮品红色格式的字符串

#### SbrightRed

```go
func (c *ColorLib) SbrightRed(msg ...any) string
```

SbrightRed 方法用于将传入的参数以亮红色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有亮红色格式的字符串

#### SbrightRedf

```go
func (c *ColorLib) SbrightRedf(format string, a ...any) string
```

SbrightRedf 方法用于将传入的参数以亮红色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有亮红色格式的字符串

#### SbrightWhite

```go
func (c *ColorLib) SbrightWhite(msg ...any) string
```

SbrightWhite 方法用于将传入的参数以亮白色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有亮白色格式的字符串

#### SbrightWhitef

```go
func (c *ColorLib) SbrightWhitef(format string, a ...any) string
```

SbrightWhitef 方法用于将传入的参数以亮白色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有亮白色格式的字符串

#### SbrightYellow

```go
func (c *ColorLib) SbrightYellow(msg ...any) string
```

SbrightYellow 方法用于将传入的参数以亮黄色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有亮黄色格式的字符串

#### SbrightYellowf

```go
func (c *ColorLib) SbrightYellowf(format string, a ...any) string
```

SbrightYellowf 方法用于将传入的参数以亮黄色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有亮黄色格式的字符串

#### Scyan

```go
func (c *ColorLib) Scyan(msg ...any) string
```

Scyan 方法用于将传入的参数以青色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有青色格式的字符串

#### Scyanf

```go
func (c *ColorLib) Scyanf(format string, a ...any) string
```

Scyanf 方法用于将传入的参数以青色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有青色格式的字符串

#### SetBlink

```go
func (c *ColorLib) SetBlink(enabled bool)
```

SetBlink 设置是否启用闪烁输出。

- 参数:
  - `enabled`: 是否启用闪烁输出（`true` - 启用，`false` - 禁用）

#### SetBold

```go
func (c *ColorLib) SetBold(enabled bool)
```

SetBold 设置是否启用粗体输出。

- 参数:
  - `enabled`: 是否启用粗体输出（`true` - 启用，`false` - 禁用）

#### SetColor

```go
func (c *ColorLib) SetColor(enabled bool)
```

SetColor 设置是否启用颜色输出。

- 参数:
  - `enabled`: 是否启用颜色输出（`true` - 启用，`false` - 禁用）

#### SetUnderline

```go
func (c *ColorLib) SetUnderline(enabled bool)
```

SetUnderline 设置是否启用下划线输出。

- 参数:
  - `enabled`: 是否启用下划线输出（`true` - 启用，`false` - 禁用）

#### Sgray

```go
func (c *ColorLib) Sgray(msg ...any) string
```

Sgray 方法用于将传入的参数以灰色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有灰色格式的字符串

#### Sgrayf

```go
func (c *ColorLib) Sgrayf(format string, a ...any) string
```

Sgrayf 方法用于将传入的参数以灰色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有灰色格式的字符串

#### Sgreen

```go
func (c *ColorLib) Sgreen(msg ...any) string
```

Sgreen 方法用于将传入的参数以绿色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有绿色格式的字符串

#### Sgreenf

```go
func (c *ColorLib) Sgreenf(format string, a ...any) string
```

Sgreenf 方法用于将传入的参数以绿色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有绿色格式的字符串

#### Smagenta

```go
func (c *ColorLib) Smagenta(msg ...any) string
```

Smagenta 方法用于将传入的参数以品红色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有品红色格式的字符串

#### Smagentaf

```go
func (c *ColorLib) Smagentaf(format string, a ...any) string
```

Smagentaf 方法用于将传入的参数以品红色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有品红色格式的字符串

#### Sred

```go
func (c *ColorLib) Sred(msg ...any) string
```

Sred 方法用于将传入的参数以红色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有红色格式的字符串

#### Sredf

```go
func (c *ColorLib) Sredf(format string, a ...any) string
```

Sredf 方法用于将传入的参数以红色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有红色格式的字符串

#### Swhite

```go
func (c *ColorLib) Swhite(msg ...any) string
```

Swhite 方法用于将传入的参数以白色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有白色格式的字符串

#### Swhitef

```go
func (c *ColorLib) Swhitef(format string, a ...any) string
```

Swhitef 方法用于将传入的参数以白色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有白色格式的字符串

#### Syellow

```go
func (c *ColorLib) Syellow(msg ...any) string
```

Syellow 方法用于将传入的参数以黄色文本形式返回。

- 参数:
  - `msg`: 可变参数，要处理的消息内容
- 返回值:
  - `string`: 带有黄色格式的字符串

#### Syellowf

```go
func (c *ColorLib) Syellowf(format string, a ...any) string
```

Syellowf 方法用于将传入的参数以黄色文本形式返回。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符
- 返回值:
  - `string`: 带有黄色格式的字符串

#### White

```go
func (c *ColorLib) White(msg ...any)
```

White 方法用于将传入的参数以白色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### Whitef

```go
func (c *ColorLib) Whitef(format string, a ...any)
```

Whitef 方法用于将传入的参数以白色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### WithBlink

```go
func (c *ColorLib) WithBlink(enabled bool) *ColorLib
```

WithBlink 启用闪烁效果（链式调用）。

- 参数:
  - `enabled`: 是否启用闪烁效果（`true` - 启用，`false` - 禁用）
- 返回值:
  - `*ColorLib`: 当前 ColorLib 对象

#### WithBold

```go
func (c *ColorLib) WithBold(enabled bool) *ColorLib
```

WithBold 设置是否启用粗体输出（链式调用）。

- 参数:
  - `enabled`: 是否启用粗体输出（`true` - 启用，`false` - 禁用）
- 返回值:
  - `*ColorLib`: 当前 ColorLib 对象

#### WithColor

```go
func (c *ColorLib) WithColor(enabled bool) *ColorLib
```

WithColor 设置是否启用颜色输出（链式调用）。

- 参数:
  - `enabled`: 是否启用颜色输出（`true` - 启用，`false` - 禁用）
- 返回值:
  - `*ColorLib`: 当前 ColorLib 对象

#### WithUnderline

```go
func (c *ColorLib) WithUnderline(enabled bool) *ColorLib
```

WithUnderline 设置是否启用下划线输出（链式调用）。

- 参数:
  - `enabled`: 是否启用下划线输出（`true` - 启用，`false` - 禁用）
- 返回值:
  - `*ColorLib`: 当前 ColorLib 对象

#### WithWriter

```go
func (c *ColorLib) WithWriter(w io.Writer) *ColorLib
```

WithWriter 创建一个使用指定输出接口的新实例（不可变设计）。

- 参数:
  - `w`: 输出接口
- 返回值:
  - `*ColorLib`: 新的 ColorLib 实例

#### Yellow

```go
func (c *ColorLib) Yellow(msg ...any)
```

Yellow 方法用于将传入的参数以黄色文本形式打印到控制台。

- 参数:
  - `msg`: 可变参数，要打印的消息内容

#### Yellowf

```go
func (c *ColorLib) Yellowf(format string, a ...any)
```

Yellowf 方法用于将传入的参数以黄色文本形式打印到控制台。

- 参数:
  - `format`: 格式化字符串，用于指定输出的格式
  - `a`: 可变参数，用于填充格式化字符串中的占位符

#### GetColor

```go
func (c *ColorLib) GetColor() bool
```

GetColor 获取颜色启用状态。

- 返回值:
  - `bool`: 颜色启用状态

#### GetBold

```go
func (c *ColorLib) GetBold() bool
```

GetBold 获取粗体启用状态。

- 返回值:
  - `bool`: 粗体启用状态

#### GetUnderline

```go
func (c *ColorLib) GetUnderline() bool
```

GetUnderline 获取下划线启用状态。

- 返回值:
  - `bool`: 下划线启用状态

#### GetBlink

```go
func (c *ColorLib) GetBlink() bool
```

GetBlink 获取闪烁启用状态。

- 返回值:
  - `bool`: 闪烁启用状态
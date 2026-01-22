# Package color

Package color 提供了颜色管理功能。该文件实现了 `ColorManager` 结构体，用于管理 ANSI 颜色代码和颜色名称之间的映射关系，支持标准颜色和亮色系列，提供颜色代码与名称的双向转换功能。

## CONSTANTS

```go
const (
	// 标准颜色 (30-37)
	Black   = 30 // 黑色
	Red     = 31 // 红色
	Green   = 32 // 绿色
	Yellow  = 33 // 黄色
	Blue    = 34 // 蓝色
	Magenta = 35 // 品红色
	Cyan    = 36 // 青色
	White   = 37 // 白色
	Gray    = 90 // 灰色

	// 亮色 (90-97) - 统一用 Bright 前缀
	BrightRed     = 91 // 亮红色
	BrightGreen   = 92 // 亮绿色
	BrightYellow  = 93 // 亮黄色
	BrightBlue    = 94 // 亮蓝色
	BrightMagenta = 95 // 亮品红色
	BrightCyan    = 96 // 亮青色
	BrightWhite   = 97 // 亮白色
)
```

## TYPES

### ColorManager

```go
type ColorManager struct {
	// Has unexported fields.
}
```

`ColorManager` 颜色管理。

### Functions

#### NewColorManager

```go
func NewColorManager() *ColorManager
```

`NewColorManager` 创建颜色管理器。

- 返回值:
  - `*ColorManager`: 颜色管理器实例

### Methods

#### GetColorCode

```go
func (cm *ColorManager) GetColorCode(name string) (int, bool)
```

`GetColorCode` 根据颜色名称获取颜色代码。

- 参数:
  - `name`: 颜色名称
- 返回值:
  - `code`: 颜色代码
  - `ok`: 是否成功找到颜色代码

#### GetColorCodeString

```go
func (cm *ColorManager) GetColorCodeString(code int) (string, bool)
```

`GetColorCodeString` 根据颜色代码获取颜色代码字符串。

- 参数:
  - `code`: 颜色代码
- 返回值:
  - `string`: 颜色代码字符串
  - `ok`: 是否成功找到颜色代码字符串

#### GetColorName

```go
func (cm *ColorManager) GetColorName(code int) (string, bool)
```

`GetColorName` 根据颜色代码获取颜色名称。

- 参数:
  - `code`: 颜色代码
- 返回值:
  - `name`: 颜色名称
  - `ok`: 是否成功找到颜色名称

#### IsColorCode

```go
func (cm *ColorManager) IsColorCode(code int) bool
```

`IsColorCode` 判断是否为颜色代码。

- 参数:
  - `code`: 颜色代码
- 返回值:
  - `ok`: 是否为颜色代码

#### IsColorName

```go
func (cm *ColorManager) IsColorName(name string) bool
```

`IsColorName` 判断是否为颜色名称。

- 参数:
  - `name`: 颜色名称
- 返回值:
  - `ok`: 是否为颜色名称
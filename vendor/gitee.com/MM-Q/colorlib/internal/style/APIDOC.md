# Package style

Package style 提供了样式配置管理功能。该文件实现了 `StyleConfig` 结构体，用于统一管理颜色、粗体、下划线、闪烁等样式配置，使用原子操作确保线程安全的样式状态管理。

Package style 提供了样式配置的工具函数。该文件实现了 `StyleConfig` 结构体的方法，包括样式状态的设置和获取功能，支持颜色、粗体、下划线、闪烁等样式属性的线程安全操作。

## TYPES

### StyleConfig

```go
type StyleConfig struct {
	// Has unexported fields.
}
```

`StyleConfig` 统一管理所有样式配置。

### Functions

#### NewStyleConfig

```go
func NewStyleConfig() *StyleConfig
```

`NewStyleConfig` 创建一个新的样式配置实例。

- 返回值:
  - `*StyleConfig`: 新创建的样式配置实例

### Methods

#### GetBlink

```go
func (s *StyleConfig) GetBlink() bool
```

`GetBlink` 获取闪烁启用状态。

- 返回值:
  - `bool`: 闪烁启用状态

#### GetBold

```go
func (s *StyleConfig) GetBold() bool
```

`GetBold` 获取加粗启用状态。

- 返回值:
  - `bool`: 加粗启用状态

#### GetColor

```go
func (s *StyleConfig) GetColor() bool
```

`GetColor` 获取颜色启用状态。

- 返回值:
  - `bool`: 颜色启用状态

#### GetUnderline

```go
func (s *StyleConfig) GetUnderline() bool
```

`GetUnderline` 获取下划线启用状态。

- 返回值:
  - `bool`: 下划线启用状态

#### SetBlink

```go
func (s *StyleConfig) SetBlink(enable bool)
```

`SetBlink` 设置闪烁启用状态。

- 参数:
  - `enable`: 布尔值，表示是否启用闪烁

#### SetBold

```go
func (s *StyleConfig) SetBold(enable bool)
```

`SetBold` 设置加粗启用状态。

- 参数:
  - `enable`: 布尔值，表示是否启用加粗

#### SetColor

```go
func (s *StyleConfig) SetColor(enable bool)
```

`SetColor` 设置颜色启用状态。

- 参数:
  - `enable`: 布尔值，表示是否启用颜色

#### SetUnderline

```go
func (s *StyleConfig) SetUnderline(enable bool)
```

`SetUnderline` 设置下划线启用状态。

- 参数:
  - `enable`: 布尔值，表示是否启用下划线

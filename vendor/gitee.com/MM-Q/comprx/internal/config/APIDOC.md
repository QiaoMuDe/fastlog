# Package config

Package config 提供了压缩器的配置管理功能，定义了压缩器的核心配置结构体，并提供了配置的创建和压缩等级转换等实用功能。

## 主要类型

- **Config**: 压缩器配置结构体

## 主要功能

- 创建默认配置
- 压缩等级转换
- 进度显示配置
- 文件过滤配置
- 路径验证配置

## 使用示例

```go
// 创建默认配置
cfg := config.New()

// 设置压缩等级
cfg.CompressionLevel = types.CompressionLevelBest

// 启用文件覆盖
cfg.OverwriteExisting = true
```

## FUNCTIONS

### GetCompressionLevel

```go
func GetCompressionLevel(level types.CompressionLevel) int
```

- **描述**: 根据配置返回对应的压缩等级
- **参数**:
  - `level`: `types.CompressionLevel` - 压缩等级
- **返回值**:
  - `int` - 压缩等级

## TYPES

### Config

```go
type Config struct {
    CompressionLevel      types.CompressionLevel // 压缩等级
    OverwriteExisting     bool                   // 是否覆盖已存在的文件
    Progress              *progress.Progress     // 进度显示
    DisablePathValidation bool                   // 是否禁用路径验证
    Filter                *types.FilterOptions   // 文件过滤配置
}
```

- **描述**: 压缩器配置

### New

```go
func New() *Config
```

- **描述**: 创建新的压缩器配置
- **返回**:
  - `*Config`: 新的压缩器配置实例

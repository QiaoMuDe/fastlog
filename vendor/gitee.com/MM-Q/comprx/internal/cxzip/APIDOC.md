# Package cxzip

Package cxzip 提供了 ZIP 格式的压缩包内容列表、压缩和解压缩功能的实现。ZIP 格式支持多种文件类型，包括普通文件、目录和符号链接，并提供完整的压缩信息。该包支持进度显示、文件过滤、路径安全验证和配置化的压缩与解压缩功能。

## ZIP 压缩包内容列表功能

### 主要功能

- **ZIP 压缩包完整文件列表获取**
- **限制数量的文件列表获取**
- **模式匹配的文件列表过滤**
- **多种文件类型支持（普通文件、目录、符号链接）**
- **完整的压缩信息统计**

### 文件类型支持

- **普通文件**: 完整的文件信息和压缩信息
- **目录**: 目录标识和权限信息
- **符号链接**: 链接目标路径读取和保存

### 压缩信息

- 原始文件大小
- 压缩后大小
- 压缩率计算
- 文件修改时间
- 文件权限模式

### 性能优化

- 限制模式下的容量预分配
- 符号链接目标的安全读取
- 错误容忍的链接目标处理

### 使用示例

```go
// 获取 ZIP 文件完整列表
info, err := cxzip.ListZip("archive.zip")

// 获取前 10 个文件信息
info, err := cxzip.ListZipLimit("archive.zip", 10)

// 获取匹配 *.go 模式的文件
info, err := cxzip.ListZipMatch("archive.zip", "*.go")
```

## ZIP 格式的解压缩功能

### 主要功能

- **ZIP 格式解压缩**
- **支持普通文件、目录、符号链接解压**
- **文件过滤功能**
- **进度显示支持**
- **路径安全验证**
- **文件覆盖控制**

### 安全特性

- 路径遍历攻击防护
- 安全的文件路径验证
- 可配置的路径验证开关

### 支持的文件类型

- **普通文件**: 完整解压文件内容
- **目录**: 创建目录结构
- **符号链接**: 重建符号链接
- **空文件**: 创建空文件

### 使用示例

```go
// 创建配置
cfg := config.New()
cfg.OverwriteExisting = true

// 解压文件
err := cxzip.Unzip("archive.zip", "output_dir", cfg)
```

## ZIP 格式的压缩功能

### 主要功能

- **ZIP 格式压缩**
- **支持文件和目录压缩**
- **符号链接和特殊文件处理**
- **文件过滤功能**
- **进度显示支持**
- **可配置的压缩等级**

### 支持的文件类型

- **普通文件**: 使用配置的压缩方法
- **目录**: 创建目录条目
- **符号链接**: 保存链接目标
- **特殊文件**: 创建占位符条目

### 使用示例

```go
// 创建配置
cfg := config.New()
cfg.CompressionLevel = types.CompressionLevelBest

// 压缩文件
err := cxzip.Zip("output.zip", "input_dir", cfg)
```

## FUNCTIONS

### ListZip

```go
func ListZip(archivePath string) (*types.ArchiveInfo, error)
```

- **描述**: 获取 ZIP 压缩包的所有文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListZipLimit

```go
func ListZipLimit(archivePath string, limit int) (*types.ArchiveInfo, error)
```

- **描述**: 获取 ZIP 压缩包指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制返回的文件数量
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListZipMatch

```go
func ListZipMatch(archivePath string, pattern string) (*types.ArchiveInfo, error)
```

- **描述**: 获取 ZIP 压缩包中匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### Unzip

```go
func Unzip(zipFilePath string, targetDir string, cfg *config.Config) error
```

- **描述**: 解压缩 ZIP 文件到指定目录
- **参数**:
  - `zipFilePath`: 要解压缩的 ZIP 文件路径
  - `targetDir`: 解压缩后的目标目录路径
  - `cfg`: 解压缩配置
- **返回**:
  - `error`: 解压缩过程中发生的错误

### Zip

```go
func Zip(dst string, src string, cfg *config.Config) error
```

- **描述**: 创建 ZIP 压缩文件
- **参数**:
  - `dst`: 生成的 ZIP 文件路径
  - `src`: 需要压缩的源路径
  - `cfg`: 压缩配置指针
- **返回**:
  - `error`: 操作过程中遇到的错误
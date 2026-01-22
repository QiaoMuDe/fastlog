# Package core

Package core 提供了压缩库的核心功能实现，实现了压缩库的主要业务逻辑，包括压缩和解压缩操作的统一接口。支持多种压缩格式的自动检测和处理，提供了配置化的压缩解压缩功能。

## 主要类型

- **Comprx**: 压缩器核心结构体

## 主要功能

- 统一的压缩和解压缩接口
- 自动压缩格式检测
- 支持多种压缩格式（ZIP、TAR、TGZ、GZIP、BZIP2、ZLIB）
- 配置化的压缩参数
- 文件和目录的智能处理

## 支持的压缩格式

- **ZIP**: `.zip`
- **TAR**: `.tar`
- **TGZ**: `.tgz`, `.tar.gz`
- **GZIP**: `.gz`
- **BZIP2**: `.bz2`, `.bzip2`（仅解压）
- **ZLIB**: `.zlib`

## 使用示例

```go
// 创建压缩器实例
comprx := core.New()

// 压缩文件
err := comprx.Pack("output.zip", "input_dir")

// 解压文件
err := comprx.Unpack("archive.zip", "output_dir")
```

## Package core 提供压缩包内容列表功能的核心实现

该文件实现了压缩包内容查看的统一接口，支持多种压缩格式的文件列表功能。提供了完整列表、限制数量列表和模式匹配列表等功能。

### 主要功能

- 列出压缩包内所有文件信息
- 支持限制返回文件数量
- 支持文件名模式匹配过滤
- 自动检测压缩格式
- 统一的错误处理

### 支持的压缩格式

- **ZIP**: `.zip`
- **TAR**: `.tar`
- **TGZ**: `.tgz`, `.tar.gz`
- **GZIP**: `.gz`
- **BZIP2**: `.bz2`, `.bzip2`
- **ZLIB**: `.zlib`

### 使用示例

```go
// 列出所有文件
info, err := core.List("archive.zip")

// 列出前10个文件
info, err := core.ListLimit("archive.zip", 10)

// 列出匹配模式的文件
info, err := core.ListMatch("archive.zip", "*.go")
```

## VARIABLES

### New

```go
var New = NewComprx
```

- **描述**: 创建压缩器实例（`NewComprx` 的别名）
- **返回**:
  - `*Comprx`: 压缩器实例

## FUNCTIONS

### List

```go
func List(archivePath string) (*types.ArchiveInfo, error)
```

- **描述**: 列出压缩包的所有文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListLimit

```go
func ListLimit(archivePath string, limit int) (*types.ArchiveInfo, error)
```

- **描述**: 列出指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制返回的文件数量
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListMatch

```go
func ListMatch(archivePath string, pattern string) (*types.ArchiveInfo, error)
```

- **描述**: 列出匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

## TYPES

### Comprx

```go
type Comprx struct {
    Config *config.Config // 压缩器配置
}
```

- **描述**: 压缩器

### NewComprx

```go
func NewComprx() *Comprx
```

- **描述**: 创建压缩器实例
- **返回**:
  - `*Comprx`: 压缩器实例

### Pack

```go
func (c *Comprx) Pack(dst string, src string) error
```

- **描述**: 压缩文件或目录
- **参数**:
  - `dst`: 目标文件路径
  - `src`: 源文件路径
- **返回**:
  - `error`: 错误信息

### Unpack

```go
func (c *Comprx) Unpack(src string, dst string) error
```

- **描述**: 解压文件
- **参数**:
  - `src`: 源文件路径
  - `dst`: 目标目录路径
- **返回**:
  - `error`: 错误信息

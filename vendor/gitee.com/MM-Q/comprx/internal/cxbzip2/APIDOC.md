# Package cxbzip2

Package cxbzip2 提供了 BZIP2 格式的压缩包内容列表功能和解压缩功能的实现。该包实现了 BZIP2 格式压缩包的文件信息获取功能，包括基本列表、限制数量列表和模式匹配列表。同时，也实现了 BZIP2 格式的单文件解压缩操作，支持进度显示和智能目标路径处理。

## 主要功能

### 压缩包内容列表功能

- **BZIP2 压缩包文件信息获取**
- **原始文件名智能推导**
- **原始文件大小计算**
- **模式匹配过滤**
- **压缩率计算**

### 特殊处理

- 智能推导原始文件名（去除 `.bz2` 或 `.bzip2` 后缀）
- 文件名缺失时使用默认后缀
- 通过完整读取计算原始文件大小
- 使用压缩文件修改时间（BZIP2 不保存原始时间）
- 使用默认文件权限（BZIP2 不保存权限信息）

### 格式限制

- BZIP2 格式不保存文件元数据
- 无原始文件名信息
- 无修改时间信息
- 无文件权限信息

### 性能优化

- 使用缓冲区池减少内存分配
- 高效的文件大小计算方法
- 错误容忍的大小估算机制

### 使用示例

```go
// 获取 BZIP2 文件信息
info, err := cxbzip2.ListBz2("archive.bz2")

// 获取匹配模式的文件信息
info, err := cxbzip2.ListBz2Match("archive.bzip2", "*.txt")

// 限制返回文件数量（对 BZIP2 无实际效果）
info, err := cxbzip2.ListBz2Limit("archive.bz2", 10)
```

## 解压缩功能

### 主要功能

- **BZIP2 格式单文件解压缩**
- **进度显示支持**
- **智能目标路径处理**
- **文件覆盖控制**
- **路径安全验证**

### 智能处理

- 目标为目录时自动生成文件名
- 自动去除 `.bz2` 和 `.bzip2` 扩展名作为文件名
- 自动创建目标文件的父目录
- 路径安全验证防止路径遍历攻击

### 解压特性

- 使用 Burrows-Wheeler 变换算法
- 高压缩比，适合文本文件
- 流式解压，内存占用低
- 支持大文件解压

### 性能优化

- 智能缓冲区大小选择
- 进度条模式下的大小预计算
- 错误容忍的大小估算机制

### 使用示例

```go
// 创建配置
cfg := config.New()
cfg.OverwriteExisting = true

// 解压文件到指定路径
err := cxbzip2.Unbz2("archive.bz2", "output.txt", cfg)

// 解压文件到目录（自动生成文件名）
err := cxbzip2.Unbz2("archive.bzip2", "output_dir/", cfg)
```

## FUNCTIONS

### ListBz2

```go
func ListBz2(archivePath string) (*types.ArchiveInfo, error)
```

- **描述**: 获取 BZ2 压缩包的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListBz2Limit

```go
func ListBz2Limit(archivePath string, limit int) (*types.ArchiveInfo, error)
```

- **描述**: 获取 BZ2 压缩包指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制返回的文件数量
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListBz2Match

```go
func ListBz2Match(archivePath string, pattern string) (*types.ArchiveInfo, error)
```

- **描述**: 获取 BZ2 压缩包中匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### Unbz2

```go
func Unbz2(bz2FilePath string, targetPath string, cfg *config.Config) error
```

- **描述**: 解压缩 BZIP2 文件
- **参数**:
  - `bz2FilePath`: 要解压缩的 BZIP2 文件路径
  - `targetPath`: 解压缩后的目标文件路径
  - `cfg`: 解压缩配置
- **返回**:
  - `error`: 解压缩过程中发生的错误

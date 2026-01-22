# Package cxtgz

Package cxtgz 提供了 TGZ (tar.gz) 格式的压缩包内容列表、压缩和解压缩功能的实现。TGZ 是 TAR 归档格式与 GZIP 压缩的组合，提供了高效的多文件压缩能力。该包支持多种文件类型的处理，包括普通文件、目录、符号链接和硬链接，同时提供完整的进度显示、文件过滤和路径安全验证功能。

## TGZ 压缩包内容列表功能

### 主要功能

- **TGZ 压缩包完整文件列表获取**
- **限制数量的文件列表获取**
- **模式匹配的文件列表过滤**
- **多种文件类型支持（普通文件、目录、符号链接、硬链接）**
- **文件元数据完整保存**

### 文件类型支持

- **普通文件**: 完整的文件信息
- **目录**: 目录标识和权限信息
- **符号链接**: 链接目标路径保存
- **硬链接**: 链接目标路径保存

### 元数据信息

- 文件名和路径
- 文件原始大小
- 修改时间
- 文件权限模式
- 文件类型标识
- 符号链接目标

### 压缩特性

- **整体压缩**: TGZ 对整个 TAR 归档进行压缩
- **单个文件压缩大小无法准确计算**
- **提供整体压缩包大小信息**

### 使用示例

```go
// 获取 TGZ 文件完整列表
info, err := cxtgz.ListTgz("archive.tar.gz")

// 获取前 10 个文件信息
info, err := cxtgz.ListTgzLimit("archive.tar.gz", 10)

// 获取匹配 *.go 模式的文件
info, err := cxtgz.ListTgzMatch("archive.tar.gz", "*.go")
```

## TGZ 格式的压缩功能

### 主要功能

- **TGZ 格式文件和目录压缩**
- **支持多种文件类型（普通文件、目录、符号链接、特殊文件）**
- **可配置的压缩等级**
- **进度显示支持**
- **文件过滤功能**
- **文件覆盖控制**
- **相对路径处理**

### 压缩流程

1. 创建 GZIP 压缩流
2. 在 GZIP 流上创建 TAR 归档流
3. 将文件按 TAR 格式写入并通过 GZIP 压缩

### 文件类型支持

- **普通文件**: 完整内容压缩
- **目录**: 创建目录条目
- **符号链接**: 保存链接目标
- **特殊文件**: 保存文件元数据

### 路径处理

- 自动转换为 TAR 标准路径格式（正斜杠）
- 保留目录结构的相对路径
- 支持单文件和目录压缩

### 使用示例

```go
// 创建配置
cfg := config.New()
cfg.CompressionLevel = types.CompressionLevelBest
cfg.OverwriteExisting = true

// 压缩目录
err := cxtgz.Tgz("archive.tar.gz", "source_dir", cfg)

// 压缩单个文件
err := cxtgz.Tgz("file.tar.gz", "single_file.txt", cfg)
```

## TGZ 格式的解压缩功能

### 主要功能

- **TGZ 格式文件和目录解压缩**
- **支持多种文件类型（普通文件、目录、符号链接、硬链接）**
- **进度显示支持**
- **路径安全验证**
- **文件过滤功能**
- **文件覆盖控制**

### 解压流程

1. 打开 TGZ 文件
2. 创建 GZIP 解压缩流
3. 在 GZIP 流上创建 TAR 读取流
4. 按 TAR 格式解析并解压文件

### 文件类型支持

- **普通文件**: 完整内容解压
- **目录**: 创建目录结构
- **符号链接**: 恢复链接关系
- **硬链接**: 创建硬链接
- **其他类型**: 跳过处理并提示

### 安全特性

- 路径遍历攻击防护
- 文件路径验证
- 可配置的路径验证开关
- 文件覆盖保护

### 性能优化

- 智能缓冲区大小选择
- 空文件特殊处理
- 进度条模式下的大小预计算

### 使用示例

```go
// 创建配置
cfg := config.New()
cfg.OverwriteExisting = true

// 解压 TGZ 文件
err := cxtgz.Untgz("archive.tar.gz", "output_dir", cfg)
```

## FUNCTIONS

### ListTgz

```go
func ListTgz(archivePath string) (*types.ArchiveInfo, error)
```

- **描述**: 获取 TGZ 压缩包的所有文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListTgzLimit

```go
func ListTgzLimit(archivePath string, limit int) (*types.ArchiveInfo, error)
```

- **描述**: 获取 TGZ 压缩包指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制返回的文件数量
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListTgzMatch

```go
func ListTgzMatch(archivePath string, pattern string) (*types.ArchiveInfo, error)
```

- **描述**: 获取 TGZ 压缩包中匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### Tgz

```go
func Tgz(dst string, src string, cfg *config.Config) error
```

- **描述**: 创建 TGZ (tar.gz) 压缩文件
- **参数**:
  - `dst`: 生成的 TGZ 文件路径
  - `src`: 需要压缩的源路径
  - `cfg`: 压缩配置指针
- **返回**:
  - `error`: 操作过程中遇到的错误

### Untgz

```go
func Untgz(tgzFilePath string, targetDir string, cfg *config.Config) error
```

- **描述**: 解压缩 TGZ 文件到指定目录
- **参数**:
  - `tgzFilePath`: 要解压缩的 TGZ 文件路径
  - `targetDir`: 解压缩后的目标目录路径
  - `cfg`: 解压缩配置
- **返回**:
  - `error`: 解压缩过程中发生的错误
# Package cxzlib

Package cxzlib 提供了 ZLIB 格式的压缩、解压缩以及压缩包内容列表功能的实现。ZLIB 格式是一种单文件压缩格式，支持多种压缩等级和高效的内存管理。该包提供了内存压缩和流式压缩功能，同时也支持文件的压缩和解压缩操作。

## ZLIB 压缩包内容列表功能

### 主要功能

- **ZLIB 压缩包文件信息获取**
- **原始文件名智能推导**
- **原始文件大小计算**
- **模式匹配过滤**
- **压缩率计算**

### 特殊处理

- 智能推导原始文件名（去除 `.zlib` 后缀）
- 文件名缺失时使用默认后缀
- 通过完整读取计算原始文件大小
- 使用压缩文件修改时间（ZLIB 不保存原始时间）
- 使用默认文件权限（ZLIB 不保存权限信息）

### 格式限制

- ZLIB 格式不保存文件元数据
- 无原始文件名信息
- 无修改时间信息
- 无文件权限信息

### 性能优化

- 使用缓冲区池减少内存分配
- 高效的文件大小计算方法
- 错误容忍的大小估算机制

### 使用示例

```go
// 获取 ZLIB 文件信息
info, err := cxzlib.ListZlib("archive.zlib")

// 获取匹配模式的文件信息
info, err := cxzlib.ListZlibMatch("archive.zlib", "*.txt")

// 限制返回文件数量（对 ZLIB 无实际效果）
info, err := cxzlib.ListZlibLimit("archive.zlib", 10)
```

## ZLIB 内存压缩和流式压缩功能

### 主要功能

- **ZLIB 内存压缩：字节数组和字符串的压缩解压**
- **ZLIB 流式压缩：支持 `io.Reader` 和 `io.Writer` 接口**
- **支持自定义压缩等级**
- **优化的内存分配策略**
- **完善的错误处理和资源管理**

### 压缩特性

- 使用 DEFLATE 压缩算法
- 包含 Adler-32 校验和
- 比 GZIP 格式更紧凑（无文件头信息）
- 支持多种压缩等级

### 性能优化

- 预分配缓冲区减少内存重分配
- 智能估算压缩后大小
- 直接字节操作避免额外拷贝
- 自动资源清理防止内存泄漏

### 使用示例

```go
// 压缩字节数据
compressed, err := cxzlib.CompressBytes(data, types.CompressionLevelBest)

// 解压字节数据
decompressed, err := cxzlib.DecompressBytes(compressed)

// 压缩字符串
compressed, err := cxzlib.CompressString("hello world", types.CompressionLevelDefault)

// 流式压缩
err := cxzlib.CompressStream(dst, src, types.CompressionLevelFast)
```

## ZLIB 解压缩功能

### 主要功能

- **ZLIB 格式单文件解压缩**
- **进度显示支持**
- **智能目标路径处理**
- **文件覆盖控制**
- **高效的缓冲区管理**

### 智能处理

- 目标为目录时自动生成文件名
- 自动去除 `.zlib` 扩展名作为文件名
- 自动创建目标文件的父目录

### 解压特性

- 使用 DEFLATE 解压算法
- 包含 Adler-32 校验和验证
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
err := cxzlib.Unzlib("archive.zlib", "output.txt", cfg)

// 解压文件到目录（自动生成文件名）
err := cxzlib.Unzlib("archive.zlib", "output_dir/", cfg)
```

## ZLIB 压缩功能

### 主要功能

- **ZLIB 格式单文件压缩**
- **可配置的压缩等级**
- **进度显示支持**
- **文件覆盖控制**
- **高效的缓冲区管理**

### 限制

- 只支持单个文件压缩
- 不支持目录压缩
- 不支持多文件打包
- 不保存文件元数据（文件名、修改时间等）

### 压缩特性

- 使用 DEFLATE 压缩算法
- 支持多种压缩等级
- 包含 Adler-32 校验和
- 比 GZIP 格式更紧凑（无文件头信息）

### 使用示例

```go
// 创建配置
cfg := config.New()
cfg.CompressionLevel = types.CompressionLevelBest

// 压缩单个文件
err := cxzlib.Zlib("output.zlib", "input.txt", cfg)
```

## FUNCTIONS

### CompressBytes

```go
func CompressBytes(data []byte, level types.CompressionLevel) (result []byte, err error)
```

- **描述**: 压缩字节数据到内存
- **参数**:
  - `data`: 要压缩的字节数据
  - `level`: 压缩级别
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息

### CompressStream

```go
func CompressStream(dst io.Writer, src io.Reader, level types.CompressionLevel) (err error)
```

- **描述**: 流式压缩数据
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器
  - `level`: 压缩级别
- **返回**:
  - `error`: 错误信息

### CompressString

```go
func CompressString(text string, level types.CompressionLevel) ([]byte, error)
```

- **描述**: 压缩字符串到内存
- **参数**:
  - `text`: 要压缩的字符串
  - `level`: 压缩级别
- **返回**:
  - `[]byte`: 厎压后的数据
  - `error`: 错误信息

### DecompressBytes

```go
func DecompressBytes(compressedData []byte) (result []byte, err error)
```

- **描述**: 从内存解压字节数据
- **参数**:
  - `compressedData`: 压缩的字节数据
- **返回**:
  - `[]byte`: 解压后的数据
  - `error`: 错误信息

### DecompressStream

```go
func DecompressStream(dst io.Writer, src io.Reader) (err error)
```

- **描述**: 流式解压数据
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器（压缩数据）
- **返回**:
  - `error`: 错误信息

### DecompressString

```go
func DecompressString(compressedData []byte) (string, error)
```

- **描述**: 从内存解压为字符串
- **参数**:
  - `compressedData`: 压缩的字节数据
- **返回**:
  - `string`: 解压后的字符串
  - `error`: 错误信息

### ListZlib

```go
func ListZlib(archivePath string) (*types.ArchiveInfo, error)
```

- **描述**: 获取 ZLIB 压缩包的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListZlibLimit

```go
func ListZlibLimit(archivePath string, limit int) (*types.ArchiveInfo, error)
```

- **描述**: 获取 ZLIB 压缩包指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制返回的文件数量
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListZlibMatch

```go
func ListZlibMatch(archivePath string, pattern string) (*types.ArchiveInfo, error)
```

- **描述**: 获取 ZLIB 压缩包中匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### Unzlib

```go
func Unzlib(zlibFilePath string, targetPath string, config *config.Config) error
```

- **描述**: 解压缩 ZLIB 文件
- **参数**:
  - `zlibFilePath`: 要解压缩的 ZLIB 文件路径
  - `targetPath`: 解压缩后的目标文件路径
  - `config`: 解压缩配置
- **返回**:
  - `error`: 解压缩过程中发生的错误

### Zlib

```go
func Zlib(dst string, src string, cfg *config.Config) error
```

- **描述**: 压缩单个文件为 ZLIB 格式
- **参数**:
  - `dst`: 生成的 ZLIB 文件路径
  - `src`: 需要压缩的源文件路径
  - `cfg`: 压缩配置指针
- **返回**:
  - `error`: 操作过程中遇到的错误
# Package cxgzip

Package cxgzip 提供了 GZIP 格式的压缩、解压缩以及压缩包内容列表功能的实现。该包支持单文件压缩和解压缩操作，提供可配置的压缩等级、进度显示、文件元数据保存和路径安全验证等功能。

## GZIP 压缩功能

### 主要功能

- **GZIP 格式单文件压缩**
- **可配置的压缩等级**
- **进度显示支持**
- **文件元数据保存（文件名、修改时间）**
- **文件覆盖控制**

### 限制

- 只支持单个文件压缩
- 不支持目录压缩
- 不支持多文件打包

### 使用示例

```go
// 创建配置
cfg := config.New()
cfg.CompressionLevel = types.CompressionLevelBest

// 压缩单个文件
err := cxgzip.Gzip("output.gz", "input.txt", cfg)
```

## GZIP 压缩包内容列表功能

### 主要功能

- **GZIP 压缩包文件信息获取**
- **原始文件名和大小计算**
- **文件修改时间获取**
- **模式匹配过滤**
- **压缩率计算**

### 特殊处理

- 自动从 GZIP 文件头获取原始文件名
- 文件名缺失时智能推导（去除 `.gz` 后缀）
- 通过完整读取计算原始文件大小
- 使用默认文件权限（GZIP 不保存权限信息）

### 性能优化

- 使用缓冲区池减少内存分配
- 高效的文件大小计算方法
- 错误容忍的大小估算机制

### 使用示例

```go
// 获取 GZIP 文件信息
info, err := cxgzip.ListGzip("archive.gz")

// 获取匹配模式的文件信息
info, err := cxgzip.ListGzipMatch("archive.gz", "*.txt")

// 限制返回文件数量（对 GZIP 无实际效果）
info, err := cxgzip.ListGzipLimit("archive.gz", 10)
```

## GZIP 内存压缩和流式压缩功能

### 主要功能

- **GZIP 内存压缩：字节数组和字符串的压缩解压**
- **GZIP 流式压缩：支持 `io.Reader` 和 `io.Writer` 接口**
- **支持自定义压缩等级**
- **优化的内存分配策略**
- **完善的错误处理和资源管理**

### 性能优化

- 预分配缓冲区减少内存重分配
- 智能估算压缩后大小
- 直接字节操作避免额外拷贝
- 自动资源清理防止内存泄漏

### 使用示例

```go
// 压缩字节数据
compressed, err := cxgzip.CompressBytes(data, types.CompressionLevelBest)

// 解压字节数据
decompressed, err := cxgzip.DecompressBytes(compressed)

// 压缩字符串
compressed, err := cxgzip.CompressString("hello world", types.CompressionLevelDefault)

// 流式压缩
err := cxgzip.CompressStream(dst, src, types.CompressionLevelFast)
```

## GZIP 解压缩功能

### 主要功能

- **GZIP 格式单文件解压缩**
- **进度显示支持**
- **文件元数据恢复（文件名、修改时间）**
- **路径安全验证**
- **文件覆盖控制**
- **智能目标路径处理**

### 安全特性

- 路径遍历攻击防护
- GZIP 文件头文件名验证
- 可配置的路径验证开关

### 智能处理

- 自动从 GZIP 文件头获取原始文件名
- 目标为目录时自动生成文件名
- 自动去除 `.gz` 扩展名作为备选文件名

### 使用示例

```go
// 创建配置
cfg := config.New()
cfg.OverwriteExisting = true

// 解压文件到指定路径
err := cxgzip.Ungzip("archive.gz", "output.txt", cfg)

// 解压文件到目录（自动生成文件名）
err := cxgzip.Ungzip("archive.gz", "output_dir/", cfg)
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
  - `[]byte`: 压缩后的数据
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

### Gzip

```go
func Gzip(dst string, src string, cfg *config.Config) error
```

- **描述**: 压缩单个文件为 GZIP 格式
- **参数**:
  - `dst`: 生成的 GZIP 文件路径
  - `src`: 需要压缩的源文件路径
  - `cfg`: 压缩配置指针
- **返回**:
  - `error`: 操作过程中遇到的错误

### ListGzip

```go
func ListGzip(archivePath string) (*types.ArchiveInfo, error)
```

- **描述**: 获取 GZIP 压缩包的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListGzipLimit

```go
func ListGzipLimit(archivePath string, limit int) (*types.ArchiveInfo, error)
```

- **描述**: 获取 GZIP 压缩包指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制返回的文件数量
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListGzipMatch

```go
func ListGzipMatch(archivePath string, pattern string) (*types.ArchiveInfo, error)
```

- **描述**: 获取 GZIP 压缩包中匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `*types.ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### Ungzip

```go
func Ungzip(gzipFilePath string, targetPath string, config *config.Config) error
```

- **描述**: 解压缩 GZIP 文件
- **参数**:
  - `gzipFilePath`: 要解压缩的 GZIP 文件路径
  - `targetPath`: 解压缩后的目标文件路径
  - `config`: 解压缩配置
- **返回**:
  - `error`: 解压缩过程中发生的错误

# Package utils

Package utils 提供了多种实用工具函数，包括缓冲区管理、文件和目录大小计算、路径处理、文件列表处理和格式化显示等功能。这些工具函数被压缩库的各个模块广泛使用，提供了统一的基础功能。

## 文件和目录大小计算

### 主要功能

- 获取单个文件的大小
- 递归计算目录的总大小
- 提供安全版本（出错返回 0）和详细版本（返回错误信息）
- 自动忽略符号链接等特殊文件
- 详细的错误分类和处理

### 错误处理

- 文件不存在错误
- 权限不足错误
- 其他系统错误
- 遍历过程中的动态错误处理

### 使用示例

```go
// 安全版本，出错时返回 0
size := utils.GetSizeOrZero("./mydir")

// 详细版本，返回错误信息
size, err := utils.GetSize("./myfile.txt")
if err != nil {
    log.Printf("获取大小失败: %v", err)
}
```

## 文件列表处理和格式化显示

### 主要功能

- 文件大小格式化显示（B、KB、MB、GB 等）
- 文件权限格式化显示
- 文件名模式匹配（支持通配符）
- 压缩包摘要信息打印
- 文件列表格式化打印
- 文件列表过滤和限制

### 显示格式

- 简洁模式：仅显示文件名
- 详细模式：显示权限、大小、时间等完整信息
- 支持符号链接目标显示
- 自动计算压缩率

### 使用示例

```go
// 格式化文件大小
sizeStr := utils.FormatFileSize(1024 * 1024) // "1.0 MB"

// 打印压缩包摘要
utils.PrintArchiveSummary(archiveInfo)

// 打印文件列表（详细模式）
utils.PrintFileList(files, true)

// 模式匹配过滤
filtered := utils.FilterFilesByPattern(files, "*.go")
```

## 文件系统操作和路径处理

### 主要功能

- 文件和目录存在性检查
- 目录创建和确保
- 动态缓冲区大小计算
- 路径安全验证和转换
- 绝对路径处理

### 安全特性

- 路径遍历攻击防护
- 绝对路径检测
- UNC 路径和协议前缀检测
- Windows 特殊路径处理

### 使用示例

```go
// 检查文件是否存在
if utils.Exists("file.txt") {
    // 文件存在
}

// 确保目录存在
err := utils.EnsureDir("output/dir")

// 验证路径安全性
safePath, err := utils.ValidatePathSimple(targetDir, filePath, false)
```

## CONSTANTS

### 文件大小格式化相关常量

```go
const (
    // SizeUnit 文件大小计算单位 (1024 字节)
    SizeUnit = 1024
    // SizeUnitStr 文件大小单位字符串 (KB, MB, GB, TB, PB, EB)
    SizeUnitStr = "KMGTPE"
)
```

### 文件处理相关常量

```go
const (
    // DefaultBufferSize 默认缓冲区大小 (32KB)
    // 用于读取压缩文件内容时的缓冲区
    DefaultBufferSize = 32 * 1024

    // DefaultFileMode 默认文件权限 (0644)
    // 用于不保存文件权限的压缩格式 (如 GZIP, BZ2)
    DefaultFileMode = 0644
)
```

### 文件扩展名相关常量

```go
const (
    // DecompressedSuffix 解压缩文件的默认后缀
    DecompressedSuffix = ".decompressed"
)
```

### 切片预分配相关常量

```go
const (
    // DefaultFileCapacity 默认文件列表初始容量
    // 适用于 TAR/TGZ 等无法预先知道文件数量的格式
    DefaultFileCapacity = 256
)
```

## FUNCTIONS

### DetectCompressFormat

```go
func DetectCompressFormat(filename string) (types.CompressType, error)
```

- **描述**: 智能检测压缩文件格式
- **参数**:
  - `filename`: 文件名
- **返回**:
  - `types.CompressType`: 检测到的压缩格式
  - `error`: 错误信息

### EnsureAbsPath

```go
func EnsureAbsPath(path, pathType string) (string, error)
```

- **描述**: 确保路径为绝对路径，如果不是则转换为绝对路径
- **参数**:
  - `path`: 待检查的路径
  - `pathType`: 路径类型描述（用于错误信息）
- **返回**:
  - `string`: 绝对路径
  - `error`: 转换过程中的错误

### EnsureDir

```go
func EnsureDir(path string) error
```

- **描述**: 检查指定路径的目录是否存在，不存在则创建
- **参数**:
  - `path`: 要检查的目录路径
- **返回**:
  - `error`: 如果创建目录成功，则返回 `nil`，否则返回错误信息

### Exists

```go
func Exists(path string) bool
```

- **描述**: 检查指定路径的文件或目录是否存在
- **参数**:
  - `path`: 要检查的路径
- **返回**:
  - `bool`: 如果文件或目录存在，则返回 `true`，否则返回 `false`

### FilterFilesByPattern

```go
func FilterFilesByPattern(files []types.FileInfo, pattern string) []types.FileInfo
```

- **描述**: 根据模式过滤文件列表
- **参数**:
  - `files`: 文件列表
  - `pattern`: 模式字符串
- **返回**:
  - `[]types.FileInfo`: 过滤后的文件列表

### FormatFileMode

```go
func FormatFileMode(mode os.FileMode) string
```

- **描述**: 格式化文件权限显示
- **参数**:
  - `mode`: 文件权限
- **返回**:
  - `string`: 格式化后的文件权限字符串

### FormatFileSize

```go
func FormatFileSize(size int64) string
```

- **描述**: 格式化文件大小显示
- **参数**:
  - `size`: 文件大小
- **返回**:
  - `string`: 格式化后的文件大小字符串

### GetSize

```go
func GetSize(path string) (int64, error)
```

- **描述**: 获取文件或目录的大小（字节）
- **参数**:
  - `path`: 文件或目录路径
- **返回**:
  - `int64`: 文件或目录的总大小（字节）
  - `error`: 错误信息
- **注意**:
  - 如果是文件，返回文件大小
  - 如果是目录，返回目录中所有文件的总大小
  - 如果路径不存在，返回错误
  - 只计算普通文件的大小，忽略符号链接等特殊文件

### GetSizeOrZero

```go
func GetSizeOrZero(path string) int64
```

- **描述**: 获取文件或目录的大小，出错时返回 0
- **参数**:
  - `path`: 文件或目录路径
- **返回**:
  - `int64`: 文件或目录的总大小（字节），出错时返回 0
- **功能**:
  - 如果是文件，返回文件大小
  - 如果是目录，返回目录中所有普通文件的总大小
  - 忽略符号链接等特殊文件
  - 发生任何错误时返回 0，不抛出异常
- **注意**:
  - 此函数为 `GetSize` 的安全版本，适用于不需要错误处理的场景
  - 如需详细错误信息，请使用 `GetSize` 函数

### LimitFiles

```go
func LimitFiles(files []types.FileInfo, limit int) []types.FileInfo
```

- **描述**: 限制文件列表数量
- **参数**:
  - `files`: 文件列表
  - `limit`: 限制数量
- **返回**:
  - `[]types.FileInfo`: 限制后的文件列表

### MatchPattern

```go
func MatchPattern(name, pattern string) bool
```

- **描述**: 文件名模式匹配，支持简单的通配符匹配: `*` 和 `?`
- **参数**:
  - `name`: 文件名
  - `pattern`: 模式字符串
- **返回**:
  - `bool`: 是否匹配成功

### PrintArchiveSummary

```go
func PrintArchiveSummary(archiveInfo *types.ArchiveInfo)
```

- **描述**: 打印压缩包摘要信息
- **参数**:
  - `archiveInfo`: 压缩包信息

### PrintFileInfo

```go
func PrintFileInfo(info types.FileInfo, showDetails bool)
```

- **描述**: 格式化打印单个文件信息
- **参数**:
  - `info`: 文件信息
  - `showDetails`: 是否显示详细信息

### PrintFileList

```go
func PrintFileList(files []types.FileInfo, showDetails bool)
```

- **描述**: 打印文件列表
- **参数**:
  - `files`: 文件列表
  - `showDetails`: 是否显示详细信息

### ValidatePathSimple

```go
func ValidatePathSimple(targetDir, filePath string, skipValidation bool) (string, error)
```

- **描述**: 验证文件路径是否安全，防止路径遍历攻击
- **参数**:
  - `targetDir`: 目标目录
  - `filePath`: 要验证的文件路径
  - `skipValidation`: 是否跳过安全验证（警告：仅在处理可信数据时使用）
- **返回**:
  - `string`: 安全的文件路径
  - `error`: 如果路径不安全，则返回错误信息
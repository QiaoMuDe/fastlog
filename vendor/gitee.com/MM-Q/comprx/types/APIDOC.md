# Package types

Package types 提供了文件过滤功能的核心类型和接口定义，以及压缩包文件信息和压缩包整体信息的数据结构。此外，还定义了 comprx 库使用的核心数据类型和常量。

## 文件过滤功能

### 主要类型

- **FileFilter**: 文件过滤器接口
- **FilterOptions**: 过滤配置选项结构体

### 主要功能

- 支持包含和排除模式的 glob 匹配
- 支持文件大小范围过滤
- 提供高性能的快速匹配算法
- 支持复杂的 glob 模式匹配
- 提供过滤条件验证功能

### 使用示例

```go
// 创建过滤选项
filter := &types.FilterOptions{
    Include: []string{"*.go", "*.md"},
    Exclude: []string{"*_test.go"},
    MaxSize: 10 * 1024 * 1024, // 10MB
}

// 检查文件是否应该跳过
shouldSkip := filter.ShouldSkipByParams("main.go", 1024, false)
```

## 压缩包文件信息和压缩包整体信息

### 主要类型

- **FileInfo**: 压缩包内单个文件的详细信息
- **ArchiveInfo**: 压缩包的整体信息和文件列表

### 主要功能

- 存储文件的基本属性（名称、大小、时间等）
- 记录压缩相关信息（原始大小、压缩后大小）
- 支持符号链接信息
- 提供压缩包统计信息

### 使用示例

```go
// 获取压缩包信息
info, err := comprx.List("archive.zip")
if err != nil {
    log.Fatal(err)
}

// 访问压缩包统计信息
fmt.Printf("文件总数: %d\n", info.TotalFiles)
fmt.Printf("原始大小: %d 字节\n", info.TotalSize)

// 遍历文件列表
for _, file := range info.Files {
    fmt.Printf("文件: %s, 大小: %d\n", file.Name, file.Size)
}
```

## 核心数据类型和常量

### 主要类型

- **ProgressStyle**: 进度条样式类型
- **CompressType**: 压缩格式类型
- **CompressionLevel**: 压缩等级类型

### 主要功能

- 定义支持的压缩格式和进度条样式
- 提供压缩格式自动检测功能
- 提供类型验证和转换方法
- 定义压缩等级常量和验证

## FUNCTIONS

### HasFilterConditions

```go
func HasFilterConditions(filter *FilterOptions) bool
```

- **描述**: 检查过滤器是否有任何过滤条件
- **参数**:
  - `filter`: 过滤配置选项
- **返回**:
  - `bool`: `true` 表示有过滤条件，`false` 表示没有

### IsSupportedCompressType

```go
func IsSupportedCompressType(ct string) bool
```

- **描述**: 判断是否受支持的压缩格式
- **参数**:
  - `ct`: 压缩格式字符串
- **返回**:
  - `bool`: 如果是受支持的压缩格式, 返回 `true`, 否则返回 `false`

### SupportedCompressTypes

```go
func SupportedCompressTypes() []string
```

- **描述**: 返回受支持的压缩格式字符串列表
- **返回**:
  - `[]string`: 受支持的压缩格式字符串列表

## TYPES

### ArchiveInfo

```go
type ArchiveInfo struct {
    Type           CompressType // 压缩包类型
    TotalFiles     int          // 总文件数
    TotalSize      int64        // 总原始大小
    CompressedSize int64        // 总压缩大小
    Files          []FileInfo   // 文件列表
}
```

- **描述**: 压缩包整体信息

### CompressType

```go
type CompressType string
```

- **描述**: 支持的压缩格式
- **压缩格式类型定义**:
  - `CompressTypeZip`: zip 压缩格式
  - `CompressTypeTar`: tar 压缩格式
  - `CompressTypeTgz`: tgz 压缩格式
  - `CompressTypeTarGz`: tar.gz 压缩格式
  - `CompressTypeGz`: gz 压缩格式
  - `CompressTypeBz2`: bz2 压缩格式
  - `CompressTypeBzip2`: bzip2 压缩格式
  - `CompressTypeZlib`: zlib 压缩格式

### 常量

```go
const (
    CompressTypeZip   CompressType = ".zip"    // zip 压缩格式
    CompressTypeTar   CompressType = ".tar"    // tar 压缩格式
    CompressTypeTgz   CompressType = ".tgz"    // tgz 压缩格式
    CompressTypeTarGz CompressType = ".tar.gz" // tar.gz 压缩格式
    CompressTypeGz    CompressType = ".gz"     // gz 压缩格式
    CompressTypeBz2   CompressType = ".bz2"    // bz2 压缩格式
    CompressTypeBzip2 CompressType = ".bzip2"  // bzip2 压缩格式
    CompressTypeZlib  CompressType = ".zlib"   // zlib 压缩格式
)
```

### String

```go
func (c CompressType) String() string
```

- **描述**: 压缩格式的字符串表示
- **返回**:
  - `string`: 压缩格式的字符串表示

### CompressionLevel

```go
type CompressionLevel int
```

- **描述**: 压缩等级类型
- **支持的压缩等级**:
  - `CompressionLevelDefault`: 默认压缩等级
  - `CompressionLevelNone`: 禁用压缩
  - `CompressionLevelFast`: 快速压缩
  - `CompressionLevelBest`: 最佳压缩
  - `CompressionLevelHuffmanOnly`: 仅使用Huffman编码

### 常量

```go
const (
    CompressionLevelDefault     CompressionLevel = -1 // 默认压缩等级(zip仅支持该等级)
    CompressionLevelNone        CompressionLevel = 0  // 禁用压缩(zip仅支持该等级)
    CompressionLevelFast        CompressionLevel = 1  // 快速压缩
    CompressionLevelBest        CompressionLevel = 9  // 最佳压缩
    CompressionLevelHuffmanOnly CompressionLevel = -2 // 仅使用Huffman编码
)
```

### SupportedCompressionLevels

```go
func SupportedCompressionLevels() []CompressionLevel
```

- **描述**: 返回所有预定义的压缩等级

### IsValid

```go
func (cl CompressionLevel) IsValid() bool
```

- **描述**: 检查压缩等级是否有效
- **有效范围**: -2 到 9

### String

```go
func (cl CompressionLevel) String() string
```

- **描述**: 返回压缩等级的字符串表示

### FileInfo

```go
type FileInfo struct {
    Name           string      // 文件名/路径
    Size           int64       // 原始大小
    CompressedSize int64       // 压缩后大小
    ModTime        time.Time   // 修改时间
    Mode           os.FileMode // 文件权限
    IsDir          bool        // 是否为目录
    IsSymlink      bool        // 是否为符号链接
    LinkTarget     string      // 符号链接目标(如果是符号链接)
}
```

- **描述**: 压缩包内文件信息

### FilterOptions

```go
type FilterOptions struct {
    Include []string // 包含模式，支持 glob 语法，只处理匹配的文件
    Exclude []string // 排除模式，支持 glob 语法，跳过匹配的文件
    MaxSize int64    // 最大文件大小（字节），0 表示无限制
    MinSize int64    // 最小文件大小（字节），默认为 0
}
```

- **描述**: 过滤配置选项
- **用于指定压缩时或解压时的文件过滤条件**:
  - `Include`: 包含模式列表，支持 glob 语法，只有匹配的文件才会被处理
  - `Exclude`: 排除模式列表，支持 glob 语法，匹配的文件会被跳过
  - `MaxSize`: 最大文件大小限制（字节），0 表示无限制，超过此大小的文件会被跳过
  - `MinSize`: 最小文件大小限制（字节），默认为 0，小于此大小的文件会被跳过

### ShouldSkipByParams

```go
func (f *FilterOptions) ShouldSkipByParams(path string, size int64, isDir bool) bool
```

- **描述**: 判断文件是否应该被跳过(通用方法，用于压缩和解压)
- **过滤逻辑**:
  1. 检查文件大小是否符合要求
  2. 如果指定了包含模式，检查文件是否匹配包含模式
  3. 检查文件是否匹配排除模式
- **参数**:
  - `path`: 文件路径
  - `size`: 文件大小（字节）
  - `isDir`: 是否为目录
- **返回**:
  - `bool`: `true` 表示应该跳过，`false` 表示应该处理

### Validate

```go
func (f *FilterOptions) Validate() error
```

- **描述**: 验证过滤器选项
- **返回**:
  - `error`: 验证错误，如果验证通过则返回 `nil`

### ProgressStyle

```go
type ProgressStyle string
```

- **描述**: 进度条样式类型
- **进度条样式类型定义**:
  - `ProgressStyleText`: 文本样式进度条 - 使用文字描述进度
  - `ProgressStyleUnicode`: Unicode样式进度条 - 使用Unicode字符绘制精美进度条
  - `ProgressStyleASCII`: ASCII样式进度条 - 使用基础ASCII字符绘制兼容性最好的进度条

### 常量

```go
const (
    ProgressStyleText    ProgressStyle = "text"    // 文本样式进度条 - 使用文字描述进度
    ProgressStyleDefault ProgressStyle = "default" // 默认进度条样式 - progress库的默认进度条样式
    ProgressStyleUnicode ProgressStyle = "unicode" // Unicode样式进度条 - 使用Unicode字符绘制精美进度条
    ProgressStyleASCII   ProgressStyle = "ascii"   // ASCII样式进度条 - 使用基础ASCII字符绘制兼容性最好的进度条
)
```

### SupportedProgressStyles

```go
func SupportedProgressStyles() []ProgressStyle
```

- **描述**: 返回所有支持的进度条样式

### IsValid

```go
func (ps ProgressStyle) IsValid() bool
```

- **描述**: 检查进度条样式是否有效

### String

```go
func (ps ProgressStyle) String() string
```

- **描述**: 返回进度条样式的字符串表示

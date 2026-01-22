# Package comprx

Package comprx 提供了一个统一的压缩和解压缩库，支持多种压缩格式，包括 ZIP、TAR、GZIP、BZIP2、ZLIB 和 TGZ。它还支持进度条显示、文件过滤、并发安全操作等高级功能。

## 主要功能

- 压缩和解压缩多种格式的文件
- 支持进度条显示
- 文件过滤功能
- 线程安全操作
- 灵活的配置选项

## 基本使用示例

```go
// 简单压缩
err := comprx.Pack("output.zip", "input_dir")

// 简单解压
err := comprx.Unpack("archive.zip", "output_dir")

// 带进度条的压缩
err := comprx.PackProgress("output.zip", "input_dir")
```

## 文件过滤功能

Package comprx 提供文件过滤功能，支持从忽略文件加载排除模式。

### 主要功能

- 从忽略文件加载排除模式
- 支持注释行和空行处理
- 自动去重排除模式
- 支持 glob 模式匹配
- 提供文件不存在时的容错处理

### 支持的忽略文件格式

- 每行一个模式
- `#` 开头的注释行
- 空行自动忽略
- 支持标准 glob 通配符

### 使用示例

```go
// 加载忽略文件，文件不存在会报错
patterns, err := comprx.LoadExcludeFromFile(".gitignore")

// 加载忽略文件，文件不存在返回空列表
patterns, err := comprx.LoadExcludeFromFileOrEmpty(".comprxignore")
```

## 压缩包内容列表和信息查看功能

Package comprx 提供查看压缩包内容的各种方法，包括列出文件信息、打印压缩包信息等。支持多种压缩格式，提供简洁和详细两种显示样式，支持文件过滤和数量限制。

### 主要功能

- 列出压缩包内的文件信息
- 打印压缩包基本信息
- 支持文件名模式匹配
- 支持限制显示文件数量
- 提供简洁和详细两种显示样式

### 使用示例

```go
// 列出压缩包内所有文件
info, err := comprx.List("archive.zip")

// 打印压缩包信息（简洁样式）
err := comprx.PrintLs("archive.zip")

// 打印匹配模式的文件（详细样式）
err := comprx.PrintLlMatch("archive.zip", "*.go")
```

## 内存中的压缩和解压缩功能

Package comprx 提供 GZIP 和 ZLIB 格式的内存压缩和流式压缩功能。支持字节数组、字符串和流式数据的压缩与解压缩操作。

### 主要功能

- GZIP 内存压缩：字节数组和字符串的压缩解压
- GZIP 流式压缩：支持 `io.Reader` 和 `io.Writer` 接口
- ZLIB 内存压缩：字节数组和字符串的压缩解压
- ZLIB 流式压缩：支持 `io.Reader` 和 `io.Writer` 接口
- 支持自定义压缩等级

### 使用示例

```go
// GZIP 压缩字符串
compressed, err := comprx.GzipString("hello world")

// ZLIB 解压字节数据
decompressed, err := comprx.UnzlibBytes(compressedData)
```

## 压缩和解压缩操作的配置选项

Package comprx 定义了 `Options` 结构体和相关的配置方法，用于控制压缩和解压缩操作的行为。支持压缩等级设置、进度条显示、文件过滤、路径验证等功能的配置。

### 主要类型

- `Options`: 压缩/解压配置选项结构体

### 主要功能

- 提供默认配置选项
- 支持链式配置方法
- 提供各种预设配置选项

### 使用示例

```go
opts := comprx.Options{
    CompressionLevel: config.CompressionLevelBest,
    OverwriteExisting: true,
    ProgressEnabled: true,
    ProgressStyle: ProgressStyleUnicode,
}
err := comprx.PackOptions("output.zip", "input_dir", opts)
```

## 文件和目录大小计算功能

Package comprx 提供计算文件或目录大小的实用函数。支持单个文件大小获取和目录递归大小计算，提供安全和详细两种版本。

### 主要功能

- 获取单个文件的大小
- 递归计算目录的总大小
- 提供安全版本（出错返回 0）和详细版本（返回错误信息）
- 自动忽略符号链接等特殊文件

### 使用示例

```go
// 安全版本，出错时返回 0
size := comprx.GetSizeOrZero("./mydir")

// 详细版本，返回错误信息
size, err := comprx.GetSize("./myfile.txt")
```

## FUNCTIONS

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

### GzipBytes

```go
func GzipBytes(data []byte) ([]byte, error)
```

- **描述**: 压缩字节数据（使用默认压缩等级）
- **参数**:
  - `data`: 要压缩的字节数据
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := GzipBytes([]byte("hello world"))
```

### GzipBytesWithLevel

```go
func GzipBytesWithLevel(data []byte, level CompressionLevel) ([]byte, error)
```

- **描述**: 压缩字节数据（指定压缩等级）
- **参数**:
  - `data`: 要压缩的字节数据
  - `level`: 压缩级别
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := GzipBytesWithLevel([]byte("hello world"), CompressionLevelBest)
```

### GzipStream

```go
func GzipStream(dst io.Writer, src io.Reader) error
```

- **描述**: 流式压缩数据（使用默认压缩等级）
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
file, _ := os.Open("input.txt")
defer file.Close()

var buf bytes.Buffer
err := GzipStream(&buf, file)
```

### GzipStreamWithLevel

```go
func GzipStreamWithLevel(dst io.Writer, src io.Reader, level CompressionLevel) error
```

- **描述**: 流式压缩数据（指定压缩等级）
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器
  - `level`: 压缩级别
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
file, _ := os.Open("input.txt")
defer file.Close()

output, _ := os.Create("output.gz")
defer output.Close()

err := GzipStreamWithLevel(output, file, CompressionLevelBest)
```

### GzipString

```go
func GzipString(text string) ([]byte, error)
```

- **描述**: 压缩字符串（使用默认压缩等级）
- **参数**:
  - `text`: 要压缩的字符串
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := GzipString("hello world")
```

### GzipStringWithLevel

```go
func GzipStringWithLevel(text string, level CompressionLevel) ([]byte, error)
```

- **描述**: 压缩字符串（指定压缩等级）
- **参数**:
  - `text`: 要压缩的字符串
  - `level`: 压缩级别
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := GzipStringWithLevel("hello world", CompressionLevelBest)
```

### List

```go
func List(archivePath string) (*ArchiveInfo, error)
```

- **描述**: 列出压缩包的所有文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `*ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListLimit

```go
func ListLimit(archivePath string, limit int) (*ArchiveInfo, error)
```

- **描述**: 列出指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制返回的文件数量
- **返回**:
  - `*ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### ListMatch

```go
func ListMatch(archivePath string, pattern string) (*ArchiveInfo, error)
```

- **描述**: 列出匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `*ArchiveInfo`: 压缩包信息
  - `error`: 错误信息

### LoadExcludeFromFile

```go
func LoadExcludeFromFile(ignoreFilePath string) ([]string, error)
```

- **描述**: 从忽略文件加载排除模式
- **参数**:
  - `ignoreFilePath`: 忽略文件路径（如 `.comprxignore`, `.gitignore`）
- **返回**:
  - `[]string`: 排除模式列表（已去重）
  - `error`: 错误信息
- **支持的文件格式**:
  - 每行一个模式
  - 支持 `#` 开头的注释行
  - 自动忽略空行
  - 支持 glob 模式匹配
  - 自动去除重复模式
- **使用示例**:

```go
patterns, err := comprx.LoadExcludeFromFile(".comprxignore")
```

### LoadExcludeFromFileOrEmpty

```go
func LoadExcludeFromFileOrEmpty(ignoreFilePath string) ([]string, error)
```

- **描述**: 从忽略文件加载排除模式，文件不存在时返回空列表
- **参数**:
  - `ignoreFilePath`: 忽略文件路径
- **返回**:
  - `[]string`: 排除模式列表，文件不存在时返回空列表
  - `error`: 错误信息（文件不存在不算错误）
- **使用示例**:

```go
patterns, err := comprx.LoadExcludeFromFileOrEmpty(".comprxignore")
```

### Pack

```go
func Pack(dst string, src string) error
```

- **描述**: 压缩文件或目录（禁用进度条） - 线程安全
- **参数**:
  - `dst`: 目标文件路径
  - `src`: 源文件路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := Pack("output.zip", "input_dir")
```

### PackOptions

```go
func PackOptions(dst string, src string, opts Options) error
```

- **描述**: 使用指定配置压缩文件或目录 - 线程安全
- **参数**:
  - `dst`: 目标文件路径
  - `src`: 源文件路径
  - `opts`: 配置选项
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
opts := Options{
    CompressionLevel: config.CompressionLevelBest,
    OverwriteExisting: true,
    ProgressEnabled: true,
    ProgressStyle: ProgressStyleUnicode,
}
err := PackOptions("output.zip", "input_dir", opts)
```

### PackProgress

```go
func PackProgress(dst string, src string) error
```

- **描述**: 压缩文件或目录（启用进度条） - 线程安全
- **参数**:
  - `dst`: 目标文件路径
  - `src`: 源文件路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := PackProgress("output.zip", "input_dir")
```

### PrintArchiveAndFiles

```go
func PrintArchiveAndFiles(archivePath string, detailed bool) error
```

- **描述**: 打印压缩包信息和所有文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintArchiveAndFilesLimit

```go
func PrintArchiveAndFilesLimit(archivePath string, limit int, detailed bool) error
```

- **描述**: 打印压缩包信息和指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintArchiveAndFilesMatch

```go
func PrintArchiveAndFilesMatch(archivePath string, pattern string, detailed bool) error
```

- **描述**: 打印压缩包信息和匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintArchiveInfo

```go
func PrintArchiveInfo(archivePath string) error
```

- **描述**: 打印压缩包本身的基本信息
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `error`: 错误信息

### PrintFiles

```go
func PrintFiles(archivePath string, detailed bool) error
```

- **描述**: 打印压缩包内所有文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintFilesLimit

```go
func PrintFilesLimit(archivePath string, limit int, detailed bool) error
```

- **描述**: 打印压缩包内指定数量的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintFilesMatch

```go
func PrintFilesMatch(archivePath string, pattern string, detailed bool) error
```

- **描述**: 打印压缩包内匹配指定模式的文件信息
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
  - `detailed`: `true`=详细样式, `false`=简洁样式(默认)
- **返回**:
  - `error`: 错误信息

### PrintInfo

```go
func PrintInfo(archivePath string) error
```

- **描述**: 打印压缩包信息和所有文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `error`: 错误信息

### PrintInfoDetailed

```go
func PrintInfoDetailed(archivePath string) error
```

- **描述**: 打印压缩包信息和所有文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `error`: 错误信息

### PrintInfoDetailedLimit

```go
func PrintInfoDetailedLimit(archivePath string, limit int) error
```

- **描述**: 打印压缩包信息和指定数量的文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
- **返回**:
  - `error`: 错误信息

### PrintInfoDetailedMatch

```go
func PrintInfoDetailedMatch(archivePath string, pattern string) error
```

- **描述**: 打印压缩包信息和匹配指定模式的文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `error`: 错误信息

### PrintInfoLimit

```go
func PrintInfoLimit(archivePath string, limit int) error
```

- **描述**: 打印压缩包信息和指定数量的文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
- **返回**:
  - `error`: 错误信息

### PrintInfoMatch

```go
func PrintInfoMatch(archivePath string, pattern string) error
```

- **描述**: 打印压缩包信息和匹配指定模式的文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `error`: 错误信息

### PrintLl

```go
func PrintLl(archivePath string) error
```

- **描述**: 打印压缩包内所有文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `error`: 错误信息

### PrintLlLimit

```go
func PrintLlLimit(archivePath string, limit int) error
```

- **描述**: 打印压缩包内指定数量的文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
- **返回**:
  - `error`: 错误信息

### PrintLlMatch

```go
func PrintLlMatch(archivePath string, pattern string) error
```

- **描述**: 打印压缩包内匹配指定模式的文件信息（详细样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `error`: 错误信息

### PrintLs

```go
func PrintLs(archivePath string) error
```

- **描述**: 打印压缩包内所有文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
- **返回**:
  - `error`: 错误信息

### PrintLsLimit

```go
func PrintLsLimit(archivePath string, limit int) error
```

- **描述**: 打印压缩包内指定数量的文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `limit`: 限制打印的文件数量
- **返回**:
  - `error`: 错误信息

### PrintLsMatch

```go
func PrintLsMatch(archivePath string, pattern string) error
```

- **描述**: 打印压缩包内匹配指定模式的文件信息（简洁样式）
- **参数**:
  - `archivePath`: 压缩包文件路径
  - `pattern`: 文件名匹配模式 (支持通配符 `*` 和 `?`)
- **返回**:
  - `error`: 错误信息

### UngzipBytes

```go
func UngzipBytes(compressedData []byte) ([]byte, error)
```

- **描述**: 解压字节数据
- **参数**:
  - `compressedData`: 压缩的字节数据
- **返回**:
  - `[]byte`: 解压后的数据
  - `error`: 错误信息
- **使用示例**:

```go
decompressed, err := UngzipBytes(compressedData)
```

### UngzipStream

```go
func UngzipStream(dst io.Writer, src io.Reader) error
```

- **描述**: 流式解压数据
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器（压缩数据）
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
compressedFile, _ := os.Open("input.gz")
defer compressedFile.Close()

output, _ := os.Create("output.txt")
defer output.Close()

err := UngzipStream(output, compressedFile)
```

### UngzipString

```go
func UngzipString(compressedData []byte) (string, error)
```

- **描述**: 解压为字符串
- **参数**:
  - `compressedData`: 压缩的字节数据
- **返回**:
  - `string`: 解压后的字符串
  - `error`: 错误信息
- **使用示例**:

```go
text, err := UngzipString(compressedData)
```

### Unpack

```go
func Unpack(src string, dst string) error
```

- **描述**: 解压文件（禁用进度条） - 线程安全
- **参数**:
  - `src`: 源文件路径
  - `dst`: 目标目录路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := Unpack("archive.zip", "output_dir")
```

### UnpackDir

```go
func UnpackDir(archivePath string, dirName string, outputDir string) error
```

- **描述**: 解压指定目录 - 线程安全
- **参数**:
  - `archivePath`: 压缩包路径
  - `dirName`: 要解压的目录名
  - `outputDir`: 输出目录路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := UnpackDir("archive.zip", "src", "output/")
```

### UnpackFile

```go
func UnpackFile(archivePath string, fileName string, outputDir string) error
```

- **描述**: 解压指定文件名 - 线程安全
- **参数**:
  - `archivePath`: 压缩包路径
  - `fileName`: 要解压的文件名
  - `outputDir`: 输出目录路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := UnpackFile("archive.zip", "config.json", "output/")
```

### UnpackMatch

```go
func UnpackMatch(archivePath string, keyword string, outputDir string) error
```

- **描述**: 解压匹配关键字的文件 - 线程安全
- **参数**:
  - `archivePath`: 压缩包路径
  - `keyword`: 匹配关键字
  - `outputDir`: 输出目录路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := UnpackMatch("archive.zip", "test", "output/")
```

### UnpackOptions

```go
func UnpackOptions(src string, dst string, opts Options) error
```

- **描述**: 使用指定配置解压文件 - 线程安全
- **参数**:
  - `src`: 源文件路径
  - `dst`: 目标目录路径
  - `opts`: 配置选项
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
opts := Options{
    OverwriteExisting: true,
    ProgressEnabled: true,
    ProgressStyle: ProgressStyleASCII,
}
err := UnpackOptions("archive.zip", "output_dir", opts)
```

### UnpackProgress

```go
func UnpackProgress(src string, dst string) error
```

- **描述**: 解压文件（启用进度条） - 线程安全
- **参数**:
  - `src`: 源文件路径
  - `dst`: 目标目录路径
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
err := UnpackProgress("archive.zip", "output_dir")
```

### UnzlibBytes

```go
func UnzlibBytes(compressedData []byte) ([]byte, error)
```

- **描述**: 解压字节数据
- **参数**:
  - `compressedData`: 压缩的字节数据
- **返回**:
  - `[]byte`: 解压后的数据
  - `error`: 错误信息
- **使用示例**:

```go
decompressed, err := UnzlibBytes(compressedData)
```

### UnzlibStream

```go
func UnzlibStream(dst io.Writer, src io.Reader) error
```

- **描述**: 流式解压数据
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器（压缩数据）
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
compressedFile, _ := os.Open("input.zlib")
defer compressedFile.Close()

output, _ := os.Create("output.txt")
defer output.Close()

err := UnzlibStream(output, compressedFile)
```

### UnzlibString

```go
func UnzlibString(compressedData []byte) (string, error)
```

- **描述**: 解压为字符串
- **参数**:
  - `compressedData`: 压缩的字节数据
- **返回**:
  - `string`: 解压后的字符串
  - `error`: 错误信息
- **使用示例**:

```go
text, err := UnzlibString(compressedData)
```

### ZlibBytes

```go
func ZlibBytes(data []byte) ([]byte, error)
```

- **描述**: 压缩字节数据（使用默认压缩等级）
- **参数**:
  - `data`: 要压缩的字节数据
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := ZlibBytes([]byte("hello world"))
```

### ZlibBytesWithLevel

```go
func ZlibBytesWithLevel(data []byte, level CompressionLevel) ([]byte, error)
```

- **描述**: 压缩字节数据（指定压缩等级）
- **参数**:
  - `data`: 要压缩的字节数据
  - `level`: 压缩级别
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := ZlibBytesWithLevel([]byte("hello world"), CompressionLevelBest)
```

### ZlibStream

```go
func ZlibStream(dst io.Writer, src io.Reader) error
```

- **描述**: 流式压缩数据（使用默认压缩等级）
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
file, _ := os.Open("input.txt")
defer file.Close()

var buf bytes.Buffer
err := ZlibStream(&buf, file)
```

### ZlibStreamWithLevel

```go
func ZlibStreamWithLevel(dst io.Writer, src io.Reader, level CompressionLevel) error
```

- **描述**: 流式压缩数据（指定压缩等级）
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器
  - `level`: 压缩级别
- **返回**:
  - `error`: 错误信息
- **使用示例**:

```go
file, _ := os.Open("input.txt")
defer file.Close()

output, _ := os.Create("output.zlib")
defer output.Close()

err := ZlibStreamWithLevel(output, file, CompressionLevelBest)
```

### ZlibString

```go
func ZlibString(text string) ([]byte, error)
```

- **描述**: 压缩字符串（使用默认压缩等级）
- **参数**:
  - `text`: 要压缩的字符串
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := ZlibString("hello world")
```

### ZlibStringWithLevel

```go
func ZlibStringWithLevel(text string, level CompressionLevel) ([]byte, error)
```

- **描述**: 压缩字符串（指定压缩等级）
- **参数**:
  - `text`: 要压缩的字符串
  - `level`: 压缩级别
- **返回**:
  - `[]byte`: 压缩后的数据
  - `error`: 错误信息
- **使用示例**:

```go
compressed, err := ZlibStringWithLevel("hello world", CompressionLevelBest)
```

## TYPES

### ArchiveInfo

```go
type ArchiveInfo = types.ArchiveInfo
```

- **描述**: 压缩包整体信息
- **字段**:
  - `Type`: 压缩包类型
  - `TotalFiles`: 总文件数
  - `TotalSize`: 总原始大小
  - `CompressedSize`: 总压缩大小
  - `Files`: 文件列表
- **使用示例**:

```go
info, err := comprx.List("archive.zip")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("压缩格式: %s\n", info.Type)
fmt.Printf("文件总数: %d\n", info.TotalFiles)
```

### CompressionLevel

```go
type CompressionLevel = types.CompressionLevel
```

- **描述**: 压缩等级类型
- **支持的压缩等级**:
  - `CompressionLevelDefault`: 默认压缩等级
  - `CompressionLevelNone`: 禁用压缩
  - `CompressionLevelFast`: 快速压缩
  - `CompressionLevelBest`: 最佳压缩
  - `CompressionLevelHuffmanOnly`: 仅使用Huffman编码

### CompressType

```go
type CompressType = types.CompressType
```

- **描述**: 压缩格式类型
- **支持的压缩格式**:
  - `CompressTypeZip`: zip 压缩格式
  - `CompressTypeTar`: tar 压缩格式
  - `CompressTypeTgz`: tgz 压缩格式
  - `CompressTypeTarGz`: tar.gz 压缩格式
  - `CompressTypeGz`: gz 压缩格式
  - `CompressTypeBz2`: bz2 压缩格式
  - `CompressTypeBzip2`: bzip2 压缩格式
  - `CompressTypeZlib`: zlib 压缩格式

### FileInfo

```go
type FileInfo = types.FileInfo
```

- **描述**: 压缩包内文件信息
- **字段**:
  - `Name`: 文件名/路径
  - `Size`: 原始大小
  - `CompressedSize`: 压缩后大小
  - `ModTime`: 修改时间
  - `Mode`: 文件权限
  - `IsDir`: 是否为目录
  - `IsSymlink`: 是否为符号链接
  - `LinkTarget`: 符号链接目标(如果是符号链接)
- **使用示例**:

```go
info, err := comprx.List("archive.zip")
for _, file := range info.Files {
    fmt.Printf("文件: %s, 大小: %d\n", file.Name, file.Size)
}
```

### FilterOptions

```go
type FilterOptions = types.FilterOptions
```

- **描述**: 过滤配置选项
- **字段**:
  - `Include`: 包含模式列表，支持 glob 语法
  - `Exclude`: 排除模式列表，支持 glob 语法
  - `MaxSize`: 最大文件大小限制(字节)，0 表示无限制
  - `MinSize`: 最小文件大小限制(字节)，默认为 0
- **使用示例**:

```go
filter := &comprx.FilterOptions{
    Include: []string{"*.go", "*.md"},
    Exclude: []string{"*_test.go"},
    MaxSize: 10 * 1024 * 1024, // 10MB
}
```

### ProgressStyle

```go
type ProgressStyle = types.ProgressStyle
```

- **描述**: 进度条样式类型
- **支持的进度条样式**:
  - `ProgressStyleText`: 文本样式进度条 - 使用文字描述进度
  - `ProgressStyleDefault`: 默认进度条样式 - progress库的默认进度条样式
  - `ProgressStyleUnicode`: Unicode样式进度条 - 使用Unicode字符绘制精美进度条
  - `ProgressStyleASCII`: ASCII样式进度条 - 使用基础ASCII字符绘制兼容性最好的进度条
- **使用示例**:

```go
opts := comprx.Options{
    ProgressStyle: comprx.ProgressStyleUnicode,
    ProgressEnabled: true,
}
```

## CONSTANTS

### 压缩等级常量

```go
const (
    CompressionLevelDefault     = types.CompressionLevelDefault     // 默认压缩等级
    CompressionLevelNone        = types.CompressionLevelNone        // 禁用压缩
    CompressionLevelFast        = types.CompressionLevelFast        // 快速压缩
    CompressionLevelBest        = types.CompressionLevelBest        // 最佳压缩
    CompressionLevelHuffmanOnly = types.CompressionLevelHuffmanOnly // 仅使用Huffman编码
)
```

### 压缩格式常量

```go
const (
    CompressTypeZip    = types.CompressTypeZip    // zip 压缩格式
    CompressTypeTar    = types.CompressTypeTar    // tar 压缩格式
    CompressTypeTgz    = types.CompressTypeTgz    // tgz 压缩格式
    CompressTypeTarGz  = types.CompressTypeTarGz  // tar.gz 压缩格式
    CompressTypeGz     = types.CompressTypeGz     // gz 压缩格式
    CompressTypeBz2    = types.CompressTypeBz2    // bz2 压缩格式
    CompressTypeBzip2  = types.CompressTypeBzip2  // bzip2 压缩格式
    CompressTypeZlib   = types.CompressTypeZlib   // zlib 压缩格式
)
```

### 进度条样式常量

```go
const (
    ProgressStyleText    = types.ProgressStyleText    // 文本样式进度条
    ProgressStyleDefault = types.ProgressStyleDefault // 默认进度条样式
    ProgressStyleUnicode = types.ProgressStyleUnicode // Unicode样式进度条
    ProgressStyleASCII   = types.ProgressStyleASCII   // ASCII样式进度条
)
```

### Options

```go
type Options struct {
    CompressionLevel      CompressionLevel // 压缩等级
    OverwriteExisting     bool                   // 是否覆盖已存在的文件
    ProgressEnabled       bool                   // 是否启用进度显示
    ProgressStyle         ProgressStyle    // 进度条样式
    DisablePathValidation bool                   // 是否禁用路径验证
    Filter                FilterOptions    // 过滤选项
}
```

- **描述**: 压缩/解压配置选项

### ASCIIProgressOptions

```go
func ASCIIProgressOptions() Options
```

- **描述**: 返回ASCII样式进度条配置选项
- **返回**:
  - `Options`: ASCII样式进度条配置选项
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", ASCIIProgressOptions())
```

### DefaultOptions

```go
func DefaultOptions() Options
```

- **描述**: 返回默认配置选项
- **返回**:
  - `Options`: 默认配置选项
- **默认配置**:
  - `CompressionLevel`: 默认压缩等级
  - `OverwriteExisting`: `false` (不覆盖已存在文件)
  - `ProgressEnabled`: `false` (不显示进度)
  - `ProgressStyle`: 文本样式
  - `DisablePathValidation`: `false` (启用路径验证)

### DefaultProgressOptions

```go
func DefaultProgressOptions() Options
```

- **描述**: 返回默认样式进度条配置选项
- **返回**:
  - `Options`: 默认样式进度条配置选项
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", DefaultProgressOptions())
```

### ForceOptions

```go
func ForceOptions() Options
```

- **描述**: 返回强制模式配置选项
- **返回**:
  - `Options`: 强制模式配置选项
- **配置特点**:
  - `OverwriteExisting`: `true` (覆盖已存在文件)
  - `DisablePathValidation`: `true` (禁用路径验证)
  - `ProgressEnabled`: `false` (关闭进度条)
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", ForceOptions())
```

### NoCompressionOptions

```go
func NoCompressionOptions() Options
```

- **描述**: 返回禁用压缩且启用进度条的配置选项
- **返回**:
  - `Options`: 禁用压缩且启用进度条的配置选项
- **配置特点**:
  - `CompressionLevel`: 无压缩 (存储模式)
  - `ProgressEnabled`: `true` (启用进度条)
  - `ProgressStyle`: 文本样式
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", NoCompressionOptions())
```

### NoCompressionProgressOptions

```go
func NoCompressionProgressOptions(style ProgressStyle) Options
```

- **描述**: 返回禁用压缩且启用指定样式进度条的配置选项
- **参数**:
  - `style`: 进度条样式
- **返回**:
  - `Options`: 禁用压缩且启用指定样式进度条的配置选项
- **配置特点**:
  - `CompressionLevel`: 无压缩 (存储模式)
  - `ProgressEnabled`: `true` (启用进度条)
  - `ProgressStyle`: 指定样式
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", NoCompressionProgressOptions(ProgressStyleUnicode))
```

### ProgressOptions

```go
func ProgressOptions(style ProgressStyle) Options
```

- **描述**: 返回带进度显示的配置选项
- **参数**:
  - `style`: 进度条样式
- **返回**:
  - `Options`: 带进度显示的配置选项

### TextProgressOptions

```go
func TextProgressOptions() Options
```

- **描述**: 返回文本样式进度条配置选项
- **返回**:
  - `Options`: 文本样式进度条配置选项
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", TextProgressOptions())
```

### UnicodeProgressOptions

```go
func UnicodeProgressOptions() Options
```

- **描述**: 返回Unicode样式进度条配置选项
- **返回**:
  - `Options`: Unicode样式进度条配置选项
- **使用示例**:

```go
err := PackOptions("output.zip", "input_dir", UnicodeProgressOptions())
```
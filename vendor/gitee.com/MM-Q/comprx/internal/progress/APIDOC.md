# Package progress

Package progress 提供了压缩和解压缩操作的进度显示功能，实现了多种样式的进度显示，包括文本模式和进度条模式。它支持压缩和解压缩过程中的实时进度反馈，并提供了统一的进度管理接口。

## 主要类型

- **Progress**: 进度显示器结构体

## 主要功能

- 支持多种进度条样式（文本、ASCII、Unicode、默认）
- 提供压缩和解压缩的进度显示
- 支持文件扫描进度显示
- 提供带进度的数据复制功能
- 自动管理进度条生命周期

## 支持的进度样式

- 文本模式：显示操作文本信息
- ASCII模式：使用ASCII字符的进度条
- Unicode模式：使用Unicode字符的精美进度条
- 默认模式：使用库默认样式的进度条

## 使用示例

```go
// 创建进度显示器
progress := progress.New()

// 开始进度显示
err := progress.Start(totalSize, "archive.zip", "正在解压...")

// 带进度的数据复制
written, err := progress.CopyBuffer(dst, src, buffer)

// 关闭进度显示
err := progress.Close()
```

## 源文件大小计算和进度显示的实用工具函数

### 主要功能

- 计算源路径中所有普通文件的总大小
- 在计算过程中显示扫描进度
- 支持文件过滤器，跳过不需要的文件
- 自动区分文件和目录处理
- 只在进度条模式下执行计算

### 性能优化

- 文本模式下跳过大小计算，直接返回0
- 支持过滤器提前跳过不需要的文件和目录
- 实时更新扫描进度条
- 错误容忍，遇到错误继续遍历

### 使用示例

```go
// 计算源文件总大小并显示进度
totalSize := progress.CalculateSourceTotalSizeWithProgress(
    srcPath,
    progressObj,
    "正在分析内容...",
    filterOptions
)
```

## FUNCTIONS

### CalculateSourceTotalSizeWithProgress

```go
func CalculateSourceTotalSizeWithProgress(srcPath string, progress *Progress, scanMessage string, filter *types.FilterOptions) int64
```

- **描述**: 计算源路径中所有普通文件的总大小并显示进度
- **参数**:
  - `srcPath`: 源路径（文件或目录）
  - `progress`: 进度显示对象
  - `scanMessage`: 扫描时显示的消息，如 "正在分析内容..."
  - `filter`: 文件过滤器，用于跳过不需要的文件
- **返回**:
  - `int64`: 普通文件的总大小（字节）
- **功能**:
  - 只在进度条模式下计算总大小，文本模式返回 0
  - 显示扫描进度条并实时更新
  - 支持单个文件和目录的大小计算
  - 只计算普通文件，忽略目录、符号链接等特殊文件
  - 应用过滤器跳过不需要处理的文件

## TYPES

### Progress

```go
type Progress struct {
    Enabled  bool                // 是否启用进度显示
    BarStyle types.ProgressStyle // 进度条样式

    // Has unexported fields.
}
```

- **描述**: 控制台进度显示器

### New

```go
func New() *Progress
```

- **描述**: 创建进度显示器
- **返回**:
  - `*Progress`: 简单进度显示器

### Adding

```go
func (s *Progress) Adding(filePath string)
```

- **描述**: 显示添加文件
- **参数**:
  - `filePath`: 文件路径

### Archive

```go
func (s *Progress) Archive(archivePath string)
```

- **描述**: 显示压缩文件信息
- **参数**:
  - `archivePath`: 压缩文件路径

### Close

```go
func (s *Progress) Close() error
```

- **描述**: 关闭进度显示，清理资源
- **返回**:
  - `error`: 清理错误
- **使用示例**:
  - `err := cfg.Progress.Close()`

### CloseBar

```go
func (s *Progress) CloseBar(bar *progressbar.ProgressBar) error
```

- **描述**: 通用进度条关闭方法
- **参数**:
  - `bar`: 进度条实例
- **返回**:
  - `error`: 清理错误

### Compressing

```go
func (s *Progress) Compressing(filePath string)
```

- **描述**: 显示压缩文件信息
- **参数**:
  - `filePath`: 文件路径

### CopyBuffer

```go
func (s *Progress) CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error)
```

- **描述**: 带进度显示的数据复制
- **参数**:
  - `dst`: 目标写入器
  - `src`: 源读取器
  - `buf`: 缓冲区
- **返回**:
  - `written`: 写入的字节数
  - `err`: 错误信息
- **使用示例**:
  - `written, err := cfg.Progress.CopyBuffer(fileWriter, zipReader, buffer, "file.txt")`

### Creating

```go
func (s *Progress) Creating(dirPath string)
```

- **描述**: 显示创建目录
- **参数**:
  - `dirPath`: 目录路径

### Inflating

```go
func (s *Progress) Inflating(filePath string)
```

- **描述**: 显示解压文件
- **参数**:
  - `filePath`: 文件路径

### IsEnabled

```go
func (s *Progress) IsEnabled() bool
```

- **描述**: 检查是否启用
- **返回**:
  - `bool`: 是否启用

### Start

```go
func (s *Progress) Start(totalSize int64, archivePath, description string) error
```

- **描述**: 开始进度显示，创建进度条
- **参数**:
  - `totalSize`: 总数据大小
  - `archivePath`: 压缩包路径（用于文本模式显示）
  - `description`: 操作描述（用于进度条模式显示，如"正在解压 xxx.zip..."）
- **返回**:
  - `error`: 初始化错误

### StartScan

```go
func (s *Progress) StartScan(description string) *progressbar.ProgressBar
```

- **描述**: 开始扫描进度显示，创建进度条
- **参数**:
  - `description`: 操作描述（如"正在计算文件大小..."）
- **返回**:
  - `bar`: 进度条实例

### Storing

```go
func (s *Progress) Storing(dirPath string)
```

- **描述**: 显示存储目录
- **参数**:
  - `dirPath`: 目录路径
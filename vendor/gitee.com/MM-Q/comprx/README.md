# ComprX - Go 压缩解压缩库

[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.24.4-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

ComprX 是一个功能强大、易于使用的 Go 语言压缩解压缩库，支持多种压缩格式，提供线程安全的操作和丰富的配置选项。

## ✨ 特性

- 🗜️ **多格式支持**: ZIP、TAR、TGZ、TAR.GZ、GZIP、ZLIB、BZ2/BZIP2
- 🔒 **线程安全**: 所有操作都是线程安全的
- 📊 **进度显示**: 支持多种样式的进度条（文本、Unicode、ASCII、默认）
- 🎛️ **灵活配置**: 支持压缩级别、覆盖设置等多种配置选项
- 🔍 **智能过滤**: 支持文件包含/排除模式、大小过滤，压缩和解压都支持
- 💾 **内存操作**: 支持 GZIP 和 ZLIB 的字节数据和字符串内存压缩/解压
- 🌊 **流式处理**: 支持 GZIP 和 ZLIB 的流式压缩和解压缩
- 📝 **简单易用**: 提供简洁的 API 接口和链式配置
- 📋 **文件列表**: 支持查看压缩包内容，支持模式匹配和数量限制
- 🎯 **忽略文件**: 支持从 .gitignore 等文件加载排除模式，自动去重和优化

## 📦 安装

```bash
go get gitee.com/MM-Q/comprx
```

## 🚀 快速开始

### 基本压缩和解压

```go
package main

import (
    "fmt"
    "gitee.com/MM-Q/comprx"
)

func main() {
    // 压缩文件或目录
    err := comprx.Pack("output.zip", "input_dir")
    if err != nil {
        fmt.Printf("压缩失败: %v\n", err)
        return
    }
    
    // 解压文件
    err = comprx.Unpack("output.zip", "output_dir")
    if err != nil {
        fmt.Printf("解压失败: %v\n", err)
        return
    }
    
    fmt.Println("操作完成!")
}
```

### 带进度条的压缩解压

```go
// 压缩时显示进度条
err := comprx.PackProgress("output.tar.gz", "large_directory")

// 解压时显示进度条
err := comprx.UnpackProgress("archive.zip", "output_dir")
```

### 自定义配置

```go
import (
    "gitee.com/MM-Q/comprx"
    "gitee.com/MM-Q/comprx/types"
)

// 创建自定义配置
opts := comprx.Options{
    CompressionLevel:  CompressionLevelBest,  // 最佳压缩
    OverwriteExisting: true,                        // 覆盖已存在文件
    ProgressEnabled:   true,                        // 启用进度条
    ProgressStyle:     ProgressStyleUnicode,  // Unicode 样式进度条
}

// 使用自定义配置压缩
err := comprx.PackOptions("output.zip", "input_dir", opts)

// 使用自定义配置解压
err := comprx.UnpackOptions("archive.zip", "output_dir", opts)
```

### 文件过滤功能

ComprX 支持强大的文件过滤功能，可以在压缩和解压时选择性处理文件：

```go
// 基本过滤配置
opts := comprx.DefaultOptions().
    WithInclude([]string{"*.go", "*.md"}).           // 只包含 Go 文件和 Markdown 文件
    WithExclude([]string{"*_test.go", "vendor/*"}).  // 排除测试文件和 vendor 目录
    WithSizeFilter(1024, 10*1024*1024).             // 只处理 1KB-10MB 的文件
    WithProgress(true)

err := comprx.PackOptions("filtered.zip", "project_dir", opts)

// 链式配置示例
opts := comprx.DefaultOptions().
    WithInclude([]string{"*.jpg", "*.png", "*.gif"}). // 只包含图片文件
    WithMaxSize(5 * 1024 * 1024).                     // 最大 5MB
    WithProgressAndStyle(true, ProgressStyleUnicode)

err := comprx.PackOptions("images.zip", "photos", opts)

// 使用 Set 方法配置
opts := comprx.DefaultOptions()
opts.SetInclude([]string{"src/*", "docs/*"})        // 只包含 src 和 docs 目录
opts.SetExclude([]string{"*.tmp", "*.log"})         // 排除临时文件和日志
opts.SetMinSize(100)                                // 最小 100 字节
opts.SetProgressAndStyle(true, ProgressStyleASCII)

err := comprx.PackOptions("source.tar.gz", "project", opts)
```

### 从忽略文件加载排除模式

```go
// 从 .gitignore 文件加载排除模式
excludePatterns := LoadExcludeFromFileOrEmpty(".gitignore")

opts := comprx.DefaultOptions().
    WithExclude(excludePatterns).
    WithProgress(true)

err := comprx.PackOptions("clean.zip", "project", opts)

// 组合多个忽略文件
gitignore := LoadExcludeFromFileOrEmpty(".gitignore")
dockerignore := LoadExcludeFromFileOrEmpty(".dockerignore")

allExcludes := append(gitignore, dockerignore...)
allExcludes = append(allExcludes, "*.tmp", "build/*") // 添加额外排除模式

opts := comprx.DefaultOptions().WithExclude(allExcludes)
```

## 🔍 过滤器功能详解

### 过滤器工作原理

过滤器采用三层过滤机制，按以下优先级顺序执行：

1. **文件大小过滤**：首先检查文件大小是否在允许范围内
2. **包含模式检查**：如果设置了包含模式，文件必须匹配至少一个包含模式
3. **排除模式检查**：如果文件匹配任何排除模式，将被跳过

### 模式匹配语法

支持标准的 glob 模式匹配：

```go
// 文件名匹配
"*.go"          // 匹配所有 .go 文件
"test_*.txt"    // 匹配以 test_ 开头的 .txt 文件

// 路径匹配
"src/*.go"      // 匹配 src 目录下的 .go 文件
"docs/**"       // 匹配 docs 目录及其子目录的所有文件

// 目录匹配
"vendor/"       // 匹配 vendor 目录
"node_modules/" // 匹配 node_modules 目录

// 复杂模式
"**/test_*.go"  // 匹配任意深度目录下以 test_ 开头的 .go 文件
```

### 实际应用场景

```go
// 场景1：只打包源代码，排除构建产物
opts := comprx.DefaultOptions().
    WithInclude([]string{"*.go", "*.md", "*.yml", "*.yaml"}).
    WithExclude([]string{"*.exe", "*.so", "*.dll", "build/*", "dist/*"}).
    WithProgress(true)

// 场景2：备份项目，排除依赖和缓存
gitignore := LoadExcludeFromFileOrEmpty(".gitignore")
opts := comprx.DefaultOptions().
    WithExclude(append(gitignore, "node_modules/*", ".git/*", "*.log")).
    WithMaxSize(100 * 1024 * 1024) // 排除超过100MB的文件

// 场景3：只打包小文件，用于快速传输
opts := comprx.DefaultOptions().
    WithInclude([]string{"*.txt", "*.json", "*.xml"}).
    WithSizeFilter(0, 1024*1024) // 只包含1MB以下的文件

// 场景4：媒体文件归档
opts := comprx.DefaultOptions().
    WithInclude([]string{"*.jpg", "*.png", "*.mp4", "*.mp3"}).
    WithMinSize(1024) // 排除小于1KB的文件（可能是缩略图）
```

## 🧠 内存压缩 API

### 字节数据压缩

```go
// 压缩字节数据
data := []byte("Hello, World!")
compressed, err := comprx.GzipBytes(data, CompressionLevelDefault)

// 解压字节数据
decompressed, err := comprx.UngzipBytes(compressed)
```

### 字符串压缩

```go
// 压缩字符串
text := "这是一个测试字符串"
compressed, err := comprx.GzipString(text, CompressionLevelBest)

// 解压为字符串
decompressed, err := comprx.UngzipString(compressed)
```

## 🌊 流式压缩 API

```go
import (
    "os"
    "bytes"
)

// 流式压缩（默认压缩级别）
file, _ := os.Open("input.txt")
defer file.Close()

var buf bytes.Buffer
err := comprx.GzipStream(&buf, file)

// 流式压缩（指定压缩级别）
output, _ := os.Create("output.gz")
defer output.Close()

err := comprx.GzipStreamWithLevel(output, file, CompressionLevelBest)

// 流式解压
compressedFile, _ := os.Open("input.gz")
defer compressedFile.Close()

outputFile, _ := os.Create("output.txt")
defer outputFile.Close()

err := comprx.UngzipStream(outputFile, compressedFile)
```

## 📋 支持的格式

| 格式 | 扩展名 | 压缩 | 解压 | 说明 |
|------|--------|------|------|------|
| ZIP | `.zip` | ✅ | ✅ | 最常用的压缩格式 |
| TAR | `.tar` | ✅ | ✅ | Unix 标准归档格式 |
| TGZ | `.tgz` | ✅ | ✅ | TAR + GZIP 压缩 |
| TAR.GZ | `.tar.gz` | ✅ | ✅ | TAR + GZIP 压缩 |
| GZIP | `.gz` | ✅ | ✅ | 单文件 GZIP 压缩 |
| BZIP2 | `.bz2`, `.bzip2` | ❌ | ✅ | 仅支持解压 |

## ⚙️ 配置选项

### 压缩级别

```go
CompressionLevelDefault     // 默认压缩级别
CompressionLevelNone        // 禁用压缩
CompressionLevelFast        // 快速压缩
CompressionLevelBest        // 最佳压缩
CompressionLevelHuffmanOnly // 仅使用 Huffman 编码
```

### 进度条样式

```go
ProgressStyleText     // 文本样式
ProgressStyleDefault  // 默认样式
ProgressStyleUnicode  // Unicode 样式: ████████████░░░░░░░░ 60%
ProgressStyleASCII    // ASCII 样式: [##########          ] 50%
```

### 过滤器选项

```go
// FilterOptions 结构体
type FilterOptions struct {
    Include []string // 包含模式，支持 glob 语法，只处理匹配的文件
    Exclude []string // 排除模式，支持 glob 语法，跳过匹配的文件
    MaxSize int64    // 最大文件大小（字节），0 表示无限制
    MinSize int64    // 最小文件大小（字节），默认为 0
}

// 过滤器配置方法
opts.SetInclude([]string{"*.go", "*.md"})           // 设置包含模式
opts.SetExclude([]string{"*_test.go", "vendor/*"})  // 设置排除模式
opts.SetSizeFilter(1024, 10*1024*1024)             // 设置大小范围 1KB-10MB
opts.SetMaxSize(5 * 1024 * 1024)                   // 设置最大文件大小 5MB
opts.SetMinSize(100)                                // 设置最小文件大小 100字节

// 链式配置方法
opts := comprx.DefaultOptions().
    WithInclude([]string{"*.jpg", "*.png"}).
    WithExclude([]string{"*.tmp"}).
    WithMaxSize(10 * 1024 * 1024)
```

### 预定义配置选项

```go
// 基础配置
comprx.DefaultOptions()           // 默认配置
comprx.ForceOptions()            // 强制模式（覆盖文件，禁用路径验证）
comprx.NoCompressionOptions()    // 无压缩模式

// 进度条配置
comprx.TextProgressOptions()     // 文本样式进度条
comprx.UnicodeProgressOptions()  // Unicode 样式进度条
comprx.ASCIIProgressOptions()    // ASCII 样式进度条
comprx.DefaultProgressOptions()  // 默认样式进度条

// 自定义进度条样式
comprx.ProgressOptions(ProgressStyleUnicode)
comprx.NoCompressionProgressOptions(ProgressStyleASCII)
```

## 🏗️ 项目结构

```
comprx/
├── comprx.go              # 主要 API 接口
├── options.go             # 配置选项和链式配置方法
├── filter.go              # 过滤器相关 API
├── list.go                # 文件列表 API
├── size.go                # 大小计算 API
├── types/                 # 类型定义
│   ├── types.go          # 基础类型定义（压缩格式、压缩级别、进度条样式）
│   ├── filter.go         # 过滤器类型和实现
│   └── list.go           # 列表相关类型
├── internal/
│   ├── core/             # 核心压缩逻辑和集成
│   ├── cxzip/            # ZIP 格式处理（压缩、解压、列表）
│   ├── cxtar/            # TAR 格式处理（压缩、解压、列表）
│   ├── cxtgz/            # TGZ 格式处理（压缩、解压、列表）
│   ├── cxgzip/           # GZIP 格式处理（压缩、解压、内存操作、流式处理）
│   ├── cxbzip2/          # BZIP2 格式处理（仅解压和列表）
│   ├── progress/         # 进度条实现和大小计算
│   └── utils/            # 工具函数（路径验证、缓冲区管理等）
└── README.md
```

## 🔧 高级功能

### 文件列表查看

```go
// 查看压缩包内容
files, err := comprx.List("archive.zip")
if err != nil {
    log.Fatal(err)
}

for _, file := range files {
    fmt.Printf("文件: %s, 大小: %d 字节, 修改时间: %s\n", 
        file.Name, file.Size, file.ModTime.Format("2006-01-02 15:04:05"))
}
```

## 🧪 测试

运行所有测试：

```bash
go test ./...
```

运行特定模块测试：

```bash
# 测试核心功能
go test ./internal/core/

# 测试各种压缩格式
go test ./internal/cxzip/
go test ./internal/cxtar/
go test ./internal/cxtgz/
go test ./internal/cxgzip/
go test ./internal/cxzlib/
go test ./internal/cxbzip2/

# 测试过滤器功能
go test ./types/ -v

# 测试忽略文件加载功能
go test -v -run TestLoadExcludeFromFile

# 测试并发安全性
go test -run TestConcurrent
```

运行性能测试：

```bash
# 过滤器性能测试
go test ./types/ -bench=BenchmarkFilter

# 忽略文件加载性能测试
go test -bench=BenchmarkLoadExcludeFromFile

# 压缩性能测试
go test ./internal/cxzip/ -bench=.
go test ./internal/cxgzip/ -bench=.
go test ./internal/cxzlib/ -bench=.
```

## 📊 性能特点

- **内存效率**: 流式处理大文件，内存占用稳定
- **并发安全**: 支持多协程同时操作不同的压缩任务
- **智能过滤**: 三层过滤机制，性能优化的文件筛选
- **进度可视**: 实时进度显示，支持大文件操作监控

## 🔄 版本兼容性

- **Go 版本**: 要求 Go 1.24.4 或更高版本
- **平台支持**: Windows、Linux、macOS
- **架构支持**: amd64、arm64

## 📄 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📞 联系

- 项目地址: [https://gitee.com/MM-Q/comprx](https://gitee.com/MM-Q/comprx)

---

**ComprX** - 让压缩解压变得简单高效！ 🚀

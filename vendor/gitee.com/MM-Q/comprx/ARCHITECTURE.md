# Comprx 压缩库架构图

## 整体架构

```mermaid
graph TB
    %% 用户接口层
    subgraph "用户接口层 (Public API)"
        API[comprx.go<br/>统一API接口]
        API --> Pack[Pack/PackProgress<br/>压缩函数]
        API --> Unpack[Unpack/UnpackProgress<br/>解压函数]
        API --> Options[PackOptions/UnpackOptions<br/>配置化函数]
        API --> Specific[UnpackFile/UnpackDir/UnpackMatch<br/>特定解压函数]
    end

    %% 核心层
    subgraph "核心层 (Core Layer)"
        Core[internal/core/comprx.go<br/>核心压缩器]
        Core --> CorePack[Pack方法<br/>格式检测与分发]
        Core --> CoreUnpack[Unpack方法<br/>格式检测与分发]
    end

    %% 配置层
    subgraph "配置层 (Configuration)"
        Config[internal/config/<br/>配置管理]
        Config --> CompLevel[压缩等级配置]
        Config --> Progress[进度条配置]
        Config --> Filter[过滤器配置]
    end

    %% 格式处理层
    subgraph "格式处理层 (Format Handlers)"
        ZIP[internal/cxzip/<br/>ZIP处理]
        TAR[internal/cxtar/<br/>TAR处理]
        TGZ[internal/cxtgz/<br/>TGZ处理]
        GZIP[internal/cxgzip/<br/>GZIP处理]
        BZIP2[internal/cxbzip2/<br/>BZIP2处理]
        ZLIB[internal/cxzlib/<br/>ZLIB处理]
    end

    %% 工具层
    subgraph "工具层 (Utilities)"
        Utils[internal/utils/<br/>通用工具]
        Utils --> Buffer[缓冲区管理]
        Utils --> Size[大小计算]
        Utils --> Validate[路径验证]
        Utils --> FileOps[文件操作]
    end

    %% 进度条层
    subgraph "进度条层 (Progress)"
        ProgressPkg[internal/progress/<br/>进度条管理]
        ProgressPkg --> ProgressBar[进度条显示]
        ProgressPkg --> SizeCalc[大小计算]
    end

    %% 类型定义层
    subgraph "类型定义层 (Types)"
        Types[types/<br/>类型定义]
        Types --> FilterTypes[过滤器类型]
        Types --> CompressTypes[压缩格式类型]
        Types --> ProgressTypes[进度条类型]
        Types --> ListTypes[列表类型]
    end

    %% 连接关系
    API --> Core
    Core --> Config
    Core --> ZIP
    Core --> TAR
    Core --> TGZ
    Core --> GZIP
    Core --> BZIP2
    Core --> ZLIB
    
    ZIP --> Utils
    TAR --> Utils
    TGZ --> Utils
    GZIP --> Utils
    BZIP2 --> Utils
    ZLIB --> Utils
    
    ZIP --> ProgressPkg
    TAR --> ProgressPkg
    TGZ --> ProgressPkg
    GZIP --> ProgressPkg
    BZIP2 --> ProgressPkg
    ZLIB --> ProgressPkg
    
    Config --> Types
    Utils --> Types
    ProgressPkg --> Types
    
    %% 外部依赖
    ProgressPkg --> ExtProgress[github.com/schollz/progressbar/v3<br/>外部进度条库]

    %% 样式定义
    classDef apiLayer fill:#e1f5fe
    classDef coreLayer fill:#f3e5f5
    classDef configLayer fill:#e8f5e8
    classDef formatLayer fill:#fff3e0
    classDef utilLayer fill:#fce4ec
    classDef progressLayer fill:#f1f8e9
    classDef typeLayer fill:#e0f2f1
    classDef external fill:#ffebee

    class API,Pack,Unpack,Options,Specific apiLayer
    class Core,CorePack,CoreUnpack coreLayer
    class Config,CompLevel,Progress,Filter configLayer
    class ZIP,TAR,TGZ,GZIP,BZIP2,ZLIB formatLayer
    class Utils,Buffer,Size,Validate,FileOps utilLayer
    class ProgressPkg,ProgressBar,SizeCalc progressLayer
    class Types,FilterTypes,CompressTypes,ProgressTypes,ListTypes typeLayer
    class ExtProgress external
```

## 数据流图

```mermaid
flowchart TD
    %% 压缩流程
    subgraph "压缩流程 (Pack Flow)"
        A1[用户调用Pack函数] --> A2[格式检测]
        A2 --> A3[参数验证]
        A3 --> A4[配置应用]
        A4 --> A5[选择对应格式处理器]
        A5 --> A6[执行压缩操作]
        A6 --> A7[进度条更新]
        A7 --> A8[完成压缩]
    end

    %% 解压流程
    subgraph "解压流程 (Unpack Flow)"
        B1[用户调用Unpack函数] --> B2[格式检测]
        B2 --> B3[文件存在性检查]
        B3 --> B4[配置应用]
        B4 --> B5[选择对应格式处理器]
        B5 --> B6[应用过滤器]
        B6 --> B7[执行解压操作]
        B7 --> B8[进度条更新]
        B8 --> B9[完成解压]
    end

    %% 配置流程
    subgraph "配置流程 (Configuration Flow)"
        C1[默认配置] --> C2[用户自定义配置]
        C2 --> C3[配置验证]
        C3 --> C4[配置合并]
        C4 --> C5[应用到处理器]
    end

    %% 样式
    classDef packFlow fill:#e3f2fd
    classDef unpackFlow fill:#f3e5f5
    classDef configFlow fill:#e8f5e8

    class A1,A2,A3,A4,A5,A6,A7,A8 packFlow
    class B1,B2,B3,B4,B5,B6,B7,B8,B9 unpackFlow
    class C1,C2,C3,C4,C5 configFlow
```

## 模块依赖关系

```mermaid
graph LR
    %% 主要模块
    Main[comprx.go] --> Core[internal/core]
    Core --> Config[internal/config]
    Core --> Formats[格式处理模块]
    
    %% 格式处理模块
    subgraph Formats [格式处理模块]
        ZIP[cxzip]
        TAR[cxtar]
        TGZ[cxtgz]
        GZIP[cxgzip]
        BZIP2[cxbzip2]
        ZLIB[cxzlib]
    end
    
    %% 所有格式处理模块都依赖工具层
    ZIP --> Utils[internal/utils]
    TAR --> Utils
    TGZ --> Utils
    GZIP --> Utils
    BZIP2 --> Utils
    ZLIB --> Utils
    
    %% 所有格式处理模块都依赖进度条
    ZIP --> Progress[internal/progress]
    TAR --> Progress
    TGZ --> Progress
    GZIP --> Progress
    BZIP2 --> Progress
    ZLIB --> Progress
    
    %% 类型定义被多个模块使用
    Config --> Types[types]
    Utils --> Types
    Progress --> Types
    Main --> Types
    
    %% 外部依赖
    Progress --> External[外部库<br/>progressbar/v3]
    
    %% 样式
    classDef mainModule fill:#bbdefb
    classDef coreModule fill:#c8e6c9
    classDef formatModule fill:#ffe0b2
    classDef utilModule fill:#f8bbd9
    classDef typeModule fill:#d1c4e9
    classDef external fill:#ffcdd2

    class Main mainModule
    class Core,Config coreModule
    class ZIP,TAR,TGZ,GZIP,BZIP2,ZLIB formatModule
    class Utils,Progress utilModule
    class Types typeModule
    class External external
```

## 关键特性

### 支持的压缩格式
- **ZIP**: .zip 文件的压缩和解压
- **TAR**: .tar 文件的压缩和解压
- **TGZ**: .tgz, .tar.gz 文件的压缩和解压
- **GZIP**: .gz 文件的压缩和解压
- **BZIP2**: .bz2, .bzip2 文件的解压（仅支持解压）
- **ZLIB**: .zlib 文件的压缩和解压

### 核心功能
- **自动格式检测**: 根据文件扩展名自动选择合适的处理器
- **进度条支持**: 可选的进度条显示，支持多种样式
- **文件过滤**: 支持包含/排除模式、文件大小过滤
- **线程安全**: 所有操作都是线程安全的
- **配置化**: 灵活的配置选项，支持压缩等级、覆盖设置等
- **错误处理**: 完善的错误处理和参数验证

### 设计模式
- **策略模式**: 不同压缩格式使用不同的处理策略
- **工厂模式**: 根据文件格式创建对应的处理器
- **配置模式**: 统一的配置管理
- **适配器模式**: 统一的API接口适配不同的压缩库
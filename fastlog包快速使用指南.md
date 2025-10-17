# fastlog 包快速使用指南

本指南将通过“创建配置→初始化日志器→记录不同类型日志”的完整流程，搭配可直接复用的代码示例，帮助你快速上手 fastlog 日志库。


## 一、前期准备：导入包
首先在代码中导入 fastlog 包，若需自定义配置，可额外导入 `time` 包（用于设置时间相关参数）。
```go
import (
    "time"

    "gitee.com/MM-Q/fastlog"
)
```


## 二、第一步：创建日志配置
fastlog 提供 **3种预设配置**（控制台/开发/生产环境）和 **自定义配置** 方式，可根据使用场景选择。

### 1. 场景1：仅控制台输出（快速调试）
适用于本地开发调试，仅输出到控制台，不写日志文件。
```go
// 创建控制台配置：默认日志级别为 DEBUG，禁用文件输出
cfg := fastlog.ConsoleConfig()
```

### 2. 场景2：开发环境配置（本地测试）
适用于本地开发，保留少量日志文件，启用详细日志格式。
```go
// 参数1：日志目录名（自定义，不存在会自动创建）
// 参数2：日志文件名（自定义）
cfg := fastlog.DevConfig("./dev_logs", "my_app_dev.log")
// 特性：日志级别 DEBUG、保留5个文件、保留7天、启用详细格式
```

### 3. 场景3：生产环境配置（线上服务）
适用于线上服务，启用日志压缩、限制文件保留时间和数量，禁用控制台输出。
```go
// 参数1：日志目录名（建议用绝对路径，如 "/var/log/my_app"）
// 参数2：日志文件名
cfg := fastlog.ProdConfig("/var/log/my_app", "my_app_prod.log")
// 特性：日志压缩、保留30天、保留24个文件、禁用控制台输出
```

### 4. 场景4：自定义配置（灵活调整）
若预设配置不满足需求，可通过 `NewFastLogConfig` 创建基础配置后，手动修改参数。
```go
// 1. 创建基础配置：默认目录 "applogs"、默认文件名 "app.log"
cfg := fastlog.NewFastLogConfig("", "") 

// 2. 手动修改配置（根据需求选择修改）
cfg.LogLevel = fastlog.INFO          // 日志级别设为 INFO（只记录 INFO 及以上）
cfg.OutputToConsole = true           // 启用控制台输出
cfg.OutputToFile = true              // 启用文件输出
cfg.MaxSize = 50                     // 单个日志文件最大 50MB（默认10MB）
cfg.Compress = true                  // 启用日志压缩
cfg.FlushInterval = 2 * time.Second  // 刷新间隔设为 2秒（默认1秒）
cfg.CallerInfo = true                // 记录调用者信息（文件名:函数名:行号）
```


## 三、第二步：初始化日志器
通过第一步创建的配置，调用 `NewFLog` 初始化日志器实例，后续所有日志操作都基于该实例。
```go
// 初始化日志器：传入配置实例
logger := fastlog.NewFLog(cfg)

// 可选：程序退出前关闭日志器（确保缓存日志写入文件）
defer logger.Close()
```


## 四、第三步：记录不同类型日志
fastlog 支持 **3种日志风格**（无占位符/带占位符/键值对），覆盖不同场景需求，且包含 5 个日志级别（DEBUG < INFO < WARN < ERROR < FATAL）。

### 1. 风格1：无占位符（简单文本）
适用于无需格式化的简单日志，直接传入任意类型参数（会自动转为字符串）。
```go
logger.Debug("这是 DEBUG 级别的简单日志")          // 调试日志（开发调试用）
logger.Info("这是 INFO 级别的简单日志", 123, "abc") // 信息日志（正常运行状态）
logger.Warn("这是 WARN 级别的简单日志", time.Now()) // 警告日志（潜在问题）
logger.Error("这是 ERROR 级别的简单日志", err)      // 错误日志（操作失败）
// logger.Fatal("这是 FATAL 级别的简单日志")        // 致命日志（会终止程序，谨慎使用）
```

### 2. 风格2：带占位符（格式化文本）
适用于需要动态拼接内容的日志，语法与标准库 `fmt` 一致（支持 `%s` `%d` `%v` 等占位符）。
```go
userID := 1001
userName := "zhangsan"
logger.Debugf("用户登录 [ID: %d, 姓名: %s]", userID, userName) // DEBUG 级格式化日志
logger.Infof("数据同步完成，共处理 %d 条记录", 500)             // INFO 级格式化日志
logger.Warnf("磁盘空间不足，剩余空间: %.2f GB", 10.5)          // WARN 级格式化日志
logger.Errorf("请求失败 [URL: %s, 状态码: %d]", "/api/user", 404) // ERROR 级格式化日志
```

### 3. 风格3：键值对（结构化日志）
适用于需要结构化分析的日志（如 JSON 格式输出时，键值对会转为 JSON 字段），通过 `fastlog.XXX` 函数创建字段。
```go
// 示例1：记录 INFO 级键值对日志
logger.InfoFields(
    "用户下单成功", // 日志主消息
    fastlog.Int("order_id", 123456),       // 整数字段
    fastlog.String("user_name", "zhangsan"),// 字符串字段
    fastlog.Float64("amount", 99.9),        // 浮点数字段
    fastlog.Time("create_time", time.Now()),// 时间字段
    fastlog.Bool("pay_success", true)       // 布尔字段
)

// 示例2：记录 ERROR 级键值对日志（附带错误信息）
err := fmt.Errorf("数据库连接超时")
logger.ErrorFields(
    "数据库操作失败",
    fastlog.String("action", "query_user"),
    fastlog.Error("error", err), // 错误字段
    fastlog.Int("retry_count", 3)
)
```


## 五、完整示例代码（可直接运行）
以下是一个包含“自定义配置→初始化→记录日志”的完整示例，复制到项目中即可测试。
```go
package main

import (
    "fmt"
    "time"

    "gitee.com/MM-Q/fastlog"
)

func main() {
    // 1. 创建自定义配置
    cfg := fastlog.NewFastLogConfig("./my_logs", "my_app.log")
    cfg.LogLevel = fastlog.INFO          // 日志级别 INFO
    cfg.OutputToConsole = true           // 启用控制台输出
    cfg.OutputToFile = true              // 启用文件输出
    cfg.MaxSize = 20                     // 单个文件最大 20MB
    cfg.CallerInfo = true                // 记录调用者信息

    // 2. 初始化日志器（defer 关闭）
    logger := fastlog.NewFLog(cfg)
    defer func() {
        if err := logger.Close(); err != nil {
            fmt.Printf("关闭日志器失败: %v\n", err)
        }
    }()

    // 3. 记录不同类型日志
    // 3.1 无占位符
    logger.Info("程序启动成功，日志目录:", cfg.LogDirName)

    // 3.2 带占位符
    userID := 1001
    logger.Infof("用户 [ID: %d] 访问系统", userID)

    // 3.3 键值对（模拟业务日志）
    err := fmt.Errorf("请求超时")
    logger.ErrorFields(
        "API 请求失败",
        fastlog.String("url", "https://api.example.com/user"),
        fastlog.Int("status_code", 504),
        fastlog.Error("error", err),
        fastlog.Time("request_time", time.Now()),
    )

    fmt.Println("日志记录完成，可查看 ./my_logs 目录下的日志文件")
}
```
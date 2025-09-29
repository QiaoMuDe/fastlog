/*
config.go - 日志配置管理模块
定义日志配置结构体及配置项的设置与获取方法，负责管理FastLog的所有可配置参数。
*/
package fastlog

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// FastLogConfig 定义一个配置结构体, 用于配置日志记录器
type FastLogConfig struct {
	LogDirName          string        // 日志目录路径
	LogFileName         string        // 日志文件名
	OutputToConsole     bool          // 是否将日志输出到控制台
	OutputToFile        bool          // 是否将日志输出到文件
	FlushInterval       time.Duration // 刷新间隔, 单位为time.Duration
	LogLevel            LogLevel      // 日志级别
	ChanIntSize         int           // 通道大小 默认10000
	LogFormat           LogFormatType // 日志格式选项
	Color               bool          // 是否启用终端颜色
	Bold                bool          // 是否启用终端字体加粗
	MaxSize             int           // 最大日志文件大小, 单位为MB, 默认10MB
	MaxAge              int           // 最大日志文件保留天数, 默认为0, 表示不做限制
	MaxFiles            int           // 最大日志文件保留数量, 默认为0, 表示不做限制
	LocalTime           bool          // 是否使用本地时间 默认使用UTC时间
	Compress            bool          // 是否启用日志文件压缩 默认不启用
	BatchSize           int           // 批处理数量
	DisableBackpressure bool          // 是否禁用背压控制, 默认false(即默认启用背压)
}

// NewFastLogConfig 创建一个新的FastLogConfig实例, 用于配置日志记录器。
//
// 参数:
//   - logDirName: 日志目录名称, 默认为"applogs"。
//   - logFileName: 日志文件名称, 默认为"app.log"。
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
func NewFastLogConfig(logDirName string, logFileName string) *FastLogConfig {
	// 返回一个新的FastLogConfig实例
	return &FastLogConfig{
		LogDirName:          logDirName,          // 日志目录名称
		LogFileName:         logFileName,         // 日志文件名称
		OutputToConsole:     true,                // 是否将日志输出到控制台
		OutputToFile:        true,                // 是否将日志输出到文件
		LogLevel:            INFO,                // 日志级别 默认INFO
		ChanIntSize:         defaultChanSize,     // 通道大小 默认10000
		FlushInterval:       normalFlushInterval, // 刷新间隔 默认500毫秒
		LogFormat:           Simple,              // 日志格式选项
		MaxSize:             defaultMaxFileSize,  // 最大日志文件大小, 单位为MB, 默认10MB
		MaxAge:              0,                   // 最大日志文件保留天数, 默认为0, 表示不做限制
		MaxFiles:            0,                   // 最大日志文件保留数量, 默认为0, 表示不做限制
		LocalTime:           true,                // 是否使用本地时间 默认使用本地时间
		Compress:            false,               // 是否启用日志文件压缩 默认不启用
		Color:               true,                // 是否启用终端颜色
		Bold:                true,                // 是否启用终端字体加粗
		BatchSize:           defaultBatchSize,    // 批处理数量
		DisableBackpressure: false,               // 是否禁用背压控制, 默认false(即默认启用背压)
	}
}

// =============================================================
// 预设配置模式函数
// =============================================================

// DevConfig 开发模式配置
// 特点：
//   - 双输出：控制台+文件同时输出，便于实时查看和持久化存储
//   - 详细信息：DEBUG级别日志，Detailed格式，包含完整调用信息
//   - 快速响应：100ms快速刷新，立即看到日志输出
//   - 彩色显示：控制台彩色+加粗，提升开发体验
//   - 短期保留：7天保留期，10个备份文件，节省开发环境存储
//
// 参数:
//   - logDirName: 日志目录名称
//   - logFileName: 日志文件名称
//
// 返回值:
//   - *FastLogConfig: 开发模式配置实例
func DevConfig(logDirName, logFileName string) *FastLogConfig {
	config := NewFastLogConfig(logDirName, logFileName)

	// 开发模式差异化设置
	config.LogLevel = DEBUG                  // 开发模式默认DEBUG级别
	config.FlushInterval = fastFlushInterval // 开发模式快速刷新
	config.LogFormat = Detailed              // 开发模式详细信息
	config.MaxAge = developmentMaxAge        // 开发模式默认7天保留
	config.MaxFiles = developmentMaxBackups  // 开发模式默认10个备份文件
	config.DisableBackpressure = true        // 开发模式禁用背压，便于调试

	return config
}

// ProdConfig 生产模式配置
// 特点：
//   - 高性能：仅文件输出，大缓冲区(20000)，慢刷新(1000ms)，大批处理(2000)
//   - 结构化：JSON格式便于日志分析系统解析和索引
//   - 合理级别：INFO级别过滤调试信息，减少日志量
//   - 长期存储：30天保留期，100个备份文件，满足审计要求
//   - 空间优化：启用压缩减少磁盘占用
//   - 无装饰：关闭颜色和加粗，纯净输出
//
// 参数:
//   - logDirName: 日志目录名称
//   - logFileName: 日志文件名称
//
// 返回值:
//   - *FastLogConfig: 生产模式配置实例
func ProdConfig(logDirName, logFileName string) *FastLogConfig {
	config := NewFastLogConfig(logDirName, logFileName)

	// 生产模式差异化设置
	config.OutputToConsole = false           // 生产模式仅文件输出
	config.ChanIntSize = largeChanSize       // 生产模式大缓冲区
	config.FlushInterval = slowFlushInterval // 生产模式慢刷新
	config.LogFormat = Json                  // 生产模式结构化
	config.Color = false                     // 生产模式无装饰
	config.Bold = false                      // 生产模式无加粗
	config.MaxAge = productionMaxAge         // 生产模式长期存储
	config.MaxFiles = productionMaxBackups   // 生产模式长期备份
	config.Compress = true                   // 生产模式压缩存储
	config.BatchSize = defaultBatchSize * 2  // 生产模式大批处理

	return config
}

// ConsoleConfig 控制台模式配置
// 特点：
//   - 纯控制台：仅控制台输出，无文件存储，适合临时调试
//   - 视觉友好：彩色+加粗显示，Basic格式简洁易读
//   - 快速响应：小缓冲区(5000)，快刷新(100ms)，小批处理(500)
//   - 轻量级：INFO级别，无文件操作开销
//   - 即时性：适合开发调试、脚本运行等临时场景
//
// 返回值:
//   - *FastLogConfig: 控制台模式配置实例
func ConsoleConfig() *FastLogConfig {
	config := NewFastLogConfig("", "")

	// 控制台模式差异化设置
	config.OutputToFile = false              // 控制台模式纯控制台
	config.ChanIntSize = smallChanSize       // 控制台模式小缓冲区
	config.FlushInterval = fastFlushInterval // 控制台模式快刷新
	config.MaxSize = 0                       // 控制台模式无文件
	config.MaxAge = 0                        // 控制台模式无文件
	config.MaxFiles = 0                      // 控制台模式无文件
	config.BatchSize = defaultBatchSize / 2  // 控制台模式小批处理
	config.DisableBackpressure = true        // 控制台模式禁用背压，便于实时查看

	return config
}

// FileConfig 文件模式配置
// 特点：
//   - 纯文件：仅文件输出，无控制台干扰，适合后台服务
//   - 结构化：BasicStructured格式，平衡可读性和解析性
//   - 中等性能：标准缓冲区和刷新间隔，平衡性能和实时性
//   - 中期存储：14天保留期，30个备份文件，适合一般业务
//   - 无装饰：关闭颜色和加粗，纯净文件输出
//
// 参数:
//   - logDirName: 日志目录名称
//   - logFileName: 日志文件名称
//
// 返回值:
//   - *FastLogConfig: 文件模式配置实例
func FileConfig(logDirName, logFileName string) *FastLogConfig {
	config := NewFastLogConfig(logDirName, logFileName)

	// 文件模式差异化设置
	config.OutputToConsole = false     // 文件模式纯文件
	config.LogFormat = BasicStructured // 文件模式结构化
	config.Color = false               // 文件模式无装饰
	config.Bold = false                // 文件模式无加粗
	config.MaxAge = fileMaxAge         // 文件模式中期存储
	config.MaxFiles = fileMaxBackups   // 文件模式中期备份

	return config
}

// SilentConfig 静默模式配置
// 特点：
//   - 最小输出：仅WARN级别，只记录警告和错误，极大减少日志量
//   - 极致性能：大缓冲区(20000)，最慢刷新(1000ms)，大批处理(2000)
//   - 高效格式：JsonSimple简化JSON，减少序列化开销
//   - 长期存储：30天保留期，50个备份文件，重要信息不丢失
//   - 空间优化：启用压缩，最大化存储效率
//   - 适用场景：高并发生产环境，性能敏感应用
//
// 参数:
//   - logDirName: 日志目录名称
//   - logFileName: 日志文件名称
//
// 返回值:
//   - *FastLogConfig: 静默模式配置实例
func SilentConfig(logDirName, logFileName string) *FastLogConfig {
	config := NewFastLogConfig(logDirName, logFileName)

	// 静默模式差异化设置
	config.OutputToConsole = false           // 静默模式纯文件
	config.LogLevel = WARN                   // 静默模式最小输出
	config.ChanIntSize = largeChanSize       // 静默模式大缓冲区
	config.FlushInterval = slowFlushInterval // 静默模式慢刷新
	config.LogFormat = JsonSimple            // 静默模式简单JSON
	config.Color = false                     // 静默模式无装饰
	config.Bold = false                      // 静默模式无加粗
	config.MaxAge = silentMaxAge             // 静默模式长期存储
	config.MaxFiles = silentMaxBackups       // 静默模式长期备份
	config.Compress = true                   // 静默模式压缩存储
	config.BatchSize = defaultBatchSize * 2  // 静默模式大批处理

	return config
}

// ========================================================================
// 内部辅助函数
// ========================================================================

// validateConfig 验证配置并设置默认值
// 发现任何不合理的配置值都会panic，确保调用者提供正确配置
func (c *FastLogConfig) validateConfig() {
	// 配置对象不能为nil
	if c == nil {
		panic("FastLogConfig cannot be nil")
	}

	// ========================================================================
	// 第一步：设置所有默认值
	// ========================================================================

	// 设置通道大小默认值
	if c.ChanIntSize == 0 {
		c.ChanIntSize = defaultChanSize
	}

	// 设置刷新间隔默认值
	if c.FlushInterval == 0 {
		c.FlushInterval = normalFlushInterval
	}

	// 设置批处理大小默认值
	if c.BatchSize == 0 {
		c.BatchSize = defaultBatchSize
	}

	// 设置日志级别默认值
	if c.LogLevel == 0 {
		c.LogLevel = INFO
	}

	// 设置日志格式默认值
	if c.LogFormat == 0 {
		c.LogFormat = Simple
	}

	// 设置文件大小默认值
	if c.MaxSize == 0 {
		c.MaxSize = defaultMaxFileSize
	}

	// 设置文件相关默认值（仅在启用文件输出时）
	if c.OutputToFile {
		if strings.TrimSpace(c.LogDirName) == "" {
			c.LogDirName = defaultLogDir
		}
		if strings.TrimSpace(c.LogFileName) == "" {
			c.LogFileName = defaultLogFileName
		}
	}

	// ========================================================================
	// 第二步：验证所有配置值
	// ========================================================================

	// 必须启用至少一种输出方式
	if !c.OutputToConsole && !c.OutputToFile {
		panic("at least one output method must be enabled: OutputToConsole or OutputToFile")
	}

	// 验证通道大小
	if c.ChanIntSize < 0 {
		panic("ChanIntSize cannot be negative")
	}
	if c.ChanIntSize > maxChanSize {
		panic(fmt.Sprintf("ChanIntSize %d exceeds maximum %d", c.ChanIntSize, maxChanSize))
	}

	// 验证刷新间隔
	if c.FlushInterval < 0 {
		panic("FlushInterval cannot be negative")
	}
	if c.FlushInterval > 0 && c.FlushInterval < minFlushInterval {
		panic(fmt.Sprintf("FlushInterval %v too small, minimum %v", c.FlushInterval, minFlushInterval))
	}
	if c.FlushInterval > maxFlushInterval {
		panic(fmt.Sprintf("FlushInterval %v exceeds maximum %v", c.FlushInterval, maxFlushInterval))
	}

	// 验证批处理大小
	if c.BatchSize < 0 {
		panic("BatchSize cannot be negative")
	}
	if c.BatchSize > maxBatchSize {
		panic(fmt.Sprintf("BatchSize %d exceeds maximum %d", c.BatchSize, maxBatchSize))
	}

	// 验证日志级别
	if c.LogLevel < DEBUG || c.LogLevel > NONE {
		panic(fmt.Sprintf("invalid LogLevel %d, must be %d-%d", c.LogLevel, DEBUG, NONE))
	}

	// 验证日志格式
	if c.LogFormat < Detailed || c.LogFormat > Custom {
		panic(fmt.Sprintf("invalid LogFormat %d, must be %d-%d", c.LogFormat, Detailed, Custom))
	}

	// 验证文件大小
	if c.MaxSize < 0 {
		panic("MaxSize cannot be negative")
	}
	if c.MaxSize > maxSingleFileSize {
		panic(fmt.Sprintf("MaxSize %d exceeds maximum %d MB", c.MaxSize, maxSingleFileSize))
	}

	// 验证保留天数
	if c.MaxAge < 0 {
		panic("MaxAge cannot be negative")
	}
	if c.MaxAge > maxRetentionDays {
		panic(fmt.Sprintf("MaxAge %d exceeds maximum %d days", c.MaxAge, maxRetentionDays))
	}

	// 验证保留文件数
	if c.MaxFiles < 0 {
		panic("MaxFiles cannot be negative")
	}
	if c.MaxFiles > maxRetentionFiles {
		panic(fmt.Sprintf("MaxFiles %d exceeds maximum %d files", c.MaxFiles, maxRetentionFiles))
	}

	// 验证文件输出相关配置
	if c.OutputToFile {
		// 验证路径安全性
		if strings.Contains(c.LogDirName, "..") {
			panic("LogDirName contains path traversal '..'")
		}
		if strings.Contains(c.LogFileName, "..") {
			panic("LogFileName contains path traversal '..'")
		}

		// 清理文件名
		c.LogDirName = cleanFileName(c.LogDirName)
		c.LogFileName = cleanFileName(c.LogFileName)
	}
}

// cleanFileName 清理文件名中的非法字符和格式问题
//
// 参数:
//   - filename: 原始文件名（可能包含路径）
//
// 返回:
//   - string: 清理后的文件名
func cleanFileName(filename string) string {
	// 处理空字符串
	if strings.TrimSpace(filename) == "" {
		return defaultLogFileName
	}

	// 1. 使用 filepath.Clean 进行路径规范化
	// 这会自动处理 "./"、多余的分隔符等问题
	cleaned := filepath.Clean(filename)

	// 2. 安全检查：移除上级目录引用
	if strings.Contains(cleaned, "..") {
		// 如果包含 ".."，只保留文件名部分
		cleaned = filepath.Base(cleaned)
	}

	// 3. 分离目录和文件名
	dir := filepath.Dir(cleaned)
	actualFileName := filepath.Base(cleaned)

	// 4. 清理文件名中的非法字符
	for _, char := range invalidFileChars {
		actualFileName = strings.ReplaceAll(actualFileName, char, charReplacement)
	}

	// 5. 移除文件名开头/结尾的点和空格
	actualFileName = strings.Trim(actualFileName, ". ")

	// 6. 处理文件名长度限制
	if len(actualFileName) > maxFileNameLength {
		ext := filepath.Ext(actualFileName)
		maxNameLen := maxFileNameLength - len(ext)
		if maxNameLen > 0 {
			actualFileName = actualFileName[:maxNameLen] + ext
		} else {
			actualFileName = actualFileName[:maxFileNameLength]
		}
	}

	// 7. 确保文件名不为空
	if actualFileName == "" {
		actualFileName = defaultLogFileName
	}

	// 8. 重新组合路径
	if dir == "." {
		cleaned = actualFileName
	} else {
		cleaned = filepath.Join(dir, actualFileName)
	}

	// 9. 最终路径长度检查
	if len(cleaned) > maxPathLength {
		cleaned = cleaned[:maxPathLength-truncateReserve] + truncatedSuffix
	}

	return cleaned
}

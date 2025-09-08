/*
config.go - 日志配置管理模块
定义日志配置结构体及配置项的设置与获取方法，负责管理FastLog的所有可配置参数。
*/
package fastlog

import (
	"path/filepath"
	"strings"
	"time"
)

// FastLogConfig 定义一个配置结构体, 用于配置日志记录器
type FastLogConfig struct {
	LogDirName      string        // 日志目录路径
	LogFileName     string        // 日志文件名
	OutputToConsole bool          // 是否将日志输出到控制台
	OutputToFile    bool          // 是否将日志输出到文件
	FlushInterval   time.Duration // 刷新间隔, 单位为time.Duration
	LogLevel        LogLevel      // 日志级别
	ChanIntSize     int           // 通道大小 默认10000
	LogFormat       LogFormatType // 日志格式选项
	Color           bool          // 是否启用终端颜色
	Bold            bool          // 是否启用终端字体加粗
	MaxLogFileSize  int           // 最大日志文件大小, 单位为MB, 默认10MB
	MaxLogAge       int           // 最大日志文件保留天数, 默认为0, 表示不做限制
	MaxLogBackups   int           // 最大日志文件保留数量, 默认为0, 表示不做限制
	IsLocalTime     bool          // 是否使用本地时间 默认使用UTC时间
	EnableCompress  bool          // 是否启用日志文件压缩 默认不启用
	BatchSize       int           // 批处理数量
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
		LogDirName:      logDirName,          // 日志目录名称
		LogFileName:     logFileName,         // 日志文件名称
		OutputToConsole: true,                // 是否将日志输出到控制台
		OutputToFile:    true,                // 是否将日志输出到文件
		LogLevel:        INFO,                // 日志级别 默认INFO
		ChanIntSize:     defaultChanSize,     // 通道大小 默认10000
		FlushInterval:   normalFlushInterval, // 刷新间隔 默认500毫秒
		LogFormat:       Simple,              // 日志格式选项
		MaxLogFileSize:  defaultMaxFileSize,  // 最大日志文件大小, 单位为MB, 默认10MB
		MaxLogAge:       0,                   // 最大日志文件保留天数, 默认为0, 表示不做限制
		MaxLogBackups:   0,                   // 最大日志文件保留数量, 默认为0, 表示不做限制
		IsLocalTime:     true,                // 是否使用本地时间 默认使用本地时间
		EnableCompress:  false,               // 是否启用日志文件压缩 默认不启用
		Color:           true,                // 是否启用终端颜色
		Bold:            true,                // 是否启用终端字体加粗
		BatchSize:       defaultBatchSize,    // 批处理数量
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
	config.LogLevel = DEBUG
	config.FlushInterval = fastFlushInterval
	config.LogFormat = Detailed
	config.MaxLogAge = developmentMaxAge
	config.MaxLogBackups = developmentMaxBackups
	
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
	config.OutputToConsole = false
	config.ChanIntSize = largeChanSize
	config.FlushInterval = slowFlushInterval
	config.LogFormat = Json
	config.Color = false
	config.Bold = false
	config.MaxLogAge = productionMaxAge
	config.MaxLogBackups = productionMaxBackups
	config.EnableCompress = true
	config.BatchSize = defaultBatchSize * 2
	
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
	config.OutputToFile = false
	config.ChanIntSize = smallChanSize
	config.FlushInterval = fastFlushInterval
	config.MaxLogFileSize = 0
	config.MaxLogAge = 0
	config.MaxLogBackups = 0
	config.BatchSize = defaultBatchSize / 2
	
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
	config.OutputToConsole = false
	config.LogFormat = BasicStructured
	config.Color = false
	config.Bold = false
	config.MaxLogAge = fileMaxAge
	config.MaxLogBackups = fileMaxBackups
	
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
	config.OutputToConsole = false
	config.LogLevel = WARN
	config.ChanIntSize = largeChanSize
	config.FlushInterval = slowFlushInterval
	config.LogFormat = JsonSimple
	config.Color = false
	config.Bold = false
	config.MaxLogAge = silentMaxAge
	config.MaxLogBackups = silentMaxBackups
	config.EnableCompress = true
	config.BatchSize = defaultBatchSize * 2
	
	return config
}

// ========================================================================
// 内部辅助函数
// ========================================================================

// fixFinalConfig 最终配置修正函数 - 在NewFastLog开始时调用
// 负责修正所有不合理的配置值, 确保系统稳定运行
// 对于无法修正的关键错误会panic
func (c *FastLogConfig) fixFinalConfig() {
	// 第一步：检查系统资源充足性
	c.checkSystemResources()

	// 第二步：修正可以修正的配置项
	c.fixFileConfig()

	// 第三步：验证关键配置，无法修正时panic
	c.validateCriticalConfig()
	c.fixPerformanceConfig()
	c.fixLogConfig()

	// 第四步：最终一致性检查
	c.validateFinalConsistency()
}

// validateCriticalConfig 验证关键配置，无法修正时panic
func (c *FastLogConfig) validateCriticalConfig() {
	// 配置对象不能为nil
	if c == nil {
		panic("FastLogConfig configuration object cannot be nil")
	}

	// 必须启用至少一种输出方式
	if !c.OutputToConsole && !c.OutputToFile {
		panic("At least one output method must be enabled: console output (OutputToConsole) or file output (OutputToFile)")
	}

	// 如果启用文件输出，目录名和文件名不能同时为空
	if c.OutputToFile {
		if strings.TrimSpace(c.LogDirName) == "" && strings.TrimSpace(c.LogFileName) == "" {
			panic("When file output is enabled, log directory name (LogDirName) and file name (LogFileName) cannot both be empty")
		}
	}
}

// checkSystemResources 检查系统资源是否充足
func (c *FastLogConfig) checkSystemResources() {
	// 检查通道大小是否会占用过多内存
	if c.ChanIntSize > maxChanSize { // 超过100万条消息
		panic("channel size too large, may cause memory overflow. Recommend setting within 1 million entries")
	}

	// 检查刷新间隔是否过小导致CPU占用过高
	if c.FlushInterval > 0 && c.FlushInterval < minFlushInterval {
		panic("refresh interval too small (less than 1 microsecond), will cause high CPU usage")
	}

	// 检查文件大小配置是否合理
	if c.MaxLogFileSize > maxSingleFileSize { // 超过10GB
		panic("single log file size too large (exceeds 10GB), may cause insufficient disk space")
	}
}

// fixFileConfig 修正文件相关配置
func (c *FastLogConfig) fixFileConfig() {
	// 只在文件输出模式下修正
	if !c.OutputToFile {
		return
	}

	// 1. 修正基本字符串字段
	if strings.TrimSpace(c.LogDirName) == "" {
		c.LogDirName = defaultLogDir
	}

	if strings.TrimSpace(c.LogFileName) == "" {
		c.LogFileName = defaultLogFileName
	}

	// 2. 清理文件名中的非法字符
	originalDir := c.LogDirName
	cleanedDir := cleanFileName(originalDir)
	if originalDir != cleanedDir {
		c.LogDirName = cleanedDir
	}

	originalFile := c.LogFileName
	cleanedFile := cleanFileName(originalFile)
	if originalFile != cleanedFile {
		c.LogFileName = cleanedFile
	}

	// 3. 修正文件轮转配置
	if c.MaxLogFileSize <= 0 {
		c.MaxLogFileSize = defaultMaxFileSize
	} else if c.MaxLogFileSize > maxSingleFileSize {
		c.MaxLogFileSize = maxSingleFileSize
	}

	if c.MaxLogAge < 0 {
		c.MaxLogAge = 0
	} else if c.MaxLogAge > maxRetentionDays { // 最多保留10年
		c.MaxLogAge = maxRetentionDays
	}

	if c.MaxLogBackups < 0 {
		c.MaxLogBackups = 0
	} else if c.MaxLogBackups > maxRetentionFiles {
		c.MaxLogBackups = maxRetentionFiles
	}
}

// fixPerformanceConfig 修正性能相关配置
func (c *FastLogConfig) fixPerformanceConfig() {
	// 修正通道大小
	if c.ChanIntSize <= 0 {
		c.ChanIntSize = defaultChanSize
	} else if c.ChanIntSize > chanSizeLimit {
		c.ChanIntSize = chanSizeLimit
	}

	// 修正刷新间隔
	if c.FlushInterval <= 0 {
		c.FlushInterval = normalFlushInterval
	} else if c.FlushInterval < normalMinFlush {
		c.FlushInterval = normalMinFlush
	} else if c.FlushInterval > maxFlushInterval {
		c.FlushInterval = maxFlushInterval
	}

	// 修正批处理数量
	if c.BatchSize <= 0 {
		c.BatchSize = defaultBatchSize
	} else if c.BatchSize > maxBatchSize {
		c.BatchSize = maxBatchSize
	}
}

// fixLogConfig 修正日志级别和格式配置
func (c *FastLogConfig) fixLogConfig() {
	// 修正日志级别
	if c.LogLevel < DEBUG || c.LogLevel > NONE {
		c.LogLevel = INFO
	}

	// 修正日志格式
	if c.LogFormat < Detailed || c.LogFormat > Custom {
		c.LogFormat = Simple
	}
}

// validateFinalConsistency 最终一致性检查
func (c *FastLogConfig) validateFinalConsistency() {
	// 检查配置组合是否合理
	if c.OutputToFile {
		// 如果启用文件输出但文件轮转配置可能导致日志丢失
		if c.MaxLogAge > 0 && c.MaxLogBackups > 0 {
			// 计算可能的最大日志保留量
			maxRetentionDays := c.MaxLogAge
			maxFiles := c.MaxLogBackups

			// 如果配置过于激进，调整为更保守的值
			if maxRetentionDays < minRetentionDays && maxFiles < minRetentionFiles {
				// 至少保留7天或5个文件
				if c.MaxLogAge > 0 && c.MaxLogAge < minRetentionDays {
					c.MaxLogAge = minRetentionDays
				}
				if c.MaxLogBackups > 0 && c.MaxLogBackups < minRetentionFiles {
					c.MaxLogBackups = minRetentionFiles
				}
			}
		}
	}

	// 检查性能配置的合理性
	if c.ChanIntSize > performanceThreshold && c.FlushInterval < performanceFlushMin {
		// 大通道配合高频刷新可能导致性能问题，调整刷新间隔
		c.FlushInterval = performanceFlushMin
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

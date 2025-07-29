// config.go - 日志配置管理模块
// 定义日志配置结构体及配置项的设置与获取方法,
// 负责管理FastLog的所有可配置参数。
package fastlog

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gitee.com/MM-Q/colorlib"
)

// FastLogConfigurer 定义日志配置器接口, 包含所有配置项的设置和获取方法
type FastLogConfigurer interface {
	// SetLogDirName 设置日志目录路径
	SetLogDirName(dirName string)
	// GetLogDirName 获取日志目录路径
	GetLogDirName() string

	// SetLogFileName 设置日志文件名
	SetLogFileName(fileName string)
	// GetLogFileName 获取日志文件名
	GetLogFileName() string

	// SetPrintToConsole 设置是否将日志输出到控制台
	SetPrintToConsole(print bool)
	// GetPrintToConsole 获取是否将日志输出到控制台的状态
	GetPrintToConsole() bool

	// SetConsoleOnly 设置是否仅输出到控制台
	SetConsoleOnly(only bool)
	// GetConsoleOnly 获取是否仅输出到控制台的状态
	GetConsoleOnly() bool

	// SetFlushInterval 设置刷新间隔
	SetFlushInterval(interval time.Duration)
	// GetFlushInterval 获取刷新间隔
	GetFlushInterval() time.Duration

	// SetLogLevel 设置日志级别
	SetLogLevel(level LogLevel)
	// GetLogLevel 获取日志级别
	GetLogLevel() LogLevel

	// SetChanIntSize 设置通道大小
	SetChanIntSize(size int)
	// GetChanIntSize 获取通道大小
	GetChanIntSize() int

	// SetLogFormat 设置日志格式选项
	SetLogFormat(format LogFormatType)
	// GetLogFormat 获取日志格式选项
	GetLogFormat() LogFormatType

	// SetNoColor 设置是否禁用终端颜色
	SetNoColor(noColor bool)
	// GetNoColor 获取是否禁用终端颜色的状态
	GetNoColor() bool

	// SetNoBold 设置是否禁用终端字体加粗
	SetNoBold(noBold bool)
	// GetNoBold 获取是否禁用终端字体加粗的状态
	GetNoBold() bool

	// SetMaxLogFileSize 设置最大日志文件大小(MB)
	SetMaxLogFileSize(size int)
	// GetMaxLogFileSize 获取最大日志文件大小(MB)
	GetMaxLogFileSize() int

	// SetMaxLogAge 设置最大日志文件保留天数
	SetMaxLogAge(age int)
	// GetMaxLogAge 获取最大日志文件保留天数
	GetMaxLogAge() int

	// SetMaxLogBackups 设置最大日志文件保留数量
	SetMaxLogBackups(backups int)
	// GetMaxLogBackups 获取最大日志文件保留数量
	GetMaxLogBackups() int

	// SetIsLocalTime 设置是否使用本地时间
	SetIsLocalTime(local bool)
	// GetIsLocalTime 获取是否使用本地时间的状态
	GetIsLocalTime() bool

	// SetEnableCompress 设置是否启用日志文件压缩
	SetEnableCompress(compress bool)
	// GetEnableCompress 获取是否启用日志文件压缩的状态
	GetEnableCompress() bool
}

// FastLogConfig 定义一个配置结构体, 用于配置日志记录器
type FastLogConfig struct {
	logDirName     string        // 日志目录路径
	logFileName    string        // 日志文件名
	printToConsole bool          // 是否将日志输出到控制台
	consoleOnly    bool          // 是否仅输出到控制台
	flushInterval  time.Duration // 刷新间隔, 单位为time.Duration
	logLevel       LogLevel      // 日志级别
	chanIntSize    int           // 通道大小 默认10000
	logFormat      LogFormatType // 日志格式选项
	noColor        bool          // 是否禁用终端颜色
	noBold         bool          // 是否禁用终端字体加粗
	maxLogFileSize int           // 最大日志文件大小, 单位为MB, 默认5MB
	maxLogAge      int           // 最大日志文件保留天数, 默认为0, 表示不做限制
	maxLogBackups  int           // 最大日志文件保留数量, 默认为0, 表示不做限制
	isLocalTime    bool          // 是否使用本地时间 默认使用UTC时间
	enableCompress bool          // 是否启用日志文件压缩 默认不启用
	setMu          sync.RWMutex  // 用于保护配置的锁 (读写锁优化并发读性能)
}

// validateFinalConfig 最终配置验证函数 - 在NewFastLog开始时调用
// 只负责检查配置并打印警告信息, 不修改配置
func (c *FastLogConfig) validateFinalConfig() {
	// 获取打印对象
	cl := colorlib.NewColorLib()
	cl.SetNoColor(true) // 禁用颜色

	// 1. 检查通道大小
	if c.GetChanIntSize() <= 0 {
		cl.PrintErrf("通道大小无效: %d (必须大于0)\n", c.GetChanIntSize())

	} else if c.GetChanIntSize() < 1000 {
		cl.PrintWarnf("通道大小过小: %d (建议至少1000, 避免日志阻塞)\n", c.GetChanIntSize())

	} else if c.GetChanIntSize() > 100000 {
		cl.PrintWarnf("通道大小过大: %d (可能导致内存占用过高)\n", c.GetChanIntSize())

	} else if c.GetChanIntSize() > 50000 {
		cl.PrintInff("通道大小较大: %d (请确保有足够内存)\n", c.GetChanIntSize())

	}

	// 2. 检查刷新间隔
	interval := c.GetFlushInterval()
	if interval <= 0 {
		cl.PrintErrf("刷新间隔无效: %v (必须大于0)\n", interval)

	} else if interval < 10*time.Millisecond {
		cl.PrintErrf("刷新间隔过短: %v (可能严重影响性能)\n", interval)

	} else if interval < 50*time.Millisecond {
		cl.PrintWarnf("刷新间隔较短: %v (可能影响性能, 建议100ms以上)\n", interval)

	} else if interval > 30*time.Second {
		cl.PrintWarnf("刷新间隔过长: %v (可能导致日志丢失)\n", interval)

	} else if interval > 10*time.Second {
		cl.PrintInff("刷新间隔较长: %v (可能影响日志实时性)\n", interval)

	}

	// 3. 检查最大文件大小
	fileSize := c.GetMaxLogFileSize()
	if fileSize <= 0 {
		cl.PrintErrf("最大文件大小无效: %dMB (必须大于0)\n", fileSize)

	} else if fileSize < 1 {
		cl.PrintWarnf("文件大小过小: %dMB (可能导致频繁切割)\n", fileSize)

	} else if fileSize > 1000 {
		cl.PrintWarnf("文件大小过大: %dMB (建议不超过1000MB)\n", fileSize)

	} else if fileSize > 500 {
		cl.PrintInff("文件大小较大: %dMB (请确保有足够磁盘空间)\n", fileSize)

	}

	// 4. 检查最大保留天数
	maxAge := c.GetMaxLogAge()
	if maxAge < 0 {
		cl.PrintErrf("最大保留天数无效: %d (不能为负数)\n", maxAge)

	} else if maxAge > 3650 {
		cl.PrintWarnf("保留天数过长: %d天 (超过10年, 建议调整)\n", maxAge)

	} else if maxAge > 365 {
		cl.PrintInff("保留天数较长: %d天 (请确保有足够磁盘空间)\n", maxAge)

	} else if maxAge > 0 && maxAge < 1 {
		cl.PrintWarnf("保留天数过短: %d天 (可能导致重要日志丢失)\n", maxAge)

	}

	// 5. 检查最大备份数量
	maxBackups := c.GetMaxLogBackups()
	if maxBackups < 0 {
		cl.PrintErrf("最大备份数量无效: %d (不能为负数)\n", maxBackups)

	} else if maxBackups > 1000 {
		cl.PrintWarnf("备份数量过多: %d个 (建议不超过1000个)\n", maxBackups)

	} else if maxBackups > 100 {
		cl.PrintInff("备份数量较多: %d个 (请确保有足够磁盘空间)\n", maxBackups)

	} else if maxBackups > 0 && maxBackups < 3 {
		cl.PrintInff("备份数量较少: %d个 (建议至少保留3个备份)\n", maxBackups)

	}

	// 6. 检查日志级别
	logLevel := c.GetLogLevel()
	if logLevel < DEBUG || logLevel > NONE {
		cl.PrintErrf("无效的日志级别: %d (有效范围: %d-%d)\n", logLevel, DEBUG, NONE)

	}

	// 7. 检查日志格式
	logFormat := c.GetLogFormat()
	if logFormat < Json || logFormat > Custom {
		cl.PrintErrf("无效的日志格式: %d (有效范围: %d-%d)\n", logFormat, Json, Custom)

	}

	// 8. 检查文件名和目录名（仅在非控制台模式下）
	if !c.GetConsoleOnly() {
		dirName := c.GetLogDirName()
		fileName := c.GetLogFileName()

		if dirName == "" {
			cl.PrintErrf("日志目录名为空 (非控制台模式下必须指定)\n")

		} else if containsInvalidChars(dirName) {
			cl.PrintErrf("日志目录名包含非法字符: '%s'\n", dirName)

		}

		if fileName == "" {
			cl.PrintErrf("日志文件名为空 (非控制台模式下必须指定)\n")

		} else if containsInvalidChars(fileName) {
			cl.PrintErrf("日志文件名包含非法字符: '%s'\n", fileName)

		}
	}

	// 9. 检查配置组合的合理性
	if maxAge > 0 && maxBackups > 0 {
		estimatedFiles := maxAge * 24 // 假设每小时一个文件
		if estimatedFiles > maxBackups*10 {
			cl.PrintInff("配置不平衡: 保留天数(%d)可能产生的文件数远超备份限制(%d)\n", maxAge, maxBackups)

		}
	}

	// 10. 检查性能相关配置
	if c.GetChanIntSize() < 1000 && interval > 1*time.Second {
		cl.PrintInff("性能配置不佳: 通道小(%d)且刷新间隔长(%v), 可能导致阻塞\n", c.GetChanIntSize(), interval)

	}
}

// fixFinalConfig 最终配置修正函数 - 在NewFastLog开始时调用
// 负责修正所有不合理的配置值, 确保系统稳定运行
func (c *FastLogConfig) fixFinalConfig() {
	// 获取打印对象
	cl := colorlib.NewColorLib()
	cl.SetNoColor(true) // 禁用颜色

	// 1. 修正基本字符串字段
	if c.GetLogDirName() == "" && !c.GetConsoleOnly() {
		c.SetLogDirName("logs")
		cl.PrintOk("修正日志目录名: 空值 -> 'logs'")
	}

	if c.GetLogFileName() == "" && !c.GetConsoleOnly() {
		c.SetLogFileName("app.log")
		cl.PrintOk("修正日志文件名: 空值 -> 'app.log'")
	}

	// 2. 清理文件名中的非法字符
	if !c.GetConsoleOnly() {
		originalDir := c.GetLogDirName()
		cleanedDir := cleanFileName(originalDir)
		if originalDir != cleanedDir {
			c.SetLogDirName(cleanedDir)
			cl.PrintOkf("修正日志目录名: '%s' -> '%s' (清理非法字符)\n", originalDir, cleanedDir)

		}

		originalFile := c.GetLogFileName()
		cleanedFile := cleanFileName(originalFile)
		if originalFile != cleanedFile {
			c.SetLogFileName(cleanedFile)
			cl.PrintOkf("修正日志文件名: '%s' -> '%s' (清理非法字符)\n", originalFile, cleanedFile)

		}
	}

	// 3. 修正通道大小
	originalChan := c.GetChanIntSize()
	if originalChan <= 0 {
		c.SetChanIntSize(10000)
		cl.PrintOkf("修正通道大小: %d -> 10000 (默认值)\n", originalChan)

	} else if originalChan > 100000 {
		c.SetChanIntSize(100000)
		cl.PrintOkf("修正通道大小: %d -> 100000 (最大值)\n", originalChan)

	}

	// 4. 修正刷新间隔
	originalInterval := c.GetFlushInterval()
	if originalInterval <= 0 {
		c.SetFlushInterval(500 * time.Millisecond)
		cl.PrintOkf("修正刷新间隔: %v -> 500ms (默认值)\n", originalInterval)

	} else if originalInterval < 10*time.Millisecond {
		c.SetFlushInterval(10 * time.Millisecond)
		cl.PrintOkf("修正刷新间隔: %v -> 10ms (最小值)\n", originalInterval)

	} else if originalInterval > 30*time.Second {
		c.SetFlushInterval(30 * time.Second)
		cl.PrintOkf("修正刷新间隔: %v -> 30s (最大值)\n", originalInterval)

	}

	// 5. 修正最大文件大小
	originalSize := c.GetMaxLogFileSize()
	if originalSize <= 0 {
		c.SetMaxLogFileSize(5)
		cl.PrintOkf("修正最大文件大小: %dMB -> 5MB (默认值)\n", originalSize)

	} else if originalSize > 1000 {
		c.SetMaxLogFileSize(1000)
		cl.PrintOkf("修正最大文件大小: %dMB -> 1000MB (最大值)\n", originalSize)

	}

	// 6. 修正最大保留天数
	originalAge := c.GetMaxLogAge()
	if originalAge < 0 {
		c.SetMaxLogAge(0)
		cl.PrintOkf("修正最大保留天数: %d -> 0 (不限制)\n", originalAge)

	} else if originalAge > 3650 {
		c.SetMaxLogAge(3650)
		cl.PrintOkf("修正最大保留天数: %d -> 3650天 (最大值)\n", originalAge)

	}

	// 7. 修正最大备份数量
	originalBackups := c.GetMaxLogBackups()
	if originalBackups < 0 {
		c.SetMaxLogBackups(0)
		cl.PrintOkf("修正最大备份数量: %d -> 0 (不限制)\n", originalBackups)

	} else if originalBackups > 1000 {
		c.SetMaxLogBackups(1000)
		cl.PrintOkf("修正最大备份数量: %d -> 1000个 (最大值)\n", originalBackups)

	}

	// 8. 修正日志级别
	originalLevel := c.GetLogLevel()
	if originalLevel < DEBUG || originalLevel > NONE {
		c.SetLogLevel(INFO)
		cl.PrintOkf("修正日志级别: %d -> %d (INFO)\n", originalLevel, INFO)

	}

	// 9. 修正日志格式
	originalFormat := c.GetLogFormat()
	if originalFormat < Json || originalFormat > Custom {
		c.SetLogFormat(Detailed)
		cl.PrintOkf("修正日志格式: %d -> %d (Detailed)\n", originalFormat, Detailed)

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
		return "default.log"
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
		actualFileName = strings.ReplaceAll(actualFileName, char, "_")
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
		actualFileName = "app.log"
	}

	// 8. 重新组合路径
	if dir == "." {
		cleaned = actualFileName
	} else {
		cleaned = filepath.Join(dir, actualFileName)
	}

	// 9. 最终路径长度检查
	if len(cleaned) > maxPathLength {
		cleaned = cleaned[:maxPathLength-10] + "_truncated"
	}

	return cleaned
}

// containsInvalidChars 检查文件名是否包含非法字符或格式问题
//
// 参数:
//   - filename: 文件名（可能包含路径）
//
// 返回:
//   - bool: 是否包含非法字符或格式问题
func containsInvalidChars(filename string) bool {
	// 检查是否为空
	if strings.TrimSpace(filename) == "" {
		return true
	}

	// 使用 filepath.Clean 清理路径
	cleanedPath := filepath.Clean(filename)

	// 检查总路径长度
	if len(cleanedPath) > maxPathLength {
		return true
	}

	// 检查是否包含上级目录引用（安全检查）
	if strings.Contains(cleanedPath, "..") {
		return true
	}

	// 提取实际文件名
	actualFileName := filepath.Base(cleanedPath)

	// 检查文件名长度
	if len(actualFileName) > maxFileNameLength {
		return true
	}

	// 检查文件名中的非法字符
	for _, char := range invalidFileChars {
		if strings.Contains(actualFileName, char) {
			return true
		}
	}

	// 检查文件名是否以点或空格开头/结尾
	if strings.HasPrefix(actualFileName, ".") || strings.HasPrefix(actualFileName, " ") ||
		strings.HasSuffix(actualFileName, ".") || strings.HasSuffix(actualFileName, " ") {
		return true
	}

	return false
}

// SetLogDirName 设置日志目录路径
func (c *FastLogConfig) SetLogDirName(dirName string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.logDirName = dirName
}

// SetLogFileName 设置日志文件名
func (c *FastLogConfig) SetLogFileName(fileName string) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.logFileName = fileName
}

// SetPrintToConsole 设置是否将日志输出到控制台
func (c *FastLogConfig) SetPrintToConsole(print bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.printToConsole = print
}

// SetConsoleOnly 设置是否仅输出到控制台
func (c *FastLogConfig) SetConsoleOnly(only bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.consoleOnly = only
}

// SetFlushInterval 设置刷新间隔
func (c *FastLogConfig) SetFlushInterval(interval time.Duration) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.flushInterval = interval
}

// SetLogLevel 设置日志级别
func (c *FastLogConfig) SetLogLevel(level LogLevel) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.logLevel = level
}

// SetChanIntSize 设置通道大小
func (c *FastLogConfig) SetChanIntSize(size int) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.chanIntSize = size
}

// SetLogFormat 设置日志格式选项
func (c *FastLogConfig) SetLogFormat(format LogFormatType) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.logFormat = format
}

// SetNoColor 设置是否禁用终端颜色
func (c *FastLogConfig) SetNoColor(noColor bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.noColor = noColor
}

// SetNoBold 设置是否禁用终端字体加粗
func (c *FastLogConfig) SetNoBold(noBold bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.noBold = noBold
}

// SetMaxLogFileSize 设置最大日志文件大小(MB)
func (c *FastLogConfig) SetMaxLogFileSize(size int) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.maxLogFileSize = size
}

// SetMaxLogAge 设置最大日志文件保留天数
func (c *FastLogConfig) SetMaxLogAge(age int) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.maxLogAge = age
}

// SetMaxLogBackups 设置最大日志文件保留数量
func (c *FastLogConfig) SetMaxLogBackups(backups int) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.maxLogBackups = backups
}

// SetIsLocalTime 设置是否使用本地时间
func (c *FastLogConfig) SetIsLocalTime(local bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.isLocalTime = local
}

// SetEnableCompress 设置是否启用日志文件压缩
func (c *FastLogConfig) SetEnableCompress(compress bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.enableCompress = compress
}

// GetLogDirName 获取日志目录路径
func (c *FastLogConfig) GetLogDirName() string {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.logDirName
}

// GetLogFileName 获取日志文件名
func (c *FastLogConfig) GetLogFileName() string {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.logFileName
}

// GetPrintToConsole 获取是否将日志输出到控制台的状态
func (c *FastLogConfig) GetPrintToConsole() bool {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.printToConsole
}

// GetConsoleOnly 获取是否仅输出到控制台的状态
func (c *FastLogConfig) GetConsoleOnly() bool {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.consoleOnly
}

// GetFlushInterval 获取刷新间隔
func (c *FastLogConfig) GetFlushInterval() time.Duration {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.flushInterval
}

// GetLogLevel 获取日志级别
func (c *FastLogConfig) GetLogLevel() LogLevel {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.logLevel
}

// GetChanIntSize 获取通道大小
func (c *FastLogConfig) GetChanIntSize() int {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.chanIntSize
}

// GetLogFormat 获取日志格式选项
func (c *FastLogConfig) GetLogFormat() LogFormatType {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.logFormat
}

// GetNoColor 获取是否禁用终端颜色的状态
func (c *FastLogConfig) GetNoColor() bool {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.noColor
}

// GetNoBold 获取是否禁用终端字体加粗的状态
func (c *FastLogConfig) GetNoBold() bool {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.noBold
}

// GetMaxLogFileSize 获取最大日志文件大小(MB)
func (c *FastLogConfig) GetMaxLogFileSize() int {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.maxLogFileSize
}

// GetMaxLogAge 获取最大日志文件保留天数
func (c *FastLogConfig) GetMaxLogAge() int {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.maxLogAge
}

// GetMaxLogBackups 获取最大日志文件保留数量
func (c *FastLogConfig) GetMaxLogBackups() int {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.maxLogBackups
}

// GetIsLocalTime 获取是否使用本地时间的状态
func (c *FastLogConfig) GetIsLocalTime() bool {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.isLocalTime
}

// GetEnableCompress 获取是否启用日志文件压缩的状态
func (c *FastLogConfig) GetEnableCompress() bool {
	c.setMu.RLock()         // 读取锁定
	defer c.setMu.RUnlock() // 读取解锁
	return c.enableCompress
}

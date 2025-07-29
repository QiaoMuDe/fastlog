package fastlog

import (
	"sync"
	"time"
)

// FastLogConfigurer 定义日志配置器接口，包含所有配置项的设置和获取方法
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

// FastLogConfig 定义一个配置结构体，用于配置日志记录器
type FastLogConfig struct {
	logDirName     string        // 日志目录路径
	logFileName    string        // 日志文件名
	printToConsole bool          // 是否将日志输出到控制台
	consoleOnly    bool          // 是否仅输出到控制台
	flushInterval  time.Duration // 刷新间隔，单位为time.Duration
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

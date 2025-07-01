package fastlog

import (
	"time"
)

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

// SetMaxBufferSize 设置最大缓冲区大小(MB)
func (c *FastLogConfig) SetMaxBufferSize(size int) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.maxBufferSize = size
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
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.logDirName
}

// GetLogFileName 获取日志文件名
func (c *FastLogConfig) GetLogFileName() string {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.logFileName
}

// GetPrintToConsole 获取是否将日志输出到控制台的状态
func (c *FastLogConfig) GetPrintToConsole() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.printToConsole
}

// GetConsoleOnly 获取是否仅输出到控制台的状态
func (c *FastLogConfig) GetConsoleOnly() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.consoleOnly
}

// GetFlushInterval 获取刷新间隔
func (c *FastLogConfig) GetFlushInterval() time.Duration {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.flushInterval
}

// GetLogLevel 获取日志级别
func (c *FastLogConfig) GetLogLevel() LogLevel {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.logLevel
}

// GetChanIntSize 获取通道大小
func (c *FastLogConfig) GetChanIntSize() int {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.chanIntSize
}

// GetLogFormat 获取日志格式选项
func (c *FastLogConfig) GetLogFormat() LogFormatType {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.logFormat
}

// GetMaxBufferSize 获取最大缓冲区大小(MB)
func (c *FastLogConfig) GetMaxBufferSize() int {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.maxBufferSize
}

// GetNoColor 获取是否禁用终端颜色的状态
func (c *FastLogConfig) GetNoColor() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.noColor
}

// GetNoBold 获取是否禁用终端字体加粗的状态
func (c *FastLogConfig) GetNoBold() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.noBold
}

// GetMaxLogFileSize 获取最大日志文件大小(MB)
func (c *FastLogConfig) GetMaxLogFileSize() int {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.maxLogFileSize
}

// GetMaxLogAge 获取最大日志文件保留天数
func (c *FastLogConfig) GetMaxLogAge() int {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.maxLogAge
}

// GetMaxLogBackups 获取最大日志文件保留数量
func (c *FastLogConfig) GetMaxLogBackups() int {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.maxLogBackups
}

// GetIsLocalTime 获取是否使用本地时间的状态
func (c *FastLogConfig) GetIsLocalTime() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.isLocalTime
}

// GetEnableCompress 获取是否启用日志文件压缩的状态
func (c *FastLogConfig) GetEnableCompress() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.enableCompress
}

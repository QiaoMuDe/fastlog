package fastlog // import "gitee.com/MM-Q/fastlog/v2"

定义全局常量变量或结构体

用于存放fastlog的方法

TYPES

// 日志记录器
type FastLog struct {
	/*  私有属性 内部使用无需修改  */
	logFilePath   string             // 日志文件路径  内部拼接的 [logDirName+logFileName]
	logChan       chan *logMessage   // 日志通道  用于异步写入日志文件
	logWait       sync.WaitGroup     // 等待组 用于等待所有goroutine完成
	fileWriter    io.Writer          // 文件写入器
	fileMu        sync.Mutex         // 文件锁 用于保护文件缓冲区的写入操作
	consoleMu     sync.Mutex         // 控制台锁 用于保护控制台缓冲区的写入操作
	consoleWriter io.Writer          // 控制台写入器
	startOnce     sync.Once          // 用于确保日志处理器只启动一次
	ctx           context.Context    // 控制刷新器的上下文
	cancel        context.CancelFunc // 控制刷新器的取消函数
	cl            *colorlib.ColorLib // 提供终端颜色输出的库

	/* logrotatex 日志文件切割 */
	logGer *logrotatex.LogRotateX // 日志文件切割器

	/* 嵌入的配置结构体 */
	config *FastLogConfig // 配置结构体

	// 双缓冲区配置
	fileBuffers      [2]*bytes.Buffer // 文件双缓冲区
	fileBufferIdx    atomic.Int32     // 当前使用的文件缓冲区索引
	consoleBuffers   [2]*bytes.Buffer // 控制台双缓冲区
	consoleBufferIdx atomic.Int32     // 当前使用的控制台缓冲区索引
	fileBufferMu     sync.Mutex       // 文件缓冲区锁
	consoleBufferMu  sync.Mutex       // 控制台缓冲区锁

	// 用于控制缓冲区刷新的锁
	flushLock sync.Mutex

	// 用于控制关闭过程的锁
	closeLock sync.Mutex
}

// NewFastLog 创建一个新的FastLog实例, 用于记录日志。
// 参数:
//   - config: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。
//
// 返回值:
//   - *FastLog: 一个指向FastLog实例的指针。
//   - error: 如果创建日志记录器失败, 则返回一个错误。
func NewFastLog(config *FastLogConfig) (*FastLog, error) {
	// 检查配置结构体是否为nil
	if config == nil {
		return nil, fmt.Errorf("FastLogConfig 为 nil")
	}

	// 声明一些配置变量
	var (
		fileWriter    io.Writer // 文件写入器
		consoleWriter io.Writer // 控制台写入器
	)

	// 如果允许将日志输出到控制台, 或者仅输出到控制台, 则初始化控制台写入器。
	if config.GetConsoleOnly() || config.GetPrintToConsole() {
		consoleWriter = os.Stdout // 控制台写入器
	} else {
		consoleWriter = io.Discard // 不输出到控制台, 直接丢弃
	}

	// 拼接日志文件路径
	var logFilePath string
	// 如果日志目录名称和日志文件名称都不为空, 并且不是仅输出到控制台, 则拼接日志文件路径。
	if config.GetLogDirName() != "" || config.GetLogFileName() != "" && !config.GetConsoleOnly() {
		logFilePath = filepath.Join(config.GetLogDirName(), config.GetLogFileName())
	}

	// 如果不是仅输出到控制台, 则初始化日志文件写入器。
	var logger *logrotatex.LogRotateX
	if !config.GetConsoleOnly() {
		// 检查日志目录是否存在, 如果不存在则创建。
		if _, checkPathErr := checkPath(config.GetLogDirName()); checkPathErr != nil {
			if mkdirErr := os.MkdirAll(config.GetLogDirName(), 0644); mkdirErr != nil {
				return nil, fmt.Errorf("创建日志目录失败: %s", mkdirErr)
			}
		}

		// 初始化日志文件切割器
		logger = &logrotatex.LogRotateX{
			Filename:   logFilePath,                // 日志文件路径,
			MaxSize:    config.GetMaxLogFileSize(), // 最大日志文件大小, 单位为MB
			MaxAge:     config.GetMaxLogAge(),      // 最大日志文件保留天数
			MaxBackups: config.GetMaxLogBackups(),  // 最大日志文件保留数量
			LocalTime:  config.GetIsLocalTime(),    // 是否使用本地时间
			Compress:   config.GetEnableCompress(), // 是否启用日志文件压缩
		}

		// 初始化文件写入器
		fileWriter = logger
	} else {
		fileWriter = io.Discard // 仅输出到控制台, 不输出到文件
	}

	// 初始化双缓冲区
	fileBuffers := [2]*bytes.Buffer{
		bytes.NewBuffer(make([]byte, config.GetMaxBufferSize())),
		bytes.NewBuffer(make([]byte, config.GetMaxBufferSize())),
	}
	consoleBuffers := [2]*bytes.Buffer{
		bytes.NewBuffer(make([]byte, config.GetMaxBufferSize())),
		bytes.NewBuffer(make([]byte, config.GetMaxBufferSize())),
	}

	// 清空缓冲区
	fileBuffers[0].Reset()
	fileBuffers[1].Reset()
	consoleBuffers[0].Reset()
	consoleBuffers[1].Reset()

	// 创建一个新的FastLog实例, 将配置和缓冲区赋值给实例。
	f := &FastLog{
		logGer:         logger,                                          // 日志文件切割器
		fileWriter:     fileWriter,                                      // 文件写入器,
		consoleWriter:  consoleWriter,                                   // 控制台写入器,
		logFilePath:    logFilePath,                                     // 日志文件路径
		cl:             colorlib.NewColorLib(),                          // 颜色库实例
		config:         config,                                          // 配置结构体
		logChan:        make(chan *logMessage, config.GetChanIntSize()), // 日志消息通道
		fileBuffers:    fileBuffers,                                     // 文件缓冲
		consoleBuffers: consoleBuffers,                                  // 控制台缓冲
	}

	// 根据noColor的值, 设置颜色库的颜色选项
	if f.config.GetNoColor() {
		f.cl.NoColor.Store(true) // 设置颜色库的颜色选项为禁用
	}

	// 根据noBold的值, 设置颜色库的字体加粗选项
	if f.config.GetNoBold() {
		f.cl.NoBold.Store(true) // 设置颜色库的字体加粗选项为禁用
	}

	// 设置缓冲区索引为0
	f.fileBufferIdx.Store(0)
	f.consoleBufferIdx.Store(0)

	// 创建 context 用于控制协程退出
	f.ctx, f.cancel = context.WithCancel(context.Background())

	// 使用 sync.Once 确保日志处理器只启动一次
	f.startOnce.Do(func() {
		go f.processLogs() // 启动日志处理器
		go f.flushBuffer() // 启动定时刷新缓冲区
	})

	// 返回FastLog实例和nil错误
	return f, nil
}
// Close 关闭FastLog实例, 并等待所有日志处理完成。
func (f *FastLog) Close() error {
	f.closeLock.Lock()
	defer f.closeLock.Unlock()

	// 打印关闭日志记录器的信息
	f.Info("stop logging...")

	// 确保只关闭一次
	var closeOnce sync.Once
	var closeErr error
	closeOnce.Do(func() {
		// 关闭日志通道
		close(f.logChan)

		// 关闭协程
		f.cancel()

		// 等待所有日志处理完成
		f.logWait.Wait()

		// 刷新剩余的日志 缓冲区1
		f.flushBufferNow()

		// 刷新剩余的日志 缓冲区2
		f.flushBufferNow()

		// 如果不是仅输出到控制台, 同时日志文件句柄不为nil, 则关闭日志文件。
		if !f.config.GetConsoleOnly() && f.logGer != nil {
			f.fileMu.Lock()
			defer f.fileMu.Unlock()
			if err := f.logGer.Close(); err != nil {
				closeErr = fmt.Errorf("关闭日志文件失败: %v", err)
			}
		}

	}) // 执行一次

	// 检查是否有错误发生
	if closeErr != nil {
		return closeErr
	}

	return nil
}
// debug 记录调试级别的日志，不支持占位符
func (l *FastLog) Debug(v ...any) {
	// 检查日志级别，如果小于等于 Debug 级别，则不记录日志。
	if DEBUG < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,        // 时间戳
		level:       DEBUG,            // 日志级别
		message:     fmt.Sprint(v...), // 日志消息
		fileName:    filename,         // 文件名
		funcName:    funcName,         // 函数名
		line:        line,             // 行号
		goroutineID: getGoroutineID(), // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}
// Debugf 记录调试级别的日志，支持占位符，格式化
func (l *FastLog) Debugf(format string, v ...any) {
	// 检查日志级别，如果小于等于 Debug 级别，则不记录日志。
	if DEBUG < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,                 // 时间戳
		level:       DEBUG,                     // 日志级别
		message:     fmt.Sprintf(format, v...), // 日志消息
		fileName:    filename,                  // 文件名
		funcName:    funcName,                  // 函数名
		line:        line,                      // 行号
		goroutineID: getGoroutineID(),          // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg
}
// Error 记录错误级别的日志，不支持占位符
func (l *FastLog) Error(v ...any) {
	// 检查日志级别，如果小于等于 Error 级别，则不记录日志。
	if ERROR < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,        // 时间戳
		level:       ERROR,            // 日志级别
		message:     fmt.Sprint(v...), // 日志消息
		fileName:    filename,         // 文件名
		funcName:    funcName,         // 函数名
		line:        line,             // 行号
		goroutineID: getGoroutineID(), // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}
// Errorf 记录错误级别的日志，支持占位符，格式化
func (l *FastLog) Errorf(format string, v ...any) {
	// 检查日志级别，如果小于等于 Error 级别，则不记录日志。
	if ERROR < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,                 // 时间戳
		level:       ERROR,                     // 日志级别
		message:     fmt.Sprintf(format, v...), // 日志消息
		fileName:    filename,                  // 文件名
		funcName:    funcName,                  // 函数名
		line:        line,                      // 行号
		goroutineID: getGoroutineID(),          // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}
// Info 记录信息级别的日志，不支持占位符
func (l *FastLog) Info(v ...any) {
	// 检查日志级别，如果小于等于 Info 级别，则不记录日志。
	if INFO < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,        // 时间戳
		level:       INFO,             // 日志级别
		message:     fmt.Sprint(v...), // 日志消息
		fileName:    filename,         // 文件名
		funcName:    funcName,         // 函数名
		line:        line,             // 行号
		goroutineID: getGoroutineID(), // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}
// Infof 记录信息级别的日志，支持占位符，格式化
func (l *FastLog) Infof(format string, v ...any) {
	// 检查日志级别，如果小于等于 Info 级别，则不记录日志。
	if INFO < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,                 // 时间戳
		level:       INFO,                      // 日志级别
		message:     fmt.Sprintf(format, v...), // 日志消息
		fileName:    filename,                  // 文件名
		funcName:    funcName,                  // 函数名
		line:        line,                      // 行号
		goroutineID: getGoroutineID(),          // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}
// Success 记录成功级别的日志，不支持占位符
func (l *FastLog) Success(v ...any) {
	// 检查日志级别，如果小于等于 Success 级别，则不记录日志。
	if SUCCESS < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,        // 时间戳
		level:       SUCCESS,          // 日志级别
		message:     fmt.Sprint(v...), // 日志消息
		fileName:    filename,         // 文件名
		funcName:    funcName,         // 函数名
		line:        line,             // 行号
		goroutineID: getGoroutineID(), // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}
// Successf 记录成功级别的日志，支持占位符，格式化
func (l *FastLog) Successf(format string, v ...any) {
	// 检查日志级别，如果小于等于 Success 级别，则不记录日志。
	if SUCCESS < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,                 // 时间戳
		level:       SUCCESS,                   // 日志级别
		message:     fmt.Sprintf(format, v...), // 日志消息
		fileName:    filename,                  // 文件名
		funcName:    funcName,                  // 函数名
		line:        line,                      // 行号
		goroutineID: getGoroutineID(),          // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}
// Warn 记录警告级别的日志，不支持占位符
func (l *FastLog) Warn(v ...any) {
	// 检查日志级别，如果小于等于 Warn 级别，则不记录日志。
	if WARN < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,        // 时间戳
		level:       WARN,             // 日志级别
		message:     fmt.Sprint(v...), // 日志消息
		fileName:    filename,         // 文件名
		funcName:    funcName,         // 函数名
		line:        line,             // 行号
		goroutineID: getGoroutineID(), // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}
// Warnf 记录警告级别的日志，支持占位符，格式化
func (l *FastLog) Warnf(format string, v ...any) {
	// 检查日志级别，如果小于等于 Warn 级别，则不记录日志。
	if WARN < l.config.GetLogLevel() {
		return
	}

	// 获取调用者的信息
	filename, funcName, line, ok := getCallerInfo(2)
	if !ok {
		filename = "unknown"
		funcName = "unknown"
		line = 0
	}

	// 获取本地时区的当前时间时间戳。
	timestamp := time.Unix(time.Now().Unix(), 0)

	// 设置结构体的属性值
	logMsg := &logMessage{
		timestamp:   timestamp,                 // 时间戳
		level:       WARN,                      // 日志级别
		message:     fmt.Sprintf(format, v...), // 日志消息
		fileName:    filename,                  // 文件名
		funcName:    funcName,                  // 函数名
		line:        line,                      // 行号
		goroutineID: getGoroutineID(),          // 协程ID
	}

	// 将日志消息发送到日志通道
	l.logChan <- logMsg

}
// 定义一个配置结构体，用于配置日志记录器
type FastLogConfig struct {
	logDirName     string        // 日志目录路径
	logFileName    string        // 日志文件名
	printToConsole bool          // 是否将日志输出到控制台
	consoleOnly    bool          // 是否仅输出到控制台
	flushInterval  time.Duration // 刷新间隔，单位为time.Duration
	logLevel       LogLevel      // 日志级别
	chanIntSize    int           // 通道大小 默认10000
	logFormat      LogFormatType // 日志格式选项
	maxBufferSize  int           // 最大缓冲区大小, 单位为MB, 默认1MB
	noColor        bool          // 是否禁用终端颜色
	noBold         bool          // 是否禁用终端字体加粗
	maxLogFileSize int           // 最大日志文件大小, 单位为MB, 默认5MB
	maxLogAge      int           // 最大日志文件保留天数, 默认为0, 表示不做限制
	maxLogBackups  int           // 最大日志文件保留数量, 默认为0, 表示不做限制
	isLocalTime    bool          // 是否使用本地时间 默认使用UTC时间
	enableCompress bool          // 是否启用日志文件压缩 默认不启用
	setMu          sync.Mutex    // 用于保护配置的锁
}

// NewFastLogConfig 创建一个新的FastLogConfig实例, 用于配置日志记录器。
// 参数:
//   - logDirName: 日志目录名称, 默认为"applogs"。
//   - logFileName: 日志文件名称, 默认为"app.log"。
//
// 返回值:
//   - *FastLogConfig: 一个指向FastLogConfig实例的指针。
func NewFastLogConfig(logDirName string, logFileName string) *FastLogConfig {
	// 如果日志目录名称为空, 则使用默认值"logs"。
	if logDirName == "" {
		logDirName = "logs"
	}

	// 如果日志文件名称为空, 则使用默认值"app.log"。
	if logFileName == "" {
		logFileName = "app.log"
	}

	// 返回一个新的FastLogConfig实例
	return &FastLogConfig{
		logDirName:     logDirName,             // 日志目录名称
		logFileName:    logFileName,            // 日志文件名称
		printToConsole: true,                   // 是否将日志输出到控制台
		consoleOnly:    false,                  // 是否仅输出到控制台
		logLevel:       INFO,                   // 日志级别 默认INFO
		chanIntSize:    10000,                  // 通道大小 增加到10000
		flushInterval:  500 * time.Millisecond, // 刷新间隔 缩短到500毫秒
		logFormat:      Detailed,               // 日志格式选项
		maxBufferSize:  1 * 1024 * 1024,        // 最大缓冲区大小 默认1MB, 单位为MB
		maxLogFileSize: 5,                      // 最大日志文件大小, 单位为MB, 默认5MB
		maxLogAge:      0,                      // 最大日志文件保留天数, 默认为0, 表示不做限制
		maxLogBackups:  0,                      // 最大日志文件保留数量, 默认为0, 表示不做限制
		isLocalTime:    false,                  // 是否使用本地时间 默认使用UTC时间
		enableCompress: false,                  // 是否启用日志文件压缩 默认不启用
		noColor:        false,                  // 是否禁用终端颜色
		noBold:         false,                  // 是否禁用终端字体加粗
	}
}
// GetChanIntSize 获取通道大小
func (c *FastLogConfig) GetChanIntSize() int {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.chanIntSize
}
// GetConsoleOnly 获取是否仅输出到控制台的状态
func (c *FastLogConfig) GetConsoleOnly() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.consoleOnly
}
// GetEnableCompress 获取是否启用日志文件压缩的状态
func (c *FastLogConfig) GetEnableCompress() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.enableCompress
}
// GetFlushInterval 获取刷新间隔
func (c *FastLogConfig) GetFlushInterval() time.Duration {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.flushInterval
}
// GetIsLocalTime 获取是否使用本地时间的状态
func (c *FastLogConfig) GetIsLocalTime() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.isLocalTime
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
// GetLogFormat 获取日志格式选项
func (c *FastLogConfig) GetLogFormat() LogFormatType {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.logFormat
}
// GetLogLevel 获取日志级别
func (c *FastLogConfig) GetLogLevel() LogLevel {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.logLevel
}
// GetMaxBufferSize 获取最大缓冲区大小(MB)
func (c *FastLogConfig) GetMaxBufferSize() int {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.maxBufferSize
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
// GetMaxLogFileSize 获取最大日志文件大小(MB)
func (c *FastLogConfig) GetMaxLogFileSize() int {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.maxLogFileSize
}
// GetNoBold 获取是否禁用终端字体加粗的状态
func (c *FastLogConfig) GetNoBold() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.noBold
}
// GetNoColor 获取是否禁用终端颜色的状态
func (c *FastLogConfig) GetNoColor() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.noColor
}
// GetPrintToConsole 获取是否将日志输出到控制台的状态
func (c *FastLogConfig) GetPrintToConsole() bool {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	return c.printToConsole
}
// SetChanIntSize 设置通道大小
func (c *FastLogConfig) SetChanIntSize(size int) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.chanIntSize = size
}
// SetConsoleOnly 设置是否仅输出到控制台
func (c *FastLogConfig) SetConsoleOnly(only bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.consoleOnly = only
}
// SetEnableCompress 设置是否启用日志文件压缩
func (c *FastLogConfig) SetEnableCompress(compress bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.enableCompress = compress
}
// SetFlushInterval 设置刷新间隔
func (c *FastLogConfig) SetFlushInterval(interval time.Duration) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.flushInterval = interval
}
// SetIsLocalTime 设置是否使用本地时间
func (c *FastLogConfig) SetIsLocalTime(local bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.isLocalTime = local
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
// SetLogFormat 设置日志格式选项
func (c *FastLogConfig) SetLogFormat(format LogFormatType) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.logFormat = format
}
// SetLogLevel 设置日志级别
func (c *FastLogConfig) SetLogLevel(level LogLevel) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.logLevel = level
}
// SetMaxBufferSize 设置最大缓冲区大小(MB)
func (c *FastLogConfig) SetMaxBufferSize(size int) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.maxBufferSize = size
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
// SetMaxLogFileSize 设置最大日志文件大小(MB)
func (c *FastLogConfig) SetMaxLogFileSize(size int) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.maxLogFileSize = size
}
// SetNoBold 设置是否禁用终端字体加粗
func (c *FastLogConfig) SetNoBold(noBold bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.noBold = noBold
}
// SetNoColor 设置是否禁用终端颜色
func (c *FastLogConfig) SetNoColor(noColor bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.noColor = noColor
}
// SetPrintToConsole 设置是否将日志输出到控制台
func (c *FastLogConfig) SetPrintToConsole(print bool) {
	c.setMu.Lock()
	defer c.setMu.Unlock()
	c.printToConsole = print
}
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

	// SetMaxBufferSize 设置最大缓冲区大小(MB)
	SetMaxBufferSize(size int)
	// GetMaxBufferSize 获取最大缓冲区大小(MB)
	GetMaxBufferSize() int

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

// 定义一个接口, 声明对外暴露的方法
type FastLogInterface interface {
	Info(v ...any)                    // 记录信息级别的日志，不支持占位符
	Warn(v ...any)                    // 记录警告级别的日志，不支持占位符
	Error(v ...any)                   // 记录错误级别的日志，不支持占位符
	Success(v ...any)                 // 记录成功级别的日志，不支持占位符
	Debug(v ...any)                   // 记录调试级别的日志，不支持占位符
	Close()                           // 关闭日志记录器
	Infof(format string, v ...any)    // 记录信息级别的日志，支持占位符，格式化
	Warnf(format string, v ...any)    // 记录警告级别的日志，支持占位符，格式化
	Errorf(format string, v ...any)   // 记录错误级别的日志，支持占位符，格式化
	Successf(format string, v ...any) // 记录成功级别的日志，支持占位符，格式化
	Debugf(format string, v ...any)   // 记录调试级别的日志，支持占位符，格式化
}

// 日志格式选项
type LogFormatType int

// 日志格式选项
const (
	Detailed LogFormatType = iota // 详细格式
	Bracket                       // 方括号格式
	Json                          // json格式
	Threaded                      // 协程格式
	Simple                        // 简约格式
	Custom                        // 自定义格式
)
// 日志级别枚举
type LogLevel int

// 定义日志级别
const (
	DEBUG   LogLevel = 10  // 调试级别
	INFO    LogLevel = 20  // 信息级别
	SUCCESS LogLevel = 30  // 成功级别
	WARN    LogLevel = 40  // 警告级别
	ERROR   LogLevel = 50  // 错误级别
	None    LogLevel = 999 // 无日志级别
)
// PathInfo 是一个结构体，用于封装路径的信息
type PathInfo struct {
	Path    string      // 路径
	Exists  bool        // 是否存在
	IsFile  bool        // 是否为文件
	IsDir   bool        // 是否为目录
	Size    int64       // 文件大小（字节）
	Mode    os.FileMode // 文件权限
	ModTime time.Time   // 文件修改时间
}


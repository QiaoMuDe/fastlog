package fastlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Level 表示日志级别
type Level int32

// 日志级别常量
const (
	DEBUG Level = iota + 1 // 调试级别 (1)
	INFO                   // 信息级别 (2)
	WARN                   // 警告级别 (3)
	ERROR                  // 错误级别 (4)
	FATAL                  // 致命级别 (5)
	PANIC                  // 恐慌级别 (6)
)

// 日志级别名称常量
const (
	LevelNameDebug = "DEBUG"
	LevelNameInfo  = "INFO"
	LevelNameWarn  = "WARN"
	LevelNameError = "ERROR"
	LevelNameFatal = "FATAL"
	LevelNamePanic = "PANIC"
)

// String 返回级别的字符串表示
func (l Level) String() string {
	switch l {
	case DEBUG:
		return LevelNameDebug
	case INFO:
		return LevelNameInfo
	case WARN:
		return LevelNameWarn
	case ERROR:
		return LevelNameError
	case FATAL:
		return LevelNameFatal
	case PANIC:
		return LevelNamePanic
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

// Enabled 检查是否启用该级别 (lvl >= l 时启用)
//
// 参数:
//   - lvl: 要检查的级别
//
// 返回:
//   - bool: 是否启用该级别
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

// ParseLevel 从字符串解析日志级别
//
// 参数:
//   - s: 要解析的字符串
//
// 返回:
//   - Level: 解析后的日志级别
//   - error: 如果解析失败
func ParseLevel(s string) (Level, error) {
	switch strings.ToUpper(s) {
	case LevelNameDebug:
		return DEBUG, nil
	case LevelNameInfo:
		return INFO, nil
	case LevelNameWarn:
		return WARN, nil
	case LevelNameError:
		return ERROR, nil
	case LevelNameFatal:
		return FATAL, nil
	case LevelNamePanic:
		return PANIC, nil
	default:
		return INFO, fmt.Errorf("unknown level: %s", s)
	}
}

// AllLevels 返回所有日志级别
//
// 返回:
//   - []Level: 包含所有日志级别的切片
func AllLevels() []Level {
	return []Level{DEBUG, INFO, WARN, ERROR, FATAL, PANIC}
}

// Entry 表示一条日志记录
type Entry struct {
	Time       time.Time // 时间戳
	Level      Level     // 日志级别
	Message    string    // 日志消息
	Caller     string    // 调用者信息: file.go:func:line
	Fields     []Field   // 键值对字段
	TimeFormat string    // 时间格式, 从 Config.TimeFormat 传递
}

// callerSkip 是 getCaller 的跳过层数常量
// 用于跳过日志库内部调用栈, 直接定位到用户的调用位置
const callerSkip = 3

// Logger 日志记录器
//
// Logger 是 FastLog 的核心日志记录器, 提供日志记录、级别控制、采样等功能。
// 支持 6 种日志级别: DEBUG, INFO, WARN, ERROR, FATAL, PANIC。
// 支持三种调用方式: 标准日志 (Info)、格式化日志 (Infof)、结构化日志 (Infow)。
//
// 必须通过 fastlog.New(cfg) 构造函数创建, 切勿直接声明空结构体使用。
// 直接声明 Logger{} 会导致内部 writer/sampler 等关键字段为 nil, 引发 panic。
//
// 使用示例:
//
//	cfg := fastlog.NewConfig("logs/app.log")
//	logger := fastlog.New(cfg)
//	defer func() { _ = logger.Close() }()
//	logger.Info("服务启动成功")
type Logger struct {
	config  *Config        // 日志配置
	writer  io.WriteCloser // 日志写入器
	sampler *Sampler       // 日志采样器, nil 表示不启用采样
	mu      sync.Mutex     // 日志记录器的互斥锁
	level   atomic.Int32   // 运行时日志级别, 支持动态调整, 初始化时从 config.Level 设置
}

// New 创建一个新的日志记录器
//
// 参数:
//   - cfg: 日志配置, 零值时使用默认配置
//
// 返回:
//   - *Logger: 新的日志记录器实例
//
// 注意: 如果配置验证失败, 会触发 panic 以便快速发现问题
func New(cfg *Config) *Logger {
	if cfg == nil {
		panic("config is nil")
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		panic(fmt.Sprintf("invalid config: %v", err))
	}

	// 克隆配置
	config := cfg.Clone()

	// 应用默认值
	if config.Level == 0 {
		config.Level = INFO
	}
	if config.Formatter == nil {
		config.Formatter = Def{}
	}
	if config.TimeFormat == "" {
		config.TimeFormat = DefaultTimeFormat
	}

	// 创建写入器
	writer := config.NewWriter()
	if writer == nil {
		// 如果未指定写入器, 则使用控制台写入器
		writer = &ConsoleWriter{w: os.Stdout}
	}

	// 创建采样器
	sampler := config.NewSampler()

	// 创建日志记录器实例
	l := &Logger{
		config:  config,         // 日志配置
		writer:  writer,         // 日志写入器
		sampler: sampler,        // 日志采样器
		level:   atomic.Int32{}, // 运行时日志级别, 初始化时从 config.Level 设置
	}

	// 以 Config.Level 作为运行时级别的初始值
	l.level.Store(int32(config.Level))

	return l
}

// SetLevel 运行时动态修改日志级别, 立即生效
//
// 参数:
//   - level: 新的日志级别
func (l *Logger) SetLevel(level Level) {
	l.level.Store(int32(level))
}

// Level 返回当前运行时日志级别
//
// 返回:
//   - Level: 当前日志级别
func (l *Logger) Level() Level {
	return Level(l.level.Load())
}

// log 记录日志的核心方法
//
// 参数:
//   - level: 日志级别
//   - msg: 日志消息
//   - fields: 日志字段
func (l *Logger) log(level Level, msg string, fields []Field) {
	// 检查日志级别是否启用, 如果未启用则直接返回
	if !Level(l.level.Load()).Enabled(level) {
		return
	}

	// 采样检查: 如果采样器存在且判定为抑制, 则直接丢弃
	if l.sampler != nil && !l.sampler.Allow(level, msg) {
		return
	}

	// 从对象池获取日志条目
	entry := GetEntry()
	defer PutEntry(entry)

	// 填充日志条目
	entry.Time = time.Now()                                     // 时间戳
	entry.Level = level                                         // 日志级别
	entry.Message = msg                                         // 日志消息
	entry.TimeFormat = l.config.TimeFormat                      // 时间格式
	entry.Fields = append(entry.Fields[:0], l.config.Fields...) // 添加配置中的字段
	entry.Fields = append(entry.Fields, fields...)              // 添加用户提供的字段

	// 记录调用者信息
	if l.config.Caller {
		entry.Caller = getCaller(callerSkip)
	}

	// 填充字段
	data, err := l.config.Formatter.Format(entry)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "format error: %v\n", err)
		return
	}

	// 写入日志
	l.mu.Lock()
	_, err = l.writer.Write(data)
	l.mu.Unlock()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "write error: %v\n", err)
	}
}

// Debug 记录调试日志
//
// 参数:
//   - msg: 日志消息
func (l *Logger) Debug(msg string) {
	l.log(DEBUG, msg, nil)
}

// Info 记录信息日志
//
// 参数:
//   - msg: 日志消息
func (l *Logger) Info(msg string) {
	l.log(INFO, msg, nil)
}

// Warn 记录警告日志
//
// 参数:
//   - msg: 日志消息
func (l *Logger) Warn(msg string) {
	l.log(WARN, msg, nil)
}

// Error 记录错误日志
//
// 参数:
//   - msg: 日志消息
func (l *Logger) Error(msg string) {
	l.log(ERROR, msg, nil)
}

// Fatal 记录致命日志并退出程序
//
// 参数:
//   - msg: 日志消息
func (l *Logger) Fatal(msg string) {
	l.log(FATAL, msg, nil)
	_ = l.Sync()
	os.Exit(1)
}

// Panic 记录恐慌日志并触发 panic
//
// 参数:
//   - msg: 日志消息
func (l *Logger) Panic(msg string) {
	l.log(PANIC, msg, nil)
	_ = l.Sync()
	panic(msg)
}

// Debugf 记录格式化的调试日志
//
// 参数:
//   - format: 格式化字符串
//   - args: 格式化参数
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, fmt.Sprintf(format, args...), nil)
}

// Infof 记录格式化的信息日志
//
// 参数:
//   - format: 格式化字符串
//   - args: 格式化参数
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(INFO, fmt.Sprintf(format, args...), nil)
}

// Warnf 记录格式化的警告日志
//
// 参数:
//   - format: 格式化字符串
//   - args: 格式化参数
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WARN, fmt.Sprintf(format, args...), nil)
}

// Errorf 记录格式化的错误日志
//
// 参数:
//   - format: 格式化字符串
//   - args: 格式化参数
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, fmt.Sprintf(format, args...), nil)
}

// Fatalf 记录格式化的致命日志并退出程序
//
// 参数:
//   - format: 格式化字符串
//   - args: 格式化参数
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(FATAL, fmt.Sprintf(format, args...), nil)
	_ = l.Sync()
	os.Exit(1)
}

// Panicf 记录格式化的恐慌日志并触发 panic
//
// 参数:
//   - format: 格式化字符串
//   - args: 格式化参数
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.log(PANIC, fmt.Sprintf(format, args...), nil)
	_ = l.Sync()
	panic(fmt.Sprintf(format, args...))
}

// Debugw 记录带字段的调试日志
//
// 参数:
//   - msg: 日志消息
//   - fields: 日志字段
func (l *Logger) Debugw(msg string, fields ...Field) {
	l.log(DEBUG, msg, fields)
}

// Infow 记录带字段的信息日志
//
// 参数:
//   - msg: 日志消息
//   - fields: 日志字段
func (l *Logger) Infow(msg string, fields ...Field) {
	l.log(INFO, msg, fields)
}

// Warnw 记录带字段的警告日志
//
// 参数:
//   - msg: 日志消息
//   - fields: 日志字段
func (l *Logger) Warnw(msg string, fields ...Field) {
	l.log(WARN, msg, fields)
}

// Errorw 记录带字段的错误日志
//
// 参数:
//   - msg: 日志消息
//   - fields: 日志字段
func (l *Logger) Errorw(msg string, fields ...Field) {
	l.log(ERROR, msg, fields)
}

// Fatalw 记录带字段的致命日志并退出程序
//
// 参数:
//   - msg: 日志消息
//   - fields: 日志字段
func (l *Logger) Fatalw(msg string, fields ...Field) {
	l.log(FATAL, msg, fields)
	_ = l.Sync()
	os.Exit(1)
}

// Panicw 记录带字段的恐慌日志并触发 panic
//
// 参数:
//   - msg: 日志消息
//   - fields: 日志字段
func (l *Logger) Panicw(msg string, fields ...Field) {
	l.log(PANIC, msg, fields)
	_ = l.Sync()
	panic(msg)
}

// Sync 同步日志到存储
//
// 返回:
//   - error: 同步过程中的错误, 如果写入器不支持同步则返回 nil
func (l *Logger) Sync() error {
	if syncer, ok := l.writer.(interface{ Sync() error }); ok {
		return syncer.Sync()
	}
	return nil
}

// Close 关闭日志记录器
//
// 返回:
//   - error: 关闭过程中的错误
func (l *Logger) Close() error {
	return l.writer.Close()
}

// getCaller 获取调用者信息
//
// 参数:
//   - skip: 跳过调用栈的层数
//
// 返回:
//   - string: 调用者信息, 格式为 "文件名:函数名:行号"
//   - error: 如果获取调用者信息失败
func getCaller(skip int) string {
	// 获取调用栈信息: pc=程序计数器, file=文件名, line=行号
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "?:?:0"
	}

	// 取文件名最后一段, 如 "path/to/main.go" → "main.go"
	file = filepath.Base(file)

	// 通过 PC 获取函数名, 获取失败时用 "?" 保底
	fnName := "?"
	if fn := runtime.FuncForPC(pc); fn != nil {
		fnName = fn.Name()
		// 取完整函数名最后一个点之后的部分, 如 "main.main" → "main"
		if i := strings.LastIndexByte(fnName, '.'); i >= 0 {
			fnName = fnName[i+1:]
		}
	}

	// 格式化调用者信息
	return fmt.Sprintf("%s:%s:%d", file, fnName, line)
}

// EntryPool 日志条目池, 用于减少内存分配
var EntryPool = sync.Pool{
	New: func() interface{} {
		return &Entry{
			Fields: make([]Field, 0, 8),
		}
	},
}

// GetEntry 从池中获取日志条目
//
// 返回:
//   - *Entry: 日志条目实例
//   - error: 如果池为空, 返回错误
func GetEntry() *Entry {
	return EntryPool.Get().(*Entry)
}

// PutEntry 将日志条目放回池中
//
// 参数:
//   - e: 要放回池的日志条目
func PutEntry(e *Entry) {
	e.Fields = e.Fields[:0]
	e.Caller = ""
	e.Message = ""
	e.Time = time.Time{}
	EntryPool.Put(e)
}

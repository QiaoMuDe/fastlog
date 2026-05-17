package fastlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Logger 日志记录器
type Logger struct {
	config  *Config        // 日志配置
	writer  io.WriteCloser // 日志写入器
	sampler *Sampler       // 日志采样器, nil 表示不启用采样
	mu      sync.Mutex     // 日志记录器的互斥锁
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

	// 创建写入器
	writer := config.NewWriter()
	if writer == nil {
		// 如果未指定写入器, 则使用控制台写入器
		writer = &ConsoleWriter{w: os.Stdout}
	}

	// 创建采样器
	sampler := config.NewSampler()

	// 返回日志记录器实例
	return &Logger{
		config:  config,
		writer:  writer,
		sampler: sampler,
	}
}

// log 记录日志的核心方法
//
// 参数:
//   - level: 日志级别
//   - msg: 日志消息
//   - fields: 日志字段
func (l *Logger) log(level Level, msg string, fields []Field) {
	if !l.config.Level.Enabled(level) {
		return
	}

	// 采样检查: 如果采样器存在且判定为抑制, 则直接丢弃
	if l.sampler != nil && !l.sampler.Allow(level, msg) {
		return
	}

	// 从对象池获取日志条目
	entry := GetEntry()
	defer PutEntry(entry)

	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg
	entry.Fields = append(entry.Fields[:0], l.config.Fields...)
	entry.Fields = append(entry.Fields, fields...)

	if l.config.Caller {
		entry.Caller = getCaller(callerSkip)
	}

	data, err := l.config.Formatter.Format(entry)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "format error: %v\n", err)
		return
	}

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
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}

	// 获取函数名最后一个点之后的部分
	fnName := fn.Name()
	if i := strings.LastIndexByte(fnName, '.'); i >= 0 {
		fnName = fnName[i+1:]
	}

	// 获取文件名
	file = filepath.Base(file)

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

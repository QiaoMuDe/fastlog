package fastlog

import (
	"fmt"
	"strings"
	"time"
)

// Level 表示日志级别
type Level int8

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

// Formatter 定义日志格式化器接口
type Formatter interface {
	// Format 将日志条目格式化为字节数组
	Format(entry *Entry) ([]byte, error)
}

// Entry 表示一条日志记录
type Entry struct {
	Time    time.Time // 时间戳
	Level   Level     // 日志级别
	Message string    // 日志消息
	Caller  string    // 调用者信息: file.go:func:line
	Fields  []Field   // 键值对字段
}

// callerSkip 是 getCaller 的跳过层数常量
// 用于跳过日志库内部调用栈, 直接定位到用户的调用位置
const callerSkip = 3

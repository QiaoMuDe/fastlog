package fastlog

import (
	"bytes"
	"fmt"
	"time"

	"github.com/goccy/go-json"
)

// Formatter 定义日志格式化器接口
type Formatter interface {
	// Format 将日志条目格式化为字节数组
	Format(entry *Entry) ([]byte, error)
}

// Def 默认格式
// 格式: 2025-01-15T10:30:45 | INFO    | main.go:main:15 - 用户登录成功
type Def struct{}

// Format 实现默认格式化器
//
// 参数:
//   - entry: 日志条目
//
// 返回:
//   - []byte: 格式化后的字节数组
//   - error: 如果格式化失败
func (f Def) Format(entry *Entry) ([]byte, error) {
	var buf bytes.Buffer

	// 时间戳
	buf.WriteString(entry.Time.Format(time.RFC3339))
	buf.WriteString(" | ")

	// 级别 (左对齐, 6字符宽度)
	_, _ = fmt.Fprintf(&buf, "%-6s", entry.Level.String())
	buf.WriteString(" | ")

	// 调用者信息
	if entry.Caller != "" {
		buf.WriteString(entry.Caller)
		buf.WriteString(" - ")
	}

	// 消息
	buf.WriteString(entry.Message)

	// 字段
	if len(entry.Fields) > 0 {
		buf.WriteByte(' ')
		for i, field := range entry.Fields {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(field.Format())
		}
	}

	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

// JSON JSON 格式
type JSON struct{}

// Format 实现 JSON 格式
//
// 参数:
//   - entry: 日志条目
//
// 返回:
//   - []byte: 格式化后的字节数组
//   - error: 如果格式化失败
func (f JSON) Format(entry *Entry) ([]byte, error) {
	// 预分配容量, 避免 rehash
	cap := 4 + len(entry.Fields) // time + level + message + caller + fields
	data := make(map[string]interface{}, cap)

	// 添加基础字段
	data["time"] = entry.Time.Format(time.RFC3339)
	data["level"] = entry.Level.String()
	data["message"] = entry.Message

	// 添加调用者信息
	if entry.Caller != "" {
		data["caller"] = entry.Caller
	}

	// 添加字段
	for _, field := range entry.Fields {
		data[field.Key()] = field.toInterface()
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// 添加换行符
	b = append(b, '\n')
	return b, nil
}

// Simple 简单格式
// 格式: 2025-01-15T10:30:45 INFO  用户登录成功
type Simple struct{}

// Format 实现简单格式
//
// 参数:
//   - entry: 日志条目
//
// 返回:
//   - []byte: 格式化后的字节数组
//   - error: 如果格式化失败
func (f Simple) Format(entry *Entry) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(entry.Time.Format(time.RFC3339))
	buf.WriteByte(' ')
	buf.WriteString(entry.Level.String())
	buf.WriteByte(' ')
	buf.WriteString(entry.Message)

	if len(entry.Fields) > 0 {
		buf.WriteByte(' ')
		for i, field := range entry.Fields {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(field.Format())
		}
	}

	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

// KV 键值对格式
// 格式: time=2025-01-15T10:30:45 level=INFO message=用户登录成功
type KV struct{}

// Format 实现键值对格式
//
// 参数:
//   - entry: 日志条目
//
// 返回:
//   - []byte: 格式化后的字节数组
//   - error: 如果格式化失败
func (f KV) Format(entry *Entry) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("time=")
	buf.WriteString(entry.Time.Format(time.RFC3339))
	buf.WriteString(" level=")
	buf.WriteString(entry.Level.String())
	buf.WriteString(" message=")
	buf.WriteString(entry.Message)

	if entry.Caller != "" {
		buf.WriteString(" caller=")
		buf.WriteString(entry.Caller)
	}

	for _, field := range entry.Fields {
		buf.WriteByte(' ')
		buf.WriteString(field.Format())
	}

	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

// Compact 极简格式
// 格式: [I] 2025-01-15 10:30:45 用户登录成功
// 特点: 级别首字母 + 完整时间，简洁易读，适合容器环境
type Compact struct{}

// Format 实现 Compact 格式
//
// 参数:
//   - entry: 日志条目
//
// 返回:
//   - []byte: 格式化后的字节数组
//   - error: 如果格式化失败
func (f Compact) Format(entry *Entry) ([]byte, error) {
	var buf bytes.Buffer

	// 级别首字母
	buf.WriteByte('[')
	buf.WriteByte(entry.Level.String()[0])
	buf.WriteString("] ")

	// 时间（年月日时分秒）
	buf.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
	buf.WriteByte(' ')

	// 消息
	buf.WriteString(entry.Message)

	// 字段（简化为 key=value 形式）
	if len(entry.Fields) > 0 {
		buf.WriteString(" | ")
		for i, field := range entry.Fields {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteString(field.Format())
		}
	}

	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

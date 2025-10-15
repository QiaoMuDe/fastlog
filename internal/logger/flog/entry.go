package flog

import (
	"bytes"
	"sync"

	"gitee.com/MM-Q/fastlog/internal/config"
	"gitee.com/MM-Q/fastlog/internal/types"
	"gitee.com/MM-Q/go-kit/pool"
)

// entryPool Entry对象池，用于重用Entry实例
var entryPool = sync.Pool{
	New: func() interface{} {
		return &Entry{
			fields: make([]Field, 0, 8), // 预分配字段切片容量
		}
	},
}

// Entry 定义了日志条目的结构体
type Entry struct {
	time   string         // 日志记录时间
	level  types.LogLevel // 日志级别
	msg    string         // 日志消息
	caller []byte         // 调用者信息
	fields []Field        // 日志字段
}

// NewEntry 创建一个新的日志条目（使用对象池优化）
//
// 参数：
//   - needFileInfo: 是否需要文件信息。
//   - level: 日志级别。
//   - msg: 日志消息。
//   - *fields: 日志字段
//
// 返回值：
//   - *Entry: 一个指向 Entry 实例的指针。
func NewEntry(needFileInfo bool, level types.LogLevel, msg string, fields ...*Field) *Entry {
	// 从对象池获取Entry实例
	e := getEntry()

	// 设置日志条目基本信息
	e.level = level
	e.msg = msg

	// 复制字段切片（避免外部修改影响池化实例）
	if len(fields) > 0 {
		for _, field := range fields {
			if field != nil {
				e.fields = append(e.fields, *field)
			}
		}
	}

	// 获取调用时间
	e.time = types.GetCachedTimestamp()

	// 仅当需要文件信息时才获取调用者信息
	if needFileInfo {
		e.caller = types.GetCallerInfo(3)
	} else {
		e.caller = []byte("")
	}

	// 返回日志条目指针
	return e
}

// buildLog 构建日志消息
//
// 参数：
//   - cfg: 日志配置。
//   - e: 日志条目。
//
// 返回值：
//   - []byte: 构建的日志消息。
func buildLog(cfg *config.FastLogConfig, e *Entry) []byte {
	if cfg == nil || e == nil {
		return []byte{}
	}

	// 根据是否需要文件信息，选择不同的Json格式
	if cfg.CallerInfo {
		return pool.WithBuf(func(b *bytes.Buffer) {
			b.Write([]byte(`{"time":"`))
			b.WriteString(e.time)
			b.Write([]byte(`","level":"`))
			b.WriteString(types.LogLevelToString(e.level))
			b.Write([]byte(`","caller":"`))
			b.Write(e.caller)
			b.Write([]byte(`","msg":"`))
			b.WriteString(e.msg)
			b.Write([]byte(`"`))

			// 添加字段
			if len(e.fields) > 0 {
				b.Write([]byte(`,`))
				for i, field := range e.fields {
					if i > 0 {
						b.Write([]byte(`,`))
					}
					b.Write([]byte(`"`))
					b.WriteString(field.Key())
					b.Write([]byte(`":"`))
					b.WriteString(field.Value())
					b.Write([]byte(`"`))
				}
				b.Write([]byte(`}`))
			} else {
				b.Write([]byte(`}`))
			}
		})
	} else {
		return pool.WithBuf(func(b *bytes.Buffer) {
			b.Write([]byte(`{"time":"`))
			b.WriteString(e.time)
			b.Write([]byte(`","level":"`))
			b.WriteString(types.LogLevelToString(e.level))
			b.Write([]byte(`","msg":"`))
			b.WriteString(e.msg)
			b.Write([]byte(`"`))

			// 添加字段
			if len(e.fields) > 0 {
				b.Write([]byte(`,`))
				for i, field := range e.fields {
					if i > 0 {
						b.Write([]byte(`,`))
					}
					b.Write([]byte(`"`))
					b.WriteString(field.Key())
					b.Write([]byte(`":"`))
					b.WriteString(field.Value())
					b.Write([]byte(`"`))
				}
				b.Write([]byte(`}`))
			} else {
				b.Write([]byte(`}`))
			}
		})
	}
}

// getEntry 从对象池获取Entry实例
//
// 返回值：
//   - *Entry: 一个重置过的Entry实例
func getEntry() *Entry {
	return entryPool.Get().(*Entry)
}

// putEntry 将Entry实例归还到对象池
//
// 参数：
//   - e: 要归还的Entry实例
func putEntry(e *Entry) {
	if e == nil {
		return
	}

	// 清空字段，避免内存泄漏
	e.time = ""
	e.level = 0
	e.msg = ""
	e.caller = []byte("")
	e.fields = e.fields[:0] // 清空切片但保留容量
	entryPool.Put(e)
}

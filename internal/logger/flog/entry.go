package flog

import (
	"strconv"

	"gitee.com/MM-Q/fastlog/internal/types"
)

// Entry 定义了日志条目的结构体
type Entry struct {
	time   string         // 日志记录时间
	level  types.LogLevel // 日志级别
	msg    string         // 日志消息
	caller string         // 调用者信息
	fields []Field        // 日志字段
}

// NewEntry 创建一个新的日志条目
//
// 参数：
//   - level: 日志级别。
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
//
// 返回值：
//   - *Entry: 一个指向 Entry 实例的指针。
func NewEntry(level types.LogLevel, msg string, fields ...Field) *Entry {
	// 初始化日志条目
	e := &Entry{level: level, msg: msg, fields: fields}

	// 获取调用时间
	e.time = types.GetCachedTimestamp()

	// 获取调用者信息
	fileName, functionName, line, ok := types.GetCallerInfo(2)
	if ok {
		e.caller = fileName + ":" + functionName + ":" + strconv.Itoa(int(line))
	} else {
		e.caller = "???"
	}

	// 返回日志条目指针
	return e
}

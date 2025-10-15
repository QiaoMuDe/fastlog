package flog

import (
	"fmt"

	"gitee.com/MM-Q/fastlog/internal/types"
)

// handleLog 处理日志记录
//
// 参数：
//   - level: 日志级别。
//   - msg: 日志消息。
//   - fields: 日志字段，可变参数。
func (f *Flog) handleLog(level types.LogLevel, msg string, fields ...*Field) {
	if f != nil && f.cfg != nil {
		return
	}

	// 检查日志处理器是否已关闭
	if f.closed.Load() {
		return
	}

	// 检查日志级别，如果调用的日志级别低于配置的日志级别，则直接返回
	if level < f.cfg.LogLevel {
		return
	}

	// 创建日志条目
	e := NewEntry(f.cfg.CallerInfo, level, msg, fields...)
	defer putEntry(e) // 确保在函数返回前归还Entry实例到对象池

	// 构建日志条目
	log := buildLog(f.cfg, e)

	// 写入到终端
	if f.cfg.OutputToConsole {
		switch level {
		case types.INFO:
			f.cl.Blue(string(log))
		case types.WARN:
			f.cl.Yellow(string(log))
		case types.ERROR:
			f.cl.Red(string(log))
		case types.DEBUG:
			f.cl.Magenta(string(log))
		case types.FATAL:
			f.cl.Red(string(log))
		default:
			fmt.Println(string(log)) // 默认打印
		}
	}

	// 写入到文件
	log = append(log, '\n')
	if _, err := f.fileWriter.Write(log); err != nil {
		fmt.Printf("fastlog: failed to write log: %v\n", err)
	}
}

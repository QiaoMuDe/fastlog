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
	if f == nil || f.cfg == nil {
		return
	}

	// 检查日志处理器是否已关闭
	if f.closed.Load() {
		return
	}

	// 检查日志级别，使用位运算判断是否应该记录该级别的日志
	if !types.ShouldLog(level, f.cfg.LogLevel) {
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
		case types.INFO_Mask:
			f.cl.Blue(string(log))
		case types.WARN_Mask:
			f.cl.Yellow(string(log))
		case types.ERROR_Mask:
			f.cl.Red(string(log))
		case types.DEBUG_Mask:
			f.cl.Magenta(string(log))
		case types.FATAL_Mask:
			f.cl.Red(string(log))
		default:
			// 对于未知级别，使用默认颜色输出
			f.cl.White(string(log))
		}
	}

	// 写入到文件
	if f.cfg.OutputToFile && f.fileWriter != nil {
		// 确保日志以换行符结尾
		if len(log) == 0 || log[len(log)-1] != '\n' {
			log = append(log, '\n')
		}
		if _, err := f.fileWriter.Write(log); err != nil {
			fmt.Printf("fastlog: failed to write log: %v\n", err)
		}
	}
}

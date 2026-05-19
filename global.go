package fastlog

import "sync"

var (
	globalLogger     *Logger   // 全局默认日志记录器
	globalLoggerOnce sync.Once // 全局日志记录器初始化一次
)

// L 返回全局默认日志记录器
//
// 全局日志记录器在第一次调用时创建, 使用 Console() 配置（DEBUG 级别, 纯控制台输出）。
// 适合快速使用和调试, 无需手动创建 Logger 实例:
//
//	fastlog.L().Info("服务启动")
//	fastlog.L().Errorw("连接失败", fastlog.Err(err))
func L() *Logger {
	globalLoggerOnce.Do(func() {
		globalLogger = New(Console())
	})
	return globalLogger
}

// Close 关闭全局默认日志记录器
//
// 如果全局日志记录器未初始化（未调用 L()）, 返回 nil。
func Close() error {
	if globalLogger != nil {
		return globalLogger.Close()
	}
	return nil
}

// Sync 同步全局默认日志记录器的缓冲区数据到存储
//
// 如果全局日志记录器未初始化（未调用 L()）或写入器不支持同步, 返回 nil。
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

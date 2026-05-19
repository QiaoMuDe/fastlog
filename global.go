package fastlog

import "sync"

var (
	globalLogger     *Logger   // 全局默认日志记录器
	globalLoggerOnce sync.Once // 全局日志记录器初始化一次
)

// L 返回全局默认日志记录器
//
// 全局日志记录器默认输出到控制台 (DEBUG 级别), 适合快速使用和调试。
// 第一次调用 L() 时才会初始化, 在此之前调用 SetDefault 替换的实例不受影响。
//
// 示例:
//
//	fastlog.L().Info("服务启动")
//	fastlog.L().Errorw("连接失败", fastlog.Err(err))
func L() *Logger {
	globalLoggerOnce.Do(func() {
		if globalLogger == nil {
			globalLogger = New(Console())
		}
	})
	return globalLogger
}

// SetDefault 设置全局默认日志记录器
//
// 参数:
//   - l: 日志记录器实例
//
// SetDefault 在首次调用 L() 之前设置, 后续 L() 返回该实例;
// SetDefault 在首次调用 L() 之后设置, 后续 L() 同样返回新实例。
//
// 示例:
//
//	fastlog.SetDefault(fastlog.New(fastlog.Prod("logs/app.log")))
//	defer fastlog.Close()
func SetDefault(l *Logger) {
	globalLogger = l
}

// Close 关闭全局默认日志记录器
//
// 如果全局日志记录器未初始化或已通过 SetDefault 设置了自定义实例, 同样生效。
// 关闭后继续通过 L() 记录日志仍可工作, 但写入会失败（默认控制台输出不受影响）。
//
// 示例:
//
//	fastlog.SetDefault(fastlog.New(fastlog.Prod("logs/app.log")))
//	defer fastlog.Close()
func Close() error {
	if globalLogger != nil {
		return globalLogger.Close()
	}
	return nil
}

// Sync 同步全局默认日志记录器的缓冲区数据到存储
//
// 在程序退出前调用, 确保所有日志都已写入磁盘, 避免日志丢失。
// 如果写入器不支持同步, 返回 nil。
//
// 示例:
//
//	fastlog.SetDefault(fastlog.New(fastlog.Prod("logs/app.log")))
//	defer fastlog.Sync()
//	defer fastlog.Close()
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

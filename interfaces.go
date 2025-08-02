/*
interfaces.go - 接口定义模块
定义处理器所需的最小依赖接口，用于打破循环依赖并提高代码的可测试性。
*/
package fastlog

import (
	"context"
	"io"
	"time"

	"gitee.com/MM-Q/colorlib"
)

// ProcessorDependencies 定义处理器所需的最小依赖接口
// 通过接口隔离原则，processor 只能访问必要的功能，避免持有完整的 FastLog 引用
type ProcessorDependencies interface {
	// GetConfig 获取日志配置
	GetConfig() *FastLogConfig

	// GetFileWriter 获取文件写入器
	GetFileWriter() io.Writer

	// GetConsoleWriter 获取控制台写入器
	GetConsoleWriter() io.Writer

	// GetColorLib 获取颜色库实例
	GetColorLib() *colorlib.ColorLib

	// GetContext 获取上下文，用于控制处理器生命周期
	GetContext() context.Context

	// GetLogChannel 获取日志消息通道
	GetLogChannel() <-chan *logMsg

	// NotifyProcessorDone 通知处理器完成工作
	NotifyProcessorDone()
}

// WriterPair 写入器对，用于批量传递写入器
type WriterPair struct {
	FileWriter    io.Writer
	ConsoleWriter io.Writer
}

// ProcessorConfig 处理器配置结构
type ProcessorConfig struct {
	BatchSize     int           // 批量处理大小
	FlushInterval time.Duration // 刷新间隔
}

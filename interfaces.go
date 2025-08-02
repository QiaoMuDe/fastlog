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

// processorDependencies 定义处理器所需的最小依赖接口
// 通过接口隔离原则，processor 只能访问必要的功能，避免持有完整的 FastLog 引用
type processorDependencies interface {
	// getConfig 获取日志配置
	getConfig() *FastLogConfig

	// getFileWriter 获取文件写入器
	getFileWriter() io.Writer

	// getConsoleWriter 获取控制台写入器
	getConsoleWriter() io.Writer

	// getColorLib 获取颜色库实例
	getColorLib() *colorlib.ColorLib

	// getContext 获取上下文，用于控制处理器生命周期
	getContext() context.Context

	// getLogChannel 获取日志消息通道
	getLogChannel() <-chan *logMsg

	// notifyProcessorDone 通知处理器完成工作
	notifyProcessorDone()
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

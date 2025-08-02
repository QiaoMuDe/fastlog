/*
processor.go - 单线程日志处理器实现
负责从日志通道接收消息、批量缓存，并根据批次大小或时间间隔触发处理，
实现日志的批量格式化和输出。
*/
package fastlog

import (
	"bytes"
	"runtime/debug"
	"time"
)

// processor 单线程日志处理器
type processor struct {
	// 依赖接口 (替代直接持有FastLog引用)
	deps ProcessorDependencies

	// 单一缓冲区 (单线程使用，无需锁)
	fileBuffer    *bytes.Buffer // 文件缓冲区
	consoleBuffer *bytes.Buffer // 控制台缓冲区

	// 批量处理配置
	batchSize     int           // 批量处理数量
	flushInterval time.Duration // 批量处理间隔
}

// newProcessor 创建新的处理器实例
// 使用依赖注入模式，避免循环依赖
//
// 参数:
//   - deps: 依赖接口 (替代直接持有FastLog引用)
//   - batchSize: 批处理条数
//   - flushInterval: 定时刷新间隔
//
// 返回:
//   - *processor: 新的处理器实例
func newProcessor(deps ProcessorDependencies, batchSize int, flushInterval time.Duration) *processor {
	return &processor{
		deps:          deps,            // 依赖接口 (替代直接持有FastLog引用)
		fileBuffer:    &bytes.Buffer{}, // 文件缓冲区
		consoleBuffer: &bytes.Buffer{}, // 控制台缓冲区
		batchSize:     batchSize,       // 批处理条数
		flushInterval: flushInterval,   // 定时刷新间隔
	}
}

// singleThreadProcessor 单线程日志处理器
// 负责从日志通道接收消息、批量缓存，并根据批次大小或时间间隔触发处理
func (p *processor) singleThreadProcessor() {
	// 添加初始化检查
	if p == nil {
		panic("processor is nil")
	}
	if p.deps == nil {
		panic("processor.deps is nil")
	}
	if p.deps.GetConfig() == nil {
		panic("processor.deps.GetConfig() is nil")
	}
	if p.fileBuffer == nil {
		panic("processor.fileBuffer is nil")
	}
	if p.consoleBuffer == nil {
		panic("processor.consoleBuffer is nil")
	}
	// 检查通道是否为nil
	if p.deps.GetLogChannel() == nil {
		panic("processor.deps.GetLogChannel() is nil")
	}

	// 初始化日志批处理缓冲区，预分配容量以减少内存分配, 容量为配置的批处理大小batchSize
	batch := make([]*logMsg, 0, p.batchSize)

	// 创建定时刷新器，间隔由flushInterval指定
	ticker := time.NewTicker(p.flushInterval)

	defer func() {
		// 捕获panic
		if r := recover(); r != nil {
			p.deps.GetColorLib().PrintErrf("日志处理器发生panic: %s\nstack: %s\n", r, debug.Stack())
		}

		// 减少等待组中的计数器。
		p.deps.NotifyProcessorDone()
	}()

	// 主循环：持续处理日志消息和定时事件
	for {
		select {
		case logMsg := <-p.deps.GetLogChannel(): // 从日志通道接收新日志消息
			// 添加消息空值检查
			if logMsg == nil {
				continue // 跳过 nil 消息
			}

			// 将日志消息添加到批处理缓冲区
			batch = append(batch, logMsg)

			// 只在满足条件时才处理: 批处理切片写满或者缓冲区到达90%阈值
			shouldFlush := len(batch) >= p.batchSize || p.shouldFlushByThreshold()

			// 检查是否需要处理(满足条件之一)
			if shouldFlush {
				p.processAndFlushBatch(batch) // 处理并刷新批处理缓冲区
				batch = batch[:0]             // 重置批处理缓冲区，准备接收新消息
			}

		case <-ticker.C: // 定时刷新事件
			// 定时刷新：处理剩余消息并刷新缓冲区
			if len(batch) > 0 {
				p.processAndFlushBatch(batch) // 处理并刷新批处理缓冲区
				batch = batch[:0]             // 重置batch
			}

		case <-p.deps.GetContext().Done(): // 上下文取消信号，表示应停止处理
			// 关闭定时器
			ticker.Stop()

			// 处理剩余的batch(如果有的话)
			if len(batch) > 0 {
				p.processAndFlushBatch(batch) // 处理并刷新批处理缓冲区
			}

			return
		}
	}
}

// processAndFlushBatch 处理并刷新日志批处理缓冲区,
// 该函数负责格式化日志消息, 将它们写入相应的缓冲区(文件或控制台),
// 然后将缓冲区内容刷新到实际的输出目标(文件或控制台)。
//
// 参数:
// - batch []*logMsg: 日志批处理缓冲区，包含一批待处理的日志消息。
func (p *processor) processAndFlushBatch(batch []*logMsg) {
	// 重置缓冲区（清空原有内容，准备接收新数据）
	p.fileBuffer.Reset()    // 重置文件缓冲区
	p.consoleBuffer.Reset() // 重置控制台缓冲区

	// 获取配置
	config := p.deps.GetConfig()

	// 遍历批处理中的所有日志消息
	for _, logMsg := range batch {
		// 格式化日志消息，包括时间戳、级别、调用者信息等
		formattedMsg := p.formatLogWithDeps(logMsg)

		// 文件输出处理：如果启用了文件输出，则将格式化后的消息写入文件缓冲区
		if config.OutputToFile {
			p.fileBuffer.WriteString(formattedMsg + "\n")
		}

		// 控制台输出处理：如果启用了控制台输出，则将带颜色的消息写入控制台缓冲区
		if config.OutputToConsole {
			// 先渲染颜色再写入缓冲区，以确保控制台输出具有颜色效果
			coloredMsg := p.addColorWithDeps(logMsg, formattedMsg)
			p.consoleBuffer.WriteString(coloredMsg + "\n")
		}
	}

	// 如果启用文件输出, 并且文件缓冲区有内容, 则将缓冲区内容写入文件
	if config.OutputToFile && p.fileBuffer.Len() > 0 {
		// 将文件缓冲区的内容一次性写入文件, 提高I/O效率
		if _, writeErr := p.deps.GetFileWriter().Write(p.fileBuffer.Bytes()); writeErr != nil {
			// 如果写入失败，记录错误信息和堆栈跟踪
			p.deps.GetColorLib().PrintErrf("写入文件失败: %s\nstack: %s\n", writeErr, debug.Stack())

			// 如果启用了控制台输出，将文件内容降级输出到控制台
			if config.OutputToConsole {
				if _, consoleErr := p.deps.GetConsoleWriter().Write(p.fileBuffer.Bytes()); consoleErr != nil {
					// 控制台输出失败时静默处理，避免影响程序运行
					// 只在调试模式下输出错误信息（如果有其他可用的错误输出渠道）
					_ = writeErr // 静默忽略控制台输出错误
				}
			}
		}
	}

	// 如果启用控制台输出, 并且控制台缓冲区有内容, 则将缓冲区内容写入控制台
	if config.OutputToConsole && p.consoleBuffer.Len() > 0 {
		// 将控制台缓冲区的内容一次性写入控制台, 提高I/O效率
		if _, writeErr := p.deps.GetConsoleWriter().Write(p.consoleBuffer.Bytes()); writeErr != nil {
			// 控制台输出失败时静默处理，避免影响程序运行
			// 只在调试模式下输出错误信息（如果有其他可用的错误输出渠道）
			_ = writeErr // 静默忽略控制台输出错误
		}
	}

	// 在这里批量回收所有对象
	for _, logMsg := range batch {
		putLogMsg(logMsg)
	}
}

// formatLogWithDeps 使用依赖接口格式化日志消息（优化版本，直接调用纯函数）
func (p *processor) formatLogWithDeps(logMsg *logMsg) string {
	// 直接调用纯函数版本，避免创建临时对象
	return formatLogMessage(p.deps.GetConfig(), logMsg)
}

// addColorWithDeps 使用依赖接口添加颜色（优化版本，直接调用纯函数）
func (p *processor) addColorWithDeps(logMsg *logMsg, formattedMsg string) string {
	// 直接调用纯函数版本，避免创建临时对象
	return addColorToMessage(p.deps.GetColorLib(), logMsg.Level, formattedMsg)
}

// shouldFlushByThreshold 检查是否应该根据缓冲区大小阈值进行刷新
// 当文件缓冲区或控制台缓冲区任一达到90%阈值时返回true
func (p *processor) shouldFlushByThreshold() bool {
	config := p.deps.GetConfig()

	// 检查文件缓冲区是否达到90%阈值
	if config.OutputToFile {
		if p.fileBuffer.Len() >= fileFlushThreshold {
			return true
		}
	}

	// 检查控制台缓冲区是否达到90%阈值
	if config.OutputToConsole {
		if p.consoleBuffer.Len() >= consoleFlushThreshold {
			return true
		}
	}

	return false
}

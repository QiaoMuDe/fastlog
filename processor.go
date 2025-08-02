/*
processor.go - 单线程日志处理器实现
负责从日志通道接收消息、批量缓存，并根据批次大小或时间间隔触发处理，
实现日志的批量格式化和输出。使用智能分层缓冲区池优化内存管理。
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
	deps processorDependencies

	// 智能分层缓冲区池 (替代固定缓冲区)
	bufferPool *smartTieredBufferPool

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
func newProcessor(deps processorDependencies, batchSize int, flushInterval time.Duration) *processor {
	return &processor{
		deps:          deps,                  // 依赖接口 (替代直接持有FastLog引用)
		bufferPool:    globalSmartBufferPool, // 智能分层缓冲区池
		batchSize:     batchSize,             // 批处理条数
		flushInterval: flushInterval,         // 定时刷新间隔
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
	if p.deps.getConfig() == nil {
		panic("processor.deps.getConfig() is nil")
	}
	if p.bufferPool == nil {
		panic("processor.bufferPool is nil")
	}
	// 检查通道是否为nil
	if p.deps.getLogChannel() == nil {
		panic("processor.deps.getLogChannel() is nil")
	}

	// 初始化日志批处理缓冲区，预分配容量以减少内存分配, 容量为配置的批处理大小batchSize
	batch := make([]*logMsg, 0, p.batchSize)

	// 创建定时刷新器，间隔由flushInterval指定
	ticker := time.NewTicker(p.flushInterval)

	defer func() {
		// 捕获panic
		if r := recover(); r != nil {
			p.deps.getColorLib().PrintErrf("日志处理器发生panic: %s\nstack: %s\n", r, debug.Stack())
		}

		// 减少等待组中的计数器。
		p.deps.notifyProcessorDone()
	}()

	// 主循环：持续处理日志消息和定时事件
	for {
		select {
		case logMsg := <-p.deps.getLogChannel(): // 从日志通道接收新日志消息
			// 添加消息空值检查
			if logMsg == nil {
				continue // 跳过 nil 消息
			}

			// 将日志消息添加到批处理缓冲区
			batch = append(batch, logMsg)

			// 只在满足条件时才处理: 批处理切片写满或者缓冲区到达90%阈值
			shouldFlush := len(batch) >= p.batchSize || p.shouldFlushByThreshold(batch)

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

		case <-p.deps.getContext().Done(): // 上下文取消信号，表示应停止处理
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

// processAndFlushBatch 处理并刷新日志批处理缓冲区（智能缓冲区优化版本）,
// 该函数负责直接将日志消息格式化到缓冲区, 避免创建中间字符串,
// 然后将缓冲区内容刷新到实际的输出目标(文件或控制台)。
//
// 参数:
// - batch []*logMsg: 日志批处理缓冲区，包含一批待处理的日志消息。
func (p *processor) processAndFlushBatch(batch []*logMsg) {
	// 🛡️ 使用defer确保对象一定会被回收
	defer func() {
		// 批量回收所有对象
		for _, logMsg := range batch {
			if logMsg != nil {
				putLogMsg(logMsg)
			}
		}

		// 如果发生panic，记录但不重新抛出
		if r := recover(); r != nil {
			p.deps.getColorLib().PrintErrf("批处理时发生panic: %v\n", r)
			// 不重新panic，保证处理器继续运行
		}
	}()

	// 完整的空指针检查
	if p == nil {
		return
	}
	if p.bufferPool == nil {
		return
	}
	if p.deps == nil {
		return
	}
	if len(batch) == 0 {
		return
	}

	// 获取配置并检查
	config := p.deps.getConfig()
	if config == nil {
		return
	}

	// 估算批次大小，用于选择合适的缓冲区
	estimatedSize := len(batch) * 200 // 假设每条日志平均200字节

	// 🎯 智能获取分层缓冲区
	var fileBuffer, consoleBuffer *bytes.Buffer

	if config.OutputToFile {
		// 获取文件缓冲区（大容量，32KB起步）
		fileBuffer = p.bufferPool.GetFileBuffer(estimatedSize)
		defer p.bufferPool.PutFileBuffer(fileBuffer)
	}

	if config.OutputToConsole {
		// 获取控制台缓冲区（小容量，8KB起步）
		consoleBuffer = p.bufferPool.GetConsoleBuffer(estimatedSize)
		defer p.bufferPool.PutConsoleBuffer(consoleBuffer)
	}

	// 遍历批处理中的所有日志消息（智能缓冲区优化版本）
	for _, logMsg := range batch {
		// 跳过空的日志消息
		if logMsg == nil {
			continue
		}

		// 估算单条日志大小
		singleLogSize := len(logMsg.Message) + 100 // 消息长度 + 格式化开销

		// 文件输出处理：智能缓冲区升级 + 直接格式化
		if config.OutputToFile && fileBuffer != nil {
			// 🚀 智能检查并升级缓冲区（32KB -> 256KB -> 1MB）
			fileBuffer = p.bufferPool.CheckAndUpgradeFileBuffer(fileBuffer, singleLogSize)
			formatLogDirectlyToBuffer(fileBuffer, config, logMsg, false, p.deps.getColorLib())
			fileBuffer.WriteByte('\n') // 添加换行符
		}

		// 控制台输出处理：智能缓冲区升级 + 直接格式化，带颜色处理
		if config.OutputToConsole && consoleBuffer != nil {
			// 🚀 智能检查并升级缓冲区（8KB -> 32KB -> 64KB）
			consoleBuffer = p.bufferPool.CheckAndUpgradeConsoleBuffer(consoleBuffer, singleLogSize)
			formatLogDirectlyToBuffer(consoleBuffer, config, logMsg, true, p.deps.getColorLib())
			consoleBuffer.WriteByte('\n') // 添加换行符
		}
	}

	// 如果启用文件输出, 并且文件缓冲区有内容, 则将缓冲区内容写入文件
	if config.OutputToFile && fileBuffer != nil && fileBuffer.Len() > 0 {
		// 将文件缓冲区的内容一次性写入文件, 提高I/O效率
		if _, writeErr := p.deps.getFileWriter().Write(fileBuffer.Bytes()); writeErr != nil {
			// 如果写入失败，记录错误信息和堆栈跟踪
			p.deps.getColorLib().PrintErrf("写入文件失败: %s\nstack: %s\n", writeErr, debug.Stack())

			// 如果启用了控制台输出，将文件内容降级输出到控制台
			if config.OutputToConsole && consoleBuffer != nil {
				if _, consoleErr := p.deps.getConsoleWriter().Write(fileBuffer.Bytes()); consoleErr != nil {
					// 控制台输出失败时静默处理，避免影响程序运行
					// 只在调试模式下输出错误信息（如果有其他可用的错误输出渠道）
					_ = writeErr // 静默忽略控制台输出错误
				}
			}
		}
	}

	// 如果启用控制台输出, 并且控制台缓冲区有内容, 则将缓冲区内容写入控制台
	if config.OutputToConsole && consoleBuffer != nil && consoleBuffer.Len() > 0 {
		// 将控制台缓冲区的内容一次性写入控制台, 提高I/O效率
		if _, writeErr := p.deps.getConsoleWriter().Write(consoleBuffer.Bytes()); writeErr != nil {
			// 控制台输出失败时静默处理，避免影响程序运行
			// 只在调试模式下输出错误信息（如果有其他可用的错误输出渠道）
			_ = writeErr // 静默忽略控制台输出错误
		}
	}
}

// shouldFlushByThreshold 检查是否应该根据缓冲区大小阈值进行刷新
// 智能版本：基于批次大小估算，而不是实际缓冲区大小
//
// 参数:
//   - batch: 当前批次的日志消息
//
// 返回值:
//   - bool: 是否应该刷新
func (p *processor) shouldFlushByThreshold(batch []*logMsg) bool {
	if len(batch) == 0 {
		return false
	}

	config := p.deps.getConfig()
	if config == nil {
		return false
	}

	// 估算当前批次的大小
	estimatedSize := len(batch) * 200 // 每条日志约200字节

	// 检查是否达到阈值
	if config.OutputToFile && estimatedSize >= fileSmallThreshold {
		return true
	}

	if config.OutputToConsole && estimatedSize >= consoleSmallThreshold {
		return true
	}

	return false
}

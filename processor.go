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
	// 日志记录器
	f *FastLog

	// 单一缓冲区（单线程使用，无需锁）
	fileBuffer    *bytes.Buffer // 文件缓冲区
	consoleBuffer *bytes.Buffer // 控制台缓冲区
	bufferSize    int           // 缓冲区大小

	// 批量处理配置
	batchSize     int           // 批量处理数量
	flushInterval time.Duration // 批量处理间隔
}

// singleThreadProcessor 单线程日志处理器
// 负责从日志通道接收消息、批量缓存，并根据批次大小或时间间隔触发处理
func (p *processor) singleThreadProcessor() {
	// 添加初始化检查
	if p == nil {
		panic("processor is nil")
	}
	if p.f == nil {
		panic("processor.f is nil")
	}
	if p.f.config == nil {
		panic("processor.f.config is nil")
	}
	if p.fileBuffer == nil {
		panic("processor.fileBuffer is nil")
	}
	if p.consoleBuffer == nil {
		panic("processor.consoleBuffer is nil")
	}
	// 检查通道是否为nil
	if p.f.logChan == nil {
		panic("processor.f.logChan is nil")
	}

	// 初始化日志批处理缓冲区，预分配容量以减少内存分配, 容量为配置的批处理大小batchSize
	batch := make([]*logMessage, 0, p.batchSize)

	// 创建定时刷新器，间隔由flushInterval指定
	ticker := time.NewTicker(p.flushInterval)

	defer func() {
		// 捕获panic
		if r := recover(); r != nil {
			p.f.cl.PrintErrf("日志处理器发生panic: %s\nstack: %s\n", r, debug.Stack())
		}

		// 减少等待组中的计数器。
		p.f.logWait.Done()
	}()

	// 主循环：持续处理日志消息和定时事件
	for {
		select {
		case logMsg := <-p.f.logChan: // 从日志通道接收新日志消息
			// 添加消息空值检查
			if logMsg == nil {
				continue // 跳过 nil 消息
			}

			// 多级背压处理: 根据通道使用率丢弃低级别日志消息
			if shouldDropLogByBackpressure(p.f.logChan, logMsg.level) {
				continue
			}

			// 将日志消息添加到批处理缓冲区
			batch = append(batch, logMsg)

			// 只在满足条件时才处理: 批处理切片写满或者缓冲区到达90%阈值
			shouldFlush := len(batch) >= p.batchSize || p.shouldFlushByThreshold()

			// 检查是否需要处理(满足条件之一)
			if shouldFlush {
				p.processBatch(batch) // 一次性处理整个批次
				p.flushBuffers()      // 刷新缓冲区到目标输出
				batch = batch[:0]     // 重置批处理缓冲区，准备接收新消息
			}

		case <-ticker.C: // 定时刷新事件
			// 定时刷新：处理剩余消息并刷新缓冲区
			if len(batch) > 0 {
				p.processBatch(batch) // 处理剩余的batch到缓冲区
				p.flushBuffers()      // 刷新缓冲区到输出
				batch = batch[:0]     // 重置batch
			} else {
				// 即使batch为空，也可能缓冲区有内容需要刷新
				p.flushBuffers()
			}

		case <-p.f.ctx.Done(): // 上下文取消信号，表示应停止处理
			// 关闭定时器
			ticker.Stop()

			// 处理剩余的batch（如果有的话）
			if len(batch) > 0 {
				p.processBatch(batch)
			}

			// 最后刷新所有缓冲区内容
			p.flushBuffers()

			return
		}
	}
}

// processBatch 批量处理日志消息
// 将一批日志消息格式化后写入对应的缓冲区(文件和控制台)
//
// 参数：
//   - batch: 待处理的日志消息切片
func (p *processor) processBatch(batch []*logMessage) {
	// 重置缓冲区（清空原有内容，准备接收新数据）
	p.fileBuffer.Reset()    // 重置文件缓冲区
	p.consoleBuffer.Reset() // 重置控制台缓冲区

	// 遍历批处理中的所有日志消息
	for _, logMsg := range batch {
		// 格式化日志消息（根据配置的格式选项）
		formattedMsg := formatLog(p.f, logMsg)

		// 文件输出处理
		if p.f.config.OutputToFile {
			// 将格式化后的消息写入文件缓冲区（附加换行符）
			p.fileBuffer.WriteString(formattedMsg + "\n")
		}

		// 控制台输出处理
		if p.f.config.OutputToConsole {
			// 添加终端颜色样式（根据日志级别）
			coloredMsg := addColor(p.f, logMsg, formattedMsg)

			// 将带颜色的消息写入控制台缓冲区（附加换行符）
			p.consoleBuffer.WriteString(coloredMsg + "\n")
		}
	}
}

// flushBuffers 将缓冲区内容刷新到输出目标
// 负责将文件缓冲区和控制台缓冲区的内容批量写入到对应的输出设备
//
// 注意：此方法不处理缓冲区重置操作，仅执行写入操作
func (p *processor) flushBuffers() {
	// 检查文件缓冲区是否有待写入内容
	if p.f.config.OutputToFile {
		if p.fileBuffer.Len() > 0 {
			// 将文件缓冲区内容写入到文件写入器
			// 使用底层字节数组避免额外内存分配
			if _, writeErr := p.f.fileWriter.Write(p.fileBuffer.Bytes()); writeErr != nil {
				p.f.cl.PrintErrf("写入文件失败: %s\nstack: %s\n", writeErr, debug.Stack())
			}
		}
	}

	// 检查控制台缓冲区是否有待写入内容
	if p.f.config.OutputToConsole {
		if p.consoleBuffer.Len() > 0 {
			// 将控制台缓冲区内容写入到控制台写入器
			// 包含已添加的颜色控制字符
			if _, writeErr := p.f.consoleWriter.Write(p.consoleBuffer.Bytes()); writeErr != nil {
				p.f.cl.PrintErrf("写入控制台失败: %s\nstack: %s\n", writeErr, debug.Stack())
			}
		}
	}
}

// shouldFlushByThreshold 检查是否应该根据缓冲区大小阈值进行刷新
// 当文件缓冲区或控制台缓冲区任一达到90%阈值时返回true
func (p *processor) shouldFlushByThreshold() bool {
	// 检查文件缓冲区是否达到90%阈值
	if p.f.config.OutputToFile {
		if p.fileBuffer.Len() >= flushThreshold {
			return true
		}
	}

	// 检查控制台缓冲区是否达到90%阈值
	if p.f.config.OutputToConsole {
		if p.consoleBuffer.Len() >= flushThreshold {
			return true
		}
	}

	return false
}

/*
processor.go - 单线程日志处理器实现
负责从日志通道接收消息、批量缓存，并根据批次大小或时间间隔触发处理，
实现日志的批量格式化和输出。
*/
package fastlog

import (
	"bytes"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

// 字符串构建器对象池，用于复用临时字符串构建器，减少内存分配
var stringBuilderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// getStringBuilder 从对象池获取字符串构建器，使用安全的类型断言
func getStringBuilder() *strings.Builder {
	// 方式1: 安全的类型断言 (推荐)
	if builder, ok := stringBuilderPool.Get().(*strings.Builder); ok {
		return builder
	}
	// 如果类型断言失败，创建新的构建器作为fallback
	return &strings.Builder{}
}

// putStringBuilder 将字符串构建器归还到对象池
func putStringBuilder(builder *strings.Builder) {
	if builder != nil {
		builder.Reset()                // 重置构建器内容
		stringBuilderPool.Put(builder) // 归还到对象池
	}
}

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
	// 完整的空指针检查
	if p == nil {
		return
	}
	if p.fileBuffer == nil || p.consoleBuffer == nil {
		return
	}
	if p.deps == nil {
		return
	}
	if len(batch) == 0 {
		return
	}

	// 重置缓冲区（清空原有内容，准备接收新数据）
	p.fileBuffer.Reset()    // 重置文件缓冲区
	p.consoleBuffer.Reset() // 重置控制台缓冲区

	// 获取配置并检查
	config := p.deps.GetConfig()
	if config == nil {
		return
	}

	// 从对象池获取字符串构建器，用于复用临时字符串构建
	builder := getStringBuilder()
	defer putStringBuilder(builder)

	// 遍历批处理中的所有日志消息
	for _, logMsg := range batch {
		// 跳过空的日志消息
		if logMsg == nil {
			continue
		}

		// 使用对象池中的构建器格式化日志消息，避免临时字符串分配
		builder.Reset() // 重置构建器，准备格式化新消息
		p.formatLogDirectlyToBuilder(builder, logMsg)

		// 检查格式化结果
		if builder.Len() == 0 {
			continue
		}

		// 获取格式化后的消息内容
		formattedMsg := builder.String()

		// 文件输出处理：如果启用了文件输出，则将格式化后的消息写入文件缓冲区
		if config.OutputToFile {
			p.fileBuffer.WriteString(formattedMsg)
			p.fileBuffer.WriteByte('\n') // 使用WriteByte避免字符串拼接
		}

		// 控制台输出处理：如果启用了控制台输出，则将带颜色的消息写入控制台缓冲区
		if config.OutputToConsole {
			// 重置构建器，准备添加颜色
			builder.Reset()
			p.addColorDirectlyToBuilder(builder, logMsg, formattedMsg)

			// 将带颜色的消息写入控制台缓冲区
			p.consoleBuffer.WriteString(builder.String())
			p.consoleBuffer.WriteByte('\n') // 使用WriteByte避免字符串拼接
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

// formatLogDirectlyToBuilder 直接将格式化的日志消息写入字符串构建器，避免临时字符串分配
func (p *processor) formatLogDirectlyToBuilder(builder *strings.Builder, logMsg *logMsg) {
	// 完整的空指针检查
	if p == nil || p.deps == nil || builder == nil || logMsg == nil {
		return
	}

	config := p.deps.GetConfig()
	if config == nil {
		return
	}

	// 直接调用纯函数版本，获取格式化字符串
	formattedMsg := formatLogMessage(config, logMsg)
	if formattedMsg != "" {
		builder.WriteString(formattedMsg)
	}
}

// addColorDirectlyToBuilder 直接将带颜色的日志消息写入字符串构建器，避免临时字符串分配
func (p *processor) addColorDirectlyToBuilder(builder *strings.Builder, logMsg *logMsg, formattedMsg string) {
	// 完整的空指针检查
	if p == nil || p.deps == nil || builder == nil || logMsg == nil || formattedMsg == "" {
		return
	}

	colorLib := p.deps.GetColorLib()
	if colorLib == nil {
		builder.WriteString(formattedMsg) // 如果没有颜色库，直接写入原始消息
		return
	}

	// 直接调用纯函数版本，获取带颜色的字符串
	coloredMsg := addColorToMessage(colorLib, logMsg.Level, formattedMsg)
	builder.WriteString(coloredMsg)
}

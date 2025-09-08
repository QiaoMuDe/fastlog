/*
processor.go - 单线程日志处理器实现
负责从日志通道接收消息、批量缓存，并根据批次大小或时间间隔触发处理，
*/
package fastlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/go-kit/pool"
)

// processor 单线程日志处理器
type processor struct {
	// 依赖接口 (替代直接持有FastLog引用)
	deps processorDependencies

	// 批量处理配置
	batchSize     int           // 批量处理数量
	flushInterval time.Duration // 批量处理间隔
	bufferSize    int           // 缓冲区大小
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
		deps:          deps,                           // 依赖接口 (替代直接持有FastLog引用)
		batchSize:     batchSize,                      // 批处理条数
		flushInterval: flushInterval,                  // 定时刷新间隔
		bufferSize:    calculateBufferSize(batchSize), // 缓冲区大小
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
	if p.deps.getLogChannel() == nil {
		panic("processor.deps.getLogChannel() is nil")
	}

	// 初始化日志批处理切片，预分配容量以减少内存分配, 容量为配置的批处理大小batchSize
	batch := make([]*logMsg, 0, p.batchSize)

	// 创建定时刷新器，间隔由flushInterval指定
	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop() // 确保定时器在函数退出时停止

	defer func() {
		// 处理剩余消息
		p.drainRemainingMessages(batch)

		// 捕获panic
		if r := recover(); r != nil {
			p.deps.getColorLib().PrintErrorf("Log handler panic: %s\nstack: %s\n", r, debug.Stack())
		}

		// 减少等待组中的计数器。
		p.deps.notifyProcessorDone()
	}()

	// 主循环：持续处理日志消息和定时事件
	for {
		select {
		// 从日志通道接收新日志消息
		case logMsg, ok := <-p.deps.getLogChannel():
			if !ok { // 检查日志通道是否关闭
				return
			}

			// 添加消息空值检查
			if logMsg == nil {
				continue // 跳过 nil 消息
			}

			// 将日志消息添加到批处理缓冲区
			batch = append(batch, logMsg)

			// 只在满足条件时才处理: 批处理切片写满
			if len(batch) >= p.batchSize {
				p.processAndFlushBatch(batch) // 处理并刷新批处理缓冲区
				batch = batch[:0]             // 重置批处理缓冲区，准备接收新消息
			}

		// 定时刷新事件
		case <-ticker.C:
			// 定时刷新: 处理剩余消息并刷新缓冲区
			if len(batch) > 0 {
				p.processAndFlushBatch(batch) // 处理并刷新批处理缓冲区
				batch = batch[:0]             // 重置batch
			}

		// 上下文取消信号，表示应停止处理
		case <-p.deps.getContext().Done():
			return
		}
	}
}

// drainRemainingMessages 用于在返回之前处理日志通道中剩余的日志消息
//
// 参数:
//   - batch: 待处理的日志消息批次
func (p *processor) drainRemainingMessages(batch []*logMsg) {
	// 先处理当前 batch
	if len(batch) > 0 {
		p.processAndFlushBatch(batch) // 处理并刷新批处理缓冲区
		batch = batch[:0]             // 重置batch
	}

	// 非阻塞地读取通道中的剩余消息
	for {
		select {
		case logMsg, ok := <-p.deps.getLogChannel():
			if !ok {
				return // 通道已关闭且清空
			}

			if logMsg != nil {
				// 添加到批处理切片
				batch = append(batch, logMsg)

				// 只在满足条件时才处理: 批处理切片写满
				if len(batch) >= p.batchSize {
					p.processAndFlushBatch(batch) // 处理并刷新批处理缓冲区
					batch = batch[:0]             // 重置批处理缓冲区，准备接收新消息
				}
			}

		default:
			// 通道中没有更多消息,处理最后的 batch
			if len(batch) > 0 {
				p.processAndFlushBatch(batch)
			}
			return
		}
	}
}

// processAndFlushBatch 处理并刷新日志批处理缓冲区(智能缓冲区优化版本),
// 该函数负责直接将日志消息格式化到缓冲区, 避免创建中间字符串,
// 然后将缓冲区内容刷新到实际的输出目标(文件或控制台)。
//
// 参数:
//   - batch []*logMsg: 日志批处理缓冲区，包含一批待处理的日志消息。
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
			p.deps.getColorLib().PrintErrorf("panic occurred during batch processing: %v\n", r)
			// Don't re-panic to keep processor running
		}
	}()

	// 完整的空指针检查
	if p == nil {
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

	// 根据配置获取文件和控制台缓冲区
	var fileBuffer, consoleBuffer *bytes.Buffer
	if config.OutputToFile {
		fileBuffer = pool.GetBufCap(p.bufferSize)
		defer pool.PutBuf(fileBuffer)
	}
	if config.OutputToConsole {
		consoleBuffer = pool.GetBufCap(p.bufferSize)
		defer pool.PutBuf(consoleBuffer)
	}

	// 遍历批处理中的所有日志消息（智能缓冲区优化版本）
	for _, logMsg := range batch {
		// 跳过空的日志消息
		if logMsg == nil {
			continue
		}

		// 直接格式化日志消息到缓冲区
		if config.OutputToFile && fileBuffer != nil {
			formatLogDirectlyToBuffer(fileBuffer, config, logMsg, false, p.deps.getColorLib())
			fileBuffer.WriteByte('\n') // 添加换行符
		}

		// 直接格式化日志消息到缓冲区
		if config.OutputToConsole && consoleBuffer != nil {
			formatLogDirectlyToBuffer(consoleBuffer, config, logMsg, true, p.deps.getColorLib())
			consoleBuffer.WriteByte('\n') // 添加换行符
		}
	}

	// 如果启用文件输出, 并且文件缓冲区有内容, 则将缓冲区内容写入文件
	if config.OutputToFile && fileBuffer != nil && fileBuffer.Len() > 0 {
		// 将文件缓冲区的内容一次性写入文件, 提高I/O效率
		if _, writeErr := p.deps.getFileWriter().Write(fileBuffer.Bytes()); writeErr != nil {
			// 如果写入失败，记录错误信息和堆栈跟踪
			p.deps.getColorLib().PrintErrorf("Failed to write file: %s\nstack: %s\n", writeErr, debug.Stack())

			// 如果启用了控制台输出，将文件内容降级输出到控制台
			if config.OutputToConsole && consoleBuffer != nil {
				_, _ = p.deps.getConsoleWriter().Write(fileBuffer.Bytes())
			}
		}
	}

	// 如果启用控制台输出, 并且控制台缓冲区有内容, 则将缓冲区内容写入控制台
	if config.OutputToConsole && consoleBuffer != nil && consoleBuffer.Len() > 0 {
		// 控制台输出失败时静默处理，避免影响程序运行
		_, _ = p.deps.getConsoleWriter().Write(consoleBuffer.Bytes())
	}
}

// formatLogDirectlyToBuffer 直接将日志消息格式化到缓冲区，避免创建中间字符串（零拷贝优化）
//
// 参数:
//   - buffer: 目标缓冲区
//   - config: 日志配置
//   - logMsg: 日志消息
//   - withColor: 是否添加颜色（用于控制台输出）
//   - colorLib: 颜色库实例（当withColor为true时使用）
func formatLogDirectlyToBuffer(buffer *bytes.Buffer, config *FastLogConfig, logMsg *logMsg, withColor bool, colorLib *colorlib.ColorLib) {
	// 检查参数有效性
	if buffer == nil || config == nil || logMsg == nil || colorLib == nil {
		return
	}

	// 如果时间戳为空，使用缓存的时间戳
	if logMsg.Timestamp == "" {
		logMsg.Timestamp = getCachedTimestamp()
	}

	// 检查关键字段是否为空，设置默认值
	if logMsg.Message == "" {
		return // 消息为空直接返回
	}
	if logMsg.FileName == "" {
		logMsg.FileName = "unknown-file"
	}
	if logMsg.FuncName == "" {
		logMsg.FuncName = "unknown-func"
	}

	// 文本格式处理: 先格式化到临时缓冲区，然后根据需要添加颜色
	tempBuffer := pool.GetBuf()
	defer pool.PutBuf(tempBuffer)

	// 根据日志格式格式化到临时缓冲区
	switch config.LogFormat {
	// JSON格式
	case Json:
		// 序列化为JSON并直接写入缓冲区
		if jsonBytes, err := json.Marshal(logMsg); err == nil {
			tempBuffer.Write(jsonBytes)
		} else {
			// JSON序列化失败时的降级处理
			fmt.Fprintf(tempBuffer,
				logFormatMap[Json],
				logMsg.Timestamp, logLevelToString(logMsg.Level), "unknown", "unknown", 0,
				fmt.Sprintf("Failed to serialize original message: %v | Original content: %s", err, logMsg.Message),
			)
		}

	// JsonSimple格式（无文件信息）
	case JsonSimple:
		// 序列化为JSON并直接写入缓冲区
		if jsonBytes, err := json.Marshal(simpleLogMsg{
			Timestamp: logMsg.Timestamp,
			Level:     logMsg.Level,
			Message:   logMsg.Message,
		}); err == nil {
			tempBuffer.Write(jsonBytes)
		} else {
			// JSON序列化失败时的降级处理
			fmt.Fprintf(tempBuffer, logFormatMap[JsonSimple],
				logMsg.Timestamp, logLevelToString(logMsg.Level), fmt.Sprintf("Failed to serialize: %v | Original: %s", err, logMsg.Message))
		}

	// 详细格式
	case Detailed:
		tempBuffer.WriteString(logMsg.Timestamp) // 时间戳
		tempBuffer.WriteString(" | ")
		levelStr := logLevelToPaddedString(logMsg.Level) // 使用预填充的日志级别字符串
		tempBuffer.WriteString(levelStr)
		tempBuffer.WriteString(" | ")
		tempBuffer.WriteString(logMsg.FileName) // 文件信息
		tempBuffer.WriteByte(':')
		tempBuffer.WriteString(logMsg.FuncName) // 函数
		tempBuffer.WriteByte(':')
		tempBuffer.WriteString(strconv.Itoa(int(logMsg.Line))) // 行号
		tempBuffer.WriteString(" - ")
		tempBuffer.WriteString(logMsg.Message) // 消息

	// 简约格式
	case Simple:
		tempBuffer.WriteString(logMsg.Timestamp) // 时间戳
		tempBuffer.WriteString(" | ")
		levelStr := logLevelToPaddedString(logMsg.Level) // 使用预填充的日志级别字符串
		tempBuffer.WriteString(levelStr)
		tempBuffer.WriteString(" | ")
		tempBuffer.WriteString(logMsg.Message) // 消息

	// 结构化格式
	case Structured:
		tempBuffer.WriteString("T:") // 时间戳
		tempBuffer.WriteString(logMsg.Timestamp)
		tempBuffer.WriteString("|L:")                    // 日志级别
		levelStr := logLevelToPaddedString(logMsg.Level) // 使用预填充的日志级别字符串
		tempBuffer.WriteString(levelStr)
		tempBuffer.WriteString("|F:") // 文件信息
		tempBuffer.WriteString(logMsg.FileName)
		tempBuffer.WriteByte(':')
		tempBuffer.WriteString(logMsg.FuncName)
		tempBuffer.WriteByte(':')
		tempBuffer.WriteString(strconv.Itoa(int(logMsg.Line)))
		tempBuffer.WriteString("|M:") // 消息
		tempBuffer.WriteString(logMsg.Message)

	// 基础结构化格式(无文件信息)
	case BasicStructured:
		tempBuffer.WriteString("T:") // 时间戳
		tempBuffer.WriteString(logMsg.Timestamp)
		tempBuffer.WriteString("|L:")                    // 日志级别
		levelStr := logLevelToPaddedString(logMsg.Level) // 使用预填充的日志级别字符串
		tempBuffer.WriteString(levelStr)
		tempBuffer.WriteString("|M:") // 消息
		tempBuffer.WriteString(logMsg.Message)

	// 简单时间格式
	case SimpleTimestamp:
		tempBuffer.WriteString(logMsg.Timestamp) // 时间戳
		tempBuffer.WriteString(" ")
		levelStr := logLevelToPaddedString(logMsg.Level) // 使用预填充的日志级别字符串
		tempBuffer.WriteString(levelStr)                 // 日志级别
		tempBuffer.WriteString(" ")
		tempBuffer.WriteString(logMsg.Message) // 消息

	// 自定义格式
	case Custom:
		tempBuffer.WriteString(logMsg.Message)

	// 默认情况
	default:
		tempBuffer.WriteString("Unrecognized log format option: ")
		fmt.Fprintf(tempBuffer, "%v", config.LogFormat)
	}

	// 根据withColor参数决定是否添加颜色
	if withColor {
		// 使用零拷贝版本：直接将带颜色的内容写入目标缓冲区(控制台)
		addColorToBuffer(buffer, colorLib, logMsg.Level, tempBuffer)
	} else {
		// 直接写入原始内容(文件)
		buffer.Write(tempBuffer.Bytes())
	}
}

// logLevelToPaddedString 将 LogLevel 转换为带填充的字符串（用于文本格式化）
//
// 参数：
//   - level: 要转换的日志级别
//
// 返回值：
//   - string: 对应的带填充的日志级别字符串（7个字符），如果 level 无效, 则返回 "UNKNOWN"
func logLevelToPaddedString(level LogLevel) string {
	// 使用预构建的带填充映射表进行O(1)查询（适用于文本格式）
	if str, exists := logLevelPaddedStringMap[level]; exists {
		return str
	}
	return "UNKNOWN"
}

// addColorToBuffer 直接将带颜色的消息写入缓冲区，避免创建中间字符串（零拷贝优化版本）
//
// 参数：
//   - buffer: 目标缓冲区
//   - cl: 颜色库实例
//   - level: 日志级别
//   - sourceBuffer: 源缓冲区（包含原始消息内容）
func addColorToBuffer(buffer *bytes.Buffer, cl *colorlib.ColorLib, level LogLevel, sourceBuffer *bytes.Buffer) {
	// 完整的空指针和参数检查
	if buffer == nil || cl == nil || sourceBuffer == nil {
		return
	}

	// 检查源缓冲区是否为空
	if sourceBuffer.Len() == 0 {
		return
	}

	// 获取源缓冲区的内容(避免String()调用的内存分配)
	sourceBytes := sourceBuffer.Bytes()
	sourceString := string(sourceBytes) // 这里仍需要一次转换，但比多次String()调用更高效

	// 根据日志级别添加颜色并直接写入目标缓冲区
	switch level {
	case INFO:
		buffer.WriteString(cl.Sblue(sourceString)) // Blue
	case WARN:
		buffer.WriteString(cl.Syellow(sourceString)) // Yellow
	case ERROR:
		buffer.WriteString(cl.Sred(sourceString)) // Red
	case DEBUG:
		buffer.WriteString(cl.Smagenta(sourceString)) // Magenta
	case FATAL:
		buffer.WriteString(cl.Sred(sourceString)) // Red
	default:
		// 如果没有匹配到日志级别，直接写入原始内容
		buffer.Write(sourceBytes)
	}
}

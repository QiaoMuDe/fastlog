package fastlog

// processLogs 日志处理器, 用于处理日志消息。
// 读取通道中的日志消息并将其处理为指定的日志格式, 然后写入到缓冲区中。
func (f *FastLog) processLogs() {
	f.logWait.Add(1) // 增加等待组中的计数器

	// 创建一个goroutine, 用于处理日志消息
	go func() {
		defer func() {
			// 减少等待组中的计数器。
			f.logWait.Done()

			// 捕获panic
			if r := recover(); r != nil {
				f.Errorf("日志处理器发生panic: %v", r)
			}
		}()

		// 初始化控制台字符串构建器
		for {
			select {
			// 监听通道, 如果通道关闭, 则退出循环
			case <-f.ctx.Done():
				// 检查通道是否为空
				if len(f.logChan) != 0 {
					// 处理通道中剩余的日志消息
					for rawMsg := range f.logChan {
						f.handleLog(rawMsg)
					}
				}
				return
			// 监听通道, 如果通道中有日志消息, 则处理日志消息
			case rawMsg, ok := <-f.logChan:
				if !ok {
					return
				}

				// 处理日志消息
				f.handleLog(rawMsg)
			}
		}
	}()
}

// handleLog 处理单条日志消息的逻辑
// 格式化日志消息, 然后写入到缓冲区中。
// 如果缓冲区大小达到90%, 则立即刷新缓冲区。

// 参数:
//   - rawMsg: 日志消息

// 返回值:
//   - 无
func (f *FastLog) handleLog(rawMsg *logMessage) {
	// 格式化日志消息
	formattedLog := formatLog(f, rawMsg)

	// 局部变量记录是否需要刷新
	var needFileFlush bool
	var needConsoleFlush bool

	// 处理文件缓冲区 - 在同一个锁内完成检查和写入
	if !f.config.GetConsoleOnly() {
		f.fileBufferMu.Lock()

		// 记录当前文件缓冲区的索引和大小
		currentFileBufferIdx := f.fileBufferIdx.Load()
		currentFileBufferSize := f.fileBuffers[currentFileBufferIdx].Len()

		// 检查是否需要刷新（写入前检查）
		needFileFlush = currentFileBufferSize >= flushThreshold

		// 写入文件缓冲区
		if f.fileBuffers[currentFileBufferIdx] != nil {
			f.fileBuffers[currentFileBufferIdx].WriteString(formattedLog + "\n")
		}
		f.fileBufferMu.Unlock()
	}

	// 处理控制台缓冲区 - 同样的逻辑
	if f.config.GetPrintToConsole() || f.config.GetConsoleOnly() {
		f.consoleBufferMu.Lock()

		// 记录当前的控制台缓冲区的索引和大小
		currentConsoleBufferIdx := f.consoleBufferIdx.Load()
		currentConsoleBufferSize := f.consoleBuffers[currentConsoleBufferIdx].Len()

		// 检查是否需要刷新（写入前检查）
		needConsoleFlush = currentConsoleBufferSize >= flushThreshold

		// 渲染颜色
		consoleLog := addColor(f, rawMsg, formattedLog)

		// 写入控制台缓冲区
		if f.consoleBuffers[currentConsoleBufferIdx] != nil {
			f.consoleBuffers[currentConsoleBufferIdx].WriteString(consoleLog + "\n")
		}
		f.consoleBufferMu.Unlock()
	}

	// 统一判断是否需要刷新（在所有写入完成后）
	if needFileFlush || needConsoleFlush {
		f.flushBufferNow()
	}
}

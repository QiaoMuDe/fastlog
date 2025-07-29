package fastlog

import "time"

// flushBuffer 定时刷新缓冲区
func (f *FastLog) flushBuffer() {
	// 新增一个等待组, 用于等待刷新缓冲区的协程完成
	f.logWait.Add(1)

	// 定义一个定时器, 用于定时刷新缓冲区
	ticker := time.NewTicker(f.config.GetFlushInterval())

	// 创建一个goroutine, 用于定时刷新缓冲区
	go func() {
		defer func() {
			// 捕获panic
			if r := recover(); r != nil {
				f.Errorf("刷新缓冲区发生panic: %v", r)
			}

			// 减少等待组中的计数器。
			f.logWait.Done()

			// 关闭定时器
			if ticker != nil {
				ticker.Stop()
			}
		}()

		// 循环监听定时器
		for {
			select {
			case <-f.ctx.Done():
				return // 当 context 被取消时, 退出协程
			case <-ticker.C:
				f.flushBufferNow() // 刷新缓冲区
			}
		}
	}()
}

// flushBufferNow 立即刷新缓冲区
// 刷新缓冲区, 并将控制台缓冲区的内容输出到控制台
// 将控制台缓冲区的内容输出到控制台后, 将控制台缓冲区清空。
func (f *FastLog) flushBufferNow() {
	f.flushLock.Lock() // 加锁, 防止并发刷新
	defer f.flushLock.Unlock()

	// 判断是否需要刷新文件缓冲区
	if !f.config.GetConsoleOnly() {
		f.fileBufferMu.Lock()
		defer f.fileBufferMu.Unlock()

		// 获取当前缓冲区索引
		currentIdx := f.fileBufferIdx.Load()

		// 检查当前缓冲区是否有内容
		if f.fileBuffers[currentIdx].Len() > 0 {
			// 切换到另一个缓冲区
			newIdx := 1 - currentIdx // 0 -> 1, 1 -> 0
			f.fileBufferIdx.Store(newIdx)

			// 将文件缓冲区中的内容写入文件
			f.fileMu.Lock()
			defer f.fileMu.Unlock()
			if _, err := f.fileWriter.Write(f.fileBuffers[currentIdx].Bytes()); err != nil {
				f.Errorf("写入文件失败: %v", err)
			}

			// 重置当前缓冲区
			f.fileBuffers[currentIdx].Reset()
		}
	}

	// 判断是否需要刷新控制台缓冲区
	if f.config.GetPrintToConsole() || f.config.GetConsoleOnly() {
		f.consoleBufferMu.Lock()
		defer f.consoleBufferMu.Unlock()

		// 获取当前缓冲区索引
		currentIdx := f.consoleBufferIdx.Load()

		// 检查当前缓冲区是否有内容
		if f.consoleBuffers[currentIdx].Len() > 0 {
			// 切换到另一个缓冲区
			newIdx := 1 - currentIdx // 0  -> 1, 1 -> 0
			f.consoleBufferIdx.Store(newIdx)

			// 将控制台缓冲区中的内容写入控制台
			f.consoleMu.Lock()
			defer f.consoleMu.Unlock()
			if _, err := f.consoleWriter.Write(f.consoleBuffers[currentIdx].Bytes()); err != nil {
				f.Errorf("写入控制台失败: %v", err)
			}

			// 重置当前缓冲区
			f.consoleBuffers[currentIdx].Reset()
		}
	}
}

// cleanupBuffers 清理缓冲区资源，释放内存
func (f *FastLog) cleanupBuffers() {
	// 清理文件缓冲区
	f.fileBufferMu.Lock()
	for i := range f.fileBuffers {
		if f.fileBuffers[i] != nil {
			f.fileBuffers[i].Reset()
			f.fileBuffers[i] = nil // 帮助GC回收
		}
	}
	f.fileBufferMu.Unlock()

	// 清理控制台缓冲区
	f.consoleBufferMu.Lock()
	for i := range f.consoleBuffers {
		if f.consoleBuffers[i] != nil {
			f.consoleBuffers[i].Reset()
			f.consoleBuffers[i] = nil // 帮助GC回收
		}
	}
	f.consoleBufferMu.Unlock()
}

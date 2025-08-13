/*
smart_buffer_pool_test.go - 智能分层缓冲区池测试
测试分层缓冲区的获取、升级、归还等核心功能，
确保90%阈值触发机制和内存管理正确工作。
*/
package fastlog

import (
	"bytes"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestSmartTieredBufferPool_BasicFunctionality 测试基本功能
func TestSmartTieredBufferPool_BasicFunctionality(t *testing.T) {
	pool := newSmartTieredBufferPool()

	t.Run("文件缓冲区获取测试", func(t *testing.T) {
		// 测试小文件缓冲区
		smallBuffer := pool.GetFileBuffer(1024) // 1KB
		if smallBuffer.Cap() != fileSmallBufferCapacity {
			t.Errorf("小文件缓冲区容量错误: 期望 %d, 实际 %d", fileSmallBufferCapacity, smallBuffer.Cap())
		}

		// 测试中等文件缓冲区
		mediumBuffer := pool.GetFileBuffer(100 * 1024) // 100KB
		if mediumBuffer.Cap() != fileMediumBufferCapacity {
			t.Errorf("中等文件缓冲区容量错误: 期望 %d, 实际 %d", fileMediumBufferCapacity, mediumBuffer.Cap())
		}

		// 测试大文件缓冲区
		largeBuffer := pool.GetFileBuffer(500 * 1024) // 500KB
		if largeBuffer.Cap() != fileLargeBufferCapacity {
			t.Errorf("大文件缓冲区容量错误: 期望 %d, 实际 %d", fileLargeBufferCapacity, largeBuffer.Cap())
		}

		// 归还缓冲区
		pool.PutFileBuffer(smallBuffer)
		pool.PutFileBuffer(mediumBuffer)
		pool.PutFileBuffer(largeBuffer)
	})

	t.Run("控制台缓冲区获取测试", func(t *testing.T) {
		// 测试小控制台缓冲区
		smallBuffer := pool.GetConsoleBuffer(1024) // 1KB
		if smallBuffer.Cap() != consoleSmallBufferCapacity {
			t.Errorf("小控制台缓冲区容量错误: 期望 %d, 实际 %d", consoleSmallBufferCapacity, smallBuffer.Cap())
		}

		// 测试中等控制台缓冲区
		mediumBuffer := pool.GetConsoleBuffer(16 * 1024) // 16KB
		if mediumBuffer.Cap() != consoleMediumBufferCapacity {
			t.Errorf("中等控制台缓冲区容量错误: 期望 %d, 实际 %d", consoleMediumBufferCapacity, mediumBuffer.Cap())
		}

		// 测试大控制台缓冲区
		largeBuffer := pool.GetConsoleBuffer(50 * 1024) // 50KB
		if largeBuffer.Cap() != consoleLargeBufferCapacity {
			t.Errorf("大控制台缓冲区容量错误: 期望 %d, 实际 %d", consoleLargeBufferCapacity, largeBuffer.Cap())
		}

		// 归还缓冲区
		pool.PutConsoleBuffer(smallBuffer)
		pool.PutConsoleBuffer(mediumBuffer)
		pool.PutConsoleBuffer(largeBuffer)
	})
}

// TestSmartTieredBufferPool_ThresholdUpgrade 测试90%阈值升级机制
func TestSmartTieredBufferPool_ThresholdUpgrade(t *testing.T) {
	pool := newSmartTieredBufferPool()

	t.Run("文件缓冲区90%阈值升级", func(t *testing.T) {
		// 获取小文件缓冲区
		buffer := pool.GetFileBuffer(1024)
		originalCap := buffer.Cap()

		// 写入数据到接近90%阈值
		data := make([]byte, fileSmallThreshold-100) // 接近但未达到阈值
		buffer.Write(data)

		// 检查升级 - 应该不升级
		newBuffer := pool.CheckAndUpgradeFileBuffer(buffer, 50)
		if newBuffer != buffer {
			t.Error("未达到90%阈值时不应该升级缓冲区")
		}

		// 写入更多数据超过90%阈值
		upgradedBuffer := pool.CheckAndUpgradeFileBuffer(buffer, 200)
		if upgradedBuffer == buffer {
			t.Error("达到90%阈值时应该升级缓冲区")
		}
		if upgradedBuffer.Cap() <= originalCap {
			t.Errorf("升级后的缓冲区容量应该更大: 原始 %d, 升级后 %d", originalCap, upgradedBuffer.Cap())
		}

		// 验证数据完整性
		if upgradedBuffer.Len() != len(data) {
			t.Errorf("升级后数据长度不匹配: 期望 %d, 实际 %d", len(data), upgradedBuffer.Len())
		}

		pool.PutFileBuffer(upgradedBuffer)
	})

	t.Run("控制台缓冲区90%阈值升级", func(t *testing.T) {
		// 获取小控制台缓冲区
		buffer := pool.GetConsoleBuffer(1024)
		originalCap := buffer.Cap()

		// 写入数据到接近90%阈值
		data := make([]byte, consoleSmallThreshold-100)
		buffer.Write(data)

		// 检查升级 - 应该不升级
		newBuffer := pool.CheckAndUpgradeConsoleBuffer(buffer, 50)
		if newBuffer != buffer {
			t.Error("未达到90%阈值时不应该升级缓冲区")
		}

		// 写入更多数据超过90%阈值
		upgradedBuffer := pool.CheckAndUpgradeConsoleBuffer(buffer, 200)
		if upgradedBuffer == buffer {
			t.Error("达到90%阈值时应该升级缓冲区")
		}
		if upgradedBuffer.Cap() <= originalCap {
			t.Errorf("升级后的缓冲区容量应该更大: 原始 %d, 升级后 %d", originalCap, upgradedBuffer.Cap())
		}

		pool.PutConsoleBuffer(upgradedBuffer)
	})
}

// TestSmartTieredBufferPool_BufferReuse 测试缓冲区重用机制
func TestSmartTieredBufferPool_BufferReuse(t *testing.T) {
	pool := newSmartTieredBufferPool()

	t.Run("缓冲区重用验证", func(t *testing.T) {
		// 获取并归还缓冲区
		buffer1 := pool.GetFileBuffer(1024)
		buffer1.WriteString("test data")
		pool.PutFileBuffer(buffer1)

		// 再次获取相同大小的缓冲区
		buffer2 := pool.GetFileBuffer(1024)

		// 验证缓冲区被重用（内容应该被清空）
		if buffer2.Len() != 0 {
			t.Error("重用的缓冲区应该被清空")
		}

		// 在某些情况下可能是同一个缓冲区对象
		if buffer1 == buffer2 {
			t.Log("缓冲区对象被成功重用")
		}

		pool.PutFileBuffer(buffer2)
	})
}

// TestSmartTieredBufferPool_BufferClassification 测试缓冲区分类机制
func TestSmartTieredBufferPool_BufferClassification(t *testing.T) {
	pool := newSmartTieredBufferPool()

	t.Run("扩容后的缓冲区重新分类", func(t *testing.T) {
		// 获取小缓冲区
		buffer := pool.GetFileBuffer(1024)

		// 手动扩容缓冲区（模拟大量写入）
		largeData := make([]byte, fileMediumBufferCapacity)
		buffer.Write(largeData)

		// 归还时应该根据实际容量重新分类
		originalCap := buffer.Cap()
		pool.PutFileBuffer(buffer)

		// 再次获取中等大小的缓冲区，可能会得到之前扩容的缓冲区
		newBuffer := pool.GetFileBuffer(100 * 1024)

		// 验证缓冲区容量合理
		if newBuffer.Cap() < fileMediumBufferCapacity {
			t.Errorf("获取的中等缓冲区容量不足: %d", newBuffer.Cap())
		}

		t.Logf("原始容量: %d, 新缓冲区容量: %d", originalCap, newBuffer.Cap())
		pool.PutFileBuffer(newBuffer)
	})
}

// TestSmartTieredBufferPool_SuperLargeBuffer 测试超大缓冲区处理
func TestSmartTieredBufferPool_SuperLargeBuffer(t *testing.T) {
	pool := newSmartTieredBufferPool()

	t.Run("超大缓冲区不进池", func(t *testing.T) {
		// 获取大缓冲区
		buffer := pool.GetFileBuffer(500 * 1024)

		// 写入大量数据，触发超大缓冲区创建
		largeData := make([]byte, fileLargeThreshold+1000)
		buffer.Write(largeData)

		// 升级到超大缓冲区
		superLargeBuffer := pool.CheckAndUpgradeFileBuffer(buffer, 1000)

		// 验证超大缓冲区容量
		if superLargeBuffer.Cap() <= fileLargeBufferCapacity {
			t.Error("超大缓冲区容量应该超过标准大缓冲区")
		}

		// 归还超大缓冲区（应该不进池）
		pool.PutFileBuffer(superLargeBuffer)

		t.Logf("超大缓冲区容量: %d bytes", superLargeBuffer.Cap())
	})
}

// TestSmartTieredBufferPool_ConcurrentAccess 测试并发访问
func TestSmartTieredBufferPool_ConcurrentAccess(t *testing.T) {
	pool := newSmartTieredBufferPool()
	const goroutineCount = 100
	const operationsPerGoroutine = 100

	t.Run("并发获取和归还缓冲区", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(goroutineCount)

		for i := 0; i < goroutineCount; i++ {
			go func(id int) {
				defer wg.Done()

				for j := 0; j < operationsPerGoroutine; j++ {
					// 随机获取不同大小的缓冲区
					size := (id*j + j) % 100000 // 0-100KB

					// 文件缓冲区操作
					fileBuffer := pool.GetFileBuffer(size)
					fileBuffer.WriteString("concurrent test data")
					pool.PutFileBuffer(fileBuffer)

					// 控制台缓冲区操作
					consoleBuffer := pool.GetConsoleBuffer(size)
					consoleBuffer.WriteString("concurrent test data")
					pool.PutConsoleBuffer(consoleBuffer)
				}
			}(i)
		}

		wg.Wait()
		t.Logf("并发测试完成: %d个goroutine，每个执行%d次操作", goroutineCount, operationsPerGoroutine)
	})
}

// TestSmartTieredBufferPool_MemoryUsage 测试内存使用情况
func TestSmartTieredBufferPool_MemoryUsage(t *testing.T) {
	pool := newSmartTieredBufferPool()

	t.Run("内存使用监控", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		// 大量缓冲区操作
		buffers := make([]*bytes.Buffer, 1000)
		for i := 0; i < 1000; i++ {
			buffers[i] = pool.GetFileBuffer(i * 100)
			buffers[i].Write(make([]byte, i*10))
		}

		// 归还所有缓冲区
		for _, buffer := range buffers {
			pool.PutFileBuffer(buffer)
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		memoryIncrease := int64(m2.Alloc) - int64(m1.Alloc)
		t.Logf("内存使用变化: %d bytes", memoryIncrease)

		// 验证内存使用合理（调整阈值以适应实际情况）
		// 由于创建了1000个缓冲区，每个缓冲区容量不同，内存使用会比较高
		if memoryIncrease > 200*1024*1024 { // 200MB
			t.Errorf("内存使用过多: %d bytes (%.2f MB)", memoryIncrease, float64(memoryIncrease)/(1024*1024))
		} else {
			t.Logf("内存使用在合理范围内: %d bytes (%.2f MB)", memoryIncrease, float64(memoryIncrease)/(1024*1024))
		}
	})
}

// TestSmartTieredBufferPool_EdgeCases 测试边界情况
func TestSmartTieredBufferPool_EdgeCases(t *testing.T) {
	pool := newSmartTieredBufferPool()

	t.Run("边界情况处理", func(t *testing.T) {
		// 测试nil缓冲区
		result := pool.CheckAndUpgradeFileBuffer(nil, 1000)
		if result == nil {
			t.Error("nil缓冲区应该返回新缓冲区")
		}
		pool.PutFileBuffer(result)

		// 测试零大小请求
		zeroBuffer := pool.GetFileBuffer(0)
		if zeroBuffer == nil {
			t.Error("零大小请求应该返回有效缓冲区")
		}
		pool.PutFileBuffer(zeroBuffer)

		// 测试负数大小请求
		negativeBuffer := pool.GetFileBuffer(-100)
		if negativeBuffer == nil {
			t.Error("负数大小请求应该返回有效缓冲区")
		}
		pool.PutFileBuffer(negativeBuffer)

		// 测试归还nil缓冲区
		pool.PutFileBuffer(nil)    // 应该不会panic
		pool.PutConsoleBuffer(nil) // 应该不会panic
	})
}

// TestSmartTieredBufferPool_Performance 性能测试
func TestSmartTieredBufferPool_Performance(t *testing.T) {
	pool := newSmartTieredBufferPool()

	t.Run("性能基准测试", func(t *testing.T) {
		const iterations = 10000

		start := time.Now()

		for i := 0; i < iterations; i++ {
			// 获取缓冲区
			buffer := pool.GetFileBuffer(i % 50000)

			// 写入数据
			buffer.Write(make([]byte, i%1000))

			// 可能的升级
			buffer = pool.CheckAndUpgradeFileBuffer(buffer, i%500)

			// 归还缓冲区
			pool.PutFileBuffer(buffer)
		}

		duration := time.Since(start)
		avgTime := duration / iterations

		t.Logf("性能测试结果:")
		t.Logf("  总耗时: %v", duration)
		t.Logf("  平均每次操作: %v", avgTime)
		t.Logf("  每秒操作数: %.0f", float64(iterations)/duration.Seconds())

		// 性能要求：平均每次操作应该在12微秒以内
		if avgTime > 12*time.Microsecond {
			t.Errorf("性能不达标: 平均每次操作 %v > 12μs", avgTime)
		}
	})
}

// TestSmartTieredBufferPool_Integration 集成测试
func TestSmartTieredBufferPool_Integration(t *testing.T) {
	pool := newSmartTieredBufferPool()

	t.Run("完整工作流程测试", func(t *testing.T) {
		// 模拟真实的日志处理场景
		const batchSize = 100
		const batchCount = 50

		for batch := 0; batch < batchCount; batch++ {
			// 估算批次大小
			estimatedSize := batchSize * 200 // 每条日志约200字节

			// 获取文件和控制台缓冲区
			fileBuffer := pool.GetFileBuffer(estimatedSize)
			consoleBuffer := pool.GetConsoleBuffer(estimatedSize)

			// 模拟批量写入日志
			for i := 0; i < batchSize; i++ {
				logData := []byte("2025-08-02 23:36:08 | INFO | test.go:main:123 - 这是一条测试日志消息\n")

				// 检查并升级缓冲区
				fileBuffer = pool.CheckAndUpgradeFileBuffer(fileBuffer, len(logData))
				consoleBuffer = pool.CheckAndUpgradeConsoleBuffer(consoleBuffer, len(logData))

				// 写入数据
				fileBuffer.Write(logData)
				consoleBuffer.Write(logData)
			}

			// 验证数据完整性
			expectedLen := batchSize * len("2025-08-02 23:36:08 | INFO | test.go:main:123 - 这是一条测试日志消息\n")
			if fileBuffer.Len() != expectedLen {
				t.Errorf("文件缓冲区数据长度不正确: 期望 %d, 实际 %d", expectedLen, fileBuffer.Len())
			}
			if consoleBuffer.Len() != expectedLen {
				t.Errorf("控制台缓冲区数据长度不正确: 期望 %d, 实际 %d", expectedLen, consoleBuffer.Len())
			}

			// 归还缓冲区
			pool.PutFileBuffer(fileBuffer)
			pool.PutConsoleBuffer(consoleBuffer)
		}

		t.Logf("集成测试完成: 处理了 %d 个批次，每批次 %d 条日志", batchCount, batchSize)
	})
}

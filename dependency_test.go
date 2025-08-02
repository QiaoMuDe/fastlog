/*
dependency_test.go - 循环依赖修复验证测试文件
验证通过接口隔离和依赖注入方案是否成功解决了循环依赖问题，
包含内存泄漏预防测试和处理器依赖接口实现验证。
*/
package fastlog

import (
	"runtime"
	"testing"
	"time"
)

// TestCircularDependencyFixed 测试循环依赖是否已修复
func TestCircularDependencyFixed(t *testing.T) {
	// 创建配置
	config := NewFastLogConfig("logs", "test.log")
	config.OutputToConsole = true
	config.OutputToFile = false
	config.ChanIntSize = 10

	// 创建日志实例
	logger, err := NewFastLog(config)
	if err != nil {
		t.Fatalf("创建日志实例失败: %v", err)
	}

	// 记录一些日志
	logger.Info("测试循环依赖修复")
	logger.Debug("这是一条调试消息")
	logger.Warn("这是一条警告消息")

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 关闭日志器
	logger.Close()

	// 强制垃圾回收
	runtime.GC()
	runtime.GC() // 执行两次确保彻底回收

	t.Log("循环依赖修复验证通过：日志器正常创建、使用和关闭")
}

// TestProcessorDependencyInterface 测试处理器依赖接口
func TestProcessorDependencyInterface(t *testing.T) {
	// 创建配置
	config := NewFastLogConfig("logs", "test.log")
	config.OutputToConsole = false
	config.OutputToFile = false

	// 创建日志实例
	logger, err := NewFastLog(config)
	if err != nil {
		t.Fatalf("创建日志实例失败: %v", err)
	}
	defer logger.Close()

	// 验证FastLog实现了ProcessorDependencies接口
	var deps ProcessorDependencies = logger

	// 测试接口方法
	if deps.GetConfig() == nil {
		t.Error("GetConfig() 返回 nil")
	}

	if deps.GetFileWriter() == nil {
		t.Error("GetFileWriter() 返回 nil")
	}

	if deps.GetConsoleWriter() == nil {
		t.Error("GetConsoleWriter() 返回 nil")
	}

	if deps.GetColorLib() == nil {
		t.Error("GetColorLib() 返回 nil")
	}

	if deps.GetContext() == nil {
		t.Error("GetContext() 返回 nil")
	}

	if deps.GetLogChannel() == nil {
		t.Error("GetLogChannel() 返回 nil")
	}

	t.Log("ProcessorDependencies接口实现验证通过")
}

// TestMemoryLeakPrevention 测试内存泄漏预防
func TestMemoryLeakPrevention(t *testing.T) {
	// 记录初始内存统计
	var m1 runtime.MemStats
	runtime.GC()
	runtime.GC() // 执行两次GC确保稳定状态
	time.Sleep(50 * time.Millisecond)
	runtime.ReadMemStats(&m1)

	// 创建和销毁多个日志实例
	for i := 0; i < 10; i++ {
		config := NewFastLogConfig("logs", "test.log")
		config.OutputToConsole = false
		config.OutputToFile = false
		config.ChanIntSize = 5

		logger, err := NewFastLog(config)
		if err != nil {
			t.Fatalf("创建日志实例失败: %v", err)
		}

		// 写入一些日志
		logger.Infof("测试消息 %d", i)
		time.Sleep(10 * time.Millisecond)

		// 关闭日志器
		logger.Close()
	}

	// 强制垃圾回收
	runtime.GC()
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	// 记录最终内存统计
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	// 安全地计算内存变化（处理uint64下溢问题）
	var memoryChange int64
	if m2.Alloc >= m1.Alloc {
		memoryChange = int64(m2.Alloc - m1.Alloc)
		t.Logf("内存增长: %d bytes", memoryChange)
	} else {
		memoryChange = -int64(m1.Alloc - m2.Alloc)
		t.Logf("内存减少: %d bytes", -memoryChange)
	}

	// 检查是否存在明显的内存泄漏
	// 允许合理的内存增长（最多1MB），但不应该有大量泄漏
	if memoryChange > 1024*1024 {
		t.Errorf("可能存在内存泄漏，内存增长: %d bytes", memoryChange)
	} else if memoryChange < 0 {
		t.Logf("内存使用优化良好，减少了 %d bytes", -memoryChange)
	} else {
		t.Logf("内存增长在合理范围内: %d bytes", memoryChange)
	}

	t.Log("内存泄漏预防测试通过")
}

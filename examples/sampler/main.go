package main

import (
	"time"

	fastlog "gitee.com/MM-Q/fastlog"
)

func main() {
	example1()

	example2()
}

// example1 采样器基础示例
func example1() {
	println("=== 采样器基础示例 ===")

	// 每 3 秒一个窗口, 前 2 条放行, 之后每 3 条放行 1 条
	logger := fastlog.New(&fastlog.Config{
		Level:             fastlog.INFO,
		OutputConsole:     true,
		SamplerTick:       3 * time.Second,
		SamplerInitial:    2,
		SamplerThereafter: 3,
	})
	defer func() { _ = logger.Close() }()

	// 模拟高并发重复日志
	for i := 0; i < 10; i++ {
		logger.Errorw("数据库连接超时", fastlog.String("db", "mysql"))
	}

	// 不同消息不受采样影响
	logger.Errorw("磁盘空间不足", fastlog.String("disk", "/dev/sda1"))
}

// example2 采样器窗口重置示例
func example2() {
	println("\n=== 采样器窗口重置示例 ===")

	logger := fastlog.New(&fastlog.Config{
		Level:             fastlog.INFO,
		OutputConsole:     true,
		SamplerTick:       1 * time.Second,
		SamplerInitial:    1,
		SamplerThereafter: 0,
	})
	defer func() { _ = logger.Close() }()

	// 第一轮: 1 秒内只放行第 1 条
	logger.Warn("缓存命中率低")
	logger.Warn("缓存命中率低")
	logger.Warn("缓存命中率低")

	// 等待窗口过期后重新计数
	time.Sleep(1100 * time.Millisecond)

	// 第二轮: 窗口重置, 第 1 条又会被放行
	logger.Warn("缓存命中率低")
	logger.Warn("缓存命中率低")
}

package main

import (
	fastlog "gitee.com/MM-Q/fastlog"
)

func main() {
	// example2 先跑, 因为 example1 的 Fatal 会退出程序
	example2()

	example1()
}

// example1 彩色日志示例
func example1() {
	println("=== 彩色日志示例 ===")

	logger := fastlog.New(&fastlog.Config{
		Level:         fastlog.DEBUG,
		OutputConsole: true,
		NoColor:       false, // 启用彩色输出
	})
	defer func() { _ = logger.Close() }()

	logger.Debug("调试信息")
	logger.Info("信息日志")
	logger.Warn("警告日志")
	logger.Error("错误日志")

	// 用 recover 捕获 Panic, 确保后续 Fatal 能执行
	func() {
		defer func() { _ = recover() }()
		logger.Panic("恐慌日志")
	}()

	logger.Fatal("致命日志")
}

// example2 禁用颜色示例
func example2() {
	println("\n=== 禁用颜色示例 (NoColor = true) ===")

	logger := fastlog.New(&fastlog.Config{
		Level:         fastlog.DEBUG,
		OutputConsole: true,
		NoColor:       true, // 禁用彩色输出
	})
	defer func() { _ = logger.Close() }()

	logger.Debug("禁用颜色后的调试信息")
	logger.Info("禁用颜色后的信息日志")
	logger.Warn("禁用颜色后的警告日志")
	logger.Error("禁用颜色后的错误日志")
}

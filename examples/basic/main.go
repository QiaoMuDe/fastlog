package main

import (
	"time"

	fastlog "gitee.com/MM-Q/fastlog"
)

func main() {
	// 示例 1: 基础使用
	example1()

	// 示例 2: 不同格式
	example2()

	// 示例 3: 键值对字段
	example3()

	// 示例 4: 配置字段和调用者信息
	example4()

	// 示例 5: 创建记录器
	example5()

	// 示例 6: 日志采样
	example6()
}

// example1 基础使用示例
func example1() {
	println("=== 示例 1: 基础使用 ===")

	logger := fastlog.New(&fastlog.Config{
		Level:         fastlog.INFO,
		OutputConsole: true,
	})
	defer func() { _ = logger.Close() }()

	logger.Info("应用启动成功")
	logger.Debug("调试信息") // 不会输出, 因为级别是 INFO
	logger.Warn("警告信息")
	logger.Error("错误信息")
}

// example2 不同格式示例
func example2() {
	println("\n=== 示例 2: 不同格式 ===")

	// JSON 格式
	jsonLogger := fastlog.New(&fastlog.Config{
		Level:         fastlog.DEBUG,
		OutputConsole: true,
		Formatter:     fastlog.JSON{},
	})
	jsonLogger.Info("JSON 格式日志")
	jsonLogger.Debugw("调试日志",
		fastlog.String("module", "test"),
		fastlog.Int("count", 42),
	)

	// 时间戳格式
	tsLogger := fastlog.New(&fastlog.Config{
		Level:         fastlog.INFO,
		OutputConsole: true,
		Formatter:     fastlog.Timestamp{},
	})
	tsLogger.Info("时间戳格式日志")

	// 键值对格式
	kvLogger := fastlog.New(&fastlog.Config{
		Level:         fastlog.INFO,
		OutputConsole: true,
		Formatter:     fastlog.KV{},
	})
	kvLogger.Info("键值对格式日志")

	// LogFmt 格式
	logfmtLogger := fastlog.New(&fastlog.Config{
		Level:         fastlog.INFO,
		OutputConsole: true,
		Formatter:     fastlog.LogFmt{},
	})
	logfmtLogger.Info("LogFmt 格式日志")
}

// example3 键值对字段示例
func example3() {
	println("\n=== 示例 3: 键值对字段 ===")

	logger := fastlog.New(&fastlog.Config{
		Level:         fastlog.DEBUG,
		OutputConsole: true,
	})

	logger.Infow("用户登录",
		fastlog.String("username", "admin"),
		fastlog.Int("user_id", 12345),
		fastlog.Bool("success", true),
		fastlog.Float64("score", 98.5),
	)

	logger.Debugw("请求详情",
		fastlog.String("method", "GET"),
		fastlog.String("path", "/api/users"),
		fastlog.Int("status", 200),
		fastlog.Duration("latency", 150*time.Millisecond),
	)
}

// example4 配置字段和调用者信息示例
func example4() {
	println("\n=== 示例 4: 配置字段和调用者信息 ===")

	logger := fastlog.New(&fastlog.Config{
		Level:         fastlog.DEBUG,
		OutputConsole: true,
		Formatter:     fastlog.Def{},
		Caller:        true,
		Fields: []fastlog.Field{
			fastlog.String("app", "myapp"),
			fastlog.String("version", "1.0.0"),
		},
	})

	logger.Info("应用启动")
	logger.Infow("用户操作", fastlog.String("action", "login"))
}

// example5 创建记录器示例
func example5() {
	println("\n=== 示例 5: 创建记录器 ===")

	logger := fastlog.New(&fastlog.Config{
		Level:         fastlog.INFO,
		OutputConsole: true,
		Formatter:     fastlog.Def{},
	})
	logger.Info("日志信息")
	logger.Infof("格式化日志: %s", "测试")
	logger.Warn("警告")

	logger.Infow("用户操作",
		fastlog.String("action", "login"),
		fastlog.String("user", "testuser"),
	)
}

// example6 日志采样示例
func example6() {
	println("\n=== 示例 6: 日志采样 ===")

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

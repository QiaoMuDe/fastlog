package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"gitee.com/MM-Q/fastlog"
)

func main() {
	// 创建日志目录
	logDir := "./levelrouter_logs"
	_ = os.MkdirAll(logDir, 0755)

	// 创建配置，启用级别路由
	cfg := fastlog.NewConfig(logDir + "/app.log")
	cfg.LevelRouter = true // 启用级别路由功能
	cfg.Level = fastlog.DEBUG
	cfg.Caller = true

	// 创建日志记录器
	logger := fastlog.New(cfg)
	defer func() { _ = logger.Close() }()

	fmt.Println("=== 级别路由功能演示 ===")
	fmt.Printf("日志目录: %s\n", logDir)
	fmt.Println("启用级别路由后，日志会同时写入:")
	fmt.Println("  1. app.log - 全量日志（所有级别）")
	fmt.Println("  2. DEBUG.log - 调试级别日志")
	fmt.Println("  3. INFO.log - 信息级别日志")
	fmt.Println("  4. WARN.log - 警告级别日志")
	fmt.Println("  5. ERROR.log - 错误级别日志")
	fmt.Println("  6. FATAL.log - 致命级别日志")
	fmt.Println("  7. PANIC.log - 恐慌级别日志")
	fmt.Println()

	// 模拟业务运行，随机生成不同级别的日志
	fmt.Println("模拟业务运行中...")
	fmt.Println()

	operations := []string{
		"用户登录",
		"查询数据库",
		"调用外部API",
		"处理订单",
		"发送邮件",
		"生成报表",
		"清理缓存",
		"备份数据",
	}

	// 模拟 20 次业务操作
	for i := 0; i < 20; i++ {
		// 模拟业务耗时 100-500ms
		processingTime := time.Duration(100+rand.Intn(400)) * time.Millisecond
		time.Sleep(processingTime)

		// 随机选择一个操作
		op := operations[rand.Intn(len(operations))]

		// 根据随机数决定日志级别，模拟不同场景
		r := rand.Intn(100)
		switch {
		case r < 50:
			// 50% 概率 - DEBUG：详细的调试信息
			logger.Debugw(fmt.Sprintf("[%s] 处理中...", op),
				fastlog.Int("iteration", i+1),
				fastlog.Duration("elapsed", processingTime),
			)

		case r < 80:
			// 30% 概率 - INFO：正常业务信息
			logger.Infow(fmt.Sprintf("[%s] 处理完成", op),
				fastlog.Int("iteration", i+1),
				fastlog.Duration("processing_time", processingTime),
			)

		case r < 95:
			// 15% 概率 - WARN：警告信息
			logger.Warnw(fmt.Sprintf("[%s] 处理较慢", op),
				fastlog.Int("iteration", i+1),
				fastlog.Duration("processing_time", processingTime),
				fastlog.String("suggestion", "建议优化性能"),
			)

		default:
			// 5% 概率 - ERROR：错误信息
			logger.Errorw(fmt.Sprintf("[%s] 处理失败", op),
				fastlog.Int("iteration", i+1),
				fastlog.Duration("processing_time", processingTime),
				fastlog.String("error", "模拟错误信息"),
			)
		}
	}

	fmt.Println()
	fmt.Println("=== 演示完成 ===")
	fmt.Println()
	fmt.Println("请查看以下日志文件:")
	fmt.Printf("  %s/app.log    - 包含所有级别的全量日志\n", logDir)
	fmt.Printf("  %s/DEBUG.log  - 仅包含 DEBUG 级别的日志\n", logDir)
	fmt.Printf("  %s/INFO.log   - 仅包含 INFO 级别的日志\n", logDir)
	fmt.Printf("  %s/WARN.log   - 仅包含 WARN 级别的日志\n", logDir)
	fmt.Printf("  %s/ERROR.log  - 仅包含 ERROR 级别的日志\n", logDir)
	fmt.Println()
	fmt.Println("级别路由的优势:")
	fmt.Println("  1. 全量日志用于完整审计和问题追踪")
	fmt.Println("  2. 按级别分离便于快速定位特定严重程度的日志")
	fmt.Println("  3. ERROR.log 可以直接用于错误监控和告警")
	fmt.Println("  4. 避免在大文件中搜索特定级别日志的低效操作")
}

package main

import (
	"fmt"

	fastlog "gitee.com/MM-Q/fastlog"
)

func main() {
	printFormats()
}

func printFormats() {
	// =========================================================
	// 1. Def 格式 (默认格式)
	// 格式: 2025-01-15 10:30:45 | INFO    | main.go:main:15 - 消息
	// =========================================================
	defLogger := fastlog.New(&fastlog.Config{
		Level:         fastlog.DEBUG,
		OutputConsole: true,
		Formatter:     fastlog.Def{},
		Caller:        true,
	})
	printHeader("Def 格式 - Info")
	for i := 0; i < 10; i++ {
		defLogger.Info(fmt.Sprintf("用户登录操作 第%d次", i+1))
	}

	printHeader("Def 格式 - Infof")
	for i := 0; i < 10; i++ {
		defLogger.Infof("用户 %s 登录 %s, 第%d次", "admin", "控制台", i+1)
	}

	printHeader("Def 格式 - Infow")
	for i := 0; i < 10; i++ {
		defLogger.Infow("用户操作",
			fastlog.String("user", "admin"),
			fastlog.Int("seq", i+1),
			fastlog.String("action", "login"),
		)
	}
	_ = defLogger.Close()

	// =========================================================
	// 2. JSON 格式
	// 格式: {"time":"...","level":"INFO","message":"...","key":"val"}
	// =========================================================
	jsonLogger := fastlog.New(&fastlog.Config{
		Level:         fastlog.DEBUG,
		OutputConsole: true,
		Formatter:     fastlog.JSON{},
	})
	printHeader("JSON 格式 - Info")
	for i := 0; i < 10; i++ {
		jsonLogger.Info(fmt.Sprintf("用户登录操作 第%d次", i+1))
	}

	printHeader("JSON 格式 - Infof")
	for i := 0; i < 10; i++ {
		jsonLogger.Infof("用户 %s 登录 %s, 第%d次", "admin", "控制台", i+1)
	}

	printHeader("JSON 格式 - Infow")
	for i := 0; i < 10; i++ {
		jsonLogger.Infow("用户操作",
			fastlog.String("user", "admin"),
			fastlog.Int("seq", i+1),
			fastlog.String("action", "login"),
		)
	}
	_ = jsonLogger.Close()

	// =========================================================
	// 3. Simple 格式
	// 格式: 2025-01-15 10:30:45 INFO 消息
	// =========================================================
	simpleLogger := fastlog.New(&fastlog.Config{
		Level:         fastlog.DEBUG,
		OutputConsole: true,
		Formatter:     fastlog.Simple{},
	})
	printHeader("Simple 格式 - Info")
	for i := 0; i < 10; i++ {
		simpleLogger.Info(fmt.Sprintf("用户登录操作 第%d次", i+1))
	}

	printHeader("Simple 格式 - Infof")
	for i := 0; i < 10; i++ {
		simpleLogger.Infof("用户 %s 登录 %s, 第%d次", "admin", "控制台", i+1)
	}

	printHeader("Simple 格式 - Infow")
	for i := 0; i < 10; i++ {
		simpleLogger.Infow("用户操作",
			fastlog.String("user", "admin"),
			fastlog.Int("seq", i+1),
			fastlog.String("action", "login"),
		)
	}
	_ = simpleLogger.Close()

	// =========================================================
	// 4. KV 格式 (键值对格式)
	// 格式: time=... level=INFO message=... key=val
	// =========================================================
	kvLogger := fastlog.New(&fastlog.Config{
		Level:         fastlog.DEBUG,
		OutputConsole: true,
		Formatter:     fastlog.KV{},
	})
	printHeader("KV 格式 - Info")
	for i := 0; i < 10; i++ {
		kvLogger.Info(fmt.Sprintf("用户登录操作 第%d次", i+1))
	}

	printHeader("KV 格式 - Infof")
	for i := 0; i < 10; i++ {
		kvLogger.Infof("用户 %s 登录 %s, 第%d次", "admin", "控制台", i+1)
	}

	printHeader("KV 格式 - Infow")
	for i := 0; i < 10; i++ {
		kvLogger.Infow("用户操作",
			fastlog.String("user", "admin"),
			fastlog.Int("seq", i+1),
			fastlog.String("action", "login"),
		)
	}
	_ = kvLogger.Close()

	// =========================================================
	// 5. Compact 格式 (极简格式, 时间格式遵循 TimeFormat)
	// 格式: [I] 2025-01-15 10:30:45 用户登录成功 | user=admin seq=1 action=login
	// =========================================================
	compactLogger := fastlog.New(&fastlog.Config{
		Level:         fastlog.DEBUG,
		OutputConsole: true,
		Formatter:     fastlog.Compact{},
	})
	printHeader("Compact 格式 - Info")
	for i := 0; i < 10; i++ {
		compactLogger.Info(fmt.Sprintf("用户登录操作 第%d次", i+1))
	}

	printHeader("Compact 格式 - Infof")
	for i := 0; i < 10; i++ {
		compactLogger.Infof("用户 %s 登录 %s, 第%d次", "admin", "控制台", i+1)
	}

	printHeader("Compact 格式 - Infow")
	for i := 0; i < 10; i++ {
		compactLogger.Infow("用户操作",
			fastlog.String("user", "admin"),
			fastlog.Int("seq", i+1),
			fastlog.String("action", "login"),
		)
	}
	_ = compactLogger.Close()
}

func printHeader(title string) {
	fmt.Printf("\n========== %s ==========\n", title)
}

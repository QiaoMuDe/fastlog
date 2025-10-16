package fastlog

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gitee.com/MM-Q/fastlog/internal/config"
	"gitee.com/MM-Q/fastlog/internal/types"
)

// TestMain 全局测试入口，控制非verbose模式下的输出重定向
func TestMain(m *testing.M) {
	flag.Parse() // 解析命令行参数
	// 保存原始标准输出和错误输出
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	var nullFile *os.File
	var err error

	// 非verbose模式下重定向到空设备
	if !testing.Verbose() {
		nullFile, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
		if err != nil {
			panic("无法打开空设备文件: " + err.Error())
		}
		os.Stdout = nullFile
		os.Stderr = nullFile
	}

	// 运行所有测试
	exitCode := m.Run()

	//清理日志目录
	if err := os.RemoveAll("logs"); err != nil {
		fmt.Printf("清理日志目录失败: %v\n", err)
	}

	// 恢复原始输出
	if !testing.Verbose() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
		_ = nullFile.Close()
	}

	os.Exit(exitCode)
}

// TestCustomFormat 测试自定义日志格式
func TestCustomFormat(t *testing.T) {
	// 创建日志配置
	cfg := config.NewFastLogConfig("logs", "custom.log")
	cfg.LogLevel = types.DEBUG
	cfg.LogFormat = types.Custom

	// 创建日志记录器
	log := NewFastLog(cfg)
	defer func() { _ = log.Close() }()

	// 定义web应用程序日志格式
	webAppLogFormat := `%s [%s] %s %s %s %d %d %s %s %dms`

	// 模拟web应用程序日志记录
	log.Errorf(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:00", "GET", "/index.html", "HTTP/1.1", 200, 1234, "Mozilla/5.0", "en-US", 100)
	log.Infof(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:01", "POST", "/login", "HTTP/1.1", 401, 0, "Mozilla/5.0", "en-US", 200)
	log.Warnf(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:03", "PUT", "/settings", "HTTP/1.1", 500, 0, "Mozilla/5.0", "en-US", 500)
	log.Debugf(webAppLogFormat, "127.0.0.1", "2023-10-01 12:00:04", "DELETE", "/logout", "HTTP/1.1", 200, 0, "Mozilla/5.0", "en-US", 50)
}

// TestNoColor 测试无颜色日志
func TestNoColor(t *testing.T) {
	// 创建日志配置
	cfg := config.NewFastLogConfig("logs", "nocolor.log")
	cfg.LogLevel = types.DEBUG
	cfg.Color = false // 禁用终端颜色

	// 创建日志记录器
	log := NewFastLog(cfg)

	defer func() { _ = log.Close() }()

	// 打印测试日志
	log.Info("测试无颜色日志")
	log.Warn("测试无颜色日志")
	log.Error("测试无颜色日志")
	log.Debug("测试无颜色日志")
}

// TestNoBold 测试无加粗日志
func TestNoBold(t *testing.T) {
	// 创建日志配置
	cfg := config.NewFastLogConfig("logs", "nobold.log")
	cfg.LogLevel = types.DEBUG
	cfg.Bold = false // 禁用终端字体加粗
	cfg.Color = false

	// 创建日志记录器
	log := NewFastLog(cfg)

	defer func() { _ = log.Close() }()

	// 打印测试日志
	log.Info("测试无加粗日志")
	log.Warn("测试无加粗日志")
	log.Error("测试无加粗日志")
	log.Debug("测试无加粗日志")
}

// TestNewFastLog_Initialization 测试日志记录器初始化
func TestNewFastLog_Initialization(t *testing.T) {
	// 测试nil配置错误 - 应该会panic
	t.Run("nil config should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("使用nil配置时应该panic")
			}
		}()
		// 这行代码应该会panic
		_ = NewFastLog(nil)
	})

	// 测试正常初始化 - 使用临时目录避免文件冲突
	t.Run("normal initialization", func(t *testing.T) {
		cfg := config.NewFastLogConfig("logs", "init_test.log")
		cfg.OutputToConsole = false // 禁用控制台输出，避免测试干扰
		cfg.OutputToFile = true     // 只测试文件输出

		log := NewFastLog(cfg)

		// 验证日志实例是否正确创建
		if log == nil {
			t.Fatal("日志实例不应为nil")
		}

		// 简单测试日志功能是否正常
		log.Info("初始化测试日志")

		// 给processor一些时间完成初始化和处理
		time.Sleep(100 * time.Millisecond)

		_ = log.Close()

		// 验证日志文件是否创建成功
		logFile := filepath.Join("logs", "init_test.log")
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Error("日志文件未创建")
		}
	})
}

// TestLogFormats 测试日志库支持的所有日志格式
func TestLogFormats(t *testing.T) {
	// 定义日志格式及其对应的文件名
	formats := map[types.LogFormatType]string{
		types.Def:       "def.log",
		types.Json:      "json.log",
		types.Timestamp: "simple_timestamp.log",
		// 注意：对于Custom格式，日志库内部不进行格式化，需要在外部格式化后传入
		types.Custom: "custom.log",
	}

	// 为每种格式创建日志记录器并测试所有日志级别
	for format, filename := range formats {
		t.Run(fmt.Sprintf("Format_%s", filename), func(t *testing.T) {
			// 创建日志配置
			cfg := config.NewFastLogConfig("logs", filename)
			cfg.LogLevel = types.DEBUG
			cfg.LogFormat = format
			cfg.OutputToConsole = true
			cfg.CallerInfo = false

			// 创建日志记录器
			log := NewFastLog(cfg)
			defer func() { _ = log.Close() }()

			// 测试所有日志级别
			// 对于Custom格式，需要外部格式化
			if format == types.Custom {
				log.Debugf("[自定义格式] [%s] %s", "DEBUG", "这是一条调试日志")
				log.Infof("[自定义格式] [%s] %s", "INFO", "这是一条信息日志")
				log.Warnf("[自定义格式] [%s] %s", "WARN", "这是一条警告日志")
				log.Errorf("[自定义格式] [%s] %s", "ERROR", "这是一条错误日志")
			} else {
				log.Debug("这是一条调试日志")
				log.Info("这是一条信息日志")
				log.Warn("这是一条警告日志")
				log.Error("这是一条错误日志")

				log.DebugF("这是一条键值对调试日志",
					String("key1", "value1"),
					Int("key2", 42),
				)

				log.InfoF("这是一条键值对信息日志",
					String("key1", "value1"),
					Int("key2", 42),
				)

				log.WarnF("这是一条键值对警告日志",
					String("key1", "value1"),
					Int("key2", 42),
				)

				log.ErrorF("这是一条键值对错误日志",
					String("key1", "value1"),
					Int("key2", 42),
				)
			}

			// 给日志处理器一些时间来处理和写入日志
			time.Sleep(100 * time.Millisecond)
		})
	}
}

// TestFatal 测试Fatal方法
func TestFatal(t *testing.T) {
	// 定义测试用例名称，用于环境变量标识
	const testName = "TestFatal"

	// 子进程模式：执行实际的Fatal调用
	if os.Getenv("TEST_MODE") == testName {
		config := config.NewFastLogConfig("test-logs", "fatal_test.log")
		if err := os.MkdirAll(config.LogDirName, 0755); err != nil {
			panic(err)
		}
		log := NewFastLog(config)
		log.Fatal("fatal_test message")
		return
	}

	// 主进程模式：启动子进程并检查结果
	cmd := exec.Command(os.Args[0], "-test.run=^"+testName+"$", "-test.v")
	cmd.Env = append(os.Environ(), "TEST_MODE="+testName)
	output, err := cmd.CombinedOutput()

	// 检查是否以预期的退出码1退出
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitCode := exitErr.ExitCode(); exitCode != 1 {
			t.Errorf("期望退出码1，实际得到%d\n输出: %s", exitCode, output)
		}
	} else if err != nil {
		t.Fatalf("测试命令执行失败: %v\n输出: %s", err, output)
	} else {
		t.Error("Fatal方法没有导致程序退出")
	}

	// 验证日志文件是否正确创建并包含预期内容
	logPath := filepath.Join("test-logs", "fatal_test.log")
	defer func() {
		// 确保测试后清理文件
		_ = os.Remove(logPath)
		_ = os.Remove("test-logs")
	}()

	if _, statErr := os.Stat(logPath); os.IsNotExist(statErr) {
		t.Errorf("日志文件未创建: %s\n子进程输出: %s", logPath, output)
		return
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Errorf("无法读取日志文件: %v", err)
	} else if !strings.Contains(string(content), "fatal_test message") {
		t.Errorf("日志内容不包含预期消息，内容: %s", content)
	}
}

// TestFatalf 测试Fatalf方法
func TestFatalf(t *testing.T) {
	// 定义测试用例名称，用于环境变量标识
	const testName = "TestFatalf"

	// 子进程模式：执行实际的Fatalf调用
	if os.Getenv("TEST_MODE") == testName {
		config := config.NewFastLogConfig("logs", "fatalf_test.log")
		log := NewFastLog(config)
		defer func() { _ = log.Close() }()

		log.Fatalf("fatalf_test %s message", "formatted")
		return
	}

	// 主进程模式：启动子进程并检查结果
	cmd := exec.Command(os.Args[0], "-test.run=^"+testName+"$", "-test.v")
	cmd.Env = append(os.Environ(), "TEST_MODE="+testName)
	output, err := cmd.CombinedOutput()

	// 检查是否以预期的退出码1退出
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitCode := exitErr.ExitCode(); exitCode != 1 {
			t.Errorf("期望退出码1，实际得到%d\n输出: %s", exitCode, output)
		}
	} else if err != nil {
		t.Fatalf("测试命令执行失败: %v\n输出: %s", err, output)
	} else {
		t.Error("Fatalf方法没有导致程序退出")
	}

	// 验证日志文件是否正确创建并包含预期内容
	logPath := filepath.Join("logs", "fatalf_test.log")
	defer func() {
		// 确保测试后清理文件
		_ = os.Remove(logPath)
		_ = os.Remove("logs")
	}()

	if _, statErr := os.Stat(logPath); os.IsNotExist(statErr) {
		t.Errorf("日志文件未创建: %s\n子进程输出: %s", logPath, output)
		return
	}

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Errorf("无法读取日志文件: %v", err)
	} else if !strings.Contains(string(content), "fatalf_test formatted message") {
		t.Errorf("日志内容不包含预期消息，内容: %s", content)
	}
}

// TestLog 测试使用flog库记录不同级别的日志
func TestLog(t *testing.T) {
	// 创建开发环境配置
	cfg := config.DevConfig("logs", "app.log")

	// 也可以自定义配置
	cfg.LogLevel = types.DEBUG
	cfg.LogFormat = types.Json
	cfg.CallerInfo = true
	cfg.OutputToFile = true

	// 创建flog实例
	logger := NewFastLog(cfg)
	defer func() { _ = logger.Close() }()

	// 记录不同级别的日志
	logger.Debug("这是一条调试日志")
	logger.Info("这是一条信息日志")
	logger.Warn("这是一条警告日志")
	logger.Error("这是一条错误日志")

	// 使用字段记录结构化日志
	logger.Info("用户登录",
		String("user_id", "12345"),
		String("username", "alice"),
		Int("age", 25),
		Bool("is_admin", false))

	// 记录带错误信息的日志
	err := errors.New("数据库连接失败")
	logger.Error("处理用户请求时发生错误",
		String("request_id", "req-12345"),
		String("url", "/api/users"),
		Error("error", err))

	// 记录带时间信息的日志
	logger.Info("任务执行完成",
		String("task_name", "data_backup"),
		Duration("duration", 5*time.Second),
		Int("records_processed", 1000))

	// 记录数值型字段
	logger.Info("系统状态监控",
		Int("cpu_usage", 75),
		Float64("memory_usage", 85.5),
		Uint64("disk_space", 1024000))

	// 使用不同日志格式
	jsonCfg := config.NewFastLogConfig("logs", "json.log")
	jsonCfg.LogFormat = types.Json
	jsonCfg.CallerInfo = true

	jsonLogger := NewFastLog(jsonCfg)
	defer func() { _ = jsonLogger.Close() }()

	jsonLogger.Info("JSON格式日志示例",
		String("service", "user-service"),
		Int("status_code", 200))
}

// TestFlogUsage 测试flog库的基本使用方法
func TestFlogUsage(t *testing.T) {
	// 创建测试配置 - 使用控制台模式避免文件写入
	cfg := config.ConsoleConfig()
	cfg.LogFormat = types.Json

	// 创建flog实例
	logger := NewFastLog(cfg)
	if logger == nil {
		t.Fatal("Failed to create flog instance")
	}
	defer func() { _ = logger.Close() }()

	// 测试不同级别的日志记录
	t.Run("BasicLogLevels", func(t *testing.T) {
		logger.Debug("这是一条调试日志")
		logger.Info("这是一条信息日志")
		logger.Warn("这是一条警告日志")
		logger.Error("这是一条错误日志")
	})

	// 测试带字段的日志记录
	t.Run("LogsWithFields", func(t *testing.T) {
		logger.Info("用户登录",
			String("user_id", "12345"),
			String("username", "alice"),
			Int("age", 25),
			Bool("is_admin", false))
	})

	// 测试带错误信息的日志记录
	t.Run("LogsWithError", func(t *testing.T) {
		err := errors.New("数据库连接失败")
		logger.Error("处理用户请求时发生错误",
			String("request_id", "req-12345"),
			String("url", "/api/users"),
			Error("error", err))
	})

	// 测试带时间信息的日志记录
	t.Run("LogsWithTimeInfo", func(t *testing.T) {
		logger.Info("任务执行完成",
			String("task_name", "data_backup"),
			Duration("duration", 5*time.Second),
			Int("records_processed", 1000))
	})

	// 测试数值型字段
	t.Run("LogsWithNumericFields", func(t *testing.T) {
		logger.Info("系统状态监控",
			Int("cpu_usage", 75),
			Float64("memory_usage", 85.5),
			Uint64("disk_space", 1024000))
	})
}

// TestFlogWithDifferentFormats 测试不同日志格式
func TestFlogWithDifferentFormats(t *testing.T) {
	formats := []types.LogFormatType{
		types.Def,
		types.Json,
		types.Timestamp,
	}

	for _, format := range formats {
		t.Run(format.String(), func(t *testing.T) {
			// 创建测试配置
			cfg := config.ConsoleConfig()
			cfg.LogFormat = format

			// 创建flog实例
			logger := NewFastLog(cfg)
			if logger == nil {
				t.Fatal("Failed to create flog instance")
			}
			defer func() { _ = logger.Close() }()

			// 记录带字段的日志
			logger.Info("测试不同日志格式",
				String("service", "test-service"),
				Int("status_code", 200),
				String("user_id", "12345"))
		})
	}
}

// TestFlogConfigurations 测试不同的配置选项
func TestFlogConfigurations(t *testing.T) {
	// 测试启用调用者信息
	t.Run("WithCallerInfo", func(t *testing.T) {
		cfg := config.ConsoleConfig()
		cfg.CallerInfo = true

		logger := NewFastLog(cfg)
		if logger == nil {
			t.Fatal("Failed to create flog instance")
		}
		defer func() { _ = logger.Close() }()

		logger.Info("启用调用者信息的日志")
	})

	// 测试不同日志级别
	t.Run("DifferentLogLevels", func(t *testing.T) {
		levels := []types.LogLevel{
			types.DEBUG,
			types.INFO,
			types.WARN,
			types.ERROR,
		}

		for _, level := range levels {
			t.Run(level.String(), func(t *testing.T) {
				cfg := config.ConsoleConfig()
				cfg.LogLevel = level

				logger := NewFastLog(cfg)
				if logger == nil {
					t.Fatal("Failed to create flog instance")
				}
				defer func() { _ = logger.Close() }()

				// 尝试记录所有级别的日志
				logger.Debug("调试信息")
				logger.Info("一般信息")
				logger.Warn("警告信息")
				logger.Error("错误信息")
			})
		}
	})
}

// BenchmarkFlogInfo 简单的性能基准测试
func BenchmarkFlogInfo(b *testing.B) {
	// 创建测试配置
	cfg := config.ConsoleConfig()
	cfg.OutputToConsole = false // 关闭控制台输出以提高性能测试准确性

	// 创建flog实例
	logger := NewFastLog(cfg)
	if logger == nil {
		b.Fatal("Failed to create flog instance")
	}
	defer func() { _ = logger.Close() }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("性能测试日志")
	}
}

// BenchmarkFlogInfoWithFields 带字段的性能基准测试
func BenchmarkFlogInfoWithFields(b *testing.B) {
	// 创建测试配置
	cfg := config.ConsoleConfig()
	cfg.OutputToConsole = false // 关闭控制台输出以提高性能测试准确性

	// 创建flog实例
	logger := NewFastLog(cfg)
	if logger == nil {
		b.Fatal("Failed to create flog instance")
	}
	defer func() { _ = logger.Close() }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("带字段的性能测试日志",
			String("key1", "value1"),
			Int("key2", 123),
			Bool("key3", true))
	}
}

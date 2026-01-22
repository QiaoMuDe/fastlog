package colorlib

import (
	"os"
	"strconv"
)

// ConcurrentTestConfig 并发测试配置
type ConcurrentTestConfig struct {
	NumGoroutines int  // goroutine数量
	NumOperations int  // 每个goroutine的操作次数
	EnableStress  bool // 是否启用压力测试模式
}

// GetTestConfig 获取测试配置
func GetTestConfig() ConcurrentTestConfig {
	config := ConcurrentTestConfig{
		NumGoroutines: 100,  // 默认100个goroutine
		NumOperations: 1000, // 默认每个goroutine执行1000次操作
		EnableStress:  false,
	}

	// 从环境变量读取配置
	if val := os.Getenv("TEST_GOROUTINES"); val != "" {
		if num, err := strconv.Atoi(val); err == nil && num > 0 {
			config.NumGoroutines = num
		}
	}

	if val := os.Getenv("TEST_OPERATIONS"); val != "" {
		if num, err := strconv.Atoi(val); err == nil && num > 0 {
			config.NumOperations = num
		}
	}

	if val := os.Getenv("STRESS_TEST"); val != "" {
		config.EnableStress = val == "1" || val == "true"
		if config.EnableStress {
			// 压力测试模式下增加并发量
			config.NumGoroutines *= 2
			config.NumOperations *= 2
		}
	}

	return config
}

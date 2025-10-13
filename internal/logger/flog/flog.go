package flog

import (
	"sync/atomic"

	"gitee.com/MM-Q/colorlib"
	"gitee.com/MM-Q/fastlog/internal/config"
	"gitee.com/MM-Q/logrotatex"
)

// Flog 日志记录器
type Flog struct {
	fileWriter *logrotatex.BufferedWriter // 带缓冲的文件写入器
	cl         *colorlib.ColorLib         // 提供终端颜色输出的库
	cfg        *config.FastLogConfig      // 嵌入的配置结构体
	closed     atomic.Bool                // 标记日志处理器是否已关闭
	fields     map[string]interface{}     // 自定义字段, 用于在日志中添加额外的上下文信息
}

// NewFlog 创建一个新的Flog实例, 用于记录日志。
//
// 参数:
//   - config: 一个指向FastLogConfig实例的指针, 用于配置日志记录器。
//
// 返回值:
//   - *Flog: 一个指向Flog实例的指针。
func NewFlog(cfg *config.FastLogConfig) *Flog {
	// 检查配置结构体是否为nil
	if cfg == nil {
		panic("FastLogConfig cannot be nil")
	}

	// 最终验证
	cfg.ValidateConfig()

	// 克隆配置结构体防止原配置被意外修改
	cfg = cfg.Clone()

	// 初始化写入器
	fileWriter := config.CreateBufferedWriter(cfg)

	// 创建一个新的Flog实例, 将配置和缓冲区赋值给实例。
	f := &Flog{
		fileWriter: fileWriter,             // 带缓冲的文件写入器
		cl:         colorlib.NewColorLib(), // 颜色库实例, 用于在终端中显示颜色
		cfg:        cfg,                    // 配置结构体
		closed:     atomic.Bool{},          // 标记日志处理器是否已关闭
	}

	// 配置设置
	f.cl.SetColor(f.cfg.Color) // 设置颜色库的颜色选项
	f.cl.SetBold(f.cfg.Bold)   // 设置颜色库的字体加粗选项
	f.closed.Store(false)      // 初始化日志处理器为未关闭状态

	// 返回Flog实例
	return f
}

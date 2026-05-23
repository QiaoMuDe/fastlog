package fastlog

import "io"

// hook 日志钩子接口（内部使用，小写不导出）
// 用于在日志输出时执行额外操作，如按级别分发到不同文件
type hook interface {
	// Fire 日志触发时调用
	// 参数:
	//   - entry: 日志条目，包含完整信息
	//   - data: 格式化后的日志数据
	// 返回:
	//   - error: 执行过程中的错误
	Fire(entry *Entry, data []byte) error

	// Levels 返回关心的日志级别列表
	// 只有这些级别的日志会触发 Fire
	Levels() []Level

	// Sync 同步日志到存储
	// 返回:
	//   - error: 同步过程中的错误
	Sync() error

	// Close 关闭钩子资源
	// 返回:
	//   - error: 关闭过程中的错误
	Close() error
}

// levelHook 按级别分发的钩子（内部使用）
// 将特定级别的日志写入指定目标
type levelHook struct {
	level  Level          // 关心的级别
	writer io.WriteCloser // 写入目标
}

// Fire 执行钩子
// 如果 entry.Level 匹配，则写入专属文件
// 参数:
//   - entry: 日志条目
//   - data: 格式化后的日志数据
//
// 返回:
//   - error: 写入过程中的错误
func (h *levelHook) Fire(entry *Entry, data []byte) error {
	if entry.Level != h.level {
		return nil // 级别不匹配，忽略
	}
	_, err := h.writer.Write(data)
	return err
}

// Levels 返回关心的级别
// 返回:
//   - []Level: 只包含一个级别
func (h *levelHook) Levels() []Level {
	return []Level{h.level}
}

// Sync 同步日志到存储
// 如果写入器支持 Sync 方法则调用，否则返回 nil
// 返回:
//   - error: 同步过程中的错误
func (h *levelHook) Sync() error {
	if h.writer != nil {
		if syncer, ok := h.writer.(interface{ Sync() error }); ok {
			return syncer.Sync()
		}
	}
	return nil
}

// Close 关闭写入器
// 返回:
//   - error: 关闭过程中的错误
func (h *levelHook) Close() error {
	if h.writer != nil {
		return h.writer.Close()
	}
	return nil
}

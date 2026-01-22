package pool

// 字节单位定义
const (
	Byte = 1 << (10 * iota) // 1 字节
	KB                      // 千字节 (1024 B)
	MB                      // 兆字节 (1024 KB)
	GB                      // 吉字节 (1024 MB)
	TB                      // 太字节 (1024 GB)
)

// CalculateBufferSize 根据文件大小动态计算最佳缓冲区大小。
// 采用分层策略，平衡内存使用和I/O性能。
//
// 参数:
//   - fileSize: 文件大小（字节）
//
// 返回:
//   - int: 计算出的最佳缓冲区大小（字节）
//
// 缓冲区分配策略:
//   - ≤ 0 或 ≤ 4KB: 使用 1KB 缓冲区，确保最小缓冲区大小
//   - 4KB - 32KB: 使用 8KB 缓冲区
//   - 32KB - 128KB: 使用 32KB 缓冲区
//   - 128KB - 512KB: 使用 64KB 缓冲区
//   - 512KB - 1MB: 使用 128KB 缓冲区
//   - 1MB - 4MB: 使用 256KB 缓冲区
//   - 4MB - 16MB: 使用 512KB 缓冲区
//   - 16MB - 64MB: 使用 1MB 缓冲区
//   - > 64MB: 使用 2MB 缓冲区
//
// 设计原则:
//   - 极小文件: 最小化内存占用
//   - 小文件: 适度缓冲，节省内存
//   - 大文件: 增大缓冲区，提升I/O吞吐量
//   - 超大文件: 限制最大缓冲区，避免过度内存消耗
func CalculateBufferSize(fileSize int64) int {
	switch {
	case fileSize <= 0: // 空文件或无效大小，使用最小1KB缓冲区
		return int(KB)
	case fileSize <= 4*KB: // 极小文件使用1KB缓冲区，避免过小缓冲区
		return int(KB)
	case fileSize < 32*KB: // 小于 32KB 的文件使用 8KB 缓冲区
		return int(8 * KB)
	case fileSize < 128*KB: // 32KB-128KB 使用 32KB 缓冲区
		return int(32 * KB)
	case fileSize < 512*KB: // 128KB-512KB 使用 64KB 缓冲区
		return int(64 * KB)
	case fileSize < 1*MB: // 512KB-1MB 使用 128KB 缓冲区
		return int(128 * KB)
	case fileSize < 4*MB: // 1MB-4MB 使用 256KB 缓冲区
		return int(256 * KB)
	case fileSize < 16*MB: // 4MB-16MB 使用 512KB 缓冲区
		return int(512 * KB)
	case fileSize < 64*MB: // 16MB-64MB 使用 1MB 缓冲区
		return int(1 * MB)
	default: // 大于 64MB 的文件使用 2MB 缓冲区
		return int(2 * MB)
	}
}

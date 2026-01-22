# Package pool

**Import Path:** `gitee.com/MM-Q/go-kit/pool`

Package pool 提供高性能对象池管理功能, 通过复用对象优化内存使用。

该包实现了基于 sync.Pool 的多种对象池, 用于减少频繁的内存分配和回收。 通过复用对象, 可以显著提升应用程序的性能, 特别是在高并发场景下。

**主要功能：**
- 字节切片对象池管理
- 动态容量对象获取
- 自动内存回收控制
- 防止内存泄漏的容量限制
- 支持多种对象类型的池化

**性能优化：**
- 使用 sync.Pool 减少 GC 压力
- 支持不同容量的对象需求
- 自动限制大对象回收
- 预热机制提升冷启动性能

Package pool 提供随机数生成器对象池功能，通过对象池优化随机数生成性能。

随机数生成器对象池专门用于复用math/rand.Rand对象， 避免频繁创建随机数生成器的开销，特别适用于ID生成、测试数据生成等场景。

Package pool 提供Timer对象池功能，通过对象池优化定时器使用。

Timer对象池专门用于复用time.Timer对象，避免频繁创建和销毁定时器的开销。
Timer的创建成本相对较高，特别是在高并发场景下，复用可以显著提升性能。

## Constants

```go
const (
	Byte = 1 << (10 * iota) // 1 字节
	KB                      // 千字节 (1024 B)
	MB                      // 兆字节 (1024 KB)
	GB                      // 吉字节 (1024 MB)
	TB                      // 太字节 (1024 GB)
)
```
字节单位定义

## Functions

### func CalculateBufferSize

```go
func CalculateBufferSize(fileSize int64) int
```

CalculateBufferSize 根据文件大小动态计算最佳缓冲区大小。 采用分层策略，平衡内存使用和I/O性能。

**参数:**
- fileSize: 文件大小（字节）

**返回:**
- int: 计算出的最佳缓冲区大小（字节）

**缓冲区分配策略:**
- ≤ 0 或 ≤ 4KB: 使用 1KB 缓冲区，确保最小缓冲区大小
- 4KB - 32KB: 使用 8KB 缓冲区
- 32KB - 128KB: 使用 32KB 缓冲区
- 128KB - 512KB: 使用 64KB 缓冲区
- 512KB - 1MB: 使用 128KB 缓冲区
- 1MB - 4MB: 使用 256KB 缓冲区
- 4MB - 16MB: 使用 512KB 缓冲区
- 16MB - 64MB: 使用 1MB 缓冲区
- > 64MB: 使用 2MB 缓冲区

**设计原则:**
- 极小文件: 最小化内存占用
- 小文件: 适度缓冲，节省内存
- 大文件: 增大缓冲区，提升I/O吞吐量
- 超大文件: 限制最大缓冲区，避免过度内存消耗

### func GetBuf

```go
func GetBuf() *bytes.Buffer
```

GetBuf 从默认缓冲区池获取默认容量的字节缓冲区

**返回值:**
- *bytes.Buffer: 容量至少为默认容量的字节缓冲区

### func GetBufCap

```go
func GetBufCap(cap int) *bytes.Buffer
```

GetBufCap 从默认缓冲区池获取指定容量的字节缓冲区

**参数:**
- cap: 缓冲区初始容量

**返回值:**
- *bytes.Buffer: 容量至少为capacity的字节缓冲区

### func GetByte

```go
func GetByte() []byte
```

GetByte 从默认字节池获取默认容量的缓冲区

**返回值:**
- []byte: 长度为默认容量, 容量至少为默认容量的缓冲区

### func GetByteCap

```go
func GetByteCap(size int) []byte
```

GetByteCap 从默认字节池获取指定容量的缓冲区

**参数:**
- size: 缓冲区容量

**返回值:**
- []byte: 长度为capacity, 容量至少为capacity的缓冲区

### func GetByteEmpty

```go
func GetByteEmpty(size int) []byte
```

GetByteEmpty 从默认字节池获取空缓冲区

**参数:**
- size: 指定容量要求

**返回值:**
- []byte: 长度为0但容量至少为capacity的缓冲区切片

### func GetRand

```go
func GetRand() *rand.Rand
```

GetRand 从池中获取随机数生成器

**返回值:**
- *rand.Rand: 随机数生成器实例

**说明:**
- 返回的生成器已经初始化了随机种子
- 使用完毕后应调用PutRand归还
- 注意：返回的生成器不是线程安全的，不要在多个goroutine间共享

### func GetRandWithSeed

```go
func GetRandWithSeed(seed int64) *rand.Rand
```

GetRandWithSeed 获取指定种子的随机数生成器

**参数:**
- seed: 随机数种子

**返回值:**
- *rand.Rand: 随机数生成器实例

**说明:**
- 返回的生成器使用指定的种子初始化
- 适用于需要可重现随机序列的场景

### func GetStr

```go
func GetStr() *strings.Builder
```

GetStr 从默认字符串池获取默认容量的字符串构建器

**返回值:**
- *strings.Builder: 容量至少为默认容量的字符串构建器

### func GetStrCap

```go
func GetStrCap(cap int) *strings.Builder
```

GetStrCap 从默认字符串池获取指定容量的字符串构建器

**参数:**
- cap: 字符串构建器初始容量

**返回值:**
- *strings.Builder: 容量至少为capacity的字符串构建器

### func GetTimer

```go
func GetTimer(duration time.Duration) *time.Timer
```

GetTimer 从池中获取定时器并设置超时时间

**参数:**
- duration: 定时器超时时间

**返回值:**
- *time.Timer: 已设置超时时间的定时器

**说明:**
- 返回的定时器已经启动，会在指定时间后触发
- 适用于超时控制场景，定时器会在指定时间后自动触发
- 使用完毕后应调用PutTimer归还

### func GetTimerEmpty

```go
func GetTimerEmpty() *time.Timer
```

GetTimerEmpty 从池中获取未启动的定时器

**返回值:**
- *time.Timer: 未启动的定时器，需要手动调用Reset设置时间

**说明:**
- 适用于需要手动控制定时器启动和停止的场景
- 定时器处于停止状态，不会自动触发
- 获取后需要调用timer.Reset(duration)启动
- 使用完毕后应调用PutTimer归还

### func PutBuf

```go
func PutBuf(buf *bytes.Buffer)
```

PutBuf 将字节缓冲区归还到默认缓冲区池

**参数:**
- buf: 要归还的字节缓冲区

**说明:**
- 该函数将字节缓冲区归还到对象池，以便后续复用。

### func PutByte

```go
func PutByte(buf []byte)
```

PutByte 将缓冲区归还到默认字节池

**参数:**
- buf: 要归还的缓冲区

**说明:**
- 该函数将缓冲区归还到对象池, 以便后续复用。

### func PutRand

```go
func PutRand(rng *rand.Rand)
```

PutRand 将随机数生成器归还到池中

**参数:**
- rng: 要归还的随机数生成器

### func PutStr

```go
func PutStr(buf *strings.Builder)
```

PutStr 将字符串构建器归还到默认字符串池

**参数:**
- buf: 要归还的字符串构建器

**说明:**
- 该函数将字符串构建器归还到对象池，以便后续复用。

### func PutTimer

```go
func PutTimer(timer *time.Timer)
```

PutTimer 将定时器归还到池中

**参数:**
- timer: 要归还的定时器

**说明:**
- 该函数会自动停止定时器并清理状态
- 归还后的定时器会被重置，可以安全复用

### func WithBuf

```go
func WithBuf(fn func(*bytes.Buffer)) []byte
```

WithBuf 使用默认容量的字节缓冲区执行函数，自动管理获取和归还

**参数:**
- fn: 使用字节缓冲区的函数

**返回值:**
- []byte: 函数执行后缓冲区的字节数据副本

**使用示例:**

```go
data := pool.WithBuf(func(buf *bytes.Buffer) {
    buf.WriteString("Hello")
    buf.WriteByte(' ')
    buf.WriteString("World")
})
```

### func WithBufCap

```go
func WithBufCap(cap int, fn func(*bytes.Buffer)) []byte
```

WithBufCap 使用指定容量的字节缓冲区执行函数，自动管理获取和归还

**参数:**
- cap: 字节缓冲区初始容量
- fn: 使用字节缓冲区的函数

**返回值:**
- []byte: 函数执行后缓冲区的字节数据副本

**使用示例:**

```go
data := pool.WithBufCap(1024, func(buf *bytes.Buffer) {
    buf.WriteString("Hello")
    buf.WriteByte(' ')
    buf.WriteString("World")
})
```

### func WithRand

```go
func WithRand[T any](fn func(*rand.Rand) T) T
```

WithRand 使用随机数生成器执行函数，自动管理获取和归还

**参数:**
- fn: 使用随机数生成器的函数

**返回值:**
- T: 函数返回的结果

**使用示例:**

```go
// 生成随机整数
num := pool.WithRand(func(rng *rand.Rand) int {
    return rng.Intn(100)
})

// 生成随机字符串
str := pool.WithRand(func(rng *rand.Rand) string {
    return fmt.Sprintf("id_%d", rng.Int63())
})
```

### func WithRandSeed

```go
func WithRandSeed[T any](seed int64, fn func(*rand.Rand) T) T
```

WithRandSeed 使用指定种子的随机数生成器执行函数，自动管理获取和归还

**参数:**
- seed: 随机数种子
- fn: 使用随机数生成器的函数

**返回值:**
- T: 函数返回的结果

**使用示例:**

```go
// 生成可重现的随机序列
nums := pool.WithRandSeed(12345, func(rng *rand.Rand) []int {
    result := make([]int, 5)
    for i := range result {
        result[i] = rng.Intn(100)
    }
    return result
})
```

### func WithStr

```go
func WithStr(fn func(*strings.Builder)) string
```

WithStr 使用默认容量的字符串构建器执行函数，自动管理获取和归还

**参数:**
- fn: 使用字符串构建器的函数

**返回值:**
- string: 函数执行后构建的字符串结果

**使用示例:**

```go
result := pool.WithStr(func(buf *strings.Builder) {
    buf.WriteString("Hello")
    buf.WriteByte(' ')
    buf.WriteString("World")
})
```

### func WithStrCap

```go
func WithStrCap(cap int, fn func(*strings.Builder)) string
```

WithStrCap 使用指定容量的字符串构建器执行函数，自动管理获取和归还

**参数:**
- cap: 字符串构建器初始容量
- fn: 使用字符串构建器的函数

**返回值:**
- string: 函数执行后构建的字符串结果

**使用示例:**

```go
result := pool.WithStrCap(64, func(buf *strings.Builder) {
    buf.WriteString("Hello")
    buf.WriteByte(' ')
    buf.WriteString("World")
})
```

## Types

### type BufPool

```go
type BufPool struct {
	// Has unexported fields.
}
```

BufPool 字节缓冲区对象池，支持自定义配置

#### func NewBufPool

```go
func NewBufPool(defCap, maxCap int) *BufPool
```

NewBufPool 创建新的字节缓冲区对象池

**参数:**
- defCap: 默认字节缓冲区容量
- maxCap: 最大回收缓冲区容量，超过此容量的缓冲区不会被回收

**返回值:**
- *BufPool: 字节缓冲区对象池实例

#### func (*BufPool) Get

```go
func (bp *BufPool) Get() *bytes.Buffer
```

Get 获取默认容量的字节缓冲区

**返回:**
- *bytes.Buffer: 容量至少为默认容量的字节缓冲区

**说明:**
- 返回的字节缓冲区已经重置为空状态，可以直接使用
- 底层容量可能大于默认容量，来自对象池的复用缓冲区

#### func (*BufPool) GetCap

```go
func (bp *BufPool) GetCap(cap int) *bytes.Buffer
```

GetCap 获取指定容量的字节缓冲区

**参数:**
- cap: 需要的字节缓冲区容量

**返回:**
- *bytes.Buffer: 容量至少为capacity的字节缓冲区

**说明:**
- 返回的字节缓冲区已经重置为空状态，可以直接使用
- 底层容量可能大于capacity，来自对象池的复用缓冲区
- 如果capacity <= 0, 返回默认容量的缓冲区

#### func (*BufPool) Put

```go
func (bp *BufPool) Put(buf *bytes.Buffer)
```

Put 归还字节缓冲区到对象池

**参数:**
- buf: 要归还的字节缓冲区

#### func (*BufPool) With

```go
func (bp *BufPool) With(fn func(*bytes.Buffer)) []byte
```

With 使用默认容量的字节缓冲区执行函数，自动管理获取和归还

**参数:**
- fn: 使用字节缓冲区的函数

**返回值:**
- []byte: 函数执行后缓冲区的字节数据副本

**说明:**
- 自动从对象池获取默认容量的字节缓冲区
- 执行用户提供的函数
- 获取缓冲区字节数据的副本
- 自动归还字节缓冲区到对象池
- 即使函数发生panic也会正确归还资源

#### func (*BufPool) WithCap

```go
func (bp *BufPool) WithCap(cap int, fn func(*bytes.Buffer)) []byte
```

WithCap 使用指定容量的字节缓冲区执行函数，自动管理获取和归还

**参数:**
- cap: 字节缓冲区初始容量
- fn: 使用字节缓冲区的函数

**返回值:**
- []byte: 函数执行后缓冲区的字节数据副本

**说明:**
- 自动从对象池获取指定容量的字节缓冲区
- 执行用户提供的函数
- 获取缓冲区字节数据的副本
- 自动归还字节缓冲区到对象池
- 即使函数发生panic也会正确归还资源

### type BytePool

```go
type BytePool struct {
	// Has unexported fields.
}
```

BytePool 字节切片对象池, 支持自定义配置

#### func NewBytePool

```go
func NewBytePool(defCap, maxCap int) *BytePool
```

NewBytePool 创建新的字节切片对象池

**参数:**
- defCap: 默认缓冲区容量
- maxCap: 最大回收缓冲区容量, 超过此容量的缓冲区不会被回收

**返回值:**
- *BytePool: 字节切片对象池实例

#### func (*BytePool) Get

```go
func (bp *BytePool) Get() []byte
```

Get 获取默认容量的缓冲区

**返回:**
- []byte: 长度为默认容量, 容量至少为默认容量的缓冲区切片

**说明:**
- 返回的缓冲区长度等于默认容量, 可以直接使用
- 底层容量可能大于默认容量, 来自对象池的复用缓冲区

#### func (*BytePool) GetCap

```go
func (bp *BytePool) GetCap(size int) []byte
```

GetCap 获取指定容量的缓冲区

**参数:**
- size: 需要的缓冲区容量

**返回:**
- []byte: 长度为capacity, 容量至少为capacity的缓冲区切片

**说明:**
- 返回的缓冲区长度等于请求的capacity, 可以直接使用
- 底层容量可能大于capacity, 来自对象池的复用缓冲区
- 如果capacity <= 0, 使用默认容量

#### func (*BytePool) GetEmpty

```go
func (bp *BytePool) GetEmpty(size int) []byte
```

GetEmpty 获取指定容量的空缓冲区

**参数:**
- size: 指定容量要求

**返回:**
- []byte: 长度为0但容量至少为capacity的缓冲区切片

**说明:**
- 适用于需要使用append操作逐步构建数据的场景
- 避免频繁的内存重新分配
- 如果capacity <= 0, 使用默认容量

#### func (*BytePool) Put

```go
func (bp *BytePool) Put(buf []byte)
```

Put 归还缓冲区到对象池

**参数:**
- buf: 要归还的缓冲区

### type StrPool

```go
type StrPool struct {
	// Has unexported fields.
}
```

StrPool 字符串构建器对象池，支持自定义配置

#### func NewStrPool

```go
func NewStrPool(defCap, maxCap int) *StrPool
```

NewStrPool 创建新的字符串构建器对象池

**参数:**
- defCap: 默认字符串构建器容量
- maxCap: 最大回收构建器容量，超过此容量的构建器不会被回收

**返回值:**
- *StrPool: 字符串构建器对象池实例

#### func (*StrPool) Get

```go
func (sp *StrPool) Get() *strings.Builder
```

Get 获取默认容量的字符串构建器

**返回:**
- *strings.Builder: 容量至少为默认容量的字符串构建器

**说明:**
- 返回的字符串构建器已经重置为空状态，可以直接使用
- 底层容量可能大于默认容量，来自对象池的复用构建器

#### func (*StrPool) GetCap

```go
func (sp *StrPool) GetCap(cap int) *strings.Builder
```

GetCap 获取指定容量的字符串构建器

**参数:**
- cap: 需要的字符串构建器容量

**返回:**
- *strings.Builder: 容量至少为capacity的字符串构建器

**说明:**
- 返回的字符串构建器已经重置为空状态，可以直接使用
- 底层容量可能大于capacity，来自对象池的复用构建器
- 如果capacity <= 0, 返回默认容量的构建器

#### func (*StrPool) Put

```go
func (sp *StrPool) Put(buf *strings.Builder)
```

Put 归还字符串构建器到对象池

**参数:**
- b: 要归还的字符串构建器

**说明:**
- nil构建器不会被回收
- 容量不超过maxCap的构建器直接重置后归还
- 容量超过maxCap的构建器会创建一个新的小容量构建器进行归还（智能缩容）

#### func (*StrPool) With

```go
func (sp *StrPool) With(fn func(*strings.Builder)) string
```

With 使用默认容量的字符串构建器执行函数，自动管理获取和归还

**参数:**
- fn: 使用字符串构建器的函数

**返回值:**
- string: 函数执行后构建的字符串结果

**说明:**
- 自动从对象池获取默认容量的字符串构建器
- 执行用户提供的函数
- 获取构建的字符串结果
- 自动归还字符串构建器到对象池
- 即使函数发生panic也会正确归还资源

#### func (*StrPool) WithCap

```go
func (sp *StrPool) WithCap(cap int, fn func(*strings.Builder)) string
```

WithCap 使用指定容量的字符串构建器执行函数，自动管理获取和归还

**参数:**
- cap: 字符串构建器初始容量
- fn: 使用字符串构建器的函数

**返回值:**
- string: 函数执行后构建的字符串结果

**说明:**
- 自动从对象池获取指定容量的字符串构建器
- 执行用户提供的函数
- 获取构建的字符串结果
- 自动归还字符串构建器到对象池
- 即使函数发生panic也会正确归还资源
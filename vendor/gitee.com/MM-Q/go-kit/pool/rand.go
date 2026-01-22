// Package pool 提供随机数生成器对象池功能，通过对象池优化随机数生成性能。
//
// 随机数生成器对象池专门用于复用math/rand.Rand对象，
// 避免频繁创建随机数生成器的开销，特别适用于ID生成、测试数据生成等场景。
package pool

import (
	"math/rand"
	"sync"
	"time"
)

// 随机数生成器池
var randPool = sync.Pool{
	New: func() interface{} {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	},
}

// GetRand 从池中获取随机数生成器
//
// 返回值:
//   - *rand.Rand: 随机数生成器实例
//
// 说明:
//   - 返回的生成器已经初始化了随机种子
//   - 使用完毕后应调用PutRand归还
//   - 注意：返回的生成器不是线程安全的，不要在多个goroutine间共享
func GetRand() *rand.Rand {
	if r := randPool.Get(); r != nil {
		if gen, ok := r.(*rand.Rand); ok {
			// 重新设置随机种子，避免复用相同的随机序列
			gen.Seed(time.Now().UnixNano())
			return gen
		}
	}
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// PutRand 将随机数生成器归还到池中
//
// 参数:
//   - rng: 要归还的随机数生成器
func PutRand(rng *rand.Rand) {
	if rng != nil {
		randPool.Put(rng)
	}
}

// GetRandWithSeed 获取指定种子的随机数生成器
//
// 参数:
//   - seed: 随机数种子
//
// 返回值:
//   - *rand.Rand: 随机数生成器实例
//
// 说明:
//   - 返回的生成器使用指定的种子初始化
//   - 适用于需要可重现随机序列的场景
func GetRandWithSeed(seed int64) *rand.Rand {
	if r := randPool.Get(); r != nil {
		if gen, ok := r.(*rand.Rand); ok {
			gen.Seed(seed)
			return gen
		}
	}
	return rand.New(rand.NewSource(seed))
}

// WithRand 使用随机数生成器执行函数，自动管理获取和归还
//
// 参数:
//   - fn: 使用随机数生成器的函数
//
// 返回值:
//   - T: 函数返回的结果
//
// 使用示例:
//
//	// 生成随机整数
//	num := pool.WithRand(func(rng *rand.Rand) int {
//	    return rng.Intn(100)
//	})
//
//	// 生成随机字符串
//	str := pool.WithRand(func(rng *rand.Rand) string {
//	    return fmt.Sprintf("id_%d", rng.Int63())
//	})
func WithRand[T any](fn func(*rand.Rand) T) T {
	rng := GetRand()
	defer PutRand(rng)
	return fn(rng)
}

// WithRandSeed 使用指定种子的随机数生成器执行函数，自动管理获取和归还
//
// 参数:
//   - seed: 随机数种子
//   - fn: 使用随机数生成器的函数
//
// 返回值:
//   - T: 函数返回的结果
//
// 使用示例:
//
//	// 生成可重现的随机序列
//	nums := pool.WithRandSeed(12345, func(rng *rand.Rand) []int {
//	    result := make([]int, 5)
//	    for i := range result {
//	        result[i] = rng.Intn(100)
//	    }
//	    return result
//	})
func WithRandSeed[T any](seed int64, fn func(*rand.Rand) T) T {
	rng := GetRandWithSeed(seed)
	defer PutRand(rng)
	return fn(rng)
}

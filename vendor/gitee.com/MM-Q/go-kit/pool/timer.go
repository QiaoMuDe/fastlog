// Package pool 提供Timer对象池功能，通过对象池优化定时器使用。
//
// Timer对象池专门用于复用time.Timer对象，避免频繁创建和销毁定时器的开销。
// Timer的创建成本相对较高，特别是在高并发场景下，复用可以显著提升性能。
package pool

import (
	"sync"
	"time"
)

// 定时器池
var timerPool = sync.Pool{
	New: func() interface{} {
		timer := time.NewTimer(time.Hour)
		if !timer.Stop() {
			<-timer.C
		}
		return timer
	},
}

// GetTimer 从池中获取定时器并设置超时时间
//
// 参数:
//   - duration: 定时器超时时间
//
// 返回值:
//   - *time.Timer: 已设置超时时间的定时器
//
// 说明:
//   - 返回的定时器已经启动，会在指定时间后触发
//   - 适用于超时控制场景，定时器会在指定时间后自动触发
//   - 使用完毕后应调用PutTimer归还
func GetTimer(duration time.Duration) *time.Timer {
	if t := timerPool.Get(); t != nil {
		if timer, ok := t.(*time.Timer); ok {
			timer.Reset(duration)
			return timer
		}
	}
	return time.NewTimer(duration)
}

// PutTimer 将定时器归还到池中
//
// 参数:
//   - timer: 要归还的定时器
//
// 说明:
//   - 该函数会自动停止定时器并清理状态
//   - 归还后的定时器会被重置，可以安全复用
func PutTimer(timer *time.Timer) {
	if timer == nil {
		return
	}

	// 停止定时器并清理channel
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}

	timerPool.Put(timer)
}

// GetTimerEmpty 从池中获取未启动的定时器
//
// 返回值:
//   - *time.Timer: 未启动的定时器，需要手动调用Reset设置时间
//
// 说明:
//   - 适用于需要手动控制定时器启动和停止的场景
//   - 定时器处于停止状态，不会自动触发
//   - 获取后需要调用timer.Reset(duration)启动
//   - 使用完毕后应调用PutTimer归还
func GetTimerEmpty() *time.Timer {
	if t := timerPool.Get(); t != nil {
		if timer, ok := t.(*time.Timer); ok {
			return timer
		}
	}
	timer := time.NewTimer(time.Hour)
	if !timer.Stop() {
		<-timer.C
	}
	return timer
}

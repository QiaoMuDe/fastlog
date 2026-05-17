package fastlog

import (
	"testing"
	"time"
)

func TestSamplerInitialAllow(t *testing.T) {
	s := NewSampler(time.Minute, 3, 10)
	msg := "test message"

	for i := 1; i <= 3; i++ {
		if !s.Allow(INFO, msg) {
			t.Errorf("第 %d 次调用 Allow() = false, 在 initial 窗口内应为 true", i)
		}
	}

	if s.Allow(INFO, msg) {
		t.Errorf("第 4 次调用 Allow() = true, 超出 initial 窗口应为 false")
	}
}

func TestSamplerThereafter(t *testing.T) {
	s := NewSampler(time.Minute, 3, 10)
	msg := "test message"

	// 前 3 条放行
	for i := 0; i < 3; i++ {
		s.Allow(INFO, msg)
	}

	// 第 4~12 条抑制 (initial + thereafter - 1 = 3 + 10 - 1 = 12)
	for i := 4; i <= 12; i++ {
		if s.Allow(INFO, msg) {
			t.Errorf("第 %d 次 Allow() 在 thereafter 间隔内应为 false", i)
		}
	}

	// 第 13 条放行 (initial + thereafter = 13)
	if !s.Allow(INFO, msg) {
		t.Errorf("第 13 次 Allow() 应为 true (每 10 条放行 1 条) ")
	}
}

func TestSamplerWindowReset(t *testing.T) {
	s := NewSampler(50*time.Millisecond, 3, 10)
	msg := "test message"

	// 填满 initial
	for i := 0; i < 3; i++ {
		s.Allow(INFO, msg)
	}

	// 超出 initial 后应该被抑制
	if s.Allow(INFO, msg) {
		t.Errorf("初始窗口内第4条应为 false")
	}

	// 等待窗口过期
	time.Sleep(60 * time.Millisecond)

	// 窗口重置后应重新从 initial 计数
	for i := 1; i <= 3; i++ {
		if !s.Allow(INFO, msg) {
			t.Errorf("窗口重置后第 %d 次 Allow() = false, 应为 true", i)
		}
	}
}

func TestSamplerIndependentLevel(t *testing.T) {
	s := NewSampler(time.Minute, 2, 5)
	msg := "test message"

	// 填满 INFO 的 initial
	s.Allow(INFO, msg)
	s.Allow(INFO, msg)

	// INFO 第 3 条被抑制
	if s.Allow(INFO, msg) {
		t.Errorf("INFO 超出 initial 应为 false")
	}

	// DEBUG 应该不受影响
	for i := 1; i <= 2; i++ {
		if !s.Allow(DEBUG, msg) {
			t.Errorf("DEBUG 第 %d 次 Allow() = false, 不应受 INFO 影响", i)
		}
	}
}

func TestSamplerThereafterZero(t *testing.T) {
	s := NewSampler(time.Minute, 3, 0)
	msg := "test message"

	for i := 0; i < 3; i++ {
		s.Allow(INFO, msg)
	}

	// thereafter=0 代表 after initial 永久抑制
	for i := 0; i < 10; i++ {
		if s.Allow(INFO, msg) {
			t.Errorf("thereafter=0 时 initial 后应永久抑制, 但第 %d 次返回 true", i+4)
		}
	}
}

func TestSamplerDifferentMessages(t *testing.T) {
	s := NewSampler(time.Minute, 1, 0)
	msg1 := "message one"
	msg2 := "message two"

	// msg1 第 1 条放行, 第 2 条抑制 (thereafter=0)
	if !s.Allow(INFO, msg1) {
		t.Errorf("msg1 第1条应为 true")
	}
	if s.Allow(INFO, msg1) {
		t.Errorf("msg1 第2条应为 false")
	}

	// msg2 第 1 条应该放行 (不同消息不同桶)
	if !s.Allow(INFO, msg2) {
		t.Errorf("不同消息 msg2 第1条应为 true")
	}
}

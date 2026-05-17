package fastlog

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

type mockWriteCloser struct {
	*bytes.Buffer
	closeErr error
	closed   bool
}

func (m *mockWriteCloser) Close() error {
	m.closed = true
	return m.closeErr
}

func newMock() *mockWriteCloser {
	return &mockWriteCloser{Buffer: &bytes.Buffer{}}
}

func TestColorWriterNoColor(t *testing.T) {
	buf := &bytes.Buffer{}
	cw := &ColorWriter{w: buf, NoColor: true}

	data := []byte("2026-01-15T10:30:45Z | INFO | hello\n")
	n, err := cw.Write(data)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != len(data) {
		t.Errorf("Write() = %d, want %d", n, len(data))
	}
	if !bytes.Equal(buf.Bytes(), data) {
		t.Errorf("NoColor output = %q, want %q", buf.Bytes(), data)
	}
}

func TestColorWriterColorPassthrough(t *testing.T) {
	buf := &bytes.Buffer{}
	cw := &ColorWriter{w: buf, NoColor: false}

	// 包含 INFO 关键字的数据, 确保内容完整通过
	data := []byte("2026-01-15T10:30:45Z | INFO | hello\n")
	_, err := cw.Write(data)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("hello")) {
		t.Errorf("Output should contain original content, got: %q", buf.Bytes())
	}
}

func TestColorWriterUnknownLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	cw := &ColorWriter{w: buf, NoColor: false}

	// 不包含任何级别关键字的数据
	data := []byte("some random log data without level\n")
	_, err := cw.Write(data)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if !bytes.Equal(buf.Bytes(), data) {
		t.Errorf("Unknown level output = %q, want %q", buf.Bytes(), data)
	}
}

func TestColorWriterClose(t *testing.T) {
	cw := NewColorWriter(false)
	if err := cw.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestDetectLevel(t *testing.T) {
	tests := []struct {
		name string
		data string
		want Level
	}{
		{"DEBUG", "[DEBUG] msg", DEBUG},
		{"INFO", "[INFO] msg", INFO},
		{"WARN", "[WARN] msg", WARN},
		{"ERROR", "[ERROR] msg", ERROR},
		{"FATAL", "[FATAL] msg", FATAL},
		{"PANIC", "[PANIC] msg", PANIC},
		// 高优先级优先 (PANIC > DEBUG)
		{"PANIC before DEBUG", "PANIC and DEBUG", PANIC},
		{"no match", "some random text", INFO},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cw := &ColorWriter{w: &bytes.Buffer{}}
			got := cw.detectLevel([]byte(tt.data))
			if got != tt.want {
				t.Errorf("detectLevel(%q) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestMultiWriterWrite(t *testing.T) {
	m1 := newMock()
	m2 := newMock()
	mw := NewMultiWriter(m1, m2)

	data := []byte("test data\n")
	n, err := mw.Write(data)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != len(data) {
		t.Errorf("Write() = %d, want %d", n, len(data))
	}

	if !bytes.Equal(m1.Bytes(), data) {
		t.Errorf("writer1 got %q, want %q", m1.Bytes(), data)
	}
	if !bytes.Equal(m2.Bytes(), data) {
		t.Errorf("writer2 got %q, want %q", m2.Bytes(), data)
	}
}

type failWriter struct {
	io.WriteCloser
}

func (f *failWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}

func (f *failWriter) Close() error {
	return nil
}

func TestMultiWriterWriteError(t *testing.T) {
	m1 := newMock()
	fail := &failWriter{}
	mw := NewMultiWriter(m1, fail)

	_, err := mw.Write([]byte("test"))
	if err == nil {
		t.Errorf("Write() should return error when a writer fails")
	}
}

func TestMultiWriterClose(t *testing.T) {
	m1 := newMock()
	m2 := newMock()
	mw := NewMultiWriter(m1, m2)

	if err := mw.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !m1.closed {
		t.Errorf("writer1 not closed")
	}
	if !m2.closed {
		t.Errorf("writer2 not closed")
	}
}

type failCloser struct {
	io.WriteCloser
	closed bool
}

func (f *failCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (f *failCloser) Close() error {
	f.closed = true
	return errors.New("close failed")
}

func TestMultiWriterCloseError(t *testing.T) {
	f1 := &failCloser{}
	f2 := &failCloser{}
	mw := NewMultiWriter(f1, f2)

	err := mw.Close()
	if err == nil {
		t.Errorf("Close() should return error when writers fail")
	}
	if !f1.closed || !f2.closed {
		t.Errorf("All writers should be closed even on error")
	}
	// errors.Join 合并多个错误
	if !strings.Contains(err.Error(), "close failed") {
		t.Errorf("Close() error should contain 'close failed', got: %v", err)
	}
}

func TestMultiWriterEmpty(t *testing.T) {
	mw := NewMultiWriter()

	n, err := mw.Write([]byte("test"))
	if err != nil {
		t.Errorf("Empty MultiWriter Write() error = %v", err)
	}
	if n != 4 {
		t.Errorf("Empty MultiWriter Write() = %d, want 4", n)
	}

	if err := mw.Close(); err != nil {
		t.Errorf("Empty MultiWriter Close() error = %v", err)
	}
}

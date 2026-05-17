package fastlog

import (
	"testing"
	"time"

	"gitee.com/MM-Q/comprx"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig("logs/app.log")

	if cfg.Level != INFO {
		t.Errorf("NewConfig().Level = %v, want INFO", cfg.Level)
	}
	if cfg.Formatter == nil {
		t.Errorf("NewConfig().Formatter should not be nil")
	}
	if !cfg.OutputConsole {
		t.Errorf("NewConfig().OutputConsole should be true")
	}
	if !cfg.OutputFile {
		t.Errorf("NewConfig().OutputFile should be true")
	}
	if cfg.LogPath != "logs/app.log" {
		t.Errorf("NewConfig().LogPath = %q, want %q", cfg.LogPath, "logs/app.log")
	}
	if cfg.SamplerTick != DefaultSamplerTick {
		t.Errorf("NewConfig().SamplerTick = %v, want %v", cfg.SamplerTick, DefaultSamplerTick)
	}
	if cfg.SamplerInitial != DefaultSamplerInitial {
		t.Errorf("NewConfig().SamplerInitial = %d, want %d", cfg.SamplerInitial, DefaultSamplerInitial)
	}
	if cfg.SamplerThereafter != DefaultSamplerThereafter {
		t.Errorf("NewConfig().SamplerThereafter = %d, want %d", cfg.SamplerThereafter, DefaultSamplerThereafter)
	}
}

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.LogPath != "logs/app.log" {
		t.Errorf("Default().LogPath = %q, want %q", cfg.LogPath, "logs/app.log")
	}
}

func TestDev(t *testing.T) {
	cfg := Dev("dev.log")
	if cfg.Level != DEBUG {
		t.Errorf("Dev().Level = %v, want DEBUG", cfg.Level)
	}
	if !cfg.Caller {
		t.Errorf("Dev().Caller should be true")
	}
	if cfg.SamplerTick != 0 {
		t.Errorf("Dev().SamplerTick should be 0 (disabled), got %v", cfg.SamplerTick)
	}
	if cfg.MaxSize != 10 {
		t.Errorf("Dev().MaxSize = %d, want 10", cfg.MaxSize)
	}
	if cfg.Compress {
		t.Errorf("Dev().Compress should be false")
	}
	if cfg.RotateByDay {
		t.Errorf("Dev().RotateByDay should be false")
	}
}

func TestProd(t *testing.T) {
	cfg := Prod("prod.log")
	if cfg.Level != WARN {
		t.Errorf("Prod().Level = %v, want WARN", cfg.Level)
	}
	if cfg.OutputConsole {
		t.Errorf("Prod().OutputConsole should be false")
	}
	if !cfg.Async {
		t.Errorf("Prod().Async should be true")
	}
	if cfg.MaxSize != 100 {
		t.Errorf("Prod().MaxSize = %d, want 100", cfg.MaxSize)
	}
	if cfg.MaxFiles != 14 {
		t.Errorf("Prod().MaxFiles = %d, want 14", cfg.MaxFiles)
	}
	if cfg.MaxAge != 14 {
		t.Errorf("Prod().MaxAge = %d, want 14", cfg.MaxAge)
	}
	if !cfg.Compress {
		t.Errorf("Prod().Compress should be true")
	}
}

func TestConsole(t *testing.T) {
	cfg := Console()
	if cfg.Level != DEBUG {
		t.Errorf("Console().Level = %v, want DEBUG", cfg.Level)
	}
	if cfg.OutputFile {
		t.Errorf("Console().OutputFile should be false")
	}
	if cfg.LogPath != "" {
		t.Errorf("Console().LogPath should be empty, got %q", cfg.LogPath)
	}
	if !cfg.OutputConsole {
		t.Errorf("Console().OutputConsole should be true")
	}
}

func TestDocker(t *testing.T) {
	cfg := Docker()
	if cfg.Level != WARN {
		t.Errorf("Docker().Level = %v, want WARN", cfg.Level)
	}
	if _, ok := cfg.Formatter.(JSON); !ok {
		t.Errorf("Docker().Formatter should be JSON{}")
	}
	if cfg.OutputFile {
		t.Errorf("Docker().OutputFile should be false")
	}
	if !cfg.OutputConsole {
		t.Errorf("Docker().OutputConsole should be true")
	}
	if cfg.LogPath != "" {
		t.Errorf("Docker().LogPath should be empty")
	}
}

func TestValidate(t *testing.T) {
	t.Run("valid console only", func(t *testing.T) {
		cfg := &Config{OutputConsole: true}
		if err := cfg.Validate(); err != nil {
			t.Errorf("Validate() error = %v, want nil", err)
		}
	})

	t.Run("valid file only", func(t *testing.T) {
		cfg := &Config{OutputFile: true, LogPath: "/tmp/test.log"}
		if err := cfg.Validate(); err != nil {
			t.Errorf("Validate() error = %v, want nil", err)
		}
	})

	t.Run("valid both", func(t *testing.T) {
		cfg := &Config{OutputConsole: true, OutputFile: true, LogPath: "/tmp/test.log"}
		if err := cfg.Validate(); err != nil {
			t.Errorf("Validate() error = %v, want nil", err)
		}
	})

	t.Run("no output", func(t *testing.T) {
		cfg := &Config{}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate() should error when no output set")
		}
	})

	t.Run("file output without path", func(t *testing.T) {
		cfg := &Config{OutputFile: true}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate() should error when OutputFile=true but no LogPath")
		}
	})

	t.Run("negative MaxSize", func(t *testing.T) {
		cfg := &Config{OutputFile: true, LogPath: "/tmp/test.log", MaxSize: -1}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate() should error when MaxSize < 0")
		}
	})

	t.Run("negative MaxFiles", func(t *testing.T) {
		cfg := &Config{OutputFile: true, LogPath: "/tmp/test.log", MaxFiles: -1}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate() should error when MaxFiles < 0")
		}
	})

	t.Run("negative MaxAge", func(t *testing.T) {
		cfg := &Config{OutputFile: true, LogPath: "/tmp/test.log", MaxAge: -1}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate() should error when MaxAge < 0")
		}
	})

	t.Run("sampler negative initial", func(t *testing.T) {
		cfg := &Config{OutputConsole: true, SamplerTick: time.Second, SamplerInitial: -1, SamplerThereafter: 10}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate() should error when SamplerInitial < 0")
		}
	})

	t.Run("sampler negative thereafter", func(t *testing.T) {
		cfg := &Config{OutputConsole: true, SamplerTick: time.Second, SamplerInitial: 3, SamplerThereafter: -1}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate() should error when SamplerThereafter < 0")
		}
	})

	t.Run("buffer size too small", func(t *testing.T) {
		cfg := &Config{OutputConsole: true, MaxBufferSize: 1000}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate() should error when MaxBufferSize < 64KB")
		}
	})

	t.Run("sync interval too short", func(t *testing.T) {
		cfg := &Config{OutputConsole: true, SyncInterval: 100 * time.Millisecond}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate() should error when SyncInterval < 500ms")
		}
	})

	t.Run("negative MaxBufferSize", func(t *testing.T) {
		cfg := &Config{OutputConsole: true, MaxBufferSize: -1}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate() should error when MaxBufferSize < 0")
		}
	})

	t.Run("negative SyncInterval", func(t *testing.T) {
		cfg := &Config{OutputConsole: true, SyncInterval: -1}
		if err := cfg.Validate(); err == nil {
			t.Errorf("Validate() should error when SyncInterval < 0")
		}
	})
}

func TestClone(t *testing.T) {
	t.Run("basic clone", func(t *testing.T) {
		cfg := NewConfig("test.log")
		clone := cfg.Clone()

		if clone.Level != cfg.Level {
			t.Errorf("Clone().Level = %v, want %v", clone.Level, cfg.Level)
		}
		if clone.LogPath != cfg.LogPath {
			t.Errorf("Clone().LogPath = %q, want %q", clone.LogPath, cfg.LogPath)
		}
	})

	t.Run("clone independence", func(t *testing.T) {
		cfg := &Config{
			Level:         INFO,
			OutputConsole: true,
			Fields:        []Field{String("app", "myapp")},
		}
		clone := cfg.Clone()

		// 修改原始配置不影响克隆
		cfg.Level = DEBUG
		cfg.Fields[0] = String("app", "changed")

		if clone.Level == DEBUG {
			t.Errorf("Clone should not be affected by original Level change")
		}
		if clone.Fields[0].Value() == "changed" {
			t.Errorf("Clone Fields should not be affected by original Fields change")
		}
	})

	t.Run("clone nil fields", func(t *testing.T) {
		cfg := &Config{OutputConsole: true}
		clone := cfg.Clone()
		// 原字段为 nil 时, 克隆后不应分配新切片
		if clone.Fields != nil {
			t.Errorf("Clone().Fields should be nil when original has nil fields")
		}
	})
}

func TestConfigNewSampler(t *testing.T) {
	t.Run("returns sampler when tick > 0", func(t *testing.T) {
		cfg := &Config{SamplerTick: time.Second, SamplerInitial: 3, SamplerThereafter: 10}
		s := cfg.NewSampler()
		if s == nil {
			t.Errorf("NewSampler() should return non-nil when SamplerTick > 0")
		}
	})

	t.Run("returns nil when tick <= 0", func(t *testing.T) {
		cfg := &Config{SamplerTick: 0}
		s := cfg.NewSampler()
		if s != nil {
			t.Errorf("NewSampler() should return nil when SamplerTick = 0")
		}
	})

	t.Run("returns nil when tick negative", func(t *testing.T) {
		cfg := &Config{SamplerTick: -time.Second}
		s := cfg.NewSampler()
		if s != nil {
			t.Errorf("NewSampler() should return nil when SamplerTick < 0")
		}
	})
}

func TestConfigNewWriter(t *testing.T) {
	t.Run("console only", func(t *testing.T) {
		cfg := &Config{OutputConsole: true}
		w := cfg.NewWriter()
		if w == nil {
			t.Errorf("NewWriter() should return non-nil for console only")
		}
	})

	t.Run("file only", func(t *testing.T) {
		cfg := &Config{OutputFile: true, LogPath: "test.log"}
		w := cfg.NewWriter()
		if w == nil {
			t.Errorf("NewWriter() should return non-nil for file only")
		}
	})

	t.Run("both", func(t *testing.T) {
		cfg := &Config{OutputConsole: true, OutputFile: true, LogPath: "test.log"}
		w := cfg.NewWriter()
		if w == nil {
			t.Errorf("NewWriter() should return non-nil for both outputs")
		}
	})

	t.Run("no output", func(t *testing.T) {
		cfg := &Config{}
		w := cfg.NewWriter()
		if w != nil {
			t.Errorf("NewWriter() should return nil when no output set")
		}
	})
}

func TestConfigCompressType(t *testing.T) {
	cfg := NewConfig("test.log")
	if cfg.CompressType != comprx.CompressTypeGz {
		t.Errorf("Default CompressType = %v, want %v", cfg.CompressType, comprx.CompressTypeGz)
	}
}

func TestValidateAllValid(t *testing.T) {
	// NewConfig 创建的配置应该能通过验证
	cfg := NewConfig("test.log")
	if err := cfg.Validate(); err != nil {
		t.Errorf("NewConfig().Validate() error = %v, want nil", err)
	}
}

func TestValidateZeroValues(t *testing.T) {
	// 零值 MaxSize/MaxFiles/MaxAge 是合法的（表示不限制）
	cfg := &Config{
		OutputConsole: true,
		OutputFile:    true,
		LogPath:       "test.log",
		MaxSize:       0,
		MaxFiles:      0,
		MaxAge:        0,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("Validate() with zero MaxSize/MaxFiles/MaxAge error = %v, want nil", err)
	}
}

package fastlog

import (
	"testing"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		name string
		l    Level
		want string
	}{
		{"DEBUG", DEBUG, "DEBUG"},
		{"INFO", INFO, "INFO"},
		{"WARN", WARN, "WARN"},
		{"ERROR", ERROR, "ERROR"},
		{"FATAL", FATAL, "FATAL"},
		{"PANIC", PANIC, "PANIC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("Level.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLevelEnabled(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		lvl   Level
		want  bool
	}{
		{"INFO enabled INFO", INFO, INFO, true},
		{"INFO enabled WARN", INFO, WARN, true},
		{"INFO enabled ERROR", INFO, ERROR, true},
		{"INFO enabled FATAL", INFO, FATAL, true},
		{"INFO enabled PANIC", INFO, PANIC, true},
		{"INFO not enabled DEBUG", INFO, DEBUG, false},

		{"DEBUG enabled all", DEBUG, PANIC, true},
		{"PANIC not enabled FATAL", PANIC, FATAL, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.Enabled(tt.lvl); got != tt.want {
				t.Errorf("Level(%d).Enabled(%d) = %v, want %v", tt.level, tt.lvl, got, tt.want)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    Level
		wantErr bool
	}{
		{"lowercase", "debug", DEBUG, false},
		{"uppercase", "DEBUG", DEBUG, false},
		{"mixed case", "Debug", DEBUG, false},
		{"info", "INFO", INFO, false},
		{"warn", "WARN", WARN, false},
		{"error", "ERROR", ERROR, false},
		{"fatal", "FATAL", FATAL, false},
		{"panic", "PANIC", PANIC, false},
		{"unknown string", "unknown", INFO, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLevel(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLevel(%q) error = %v, wantErr %v", tt.s, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestAllLevels(t *testing.T) {
	levels := AllLevels()
	want := []Level{DEBUG, INFO, WARN, ERROR, FATAL, PANIC}
	if len(levels) != len(want) {
		t.Errorf("AllLevels() length = %d, want %d", len(levels), len(want))
		return
	}
	for i := range want {
		if levels[i] != want[i] {
			t.Errorf("AllLevels()[%d] = %v, want %v", i, levels[i], want[i])
		}
	}
}

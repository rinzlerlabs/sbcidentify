package test

import (
	"log/slog"
	"testing"

	"github.com/rinzlerlabs/sbcidentify"
	"github.com/rinzlerlabs/sbcidentify/boardtype"
)

func TestShouldSkip(t *testing.T) {
	sbcidentify.SetLogLevel(slog.LevelDebug)
	tests := []struct {
		name       string
		setup      func() *test
		shouldSkip bool
		msg        string
	}{
		{
			name: "No requirements",
			setup: func() *test {
				return Test()
			},
			shouldSkip: false,
			msg:        "Test has no requirements",
		},
		{
			name: "Requires SBC but no SBC present",
			setup: func() *test {
				return Test().RequiresSbc()
			},
			shouldSkip: true,
			msg:        "Test requires physical SBC",
		},
		{
			name: "Requires specific board type but different board type present",
			setup: func() *test {
				return Test().RequiresBoardType(boardtype.RaspberryPi3B)
			},
			shouldSkip: true,
			msg:        "Test requires board type",
		},
		{
			name: "Requires root but not running as root",
			setup: func() *test {
				return Test().RequiresRoot()
			},
			shouldSkip: true,
			msg:        "Test requires root",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup().ShouldSkip(t)
			if tt.shouldSkip {
				t.Errorf("%v should have skipped, but did not", tt.name)
			}
		})
	}
}

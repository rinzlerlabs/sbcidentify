package test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/rinzlerlabs/sbcidentify/boardtype"
)

func TestShouldSkip(t *testing.T) {
	prev := slog.Default()
	t.Cleanup(func() { slog.SetDefault(prev) })
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
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

// Just making sure that errors.Join works as expected with nil values, since we use it in GetBoardType.
func TestErrorsJoinWithNil(t *testing.T) {
	err := errors.Join(nil, nil)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	var ErrUnknownBoard error = errors.New("unknown board")
	err = errors.Join(nil, ErrUnknownBoard)
	if errors.Is(err, ErrUnknownBoard) == false {
		t.Errorf("Expected ErrUnknownBoard, got %v", err)
	}

	err = errors.Join(ErrUnknownBoard, nil)
	if errors.Is(err, ErrUnknownBoard) == false {
		t.Errorf("Expected ErrUnknownBoard, got %v", err)
	}
}

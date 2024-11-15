package raspberrypi

import (
	"log/slog"
	"os"
	"os/exec"
	"testing"

	"github.com/thegreatco/sbcidentify/identifier"
)

func setup(t *testing.T) *slog.Logger {
	t.Helper()
	execLookPath = exec.LookPath
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	id := NewRaspberryPiIdentifier(logger)
	board, err := id.GetBoardType()
	if err != nil && err != identifier.ErrCannotIdentifyBoard {
		t.Fatalf("GetBoardType() failed: %v", err)
	}
	if board.GetManufacturer() != "Raspberry Pi" {
		t.Skip("Not a Raspberry Pi")
	}
	return logger
}

func TestGetInstalledRAM(t *testing.T) {
	logger := setup(t)
	ram, err := getInstalledRAM(logger)
	if err != nil {
		t.Fatalf("getInstalledRAM() failed: %v", err)
	}
	t.Logf("RAM: %dMB", ram)

	execLookPath = func(string) (string, error) {
		return "", exec.ErrNotFound
	}
	_, err = getInstalledRAM(logger)
	if err != ErrVcgencmdNotFound {
		t.Fatalf("getInstalledRAM() returned error %v, expected %v", err, ErrVcgencmdNotFound)
	}
}

func TestParseVcgencmdMemoryOutput(t *testing.T) {
	logger := setup(t)

	tests := []struct {
		input  string
		output int
		err    error
	}{
		{"total_mem", 0, ErrInvalidMeminfo},
		{"total_mem=", 0, ErrInvalidMeminfo},
		{"total_mem=foo", 0, ErrInvalidMeminfo},
		{"", 0, ErrInvalidMeminfo},
		{"total_mem=2048MB", 0, ErrInvalidMeminfo},
		{"total_mem=1024", 1024, nil},
		{"total_mem=1024\n", 1024, nil},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			ram, err := parseVcgencmdMemoryOutput(logger, test.input)
			if err != test.err {
				t.Fatalf("parseVcgencmdMemoryOutput() returned error %v, expected %v", err, test.err)
			}
			if ram != test.output {
				t.Fatalf("parseVcgencmdMemoryOutput() returned %d, expected %d", ram, test.output)
			}
		})
	}
}

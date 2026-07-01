package raspberrypi

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"testing"

	"github.com/rinzlerlabs/sbcidentify/boardtype"
	"github.com/rinzlerlabs/sbcidentify/identifier"
)

func setup(t *testing.T) (*slog.Logger, identifier.BoardIdentifier) {
	t.Helper()
	execLookPath = exec.LookPath
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	id := NewRaspberryPiIdentifier(logger)
	board, err := id.GetBoardType()
	if err == ErrCannotIdentifyBoard || err == identifier.ErrCannotIdentifyBoard {
		t.Skip("Not a Raspberry Pi")
	}
	if err != nil {
		t.Fatalf("GetBoardType() failed: %v", err)
	}
	if board.GetManufacturer() != "Raspberry Pi" {
		t.Skip("Not a Raspberry Pi")
	}
	return logger, id
}

func TestGetInstalledRAM(t *testing.T) {
	logger, _ := setup(t)
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
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

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

func TestIsBoardType(t *testing.T) {
	tests := []struct {
		Want     boardtype.SBC
		Have     boardtype.SBC
		expected bool
	}{
		// Root hierarchy
		{boardtype.RaspberryPi, boardtype.RaspberryPi4B8GB, true},
		{boardtype.RaspberryPi, boardtype.RaspberryPi5B8GB, true},
		// Pi 3 family
		{boardtype.RaspberryPi3B, boardtype.RaspberryPi4B, false},
		{boardtype.RaspberryPi3B, boardtype.RaspberryPi3BPlus, true},
		{boardtype.RaspberryPi3BPlus, boardtype.RaspberryPi3B, false},
		// Pi 4 family
		{boardtype.RaspberryPi4, boardtype.RaspberryPi4B8GB, true},
		{boardtype.RaspberryPi4B, boardtype.RaspberryPi4B8GB, true},
		{boardtype.RaspberryPi4B8GB, boardtype.RaspberryPi4B, false},
		{boardtype.RaspberryPi4B4GB, boardtype.RaspberryPi4B8GB, false},
		// Pi 5 family
		{boardtype.RaspberryPi5B, boardtype.RaspberryPi5B4GB, true},
		{boardtype.RaspberryPi5B4GB, boardtype.RaspberryPi5B, false},
		// CM5 Lite hierarchy
		{boardtype.RaspberryPi5B, boardtype.RaspberryPiCM5Lite4GB, true},
		{boardtype.RaspberryPiCM5Lite, boardtype.RaspberryPiCM5Lite4GB, true},
		{boardtype.RaspberryPiCM5Lite4GB, boardtype.RaspberryPiCM5Lite, false},
		// CM5 Lite is not CM5
		{boardtype.RaspberryPiCM5Lite, boardtype.RaspberryPiCM54GB, false},
		{boardtype.RaspberryPiCM54GB, boardtype.RaspberryPiCM5Lite4GB, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Want_%v_Have_%v", test.Want.GetPrettyName(), test.Have.GetPrettyName()), func(t *testing.T) {
			if test.Have.IsBoardType(test.Want) != test.expected {
				t.Fatalf("IsBoardType() returned %v, expected %v", !test.expected, test.expected)
			}
		})
	}
}

func TestGetSOC(t *testing.T) {
	tests := []struct {
		board    boardtype.SBC
		expected string
	}{
		{boardtype.RaspberryPi3B, "BCM2837"},
		{boardtype.RaspberryPi3BPlus, "BCM2837"},
		{boardtype.RaspberryPi4B8GB, "BCM2711"},
		{boardtype.RaspberryPiCM48GB, "BCM2711"},
		{boardtype.RaspberryPi4400, "BCM2711"},
		{boardtype.RaspberryPi5B8GB, "BCM2712"},
		{boardtype.RaspberryPiCM58GB, "BCM2712"},
		{boardtype.RaspberryPiCM5Lite4GB, "BCM2712"},
		{boardtype.RaspberryPiCM5Lite, "BCM2712"},
		{boardtype.RaspberryPi, ""},
	}
	for _, test := range tests {
		t.Run(test.board.GetPrettyName(), func(t *testing.T) {
			if test.board.GetSOC() != test.expected {
				t.Fatalf("GetSOC() returned %q, expected %q", test.board.GetSOC(), test.expected)
			}
		})
	}
}

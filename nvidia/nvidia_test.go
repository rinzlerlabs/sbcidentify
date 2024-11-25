package nvidia

import (
	"fmt"
	"os"
	"testing"

	"log/slog"

	"github.com/stretchr/testify/assert"
	"github.com/thegreatco/sbcidentify/boardtype"
	"github.com/thegreatco/sbcidentify/identifier"
)

func setup(t *testing.T) *slog.Logger {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	id := NewNvidiaIdentifier(logger)
	board, err := id.GetBoardType()
	if err != nil && err != identifier.ErrCannotIdentifyBoard {
		t.Fatalf("GetBoardType() failed: %v", err)
	}
	if board.GetManufacturer() != "NVIDIA" {
		t.Skip("Not an NVIDIA board")
	}
	return logger
}

func TestParseModuleName(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	f, e := getModuleNameFromDtsFilename(logger, "/dvs/git/dirty/git-master_linux/kernel/kernel-5.10/arch/arm64/boot/dts/../../../../../../hardware/nvidia/platform/t23x/p3768/kernel-dts/tegra234-p3767-0003-p3768-0000-a0.dts")
	if e != nil {
		t.Fatalf("getModuleNameFromDtsFilename() failed: %v", e)
	}
	if f != "tegra234-p3767-0003-p3768-0000-a0" {
		t.Fatalf("getModuleNameFromDtsFilename() returned %s, expected tegra234-p3767-0003-p3768-0000-a0", f)
	}
}

func TestParseModelNameFromModuleName(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	f, e := getModuleModelFromModuleName(logger, "tegra234-p3767-0003-p3768-0000-a0")
	assert.NoError(t, e)
	assert.Equal(t, "p3767-0003", f)
}

func TestIsBoardType(t *testing.T) {
	tests := []struct {
		Want     boardtype.SBC
		Have     boardtype.SBC
		expected bool
	}{
		{JetsonAGXOrin, JetsonAGXOrin64GB, true},
		{JetsonAGXOrin64GB, JetsonAGXOrin, false},
		{JetsonOrinNano, JetsonOrinNano8GB, true},
		{JetsonOrinNano, JetsonAGXOrin, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Want_%v_Have_%v", test.Want.GetPrettyName(), test.Have.GetPrettyName()), func(t *testing.T) {
			if test.Have.IsBoardType(test.Want) != test.expected {
				t.Fatalf("IsBoardType() returned %v, expected %v", !test.expected, test.expected)
			}
		})
	}
}

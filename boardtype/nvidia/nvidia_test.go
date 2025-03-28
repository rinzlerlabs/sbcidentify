package nvidia

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rinzlerlabs/sbcidentify/boardtype"
)

func TestParseModuleName(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	f, e := getModuleNameFromDtsFilename(logger, "/dvs/git/dirty/git-master_linux/kernel/kernel-5.10/arch/arm64/boot/dts/../../../../../../hardware/nvidia/platform/t23x/p3768/kernel-dts/tegra234-p3767-0003-p3768-0000-a0.dts")
	require.NoError(t, e)
	require.Equal(t, "tegra234-p3767-0003-p3768-0000-a0", f)
	var detectedType boardtype.SBC
	f, e = getModuleNameFromDtsFilename(logger, "/nv-public/nv-platform/tegra234-p3768-0000+p3767-0000-nv-dsboard-ornx.dts")
	require.NoError(t, e)
	require.Equal(t, "tegra234-p3768-0000+p3767-0000-nv-board-orin", f)
	for _, m := range jetsonModulesByModelNumber {
		if strings.Contains(f, m.Model) {
			detectedType = m.Type
		}
	}
	require.Equal(t, boardtype.JetsonOrinNX16GB, detectedType)
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
		{boardtype.NVIDIA, boardtype.JetsonAGXOrin64GB, true},
		{boardtype.Jetson, boardtype.JetsonAGXOrin64GB, true},
		{boardtype.JetsonAGXOrin, boardtype.JetsonAGXOrin64GB, true},
		{boardtype.JetsonAGXOrin64GB, boardtype.JetsonAGXOrin, false},
		{boardtype.JetsonOrinNano, boardtype.JetsonOrinNano8GB, true},
		{boardtype.JetsonOrinNano, boardtype.JetsonAGXOrin, false},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Want_%v_Have_%v", test.Want.GetPrettyName(), test.Have.GetPrettyName()), func(t *testing.T) {
			if test.Have.IsBoardType(test.Want) != test.expected {
				t.Fatalf("IsBoardType() returned %v, expected %v", !test.expected, test.expected)
			}
		})
	}
}

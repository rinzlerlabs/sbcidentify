package nvidia

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
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

	// The + format (carrier+module combined DTS) strips to the base filename without extension.
	// getModuleModelFromModuleName cannot extract a usable model number from the + format,
	// so these boards fall back to DTBM detection on real hardware.
	f, e = getModuleNameFromDtsFilename(logger, "/nv-public/nv-platform/tegra234-p3768-0000+p3767-0000-nv-dsboard-ornx.dts")
	require.NoError(t, e)
	require.Equal(t, "tegra234-p3768-0000+p3767-0000-nv-dsboard-ornx", f)
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
		// AGX Orin hierarchy
		{boardtype.NVIDIA, boardtype.JetsonAGXOrin64GB, true},
		{boardtype.Jetson, boardtype.JetsonAGXOrin64GB, true},
		{boardtype.JetsonAGXOrin, boardtype.JetsonAGXOrin64GB, true},
		{boardtype.JetsonAGXOrin64GB, boardtype.JetsonAGXOrin, false},
		{boardtype.JetsonAGXOrin32GB, boardtype.JetsonAGXOrin64GB, false},
		// Orin NX — goes through JetsonOrin
		{boardtype.JetsonOrin, boardtype.JetsonOrinNX16GB, true},
		{boardtype.JetsonOrinNX, boardtype.JetsonOrinNX16GB, true},
		{boardtype.JetsonOrinNX16GB, boardtype.JetsonOrinNX8GB, false},
		// Orin Nano
		{boardtype.JetsonOrinNano, boardtype.JetsonOrinNano8GB, true},
		{boardtype.JetsonOrinNano, boardtype.JetsonAGXOrin, false},
		// AGX Orin is not Orin NX
		{boardtype.JetsonAGXOrin, boardtype.JetsonOrinNX16GB, false},
		// Xavier hierarchy
		{boardtype.Jetson, boardtype.JetsonAGXXavier32GB, true},
		{boardtype.JetsonAGXXavier, boardtype.JetsonAGXXavier32GB, true},
		{boardtype.JetsonAGXXavier32GB, boardtype.JetsonAGXXavier, false},
		{boardtype.JetsonXavierNX, boardtype.JetsonXavierNX16GB, true},
		// Cross-family false
		{boardtype.JetsonXavier, boardtype.JetsonOrinNX16GB, false},
		{boardtype.JetsonOrin, boardtype.JetsonAGXXavier32GB, false},
		// TX2 family
		{boardtype.Jetson, boardtype.JetsonTX2, true},
		{boardtype.JetsonTX2, boardtype.JetsonTX2i, true},
		{boardtype.JetsonTX2, boardtype.JetsonTX24GB, true},
		{boardtype.JetsonTX2, boardtype.JetsonTX1, false},
		// Nano family
		{boardtype.JetsonNano, boardtype.JetsonNano4GB, true},
		{boardtype.JetsonNano, boardtype.JetsonNanoDeveloperKit, true},
		{boardtype.JetsonNano, boardtype.JetsonTX1, false},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Want_%v_Have_%v", test.Want.GetPrettyName(), test.Have.GetPrettyName()), func(t *testing.T) {
			if test.Have.IsBoardType(test.Want) != test.expected {
				t.Fatalf("IsBoardType() returned %v, expected %v", !test.expected, test.expected)
			}
		})
	}
}

func writeMemInfo(t *testing.T, totalKB int) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "meminfo")
	content := fmt.Sprintf("MemTotal:       %d kB\nMemFree:        1000000 kB\n", totalKB)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	memInfoPath = path
	t.Cleanup(func() { memInfoPath = "/proc/meminfo" })
}

func TestGetInstalledRAMMB(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	writeMemInfo(t, 32505856) // 31744 MB
	ram, err := getInstalledRAMMB(logger)
	require.NoError(t, err)
	assert.Equal(t, 31744, ram)

	// missing file → error
	memInfoPath = "/nonexistent/meminfo"
	t.Cleanup(func() { memInfoPath = "/proc/meminfo" })
	_, err = getInstalledRAMMB(logger)
	assert.Error(t, err)
}

func TestRefineByInstalledRAM(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name    string
		totalKB int
		board   boardtype.SBC
		want    boardtype.SBC
	}{
		// AGX Orin Developer Kit: p3701-0000 with 32GB RAM.
		// /proc/meminfo shows ~30-31 GB after NVIDIA carveouts.
		{"AGX Orin 32GB devkit", 32505856, boardtype.JetsonAGXOrin, boardtype.JetsonAGXOrin32GB},
		// AGX Orin Developer Kit with 64GB RAM (~60 GB visible to OS).
		{"AGX Orin 64GB devkit", 62914560, boardtype.JetsonAGXOrin, boardtype.JetsonAGXOrin64GB},
		// Board already has specific RAM — no refinement needed.
		{"already specific", 32505856, boardtype.JetsonOrinNX16GB, boardtype.JetsonOrinNX16GB},
		// Generic board with no candidates (Jetson Nano) — stays generic.
		{"no candidates", 4194304, boardtype.JetsonNano, boardtype.JetsonNano},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.totalKB > 0 {
				writeMemInfo(t, tt.totalKB)
			}
			got := refineByInstalledRAM(logger, tt.board)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetSOC(t *testing.T) {
	tests := []struct {
		board    boardtype.SBC
		expected string
	}{
		{boardtype.JetsonOrinNX16GB, "Tegra234"},
		{boardtype.JetsonOrinNX8GB, "Tegra234"},
		{boardtype.JetsonOrinNano8GB, "Tegra234"},
		{boardtype.JetsonAGXOrin64GB, "Tegra234"},
		{boardtype.JetsonAGXXavier16GB, "Tegra194"},
		{boardtype.JetsonXavierNX8GB, "Tegra194"},
		{boardtype.ClaraAGX, "Tegra194"},
		{boardtype.JetsonTX2, "Tegra186"},
		{boardtype.JetsonTX2i, "Tegra186"},
		{boardtype.JetsonTX24GB, "Tegra186"},
		{boardtype.JetsonNano4GB, "Tegra210"},
		{boardtype.JetsonTX1, "Tegra210"},
		{boardtype.ShieldTV, "Tegra210"},
		{boardtype.NVIDIA, ""},
		{boardtype.Jetson, ""},
	}
	for _, test := range tests {
		t.Run(test.board.GetPrettyName(), func(t *testing.T) {
			assert.Equal(t, test.expected, test.board.GetSOC())
		})
	}
}

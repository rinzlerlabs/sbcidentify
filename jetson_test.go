package sbcidentify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseModuleName(t *testing.T) {
	f, e := getModuleNameFromDtsFilename("/dvs/git/dirty/git-master_linux/kernel/kernel-5.10/arch/arm64/boot/dts/../../../../../../hardware/nvidia/platform/t23x/p3768/kernel-dts/tegra234-p3767-0003-p3768-0000-a0.dts")
	assert.NoError(t, e)
	assert.Equal(t, "tegra234-p3767-0003-p3768-0000-a0", f)
}

func TestParseModelNameFromModuleName(t *testing.T) {
	f, e := getModuleModelFromModuleName("tegra234-p3767-0003-p3768-0000-a0")
	assert.NoError(t, e)
	assert.Equal(t, "p3767-0003", f)
}

package boardtype_test

import (
	"testing"

	"github.com/rinzlerlabs/sbcidentify/boardtype"
	"github.com/stretchr/testify/assert"
)

func TestGetPrettyName(t *testing.T) {
	tests := []struct {
		board    boardtype.BoardType
		expected string
	}{
		{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin NX", RAM: 16384}, "NVIDIA Jetson Orin NX 16GB"},
		{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin NX", RAM: 8192}, "NVIDIA Jetson Orin NX 8GB"},
		{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "3A+", RAM: 512}, "Raspberry Pi Raspberry Pi 3A+ 512MB"},
		{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin", RAM: 0}, "NVIDIA Jetson Orin"},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			assert.Equal(t, test.expected, test.board.GetPrettyName())
		})
	}
}

func TestIsBoardTypeHierarchy(t *testing.T) {
	root := boardtype.BoardType{Manufacturer: "X", Model: "M", SubModel: ""}
	child := boardtype.BoardType{Manufacturer: "X", Model: "M", SubModel: "S", BaseModel: &root}
	grandchild := boardtype.BoardType{Manufacturer: "X", Model: "M", SubModel: "S", RAM: 1024, BaseModel: &child}
	other := boardtype.BoardType{Manufacturer: "Y", Model: "M", SubModel: ""}

	// grandchild matches itself and all ancestors
	assert.True(t, grandchild.IsBoardType(grandchild))
	assert.True(t, grandchild.IsBoardType(child))
	assert.True(t, grandchild.IsBoardType(root))
	// ancestors do not match descendants
	assert.False(t, root.IsBoardType(grandchild))
	assert.False(t, child.IsBoardType(grandchild))
	// different manufacturer never matches
	assert.False(t, grandchild.IsBoardType(other))
	assert.False(t, root.IsBoardType(other))
}

func TestGetSOCOnRealBoardVars(t *testing.T) {
	assert.Equal(t, "Tegra234", boardtype.JetsonOrinNX16GB.GetSOC())
	assert.Equal(t, "Tegra194", boardtype.JetsonAGXXavier16GB.GetSOC())
	assert.Equal(t, "Tegra186", boardtype.JetsonTX2.GetSOC())
	assert.Equal(t, "Tegra210", boardtype.JetsonNano4GB.GetSOC())
	assert.Equal(t, "BCM2837", boardtype.RaspberryPi3B.GetSOC())
	assert.Equal(t, "BCM2711", boardtype.RaspberryPi4B8GB.GetSOC())
	assert.Equal(t, "BCM2712", boardtype.RaspberryPi5B8GB.GetSOC())
	assert.Equal(t, "BCM2712", boardtype.RaspberryPiCM5Lite4GB.GetSOC())
	assert.Equal(t, "", boardtype.NVIDIA.GetSOC())
	assert.Equal(t, "", boardtype.RaspberryPi.GetSOC())
}

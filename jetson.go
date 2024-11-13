package sbcidentify

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	BoardTypeJetsonOrinNX16GB              BoardType = "NVIDIA Jetson Orin NX 16GB RAM"
	BoardTypeJetsonOrinNX8GB               BoardType = "NVIDIA Jetson Orin NX 8GB RAM"
	BoardTypeJetsonOrinNano8GB             BoardType = "NVIDIA Jetson Orin Nano 8GB RAM"
	BoardTypeJetsonOrinNano4GB             BoardType = "NVIDIA Jetson Orin Nano 4GB RAM"
	BoardTypeJetsonOrinNanoDeveloperKit    BoardType = "NVIDIA Jetson Orin Nano Developer kit"
	BoardTypeJetsonAGXOrin                 BoardType = "NVIDIA Jetson AGX Orin"
	BoardTypeJetsonAGXOrin32GB             BoardType = "NVIDIA Jetson AGX Orin 32GB RAM"
	BoardTypeJetsonAGXOrin64GB             BoardType = "NVIDIA Jetson AGX Orin 64GB RAM"
	BoardTypeJetsonXavierNXDeveloperKit    BoardType = "NVIDIA Jetson Xavier NX Developer kit"
	BoardTypeJetsonXavierNX8GB             BoardType = "NVIDIA Jetson Xavier NX 8GB RAM"
	BoardTypeJetsonXavierNX16GB            BoardType = "NVIDIA Jetson Xavier NX 16GB RAM"
	BoardTypeJetsonAGXXavier8GB            BoardType = "NVIDIA Jetson AGX Xavier 8GB RAM"
	BoardTypeJetsonAGXXavier16GB           BoardType = "NVIDIA Jetson AGX Xavier 16GB RAM"
	BoardTypeJetsonAGXXavier32GB           BoardType = "NVIDIA Jetson AGX Xavier 32GB RAM"
	BoardTypeJetsonAGXXavier64GB           BoardType = "NVIDIA Jetson AGX Xavier 64GB RAM"
	BoardTypeJetsonAGXXavierIndustrial32GB BoardType = "NVIDIA Jetson AGX Xavier Industrial 32GB RAM"
	BoardTypeJetsonNano2GB                 BoardType = "NVIDIA Jetson Nano 2GB RAM"
	BoardTypeJetsonNano16GbEMMC            BoardType = "NVIDIA Jetson Nano module 16GB eMMC"
	BoardTypeJetsonNano4GB                 BoardType = "NVIDIA Jetson Nano 4GB RAM"
	BoardTypeJetsonTX2NX                   BoardType = "NVIDIA Jetson TX2 NX"
	BoardTypeJetsonTX24GB                  BoardType = "NVIDIA Jetson TX2 4GB RAM"
	BoardTypeJetsonTX2i                    BoardType = "NVIDIA Jetson TX2i"
	BoardTypeJetsonTX2                     BoardType = "NVIDIA Jetson TX2"
	BoardTypeJetsonTX1                     BoardType = "NVIDIA Jetson TX1"
	BoardTypeJetsonTK1                     BoardType = "NVIDIA Jetson TK1"
	BoardTypeClaraAGX                      BoardType = "NVIDIA Clara AGX"
)

type jetson struct {
	Model string
	Type  BoardType
}

var jetsonModules = []jetson{
	{"p3767-0000", BoardTypeJetsonOrinNX16GB},
	{"p3767-0001", BoardTypeJetsonOrinNX8GB},

	{"p3767-0003", BoardTypeJetsonOrinNano8GB},
	{"p3767-0004", BoardTypeJetsonOrinNano4GB},
	{"p3767-0005", BoardTypeJetsonOrinNanoDeveloperKit},

	{"p3701-0000", BoardTypeJetsonAGXOrin},
	{"p3701-0004", BoardTypeJetsonAGXOrin32GB},
	{"p3701-0005", BoardTypeJetsonAGXOrin64GB},

	{"p3668-0000", BoardTypeJetsonXavierNXDeveloperKit},
	{"p3668-0001", BoardTypeJetsonXavierNX8GB},
	{"p3668-0003", BoardTypeJetsonXavierNX16GB},

	{"p2888-0001", BoardTypeJetsonAGXXavier16GB},
	{"p2888-0003", BoardTypeJetsonAGXXavier32GB},
	{"p2888-0004", BoardTypeJetsonAGXXavier32GB},
	{"p2888-0005", BoardTypeJetsonAGXXavier64GB},
	{"p2888-0006", BoardTypeJetsonAGXXavier8GB},
	{"p2888-0008", BoardTypeJetsonAGXXavierIndustrial32GB},

	{"p3448-0000", BoardTypeJetsonNano4GB},
	{"p3448-0002", BoardTypeJetsonNano16GbEMMC},
	{"p3448-0003", BoardTypeJetsonNano2GB},

	{"p3636-0001", BoardTypeJetsonTX2NX},
	{"p3509-0000", BoardTypeJetsonTX2NX},

	{"p3489-0888", BoardTypeJetsonTX24GB},
	{"p3489-0000", BoardTypeJetsonTX2i},
	{"p3310-1000", BoardTypeJetsonTX2},

	{"p2180-1000", BoardTypeJetsonTX1},

	{"r375-0001", BoardTypeJetsonTK1},

	{"p3904-0000", BoardTypeClaraAGX},
}

type jetsonIdentifier struct{}

func (r jetsonIdentifier) GetBoardType() (BoardType, error) {
	dtsFilename, err := getDtsFilename()
	if err != nil {
		return BoardTypeUnknown, err
	}
	moduleName, err := getModuleNameFromDtsFilename(dtsFilename)
	if err != nil {
		return BoardTypeUnknown, err
	}
	moduleModel, err := getModuleModelFromModuleName(moduleName)
	if err != nil {
		return BoardTypeUnknown, err
	}
	for _, m := range jetsonModules {
		if m.Model == moduleModel {
			return m.Type, nil
		}
	}
	return BoardTypeUnknown, ErrUnknownBoard
}

func getDtsFilename() (string, error) {
	s, e := os.ReadFile("/proc/device-tree/nvidia,dtsfilename")
	if e != nil {
		return "", e
	}
	return string(s), nil
}

func getModuleNameFromDtsFilename(dtsFilename string) (string, error) {
	filename := filepath.Base(dtsFilename)
	return strings.TrimSuffix(filename, filepath.Ext(filename)), nil
}

func getModuleModelFromModuleName(moduleName string) (string, error) {
	parts := strings.Split(moduleName, "-")
	if len(parts) >= 4 {
		return strings.Join(parts[1:3], "-"), nil
	}
	return "", ErrUnknownBoard
}

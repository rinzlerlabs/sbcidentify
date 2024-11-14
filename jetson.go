package sbcidentify

import (
	"errors"
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
	BoardTypeJetsonAGXXavier               BoardType = "NVIDIA Jetson AGX Xavier"
	BoardTypeJetsonAGXXavier8GB            BoardType = "NVIDIA Jetson AGX Xavier 8GB RAM"
	BoardTypeJetsonAGXXavier16GB           BoardType = "NVIDIA Jetson AGX Xavier 16GB RAM"
	BoardTypeJetsonAGXXavier32GB           BoardType = "NVIDIA Jetson AGX Xavier 32GB RAM"
	BoardTypeJetsonAGXXavier64GB           BoardType = "NVIDIA Jetson AGX Xavier 64GB RAM"
	BoardTypeJetsonAGXXavierIndustrial32GB BoardType = "NVIDIA Jetson AGX Xavier Industrial 32GB RAM"
	BoardTypeJetsonNanoDeveloperKit        BoardType = "NVIDIA Jetson Nano Developer Kit"
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
	BoardTypeShieldTV                      BoardType = "NVIDIA Shield TV"
)

// NVIDIA Jetson AGX Orin Developer Kit
type jetson struct {
	Model string
	Type  BoardType
}

var (
	ErrDtsFileDoesNotExist = errors.New("DTS file does not exist")
)

var jetsonModulesByModelNumber = []jetson{
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
	{"p2972-0000", BoardTypeJetsonAGXXavier},

	{"p2771-0000", BoardTypeJetsonTX2},

	{"p3448-0000", BoardTypeJetsonNano4GB},
	{"p3448-0002", BoardTypeJetsonNano16GbEMMC},
	{"p3448-0003", BoardTypeJetsonNano2GB},
	{"p3450-0000", BoardTypeJetsonNanoDeveloperKit},

	{"p3636-0001", BoardTypeJetsonTX2NX},
	{"p3509-0000", BoardTypeJetsonTX2NX},

	{"p3489-0888", BoardTypeJetsonTX24GB},
	{"p3489-0000", BoardTypeJetsonTX2i},
	{"p3310-1000", BoardTypeJetsonTX2},

	{"p2180-1000", BoardTypeJetsonTX1},
	{"p2371-2180", BoardTypeJetsonTX1},

	{"p2894-0050", BoardTypeShieldTV},

	{"r375-0001", BoardTypeJetsonTK1},

	{"p3904-0000", BoardTypeClaraAGX},
}

var jetsonModulesByDeviceTreeBaseModel = []jetson{
	{"NVIDIA Jetson Orin NX Engineering Reference Developer Kit", BoardTypeJetsonOrinNX16GB},
	{"NVIDIA Jetson Orin Nano Developer Kit", BoardTypeJetsonOrinNanoDeveloperKit},
	{"NVIDIA Jetson TX2 Developer Kit", BoardTypeJetsonTX2},
	{"NVIDIA Jetson TX2", BoardTypeJetsonTX2},
	{"NVIDIA Jetson TX2 NX Developer Kit", BoardTypeJetsonTX2NX},
	{"NVIDIA Jetson AGX Xavier", BoardTypeJetsonAGXXavier},
	{"NVIDIA Jetson AGX Xavier Developer Kit", BoardTypeJetsonAGXXavier},
	{"NVIDIA Jetson Xavier NX Developer Kit (SD-card)", BoardTypeJetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX Developer Kit (eMMC)", BoardTypeJetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX (SD-card)", BoardTypeJetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX (eMMC)", BoardTypeJetsonXavierNX8GB},
	{"NVIDIA Jetson TX1", BoardTypeJetsonTX1},
	{"NVIDIA Jetson TX1 Developer Kit", BoardTypeJetsonTX1},
	{"NVIDIA Shield TV", BoardTypeShieldTV},
	{"NVIDIA Jetson Nano Developer Kit", BoardTypeJetsonNanoDeveloperKit},
	{"NVIDIA Jetson AGX Orin Developer Kit", BoardTypeJetsonAGXOrin},
	{"NVIDIA Jetson AGX Orin", BoardTypeJetsonAGXOrin},
}

type jetsonIdentifier struct{}

func (r jetsonIdentifier) GetBoardType() (BoardType, error) {
	boardType, err := getBoardTypeFromModuleModel()
	if err == ErrUnknownBoard || err == ErrDtsFileDoesNotExist {
		return getBoardTypeByDeviceTreeBaseModel()
	} else if err != nil {
		return BoardTypeUnknown, err
	} else {
		return boardType, nil
	}
}

func getBoardTypeFromModuleModel() (BoardType, error) {
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
	for _, m := range jetsonModulesByModelNumber {
		if m.Model == moduleModel {
			return m.Type, nil
		}
	}
	return BoardTypeUnknown, ErrUnknownBoard
}

func getBoardTypeByDeviceTreeBaseModel() (BoardType, error) {
	dtbm, err := getDeviceTreeBaseModel()
	if err != nil {
		return BoardTypeUnknown, err
	}
	for _, m := range jetsonModulesByDeviceTreeBaseModel {
		if strings.Contains(dtbm, m.Model) {
			return m.Type, nil
		}
	}
	return BoardTypeUnknown, ErrUnknownBoard
}

func getDtsFilename() (string, error) {
	if _, err := os.Stat("/proc/device-tree/nvidia,dtsfilename"); os.IsNotExist(err) {
		return "", ErrDtsFileDoesNotExist
	}
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

package nvidia

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/thegreatco/sbcidentify/boardtype"
	"github.com/thegreatco/sbcidentify/identifier"
)

func init() {
	identifier.RegisterBoardIdentifier(NewNvidiaIdentifier)
}

const (
	dtsFileName = "/proc/device-tree/nvidia,dtsfilename"
)

type NVIDIA struct{ boardtype.BoardType }
type Jetson struct{ NVIDIA }
type JetsonOrin struct{ Jetson }
type JetsonOrinNX struct{ JetsonOrin }
type JetsonOrinNano struct{ JetsonOrin }
type JetsonAGX struct{ Jetson }
type JetsonAGXXavier struct{ JetsonAGX }
type JetsonXavier struct{ Jetson }
type JetsonXavierNX struct{ JetsonXavier }
type JetsonNano struct{ Jetson }
type JetsonTX2 struct{ Jetson }
type ClaraAGX struct{ NVIDIA }
type ShieldTV struct{ NVIDIA }

var (
	jetsonOrinNX16GB              = JetsonOrinNX{JetsonOrin{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin NX", RAM: 16384}}}}}
	jetsonOrinNX8GB               = JetsonOrinNX{JetsonOrin{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin NX", RAM: 8192}}}}}
	jetsonOrinNano8GB             = JetsonOrinNano{JetsonOrin{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin Nano", RAM: 8192}}}}}
	jetsonOrinNano4GB             = JetsonOrinNano{JetsonOrin{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin Nano", RAM: 4096}}}}}
	jetsonOrinNanoDeveloperKit    = JetsonOrinNano{JetsonOrin{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin Nano Developer Kit", RAM: 8192}}}}}
	jetsonAGXOrin                 = JetsonAGX{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Orin", RAM: 0}}}}
	jetsonAGXOrin32GB             = JetsonAGX{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Orin", RAM: 32768}}}}
	jetsonAGXOrin64GB             = JetsonAGX{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Orin", RAM: 65536}}}}
	jetsonXavierNXDeveloperKit    = JetsonXavierNX{JetsonXavier{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier NX Developer Kit", RAM: 0}}}}}
	jetsonXavierNX8GB             = JetsonXavierNX{JetsonXavier{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier NX", RAM: 8192}}}}}
	jetsonXavierNX16GB            = JetsonXavierNX{JetsonXavier{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier NX", RAM: 16384}}}}}
	jetsonAGXXavier8GB            = JetsonAGXXavier{JetsonAGX{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 8192}}}}}
	jetsonAGXXavier               = JetsonAGXXavier{JetsonAGX{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 0}}}}}
	jetsonAGXXavier16GB           = JetsonAGXXavier{JetsonAGX{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 16384}}}}}
	jetsonAGXXavier32GB           = JetsonAGXXavier{JetsonAGX{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 32768}}}}}
	jetsonAGXXavier64GB           = JetsonAGXXavier{JetsonAGX{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 65536}}}}}
	jetsonAGXXavierIndustrial32GB = JetsonAGXXavier{JetsonAGX{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier Industrial", RAM: 32768}}}}}
	jetsonNanoDeveloperKit        = JetsonNano{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano Developer Kit", RAM: 0}}}}
	jetsonNano2GB                 = JetsonNano{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano", RAM: 2048}}}}
	jetsonNano16GbEMMC            = JetsonNano{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano", RAM: 0}}}}
	jetsonNano4GB                 = JetsonNano{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano", RAM: 4096}}}}
	jetsonTX2NX                   = JetsonTX2{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2 NX", RAM: 0}}}}
	jetsonTX24GB                  = JetsonTX2{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2", RAM: 4096}}}}
	jetsonTX2i                    = JetsonTX2{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2i", RAM: 0}}}}
	jetsonTX2                     = JetsonTX2{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2", RAM: 0}}}}
	jetsonTX1                     = JetsonTX2{Jetson{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX1", RAM: 0}}}}
	claraAGX                      = ClaraAGX{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Clara", SubModel: "AGX", RAM: 0}}}
	shieldTV                      = ShieldTV{NVIDIA{boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Shield", SubModel: "TV", RAM: 0}}}
)

// NVIDIA Jetson AGX Orin Developer Kit
type jetson struct {
	Model string
	Type  boardtype.SBC
}

var (
	ErrDtsFileDoesNotExist = errors.New("DTS file does not exist")
	ErrCannotIdentifyBoard = errors.New("cannot identify NVIDIA board")
)

var jetsonModulesByModelNumber = []jetson{
	{"p3767-0000", jetsonOrinNX16GB},
	{"p3767-0001", jetsonOrinNX8GB},

	{"p3767-0003", jetsonOrinNano8GB},
	{"p3767-0004", jetsonOrinNano4GB},
	{"p3767-0005", jetsonOrinNanoDeveloperKit},

	{"p3701-0000", jetsonAGXOrin},
	{"p3701-0004", jetsonAGXOrin32GB},
	{"p3701-0005", jetsonAGXOrin64GB},

	{"p3668-0000", jetsonXavierNXDeveloperKit},
	{"p3668-0001", jetsonXavierNX8GB},
	{"p3668-0003", jetsonXavierNX16GB},

	{"p2888-0001", jetsonAGXXavier16GB},
	{"p2888-0003", jetsonAGXXavier32GB},
	{"p2888-0004", jetsonAGXXavier32GB},
	{"p2888-0005", jetsonAGXXavier64GB},
	{"p2888-0006", jetsonAGXXavier8GB},
	{"p2888-0008", jetsonAGXXavierIndustrial32GB},
	{"p2972-0000", jetsonAGXXavier},

	{"p2771-0000", jetsonTX2},

	{"p3448-0000", jetsonNano4GB},
	{"p3448-0002", jetsonNano16GbEMMC},
	{"p3448-0003", jetsonNano2GB},
	{"p3450-0000", jetsonNanoDeveloperKit},

	{"p3636-0001", jetsonTX2NX},
	{"p3509-0000", jetsonTX2NX},

	{"p3489-0888", jetsonTX24GB},
	{"p3489-0000", jetsonTX2i},
	{"p3310-1000", jetsonTX2},

	{"p2180-1000", jetsonTX1},
	{"p2371-2180", jetsonTX1},

	{"p2894-0050", shieldTV},

	{"p3904-0000", claraAGX},
}

var jetsonModulesByDeviceTreeBaseModel = []jetson{
	{"NVIDIA Jetson Orin NX Engineering Reference Developer Kit", jetsonOrinNX16GB},
	{"NVIDIA Jetson Orin Nano Developer Kit", jetsonOrinNanoDeveloperKit},
	{"NVIDIA Jetson TX2 Developer Kit", jetsonTX2},
	{"NVIDIA Jetson TX2", jetsonTX2},
	{"NVIDIA Jetson TX2 NX Developer Kit", jetsonTX2NX},
	{"NVIDIA Jetson AGX Xavier", jetsonAGXXavier},
	{"NVIDIA Jetson AGX Xavier Developer Kit", jetsonAGXXavier},
	{"NVIDIA Jetson Xavier NX Developer Kit (SD-card)", jetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX Developer Kit (eMMC)", jetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX (SD-card)", jetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX (eMMC)", jetsonXavierNX8GB},
	{"NVIDIA Jetson TX1", jetsonTX1},
	{"NVIDIA Jetson TX1 Developer Kit", jetsonTX1},
	{"NVIDIA Shield TV", shieldTV},
	{"NVIDIA Jetson Nano Developer Kit", jetsonNanoDeveloperKit},
	{"NVIDIA Jetson AGX Orin Developer Kit", jetsonAGXOrin},
	{"NVIDIA Jetson AGX Orin", jetsonAGXOrin},
}

type jetsonIdentifier struct {
	logger *slog.Logger
}

func NewNvidiaIdentifier(logger *slog.Logger) identifier.BoardIdentifier {
	logger.Debug("initializing Jetson identifier")
	newLogger := logger.With(slog.String("source", "NVIDIA"))
	return jetsonIdentifier{
		logger: newLogger,
	}
}

func (r jetsonIdentifier) Name() string {
	return "Jetson Identifier"
}

func (r jetsonIdentifier) GetBoardType() (boardtype.SBC, error) {
	boardType, err := getBoardTypeFromModuleModel(r.logger)
	if err == ErrDtsFileDoesNotExist {
		r.logger.Debug("DTS file does not exist, falling back to device tree base model")
		return getBoardTypeByDeviceTreeBaseModel(r.logger)
	} else if err == identifier.ErrCannotIdentifyBoard {
		r.logger.Debug("unknown board, falling back to device tree base model")
		boardType, err = getBoardTypeByDeviceTreeBaseModel(r.logger)
		if err == identifier.ErrCannotIdentifyBoard {
			r.logger.Debug("unknown board")
			return boardtype.BoardTypeUnknown, ErrCannotIdentifyBoard
		} else if err != nil {
			r.logger.Debug("error getting board type", slog.Any("error", err))
			return boardtype.BoardTypeUnknown, err
		} else {
			r.logger.Debug("board type", slog.String("type", string(boardType.GetPrettyName())))
			return boardType, nil
		}
	} else if err != nil {
		r.logger.Debug("error getting board type", slog.Any("error", err))
		return boardtype.BoardTypeUnknown, err
	} else {
		r.logger.Debug("board type", slog.String("type", string(boardType.GetPrettyName())))
		return boardType, nil
	}
}

func getBoardTypeFromModuleModel(logger *slog.Logger) (boardtype.SBC, error) {
	dtsFilename, err := getDtsFile(logger)
	if err != nil {
		return boardtype.BoardTypeUnknown, err
	}
	moduleName, err := getModuleNameFromDtsFilename(logger, dtsFilename)
	if err != nil {
		return boardtype.BoardTypeUnknown, err
	}
	moduleModel, err := getModuleModelFromModuleName(logger, moduleName)
	if err != nil {
		return boardtype.BoardTypeUnknown, err
	}
	for _, m := range jetsonModulesByModelNumber {
		if m.Model == moduleModel {
			return m.Type, nil
		}
	}
	return boardtype.BoardTypeUnknown, identifier.ErrCannotIdentifyBoard
}

func getBoardTypeByDeviceTreeBaseModel(logger *slog.Logger) (boardtype.SBC, error) {
	dtbm, err := identifier.GetDeviceTreeBaseModel(logger)
	if err != nil {
		return boardtype.BoardTypeUnknown, err
	}
	for _, m := range jetsonModulesByDeviceTreeBaseModel {
		if strings.Contains(dtbm, m.Model) {
			return m.Type, nil
		}
	}
	logger.Debug("device tree base model does not match any boards", slog.String("model", dtbm))
	return boardtype.BoardTypeUnknown, ErrCannotIdentifyBoard
}

func getDtsFile(logger *slog.Logger) (string, error) {
	if _, err := os.Stat(dtsFileName); os.IsNotExist(err) {
		logger.Debug("DTS file does not exist", slog.Any("error", err))
		return "", ErrDtsFileDoesNotExist
	}
	s, e := os.ReadFile(dtsFileName)
	if e != nil {
		logger.Debug("cannot read DTS file", slog.Any("error", e))
		return "", e
	}
	str := string(s)
	logger.Debug("DTS file", slog.String("filename", str))
	return str, nil
}

func getModuleNameFromDtsFilename(logger *slog.Logger, dtsFilename string) (string, error) {
	filename := filepath.Base(dtsFilename)
	ret := strings.TrimSuffix(filename, filepath.Ext(filename))
	logger.Debug("module name", slog.String("name", ret))
	return ret, nil
}

func getModuleModelFromModuleName(logger *slog.Logger, moduleName string) (string, error) {
	parts := strings.Split(moduleName, "-")
	if len(parts) >= 4 {
		ret := strings.Join(parts[1:3], "-")
		logger.Debug("module model", slog.String("model", ret))
		return ret, nil
	}
	logger.Debug("error parsing module name", slog.String("name", moduleName))
	return "", identifier.ErrCannotIdentifyBoard
}

package nvidia

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/rinzlerlabs/sbcidentify/boardtype"
	"github.com/rinzlerlabs/sbcidentify/identifier"
)

func init() {
	identifier.RegisterBoardIdentifier(NewNvidiaIdentifier)
}

const (
	dtsFileName = "/proc/device-tree/nvidia,dtsfilename"
)

var (
	NVIDIA                        = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "", SubModel: "", RAM: 0}
	Jetson                        = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "", RAM: 0, BaseModel: &NVIDIA}
	JetsonOrin                    = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin", RAM: 0, BaseModel: &Jetson}
	JetsonXavier                  = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier", RAM: 0, BaseModel: &Jetson}
	JetsonOrinNX                  = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin NX", RAM: 0, BaseModel: &JetsonOrin}
	JetsonOrinNX16GB              = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin NX", RAM: 16384, BaseModel: &JetsonOrinNX}
	JetsonOrinNX8GB               = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin NX", RAM: 8192, BaseModel: &JetsonOrinNX}
	JetsonOrinNano                = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin Nano", RAM: 0, BaseModel: &JetsonOrin}
	JetsonOrinNano8GB             = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin Nano", RAM: 8192, BaseModel: &JetsonOrinNano}
	JetsonOrinNano4GB             = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin Nano", RAM: 4096, BaseModel: &JetsonOrinNano}
	JetsonOrinNanoDeveloperKit    = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Orin Nano Developer Kit", RAM: 8192, BaseModel: &JetsonOrinNano}
	JetsonAGXOrin                 = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Orin", RAM: 0, BaseModel: &Jetson}
	JetsonAGXOrin32GB             = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Orin", RAM: 32768, BaseModel: &JetsonAGXOrin}
	JetsonAGXOrin64GB             = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Orin", RAM: 65536, BaseModel: &JetsonAGXOrin}
	JetsonXavierNX                = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier NX", RAM: 0, BaseModel: &Jetson}
	JetsonXavierNXDeveloperKit    = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier NX Developer Kit", RAM: 0, BaseModel: &JetsonXavierNX}
	JetsonXavierNX8GB             = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier NX", RAM: 8192, BaseModel: &JetsonXavierNX}
	JetsonXavierNX16GB            = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Xavier NX", RAM: 16384, BaseModel: &JetsonXavierNX}
	JetsonAGXXavier               = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 0, BaseModel: &Jetson}
	JetsonAGXXavier8GB            = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 8192, BaseModel: &JetsonAGXXavier}
	JetsonAGXXavier16GB           = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 16384, BaseModel: &JetsonAGXXavier}
	JetsonAGXXavier32GB           = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 32768, BaseModel: &JetsonAGXXavier}
	JetsonAGXXavier64GB           = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier", RAM: 65536, BaseModel: &JetsonAGXXavier}
	JetsonAGXXavierIndustrial32GB = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "AGX Xavier Industrial", RAM: 32768, BaseModel: &JetsonAGXXavier}
	JetsonNano                    = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano", RAM: 0, BaseModel: &Jetson}
	JetsonNanoDeveloperKit        = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano Developer Kit", RAM: 0, BaseModel: &JetsonNano}
	JetsonNano2GB                 = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano", RAM: 2048, BaseModel: &JetsonNano}
	JetsonNano16GbEMMC            = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano", RAM: 0, BaseModel: &JetsonNano}
	JetsonNano4GB                 = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "Nano", RAM: 4096, BaseModel: &JetsonNano}
	JetsonTX2NX                   = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2 NX", RAM: 0, BaseModel: &Jetson}
	JetsonTX24GB                  = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2", RAM: 4096, BaseModel: &JetsonTX2}
	JetsonTX2i                    = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2i", RAM: 0, BaseModel: &JetsonTX2}
	JetsonTX2                     = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX2", RAM: 0, BaseModel: &Jetson}
	JetsonTX1                     = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Jetson", SubModel: "TX1", RAM: 0, BaseModel: &Jetson}
	ClaraAGX                      = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Clara", SubModel: "AGX", RAM: 0, BaseModel: &NVIDIA}
	ShieldTV                      = boardtype.BoardType{Manufacturer: "NVIDIA", Model: "Shield", SubModel: "TV", RAM: 0, BaseModel: &NVIDIA}
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
	{"p3767-0000", JetsonOrinNX16GB},
	{"p3767-0001", JetsonOrinNX8GB},

	{"p3767-0003", JetsonOrinNano8GB},
	{"p3767-0004", JetsonOrinNano4GB},
	{"p3767-0005", JetsonOrinNanoDeveloperKit},

	{"p3701-0000", JetsonAGXOrin},
	{"p3701-0004", JetsonAGXOrin32GB},
	{"p3701-0005", JetsonAGXOrin64GB},

	{"p3668-0000", JetsonXavierNXDeveloperKit},
	{"p3668-0001", JetsonXavierNX8GB},
	{"p3668-0003", JetsonXavierNX16GB},

	{"p2888-0001", JetsonAGXXavier16GB},
	{"p2888-0003", JetsonAGXXavier32GB},
	{"p2888-0004", JetsonAGXXavier32GB},
	{"p2888-0005", JetsonAGXXavier64GB},
	{"p2888-0006", JetsonAGXXavier8GB},
	{"p2888-0008", JetsonAGXXavierIndustrial32GB},
	{"p2972-0000", JetsonAGXXavier},

	{"p2771-0000", JetsonTX2},

	{"p3448-0000", JetsonNano4GB},
	{"p3448-0002", JetsonNano16GbEMMC},
	{"p3448-0003", JetsonNano2GB},
	{"p3450-0000", JetsonNanoDeveloperKit},

	{"p3636-0001", JetsonTX2NX},
	{"p3509-0000", JetsonTX2NX},

	{"p3489-0888", JetsonTX24GB},
	{"p3489-0000", JetsonTX2i},
	{"p3310-1000", JetsonTX2},

	{"p2180-1000", JetsonTX1},
	{"p2371-2180", JetsonTX1},

	{"p2894-0050", ShieldTV},

	{"p3904-0000", ClaraAGX},
}

var jetsonModulesByDeviceTreeBaseModel = []jetson{
	{"NVIDIA Jetson Orin NX Engineering Reference Developer Kit", JetsonOrinNX16GB},
	{"NVIDIA Jetson Orin Nano Developer Kit", JetsonOrinNanoDeveloperKit},
	{"NVIDIA Jetson TX2 Developer Kit", JetsonTX2},
	{"NVIDIA Jetson TX2", JetsonTX2},
	{"NVIDIA Jetson TX2 NX Developer Kit", JetsonTX2NX},
	{"NVIDIA Jetson AGX Xavier", JetsonAGXXavier},
	{"NVIDIA Jetson AGX Xavier Developer Kit", JetsonAGXXavier},
	{"NVIDIA Jetson Xavier NX Developer Kit (SD-card)", JetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX Developer Kit (eMMC)", JetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX (SD-card)", JetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX (eMMC)", JetsonXavierNX8GB},
	{"NVIDIA Jetson TX1", JetsonTX1},
	{"NVIDIA Jetson TX1 Developer Kit", JetsonTX1},
	{"NVIDIA Shield TV", ShieldTV},
	{"NVIDIA Jetson Nano Developer Kit", JetsonNanoDeveloperKit},
	{"NVIDIA Jetson AGX Orin Developer Kit", JetsonAGXOrin},
	{"NVIDIA Jetson AGX Orin", JetsonAGXOrin},
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
		boardType, err = getBoardTypeByDeviceTreeBaseModel(r.logger)
		if err == identifier.ErrCannotIdentifyBoard {
			r.logger.Debug("unknown board")
			return nil, ErrCannotIdentifyBoard
		} else if err != nil {
			r.logger.Debug("error getting board type", slog.Any("error", err))
			return nil, err
		} else {
			r.logger.Debug("board type", slog.String("type", string(boardType.GetPrettyName())))
			return boardType, nil
		}
	} else if err == identifier.ErrCannotIdentifyBoard {
		r.logger.Debug("unknown board, falling back to device tree base model")
		boardType, err = getBoardTypeByDeviceTreeBaseModel(r.logger)
		if err == identifier.ErrCannotIdentifyBoard {
			r.logger.Debug("unknown board")
			return nil, ErrCannotIdentifyBoard
		} else if err != nil {
			r.logger.Debug("error getting board type", slog.Any("error", err))
			return nil, err
		} else {
			r.logger.Debug("board type", slog.String("type", string(boardType.GetPrettyName())))
			return boardType, nil
		}
	} else if err != nil {
		r.logger.Debug("error getting board type", slog.Any("error", err))
		return nil, err
	} else {
		r.logger.Debug("board type", slog.String("type", string(boardType.GetPrettyName())))
		return boardType, nil
	}
}

func getBoardTypeFromModuleModel(logger *slog.Logger) (boardtype.SBC, error) {
	dtsFilename, err := getDtsFile(logger)
	if err != nil {
		return nil, err
	}
	moduleName, err := getModuleNameFromDtsFilename(logger, dtsFilename)
	if err != nil {
		return nil, err
	}
	moduleModel, err := getModuleModelFromModuleName(logger, moduleName)
	if err != nil {
		return nil, err
	}
	for _, m := range jetsonModulesByModelNumber {
		if m.Model == moduleModel {
			return m.Type, nil
		}
	}
	return nil, identifier.ErrCannotIdentifyBoard
}

func getBoardTypeByDeviceTreeBaseModel(logger *slog.Logger) (boardtype.SBC, error) {
	dtbm, err := identifier.GetDeviceTreeBaseModel(logger)
	if err != nil {
		return nil, err
	}
	for _, m := range jetsonModulesByDeviceTreeBaseModel {
		if strings.Contains(dtbm, m.Model) {
			return m.Type, nil
		}
	}
	logger.Debug("device tree base model does not match any boards", slog.String("model", dtbm))
	return nil, ErrCannotIdentifyBoard
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

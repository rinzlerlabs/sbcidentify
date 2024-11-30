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
	{"p3767-0000", boardtype.JetsonOrinNX16GB},
	{"p3767-0001", boardtype.JetsonOrinNX8GB},

	{"p3767-0003", boardtype.JetsonOrinNano8GB},
	{"p3767-0004", boardtype.JetsonOrinNano4GB},
	{"p3767-0005", boardtype.JetsonOrinNanoDeveloperKit},

	{"p3701-0000", boardtype.JetsonAGXOrin},
	{"p3701-0004", boardtype.JetsonAGXOrin32GB},
	{"p3701-0005", boardtype.JetsonAGXOrin64GB},

	{"p3668-0000", boardtype.JetsonXavierNXDeveloperKit},
	{"p3668-0001", boardtype.JetsonXavierNX8GB},
	{"p3668-0003", boardtype.JetsonXavierNX16GB},

	{"p2888-0001", boardtype.JetsonAGXXavier16GB},
	{"p2888-0003", boardtype.JetsonAGXXavier32GB},
	{"p2888-0004", boardtype.JetsonAGXXavier32GB},
	{"p2888-0005", boardtype.JetsonAGXXavier64GB},
	{"p2888-0006", boardtype.JetsonAGXXavier8GB},
	{"p2888-0008", boardtype.JetsonAGXXavierIndustrial32GB},
	{"p2972-0000", boardtype.JetsonAGXXavier},

	{"p2771-0000", boardtype.JetsonTX2},

	{"p3448-0000", boardtype.JetsonNano4GB},
	{"p3448-0002", boardtype.JetsonNano16GbEMMC},
	{"p3448-0003", boardtype.JetsonNano2GB},
	{"p3450-0000", boardtype.JetsonNanoDeveloperKit},

	{"p3636-0001", boardtype.JetsonTX2NX},
	{"p3509-0000", boardtype.JetsonTX2NX},

	{"p3489-0888", boardtype.JetsonTX24GB},
	{"p3489-0000", boardtype.JetsonTX2i},
	{"p3310-1000", boardtype.JetsonTX2},

	{"p2180-1000", boardtype.JetsonTX1},
	{"p2371-2180", boardtype.JetsonTX1},

	{"p2894-0050", boardtype.ShieldTV},

	{"p3904-0000", boardtype.ClaraAGX},
}

var jetsonModulesByDeviceTreeBaseModel = []jetson{
	{"NVIDIA Jetson Orin NX Engineering Reference Developer Kit", boardtype.JetsonOrinNX16GB},
	{"NVIDIA Jetson Orin Nano Developer Kit", boardtype.JetsonOrinNanoDeveloperKit},
	{"NVIDIA Jetson TX2 Developer Kit", boardtype.JetsonTX2},
	{"NVIDIA Jetson TX2", boardtype.JetsonTX2},
	{"NVIDIA Jetson TX2 NX Developer Kit", boardtype.JetsonTX2NX},
	{"NVIDIA Jetson AGX Xavier", boardtype.JetsonAGXXavier},
	{"NVIDIA Jetson AGX Xavier Developer Kit", boardtype.JetsonAGXXavier},
	{"NVIDIA Jetson Xavier NX Developer Kit (SD-card)", boardtype.JetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX Developer Kit (eMMC)", boardtype.JetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX (SD-card)", boardtype.JetsonXavierNXDeveloperKit},
	{"NVIDIA Jetson Xavier NX (eMMC)", boardtype.JetsonXavierNX8GB},
	{"NVIDIA Jetson TX1", boardtype.JetsonTX1},
	{"NVIDIA Jetson TX1 Developer Kit", boardtype.JetsonTX1},
	{"NVIDIA Shield TV", boardtype.ShieldTV},
	{"NVIDIA Jetson Nano Developer Kit", boardtype.JetsonNanoDeveloperKit},
	{"NVIDIA Jetson AGX Orin Developer Kit", boardtype.JetsonAGXOrin},
	{"NVIDIA Jetson AGX Orin", boardtype.JetsonAGXOrin},
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

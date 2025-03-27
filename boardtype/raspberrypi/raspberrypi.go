package raspberrypi

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"log/slog"

	"github.com/rinzlerlabs/sbcidentify/boardtype"
	"github.com/rinzlerlabs/sbcidentify/identifier"
)

func init() {
	identifier.RegisterBoardIdentifier(NewRaspberryPiIdentifier)
}

var (
	ErrVcgencmdNotFound    = errors.New("vcgencmd not found")
	ErrInvalidMeminfo      = errors.New("invalid meminfo")
	ErrCannotIdentifyBoard = errors.New("cannot identify Raspberry Pi board")
	execLookPath           = exec.LookPath
)

type raspberryPi struct {
	Model    string
	Memory   int
	Type     boardtype.SBC
	Fallback boardtype.SBC
}

var raspberryPiModels = []raspberryPi{
	{"Raspberry Pi 3 Model B", 1024, boardtype.RaspberryPi3B, boardtype.RaspberryPi3B},
	{"Raspberry Pi 3 Model A", 512, boardtype.RaspberryPi3APlus, boardtype.RaspberryPi3APlus},
	{"Raspberry Pi 3 Model B", 1024, boardtype.RaspberryPi3BPlus, boardtype.RaspberryPi3BPlus},
	{"Raspberry Pi 4 Model B", 1024, boardtype.RaspberryPi4B1GB, boardtype.RaspberryPi4B},
	{"Raspberry Pi 4 Model B", 2048, boardtype.RaspberryPi4B2GB, boardtype.RaspberryPi4B},
	{"Raspberry Pi 4 Model B", 4096, boardtype.RaspberryPi4B4GB, boardtype.RaspberryPi4B},
	{"Raspberry Pi 4 Model B", 8192, boardtype.RaspberryPi4B8GB, boardtype.RaspberryPi4B},
	{"Raspberry Pi 400", 4096, boardtype.RaspberryPi4400, boardtype.RaspberryPi4400},
	{"Raspberry Pi 5 Model B", 2048, boardtype.RaspberryPi5B2GB, boardtype.RaspberryPi5B},
	{"Raspberry Pi 5 Model B", 4096, boardtype.RaspberryPi5B4GB, boardtype.RaspberryPi5B},
	{"Raspberry Pi 5 Model B", 8192, boardtype.RaspberryPi5B8GB, boardtype.RaspberryPi5B},
}

func NewRaspberryPiIdentifier(logger *slog.Logger) identifier.BoardIdentifier {
	logger.Debug("initializing Raspberry Pi identifier")
	newLogger := logger.With(slog.String("source", "RaspberryPiIdentifier"))
	return &raspberryPiIdentifier{logger: newLogger}
}

type raspberryPiIdentifier struct {
	logger *slog.Logger
}

func (r raspberryPiIdentifier) Name() string {
	return "Raspberry Pi Identifier"
}

func (r raspberryPiIdentifier) GetBoardType() (boardtype.SBC, error) {
	r.logger.Debug("getting board type")
	dtbm, err := identifier.GetDeviceTreeBaseModel(r.logger)
	if err == identifier.ErrCannotIdentifyBoard {
		dtbm, err = identifier.GetDeviceTreeModel(r.logger)
		if err == identifier.ErrCannotIdentifyBoard {
			return nil, ErrCannotIdentifyBoard
		} else if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	r.logger.Debug("device tree model", slog.String("model", dtbm))
	subModels := make([]raspberryPi, 0)
	for _, m := range raspberryPiModels {
		if strings.Contains(dtbm, m.Model) {
			subModels = append(subModels, m)
		}
	}
	if len(subModels) == 0 {
		return nil, ErrCannotIdentifyBoard
	}
	ramMb, err := getInstalledRAM(r.logger)
	if err == ErrVcgencmdNotFound {
		r.logger.Debug("vcgencmd not found, using fallback", slog.String("model", dtbm), slog.Int("ram", ramMb), slog.Any("fallback", subModels[0].Fallback))
		return subModels[0].Fallback, nil
	} else if err != nil {
		return nil, err
	}
	for _, m := range subModels {
		if m.Memory == ramMb {
			return m.Type, nil
		}
	}
	r.logger.Debug("no matching model found, using fallback", slog.String("model", dtbm), slog.Int("ram", ramMb), slog.Int("subModels", len(subModels)), slog.Any("subModels", subModels), slog.Any("fallback", subModels[0].Fallback))
	return subModels[0].Fallback, nil
}

func getInstalledRAM(logger *slog.Logger) (int, error) {
	if _, err := execLookPath("vcgencmd"); err != nil {
		logger.Debug("vcgencmd not found", slog.Any("error", err))
		return 0, ErrVcgencmdNotFound
	}
	out, err := exec.Command("vcgencmd", "get_config", "total_mem").Output()
	if err != nil {
		return 0, err
	}
	output := strings.TrimSpace(string(out))
	return parseVcgencmdMemoryOutput(logger, output)
}

func parseVcgencmdMemoryOutput(logger *slog.Logger, output string) (int, error) {
	logger.Debug("vcgencmd output", slog.String("output", output))
	parts := strings.Split(output, "=")
	if len(parts) != 2 {
		return 0, ErrInvalidMeminfo
	}
	installedRam, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		logger.Debug("Failed to parse RAM", slog.String("output", output), slog.Any("error", err))
		return 0, ErrInvalidMeminfo
	}
	logger.Debug("Parsed RAM", slog.Int("total_mem", installedRam))
	return installedRam, nil
}

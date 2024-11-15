package raspberrypi

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"log/slog"

	"github.com/thegreatco/sbcidentify/boardtype"
	"github.com/thegreatco/sbcidentify/identifier"
)

func init() {
	identifier.RegisterBoardIdentifier(NewRaspberryPiIdentifier)
}

type RaspberryPi struct{ boardtype.BoardType }
type RaspberryPi3 struct{ RaspberryPi }
type RaspberryPi4 struct{ RaspberryPi }
type RaspberryPi5 struct{ RaspberryPi }

var (
	RaspberryPi3B     = RaspberryPi3{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "3B", RAM: 1024}}}
	RaspberryPi3APlus = RaspberryPi3{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "3A+", RAM: 512}}}
	RaspberryPi3BPlus = RaspberryPi3{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "3B+", RAM: 1024}}}
	RaspberryPi4B1GB  = RaspberryPi4{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4B", RAM: 1024}}}
	RaspberryPi4B2GB  = RaspberryPi4{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4B", RAM: 2048}}}
	RaspberryPi4B4GB  = RaspberryPi4{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4B", RAM: 4096}}}
	RaspberryPi4B8GB  = RaspberryPi4{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4B", RAM: 8192}}}
	RaspberryPi4B     = RaspberryPi4{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4B", RAM: 0}}}
	RaspberryPi4400   = RaspberryPi4{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "4 400", RAM: 4096}}}
	RaspberryPi5B     = RaspberryPi5{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "5B", RAM: 0}}}
	RaspberryPi5B2GB  = RaspberryPi5{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "5B", RAM: 2048}}}
	RaspberryPi5B4GB  = RaspberryPi5{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "5B", RAM: 4096}}}
	RaspberryPi5B8GB  = RaspberryPi5{RaspberryPi{boardtype.BoardType{Manufacturer: "Raspberry Pi", Model: "Raspberry Pi", SubModel: "5B", RAM: 8192}}}
)

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
	{"Raspberry Pi 3 Model B", 1024, RaspberryPi3B, RaspberryPi3B},
	{"Raspberry Pi 3 Model A", 512, RaspberryPi3APlus, RaspberryPi3APlus},
	{"Raspberry Pi 3 Model B", 1024, RaspberryPi3BPlus, RaspberryPi3BPlus},
	{"Raspberry Pi 4 Model B", 1024, RaspberryPi4B1GB, RaspberryPi4B},
	{"Raspberry Pi 4 Model B", 2048, RaspberryPi4B2GB, RaspberryPi4B},
	{"Raspberry Pi 4 Model B", 4096, RaspberryPi4B4GB, RaspberryPi4B},
	{"Raspberry Pi 4 Model B", 8192, RaspberryPi4B8GB, RaspberryPi4B},
	{"Raspberry Pi 400", 4096, RaspberryPi4400, RaspberryPi4400},
	{"Raspberry Pi 5 Model B", 2048, RaspberryPi5B2GB, RaspberryPi5B},
	{"Raspberry Pi 5 Model B", 4096, RaspberryPi5B4GB, RaspberryPi5B},
	{"Raspberry Pi 5 Model B", 8192, RaspberryPi5B8GB, RaspberryPi5B},
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

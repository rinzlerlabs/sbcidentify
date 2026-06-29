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
	var subModels []raspberryPi
	for _, m := range raspberryPiModels {
		if strings.Contains(dtbm, m.Model) {
			subModels = append(subModels, m)
		}
	}
	if len(subModels) > 1 {
		best := 0
		for _, m := range subModels {
			if len(m.Model) > best {
				best = len(m.Model)
			}
		}
		filtered := subModels[:0]
		for _, m := range subModels {
			if len(m.Model) == best {
				filtered = append(filtered, m)
			}
		}
		subModels = filtered
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
	var ramMatches []raspberryPi
	for _, m := range subModels {
		if m.Memory == ramMb {
			ramMatches = append(ramMatches, m)
		}
	}
	if len(ramMatches) > 1 {
		names := make([]string, len(ramMatches))
		for i, m := range ramMatches {
			names[i] = m.Type.GetPrettyName()
		}
		r.logger.Warn("multiple RAM matches, using first", slog.String("model", dtbm), slog.Int("ram_mb", ramMb), slog.Any("matches", names))
	}
	if len(ramMatches) > 0 {
		return ramMatches[0].Type, nil
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

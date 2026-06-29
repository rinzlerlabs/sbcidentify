package nvidia

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
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

type jetson struct {
	Model string
	Type  boardtype.SBC
}

var (
	ErrDtsFileDoesNotExist = errors.New("DTS file does not exist")
	ErrCannotIdentifyBoard = errors.New("cannot identify NVIDIA board")
)

// memInfoPath is a var so tests can swap it for a temp file.
var memInfoPath = "/proc/meminfo"

// jetsonRAMRefinementCandidates lists boards where detection returns a generic
// result (RAM==0) and installed RAM can narrow it down. Currently covers the
// AGX Orin Developer Kit, which uses module ID p3701-0000 for all RAM configs.
var jetsonRAMRefinementCandidates = []boardtype.SBC{
	boardtype.JetsonAGXOrin32GB,
	boardtype.JetsonAGXOrin64GB,
}

func getInstalledRAMMB(logger *slog.Logger) (int, error) {
	data, err := os.ReadFile(memInfoPath)
	if err != nil {
		return 0, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "MemTotal:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return 0, fmt.Errorf("unexpected MemTotal format: %s", line)
		}
		kb, err := strconv.Atoi(fields[1])
		if err != nil {
			return 0, fmt.Errorf("failed to parse MemTotal: %w", err)
		}
		return kb / 1024, nil
	}
	return 0, fmt.Errorf("MemTotal not found in %s", memInfoPath)
}

// refineByInstalledRAM upgrades a generic board (RAM==0) to the most specific
// known variant by reading actual installed RAM from /proc/meminfo and picking
// the closest match from jetsonRAMRefinementCandidates.
func refineByInstalledRAM(logger *slog.Logger, board boardtype.SBC) boardtype.SBC {
	if board.GetRAM() != 0 {
		return board
	}
	ram, err := getInstalledRAMMB(logger)
	if err != nil {
		logger.Debug("cannot read installed RAM, using generic board type", slog.Any("error", err))
		return board
	}
	logger.Debug("installed RAM", slog.Int("mb", ram))

	var best boardtype.SBC
	bestDiff := int(^uint(0) >> 1)
	for _, candidate := range jetsonRAMRefinementCandidates {
		if candidate.GetRAM() == 0 || !candidate.IsBoardType(board) {
			continue
		}
		diff := candidate.GetRAM() - ram
		if diff < 0 {
			diff = -diff
		}
		if diff < bestDiff {
			best = candidate
			bestDiff = diff
		}
	}
	if best != nil {
		logger.Debug("refined board type by installed RAM", slog.String("type", best.GetPrettyName()))
		return best
	}
	return board
}

// collectMatches returns all unique board types from table whose Model is a
// substring of input. Deduplicates by pretty name so overlapping patterns that
// resolve to the same board don't produce a spurious warning.
func collectMatches(logger *slog.Logger, input string, table []jetson) []boardtype.SBC {
	seen := make(map[string]bool)
	var matches []boardtype.SBC
	for _, m := range table {
		if strings.Contains(input, m.Model) {
			key := m.Type.GetPrettyName()
			if !seen[key] {
				seen[key] = true
				matches = append(matches, m.Type)
			}
		}
	}
	if len(matches) > 1 {
		names := make([]string, len(matches))
		for i, m := range matches {
			names[i] = m.GetPrettyName()
		}
		logger.Warn("multiple board matches, using first", slog.String("input", input), slog.Any("matches", names))
	}
	return matches
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
	} else if err == identifier.ErrCannotIdentifyBoard {
		r.logger.Debug("unknown board from DTS, falling back to device tree base model")
		boardType, err = getBoardTypeByDeviceTreeBaseModel(r.logger)
	}
	if err == identifier.ErrCannotIdentifyBoard || err == ErrCannotIdentifyBoard {
		r.logger.Debug("unknown board")
		return nil, ErrCannotIdentifyBoard
	}
	if err != nil {
		r.logger.Debug("error getting board type", slog.Any("error", err))
		return nil, err
	}
	boardType = refineByInstalledRAM(r.logger, boardType)
	r.logger.Debug("board type", slog.String("type", boardType.GetPrettyName()))
	return boardType, nil
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
	matches := collectMatches(logger, moduleModel, jetsonModulesByModelNumber)
	if len(matches) == 0 {
		return nil, identifier.ErrCannotIdentifyBoard
	}
	return matches[0], nil
}

func getBoardTypeByDeviceTreeBaseModel(logger *slog.Logger) (boardtype.SBC, error) {
	dtbm, err := identifier.GetDeviceTreeBaseModel(logger)
	if err != nil {
		return nil, err
	}
	matches := collectMatches(logger, dtbm, jetsonModulesByDeviceTreeBaseModel)
	if len(matches) == 0 {
		logger.Debug("device tree base model does not match any boards", slog.String("model", dtbm))
		return nil, ErrCannotIdentifyBoard
	}
	return matches[0], nil
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

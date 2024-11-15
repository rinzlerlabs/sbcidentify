package sbcidentify

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	procDeviceTreeModelFile     = "/proc/device-tree/model"
	firmwareDeviceTreeModelFile = "/sys/firmware/devicetree/base/model"
	socIdFile                   = "/sys/devices/soc0/soc_id"
	dtsFileName                 = "/proc/device-tree/nvidia,dtsfilename"
)

var (
	ErrInvalidMeminfo      = errors.New("invalid meminfo")
	ErrCannotIdentifyBoard = errors.New("cannot identify board")
	ErrVcgencmdNotFound    = errors.New("vcgencmd not found")
)

func getDeviceTreeBaseModel() (string, error) {
	if _, err := os.Stat(firmwareDeviceTreeModelFile); err != nil {
		logger.Debug("cannot read firmware device tree model file", slog.Any("error", err))
		return "", ErrCannotIdentifyBoard
	}
	c, err := os.ReadFile(firmwareDeviceTreeModelFile)
	if err != nil {
		logger.Debug("cannot read firmware device tree model file", slog.Any("error", err))
		return "", ErrCannotIdentifyBoard
	}
	str := strings.TrimSuffix(strings.TrimSpace(string(c)), "\x00")
	logger.Debug("firmware device tree model", slog.String("model", str))
	return str, nil
}

func getDeviceTreeModel() (string, error) {
	if _, err := os.Stat(procDeviceTreeModelFile); err != nil {
		logger.Debug("cannot read proc device tree model file", slog.Any("error", err))
		return "", ErrCannotIdentifyBoard
	}
	c, err := os.ReadFile(procDeviceTreeModelFile)
	if err != nil {
		logger.Debug("cannot read proc device tree model file", slog.Any("error", err))
		return "", ErrCannotIdentifyBoard
	}
	str := strings.TrimSpace(string(c))
	logger.Debug("proc device tree model", slog.String("model", str))
	return str, nil
}

// func getSoCId() (int, error) {
// 	c, err := os.ReadFile("/sys/devices/soc0/soc_id")
// 	if err != nil {
// 		return 0, err
// 	}
// 	str := string(c)
// 	return strconv.Atoi(str)
// }

func getInstalledRAM() (int, error) {
	if _, err := exec.LookPath("vcgencmd"); err != nil {
		logger.Debug("vcgencmd not found", slog.Any("error", err))
		return 0, ErrVcgencmdNotFound
	}
	out, err := exec.Command("vcgencmd", "get_config", "total_mem").Output()
	if err != nil {
		return 0, err
	}
	output := strings.TrimSpace(string(out))
	logger.Debug("vcgencmd output", slog.String("output", output))
	parts := strings.Split(output, "=")
	if len(parts) != 2 {
		return 0, ErrInvalidMeminfo
	}
	installedRam, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	logger.Debug("Parsed RAM", slog.Int("total_mem", installedRam))
	return installedRam, nil
}

func getDtsFile() (string, error) {
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

func getModuleNameFromDtsFilename(dtsFilename string) (string, error) {
	filename := filepath.Base(dtsFilename)
	ret := strings.TrimSuffix(filename, filepath.Ext(filename))
	logger.Debug("module name", slog.String("name", ret))
	return ret, nil
}

func getModuleModelFromModuleName(moduleName string) (string, error) {
	parts := strings.Split(moduleName, "-")
	if len(parts) >= 4 {
		ret := strings.Join(parts[1:3], "-")
		logger.Debug("module model", slog.String("model", ret))
		return ret, nil
	}
	logger.Debug("error parsing module name", slog.String("name", moduleName))
	return "", ErrUnknownBoard
}

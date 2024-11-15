package identifier

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const (
	procDeviceTreeModelFile     = "/proc/device-tree/model"
	firmwareDeviceTreeModelFile = "/sys/firmware/devicetree/base/model"
	socIdFile                   = "/sys/devices/soc0/soc_id"
)

var (
	ErrCannotIdentifyBoard = errors.New("cannot identify board")
)

func GetDeviceTreeBaseModel(logger *slog.Logger) (string, error) {
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

func GetDeviceTreeModel(logger *slog.Logger) (string, error) {
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

func GetSoCId() (int, error) {
	c, err := os.ReadFile("/sys/devices/soc0/soc_id")
	if err != nil {
		return 0, err
	}
	str := string(c)
	return strconv.Atoi(str)
}

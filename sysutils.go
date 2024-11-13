package sbcidentify

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	ErrInvalidMeminfo = errors.New("invalid meminfo")
)

func getDeviceTreeBaseModel() (string, error) {
	c, err := os.ReadFile("/sys/firmware/devicetree/base/model")
	if err != nil {
		return "", err
	}
	str := string(c)
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
	out, err := exec.Command("vcgencmd", "get_config", "total_mem").Output()
	if err != nil {
		return 0, err
	}
	output := strings.TrimSpace(string(out))
	parts := strings.Split(output, "=")
	if len(parts) != 2 {
		return 0, ErrInvalidMeminfo
	}
	return strconv.Atoi(parts[1])
}

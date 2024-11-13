package sbcidentify

import (
	"os"
	"strings"
)

type BoardType string

var (
	BoardTypeUnknown      BoardType = "Unknown"
	BoardTypeRaspberryPi4 BoardType = "Raspberry Pi 4"
	BoardTypeRaspberryPi5 BoardType = "Raspberry Pi 5"
)

func GetBoardType() (BoardType, error) {
	c, err := os.ReadFile("/sys/firmware/devicetree/base/model")
	if err != nil {
		return BoardTypeUnknown, err
	}
	str := string(c)
	if strings.HasPrefix(str, string(BoardTypeRaspberryPi4)) {
		return BoardTypeRaspberryPi4, nil
	}

	if strings.HasPrefix(str, string(BoardTypeRaspberryPi5)) {
		return BoardTypeRaspberryPi5, nil
	}

	return BoardTypeUnknown, nil
}

func IsBoardType(boardType BoardType) bool {
	board, err := GetBoardType()
	if err != nil {
		return false
	}
	return board == boardType
}

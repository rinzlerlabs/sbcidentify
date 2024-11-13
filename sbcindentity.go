package sbcidentify

import (
	"os"
	"strings"
)

type BoardType string

var (
	BoardTypeUnknown        BoardType = "Unknown"
	BoardTypeRaspberryPi4   BoardType = "Raspberry Pi 4"
	BoardTypeRaspberryPi5   BoardType = "Raspberry Pi 5"
	BoardTypeJetsonNano     BoardType = "NVIDIA Jetson Nano"
	BoardTypeJetsonOrinNano BoardType = "NVIDIA Orin Nano"
	BoardTypeJersonOrin     BoardType = "NVIDIA AGX Orin"
)

var boards = []BoardType{
	BoardTypeRaspberryPi4,
	BoardTypeRaspberryPi5,
	BoardTypeJetsonOrinNano,
}

func GetBoardType() (BoardType, error) {
	c, err := os.ReadFile("/sys/firmware/devicetree/base/model")
	if err != nil {
		return BoardTypeUnknown, err
	}
	str := string(c)
	for _, board := range boards {
		if strings.HasPrefix(str, string(board)) {
			return board, nil
		}
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

package sbcidentify

import (
	"errors"
)

const (
	BoardTypeUnknown BoardType = "Unknown"
)

var (
	ErrUnknownBoard error = errors.New("unknown board")
)

type boardIdentifier interface {
	GetBoardType() (BoardType, error)
}

var boardIdentifiers = []boardIdentifier{
	raspberryPiIdentifier{},
	jetsonIdentifier{},
}

func GetBoardType() (BoardType, error) {
	for _, identifier := range boardIdentifiers {
		board, err := identifier.GetBoardType()
		if err != nil {
			continue
		}
		return board, nil
	}
	return BoardTypeUnknown, ErrUnknownBoard
}

func IsBoardType(boardType BoardType) bool {
	board, err := GetBoardType()
	if err != nil {
		return false
	}
	return board == boardType
}

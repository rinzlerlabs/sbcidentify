package sbcidentify

import (
	"errors"
	"log/slog"
	"os"
)

const (
	BoardTypeUnknown BoardType = "Unknown"
)

var (
	ErrUnknownBoard error          = errors.New("unknown board")
	logLevel        *slog.LevelVar = new(slog.LevelVar)
	logger          *slog.Logger   = slog.New(NewLogHandler(os.Stderr, &HandlerConfig{Level: logLevel}))
)

func SetLogLevel(level slog.Level) {
	logLevel.Set(level)
}

func SetLogger(l *slog.Logger) {
	logger = l
}

type boardIdentifier interface {
	GetBoardType() (BoardType, error)
}

var boardIdentifiers = []boardIdentifier{
	jetsonIdentifier{},
	raspberryPiIdentifier{},
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

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
	logger          *slog.Logger   = slog.New(NewLogHandler(os.Stderr, &HandlerConfig{Level: logLevel})).With("source", "sbcidentify")
)

func SetLogLevel(level slog.Level) {
	logLevel.Set(level)
}

func SetLogger(l *slog.Logger) {
	logger = l
}

type boardIdentifier interface {
	Name() string
	GetBoardType() (BoardType, error)
}

func GetBoardType() (BoardType, error) {
	boardIdentifiers := []boardIdentifier{
		NewJetsonIdentifier(logger),
		NewRaspberryPiIdentifier(logger),
	}
	var final error
	for _, identifier := range boardIdentifiers {
		board, err := identifier.GetBoardType()
		if err != nil {
			final = errors.Join(final, err)
			continue
		}
		return board, nil
	}
	return BoardTypeUnknown, final
}

func IsBoardType(boardType BoardType) bool {
	board, err := GetBoardType()
	if err != nil {
		return false
	}
	return board == boardType
}

package sbcidentify

import (
	"errors"
	"log/slog"
	"os"

	"github.com/thegreatco/sbcidentify/boardtype"
	"github.com/thegreatco/sbcidentify/identifier"
	_ "github.com/thegreatco/sbcidentify/nvidia"
	_ "github.com/thegreatco/sbcidentify/raspberrypi"
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

func GetBoardType() (boardtype.SBC, error) {
	boardIdentifiers := identifier.BuildIdentifiers(logger)
	var final error
	for _, identifier := range boardIdentifiers {
		board, err := identifier.GetBoardType()
		if err != nil {
			final = errors.Join(final, err)
			continue
		}
		return board, nil
	}
	return nil, final
}

func IsBoardType(boardType boardtype.SBC) bool {
	board, err := GetBoardType()
	if err != nil {
		return false
	}
	return board == boardType
}

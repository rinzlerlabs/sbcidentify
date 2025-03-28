package sbcidentify

import (
	"errors"
	"log/slog"
	"os"

	"github.com/rinzlerlabs/sbcidentify/boardtype"
	"github.com/rinzlerlabs/sbcidentify/identifier"

	_ "github.com/rinzlerlabs/sbcidentify/boardtype/nvidia"
	_ "github.com/rinzlerlabs/sbcidentify/boardtype/raspberrypi"
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
	if len(boardIdentifiers) == 0 {
		panic("no board identifiers found")
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
	return nil, final
}

func IsBoardType(boardType boardtype.SBC) bool {
	board, err := GetBoardType()
	if err != nil {
		return false
	}
	if board == nil {
		logger.Debug("board is nil, this is unexpected")
		return false
	}
	return board.IsBoardType(boardType)
}

func IsRaspberryPi() bool {
	return IsBoardType(boardtype.RaspberryPi)
}

func IsNvidia() bool {
	return IsBoardType(boardtype.NVIDIA)
}

func IsJetson() bool {
	return IsBoardType(boardtype.Jetson)
}

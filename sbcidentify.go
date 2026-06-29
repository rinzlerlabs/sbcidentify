package sbcidentify

import (
	"errors"
	"log/slog"
	"sync/atomic"

	"github.com/rinzlerlabs/sbcidentify/boardtype"
	"github.com/rinzlerlabs/sbcidentify/identifier"

	_ "github.com/rinzlerlabs/sbcidentify/boardtype/nvidia"
	_ "github.com/rinzlerlabs/sbcidentify/boardtype/raspberrypi"
)

var (
	ErrUnknownBoard error = errors.New("unknown board")
	logger          atomic.Pointer[slog.Logger]
)

// SetLogger gives the library a specific logger. If never called, all internal
// logging goes through slog.Default(), so the application's slog.SetDefault
// configuration is automatically respected with no library-side setup required.
func SetLogger(l *slog.Logger) {
	logger.Store(l)
}

func getLogger() *slog.Logger {
	if l := logger.Load(); l != nil {
		return l
	}
	return slog.Default()
}

func GetBoardType() (boardtype.SBC, error) {
	boardIdentifiers := identifier.BuildIdentifiers(getLogger())
	if len(boardIdentifiers) == 0 {
		panic("no board identifiers found")
	}
	var errs error
	var matches []boardtype.SBC
	for _, id := range boardIdentifiers {
		board, err := id.GetBoardType()
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		matches = append(matches, board)
	}
	if len(matches) == 0 {
		return nil, errs
	}
	if len(matches) > 1 {
		names := make([]string, len(matches))
		for i, m := range matches {
			names[i] = m.GetPrettyName()
		}
		getLogger().Warn("multiple board identifiers matched, using first", slog.Any("matches", names))
	}
	return matches[0], nil
}

func IsBoardType(boardType boardtype.SBC) bool {
	board, err := GetBoardType()
	if err != nil {
		return false
	}
	if board == nil {
		getLogger().Debug("board is nil, this is unexpected")
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

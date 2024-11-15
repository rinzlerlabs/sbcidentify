package identifier

import (
	"log/slog"

	boardType "github.com/thegreatco/sbcidentify/boardtype"
)

type BoardIdentifier interface {
	Name() string
	GetBoardType() (boardType.SBC, error)
}

var identifiers []func(*slog.Logger) BoardIdentifier = make([]func(*slog.Logger) BoardIdentifier, 0)

func RegisterBoardIdentifier(identifier func(*slog.Logger) BoardIdentifier) {
	identifiers = append(identifiers, identifier)
}

func BuildIdentifiers(logger *slog.Logger) []BoardIdentifier {
	ids := make([]BoardIdentifier, 0)
	for _, id := range identifiers {
		ids = append(ids, id(logger))
	}
	return ids
}

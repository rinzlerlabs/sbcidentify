package test

import (
	"os"
	"testing"

	"github.com/rinzlerlabs/sbcidentify"
	"github.com/rinzlerlabs/sbcidentify/boardtype"
)

type test struct {
	requiresRoot      bool
	requiredBoardType *boardtype.SBC
	requiresSbc       bool
}

func Test() *test {
	return &test{}
}

func (t *test) RequiresRoot() *test {
	t.requiresRoot = true
	return t
}

func (t *test) RequiresBoardType(boardType boardtype.SBC) *test {
	t.requiredBoardType = &boardType
	return t
}

func (t *test) RequiresSbc() *test {
	t.requiresSbc = true
	return t
}

func (t *test) ShouldSkip(test *testing.T) {
	if t.requiresSbc {
		boardType, err := sbcidentify.GetBoardType()
		if err != nil {
			test.Log(err)
		}
		if boardType == nil {
			test.Skip("Test requires physical SBC, not running on SBC")
		}
	}

	if t.requiredBoardType != nil && !sbcidentify.IsBoardType(*t.requiredBoardType) {
		test.Skipf("Test requires board type %v", (*t.requiredBoardType).GetPrettyName())
	}
	if t.requiresRoot && !IsRoot() {
		test.Skip("Test requires root, not running as root")
	}
}

func IsRoot() bool {
	return os.Geteuid() == 0
}

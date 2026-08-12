package os

import (
	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/context"
)

func OpenFile(ctx *context.Context, path string, rw bits.Flags) (Handle, bool) {
	return OpenOrCreateFile(ctx, path, rw, 0, 0)
}

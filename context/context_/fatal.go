package context_

import (
	"os"

	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/fmt/fmt_"
)

func Fatal(ctx *context.Context, msg string) {
	fmt_.Eprintln(ctx, ctx.Fmt.Reset().S(msg).S(": ").S(ctx.Error()))
	os.Exit(1)
}

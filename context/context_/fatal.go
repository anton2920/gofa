package context_

import (
	"os"

	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/fmt/fmt_"
)

func Fatal(ctx *context.Context, msg string) {
	ctx.Fmt.Reset().S(msg)
	if len(ctx.Error()) > 0 {
		ctx.Fmt.S(": ").S(ctx.Error())
	}
	if ctx.ErrorCode() > 0 {
		ctx.Fmt.S(" (code=").D(ctx.ErrorCode()).S(")")
	}
	fmt_.Eprintln(ctx, &ctx.Fmt)
	os.Exit(1)
}

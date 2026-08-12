package context_

import "github.com/anton2920/gofa/context"

func Must(ctx *context.Context, ok bool, msg string) {
	if !ok {
		Fatal(ctx, msg)
	}
}

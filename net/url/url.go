package url

import (
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/trace/trace_"
)

type URL struct {
	Path     Path
	RawQuery string

	Query Values
}

func (u *URL) ParseQuery(ctx *context.Context) bool {
	t := trace_.Begin("")

	if len(u.Query.Keys) > 0 {
		trace_.End(t)
		return true
	}

	ok := ParseQuery(ctx, &u.Query, u.RawQuery)

	trace_.End(t)
	return ok
}

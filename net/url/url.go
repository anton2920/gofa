package url

import (
	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/trace"
)

type URL struct {
	Path     Path
	RawQuery string

	Query Values
}

func (u *URL) ParseQuery(arena *alloc.Arena) error {
	t := trace.Begin("")

	if len(u.Query.Keys) > 0 {
		trace.End(t)
		return nil
	}

	err := ParseQuery(arena, &u.Query, u.RawQuery)

	trace.End(t)
	return err
}

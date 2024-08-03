package url

import "github.com/anton2920/gofa/prof"

type URL struct {
	Path     string
	RawQuery string

	Query Values
}

func (u *URL) ParseQuery() error {
	p := prof.Begin("")

	if len(u.Query.Keys) != 0 {
		prof.End(p)
		return nil
	}

	err := ParseQuery(&u.Query, u.RawQuery)

	prof.End(p)
	return err
}

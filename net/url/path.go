package url

import (
	"fmt"

	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
)

type Path string

/* Match returns `true` if `p` matches format described in `format`. Additionally it slices `p` with the length of the matched string. */
func (p *Path) Match(format string, args ...interface{}) bool {
	t := trace.Begin("")

	var narg int
	var ok bool

	path := string(*p)
	for {
		percent := strings.FindChar(format, '%')
		if percent == -1 {
			const ellipsis = "..."

			if !strings.EndsWith(format, ellipsis) {
				ok = path == format
			} else {
				format = format[:len(format)-len(ellipsis)]
				ok = strings.StartsWith(path, format)
			}
			if ok {
				*p = Path(path[len(format):])
			}

			trace.End(t)
			return ok
		}

		match := format[:percent]
		if !strings.StartsWith(path, match) {
			trace.End(t)
			return false
		}
		path = path[len(match):]
		format = format[len(match):]

		slashP := strings.FindChar(path, '/')
		if slashP == -1 {
			slashP = len(path)
		}

		slashF := strings.FindChar(format, '/')
		if slashF == -1 {
			slashF = len(format)
		}

		n, err := fmt.Sscanf(path[:slashP], format[:slashF], args[narg:]...)
		if (n == 0) && (err != nil) {
			trace.End(t)
			return false
		}
		narg += n

		path = path[slashP:]
		format = format[2:]
	}
}

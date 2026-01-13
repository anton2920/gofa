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

	const ellipsis = "..."

	var narg int
	var ok bool

	path := string(*p)
	for {
		percent := strings.FindChar(format, '%')
		if percent == -1 {
			if !strings.EndsWith(format, ellipsis) {
				ok = path == format
			} else {
				format = format[:len(format)-len(ellipsis)]
				ok = strings.StartsWith(path, format)
			}

			//runtime.Breakpoint()
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

		/* NOTE(anton2920): this assumes that format strings are /%[a-z]/. */
		nextF := percent - len(match) + 2

		var nextP int
		if nextF >= len(format) {
			nextP = len(path)
		} else {
			nextP = strings.FindChar(path, format[nextF])
			if (nextP == -1) || (strings.StartsWith(format[nextF:], ellipsis)) {
				nextP = strings.FindChar(path, '/')
				if nextP == -1 {
					nextP = len(path)
				}
			}
		}

		/* TODO(anton2920): moved to heap: fmt.str. */
		n, err := fmt.Sscanf(path[:nextP], format[:nextF], args[narg:]...)
		if (n == 0) && (err != nil) {
			trace.End(t)
			return false
		}
		narg += n

		path = path[nextP:]
		format = format[nextF:]
	}
}

package errors

import "github.com/anton2920/gofa/trace"

func Or(errs ...error) error {
	t := trace.Begin("")

	for i := 0; i < len(errs); i++ {
		if errs[i] != nil {
			trace.End(t)
			return errs[i]
		}
	}

	trace.End(t)
	return nil
}

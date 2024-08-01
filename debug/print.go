//go:build gofadebug

package debug

import (
	"fmt"
	"os"
)

func Printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	if format[len(format)-1] != '\n' {
		os.Stderr.Write([]byte{'\n'})
	}
}

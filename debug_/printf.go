//go:build gofadebug
// +build gofadebug

package debug_

import (
	"fmt"

	"github.com/anton2920/gofa/bytes"
)

func Printf(format string, args ...interface{}) (int, error) {
	buf := make([]byte, 4096)
	n := copy(buf, format)
	if buf[n-1] != '\n' {
		buf[n] = '\n'
	}
	return fmt.Printf(bytes.AsString(buf), args...)
}

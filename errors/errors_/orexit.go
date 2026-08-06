package errors_

import (
	"github.com/anton2920/gofa/fmt"
	"github.com/anton2920/gofa/fmt/fmt_"
	"github.com/anton2920/gofa/os"
)

func OrExit(msg string, ierr error) {
	if ierr != nil {
		var f fmt.Formatter
		f.InitWithByteSlice(make([]byte, 4096))

		fmt_.Eprintln(f.S(msg).S(ierr.Error()))
		os.Exit(1)
	}
}

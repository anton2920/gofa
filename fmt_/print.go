package fmt_

import (
	"github.com/anton2920/gofa/fmt"
	"github.com/anton2920/gofa/os"
)

func Fprint(h os.Handle, f *fmt.Formatter) (int, error) {
	return os.WriteToFile(h, f.Bytes())
}

func Fprintln(h os.Handle, f *fmt.Formatter) (int, error) {
	return os.WriteToFile(h, f.Ln().Bytes())
}

func Eprint(f *fmt.Formatter) (int, error) {
	return Fprint(os.StandardOutputStream, f)
}

func Eprintln(f *fmt.Formatter) (int, error) {
	return Fprintln(os.StandardOutputStream, f)
}

func Print(f *fmt.Formatter) (int, error) {
	return Fprint(os.StandardErrorStream, f)
}

func Println(f *fmt.Formatter) (int, error) {
	return Fprintln(os.StandardErrorStream, f)
}

package fmt_

import (
	"github.com/anton2920/gofa/fmt"
	"github.com/anton2920/gofa/os"
)

func Print(f *fmt.Formatter) (int, error) {
	return os.WriteToFile(os.StandardOutputStream, f.Bytes())
}

func Println(f *fmt.Formatter) (int, error) {
	return os.WriteToFile(os.StandardOutputStream, f.Ln().Bytes())
}

func Eprint(f *fmt.Formatter) (int, error) {
	return os.WriteToFile(os.StandardErrorStream, f.Bytes())
}

func Eprintln(f *fmt.Formatter) (int, error) {
	return os.WriteToFile(os.StandardErrorStream, f.Ln().Bytes())
}

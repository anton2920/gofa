package fmt_

import (
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/fmt"
	"github.com/anton2920/gofa/os"
)

func Fprint(ctx *context.Context, h os.Handle, f *fmt.Formatter) (int, bool) {
	return os.WriteToFile(ctx, h, f.Bytes())
}

func Fprintln(ctx *context.Context, h os.Handle, f *fmt.Formatter) (int, bool) {
	return os.WriteToFile(ctx, h, f.Ln().Bytes())
}

func Eprint(ctx *context.Context, f *fmt.Formatter) (int, bool) {
	return Fprint(ctx, os.StandardErrorStream, f)
}

func Eprintln(ctx *context.Context, f *fmt.Formatter) (int, bool) {
	return Fprintln(ctx, os.StandardErrorStream, f)
}

func Print(ctx *context.Context, f *fmt.Formatter) (int, bool) {
	return Fprint(ctx, os.StandardOutputStream, f)
}

func Println(ctx *context.Context, f *fmt.Formatter) (int, bool) {
	return Fprintln(ctx, os.StandardOutputStream, f)
}

package fmt_

import (
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/fmt"
	"github.com/anton2920/gofa/os"
)

func Fscan(ctx *context.Context, f os.Handle, s *fmt.Scanner) *fmt.Scanner {
	s.Reset()
	n, ok := os.ReadFromFile(ctx, f, s.Buffer[s.Pos:])
	if !ok {
		n = 0
	}
	s.MaxPos = n
	return s
}

func Scan(ctx *context.Context, s *fmt.Scanner) *fmt.Scanner {
	return Fscan(ctx, os.StandardInputStream, s)
}

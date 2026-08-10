package ints

import "github.com/anton2920/gofa/debug/debug_"

func AlignUpPow2(n int, align int) int {
	debug_.Assert((align&(align-1)) == 0, "alignment must be a power of two")
	return int((n + (align - 1)) & ^(align - 1))
}

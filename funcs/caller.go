package funcs

import (
	"unsafe"

	"github.com/anton2920/gofa/pointers"
)

/* GetCallerPC returns a value of %IP register that is going to be used by RET instruction. arg0 is the address of the first agrument function of interest accepts. */
//go:nosplit
func GetCallerPC(arg0 unsafe.Pointer) uintptr {
	return *(*uintptr)(pointers.Sub(arg0, unsafe.Sizeof(arg0)))
}

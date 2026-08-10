//go:build gofadebug
// +build gofadebug

package debug_

//go:nosplit
func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

//go:nosplit
func AssertZero(n int, msg string) {
	if n != 0 {
		panic(msg)
	}
}

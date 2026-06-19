//go:build !gofadebug
// +build !gofadebug

package debug

//go:nosplit
func Assert(_ bool, _ string) {}

//go:nosplit
func AssertZero(_ int, _ string) {}

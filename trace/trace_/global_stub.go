//go:build !gofatrace
// +build !gofatrace

package trace_

//go:nosplit
func BeginProfile() {}

//go:nosplit
func Begin(_ string) int {
	return 0
}

//go:nosplit
func End(_ int) {}

//go:nosplit
func EndAndPrintProfile() {}

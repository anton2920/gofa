//go:build gofadebug
// +build gofadebug

package debug_

import "runtime"

func init() {
	Println("[debug_]: Go version ", runtime.Version())
}

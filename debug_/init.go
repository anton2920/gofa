//go:build gofadebug
// +build gofadebug

package debug_

import "runtime"

func init() {
	Printf("[debug_]: Go version %s", runtime.Version())
}

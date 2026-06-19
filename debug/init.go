package debug

import "runtime"

func init() {
	Printf("[debug]: %s\n", runtime.Version())
}

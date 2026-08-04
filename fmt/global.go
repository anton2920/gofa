package fmt

import "sync/atomic"

var fs [16]Formatter
var current int32

func init() {
	for i := 0; i < len(fs); i++ {
		fs[i].Buffer = make([]byte, 1024)
	}
}

/* NOTE(anton2920): there's a risk of a race, if more than 'len(fs)' goroutines try to use it simultaneously. */
func Reset() *Formatter {
	n := atomic.AddInt32(&current, 1)
	return fs[n&int32(len(fs)-1)].Reset()
}

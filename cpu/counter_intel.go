//go:build 386 || amd64
// +build 386 amd64

package cpu

import "github.com/anton2920/gofa/cpu/intel"

func ReadPerformanceCounter() Cycles {
	return Cycles(intel.RDTSC())
}

//go:build i386 || amd64
// +build i386 amd64

package cpu

import "github.com/anton2920/gofa/cpu/intel"

func GetPerformanceCounter() Cycles {
	counter := intel.RDTSC()
	return Cycles(counter)
}

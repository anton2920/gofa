//go:build 386 || amd64
// +build 386 amd64

package cpu

import "github.com/anton2920/gofa/cpu/intel"

func ReadMemoryBarrier() {
	intel.LFENCE()
}

func WriteMemoryBarrier() {
	intel.SFENCE()
}

func TotalMemoryBarrier() {
	intel.MFENCE()
}

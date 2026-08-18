//go:build 386 || amd64
// +build 386 amd64

package cpu

import "github.com/anton2920/gofa/cpu/intel"

const RandomRetryLimit = 10

func ReadRandomUint16() (uint16, bool) {
	for i := 0; i < RandomRetryLimit; i++ {
		if n, ok := intel.RDRANDW(); ok {
			return n, true
		}
	}
	return 0, false
}

func ReadRandomUint32() (uint32, bool) {
	for i := 0; i < RandomRetryLimit; i++ {
		if n, ok := intel.RDRANDL(); ok {
			return n, true
		}
	}
	return 0, false
}

func ReadRandomUint64() (uint64, bool) {
	for i := 0; i < RandomRetryLimit; i++ {
		if n, ok := intel.RDRANDQ(); ok {
			return n, true
		}
	}
	return 0, false
}

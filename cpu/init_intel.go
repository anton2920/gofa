//go:build 386 || amd64
// +build 386 amd64

package cpu

import "github.com/anton2920/gofa/cpu/intel"

var CPUHz Cycles

func init() {
	CPUHz = Cycles(intel.CPUHz)
}

//go:build i386 || amd64
// +build i386 amd64

package cpu

import (
	"github.com/anton2920/gofa/cpu/intel"
	"github.com/anton2920/gofa/debug"
)

var CPUHz Cycles

func init() {
	CPUHz = Cycles(intel.CPUHz)
	if CPUHz != 0 {
		debug.Printf("[cpu]: CPU Frequency %dHz", CPUHz)
	}
}

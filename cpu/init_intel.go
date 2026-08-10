//go:build 386 || amd64
// +build 386 amd64

package cpu

import (
	"github.com/anton2920/gofa/cpu/intel"
	"github.com/anton2920/gofa/debug/debug_"
)

var CPUHz Cycles

func init() {
	CPUHz = Cycles(intel.CPUHz)
	if CPUHz > 0 {
		debug_.Printf("[cpu]: CPU frequency: %dHz", CPUHz)
	}
}

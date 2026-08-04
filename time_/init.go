package time_

import (
	"github.com/anton2920/gofa/cpu"
	"github.com/anton2920/gofa/debug_"
	"github.com/anton2920/gofa/time"
)

func init() {
	if cpu.CPUHz == 0 {
		const (
			millisecondsToWait = 10
			osHz               = int64(time.Second)
			osWaitTime         = osHz * millisecondsToWait / (time.Second / time.Millisecond)
		)
		var osElapsed int64

		osStart := Now()
		cpuStart := cpu.ReadPerformanceCounter()

		for osElapsed < osWaitTime {
			osElapsed = Now() - osStart
		}

		cpuEnd := cpu.ReadPerformanceCounter()
		cpuElapsed := int64(cpuEnd - cpuStart)
		cpu.CPUHz = cpu.Cycles(cpuElapsed * osHz / osElapsed)

		debug_.Printf("[time]: CPU Frequency %dHz", cpu.CPUHz)
	}
}

package time

import (
	"github.com/anton2920/gofa/cpu"
	"github.com/anton2920/gofa/debug"
	"github.com/anton2920/gofa/syscall"
)

func init() {
	if cpu.CPUHz == 0 {
		const osHz = int64(Second)
		const millisecondsToWait = 10

		cpuStart := cpu.ReadPerformanceCounter()
		osStart := Now()

		osEnd := Now()
		osElapsed := int64(0)
		osWaitTime := osHz * millisecondsToWait / (Second / Millisecond)

		for osElapsed < osWaitTime {
			osEnd = Now()
			osElapsed = osEnd - osStart
		}

		cpuEnd := cpu.ReadPerformanceCounter()
		cpuElapsed := int64(cpuEnd - cpuStart)
		cpu.CPUHz = cpu.Cycles(cpuElapsed * osHz / osElapsed)

		debug.Printf("[time]: CPU Frequency %dHz", cpu.CPUHz)
	}
}

/* Now returns current wallclock time, nanosecond resolution. */
func Now() int64 {
	var tp syscall.Timespec
	syscall.ClockGettime(syscall.CLOCK_REALTIME, &tp)
	return tp.Sec*int64(Second) + tp.Nsec
}

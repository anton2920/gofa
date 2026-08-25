package time_

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/cpu"
	"github.com/anton2920/gofa/debug/debug_"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
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

		osStart := NowInNanoseconds()
		cpuStart := cpu.ReadPerformanceCounter()

		for osElapsed < osWaitTime {
			osElapsed = NowInNanoseconds() - osStart
		}

		cpuEnd := cpu.ReadPerformanceCounter()
		cpuElapsed := int64(cpuEnd - cpuStart)
		cpu.CPUHz = cpu.Cycles(cpuElapsed * osHz / osElapsed)

		buf := make([]byte, ints.Bufsize)
		n := slices.PutUint64(buf, uint64(cpu.CPUHz))
		debug_.Println("[time_]: CPU frequency: ", bytes.AsString(buf[:n]), "Hz")
	}
}

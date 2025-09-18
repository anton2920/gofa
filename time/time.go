package time

import (
	"github.com/anton2920/gofa/debug"
	"github.com/anton2920/gofa/intel"
	"github.com/anton2920/gofa/syscall"
)

const (
	MsecPerSec = 1000
	UsecPerSec = MsecPerSec * 1000
	NsecPerSec = UsecPerSec * 1000
)

const (
	OneHour  = 60 * 60
	OneDay   = OneHour * 24
	OneWeek  = OneDay * 7
	OneMonth = OneDay * 30
	OneYear  = OneDay * 365
)

func init() {
	if intel.CPUHz == 0 {
		const osHz = int64(1 * NsecPerSec)
		const millisecondsToWait = 10

		cpuStart := intel.RDTSC()
		osStart := UnixNsec()

		osEnd := UnixNsec()
		osElapsed := int64(0)
		osWaitTime := osHz * millisecondsToWait / MsecPerSec

		for osElapsed < osWaitTime {
			osEnd = UnixNsec()
			osElapsed = osEnd - osStart
		}

		cpuEnd := intel.RDTSC()
		cpuElapsed := int64(cpuEnd - cpuStart)
		intel.CPUHz = intel.Cycles(cpuElapsed * osHz / osElapsed)
		debug.Printf("[gofa/time]: CPU Frequency %dHz", intel.CPUHz)
	}
}

func Unix() int64 {
	var tp syscall.Timespec
	syscall.ClockGettime(syscall.CLOCK_REALTIME, &tp)
	return tp.Sec
}

func UnixNsec() int64 {
	var tp syscall.Timespec
	syscall.ClockGettime(syscall.CLOCK_REALTIME, &tp)
	return tp.Sec*1000000000 + tp.Nsec
}

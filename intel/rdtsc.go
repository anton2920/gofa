package intel

import "github.com/anton2920/gofa/time"

type Cycles uint64

var CpuHz Cycles

//go:nosplit
func RDTSC() Cycles

func init() {
	const osHz = int64(1 * time.NsecPerSec)
	const millisecondsToWait = 10

	cpuStart := RDTSC()
	osStart := time.UnixNsec()

	osEnd := time.UnixNsec()
	osElapsed := int64(0)
	osWaitTime := osHz * millisecondsToWait / time.MsecPerSec
	for osElapsed < osWaitTime {
		osEnd = time.UnixNsec()
		osElapsed = osEnd - osStart
	}

	cpuEnd := RDTSC()
	cpuElapsed := int64(cpuEnd - cpuStart)
	CpuHz = Cycles(cpuElapsed * osHz / osElapsed)
}

func (c Cycles) ToNsec() int64 {
	return int64(c * time.NsecPerSec / CpuHz)
}

func (c Cycles) ToUsec() int64 {
	return int64(c * time.UsecPerSec / CpuHz)
}

func (c Cycles) ToMsec() int64 {
	return int64(c * time.MsecPerSec / CpuHz)
}

func (c Cycles) ToSec() int64 {
	return int64(c / CpuHz)
}

package time

type Cycles uint64

var CpuHz Cycles

func init() {
	const osHz = int64(1 * NsecPerSec)
	const millisecondsToWait = 10

	cpuStart := RDTSC()
	osStart := UnixNsec()

	osEnd := UnixNsec()
	osElapsed := int64(0)
	osWaitTime := osHz * millisecondsToWait / MsecPerSec
	for osElapsed < osWaitTime {
		osEnd = UnixNsec()
		osElapsed = osEnd - osStart
	}

	cpuEnd := RDTSC()
	cpuElapsed := int64(cpuEnd - cpuStart)
	CpuHz = Cycles(cpuElapsed * osHz / osElapsed)
}

func RDTSC() Cycles

func (c Cycles) ToNsec() int64 {
	return int64(c * NsecPerSec / CpuHz)
}

func (c Cycles) ToUsec() int64 {
	return int64(c * UsecPerSec / CpuHz)
}

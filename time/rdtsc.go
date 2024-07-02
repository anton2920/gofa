package time

type Cycles uint64

var CpuFreq Cycles

func init() {
	const millisecondsToWait = 1
	const osFreq = int64(1 * NsecPerSec)

	cpuStart := RDTSC()
	osStart := UnixNsec()

	osEnd := UnixNsec()
	osElapsed := int64(0)
	osWaitTime := osFreq * millisecondsToWait / 1000
	for osElapsed < osWaitTime {
		osEnd = UnixNsec()
		osElapsed = osEnd - osStart
	}

	cpuEnd := RDTSC()
	cpuElapsed := int64(cpuEnd - cpuStart)
	CpuFreq = Cycles(osFreq * cpuElapsed / osElapsed)
}

func RDTSC() Cycles

func (c Cycles) ToNsec() int64 {
	return int64(c * NsecPerSec / CpuFreq)
}

func (c Cycles) ToUsec() int64 {
	return int64(c * UsecPerSec / CpuFreq)
}

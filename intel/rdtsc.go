package intel

type Cycles uint64

//go:nosplit
func RDTSC() Cycles

func (c Cycles) ToNsec() int64 {
	return int64(c * 1000000000 / CPUHz)
}

func (c Cycles) ToUsec() int64 {
	return int64(c * 1000000 / CPUHz)
}

func (c Cycles) ToMsec() int64 {
	return int64(c * 1000 / CPUHz)
}

func (c Cycles) ToSec() int64 {
	return int64(c / CPUHz)
}

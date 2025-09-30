package cpu

type Cycles uint64

func (c Cycles) ToNanoseconds() int64 {
	return int64(c * 1000000000 / CPUHz)
}

func (c Cycles) ToMicroseconds() int64 {
	return int64(c * 1000000 / CPUHz)
}

func (c Cycles) ToMilliseconds() int64 {
	return int64(c * 1000 / CPUHz)
}

func (c Cycles) ToSeconds() int64 {
	return int64(c / CPUHz)
}

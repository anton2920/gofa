package cpu

type Cycles uint64

func (c Cycles) ToNanoseconds() float64 {
	return float64(float64(c) * 1000000000 / float64(CPUHz))
}

func (c Cycles) ToNanosecondsTruncated() int64 {
	return int64(c * 1000000000 / CPUHz)
}

func (c Cycles) ToMicroseconds() float64 {
	return float64(float64(c) * 1000000 / float64(CPUHz))
}

func (c Cycles) ToMicrosecondsTruncated() int64 {
	return int64(c * 1000000 / CPUHz)
}

func (c Cycles) ToMilliseconds() float64 {
	return float64(float64(c) * 1000 / float64(CPUHz))
}

func (c Cycles) ToMillisecondsTruncated() int64 {
	return int64(c * 1000 / CPUHz)
}

func (c Cycles) ToSeconds() float64 {
	return float64(float64(c) / float64(CPUHz))
}

func (c Cycles) ToSecondsTruncated() int64 {
	return int64(c / CPUHz)
}

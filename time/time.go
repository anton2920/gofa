package time

import "github.com/anton2920/gofa/syscall"

const (
	MsecPerSec = 1000
	UsecPerSec = MsecPerSec * 1000
	NsecPerSec = UsecPerSec * 1000
)

func Unix() int {
	var tp syscall.Timespec
	syscall.ClockGettime(syscall.CLOCK_REALTIME, &tp)
	return int(tp.Sec)
}

func UnixNsec() int64 {
	var tp syscall.Timespec
	syscall.ClockGettime(syscall.CLOCK_REALTIME, &tp)
	return tp.Sec*1_000_000_000 + tp.Nsec
}

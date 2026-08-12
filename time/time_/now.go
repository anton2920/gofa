package time_

import (
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/time"
)

func NowInSeconds() int {
	t, _ := os.GetCurrentTime()
	return t.Seconds
}

func NowInMilliseconds() int64 {
	return NowInNanoseconds() / 1000000
}

func NowInMicroseconds() int64 {
	return NowInNanoseconds() / 1000
}

func NowInNanoseconds() int64 {
	t, _ := os.GetCurrentTime()
	return int64(t.Seconds)*time.Second + int64(t.Nanoseconds)
}

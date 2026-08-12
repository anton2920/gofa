package time_

import (
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/time"
)

func NowInSeconds() int {
	var t os.SecondsWithNanoseconds
	_ = os.GetCurrentTime(&context.Context{}, &t)
	return t.Seconds
}

func NowInMilliseconds() int64 {
	return NowInNanoseconds() / 1000000
}

func NowInMicroseconds() int64 {
	return NowInNanoseconds() / 1000
}

func NowInNanoseconds() int64 {
	var t os.SecondsWithNanoseconds
	_ = os.GetCurrentTime(&context.Context{}, &t)
	return int64(t.Seconds)*time.Second + int64(t.Nanoseconds)
}

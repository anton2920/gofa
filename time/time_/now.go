package time_

import (
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/time"
)

/* Now returns current wallclock time, nanosecond resolution. */
func Now() int64 {
	t, _ := os.GetCurrentTime()
	return int64(t.Seconds)*time.Second + int64(t.Nanoseconds)
}

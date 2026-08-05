package time_

import (
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/time"
)

func Sleep(t int64) error {
	var sn os.SecondsWithNanoseconds
	sn.Seconds = int(t / time.Second)
	sn.Nanoseconds = int(t % time.Second)
	return os.BlockForSpecifiedAmountOfTime(sn)
}

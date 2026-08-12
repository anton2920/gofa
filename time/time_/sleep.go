package time_

import (
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/time"
)

func Sleep(t int64) {
	var sn os.SecondsWithNanoseconds
	sn.Seconds = int(t / time.Second)
	sn.Nanoseconds = int(t % time.Second)
	_ = os.BlockForSpecifiedAmountOfTime(&context.Context{}, sn)
}

package log_

import (
	"sync/atomic"
	"unsafe"

	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/time_"
)

var fs [16]log.Formatter
var current int32

var LogLevel log.Level

func init() {
	for i := 0; i < len(fs); i++ {
		fs[i].Fmt.Buffer = make([]byte, 1024)
	}
}

/* NOTE(anton2920): there's a risk of a race, if more than 'len(fs)' goroutines try to use it simultaneously. */
func Reset() *log.Formatter {
	n := atomic.AddInt32(&current, 1)
	return fs[n&int32(len(fs)-1)].Reset(time_.Now())
}

/* NOTE(anton2920): there's a risk of a race, if more than 'len(fs)' goroutines try to use it simultaneously. */
func Log(level log.Level) *log.Formatter {
	logLevel := log.Level(atomic.LoadInt32((*int32)(unsafe.Pointer(&LogLevel))))
	if level >= logLevel {
		n := atomic.AddInt32(&current, 1)
		f := &fs[n&int32(len(fs)-1)]
		f.MinimumLevel = logLevel
		return f.Log(level, time_.Now())
	}
	return nil
}

func Debug() *log.Formatter {
	return Log(log.LevelDebug)
}

func Info() *log.Formatter {
	return Log(log.LevelInfo)
}

func Warn() *log.Formatter {
	return Log(log.LevelWarn)
}

func Error() *log.Formatter {
	return Log(log.LevelError)
}

func Fatal() *log.Formatter {
	return Log(log.LevelFatal)
}

func Panic() *log.Formatter {
	return Log(log.LevelPanic)
}

func SetLevel(level log.Level) log.Level {
	return log.Level(atomic.SwapInt32((*int32)(unsafe.Pointer(&LogLevel)), int32(level)))
}

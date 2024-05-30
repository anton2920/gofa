package log

import (
	"fmt"
	"os"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

type Level int32

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

var Level2String = [...]interface{}{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
	LevelPanic: "PANIC",
}

var LogLevel Level = LevelInfo

func Logf(level Level, format string, args ...interface{}) {
	if (level < LevelFatal) && (level < Level(atomic.LoadInt32((*int32)(unsafe.Pointer(&LogLevel))))) {
		return
	}

	buffer := make([]byte, 0, 512)
	buffer = time.Now().AppendFormat(buffer, "2006/01/02 15:04:05")
	buffer = fmt.Appendf(buffer, " %5s ", Level2String[level])
	buffer = fmt.Appendf(buffer, format, args...)
	if format[len(format)-1] != '\n' {
		buffer = append(buffer, '\n')
	}

	switch level {
	default:
		/* TODO(anton2920): are race-conditions possible? */
		syscall.Write(2, buffer)
	case LevelFatal:
		syscall.Write(2, buffer)
		os.Exit(1)
	case LevelPanic:
		panic(string(buffer))
	}
}

func Debugf(format string, args ...interface{}) {
	Logf(LevelDebug, format, args...)
}

func Infof(format string, args ...interface{}) {
	Logf(LevelInfo, format, args...)
}

func Warnf(format string, args ...interface{}) {
	Logf(LevelWarn, format, args...)
}

func Errorf(format string, args ...interface{}) {
	Logf(LevelError, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Logf(LevelFatal, format, args...)
}

func Panicf(format string, args ...interface{}) {
	Logf(LevelPanic, format, args...)
}

func SetLevel(new Level) Level {
	return Level(atomic.SwapInt32((*int32)(unsafe.Pointer(&LogLevel)), int32(new)))
}

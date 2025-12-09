package log

import (
	"bytes"
	"fmt"
	"os"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/anton2920/gofa/pointers"
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

	var buffer bytes.Buffer

	buffer.WriteString(time.Now().Format("2006/01/02 15:04:05"))
	fmt.Fprintf((*bytes.Buffer)(pointers.Noescape(unsafe.Pointer(&buffer))), " %5s ", Level2String[level])
	fmt.Fprintf((*bytes.Buffer)(pointers.Noescape(unsafe.Pointer(&buffer))), format, args...)
	if format[len(format)-1] != '\n' {
		buffer.WriteRune('\n')
	}

	switch level {
	default:
		/* TODO(anton2920): are race-conditions possible? */
		syscall.Write(2, buffer.Bytes())
	case LevelFatal:
		syscall.Write(2, buffer.Bytes())
		os.Exit(1)
	case LevelPanic:
		panic(buffer.String())
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

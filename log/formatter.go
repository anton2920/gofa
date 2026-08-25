package log

import "github.com/anton2920/gofa/fmt"

type Formatter struct {
	Fmt fmt.Formatter

	MinimumLevel Level
	CurrentLevel Level
}

func (f *Formatter) InitWithByteSlice(buf []byte) {
	f.Fmt.InitWithByteSlice(buf)
}

func (f *Formatter) InitWithMinimumLevelAndByteSlice(level Level, buf []byte) {
	f.MinimumLevel = level
	f.InitWithByteSlice(buf)
}

func (f *Formatter) Reset(t int64) *Formatter {
	f.Fmt.Reset()
	f.CurrentLevel = LevelLast
	return f.DateTime(t).S(" ")
}

func (f *Formatter) Log(level Level, t int64) *Formatter {
	f.Fmt.Reset()
	f.CurrentLevel = level

	if f.CurrentLevel >= f.MinimumLevel {
		f.DateTime(t).S(" ").W(5).S(Level2String[f.CurrentLevel]).S(" ")
		return f
	}

	return nil
}

func (f *Formatter) Debug(t int64) *Formatter {
	return f.Log(LevelDebug, t)
}

func (f *Formatter) Info(t int64) *Formatter {
	return f.Log(LevelInfo, t)
}

func (f *Formatter) Warn(t int64) *Formatter {
	return f.Log(LevelWarn, t)
}

func (f *Formatter) Error(t int64) *Formatter {
	return f.Log(LevelError, t)
}

func (f *Formatter) Fatal(t int64) *Formatter {
	return f.Log(LevelFatal, t)
}

func (f *Formatter) Panic(t int64) *Formatter {
	return f.Log(LevelPanic, t)
}

func (f *Formatter) String() string {
	if f != nil {
		return f.Fmt.String()
	}
	return ""
}

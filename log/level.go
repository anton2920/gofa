package log

type Level int32

const (
	LevelDebug = Level(iota - 1)
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
	LevelLast
)

func Level2String(lvl Level) string {
	switch lvl {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelPanic:
		return "PANIC"
	}
	return ""
}

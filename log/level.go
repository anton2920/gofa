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

var Level2String = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
	LevelPanic: "PANIC",
}

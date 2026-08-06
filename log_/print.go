package log_

import (
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/os"
)

var panicMsg interface{} = "caused by log_.Panic"

func Println(f *log.Formatter) {
	if f != nil {
		h := os.StandardOutputStream
		if f.CurrentLevel > log.LevelWarn {
			h = os.StandardErrorStream
		}
		os.WriteToFile(h, f.Fmt.Ln().Bytes())
		switch f.CurrentLevel {
		case log.LevelFatal:
			os.Exit(1)
		case log.LevelPanic:
			panic(panicMsg)
		}
	}
}

package logs

import (
	"log"
	"os"
)

type Level int

const (
	LevelClose Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelAll
)

var levelStrMap = map[Level]string{
	LevelClose: "",
	LevelError: "[Error]",
	LevelWarn:  "[Warn]",
	LevelInfo:  "[Info]",
	LevelDebug: "[Debug]",
	LevelAll:   "",
}

// simple log
type LogFactory struct {
	Level Level
}

var Log = &LogFactory{Level: LevelAll}

func logg(level Level, v ...any) {
	if Log.Level < level {
		return
	}
	if level == LevelError {
		if os.Getenv("GO_DECORATOR_LOG_LEVEL") == "debug" {
			log.Panicln(append([]any{levelStrMap[level]}, v...)...)
			return
		}
		log.Println(append([]any{levelStrMap[level]}, v...)...)
		os.Exit(2)
		return
	}
	log.Println(append([]any{levelStrMap[level]}, v...)...)
}

func Debug(v ...any) {
	logg(LevelDebug, v...)
}

func Info(v ...any) {
	logg(LevelInfo, v...)
}

func Warn(v ...any) {
	logg(LevelWarn, v...)
}

func Error(v ...any) {
	logg(LevelError, v...)
}

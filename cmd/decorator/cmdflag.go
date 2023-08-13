package main

import (
	"flag"
	"github.com/dengsgo/go-decorator/cmd/logs"
)

type CmdFlag struct {
	Level   string // -decor.level
	TempDir string // -decor.tempDir

	// go build args
	toolPath  string
	chainName string
	chainArgs []string
}

func initUseFlag() {
	flag.StringVar(&cmdFlag.Level,
		"decor.level",
		"warn",
		"output log level. all/debug/info/warn/error/close")
	flag.StringVar(&cmdFlag.TempDir,
		"decor.tempDir",
		"",
		"tool workspace dir. default same as go build workspace")
	flag.Parse()
	switch cmdFlag.Level {
	case "all":
		logs.Log.Level = logs.LevelAll
	case "debug":
		logs.Log.Level = logs.LevelDebug
	case "info":
		logs.Log.Level = logs.LevelInfo
	case "warn":
		logs.Log.Level = logs.LevelWarn
	case "error", "":
		logs.Log.Level = logs.LevelError
	case "close":
		logs.Log.Level = logs.LevelClose
	}
	if cmdFlag.TempDir != "" {
		tempDir = cmdFlag.TempDir // TODO check
	}
}

var cmdFlag = &CmdFlag{}

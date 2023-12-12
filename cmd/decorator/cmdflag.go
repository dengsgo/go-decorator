package main

import (
	"flag"
	"fmt"
	"github.com/dengsgo/go-decorator/cmd/logs"
	"log"
	"os"
	"strings"
)

const version = `v0.21.0 beta`
const opensourceUrl = `https://github.com/dengsgo/go-decorator`

type CmdFlag struct {
	Level     string // -d.log
	TempDir   string // -d.tempDir
	ClearWork bool   // -d.clearWork
	Version   string // -version

	// go build args
	toolPath  string
	chainName string
	chainArgs []string
}

func initUseFlag() {
	flag.StringVar(&cmdFlag.Level,
		"d.log",
		"warn",
		"output log level. all/debug/info/warn/error/close")
	flag.StringVar(&cmdFlag.TempDir,
		"d.tempDir",
		"",
		"tool workspace dir. default same as go build workspace")
	flag.BoolVar(&cmdFlag.ClearWork,
		"d.clearWork",
		true,
		"empty workspace when compilation is complete")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(),
			"decorator [-d.log] [-d.tempDir] chainToolPath chainArgs\n")
		flag.PrintDefaults()
	}
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
	log.SetPrefix("decorator: ")
	if logs.Log.Level < logs.LevelDebug {
		log.SetFlags(0)
	}
	if cmdFlag.TempDir != "" {
		tempDir = cmdFlag.TempDir // TODO check
	}
	cmdFlag.toolPath = os.Args[0]
	goToolDir := os.Getenv("GOTOOLDIR")
	if goToolDir == "" {
		logs.Info("env key `GOTOOLDIR` not found")
	}
	if len(os.Args) < 2 {
		fmt.Fprintf(flag.CommandLine.Output(),
			"decorator %s , %s\n", version, opensourceUrl)
		os.Exit(0)
	}
	for i, arg := range os.Args[1:] {
		if goToolDir != "" && strings.HasPrefix(arg, goToolDir) {
			cmdFlag.chainName = arg
			if len(os.Args[1:]) > i+1 {
				cmdFlag.chainArgs = os.Args[i+2:]
			}
			break
		}
	}
}

var cmdFlag = &CmdFlag{}

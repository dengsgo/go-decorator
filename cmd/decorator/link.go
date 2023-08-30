package main

import (
	"github.com/dengsgo/go-decorator/cmd/logs"
	"os"
	"path/filepath"
	"strings"
)

func link(args []string) {
	var cfg string
	buildmode := false
	for _, arg := range args {
		if arg == "-buildmode=exe" ||
			// windows
			arg == "-buildmode=pie" {
			buildmode = true
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		if strings.Contains(arg, filepath.Join("b001", "importcfg.link")) {
			cfg = arg
		}
	}
	logs.Debug("cfg", cfg)
	if !buildmode || cfg == "" {
		return
	}
	workspaceCleaner = func() {
		_ = os.RemoveAll(tempDir)
	}
}

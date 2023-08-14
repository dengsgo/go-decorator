package main

import (
	"github.com/dengsgo/go-decorator/cmd/logs"
	"os"
	"path/filepath"
	"strings"
)

func link(args []string) {
	//-importcfg $WORK/b001/importcfg.link
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
	if !buildmode {
		return
	}
	if cfg == "" {
		logs.Error("link -importcfg cfg not found")
	}
	bs, err := os.ReadFile(cfg)
	if err != nil {
		logs.Error("link ReadFile cfg fail", err)
	}
	sf := filepath.Join(tempDir, "sharedPackagefile.txt")
	workspaceCleaner = func() {
		_ = os.Remove(sf)
		_ = os.RemoveAll(tempDir)
	}
	newBs, err := os.ReadFile(sf)
	if err != nil {
		logs.Debug("link ReadFile sharedPackagefile.txt err", err)
		logs.Debug("nothing todo link, process ignore")
		return
	}
	// TODO 去重
	bs = append(bs, newBs...)
	os.WriteFile(cfg, bs, 0777)
}

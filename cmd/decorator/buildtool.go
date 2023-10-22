package main

import (
	"github.com/dengsgo/go-decorator/cmd/logs"
	"github.com/dengsgo/go-decorator/decor"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	biSymbol             = "\n\t"
	decoratorScanFlag    = "//go:decor "
	decoratorPackagePath = "github.com/dengsgo/go-decorator/decor"
)

var (
	tempDir       = path.Join(os.TempDir(), "gobuild_decorator_works")
	tempGenDir    = tempDir
	projectDir, _ = os.Getwd()
	exitDo        = func() {}
)

func inits() {
	initUseFlag()
	initTempDir()
}

func initTempDir() {
	if err := os.MkdirAll(tempDir, 0777); err != nil {
		logs.Error("Init() fail, os.MkdirAll tempDir", err)
	}
}

func main() {
	inits()
	logs.Debug("os.Args", os.Args)
	logs.Debug("os.Env", os.Environ())
	if cmdFlag.chainName == "" {
		logs.Error("currently not in a compilation chain environment and cannot be used")
	}
	logs.Debug("cmdFlag", cmdFlag)
	chainName := cmdFlag.chainName
	chainArgs := cmdFlag.chainArgs
	toolName := filepath.Base(chainName)

	var err error
	switch strings.TrimSuffix(toolName, ".exe") {
	case "compile":
		err = compile(chainArgs)
	case "link":
		link(chainArgs)
		defer func() {
			logs.Debug("exitDo() begin")
			exitDo()
			logs.Debug("exitDo() end")
		}()
	}

	if err != nil {
		logs.Error(err)
	}
	// build
	cmd := exec.Command(chainName, chainArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if cmd.Run() != nil {
		//logs.Error("run toolchain err", chainName, err)
	}
}

//go:decor logging
func test(v ...string) string {
	return ""
}

func logging(ctx *decor.Context, s string, a int, b bool) {
	ctx.TargetDo()
}

// ###############################

//func myFuncDecor(a int, b string) (_decorGenOut1 int, _decorGenOut2 int) {
//	decor := &DecorContext{
//		WarpFuncIn:  []any{a, b},
//		WarpFuncOut: []any{_decorGenOut1, _decorGenOut2},
//	}
//	decor.Func = func() {
//		decor.WarpFuncOut[0], decor.WarpFuncOut[1] = func(a int, b string) (int, int) {
//			log.Println("Func: myFunc", b)
//			return a, a + 1
//		}(decor.WarpFuncIn[0].(int), decor.WarpFuncIn[1].(string))
//	}
//	logging(decor)
//	return decor.WarpFuncOut[0].(int), decor.WarpFuncOut[1].(int)
//}

package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"log"
)

func main() {
	log.SetFlags(0)
	log.Println("MAIN call myFunc1()")
	myFunc1()
	log.Println("MAIN call myFunc2UseMultipleDecor()")
	myFunc2UseMultipleDecor()
	log.Println("MAIN call myFunc2UseMultipleDecor()")
	myFunc3HaveSigParams(100, "this is test case", "")
}

//go:decor decorHandlerFunc
func myFunc1() {
	log.Println("call myFunc1")
}

//go:decor decorHandlerFunc
//go:decor yetDecorHandlerFunc
func myFunc2UseMultipleDecor() {
	log.Println("call myFunc1")
}

// This is a function definition consisting of input and output parameters.
// In the processing method of the decorator, they can be easily obtained
// through the context of ctx.TargetIn and ctx.TargetOut.
//
// Furthermore, you can modify the input parameters before ctx.TargetDo()
// and the output parameters after ctx.TargetDo()
//
//go:decor decorHandlerFunc
func myFunc3HaveSigParams(i int, s string, v ...any) (bool, *int64) {
	return true, new(int64)
}

// decor handler functions

func decorHandlerFunc(ctx *decor.Context) {
	log.Println("call decorHandlerFunc in", ctx.Kind, ctx.TargetIn)
	ctx.TargetDo()
	log.Println("call decorHandlerFunc out", ctx.TargetOut)
}

func yetDecorHandlerFunc(ctx *decor.Context) {
	log.Println("call yetDecorHandlerFunc in", ctx.Kind, ctx.TargetIn)
	ctx.TargetDo()
	log.Println("call yetDecorHandlerFunc out", ctx.TargetOut)
}

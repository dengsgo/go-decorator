package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"log"
)

func dontExec(ctx *decor.Context) {
	log.Println("decorator 'dontExec' be called, but nothing todo")
}

func execIt(ctx *decor.Context) {
	log.Println("decorator 'exec' be called, begin")
	ctx.TargetDo()
	log.Println("decorator 'exec' be called, end")
}

//go:decor dontExec
func Add(a, b int64) int64 {
	return a + b
}

//go:decor execIt
func empty() string {
	return "empty() result is here"
}

func main() {
	addResult := Add(1, 1)
	log.Println("use decor dontExec: addResult(1, 1) =", addResult)
	sep()
	log.Println("use decor execIt: empty() =", empty())
}

func sep() {
	log.Println("================")
}

func init() {
	log.SetFlags(0)
}

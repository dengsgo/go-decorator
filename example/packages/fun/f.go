package fun

import (
	"github.com/dengsgo/go-decorator/decor"
	_ "github.com/dengsgo/go-decorator/example/packages/fun1"
	"log"
)

func DecorHandlerFunc(ctx *decor.Context) {
	log.Println("call fun.decorHandlerFunc in", ctx.Kind, ctx.TargetIn)
	ctx.TargetDo()
	Ts()
	log.Println("call fun.decorHandlerFunc out", ctx.TargetOut)
}

//go:decor fun1.DecorHandlerFunc
func Ts() {
	log.Println("call Ts yeah")
}

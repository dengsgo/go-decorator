package fun1

import (
	"github.com/dengsgo/go-decorator/decor"
	"log"
)

func DecorHandlerFunc(ctx *decor.Context) {
	log.Println("call fun1.decorHandlerFunc in", ctx.Kind, ctx.TargetIn)
	ctx.TargetDo()
	log.Println("call fun1.decorHandlerFunc out", ctx.TargetOut)
}

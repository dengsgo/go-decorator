package fun

import (
	"github.com/dengsgo/go-decorator/decor"
	"log"
)

func DecorHandlerFunc(ctx *decor.Context) {
	log.Println("call decorHandlerFunc in", ctx.Kind, ctx.TargetIn)
	ctx.TargetDo()
	log.Println("call decorHandlerFunc out", ctx.TargetOut)
}

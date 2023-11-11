package externala

import (
	"github.com/dengsgo/go-decorator/decor"
	_ "github.com/dengsgo/go-decorator/example/usages/externala/deepexternal"
	"github.com/dengsgo/go-decorator/example/usages/g"
)

func OnlyPrintSelf(ctx *decor.Context) {
	g.PrintfLn("the target use [externala.OnlyPrintSelf] decorator")
	ctx.TargetDo()
	s := UseDeepExternalDecor()
	g.PrintfLn(s)
}

//go:decor deepexternal.FixedStringWhenReturnString
func UseDeepExternalDecor() string {
	return "UseDeepExternalDecor return string, It will be modified by the decorator"
}

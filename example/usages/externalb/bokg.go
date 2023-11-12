package externalb

import "github.com/dengsgo/go-decorator/decor"

func MathIntegerPlus(a, b int) int {
	return a + b
}

func DoubleIntegerValue(ctx *decor.Context) {
	ctx.TargetDo()
	if len(ctx.TargetOut) == 0 {
		return
	}
	switch firstValue := ctx.TargetOut[0].(type) {
	case int:
		ctx.TargetOut[0] = firstValue * 2
	}
}

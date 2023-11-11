package deepexternal

import "github.com/dengsgo/go-decorator/decor"

// 如果目标函数有返回值并且第一个是 string类型，将它修改成固定的字符串
func FixedStringWhenReturnString(ctx *decor.Context) {
	ctx.TargetDo()
	if len(ctx.TargetOut) > 0 && func() bool {
		_, ok := ctx.TargetOut[0].(string)
		return ok
	}() {
		ctx.TargetOut[0] = "Return String By [externala.deepexternal.FixedStringWhenReturnString]"
	}
}

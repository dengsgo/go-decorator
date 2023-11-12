package main

import (
	"fmt"
	"github.com/dengsgo/go-decorator/decor"
)

// 这个文件演示使用带有参数的装饰器用法，和 lint 的用法。
// lint 在编译时会验证目标函数的传参，如果不过约束无法通过编译。
//
// required: 参数必传
// nonzero: 参数不能是零值
// 通过文档 Guide.md 查看用法。

// This is a decorator function with parameters.
// It checks if the first element of ctx.TargetOut is a string, and if it is, it replaces that element
// with a formatted string that includes the values of the input parameters.
// Use `go:decor-lint` to add call constraints, such as which arguments are "required" and so on.
// If you don't meet the constraints, you will get an error at compile time.
//
//go:decor-lint required: {msg, repeat, count, f}
//go:decor-lint nonzero: {msg, count, f}
func hit(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	ctx.TargetDo()
	if len(ctx.TargetOut) >= 1 &&
		func() bool {
			_, ok := ctx.TargetOut[0].(string)
			return ok
		}() {
		ctx.TargetOut[0] = fmt.Sprintf("hit received: msg=%s, count=%d, repeat=%t, f=%f, opt=%s\n",
			msg, count, repeat, f, opt)
	}
}

// The function has a decorator called hit with some arguments.
// The decorator is applied to the function using a comment with the go:decor directive.
// The decorator is expected to modify the behavior of the function in some way.
// The function itself does not have any implementation and returns an empty string.
//
//go:decor hit#{msg: "message from decor", repeat: true, count: 10, f:1}
func useArgsDecor() (s string) {
	return
}

// =============================================
// ========== 下面演示更多 lint 的用法 ===========
// =============================================

//go:decor-lint required: {msg, repeat, count: {gte:5, lte:20}, f: {gt:0}}
func hitUseRequiredLint(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	hit(ctx, msg, count, repeat, f, opt)
}

//go:decor-lint nonzero: {msg, count, f}
func hitUseNonzeroLint(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	hit(ctx, msg, count, repeat, f, opt)
}

//go:decor-lint required: {msg, repeat, count: {gte:5, lte:20}, f: {gt:-10}}
//go:decor-lint nonzero: {msg, count, f}
func hitBothUseLint(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	hit(ctx, msg, count, repeat, f, opt)
}

//go:decor-lint required: {msg}
//go:decor-lint required: {msg: {"hello", "world"}}
//go:decor-lint required: {count: {gte:100, lte:200}}
//go:decor-lint nonzero: { f }
func hitUseMultilineLint(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	hit(ctx, msg, count, repeat, f, opt)
}

//go:decor hitUseRequiredLint#{msg: "你好", repeat: false, count: 10, f:1}
func useHitUseRequiredLint() (s string) {
	return
}

//go:decor hitUseNonzeroLint#{msg: "你好", count: 150, f:1}
func useHitUseNonzeroLint() (s string) {
	return
}

//go:decor hitBothUseLint#{msg: "message from decor, useHitBothUseLint", repeat: true, count: 10, f:1}
func useHitBothUseLint() (s string) {
	return
}

//go:decor hitUseMultilineLint#{msg: "hello", repeat: true, count: 150, f:1}
func useHitUseMultilineLintDecor() (s string) {
	return
}

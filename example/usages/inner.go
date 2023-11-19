package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"github.com/dengsgo/go-decorator/example/usages/g"
)

// 这个文件演示了使用当前包内的装饰器的用法。
// 如下代码所示，只需要 //go:decor dumpDecorContext 即可。
// 无需考虑装饰器所在包，因为包和当前包时同一个。

//go:decor dumpDecorContext
func useScopeInnerDecor(s string, i int) string {
	return "useLocalScopeDecor concat: " + s
}

// 这个函数实现了 func(*decor.Context [, any]) 类型，所以他是个装饰器，
// 可以通过 //go:decor 指令用在任意函数和方法上。
//
// 它的作用是在目标方法执行前和执行后分别打印出他们的 ctx 内容，
// 包括目标函数的类型、输入、输出、执行的次数。
func dumpDecorContext(ctx *decor.Context) {
	tpl := " dumpDecorContext: Kind: %+v, TargetName: %+v, Receiver: %+v, TargetIn: %+v, TargetOut: %+v, doRef: %+v\n"
	g.Printf("=>"+tpl, ctx.Kind, ctx.TargetName, ctx.Receiver, ctx.TargetIn, ctx.TargetOut, ctx.DoRef())
	ctx.TargetDo()
	g.Printf("<="+tpl, ctx.Kind, ctx.TargetName, ctx.Receiver, ctx.TargetIn, ctx.TargetOut, ctx.DoRef())
}

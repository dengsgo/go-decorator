package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"github.com/dengsgo/go-decorator/example/usages/g"
)

// 下面演示`type T types` 和它的方法同时使用装饰器。
// 方法的装饰器会先执行，然后再执行类型的装饰器。
// 比如下面的 multipleStructStandType 的 sayHello 方法，类型使用了装饰器 dumpDecorTextMore，
// 方法 使用了 dumpDecorText、dumpDecorTextAgain 这两个装饰器。那么执行的顺序为：
// dumpDecorText
// dumpDecorTextAgain
// dumpDecorTextMore
//
// 提示：不推荐同时使用多个装饰器装饰目标函数！这会增加开发者阅读代码的难度。

//go:decor dumpDecorTextMore#{text: "from type multipleStructStandType struct{}"}
type multipleStructStandType struct{}

//go:decor dumpDecorText#{text: "from sayHello()"}
//go:decor dumpDecorTextAgain#{text: "from sayHello()"}
func (m multipleStructStandType) sayHello() string {
	return "hello, decorator"
}

type (
	//go:decor dumpDecorTextMore#{text: "from multipleStructWrapType struct{}"}
	multipleStructWrapType struct {
	}
)

//go:decor dumpDecorText#{text: "from sayNiHao()"}
//go:decor dumpDecorTextAgain#{text: "from sayNiHao()"}
func (m multipleStructWrapType) sayNiHao() string {
	return "hello, decorator"
}

//go:decor-lint nonzero: {text}
func dumpDecorText(ctx *decor.Context, text string) {
	g.PrintfLn("dumpDecorText: TargetName: %+v, text: %+v", ctx.TargetName, text)
	ctx.TargetDo()
}

//go:decor-lint nonzero: {text}
func dumpDecorTextAgain(ctx *decor.Context, text string) {
	g.PrintfLn("dumpDecorTextAgain: TargetName: %+v, text: %+v", ctx.TargetName, text)
	ctx.TargetDo()
}

//go:decor-lint nonzero: {text}
func dumpDecorTextMore(ctx *decor.Context, text string) {
	g.PrintfLn("dumpDecorTextMore: TargetName: %+v, text: %+v", ctx.TargetName, text)
	ctx.TargetDo()
}

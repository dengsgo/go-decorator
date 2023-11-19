package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"github.com/dengsgo/go-decorator/example/usages/g"
)

// 这个文件演示带有 Receiver 的方法如何装饰器。
// 无论 Receiver 是结构体自己还是指针，用法和普通的函数没有人任何区别，
// 装饰器上下文 ctx.Kind 的值为 decor.KMethod , 代表当前的目标函数是个方法。

type methodTestPointerStruct struct{}

//go:decor dumpDecorContext
func (m *methodTestPointerStruct) doSomething(s string) string {
	return "*methodTestPointerStruct.recPointerDoSomething: " + s
}

type methodTestRawStruct struct{}

//go:decor dumpDecorContext
func (m methodTestRawStruct) doSomething(s string) string {
	return "methodTestRawStruct.recPointerDoSomething: " + s
}

type methodTestStruct struct{}

//go:decor validCtxReceiver
func (m *methodTestStruct) todo() {}

func validCtxReceiver(ctx *decor.Context) {
	if ctx.Kind == decor.KMethod &&
		ctx.TargetName == "todo" && func() bool {
		_, ok := ctx.Receiver.(*methodTestStruct)
		return ok
	}() {
		g.Printf("validCtxReceiver OK")
	} else {
		g.Printf("validCtxReceiver FAIL")
	}
}

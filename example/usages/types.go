package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"github.com/dengsgo/go-decorator/example/usages/g"
)

// 这个文件演示给 `type T types` 类型声明添加注释//go:decor F，
// decorator 会自动使用装饰器 F 装饰代理以 T 或者 *T 为接收者的所有方法.
// 比如下面的
// //go:decor dumpTargetType
// type structType struct {
//	 name string
// }
//
// 添加注释//go:decor dumpTargetType，
// structType 的方法集 Name、StrName、empty 会自动被装饰器 dumpTargetType 代理装饰。
// 方法的接收者可以是值接收者，也可以是指针接收者，都会被自动装饰。

//go:decor dumpTargetType
type structType struct {
	name string
}

func (s *structType) Name() string {
	g.PrintfLn("structType: %v", s.name)
	return s.name
}

func (s *structType) StrName(name string) {
	s.name = name
}

func (s *structType) empty() {}

//go:decor dumpTargetType
type varIntType int

func (v varIntType) value() int {
	return int(v)
}

func (v varIntType) zeroSelf() {
	v = 0
}

//go:decor dumpTargetType
type VarStringType string

func (v VarStringType) value() string {
	return string(v)
}

//go:decor dumpTargetType
type nonMethodType struct{}

//go:decor dumpTargetType
type otherFileDefMethodType struct{}

//go:decor dumpTargetType
type genericType[T any] struct {
	t T
}

func (g *genericType[T]) value() T {
	return g.t
}

func dumpTargetType(ctx *decor.Context) {
	g.PrintfLn("dumpTargetType say: Receiver: %+v, TargetName: %+v", ctx.Receiver, ctx.TargetName)
	ctx.TargetDo()
}

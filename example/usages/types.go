package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"github.com/dengsgo/go-decorator/example/usages/g"
)

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

func (v *VarStringType) value() string {
	return string(*v)
}

//go:decor dumpTargetType
type nonMethodType struct{}

//go:decor dumpTargetType
type otherFileDefMethodType struct{}

func dumpTargetType(ctx *decor.Context) {
	g.PrintfLn("dumpTargetType say: Receiver: %v, TargetName: %+v", ctx.Receiver, ctx.TargetName)
	ctx.TargetDo()
}

package main

import (
	"github.com/dengsgo/go-decorator/example/usages/g"
	"strings"
	"testing"
)

func TestMultipleStructStandType_sayHello(t *testing.T) {
	m := multipleStructStandType{}
	m.sayHello()
	out := strings.TrimSpace(g.TestBuffers.String())
	r := `dumpDecorText: TargetName: sayHello, text: from sayHello()
dumpDecorTextAgain: TargetName: sayHello, text: from sayHello()
dumpDecorTextMore: TargetName: sayHello, text: from type multipleStructStandType struct{}`
	if out != r {
		t.Fatalf("TestNonMethodType fail, out : %s, \nshould : %s", out, r)
	}
	g.ResetTestBuffers()
}

func TestMultipleStructWrapType_sayNiHao(t *testing.T) {
	m := multipleStructWrapType{}
	m.sayNiHao()
	out := strings.TrimSpace(g.TestBuffers.String())
	r := `dumpDecorText: TargetName: sayNiHao, text: from sayNiHao()
dumpDecorTextAgain: TargetName: sayNiHao, text: from sayNiHao()
dumpDecorTextMore: TargetName: sayNiHao, text: from multipleStructWrapType struct{}`
	if out != r {
		t.Fatalf("TestNonMethodType fail, out : %s, \nshould : %s", out, r)
	}
	g.ResetTestBuffers()
}

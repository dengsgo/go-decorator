package main

import (
	"github.com/dengsgo/go-decorator/example/usages/g"
	"log"
	"strings"
	"testing"
)

func TestMethodTestPointerStruct_doSomething(t *testing.T) {
	m := &methodTestPointerStruct{}
	m.doSomething("main called")
	out := strings.TrimSpace(g.TestBuffers.String())
	s := `=> dumpDecorContext: Kind: 1, TargetName: doSomething, Receiver: &{}, TargetIn: [main called], TargetOut: [], doRef: 0
<= dumpDecorContext: Kind: 1, TargetName: doSomething, Receiver: &{}, TargetIn: [main called], TargetOut: [*methodTestPointerStruct.recPointerDoSomething: main called], doRef: 1`
	log.Println(out)
	if out != s {
		t.Fatalf("TestDoSomething methodTestPointerStruct fail")
	}
	g.ResetTestBuffers()
}

func TestMethodTestRawStruct_doSomething(t *testing.T) {
	m := methodTestRawStruct{}
	m.doSomething("main called")
	out := strings.TrimSpace(g.TestBuffers.String())
	s := `=> dumpDecorContext: Kind: 1, TargetName: doSomething, Receiver: {}, TargetIn: [main called], TargetOut: [], doRef: 0
<= dumpDecorContext: Kind: 1, TargetName: doSomething, Receiver: {}, TargetIn: [main called], TargetOut: [methodTestRawStruct.recPointerDoSomething: main called], doRef: 1`
	log.Println(out)
	if out != s {
		t.Fatalf("TestDoSomething methodTestRawStruct fail")
	}
	g.ResetTestBuffers()
}

func TestMethodTestStruct_todo(t *testing.T) {
	m := &methodTestStruct{}
	m.todo()
	s := `validCtxReceiver OK`
	if s != g.TestBuffers.String() {
		t.Fatalf("TestMethodTestStruct_todo fail")
	}
	g.ResetTestBuffers()
}

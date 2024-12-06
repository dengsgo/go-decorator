package main

import (
	"strings"
	"testing"

	"github.com/dengsgo/go-decorator/example/usages/g"
)

func TestUnderscoresParamIn1Out1(t *testing.T) {
	underscoresParamIn1Out1("")
	out := `=> dumpDecorContext: Kind: 0, TargetName: underscoresParamIn1Out1, Receiver: <nil>, TargetIn: [], TargetOut: [0], doRef: 0
<= dumpDecorContext: Kind: 0, TargetName: underscoresParamIn1Out1, Receiver: <nil>, TargetIn: [], TargetOut: [1], doRef: 1`
	if strings.TrimSpace(g.TestBuffers.String()) != strings.TrimSpace(out) {
		t.Fatalf("TestUnderscoresParamIn1Out1 fail, out not match. \nshould: %+v\n, but: %+v", out, g.TestBuffers.String())
	}
	g.ResetTestBuffers()
}

func TestUnderscoresParamInHyperOut1(t *testing.T) {
	underscoresParamInHyperOut1("", 10086)
	out := `=> dumpDecorContext: Kind: 0, TargetName: underscoresParamInHyperOut1, Receiver: <nil>, TargetIn: [ 10086], TargetOut: [0], doRef: 0
<= dumpDecorContext: Kind: 0, TargetName: underscoresParamInHyperOut1, Receiver: <nil>, TargetIn: [ 10086], TargetOut: [1], doRef: 1`
	if strings.TrimSpace(g.TestBuffers.String()) != strings.TrimSpace(out) {
		t.Fatalf("TestUnderscoresParamInHyperOut1 fail, out not match. \nshould: %+v\n, but: %+v", out, g.TestBuffers.String())
	}
	g.ResetTestBuffers()
}

func TestUnderscoresParamIn2Out2(t *testing.T) {
	underscoresParamIn2Out2("", 10010)
	out := `=> dumpDecorContext: Kind: 0, TargetName: underscoresParamIn2Out2, Receiver: <nil>, TargetIn: [ 10010], TargetOut: [0 0], doRef: 0
<= dumpDecorContext: Kind: 0, TargetName: underscoresParamIn2Out2, Receiver: <nil>, TargetIn: [ 10010], TargetOut: [1 1], doRef: 1`
	if strings.TrimSpace(g.TestBuffers.String()) != strings.TrimSpace(out) {
		t.Fatalf("TestUnderscoresParamIn2Out2 fail, out not match. \nshould: %+v\n, but: %+v", out, g.TestBuffers.String())
	}
	g.ResetTestBuffers()
}

func TestUnderscoresParamIn0Out2(t *testing.T) {
	underscoresParamIn0Out2()
	out := `=> dumpDecorContext: Kind: 0, TargetName: underscoresParamIn0Out2, Receiver: <nil>, TargetIn: [], TargetOut: [0 0], doRef: 0
<= dumpDecorContext: Kind: 0, TargetName: underscoresParamIn0Out2, Receiver: <nil>, TargetIn: [], TargetOut: [1 1], doRef: 1`
	if strings.TrimSpace(g.TestBuffers.String()) != strings.TrimSpace(out) {
		t.Fatalf("TestUnderscoresParamIn0Out2 fail, out not match. \nshould: %+v\n, but: %+v", out, g.TestBuffers.String())
	}
	g.ResetTestBuffers()
}

func TestUnderscoresParamIn2Out0(t *testing.T) {
	underscoresParamIn2Out0(1, 2)
	out := `=> dumpDecorContext: Kind: 0, TargetName: underscoresParamIn2Out0, Receiver: <nil>, TargetIn: [1 2], TargetOut: [], doRef: 0
<= dumpDecorContext: Kind: 0, TargetName: underscoresParamIn2Out0, Receiver: <nil>, TargetIn: [1 2], TargetOut: [], doRef: 1`
	if strings.TrimSpace(g.TestBuffers.String()) != strings.TrimSpace(out) {
		t.Fatalf("TestUnderscoresParamIn2Out0 fail, out not match. \nshould: %+v\n, but: %+v", out, g.TestBuffers.String())
	}
	g.ResetTestBuffers()
}

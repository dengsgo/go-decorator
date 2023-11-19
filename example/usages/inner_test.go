package main

import (
	"github.com/dengsgo/go-decorator/example/usages/g"
	"os"
	"strings"
	"testing"
)

func setup() {
	g.InTest = true
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestUseScopeInnerDecor(t *testing.T) {
	cas := []struct {
		ins string
		ini int

		out string
	}{
		{
			"hello, world", 100,
			`=> dumpDecorContext: Kind: 0, TargetName: useScopeInnerDecor, Receiver: <nil>, TargetIn: [hello, world 100], TargetOut: [], doRef: 0
<= dumpDecorContext: Kind: 0, TargetName: useScopeInnerDecor, Receiver: <nil>, TargetIn: [hello, world 100], TargetOut: [useLocalScopeDecor concat: hello, world], doRef: 1`,
		},
		{
			"hello,中国", 999,
			`=> dumpDecorContext: Kind: 0, TargetName: useScopeInnerDecor, Receiver: <nil>, TargetIn: [hello,中国 999], TargetOut: [], doRef: 0
<= dumpDecorContext: Kind: 0, TargetName: useScopeInnerDecor, Receiver: <nil>, TargetIn: [hello,中国 999], TargetOut: [useLocalScopeDecor concat: hello,中国], doRef: 1`,
		},
	}
	for i, v := range cas {
		_ = useScopeInnerDecor(v.ins, v.ini)
		s := g.TestBuffers.String()
		s = strings.ReplaceAll(s, "\r\n", "\n")
		v.out = strings.ReplaceAll(v.out, "\r\n", "\n")
		if strings.TrimSpace(s) != strings.TrimSpace(v.out) {
			t.Fatalf("TestUseScopeInnerDecor useScopeInnerDecor fail, out not match. index: %+v\n case: %+v\nshould: %+v\n",
				i, v, s)
		}
		g.ResetTestBuffers()
	}
}

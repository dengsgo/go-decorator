package main

import (
	_ "github.com/dengsgo/go-decorator/decor"
	_ "github.com/dengsgo/go-decorator/example/usages/externala"
	"github.com/dengsgo/go-decorator/example/usages/g"
	"strings"
	"testing"
)

func TestUseExternalaDecor(t *testing.T) {
	out := `the target use [externala.OnlyPrintSelf] decorator
Return String By [externala.deepexternal.FixedStringWhenReturnString]`
	useExternalaDecor()
	if strings.TrimSpace(g.TestBuffers.String()) != strings.TrimSpace(out) {
		t.Fatalf("TestUseExternalaDecor fail")
	}
	g.ResetTestBuffers()
}

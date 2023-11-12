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

func TestPlus(t *testing.T) {
	cas := []struct {
		a,
		b,
		r int
	}{
		{2, 3, 10},
	}
	for i, v := range cas {
		num := plus(v.a, v.b)
		if num != v.r {
			t.Fatalf("TestPlus fail case %+v: plus(%+v, %+v) = %+v, but got %v",
				i, v.a, v.b, num, v.r)
		}
	}
	g.ResetTestBuffers()
}

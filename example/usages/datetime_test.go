package main

import (
	_ "github.com/dengsgo/go-decorator/decor"
	_ "github.com/dengsgo/go-decorator/example/usages/externala"
	"github.com/dengsgo/go-decorator/example/usages/g"
	"strings"
	"testing"
)

func TestDatetime(t *testing.T) {
	out := `logging print target in [1692450000]
logging print target out [2023-08-19 21:00:00 +0800 CST]
datetime(1692450000)=2023-08-19 21:00:00 +0800 CST`
	_t := 1692450000
	s := datetime(_t)
	g.Printf("datetime(%d)=%s\n", _t, s)
	if strings.TrimSpace(g.TestBuffers.String()) != strings.TrimSpace(out) {
		t.Fatalf("TestDatetime fail")
	}
	g.ResetTestBuffers()
}

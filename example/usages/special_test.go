package main

import (
	_ "github.com/dengsgo/go-decorator/decor"
	_ "github.com/dengsgo/go-decorator/example/usages/externala"
	"github.com/dengsgo/go-decorator/example/usages/g"
	"strings"
	"testing"
)

func TestNilInAndOut(t *testing.T) {
	out := `logging print target in [<nil> <nil> test]
logging print target out [false <nil>]`
	_, _ = nilInAndOut(nil, nil, "test")
	if strings.TrimSpace(g.TestBuffers.String()) != strings.TrimSpace(out) {
		t.Fatalf("TestNilInAndOut fail")
	}
	g.ResetTestBuffers()
}

func TestEllipsisIn(t *testing.T) {
	out := `logging print target in [0 [hello world !]]
logging print target out [[2024 1 1]]`
	_ = ellipsisIn(0, "hello", "world", "!")
	if strings.TrimSpace(g.TestBuffers.String()) != strings.TrimSpace(out) {
		t.Fatalf("TestEllipsisIn fail")
	}
	g.ResetTestBuffers()
}

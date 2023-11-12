package main

import (
	_ "github.com/dengsgo/go-decorator/decor"
	_ "github.com/dengsgo/go-decorator/example/usages/externala"
	"github.com/dengsgo/go-decorator/example/usages/g"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("SumInt", func(t *testing.T) {
		num := Sum(1, 2, 3, 4, 5, 6, 7, 8, 9)
		if num != 45 {
			t.Fatalf("TestSum fail case name: %s", t.Name())
		}
	})

	t.Run("SumFloat", func(t *testing.T) {
		num := Sum(4.5, 100.0, 990.0)
		if num != 1094.5 {
			t.Fatalf("TestSum fail case name: %s", t.Name())
		}
	})
	g.ResetTestBuffers()
}

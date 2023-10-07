package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"log"
)

func logging(ctx *decor.Context) {
	log.Println("logging print target in", ctx.TargetIn)
	ctx.TargetDo()
	log.Println("logging print target out", ctx.TargetOut)
}

//go:decor logging
func plus[T int8 | int16 | int | int32 | int64 | float32 | float64](a, b T) T {
	return a + b
}

func main() {
	log.SetFlags(0)
	log.Println("plus(1, 10) = ", plus(1, 10))
}

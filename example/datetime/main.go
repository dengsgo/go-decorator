package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"log"
	"time"
)

func logging(ctx *decor.Context) {
	log.Println("logging print target in", ctx.TargetIn)
	ctx.TargetDo()
	log.Println("logging print target out", ctx.TargetOut)
}

// Convert timestamp to string date format.
//
//go:decor logging
func datetime(timestamp int) string {
	return time.Unix(int64(timestamp), 0).String()
}

func main() {
	t := 1692450000
	s := datetime(t)
	log.Printf("datetime(%d)=%s\n", t, s)
}

func init() {
	log.SetFlags(0)
}

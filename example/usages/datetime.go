package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"github.com/dengsgo/go-decorator/example/usages/g"
	"time"
)

func logging(ctx *decor.Context) {
	g.PrintfLn("logging print target in %+v", ctx.TargetIn)
	ctx.TargetDo()
	g.PrintfLn("logging print target out %+v", ctx.TargetOut)
}

// Convert timestamp to string date format.
//
//go:decor logging
func datetime(timestamp int) string {
	return time.Unix(int64(timestamp), 0).String()
}

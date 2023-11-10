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

	sdt := &Student{"tim"}
	log.Printf("%s\n", sdt.Name())
	sdt = &Student{"cook"}
	log.Printf("%s\n", sdt.Name())
}

func init() {
	log.SetFlags(0)
	time.Local = time.FixedZone("CST", 8*3600)
}

// ============ type method ===========

type Student struct {
	name string
}

//go:decor update
func (s *Student) Name() string {
	return s.name
}

func update(ctx *decor.Context) {
	ctx.TargetDo()
	if len(ctx.TargetOut) >= 1 &&
		func() bool {
			s, ok := ctx.TargetOut[0].(string)
			return ok && s == "tim"
		}() {
		ctx.TargetOut[0] = "mike"
	}
}

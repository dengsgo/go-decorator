package main

import (
	"fmt"
	"github.com/dengsgo/go-decorator/decor"
	"log"
)

// This is a decorator function with parameters.
// It checks if the first element of ctx.TargetOut is a string, and if it is, it replaces that element
// with a formatted string that includes the values of the input parameters.
// Use `go:decor-lint` to add call constraints, such as which arguments are "requirbu" and so on.
// If you don't meet the constraints, you will get an error at compile time.
//
//go:decor-lint required: {msg, count, repeat, f}
//go:decor-lint nonzero: {msg, count, f}
func hit(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	ctx.TargetDo()
	if len(ctx.TargetOut) >= 1 &&
		func() bool {
			_, ok := ctx.TargetOut[0].(string)
			return ok
		}() {
		ctx.TargetOut[0] = fmt.Sprintf("hit received: msg=%s, count=%d, repeat=%t, f=%f, opt=%s\n",
			msg, count, repeat, f, opt)
	}
}

// The function has a decorator called hit with some arguments.
// The decorator is applied to the function using a comment with the go:decor directive.
// The decorator is expected to modify the behavior of the function in some way.
// The function itself does not have any implementation and returns an empty string.
//
//go:decor hit#{msg: "message from decor", repeat: true, count: 10, f:1}
func useArgsDecor() (s string) {
	return
}

func main() {
	s := useArgsDecor()
	log.Printf("useArgsDecor()=%s\n", s)
}

func init() {
	log.SetFlags(0)
}

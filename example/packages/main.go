package main

import "log"

// The file uses a decorator and needs to be imported anonymously using the package import method.
// Like this:
// import _ "github.com/dengsgo/go-decorator/decor"
//
// If the decorator used from another package and the package has not been imported,
// it needs to be imported anonymously.
// Just like the Fun1 package below:
// _ "github.com/dengsgo/go-decorator/example/packages/fun1"
//
// Why do I need anonymous imports? This is because the go specification does not allow
// the introduction of packages that are not used.
// If the package has been imported and used normally, there is no need to import anonymously again
import (
	_ "github.com/dengsgo/go-decorator/decor"
	"github.com/dengsgo/go-decorator/example/packages/fun"
	_ "github.com/dengsgo/go-decorator/example/packages/fun1"
)

func main() {
	log.SetFlags(0)
	log.Println("MAIN call myFunc1()")
	myFunc1()
	fun.Ts()
}

//go:decor fun.DecorHandlerFunc
func myFunc1() {
	log.Println("call myFunc1")
}

//go:decor fun.DecorHandlerFunc
func myFunc2() {
	log.Println("call myFunc1")
}

//go:decor fun1.DecorHandlerFunc
func myFunc3() {
	log.Println("call myFunc1")
}

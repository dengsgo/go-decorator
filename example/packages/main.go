package main

// Here, use go:decor to import the package where the decorator function is located.
// Multiple packages can be imported in multiple lines, just like import statements
//
import (
	_ "github.com/dengsgo/go-decorator/decor"
	"github.com/dengsgo/go-decorator/example/packages/fun"
	_ "github.com/dengsgo/go-decorator/example/packages/fun1"
	"log"
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

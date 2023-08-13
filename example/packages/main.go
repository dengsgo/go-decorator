package main

// Here, use go:decor to import the package where the decorator function is located.
// Multiple packages can be imported in multiple lines, just like import statements
//
//go:decor github.com/dengsgo/go-decorator/example/packages/fun
import (
	"log"
)

func main() {
	log.Println("MAIN call myFunc1()")
	myFunc1()
}

//go:decor fun.DecorHandlerFunc
func myFunc1() {
	log.Println("call myFunc1")
}

# go-decorator

[中文](README.zh_cn.md) | English (Translated from Chinese document)

> Don`t apply it in production environment !!!  
> Project is still under iterative development and is open for testing only. 

`go-decorator` is a compilation chaining tool that implements the decorator feature of the Go language.

In the go code, use the `//go:decorator decoratorfunctionName` to annotate functions, and you can use the `decoratorfunctionName` to quickly complete template code injection, change function behavior, control logical flow, and more.

The usage scenarios of decorators can be compared to other languages, such as Python and TypeScript.

`go-decorator` is a decorator injection performed during compilation, so it does not damage the source files of your project, nor does it generate new `.go` files or other unnecessary files in the project. It is completely different from the `go:generate` generative expression.

## Install

The project is still actively developing, and the best installation method currently is source code compilation (binary distribution will be provided after stabilization).

```shell
$ git clone https://github.com/dengsgo/go-decorator.git
$ cd go-decorator/cmd/decorator
$ go build
```

Successfully compiled will result in the `decorator` binary file. Add this file path to your environment variable `Path` for future calls.

You can also directly `go install`:
```shell
$ go install github.com/dengsgo/go-decorator/cmd/decorator@latest
```

Run `decorator --help` and the help message `decorator` appears to indicate successful installation.

Note: Please update frequently to install the latest version. Get the best experience.

## Usage

Simply add the `-toolexec decorator` parameter to the `go build` command.

If you usually use `go build`, now all you need to do is add the toolchain parameter to become `go build -toolexec decorator`. Everything else is the same as before, no changes need to be made!

The other subcommands of go are also used in the same way.

## Code

Introducing decorator dependencies in your project (must be a `go.mod` project):

```shell
$ go get -u github.com/dengsgo/go-decorator
```

Write similar code:

```go
package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"log"
)

func main() {
	// Call your function normally.
	// Since this is a function that declares the use of decorator logging,
	// The decorator compilation chain injects calls to the decorator method logging during code compilation.
	// So after compiling and running using the above method, you will get the following output:
	// 
	// 2023/08/13 20:26:30 decorator function logging in []
	// 2023/08/13 20:26:30 this is a function: myFunc
	// 2023/08/13 20:26:30 decorator function logging out []
	// 
	// Instead of just one sentence output from myFunc itself.
	// That is to say, the behavior of this method has been changed through the decorator!
	myFunc() 
}

// Declare that the function will be decorated using the decorator logging by using the go: decor annotation.
//
//go:decor logging
func myFunc() {
	log.Println("this is a function: myFunc")
}

// This is a regular function. 
// But it implements the func (*decor.Context) type, so it is still a decorator method,
// You can use this decorator on other functions.
// In the function, ctx is the decorator context, and the input and output parameters of the target function 
// can be obtained through ctx and the execution of the target method.
// If ctx.TargetDo() is not executed in the function, it means that the target function will not execute,
// Even if you call the decorated target function in your code! At this point, the objective function returns zero values.
// Before ctx.TargetDo(), ctx.TargetIn can be modified to change the input parameter values.
// After ctx.TargetDo(), you can modify ctx.TargetOut to change the return value.
// Only the values of the input and output parameters can be changed. Don't try to change their type and quantity, as this will trigger runtime panic!!!
func logging(ctx *decor.Context) {
	log.Println("decorator function logging in", ctx.TargetIn)
	ctx.TargetDo()
	log.Println("decorator function logging out", ctx.TargetOut)
}

```

## Example


[example](example) This directory demonstrates how to write code correctly to use the `decorator` tool.

[**single**](example/single): This is a single file example where both the decorator definition and the decorated function are located within the same package. In this case, there is no need to consider the issue of importing dependent packages, just use the example code.

[**packages**](example/packages)：The example of this project is that the decorator definition and the decorated function are not in the same package, and anonymous package import is required.

More examples and manuals are gradually being added.

## Requirement

- go1.18+  
- go.mod project

## Issue

If you find any problems, you can provide feedback here. [GitHub Issues](https://github.com/dengsgo/go-decorator/issues)  

## Contribute

The project is still under development and due to frequent changes, external contributions are temporarily not accepted. Welcome to submit a Pull Request after stabilization
# go-decorator

[![Actions](https://github.com/dengsgo/go-decorator/workflows/Go/badge.svg)](https://github.com/dengsgo/go-decorator/actions)  [![Go Report Card](https://goreportcard.com/badge/github.com/dengsgo/go-decorator)](https://goreportcard.com/report/github.com/dengsgo/go-decorator)  [![godoc.org](https://godoc.org/github.com/dengsgo/go-decorator/decor?status.svg)](https://godoc.org/github.com/dengsgo/go-decorator/decor)  [![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/dengsgo/go-decorator/decor)  


[中文](README.zh_cn.md) | English (Translated from Chinese document)

> Don`t apply it in production environment !!!  
> Project is still under iterative development and is open for testing only. 

`go-decorator` is a middleware tool for Go language compilers that enables non-intrusive decorator usage via annotations.


## Feature

- Use `//go:decor decoratorfunctionName` to annotate a function with the decorator `decoratorfunctionName` for quick sample code injection, non-intrusive changes to function behavior, and control of logic flow;  
- Define a function of type `func(*decor.Context)`, which can be used as a decorator for any top-level function.  
- support the use of multiple (line) `//go:decor` decorators to decorate functions.  
- Provide helpful error hints to detect problems at compile time and give the cause and line number of the error (e.g. undefined decorator or unreferenced package, etc.).  
- The target function is only executed at compile time and does not affect the performance of the compiled program, and without reflection operations.  
- It provides a basic usage guide.  

Decorator usage can be similar to other languages such as Python, TypeScript. (Ideal for caching, forensics, logging, and other scenarios, as a aid to free up duplicate coding).

`go-decorator` is a compile-time decorator injection technique. Using it does not affect your project's source files and does not generate additional `.go` files or other redundant files in your project. This injection method is very different from the `go:generate` generation method.
## Guide

查看： [中文文档](GUIDE.zh_cn.md#使用引导)  |  [English Guide](GUIDE.md#guide)  |  More

## Install

Install via `go install`.
```shell
$ go install github.com/dengsgo/go-decorator/cmd/decorator@latest
```

Run `decorator` and it will show you the `decorator` version.

Note: Please update frequently to install the latest version. Get the best experience.

## Usage

`decorator` is `go`'s compilation chaining tool, which relies on the `go` command to invoke it and compile the code.

Simply add the `-toolexec decorator` parameter to the `go build` command.

If you usually use `go build`, now all you need to do is add the toolchain parameter to become `go build -toolexec decorator`. Everything remains the same, no changes need to be made!

The other subcommands of go are also utilized in the same manner.

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

| Project                           | Notes                                                                                                                                                                                                                                              |
|-----------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [**single**](example/single)      | This is a single file example where both the decorator definition and the decorated function are located within the same package. In this case, there is no need to consider the issue of importing dependent packages, just use the example code. | 
| [**packages**](example/packages)  | The example of this project is that the decorator definition and the decorated function are not in the same package, and anonymous package import is required.                                                                                     |
| [**datetime**](example/datetime)  | The complete code used in the demonstration example in the Guide                                                                                                                                                                                   |
| [**emptyfunc**](example/emptyfunc) | The difference between calling and not calling `TargetDo()` in the demo decorator                                                                                                                                                                  |


See more [Guide](#guide) .

## Requirement

- go1.18+  
- go.mod project

## Issue

If you find any problems, you can provide feedback here. [GitHub Issues](https://github.com/dengsgo/go-decorator/issues)  

## Contribute

The project is still under development and due to frequent changes, external contributions are temporarily not accepted. Welcome to submit a Pull Request after stabilization.

## TODO

- More documents.
- IDE friendly tool support.
- better performance.
- More testing coverage.
- More clear error reminders.
- More...
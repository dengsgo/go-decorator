# go-decorator

[![Actions](https://github.com/dengsgo/go-decorator/workflows/Go/badge.svg)](https://github.com/dengsgo/go-decorator/actions)  [![Go Report Card](https://goreportcard.com/badge/github.com/dengsgo/go-decorator)](https://goreportcard.com/report/github.com/dengsgo/go-decorator)  [![godoc.org](https://godoc.org/github.com/dengsgo/go-decorator/decor?status.svg)](https://godoc.org/github.com/dengsgo/go-decorator/decor)  [![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/dengsgo/go-decorator/decor)  


[中文](README.zh_cn.md) | English 

> Beta version, use with caution in production environments. Click ⭐Star to follow the project's progress. 

`go-decorator`, Go 便捷使用装饰器的工具，装饰器能够切面 (AOP)、代理 (Proxy) 任意的函数和方法，提供观察和控制函数的能力。

`go-decorator` is a tool that allows Go to easily use decorators. The decorator can slice aspect (AOP) and proxy any function and method, providing the ability to observe and control functions.


## Feature

- Add the comment `//go:decor F` to use the decorator (`F` is the decorator function) to quickly complete the logic such as "boilerplate code injection, non-intrusive function behavior change, control logic flow";  
- You can freely define functions as decorators and apply them to any top-level function or method;  
- Support the use of multiple (line) `//go:decor` decorator decorate the functions;
- Support comment `type T types` type declaration, decorator will automatically decorate proxy all methods with `T` or `*T` as the receiver;  
- The decorator supports optional parameters, which brings more possibilities to development.  
- Support compile-time lint verification to ensure the robustness of Go compiled code.   
- Provide helpful error hints to detect problems at compile time and give the cause and line number of the error (e.g. undefined decorator or unreferenced package, etc.).  
- Enhancing the target function only at compile time does not degrade the performance of the compiled program, nor does it have reflection operations.

Decorator usage can be similar to other languages such as Python, TypeScript. (Ideal for caching, forensics, logging, and other scenarios, as an aid to free up duplicate coding).

> `go-decorator` is a compile-time decorator injection technique. Using it does not affect your project's source files and does not generate additional `.go` files or other redundant files in your project. This injection method is very different from the `go:generate` generation method.


## Guide

查看： [中文文档](GUIDE.zh_cn.md#使用引导)  |  [English Guide](GUIDE.md#guide)  

## Install

Install via `go install`.
```shell
$ go install github.com/dengsgo/go-decorator/cmd/decorator@latest
```

Run `decorator` and it will show you the `decorator` version.

```shell
$ decorator
decorator 0.15.0 beta , https://github.com/dengsgo/go-decorator
```

Tip: Run the above installation command frequently to install the latest version for bug fixes, enhanced experience, and more new features.

## Usage

`decorator` relies on the native `go` command to call it, just add the `-toolexec decorator` parameter to the subcommand of `go`.
For example:

|Native Command|Use `decorator`|
|--------|--------|
| `go build` | `go build -toolexec decorator` |
| `go run main.go` | `go run -toolexec decorator main.go` |
| `go test -v` | `go test -toolexec decorator -v` |
| `go install` | `go install -toolexec decorator` |
| `go ... -flags...` | `go ... -toolexec decorator -flags...` |


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
// In the function, ctx is the decorator context, and the function name, input and output 
// parameters, and execution of the target function can be obtained through ctx.
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

Decorator usage with parameters:

```go
package main

import (
	"github.com/dengsgo/go-decorator/decor"
)

func main()  {
	optionalParametersFuncDemo()
}

// It is also a decorator, unlike a normal decorator, which allows the objective function to provide an additional parameter 'level'.
// Decorator provides lint syntax for developers to force verification when compiling code. For example:
// `Required` requires that the target function must pass a value to this field;
// `nonzero` requires that the objective function cannot pass values in time and space.
// If the compilation test fails, the compilation will fail.
// The usage is as follows: (See Guide.md for more usage):
//
//go:decor-lint required: {level}
//go:decor-lint nonzero: {level}
func levelLogging(ctx *decor.Context, level string)  {
	if level == "debug" {
		// to do something
	}
	ctx.TargetDo()
}

// This method uses the decorator levelLogging and passes the "level" parameter value "debug" to the decorator.
// 
//go:decor levelLogging#{level: "debug"}
func optionalParametersFuncDemo()  {
	// function code
}
```

Add a comment to the' `type T types` type declaration `//go:decor F`, and the decorator will automatically use the decorator `F` to decorate all methods that have `T` or `*T` as receiver:

```go
package main

import (
	"github.com/dengsgo/go-decorator/decor"
)

// add comments//go: decor dumpTargetType,
// The structType method sets Name, StrName, and empty are automatically decorated by the decorator dumpTargetType proxy.
// The receiver of a method can be either a value receiver or a pointer receiver, and is automatically decorated.

//go:decor dumpTargetType
type structType struct {
	name string
}

func (s *structType) Name() string {
	return s.name
}

func (s *structType) StrName(name string) {
	s.name = name
}

func (s *structType) empty() {}
```


## Example

[example/usages](example/usages) This directory demonstrates how to write code correctly to use the `decorator` tool.

```go
func main() {
	section("inner.go")
	// Example: Using a function with a package decorator
	useScopeInnerDecor("hello, world", 100)

	section("external.go")
	// Example: Functions using other wrapper decorators 
	useExternalaDecor()
	g.PrintfLn("plus(2, 3) = %+v", plus(2, 3))

	section("datetime.go")
	// Example: Code demonstrating the use of decorators in document Guide.md
	{
		t := 1692450000
		s := datetime(t)
		g.Printf("datetime(%d)=%s\n", t, s)
	}

	section("genericfunc.go")
	// Example: Generic functions using decorators
	g.PrintfLn("Sum(1, 2, 3, 4, 5, 6, 7, 8, 9) = %+v", Sum(1, 2, 3, 4, 5, 6, 7, 8, 9))

	section("method.go")
	// Example: methods using decorators
	{
		m := &methodTestPointerStruct{}
		m.doSomething("main called")
	}
	{
		m := methodTestRawStruct{}
		m.doSomething("main called")
	}

	section("withdecorparams.go")
	// Example: How to pass values when using a decorator with parameters 
	g.PrintfLn("useArgsDecor() = %+v", useArgsDecor())
	// Example: How to use Lint to constrain and validate the parameters of the objective function during compilation in a decorator
	g.Printf("useHitUseRequiredLint() = %+v", useHitUseRequiredLint())
	g.Printf("useHitUseNonzeroLint() = %+v", useHitUseNonzeroLint())
	g.Printf("useHitBothUseLint() = %+v", useHitBothUseLint())
	g.Printf("useHitUseMultilineLintDecor() = %+v", useHitUseMultilineLintDecor())
}

```

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
# Guide

- [Guide](#guide)
  - [Requirement](#requirement)
  - [Install](#install)
  - [Usage](#usage)
    - [Adding parameters](#adding-parameters)
    - [Add a dependency](#add-a-dependency)
    - [Destination functions and decorators](#destination-functions-and-decorators)
    - [Customizing decorators](#customizing-decorators)
    - [Using decorators](#using-decorators)
    - [Using multiple decorators](#using-multiple-decorators)
  - [Context](#context)
    - [ctx.Kind](#ctxkind)
    - [ctx.TargetIn](#ctxtargetin)
    - [ctx.TargetOut](#ctxtargetout)
    - [ctx.TargetDo()](#ctxtargetdo)
    - [ctx.DoRef()](#ctxdoref)
  - [Package references](#package-references)
  - [Conditions and restrictions](#conditions-and-restrictions)
  - [Development and Debugging](#development-and-debugging)
  - [Performance](#performance)
  - [More](#more)


`go-decorator` is a compilation chain tool that implements the Go language decorator feature, allowing annotations to be used to apply decorators to any function.

## Requirement

- go1.18+
- go.mod project

## Install

```shell
$ go install github.com/dengsgo/go-decorator/cmd/decorator@latest
```

Run `decorator` to display the current version.

```shell
$ decorator
decorator 0.10.0 beta , https://github.com/dengsgo/go-decorator
```

Tip: Run the above installation command frequently to install the latest version for bug fixes, enhanced experience, and more new features.

## Usage

`decorator` is `go`'s compilation chaining tool, which relies on the `go` command to invoke it and compile the code.

### Adding parameters

`decorator` relies on the native `go` command to call it, just add the `-toolexec decorator` parameter to the subcommand of `go`.
For example:

|Native Command|Use `decorator`|
|--------|--------|
| `go build` | `go build -toolexec decorator` |
| `go run main.go` | `go run -toolexec decorator main.go` |
| `go test -v` | `go test -toolexec decorator -v` |
| `go install` | `go install -toolexec decorator` |
| `go ... -flags...` | `go ... -toolexec decorator -flags...` |


### Add a dependency

In your project root directory, add the `go-decorator` dependency.

```shell
$ go get -u github.com/dengsgo/go-decorator
```

### Understand destination functions and decorators

> Target functions: functions that use a decorator, also known as decorated functions.  
> For example, if a function A uses a decorator B to decorate itself, A is the target function.

Decorators are also functions. When code is run to the target function, it doesn't actually execute it, but runs the decorator it uses. The actual target function logic is wrapped into the decorator and allows the decorator to control it.

### Customizing decorators

A decorator is an ordinary `go` Top-level Function of type `func(*decor.Context [, ...any])`. As long as the function satisfies this type, it is a legal decorator and can be used to decorate other functions in the project code.

For example, here's a logging decorator that prints the arguments of the called function:

```go
package main

import "github.com/dengsgo/go-decorator/decor"

func logging(ctx *decor.Context) {
	log.Println("logging print target in", ctx.TargetIn)
	ctx.TargetDo()
	log.Println("logging print target out", ctx.TargetOut)
}
```

This function `logging` is a legal decorator that can be used on any first-class function.

For the `ctx *decor.Context` argument, jump here [Context](#context).

### Using decorators

Decorators can be used on any first-level function by annotating `//go:decor`.

For example, we have a function `datetime`, which converts a timestamp to a string date format. The `logging` decorator can be used to print the function in and out by `//go:decor logging`:

```go
// Omitted code ...

// Convert timestamp to string date format.
//
//go:decor logging
func datetime(timestamp int) string {
	return time.Unix(int64(timestamp), 0).String()
}

// Omitted code ...
```

`datetime` is recognized at compile time and injected into `logging` calls. When the `datetime` function is called elsewhere, `logging` is automatically executed.

For example, if we call `datetime` in the `main` entry function.

```go
func main() {
    t := 1692450000
    s := datetime(t)
    log.Printf("datetime(%d)=%s\n", t, s)
}
```

Compile, run with the following command.

```shell
$ go build -toolexec decorator
$ . /datetime
```

The following output will be seen:

```shell
2023/08/19 21:12:21 logging print target in [1692450000]
2023/08/19 21:12:21 logging print target out [2023-08-19 21:00:00 +0800 CST]
2023/08/19 21:12:21 datetime(1692450000)=2023-08-19 21:00:00 +0800 CST
```

Only the `datetime` function is called in our code, but you can see that the logging decorator used is also executed!

The full code can be seen in the [example/datetime](example/datetime).

### Using multiple decorators

`decorator` allows multiple decorators to be used at the same time to decorate the target function.

Multiple decorators can be used with multiple `//go:decor `.

For example, the following `datetime` uses 3 decorators, `logging`, `appendFile`, and `timeFollowing`:

```go
// Omitted code ...

// Convert timestamp to string date format.
//
//go:decor logging
//go:decor appendFile
//go:decor timeFollowing
func datetime(timestamp int64) string {
    return time.Unix(timestamp, 0).String()
}

// Omitted code ...
```

If more than one decorator is used, the decorator execution is prioritized from top to bottom, i.e. the one defined first is executed first. In the above decorator, the order of execution is `logging` -> `appendFile` -> `timeFollowing`.

The use of multiple decorators may result in less readable code and increase the cost of understanding the logic flow, especially if the decorator itself is particularly complex. This is not recommended.


### Decorator with additional parameters

As the name suggests, decorators allow for defining additional parameters in addition to the first parameter `*decor.Context`, such as:

```go
package main
import "github.com/dengsgo/go-decorator/decor"

func hit(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	// code...
} 
```

The `hit` function is a legitimate decorator with optional parameters, which allows the target function to pass in the corresponding value when calling, and the hit function can obtain the parameter value of the target function.

The following parameter types are allowed:

| types  | keyword|     
|-----|-----|
| Integer  | int,int8,int16.int32,int64,unit,unit8,unit16,unit32,unit64 |
| Float | float32,float64 |
| String | string |
| Boolean | bool |

If it exceeds the above types, it cannot be compiled.


### Using Decorators with Parameters

Use the `//go:decor function#{}` method to pass parameters to the decorator. Compared to non parametric calls, there is an additional section called `#{}`, which we refer to as the parameter field.

The parameter field starts with a `#` identifier, followed by key value pairs such as `{key: value, key1: value1}`. The key is the formal parameter name of the decorator, and the value is the String, Boolean value, or Numerical value to be passed.

For example, we need to call the `hit` decorator defined above:

```go
package main
import "github.com/dengsgo/go-decorator/decor"

func hit(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	// code...
}

//go:decor hit#{msg: "message from decor", repeat: true, count: 10, f:1}
func useArgsDecor()  {}
```

The `decorator` will automatically pass the `{msg: "message from decor", repeat: true, count: 10, f: 1}` parameters to the `decorator` according to their formal parameter names during compilation.

The order of parameters in the parameter field is independent of the formal parameter order of the decorator, and you can organize the code according to your own habits.

When there is no corresponding formal parameter value in the parameter field, such as `opt`  above, the corresponding type's zero value will be passed by default.

### Decorator constraints and validation

`decorator` allows the use of annotations `//go:decor-lint linter: {}` on decorators to add decorator constraints. This constraint can be used at compile time to verify whether the call to the target function is legal.

Currently, there are two built-in decorator constraints:

#### required

Validation parameters must be passed. For example:

```go
//go:decor-lint required: {msg, count, repeat, f}
func hit(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	// code...
}
```

The four parameters, `msg`, `count`, `repeat`, and `f`, require that the target function must be passed during invocation, otherwise compilation cannot pass.

Not only that, `required` also supports validation of enumerations and scopes. For example:

**Enumeration Value Restrictions**:

`//Go: decor int required: {msg: {"hello", "world", "yeah"}, count, repeat, f}`: The argument to 'msg' must be one of the three values`"hello", "world", "yeah"`.

**Scope limitations**:

`//Go: decor int required: {msg: {gte: 8, lte: 24}, count, repeat, f}`: The string length range for 'msg' is required to be between '[8,24]'.

There are currently four supported scope directives:

| 范围指令  | 说明 |
|-------|----|
| `gte` | `>=` |
| `gt`  | `>`  |
| `lte` | `<=` |
| `lt`  | `<`  |

#### nonzero

The validation parameter value cannot be zero. For example:

```go
//go:decor-lint nonzero: {msg, count, f}
func hit(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	// code...
}
```

The three parameters `msg`, `count`, and `f` require the target function to pass values that cannot be zero when called.

> You can add '//go:decor-lint' rule constraints multiple times on the decorator, which means that the target function must all meet these constraints when calling the decorator in order to compile properly.

## Context

`ctx *decor.Context` is the entry parameter of the decorator function, which is the context of the target function (i.e., the function that uses this decorator, also known as the decorated function).

This context can be used in the decorator to modify the in- and out-parameters of the target function, adjust the execution logic, and so on.

### ctx.Kind

Target function type, currently only `decor.KFunc`, function type.

### ctx.TargetIn

The list of inputs to the target function. It is a []any slice, where the type of each element corresponds to the type of the target function's entry parameter. If the target function has no in-parameters, the list is empty.

This slice is used by `ctx.TargetDo()` as an input to the real call, so changing its element values modifies the input to the target function. Changes are only valid before the `ctx.TargetDo()` call.

### ctx.TargetOut

A list of the out parameters of the target function. It is a []any slice, where the type of each element matches the type of the target function's output. If the target function has no outgoing parameters, the list is empty.

This slice is used by `ctx.TargetDo()` to receive the result of a real call, so changing the values of its elements modifies the arguments of the target function. Changes are only valid after a `ctx.TargetDo()` call.

### ctx.TargetDo()

Executes the target function. It is a parameterless wrapper around the target function, and calling it actually executes the target function logic.

It gets the target function input from `ctx.TargetIn`, executes the target function code, and assigns the result to `ctx.TargetOut`.

If `ctx.TargetDo()` is not executed in the decorator, it means that the real logic of the target function will not be executed, and the result of the call to the target function will be zero-value (without modifying ctx.TargetOut).  

### ctx.DoRef()

`DoRef()` gets the number of times an anonymous wrapper class has been executed.

Usually, it shows the number of times `TargetDo()` was called in the decorator function.

> Be careful when writing decorator code, be sure to assert the type of the element values of ctx.TargetIn, ctx.TargetOut, any incorrectly-typed assignments will generate a runtime panic.  
> Do not change ctx.TargetIn, ctx.TargetOut values (assign/append/delete, etc.), this will cause a serious error panic on ctx.TargetDo() calls.

## Package references

In the `datetime` [example/datetime](example/datetime) example above, our decorator and target function are in a package, and we don't need to think about packages.

Package references need to be considered when we have many packages.

The go specification prevents importing packages that aren't used by the code in the current file, which means that comments like `//go:decor` don't really import packages, so we need to use an anonymous package import to import the corresponding package. Like this `import _"path/to/your/package"`.

There are two cases where you need to import a package anonymously:

One, the function of the package uses decorator annotations, but does not import the package ``github.com/dengsgo/go-decorator/decor``, which requires us to add an anonymous import package:

```go
import _ "github.com/dengsgo/go-decorator/decor"
```

Second, if package (A) references a decorator of another package (B), and B is not `imported` by A, we need to import it using the anonymous import package.

For example:

```go
package main

// other imports
import _ "github.com/dengsgo/go-decorator/example/packages/fun1"

//go:decor fun1.DecorHandlerFunc
func test() {
	//...
}
```

Of course, if the package is already used by other code in the file and has already been imported, then there is no need to import it anonymously.

For a complete example check out the [example/packages](example/packages) .

## Conditions and restrictions

The following conditions require attention:

- The scope of the target function using the decorator **is limited to within the current project**. Other libraries that depend on it **cannot** be decorated even with `//go:decor`.

For example, if your project module name is `a/b/c`, then `//go:decor` will only work in `a/b/c` and its subpackages (`a/b/c/d` works, `a/m/` does not).

But `//go:decor` can use decorators from any package, with no scope restrictions.

- **Can't** use the same decorator repeatedly on the same target function at the same time;  
- **Can't** apply a decorator to a decorator function;  
- After upgrading `decorator` or adjusting compilation parameters it may be necessary to append the `-a` parameter to the go command **to force compilation** once to overwrite the old compilation cache.  

## Development and Debugging

The `decorator` is used by the go compiler as a link in the go compilation chain and is loaded at compile time. It is compatible with the go compilation chain and does not cause side effects.

The only thing you need to change in your development process is to add the `-toolexec decorator` parameter to the go commands you use, but everything else is exactly the same, so it doesn't feel like a change.

You can also remove this parameter at any time. Drop the project's use of the go decorator. Even if you keep the `//go:decor` comment in your code, it has no side effect (because it's just a meaningless comment to the standard toolchain).

The same applies to debugging.

For example, in vscode, edit `launch.json`.

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch file",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${file}",
            "buildFlags": "-toolexec decorator"
        }
    ]
}
```

Add the line `"buildFlags":"-toolexec decorator"` to enable decorator compilation for `decorator`.

Then just breakpoint and debug normally.

> The debugging experience will continue to improve, so please let me know if you find any problems! [Issues](https://github.com/dengsgo/go-decorator/issues)。

## Performance

Although `decorator` does extra processing on the target function at compile time, it only builds the necessary context parameters, with no extra overhead and no reflection. Performance is almost identical to calling the decorator function directly from the original go code.

// TODO provides a comparison of performance metrics

## More

// TODO